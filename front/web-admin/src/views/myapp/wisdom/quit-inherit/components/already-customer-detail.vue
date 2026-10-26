<template>
  <!-- 已分配好友情况 -->
  <div class="already-detail">
    <div class="box-in">
      <div class="guide-left">
        <div class="th-item th-header">
          <div class="td">所属门店</div>
          <div class="td">导购姓名</div>
        </div>
        <div class="t-body">
          <div
            class="th-item"
            :class="[
              takeover_userid === item.takeover_userid ? 'guide-on' : ''
            ]"
            v-for="(item, index) in guide_list"
            :key="index"
            @click="changeGuide(item.takeover_userid)"
          >
            <div
              class="td"
              :title="item.takeover_department_name"
            >
              {{ item.takeover_department_name }}
            </div>
            <div
              class="td"
              :title="item.takeover_name"
            >
              {{ item.takeover_name }}
            </div>
          </div>
        </div>
        <layout-empty
          class="empty-com"
          v-if="!guide_list.length"
          :message="'暂无导购'"
        ></layout-empty>
        <div class="page-view">
          <el-pagination
            layout="prev, pager, next"
            :total="metal.total"
            :current-page="metal.page"
            @current-change="handleLeftChange"
          >
          </el-pagination>
        </div>
      </div>
      <div class="line"></div>
      <div class="customer-right">
        <el-table
          :data="customer_list"
          height="400px"
          stripe
        >
          <el-table-column
            prop="external_nickname"
            label="分配顾客好友"
            show-overflow-tooltip
          >
            <template slot-scope="scope">
              <img
                :src="scope.row.external_avatar"
                class="external-avatar"
              />
              <span>{{ scope.row.external_nickname }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="take_over_time"
            label="接替时间"
            width="190"
          >
            <template slot-scope="scope">
              <span>{{
                scope.row.take_over_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="external_status_name"
            label="客户接替状态"
            width="130"
          >
            <template slot-scope="scope">
              <span :class="[
                  scope.row.external_status_name === '接替完毕'
                    ? 'success-status'
                    : scope.row.external_status_name === '等待接替'
                    ? 'info-status'
                    : 'danger-status'
                ]">{{ scope.row.external_status_name }}</span>
            </template>
          </el-table-column>
          <div
            slot="empty"
            v-if="!customer_list.length"
          >
            <layout-empty :message="'暂无好友'"></layout-empty>
          </div>
        </el-table>

        <div class="page-view">
          <el-pagination
            layout="prev, pager, next"
            :total="meta.total"
            :current-page="meta.page"
            @current-change="handleCurrentChange"
          >
          </el-pagination>
        </div>
      </div>
    </div>
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
  private customer_list: Array<any> = []; // 右侧顾客列表
  private guide_list: Array<any> = []; // 导购列表
  private takeover_userid = 0; //接替人的user_id
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0,
  };
  // 左边分页
  private metal = {
    size: 10,
    page: 1,
    total: 0,
  };
  //内部绑定函数（获取当前第几页）
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.takeoverExternalList();
  }
  //内部绑定函数（获取当前第几页）
  private handleLeftChange(val: number) {
    this.metal.page = val;
    this.takeoverList();
  }
  // 已分配导购列表
  private async takeoverList() {
    await Api.takeoverList({
      page: this.metal.page,
      size: this.metal.size,
      handover_userid: this.handoverUserid,
    }).then((res: any) => {
      if (res.code === 200) {
        this.guide_list = res.data.list;
        if (
          this.guide_list &&
          this.guide_list.length &&
          !this.takeover_userid
        ) {
          this.takeover_userid = this.guide_list[0].takeover_userid;
        }
        this.metal.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 已分配外部联系人列表
  private async takeoverExternalList() {
    if (!this.takeover_userid) {
      return;
    }
    await Api.takeoverExternalList({
      page: this.meta.page,
      size: this.meta.size,
      takeover_userid: this.takeover_userid,
      handover_userid: this.handoverUserid,
    }).then((res: any) => {
      if (res.code === 200) {
        this.customer_list = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 切换导购
  private changeGuide(takeoverUserid: any) {
    this.takeover_userid = takeoverUserid;
    this.meta.page = 1;
    this.takeoverExternalList();
  }
  // 生命周期
  async mounted() {
    await this.takeoverList();
    await this.takeoverExternalList();
  }
}
</script>
<style lang="scss" scoped>
.already-detail {
  .box-in {
    display: flex;
    .guide-left {
      position: relative;
      width: 25%;
      .t-body {
        height: 350px;
        overflow: auto;
        background-color: #f7f7f7;
        padding-bottom: 10px;
      }
      .empty-com {
        position: absolute;
        top: 0;
        left: 0;
      }
      .th-item {
        display: flex;
        line-height: 36px;
        color: #606266;
        padding: 0 15px;
        &:hover {
          background-color: #ebebeb;
        }
        &.th-header {
          background: #ebebf1;
          font-weight: 600;
        }
        .td {
          text-align: center;
          width: 120px;
          cursor: pointer;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        &.guide-on {
          position: relative;
          .td {
            font-weight: 600;
            color: #000000;
          }
          &::before {
            content: "";
            display: block;
            width: 2px;
            height: 14px;
            position: absolute;
            background-color: #000000;
            top: 50%;
            transform: translateY(-50%);
            left: 0;
          }
        }
      }
    }
    .line {
      width: 2%;
    }
    .customer-right {
      width: 73%;
      .success-status {
        color: $color-success;
      }
      .danger-status {
        color: $color-danger;
      }
    }
    .page-view {
      padding-top: 10px;
      text-align: center;
    }
  }
  .external-avatar {
    width: 14px;
    height: 14px;
    margin-right: 10px;
    transform: translateY(2px);
  }
}
.t-body,
/deep/.el-table--scrollable-y .el-table__body-wrapper {
  /*修改滚动条样式*/
  &::-webkit-scrollbar {
    width: 5px;
    height: 10px;
    /**/
  }
  &::-webkit-scrollbar-track {
    background: rgb(239, 239, 239);
    border-radius: 2px;
  }
  &::-webkit-scrollbar-thumb {
    background: #bfbfbf;
    border-radius: 5px;
  }
  &::-webkit-scrollbar-thumb:hover {
    background: #333;
  }
  &::-webkit-scrollbar-corner {
    background: #179a16;
  }
}
</style>
