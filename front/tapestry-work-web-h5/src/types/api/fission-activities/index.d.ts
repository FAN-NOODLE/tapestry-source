// 裂变活动列表接口入参类型
interface IFissionActivityList extends IMeta {
  search_text: string; // 活动名称搜索条件
  status: number; // 活动状态：1-上架中，2-已下架
}

// 分页类型
interface IMetas {
  page: number; // 页码
  size: number; // 每页数据量
  total: number; // 数据总量
}
// 数据统计类型
interface IInfo {
  user_num: number; // 发起人数
  fission_num: number; // 好友裂变数
  run_away_num: number; // 流失人数
  join_group_num: number; // 群成员裂变数
  leave_group_num: number; // 群成员退群数
}
// 行内数据类型
interface IList {
  id: number; // 活动ID
  name: string; // 活动名称
  desc: string; // 活动描述
  status: number; // 活动状态：1-未开始  2-活动已下架 3-已结束  4-活动进行中
  html_url: string; // h5页面地址（link）
  img_url: string; // 默认图片Url地址
}

// 下拉框类型(活动状态内容)
interface ISelect {
  text: string;
  value: number;
}
