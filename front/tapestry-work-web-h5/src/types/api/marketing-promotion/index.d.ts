// 素材库列表参数接口
interface ISourceList extends IMeta {
  type: string; // file、image、text、video、link
  keyword: string;
  group_id: number;
}
interface IGoodlist extends IMeta {
  sku_name: string;
}
interface IWxCouponList extends IMeta {
  search_text: string;
}

interface IPosterParams {
  type: number;
  operator_type: number;
  params: string;
}

interface IMakeMediaParams {
  source_id: number;
  source_type: string;
  source_url: string;
  source_name: string;
}
interface ISource {
  /**视频的cos地址 */
  url: string;
  /**是否原图，0否，1是	 */
  is_origin: number;
  /** 名称 */
  name: string;
  /** 后缀 */
  ext: string;
  /** 大小*/
  size: number;
}

/**保存素材 */
interface ISaveMaterial {
  /**素材有效期开始时间 */
  start_time: number;

  /**素材有效期结束时间 */
  end_time: number;

  /**素材的标签 */
  category_ids?: number[];

  /**素材类型,1文本素材，2图片素材，5视频素材 */
  source_type: number;
  /**素材的分组 */
  group_id?: number;
  /**图片/视频素材的资源信息 */
  sources?: Array<ISource>;
  /**文本素材的内容 */
  content?: string;
  /**视频名称 */
  video_name?: string;
  // 展示类型
  show_type?: any
}
