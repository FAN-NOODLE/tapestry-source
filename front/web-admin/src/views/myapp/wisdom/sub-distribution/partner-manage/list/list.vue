<template>
  <div class="manage-list">
    <layout-list>
      <!-- 搜索 -->
      <template slot="search">
        <search-item text="合伙人：">
          <el-select
            v-model="search_type"
            class="search-select"
            size="small"
            @change="selectChange"
          >
            <el-option
              v-for="item in search_types"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            >
            </el-option>
          </el-select>
          <el-input
            :placeholder="`请输入${search_types[search_type - 1].name}`"
            v-model="search_name"
            class="name-input search-input"
            size="small"
          >
          </el-input>
        </search-item>
        <search-item text="所属导购：">
          <search-select-member
            :selectedList="hign_choose"
            @visible-change="highSearch"
            @remove-member="removeGuide"
          >
          </search-select-member>
        </search-item>
        <search-item text="状态：">
          <el-select
            class="name-select"
            size="medium"
            v-model="partner_type"
            clearable
          >
            <el-option
              v-for="item in partner_types"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            ></el-option>
          </el-select>
        </search-item>
        <search-item text="佣金：">
          <el-select
            class="search-select"
            size="medium"
            v-model="commission_type"
            clearable
          >
            <el-option
              v-for="item in commission_types"
              :key="item.key"
              :label="item.name"
              :value="item.key"
            >
            </el-option>
          </el-select>
          <div class="com-section">
            <el-input
              size="medium"
              placeholder="最低"
              class="name-input"
              v-model="min_input"
              onkeyup="value=value.replace(/[^\d^\.]+/g,'').replace('.','$#$').replace(/\./g,'').replace('$#$','.')"
            ></el-input>
            <span>-</span>
            <el-input
              size="medium"
              placeholder="最高"
              class="name-input"
              v-model="max_input"
              onkeyup="value=value.replace(/[^\d^\.]+/g,'').replace('.','$#$').replace(/\./g,'').replace('$#$','.')"
            ></el-input>
          </div>
        </search-item>
        <div class="search-button">
          <el-button type="primary" size="small" @click="search">
            查询
          </el-button>
          <el-button class="m-l-20" size="small" @click="reset">重置</el-button>
          <div class="data-button">
            <span>数据更新时间：</span>
            <span>{{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
            <el-button
              type="text"
              @click="reload"
              size="mini"
              class="mar-l-10 el-icon-refresh"
              >更新</el-button
            >
          </div>
        </div>

        <!-- 选择导购 -->
        <select-member
          v-if="member_show"
          :visible="member_show"
          dialog_title="选择导购"
          choice_name="已选导购"
          :selectArr="hign_choose"
          :tabName="USER.version == 3 ? 3 : 1"
          :is_show_guide="true"
          @close="member_show = false"
          @submit="submitDialog"
        ></select-member>
      </template>

      <!-- 操作 -->
      <template slot="operate">
        <el-button
          size="small"
          @click="exportFile"
          :disabled="list.length === 0"
        >
          导出Excel
        </el-button>
      </template>

      <!-- 表格 -->
      <template slot="table">
        <el-table
          v-if="isLoaded"
          :data="list"
          ref="table"
          @sort-change="sortChange"
          stripe
          :header-cell-style="{
            fontSize: '14px',
            color: '#606266',
            fontWeight: '400',
            padding: '11px 0 10px 0'
          }"
        >
          <el-table-column
            label="合伙人ID"
            prop="partner_sn"
            min-width="150"
            fixed="left"
          ></el-table-column>
          <el-table-column
            label="合伙人"
            min-width="180"
            prop="partner_name"
            fixed="left"
          >
            <template slot-scope="scope">
              <div class="nameBox" v-if="scope.row.partner_name">
                <img
                  :src="
                    scope.row.partner_avatar
                      ? scope.row.partner_avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                />
                <span class="m-l-5">{{ scope.row.partner_name }}</span>
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="partner_mobile"
            label="手机号"
            min-width="150"
            size="small"
          >
          </el-table-column>
          <el-table-column label="所属导购" prop="worker_name" min-width="300">
            <template slot-scope="scope">
              <div class="nameBox" v-if="scope.row.worker_name">
                <img
                  :src="
                    scope.row.worker_avatar
                      ? scope.row.worker_avatar
                      : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
                  "
                  alt=""
                />
                <span class="m-l-5">{{ scope.row.worker_name }}</span>
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            label="待结算佣金(元)"
            min-width="210"
            prop="for_withdrawal"
            sortable="custom"
            align="right"
          >
            <div slot-scope="{ row }">
              {{ row.for_withdrawal / 100 }}
            </div>
          </el-table-column>
          <el-table-column
            label="可提现佣金(元)"
            min-width="210"
            prop="can_withdrawal"
            sortable="custom"
            align="right"
          >
            <div slot-scope="{ row }">
              {{ row.can_withdrawal / 100 }}
            </div>
          </el-table-column>
          <el-table-column
            label="提现中佣金(元)"
            min-width="210"
            prop="withdrawal_of"
            sortable="custom"
            align="right"
          >
            <div slot-scope="{ row }">
              {{ row.withdrawal_of / 100 }}
            </div>
          </el-table-column>
          <el-table-column
            label="已提现佣金(元)"
            min-width="210"
            prop="have_withdrawal"
            sortable="custom"
            align="right"
          >
            <div slot-scope="{ row }">
              {{ row.have_withdrawal / 100 }}
            </div>
          </el-table-column>
          <el-table-column label="注册时间" min-width="210" sortable="custom">
            <template slot-scope="scope">
              {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="状态" min-width="100" prop="partner_status">
            <template slot-scope="scope">
              <el-tag
                v-if="scope.row.partner_status === 1"
                type="success"
                size="small"
              >
                启用
              </el-tag>
              <el-tag
                v-if="scope.row.partner_status === 0"
                type="danger"
                size="small"
              >
                禁用
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            min-width="150"
            fixed="right"
            class="switch"
          >
            <template slot-scope="scope">
              <el-button
                @click.native.prevent="goDetail(scope.row)"
                type="text"
                class="operation-success"
              >
                查看明细
              </el-button>
              <span
                v-if="scope.row.partner_status === 0"
                @click="start(scope.row)"
                class="start"
              >
                启用
              </span>
              <span
                v-if="scope.row.partner_status === 1"
                class="stop"
                @click="stop(scope.row)"
              >
                禁用
              </span>
            </template>
          </el-table-column>
          <div slot="empty">
            <layout-empty
              v-if="!condition"
              :message="'暂无分销合伙人'"
            ></layout-empty>
            <layout-empty
              v-if="condition"
              empty_type="tag"
              :message="'没有找到相关分销合伙人'"
            ></layout-empty>
          </div>
        </el-table>

        <!-- 禁用弹框 -->
        <el-dialog
          title="禁用原因"
          width="500px"
          label-position="top"
          :visible.sync="dialog_visible"
          @close="closeDialog"
        >
          <div class="content">
            <el-form label-position="left" label-width="250">
              <el-form-item
                label="你确认要禁用该分销合伙人吗？请将禁用原因告知合伙人吧"
              >
                <el-input
                  placeholder="请输入禁用原因"
                  v-model="forbidden_reason"
                  type="textarea"
                  size="small"
                  maxlength="100"
                  :autosize="{ minRows: 6 }"
                  show-word-limit="show-word-limit"
                ></el-input>
              </el-form-item>
            </el-form>
            <div class="button">
              <el-button size="small" @click="closeDialog">取消</el-button>
              <el-button size="small" type="primary" @click="save">
                确认
              </el-button>
            </div>
          </div>
        </el-dialog>
      </template>

      <!-- 分页 -->
      <template slot="footer">
        <el-pagination
          v-if="isLoaded && list && list.length > 0"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 30, 40]"
          :current-page="meta.page"
          :page-size="meta.size"
          :total="meta.total"
          @current-change="handleCurrentPageChange"
          @size-change="handleSizeChange"
        ></el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import SelectMember from "@/components/sy-components/select-member.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import SearchSelectMember from "@/components/sy-components/search-select-member.vue";
import partnerApi from "@/api/wisdom/partner";
import UserModule from "@/store/modules/user";
import ShareModule from "@/shared/index";

@Component({
  components: {
    SelectMember,
    LayoutEmpty,
    SearchSelectMember
  }
})
export default class extends Vue {
  private search_type = 1; // 搜索类型 1合伙人姓名 2合伙人ID
  private search_types = [
    { name: "合伙人姓名", key: 1 },
    { name: "合伙人ID", key: 2 }
  ];
  private search_name = ""; // 搜索活动名称
  private operate_status: any = {};
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private updateTime = 0;
  private timerId = 0;
  private commission_type = ""; //佣金状态
  private min_input = ""; //最低佣金
  private max_input = ""; //最高佣金
  private partner_type = 0; //默认合伙人状态
  private partner_types = [
    { name: "全部", key: 0 },
    { name: "启用", key: 1 },
    { name: "禁用", key: 2 }
  ]; //合伙人状态
  private commission_types = [
    { name: "待结算佣金", key: 1 },
    { name: "可提现佣金", key: 2 },
    { name: "已提现佣金", key: 3 }
  ];
  private hign_choose: Array<any> = [];
  private isLoaded = false; // 数据是否加载完成
  private order_desc = ""; //排序方式
  private order_type = -1; //排序类型0注册时间1待提现2可提现 3已提现4提现中
  private sortArr = [
    "created_at",
    "for_withdrawal",
    "can_withdrawal",
    "have_withdrawal",
    "withdrawal_of"
  ];
  private list: Array<any> = []; // 列表数据
  private forbidden_reason = ""; //禁用原因
  private dialog_visible = false;
  private listId = [];
  private guide_workers: Array<any> = [];
  private is_show_sync = false;
  private condition = false;

  // 接口数据
  private member_show = false;

  // methods
  // 下拉框改变的时候
  private selectChange() {
    this.search_name = "";
  }

  // 合伙人启用操作
  private start(row: any) {
    const params = {
      id: row.id,
      switch: row.partner_status === 1 ? 2 : 1
    };
    partnerApi.partnerSwitchApi(params).then((res: any) => {
      this.partner_type = res.switch;
      if (res.code == 200) {
        this.$message.success("启用成功");
        this.getList();
      }
    });
  }

  // 合伙人禁用操作
  private stop(row: any) {
    this.dialog_visible = true;
    this.listId = row.id;
    this.operate_status = row;
  }
  private save() {
    if (!this.forbidden_reason) {
      this.$message.warning("请输入禁用原因");
      return;
    }
    const params = {
      action: "refuse",
      id: this.listId,
      switch: this.operate_status.partner_status === 1 ? 2 : 1,
      forbidden_reason: this.forbidden_reason
    };
    partnerApi.partnerSwitchApi(params).then(() => {
      this.$message.success("禁用成功!");
      this.closeDialog();
      this.getList();
    });
  }
  private closeDialog() {
    this.dialog_visible = false;
    this.forbidden_reason = "";
  }

  //刷新
  private reload() {
    this.list = [];
    this.reset();
    this.getList();
    this.updata();
  }
  private updata() {
    const date = new Date().getTime();
    //返回格式化后的日期
    this.updateTime = date / 1000;
    console.log(this.updateTime);
  }

  // 获取列表数据
  private getList() {
    const params = {
      search_type: this.search_type,
      search_value: this.search_name,
      partner_status: this.partner_type, // 合伙人状态
      commission: this.commission_type, //佣金状态
      commission_min: this.min_input, //最低佣金
      commission_max: this.max_input, //最高佣金
      worker_ids: this.guide_workers,
      orderby: this.order_desc,
      order_by_type: this.order_type,
      page: this.meta.page, //  页数
      size: this.meta.size // 请求数量
    };
    partnerApi.partnerListApi(params).then((res: any) => {
      console.log(res);
      const data = { ...res.data };
      this.list = data.list;
      const { total, page, size } = data;
      this.meta = { total, page, size };
      const time = Date.now() / 1000;
      this.updateTime = time;
      this.isLoaded = true;
    });
  }

  // 排序
  private sortChange(column: any) {
    // ascending 升序 descending 降序
    // 排序方式desc倒序 asc正序
    //order_type 排序类型0注册时间1待提现2可提现 3已提现4提现中
    if (column.order == null) {
      this.order_type = -1;
      this.order_desc = "";
    } else {
      this.order_type = this.sortArr.indexOf(column.prop);
      // ascending 升序 descending 降序
      this.order_desc = column.order === "ascending" ? "asc" : "desc";
    }
    this.getList();
  }

  // 每页显示条数
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }

  // 当前页
  private handleCurrentPageChange(val: number) {
    this.meta.page = val;
    this.getList();
  }

  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }

  // 员工输入框即高级搜索
  private highSearch() {
    this.member_show = true;
  }
  // 高级搜索完成
  private submitDialog(list: any) {
    this.hign_choose = list;
    this.guide_workers = list.map((v: any) => v.id);
    // 调用数据
    this.getList();
  }

  // 删除已选导购
  private removeGuide(val: any) {
    const idx = this.hign_choose.findIndex(item => item.id === val);
    idx > -1 && this.hign_choose.splice(idx, 1);
    this.guide_workers = this.hign_choose.map((v: any) => v.id);
    // 调用数据
    this.getList();
  }

  // 查询结算记录
  private search() {
    this.condition = true;
    this.meta.page = 1;
    (this.$refs["table"] as any).clearSort();
    if (
      this.commission_type != "" &&
      (this.min_input == "" || this.max_input == "")
    ) {
      this.$message.warning("请填写佣金范围");
      return;
    } else if (
      this.commission_type == "" &&
      (this.min_input != "" || this.max_input != "")
    ) {
      this.$message.warning("请选择佣金状态");
      return;
    }
    // 调用查询数据
    this.getList();
  }
  // 重置
  private reset() {
    this.search_type = 1;
    this.search_name = "";
    this.partner_type = 0;
    this.hign_choose = [];
    this.guide_workers = [];
    this.meta.page = 1;
    this.commission_type = "";
    (this.$refs["table"] as any).clearSort();
    this.min_input = "";
    this.max_input = "";
    this.order_desc = ""; //排序方式
    this.order_type = -1;
    // 调用查询数据
    this.getList();
  }
  // 导出excel
  private exportFile() {
    const params = {
      search_type: this.search_type,
      search_value: this.search_name,
      partner_status: this.partner_type, // 合伙人状态
      commission: this.commission_type, //佣金状态
      commission_min: this.min_input, //最低佣金
      commission_max: this.max_input, //最高佣金
      orderby: this.order_desc,
      order_by_type: this.order_type,
      page: this.meta.page, //  页数
      size: this.meta.size // 请求数量
    };
    partnerApi.partnerListExportApi(params).then(() => {
      this.$confirm(
        "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
        "文件下载提示",
        {
          distinguishCancelAndClose: true,
          confirmButtonText: "前往下载中心",
          cancelButtonText: "取消",
          type: "success"
        }
      ).then(() => {
        ShareModule.getShared().pushState("myapp_download-center");
      });
    });
  }

  // 去往合伙人明细
  private goDetail(params: any) {
    this.$router.push({
      name: "myapp_wisdom_sub-distribution_partner-manage_detail",
      query: {
        id: params.id
      }
    });
  }
  // 生命周期
  private mounted() {
    this.getList();
  }
}
</script>
<style lang="scss" scoped>
/deep/.search-item {
  display: inline-flex;
  align-items: center;
}

.search-item /deep/ .el-input--small .el-input__inner {
  height: 36px;
  line-height: 36px;
}

.search-item .search-select {
  width: 115px;
}

.search-item .name-input {
  width: 165px;
}

.operate-btn {
  margin: 20px 0 0;

  .reset {
    color: #303133;
    background-color: #fff;
    border-color: #c7cad0;
    margin-left: 15px;
  }
}

.name-input /deep/ .el-input__inner {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
  text-align: center;
}

.search-input /deep/ .el-input__inner {
  text-align: left;
}

.data-button {
  display: inline-block;
  margin-left: 20px;
  color: #666666;
  font-size: 12px;
}

.search-select /deep/ .el-input__inner {
  border-right: none;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.banlance-list /deep/ .buttons-content {
  display: none;
}

.banlance-list /deep/ .operate-content {
  display: none;
}

.banlance-list /deep/ .slide-box {
  display: none;
}

.banlance-list /deep/ .search-content {
  padding: 20px;
}

.button {
  text-align: right;
}

.el-select-dropdown__item.selected {
  font-weight: normal;
  height: auto !important;
}

.export-btn {
  margin-bottom: 20px;
}

.time {
  display: inline-block;
  margin-left: 20px;
}

.com-section {
  border: 1px solid #dcdfe6;
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
  background-color: #fff;
}

.com-section /deep/ .el-input {
  width: 77px;
  display: inline-flex;
}

.com-section /deep/ .el-input .el-input__inner {
  border: none;
  height: 34px;
}

.select-input /deep/ .el-input__inner:first-child {
  width: 120px;
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}

.clearbox {
  display: flex;

  .moneyStatus {
    background-color: $color-background;
    font-size: 14px;
    width: 203px;
    height: 55px;
    margin: 10px 20px 20px 0;
    padding: 10px;
    border-radius: 5px;

    span {
      display: block;
      font-size: 20px;
      font-weight: 700;
    }

    .spa {
      margin-bottom: 10px;
      font-weight: 400;
      font-size: 14px;
    }
  }
}

.nameBox {
  display: flex;
  align-items: center;

  img {
    width: 30px;
    height: 30px;
    border-radius: 50%;
  }

  span {
    width: 85px;
    flex-shrink: 0;
  }
}

/deep/.el-form-item__label {
  line-height: 21px;
  margin-bottom: 10px;
}

.start {
  cursor: pointer;
  color: #708df5;
}
.stop {
  cursor: pointer;
  color: #eb3638;
}

.start,
.stop {
  margin-left: 5px;
}

.el-scrollbar {
  .el-scrollbar__view {
    .el-select-dropdown__item {
      height: auto;
    }
  }
}

.el-select-dropdown__item.selected {
  font-weight: normal;
  height: auto !important;
}

.el-icon-refresh-right {
  color: #708df5;
  cursor: pointer;
}

/deep/.el-tree--highlight-current
  .el-tree-node.is-current
  > .el-tree-node__content {
  color: #708df5 !important;
}
</style>
