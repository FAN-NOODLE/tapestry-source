
interface IStatusItem {
  label: string;
  value: number;
  key?: string;
  rate_key?: string
}
interface IClassUrlItem {
  class_name: string;
  url: string;
}
interface ITemplateItem {
  template_name: string;
  is_required: boolean;
}

interface IDisplayTaskItem {
  /** id */
  id: number;
  /** 任务名称 */
  display_name: string;
  /** 任务内容及描述 */
  content: string;
  /** 活动开始时间 */
  start_time: number;
  /** 活动结束时间 */
  end_time: number;
  /** 任务状态: 0-未开始，1-进行中，2-已结束 */
  status: number;
  /** 审核状态: 0-待提交，1-待审核，2-已通过，3-不通过 */
  finish_status: number;
  /** 评分 */
  grade: number;
  /** 门店ID */
  store_id: number;
  /** 任务的课程学习 */
  class_url?: IClassUrlItem[];
  /** 任务的模板内容设置 */
  template_info?: ITemplateItem[];
}

interface ICommentItem {
  /** 备注id */
  id?: number;
  /** 备注的员工id */
  member_id: number;
  /** 备注的员工工号 */
  member_user_id: string;
  /** 备注的员工昵称 */
  member_name: string;
  /** 员工头像 */
  avatar: string;
  /** 是否为VM人员，0为否，1为是 */
  is_vm: number;
  /** 是否自己备注，0否，1是 */
  is_self: number;
  /** 创建时间 */
  comment_time: number;
  /** 备注内容 */
  content: {
    text?: string;
    image_url?: Partial<IImageItem>[];
  },
  /** 门店编码 */
  store_code: string;
  /** 0 默认展示，1为删除，2为编辑 */
  action: number;
  /** 回复备注的员工昵称 */
  reply_name?: string;
}
interface IImageItem {
  /** 图片id */
  id?: number;
  /** 图片地址 */
  url: string;
  /** 图片的标识名称 */
  identifier?: string;
  /** 是否通过 0 通过 1 未通过 */
  status?: number; 
  /** VM的备注数 */
  vm_num?: number;
  /** 备注的总数 */
  num?: number;
  /** 是否查看过VM备注的气泡 */
  is_look?: number;
  /** 备注 */
  comment?: ICommentItem[];
  /** 是否为草稿 */
  is_draft?: number;
  /** 上传内容 */
  _content?: string;
  /** 上传状态 uploading | done | error |*/
  _upload_status?: string;
  /** 审核状态:  0 | 1  */
  _status: number;
  /** 上传状态文案 */
  _message?: string;
}

interface ITempItem {
  is_required: boolean;
  template_name: string;
  image_list?: IImageItem[]
}

interface IStoreInfo {
  /** 分配表id */
  allot_id: number;
  /** 任务完成状态,0为未提交，1为待审核，2为未通过，3为已通过 */
  finish_status: number;
  /** VM陈列的图片 */
  image_list: IImageItem[];
  /** 所属VM 导购ID */
  owner_member_id: number;
  /** 所属VM 导购昵称 */
  owner_member_name: string;
  /** 所属VM 导购工号 */
  owner_member_user_id: string;
  /** 所属VM 初审人员ID */
  special_user_id: string;
  /** 所属VM 初审人员名称 */
  special_name: string;
  /** 所属VM 是否初审 true false */
  is_special_check: boolean;
  /** 是否需要VM Corporate审核 */
  is_corporate_check: boolean;
  /** 评分 */
  grade: number;
  /** 整体评价 */
  review_remark: string;
  /** 门店编码 */
  store_code: string;
  /** 门店名称 */
  store_name: string;
  /** 门店ID */
  store_id?: string;
  /** 提交时间 */
  submit_time: number;
  /** 审核时间 */
  review_time: number;
}

interface ITaskInfo {
  id: number;
  /** 课程学习数据 */
  class_url: IClassUrlItem[],
  /** 任务描述 */
  content: string;
  /** 任务创建时间 */
  created_at: number;
  /** 任务名称 */
  display_name: string;
  /** 任务结束时间 */
  end_time: number;
  /** 上下架状态 */
  show_type?: number;
  /** 任务开始时间 */
  start_time: number;
  /** 任务状态： 0 未开始 1 进行中 2 已结束 */
  status: number;
  /** 指定门店范围 */
  store_type?: number;
  /** 任务模板数据 */
  template_info: ITempItem[]
}

interface IDisplayTaskDetail {
  store_info: IStoreInfo,
  task_info: ITaskInfo
}

interface IVmItem {
  /** 导购ID */
  member_id: number;
  /** 导购工号 */
  user_id: string;
  /** 导购昵称 */
  member_name: string;
}