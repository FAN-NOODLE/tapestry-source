<template>
  <div class="applets-materials">
    <div class="title-btns">
      <div>
        <el-checkbox
          :indeterminate="isIndeterminate"
          v-model="checkAll"
          @change="handleCheckAllChange"
          class="all-check"
          >全选</el-checkbox
        >
        <span class="selected" v-if="checkedMaterials.length"
          >已选：{{ checkedMaterials.length }}</span
        >
        <el-button class="move-group" @click="moveGroupFn" size="small"
          >移动分组</el-button
        >
        <el-button
          class="remove"
          @click="deletesFn(0, '批量删除素材')"
          size="small"
          >批量删除</el-button
        >
        <el-button class="remove" @click="lotUpdateStatus(0)" size="small"
          >批量启用</el-button
        >
        <el-button class="remove" @click="lotUpdateStatus(1)" size="small"
          >批量禁用</el-button
        >
      </div>
      <!-- v-if="filterObj.children == null" -->
      <el-button
        size="small"
        type="primary"
        class="onload-pic"
        @click="addAppletsFn"
        v-if="filterObj.source_group_id != -2"
        >新增小程序</el-button
      >
    </div>
    <div class="content clearfix" v-if="appletsList.length">
      <el-checkbox-group
        v-model="checkedMaterials"
        @change="handleCheckedMaterialsChange"
      >
        <div class="card" v-for="(item, index) in appletsList" :key="index">
          <h3 class="applets-title">
            <span
              ><img
                src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                alt=""/></span
            >{{ item.app_name }}
          </h3>
          <p class="card-title">
            {{ item.title }}
          </p>
          <img class="img" :src="item.cover_url" alt="" />
          <div class="check-box">
            <el-checkbox :label="item.id" :key="index"><br /></el-checkbox>
            <span>{{
              item.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span>
          </div>
          <div class="btns">
            <el-button
              type="text"
              size="small"
              :style="{ color: item.source_status === 0 ? '#F56C6C' : '' }"
              @click="updateStatus(item)"
              >{{ item.source_status === 0 ? "禁用" : "启用" }}</el-button
            >
            <el-button size="small" type="text" @click="editFn(item.id)"
              >编辑</el-button
            >
            <el-button
              size="small"
              type="text"
              @click="deleteItem(item.id)"
              class="del-btn"
              >删除</el-button
            >
            <el-button
              v-if="ep.isCoach && filterObj.source_group_id !== -2"
              type="text"
              size="small"
              :disabled="item.is_top === 0 && !item.can_top"
              @click="handleTop(item)"
              >{{ item.is_top === 0 ? "置顶" : "取消置顶" }}</el-button
            >
          </div>
        </div>
      </el-checkbox-group>
    </div>
    <div v-else class="empty-show">
      <div
        v-if="
          filterObj.keyword ||
            filterObj.create_end != 0 ||
            filterObj.create_start != 0 ||
            filterObj.just_me != 0
        "
      >
        <img src="@/assets/images/empty-search.png" alt="" />
        <p>
          没有找到相关小程序素材
        </p>
      </div>
      <div v-else>
        <img src="@/assets/images/empty-box.png" alt="" />
        <p>暂无素材</p>
      </div>
    </div>
    <div class="paging" v-if="appletsList.length">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page.sync="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      >
      </el-pagination>
    </div>
    <movegroup
      v-if="moveGtoupacom"
      :multipleSelection="checkedMaterials"
      @closeDialog="closemoveFn"
    ></movegroup>
    <small-mterial-dialog
      v-if="appletsDialog"
      @submit="closeAppletsDialogFn"
      :groupId="filterObj.source_group_id"
      :groupName="filterObj.name"
      :appletsDialog="appletsDialog"
      :editObj="editObj"
    ></small-mterial-dialog>
    <el-dialog title="已过期素材" :visible.sync="overdueDialog">
      <el-alert
        title="以下素材有效期已过，无法操作批量启用，请重新设置素材有效期后，启用素材"
        type="warning"
        :closable="false"
      >
      </el-alert>
      <el-table style="margin-top:15px" height="450" stripe :data="overdueData">
        <el-table-column
          show-overflow-tooltip
          property="app_name"
          label="小程序名称"
        ></el-table-column>
        <el-table-column label="图片" min-width="80">
          <template slot-scope="scope">
            <img
              style="width: 32px;height: 32px;object-fit: cover;"
              :src="scope.row.cover_url"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import movegroup from "../components/move-group.vue";
import API from "@/api/enterprise/material-library";
import smallMterialDialog from "../components/small-mterial-dialog.vue"; //小程序素材库
import PerModule from "@/store/modules/permission";
@Component({
  components: {
    movegroup,
    smallMterialDialog
  }
})
export default class extends Vue {
  @Prop() filterObj!: any; // 筛选条件
  /*************************** data start ***************************/
  private checkAll = false; //复选框全选
  private isIndeterminate = false; //复选框全选
  private checkedMaterials: any = []; //选中的复选框
  private multipleSelectionData: Array<any> = [];
  private overdueData: Array<any> = []; // 已过期的素材
  private overdueDialog = false; // 已过期的素材的弹窗
  private appletsList = []; //小程序列表

  private moveGtoupacom = false; //移动分组
  private appletsDialog = false; //小程序
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private multipleSelection = [];
  private editObj = {}; //编辑id等信息

  get ep() {
    return PerModule.enterprise
  }

  /*************************** method start ***************************/
  //全选的方法
  private handleCheckAllChange(val: any) {
    const arr: any = [];
    this.appletsList.forEach((item: any) => {
      arr.push(item.id);
    });
    this.checkedMaterials = val ? arr : [];
    this.multipleSelectionData = val ? this.appletsList : [];
    this.isIndeterminate = false;
  }
  //单选的方法
  private handleCheckedMaterialsChange(value: any) {
    const checkedCount = value.length;
    this.checkAll = checkedCount === this.appletsList.length;
    this.isIndeterminate =
      checkedCount > 0 && checkedCount < this.appletsList.length;
    this.multipleSelectionData = this.appletsList.filter((item: any) => {
      return value.indexOf(item.id) > -1;
    });
  }
  //新增小程序
  private addAppletsFn() {
    this.editObj = { type: "create", id: this.filterObj.source_group_id };
    this.appletsDialog = true;
  }
  // 编辑小程序
  private editFn(id: number) {
    this.editObj = { type: "edit", id: id };
    this.appletsDialog = true;
  }
  // 小程序素材弹框关闭
  private closeAppletsDialogFn(obj?: any) {
    if (obj && obj.save) {
      this.getAppletsListFn();
    }
    this.appletsDialog = false;
    // 判断是否从新增子任务页面跳转进来
    // 清除地址栏携带参数
    if (this.$route.query.isAddTask) {
      this.$router.push({ query: {} });
    }
  }
  private closemoveFn(val?: boolean) {
    if (val) {
      this.getAppletsListFn();
    }
    this.moveGtoupacom = false;
  }
  //移动分组
  private moveGroupFn() {
    if (!this.checkedMaterials.length) {
      this.$message.warning("请选择要移动分组的素材");
      return false;
    }
    this.moveGtoupacom = true;
  }
  //删除
  private deletesFn(id?: number, title?: any) {
    if (!id && !this.checkedMaterials.length) {
      this.$message.warning("请选择要删除的素材");
      return false;
    }
    this.$confirm("你确认要删除选择的素材吗？", title, {
      confirmButtonText: "确定删除",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        if (this.filterObj.source_group_id === -2) {
          this.getDeleteFn(id);
        } else {
          this.moveRecycleBin(id);
        }
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  // 批量移动到回收站
  private moveRecycleBin(id?: number) {
    let ids: any = [];
    if (id) {
      ids.push(id);
    } else {
      ids = this.checkedMaterials;
    }
    API.moveGroupApi({
      group_id: -2,
      source_ids: ids
    }).then((res: any) => {
      if (res.code == 200) {
        this.getAppletsListFn();
        this.$message.success("素材已移动到回收站");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //删除
  private deleteItem(id?: number) {
    this.$confirm("确定删除该素材吗?", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    }).then(() => {
      if (this.filterObj.source_group_id === -2) {
        this.getDeleteFn(id);
      } else {
        this.moveRecycleBin(id);
      }
    });
  }
  private getDeleteFn(id?: number) {
    let ids: any = [];
    if (id) {
      ids.push(id);
    } else {
      ids = this.checkedMaterials;
    }
    API.materialDeleteApi({ ids: ids }).then((res: any) => {
      if (res.code == 200) {
        this.getAppletsListFn();
        this.$message({
          type: "success",
          message: "删除成功!"
        });
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 分页
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getAppletsListFn();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getAppletsListFn();
  }
  //获取列表
  private getAppletsListFn(isPage?: boolean) {
    if (isPage) {
      this.meta.page = 1;
      this.meta.size = 10;
    }
    const params = {
      page: this.meta.page, //页码
      size: this.meta.size, //每页数量
      type: 4
    };
    Object.assign(params, this.filterObj);
    API.consoleList(params).then((res: any) => {
      if (res.code == 200) {
        this.resetCheckFn();
        const data = res.data;
        this.appletsList = data.list || [];
        this.meta.total = data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //重置多选
  private resetCheckFn() {
    this.checkedMaterials = [];
    this.checkAll = false;
  }
  // 更新状态
  private updateStatus(row: any) {
    const nowTime = new Date().getTime();
    const endTime = row.end_time * 1000;
    let title = "";
    let tipCont = "";
    let isApi = true;
    if (row.source_status === 0) {
      title = "确定禁用该素材吗？";
      tipCont =
        "禁用后有效期内该素材将不在小程序端展示，后台其他模块也将无法调用";
    } else {
      if (nowTime >= endTime) {
        title = "已过期素材";
        tipCont = "该素材有效期已过，请重新设置后再启用";
        isApi = false;
      } else {
        title = "确定启用该素材吗？";
        tipCont =
          "启用后有效期内该素材将在小程序端展示，并支持后台其他模块调用";
      }
    }

    this.$confirm(tipCont, title, {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    }).then(() => {
      if (isApi) {
        this.updateStatusApi([row.id], row.source_status === 0 ? 1 : 0);
      } else {
        // 已过期提示点击确定打开编辑弹窗
        this.editFn(row.id);
      }
    });
  }
  private updateStatusApi(ids: any, sourceStatus: any) {
    API.updateStatus({
      ids: ids,
      source_type: 4,
      source_status: sourceStatus
    }).then((res: any) => {
      if (res.code === 200) {
        if (sourceStatus === 0) {
          this.$message.success("已启用");
        } else {
          this.$message.success("已禁用");
        }
        this.getAppletsListFn();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  /**
   * 置顶操作
   * @param row 
   */
  async handleTop(row) {
    
    const message = row.is_top === 1 ?`确定取消该素材的置顶效果吗？`: '确定置顶该素材吗？';
    await this.$confirm(message, '提示', {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })

    const res:any = await API.setTopApi({
      source_id: row.id,
      source_type: row.type,
      is_top: row.is_top === 1 ? 2 : 1
    })

    if (res.code === 200) {
      this.getAppletsListFn();
    } else {
      this.$message.error(res.message);
    }

  }
  // 批量更新状态
  private lotUpdateStatus(sourceStatus: any) {
    if (this.multipleSelectionData.length == 0) {
      this.$message.warning("请选择要操作的素材");
      return;
    }
    // 批量设置启用状态，判断选中的素材中有没有已过期的素材
    if (sourceStatus === 0) {
      const nowTime = new Date().getTime();
      this.overdueData = this.multipleSelectionData.filter(
        (item: any) => nowTime >= item.end_time * 1000
      );
      if (this.overdueData.length) {
        this.overdueDialog = true;
      } else {
        this.lotUpdateStatusTip(this.checkedMaterials, sourceStatus);
      }
    } else {
      this.lotUpdateStatusTip(this.checkedMaterials, sourceStatus);
    }
  }
  private lotUpdateStatusTip(ids: any, sourceStatus: any) {
    let title = "";
    let tipCont = "";
    if (sourceStatus === 1) {
      title = "确定禁用这些素材吗？";
      tipCont =
        "禁用后有效期内这些素材将不在小程序端展示，后台其他模块也将无法调用";
    } else {
      title = "确定启用这些素材吗？";
      tipCont =
        "启用后有效期内这些素材将在小程序端展示，并支持后台其他模块调用";
    }
    this.$confirm(tipCont, title, {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    }).then(() => {
      this.updateStatusApi(ids, sourceStatus);
    });
  }
  /*************************** 生命周期 start ***************************/
  mounted() {
    this.getAppletsListFn();
    // 判断是否从新增子任务页面跳转进来
    if (this.$route.query.isAddTask) {
      this.addAppletsFn();
    }
  }
}
</script>
<style lang="scss" scoped>
.title-btns {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  .all-check,
  .selected {
    margin-right: 20px;
    color: $color-text-main;
  }
}
.content {
  margin-right: -20px;
  background: #ffffff;
  .card {
    float: left;
    width: 155px;
    padding: 0 10px;
    border-radius: 4px;
    border: 1px solid $color-second-line;
    margin-right: 20px;
    margin-bottom: 20px;
    .applets-title {
      line-height: 20px;
      color: #707277;
      margin: 10px 0 5px 0;
      font-size: 12px;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      span {
        float: left;
        width: 20px;
        height: 20px;
        background: #d8d8d8;
        border-radius: 100%;
        overflow: hidden;
        margin-right: 5px;
        img {
          width: 100%;
          height: 100%;
        }
      }
    }
    .card-title {
      line-height: 20px;
      color: $color-text-main;
      font-size: 14px;
      min-height: 40px;
      text-overflow: -o-ellipsis-lastline;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    .img {
      width: 100%;
      height: 124px;
      border-radius: 4px;
      margin: 10px 0;
      vertical-align: top;
      object-fit: cover;
    }
    .check-box {
      overflow: hidden;
      margin-bottom: 10px;
      span {
        color: red;
        float: right;
        font-size: 12px;
        color: #707277;
      }
    }
    .btns {
      border-top: 1px solid $color-second-line;
      text-align: right;
      margin: 2px 0;
    }
  }
  .del-btn {
    margin-left: 20px;
    color: #303133 !important;
  }
}
.paging {
  text-align: right;
}
.empty-show {
  height: 360px;
  margin-top: 120px;
  text-align: center;
  p {
    color: #909399;
    line-height: 60px;
  }
}
</style>
