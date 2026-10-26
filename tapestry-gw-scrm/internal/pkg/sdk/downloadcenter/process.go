package downloadcenter

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

func (dc downloadCenter) Process(ctx ctx.Context, request ProcessRequest) (response ProcessResponse, err error) {
	err = dc.http.Post(ctx, dc.makeUrl(UrlProcess), request, &response)
	return
}
