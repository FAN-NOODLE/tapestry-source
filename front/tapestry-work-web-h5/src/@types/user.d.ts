interface ActivityQrResponse {
  error_code: number;
  activity_name: string;
  banner_url: string;
  qr_code_member_info: QrCodeMemberInfo;
  store_list: Store[];
}

interface QrCodeMemberInfo {
  member_id: number;
  user_id: string;
  qr_code: string;
  store_name: string;
  third_id: number;
  sa_name: string;
  avatar: string;
}

interface Store {
  third_id: number;
  store_code: string;
  store_name: string;
  store_address: string;
  city: string;
  city_pin_yin: string;
  distance: number;
}