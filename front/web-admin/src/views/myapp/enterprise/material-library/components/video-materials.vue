<template>
  <div class="video-material">
    <div class="title">
      <div>
        <span class="m-r-20" v-if="selectCount > 0"
          >已选：{{ selectCount }}</span
        >
        <el-button size="small" class="move-group" @click="openMoveGrop"
          >移动分组</el-button
        >
        <el-button size="small" class="remove" @click="lotRemove"
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
        @click="addVideo(0)"
        v-if="filterObj.source_group_id != -2"
        >上传视频</el-button
      >
    </div>
    <!-- stripe  是否为斑马纹 table -->
    <el-table
      :data="tableList"
      ref="multipleTable"
      stripe
      @selection-change="selectChange"
      @sort-change="sortChange"
      height="500"
    >
      <el-table-column type="selection" width="55"></el-table-column>
      <el-table-column label="视频" v-if="tableList.length > 0" min-width="290">
        <template slot-scope="scope">
          <div class="video">
            <div>
              <video
                class="video-content"
                :src="scope.row.url"
                @click="openVideo(scope.row.url)"
              />
              <img
                @click="openVideo(scope.row.url)"
                class="video-play"
                src="@/assets/images/video-play.png"
                alt=""
              />
            </div>
            <el-popover
              placement="top-start"
              title=""
              width="450"
              trigger="hover"
              :content="scope.row.name"
            >
              <div slot="reference" class="text-normal m-l-5">
                <span class="text-span">{{ scope.row.name }}</span>
              </div>
            </el-popover>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="视频" v-else min-width="150"> </el-table-column>
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
              style="margin-right: 5px;"
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
        label="关联商品"
        min-width="320"
        v-if="false"
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
      <el-table-column label="大小" prop="size" width="150">
        <template slot-scope="scope">
          <span>{{ (scope.row.size / 1024 / 1024).toFixed(2) }}MB</span>
        </template>
      </el-table-column>
      <el-table-column
        label="创建时间"
        prop="created_at"
        width="190"
        show-overflow-tooltip
        sortable
      >
        <template slot-scope="scope">
          <span>{{
            scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="start_time" sortable label="有效期" width="190">
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
      <el-table-column label="操作" width="300" fixed="right">
        <template slot-scope="scope">
          <el-button @click="addVideo(scope.row.id)" type="text" size="small"
            >编辑</el-button
          >
          <el-button
            type="text"
            size="small"
            @click="delItem(scope.row.id)"
            style="color:#303133;"
            >删除</el-button
          >
          <el-button
            type="text"
            size="small"
            @click="downloadFn(scope.row.url, scope.row.name)"
            >下载</el-button
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
      <div
        slot="empty"
        v-if="!emptyFlag && tableList.length == 0"
        style="margin-top: 50px;"
      >
        <img class="empty-img" src="@/assets/images/empty-box.png" alt="" />
        <div class="empty-box">
          <div class="desc">暂无素材</div>
        </div>
      </div>
      <div
        slot="empty"
        v-if="emptyFlag && tableList.length == 0"
        style="margin-top: 50px;"
      >
        <img class="empty-img" src="@/assets/images/empty-search.png" alt="" />
        <div class="empty-box">
          <div class="desc">没有找到相关视频素材</div>
        </div>
      </div>
    </el-table>
    <div class="footer" v-if="tableList.length > 0">
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
      :multipleSelection="this.multipleSelection"
    ></movegroup>
    <!-- 新增视频弹框 -->
    <add-video
      v-if="isShow"
      @closeDialog="closeAdd"
      :sourceId="source_id"
      :groupId="filterObj.source_group_id"
      :groupName="filterObj.name"
    ></add-video>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="is_open_video"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
    <el-dialog title="已过期素材" :visible.sync="overdueDialog">
      <el-alert
        title="以下素材有效期已过，无法操作批量启用，请重新设置素材有效期后，启用素材"
        type="warning"
        :closable="false"
      >
      </el-alert>
      <el-table style="margin-top:15px" height="450" stripe :data="overdueData">
        <el-table-column label="视频" min-width="290">
          <template slot-scope="scope">
            <div class="video">
              <div>
                <video
                  class="video-content"
                  :src="scope.row.url"
                  @click="openVideo(scope.row.url)"
                />
                <img
                  @click="openVideo(scope.row.url)"
                  class="video-play"
                  src="@/assets/images/video-play.png"
                  alt=""
                />
              </div>
              <el-popover
                placement="top-start"
                title=""
                width="450"
                trigger="hover"
                :content="scope.row.name"
              >
                <div slot="reference" class="text-normal m-l-5">
                  <span class="text-span">{{ scope.row.name }}</span>
                </div>
              </el-popover>
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
import addVideo from "../components/add-video-material.vue";
import openVideo from "./com-open-video.vue";
import COSSDK from "cos-js-sdk-v5";
import { isNull } from "@/utils/tools";
import { downLoad } from "@/utils/down-load";
import Api from "@/api/enterprise/material-library";
import { ApiResult } from "@/types";
import {
  Credential,
  getFileCredentialApi
} from "@/api/enterprise/sensitive-word";
import PerModule from "@/store/modules/permission";

@Component({
  components: {
    movegroup,
    addVideo,
    openVideo
  }
})
export default class extends Vue {
  @Prop() filterObj!: any; // 筛选条件
  //------------------------------------------------------data-----------------------------------------------------------/
  // 当前页
  private currentPage = 1;
  // 是否全选
  private checked = false;
  private moveGtoupacom = false;
  private selectCount = 0; //已选条数
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private order_prop = ""; //要排序的 字段
  private sort_tag = ""; //sort （desc 倒序 asc 正序 ）
  private multipleSelection = [];
  private multipleSelectionData: Array<any> = []; // 选中的数据
  private overdueData: Array<any> = []; // 已过期的素材
  private overdueDialog = false; // 已过期的素材的弹窗
  private tableList = [];
  private emptyFlag = false;
  private isShow = false; //上传视频弹框是否展示
  private videoUrl = ""; //视频url
  private is_open_video = false; //是否打开视频
  private cos: any = null;
  private credential?: Credential;
  private source_id = 0; // 素材ID

  get ep() {
    return PerModule.enterprise
  }
  get relationGoodsByRow() {
    return (row) => {
      return row.sku_list? row.sku_list.map((item: any) => item) : [];
    };
  }
  mounted() {
    this.getlist();
    this.initCos();
    // 判断是否从新增子任务页面跳转进来
    if (this.$route.query.isAddTask) {
      this.addVideo(0);
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
  selectChange(val: any) {
    this.selectCount = val.length;
    this.multipleSelection = val.map((item: any) => {
      return item.id;
    });
    this.multipleSelectionData = val;
  }
  // 点击移动分组
  openMoveGrop() {
    if (this.multipleSelection.length == 0) {
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
  // 编辑
  private editItem() {
    console.log("编辑");
  }
  // 下载
  private downloadFn(url: string, name: string) {
    // const sliceUrl = url.slice(29);
    // // url = uploads/2020/11/16/MDCRncszxDmRnkQxXzS4SE4fRteA43WW.mp4
    // const Url = this.cos.getObjectUrl({
    //   Bucket: this.credential?.bucket,
    //   Region: this.credential?.region /* 存储桶所在地域，必须字段 */,
    //   Key: sliceUrl
    // });
    downLoad(url, name);
  }
  // 删除
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
  // 预览视频
  private openVideo(url: string) {
    this.videoUrl = url;
    this.is_open_video = true;
  }
  // 关闭预览视频
  private cancelVideo() {
    this.is_open_video = false;
  }
  // 点击批量删除
  private lotRemove() {
    if (this.multipleSelection.length == 0) {
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
          this.getDeleteFn(this.multipleSelection);
        } else {
          this.moveRecycleBin(this.multipleSelection);
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
  handleSizeChange(val: number) {
    this.meta.size = val;
    this.meta.page = 1;
    this.getlist();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getlist();
  }

  // 列表数据
  private getlist(page?: number) {
    if (page) {
      this.meta.page = page;
      this.meta.size = 10;
    }
    const params = {
      page: this.meta.page, //页码
      size: this.meta.size, //每页数量
      order: this.order_prop, //要排序的 字段
      sort: this.sort_tag, //sort （desc 倒序 asc 正序 ）
      type: 5
    };
    Object.assign(params, this.filterObj);
    Api.consoleList(params).then((res: any) => {
      if (res.code == 200) {
        this.multipleSelection = [];
        this.checked = false;
        const data = res.data.list || [];
        const { size, page, total } = res.data;
        if (isNull(data) || data.length <= 0) {
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
          // if (this.meta.page != 1) {
          //   this.meta.page -= 1;
          //   this.getlist();
          // }
          return;
        }
        this.meta = { size, page, total };
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
  // 新增视频
  private addVideoOld(id?: any) {
    this.source_id = id;
    this.isShow = true;
  }
  private addVideo(id?: any) {
    // this.source_id = id;
    // this.isShow = true;

    this.$router.push({ 
      name: "myapp_enterprise_material-library_create-video", 
      query: { 
        id: id|| undefined,
        groupId: this.filterObj.source_group_id,
      } 
    });

    localStorage.setItem("activTab", "5");
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
        this.addVideo(row.id);
      }
    });
  }
  private updateStatusApi(ids: any, sourceStatus: any) {
    Api.updateStatus({
      ids: ids,
      source_type: 5,
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
}
</script>
<style lang="scss" scoped>
.video-material {
  // /deep/ .el-table__body-wrapper {
  //   height: 480px;
  //   overflow-y: auto;
  // }
  .title {
    margin-bottom: 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .all-check {
    margin-right: 10px;
  }
  .video {
    display: flex;
    vertical-align: middle;
    position: relative;
    &:hover {
      cursor: pointer;
    }
    .video-content {
      width: 34px;
      height: 34px;
      border-radius: 4px;
      vertical-align: middle;
    }
    .video-play {
      position: absolute;
      left: 10px;
      top: 50%;
      transform: translateY(-50%);
    }
    .text-span {
      align-items: center;
      max-width: 300px;
      text-overflow: -o-ellipsis-lastline;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    .text-normal {
      height: 100%;
      display: flex;
      align-items: center;
    }
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
}
</style>
