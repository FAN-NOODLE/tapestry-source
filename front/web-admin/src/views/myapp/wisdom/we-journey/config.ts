/** weJourney 默认参数配置 */
export const weJourneyConfig = [
  {
    name: "导购积分",
    list: [
      {
        icon: "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/wej1.png",
        name: "积分规则设置",
        describe: "设置不同维度积分规则及能力图谱规则",
        url: "/myapp/wisdom/we-journey/integral-set-list"
      },
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/wej2.png",
        name: "积分兑换数据",
        describe: "导购积分兑领明细，可手动补发积分",
        url: "/myapp/wisdom/we-journey/integral-exchange"
      },
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/wej3.png",
        name: "等级积分规则说明",
        describe: "让员工了解积分及等级升级的规则",
        url: "/myapp/wisdom/we-journey/integral-explain"
      }
    ]
  },
  {
    name: "导购等级",
    list: [
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/wej4.png",
        name: "等级规则设置",
        describe: "设置不同导购等级的升级积分条件",
        url: "/myapp/wisdom/we-journey/grade-set-list"
      },
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/wej5.png",
        name: "导购等级数据",
        describe: "导购积分、等级、奖章数据的汇总展示",
        url: "/myapp/wisdom/we-journey/grade-statistics"
      }
    ]
  },
  {
    name: "导购奖章",
    list: [
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/wej7.png",
        name: "奖章规则配置",
        describe: "设置奖章样式及奖章规则配置",
        url: "/myapp/wisdom/we-journey/medal-set-list"
      },
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/wej8.png",
        name: "奖章领取数据",
        describe: "导购奖章获取明细，可手动补发奖章",
        url: "/myapp/wisdom/we-journey/medal-statistics"
      }
    ]
  }
];

/** 积分规则类型 */
export type integralType = {
  id: number; // 维度类型
  relationId: number; // 关联能力图谱id
  relationStatus?: number; // 关联能力图谱开启状态
  listId?: number; // 关联列表id
  icon: string; // icon 图标
  describe: string; // 描述
  name: string; // 类型名称
};

/** 积分规则维度 */
export const integralConfig: Array<integralType> = [
  {
    id: 1, // 维度id
    relationId: 0,
    listId: 0,
    relationStatus: 2,
    icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/in1.png",
    name: "邀约维度",
    describe: "可配置邀约维度下的类目规则"
  },
  {
    id: 2, // 维度id
    relationId: 0,
    listId: 0,
    relationStatus: 2,
    icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/in2.png",
    name: "任务维度",
    describe: "可配置任务维度下的类目规则"
  },
  {
    id: 3, // 维度id
    relationId: 0,
    listId: 0,
    relationStatus: 2,
    icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/in3.png",
    name: "客户互动维度",
    describe: "可配置客户互动维度下的类目规则"
  },
  {
    id: 4, // 维度id
    relationId: 0,
    listId: 0,
    relationStatus: 2,
    icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/in4.png",
    name: "培训维度",
    describe: "可配置培训维度下的类目规则"
  },
  {
    id: 5, // 维度id
    relationId: 0,
    listId: 0,
    relationStatus: 2,
    icon:"https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/we-journey/in5.png",
    name: "线上能力",
    describe: "可配置线上能力维度下的类目规则"
  }
];

/** 积分维度下拉 */
export const interfaceDimension = [
  { value: 0, label: "全部" },
  { value: 1, label: "邀约维度" },
  { value: 2, label: "任务维度" },
  { value: 4, label: "培训维度" },
  { value: 3, label: "客户互动维度" },
  { value: 5, label: "线上能力维度" }
];

/** 积分来源 */
export const integralSource = [
  { value: 0, label: "全部" },
  { value: 1, label: "规则统计" },
  { value: 2, label: "手动补发" },
  { value: 3, label: "奖章奖励" }
];

/** 规则类目类型 */
export type categoryType = {
  value: number; // 值
  label: string; // 名称
  ruleType: number; // 规则类型 1. 次数 2. 条件
  disabled: boolean; // 是否置灰
  isAstrict: boolean; // 是否有上限
  astrictTips?: string; // 上限制文案
};

/** 积分规则-规则类目 */
export const integralCategory: Array<categoryType> = [
  {
    value: 1,
    label: "邀约数（每次）",
    ruleType: 1,
    disabled: false,
    isAstrict: true,
    astrictTips: "导购每天该类型积分最多可增加"
  },
  {
    value: 2,
    label: "邀约到店率（每周）",
    ruleType: 2,
    disabled: false,
    isAstrict: false
  },
  {
    value: 3,
    label: "邀约成功率（每月）",
    ruleType: 2,
    disabled: false,
    isAstrict: false
  },
  {
    value: 4,
    label: "邀约成功（每次）",
    ruleType: 1,
    disabled: false,
    isAstrict: false
  },
  {
    value: 5,
    label: "邀约占比（每月）",
    ruleType: 2,
    disabled: false,
    isAstrict: false
  },
  {
    value: 6,
    label: "邀约到店确认（每次）",
    ruleType: 1,
    disabled: false,
    isAstrict: true,
    astrictTips: "导购每天该类型积分最多可增加"
  },
  {
    value: 17,
    label: "勾选【新客 老客】按钮（每次）",
    ruleType: 1,
    disabled: false,
    isAstrict: true
  },
  {
    value: 18,
    label: "勾选【是否到店购买】按钮（每次）",
    ruleType: 1,
    disabled: false,
    isAstrict: true
  }
];

/** 任务规则-规则类目 */
export const taskCategory: Array<categoryType> = [
  {
    value: 7,
    label: "任务完成（每次）",
    ruleType: 1,
    disabled: false,
    isAstrict: false
  },
  {
    value: 8,
    label: "任务完成率（每月）",
    ruleType: 2,
    disabled: false,
    isAstrict: false
  },
  {
    value: 9,
    label: "任务触达（每人）",
    ruleType: 1,
    disabled: false,
    isAstrict: true
  }
];

/** 客户互动规则-规则类目 */
export const customerCategory: Array<categoryType> = [
  {
    value: 10,
    label: "新增好友数（每人）",
    ruleType: 1,
    disabled: false,
    isAstrict: true
  },
  {
    value: 11,
    label: "客户流失数（每人）",
    ruleType: 2,
    disabled: false,
    isAstrict: false
  },
  {
    value: 12,
    label: "聊天总数（每月）",
    ruleType: 2,
    disabled: false,
    isAstrict: false
  },
  {
    value: 13,
    label: "已回复聊天占比（每月）",
    ruleType: 2,
    disabled: false,
    isAstrict: false
  },
  {
    value: 14,
    label: "完成播种（每次）",
    ruleType: 1,
    disabled: false,
    isAstrict: true,
    astrictTips: "导购每天该类型积分最多可增加"
  }
];

/** 培训维度规则-规则类目 */
export const trainCategory: Array<categoryType> = [
  {
    value: 15,
    label: "积分等比兑换",
    ruleType: 1,
    disabled: false,
    isAstrict: false
  }
];

/** 线上能力规则-规则类目 */
export const onLineCategory: Array<categoryType> = [
  {
    value: 16,
    label: "线上销售业绩率（每月）",
    ruleType: 1,
    disabled: false,
    isAstrict: false
  },
  // todo: 增加线上销售金额
  {
    value: 19,
    label: "线上销售金额（每月）",
    ruleType: 1,
    disabled: false,
    isAstrict: false
  }
];

export const TYPE_CATEGORY: Array<Array<categoryType>> = [
  integralCategory,
  taskCategory,
  customerCategory,
  trainCategory,
  onLineCategory
];

/** option 类型 */
export type OptionType = {
  value: number;
  label: string;
};

/** 高于-达到 枚举值 */
export const condition = [
  { value: 2, label: "高于" },
  { value: 3, label: "达到" }
];

/** 平均值-百分比 枚举值 */
export const condition1 = [
  { value: 2, label: "平均值" },
  { value: 1, label: "百分比" }
];

/** 平均值-百分比 枚举值 */
export const condition2 = [
  { value: 2, label: "平均值" },
  { value: 1, label: "数值" }
];

/** 指标规则-目标率 枚举值 */
export const condition3 = [
  { value: 1, label: "top榜" },
  { value: 2, label: "高于" },
  { value: 3, label: "达到" },
  { value: 4, label: "低于" }
];

/** 指标规则-能力指标 枚举值 */
export const condition4 = [
  { value: 1, label: "SS", disabled: false },
  { value: 2, label: "S", disabled: false },
  { value: 3, label: "A+", disabled: false },
  { value: 4, label: "A", disabled: false },
  { value: 5, label: "A-", disabled: false }
];

/** 邀约目标达成率 枚举值  time_rule 时间规则: 1每次(人) 2每天 3每周 4每月 5每季度 6每年*/
export const condition5 = [
  { value: 1, label: "邀请目标达成率（每月）", time_rule: 4 },
  { value: 2, label: "任务完成率（每月）", time_rule: 4 },
  { value: 3, label: "总客户互动积分排行榜（每月）", time_rule: 4 },
  { value: 4, label: "总培训积分排行榜（每月）", time_rule: 4 },
  { value: 5, label: "目标达成率（每月）", time_rule: 4 }
];

/** 能力指标 */
export const abilityOptions = [
  { value: 1, label: "邀约能力", disabled: false },
  { value: 2, label: "任务能力", disabled: false },
  { value: 3, label: "学习能力", disabled: false },
  { value: 4, label: "客户互动能力", disabled: false },
  { value: 5, label: "线上销售能力", disabled: false }
];

/** 奖章来源 */
export const sourceOptions = [
  { value: 1, label: "手动上传" },
  { value: 2, label: "规则统计" }
];

/** 奖章类型 */
export const medalTypeOptions = [
  { value: 1, label: "普通奖章" },
  { value: 2, label: "隐藏奖章" }
];

/** 积分奖励规则类型 */
export type ItemType = {
  /** 1:次数条件 2:高于3：达到 */
  rule_if: number;
  rule_if1: number;
  /** 1.加积分 2.减积分 */
  integral_type: number;
  /** 规则值 */
  rule_value?: string;
  /** 积分 */
  integral?: string;
  /** store_type:0默认 1奥莱 2零售 */
  store_type?: number;
};

/** 奖章规则-其它奖励 */
export const otherOptions = [
  { value: 1, label: "无其他奖励" },
  { value: 2, label: "仅奖励积分" },
  { value: 3, label: "仅奖励实物奖品" },
  { value: 4, label: "奖励积分 + 实物奖品" }
];

/** 奖章规则-是否有效 */
export const validOptions = [
  { value: 1, label: "有效" },
  { value: 2, label: "失效" }
];

/** 奖章规则-启用禁用 */
export const enableOptions = [
  { value: 1, label: "启用" },
  { value: 2, label: "禁用" }
];

/** 能力图谱提示文案 */
export const atlasTips = [
  "注：邀约目标达成率 = 邀约金额 / 个人销售金额(CS2000) x 100%",
  "注：任务完成率 = 完成子任务数 / 总子任务数 x 100%",
  "",
  "",
  "注:目标达成率 = 线上平台销售额1个人销售金额 (CS2000+Link to EC) X 100%"
];
// todo:每天，每周，每月，每年下拉
export const integralType = [
  { label: "每天", value: 1 },
  { label: "每周", value: 2 },
  { label: "每月", value: 3 },
  { label: "每年", value: 4 }
];
// todo: 互动培训维度增加top百分比
export const condition3_1 = [
  { value: 1, label: "top榜" },
  { value: 5, label: "top榜占比" },
  { value: 2, label: "高于" },
  { value: 3, label: "达到" },
  { value: 4, label: "低于" },
];
/** 金币兑礼默认参数配置 */
export const coinsConfig = [
  {
    name: "导购金币",
    list: [
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/v4.0.8/u69.png",
        name: "金币明细数据",
        describe: "导购金币明细，可手动补发金币",
        url: "/myapp/wisdom/we-journey/coins-index/gold-statistics"
      },
      {
        icon:"https://cdn-cos.coach.com.cn/uploads/2023/v4.0.8/u31.png",
        name: "兑礼中心入口/金币设置",
        describe: "设置兑礼中心入口开关、金币过期时间",
        url: "/myapp/wisdom/we-journey/coins-index/gold-set"
      }
    ]
  }
];