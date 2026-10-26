package downloadcenter

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

func (dc downloadCenter) Delete(ctx ctx.Context, request DeleteRequest) (response DeleteResponse, err error) {
	err = dc.http.Post(ctx, dc.makeUrl(UrlDelete), request, &response)
	return
}
