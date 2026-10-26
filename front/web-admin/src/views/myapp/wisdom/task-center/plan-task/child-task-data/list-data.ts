export const topOneList = [
  {
    name: "子任务数",
    desc: "当前任务包含的子任务数量",
    key: "period_subtasks_number"
  },
  {
    name: "子任务完成数",
    desc: "当前任务中，所有触达导购都完成了任务，即算“已完成子任务”",
    key: "subtasks_finish_number"
  },
  {
    name: "子任务完成率",
    desc: "子任务完成数 / 子任务数量 * 100%",
    key: "subtasks_finish_percentage"
  }
];
export const topTwoList = [
  {
    name: "任务总数 (件）",
    desc: "所有导购接收到的所有子任务数量",
    key: "task_number"
  },
  {
    name: "已完成任务数 (件）",
    desc: "所有导购已完成的 所有子任务数量",
    key: "task_finish_number"
  },
  {
    name: "任务完成率",
    desc: "已完成任务数 / 任务总数 * 100%",
    key: "task_finish_percentage"
  }
];
export const topThreeList = [
  {
    name: "触达导购总人数",
    desc: "以“用户画像”为维度的子任务中，收到任务的导购数量（去重）",
    key: "reach_guide_number"
  },
  {
    name: "匹配客户总人数（不去重）",
    desc: "以“用户画像”为维度的子任务中，触达的客户人数（不去重）",
    key: "match_customers_number"
  },
  {
    name: "匹配客户总人数（去重）",
    desc: "以“用户画像”为维度的子任务中，触达的客户人数（去重）",
    key: "match_customers_number2"
  },
  {
    name: "触达群主总人数（去重）",
    desc: "以“社群画像”为维度的子任务中，触达的群主人数（去重）",
    key: "reach_lord_number"
  },
  {
    name: "匹配社群总数",
    desc: "以“社群画像”为维度的子任务中，触达的社群数量",
    key: "match_group_number"
  }
];
