<template>
  <div>
    <el-dialog
      title="选择卡等"
      :visible="isShow"
      @close="closeDialogFn"
      :close-on-click-modal="false"
    >
      <div class="table">
        <div class="table-cell">
          <div class="left-box">名称</div>
          <div class="right-box">操作</div>
        </div>
        <div
          class="table-content-row"
          v-for="(item, index) in list"
          :key="index"
        >
          <div class="left-box">{{ item.name }}</div>
          <div class="right-box">
            <el-checkbox v-model="item.isSelect"></el-checkbox>
          </div>
        </div>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialogFn">取 消</el-button>
        <el-button size="small" type="primary" @click="submitFn"
          >保 存</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import API from "@/api/enterprise/drainage-coach";
@Component
export default class extends Vue {
  private isShow = false;
  private list: Array<any> = [];

  private mounted() {
    this.getList();
  }

  private async getList() {
    const res: any = await API.drainageGetDetail({});
    if (res.code === 200) {
      this.list = [
        {
          id: 0,
          name: "非会员",
          isSelect: false
        }
      ];
      const data = res.data || [];
      data.forEach((item: any) => {
        this.list.push({
          id: item.id,
          name: item.name,
          isSelect: false
        });
      });
      this.isShow = true;
    } else {
      this.$message.error(res.message);
    }
  }

  private closeDialogFn() {
    this.$emit("channelDialog");
  }

  private submitFn() {
    const arr = this.list.filter(item => item.isSelect);
    console.log(arr);
    if (arr.length > 0) {
      this.$emit("submitCard", arr);
    } else {
      this.$message.warning("请勾选一个卡等");
    }
  }
}
</script>
<style lang="scss" scoped>
/deep/ .layout-list {
  padding: 0;
}
.table {
  // border: 1px solid #ebebeb;
  .table-cell {
    display: flex;
    height: 42px;
    line-height: 42px;
    text-align: center;
    background-color: #d7d7d7;
    .left-box {
      flex: 1;
    }
    .right-box {
      flex: 1;
    }
  }
  .table-content-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    border-bottom: 1px solid #ebebeb;
    height: 38px;
    line-height: 38px;
    text-align: center;
    .right-box {
      border-left: 1px solid #ebebeb;
      border-right: 1px solid #ebebeb;
    }
    .left-box {
      border-left: 1px solid #ebebeb;
    }
  }
}
</style>
