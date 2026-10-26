<template>
  <div class="counter">
    <layout-form>
      <template slot="content">
        <el-dialog
          :close-on-click-modal="false"
          :title="titleType == 1 ? '被删除提醒消息配置' : '删除提醒消息配置'"
          width="1000px"
          hight="620px"
          :visible.sync="visibleDialog"
          @close="closeEditappDialog()"
        >
          <div class="module-box">
            <div class="module-left">
              <div class="tit">{{ titleType == 1 ? "被" : "" }}删除提醒</div>
              <div class="tit-mar">更多通知人</div>
            </div>
            <div class="module-right">
              <div class="row">
                <el-switch
                  v-model="switchValue"
                  :active-text="switchValue ? '已开启' : '已关闭'"
                >
                </el-switch>
              </div>
              <div class="desc" style="margin-bottom: 36px;">
                {{ titleType === 1 ? '开启后，当删除客户微信时，删除的企业员工将在企业微信收到一条消息提醒' : '开启后，当删除客户微信时，删除的企业员工将在企业微信收到一条消息提醒' }}
              </div>
              <div class="row" v-if="tableData[0].stat_txt_obj.length == 0">
                <el-button type="primary" size="small" @click="selectStaffbut"
                  >选择员工</el-button
                >
              </div>
              <div class="desc">
                当删除客户微信时，选择的更多通知人将在企业微信收到一条消息提醒
              </div>
            </div>
          </div>
          <div class="table-box" v-if="tableData[0].stat_txt_obj.length > 0">
            <!-- stripe  是否为斑马纹 table -->
            <el-table :data="tableData" style="width: 100%">
              <el-table-column label="员工" prop="stat_txt_obj">
                <template slot-scope="{ row }">
                  <el-tag
                    v-for="(item, index) in row.stat_txt_obj"
                    :key="index"
                    class="el-staff-tag"
                  >
                    <img
                      src="~@/assets/images/staff-icon.png"
                      alt=""
                      class="table-text-img"
                    />
                    <span>{{ item.name }}</span>
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="110" fixed="right">
                <template>
                  <el-button
                    type="text"
                    size="small"
                    @click.native.prevent="editStaff"
                  >
                    修改
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
          <span slot="footer" class="dialog-footer">
            <el-button @click="cancel()">取 消</el-button>
            <el-button type="primary" @click="save()">确认</el-button>
          </span>
        </el-dialog>
        <selectMember
          v-if="selectStaff"
          :hasSelectArr="hasSelectArr"
          @close="close"
          @submit="getMember"
        />
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import selectMember from "@/components/business/select-member/index.vue";
import Api from "@/api/enterprise/enterprise";
@Component({
  components: {
    selectMember
  }
})
export default class extends Vue {
  @Prop() titleType!: number;
  private visibleDialog = true;
  private switchValue = true;
  private selectStaff = false;
  private listId = 0;
  private tableData: Array<any> = [
    {
      stat_txt_obj: []
    }
  ];
  // 发送给选人弹框
  private hasSelectArr = {};
  private stat_txt = {
    member: []
  };
  private qy_stat_txt = {
    member: []
  };
  // 取消，关闭弹框
  cancel() {
    this.$emit("close", false);
  }
  save() {
    this.updateMessage();
  }
  // 关闭
  closeEditappDialog() {
    this.$emit("close", false);
  }
  selectStaffbut() {
    this.hasSelectArr = [];
    this.selectStaff = true;
  }
  editStaff() {
    const arr = this.tableData[0].stat_txt_obj;
    arr.map((item: any) => {
      item.is_dp = false;
    });
    this.hasSelectArr = arr;
    this.selectStaff = true;
  }
  private getMember(list: any) {
    this.tableData[0].stat_txt_obj = list;
    this.selectStaff = false;
  }
  // 关闭选人弹框
  close() {
    this.selectStaff = false;
  }
  // 获取员工列表
  getStaffList() {
    Api.messageListApi({ send_type: this.titleType }).then((res: any) => {
      if (res.code == 200) {
        if (Object.keys(res.data).length > 0) {
          this.tableData[0].stat_txt_obj = res.data.stat_txt_obj.member;
          this.listId = res.data.id;
          this.switchValue = res.data.status == 1 ? true : false;
        } else {
          this.switchValue = res.data.status = true;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击确定更新
  updateMessage() {
    const tData = this.tableData[0].stat_txt_obj;
    tData.map((item: any) => {
      const obj = {
        name: item.name,
        user_id: item.user_id
      };
      (this.stat_txt.member as any).push(obj);
    });
    const arr: Array<any> = [];
    tData.map((item: any) => {
      arr.push(item.user_id);
    });
    (this.qy_stat_txt.member as any) = arr;
    const params = {
      send_type: this.titleType,
      id: this.listId,
      status: this.switchValue ? 1 : 2,
      stat_txt: JSON.stringify(this.stat_txt),
      qy_stat_txt: JSON.stringify(this.qy_stat_txt)
    };
    Api.updateMessageApi(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("配置成功");
        this.$emit("close", false);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private mounted() {
    this.getStaffList();
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-table__row td {
  text-align: left !important;
}
/deep/ tr .cell {
  text-align: left !important;
}
/deep/ tr td .cell {
  height: 84px !important;
}
/deep/ tr td:first-child .cell {
  overflow-y: scroll !important;
}
// /deep/ tr td:last-child .cell {
//   line-height: 80px !important;
// }
.module-box {
  display: flex;
  margin-top: 10px;
  .row {
    font-size: 14px;
    margin-bottom: 20px;
  }
  .module-left {
    text-align: right;
    margin-right: 20px;
    .tit-mar {
      margin-top: 78px;
    }
  }
  .module-right {
    flex: 1;
    .desc {
      color: #909399;
    }
  }
}
.table-box {
  /deep/.el-table__row {
    cursor: pointer;
  }
  margin-left: 88px;
  margin-top: 30px;
  .el-staff-tag {
    margin-right: 10px;
    margin-bottom: 10px;
  }
  .table-text-img {
    display: inline-block;
    width: 14px;
    height: 14px;
    vertical-align: middle;
    margin-right: 5px;
  }
}
</style>
