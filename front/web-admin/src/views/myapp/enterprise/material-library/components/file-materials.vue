<template>
  <div class="file-materials">
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
        type="primary"
        class="onload-pic"
        @click="onloadPicFn"
        size="small"
        v-if="filterObj.source_group_id != -2"
        >上传文件</el-button
      >
    </div>
    <div class="content clearfix" v-if="fileList.length">
      <el-checkbox-group
        v-model="checkedMaterials"
        @change="handleCheckedMaterialsChange"
      >
        <div class="card" v-for="(item, index) in fileList" :key="index">
          <div class="file-type" :class="`file-bg${item.file_type}`">
            {{ item.file_type != 5 ? customIcon[item.file_type - 1] : "" }}
            <i
              v-if="item.file_type == 5"
              class="icon el-icon-document-remove"
              aria-hidden="true"
            ></i>
          </div>
          <div class="file-title">
            <span>
              <el-checkbox :label="item.id" :key="index"
                ><i aria-hidden="true"></i
              ></el-checkbox>
            </span>
            <el-popover
              placement="top-start"
              v-if="item.name.length >= 8"
              title=""
              width="200px"
              trigger="hover"
              :content="item.name"
            >
              <h3 slot="reference">{{ item.name.slice(0, 8) }}...</h3>
            </el-popover>
            <h3 v-else>{{ item.name }}</h3>
          </div>
          <div class="date-time">
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
            <el-button
              size="small"
              type="text"
              class="unload"
              @click="downloadFn(item.url, item.name)"
              >下载</el-button
            >
            <el-button
              size="small"
              type="text"
              class="unload"
              @click="editFn(item.id)"
              >编辑</el-button
            >
            <el-button
              @click="deleteItem(item.id)"
              size="small"
              type="text"
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
          没有找到相关文件素材
        </p>
      </div>
      <div v-else>
        <img src="@/assets/images/empty-box.png" alt="" />
        <p>暂无素材</p>
      </div>
    </div>
    <div class="paging" v-if="fileList.length">
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
      @closeDialog="closemoveFn"
      :multipleSelection="checkedMaterials"
    ></movegroup>
    <onload-file
      v-if="upFile"
      @closeDialog="closeUploadFn"
      @save="fileListFn"
      :groupId="filterObj.source_group_id"
      :groupName="filterObj.name"
      :sourceId="source_id"
      @isupdate="isupdateFn"
    ></onload-file>
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
          property="name"
          label="文件名称"
        ></el-table-column>
        <el-table-column label="文件" min-width="80">
          <template slot-scope="scope">
            <div class="file-type-row" :class="`file-bg${scope.row.file_type}`">
              {{
                scope.row.file_type != 5
                  ? customIcon[scope.row.file_type - 1]
                  : ""
              }}
              <i
                v-if="scope.row.file_type == 5"
                class="icon el-icon-document-remove"
                aria-hidden="true"
              ></i>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import movegroup from "../components/move-group.vue";
import onloadFile from "../components/onload-file.vue";
import API from "@/api/enterprise/material-library";
import { downLoad } from "@/utils/down-load";
import COSSDK from "cos-js-sdk-v5";
import { ApiResult } from "@/types";
import {
  Credential,
  getFileCredentialApi
} from "@/api/enterprise/sensitive-word";
import PerModule from "@/store/modules/permission";
@Component({
  components: {
    movegroup,
    onloadFile
  }
})
export default class extends Vue {
  @Prop() filterObj!: any; // 筛选条件
  //------------------------------------------------------data-----------------------------------------------------------/
  private checkAll = false; //复选框全选
  private isIndeterminate = false; //复选框全选
  private checkedMaterials: any = []; //选中的复选框
  private multipleSelectionData: Array<any> = [];
  private overdueData: Array<any> = []; // 已过期的素材
  private overdueDialog = false; // 已过期的素材的弹窗
  private fileList = []; //文件列表
  private moveGtoupacom = false; //移动分组
  private upFile = false; //上传文件
  private source_id = 0; // 素材ID

  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private customIcon = ["W", "X", "P", "PDF"]; //自定义icon
  private cos: any = null;
  private credential?: Credential;

  get ep() {
    return PerModule.enterprise
  }

  private mounted() {
    this.fileListFn();
    this.initCos();
    if (this.$route.query.isAddTask) {
      this.onloadPicFn();
    }
  }

  //----------------------------------------------------methods-----------------------------------------------------------/
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
  //下载
  private downloadFn(url: string, name: string) {
    // const sliceUrl = url.slice(29);
    // url = uploads/2020/11/16/MDCRncszxDmRnkQxXzS4SE4fRteA43WW.mp4
    // const Url = this.cos.getObjectUrl({
    //   Bucket: this.credential?.bucket,
    //   Region: this.credential?.region /* 存储桶所在地域，必须字段 */,
    //   Key: sliceUrl
    // });
    downLoad(url, name);
  }
  //编辑
  private editFn(id: number) {
    this.source_id = id;
    this.upFile = true;
  }
  //全选的方法
  private handleCheckAllChange(val: any) {
    const arr: any = [];
    this.fileList.forEach((item: any) => {
      arr.push(item.id);
    });
    this.checkedMaterials = val ? arr : [];
    this.multipleSelectionData = val ? this.fileList : [];
    this.isIndeterminate = false;
  }
  //单选的方法
  private handleCheckedMaterialsChange(value: any) {
    const checkedCount = value.length;
    this.checkAll = checkedCount === this.fileList.length;
    this.isIndeterminate =
      checkedCount > 0 && checkedCount < this.fileList.length;
    this.multipleSelectionData = this.fileList.filter((item: any) => {
      return value.indexOf(item.id) > -1;
    });
  }

  //打开上传文件弹框
  private onloadPicFn() {
    this.source_id = 0;
    this.upFile = true;
  }
  // 关闭上传文件弹框
  private closeUploadFn() {
    this.upFile = false;
    if (this.$route.query.isAddTask) {
      this.$router.push({ query: {} });
    }
  }
  private closemoveFn(val?: boolean) {
    if (val) {
      this.fileListFn();
    }
    this.moveGtoupacom = false;
  }
  // 点击移动分组
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
        this.fileListFn();
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
        this.fileListFn();
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
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.fileListFn();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.fileListFn();
  }
  //获取文件列表
  private fileListFn(isPage?: boolean) {
    if (isPage) {
      this.meta.page = 1;
      this.meta.size = 10;
    }
    const params = {
      page: this.meta.page, //页码
      size: this.meta.size, //每页数量
      type: 6
    };
    Object.assign(params, this.filterObj);
    API.consoleList(params).then((res: any) => {
      if (res.code == 200) {
        this.resetCheckFn();
        const data = res.data;
        this.fileList = data.list || [];
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
  //给父级传值
  private isupdateFn(val: string) {
    this.$emit("isupdate", val);
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
      source_type: 6,
      source_status: sourceStatus
    }).then((res: any) => {
      if (res.code === 200) {
        if (sourceStatus === 0) {
          this.$message.success("已启用");
        } else {
          this.$message.success("已禁用");
        }
        this.fileListFn();
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
      this.fileListFn();
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
    width: 176px;
    padding: 0 10px;
    border-radius: 4px;
    border: 1px solid $color-second-line;
    margin-right: 15px;
    margin-bottom: 15px;
    .file-title {
      // line-height: 20px;
      height: 20px;
      color: #707277;
      display: flex;
      align-items: center;
      justify-content: space-between;
      span {
        flex: 0 0 14px;
        height: 20px;
        margin-right: 5px;
      }
      h3 {
        flex: 1;
        font-size: 14px;
        color: $color-text-main;
        text-align: right;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        height: 20px;
        line-height: 20px;
      }
    }
    .file-type {
      width: 100%;
      height: 155px;
      line-height: 155px;
      border-radius: 4px;
      margin: 10px 0;
      color: #fff;
      text-align: center;
      font-size: 40px;
      .icon {
        font-size: 40px;
      }
    }
    .date-time {
      text-align: right;
      font-size: 12px;
      color: #707277;
      line-height: 20px;
      margin: 5px 0 10px 0;
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
.file-bg1 {
  background: #25589a;
}
.file-bg2 {
  background: #02843f;
}
.file-bg3 {
  background: #d34729;
}
.file-bg4 {
  background: #0e88e7;
}
.file-bg5 {
  background: #c9c9c9;
}
.file-type-row {
  width: 20px;
  height: 20px;
  color: #fff;
  text-align: center;
  line-height: 20px;
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
