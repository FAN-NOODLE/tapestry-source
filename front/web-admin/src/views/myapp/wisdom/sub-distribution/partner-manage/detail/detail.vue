<template>
  <layout-form v-if="isLoaded" class="detail">
    <template class="wrap" slot="content">
      <div class="settlementInfo mar-b-30">
        <el-row class="title base">
          <el-col :span="24">合伙人基本信息</el-col>
        </el-row>
        <el-row class="info" :gutter="20">
          <el-col class="mar-t-5 mess" :span="6">
            <span class="info-tag">合伙人</span>
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
            <span v-text="infoList.partner_name"></span>
          </el-col>
          <el-col class="" :span="6">
            <span class="info-tag">合伙人ID</span>
            <span v-text="infoList.partner_sn"></span>
          </el-col>
          <el-col class="" :span="6">
            <span class="info-tag">手机号</span>
            <span v-text="infoList.partner_mobile"></span>
          </el-col>
          <el-col :span="6">
            <span class="info-tag">所属导购</span>
            <img
              :src="
                infoList.worker_avatar
                  ? infoList.worker_avatar
                  : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
              "
              alt=""
              width="30"
              height="30"
              class="table-text-img"
            />
            <span v-text="infoList.worker_name"></span>
          </el-col>
          <el-col class="" :span="6">
            <span class="info-tag">注册时间</span>
            <span>{{
              infoList.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span>
          </el-col>
          <el-col class="" :span="6">
            <span class="info-tag">状态</span>
            <el-tag
              v-if="infoList.partner_status === 1"
              type="success"
              size="small"
            >
              启用
            </el-tag>
            <el-tag
              v-if="infoList.partner_status === 0"
              type="danger"
              size="small"
            >
              禁用
            </el-tag>
          </el-col>
          <el-col
            class="forbidden"
            :span="12"
            v-if="infoList.partner_status === 0"
          >
            <span class="info-tag">禁用原因</span>
            <span v-text="infoList.forbidden_reason"></span>
          </el-col>
        </el-row>
        <div class="line"></div>
      </div>
      <div class="distributionGoods">
        <div class="title mar-b-10">合伙人余额明细</div>
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
          class="mar-t-15 list-last-line-none min-height-100 m-b-20"
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
                  type="warning"
                  size="small"
                >
                  待结算
                </el-tag>
                <el-tag
                  v-if="scope.row.settlement_status === 1"
                  type="success"
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
          <el-table-column label="订单状态" min-width="100">
            <template slot-scope="scope">
              <el-tag
                v-if="scope.row.order_status === 20"
                type="info"
                size="small"
              >
                未支付
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 21"
                type="info"
                size="small"
              >
                已取消
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 23"
                type="info"
                size="small"
              >
                待支付
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 25"
                type="info"
                size="small"
              >
                支付中
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 30"
                type="info"
                size="small"
              >
                待确认
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 40"
                type="info"
                size="small"
              >
                待发货
              </el-tag>
              <el-tag
                v-if="
                  scope.row.order_status === 41 &&
                    scope.row.order_status === 61 &&
                    scope.row.order_status === 81
                "
                type="info"
                size="small"
              >
                售后中(待发货,仅退款)
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 60"
                type="text"
                size="small"
              >
                待收货
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 80"
                type="info"
                size="small"
              >
                已签收
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 83"
                type="info"
                size="small"
              >
                部分退款
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 100"
                type="info"
                size="small"
              >
                已完成
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 200"
                type="info"
                size="small"
              >
                退款中
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 201"
                type="info"
                size="small"
              >
                退款失败
              </el-tag>
              <el-tag
                v-if="scope.row.order_status === 205"
                type="info"
                size="small"
              >
                已退款
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="订单支付时间" min-width="160">
            <template slot-scope="scope">
              <span>{{
                scope.row.order_pay_time | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
            </template>
          </el-table-column>
          <el-table-column
            label="分销合伙人佣金（元）"
            min-width="160"
            align="right"
          >
            <template slot-scope="scope">
              <span>{{ scope.row.partner_commission / 100 }}</span>
            </template>
          </el-table-column>
          <div slot="empty">
            <layout-empty :message="'还没有余额记录哦~'"></layout-empty>
          </div>
        </el-table>
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
      </div>
    </template>
    <template slot="footer">
      <!-- <div class="line down-line"></div> -->
      <el-button size="mini" @click="backList()" slot="back">返回 </el-button>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import partnerApi from "@/api/wisdom/partner";
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
  private detailList = []; //合伙人明细
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
  private backList() {
    this.$router.go(-1);
  }

  // 获取明细数据
  private getdetailList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      id: this.$route.query.id
    };
    console.log(params, "明细参数互相学习");

    partnerApi.partnerDetailApi(params).then(res => {
      console.log(res, "合伙人明细休息休息");
      this.detailList = res.data.list;
      this.infoList = res.data.partner;
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
      page: this.meta.page,
      size: this.meta.size,
      id: this.$route.query.id
    };
    partnerApi.partnerDetailExportApi(params).then(() => {
      this.$confirm(
        "文件已进入下载列表，请前往【下载中心】查看及下载文件。",
        "文件下载提示",
        {
          distinguishCancelAndClose: true,
          confirmButtonText: "前往下载中心",
          cancelButtonText: "取消",
          type: "info"
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

.detail /deep/ .footer-content .el-pagination {
  float: none;
}

#crs-isg .el-table_row td {
  border-bottom: 0px !important;
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
// .forbidden {
//   width: ;
// }

.goBack,
.base {
  display: inline-block;
  vertical-align: middle;
}

.base {
  margin-top: -7px;
}

.info-tag {
  color: #909399;
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
  span {
    margin-right: 10px;
    vertical-align: text-bottom;
  }
}
.info .el-col {
  margin-bottom: 20px;
}

.info .mess {
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.line {
  width: 100%;
  height: 1px;
  margin: 0 0 20px;
  background-color: #dcdfe6;
}

.down-line {
  margin: 20px 0;
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
.m-b-20 {
  margin-bottom: 20px;
}
</style>
