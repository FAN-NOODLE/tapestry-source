import { MessageBox, Message } from "element-ui";
import Api from "@/api/setting/Account";
// import GoodsApi from "@/api/goods-manage/goods";
// 启用vs禁用
export function disOrOpenAccount(id: any, isOpen: any) {
  console.log("启用禁用=====" + id);
  return new Promise(suc => {
    Api.accountStatus({
      id: id,
      status: isOpen ? 1 : 2
    }).then(() => {
      if (isOpen) {
        Message.success("启用成功");
      } else {
        Message.success("禁用成功");
      }
      suc({});
    });
  });
}
// 删除
export function deleteAccount(id: number) {
  return new Promise(suc => {
    MessageBox.confirm("确认要删除该账号吗？删除后将无法恢复", "删除账号", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    }).then(() => {
      Api.accountDel({
        id: id
      }).then(() => {
        Message.success("删除成功");
        suc({});
      });
    });
  });
}
