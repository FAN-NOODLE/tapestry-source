<!--参与活动的弹窗-->
<template>
  <div class="staff-popup">
    <van-popup
      v-model="showStaffPopup"
      class="popup-box"
      position="bottom"
      closeable
      @click-overlay="closePopup"
      @click-close-icon="closePopup"
    >
      <div class="popup-title">参与活动</div>
      <div class="staff-list-box">
        <van-list
          v-if="activityList.length > 0"
          v-model="loading"
          :finished="noMore"
          finished-text="没有更多了~"
          @load="getMoreOrder"
        >
          <div
            class="item-staff"
            v-for="(item, index) in activityList"
            :key="index"
          >
            <div class="item-staff-right">
              <div class="item-staff-name">
                <span>{{ item.name }}</span>
              </div>
              <div class="item-staff-date">参与时间: {{ item.created_at }}</div>
            </div>
            <div class="item-staff-left">
              <div
                class="status-item"
                :class="[item.current_staus === 1 ? 'status-1' : '']"
              >
                {{ item.status }}
              </div>
            </div>
          </div>
        </van-list>
        <!-- 列表缺省-->
        <empty-page
          v-if="activityList.length == 0 && !loading"
          picture="noData"
          paddingTop="60"
          paddingBottom="0"
          text="暂无参与活动"
        ></empty-page>
      </div>
    </van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import { Popup } from "vant";
import emptyPage from "@/components/base/empty-page.vue";
import { getJoinActivity } from "@/api/customer-portrait";
@Component({
  components: {
    emptyPage,
    [Popup.name]: Popup,
  },
})
export default class extends Vue {
  @Prop({ default: "", type: String }) unionId!: string;
  private showStaffPopup = true;
  private activityList = [];
  private loading = false;
  private noMore = false;
  private meta = {
    page: 1,
    size: 20,
    total: 0,
  };
  /*************************** 生命周期 start ***************************/
  private created() {
    this.userMemberList();
  }
  /*************************** method start ***************************/
  // 活动列表
  private async userMemberList() {
    this.loading = true;
    const params = {
      union_id: this.unionId,
      page: this.meta.page,
      size: this.meta.size,
    };
    const { code, message = "", data }: IApiResult = await getJoinActivity(
      params
    );
    if (code === 200) {
      data.list.forEach((item: any) => {
        item.created_at =
          item.created_at &&
          this.$helper.dateFormat(item.created_at, "yyyy-MM-dd hh:mm");
        if (item.status === "进行中") {
          item.current_staus = 1;
        } else {
          item.current_staus = 0;
        }
      });
      this.activityList = this.activityList.concat(data.list);
      this.meta.total = data.total;
      this.loading = false;
      this.meta.page++;
      if (
        data.total == this.activityList.length ||
        !data.list ||
        (data.list && data.list.length == 0)
      ) {
        this.noMore = true;
      }
    } else {
      this.noMore = true;
      this.loading = false;
      this.$toast.fail(message);
    }
  }
  // 页面滑到底部触发加载更多
  private getMoreOrder() {
    this.userMemberList();
  }
  @Emit("close")
  private closePopup() {
    return;
  }
}
</script>
<style lang="scss" scoped>
.staff-popup {
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
    .staff-list-box {
      padding: 15px 15px 0 15px;
      height: 348px;
      overflow-y: auto;
      z-index: 1;
      .item-staff {
        display: flex;
        align-items: center;
        height: 40px;
        margin-bottom: 15px;
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
        .item-staff-right {
          width: 290px;
          .item-staff-name {
            max-width: 220px;
            width: min-content;
            color: #272727;
            font-size: 14px;
            line-height: 22px;
            font-weight: 600;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            position: relative;
            padding-right: 75px;
            .item-group-type {
              border: solid 1px #0095f6;
              width: 62px;
              height: 20px;
              font-size: 12px;
              line-height: 20px;
              text-align: center;
              color: #0095f6;
              border-radius: 2px;
              position: absolute;
              right: 3px;
              top: 1px;
            }
          }
          .item-staff-date {
            font-size: 12px;
            color: #808080;
            margin-top: 5px;
            line-height: 17px;
          }
        }
        .item-staff-left {
          width: 80px;
          text-align: center;
          .status-item {
            padding: 4px 0;
            font-size: 12px;
            background: #999;
            border-radius: 2px;
            color: #333333;
            opacity: 0.8;
          }
          .status-1 {
            color: #388df0;
          }
        }
      }
      .staff-no-more {
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
