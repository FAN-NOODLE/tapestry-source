
/** 任务状态选项 */
export const statusOptions: IStatusItem[] = [
  { label: '全部任务', value: -1 },
  { label: '待开始', value: 0 },
  { label: '进行中', value: 1 },
  { label: '已结束', value: 2 },
]

/** 任务审核状态选项 */
export const auditStatusOptions: IStatusItem[] = [
  { label: '全部', value: -1, key: 'all_num', rate_key: 'all_rate' },
  { label: '待提交', value: 0, key: 'wait_num', rate_key: 'wait_rate' },
  { label: '待审核', value: 1, key: 'check_num', rate_key: 'check_rate' },
  { label: '待复审', value: 4, key: 'double_check_num', rate_key: 'double_check_rate' },
  { label: '未通过', value: 2, key: 'fail_num', rate_key: 'fail_rate' },
  { label: '已通过', value: 3, key: 'success_num', rate_key: 'success_rate' },
]

const sortOptions = JSON.parse(JSON.stringify(auditStatusOptions))
sortOptions.sort((a, b) => a.value - b.value)

export const statusName = sortOptions.filter(item => item.value>=0).map(item => item.label + '门店数')