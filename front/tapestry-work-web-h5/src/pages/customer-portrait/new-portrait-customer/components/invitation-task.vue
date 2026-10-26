<template>
  <div class="binding-relation">
    <List
      v-if="inviteList.length > 0"
      v-model="showLoading"
      @load="getList"
      :finished="!this.meta.isMore"
      finished-text="已经到底了~"
    >
      <div
        v-for="(item, index) in inviteList"
        :key="index"
        :class="{
          'relation-item': true,
          'item-bg': item.status === 1,
        }"
      >
        <invite-list :item="item" :role="role" :index="index" />
      </div>
    </List>
    <empty-page
      v-if="inviteList.length === 0"
      paddingTop="10"
      :picture="'emptyList'"
      :text="'暂无邀约记录'"
    ></empty-page>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import inviteList from "@/components/invite-list/index.vue";
import emptyPage from "@/components/base/empty-page.vue";
import { orderListApi } from "@/api/guide-invite/index";
import { List } from "vant";
@Component({
  components: {
    inviteList,
    emptyPage,
    List,
  },
})
export default class InvitationTask extends Vue {
  @Prop(Object) info!: any;
  role: Number = 3;
  showLoading: boolean = true;
  isError: boolean = false;
  inviteList: any[] = [];
  meta: any = {
    page: 0,
    pageSize: 10,
    isMore: true,
  };
  mounted() {
    this.getList();
  }
  /*    普通事件开始  */
  async getList(page = true) {
    if (!this.meta.isMore) return;
    this.showLoading = true;
    if (page) {
      this.meta.page += 1;
    }

    const params = {
      page: this.meta.page,
      size: this.meta.pageSize,
      external_userid: this.info.external_userid,
    };

    try {
      const { code, data, message }: any = await orderListApi(params);
      if (code !== 200) {
        this.$toast(message);
        return;
      }
      const { list, total } = data;
      if (list) {
        const newInviteList = page ? this.inviteList.concat(list) : list;
        this.inviteList = newInviteList;
        this.meta.total = total;
        this.meta.isMore = newInviteList.length < total;
        this.showLoading = false;
      } else {
        this.meta.isMore = false;
        this.showLoading = false;
      }
    } catch (error) {
      console.log(error);
      this.meta.isMore = Math.max(1, this.meta.page - 1);
      this.showLoading = false;
    }
  }
}
</script>
<style lang="sass" scoped>
.binding-relation
  padding: 15px 15px 0

  .relation-item
    width: 345px
    background: #FFFFFF
    padding: 15px
    margin-bottom: 15px
    position: relative

    &::after
      content: ""
      clear: both
      width: 346px
      height: 1px
      background: #F0F0F0
      position: absolute
      top: 46px
      left: 0

  .item-bg
    background-image: url("https://cdn-cos.coach.com.cn/uploads/2022/v3.0.1/weclient/list_bg.png")
    background-repeat: no-repeat
    background-size: 200px 200px
    background-position: bottom right

  .bottom-tip, .loading-box
    width: 100%
    font-size: 14px
    color: #96979E
    text-align: center
</style>
