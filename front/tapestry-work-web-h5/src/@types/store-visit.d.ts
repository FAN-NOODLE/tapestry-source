declare namespace StoreVisit {

  interface IMemberItem {
    member_id: number;
    name: string;
    email: string;
    dept_name: string;
    user_id: string;
    is_new?: boolean
    third_id?: number
  }

  interface IDetail {
    id: number;
    date: string;
    time_span: string;
    third_id: number;
    store_name: string;
    store_code: string;
    remark: string;
    /** 1已预约 2巡店中 3巡店完成 */
    status: number;
    creator: IMemberItem;
    store_manager: IMemberItem[];
    visitor_member: IMemberItem[];
    copy_member: IMemberItem[];
    to_store_time: number;
    visitor_is_show: number;
    store_is_show: number;
    overview: string;
    content: string;
    pdf_url: string;
    updated_at: number;
    is_report: number;
    [key: string]: any
  }


  
  interface IMediaItem {
    id: number;
    created_at: number;
    updated_at: number;
    deleted_at: number;
    visit_id: number;
    member_id: number;
    rel_id: number;
    rel_type: number;
    media_type: number;
    url: string;
    video_second: number;
    [key: string]:any
  }

  interface IFeedbackRecordItem {
    id: number;
    visit_id: number;
    type_name: string;
    content: string;
    media_list: Array<IMediaItem>;
    [key: string]:any
  }


  
  interface IMemberInfo {
    member_id: number;
    name: string;
    email: string;
    dept_name: string;
    user_id: string;
    avatar: string;
  }

  interface ICommentItem {
    id: number;
    visit_id: number;
    member_info: IMemberInfo;
    content: string;
    parent_id: number;
    reply_member_info: IMemberInfo;
    is_del: number;
    media_list: IMediaItem[];
    created_at: number;
    children: ICommentItem[];
  }
  
  interface IMediaItem {
    media_type: number;
    url: string;
    video_second: number;
    video_image_url: string
  }

}