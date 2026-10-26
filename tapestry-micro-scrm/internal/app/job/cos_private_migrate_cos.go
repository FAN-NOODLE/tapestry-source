package appJob

import (
	"errors"
	"fmt"

	"github.com/tencentyun/cos-go-sdk-v5"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type cosPrivateObjectMeta struct {
	Size        int64
	notFoundErr error
}

type cosPrivateObjectStore interface {
	Stat(ctx.Context, string, string) (cosPrivateObjectMeta, bool, error)
	Copy(ctx.Context, string, string, string, string) error
}

var (
	errCosPrivateSourceObjectNotFound     = errors.New("COS source object does not exist")
	errCosPrivateTargetObjectVerification = errors.New("COS target object verification failed")
)

type cosPrivateObjectRowError struct {
	err error
}

func (e *cosPrivateObjectRowError) Error() string {
	return e.err.Error()
}

func (e *cosPrivateObjectRowError) Unwrap() error {
	return e.err
}

func (e *cosPrivateObjectRowError) Is(target error) bool {
	return target == errCosPrivateSourceObjectNotFound
}

func isCosPrivateObjectRowError(err error) bool {
	var rowErr *cosPrivateObjectRowError
	return errors.As(err, &rowErr)
}

type cosPrivateTencentClient interface {
	Head(ctx.Context, string) (int64, error)
	Copy(ctx.Context, string, string) error
}

type cosPrivateTencentSDKClient struct {
	client *cos.Client
}

func (c cosPrivateTencentSDKClient) Head(requestCtx ctx.Context, key string) (int64, error) {
	response, err := c.client.Object.Head(requestCtx, key, nil)
	if err != nil {
		return 0, err
	}
	return response.ContentLength, nil
}

func (c cosPrivateTencentSDKClient) Copy(requestCtx ctx.Context, key, copySource string) error {
	_, _, err := c.client.Object.Copy(requestCtx, key, copySource, nil)
	return err
}

type cosPrivateTencentStore struct {
	sourceBucket, targetBucket string
	source, target             cosPrivateTencentClient
}

func newCosPrivateTencentStore(sourceBucket, targetBucket string) *cosPrivateTencentStore {
	return &cosPrivateTencentStore{
		sourceBucket: sourceBucket,
		targetBucket: targetBucket,
		source:       cosPrivateTencentSDKClient{client: newCosClient(sourceBucket)},
		target:       cosPrivateTencentSDKClient{client: newCosClient(targetBucket)},
	}
}

func (s cosPrivateTencentStore) Stat(requestCtx ctx.Context, bucket, key string) (cosPrivateObjectMeta, bool, error) {
	client, err := s.clientFor(bucket)
	if err != nil {
		return cosPrivateObjectMeta{}, false, err
	}

	size, err := client.Head(requestCtx, key)
	if err == nil {
		return cosPrivateObjectMeta{Size: size}, true, nil
	}
	if cos.IsNotFoundError(err) {
		return cosPrivateObjectMeta{notFoundErr: err}, false, nil
	}
	return cosPrivateObjectMeta{}, false, err
}

func (s cosPrivateTencentStore) Copy(requestCtx ctx.Context, sourceBucket, targetBucket, sourceKey, targetKey string) error {
	if _, err := s.clientFor(sourceBucket); err != nil {
		return err
	}
	targetClient, err := s.clientFor(targetBucket)
	if err != nil {
		return err
	}
	return targetClient.Copy(requestCtx, targetKey, cosPrivateCopySource(sourceBucket, sourceKey))
}

func (s cosPrivateTencentStore) clientFor(bucket string) (cosPrivateTencentClient, error) {
	switch bucket {
	case s.sourceBucket:
		return s.source, nil
	case s.targetBucket:
		return s.target, nil
	default:
		return nil, fmt.Errorf("unexpected COS bucket %q", bucket)
	}
}

func cosPrivateCopySource(bucket, key string) string {
	return cosOriginHost(bucket) + "/" + key
}

type cosPrivateCopier struct {
	store                      cosPrivateObjectStore
	sourceBucket, targetBucket string
}

func (c *cosPrivateCopier) ensureCopied(requestCtx ctx.Context, key string) error {
	return c.ensureCopiedTo(requestCtx, c.sourceBucket, key, c.targetBucket, key)
}

func (c *cosPrivateCopier) ensurePublicCopied(requestCtx ctx.Context, sourceKey, targetKey string) error {
	return c.ensureCopiedTo(requestCtx, c.sourceBucket, sourceKey, c.sourceBucket, targetKey)
}

func (c *cosPrivateCopier) ensureCopiedTo(requestCtx ctx.Context, sourceBucket, sourceKey, targetBucket, targetKey string) error {
	source, sourceExists, err := c.store.Stat(requestCtx, sourceBucket, sourceKey)
	if err != nil {
		return fmt.Errorf("stat source object %q: %w", sourceKey, err)
	}
	if !sourceExists {
		cause := source.notFoundErr
		if cause == nil {
			cause = errCosPrivateSourceObjectNotFound
		}
		return &cosPrivateObjectRowError{err: fmt.Errorf("source object %q does not exist: %w", sourceKey, cause)}
	}

	target, targetExists, err := c.store.Stat(requestCtx, targetBucket, targetKey)
	if err != nil {
		return fmt.Errorf("stat target object %q: %w", targetKey, err)
	}
	if targetExists && target.Size == source.Size {
		return nil
	}

	if err := c.store.Copy(requestCtx, sourceBucket, targetBucket, sourceKey, targetKey); err != nil {
		return fmt.Errorf("copy object %q to %q: %w", sourceKey, targetKey, err)
	}

	target, targetExists, err = c.store.Stat(requestCtx, targetBucket, targetKey)
	if err != nil {
		return fmt.Errorf("verify target object %q: %w", targetKey, err)
	}
	if !targetExists {
		return fmt.Errorf("%w: target object %q does not exist after copy", errCosPrivateTargetObjectVerification, targetKey)
	}
	if target.Size != source.Size {
		return fmt.Errorf("%w: target object %q size %d does not match source size %d", errCosPrivateTargetObjectVerification, targetKey, target.Size, source.Size)
	}
	return nil
}
