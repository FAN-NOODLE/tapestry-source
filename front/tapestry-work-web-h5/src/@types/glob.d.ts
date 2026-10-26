/**
 * 通用接口返回数据类型定义
 */
interface IApiResult<T = any> {
  code?: number;
  message?: string;
  data?: T;
  request_id?: string;
}
/**
 * 分页参数类型定义
 */
interface IMeta {
  page?: number; //当前页
  size?: number; //每页数量
  total?: number; //总数
}
/**
 * 分页接口接收数据类型定义
 */
interface IPaginationData<T = any> extends IMeta {
  list?: Array<T>;
}
interface ICyObject {
  [propName: string]: any;
}

interface Window {
  wx: {
    config: (config: unknown) => void;
    agentConfig: (config: unknown) => void;
    ready: (config: unknown) => void;
    error: (config: unknown) => void;
    previewImage: (config: unknown) => void;
    invoke: (config1: unknown, config2?: unknown, config3?: unknown) => void;
    onHistoryBack: (config?: unknown) => void;
    closeWindow: (config: unknown) => void;
    openEnterpriseChat: (config: unknown) => void;
    openUserProfile: (config: unknown) => void;
    chooseImage: (config: unknown) => void;
    checkJsApi: (config: unknown) => void;
    getLocalImgData: (config: unknown) => void;
    setClipboardData: (config: unknown) => void;
    onNetworkStatusChange: (config: unknown) => void;
    getNetworkType: (config: unknown) => void;
    previewFile: (config: unknown) => void;
    getClipboardData: (config: unknown) => void;
  };
  WeixinJSBridge: any
  MSStream: any
}
// 扩展 sensors 类型定义，添加 registerPropertyPlugin 方法
// declare module 'sa-sdk-javascript' {
//   function registerPropertyPlugin(data: any): void;
  
// }
