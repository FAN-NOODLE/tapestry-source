package user

import "strings"

const (
	AppCodeWebAdmin            = "web_admin"
	AppCodeStoreOperationAdmin = "store_operation_admin"
	ClientConsole              = "console"
	ClientStoreOperationAdmin  = "store_operation_admin"
)

// NormalizeAppCode keeps existing callers on the old admin application while
// rejecting unknown application codes before they reach the authority service.
func NormalizeAppCode(appCode string) (string, bool) {
	appCode = strings.TrimSpace(appCode)
	if appCode == "" {
		return AppCodeWebAdmin, true
	}
	switch appCode {
	case AppCodeWebAdmin, AppCodeStoreOperationAdmin:
		return appCode, true
	default:
		return "", false
	}
}

// ClientForAppCode selects the account system checked by micro-authority.
func ClientForAppCode(appCode string) string {
	if appCode == AppCodeStoreOperationAdmin {
		return ClientStoreOperationAdmin
	}
	return ClientConsole
}
