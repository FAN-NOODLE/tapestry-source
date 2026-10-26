export interface IPublicParam {
  code: number;
  data?: any;
  message: string;
}
// 上传图片/视频之前返回参数
export interface IUploadImageBeforeParam {
  lastModified?: number;
  lastModifiedDate?: any;
  name?: string | any;
  size?: number| any;
  type?: string;
  uid?: number;
  webkitRelativePath?:string;
}