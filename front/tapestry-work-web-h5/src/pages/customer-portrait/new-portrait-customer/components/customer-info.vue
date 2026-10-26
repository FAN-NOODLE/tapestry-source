<template>
  <div class="box">
    <!-- 尺寸、生日 -->
    <div class="column-2" @click="openSelPopup(1)">
      <div class="box-image" v-if="isSW">
        <img class="icon m-r" :src="renderImg('2023/v3.4/icon_2.png')" />
        <div class="label">{{ renderUserInfo("shoe_size") }}</div>
      </div>
      <div class="box-image" v-if="isKS">
        <img
          class="icon m-r"
          :src="renderImg('2022/v3.0.2/weclient/card-number-icon.png')"
        />
        <div class="label">{{ info.customer_code || "--" }}</div>
        <img
          class="icon m-l"
          :src="renderImg('2023/v3.2.0.1/copy.png')"
          @click="copyMemberCard"
        />
      </div>
      <div class="box-image">
        <img class="icon m-r" :src="renderImg('2023/v3.4/icon_1.png')" />
        <div class="label" v-if="isSW">
          {{ renderUserInfo("month") }}-{{ renderUserInfo("day") }}
        </div>
        <div class="label" v-if="isKS">
          {{ info.birthday_str }}
        </div>
        <img
          v-if="isSW"
          class="arrow-right"
          :src="renderImg('2022/v3.0.2/weclient/arrow-right-simple.png')"
        />
      </div>
    </div>

    <!-- 地址 -->
    <div class="box-image address" @click="openSelPopup(2)">
      <img class="icon m-r" :src="renderImg('2023/v3.4/address.png')" />
      <div class="label" v-if="isSW">{{ renderUserInfo("address") }}</div>
      <div class="label" v-if="isKS">{{ info.address }}</div>
      <img
        v-if="isSW"
        class="arrow-right"
        :src="renderImg('2022/v3.0.2/weclient/arrow-right-simple.png')"
      />
    </div>

    <div class="divider"></div>

    <!-- 订单信息 -->
    <div :class="[isFold ? 'more-height' : '']">
      <div
        class="column-second-2"
        v-for="(item, index) in orderList"
        :key="index"
      >
        <div class="second-label">{{ item.label }}</div>
        <div class="label label-r">
          {{ item.value ? `${item.unit} ${item.value}` : "--" }}
        </div>
      </div>
    </div>

    <!-- 收起、展开 -->
    <div class="box-more" @click="setEnable">
      <div class="more-text">{{ !isFold ? "收起" : "展开" }}</div>
      <div class="triangle-icon" :class="!isFold ? 'rotate-icon' : ''"></div>
    </div>

    <!-- 选择尺码&日期弹层 -->
    <van-popup v-model="sizeOptions.showSelSizeBir" :closeable="true">
      <div class="size-bir-view">
        <div class="title">选择尺码&日期</div>
        <div class="picker-cont">
          <div class="picker-div size">
            <div class="label">尺码</div>
            <div class="picker-in">
              <van-picker
                :default-index="sizeOptions.size_i"
                title="尺码"
                @change="sizeOptions.onChangeSize"
                :columns="sizeOptions.size_columns"
              />
            </div>
          </div>
          <div class="picker-div">
            <div class="label">月份</div>
            <van-picker
              :default-index="sizeOptions.month_i"
              @change="sizeOptions.onChangeMonth"
              title="月份"
              :columns="sizeOptions.month_columns"
            />
          </div>
          <div class="picker-div">
            <div class="label">日期</div>
            <van-picker
              :default-index="sizeOptions.day_i"
              @change="sizeOptions.onChangeDay"
              title="日期"
              :columns="sizeOptions.day_columns"
            />
          </div>
        </div>
        <div @click="confirmBirSize" class="confirm-btn">确认</div>
      </div>
    </van-popup>

    <!-- 选择地区 弹窗 -->
    <van-popup v-model="addressOptions.showAddressPicker" :closeable="true">
      <div class="size-bir-view">
        <div class="title">选择地区</div>
        <van-picker
          :default-index="addressOptions.address_i"
          :columns="addressOptions.address_columns"
          @change="addressOptions.onAddressChange"
        />
        <div @click="confirmAddress" class="confirm-btn">确认</div>
      </div>
    </van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import { saveInfoAdd } from "@/api/customer-portrait";
import { Popup, Picker } from "vant";
import { InfoType } from "../index.vue";
import orderMapList from "../data/order-list";
import Helper from "@/utils/helper";
import wx from "weixin-js-sdk";
/** 地区 */
const ADDRESS = [
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

/** 尺码 */
const SIZE = [
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
];

const ORDER_LIST = Helper.isSW() ? orderMapList.sw : orderMapList.ks;

/** 生成订单信息 */
function generateOrder(obj = {}) {
  for (let key of ORDER_LIST.keys()) {
    if (obj[key]) {
      const result = ORDER_LIST.get(key);
      result.value = obj[key];
      ORDER_LIST.set(key, result);
    }
  }
  return [...ORDER_LIST.values()];
}

@Component({
  components: {
    [Popup.name]: Popup,
    [Picker.name]: Picker,
  },
})
export default class CustomerInfo extends Vue {
  /** 顾客详情 */
  @Prop() info!: InfoType;

  /** 外部联系人id */
  @Prop() openId!: string;

  /** 订单列表 */
  private orderList = [];

  /** 展开收起状态 */
  private isFold = true;
  // 是否是sw
  private isSW = Helper.isSW();
  // 是否是ks
  private isKS = Helper.isKS();
  /** 地区属性 */
  private addressOptions = {
    showAddressPicker: false,
    address_i: 0,
    address_columns: ADDRESS,
    onAddressChange: (e: HTMLFormElement, val: string, index: number) => {
      this.addressOptions.address_i = index;
    },
  };

  /** 尺码属性 */
  private sizeOptions = {
    showSelSizeBir: false,
    size_i: 0,
    month_i: 0,
    day_i: 0,
    size_columns: SIZE,
    month_columns: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12],
    day_columns: [],
    onChangeSize: (e: HTMLFormElement, val: string, index: number) => {
      this.sizeOptions.size_i = index;
    },
    onChangeMonth: (e: HTMLFormElement, val: string, index: number) => {
      this.sizeOptions.month_i = index;
      this.sizeOptions.day_i = 0;
      this.getDayColumns();
    },
    onChangeDay: (e: HTMLFormElement, val: string, index: number) => {
      this.sizeOptions.day_i = index;
    },
  };

  @Watch("info")
  changeInfo(value) {
    this.getDayColumns();
    this.orderList = generateOrder(value || {});
  }

  private created() {
    this.getDayColumns();
    this.orderList = generateOrder(this.info || {});
  }

  /** 获取基本信息, 非空处理 */
  private renderUserInfo(attr) {
    if (!this.info.user_info_add) return "--";

    return this.info.user_info_add[attr] || "--";
  }

  /** 打开选择日期尺码弹层 */
  private openSelPopup(index) {
    if (Helper.isKS()) {
      return;
    }
    if (!this.info.customer_code) {
      this.$toast("顾客还未注册会员，无法修改");
      return;
    }
    this.addressOptions.address_i =
      this.addressOptions.address_columns.findIndex(
        (item) =>
          item ===
          (this.info.user_info_add ? this.info.user_info_add.address : "")
      );
    this.sizeOptions.size_i = this.sizeOptions.size_columns.findIndex(
      (item) =>
        item ===
        (this.info.user_info_add ? this.info.user_info_add.shoe_size : "")
    );
    this.sizeOptions.month_i = this.sizeOptions.month_columns.findIndex(
      (item) =>
        item === (this.info.user_info_add ? this.info.user_info_add.month : "")
    );
    this.sizeOptions.day_i = this.sizeOptions.day_columns.findIndex(
      (item) =>
        item === (this.info.user_info_add ? this.info.user_info_add.day : "")
    );
    if (index === 2) {
      this.addressOptions.address_i =
        this.addressOptions.address_i === -1
          ? 0
          : this.addressOptions.address_i;
      this.addressOptions.showAddressPicker = true;
    } else {
      this.sizeOptions.showSelSizeBir = true;
      this.sizeOptions.size_i =
        this.sizeOptions.size_i === -1 ? 0 : this.sizeOptions.size_i;
      this.sizeOptions.month_i =
        this.sizeOptions.month_i === -1 ? 0 : this.sizeOptions.month_i;
      this.sizeOptions.day_i =
        this.sizeOptions.day_i === -1 ? 0 : this.sizeOptions.day_i;
    }
  }

  /** 获取日期数据 */
  private getDayColumns() {
    // const year = new Date().getFullYear();
    const num = new Date(
      2000,
      this.sizeOptions.month_columns[this.sizeOptions.month_i],
      0
    ).getDate();
    let columns = [];
    for (let i = 0; i < num; i++) {
      columns.push(i + 1);
    }
    this.sizeOptions.day_columns = columns;
  }

  /** 尺寸点击确认 */
  private confirmBirSize() {
    this.saveInfoAdd();
    this.sizeOptions.showSelSizeBir = false;
  }

  /** 地址点击确认  */
  private confirmAddress() {
    this.saveInfoAdd();
    this.addressOptions.showAddressPicker = false;
  }

  /** 编辑客户地址、生日、鞋码 */
  private saveInfoAdd() {
    this.$toast.loading();

    const { address_columns, address_i } = this.addressOptions;
    const { size_columns, size_i, month_columns, day_columns, month_i, day_i } =
      this.sizeOptions;

    saveInfoAdd({
      external_userid: this.openId,
      enterprise_user_id: this.info.user_id,
      address: address_columns[address_i] || "",
      shoe_size: size_columns[size_i] || "",
      month: month_columns[month_i] || 0,
      day: day_columns[day_i] || 0,
    }).then((res: IApiResult) => {
      if (res.code == 200) {
        this.$emit("renovate");
      } else {
        this.$toast.fail(res.message);
      }
    });
  }

  /** 渲染图片 */
  private renderImg(v) {
    const baseUrl =
      "https://cdn-cos.coach.com.cn/uploads/";
    return baseUrl + v;
  }

  /** 折叠、展开 */
  private setEnable() {
    this.isFold = !this.isFold;
  }
  // 复制会员卡号
  copyMemberCard() {
    const { customer_code } = this.info;
    window.wx.invoke("setClipboardData", {
      data: customer_code,
    });
  }
}
</script>
<style lang="scss" scoped>
.box {
  padding: 15px 10px;
  background: #f3f4f6;

  &-image {
    display: flex;
    align-items: center;
  }

  .arrow-right {
    margin-left: auto;
    width: 9px;
    height: 14px;
  }

  .icon {
    width: 14px;
    height: 14px;
  }

  .m-l {
    margin-left: 10px;
  }

  .m-r {
    margin-right: 10px;
  }

  .label {
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #000000;
    line-height: 16px;
  }
  .second-label {
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #96979e;
    line-height: 16px;
  }
  .address {
    margin-top: 15px;
  }

  .divider {
    height: 1px;
    border: 1px dashed #e5e7eb;
    margin: 15px 0;
  }
  .column-2 {
    display: grid;
    grid-template-columns: 60% 40%;
    overflow: hidden;
  }

  .column-second-2 {
    display: grid;
    grid-template-columns: 50% 50%;
    overflow: hidden;
    margin-top: 16px;
    .label-r {
      text-align: end;
    }
  }

  .column-second-2:first-child {
    margin-top: 0;
  }

  &-more {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 20px;
    .more-text {
      font-size: 14px;
      font-weight: 400;
      color: #000000;
      height: 17px;
      line-height: 17px;
      margin-right: 5px;
    }

    .triangle-icon {
      margin-top: 7px;
      display: inline-block;
      height: 0;
      border-top: 4px solid #a6a6a6;
      border-bottom: 4px solid transparent;
      border-left: 4px solid transparent;
      border-right: 4px solid transparent;
    }
    .rotate-icon {
      transform: rotate(180deg);
      margin-top: -4px;
    }
  }
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
    font-size: 14px;
  }
}

.more-height {
  max-height: 55px;
  overflow: hidden;
}
</style>
