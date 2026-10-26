<template>
  <div class="sub-setting">
    <layout-form>
      <template slot="content">
        <el-form :rules="rules" ref="form" :model="form" @submit.native.prevent>
          <el-form-item class="open-switch">
            <h3 class="h3">全局设置</h3>
            <div class="set-model">
              <span class="sub-title">二级分销</span>
              <el-switch
                v-model="form.partner_status"
                active-color="#13ce66"
                :active-value="1"
                :inactive-value="0"
                inactive-color="#f0f0f2"
                class="switch"
                @change="judgeSet"
              >
              </el-switch>
              <span v-if="form.partner_status" class="sub-title">开启</span>
              <span v-if="!form.partner_status" class="sub-title">关闭</span>
              <p class="desc">开启后在导购端展示“二级分销”功能入口</p>
            </div>
            <el-dialog
              title="未开启裂变分享计算业绩"
              :visible.sync="dialogVisible"
              width="30%"
            >
              <span>你需要在导购业绩设置页开启该功能</span>
              <span slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="goSet">去开启</el-button>
              </span>
            </el-dialog>
          </el-form-item>
          <el-form-item class="commission-set" v-model="form.commission_rate">
            <div class="h3 flex-2">
              <h3>二级佣金设置</h3>
              <div class="desc-2">
                分销合伙人直接绑定的顾客下单，则合伙人获得直接佣金，所属导购获得间接佣金
              </div>
            </div>

            <div class="set-model">
              <span class="sub-title method">分佣方式</span>
              <el-radio-group v-model="form.commission_type">
                <div class="base-set">
                  <el-radio
                    size="small"
                    class="radio"
                    :label="1"
                    @change="change_type()"
                  >
                    <span>基于分销商品的实付价格</span>
                    <p class="desc choice-desc">
                      企业付出的佣金=推广佣金+间接佣金
                    </p>
                  </el-radio>
                  <div v-if="form.commission_type === 1" class="set-info">
                    <p class="set-des">
                      <span class="sub-title">分销合伙人佣金</span>
                      <span class="explain"
                        >分销商品实付价格 × 推广佣金比例</span
                      >
                    </p>
                    <p class="guide-input">
                      <el-form-item
                        prop="commission_rate"
                        label="合伙人所属导购佣金"
                      >
                        <el-input
                          width-size="large"
                          class="search-input"
                          v-model="form.commission_rate"
                          placeholder="请输入数字"
                          @keyup.native="limitInput(form, 'commission_rate')"
                        >
                          <template slot="prepend">分销商品实付价格的</template>
                          <div slot="append">%</div>
                        </el-input>
                      </el-form-item>
                    </p>
                  </div>
                </div>
                <div class="base-bill">
                  <el-radio
                    size="small"
                    class="radio"
                    :label="2"
                    @change="change_type()"
                  >
                    <span>基于分销商品的推广佣金</span>
                    <p class="desc choice-desc">企业付出的佣金=推广佣金</p>
                  </el-radio>
                  <div v-if="form.commission_type === 2" class="set-info">
                    <!-- <p class="set-des">
                      <span class="sub-title count">计算公式：</span>
                      <span class="count-info">导购提成+合伙人佣金=推广佣金（全员分销商品里设置的推广佣金）</span>
                    </p> -->
                    <p class="set-input">
                      <el-form-item
                        prop="commission_rate"
                        label="分销合伙人佣金"
                      >
                        <el-input
                          width-size="large"
                          class="search-input"
                          v-model="form.commission_rate"
                          placeholder="请输入数字"
                          @keyup.native="limitInput(form, 'commission_rate')"
                        >
                          <template slot="prepend">推广佣金的</template>
                          <div slot="append">%</div>
                        </el-input>
                      </el-form-item>
                    </p>
                    <p class="deduct">
                      <span class="sub-title count-deduct"
                        >合伙人所属导购佣金</span
                      >
                      <span class="explain"
                        >推广佣金的{{
                          100 - form.commission_rate || "-- "
                        }}%</span
                      >
                    </p>
                  </div>
                </div>
              </el-radio-group>
            </div>
          </el-form-item>
          <el-form-item class="withdraw-switch">
            <h3 class="h3">提现设置</h3>
            <div class="set-model withdraw">
              <span class="sub-title">分销合伙人可提现</span>
              <el-switch
                v-model="form.withdraw_status"
                active-color="#13ce66"
                inactive-color="#f0f0f2"
                :active-value="1"
                :inactive-value="0"
                class="switch"
              >
              </el-switch>
              <span v-if="form.withdraw_status" class="sub-title">开启</span>
              <span v-if="!form.withdraw_status" class="sub-title">关闭</span>
              <p class="desc">开启后分销合伙人可进行提现操作</p>
            </div>
          </el-form-item>
          <h3 class="h3">二级分销公告设置</h3>
          <el-form-item
            class="withdraw-switch"
            label="公告标题"
            prop="distribution_title"
            required
          >
            <cy-input
              type="text"
              placeholder="请输入公告标题"
              v-model.trim="form.distribution_title"
              clearable
              width-size="medium"
              maxlength="30"
              :show-word-limit="true"
            ></cy-input>
          </el-form-item>
          <el-form-item
            class="withdraw-switch"
            label="公告内容"
            prop="distribution_text"
            required
          >
            <cy-input
              type="textarea"
              placeholder="请输入公告内容"
              v-model.trim="form.distribution_text"
              :rows="6"
              width-size="medium"
              maxlength="200"
              :show-word-limit="true"
            ></cy-input>
          </el-form-item>
        </el-form>
      </template>
      <template slot="footer">
        <el-button type="primary" size="small" @click="saveData">
          保存
        </el-button>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Mixins } from "vue-property-decorator";
import subDistributionApi from "@/api/wisdom/SubDistribution";
import NumberInput from "@/mixins/NumberInput";

@Component
export default class extends Mixins(NumberInput) {
  private subsetting = false; //二级分销设置开关
  private withdraw = true; //分销合伙人提现开关
  private form: any = {
    partner_status: 0, // 二级分销是否开启  0 关闭 1开启 默认0
    commission_type: 1, //1 基于商品分销 2 基于推广佣金总额的分账
    commission_rate: "", //佣金(0-100)
    withdraw_status: 1, // 0:不可提现 1:可以提现
    distribution_title: "", // 分销公告
    distribution_text: "" // 分销描述
  };
  private base_set = true; // 基于类别的显示
  private dialogVisible = false;
  private fission_type: number | string = ""; //裂变分享计算业绩 1:不计算业绩 2:计算业绩

  // methods
  // 单选框改变时清空输入框的值
  private change_type() {
    this.form.commission_rate = "";
  }
  // 设置开关时判断裂变分享是否开启
  private judgeSet() {
    if (this.fission_type == 1 && this.form.partner_status == 1) {
      this.dialogVisible = !this.dialogVisible;
      this.form.partner_status = 0;
    }
  }
  // 去导购业绩设置页面
  private goSet() {
    this.$router.push({
      name: "myapp_wisdom_guide-manage_guide-set"
    });
  }

  // 校验
  // 限制小数点位数2与值范围0~100
  private limitInput(row: any, key: any) {
    const value = row[key];
    row[key] = value
      .replace(/[^\d.]/g, "")
      .replace(/^(\d+)\.(\d\d).*$/, "$1.$2");
  }

  // 校验输入数字范围
  private Input_rule(rule: any, value: any, callback: any) {
    if (value > 100 || value < 0) {
      return callback(new Error("请输入0-100以内的数字"));
    } else {
      callback();
    }
  }
  private rules = {
    distribution_title: [{ required: true, message: "请输入公告标题" }],
    distribution_text: [{ required: true, message: "请输入公告内容" }],
    commission_rate: [{ validator: this.Input_rule }]
  };

  // 二级分销查看
  private getSubdistribution() {
    subDistributionApi.subDistributionCheckApi().then((res: any) => {
      this.form.partner_status = res.data.partner_status || 0;
      this.form.commission_type = res.data.commission_type || 1;
      this.form.withdraw_status =
        res.data.withdraw_status || res.data.withdraw_status == 0
          ? res.data.withdraw_status
          : 1;
      this.form.commission_rate = res.data.commission_rate / 100 || "";
      this.fission_type = res.data.fission_type;
      this.form.distribution_title = res.data.distribution_title;
      this.form.distribution_text = res.data.distribution_text;
    });
  }
  // 二级分销设置
  // 保存前校验表单
  private saveData() {
    (this.$refs["form"] as any).validate((valid: any) => {
      if (valid) {
        if (
          this.form.commission_rate == "" ||
          this.form.commission_rate == undefined
        ) {
          this.$message.warning("二级佣金未设置");
          return false;
        }
        const params = {
          partner_status: this.form.partner_status,
          commission_type: this.form.commission_type,
          withdraw_status: this.form.withdraw_status,
          commission_rate: this.form.commission_rate * 100,
          distribution_title: this.form.distribution_title,
          distribution_text: this.form.distribution_text
        };
        this.saveApi(params);
      }
    });
  }
  // 调用保存数据接口
  private saveApi(params: any) {
    subDistributionApi["subDistributionSettingApi"](params).then((res: any) => {
      this.$message.success("保存成功");
    });
  }
  // 生命周期
  private mounted() {
    this.getSubdistribution();
  }
}
</script>
<style lang="scss" scoped>
.sub-setting {
  color: #303133;
}
.title {
  font-size: 18px;
  font-weight: bolder;
}
.h3 {
  height: 40px;
  line-height: 40px;
  background: #ebeef5;
  margin-bottom: 20px;
  font-size: 18px;
  font-family: PingFangSC-Semibold, PingFang SC;
  font-weight: 600;
  color: #303133;
  padding-left: 25px;
}
.flex-2 {
  display: flex;
  align-items: center;
  .desc-2 {
    margin-left: 10px;
    font-size: 12px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #909399;
    line-height: 20px;
  }
}
.set-model {
  margin: 0 0 0 55px;
  .sub-title {
    font-size: 14px;
    margin-right: 20px;
  }
  .switch {
    margin-right: 10px;
  }
  .desc {
    color: #909399;
    font-size: 12px;
    line-height: 18px;
    margin: 0 0 0 80px;
  }
}
.withdraw-switch {
  .desc {
    margin: 0 0 0 130px;
  }
}
.commission-set .desc {
  color: #909399;
  font-size: 12px;
  margin: 0 0 10px;
}
.commission-set .choice-desc {
  margin: 10px 0 0 24px;
}
.sub-title,
/deep/ .el-radio-group {
  display: inline-block;
  vertical-align: top;
  font-size: 15px;
  color: #606266;
}

.method {
  margin-top: -11px;
}

.set-des {
  margin: 20px 0 20px 50px;
}
.search-input {
  width: 300px;
}
.base-bill {
  margin: 20px 0;
}
.set-input,
.guide-input {
  margin-left: 25px;
  width: 490px;
}
.set-input {
  margin: 15px 0 0 48px;
}
.set-input span:nth-child(1) {
  width: 100px;
  text-align: right;
}
.set-input span {
  vertical-align: middle;
}
.set-input /deep/ .cy-form-item__label {
  font-size: 14px;
  margin-right: 20px;
  color: #909399;
  padding: 0;
}
.guide-input /deep/ .cy-form-item__label {
  width: 100px;
  text-align: right;
  margin-right: 20px;
  padding: 0;
}
.guide-input /deep/ .el-input__inner {
  width: 130px;
}
.explain {
  font-size: 15px;
}
.deduct {
  margin: 20px 0 0 20px;
}
.count-deduct {
  width: 126px;
  text-align: right;
  color: #909399;
}
.set-des .count,
.count-info {
  vertical-align: middle;
  margin: 0;
  font-size: 12px;
  color: #909399;
}
.sub-notice {
  margin-bottom: 20px;
}
.search-input /deep/ .el-input-group__prepend,
.search-input /deep/ .el-input-group__append {
  background-color: #ebeef4 !important;
}
.radio {
  display: block;
}
.withdraw {
  margin: 0;
}
</style>
