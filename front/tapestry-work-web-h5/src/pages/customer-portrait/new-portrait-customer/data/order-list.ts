const orderMapList = {
  sw: new Map([
    [
      "last_purchase_store_name",
      { label: "最近购买柜台：", value: "", unit: "" },
    ],
    [
      "last_12m_purchase_amount",
      { label: "近12月消费金额：", value: "", unit: "￥" },
    ],
    [
      "upgrade_need_amount",
      { label: "升下一级所需金额：", value: "", unit: "￥" },
    ],
    ["upgrade_need_points", { label: "升级所需积分：", value: "", unit: "" }],
    ["total_valid_points", { label: "可用积分：", value: "", unit: "" }],
    ["last_purchase_date_str", { label: "最近购买时间：", value: "", unit: "" }],
    ["upgrade_date_time", { label: "最近升级时间：", value: "", unit: "" }],
    ["expiry_date", { label: "会籍到期时间：", value: "", unit: "" }],
  ]),
  ks: new Map([
    [
      "last_purchase_store_name",
      { label: "最近购买柜台：", value: "", unit: "" },
    ],
    ["upgrade_need_points", { label: "升级所需积分：", value: "", unit: "" }],
    ["total_valid_points", { label: "可用积分：", value: "", unit: "" }],
    [
      "last_36m_purchase_amount",
      { label: "近3年消费金额：", value: "", unit: "￥" },
    ],
    [
      "last_purchase_date_str",
      { label: "最近购买时间：", value: "", unit: "" },
    ],
    ["upgrade_date_time", { label: "最近升级时间：", value: "", unit: "" }],
    ["expiry_date", { label: "会籍到期时间：", value: "", unit: "" }],
  ]),
};
export default orderMapList;
