<template>
  <div class="customer-detail">
    <van-nav-bar title="客户详情" left-arrow @click-left="goBack" />
    <div class="customer-top">
      <img
        class="customer-img"
        :src="info.avatar || require('@/assets/default/default-person.png')"
        alt=""
      />
      <div class="customer-name-box">
        <div class="customer-name">{{ info.nick_name }}</div>
        <img
          class="customer-sex-icon"
          :src="customSexList[info.gender || 0]"
          alt=""
        />
      </div>
    </div>
    <div class="border-line"></div>
    <div class="customer-middle">
      <van-form>
        <van-field
          v-model="info.age"
          name="年龄"
          maxlength="3"
          label="年龄"
          type="digit"
          placeholder="请输入年龄"
          clearable
        />
        <van-field
          v-model="info.birthday"
          name="生日"
          label="生日"
          placeholder="请选择生日"
          clearable
          @focus="closeMobileKeypad"
          @click="openPop"
        />
        <van-field name="性别" label="性别" placeholder="请选择">
          <template #input>
            <van-radio-group v-model="info.gender" direction="horizontal">
              <van-radio :name="1">
                男
                <template #icon="props">
                  <img
                    class="img-icon"
                    :src="
                      props.checked
                        ? require('@/assets/customer-detail/active-icon.png')
                        : require('@/assets/customer-detail/inactive-icon.png')
                    "
                  />
                </template>
              </van-radio>
              <van-radio :name="2">
                女
                <template #icon="props">
                  <img
                    class="img-icon"
                    :src="
                      props.checked
                        ? require('@/assets/customer-detail/active-icon.png')
                        : require('@/assets/customer-detail/inactive-icon.png')
                    "
                  />
                </template>
              </van-radio>
              <van-radio :name="0">
                未知
                <template #icon="props">
                  <img
                    class="img-icon"
                    :src="
                      props.checked
                        ? require('@/assets/customer-detail/active-icon.png')
                        : require('@/assets/customer-detail/inactive-icon.png')
                    "
                  />
                </template>
              </van-radio>
            </van-radio-group>
          </template>
        </van-field>
        <!--  -->
        <van-field
          v-model="info.mobile"
          name="电话"
          label="电话"
          type="tel"
          maxlength="11"
          clearable
          placeholder="请输入电话"
        />
        <van-field
          v-model="info.email"
          name="邮件"
          label="邮件"
          maxlength="128"
          placeholder="请输入邮件"
          clearable
        />
        <van-field
          v-model="info.qq"
          name="QQ"
          label="QQ"
          type="number"
          maxlength="20"
          placeholder="请输入QQ"
          clearable
        />
        <van-field
          v-model="info.describe"
          name="描述"
          label="描述"
          maxlength="128"
          placeholder="请输入描述"
          clearable
        />
      </van-form>
    </div>
    <div class="consignee-list-box">
      <div class="customer-consignee">收货人信息</div>
      <div
        class="item-consignee"
        v-for="(item, index) in addressList"
        :key="index"
      >
        <div class="consignee-name-box">
          <div class="consignee-name">{{ item.consignee_name }}</div>
          <div class="consignee-phone">{{ item.consignee_mobile }}</div>
        </div>
        <div class="consignee-address">
          {{ item.province }}{{ item.city }}{{ item.area }}{{ item.address }}
        </div>
        <div class="copy-box">
          <div
            class="copy-btn"
            @click="
              copyOrderNo(
                item.consignee_name +
                  item.consignee_mobile +
                  item.province +
                  item.city +
                  item.area +
                  item.address
              )
            "
          >
            复制
          </div>
        </div>
      </div>
    </div>
    <div class="border-line"></div>
    <div class="save-box">
      <div class="save-btn" @click="saveCustomer">保存</div>
    </div>
    <!-- 日历选择器-->
    <van-popup
      v-model="showCalendar"
      position="bottom"
      class="datetime-box"
      :style="{ height: '30%' }"
      ><van-datetime-picker
        v-model="currentDate"
        type="date"
        :min-date="minDate"
        :max-date="maxDate"
        @confirm="onConfirm"
        @cancel="showCalendar = false"
    /></van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import {
  customerDetailApi,
  customerAddrApi,
  customerRemarkSaveApi,
} from "@/api/customer-portrait";
import {
  Form,
  Field,
  RadioGroup,
  Radio,
  DatetimePicker,
  NavBar,
  Popup,
  NumberKeyboard,
} from "vant";
@Component({
  components: {
    [NavBar.name]: NavBar,
    [Form.name]: Form,
    [Field.name]: Field,
    [RadioGroup.name]: RadioGroup,
    [Radio.name]: Radio,
    [DatetimePicker.name]: DatetimePicker,
    [Popup.name]: Popup,
    [NumberKeyboard.name]: NumberKeyboard,
  },
})
export default class extends Vue {
  /*************************** data start ***************************/
  private customSexList = {
    "0": require("@/assets/customer-portrait/unknown-sex.png"),
    "1": require("@/assets/customer-portrait/man-sex.png"),
    "2": require("@/assets/customer-portrait/woman-sex.png"),
  };
  private currentDate = new Date();
  private minDate = new Date(1950, 1, 1);
  private maxDate = new Date();
  private info = {
    user_id: 0,
    avatar: "",
    nick_name: "",
    age: "",
    gender: 0,
    birthday: "",
    mobile: "",
    email: "",
    qq: "",
    describe: "",
  };
  private addressList = [];
  private showCalendar = false;
  public get resCustomerDetail(): any {
    return this.$store.getters["customer/getResCustomerDetail"];
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    if (this.resCustomerDetail) {
      this.initData(this.resCustomerDetail);
    } else {
      this.customerDetail();
    }

    // this.getcustomAddr();
  }
  /*************************** method start ***************************/
  private goBack() {
    this.$router.go(-1);
  }
  private copyOrderNo(msg: string) {
    this.$helper.copyMessage(msg);
  }
  private async saveCustomer() {
    const param: ICustomerRemarkSaveParams = {
      user_id: this.info.user_id, // 用户ID
      age: Number(this.info.age),
      birthday: this.info.birthday
        ? (new Date(this.info.birthday).getTime() / 1000) | 0
        : 0,
      gender: this.info.gender,
      mobile: this.info.mobile,
      email: this.info.email,
      qq: this.info.qq,
      describe: this.info.describe,
    };
    const { code, message = "" }: IApiResult = await customerRemarkSaveApi(
      param
    );
    if (code === 200) {
      this.$store.commit(
        "customer/setResCustomerDetail",
        Object.assign(this.resCustomerDetail, param)
      );
      this.$router.go(-1);
    } else {
      this.$toast.fail(message);
    }
  }
  //获取客户地址
  private async getcustomAddr() {
    const {
      code,
      message = "",
      data,
    }: IApiResult = await customerAddrApi({
      user_id:
        Number(this.$route.query.customerId) ||
        this.$store.getters["user/getCustomerId"],
      page: 1,
      size: 100,
    });
    if (code === 200) {
      this.addressList = data.list || [];
    } else {
      this.$toast.fail(message);
    }
  }
  private openPop() {
    this.showCalendar = true;
  }
  private initData(data) {
    if (!data.age) {
      data.age = "";
    }
    if (!data.birthday) {
      data.birthday = "";
    }
    this.info = data;
    if (data.birthday) {
      this.info.birthday = this.$helper.dateFormat(data.birthday, "yyyy-MM-dd");
      this.currentDate = new Date(data.birthday * 1000);
    }
  }
  // 客户详情接口
  private async customerDetail() {
    const {
      code,
      message = "",
      data,
    }: IApiResult = await customerDetailApi(
      this.$store.getters["user/getExternalId"]
    );
    if (code === 200) {
      this.initData(data);
      this.$store.commit("customer/setResCustomerDetail", data);
    } else {
      this.$toast.fail(message);
    }
  }
  private formatDate(date) {
    return `${date.getMonth() + 1}/${date.getDate()}`;
  }
  private onConfirm(date) {
    this.showCalendar = false;
    this.info.birthday = this.$helper.dateFormat(
      date.getTime() / 1000,
      "yyyy-MM-dd"
    );
  }
  private closeMobileKeypad() {
    try {
      (document.activeElement as any).blur();
    } catch (e) {
      console.log(e);
    }
  }
}
</script>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}
.datetime-box {
  height: auto !important;
}
::v-deep {
  .van-nav-bar .van-icon {
    color: #999;
  }
  .van-field {
    margin-bottom: 20px;
  }
  .van-field__label {
    color: #808080;
    font-size: 14px;
    margin-right: 10px;
    width: 28px;
  }
  .van-field__control {
    color: #272727;
    font-size: 14px;
  }
  .van-cell {
    padding: 0;
  }
  .van-cell::after {
    border-bottom: none;
  }
  .van-radio__label {
    font-size: 14px;
    color: #272727;
    margin-left: 3px;
  }
  .van-radio {
    margin-right: 20px;
  }
}

.m-t-10 {
  margin-top: 10px;
}
.img-icon {
  width: 12px;
  height: 12px;
}
.customer-detail {
  background-color: #ffffff;
  overflow-x: hidden;
  .customer-top {
    height: 80px;
    width: 100%;
    display: flex;
    align-items: center;
    padding: 20px 15px;
    background-color: #ffffff;
    .customer-img {
      width: 40px;
      height: 40px;
      border-radius: 4px;
    }
    .customer-name-box {
      margin-left: 10px;
      .customer-name {
        font-size: 16px;
        color: #272727;
        width: 295px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-weight: 600;
      }
      .customer-sex-icon {
        width: 15px;
        height: 15px;
        margin-top: 5px;
        display: block;
      }
    }
  }
  .border-line {
    width: calc(100vw - 15px);
    height: 1px;
    background: #ebebeb;
    position: relative;
    left: 15px;
  }
  .customer-middle {
    padding: 19px 15px 5px 15px;
    background-color: #ffffff;
  }
  .consignee-list-box {
    background: #ffffff;
    padding: 0 15px 20px 15px;
    .customer-consignee {
      color: #808080;
      margin-right: 10px;
      font-size: 14px;
      line-height: 20px;
    }
    .item-consignee {
      width: 100%;
      min-height: 111px;
      background-color: #f4f5f8;
      border-radius: 2px;
      margin-top: 10px;
      padding: 10px;
      position: relative;
      .consignee-name-box {
        font-size: 14px;
        display: flex;
        align-items: center;
        .consignee-name {
          padding-right: 10px;
          color: #272727;
          line-height: 20px;
          max-width: 60%;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .consignee-phone {
          color: #272727;
          line-height: 20px;
          max-width: 40%;
        }
      }
      .consignee-address {
        font-size: 14px;
        line-height: 20px;
        color: #808080;
        margin-top: 4px;
        border-bottom: solid 1px #ebebeb;
        padding-bottom: 10px;
        word-break: break-all;
      }
      .copy-box {
        height: 36px;
        position: relative;
        .copy-btn {
          cursor: pointer;
          position: absolute;
          top: 10px;
          right: 0;
          width: 75px;
          height: 26px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #272727;
          font-size: 14px;
          border-radius: 4px;
          border: solid 1px #b9bec4;
        }
      }
    }
  }
  .save-box {
    height: 72px;
    position: relative;
    background: #ffffff;
    .save-btn {
      cursor: pointer;
      position: absolute;
      top: 20px;
      right: 15px;
      width: 100px;
      height: 28px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #000000;
      border: 1px solid #000000;
      font-size: 14px;
      border-radius: 4px;
    }
  }
}
</style>
