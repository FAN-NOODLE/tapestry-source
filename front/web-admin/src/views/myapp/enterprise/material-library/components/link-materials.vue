<template>
  <div class="text-materials">
    <div class="title-btns">
      <div>
        <el-checkbox
          v-model="checkAll"
          @change="handleCheckAllChange"
          class="all-check"
          >全选</el-checkbox
        >
        <span class="selected" v-if="checkedMaterials.length"
          >已选：{{ checkedMaterials.length }}</span
        >
        <el-button class="move-group" @click="openMoveGrop" size="small"
          >移动分组</el-button
        >
        <el-button class="remove" @click="deleteFn" size="small"
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
        v-if="filterObj.source_group_id != -2"
        type="primary"
        class="onload-pic"
        @click="addBtn"
        size="small"
        >新增链接</el-button
      >
    </div>
    <div class="content clearfix">
      <div v-if="linkList.length == 0" class="empty-content">
        <img
          class="empty-img"
          v-if="!emptyFlag"
          src="@/assets/images/empty-box.png"
          alt=""
        />
        <div class="empty-box" v-if="!emptyFlag">
          <div class="desc">暂无素材</div>
        </div>
        <img
          class="empty-img"
          v-if="emptyFlag"
          src="@/assets/images/empty-search.png"
          alt=""
        />
        <div class="empty-box" v-if="emptyFlag">
          <div class="desc">没有找到相关链接素材</div>
        </div>
      </div>
      <div v-else>

        <el-checkbox-group
          v-model="checkedMaterials"
          @change="handleCheckedMaterialsChange"
        >
          <link-material-item
            v-for="(item, index) in linkList"
            :key="item.id"
            :item="{...item}"
            :index="index"
            :filter-obj="filterObj"
            @detail="detailItem"
            @edit="editItem"
            @delete="delItem"
            @update-status="updateStatus"
            @top="handleTop"
          ></link-material-item>
        </el-checkbox-group>
      </div>
    </div>
    <div class="footer" v-if="linkList.length > 0">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :page-sizes="[10, 20, 30, 40]"
        :current-page="meta.page"
        :page-size="meta.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      >
      </el-pagination>
    </div>
    <!-- 新增链接弹框 -->
    <add-link
      v-if="isShow"
      :groupId="filterObj.source_group_id"
      :groupName="filterObj.name"
      :value="value"
      :title="title"
      @closeDialog="closeAdd"
    ></add-link>
    <!-- 移动分组弹框 -->
    <movegroup
      v-if="moveGtoupacom"
      @closeDialog="closemove"
      :multipleSelection="this.checkedMaterials"
    ></movegroup>
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
          property="title"
          label="标题"
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
import addLink from "../components/add-link-material.vue";
import linkMaterialItem from "./link-material-item.vue";
import Api from "@/api/enterprise/material-library";
import { isNull } from "@/utils/tools";
import PerModule from "@/store/modules/permission";

@Component({
  components: {
    movegroup,
    addLink,
    linkMaterialItem
  }
})
export default class extends Vue {
  @Prop() filterObj!: any; // 筛选条件
  //------------------------------------------------------data-----------------------------------------------------------/
  private checkAll = false; //复选框全选
  private checkedMaterials: any = []; //选中的复选框
  private multipleSelectionData: Array<any> = [];
  private overdueData: Array<any> = []; // 已过期的素材
  private overdueDialog = false; // 已过期的素材的弹窗
  // 搜索区域
  private search = {
    keyWord: ""
  };
  // -----------------------
  private moveGtoupacom = false; //移动分组弹框是否展示
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 10
  };
  private emptyFlag = false; //是否为空素材
  private linkList: any = []; //素材列表
  private isShow = false; //新增弹窗是否展示
  private title = "新增链接";
  private value = {}; //编辑的素材内容
  private multipleSelection = [];

  get ep() {
    return PerModule.enterprise
  }
  //----------------------------------------------------methods-----------------------------------------------------------/
  // 点击移动分组
  openMoveGrop() {
    console.log(this.checkedMaterials);
    
    if (this.checkedMaterials.length == 0) {
      this.$message({
        message: "请选择要移动分组的素材",
        type: "warning"
      });
      return;
    }
    this.moveGtoupacom = true;
  }
  // 关闭移动分组弹窗
  private closemove(save: boolean) {
    this.moveGtoupacom = false;
    if (save == true) {
      this.getlist();
    }
  }
  //全选的方法
  private handleCheckAllChange(val: any) {
    const arr: any = [];
    this.linkList.forEach((item: any) => {
      arr.push(item.id);
    });
    this.checkedMaterials = val ? arr : [];
    this.multipleSelectionData = val ? this.linkList : [];
  }
  //单选的方法
  private handleCheckedMaterialsChange(value: any) {
    const checkedCount = value.length;
    this.checkAll = checkedCount === this.linkList.length;
    this.multipleSelectionData = this.linkList.filter((item: any) => {
      return value.indexOf(item.id) > -1;
    });
  }
  // 编辑每一项
  private editItem(value: any) {
    // this.isShow = true;
    // this.value = value;
    // this.title = "编辑链接";

    this.$router.push({ 
      name: "myapp_enterprise_material-library_create-link", 
      query: { 
        id: value.id|| undefined,
        groupId: this.filterObj.source_group_id,
      } 
    });
  }
  // 详情
  private detailItem(url: string) {
    window.open(url, "_blank");
  }
  // 删除每一项
  private delItem(id: number) {
    this.$confirm("确定删除该素材吗?", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    }).then(() => {
      if (this.filterObj.source_group_id === -2) {
        Api.materialDeleteApi({
          ids: [id]
        }).then((res: any) => {
          if (res.code == 200) {
            this.$message.success("删除成功");
            this.getlist();
          } else {
            this.$message.error(res.message);
          }
        });
      } else {
        this.moveRecycleBin([id]);
      }
    });
  }
  // 点击批量删除
  private deleteFn() {
    if (this.checkedMaterials.length == 0) {
      this.$message({
        message: "请选择要删除的素材",
        type: "warning"
      });
      return;
    }
    this.$confirm("你确认要删除选择的素材吗？", "批量删除素材", {
      confirmButtonText: "确定删除",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        if (this.filterObj.source_group_id === -2) {
          this.getDeleteFn(this.checkedMaterials);
        } else {
          this.moveRecycleBin(this.checkedMaterials);
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
  private moveRecycleBin(ids: any) {
    Api.moveGroupApi({
      group_id: -2,
      source_ids: ids
    }).then((res: any) => {
      if (res.code == 200) {
        this.getlist();
        this.$message.success("素材已移动到回收站");
      } else {
        this.$message.error(res.message);
      }
    });
  }

  //批量删除
  private getDeleteFn(ids: any) {
    Api.materialDeleteApi({ ids: ids }).then((res: any) => {
      if (res.code == 200) {
        this.$message({
          type: "success",
          message: "删除成功!"
        });
        this.getlist();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.size = val;
    this.meta.page = 1;
    this.getlist();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getlist();
  }
  // ------------------------------方法------------------------------
  // 列表数据
  private getlist(page?: number) {
    if (page) {
      this.meta.page = page;
      this.meta.size = 10;
    }
    const params = {
      page: this.meta.page, //页码
      size: this.meta.size, //每页数量
      type: 3
    };
    Object.assign(params, this.filterObj);
    Api.consoleList(params).then((res: any) => {
      if (res.code == 200) {
        this.checkedMaterials = [];
        this.checkAll = false;
        const data = res.data.list || [];
        const { size, page, total } = res.data;
        if (data.length <= 0 || isNull(data)) {
          if (
            this.filterObj.just_me != 0 ||
            this.filterObj.keyword != "" ||
            this.filterObj.create_start != 0 ||
            this.filterObj.create_end != 0
          ) {
            this.emptyFlag = true;
          } else {
            this.emptyFlag = false;
          }
          this.linkList = [];
          // if (this.meta.page != 1) {
          //   this.meta.page -= 1;
          //   this.getlist();
          // }
          return;
        }
        this.meta = { size, page, total };
        this.linkList = data;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 新增链接
  private addBtn() {
    // this.isShow = true;
    // this.title = "新增链接";
    // this.value = {};

    this.$router.push({ 
      name: "myapp_enterprise_material-library_create-link", 
      query: { 
        groupId: this.filterObj.source_group_id,
      } 
    });

    localStorage.setItem("activTab", "3");
  }
  // 关闭新增弹框
  private closeAdd() {
    this.isShow = false;
    this.getlist();
    // 判断是否从新增子任务页面跳转进来
    // 清除地址栏携带参数
    if (this.$route.query.isAddTask) {
      this.$router.push({ query: {} });
    }
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
        this.editItem(row);
      }
    });
  }
  private updateStatusApi(ids: any, sourceStatus: any) {
    Api.updateStatus({
      ids: ids,
      source_type: 3,
      source_status: sourceStatus
    }).then((res: any) => {
      if (res.code === 200) {
        if (sourceStatus === 0) {
          this.$message.success("已启用");
        } else {
          this.$message.success("已禁用");
        }
        this.getlist();
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

    const res:any = await Api.setTopApi({
      source_id: row.id,
      source_type: row.type,
      is_top: row.is_top === 1 ? 2 : 1
    })

    if (res.code === 200) {
      this.getlist();
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
  mounted() {
    this.getlist();
    // 判断是否从新增子任务页面跳转进来
    if (this.$route.query.isAddTask) {
      this.addBtn();
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
  // height: 516px;
  // overflow-y: auto;
  // display: flex;
  .empty-content {
    // align-self: center;
    // margin-left: 450px;
    height: 360px;
    margin-top: 120px;
    text-align: center;
    color: #909399;
  }
}
.footer {
  margin-top: 20px;
  text-align: right;
}
</style>
