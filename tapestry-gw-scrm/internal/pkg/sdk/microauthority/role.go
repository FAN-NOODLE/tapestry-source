package microauthority

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

const configurablePermissionsURL = "%s/micro-authority/v1/rbac/role/configurable-permissions"

func ConfigurablePermissions(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(configurablePermissionsURL), request)
}
