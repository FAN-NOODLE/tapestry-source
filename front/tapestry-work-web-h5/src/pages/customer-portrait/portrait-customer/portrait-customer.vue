<template>
  <div class="container">
    <div class="part-container" style="padding-bottom: 15px">
      <!-- 头部信息 -->
      <div class="part-container-row">
        <div class="person-row pointer">
          <div class="image">
            <van-image
              :src="customerDetail.avatar"
              width="100%"
              fill
              round
              height="100%"
            />
          </div>
          <div class="info-content">
            <div class="left">
              <p>
                <img
                  class="sex"
                  :src="
                    customerDetail.sex === 1
                      ? require('@/assets/customer-portrait/man-sex.png')
                      : customerDetail.sex === 2
                      ? require('@/assets/customer-portrait/woman-sex.png')
                      : require('@/assets/customer-portrait/unknown-sex.png')
                  "
                />
                <span>{{ customerDetail.nick_name }}</span>
              </p>
              <div class="c-name" @click="editCustomer">
                <!-- 姓名 -->
                <span v-if="customerDetail.name">{{
                  customerDetail.name
                }}</span>
                <span v-else>{{
                  customerDetail.remark || customerDetail.nick_name
                }}</span>
                <img
                  v-if="!customerDetail.name"
                  :src="require('@/assets/customer-portrait/edit-ic.png')"
                />
              </div>
              <!-- 编辑备注名 -->
              <van-popup v-model="showEditRemark" :closeable="true">
                <div class="edit-pop">
                  <div class="tit">编辑备注名</div>
                  <van-field
                    :maxlength="15"
                    show-word-limit
                    v-model="remark_key"
                    clearable
                    autofocus
                    placeholder="请输入备注名"
                  />
                  <div class="edit-btn" @click="updateRemark">确定</div>
                </div>
              </van-popup>
              <div class="grade">
                <!-- 会员等级 -->
                <span>{{ customerDetail.customer_type || "-" }}</span>
                <!-- CTL类型 -->
                <span v-if="!isKS() && !isSW()">{{
                  customerDetail.client_type || "-"
                }}</span>
              </div>
            </div>
            <div class="c-mobile">
              <!-- :href="`tel:${customerDetail.mobile}`" -->
              <a
                v-if="customerDetail.mobile"
                @click="makePhoneCall"
              >
                <img
                  :src="require('@/assets/customer-portrait/phone-ic.png')"
                />
              </a>
              <a @click="phoneTip" v-else href="javascript:;">
                <img
                  :src="require('@/assets/customer-portrait/phone-ic.png')"
                />
              </a>
            </div>
          </div>
        </div>
      </div>
      <!-- 基础信息 -->
      <!-- SW 基础信息 -->
      <div class="part-container-basic" v-if="isSW()">
        <div
          class="part-container-basic-info sw-info-li"
          v-for="(item, index) in swBasicinfoList"
          :key="index"
          @click="openSelPopup(index)"
        >
          <img
            v-if="index > 0"
            :src="require('@/assets/customer-portrait/shequn_arrow_right.png')"
            class="arrow_right"
          />
          <img :src="item.iconUrl" :alt="item.title" />&nbsp;
          <span v-if="item.key2 && customerDetail.user_info_add">{{
            (customerDetail.user_info_add[item.key] || "-") +
            " - " +
            (customerDetail.user_info_add[item.key2] || "")
          }}</span>
          <span v-else-if="customerDetail.user_info_add">
            {{ customerDetail.user_info_add[item.key] || "--" }}
          </span>
        </div>
        <!-- 选择尺码&日期弹层 -->
        <van-popup v-model="showSelSizeBir" :closeable="true">
          <div class="size-bir-view">
            <div class="title">选择尺码&日期</div>
            <div class="picker-cont">
              <div class="picker-div size">
                <div class="label">尺码</div>
                <div class="picker-in">
                  <van-picker
                    :default-index="size_i"
                    title="尺码"
                    @change="onChangeSize"
                    :columns="size_columns"
                  />
                </div>
              </div>
              <div class="picker-div">
                <div class="label">月份</div>
                <van-picker
                  :default-index="month_i"
                  @change="onChangeMonth"
                  title="月份"
                  :columns="month_columns"
                />
              </div>
              <div class="picker-div">
                <div class="label">日期</div>
                <van-picker
                  :default-index="day_i"
                  @change="onChangeDay"
                  title="日期"
                  :columns="day_columns"
                />
              </div>
            </div>
            <div @click="confirmBirSize" class="confirm-btn">确认</div>
          </div>
        </van-popup>
        <!-- 选择地区 弹窗 -->
        <van-popup v-model="showAddressPicker" :closeable="true">
          <div class="size-bir-view">
            <div class="title">选择地区</div>
            <van-picker
              :default-index="address_i"
              :columns="address_columns"
              @change="onAddressChange"
            />
            <div @click="confirmAddress" class="confirm-btn">确认</div>
          </div>
        </van-popup>
      </div>
      <!-- 其他企业 -->
      <div class="part-container-basic" v-else>
        <div
          class="part-container-basic-info"
          v-for="(item, index) in basicinfoList"
          :key="index"
        >
          <img :src="item.iconUrl" :alt="item.title" />&nbsp;
          {{ customerDetail[item.key] || "-" }}
        </div>
      </div>
      <!-- 客户标签 -->
      <div class="part-container-row" style="padding-bottom: 0" v-if="!isKS()">
        <div class="part-row" style="margin-top: 15px">
          <div class="part-row-left">
            <span class="part-row-left-title">客户标签</span>
            <van-popover
              placement="right-start"
              trigger="click"
              v-model="tagShowPop"
              :actions="[
                { text: '客户标签仅提供推荐建议，导购无法编辑客户标签。' },
              ]"
            >
              <template #reference>
                <img
                  class="tip-icon"
                  :src="require('@/assets/customer-portrait/desc-ic-o.png')"
                /> </template
            ></van-popover>
          </div>
        </div>
        <div class="labels">
          <div
            v-for="(item, index) in customerDetail.tagsList"
            :key="index"
            class="label-item"
          >
            {{ item.name }}
            <van-popover
              placement="right-start"
              trigger="click"
              v-model="item.isShowPop"
            >
              <div class="tags-desc">
                {{ item.third_tag_description || "暂无说明" }}
              </div>
              <template #reference>
                <img
                  class="tip-icon"
                  style="margin-left: 5px; vertical-align: middle"
                  :src="require('@/assets/customer-portrait/desc-ic.png')"
                /> </template
            ></van-popover>
          </div>
        </div>
      </div>
    </div>
    <!-- 分隔 -->
    <div class="part-container-split"></div>
    <!-- 客户印象 -->
    <div class="part-container" v-if="!isKS()">
      <div class="part-container-row">
        <div class="part-row">
          <div class="part-row-left">
            <span class="part-row-left-title">客户印象</span>
            <van-popover
              placement="right-start"
              trigger="click"
              v-model="impressionShowPop"
            >
              <div class="impression-desc">
                <p>为顾客添加或移除印象标签</p>
                <p>1、点击客户印象标签可以为顾客添加或移除印象标签</p>
                <p>2、已添加的标签以深色底色展示</p>
                <p>
                  3、所有导购都可以为顾客添加或移除印象标签，页面展示最后一次编辑
                </p>
              </div>
              <template #reference>
                <img
                  class="tip-icon"
                  :src="
                    require('@/assets/customer-portrait/desc-ic-o.png')
                  " /></template
            ></van-popover>
          </div>
        </div>
        <div class="c-impression">
          <div
            class="c-impression-item"
            v-for="(item, index) in customerDetail.user_impression_tags"
            :key="index"
          >
            <div class="title">
              {{ item.group_name }}
              <img
                class="arrow-icon"
                @click="triggerTags(index)"
                :src="
                  require(`@/assets/customer-portrait/${
                    !item.isShowTags ? 'arrow-up' : 'arrow-down'
                  }.png`)
                "
              />
            </div>
            <div class="labels" v-if="item.isShowTags">
              <div
                v-for="(label, lindex) in item.impression_tags_list"
                :key="lindex"
                class="label-item"
                :class="{ 'no-item': !label.is_have }"
                @click="tagsImpressionOperate(label, item.impression_tags_list)"
              >
                {{ label.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 分隔 -->
    <div class="part-container-split"></div>

    <div class="part-container">
      <!-- 下边标签 -->
      <div class="part-container-row">
        <!-- 信息列表 -->
        <div
          class="service-info"
          v-for="(item, index) in serviceList"
          :key="index"
          @click="toOrderDetail(item)"
        >
          <span class="key">{{ item.title }}</span>
          <div class="value">
            {{ customerDetail[item.key] || "-" }}
            <img
              v-if="item.key === 'last_12m_purchase_amount'"
              :src="require('@/assets/customer-portrait/to-ic.png')"
            />
          </div>
        </div>
      </div>
    </div>
    <!-- 分隔 -->
    <div class="part-container-split"></div>
    <div class="part-container">
      <div class="part-container-row">
        <div class="part-skip">
          <div
            class="part-skip-button"
            v-for="(item, index) in skipList"
            :key="index"
            @click="toSkipDetail(item)"
          >
            <div class="skip-button-content">
              <img :src="item.iconUrl" :alt="item.title" />
              <p>
                <span class="title-text">{{ item.title }}</span>
                <van-icon name="arrow" />
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Image, Loading, Icon, Tag, Popup, Field, Picker } from "vant";
import {
  customerDetailApi,
  tagsImpressionOperate,
  updateRemark,
  saveInfoAdd,
  // getOrderApi,
  // getOrderDetail,
} from "@/api/customer-portrait/index";
@Component({
  components: {
    [Image.name]: Image,
    [Loading.name]: Loading,
    [Icon.name]: Icon,
    [Tag.name]: Tag,
    [Popup.name]: Popup,
    [Field.name]: Field,
    [Picker.name]: Picker,
  },
})
export default class extends Vue {
  private tagShowPop = false;
  private impressionShowPop = false;
  private customerDetail: any = {};
  private showEditRemark = true;
  private showSelSizeBir = false; //选择尺码生日弹层
  private tagsLoading = false; // 客户印象标签接口是否请求中
  private remark_key = "";
  private size_columns = [
    "340",
    "345",
    "350",
    "355",
    "360",
    "365",
    "370",
    "375",
    "380",
    "385",
    "390",
    "395",
    "400",
    "405",
    "410",
    "415",
    "420",
    "425",
    "430",
    "435",
    "440",
    "445",
  ]; // 尺码数据
  private day_columns = []; // 日期数据
  private month_columns = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]; // 月份数据
  private address_i = 0; // 地址picker 下标
  private size_i = 0; // 尺码picker 下标
  private month_i = 0; // 月份picker 下标
  private day_i = 0; // 日期picker 下标
  private showAddressPicker = false;
  // 地址数据
  private address_columns = [
    "北京市",
    "天津市",
    "河北省",
    "山西省",
    "内蒙古自治区",
    "辽宁省",
    "吉林省",
    "黑龙江省",
    "上海市",
    "江苏省",
    "浙江省",
    "安徽省",
    "福建省",
    "江西省",
    "山东省",
    "河南省",
    "湖北省",
    "湖南省",
    "广东省",
    "广西壮族自治区",
    "海南省",
    "重庆市",
    "四川省",
    "贵州省",
    "云南省",
    "西藏自治区",
    "陕西省",
    "甘肃省",
    "青海省",
    "宁夏回族自治区",
    "新疆维吾尔族自治区",
    "台湾省",
    "香港特别行政区",
    "澳门特别行政区",
  ];
  // 是否KS企业
  private isKS() {
    // 企业hash：KS-uat 522474ecff3b03fb41d8360eb259da58
    // 企业hash：KS-prod 0db7b1ee65bb6b547e687c9ae54f4f56
    let hashArrs = [
      "522474ecff3b03fb41d8360eb259da58",
      "0db7b1ee65bb6b547e687c9ae54f4f56",
    ];
    return hashArrs.some((ele) => {
      return location.href.indexOf(ele) > -1;
    });
    // return true;
  }
  // 是否SW企业
  private isSW() {
    // 企业hash：SW-uat 9c217b24b3fc5dcf685a7c566d191272
    // 企业hash：SW-prod 7cdd2af2e96d8617437846067440a3ec
    let hashArrs = [
      "9c217b24b3fc5dcf685a7c566d191272",
      "7cdd2af2e96d8617437846067440a3ec",
    ];
    return hashArrs.some((ele) => {
      return location.href.indexOf(ele) > -1;
    });
    // return true;
  }
  private basicinfoList = [
    {
      title: "会员卡号",
      key: "membership",
      iconUrl: require("@/assets/customer-portrait/icon-card.png"),
    },
    {
      title: "生日",
      key: "birthday_str",
      iconUrl: require("@/assets/customer-portrait/icon-birthday.png"),
    },
  ];
  private swBasicinfoList = [
    {
      title: "尺码",
      key: "shoe_size",
      iconUrl: require("@/assets/customer-portrait/icon-size.png"),
    },
    {
      title: "生日",
      key: "month",
      key2: "day",
      iconUrl: require("@/assets/customer-portrait/icon-birthday.png"),
    },
    {
      title: "地址",
      key: "address",
      iconUrl: require("@/assets/customer-portrait/icon-site.png"),
    },
  ];
  private serviceList = [
    {
      title: "会籍到期时间:",
      key: "expiry_date",
    },
    {
      title: "最近升级时间:",
      key: "upgrade_date_time",
    },
    {
      title: "可用积分:",
      key: "total_valid_points",
    },
    {
      title: "升级所需积分:",
      key: "upgrade_need_points",
    },
    {
      title: "过去12个月消费金额:",
      key: "last_12m_purchase_amount",
    },
  ];
  private skipList = [
    // {
    //   title: "积分明细",
    //   iconUrl: require("@/assets/customer-portrait/integral-icon.png"),
    // },

    {
      title: "购物订单",
      iconUrl: require("@/assets/customer-portrait/order-icon.png"),
    },
    {
      title: "卡券详情",
      iconUrl: require("@/assets/customer-portrait/coupon-icon.png"),
    },
    // {
    //   title: "客户足迹",
    //   iconUrl: require("@/assets/customer-portrait/footprint-icon.png"),
    // },
  ];
  public get resCustomerDetail(): any {
    return this.$store.getters["customer/getResCustomerDetail"];
  }
  /*********************** 生命周期 *************************/
  private created() {
    if (
      process.env.NODE_ENV == "test" ||
      process.env.NODE_ENV == "development"
    ) {
      this.$store.commit(
        "user/setExternalId",
        String("wm9hx_EQAAneXgyMOyEG6wYnMtbo82QA")
      );
    }
    this.getUserId();
  }
  private mounted() {
    if (this.isKS()) {
      this.serviceList.splice(4, 1, {
        title: "过去3年消费金额:",
        key: "last_36m_purchase_amount",
      });
      this.basicinfoList.push({
        title: "地址",
        key: "address",
        iconUrl: require("@/assets/customer-portrait/icon-site.png"),
      });
    }
    if (this.isSW()) {
      this.getDayColums();
      // SW企业隐藏卡券详情
      this.skipList.splice(1, 1);
    }
  }
  /*********************** methods *************************/

  public get ExternalId(): string {
    return this.$store.getters["user/getExternalId"];
  }

  async getUserId(): Promise<void> {
    if (this.resCustomerDetail) {
      this.initData(this.resCustomerDetail);
    }
    this.initPage(this.ExternalId);
  }
  private initData(data) {
    const tagsList = [];
    if (data.user_tags && data.user_tags.length) {
      data.user_tags.forEach((item: any) => {
        if (item.tags_list && item.tags_list.length) {
          item.tags_list.forEach((tag) => {
            tag.isShowPop = false;
            tagsList.push(tag);
          });
        }
      });
    }
    console.log("tagsList:", tagsList);
    if (data.user_impression_tags && data.user_impression_tags.length) {
      data.user_impression_tags.forEach((item) => {
        item.isShowTags = true;
      });
    }
    this.customerDetail = { ...data, tagsList: tagsList };
    this.$store.commit("user/setCustomerId", Number(data.user_id));
  }
  async initPage(userId: string): Promise<void> {
    const { code, data, message }: IApiResult = await customerDetailApi(userId);
    if (code == 200) {
      this.initData(data);
      this.$store.commit("customer/setResCustomerDetail", data);
      this.$toast.clear();
      // 客户印象操作接口请求标示
      this.tagsLoading = false;
    } else {
      this.$toast.fail(message);
    }
  }
  // 订单跳转详情
  private toOrderDetail(item: any) {
    if (item.key == "last_12m_purchase_amount") {
      let query = {
        title: "购物订单",
        union_id: this.customerDetail.union_id,
      };
      this.$set(query, "union_id", this.customerDetail.union_id);
      this.$router.push({
        path: "/portrait-customer-details",
        query: query,
      });
    }
  }
  // 卡劵/订单跳转详情
  private toSkipDetail(item: any) {
    let query = {
      title: item.title,
      union_id: this.customerDetail.union_id,
    };
    if (item.title === "积分明细") {
      const evalid_points = this.customerDetail.evalid_points;
      this.$set(query, "evalid_points", evalid_points);
    }
    if (item.title === "购物订单") {
      this.$set(query, "union_id", this.customerDetail.union_id);
    }
    if (item.title === "卡券详情") {
      this.$set(query, "union_id", this.customerDetail.union_id);
    }
    this.$router.push({
      path: "/portrait-customer-details",
      query: query,
    });
  }
  //操作印象标签关系
  private tagsImpressionOperate(tag: any, tagArr: any) {
    if (!this.customerDetail.customer_code) {
      this.$toast("顾客还未注册会员，无法使用标签功能");
      return;
    }
    if (this.tagsLoading) return;
    let add_tag_ids = [];
    let del_tag_ids = [];
    tagArr.forEach((item: any) => {
      if (item.is_have) {
        del_tag_ids.push(item.tag_id);
      }
    });
    if (!tag.is_have) {
      add_tag_ids = [tag.tag_id];
    }
    this.$toast.loading();
    this.tagsLoading = true;
    tagsImpressionOperate({
      enterprise_user_id: this.customerDetail.user_id,
      add_tag_ids: add_tag_ids,
      del_tag_ids: del_tag_ids,
    }).then((res: any) => {
      if (res.code == 200) {
        this.initPage(this.ExternalId);
      } else {
        this.$toast.fail(res.message);
        this.tagsLoading = false;
      }
    });
  }
  // 编辑
  private editCustomer() {
    if (this.customerDetail.name) return;
    this.showEditRemark = true;
    this.remark_key = this.customerDetail.remark || "";
  }
  // 打开SW选择日期尺码弹层
  private openSelPopup(index) {
    if (!this.customerDetail.customer_code) {
      this.$toast("顾客还未注册会员，无法修改");
      return;
    }
    this.address_i = this.address_columns.findIndex(
      (item) =>
        item ===
        (this.customerDetail.user_info_add
          ? this.customerDetail.user_info_add.address
          : "")
    );
    this.size_i = this.size_columns.findIndex(
      (item) =>
        item ===
        (this.customerDetail.user_info_add
          ? this.customerDetail.user_info_add.shoe_size
          : "")
    );
    this.month_i = this.month_columns.findIndex(
      (item) =>
        item ===
        (this.customerDetail.user_info_add
          ? this.customerDetail.user_info_add.month
          : "")
    );
    this.day_i = this.day_columns.findIndex(
      (item) =>
        item ===
        (this.customerDetail.user_info_add
          ? this.customerDetail.user_info_add.day
          : "")
    );
    if (index === 2) {
      this.address_i = this.address_i === -1 ? 0 : this.address_i;
      this.showAddressPicker = true;
    } else {
      this.showSelSizeBir = true;
      this.size_i = this.size_i === -1 ? 0 : this.size_i;
      this.month_i = this.month_i === -1 ? 0 : this.month_i;
      this.day_i = this.day_i === -1 ? 0 : this.day_i;
    }
  }
  // 选择日期尺码确定
  private confirmBirSize() {
    this.saveInfoAdd();
    this.showSelSizeBir = false;
  }
  // 尺码选择
  private onChangeSize(e: any, val: any, index: any) {
    this.size_i = index;
  }
  // 月份选择
  private onChangeMonth(e: any, val: any, index: any) {
    this.month_i = index;
    this.day_i = 0;
    this.getDayColums();
  }
  // 日期选择
  private onChangeDay(e: any, val: any, index: any) {
    this.day_i = index;
  }
  // 获取日期数据
  private getDayColums() {
    // const year = new Date().getFullYear();
    const num = new Date(2000, this.month_columns[this.month_i], 0).getDate();
    let columns = [];
    for (let i = 0; i < num; i++) {
      columns.push(i + 1);
    }
    this.day_columns = columns;
  }
  // 选择地址
  private onAddressChange(e: any, val: any, index: any) {
    this.address_i = index;
  }
  // 选择地址 确定
  private confirmAddress() {
    this.saveInfoAdd();
    this.showAddressPicker = false;
  }
  // 编辑客户地址、生日、鞋码
  private saveInfoAdd() {
    this.$toast.loading();
    saveInfoAdd({
      enterprise_user_id: this.customerDetail.user_id,
      address: this.address_columns[this.address_i] || "",
      shoe_size: this.size_columns[this.size_i] || "",
      month: this.month_columns[this.month_i] || 0,
      day: this.day_columns[this.day_i] || 0,
    }).then((res: any) => {
      if (res.code == 200) {
        this.initPage(this.ExternalId);
      } else {
        this.$toast.fail(res.message);
      }
    });
  }
  // 客户印象 标签收起/展开
  private triggerTags(index: any) {
    let target = this.customerDetail.user_impression_tags;
    target[index].isShowTags = !target[index].isShowTags;
    this.$set(this.customerDetail, "user_impression_tags", target);
  }
  // 客户未提供电话时提示
  private phoneTip() {
    this.$toast("顾客不接受电话沟通");
  }
  // 修改用户备注
  private updateRemark() {
    if (!this.remark_key) {
      this.$toast("请输入备注名");
      return;
    }
    updateRemark({
      external_userid: this.ExternalId,
      remark: this.remark_key,
    }).then((res: any) => {
      if (res.code == 200) {
        this.showEditRemark = false;
        this.initPage(this.ExternalId);
      } else {
        this.$toast.fail(res.message);
      }
    });
  }

  private makePhoneCall() {
    const { mobile } = this.customerDetail
    const self = this
    window.wx.invoke('makePhoneCall',{ phoneNumber: mobile }, function () {
      // 神策数据上报 - 拨打顾客电话
      self.$track.trackCallCustomer({
        user_external_id: self.ExternalId
      })
    } )
  }
}
</script>
<style lang="scss" scoped>
.add-btn {
  width: 58px;
  height: 26px;
  border-radius: 4px;
  padding: 0;
  border: solid 1px #b9bec4;
  font-size: 14px;
  color: #272727;
  line-height: 24px;
}

.container {
  .pointer {
    cursor: pointer;
  }
  ::v-deep .van-tab.van-tab--active {
    font-weight: 600;
  }
  font-size: 14px;
  background: #ffffff;
  padding-bottom: 20px;
  .part-container {
    padding: 20px 15px 0px;
    .part-container-row:last-child {
      border: 0;
    }
    .part-container-row {
      padding-bottom: 15px;
      .tip-icon {
        width: 20px;
        height: 20px;
        margin-left: 7px;
      }
      .person-row {
        display: flex;
        .image {
          flex: 0 0 60px;
          height: 60px;
          width: 60px;
        }
        .info-content {
          flex: 1;
          display: flex;
          font-size: 16px;
          padding-left: 33px;
          position: relative;
          .left {
            width: 237px;
            position: relative;
            p {
              line-height: 25px;
              color: #272727;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
              font-weight: 500;
              font-size: 18px;
              margin-top: 5px;
            }
            .sex {
              position: absolute;
              left: -23px;
              top: 8px;
              height: 18px;
              width: 18px;
            }
            .c-name {
              font-size: 14px;
              color: #96979e;
              padding-top: 8px;
              font-weight: 400;
              line-height: 17px;
              display: flex;
              span {
                display: block;
                white-space: nowrap;
                text-overflow: ellipsis;
                overflow: hidden;
                max-width: 240px;
                font-size: 12px;
                margin-right: 4px;
              }
              img {
                width: 16px;
                height: 16px;
              }
            }
            .grade {
              font-size: 14px;
              margin-top: 8px;
              font-weight: 500;
              color: #daba82;
              line-height: 20px;
              display: flex;
              justify-content: space-between;
              span {
                display: block;
                width: 50%;
              }
            }
          }
          .c-mobile {
            position: absolute;
            right: 0;
            top: 5px;
            a {
              img {
                display: block;
                width: 24px;
                height: 24px;
                border-radius: 100%;
              }
            }
          }
        }
      }
      .service-info {
        font-size: 14px;
        display: flex;
        margin-bottom: 15px;
        justify-content: space-between;
        align-items: center;
        color: #33363c;
        .key {
          flex: 0 0 141px;
          color: #96979e;
          line-height: 20px;
          font-weight: 400;
        }
        .value {
          img {
            width: 14px;
            height: 14px;
            margin-left: 5px;
            transform: translateY(2px);
          }
        }
      }
      .part-row {
        display: flex;
        justify-content: space-between;
        padding-right: 15px;
        .part-row-right {
          .img {
            touch-action: none;
            width: 6px;
            height: 11px;
          }
          span {
            vertical-align: top;
            padding-right: 10px;
            color: #96979e;
          }
        }
        .part-row-left {
          .integral {
            color: #4a4a4a !important;
            font-weight: 400 !important;
          }
          .part-row-left-title {
            font-size: 14px;
            font-weight: 500;
            color: #33363c;
            vertical-align: super;
          }
          .van-icon-warning-o {
            margin-left: 5px;
            transform: translateY(-4px);
          }
          .part-row-left-image {
            height: 20px;
            width: 20px;
            margin-left: 5px;
            position: relative;
            top: -1px;
          }
        }
      }
      .labels {
        display: flex;
        flex-wrap: wrap;
        .label-item {
          background: #4a4a4a;
          color: #ffffff;
          font-size: 14px;
          height: 30px;
          line-height: 30px;
          padding-left: 10px;
          padding-right: 5px;
          margin-right: 10px;
          margin-top: 10px;
          max-width: 344px;
          min-width: 48px;
          text-align: center;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }

      // 跳转按钮
      .part-skip {
        display: flex;
        flex-wrap: wrap;
        .part-skip-button {
          width: calc(50% - 5px);
          height: 117px;
          background: #f6f7f9;
          display: flex;
          justify-content: center;
          align-items: center;
          border-radius: 2px;
          .skip-button-content {
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            font-size: 14px;
            color: #33363c;
            font-weight: 400;
            line-height: 20px;
            .title-text {
              margin-right: 5px;
            }
            img {
              width: 44px;
              height: 44px;
              margin-bottom: 5px;
            }
          }
        }
        .part-skip-button:nth-child(odd) {
          margin-right: 10px;
        }
      }
      .c-impression {
        width: 345px;
        padding: 15px;
        margin-top: 10px;
        background: #f3f4f6;
        .c-impression-item {
          padding-bottom: 14px;
          margin-bottom: 14px;
          border-bottom: 1px dashed #cccccc;
          .title {
            font-size: 16px;
            font-weight: 400;
            color: #33363c;
            padding-bottom: 5px;
            position: relative;
            padding-right: 20px;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            .arrow-icon {
              position: absolute;
              right: 0;
              top: 3px;
              width: 14px;
              height: 15px;
            }
          }
          .label-item {
            padding: 0 10px;
          }
          .no-item {
            background: transparent;
            color: #4a4a4a;
            border: 1px solid #4a4a4a;
          }
          &:last-child {
            margin-bottom: 0;
            border-bottom: 0 none;
          }
        }
      }
    }
    // 基本信息
    .part-container-basic {
      display: flex;
      flex-wrap: wrap;
      width: 345px;
      // height: 46px;
      line-height: 46px;
      background: #f3f4f6;
      .part-container-basic-info {
        padding-left: 10px;
        width: 40%;

        &:first-child {
          width: 60%;
        }
        &:nth-child(3) {
          width: 100%;
        }
        &.sw-info-li {
          position: relative;
          .arrow_right {
            position: absolute;
            right: 15px;
            top: 50%;
            transform: translateY(-50%);
            height: 12px;
            width: 12px;
          }
          width: 50%;
          &:nth-child(3) {
            width: 100%;
          }
        }
        img {
          width: 14px;
          height: 14px;
        }
        .value {
          line-height: 20px;
          white-space: normal;
          word-break: break-all;
          color: #33363c;
        }
        .overflow {
          height: 20px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
  .part-container-split {
    height: 10px;
    background: #f7f7f7;
  }
  .part-container-content {
    background: #ffffff;
  }
}
.tags-desc {
  max-width: 200px;
  padding: 10px;
  font-size: 12px;
  span {
    margin-right: 10px;
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
.impression-desc {
  width: 240px;
  font-size: 12px;
  padding: 15px;
}
.size-bir-view {
  width: 320px;
  .title {
    text-align: center;
    font-size: 16px;
    line-height: 50px;
  }
  .picker-cont {
    display: flex;
    .picker-div {
      flex: 3;
      text-align: center;
      .picker-in {
        border-right: 1px solid #f0f0f0;
      }
    }
  }
  .confirm-btn {
    width: 70%;
    background-color: #4a4a4a;
    color: #ffffff;
    line-height: 45px;
    text-align: center;
    margin: 20px auto;
  }
}
</style>
