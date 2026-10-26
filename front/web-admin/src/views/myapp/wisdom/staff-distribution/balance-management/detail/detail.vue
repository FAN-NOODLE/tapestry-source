<template>
  <layout-list v-if="isLoaded" class="detail">
    <template class="wrap" slot="tabs">
      <div class="settlementInfo mar-b-30">
        <el-row class="title base">
          <el-col :span="24">导购基本信息</el-col>
        </el-row>
        <el-row class="info" :gutter="20">
          <el-col class="mar-t-5 mess">
            <span class="staf-title">导购</span>
            <img
              :src="
                infoList.avatar
                  ? infoList.avatar
                  : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
              "
              alt=""
              width="30"
              height="30"
              class="table-text-img"
            />
            <span v-text="infoList.name"></span>
          </el-col>
          <el-col v-if="USER.version !== 3">
            <span class="staf-title">导购所在门店</span>
            <el-popover
              v-if="infoList.store_name && infoList.store_name.length > 20"
              ,placement="top"
              trigger="hover"
              width="150"
            >
              <div class="skuname character-hiding-w" slot="reference">
                {{ infoList.store_name }}
              </div>
              <span>{{ infoList.store_name }}</span>
            </el-popover>
            <span v-else>{{ infoList.store_name || "-" }}</span>
          </el-col>
          <el-col>
            <span class="staf-title">导购所在部门</span>
            <el-tooltip
              v-if="tooltip"
              effect="dark"
              :content="infoList.dpName"
              placement="top"
            ></el-tooltip>
            <span>{{ infoList.dpName || "-" }}</span>
          </el-col>
        </el-row>
        <div class="line"></div>
      </div>
      <div class="distributionGoods">
        <div class="title mar-b-10">导购余额明细</div>
        <el-button
          :disabled="detailList.length == 0"
          size="small"
          @click="exportList"
          class="export"
          >导出Excel
        </el-button>
        <el-table
          :data="detailList"
          stripe
          class="mar-t-15 list-last-line-none min-height-100"
          :header-cell-style="{
            fontSize: '14px',
            color: '#606266',
            fontWeight: '400',
            padding: '11px 0 10px 0'
          }"
        >
          <el-table-column
            prop="id"
            label="结算单号"
            min-width="150"
          ></el-table-column>
          <el-table-column label="结算状态" min-width="120">
            <template slot-scope="scope">
              <span>
                <el-tag
                  v-if="scope.row.settlement_status === 0"
                  type="success"
                  size="small"
                >
                  待结算
                </el-tag>
                <el-tag
                  v-if="scope.row.settlement_status === 1"
                  type="info"
                  size="small"
                >
                  已结算
                </el-tag>
                <el-tag
                  v-if="scope.row.settlement_status === 2"
                  type="info"
                  size="small"
                >
                  已取消
                </el-tag>
              </span>
            </template>
          </el-table-column>
          <el-table-column
            label="订单编号"
            min-width="150"
            prop="order_no"
          ></el-table-column>
          <el-table-column label="订单状态" min-width="120">
            <template slot-scope="scope">
              <span
                v-if="scope.row.order_status === 20"
                type="success"
                size="small"
              >
                未支付
              </span>
              <span
                v-if="scope.row.order_status === 21"
                type="success"
                size="small"
              >
                已取消
              </span>
              <span
                v-if="scope.row.order_status === 23"
                type="success"
                size="small"
              >
                待支付
              </span>
              <span
                v-if="scope.row.order_status === 25"
                type="success"
                size="small"
              >
                支付中
              </span>
              <span
                v-if="scope.row.order_status === 30"
                type="success"
                size="small"
              >
                待确认
              </span>
              <span
                v-if="scope.row.order_status === 40"
                type="success"
                size="small"
              >
                待发货
              </span>
              <span
                v-if="
                  scope.row.order_status === 41 &&
                    scope.row.order_status === 61 &&
                    scope.row.order_status === 81
                "
                type="success"
                size="small"
              >
                售后中(待发货,仅退款)
              </span>
              <span
                v-if="scope.row.order_status === 60"
                type="text"
                size="small"
              >
                待收货
              </span>
              <span
                v-if="scope.row.order_status === 80"
                type="success"
                size="small"
              >
                已签收
              </span>
              <span
                v-if="scope.row.order_status === 83"
                type="success"
                size="small"
              >
                部分退款
              </span>
              <span
                v-if="scope.row.order_status === 100"
                type="success"
                size="small"
              >
                已完成
              </span>
              <span
                v-if="scope.row.order_status === 200"
                type="success"
                size="small"
              >
                退款中
              </span>
              <span
                v-if="scope.row.order_status === 201"
                type="success"
                size="small"
              >
                退款失败
              </span>
              <span
                v-if="scope.row.order_status === 205"
                type="success"
                size="small"
              >
                已退款
              </span>
            </template>
          </el-table-column>
          <el-table-column label="订单支付时间" min-width="130">
            <template slot-scope="scope">
              <span>{{
                scope.row.order_pay_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
            </template>
          </el-table-column>
          <el-table-column label="推广佣金（元）" min-width="140">
            <template slot-scope="scope">
              <span>{{ scope.row.promotion_commission / 100 }}</span>
            </template>
          </el-table-column>
          <el-table-column
            v-if="USER.version !== 3"
            label="服务佣金（元）"
            min-width="140"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.service_commission / 100 }}</span>
            </template>
          </el-table-column>
          <div slot="empty">
            <layout-empty :message="'暂无余额记录'"></layout-empty>
          </div>
        </el-table>
      </div>
    </template>
    <template slot="footer">
      <el-pagination
        v-if="detailList.length !== 0"
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
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/balanceManage";
import LayoutForm from "@/components/public/layout-form/index.vue";
import LayoutEmpty from "@/components/sy-components/layout-empty.vue";
import LayoutList from "@/components/public/layout-list/index.vue";
import UserModule from "@/store/modules/user";
import ShareModule from "@/shared/index";
@Component({
  components: {
    LayoutForm,
    LayoutEmpty,
    LayoutList
  }
})
export default class extends Vue {
  // data数据
  private detailList = []; //余额明细
  private infoList: object = {}; //导购基本信息
  private tooltip = false;
  private meta = {
    total: 40, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private isLoaded = false; // 数据是否加载完成

  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }

  // methods
  // 返回列表页
  private goBack() {
    this.$router.go(-1);
  }

  // 获取明细数据
  private getdetailList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      worker_id: this.$route.query.id
    };
    Api.balanceDetailApi(params).then(res => {
      this.detailList = res.data.list;
      this.infoList = res.data.member;
      const { total, page, size } = res.data;
      this.meta.page = page;
      this.meta.size = size;
      if (page === 1) this.meta.total = total;
      this.isLoaded = true;
    });
  }
  // 每页数据
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getdetailList();
  }
  // 第几页
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getdetailList();
  }
  // 导出列表
  private exportList() {
    const params = {
      worker_id: this.$route.query.id,
      version: this.USER.version
    };
    Api.balanceDetailExportApi(params).then(() => {
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

  // 生命周期
  private mounted() {
    this.getdetailList();
  }
}
</script>
<style lang="scss" scoped>
.detail {
  color: #303133;
}
.mar-t-5 {
  margin-top: -3px;
}

.skuname {
  width: 200px;
  display: inline-block;
  vertical-align: text-bottom;
}

.clearfix:after {
  content: "";
  display: block;
  height: 0;
  clear: both;
  visibility: hidden;
}

.clearfix {
  zoom: 1 span {
    float: left;
    padding: 0 5px;
    line-height: 20px;
    color: #409eff;
    border: 1px solid #409eff;
    border-radius: 3px;
    background-color: rgba(64, 158, 255, 0.1);
    margin: 5px;
  }
}

.layout-list {
  padding: 20px;
}

.table-text-img {
  border-radius: 100%;
  vertical-align: middle;
  margin-right: 4px;
}

.goBack {
  font-size: 14px;
  font-weight: normal;
  cursor: pointer;
  margin-right: 20px;
  padding-right: 20px;
  border-right: 1px solid #bfbfbf;
  transform: translateY(-5px);
}

.goBack,
.base {
  display: inline-block;
  vertical-align: middle;
}

.staf-title {
  color: #606266;
}

.base {
  margin-top: -7px;
}

/deep/.content-block {
  padding: 20px;
}

/deep/.search-content {
  display: none;
}

/deep/.buttons-content {
  display: none;
}

.detail /deep/ .operate-content {
  display: none;
}

.detail /deep/ .slide-box {
  display: none;
}

.title {
  font-weight: 700;
  font-size: 18px;
  line-height: 26px;
}

.export {
  margin: 10px 0 20px;
}

.info {
  margin-top: 20px;
  font-size: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;

  span {
    margin-right: 10px;
    display: inline-block;
    vertical-align: text-bottom;
  }
}

.info .mess {
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.mess {
  width: 50%;
}

.line {
  width: 100%;
  height: 1px;
  margin: 20px 0;
  background-color: #dcdfe6;
}

.imgBox {
  display: inline-block;
  vertical-align: middle;
}

.nameBox {
  display: flex;

  align-items: center img {
    width: 30px;
    height: 30px;
    border-radius: 50%;
  }

  span {
    width: 85px;
    flex-shrink: 0;
  }
}

.settlementBox {
  display: flex;

  align-items: center img {
    width: 45px;
    height: 45px;
  }

  span {
    color: #909399;
  }
}

.page-right {
  text-align: right;
}
</style>
