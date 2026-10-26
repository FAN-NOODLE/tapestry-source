<template>
  <div class="detailStratage">
    <layout-form>
      <template slot="content">
        <cy-form ref="cform" label-width="76px" :rules="rules">
          <p class="headline ">基础信息</p>
          <cy-form-item label="策略名称">
            {{ ruleForm.strategy_name || "加粉短信" }}
          </cy-form-item>
          <cy-form-item label="策略人群">
            成交且未添加好友
            <el-popover
              placement="top-start"
              trigger="hover"
              content="已经成交且未添加导购企业微信好友的顾客，不包含历史成交客户"
            >
              <svg-icon
                slot="reference"
                class="help-icon"
                icon-class="group-help"
              ></svg-icon>
            </el-popover>
          </cy-form-item>
          <cy-form-item label="执行方式">
            <p>长期执行</p>
            <div>
              计划执行时间
              {{
                (ruleForm.time[0] / 1000) | formatDate("yyyy-mm-dd hh:mm:ss")
              }}
              -
              {{
                (ruleForm.time[1] / 1000) | formatDate("yyyy-mm-dd hh:mm:ss")
              }}
            </div>
          </cy-form-item>
          <div class="splitLine"></div>
          <p class="headline">
            条件设置<span class="mark"
              >作为策略人群的条件补充，例如：成交金额50-200，则表示对成交金额50-200的成交且未添好友的客户发送加粉短信</span
            >
          </p>
          <div class="condition-setting">
            <div class="radio-height">
              <el-radio v-model="radioType" label="1" disabled
                >不设置条件</el-radio
              >
            </div>
            <div class="radio-height">
              <el-radio v-model="radioType" label="2" disabled
                >购买金额</el-radio
              >
              <span v-if="radioType == 2"
                >最低{{ ruleForm.price.minPrice }}元-最高{{
                  ruleForm.price.maxPrice
                }}元</span
              >
            </div>
            <div class="radio-height">
              <el-radio v-model="radioType" label="3" disabled
                >购买商品</el-radio
              >
            </div>
            <div
              class="shopBoxContainer"
              v-if="selectedshopList.length && radioType == 3"
            >
              <el-table :data="selectedshopList" max-height="500">
                <el-table-column prop="id" label="SPU ID"></el-table-column>
                <el-table-column prop="goods" label="商品">
                  <template slot-scope="scope">
                    <div class="aligns-center">
                      <!--图片-->
                      <el-image fit="cover"
                        style="width: 50px; height: 50px"
                        :src="
                          scope.row.spu_source ? scope.row.spu_source.url : ''
                        "
                      ></el-image>
                      <el-tooltip
                        :disabled="
                          scope.row.spu_name.length > 33 ? false : true
                        "
                        :content="scope.row.spu_name"
                        placement="bottom"
                        effect="light"
                      >
                        <span class="m-l-10 goods-name">{{
                          scope.row.spu_name
                        }}</span>
                      </el-tooltip>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
          <div class="splitLine"></div>
          <p class="headline">触达设置</p>
          <div class="edit-part">
            <dl>
              <dt class="text-normal">触达渠道</dt>
              <dd>
                <ul class="msg">
                  <li>短信</li>
                  <li>
                    加粉短信模板
                    <span class="setting">(已设置)</span>
                  </li>
                </ul>
              </dd>
            </dl>
            <dl>
              <dt class="text-normal">触达规则</dt>
              <dd>通过短信渠道触达</dd>
            </dl>
          </div>
          <div class="splitLine"></div>
          <p class="headline">
            异常提醒
            <span class="mark"
              >策略将于24小时之后到期或短信条数不足100条的情况下通过短信提醒</span
            >
          </p>
          <cy-form-item label="提醒手机号" phone class="form-phone">
            <span class="m-l-10">{{
              ruleForm.phone ? ruleForm.phone : "未填写"
            }}</span>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="goEdit()">编辑</el-button>
        <el-button @click="goBackFn">返回</el-button>
      </template>
    </layout-form>
    <select-shop
      v-if="shopListDialog"
      @addShopsDialog="getSelectShops"
    ></select-shop>
  </div>
</template>
<script lang="ts">
import { Component, Mixins } from "vue-property-decorator";
import selectShop from "../../components/select-shop.vue";
import selectedShops from "../../components/selected-shops.vue";
import CheckNum from "@/mixins/NumberInput";
import API from "@/api/customer";

@Component({
  components: {
    selectShop,
    selectedShops
  }
})
export default class extends Mixins(CheckNum) {
  /*************************** data start ***************************/
  private radioType = "1"; //radio 类型
  private isForbidden = false; //是否禁用
  private shopListDialog = false; //选择商品弹框是否展示
  private selectedshopList: any = []; //选择的商品列表
  private ruleForm = {
    strategy_name: "",
    phone: null, //提醒手机号码
    time: [],
    price: {
      minPrice: 0,
      maxPrice: 0
    }
  };
  private settingStatus = 0; //短信设置
  private infoId = 0; //短信模板id
  private channels: Array<any> = []; //短信相关
  private customer_group_id = 1; //客户分组ID
  private strategy_type = 2; //策略类型: 1=客户策略，2=内置策略
  private strategy_method = 1; //执行方式: 1=长期执行，2=立即执行，3=定时执行
  private strategy_policy = 1; //触达规则：1=通知一个渠道，2=通知所有渠道
  private rules = {};

  /*************************** method start ***************************/
  // 详情
  private getDetail() {
    API.marketStrategyDetailApi({ id: Number(this.$route.query.id) }).then(
      (res: any) => {
        if (res.code === 200) {
          const data = res.data;
          data.time = data.start_time
            ? [data.start_time * 1000, data.end_time * 1000]
            : [];
          this.radioType = String(data.inner_condition);
          data.price = {
            minPrice: data.start_point / 100,
            maxPrice: data.end_point / 100
          };
          data.phone = data.mobile;
          data.relations.map((item: any) => {
            item.spu_source = {
              url: item.relation_detail.img_url
            };
            item.spu_name = item.relation_detail.name;
            item.id = item.relation_id;
          });
          this.ruleForm = data;
          this.infoId = data.channels.length ? data.channels[0].junior_id : 0;
          this.channels = data.channels;
          this.customer_group_id = data.customer_group_id;
          this.strategy_type = data.strategy_type;
          this.strategy_method = data.strategy_method;
          this.strategy_policy = data.strategy_policy;
          this.selectedshopList = data.relations.length ? data.relations : [];
          console.log(data.relations);
        } else {
          this.$message.error(res.message);
        }
      }
    );
  }

  //取消-返回上一页
  private goBackFn() {
    this.$router.go(-1);
  }
  private goEdit() {
    this.$router.push({
      name: "myapp_customer_auto-market_insert-market-strategy_edit",
      query: { id: this.$route.query.id }
    });
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped>
.detailStratage {
  .aligns-center {
    display: flex;
    align-items: center;
    .goods-name {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }
  .line {
    margin-left: 20px;
    margin-right: 20px;
  }
  .headline {
    color: #292a30;
    font-size: 18px;
    line-height: 26px;
    padding-bottom: 18px;
    font-family: PingFangSC-Semibold, PingFang SC;
    font-weight: 600;
    .mark {
      color: #909399;
      font-size: 14px;
      font-weight: normal;
      margin-left: 10px;
    }
  }
  .splitLine {
    height: 0;
    border-top: 1px solid #ebeef5;
    margin-bottom: 20px;
  }
  .condition-setting {
    .radio-height {
      line-height: 20px;
      margin-bottom: 20px;
    }
    .sub-radio-content {
      padding-left: 26px;
      padding-bottom: 20px;
    }
    .goods-price {
      width: 500px;
      display: flex;
      align-items: center;
      dd {
        flex: 1;
      }
      dd:nth-child(2) {
        flex: 0 0 80px;
        text-align: center;
      }
    }
    .shopBoxContainer {
      margin-left: 50px;
      margin-right: 550px;
      margin-bottom: 20px;
    }
  }
  .edit-part {
    dl {
      display: flex;
      dt,
      dd {
        line-height: 36px;
        margin-bottom: 30px;
      }
      dt {
        flex: 0 0 74px;
        margin-right: 20px;
        text-align: left;
      }
      dd {
        flex: 1;
      }
      dt.mark {
        position: relative;
      }
      dt.mark::before {
        position: absolute;
        width: 10px;
        height: 10px;
        line-height: 10px;
        content: "*";
        color: red;
        left: 6px;
        top: 13px;
      }
      .msg {
        width: 500px;
        border: 1px solid $color-second-line;
        color: $color-text-normal;
        padding: 20px;
        li {
          line-height: 20px;
          overflow: hidden;
          position: relative;
          .setting {
            color: #48ba9b;
          }
          .refresh {
            position: absolute;
            right: 90px;
            top: 0;
          }
          .link {
            position: absolute;
            right: 0;
            top: 0;
            color: $color-main;
            text-decoration: none;
          }
        }
        li:nth-child(1) {
          font-weight: bold;
          margin-bottom: 20px;
        }
      }
    }
  }
  .edit-part:last-child {
    border-bottom: none;
  }
}
li {
  list-style: none;
}
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
//覆盖日历高度
/deep/.el-range-editor--small.el-input__inner {
  height: 36px;
}
//button 字体
/deep/.el-button--small {
  font-size: 14px;
}
//覆盖input长度
/deep/.cy-input-medium {
  width: 100% !important;
}
/deep/.m-l-0 .cy-form-item__content {
  margin-left: 20px !important;
}
/deep/.el-form-item__content {
  line-height: 40px !important;
}
/deep/.in-item {
  margin-bottom: 5px;
  .label {
    padding-left: 0 !important;
  }
}
/deep/.form-phone .el-form-item__label {
  width: 90px !important;
  padding-left: 0 !important;
}
.text-normal {
  color: $color-text-normal;
}
</style>
