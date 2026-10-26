<template>
  <div class="text-materials">
    <div class="title">
      <div>
        <el-button
          class="move-group"
          @click="openMoveGrop"
          size="small"
        >移动分组</el-button>
        <el-button
          class="remove"
          @click="lotRemove"
          size="small"
        >批量删除</el-button>
        <el-button
          class="remove"
          @click="lotUpdateStatus(0)"
          size="small"
        >批量启用</el-button>
        <el-button
          class="remove"
          @click="lotUpdateStatus(1)"
          size="small"
        >批量禁用</el-button>
      </div>
      <!-- v-if="filterObj.children == null" -->
      <el-button
        type="primary"
        class="onload-pic"
        @click="onloadPic()"
        size="small"
        v-if="filterObj.source_group_id != -2"
      >上传图片</el-button>
    </div>
    <el-table
      tooltip-effect="light"
      :data="tableList"
      ref="multipleTable"
      stripe
      @selection-change="selectChange"
      @sort-change="sortChange"
      height="500"
    >
      <el-table-column
        fixed="left"
        type="selection"
        width="55"
      ></el-table-column>
      <el-table-column
        label="图片名称"
        prop="name"
        min-width="150"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column
        label="图片"
        prop="content"
        min-width="80"
      >
        <template slot-scope="scope">
          <div class="image-container">
            <el-image
              class="table-pic"
              :src="scope.row.url"
              :preview-src-list="getPreviewList(scope.row)"
            />
            <div
              class="pic-num"
              v-if="scope.row.merge_image && scope.row.merge_image.length > 1 && epName === 'coach'"
            >{{ scope.row.merge_image.length }}张</div>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        label="素材标签"
        min-width="320"
      >
        <template slot-scope="scope">
          <div v-if="!scope.row.category_arr.length">-</div>
          <div v-else-if="scope.row.category_arr.length <= 3">
            <el-tag
              v-for="(item, index) in scope.row.category_arr"
              :key="index"
              closable
              @close="removeCategory(scope.row, index)"
              type="info"
              style="margin-right: 10px;"
            >
              <div
                :title="item"
                class="tag-in"
              >{{ item }}</div>
            </el-tag>
          </div>
          <div v-else>
            <el-tag
              v-for="(item, index) in scope.row.category_arr.slice(0, 3)"
              :key="index"
              type="info"
              closable
              @close="removeCategory(scope.row, index)"
              style="margin-right: 10px;"
            >
              <div
                :title="item"
                class="tag-in"
              >{{ item }}</div>
            </el-tag>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <el-tag
                v-for="(item, index) in scope.row.category_arr"
                :key="index"
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
        label="关联商品"
        min-width="320"
      >
        <template slot-scope="scope">
          <el-tag v-if="!relationGoodsByRow(scope.row).length" type="danger">未关联商品</el-tag>
          <div v-else-if="relationGoodsByRow(scope.row).length <= 3">
            <el-tag
              v-for="(item, index) in relationGoodsByRow(scope.row)"
              :key="index"
              type="info"
              style="margin-right: 10px;"
            >
              <div
                :title="item"
                class="tag-in"
              >{{ item }}</div>
            </el-tag>
          </div>
          <div v-else>
            <el-tag
              v-for="(item, index) in relationGoodsByRow(scope.row).slice(0, 3)"
              :key="index"
              type="info"
              style="margin-right: 10px;"
            >
              <div
                :title="item"
                class="tag-in"
              >{{ item }}</div>
            </el-tag>
            <el-popover
              placement="top-start"
              trigger="hover"
            >
              <el-tag
                v-for="(item, index) in relationGoodsByRow(scope.row)"
                :key="index"
                style="margin:5px;"
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
        label="状态"
        prop="source_status"
        min-width="100"
      >
        <template slot-scope="scope">
          {{ scope.row.source_status === 0 ? "已启用" : "已禁用" }}
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
      <el-table-column
        prop="admin_name"
        label="上传人"
        min-width="130"
      >
      </el-table-column>
      <el-table-column
        label="操作"
        fixed="right"
        width="300"
      >
        <template slot-scope="scope">
          <el-button
            @click="onloadPic(scope.row)"
            type="text"
            size="small"
          >编辑</el-button>
          <el-button
            @click="delItem(scope.row.id)"
            type="text"
            size="small"
            style="color:#303133;"
          >删除</el-button>
          <el-button
            type="text"
            size="small"
            @click="onload(scope.row)"
          >下载</el-button>
          <el-button
            type="text"
            size="small"
            class="m-l-20"
            :style="{ color: scope.row.source_status === 0 ? '#F56C6C' : '' }"
            @click="updateStatus(scope.row)"
          >{{ scope.row.source_status === 0 ? "禁用" : "启用" }}</el-button>

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
      <div
        slot="empty"
        v-if="!emptyFlag"
        style="margin-top: 50px;"
      >
        <img
          class="empty-img"
          src="@/assets/images/empty-box.png"
          alt=""
        />
        <div class="empty-box">
          <div class="desc">暂无文本素材</div>
        </div>
      </div>
      <div
        slot="empty"
        v-if="emptyFlag"
        style="margin-top: 50px;"
      >
        <img
          class="empty-img"
          src="@/assets/images/empty-search.png"
          alt=""
        />
        <div class="empty-box">
          <div class="desc">没有找到相关图片素材</div>
        </div>
      </div>
    </el-table>
    <div
      class="footer"
      v-if="meta.total > 0"
    >
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
      :multipleSelection="checkedMaterials"
    ></movegroup>
    <onload-picture
      :groupId="filterObj.source_group_id"
      :name="filterObj.name"
      :editTextDate="editTextDate"
      :dialogType="dialogType"
      v-if="uppic"
      @closeDialog="closepic"
    ></onload-picture>
    <el-dialog
      title="已过期素材"
      :visible.sync="overdueDialog"
    >
      <el-alert
        title="以下素材有效期已过，无法操作批量启用，请重新设置素材有效期后，启用素材"
        type="warning"
        :closable="false"
      >
      </el-alert>
      <el-table
        style="margin-top:15px"
        height="450"
        stripe
        :data="overdueData"
      >
        <el-table-column
          label="图片"
          min-width="80"
        >
          <template slot-scope="scope">
            <img
              class="table-pic"
              :src="scope.row.url"
              style="cursor: pointer;width: 50px; height: 50px;object-fit: cover;"
              @click="onPreview(scope.$index)"
            />
          </template>
        </el-table-column>
        <el-table-column
          show-overflow-tooltip
          property="name"
          label="图片名称"
        ></el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import movegroup from "../components/move-group.vue";
import onloadPicture from "../components/onload-picture.vue";
import Api from "@/api/enterprise/material-library";
import { isNull } from "@/utils/tools";
import { downLoad } from "@/utils/down-load";
import COSSDK from "cos-js-sdk-v5";
import { ApiResult } from "@/types";
import PerModule from "@/store/modules/permission";
import {
  Credential,
  getFileCredentialApi
} from "@/api/enterprise/sensitive-word";
@Component({
  components: {
    movegroup,
    onloadPicture,
  }
})
export default class extends Vue {
  @Prop() filterObj?: any; // 筛选条件
  //------------------------------------------------------data-----------------------------------------------------------/
  private tableList: any = []; // 表格数据
  private checkAll = false; //复选框全选
  // private isIndeterminate = true; //复选框全选
  private checkedMaterials: any = []; //选中的复选框
  private multipleSelectionData: Array<any> = []; // 选中的数据
  private overdueData: Array<any> = []; // 已过期的素材
  private overdueDialog = false; // 已过期的素材的弹窗
  // 点击当前图片的索引
  private currectIndex = 0;
  // 是否全选
  private checked = false;
  private emptyFlag = false;
  private moveGtoupacom = false; //移动分组
  private uppic = false; //上传图片
  private editTextDate = {};
  private order_prop = ""; //要排序的 字段
  private sort_tag = ""; //sort （desc 倒序 asc 正序 ）
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private cos: any = null;
  private credential?: Credential;
  private dialogType = 1; // 弹窗类型 1上传 ，2编辑
  // todo: 当前品牌
  private epName = PerModule.getEpName;

  get ep() {
    return PerModule.enterprise
  }
  get relationGoodsByRow() {
    return (row) => {
      // 区分视频素材和图片素材
      return  row.sku_list? row.sku_list.map((item: any) => item) : [];
    };
  }
  mounted() {
    this.getlist();
    this.initCos();
    // 判断是否从新增子任务页面跳转进来
    if (this.$route.query.isAddTask) {
      this.onloadPic();
    }
  }
  //----------------------------------------------------methods-----------------------------------------------------------/
  // todo: 获取图片预览
  private getPreviewList(row: any) {
    if (row.merge_image && row.merge_image.length > 0 && this.epName === "coach") {
      return row.merge_image.map((item: any) => item.url);
    } else {
      return [row.url];
    }
  }
  // 初始化腾讯对象存储Cos
  private async initCos() {
    const {
      code,
      message = "",
      data
    }: ApiResult<Credential> = await getFileCredentialApi();
    if (code !== 200) {
      this.$message(message);
      return false;
    }
    this.credential = data;
    this.cos = new COSSDK({
      getAuthorization: (options: any, callback: Function) => {
        console.log(options);
        callback({
          TmpSecretId: data?.secret_id,
          TmpSecretKey: data?.secret_key,
          XCosSecurityToken: data?.session_token,
          StartTime: data?.start_time,
          ExpiredTime: data?.expired_time
        });
      }
    });

    return true;
  }
  //todo: 合并素材下载
  private onload(row: any) {
    if(row.merge_image && row.merge_image.length > 0 && this.epName === 'coach'){
      for(let i = 0; i < row.merge_image.length; i++){
        downLoad(row.merge_image[i].url, row.merge_image[i].name);
      }
    }else{
      downLoad(row.url, row.name);
    }
  }
  //全选的方法
  private selectChange(val: any) {
    const arr: any = [];
    val.forEach((item: any) => {
      arr.push(item.id);
    });
    this.checkedMaterials = val ? arr : [];
    this.multipleSelectionData = val || [];
  }
  onloadPicOld(row?: any) {
    if (row) {
      this.editTextDate = row;
      this.dialogType = 2;
    } else {
      this.editTextDate = {};
      this.dialogType = 1;
    }
    this.uppic = true;


  }
  //  打开上传图片弹框
  onloadPic(row?: any) {
    this.$router.push({ 
      name: "myapp_enterprise_material-library_create-image", 
      query: { 
        id: row ? row.id : undefined,
        groupId: this.filterObj.source_group_id,
      } 
    });

    localStorage.setItem("activTab", "2");
  }
  // 关闭上传图片弹框
  private closepic() {
    this.uppic = false;
    this.getlist();
    // 判断是否从新增子任务页面跳转进来
    // 清除地址栏携带参数
    if (this.$route.query.isAddTask) {
      this.$router.push({ query: {} });
    }
  }
  private closemove(save: boolean) {
    this.moveGtoupacom = false;
    if (save) {
      this.getlist();
    }
  }
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
  // 删除
  private delItem(id: number) {
    this.$confirm("确定删除该素材吗?", "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    }).then(() => {
      const arr: Array<number> = [];
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
    }).then(() => {
      if (this.filterObj.source_group_id === -2) {
        this.getDeleteFn(this.checkedMaterials);
      } else {
        this.moveRecycleBin(this.checkedMaterials);
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
  getDeleteFn(ids: any) {
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
  // 列表数据
  private getlist(page?: number) {
    this.checkedMaterials = [];
    this.checkAll = false;
    if (page) {
      this.meta.page = page;
      this.meta.size = 10;
    }
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      order: this.order_prop, //要排序的 字段
      sort: this.sort_tag, //sort （desc 倒序 asc 正序 ）
      type: 2
    };
    // 对象合并
    Object.assign(params, this.filterObj);
    Api.consoleList(params).then((res: any) => {
      if (res.code == 200) {
        const data = res.data.list || [];
        this.meta.total = res.data.total;
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
    this.checkedMaterials = [];
    this.checkAll = false;
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
        const ids = row.merge_image.length > 0 ? row.merge_image.map(item=>item.id) : [row.id];
        this.updateStatusApi(ids, row.source_status === 0 ? 1 : 0);
      } else {
        // 已过期提示点击确定打开编辑弹窗
        this.onloadPic(row);
      }
    });
  }
  private updateStatusApi(ids: any, sourceStatus: any) {
    Api.updateStatus({
      ids: ids,
      source_type: 2,
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
}
</script>
<style lang="scss" scoped>
.title {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .selected {
    margin-right: 20px;
  }
}

.all-check {
  margin-right: 20px;
}

.contant {
  .cont-item {
    width: 145px;
    padding: 10px 10px 0;
    background: #ffffff;
    border-radius: 4px;
    border: 1px solid #dcdfe6;
    margin-right: 20px;
    margin-bottom: 15px;
    float: left;
  }

  .check-box {
    display: flex;
    justify-content: space-between;
  }

  img {
    width: 100%;
    height: 145px;
    cursor: pointer;
  }

  .pic-name {
    color: #303133;
    text-align: right;
    width: 100%;
  }

  .line-top-box {
    text-align: right;
    font-size: 12px;
    color: #707277;
  }

  .line {
    width: 155px;
    height: 1px;
    background: #ebeef5;
    margin-top: 10px;
  }
}

.footer {
  margin-top: 20px;
  text-align: right;
}

.desc {
  color: #909399;
}

.text-span {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  word-break: break-all;
  max-width: 115px;
}

.tag-in {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  transform: translateY(-2px);
}


.image-container {
  position: relative;
  .table-pic {
    width: 50px;
    height: 50px;
    object-fit: cover;
    cursor: pointer;
  }
  .pic-num{
    position: absolute;
    left:2px;
    top: 5px;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 30px;
    height: 15px;
    background: rgba(0,0,0,.5);
    color: #fff;
    font-size: 12px;
    border-radius: 5px;
  }
}
</style>
