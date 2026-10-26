/**
 * 通用接口返回数据类型定义
 */
export interface ApiResult<T = any> {
  code?: number;
  message?: string;
  data?: T;
  request_id?: string;
}
