// 分流hash
export function getLoginHash() {
  // 线上环境 正式、测试
  let hash = localStorage.getItem("LOGIN_HASH") || "";
  if (hash) {
    hash = JSON.parse(hash);
  }
  return hash;
}

// nginx代理地址获取
/**
 * @param name
 * @
 *
 * **/
export function getProxyUrl(name: string): string {
  // 开发环境
  // return (
  //   "/saas/" +
  //   (process.env.NODE_ENV === "production" ? hash : "") +
  //   (isNewDevops ? "/" : "/") +
  //   name
  // )
  return `/saas/${name}`;
}

// 获取对应环境domain
export function getDomain(name: string): string {
  const APIS = [
    {
      name: "DEV",
      origin: "tapestry-dev.vchangyi.com",
      domain: "https://tapestry-dev-h5.vchangyi.com",
    },
    {
      name: "UAT",
      origin: "tapestry-web-uat.coach.com.cn",
      domain: "https://tapestry-h5-uat.coach.com.cn",
    },
    {
      name: "PROD",
      origin: "we.coach.com.cn",
      domain: "https://we-h5.coach.com.cn",
    },
    {
      name: "SW_UAT",
      origin: "tapestry-web-uat.stuartweitzman.cn",
      domain: "https://tapestry-h5-uat.stuartweitzman.cn",
    },
    {
      name: "SW_PROD",
      origin: "we.stuartweitzman.cn",
      domain: "https://we-h5.stuartweitzman.cn",
    },
    {
      name: "KS_PROD",
      origin: "we.katespade.cn",
      domain: "https://we-h5.katespade.cn",
    },
    {
      name: "KS_UAT",
      origin: "tapestry-web-uat.katespade.cn",
      domain: "https://tapestry-h5-uat.katespade.cn",
    },
  ];
  const API = APIS.find((item) => {
    return location.href.indexOf(item.origin) > -1;
  });
  const local = {
    store: "//localhost:20001",
    scrm: "//localhost:20000",
    isg: "//localhost:30000",
    shop: "//localhost:40000",
  };
  const domain = API ? API.domain : `${local[name]}`;
  return domain;
}
