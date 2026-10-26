<template>
  <div class="editStratage">
    <layout-form>
      <template slot="content">
        <cy-form
          ref="cform"
          :rules="rules"
          :model="ruleForm"
          label-width="76px"
        >
          <p class="headline ">基础信息</p>
          <cy-form-item label="策略名称" prop="strategy_name">
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
          <cy-form-item label="执行方式" prop="time" class="form-time" required>
            <p>长期执行</p>
            <div>
              <span>计划执行时间</span>
              <el-date-picker
                class="data-picker"
                type="datetimerange"
                v-model="ruleForm.time"
                start-placeholder="请选择开始时间"
                end-placeholder="请选择结束时间"
                size="small"
                format="yyyy-MM-dd HH:mm:ss"
                value-format="timestamp"
                :picker-options="pickerOptions"
              ></el-date-picker>
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
              <el-radio
                v-model="radioType"
                label="1"
                @change="conditionChangeFn"
                >不设置条件</el-radio
              >
            </div>
            <div class="radio-height">
              <el-radio
                v-model="radioType"
                label="2"
                @change="conditionChangeFn"
                >购买金额</el-radio
              >
            </div>
            <cy-form-item
              label=""
              prop="price"
              v-if="radioType == 2"
              class="m-l-0"
              required
            >
              <cy-input
                size="medium"
                placeholder="请输入金额"
                v-model="ruleForm.price.minPrice"
                @input="
                  _checkFloat(
                    2,
                    ruleForm.price.minPrice,
                    ruleForm.price,
                    'minPrice'
                  )
                "
              >
                <template slot="prepend">最低</template>
                <template slot="append">元</template>
              </cy-input>
              <span class="line">——</span>
              <cy-input
                size="medium"
                placeholder="请输入金额"
                v-model="ruleForm.price.maxPrice"
                @input="
                  _checkFloat(
                    2,
                    ruleForm.price.maxPrice,
                    ruleForm.price,
                    'maxPrice'
                  )
                "
              >
                <template slot="prepend">最高</template>
                <template slot="append">元</template>
              </cy-input>
            </cy-form-item>
            <div class="radio-height">
              <el-radio
                v-model="radioType"
                label="3"
                @change="conditionChangeFn"
                >购买商品</el-radio
              >
            </div>
            <div class="sub-radio-content goods-msg" v-if="radioType == 3">
              <el-button type="primary" size="small" plain @click="openShops"
                >选择商品</el-button
              >
            </div>
            <div
              class="shopBoxContainer"
              v-if="selectedshopList.length && radioType == 3"
            >
              <selected-shops
                :list="selectedshopList"
                @getShopChange="delselectShops"
              ></selected-shops>
            </div>
          </div>
          <div class="splitLine"></div>
          <p class="headline">触达设置</p>
          <div class="edit-part">
            <dl>
              <dt>触达渠道</dt>
              <dd>
                <ul class="msg">
                  <li>短信</li>
                  <li>
                    加粉短信模板
                    <span class="setting" v-if="settingStatus == 5"
                      >(已设置)</span
                    >
                    <el-link
                      v-else
                      class="refresh"
                      type="primary"
                      :underline="false"
                      @click="refreshFn"
                      >更新模板</el-link
                    >
                    <el-link
                      class="link"
                      type="primary"
                      icon="el-icon-edit"
                      :underline="false"
                      @click="goSetting"
                      >{{ settingStatus == 5 ? "修改" : "完善" }}参数</el-link
                    >
                  </li>
                </ul>
              </dd>
            </dl>
            <dl>
              <dt>触达规则</dt>
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
          <cy-form-item label="提醒手机号" prop="phone" class="form-phone">
            <cy-input
              class="data-picker"
              v-model="ruleForm.phone"
              maxlength="11"
              minlength="0"
              placeholder="请输入提醒手机号"
              size="medium"
              @input="_checkPositive(ruleForm, 'phone', ruleForm.phone)"
            ></cy-input>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="saveData(1)" size="small"
          >保存并启用</el-button
        >
        <el-button type="primary" @click="saveData()" size="small"
          >保存</el-button
        >
        <el-button @click="goBackFn" size="small">取消</el-button>
      </template>
    </layout-form>
    <select-shop
      v-if="shopListDialog"
      @addShopsDialog="getSelectShops"
    ></select-shop>
    <el-dialog
      title="离开页面提示"
      :visible.sync="isCancelDialog"
      width="420px"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        你确认放弃编辑吗？放弃后编辑的信息将不会被保存</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goBack" size="small">取消</el-button>
        <el-button type="primary" @click="saveGoBackFn" size="small"
          >确认</el-button
        >
      </span>
    </el-dialog>
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
    phone: "", //提醒手机号码
    time: [],
    price: {
      minPrice: "",
      maxPrice: ""
    }
  };
  private settingStatus = 0; //短信设置
  private infoId = 0; //短信模板id
  private channels: Array<any> = []; //短信相关
  private customer_group_id = 1; //客户分组ID
  private strategy_type = 2; //策略类型: 1=客户策略，2=内置策略
  private strategy_method = 1; //执行方式: 1=长期执行，2=立即执行，3=定时执行
  private strategy_policy = 1; //触达规则：1=通知一个渠道，2=通知所有渠道
  private status = 0; //启用or停用状态
  private isCancelDialog = false; //取消按钮是否弹框
  // 检查活动时间
  private checkTime = (rule: any, value: any, callback: any) => {
    const date = new Date();
    if (!value || value.length == 0) {
      return callback(new Error("计划执行时间不能为空"));
    }

    if (value[0] < date.getTime()) {
      return callback(new Error("开始时间需要晚于当前时间"));
    }
    if (value[0] >= value[1]) {
      return callback(new Error("开始时间需要早于结束时间"));
    }
    callback();
  };
  // 检查价格
  private checkPrice = (rule: any, value: any, callback: any) => {
    const { minPrice, maxPrice } = value;
    if (Number(minPrice) >= Number(maxPrice)) {
      return callback(new Error("最低金额不可大于最高金额"));
    }
    callback();
  };
  private rules = {
    time: [{ validator: this.checkTime, trigger: "change" }],
    price: [{ validator: this.checkPrice, trigger: "change" }]
  };
  // 日期选择器今天之前的时间置灰
  private pickerOptions = {
    disabledDate(time: any) {
      return time.getTime() < Date.now() - 8.64e7;
    }
  };
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
            minPrice: data.start_point == 0 ? 0 : data.start_point / 100,
            maxPrice: data.end_point == 0 ? 0 : data.end_point / 100
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
          this.infoId = data.channels.length
            ? data.channels[0].junior_id
            : this.infoId;
          this.channels = data.channels;
          this.customer_group_id = data.customer_group_id;
          this.strategy_type = data.strategy_type;
          this.strategy_method = data.strategy_method;
          this.strategy_policy = data.strategy_policy;
          this.selectedshopList = data.relations.length ? data.relations : [];
          this.status = data.status;
          this.refreshFn();
        } else {
          this.$message.error(res.message);
        }
      }
    );
  }

  //取消-返回上一页
  private goBackFn() {
    this.isCancelDialog = true;
    // this.$router.push({
    //   name: "myapp_customer_auto-market_insert-market-strategy_list"
    // });
  }
  //确认取消
  private saveGoBackFn() {
    this.isCancelDialog = false;
    this.$router.push({
      name: "myapp_customer_auto-market_insert-market-strategy_list"
    });
  }
  //返回上一页
  private goBack() {
    this.isCancelDialog = false;
  }
  //条件设置单选框
  private conditionChangeFn(type: string) {
    this.radioType = type;
  }

  // 购买商品-选择商品
  openShops() {
    this.shopListDialog = true;
  }
  private ArrHasItem(arr: any, item: any) {
    const has: number = arr.findIndex((v: any) => {
      return v.id == item.id;
    });
    if (has > -1) {
      return true;
    }
    return false;
  }
  private getSelectShops(rows: any) {
    this.shopListDialog = !this.shopListDialog;
    if (rows) {
      rows.forEach((e: any) => {
        const has: boolean = this.ArrHasItem(this.selectedshopList, e);
        if (!has) {
          this.selectedshopList.push(e);
        }
      });
      // this.selectedshopList = this.selectedshopList.length;
    }
    console.log(this.selectedshopList);
  }
  // 删除商品
  delselectShops(result: any) {
    console.log(result);
    this.selectedshopList = result;
    // this.selectedShopNum = result.length;
  }
  // 跳转设置参数
  private goSetting() {
    const routeData = this.$router.resolve({
      path: "/myapp/customer/auto-market/sms-mgt/set-param",
      query: {
        id: String(this.infoId)
      }
    });
    window.open(routeData.href, "_blank");
  }
  private refreshFn() {
    API.smsTemplateInfoApi({ id: this.infoId }).then((res: any) => {
      if (res.code === 200) {
        if (res.data) {
          const { status } = res.data;
          this.settingStatus = status;
          //status	1：未审核、2：审核中、3：审核失败、4：审核成功、5：已赋值
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击保存
  private saveData(type: any) {
    const data = JSON.parse(JSON.stringify(this.ruleForm));
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        const params: any = {
          id: data.id,
          strategy_name: data.strategy_name,
          start_time: data.time[0] / 1000,
          end_time: data.time[1] / 1000,
          inner_condition: Number(this.radioType),
          channels: this.channels,
          customer_group_id: this.customer_group_id,
          strategy_type: this.strategy_type,
          strategy_method: this.strategy_method,
          strategy_policy: this.strategy_policy,
          is_finished: 2
        };
        params.status = type == 1 ? 2 : this.status;
        if (data.phone) {
          params.mobile = data.phone;
        }
        if (Number(this.radioType) == 2) {
          params.start_point = Number(data.price.minPrice) * 100;
          params.end_point = Number(data.price.maxPrice) * 100;
        }
        if (Number(this.radioType) == 3) {
          params.relations = this.selectedshopList.map((item: any) => {
            return { relation_id: item.id };
          });
          if (params.relations.length < 1) {
            this.$message.error("请选择商品！");
            return;
          }
        }
        if (this.settingStatus != 5) {
          this.$message.error("短信参数未完善！");
          return;
        }
        API.marketStrategyRecordEditSaveApi(params).then((res: any) => {
          if (res.code === 200) {
            this.$message.success("保存成功");
            this.$router.push({
              name: "myapp_customer_auto-market_insert-market-strategy_list"
            });
          } else {
            this.$message.error(res.message);
          }
        });
      } else {
        return;
      }
    });
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped>
.editStratage {
  .line {
    margin-left: 20px;
    margin-right: 20px;
  }
  .headline {
    color: #292a30;
    font-size: 18px;
    line-height: 26px;
    padding-bottom: 18px;
    font-size: 18px;
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
      margin-left: 26px;
      margin-right: 400px;
      margin-bottom: 20px;
    }
  }
  .edit-part {
    border-bottom: 1px solid $color-third-line;
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
        text-align: right;
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
  .data-picker {
    margin-left: 20px;
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
/deep/.m-l-0 .el-form-item__content {
  margin-left: 20px !important;
}
/deep/.form-time .el-form-item__content .error {
  margin-left: 103px !important;
}
/deep/.in-item {
  width: 100%;
  margin-bottom: 5px;
  .label {
    padding-left: 0 !important;
  }
}
/deep/.el-form-item__content {
  line-height: 40px !important;
}
/deep/.form-phone .el-form-item__label {
  width: 90px !important;
}
</style>
