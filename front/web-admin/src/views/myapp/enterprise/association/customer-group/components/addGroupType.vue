<template>
  <el-dialog
    title="管理群类型"
    :close-on-click-modal="false"
    :visible.sync="dialogVisible"
    width="600px"
    @close="closeDialog"
  >
    <div>
      <div class="header-search">
        <div>
          <span>群类型名称：</span>
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入群类型名称"
            clearable
            v-model="typeName"
            @change="searchType"
          ></cy-input>
        </div>
        <el-button type="primary" @click="addNewType">新增群类型</el-button>
      </div>
      <div class="table m-t-30 m-b-30">
        <el-table :data="showTable" max-height="300" style="width:100%;">
          <el-table-column label="群类型名称" min-width="100">
            <template slot-scope="{ row }">
              <div
                v-if="row.status && row.status == 1"
                style="margin-bottom: 2px;"
              >
                <cy-input
                  size="small"
                  v-model="row.name"
                  placeholder="请输入群类型名称"
                  maxlength="10"
                  showWordLimit
                ></cy-input>
                <div class="color-danger">{{ row.error }}</div>
              </div>
              <div v-else>{{ row.name }}</div>
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template slot-scope="scope">
              <div v-if="scope.row.status && scope.row.status == 1">
                <el-button type="text" @click="saveChange(scope.row)"
                  >保存</el-button
                >
                <el-button type="text" @click="cancelSave(scope)"
                  >取消</el-button
                >
              </div>
              <div v-else>
                <el-button type="text" @click="edit(scope.row)">编辑</el-button>
                <el-popover
                  placement="top-start"
                  width="200"
                  trigger="hover"
                  content="该群类型下存在分组，不可删除"
                  :disabled="!(scope.row.g_num && scope.row.g_num > 0)"
                >
                  <el-button
                    type="text"
                    slot="reference"
                    :class="[
                      !(scope.row.g_num && scope.row.g_num > 0)
                        ? 'color-danger'
                        : 'color-disable',
                      'm-l-20'
                    ]"
                    @click="delItem(scope.row)"
                    >删除</el-button
                  >
                </el-popover>
              </div>
            </template>
          </el-table-column>
        </el-table>
        <div class="clearfix">
          <el-pagination
            class="float-right m-t-20"
            @current-change="handleCurrentChange"
            :current-page.sync="meta.page"
            :page-size="5"
            layout="total, prev, pager, next, jumper"
            :total="meta.total"
          >
          </el-pagination>
        </div>
      </div>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import { chatTypeList, chatTypeManage } from "@/api/enterprise/association";
@Component({})
export default class extends Vue {
  /************data**************/
  private dialogVisible = true;
  private typeName = "";
  private groupTypeList: any = [];
  private meta = {
    total: 0,
    page: 1
  };
  private tempList = [];
  /************methods**************/
  // 获取社群类型列表
  private getChatType() {
    chatTypeList().then(res => {
      console.log(res);
      res.data = res.data || [];
      this.groupTypeList = res.data;
      this.tempList = JSON.parse(JSON.stringify(res.data));
      this.meta.total = res.data.length;
    });
  }
  // 换页
  private handleCurrentChange(e: number) {
    this.meta.page = e;
  }
  // 新增类型
  private addNewType() {
    this.groupTypeList.unshift({
      status: "1",
      name: ""
    });
  }
  // 搜索
  private searchType(val: string) {
    console.log(val, "st");
    const temp = JSON.parse(JSON.stringify(this.tempList));
    if (val) {
      this.groupTypeList = temp.filter((v: any) => v.name.indexOf(val) != -1);
    } else {
      this.groupTypeList = temp;
    }
    this.meta.page = 1;
    this.meta.total = this.groupTypeList.length;
  }
  // 保存
  private saveChange(row: any) {
    console.log(row, "324");
    if (
      this.groupTypeList.filter((v: any) => v.name.trim() == row.name.trim())
        .length >= 2
    ) {
      row.error = "该名称已存在，请重新输入";
      this.groupTypeList.splice(0, 0); // 触发视图更新
      return;
    }
    if (row.id) {
      chatTypeManage({
        opt: 2,
        name: row.name,
        id: row.id
      }).then((res: any) => {
        if (res.code == 200) {
          this.getChatType();
          this.$message.success("编辑成功");
        } else {
          this.$message.error(res.message);
        }
      });
    } else {
      chatTypeManage({
        opt: 1,
        name: row.name
      }).then((res: any) => {
        if (res.code == 200) {
          this.getChatType();
          this.$message.success("新增成功");
        } else {
          this.$message.error(res.message);
        }
      });
    }
  }
  // 删除
  private delItem(row: any) {
    if (row.g_num && row.g_num > 0) return;
    chatTypeManage({
      opt: 3,
      id: row.id
    }).then((res: any) => {
      if (res.code == 200) {
        this.getChatType();
        this.$message.success("删除成功");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("close");
  }
  // 取消
  cancelSave(scope: any) {
    console.log(scope);
    if (scope.row.id) {
      this.getChatType();
    } else {
      this.groupTypeList.splice(scope.$index, 1);
    }
  }
  // 编辑
  edit(row: any) {
    console.log(row);
    row.status = "1";
    this.groupTypeList.splice(0, 0); // 触发视图更新
  }
  // 生命周期
  private mounted() {
    this.getChatType();
  }
  // 计算属性
  private get showTable() {
    return this.groupTypeList.slice(
      (this.meta.page - 1) * 5,
      this.meta.page * 5
    );
  }
}
</script>
<style lang="scss" scoped>
.header-search {
  display: flex;
  justify-content: space-between;
}
.m-l-20 {
  margin-left: 20px;
}
.float-right {
  float: right;
}
.color-disable {
  color: #909399;
}
</style>
