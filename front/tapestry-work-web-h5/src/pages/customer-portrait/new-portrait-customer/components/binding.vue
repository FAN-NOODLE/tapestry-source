<template>
  <div class="binding-relation">
    <div class="box">
      <div
        v-for="item in relationList"
        :key="item.member_id"
        class="relation-item"
      >
        <image class="avatar" :src="item.thumb_avatar" />
        <div class="relation-detail">
          <div class="top">
            <div class="title">{{ item.sa_user_name }}</div>
            <div v-if="item.is_key_sa === 1" class="tag">主顾问</div>
          </div>
          <div class="center">添加时间：{{ item.created_at }}</div>
          <div class="bottom">所属店铺：{{ item.department_name }}</div>
        </div>
      </div>
    </div>

    <!-- 加载状态展示 -->
    <div class="bottom-tip" v-if="relationList.length > 0">已经到底了~</div>
    <empty-page
      v-else
      class="empty"
      picture="emptyList"
      paddingTop="70"
      paddingBottom="60"
      text="暂无内容"
    ></empty-page>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { userBindRelationApi } from "@/api/customer-portrait";
import emptyPage from "@/components/base/empty-page.vue";

@Component({
  components: {
    emptyPage,
  },
})
export default class Binding extends Vue {
  /** 外部联系人id */
  @Prop() openId!: string;

  /** 列表详情 */
  private relationList = [];

  /** 获取列表信息 */
  async getRelationList() {
    this.$toast.loading({
      message: "加载中...",
      forbidClick: true,
    });
    const params = { external_userid: this.openId };

    try {
      let { code, data, message }: IApiResult = await userBindRelationApi(
        params
      );
      this.$toast.clear()
      if (code !== 200) return this.$toast.fail(message);

      const newList = data.map((i) => ({
        ...i,
        created_at: this.$helper.dateFormat(i.created_at, "yyyy/MM/dd hh:mm:ss"),
      }));
      this.relationList = newList;
    } catch (error) {
      this.$toast.clear()
      console.log(error);
    }
  }

  private created() {
    this.getRelationList();
  }
}
</script>
<style lang="scss" scoped>
.binding-relation {
  
  .box {
    padding: 12px 15px 0;
  }

  .empty {
    background: #fff;
  }

  .relation-item {
    height: 109px;
    padding: 15px;
    display: flex;
    background-color: #fff;
    margin-bottom: 15px;

    .avatar {
      width: 28px;
      height: 28px;
      border-radius: 50%;
      margin-right: 10px;
    }

    .relation-detail {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: space-between;

      .top {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .title {
          font-size: 16px;
          font-family: PingFangSC-Medium, PingFang SC;
          font-weight: 500;
          color: #000000;
        }

        .tag {
          width: 54px;
          height: 24px;
          display: flex;
          justify-content: center;
          align-items: center;
          background-color: #f4f5f8;
          font-size: 12px;
          font-family: PingFangSC-Regular, PingFang SC;
          font-weight: 400;
          color: #676767;
        }
      }

      .center,
      .bottom {
        font-size: 12px;
        font-family: PingFangSC-Regular, PingFang SC;
        font-weight: 400;
        color: #676767;
      }
    }
  }

  .bottom-tip,
  .loading-box {
    width: 100%;
    font-size: 14px;
    color: #96979e;
    text-align: center;
  }
}
</style>
