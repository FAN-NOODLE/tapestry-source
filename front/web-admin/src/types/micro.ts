export type IMicroValue = "store" | "shop" | "scrm" | "isg" | "cdp" | "404";

export interface IModuleConfig {
  key: string;
  value: IMicroValue;
  child: Array<IModuleConfig>;
}
export interface IRoutes {
  key: string;
  query: object;
  type: "history" | "new";
}
