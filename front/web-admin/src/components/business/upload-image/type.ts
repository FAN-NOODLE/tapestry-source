export type ImgType = "single" | "batch";
export interface IFileParam {
  _url?: string; //图片的url
  url?: string; //图片的url
  ext?: string; //文件后缀类型
  size?: number; //文件大小
  name?: string; //文件名
  hash?: string; //文件hash
  is_hover?: boolean;
  progress?: number;
  is_ai?: boolean;
  file?: File;
  relation_goods?: any[]
}
