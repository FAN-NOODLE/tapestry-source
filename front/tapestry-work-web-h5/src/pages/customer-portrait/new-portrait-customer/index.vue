<template>
  <div class="new-portrait-customer" v-if="!isLoading">
    <div class="header">
      <!-- 用户信息 -->
      <div :class="['user-box', !(isSW && info.typeIcon) && 'box-b-m']">
        <img class="avatar" :src="info.avatar" />
        <div class="content">
          <div class="content-top">
            <span class="name">{{ info.name || info.nick_name }}</span>
            <img
              class="level"
              :src="info.level"
              v-if="info.level"
              :style="{ width: isSW ? '120px' : '86px' }"
            />
            <!--  :href="`tel:${info.mobile}`" -->
            <a v-if="info.mobile" @click="makePhoneCall">
              <img class="phone" :src="baseImgs.phone" />
            </a>
          </div>
          <div
            class="content-bottom"
            @click="remarkOptions.showEditRemark = true"
          >
            <div>{{ info.remark }}</div>
            <img :src="baseImgs.edit" alt="" />
          </div>
        </div>
      </div>

      <img
        v-if="isSW && info.typeIcon"
        class="type-icon"
        :src="info.typeIcon"
      />

      <!-- 详细信息 -->
      <customer-info
        :info="info"
        :openId="ExternalId"
        @renovate="renovate"
      ></customer-info>

      <!-- 消费记录 -->
      <div class="flex-around">
        <div class="module-box" @click="toShopDetail">
          <img class="image" :src="baseImgs.shop" />
          <div class="title">
            <div class="word">消费记录</div>
            <img class="icon" :src="baseImgs.arrowRight" />
          </div>
        </div>
        <!-- 卡券详情 -->
        <div
          class="module-box"
          @click="goCardDetail"
          style="margin-right: 0"
          
        >
          <img
            class="image"
            src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/coupon-icon.png"
          />
          <div class="title">
            <div class="word">卡券详情</div>
            <img
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/arrow-right-simple.png"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 卡片切换 -->
    <van-tabs v-model="active" v-if="info" class="tabs">
      <van-tab v-for="(item, index) in tabs" :key="item.value">
        <template #title>
          <div class="tab-title">
            <span>{{ item.label }}</span>
            <span class="count" v-if="item.count">({{ item.count }})</span>
          </div>
        </template>
        <component
          v-if="active === index"
          :is="item.value"
          :info="info"
          :openId="ExternalId"
          @updateCount="getNum(ExternalId)"
        />
      </van-tab>
    </van-tabs>

    <!-- 编辑备注名 -->
    <van-popup v-model="remarkOptions.showEditRemark" :closeable="true">
      <div class="edit-pop">
        <div class="tit">编辑备注名</div>
        <van-field
          :maxlength="15"
          show-word-limit
          v-model="remarkOptions.remark_key"
          clearable
          autofocus
          placeholder="请输入备注名"
        />
        <div class="edit-btn" @click="updateRemark">确定</div>
      </div>
    </van-popup>
  </div>
</template>

<script lang="ts">
import {
  getCustomerDetailApi,
  userDetailTipsNumApi,
  updateRemark,
  getUserAuthApi,
} from "@/api/customer-portrait";
import { Component, Vue } from "vue-property-decorator";
import { Popup, Field } from "vant";
import Binding from "./components/binding.vue";
import CustomerInfo from "./components/customer-info.vue";
import Invitation from "./components/invitation-task.vue";
import MemberInfo from "./components/member-info.vue";
import SelfTask from "./components/self-task.vue";
import Helper from "@/utils/helper";
import customLoading from "@/components/loading/index.vue";
const BASE_URL =
  "https://cdn-cos.coach.com.cn/uploads/";

/** 获取等级图标 */
function getLevelImage(v) {
  let pngNameMap = {};
  if (Helper.isKS()) {
    pngNameMap = {
      "": "2022/v3.0.2/weclient/level-5.png",
      "0": "2023/v3.4-ks/level-0.png",
      "1": "2023/v3.4-ks/level-1.png",
      "2": "2022/v3.0.2/weclient/level-2.png",
      "3": "2023/v3.4-ks/level-3.png",
      "4": "2022/v3.0.2/weclient/level-0.png ",
    };
  } else if (Helper.isSW()) {
    pngNameMap = {
      "1": "2023/v3.4/vip-detail-level.png",
      "2": "2023/v3.4/vip-detail-level.png",
      "3": "2023/v3.4/gold-detail-level.png",
      "4": "2023/v3.4/new-detail-star.png",
    };
  }

  const url = pngNameMap[v.toString()];
  return url ? BASE_URL + url : "";
}

/** 获取新老客标识 */
function getUserTag(v) {
  var pngNameMap = {
    "1": "2023/v3.4/custom_type_1.png",
    "2": "2023/v3.4/custom_type_2.png",
    "3": "2023/v3.4/custom_type_4.png",
    "4": "2023/v3.4/custom_type_3.png",
  };
  return pngNameMap[v.toString()];
}

/** 客户详情内容 */
export type InfoType = {
  /** 是否好友 */
  isRelation: boolean;
  /** 会员卡号 */
  customer_code: number;
  /** 用户id */
  user_id: number;
  /** 印象标签 */
  user_impression_tags: Array<any>;
  /** 用户信息 */
  user_info_add: {
    address: string;
    day: number;
    month: number;
    shoe_size: string;
  };
  union_id: string;
  work_external_user_id: string;
};

@Component({
  components: {
    Binding,
    Invitation,
    MemberInfo,
    SelfTask,
    CustomerInfo,
    [Popup.name]: Popup,
    [Field.name]: Field,
    customLoading,
  },
})
export default class extends Vue {
  /** 是否加载中 */
  private isLoading: boolean = false;
  /** 切换卡片 */
  private tabs = [];

  /** 卡片选中value */
  private active: number = 0;
  private role: number = 0;
  /** 用户信息 */
  private info: any = "";
  /** 修改备注信息 */
  private remarkOptions = {
    showEditRemark: false,
    remark_key: "",
  };

  /** 获取用户id */
  public get ExternalId(): string {
    return this.$store.getters["user/getExternalId"]
    // || "wmUDlGTwAAEmoxILGWDOEfvKtV6iI57g"
  }
  // 是否是sw
  private isSW = Helper.isSW();
  // 是否是ks
  private isKS = Helper.isKS();
  /** 基础图标 */
  private baseImgs = {
    phone: BASE_URL + "2022/v3.0.2/weclient/phone-left.png",
    edit: BASE_URL + "2022/v3.0.2/weclient/edit-icon.png",
    shop: BASE_URL + "2022/v3.0.2/weclient/shopping-icon.png",
    arrowRight: BASE_URL + "2022/v3.0.2/weclient/arrow-right-simple.png",
  };

  /** 生命周期 */
  private created() {
    // 获取地址栏参数nav
    const nav = this.$route.query.nav;
    if (nav) {
      this.active = Number(nav);
    }
    this.getInfo(this.ExternalId);
    this.getAuth();
    // 神策数据上报 - 浏览标签
    this.$track.trackViewLabel({
      user_external_id: this.ExternalId,
    });
  }
  // 获取用户权限
  private async getAuth() {
    try {
      const { data, code, message }: IApiResult = await getUserAuthApi({
        is_new: 0,
      });
      if (code !== 200) return;
      this.role = data.role;
      if (this.role === 5) {
        this.tabs = [
          { label: "会员信息", value: "memberInfo", count: 0 },
          { label: "邀约记录", value: "invitation", count: 0 },
          { label: "自设任务", value: "selfTask", count: 0 },
          { label: "绑定关系", value: "binding", count: 0 },
        ];
      } else {
        this.tabs = [
          { label: "会员信息", value: "memberInfo", count: 0 },
          { label: "邀约记录", value: "invitation", count: 0 },
          { label: "自设任务", value: "selfTask", count: 0 },
        ];
      }
      this.getNum(this.ExternalId);
    } catch (error) {
      console.log(error);
    }
  }

  /** 刷新页面 */
  private renovate() {
    this.getInfo(this.ExternalId);
  }

  /** 组合标签 */
  private makeTag(tags) {
    if (!(tags && tags.length)) return [];

    const tagsList = [];

    tags.forEach((item) => {
      if (item.tags_list && item.tags_list.length) {
        item.tags_list.forEach((tag) => {
          tag.isShowPop = false;
          tagsList.push(tag);
        });
      }
    });
    return tagsList;
  }

  /**
   * 获取顾客邀约和自设任务数量
   */
  async getNum(userId) {
    try {
      // 接口报错
      const { data, code, message }: IApiResult = await userDetailTipsNumApi({
        external_userid: userId,
      });

      if (code !== 200) return this.$toast.fail(message);

      const { reserve_num = 0, sow_task_num = 0 } = data;

      this.tabs[1].count = reserve_num;
      this.tabs[2].count = sow_task_num;
    } catch (error) {
      console.log(error);
    }
  }

  /** 获取用户信息 */
  private async getInfo(userId) {
    this.isLoading = true;
    this.$toast.loading({ forbidClick: true, message: "加载中..." });
    const { code, data, message }: IApiResult = await getCustomerDetailApi({
      external_userid: userId,
      is_work: false,
    });
    this.isLoading = false;
    this.$toast.clear();
    if (code !== 200) return this.$toast.fail(message);

    // 印象标签, 添加折叠标识
    const user_impression_tags = data.user_impression_tags
      ? data.user_impression_tags.map((item) => ({ ...item, expand: true }))
      : [];

    this.info = {
      ...data,
      level: getLevelImage(data.customer_card),
      typeIcon:
        data.customer_type == 0
          ? ""
          : BASE_URL + getUserTag(data.customer_type),
      // 客户标签
      tagList: this.makeTag(data.user_tags),
      // 印象标签, 添加折叠标识
      user_impression_tags,
      isRelation: data.is_relation === 1,
    };
    this.remarkOptions.remark_key = data.remark;
  }

  /** 修改用户备注 */
  private updateRemark() {
    if (!this.remarkOptions.remark_key) {
      this.$toast("请输入备注名");
      return;
    }
    updateRemark({
      external_userid: this.ExternalId,
      remark: this.remarkOptions.remark_key,
    }).then((res: IApiResult) => {
      if (res.code == 200) {
        this.remarkOptions.showEditRemark = false;
        this.renovate();
      } else {
        this.$toast.fail(res.message);
      }
    });
  }

  /** 卡劵/订单跳转详情 */
  private toShopDetail() {
    let query = {
      title: "消费记录",
      union_id: this.info.union_id,
    };
    this.$router.push({
      path: "/portrait-customer-details",
      query: query,
    });
  }
  // 卡券详情
  private goCardDetail() {
    let query = {
      title: "卡券详情",
      union_id: this.info.union_id,
    };
    this.$router.push({
      path: "/portrait-customer-details",
      query: query,
    });
  }

  private makePhoneCall() {
    const { mobile } = this.info
    const self = this
    //使用window.location.href实现功能
    window.location.href = `tel:${mobile}`
    self.$track.trackCallCustomer({
      user_external_id: self.ExternalId
    })
  }
}
</script>
<style lang="scss" scoped>
::v-deep .van-tabs__line {
  background-color: #4a4a4a !important;
}

.new-portrait-customer {
  .header {
    padding: 15px;
    background: #fff;
  }

  .type-icon {
    width: 48px;
    height: 18px;
    margin-left: 75px;
    display: block;
    margin-bottom: 15px;
  }
}

.box-b-m {
  margin-bottom: 30px;
}

.user-box {
  display: grid;
  grid-template-columns: 60px 1fr;
  grid-column-gap: 15px;

  .avatar {
    width: 60px;
    height: 60px;
    border-radius: 50%;
  }

  .content {
    .content-top {
      position: relative;
      display: flex;
      margin-top: 5px;

      .name {
        font-weight: 500;
        color: #33363c;
        font-size: 16px;
      }

      .level {
        width: 120px;
        height: 24px;
        margin-right: 15px;
        margin-left: 5px;
      }

      .phone {
        position: absolute;
        right: 0;
        width: 20px;
        height: 20px;
        margin-left: auto;
      }
    }

    .content-bottom {
      margin-top: 5px;
      font-weight: 400;
      color: #96979e;
      line-height: 17px;
      font-size: 12px;
      display: flex;

      img {
        width: 16px;
        height: 16px;
        margin-left: 5px;
      }
    }
  }
}

.flex-around {
  display: flex;
  justify-content: space-between;

  &::last-child {
    margin-right: 0;
  }

  .module-box {
    position: relative;
    flex: 1;
    margin-right: 2%;
    height: 46px;
    margin-top: 15px;
    background-color: #f6f7f9;
    display: flex;
    align-items: center;

    .image {
      padding-left: 6px;
      width: 28px;
      height: 22px;
    }

    .title {
      display: flex;
      align-items: center;
      padding-left: 6px;

      .word {
        font-size: 14px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        color: #33363c;
        margin-right: 5px;
      }

      .icon {
        position: absolute;
        right: 20px;
        width: 9px;
        height: 14px;
      }
    }
  }
}

.tabs {
  margin-bottom: 15px;

  ::v-deep .van-tabs__nav {
    padding-left: 10px;
    padding-right: 10px;
  }

  .tab-title {
    display: flex;
    align-items: center;
    font-size: 12px;
  }

  .count {
    display: inline-block;
    color: #000;
    font-size: 12px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    border-radius: 6px;
    height: 12px;
    min-width: 12px;
    padding: 0 3px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin-left: 5px;
  }
}

.edit-pop {
  padding: 0 15px;
  width: 310px;

  .tit {
    font-size: 15px;
    font-weight: 500;
    color: #000000;
    line-height: 50px;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 15px;
  }

  .edit-btn {
    font-size: 14px;
    text-align: center;
    width: 250px;
    line-height: 36px;
    border: 1px solid #4a4a4a;
    color: #4a4a4a;
    margin: 15px auto 20px;
  }
}
</style>
