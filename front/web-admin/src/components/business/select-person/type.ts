// 元素节点
export type Member = {
  name: string; // 用户、部门名称
  id: string | number; // 用户、部门id
  leaf: boolean; // 是否子节点
  is_dp: boolean; // 是否部门
  mem_subscribe_status: boolean; // 是否加入企业， 部门默认为true
  third_id?: number;
  checked?: boolean; // 用于搜索列表
  store_type?: number;
};

// 部门节点
export type DpType = {
  dp_id: number; // 部门id
  dp_name: number; // 部门名称
  third_id: number; // 用于查询层级
  store_type?: number; //门店类型
};

// 员工节点
export type PersonType = {
  user_name: string; // 用户名称
  user_id: string; // 用户id
  member_id: number; // 主键id
  mem_subscribe_status: number; // 是否加入企业 1 加入
};

// 部门搜索条件
export type DpParams = {
  is_top?: number; // tree 层级
  third_id?: number;
  dp_name?: string; // 搜索名称
};

// 员工搜索条件
export type PersonParams = {
  page: number;
  size: number;
  search: {
    dp_id?: string | number;
    mem_username: string;
    userid: string;
    mobile: string;
  };
};

// tree 类型
export type treeApiType = {
  expand_key: Array<number>; // 默认展开
  setNode: () => void; // 设置选中节点
  append: (node: Member, target: Member) => void; // 往节点中追加节点
};

export type guideNameListType = {
  [userid in "userid" | "user_id"]: string;
} & {
  id: number; // 主键id
  name: string; // 用户名称
};

// emit 返回值类型
export type callBackType = {
  guide_name_list: Array<guideNameListType | Member>;
  rawData: Array<Member>; // 原始数据
};

// 组件使用
export type selectPerson = {
  visible: boolean;
  rawData: Array<Member | guideNameListType>; // 原始数据
  selectedData: Array<Member | guideNameListType>; // 选中ids
  submit: (data: callBackType) => void;
};
