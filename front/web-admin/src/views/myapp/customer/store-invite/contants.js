// 顾客状态类型
const customerStatusList = [
  {
    key: 0,
    label: "--"
  },
  {
    key: 1,
    label: "新客"
  },
  {
    key: 2,
    label: "老客"
  },
  {
    key: 3,
    label: "vip客"
  },
  {
    key: 4,
    label: "沉睡客"
  }
];
// 预约单来源
const orderOriginList = [
  {
    key: 0,
    label: "全部"
  },
  {
    key: 1,
    label: "到店预约"
  },
  {
    key: 2,
    label: "邀请函"
  }
];
// 是否到店购买 & 是否评价 & 是否成单
const ifOrNotLists = [
  {
    key: 0,
    label: "全部"
  },
  {
    key: 1,
    label: "是"
  },
  {
    key: 2,
    label: "否"
  },
];
// 是否到店购买
const isStoreOptions = [
  {
    key: -1,
    label: "全部"
  },
  {
    key: 1,
    label: "是"
  },
  {
    key: 2,
    label: "否"
  },
  {
    key: 0,
    label: "未选择"
  }
];
const statusList = [
  {
    key: -1,
    label: "全部"
  },
  {
    key: 0,
    label: "进行中"
  },
  {
    key: 3,
    label: "已取消"
  },
  {
    key: 4,
    label: "已过期"
  },
  {
    key: 2,
    label: "邀约完成"
  },
  {
    key: 1,
    label: "邀约成功"
  }
];
// todo: sw单独邀约状态
const swStatusList = [
  {
    key: -1,
    label: "全部"
  },
  {
    key: 0,
    label: "邀约生效中"
  },
  {
    key: 3,
    label: "已取消"
  },
  {
    key: 4,
    label: "邀约匹配失败"
  },
  {
    key: 2,
    label: "邀约完成"
  },
  {
    key: 1,
    label: "邀约成功"
  }
];
const statusMap = {
  0: "进行中",
  1: "邀约成功",
  2: "邀约完成",
  3: "已取消",
  4: "已过期"
};
// todo:sw单独状态map
const swStatusMap = {
  0: "邀约生效中",
  1: "邀约成功",
  2: "邀约完成",
  3: "已取消",
  4: "邀约匹配失败"
};
export {
  customerStatusList,
  orderOriginList,
  ifOrNotLists,
  statusList,
  swStatusList,
  statusMap,
  swStatusMap,
  isStoreOptions
};
