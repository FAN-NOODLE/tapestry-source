<template>
  <div>
    <el-dialog
      title="员工删除客户提醒通知配置"
      :visible="true"
      @close="closeDialog"
      :close-on-click-modal="false"
    >
      <div class="dialog-body">
        <p>
          选择当前部门人员在删除客户时，自动通知更多人在企业微信收到一条消息提醒
        </p>
        <el-button
          size="small"
          class="select-btn"
          type="primary"
          @click="openSelectStaff"
        >
          选择员工
        </el-button>
        <el-table stripe :data="list" class="mar-t-20">
          <el-table-column prop="dp_name" label="部门"></el-table-column>
          <el-table-column prop="member_name" label="员工"></el-table-column>
          <el-table-column label="员工删除客户">
            <template slot-scope="scope">
              <el-checkbox
                v-model="scope.row.select1"
                @change="change1(scope.$index)"
              ></el-checkbox>
            </template>
          </el-table-column>
          <el-table-column label="客户删除员工">
            <template slot-scope="scope">
              <el-checkbox
                v-model="scope.row.select2"
                @change="change2(scope.$index)"
              ></el-checkbox>
            </template>
          </el-table-column>
          <el-table-column prop="zip" label="操作">
            <template slot-scope="scope">
              <span @click="deleteRow(scope.$index)" class="c-blue">删除 </span>
            </template>
          </el-table-column>
          <div slot="empty" class="table-empty-label">
            <img src="@/assets/images/empty-box.png" alt="" />
            <p>暂无数据</p>
          </div>
        </el-table>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button type="primary" @click="confirm">确 定</el-button>
      </div>
    </el-dialog>
    
    <!-- 选择员工弹窗-->
    <selectMember
      v-if="showSelectStaff"
      :hasSelectArr="staff"
      :limiteNum="100"
      dialogTitle="选择员工"
      @submit="getMember"
      @close="showSelectStaff = false"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import SelectMember from "@/components/business/select-member/index.vue";
import Api from "@/api/enterprise/enterprise";
@Component({
  components: {
    SelectMember
  }
})
export default class extends Vue {
  @Prop() selectdpId!: number;
  /*** data ***/
  private showSelectStaff = false;
  private selectListId: Array<string | number> = [];
  private list: Array<any> = [];
  private staff: Array<any> = [];
  // 生命周期
  private mounted() {
    this.getList();
  }
  // 选人组件回调
  private getMember(memberList: Array<any>) {
    // 先判断是否有通过组件删除的数据
    for (let i = this.list.length - 1; i >= 0; i--) {
      const result_one = memberList.find((item: any) => {
        return item.user_id === this.list[i].user_id;
      });
      if (!result_one) {
        this.list.splice(i, 1);
      }
    }
    // 将memberList中list已有的数据删除
    for (let i = memberList.length - 1; i >= 0; i--) {
      const result_one = this.list.find((item: any) => {
        return item.user_id === memberList[i].user_id;
      });
      if (result_one) {
        memberList.splice(i, 1);
      }
    }
    const selectList: Array<any> = [];
    memberList.forEach((item: any) => {
      selectList.push({
        dp_name: item.dept_name,
        user_id: item.user_id,
        member_name: item.user_name,
        select1: true,
        select2: true
      });
    });
    this.list = this.list.concat(selectList);
    this.showSelectStaff = false;
  }

  private getList() {
    const params = {
      dp_id: this.selectdpId
    };

    Api.departmentsUserLostConfig(params).then((res: any) => {
      if (res.code === 200) {
        this.list = res.data;
        if (this.list && this.list.length > 0) {
          this.list.forEach((item: any) => {
            const send_type = item.send_type;
            if (send_type === 0) {
              item.select1 = true;
              item.select2 = true;
            } else if (send_type === 1) {
              item.select1 = false;
              item.select2 = true;
            } else if (send_type === 2) {
              item.select1 = true;
              item.select2 = false;
            }
          });
        }
        console.log("this.list", this.list);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  private openSelectStaff() {
    this.staff = [];
    this.list.forEach((item: any) => {
      if (item.user_id) {
        this.staff.push({
          name: item.member_name,
          user_id: item.user_id,
          dp_name: item.dp_name
        });
      }
    });
    this.showSelectStaff = true;
  }

  private deleteRow(index: number) {
    this.list.splice(index, 1);
    console.log("this.list", this.list);
  }

  private closeDialog() {
    this.$emit("closeDialog");
  }

  private confirm() {
    const userList: any = [];
    let isError = false;
    this.list.forEach((item: any) => {
      let send_type;
      if (item.select1 === true && item.select2 === true) {
        send_type = 0;
      } else if (item.select1 === false && item.select2 === true) {
        send_type = 1;
      } else if (item.select1 === true && item.select2 === false) {
        send_type = 2;
      } else if (item.select1 === false && item.select2 === false) {
        isError = true;
      }
      userList.push({
        user_id: item.user_id,
        send_type: send_type
      });
    });
    if (isError) {
      this.$message.warning("每一行至少勾选一个");
      return;
    }
    const params = {
      dp_id: this.selectdpId,
      user_ids: userList
    };
    Api.userLostConfig(params)
      .then(() => {
        this.closeDialog();
      })
      .catch((err: any) => {
        this.$message.error(err.message);
      });
  }
  private change1(index: number) {
    this.list[index].select1 = !this.list[index].select1;
    this.$set(this.list, index, this.list[index]);
  }

  private change2(index: number) {
    this.list[index].select2 = !this.list[index].select2;
    this.$set(this.list, index, this.list[index]);
  }
}
</script>

<style lang="scss" scoped>
.select-btn {
  margin: 20px 0;
}
.c-blue {
  color: #708df5;
  font-size: 12px;
  cursor: pointer;
}
/deep/ .el-table__body-wrapper {
  max-height: 300px !important;
  overflow-y: scroll;
}
/deep/ .table-empty-label img {
  position: relative;
  top: 40px;
}
</style>
