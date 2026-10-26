<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <div>
          <search-item text="导购：">
            <el-input
              :placeholder="`请输入${type_list[gudie_type].name}`"
              class="search-input"
              size="medium"
              clearable
              v-model="search_key"
              @keyup.enter.native="search"
            >
              <el-select
                slot="prepend"
                v-model="gudie_type"
                class="input-with-select"
              >
                <el-option
                  v-for="item in type_list"
                  :value="item.key"
                  :label="item.name"
                  :key="item.key"
                ></el-option>
              </el-select>
            </el-input>
          </search-item>
          <search-item text="所在部门：">
            <el-select
              class="search-select"
              :value="select_dp.name"
              placeholder="请选择所属部门"
              ref="selectBlur"
              size="medium"
              clearable
              @change="changeSearch"
            >
              <el-option :value="select_dp.name" :label="select_dp.name">
                <el-tree
                  class="treeStyle"
                  :props="dp_props"
                  :load="loadNode"
                  node-key="dp_id"
                  :expand-on-click-node="true"
                  :check-strictly="true"
                  :check-on-click-node="true"
                  @node-click="handleNodeClick"
                  lazy
                  ref="tree"
                  @current-change="blurFn"
                ></el-tree>
              </el-option>
            </el-select>
          </search-item>

          <search-item text="导购类型：">
            <cy-select
              v-model="search_active_type"
              size="medium"
              placeholder="请选择导购类型"
              @change="changeGuideType"
              clearable
            >
              <cy-option
                v-for="item in active_type_list"
                :value="item.id"
                :label="item.name"
                :key="item.id"
              >
              </cy-option>
            </cy-select>
          </search-item>
          <search-item text="状态：">
            <cy-select
              v-model="search_status"
              size="medium"
              placeholder="请选择在职状态"
              clearable
            >
              <cy-option
                v-for="item in status_list"
                :value="item.value"
                :label="item.label"
                :key="item.value"
              >
              </cy-option>
            </cy-select>
          </search-item>
          <div>
            <el-button class="" size="small" type="primary" @click="search">
              查询
            </el-button>
            <el-button class="mar-l-20" size="small" @click="reset">
              重置
            </el-button>
            <div class="data-button">
              <span>数据更新时间：</span>
              <span>{{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
              <!-- <el-button
                type="text"
                @click="synchroData"
                size="medium"
                class="mar-l-10"
                >更新</el-button
              >
              <el-button
                size="small"
                type="primary"
                @click="batchAddGuide"
                icon="el-icon-refresh"
                class="mar-l-10"
                >手动更新</el-button
              > -->
            </div>
          </div>
        </div>
      </template>
      <template slot="table">
        <el-table
          stripe
          :data="list"
          class="mar-t-20"
          v-if="page_init"
          tooltip-effect="light"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="index"> </el-table-column>

          <el-table-column
            prop="mem_userid"
            label="导购ID"
            width="180"
          ></el-table-column>
          <el-table-column label="导购" min-width="160">
            <div slot-scope="scope">
              <img
                :src="
                  scope.row.mem_avatar
                    ? scope.row.mem_avatar
                    : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                "
                alt=""
                class="table-text-img"
              /><span v-text="scope.row.mem_name"></span>
            </div>
          </el-table-column>
          <el-table-column label="类型" min-width="120">
            <div slot-scope="scope">
              <span v-if="scope.row.type === 3">企业导购</span
              ><span v-if="scope.row.type === 4">企业管理</span>
            </div>
          </el-table-column>

          <el-table-column
            label="职务"
            min-width="150"
            key="position"
            :show-overflow-tooltip="true"
          >
            <template slot-scope="scope">{{
              scope.row.position || "--"
            }}</template>
          </el-table-column>

          <el-table-column
            prop="dp_name"
            label="导购所在部门"
            min-width="160"
          ></el-table-column>

          <el-table-column
            key="range_info"
            label="管理范围"
            min-width="200"
            :show-overflow-tooltip="true"
          >
            <template slot-scope="scope">
              <span
                v-if="scope.row.range_info.length > 0"
                class="range-info-style"
              >
                {{ scope.row.range_info.toString() }}
              </span>
              <div v-else>--</div>
            </template>
          </el-table-column>

          <el-table-column prop="worker_status" label="状态" width="100">
            <template slot-scope="scope">
              {{ scope.row.worker_status === 1 ? "在职" : "离职" }}
            </template>
          </el-table-column>
          <el-table-column prop="mem_mobile" label="手机号" width="160">
            <div slot-scope="scope">
              {{
                scope.row.mem_mobile
                  ? scope.row.mem_mobile.replace(
                      /(\d{3})(\d{4})(\d{4})/,
                      "$1****$3"
                    )
                  : ""
              }}
            </div>
          </el-table-column>

          <el-table-column label="创建时间" width="190">
            <template slot-scope="scope"
              ><span>{{
                parseInt(scope.row.created_at)
                  | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span></template
            >
          </el-table-column>

          <el-table-column prop="zip" label="操作" width="150" fixed="right">
            <template slot-scope="scope">
              <span
                v-if="!scope.row.qr_code"
                @click="createIamge(scope.row)"
                class="c-main"
              >
                生成海报
              </span>
              <span
                v-if="scope.row.qr_code"
                @click="downloadImage(scope.row)"
                class="c-main "
                >下载海报
              </span>

              <el-dropdown placement="bottom" class="mar-l-10 pad-l-10">
                <span class="c-blue"
                  >更多<em class="el-icon-arrow-down el-icon--right"></em
                ></span>
                <el-dropdown-menu slot="dropdown">
                  <!-- 企业管理导购 显示 -->
                  <el-dropdown-item
                    v-if="scope.row.type === 4"
                    type="text"
                    @click.native="getRangeList(scope.row)"
                    class="operation-success"
                    :disabled="scope.row.worker_status === 2"
                    :class="[
                      scope.row.worker_status === 2
                        ? 'disable-color'
                        : 'operate-hover',
                    ]"
                    >编辑管理范围</el-dropdown-item
                  >
                  <el-dropdown-item
                    v-if="scope.row.type === 3 || scope.row.type === 4"
                    @click.native="translateInfo(scope.row)"
                    type="text"
                    class="operation-success "
                    :class="[
                      scope.row.worker_status === 2
                        ? 'disable-color'
                        : 'operate-hover',
                    ]"
                    :disabled="scope.row.worker_status === 2"
                    >{{
                      scope.row.type === 3 ? "转换为企业管理" : "转换为企业导购"
                    }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
            </template>
          </el-table-column>

          <div slot="empty">
            <layout-empty
              empty_type="tag"
              v-if="search_key"
              :message="'没有找到相关导购'"
            ></layout-empty>
            <layout-empty
              v-if="!search_key"
              :message="'暂无导购'"
            ></layout-empty>
          </div>
        </el-table>

        <!-- 选择人员组件 -->
        <choice-staff
          v-if="member_show"
          :visible="member_show"
          :dialog_name="
            isEditGuide
              ? '编辑管理范围'
              : `添加${
                  dialog_type === 3
                    ? '企业导购'
                    : dialog_type === 2
                    ? '社群导购'
                    : '企业管理'
                }`
          "
          choice_name="已选员工"
          :remind_name="remind_name"
          @close="closeModal"
          @submit="getChildData"
          :AllExpandClick="true"
          :IsChoiceChild="true"
          :isEdit="isEditGuide"
          :cur_step="cur_step"
          :checked_one="checked_one"
          :isGuide="true"
          :selectedRange="editedGuide.rangeList"
          :selectedRangeType="editedGuide.rangeType"
        ></choice-staff>
      </template>

      <template slot="footer">
        <el-pagination
          v-if="list && list.length !== 0 && page_init"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          class="page-right font-weight-400"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import SelectInput from "@/components/business/select-input/index.vue";
import ChoiceStaff from "./model/choice-staff.vue";
import Api from "@/api/wisdom/Guide";
import personApi from "@/api/wisdom/Personnel";
import { Throttle } from "@/utils/DebounceThrottle";
import UserModule from "@/store/modules/user";
import ShareModule from "@/shared/index";

import ApiSetting from "@/api/setting/setting";

@Component({
  components: {
    LayoutList,
    LayoutEmpty,
    ChoiceStaff,
    SelectInput,
  },
})
export default class extends Vue {
  // data数据
  private dp_props = {
    label: "dp_name",
    children: [],
    isLeaf: "leaf",
  };
  private epNameDetail = ""; //品牌名称
  private updateLock = false; // 更新锁
  private updateTime = Date.now() / 1000;
  private isShowTree = false;
  private dialog_type = 1;
  private remind_name = "";
  private isEditGuide = false; // 是否正在编辑导购
  private isUpgrade = false; // 是否正在进行转化为企业管理
  private editedGuide: any = {}; // 正在编辑的导购信息
  private cur_step = 0;
  private checked_one = 0;
  private is_show_sync = false;
  private timerId = 0;
  private text: any = "";
  private item: object = {
    date: "2016-05-03",
    name: "王小虎",
    province: "上海",
    city: "普陀区",
    address: "上海市普陀区金沙江路 1518 弄",
    zip: 200333,
  };
  // 页面初始化设置
  private page_init = false;
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private list: Array<any> = [];
  private type_list: Array<any> = [
    { name: "导购ID", key: 0, keyword: "worker_no" },
    { name: "导购名称", key: 1, keyword: "mem_name" },
  ];
  private member_show = false;
  private search_key = "";
  private gudie_type = 1;
  private search_active_type = "";
  private search_status = ""; // 选中在职状态
  private status_list: Array<any> = [
    { label: "在职", value: 1 },
    { label: "离职", value: 2 },
  ];
  private guide_text = "";

  private active_type_list = [
    { name: "企业管理", id: 4 },
    { name: "企业导购", id: 3 },
  ];
  // 所选中导购
  private guideList = [];
  private select_dp = { name: "", value: "" };

  private throttle = new Throttle();
  private openControlLimit = false; // 编辑管理范围

  get USER() {
    return UserModule.getUserInfo;
  }

  // 生命周期
  private beforeDestroy() {
    this.timerId && clearTimeout(this.timerId);
  }
  private mounted() {
    this.getEpConfig();
    this.fetchList();
  }
  // 方法
  private handleSelectionChange(val: any) {
    console.log(val);
    this.guideList = val;
  }
  private changeGuideType() {
    console.log(this.search_active_type);
  }
  private changeSearch() {
    this.select_dp.value = "";
    this.select_dp.name = "";
    this.search();
    // this.$refs.tree.setCurrentKey(null);
  }
  // 手动点击同步数据
  private synchroData() {
    if (this.updateLock) {
      this.$message("努力同步中，请耐心等待");
      return;
    }
    this.updateLock = true;
    personApi.syncOrg({}).then(
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      () => {
        this.getAutoStatus();
      },
      (err) => {
        console.error("syncOrg", err);
        this.updateLock = false;
      }
    );
  }
  // 手动同步数据
  private batchAddGuide() {
    Api.batchAddGuide().then((res: any) => {
      if (res.code === 200) {
        this.$message.success("已同步");
        this.meta.page = 1;
        this.fetchList();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 自动同步数据
  private getAutoStatus(first = 0) {
    personApi
      .syncStatus({})
      .then((res) => {
        const data = res.data.status;
        this.is_show_sync = data === "processing" ? true : false;
        if (this.is_show_sync) {
          first == 0 && this.$message.success("操作成功，更新中");
          this.timerId = setTimeout(() => {
            this.getAutoStatus(1);
          }, 8000) as unknown as number;
        } else {
          this.updateLock = false;
          this.$message.success("更新成功");
          this.fetchList();
        }
      })
      .catch((err) => {
        this.$message.error(err);
        this.updateLock = false;
      });
  }

  // 获取同步状态
  private getStatus() {
    personApi
      .syncStatus({})
      .then((res) => {
        const data = res.data.status;
        this.is_show_sync = data === "processing" ? true : false;
        if (!this.is_show_sync && this.timerId) {
          this.updateLock = false;
          this.$message.success("更新成功");
          this.fetchList();
          window.clearInterval(this.timerId);
        }
      })
      .catch((err) => {
        this.$message.error(err);
        this.updateLock = false;
      });
  }
  // 打开添加导购弹框
  private addGuide(type: number, isEdit?: boolean) {
    this.dialog_type = type;
    if (type === 4) {
      this.cur_step = 1;
      this.checked_one = 1;
      this.remind_name =
        "管理人员可以作为普通导购进行推广营销，也可以查看管理范围内的罗盘数据";
      this.isEditGuide = !!isEdit;
    }
    this.member_show = true;
  }
  // 转换为企业管理/导购
  private translateInfo(info: any) {
    this.isUpgrade = true;
    if (info.type === 3) {
      this.editedGuide = {
        ...info,
        rangeList: [],
        rangeType: 4,
      };

      this.addGuide(4, true);
    } else {
      const params = {
        type: 3,
        name: info.mem_name,
        mobile: info.mem_mobile,
      };
      this.upgradeAdmin(params);
    }
  }

  // 获取企业管理导购下的门店或部门
  private getRangeList(info: any) {
    this.editedGuide = info; // 正在编辑的导购信息
    const params = {
      guide_mobile: info.mem_mobile,
    };
    Promise.all([
      Api.wisdomGuideStoreList(params),
      Api.wisdomGuideDpList(params),
    ])
      .then((res: any) => {
        let rangeList: Array<any> = [];
        let rangeType = 2;
        if (res[0].data.list.length) {
          rangeList = res[0].data.list;
          rangeType = res[0].data.range_type;
        } else if (res[1].data.list.length) {
          rangeList = res[1].data.list;
          rangeType = res[1].data.range_type;
        } else {
          rangeList = [];
          rangeType = 2;
        }
        this.$set(this.editedGuide, "rangeList", rangeList);
        this.$set(this.editedGuide, "rangeType", rangeType);
        this.addGuide(4, true);
      })
      .catch(() => {
        this.$set(this.editedGuide, "rangeList", []);
        this.$set(this.editedGuide, "rangeType", 2);
        this.addGuide(4, true);
      });
  }
  // 获取企业
  private getEpConfig() {
    ApiSetting.getEpConfig().then((res) => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epNameDetail = name;
      // 缓存企业名称，给详情页面用
      localStorage.setItem("epNameDetail", name);
    });
  }

  // 选择人员完成
  private getChildData({
    isEdit,
    memberList,
    rangeList,
    rangeType,
  }: {
    isEdit: boolean;
    memberList: Array<any>;
    rangeList: Array<any>;
    rangeType: string | number;
  }) {
    const params = {};
    console.log(this.isUpgrade);
    if (this.isUpgrade) {
      // 转换为企业管理
      params["type"] = 4;
      params["name"] = this.editedGuide.mem_name;
      params["mobile"] = this.editedGuide.mem_mobile;
      params["range"] = rangeList
        ? rangeList.map((v: any) => v.dp_third_id)
        : [];
      params["range_type"] = 4;
      this.upgradeAdmin(params);
    } else if (isEdit) {
      // 编辑
      params["guide_mobile"] = this.editedGuide.mem_mobile;
      params["range"] = rangeList
        ? rangeList.map((v: any) => v.dp_third_id)
        : [];
      params["range_type"] = rangeType;
      this.editGuideRequest(params);
    } else {
      if (!memberList || !memberList.length) return;
      // 添加
      params["type"] = this.dialog_type;
      params["member"] = memberList;
      if (this.cur_step === 1) {
        params["range"] = rangeList ? rangeList.map((v: any) => v.id) : [];
        params["range_type"] = rangeType;
      }
      this.addGuideRequest(params);
    }
  }

  private addGuideRequest(params: any) {
    Api.wisdomGuideAdd(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("添加成功");
        this.closeModal();
        this.fetchList();
      } else {
        this.$message.error(res.message);
      }
    });
  }

  private editGuideRequest(params: any) {
    Api.wisdomGuideEdit(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("编辑成功");
        this.closeModal();
        this.fetchList();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 转换为企业管理提交
  private upgradeAdmin(params: any) {
    Api.upgradeAdmin(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("操作成功");
        this.closeModal();
        this.fetchList();
      } else {
        this.$message.error(res.message);
      }
      // 移除转换为企业管理弹窗标示
      this.isUpgrade = false;
    });
  }
  // 关闭添加导购弹框
  private closeModal() {
    this.member_show = false;
    this.isEditGuide = false;
    this.editedGuide = {};
  }
  // 更新导购信息
  private updateList() {
    if (this.guideList.length === 0) {
      this.$message.warning("请先选择要更新信息的导购");
    }
  }
  // 导出导购列表
  private exportGuide() {
    const params = {
      search: {
        dp_id: this.select_dp.value,
        type: this.search_active_type,
      },
      route_url: "wisdom-guide/list",
      tag: "wisdom-guide-list",
      version: this.USER.version,
    };
    params.search[this.type_list[this.gudie_type].key] = this.search_key;
    Api.wisdomGuideExport(params).then(() => {
      this.$confirm(
        "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
        "文件下载提示",
        {
          distinguishCancelAndClose: true,
          confirmButtonText: "前往下载中心",
          cancelButtonText: "取消",
          type: "success",
        }
      ).then(() => {
        ShareModule.getShared().pushState("myapp_download-center");
      });
    });
  }

  // 生成海报
  private createIamge(row: any) {
    this.$confirm(`你确认要生成海报吗`, "提示", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
    }).then(() => {
      Api.workerPoster({ id: row.id }).then(() => {
        this.$message.success("操作成功!");
        this.fetchList();
      });
    });
  }

  // 搜索
  private search() {
    this.meta.page = 1;
    this.fetchList();
  }
  // 获取导购列表
  private fetchList() {
    const params = {
      search: {
        dp_third_id: this.select_dp.value,
        type: this.search_active_type,
        filter_type: 1,
        worker_status: this.search_status,
      },
      page: this.meta.page,
      size: this.meta.size,
    };
    params.search[this.type_list[this.gudie_type].keyword] = this.search_key;
    console.log(params.search);
    Api.wisdomGuideList(params).then((res) => {
      const { list, page, size, total } = res.data;
      this.list = list;
      this.meta = { page, size, total };
      this.page_init = true;
    });
  }

  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.fetchList();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.fetchList();
  }

  // 选择部门
  private handleNodeClick(node: any) {
    this.select_dp.name = node.dp_name;
    this.select_dp.value = node.dp_third_id;
    this.search();
  }
  // 关闭筛选
  private blurFn() {
    (this.$refs.selectBlur as any).blur();
  }
  //动态加载节点
  private async loadNode(node: any, resolve: any) {
    if (node.level === 0) {
      const {
        data: { list },
      } = await personApi.depListApi({});
      list.forEach((item: any) => {
        if (!item.is_child_department) item.leaf = true;
      });
      resolve(list);
    }
    if (node.level >= 1) {
      this.getIndex(node, resolve);
    }
  }

  //判断部门是否有子节点
  private setChecked(arr: any) {
    arr.forEach((item: any) => {
      if (!item.is_child_department) item.leaf = true;
    });
  }

  // 加载叶子节点数据函数
  private getIndex(node: any, resolve: any) {
    const params = { dp_id: node.data.dp_id };
    personApi.depListApi(params).then((res) => {
      this.setChecked(res.data.list);
      resolve(res.data.list);
    });
  }

  // 下载图片
  downloadImage(item: any) {
    console.log("下载海报");
    const _this = this;
    const image = new Image();
    // 解决跨域 Canvas 污染问题
    image.setAttribute("crossOrigin", "anonymous");
    image.onload = function() {
      const canvas = document.createElement("canvas");
      canvas.width = image.width;
      canvas.height = image.height;
      const context: any = canvas.getContext("2d");
      context.drawImage(image, 0, 0, image.width, image.height);
      const url = canvas.toDataURL("image/png"); //得到图片的base64编码数据
      const blob = _this.dataURLtoBlob(url);
      const objurl = URL.createObjectURL(blob);
      const a = document.createElement("a"); // 生成一个a元素
      const event = new MouseEvent("click"); // 创建一个单击事件
      a.download = item.mem_name || "photo"; // 设置图片名称
      a.href = objurl; // 将生成的URL设置为a.href属性
      a.dispatchEvent(event); // 触发a的单击事件
    };
    image.src = item.qr_code;
  }

  // base64转换成blob
  private dataURLtoBlob(dataurl: any) {
    // eslint-disable-next-line prefer-const
    let arr = dataurl.split(","),
      // eslint-disable-next-line prefer-const
      mime = arr[0].match(/:(.*?);/)[1],
      // eslint-disable-next-line prefer-const
      bstr = atob(arr[1]),
      n: any = bstr.length,
      // eslint-disable-next-line prefer-const
      u8arr = new Uint8Array(n);
    // eslint-disable-next-line no-const-assign
    while (n--) {
      u8arr[n] = bstr.charCodeAt(n);
    }
    return new Blob([u8arr], { type: mime });
  }

  // 重置
  private reset() {
    this.search_key = "";
    this.gudie_type = 1;
    this.select_dp.value = "";
    this.select_dp.name = "";
    this.search_active_type = "";
    this.search_status = "";
    this.search();
  }
}
</script>
<style lang="scss" scoped>
.supernatant {
  color: #303133;
  /deep/.el-button--text {
    color: #6f3112;
  }
}
.mar-l-15 {
  margin-left: 15px;
}

.mar-l-30 {
  margin-left: 30px;
}

.mar-r-20 {
  margin-right: 20px;
}

.mar-l-20 {
  margin-left: 20px;
}

.mar-t-20 {
  margin-top: 20px;
}

.mar-b-20 {
  margin-bottom: 20px;
}

.mar-r-15 {
  margin-right: 15px;
}

.mar-t-15 {
  margin-top: 15px;
}
.c-main {
  color: #303133;
  cursor: pointer;
}
.c-blue {
  color: #6f3112;
  cursor: pointer;
}

.select-input-title {
  margin-right: 20px;
  display: inline-block;
  line-height: 14px;
  font-size: 14px;
  // font-weight: 700;
  font-style: Hiragino Sans GB;
}

.select-input {
  display: inline-block;
}

.enter-input {
  display: inline-block;
}
.operate-hover {
  color: #6f3112;
}
.operation-error {
  color: #f56c6c !important;
}
.operate-hover:hover {
  background-color: #f5f6fa;
  color: #606266;
}

.data-button {
  display: inline-block;
  margin-left: 30px;
  color: #666666;
  font-size: 14px;
}

.input-with-select {
  width: 110px;
}

.search-input {
  width: 280px !important;
}

.top-select-bc {
  background: #eeeeee;
  overflow: hidden;
}

.top-table {
  padding: 15px;
  background: #eeeeee;
  display: inline-block;
  width: 246px;
  margin-right: 20px;

  .remind-info {
    margin: 15px 0;
    color: #666666;
    font-size: 12px;
  }
}

.treeStyle {
  .el-tree-node.is-current.is-focusable.is-checked {
    color: #409eff !important;
  }
}

.el-select-dropdown__item.selected {
  font-weight: normal !important;
  height: auto !important;
}

.el-icon-refresh-right {
  color: #409eff;
  cursor: pointer;
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
  width: 30px;
  height: 30px
}

.search-select {
  width: 280px;
}

/deep/.layout-list .slide-box {
  display: none;
}
.el-scrollbar {
  .el-scrollbar__view {
    .el-select-dropdown__item {
      font-weight: normal !important;
      height: auto;
    }
  }
}
.mar-l-10 {
  margin-left: 10px;
}

.mar-r-10 {
  margin-right: 10px;
}

.update-time {
  color: #666666;
  margin-left: 20px;

  .update-btn {
    color: #409eff;
    cursor: pointer;
  }
}

/deep/ .cell {
  padding-left: 10px !important;
  padding-right: 10px !important;
}
</style>
