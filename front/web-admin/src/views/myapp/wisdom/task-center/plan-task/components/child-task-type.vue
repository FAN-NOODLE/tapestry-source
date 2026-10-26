<template>
  <div class="select-materials-dialog">
    <el-dialog
      :visible.sync="dialogVisible"
      width="800px"
      @close="createTaskType"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="grid-content border-right">
            <div class="grid-header">
              <span>{{ep.isCoach ? '创建企微沟通任务':'创建沟通任务'}}</span>
            </div>
            <div class="grid-desp" style="height: 150px;margin-left: 10px;">
              <div class="grid-desp-title">导购通过任务发送消息至指定的客户聊天框</div>
              <div class="grid-desp-content">
                <div>此任务类型可创建：</div>
                <div v-if="ep.isCoach">1）一键转发任务素材给好友顾客的企微沟通任务(企微一键转发)</div>
                <div v-if="ep.isCoach">2）灵活筛选顾客转发&单独沟通的企微沟通任务（灵活沟通转发）</div>
                <div v-if="!ep.isCoach">1）一键转发任务素材给客户的沟通任务（企微一键群发）</div>
                <div v-if="!ep.isCoach">2）导购1v1沟通客户务&打电话的沟通任务（标签任务群发）</div>
              </div>
            </div>
            <div class="grid-footer">
              <el-button
                type="primary"
                @click="createTaskType(0, '沟通任务')"
              >立即创建</el-button>
            </div>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="grid-content">
            <div class="grid-header">
              <span>创建朋友圈任务</span>
            </div>
            <div class="grid-desp" style="margin-left: 30px;">
              <div class="grid-desp-title">导购通过任务发送消息至指定客户的朋友圈</div>
              <div class="grid-desp-content">
                <div>此任务类型可创建：</div>
                <div>1）一键转发任务素材到客户朋友圈的任务（企微一键群发）</div>
                <div>2）通过标签筛选客户可见范围转发朋友圈的任务（标签任务群发）</div>
              </div>
            </div>
            <div class="grid-footer">
              <el-button
                type="primary"
                @click="createTaskType(1, '朋友圈任务')"
              >立即创建</el-button>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import PerModule from "@/store/modules/permission";
@Component
export default class extends Vue {
  dialogVisible = true;
  get ep() {
    return PerModule.enterprise
  }

  // 创建
  createTaskType(type: number, typeName: string) {
    this.$emit("createTaskType", {
      type: type,
      typeName: typeName
    });
  }
}
</script>
<style lang="scss" scoped>
.select-materials-dialog {
  .border-right {
    border-right: 1px solid #000;
  }

  .grid-content {
    height: 250px;
    text-align: center;

    .grid-header {
      font-size: 16px;
      font-weight: bold;
    }

    .grid-desp {
      text-align: left;
      line-height: 30px;
      padding: 20px 0;
      color: #999;
      &-title{
        color: #AC750B;
      }
      &-content{
        font-size: 12px;
      }
    }
  }

  /deep/ .el-dialog {
    height: 350px;
  }
}
</style>
