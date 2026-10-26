<template>
  <!-- 待分配好友情况 -->
  <div class="not-customer-detail">
    <div class="box-in">
      <div class="customer-ul">
        <div
          v-for="(item, index) in user_list"
          :key="index"
          class="customer-li"
        >
          <img
            :src="item.external_avatar"
            class="avatar"
          />
          {{ item.external_nickname }}
        </div>
        <layout-empty
          v-if="!user_list.length"
          message="暂无待分配好友"
        ></layout-empty>
      </div>
      <div
        class="page-view"
        v-if="meta.total > 10"
      >
        <el-pagination
          layout="prev, pager, next"
          :total="meta.total"
          :current-page="meta.page"
          @current-change="handleCurrentChange"
        >
        </el-pagination>
      </div>
    </div>
    <el-button
      v-if="user_list.length > 0"
      class="handel-btn"
      type="primary"
      size="small"
      @click="customerDistribution"
    >分配好友</el-button>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import Api from "@/api/wisdom/quit-inherit";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
@Component({
  components: {
    LayoutEmpty,
  },
})
export default class extends Vue {
  @Prop({ default: "" }) handoverUserid!: string;
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0,
  };
  private user_list = [];
  //内部绑定函数（获取当前第几页）
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.allotExternalList();
  }
  // 分配好友
  private customerDistribution() {
    this.$emit("customerDistribution");
  }
  private allotExternalList() {
    Api.allotExternalList({
      page: this.meta.page,
      size: this.meta.size,
      search_name: "",
      handover_userid: this.handoverUserid,
    }).then((res: any) => {
      if (res.code === 200) {
        this.user_list = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 生命周期
  mounted() {
    this.allotExternalList();
  }
}
</script>
<style lang="scss" scoped>
.not-customer-detail {
  text-align: right;
  .box-in {
    text-align: left;
    border: 1px solid #e4e7ed;
    padding: 20px;
    margin-bottom: 20px;
    .customer-ul {
      min-height: 300px;
      overflow: auto;
      .customer-li {
        margin-bottom: 10px;
        .avatar {
          width: 14px;
          height: 14px;
          transform: translateY(1px);
          margin-right: 5px;
        }
      }
    }

    .page-view {
      padding-top: 10px;
      text-align: center;
    }
  }
}
</style>
