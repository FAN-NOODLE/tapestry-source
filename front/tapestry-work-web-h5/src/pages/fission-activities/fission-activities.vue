<template>
  <div class="fission-activities-page">
    <!-- 统计模块 -->
    <div class="statistics-box">
      <div class="statistics-label">统计</div>
      <div class="first-card">
        <div class="item-statistics">
          <div class="item-num">{{ info.user_num || 0 }}</div>
          <div class="item-title">发起人数</div>
        </div>
        <div class="item-statistics">
          <div class="item-num">{{ info.fission_num || 0 }}</div>
          <div class="item-title">好友裂变数</div>
        </div>
        <div class="item-statistics">
          <div class="item-num">{{ info.run_away_num || 0 }}</div>
          <div class="item-title">好友流失数</div>
        </div>
      </div>
      <div class="second-card m-t-10" v-if="isShowSencond">
        <div class="item-statistics">
          <div class="item-num">{{ info.leave_group_num || 0 }}</div>
          <div class="item-title">群成员退群数</div>
        </div>
        <div class="item-statistics">
          <div class="item-num">{{ info.join_group_num || 0 }}</div>
          <div class="item-title">群成员裂变数</div>
        </div>
        <div class="item-statistics unshow">
          <div class="item-num">{{ 1 || 0 }}</div>
          <div class="item-title">群成员裂变数</div>
        </div>
      </div>
      <div class="fold-btn" @click="changeFn">
        <span>{{ isShowSencond ? "收起" : "展开" }}</span>
        <img
          :src="
            isShowSencond
              ? require('@/assets/customer-portrait/arrow_up.png')
              : require('@/assets/customer-portrait/arrow_down.png')
          "
          alt=""
        />
      </div>
    </div>
    <!-- 输入框模块 -->
    <div class="input-select-box">
      <van-search
        v-model="activityName"
        clearable
        placeholder="请输入活动名称"
        @input="search"
      >
        <template #left-icon>
          <img
            src="@/assets/fission-activities/search.png"
            alt=""
            class="search-img"
          />
        </template>
      </van-search>
      <van-dropdown-menu active-color="#0095F6">
        <van-dropdown-item
          v-model="groupVal"
          :options="groupOption"
          @change="search"
        />
      </van-dropdown-menu>
    </div>
    <!-- 分割 -->
    <div class="partition-line"></div>
    <!-- 列表模块 -->
    <div class="table-box">
      <van-list
        v-model="isLoading"
        :finished="isFinished"
        :finished-text="list.length ? '没有更多了~' : ''"
        @load="getTableData"
        :immediate-check="false"
        offset="100"
      >
        <div
          v-for="(item, index) in list"
          :key="index"
          :class="{
            'link-box': item.status === 1 || item.status === 4,
            'link-boxs': item.status === 2 || item.status === 3,
          }"
        >
          <div class="title">
            {{ item.name }}
          </div>
          <div class="count">
            <div class="count-title">{{ desc }}</div>
            <img :src="item.img_url" alt="" />
            <div class="count-status">
              <van-tag
                plain
                class="success-tag width-size"
                v-if="item.status === 4"
                >进行中</van-tag
              >
              <van-tag
                plain
                class="close-tag width-size"
                v-else-if="item.status === 2"
                >已下架</van-tag
              >
              <van-tag
                plain
                class="close-tag width-size"
                v-else-if="item.status === 3"
                >已结束</van-tag
              >
              <van-tag plain class="unstart-tag width-size" v-else
                >未开始</van-tag
              >
            </div>
          </div>
          <van-divider v-if="item.status === 1 || item.status === 4" />
          <div class="footer" v-if="item.status === 1 || item.status === 4">
            <van-button type="info" @click="sendChatMessage(item)"
              >发送</van-button
            >
          </div>
        </div>
        <empty-page
          v-if="list.length == 0 && !isLoading"
          picture="noData"
          paddingTop="60"
          paddingBottom="60"
          :text="activityName || groupVal ? '没有找到相关活动' : '暂无裂变活动'"
        ></empty-page>
      </van-list>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import emptyPage from "@/components/base/empty-page.vue";
import { getFissionActivityListApi } from "@/api/fission-activities";
import wx from 'weixin-js-sdk'
@Component({
  components: {
    emptyPage,
  },
})
export default class extends Vue {
  /******************************************** data开始 ********************************************/
  private isShowSencond = false; // 是否展开第二个卡片
  private activityName = ""; // 活动名称
  private groupVal = 0; // 活动状态
  private desc = "邀请好友助力，好礼兑不停";
  // 活动状态内容
  private groupOption: Array<ISelect> = [
    { text: "全部", value: 0 },
    { text: "未开始", value: 1 },
    { text: "进行中", value: 4 },
    { text: "已下架", value: 2 },
    { text: "已结束", value: 3 },
  ];
  private isLoading = false; // 结束加载状态
  private isFinished = false; // 下拉结束
  private info: IInfo = {
    user_num: 0, // 发起人数
    fission_num: 0, // 好友裂变数
    run_away_num: 0, // 流失人数
    join_group_num: 0, // 群成员裂变数
    leave_group_num: 0, // 群成员退群数
  }; // 统计数据
  private list: Array<IList> = []; // 列表数据
  // 分页
  private meta: IMeta = {
    page: 1,
    size: 20,
    total: 0,
  };
    private  async created () { // 当前需要直接退出的页面
    window.history.pushState('_android_back', '', location.href)
    window.addEventListener("popstate", this.popstate, false)
  }
  private  async  destroyed () { // 记得清除，不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false)
  }

  /******************************************** 生命周期 ********************************************/
  private mounted() {
    this.getTableData();
  }
  /******************************************* methods开始 ******************************************/
  // 点击切换卡片显隐
  private changeFn() {
    this.isShowSencond = !this.isShowSencond;
  }
  // 查询
  // 防抖会合并事件且不会去触发事件，当一定时间内没有触发再这个事件时，才真正去触发事件，即取到最后的结果
  private search = this.$helper.debounce((value: string) => {
    console.log(value);
    this.meta.page = 1;
    this.meta.size = 20;
    this.isFinished = false;
    this.list = [];
    if (!this.isLoading) {
      this.getTableData("search");
    }
  }, 500);
  // 列表数据
  private async getTableData(type?: string) {
    // 开启加载状态
    this.isLoading = true;
    // 请求接口
    const params: IFissionActivityList = {
      search_text: this.activityName,
      status: this.groupVal,
      page: this.meta.page,
      size: this.meta.size,
    };
    const {
      code,
      message = "",
      data,
    }: IApiResult = await getFissionActivityListApi(params);
    if (code === 200) {
      if (data.list) {
        this.info = data.statistics; // 统计数据
        type && type === "search"
          ? (this.list = data.list)
          : (this.list = this.list.concat(data.list)); // 列表数据(type存在时为搜索，反之为加载)
        this.meta.total = data.total; // 接收数据总数
        this.isLoading = false; // 关闭加载状态
        this.meta.page && this.meta.page++; // 分页 页数自增
        this.isFinished = this.list.length >= data.total; // 判断是否下拉更新完成
      }
      // 判断列表总数据大小是否和获取到的列表数据大小相等，
      // 当接口未返回数据，接口数据返回且数据长度为0的情况时，意味着没有更多的数据了
      // 即页面展示列表获取道了所有数据，完成下拉加载，没有更多的数据了
      if (
        data.total == this.list.length ||
        !data.list ||
        (data.list && data.list.length == 0)
      ) {
        this.isFinished = true; // 下拉更新完成
      }
    } else {
      this.isLoading = false; // 请求失败，没有加载新的数据，加载状态完成
      this.isFinished = true; // 没有数据，下拉更新完成
      this.$toast.fail(message); // 提示错误
    }
  }
  // 发送
  private sendChatMessage(item: IList) {
    // tost提示
    this.$toast.loading({
      message: "生成中...",
      forbidClick: true,
    });
    // 参数(link)
    const message = {
      msgtype: "news",
      news: {
        link: item.html_url,
        title: item.name,
        desc: this.desc,
        imgUrl: item.img_url,
      },
    };
    // 调用企微接口
    window.wx.invoke("sendChatMessage", message, function (res) {
      console.log(res);
      this.$toast.clear(); // 清除tost提示
      if (res.err_msg == "sendChatMessage:ok") {
        //发送成功
        console.log("成功回调");
      }
    });
  }
  private popstate () {
    wx.closeWindow() // 微信网页退出
  }
}
</script>
<style lang="scss" scoped>
.m-t-10 {
  margin-top: 10px;
}
.m-t-15 {
  margin-top: 15px;
}
.fission-activities-page {
  min-height: 100vh;
  background: white;
}
.statistics-box {
  padding: 15px 15px;
  .statistics-label {
    font-size: 14px;
    line-height: 20px;
    color: #808080;
    margin-bottom: 5px;
  }
  .first-card {
    background: #f4f5f8;
    height: 73px;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 15px;
    .item-statistics {
      text-align: center;
      .item-num {
        color: #333333;
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 4px;
        .currency {
          font-size: 12px;
        }
      }
      .item-title {
        font-size: 12px;
        line-height: 17px;
        color: #999999;
      }
    }
  }
  .second-card {
    background: #f4f5f8;
    height: 73px;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 15px;
    .item-statistics {
      text-align: center;
      .item-num {
        color: #333333;
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 4px;
        .currency {
          font-size: 12px;
        }
      }
      .item-title {
        font-size: 12px;
        line-height: 17px;
        color: #999999;
      }
    }
    .unshow {
      visibility: hidden;
    }
  }
  .fold-btn {
    font-size: 12px;
    color: #999999;
    display: flex;
    align-items: center;
    height: 37px;
    width: 100%;
    justify-content: center;
    cursor: pointer;
    img {
      height: 15px;
      width: 15px;
      margin-left: 4px;
    }
  }
}
.input-select-box {
  padding: 15px 15px;
  padding-top: 0px;
  display: flex;
  align-items: center;
  // height: 36px;
  // margin-bottom: 6px;
  background: #ffffff;
  ::v-deep .van-ellipsis {
    color: #a6a6a6 !important;
  }
  ::v-deep.van-search {
    width: 50%;
    margin-right: 10px;
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
    flex: 1;
  }
  ::v-deep.van-dropdown-menu__title {
    font-size: 14px;
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
.partition-line {
  width: 100%;
  height: 6px;
  background: #f5f5f5;
}
.table-box {
  margin-top: 15px;
  margin-bottom: 15px;
  .link-box {
    width: 345px;
    min-height: 172px;
    max-height: 192px;
    padding: 10px;
    margin: 0 auto;
    border-radius: 4px;
    border: 1px solid#EBEBEB;
    margin-bottom: 15px;
    .title {
      width: 325px;
      min-height: 20px;
      max-height: 40px;
      font-size: 14px;
      color: #000000;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    .count {
      margin-top: 5px;
      display: flex;
      flex-wrap: wrap;
      .count-title {
        width: 240px;
        height: 39.1px;
        font-size: 14px;
        color: hsla(0, 0%, 65%, 1);
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
      img {
        width: 80px;
        height: 80px;
      }
    }
    .van-divider {
      height: 1px;
      margin-top: 10px;
      margin-bottom: 10px;
      border-color: hsla(0, 0%, 92%, 1);
    }
    .footer {
      height: 26px;
      display: flex;
      align-items: center;
      justify-content: flex-end;
      &.coupon {
        height: 56px;
      }
      .van-button {
        width: 74px;
        height: 26px;
        background: linear-gradient(130deg, #0285ec 0%, #0dadf6 100%);
        border-radius: 4px;
      }
    }
    &:hover {
      border-color: #0095f6;
    }
  }
  .link-boxs {
    width: 345px;
    min-height: 126px;
    max-height: 146px;
    padding: 10px;
    margin: 0 auto;
    border-radius: 4px;
    border: 1px solid#EBEBEB;
    margin-bottom: 15px;
    .title {
      width: 325px;
      min-height: 20px;
      max-height: 40px;
      font-size: 14px;
      color: #000000;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    .count {
      margin-top: 5px;
      display: flex;
      flex-wrap: wrap;
      height: 80px;
      .count-title {
        width: 240px;
        height: 40px;
        font-size: 14px;
        color: hsla(0, 0%, 65%, 1);
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
      img {
        width: 80px;
        height: 80px;
      }
    }
    .van-divider {
      height: 1px;
      margin-top: 10px;
      margin-bottom: 10px;
      border-color: hsla(0, 0%, 92%, 1);
    }
    .footer {
      height: 26px;
      display: flex;
      align-items: center;
      justify-content: flex-end;
      &.coupon {
        height: 56px;
      }
      .van-button {
        width: 74px;
        height: 26px;
        background: linear-gradient(130deg, #0285ec 0%, #0dadf6 100%);
        border-radius: 4px;
      }
    }
    &:hover {
      border-color: #0095f6;
    }
  }
}
.success-tag {
  color: #77d034;
  border-color: #77d034;
}
.close-tag {
  color: #a6a6a6;
  border-color: #a6a6a6;
}
.unstart-tag {
  color: #ff8f00;
  border-color: #ff8f00;
}
.width-size {
  // width: 50px;
  height: 20px;
  font-size: 12px;
  padding: 0px 7px;
  white-space: nowrap;
}
.count-status {
  height: 20px;
  width: 150px;
  display: flex;
  margin-top: -20px;
}
.search-img {
  width: 14px;
  height: 14px;
}
::v-deep .van-field__left-icon {
  margin-top: 3px;
}
</style>
