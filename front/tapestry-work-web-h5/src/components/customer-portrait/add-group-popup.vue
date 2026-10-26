<!--加入的社群弹窗-->
<template>
  <div class="group-popup">
    <van-popup
      v-model="showGroupPopup"
      class="popup-box"
      position="bottom"
      closeable
      @click-overlay="closePopup"
      @click-close-icon="closePopup"
    >
      <div class="popup-title">加入的社群</div>
      <div class="group-list-box">
        <van-list
          v-if="groupList.length > 0"
          v-model="loading"
          :finished="noMore"
          finished-text="没有更多了~"
          @load="getMoreOrder"
        >
          <div
            class="item-group"
            v-for="(item, index) in groupList"
            :key="index"
          >
            <div class="item-img">
              <img src="@/assets/customer-portrait/group_icon.png" alt="" />
            </div>
            <div class="item-group-right">
              <div
                class="item-group-type"
                :class="{ 'not-in-group': item.member_is_join == 0 }"
              >
                {{ item.member_is_join == 0 ? "我没在群中" : "我已在群中" }}
              </div>
              <div class="name-box">
                <div class="item-group-name">
                  <span class="name">{{
                    item.group_chat_name || "未命名"
                  }}</span>
                  <div
                    class="copy-box"
                    @click="copyItem(item.group_chat_name)"
                    @mouseover.prevent="selectStyle(item)"
                    @mouseout.prevent="outStyle(item)"
                  >
                    <img
                      v-if="item.group_chat_name"
                      class="copy-icon"
                      src="@/assets/customer-portrait/copy_icon.png"
                      alt=""
                    />
                    <div class="tool-tip" v-if="item.isHover">复制群名称</div>
                  </div>
                </div>
              </div>
              <div class="item-group-date">群主: {{ item.member_name }}</div>
              <div class="item-group-date">入群时间: {{ item.join_time }}</div>
            </div>
          </div>
        </van-list>
        <!-- 列表缺省-->
        <empty-page
          v-if="groupList.length == 0 && !loading"
          picture="noData"
          paddingTop="60"
          paddingBottom="0"
          text="暂无社群"
        ></empty-page>
      </div>
    </van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Emit } from "vue-property-decorator";
import { Popup } from "vant";
import emptyPage from "@/components/base/empty-page.vue";
import { customerGroupListApi } from "@/api/customer-portrait";
@Component({
  components: {
    emptyPage,
    [Popup.name]: Popup,
  },
})
export default class extends Vue {
  private showGroupPopup = true;
  private groupList = []; // 列表
  private loading = false;
  private noMore = false;
  private ispC = false;
  private meta = {
    page: 1,
    size: 20,
    total: 0,
  };
  // 是否为PC端
  private monitorModel() {
    let userAgentInfo = navigator.userAgent;
    let Agents = [
      "Android",
      "iPhone",
      "SymbianOS",
      "Windows Phone",
      "iPad",
      "iPod",
    ];
    let flag = true;
    for (var v = 0; v < Agents.length; v++) {
      if (userAgentInfo.indexOf(Agents[v]) > 0) {
        flag = false;
        break;
      }
    }
    this.ispC = flag;
  }
  // 鼠标划入事件
  private selectStyle(item: any) {
    if (!this.ispC) return;
    const _this = this;
    this.$nextTick(function () {
      this.groupList.forEach(function (item: any) {
        _this.$set(item, "isHover", false);
      });
      this.$set(item, "isHover", true);
    });
  }
  // 鼠标划出事件
  private outStyle(item: any) {
    if (!this.ispC) return;
    this.$set(item, "isHover", false);
  }
  // 点击复制按钮
  private copyItem(data: string) {
    this.$helper.copyMessage(data);
  }
  // 群聊列表接口
  private async customerGroupList() {
    this.loading = true;
    const params = {
      user_id: this.$store.getters["user/getCustomerId"] || 121668888497026,
      page: this.meta.page,
      size: this.meta.size,
      status: 1
    };
    const { code, message = "", data }: IApiResult = await customerGroupListApi(
      params
    );
    if (code === 200) {
      data.list.forEach((item: any) => {
        item.join_time =
          item.join_time &&
          this.$helper.dateFormat(item.join_time, "yyyy-MM-dd hh:mm");
      });
      this.groupList = this.groupList.concat(data.list);
      this.meta.total = data.total;
      this.loading = false;
      this.meta.page++;
      if (
        data.total == this.groupList.length ||
        !data.list ||
        (data.list && data.list.length == 0)
      ) {
        this.noMore = true;
      }
    } else {
      this.loading = false;
      this.noMore = true;
      this.$toast.fail(message);
    }
  }
  // 页面滑到底部触发加载更多
  private getMoreOrder() {
    this.customerGroupList();
  }
  private created() {
    this.customerGroupList();
    this.monitorModel();
  }
  @Emit("close")
  private closePopup() {
    return;
  }
}
</script>
<style lang="scss" scoped>
::v-deep .toast-copy {
  padding: 15px;
  height: 78px;
}
.group-popup {
  .popup-box {
    border-radius: 4px 4px 0 0;
    ::v-deep .van-popup__close-icon {
      color: #808080;
      font-size: 14px;
      z-index: 2;
    }
    .popup-title {
      height: 50px;
      line-height: 50px;
      width: 100%;
      text-align: center;
      color: #272727;
      font-size: 16px;
      font-weight: 600;
      border-bottom: solid 1px #ebebeb;
      position: sticky;
      top: 0;
      z-index: 2;
    }
    .group-list-box {
      padding: 0 0 0 15px;
      height: 348px;
      overflow-y: auto;
      z-index: 1;
      .item-group {
        display: flex;
        align-items: flex-start;
        min-height: 40px;
        padding-bottom: 15px;
        padding-top: 15px;
        border-bottom: solid 1px #ebebeb;
        position: relative;
        .item-group-type {
          border: solid 1px #74bf3c;
          width: 74px;
          height: 20px;
          font-size: 12px;
          line-height: 20px;
          text-align: center;
          color: #74bf3c;
          border-radius: 2px;
          position: absolute;
          right: 15px;
          top: 15px;
        }
        .not-in-group {
          border: solid 1px #f37e0f;
          color: #f37e0f;
        }
        .item-img {
          width: 40px;
          height: 40px;
          border-radius: 2px;
          overflow: hidden;
          margin-right: 10px;
          flex-shrink: 0;
          img {
            width: 100%;
            height: 100%;
          }
        }
        .item-group-right {
          // width: 290px;
          .name-box {
            position: relative;
            width: min-content;
          }
          .item-group-name {
            max-width: 182px;
            width: min-content;
            color: #272727;
            font-size: 14px;
            line-height: 20px;
            font-weight: 600;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            padding-right: 20px;
            .copy-box {
              cursor: pointer;
              width: 20px;
              height: 20px;
              position: absolute;
              top: 0;
              right: 0;
              .copy-icon {
                width: 13px;
                height: 15px;
                position: absolute;
                left: 0;
                right: 0;
                top: 0;
                bottom: 0;
                margin: auto;
              }
            }
            .tool-tip {
              height: 21px;
              line-height: 21px;
              padding: 0 5px;
              font-size: 12px;
              color: #ffffff;
              background: rgba(0, 0, 0, 0.8);
              width: min-content;
              position: absolute;
              top: -25px;
              right: -25px;
              border-radius: 4px;
              font-weight: 500;
            }
          }
          .item-group-date {
            font-size: 12px;
            color: #808080;
            margin-top: 7px;
            line-height: 17px;
          }
        }
      }
      .item-group:last-child {
        border: none;
      }
      .group-no-more {
        text-align: center;
        font-size: 12px;
        color: #c8c8c8;
        line-height: 17px;
        padding-bottom: 20px;
        padding-top: 10px;
      }
      .no-more-fix {
        position: absolute;
        bottom: 0;
        width: 100%;
      }
    }
  }
}
</style>
