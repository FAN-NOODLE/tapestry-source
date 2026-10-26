package source

import (
	"errors"
	"testing"

	"micro-scrm/internal/app"
	"micro-scrm/internal/app/models"
	sourcemodel "micro-scrm/internal/app/models/source"
	coachkos "micro-scrm/internal/pkg/sdk/coach_kos"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func TestCoachKOS封面上传后缀统一为JPG(t *testing.T) {
	testCases := []string{
		"https://third.example.com/cover.webp",
		"https://third.example.com/cover.png",
		"https://third.example.com/cover",
	}

	for _, rawURL := range testCases {
		if got := coachKOSCoverExtension(rawURL, []byte("image data")); got != ".jpg" {
			t.Fatalf("原始地址=%q，上传后缀=%q，期望=.jpg", rawURL, got)
		}
	}
}

func TestCoachKOS封面为空时不转存(t *testing.T) {
	uploadCalls := 0
	service := &CoachKOSSyncService{
		uploadCover: func(app.Ctx, string) (string, error) {
			uploadCalls++
			return "", nil
		},
	}

	got := service.prepareCoachKOSCoverURL(ctx.New(), 1001, coachkos.PostResult{PostID: "post-1"})
	if got != "" || uploadCalls != 0 {
		t.Fatalf("封面地址=%q，上传次数=%d，期望空地址且不上传", got, uploadCalls)
	}
}

func TestCoachKOS封面转存失败时使用原始地址(t *testing.T) {
	originalURL := "https://third.example.com/cover.webp"
	service := &CoachKOSSyncService{
		getKOSByPostID: func(app.Ctx, int64, string) (sourcemodel.LinkKOSModel, error) {
			return sourcemodel.LinkKOSModel{}, nil
		},
		uploadCover: func(app.Ctx, string) (string, error) {
			return "", errors.New("upload failed")
		},
	}

	got := service.prepareCoachKOSCoverURL(ctx.New(), 1001, coachkos.PostResult{PostID: "post-1", CoverURL: originalURL})
	if got != originalURL {
		t.Fatalf("封面地址=%q，期望回退到 %q", got, originalURL)
	}
}

func TestCoachKOS原始封面未变化时复用已处理地址(t *testing.T) {
	originalURL := "https://third.example.com/cover.webp"
	storedURL := "https://cdn.example.com/uploads/cover.webp?imageMogr2/quality/80/format/jpg"
	uploadCalls := 0
	service := &CoachKOSSyncService{
		cdnHost: "cdn.example.com",
		getKOSByPostID: func(app.Ctx, int64, string) (sourcemodel.LinkKOSModel, error) {
			return sourcemodel.LinkKOSModel{
				MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 31}},
				SourceID:       21,
				CoverURL:       originalURL,
			}, nil
		},
		getLinkByID: func(app.Ctx, int64, int64) (sourcemodel.LinkModel, error) {
			return sourcemodel.LinkModel{CoverURL: storedURL}, nil
		},
		uploadCover: func(app.Ctx, string) (string, error) {
			uploadCalls++
			return "", nil
		},
	}

	got := service.prepareCoachKOSCoverURL(ctx.New(), 1001, coachkos.PostResult{PostID: "post-1", CoverURL: originalURL})
	if got != storedURL || uploadCalls != 0 {
		t.Fatalf("封面地址=%q，上传次数=%d，期望复用 %q", got, uploadCalls, storedURL)
	}
}

func TestCoachKOS原始封面变化时重新上传并添加压缩参数(t *testing.T) {
	oldURL := "https://third.example.com/old.webp"
	newURL := "https://third.example.com/new.webp"
	service := &CoachKOSSyncService{
		cdnHost: "cdn.example.com",
		getKOSByPostID: func(app.Ctx, int64, string) (sourcemodel.LinkKOSModel, error) {
			return sourcemodel.LinkKOSModel{
				MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 31}},
				SourceID:       21,
				CoverURL:       oldURL,
			}, nil
		},
		getLinkByID: func(app.Ctx, int64, int64) (sourcemodel.LinkModel, error) {
			return sourcemodel.LinkModel{CoverURL: "https://cdn.example.com/uploads/old.webp?imageMogr2/quality/80/format/jpg"}, nil
		},
		uploadCover: func(_ app.Ctx, rawURL string) (string, error) {
			if rawURL != newURL {
				t.Fatalf("上传地址=%q，期望=%q", rawURL, newURL)
			}
			return "https://cdn.example.com/uploads/new.webp", nil
		},
	}

	got := service.prepareCoachKOSCoverURL(ctx.New(), 1001, coachkos.PostResult{PostID: "post-1", CoverURL: newURL})
	want := "https://cdn.example.com/uploads/new.webp?imageMogr2/quality/80/format/jpg"
	if got != want {
		t.Fatalf("封面地址=%q，期望=%q", got, want)
	}
}
