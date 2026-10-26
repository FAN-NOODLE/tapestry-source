<template>
  <div class="text-materials">
    <div class="title">
      <div>
        <span class="cleck" v-if="selectCount > 0"
          >已选：{{ selectCount }}</span
        >
        <el-button class="move-group" @click="openMoveGrop" size="small"
          >移动分组</el-button
        >
        <el-button class="remove" @click="lotRemove" size="small"
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
        @click="addTextmat(1)"
        size="small"
        >新增文本素材</el-button
      >
    </div>
    <!-- stripe  是否为斑马纹 table -->
    <el-table
      tooltip-effect="light"
      :data="tableList"
      ref="multipleTable"
      stripe
      @selection-change="selectChange"
      @sort-change="sortChange"
      height="500"
    >
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column
        label="文本内容"
        prop="content"
        v-if="tableList.length > 0"
        min-width="260"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column
        label="文本内容"
        prop="content"
        v-else
        min-width="100"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column label="状态" prop="source_status" min-width="100">
        <template slot-scope="scope">
          {{ scope.row.source_status === 0 ? "已启用" : "已禁用" }}
        </template>
      </el-table-column>
      <el-table-column label="素材标签" min-width="320">
        <template slot-scope="scope">
          <div v-if="!scope.row.category_arr.length">-</div>
          <div v-else-if="scope.row.category_arr.length <= 3">
            <el-tag
              v-for="(item, index) in scope.row.category_arr"
              :key="index"
              closable
              type="info"
              style="margin-right: 10px;"
              @close="removeCategory(scope.row, index)"
            >
              <div :title="item" class="tag-in">{{ item }}</div>
            </el-tag>
          </div>
          <div v-else>
            <el-tag
              v-for="(item, index) in scope.row.category_arr.slice(0, 3)"
              :key="index"
              type="info"
              closable
              style="margin-right: 10px;"
              @close="removeCategory(scope.row, index)"
            >
              <div :title="item" class="tag-in">{{ item }}</div>
            </el-tag>
            <el-popover placement="top-start" trigger="hover">
              <el-tag
                v-for="(item, index) in scope.row.category_arr"
                :key="index"
                size="small"
                style="margin:5px;"
                closable
                @close="removeCategory(scope.row, index)"
              >
                {{ item }}
              </el-tag>
              <el-button
                type="text"
                slot="reference"
                icon="el-icon-more"
                class="el-popover-ref"
              ></el-button>
            </el-popover>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        label="创建时间"
        prop="created_at"
        width="190"
        sortable
        show-overflow-tooltip
      >
        <template slot-scope="scope">
          <span>{{
            scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="start_time"
        sortable
        label="有效期"
        min-width="190"
      >
        <template slot-scope="scope">
          <div>
            {{ scope.row.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </div>
          <div>
            {{ scope.row.end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="admin_name" label="上传人" min-width="130">
      </el-table-column>
      <el-table-column label="操作" fixed="right" width="280">
        <template slot-scope="scope">
          <el-button type="text" size="small" @click="addTextmat(scope.row)"
            >编辑</el-button
          >
          <el-button
            @click="del(scope.row.id)"
            type="text"
            size="small"
            style="color:#303133;"
            >删除</el-button
          >
          <el-button
            type="text"
            size="small"
            :style="{ color: scope.row.source_status === 0 ? '#F56C6C' : '' }"
            @click="updateStatus(scope.row)"
            >{{ scope.row.source_status === 0 ? "禁用" : "启用" }}</el-button
          >
          <el-button
            v-if="ep.isCoach && filterObj.source_group_id !== -2"
            type="text"
            size="small"
            :disabled="scope.row.is_top === 0 && !scope.row.can_top"
            @click="handleTop(scope.row)"
            >{{ scope.row.is_top === 0 ? "置顶" : "取消置顶" }}</el-button
          >
        </template>
      </el-table-column>
      <div slot="empty" v-if="!emptyFlag" style="margin-top: 50px;">
        <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
        <div class="empty-box">
          <div class="desc">暂无文本素材</div>
        </div>
      </div>
      <div slot="empty" v-if="emptyFlag" style="margin-top: 50px;">
        <img class="empty-img" src="@/assets/images/empty-search.png" alt="" />
        <div class="empty-box">
          <div class="desc">没有找到相关文本素材</div>
        </div>
      </div>
    </el-table>
    <div v-if="meta.total > 0" class="footer">
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
      @closeDialog="closemove"
      :multipleSelection="multipleSelection"
    ></movegroup>
    <!-- <move-text-mat
      v-if="moveGtoupacom"
      @closeDialog="closemove"
      :multipleSelection="multipleSelection"
    ></move-text-mat> -->
    <add-text-mat
      v-if="addtexmat"
      @closeDialog="closetex"
      :groupId="filterObj.source_group_id"
      :groupName="filterObj.name"
      :editTextDate="editTextDate"
      :title="textTitle"
    ></add-text-mat>
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
          property="content"
          label="文本素材"
        ></el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import movegroup from "../components/move-group.vue";
import addTextMat from "../components/add-text-mat.vue";
import moveTextMat from "../components/move-text-mat.vue";
import Api from "@/api/enterprise/material-library";
import { isNull } from "@/utils/tools";
import PerModule from "@/store/modules/permission";
@Component({
  components: {
    movegroup,
    addTextMat,
    moveTextMat
  }
})
export default class extends Vue {
  @Prop() filterObj?: any; // 筛选条件
  //------------------------------------------------------data-----------------------------------------------------------/
  private tableList: any = []; // 表格数据
  // 是否全选
  private checked = false;
  private moveGtoupacom = false;
  private emptyFlag = false;
  private addtexmat = false;
  private id = 0;
  private selectCount = 0; //已选条数
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private order_prop = ""; //要排序的 字段
  private sort_tag = ""; //sort （desc 倒序 asc 正序 ）
  private editTextDate = {};
  private textTitle = "";
  private multipleSelection: Array<number> = [];
  private multipleSelectionData: Array<any> = [];
  private overdueData: Array<any> = []; // 已过期的素材
  private overdueDialog = false; // 已过期的素材的弹窗

  get ep() {
    return PerModule.enterprise
  }
  //----------------------------------------------------methods-----------------------------------------------------------/
  private closetex() {
    this.addtexmat = false;
    this.getlist();
    // 判断是否从新增子任务页面跳转进来
    // 清除地址栏携带参数
    if (this.$route.query.isAddTask) {
      this.$router.push({ query: {} });
    }
  }
  private addTextmat(row: any) {
    if (row != 1) {
      this.textTitle = "编辑文本素材";
      this.editTextDate = row;
    } else {
      this.textTitle = "新增文本素材";
    }
    this.addtexmat = true;
  }

  private closemove(save: boolean) {
    this.moveGtoupacom = false;
    if (save == true) {
      this.getlist();
    }
  }
  // @Watch("groupId")
  // getGroupId(newwal: any) {
  //   this.id = newwal;
  //   this.getlist();
  // }
  // 列表数据
  getlist(page?: number) {
    if (page) {
      this.meta.page = page;
      this.meta.size = 10;
    }
    const params = {
      page: this.meta.page, //页码
      size: this.meta.size, //每页数量
      order: this.order_prop, //要排序的 字段
      sort: this.sort_tag, //sort （desc 倒序 asc 正序 ）
      type: 1
    };
    // 对象合并
    Object.assign(params, this.filterObj);
    Api.consoleList(params).then((res: any) => {
      if (res.code == 200) {
        const data = res.data.list || [];
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
          this.meta.total = res.data.total;
          this.tableList = [];
          return;
        }
        this.tableList = data.map(item => {
          return {
            ...item,
            category_ids_arr: item.category_ids
              ? item.category_ids.split(",")
              : [],
            category_arr: item.category_names
              ? item.category_names.split(",")
              : []
          };
        });
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击移动分组
  openMoveGrop() {
    if (this.selectCount == 0) {
      this.$message({
        message: "请选择要移动分组的素材",
        type: "warning"
      });
      return;
    }
    this.moveGtoupacom = true;
  }
  // 当选择项发生变化时会触发该事件，数组
  selectChange(val: any) {
    this.multipleSelection = [];
    this.multipleSelectionData = val || [];
    this.selectCount = val.length;
    val.map((item: any) => {
      this.multipleSelection.push(item.id);
    });
  }
  // 点击全选
  allClick() {
    // toggleAllSelection	用于多选表格，切换所有行的选中状态
    (this.$refs.multipleTable as any).toggleAllSelection();
  }

  // 删除
  private del(id: number) {
    this.$confirm("确定删除该素材吗?", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    }).then(() => {
      const arr = [];
      arr.push(id);
      if (this.filterObj.source_group_id === -2) {
        this.getDeleteFn(arr);
      } else {
        this.moveRecycleBin(arr);
      }
    });
  }
  // 点击批量删除
  private lotRemove() {
    if (this.selectCount == 0) {
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
    }).then(() => {
      if (this.filterObj.source_group_id === -2) {
        this.getDeleteFn(this.multipleSelection);
      } else {
        this.moveRecycleBin(this.multipleSelection);
      }
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
  //删除
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
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getlist();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getlist();
  }
  // 排序
  private sortChange(obj: any) {
    if (!obj.order) {
      this.order_prop = "";
      this.sort_tag = "";
    } else {
      this.order_prop = obj.prop;
      if (obj.order === "descending") {
        this.sort_tag = "desc";
      } else if (obj.order === "ascending") {
        this.sort_tag = "asc";
      }
    }
    this.meta.page = 1;
    this.getlist();
  }
  // 移除素材标签
  private removeCategory(row: any, index: any) {
    Api.removeCategory({
      source_id: row.id,
      category_id: Number(row.category_ids_arr[index])
    }).then((res: any) => {
      if (res.code == 200) {
        this.getlist();
      } else {
        this.$message.error(res.message);
      }
    });
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
        this.addTextmat(row);
      }
    });
  }
  private updateStatusApi(ids: any, sourceStatus: any) {
    Api.updateStatus({
      ids: ids,
      source_type: 1,
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
    if (this.selectCount == 0) {
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
        this.lotUpdateStatusTip(this.multipleSelection, sourceStatus);
      }
    } else {
      this.lotUpdateStatusTip(this.multipleSelection, sourceStatus);
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
    console.log("filterObj :>> ", this.filterObj);
    this.getlist();
    // 判断是否从新增子任务页面跳转进来
    if (this.$route.query.isAddTask) {
      this.addTextmat(1);
    }
  }
}
</script>
<style lang="scss" scoped>
.text-span {
  display: inline-block;
  text-align: left;
  width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.edit-class {
  margin-right: 10px;
}
.title {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  .cleck {
    margin-right: 15px;
  }
}
.all-check {
  margin-right: 10px;
}
.footer {
  margin-top: 20px;
  text-align: right;
}
.tag-in {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  transform: translateY(-2px);
}
</style>
