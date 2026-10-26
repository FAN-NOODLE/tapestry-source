import COS from "cos-js-sdk-v5";

export enum IUploadType {
  /** 默认 */
  DEFAULT = 'default',
  /** oss */
  OSS = 'oss',
  /** cos */
  COS = 'cos',
}

export interface IUploadOptions {
  /** 上传的类型 */
  type?: IUploadType;
  /** 上传的地址 */
  path?: string;
  /** 文件名称 */
  filename?: string;
  /** 上传指定目录 */
  folder?: string;
  
  ACL?: COS.ObjectACL;
  /** 上传进度 */
  onProgress?: <T>(progress: T) => void;
}

export interface IUploadResult {
  /** 上传的地址 */
  url: string;
  /** 上传的文件名 */
  name?: string;
  /** 上传的文件大小 */
  size?: number;
  /** 上传的文件类型 */
  type?: string;
  /** 上传的文件后缀 */
  ext?: string;
  /** 封面图（视频有） */
  duration?: number;
  cover_url?: string;
}