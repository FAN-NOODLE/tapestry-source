import { Meta } from "./meta";

/**
 * 分页接口接收数据类型定义
 */
export interface PaginationData<T = any> extends Meta {
  list?: Array<T>;
}

/**
 * 分页数据默认值
 */
export const PaginationDataDefault: PaginationData = {
  page: 1,
  size: 10,
  total: 0,
  list: []
};
