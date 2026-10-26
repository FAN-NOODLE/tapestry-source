<template>
  <el-dialog
    title="选择企业"
    :visible.sync="show"
    :before-close="handleClose"
    :close-on-click-modal="false"
    width="830px"
  >
    <div class="dialog-body">
      <el-scrollbar
        wrap-class="scrollbar-wrapper"
        style="height:440px;width: 100% !important;"
      >
        <div class="list">
          <div
            class="item"
            v-for="(item, index) in listData"
            :key="item.id"
            :class="[
              item.checked ? 'checked' : '',
              item.has_expire ? 'expire' : 'normal',
            ]"
            @click="handleSelect(item, index)"
          >
            <div class="label" v-show="item.checked">
              <i
                class="el-icon-upload-success el-icon-check icon-success"
                aria-hidden="true"
              ></i>
            </div>
            <div class="item-info">
              <img v-if="item.logo" :src="item.logo" alt="" />
              <div class="logo" v-if="!item.logo">LOGO</div>
              <div class="title-box">
                <div class="title character-hiding-o">{{ item.name }}</div>
                <div class="ex-btn" v-if="item.has_expire">已过试用期</div>
              </div>
            </div>
            <!-- <p class="m-t-10 character-hiding-o">
              小程序商城：{{ item.shop_name }}
            </p> -->
          </div>
        </div>
      </el-scrollbar>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { IEnterpriseItem, ITokenData } from "./interface";
import { getProxyUrl } from "@/utils/proxy-tool";

import Api from "@/api/public/public";
@Component
export default class SelectEnterprise extends Vue {
  @Prop({ default: false, type: Boolean }) show!: boolean;
  @Prop(Function) onClose!: Function;
  /**
   * data
   * */

  private listData: Array<IEnterpriseItem> = [];
  /**
   * methods
   * */

  private handleClose() {
    this.onClose();
  }
  private handleSelect(item: IEnterpriseItem, index: number) {
    if (item.has_expire) {
      return;
    }
    this.listData.forEach((item: IEnterpriseItem) => {
      this.$set(item, "checked", false);
    });
    this.listData[index].checked = true;
    //登陆页面选择企业逻辑 TODO:
    this.getEnterpriseToken(item);
  }

  private async getEnterpriseToken(target: IEnterpriseItem) {
    const { admin_id, enterprise_id } = target;
    localStorage.setItem("REAL_HASH", JSON.stringify(target.hash));
    const { data } = await Api.getEnterpriseToken({ admin_id, enterprise_id });

    this.authorizationNew(data);
  }
  private authorizationNew(data: ITokenData) {
    // 将用户信息存储至localstorage
    // 将用户信息存储至localstorage
    const {
      hash,
      admin_list,
      token,
      permissions_tree,
      role_list,
      permissions_list,
      js_api_config,
    } = data;
    if (admin_list[0].serve_type === 3) {
      this.$message({
        showClose: false,
        message: "该企业使用权限已到期，请联系管理员重新开通",
        type: "error",
      });
      return;
    } else {
      localStorage.setItem("LOGIN_HASH", JSON.stringify(hash));
      localStorage.setItem("ENTERPRISE_ID", admin_list[0].enterprise_id);
      localStorage.setItem("TOKEN", JSON.stringify(token));
      localStorage.setItem("USERINFO", JSON.stringify(admin_list[0]));
      localStorage.setItem("KEYS", JSON.stringify(permissions_list));
      localStorage.setItem("PERMISSION", JSON.stringify(permissions_tree));
      localStorage.setItem("ROLES", JSON.stringify(role_list));
      localStorage.setItem("JSAPICONFIG", JSON.stringify(js_api_config));

      window.location.href = `${window.location.origin}${getProxyUrl("store")}`;
    }
  }
  private async getData() {
    const enterprise_id: number = JSON.parse(
      localStorage.getItem("ENTERPRISE_ID") || "0"
    ) as number;
    const {
      data: { enterprise_admin_list },
    } = await Api.getEnterprise();
    const time: number = new Date().getTime() / 1000;
    this.listData = enterprise_admin_list;
    if (this.listData) {
      this.listData.forEach((item: IEnterpriseItem) => {
        if (item.enterprise_id == enterprise_id) {
          this.$set(item, "checked", true);
        }
        if (time > item.account_expire_at) {
          this.$set(item, "has_expire", true);
        }
      });
    }
  }
  mounted() {
    this.getData();
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .el-dialog__wrapper {
    z-index: 2499 !important;
  }
  .el-dialog {
    z-index: 2500 !important;
  }
  .scrollbar-wrapper {
    overflow-x: hidden !important;
    margin-top: 15px;
  }
  .el-dialog__body {
    padding-top: 0;
  }
  .el-dialog__title {
    font-weight: 700;
  }
}
.logo {
  box-shadow: 0 0 0 1px $color-placeholder;
  border-radius: 50%;
  min-width: 50px;
  min-height: 50px;
  max-width: 50px;
  max-height: 50px;
  color: #999;
  font-size: 12px;
  font-weight: 700;
  line-height: 50px;
  text-align: center;
}
.dialog-body {
  height: 440px;
}
.list {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  flex-direction: row;
  padding: 20px;
  .item {
    position: relative;
    overflow: hidden;
    width: 43.5%;
    border-radius: 4px;
    margin-bottom: 20px;
    cursor: pointer;
    padding: 20px;
    box-shadow: 0 0 0 1px #ebeef5;
    font-size: 16px;
    .label {
      width: 56px;
      height: 36px;
      background-color: $color-main;
      color: #ffffff;
      transform: rotate(45deg);
      line-height: 54px;
      text-align: center;
      position: absolute;
      right: -24px;
      top: -12px;
      .icon-success {
        transform: rotate(-45deg);
      }
    }
    .item-info {
      display: flex;
      line-height: 45px;
      img {
        min-width: 50px;
        min-height: 50px;
        max-width: 50px;
        max-height: 50px;
        border-radius: 50%;
      }
      .title {
        color: $color-text-main;
        margin-left: 10px;
        font-weight: 700;
        margin-top: 5px;
      }
    }
  }
}
.page-right {
  text-align: right;
}
.checked {
  box-shadow: 0 0 0 1px $color-main !important;
}
.title-box {
  display: flex;
  justify-content: space-between;
  width: 100%;
  .ex-btn {
    padding: 0 10px;
    background-color: $color-placeholder;
    color: #fff;
    border-radius: 4px;
    height: 26px;
    line-height: 26px;
    font-size: 12px;
    margin-top: 10px;
  }
}
.expire {
  color: #efeff4 !important;
  cursor: not-allowed !important;
  .title {
    color: $color-placeholder !important;
  }
}
.normal {
  &:after {
    box-shadow: 0 0 0 1px $color-placeholder;
  }
}
.character-hiding-o {
  color: #909399;
}
</style>
