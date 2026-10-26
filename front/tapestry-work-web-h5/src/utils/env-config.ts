const APIS = [
  {
    name: "LOCAL",
    origin: "localhost:",
    domain: "https://tapestry-api-uat.coach.com.cn",
    server_url: "https://analytics-sink.coach.com.cn/sa?project=default",
    app_name: "Coach_Wecom",
    robot_key:"ab4ce55e-264d-422b-92f3-8b20cb8d8797",
    welearning_url: "https://ant.strong365.com",
    welearning_token: "REDACTED_FOR_SOURCE_SHARING",
    welearning_company_code: "wejourney_test",
    cosUrl: 'https://cdn-cos-uat.coach.com.cn/static'
    // welearning_url: "https://cust.strong365.com",
    // welearning_token: "REDACTED_FOR_SOURCE_SHARING",
    // welearning_company_code: "coach",
  },
  {
    name: "DEV",
    origin: "tapestry-dev.vchangyi.com",
    domain: "https://tapestry-api-uat.coach.com.cn"
  },
  {
    name: "UAT",
    origin: "tapestry-web-uat.coach.com.cn",
    domain: "https://tapestry-api-uat.coach.com.cn",
    server_url: "https://analytics-sink.coach.com.cn/sa?project=default",
    app_name: "Coach_Wecom",
    robot_key:"ab4ce55e-264d-422b-92f3-8b20cb8d8797",
    // welearning_url: "https://ant.strong365.com",
    // welearning_token: "REDACTED_FOR_SOURCE_SHARING",
    // welearning_company_code: "wejourney_test",
    welearning_url: "https://cust.strong365.com",
    welearning_token: "REDACTED_FOR_SOURCE_SHARING",
    welearning_company_code: "coach",
    appid: 'wxbd40b04da240ef2b',
    cosUrl: "https://cdn-cos.coach.com.cn/static"
  },
  {
    name: "PROD",
    origin: "we.coach.com.cn",
    domain: "https://webk.coach.com.cn",
    server_url: "https://analytics-sink.coach.com.cn/sa?project=production",
    app_name: "Coach_Wecom",
    robot_key:"daa66549-17ac-40b0-9622-40ca8df46615",
    welearning_url: "https://cust.strong365.com",
    welearning_token: "REDACTED_FOR_SOURCE_SHARING",
    welearning_company_code: "coach",
    appid: 'wx19ba52f691c4d289',
    cosUrl: "https://cdn-cos.coach.com.cn/static"
  },
  {
    name: "SW_UAT",
    origin: "tapestry-web-uat.stuartweitzman.cn",
    domain: "https://tapestry-api-uat.stuartweitzman.cn",
    server_url: "https://analytics-sink.coach.com.cn/sa?project=sw_test",
    app_name: "SW_Wecom",
    robot_key:"ab4ce55e-264d-422b-92f3-8b20cb8d8797"
  },
  {
    name: "SW_PROD",
    origin: "we.stuartweitzman.cn",
    domain: "https://webk.stuartweitzman.cn",
    server_url: "https://analytics-sink.coach.com.cn/sa?project=sw_prd",
    app_name: "SW_Wecom",
    robot_key:"daa66549-17ac-40b0-9622-40ca8df46615"
  },
  {
    name: "KS_PROD",
    origin: "we.katespade.cn",
    domain: "https://webk.katespade.cn",
    server_url: "https://analytics-sink.coach.com.cn/sa?project=ks_prd",
    app_name: "KS_Wecom",
    robot_key:"daa66549-17ac-40b0-9622-40ca8df46615"
  },
  {
    name: "KS_UAT",
    origin: "tapestry-web-uat.katespade.cn",
    domain: "https://tapestry-api-uat.katespade.cn",
    server_url: "https://analytics-sink.coach.com.cn/sa?project=ks_test",
    app_name: "KS_Wecom",
    robot_key:"ab4ce55e-264d-422b-92f3-8b20cb8d8797"
  },
];
const API = APIS.find((item) => {
  console.log(location.href, item.origin);
  return location.href.indexOf(item.origin) > -1;
});
console.log("123", API);
const Domain = API ? API.domain : APIS[0].domain;
const VUE_APP_API_ROOT = `${Domain}/gw-work-scrm/`;
const VUE_APP_API_VERSION = "gw-work-scrm/v1";
export { VUE_APP_API_ROOT, VUE_APP_API_VERSION, API };
