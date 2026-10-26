export const WEB_ADMIN_APP_CODE = "web_admin";
export const STORE_OPERATION_ADMIN_APP_CODE = "store_operation_admin";

export const ROLE_APP_NAMES: Record<string, string> = {
  [WEB_ADMIN_APP_CODE]: "WeClient管理后台",
  [STORE_OPERATION_ADMIN_APP_CODE]: "门店运营协同平台",
};

export const ACCOUNT_SYSTEM_APP_CODES: Record<number, string> = {
  1: WEB_ADMIN_APP_CODE,
  4: STORE_OPERATION_ADMIN_APP_CODE,
};

export function roleAppName(appCode: string): string {
  return ROLE_APP_NAMES[appCode] || "WeClient管理后台";
}
