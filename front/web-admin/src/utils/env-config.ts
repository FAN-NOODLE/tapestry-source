const APIS = [
  {
    name: "LOCAL",
    origin: "localhost/",
    domain: "https://tapestry-api-uat.coach.com.cn",
  },
  {
    name: "DEV",
    origin: "tapestry-dev.vchangyi.com",
    domain: "https://tapestry-dev.vchangyi.com",
  },
  {
    name: "UAT",
    origin: "tapestry-web-uat.coach.com.cn",
    domain: "https://tapestry-api-uat.coach.com.cn",
  },
  {
    name: "PROD",
    origin: "we.coach.com.cn",
    domain: "https://webk.coach.com.cn",
  },
  {
    name: "SW_UAT",
    origin: "tapestry-web-uat.stuartweitzman.cn",
    domain: "https://tapestry-api-uat.stuartweitzman.cn",
  },
  {
    name: "SW_PROD",
    origin: "we.stuartweitzman.cn",
    domain: "https://webk.stuartweitzman.cn",
  },
  {
    name: "KS_PROD",
    origin: "we.katespade.cn",
    domain: "https://webk.katespade.cn",
  },
  {
    name: "KS_UAT",
    origin: "tapestry-web-uat.katespade.cn",
    domain: "https://tapestry-api-uat.katespade.cn",
  },
];
const API = APIS.find((item) => {
  return location.href.indexOf(item.origin) > -1;
});

const Domain = API ? API.domain : APIS[0].domain;
const VUE_APP_API_ROOT_RE = `${Domain}/gw-console/gw-console/v1/`;
const VUE_APP_API_ROOT = `${Domain}/gw-scrm/gw-scrm/`;
const VUE_APP_API_VERSION = `${Domain}/gw-scrm/gw-scrm/`;
const VUE_APP_API = `${Domain}/`;
export { API, VUE_APP_API_ROOT, VUE_APP_API, VUE_APP_API_ROOT_RE, VUE_APP_API_VERSION };
