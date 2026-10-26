package downloadcenter

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

func (dc downloadCenter) Retry(ctx ctx.Context, request RetryRequest) (response RetryResponse, err error) {
	err = dc.http.Post(ctx, dc.makeUrl(UrlRetry), request, &response)
	return
}
