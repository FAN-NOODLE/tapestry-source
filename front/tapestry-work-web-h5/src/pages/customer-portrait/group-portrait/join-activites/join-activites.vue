<template>
 <div id="join-activites">
    <van-nav-bar title="参加活动" left-arrow @click-left="$router.go(-1)" />
    <div style="padding-top: 45px;">
      <div v-if="list.length > 0" class="activites-list">
        <div class="item-list" v-for="item in list" :key="item.external_userid">
          <div class="item-right">
            <div class="item-name-box">
              <div class="item-name">{{ item.name || "-" }}</div>
              <div class="item-form">
                <div class="item-form-label">参与时间:</div>
                <div class="item-form-content">{{ item.join_time }}</div>
              </div>
            </div>
            <div class="item-action">
              {{item.status == 1 ? '进行中' : '已结束'}}
            </div>
          </div>
        </div>
      </div>
      <div v-else-if="!loading">
        <empty-page
          picture="emptyList"
          paddingTop="87"
          paddingBottom="60"
          text="暂无参与活动"
        ></empty-page>
      </div>
    </div>
 </div>
</template>
<script lang="ts">
import { NavBar } from "vant";
import { Component, Vue } from "vue-property-decorator";
import { chatGroupActInfo } from "@/api/customer-portrait";
import emptyPage from "@/components/base/empty-page.vue";
@Component({
  components: {
    [NavBar.name]: NavBar,
    emptyPage
  },
})
export default class extends Vue {
  private loading = false;
  private list = [];
  public get groupChatId(): string {
    return this.$store.getters["common/getGroupChatId"];
  }
  public get brand(): String {
    return this.$store.getters['user/getEnterpriseMark'];
  }
  private created() {
    if (this.brand !== 'HR') document.body.classList.add('other')

    if (this.$route.query.user_id) this.getActsListByMember()
  }
  // 群成员参与活动列表
  private async getActsListByMember() {
    console.log('请求数据', this.groupChatId)
    this.$toast.loading({ message: '加载中...', forbidClick: true});
    this.loading = true
    const {
      code,
      message = "",
      data,
    }: IApiResult = await chatGroupActInfo({chat_id: this.groupChatId, user_id: String(this.$route.query.user_id)});
    this.$toast.clear()
    this.loading = false
    if (code === 200) {
      if (data && data.length) {
        data.forEach((item: any) => {
          item.join_time =
            item.join_time &&
            this.$helper.dateFormat(item.join_time, "yyyy-MM-dd hh:mm");
        });
        this.list = data;
      }
    } else {
      this.$toast.fail(message);
    }
  }
}
</script>
<style lang="scss" scoped>
#join-activites{
  padding-bottom: 15px;
  ::v-deep{
    .van-nav-bar .van-icon {
      color: #000000;
    }
    .van-nav-bar__title {
      font-weight: bold;
      color: #000000;
      font-size: 17px;
    }
    .van-nav-bar__content {
      position: fixed;
      top: 0px;
      width: 100%;
      background: #ffffff;
    }
  }
  .activites-tit{
    height: 22px;
    font-size: 16px;
    font-weight: 400;
    color: #808080;
    line-height: 22px;
    margin: 15px 15px 0;
    position: relative;
    .activites-num{
      font-weight: 600;
      color: #33363C;
      position: absolute;
      right: 0;
    }
  }
  .activites-list {
    margin: 15px;
    background: #ffffff;
    box-shadow: 0px 2px 19px 0px rgba(235, 228, 217, 0.39);
    border-radius: 12px;
    border: 1px solid #ffffff;
    padding: 0 15px;
  }
  .item-list {
    display: flex;
    background: #ffffff;
    flex-shrink: 0;
    padding: 15px 0;
    border-bottom: 1px solid #f0f0f0;
    .item-index {
      font-size: 12px;
      padding: 0 12px;
      margin-bottom: 10px;
      display: flex;
      justify-content: center;
      align-items: center;
    }
    .item-img {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      overflow: hidden;
      margin-right: 8px;
      img {
        width: 100%;
        height: 100%;
        border-radius: 50%;
      }
    }
    .item-right {
      flex: 1;
      display: flex;
      justify-content: space-between;
      align-items: center;
      .item-name-box {
        display: flex;
        justify-content: space-between;
        flex-direction: column;
        .item-name {
          font-size: 14px;
          line-height: 20px;
          font-weight: 600;
          color: #33363c;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
      .item-action {
        width: 68px;
        height: 28px;
        background: #caa361;
        border-radius: 6px;
        font-size: 12px;
        line-height: 28px;
        color: #33363c;
        cursor: pointer;
        text-align: center;
      }
      .item-form {
        display: flex;
        align-items: center;
        margin-top: 6px;
        position: relative;
        font-size: 12px;
        line-height: 17px;
        color: #33363C;
        .item-form-label {
          margin-right: 5px;
        }
        .activity {
          font-size: 14px;
        }
      }
    }
    &:last-child{
      border: 0;
    }
  }
}
</style>
<style scoped lang="scss">
.other #join-activites{
  .item-action{
    background: #5B89FE !important;
    color: #ffffff !important;
  }
}
</style>
