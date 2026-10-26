<template>
  <div class="marketing-promotion-page">
    <van-tabs v-model="active" sticky @change="handleTabChange">
      <van-tab title="商品" name="goods">
        <div class="outline-box">
          <van-search
            v-model="keyWord"
            @input="onGoodsSearch"
            @clear="onGoodsSearch"
            placeholder="请输入搜索关键词"
          />
          <van-list
            v-model="goodsLoading"
            :finished="goodFinished"
            :finished-text="goodsListD.length ? '没有更多了' : ''"
            @load="goodsOnLoad"
          >
            <div class="good" v-for="item in goodsListD" :key="item.id">
              <div class="link-box">
                <img :src="item.sku_picture" alt="商品图片" />
                <div class="link-rig-box">
                  <div class="sku-name moreLlipsis">
                    {{ item.sku_name }}
                  </div>
                  <div class="attr ellipsis">{{ item.attr_str }}</div>
                  <div class="pri-cou">
                    <span>
                      <span class="pro">价格</span>
                      <span class="val">¥{{ item.sell_price / 100 }}</span>
                    </span>
                    <span>
                      <span class="pro">库存</span>
                      <span class="val">{{
                        item.stock > 10000 ? item.stock / 10000 : item.stock
                      }}</span>
                      <span class="pro" v-if="item.stock > 10000">万+</span>
                    </span>
                  </div>
                </div>
              </div>
              <div class="btn-box">
                <van-button
                  type="default"
                  class="m-r-10 resizeStyle"
                  @click="createPannelPic(item, 1)"
                  >发送海报</van-button
                >
                <van-button type="info" @click="sendChatMessage(item)"
                  >发送小程序</van-button
                >
              </div>
            </div>
          </van-list>
          <empty-page
            v-if="!goodsListD.length && !goodsLoading"
            picture="noData"
            paddingTop="60"
            paddingBottom="265"
            :text="keyWord ? '没有找到相关商品' : '暂无商品'"
          ></empty-page>
        </div>
      </van-tab>
      <van-tab title="活动" name="active">
        <div class="outline-box">
          <div class="header">
            <van-search
              v-model="activeSearchValue"
              clearable
              placeholder="请输入商品名称"
              @input="onActiveSearch"
            />
            <van-dropdown-menu active-color="#0095F6">
              <van-dropdown-item
                v-model="activeGroupValue"
                @change="handleGroupChange"
                :options="activeGroupOption"
              />
            </van-dropdown-menu>
          </div>
          <van-list
            v-model="activeLoading"
            :finished="activeFinished"
            :finished-text="activiListD.length ? '没有更多了' : ''"
            @load="activeOnLoad"
          >
            <div class="good" v-for="(item, index) in activiListD" :key="index">
              <div class="link-box">
                <img :src="item.image" alt="" />
                <div class="link-rig-box">
                  <div class="sku-name moreLlipsis">
                    {{ item.sku_name }}
                  </div>
                  <div class="attr ellipsis">{{ item.attr_str }}</div>
                  <div class="pri-cou">
                    <span class="flex-center">
                      <span class="tag">
                        {{ item.activity_type | avtiveName(item) }}
                      </span>
                      <span class="val">¥{{ item.activity_price / 100 }}</span>
                    </span>
                    <span>
                      <span class="pro">售卖价</span>
                      <span class="val">¥{{ item.sell_price / 100 }}</span>
                    </span>
                  </div>
                </div>
              </div>
              <div class="btn-box">
                <van-button
                  type="default"
                  class="m-r-10 resizeStyle"
                  @click="createPannelPic(item, item.activity_type)"
                  >发送海报</van-button
                >
                <van-button type="info" @click="sendChatMessage(item)"
                  >发送小程序</van-button
                >
              </div>
            </div>
          </van-list>
          <empty-page
            v-if="!activiListD.length && !activeLoading"
            picture="noData"
            paddingTop="60"
            paddingBottom="265"
            :text="
              activeSearchValue || activeGroupValue
                ? '没有找到相关活动'
                : '暂无活动'
            "
          ></empty-page>
        </div>
      </van-tab>
      <van-tab title="微信优惠劵" name="coupon">
        <div class="outline-box">
          <van-search
            v-model="coupon"
            @input="onCouponSearch"
            @clear="onCouponSearch"
            placeholder="请输入优惠券批次号"
          />
          <van-list
            v-model="couponLoading"
            :finished="couponFinished"
            :finished-text="couponListD.length ? '没有更多了' : ''"
            @load="couponOnLoad"
          >
            <div class="coupon-list">
              <div
                class="coupon-item"
                v-for="(item, index) in couponListD"
                :key="index"
              >
                <div class="coupon-box">
                  <div class="coupon-left">
                    <div class="pri">¥{{ item.coupon_value / 100 }}</div>
                    <div class="man">每人限领{{ item.max_quota }}张</div>
                  </div>
                  <div class="coupon-right">
                    <p class="p1">
                      {{ item.coupon_name }}
                    </p>
                    <p class="p3">批次号: {{ item.coupon_stock_id }}</p>
                    <p class="p3">
                      有效期:
                      {{ item.start_time | dateFormat("yyyy.MM.dd") }}-{{
                        item.end_time | dateFormat("yyyy.MM.dd")
                      }}
                    </p>

                    <div class="part" @click="open(item)">
                      <span>{{
                        item.is_all == 1 ? "全部商品可用" : "部分商品可用"
                      }}</span>
                      <img
                        v-if="item.is_all != 1"
                        src="../../assets/coopon-arrow.png"
                        alt=""
                      />
                    </div>
                  </div>
                </div>
                <div class="btn-box coupon">
                  <van-button
                    type="default"
                    class="m-r-10 resizeStyle"
                    @click="createPannelPic(item, 2)"
                    >发送海报</van-button
                  >
                  <van-button type="info" @click="sendChatMessage(item)"
                    >发送小程序</van-button
                  >
                </div>
              </div>
            </div>
          </van-list>
          <empty-page
            v-if="!couponListD.length && !couponLoading"
            picture="noData"
            paddingTop="60"
            paddingBottom="265"
            :text="coupon ? '没有找到相关优惠劵' : '暂无优惠劵'"
          ></empty-page>
        </div>
        <available-Product
          v-if="showavailableProduct"
          :coupon-id="couponId"
          @close="showavailableProduct = false"
        ></available-Product>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import availableProduct from "@/components/marketing/available-product.vue";
import emptyPage from "@/components/base/empty-page.vue";
import {
  goodListApi,
  activityListApi,
  wxCouponListApi,
  downloadPoster,
} from "@/api/marketing-promotion/index";
import helper from "@/utils/helper";
@Component({
  components: {
    emptyPage,
    availableProduct,
  },
  filters: {
    dateFormat: helper.dateFormat,
    avtiveName: (value, data) => {
      const activeTypeStr = {
        "14": "拼团",
        "15": "砍价",
        "12": "限时购",
        "6": "助力付邮",
        "7": "助力买赠",
        "13": "助力抽奖",
      };
      if (value === 14) {
        return `${data.need_num}人团`;
      }
      return activeTypeStr[value];
    },
  },
})
export default class extends Vue {
  private active = 0;
  private coupon = "";
  private showavailableProduct = false;
  // 商品
  private keyWord = "";
  private goodsListD: Array<any> = [];
  private goodsLoading = false;
  private goodFinished = false;
  private goodsPage = 1;
  // 活动
  private activeSearchValue = "";
  private activeGroupValue = 0;
  private activeGroupOption: Array<any> = [
    { text: "全部活动类型", value: 0 },
    // { text: "拼团", value: 14 },
    // { text: "砍价", value: 15 },
    // { text: "限时购", value: 12 },
    // { text: "助力付邮", value: 6 },
    // { text: "助力买赠", value: 7 },
    { text: "助力抽奖", value: 13 },
  ];
  private activiListD: Array<any> = [];
  private activeLoading = false;
  private activeFinished = false;
  private activePage = 1;
  // 优惠劵
  private couponListD: Array<any> = [];
  private couponLoading = false;
  private couponFinished = false;
  private couponPage = 1;
  private couponId = 0;

  private handleTabChange(name) {
    switch (name) {
      case "goods":
        this.goodsOnLoad(true);
        break;
      case "active":
        this.activeOnLoad(true);
        break;
      default:
        this.couponOnLoad(true);
        break;
    }
  }

  // 商品列表
  private async goodsOnLoad(isfirst?: boolean) {
    if (isfirst) {
      this.goodsPage = 1;
      this.goodsListD = [];
      this.goodsLoading = true;
      this.goodFinished = false;
    }
    if (this.goodsPage == 2) {
      this.goodsPage = 3;
    }
    const params = {
      sku_name: this.keyWord,
      page: this.goodsPage,
      size: this.goodsPage == 1 ? 20 : 10,
    };
    const { code, message, data }: IApiResult = await goodListApi(params);
    this.goodsLoading = false;
    if (code !== 200) {
      console.log(message);
      this.goodFinished = true;
      return;
    }
    const list = data.list.map((item) => ({
      ...item,
      attr_str: item.attributes
        .map((item) => item.attribute_item_value)
        .join(";"),
    }));
    if (isfirst) {
      this.goodsListD = list;
    } else {
      this.goodsListD = [...this.goodsListD, ...list];
    }
    this.goodFinished = this.goodsListD.length >= data.total;
    this.goodsPage++;
  }
  //搜索商品
  private onGoodsSearch = this.$helper.debounce(() => {
    this.goodsOnLoad(true);
  }, 500);

  private async activeOnLoad(isfirst?: boolean) {
    if (isfirst) {
      this.activePage = 1;
      this.activiListD = [];
      this.activeLoading = true;
      this.activeFinished = false;
    }
    if (this.activePage == 2) {
      this.activePage = 3;
    }
    const params = {
      activity_type: this.activeGroupValue,
      goods_name: this.activeSearchValue,
      page: this.activePage,
      size: this.activePage == 1 ? 20 : 10,
    };
    const { code, message, data }: IApiResult = await activityListApi(params);
    this.activeLoading = false;
    if (code !== 200) {
      console.log(message);
      this.activeFinished = true;
      return;
    }
    const list = data.list.map((item) => ({
      ...item,
      attr_str: item.attributes
        .map((item) => item.attribute_item_value)
        .join(";"),
    }));
    if (isfirst) {
      this.activiListD = list;
    } else {
      this.activiListD = [...this.activiListD, ...list];
    }
    this.activeFinished = this.activiListD.length >= data.total;
    this.activePage++;
  }

  private onActiveSearch = this.$helper.debounce(() => {
    this.activeOnLoad(true);
  }, 500);

  private handleGroupChange() {
    this.activeOnLoad(true);
  }

  // 优惠劵列表
  private async couponOnLoad(isfirst?: boolean) {
    if (isfirst) {
      this.couponPage = 1;
      this.couponListD = [];
      this.couponLoading = true;
      this.couponFinished = false;
    }
    if (this.couponPage == 2) {
      this.couponPage = 3;
    }
    const params = {
      search_text: this.coupon,
      page: this.couponPage,
      size: this.couponPage === 1 ? 20 : 10,
    };
    const { code, message, data }: IApiResult = await wxCouponListApi(params);
    this.couponLoading = false;
    if (code !== 200) {
      console.log(message);
      this.couponFinished = true;
      return;
    }
    if (isfirst) {
      this.couponListD = data.list;
    } else {
      this.couponListD = [...this.couponListD, ...data.list];
    }
    this.couponFinished = this.couponListD.length >= data.total;
    this.couponPage++;
  }

  //搜索优惠卷
  private onCouponSearch = this.$helper.debounce(() => {
    this.couponOnLoad(true);
  }, 500);

  private open(item) {
    if (item.is_all == 1) {
      return;
    }
    this.couponId = item.coupon_id;
    this.showavailableProduct = true;
  }
  // 发送海报
  private async createPannelPic(item, type) {
    this.$toast.loading({
      message: "生成中...",
      forbidClick: true,
    });
    const params = { type, operator_type: 1, params: JSON.stringify(item) };
    const { code, message, data }: IApiResult = await downloadPoster(params);
    if (code !== 200) {
      console.log(message);
      return;
    }
    window.wx.invoke(
      "sendChatMessage",
      {
        msgtype: "image",
        image: { mediaid: data.object },
      },
      (res) => {
        this.$toast.clear();
        console.log(res);
      }
    );
  }
  private async sendChatMessage(item) {
    this.$toast.loading({
      message: "发送中...",
      forbidClick: true,
    });
    const title = item.sku_name ? item.sku_name : item.coupon_name;
    const img = item.sku_picture
      ? item.sku_picture
      : item.image
      ? item.image
      : item.default_picture;
    console.log({
      msgtype: "miniprogram",
      miniprogram: {
        msgtype: "miniprogram",
        appid: item.share_app_id,
        title: title,
        imgUrl: img,
        page: item.share_path,
      },
    });
    let t = this;
    window.wx.invoke(
      "sendChatMessage",
      {
        msgtype: "miniprogram",
        miniprogram: {
          appid: item.share_app_id,
          title: title,
          imgUrl: img,
          page: item.share_path,
        },
      },
      (res) => {
        if (res.err_msg !== "sendChatMessage:ok") {
          if (res.err_msg == "sendChatMessage:cancel") {
            return;
          }
          //发送成功
          t.$toast.fail("企微官方后台未关联小程序");
        }
        // t.$toast.clear();
        console.log(res);
      }
    );
  }
}
</script>
<style lang="scss" scoped>
.marketing-promotion-page {
  min-height: 100vh;
  background: white;
  .m-r-10 {
    margin-right: 10px;
  }
}
::v-deep.van-tabs--line .van-tabs__wrap {
  border-bottom: 1px solid #ebebeb;
}
.van-search {
  padding: 0;
  margin-bottom: 15px;
  background: #f4f5f8;
  border-radius: 4px;
}
.header {
  padding: 0 14px;
  display: flex;
  align-items: center;
  // height: 36px;
  // margin-bottom: 6px;
  background: #ffffff;
  ::v-deep.van-search {
    width: 50%;
    padding: 0px;
  }
  ::v-deep.van-search__content {
    background-color: #f4f5f8;
    border-radius: 4px;
  }
  ::v-deep.van-dropdown-menu {
    display: flex;
    align-items: center;
    height: 36px;
    margin-left: 10px;
    flex: 1;
  }
  ::v-deep.van-ellipsis {
    color: #333333;
  }
  ::v-deep.van-dropdown-menu__bar {
    width: 100%;
    height: 100%;
    box-shadow: none;
    background: #f4f5f8;
    border: none;
    border-radius: 4px;
  }
}
// 一行显示
.ellipsis {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
// 两行显示
.moreLlipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}
.outline-box {
  padding: 15px;
  background: white;
  ::v-deep {
    .van-search {
      margin: 0;
    }
  }
  .good {
    padding: 10px 10px 0 10px;
    width: 345px;
    background: #ffffff;
    border-radius: 2px;
    border: 1px solid #ebebeb;
    margin: 15px 0;
  }
  .link-box {
    padding: 10px;
    width: 325px;
    height: 100px;
    background: #f4f5f8;
    border-radius: 2px;
    display: flex;
    img {
      width: 72px;
      height: 80px;
      margin-right: 15px;
    }
    .link-rig-box {
      width: 215px;
      height: 80px;
      .pri-cou {
        display: flex;
        justify-content: space-between;
        align-items: center;
        color: #333333;
        height: 16px;
        line-height: 16px;
        .flex-center {
          display: flex;
          align-items: center;
        }
        .tag {
          padding: 3px;
          height: 16px;
          line-height: 10px;
          text-align: center;
          background: linear-gradient(130deg, #0285ec 0%, #0dadf6 100%);
          border-radius: 2px;
          margin-right: 5px;
          color: white;
          font-size: 10px;
        }
        .pro {
          font-size: 12px;
          margin-right: 5px;
        }
        .val {
          font-size: 14px;
          font-weight: 500;
        }
      }
      .sku-name {
        font-size: 12px;
        font-weight: 500;
        color: #333333;
        height: 36px;
        line-height: 18px;
        margin-bottom: 4px;
      }
      .attr {
        font-weight: 400;
        color: #96979e;
        font-size: 12px;
        line-height: 16px;
        margin-bottom: 6px;
      }
    }
  }
  .btn-box {
    height: 46px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    &.coupon {
      height: 56px;
    }
    .van-button {
      height: 26px;
      border-radius: 4px;
    }
    .van-button:first-child {
      font-size: 14px;
      background: #fff;
      border: 1px solid #b9bec4;
      color: #333;
      font-weight: normal;
      padding-left: 16px;
      padding-right: 16px;
    }
    .van-button:last-child {
      background: linear-gradient(130deg, #0285ec 0%, #0dadf6 100%);
    }
  }
}
.coupon-item {
  height: 191px;
  padding: 18px 0 0 0;
  position: relative;
  &::before {
    content: " ";
    height: 1px;
    background: #ebebeb;
    width: calc(100% + 15px);
    position: absolute;
    bottom: 0;
    left: 0;
  }
  &:last-child {
    &::before {
      height: 0;
    }
  }
  .coupon-box {
    width: 345px;
    height: 115px;
    background-image: url("../../assets/coupon.png");
    background-repeat: no-repeat;
    background-size: 100% 100%;
    display: flex;
    .coupon-left {
      width: 115px;
      padding: 33px 0 0 0;
      text-align: center;
      .pri {
        font-size: 22px;
        font-weight: 500;
        color: #0387ed;
        line-height: 30px;
      }
      .man {
        font-size: 14px;
        font-weight: 400;
        color: #0387ed;
        line-height: 20px;
      }
      .max {
        font-size: 12px;
        line-height: 20px;
        font-weight: 400;
        color: rgba($color: #0387ed, $alpha: 0.7);
      }
    }
    .coupon-right {
      width: 230px;
      padding: 9px 10px 0;
      position: relative;
      p {
        margin-bottom: 3px;
        color: #0387ed;
        font-size: 12px;
        font-weight: 400;
        line-height: 20px;
      }
      .p1 {
        font-weight: 500;
        color: #0387ed;
      }
      .p3 {
        opacity: 0.7;
      }
      .part {
        display: flex;
        justify-content: space-between;
        align-items: center;
        height: 30px;
        position: absolute;
        bottom: 2px;
        left: 10px;
        line-height: 30px;
        width: calc(100% - 20px);
        span {
          opacity: 0.7;
          font-size: 12px;
          color: #0387ed;
        }
        img {
          display: inline-block;
          width: 5px;
          height: 8px;
        }
      }
    }
  }
}
</style>
