<!-- 类目规则 -->
<template>
  <el-form :model="formData" :rules="rules" :inline="true" ref="ruleForm">
    <!-- 次数规则 -->
    <div>
      <span>培训平台积分与导购积分兑换比为：</span>
      <el-form-item prop="rule_value" :rules="rules.rule_value">
        <el-input
          :disabled="disabled"
          class="bg-input"
          v-model="formData.rule_value"
          oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
          placeholder="请输入大于0的正整数"
        ></el-input>
      </el-form-item>

      <span>：</span>
      <el-form-item prop="integral" :rules="rules.integral">
        <el-input
          :disabled="disabled"
          class="bg-input"
          v-model="formData.integral"
          oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
          placeholder="请输入大于0的正整数"
        ></el-input>
      </el-form-item>
    </div>
    <div class="tips">
      当前转换比为：{{ formData.rule_value }} 培训平台积分 兑换
      {{ formData.integral }}
      导购积分（兑换后积分会做四舍五入保留一位小数展示在积分明细）
    </div>
  </el-form>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import { ItemType } from "../../config";

@Component({})
export default class TrainRule extends Vue {
  /** 奖励规则类型 15.积分等比兑换 */
  @Prop() type!: number;

  /** 详情 */
  @Prop() info!: Array<ItemType>;

  /** 是否置灰 */
  @Prop({ default: false }) disabled?: boolean;
  /** todo: store_type:0默认 1奥莱 2零售 */
  @Prop({ default: 0 }) storeType?: number;
  /** 生成默认表单数据 */
  generateData(): ItemType {
    return {
      rule_if: 1,
      rule_if1: 0,
      integral_type: 1,
      rule_value: "",
      integral: "",
      store_type: this.storeType
    };
  }

  /** 监听奖励规则类型 变化 */
  @Watch("type", { immediate: true })
  changeType() {
    // 重置表单数据，防止数据缓存
    this.formData = this.generateData();
    // 先清空校验
    this.$nextTick(() => {
      (this.$refs["ruleForm"] as HTMLFormElement).clearValidate();
    });
  }

  /** 监听奖励规则类型 变化 */
  @Watch("info", { immediate: true, deep: true })
  changeInfo() {
    if (this.info && this.info.length > 0) {
      this.formData = this.info[0];
    } else {
      this.formData = this.generateData();
    }
  }

  /** 校验规则 */
  private rules = {
    integral: [
      { required: true, message: "请输入大于0的正整数", trigger: "blur" }
    ],
    rule_value: [
      { required: true, message: "请输入大于0的正整数", trigger: "blur" }
    ],
    rule_if: [{ required: true, message: "请选择条件", trigger: "change" }],
    rule_if1: [{ required: true, message: "请选择条件", trigger: "change" }]
  };

  /** 表单数据 */
  private formData = this.generateData();

  /** 供父组件获取提交数据 */
  public getData() {
    // 格式化数组
    return [this.formData].map(item => ({
      ...item,
      rule_value: Number(item.rule_value),
      integral: Number(item.integral)
    }));
  }
}
</script>
<style lang="scss" scoped>
.icon {
  cursor: pointer;
  font-size: 20px;
  vertical-align: middle;
  margin-left: 10px;
}

.tips {
  color: #909399;
  font-size: 12px;
}

.disabled-icon {
  color: #909399;
  cursor: not-allowed;
}

/deep/.el-form-item {
  margin-right: 0px !important;
}

.bg-input {
  width: 170px !important;
  margin: 0 10px;
}
.sm-select {
  margin: 0 5px;
  /deep/.el-input {
    width: 100px !important;
  }
}

.to-store-rate {
  margin-bottom: 10px;
}
</style>
