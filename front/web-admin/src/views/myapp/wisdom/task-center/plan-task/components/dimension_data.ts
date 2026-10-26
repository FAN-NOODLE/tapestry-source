// 标签任务
export const tag_task = [
  {
    name: "导入顾客名单",
    value: 3,
    choose_staff: false, // 是否需要选择员工的选项
    down_url: "ba-performance/plan-task-download",
    down_text: "下载导入客户模板",
    has_condition: false,
    // 是否显示任务计划状态
    show_task_status: true
  },
  {
    name: "导入顾客&导购匹配关系",
    value: 4,
    choose_staff: false, // 是否需要选择员工的选项
    down_url: "ba-performance/plan-task-download",
    down_text: "下载客户&SA导入模板",
    has_condition: false
  },
  {
    name: "用户画像",
    value: 2,
    choose_staff: true, // 是否需要选择员工的选项
    down_url: "",
    has_condition: true,
    condition: [
      {
        name: "客户标签",
        value: 1
      }
    ]
  },
  {
    name: "社群数据",
    value: 1,
    choose_staff: true, // 是否需要选择员工的选项
    down_url: "", // 下载的模板地址顺带判断是否需要下载导入
    has_condition: true, // 判断是否需要选择条件内容
    condition: [
      {
        name: "群人数",
        value: 2
      },
      {
        name: "群内新增人数",
        value: 3
      },
      {
        name: "退群人数",
        value: 7
      }
    ]
  }
];
// 企业群发任务
export const group_task = [
  {
    name: "导入顾客名单",
    value: 3,
    choose_staff: false, // 是否需要选择员工的选项
    down_url: "ba-performance/plan-task-download",
    down_text: "下载导入客户模板",
    has_condition: false
  },
  {
    name: "导入顾客&导购匹配关系",
    value: 4,
    choose_staff: false, // 是否需要选择员工的选项
    down_url: "ba-performance/plan-task-download",
    down_text: "下载客户&SA导入模板",
    has_condition: false
  },
  {
    name: "用户画像",
    value: 2,
    choose_staff: true, // 是否需要选择员工的选项
    down_url: "",
    has_condition: true,
    condition: [
      {
        name: "客户标签",
        value: 1
      },
      {
        name: "全部客户",
        value: 9
      }
    ]
  }
];
// 任务类型
export const task_type = [
  {
    name: "沟通任务",
    value: 0
  },
  {
    name: "朋友圈任务",
    value: 1
  }
];
// 任务群发类型
export const task_group_type = [
  {
    name: "灵活沟通转发",
    value: 2
  },
  {
    name: "企微一键群发",
    value: 1
  },
  {
    name: "标签任务群发",
    value: 0
  }
];
