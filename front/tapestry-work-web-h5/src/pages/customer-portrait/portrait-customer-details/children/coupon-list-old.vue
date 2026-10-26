<template>
  <!-- <empty-page
    picture="noData"
    paddingTop="60"
    v-if="steps.length == 0"
    paddingBottom="60"
    text="暂无卡券"
  ></empty-page>
  <div class="container" v-else> -->
  <div class="container">
    <van-tabs
      v-model="activeTab"
      title-active-color="#0095F6"
      @click="onTabClick"
      border
      ref="tabs"
    >
      <van-tab
        v-for="(item, index) in tabList"
        :title="item.title"
        :key="index"
      >
        <!-- 状态选择栏 -->
        <div class="top-state">
          <div
            :class="[
              'allStateType',
              stateDefault == item.id ? 'state-active' : '',
            ]"
            v-for="item in stateList"
            :key="item.id"
            @click="toSelectState(item.id)"
          >
            {{ item.title }}
          </div>
        </div>
        <div class="all-coupon-list">
          <div
            class="coupon-content"
            :style="{
              backgroundImage:
                'url(' +
                (2 === 1 ? notImgUrl : 2 === 3 ? useImgUrl : pastImgUrl) +
                ')',
              backgroundSize: '100% 100%',
            }"
          >
            <div
              class="coupon-content-left"
              :style="{ color: 1 === 1 ? '#a8a8a8' : '#909090' }"
            >
              <p :style="{ color: 1 === 1 ? '#a8a8a8' : '#909090' }">
                <span>￥</span>500
              </p>
              <p>满2000可用</p>
            </div>
            <div
              class="coupon-content-right"
              :style="{ color: 1 === 2 ? '#a8a8a8' : '#909090' }"
            >
              <p :style="{ color: 1 === 2 ? '#a8a8a8' : '#909090' }">券名称</p>
              <p>仅官网可用</p>
              <p>有效期至：2020.12.31</p>
            </div>
          </div>
        </div>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { couponListApi } from "@/api/customer-portrait";
import emptyPage from "@/components/base/empty-page.vue";
@Component({
  components: {
    emptyPage,
  },
})
export default class extends Vue {
  private activeTab = 1;
  private couponList = [];
  private tabList = [
    {
      title: "礼遇卡券",
      name: 1,
    },
    {
      title: "美容服务券",
      name: 2,
    },
    {
      title: "这是其他券",
      name: 3,
    },
    {
      title: "这是其他券",
      name: 4,
    },
    {
      title: "这是其他券",
      name: 5,
    },
    {
      title: "这是其他券",
      name: 6,
    },
  ];
  private stateList = [
    {
      title: "全部",
      id: "",
    },
    {
      title: "可使用",
      id: "0",
    },
    {
      title: "已使用",
      id: "3",
    },
    {
      title: "已过期",
      id: "4",
    },
  ];
  private stateDefault = "";
  private notImgUrl = require("@/assets/customer-portrait/not-coupon.png");
  private useImgUrl = require("@/assets/customer-portrait/use-coupon.png");
  private pastImgUrl = require("@/assets/customer-portrait/past-coupon.png");
  private mounted() {
    // if (this.$store.getters["user/getCustomerId"]) {
    //   this.getfootList();
    // }
  }
  private toSelectState(id: string) {
    this.stateDefault = id;
  }
  private onTabClick(name: number, title: string) {
    this.stateDefault = "";
    console.log(title);
  }
  private formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy-MM-dd hh:mm:ss");
  }
  async getfootList(): Promise<void> {
    const { code, data, message }: IApiResult = await couponListApi({
      consumerID: "9999999",
      transType: this.stateDefault,
    });
    if (code == 200) {
      console.log(data);
      // data.list.length > 0 && (this.steps = data.list);
    } else {
      this.$toast.fail(message);
    }
  }
}
</script>
<style lang="scss" scoped>
div.container {
  position: fixed;
  height: 100%;
  width: 100%;
  background-color: #ffffff;
  .top-state {
    display: flex;
    font-size: 14px;
    font-weight: 400;
    padding: 20px 15px;
    .allStateType {
      padding: 5px 15px;
      box-sizing: border-box;
      border-radius: 2px;
      margin-right: 10px;
      color: #808080;
      background: #f6f7f9;
    }
    .state-active {
      color: #5b89fe;
      background: #eef3fe;
    }
  }
  .all-coupon-list {
    padding: 5px 15px;
    .coupon-content {
      width: 100%;
      height: 90px;
      display: flex;
      .coupon-content-left {
        width: 28% !important;
        text-align: center;
        font-size: 12px;
        padding: 18px 0px;
        color: #555555;
        p:first-child {
          font-size: 26px;
          font-weight: bold;
          color: #272727;
          span {
            font-size: 12px;
          }
        }
      }
      .coupon-content-right {
        padding: 12px 0px 12px 19px;
        font-size: 12px;
        color: #555555;
        line-height: 20px;
        p:first-child {
          font-size: 16px;
          color: #33363c;
          font-weight: bold;
          line-height: 28px;
        }
      }
    }
  }
}
</style>
