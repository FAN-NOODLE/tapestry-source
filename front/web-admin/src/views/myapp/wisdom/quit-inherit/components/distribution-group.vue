<template>
  <!-- 分配社群 -->
  <el-dialog
    title="分配导购社群"
    :visible.sync="dialogVisible"
    width="60%"
    @close="cancel"
  >
    <el-steps
      :space="400"
      :active="steps_active"
      finish-status="success"
    >
      <el-step title="选择需要分配的社群"></el-step>
      <el-step title="选择新的群主导购"></el-step>
    </el-steps>
    <div
      v-if="steps_active === 0"
      class="step-first"
    >
      <div class="cont-box left-box">
        <div class="title">待分配社群列表</div>
        <div class="box-in">
          <div
            class="check-box"
            v-if="group_list.length"
          >
            <el-checkbox
              :indeterminate="isIndeterminate"
              v-model="checkAll"
              @change="handleCheckAllChange"
            >全选</el-checkbox>
            <div class="checkbox-group">
              <el-checkbox-group
                v-model="checked_data"
                @change="handleCheckedItemChange"
              >
                <div
                  v-for="(item, index) in group_list"
                  :key="index"
                  class="check-item"
                >
                  <el-checkbox :label="item.chat_id">
                    {{ item.name }}
                  </el-checkbox>
                </div>
              </el-checkbox-group>
            </div>
          </div>
          <layout-empty
            v-else
            :message="'暂无待分配社群'"
          ></layout-empty>
        </div>
      </div>
      <div class="center-box">
        <i class="el-icon-d-arrow-right"></i>
      </div>
      <div class="cont-box right-box">
        <div class="title">已选择社群列表</div>
        <div class="box-in">
          <ul class="right-ul">
            <li
              v-for="(item, index) in checked_group()"
              :key="index"
              class="sel-li"
            >
              <div class="group-name">
                {{ item.name }}
              </div>
              <i
                @click="delCheckedData(index)"
                class="el-icon-circle-close"
              ></i>
            </li>
          </ul>
        </div>
      </div>
    </div>
    <div
      v-else
      class="step-second"
    >
      <selectStaff @selChange="selStaffChange"></selectStaff>
    </div>
    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button
        size="small"
        @click="cancel"
      >取 消</el-button>
      <template v-if="steps_active === 0">
        <el-button
          @click="handelStep(1)"
          type="primary"
          size="small"
        >下一步</el-button>
      </template>
      <template v-else>
        <el-button
          size="small"
          @click="handelStep(0)"
        >上一步</el-button>
        <el-button
          size="small"
          type="primary"
          @click="transferCustomerGroup"
        >确 认</el-button>
      </template>
    </span>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import selectStaff from "./select-staff.vue";
import Api from "@/api/wisdom/quit-inherit";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
@Component({
  components: { selectStaff, LayoutEmpty },
})
export default class extends Vue {
  @Prop({ default: "" }) handoverUserid!: string;
  private dialogVisible = true;
  private steps_active = 0;
  private name_key = "";
  private checkAll = false;
  private checked_data: Array<any> = [];
  private group_list: Array<any> = [];
  private takeover_userid = "";
  //分页
  private meta = {
    size: 9999,
    page: 1,
    total: 0,
  };
  private isIndeterminate = false;
  // 右侧已勾选数据
  private checked_group() {
    const list = this.group_list.filter((item: any) => {
      return this.checked_data.indexOf(item.chat_id) > -1;
    });
    return list;
  }
  // 下一步/上一步
  private handelStep(step: number) {
    if (step === 1 && !this.checked_data.length) {
      this.$message.warning("请选择需要分配的社群");
      return;
    }
    this.steps_active = step;
  }

  // 全选
  private handleCheckAllChange(val: any) {
    if (val) {
      this.checked_data = this.group_list.map((item: any) => item.chat_id);
    } else {
      this.checked_data = [];
    }
    this.isIndeterminate = false;
  }
  // 勾选框勾选
  private handleCheckedItemChange(value: any) {
    this.checkAll = value.length === this.group_list.length;
    this.isIndeterminate =
      value.length > 0 && value.length < this.group_list.length;
  }
  private cancel() {
    this.$emit("cancel");
  }
  // 待分配群列表
  private allotGroupList() {
    Api.allotGroupList({
      page: this.meta.page,
      size: this.meta.size,
      search_name: "",
      handover_userid: this.handoverUserid,
    }).then((res: any) => {
      if (res.code === 200) {
        this.group_list = res.data.list || [];
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 确认提交
  transferCustomerGroup() {
    if (!this.takeover_userid) {
      this.$message.warning("请选择新的群主导购");
      return;
    }
    Api.transferCustomerGroup({
      chat_id_list: this.checked_data,
      takeover_userid: this.takeover_userid,
      handover_userid: this.handoverUserid,
    }).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("提交成功");
        this.cancel();
        this.$emit("submit");
      } else if (res.code === 40098) {
        this.$message.error("所选导购尚未完成实名认证");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 选择分配的导购
  private selStaffChange(val: any) {
    if (val && val.length) {
      this.takeover_userid = val[0].user_id;
    } else {
      this.takeover_userid = "";
    }
  }
  // 删除已选择的好友
  private delCheckedData(index: number) {
    this.checked_data.splice(index, 1);
    // 全选勾选状态
    this.checkAll = this.checked_data.length === this.group_list.length;
    this.isIndeterminate =
      this.checked_data.length > 0 &&
      this.checked_data.length < this.group_list.length;
  }
  // 生命周期
  mounted() {
    this.allotGroupList();
  }
}
</script>
<style lang="scss" scoped>
.step-first {
  display: flex;
  .cont-box {
    width: 400px;
    .title {
      padding: 10px 0;
    }
    .box-in {
      border: 1px solid #ccc;
      padding: 10px;

      .check-box {
        padding-top: 10px;
        .checkbox-group {
          height: 265px;
          overflow: auto;
          margin-top: 5px;
          .check-item {
            margin-bottom: 5px;
            .avatar {
              width: 14px;
              height: 14px;
              transform: translateY(1px);
              margin-right: 5px;
            }
          }
        }
      }
    }
  }
  .center-box {
    display: flex;
    padding: 0 15px;
    align-items: center;
    justify-content: center;
  }
  .right-box {
    .box-in {
      height: 300px;
      overflow: auto;
    }
    .right-ul {
      list-style: none;
      margin: 0;
      padding: 10px;
      .sel-li {
        margin-bottom: 10px;
        position: relative;
        .group-name {
          width: 300px;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }
        .el-icon-circle-close {
          position: absolute;
          right: 0;
          top: 50%;
          transform: translateY(-50%);
          cursor: pointer;
          font-size: 16px;
          &:hover {
            color: #000000;
          }
        }
      }
    }
  }
}
.step-second {
  padding-top: 10px;
}
.box-in,
.checkbox-group {
  .checkbox-group {
    height: 230px;
    overflow: auto;
    /*修改滚动条样式*/
    &::-webkit-scrollbar {
      width: 10px;
      height: 10px;
      /**/
    }
    &::-webkit-scrollbar-track {
      background: rgb(239, 239, 239);
      border-radius: 2px;
    }
    &::-webkit-scrollbar-thumb {
      background: #bfbfbf;
      border-radius: 10px;
    }
    &::-webkit-scrollbar-thumb:hover {
      background: #333;
    }
    &::-webkit-scrollbar-corner {
      background: #179a16;
    }
  }
}
</style>
