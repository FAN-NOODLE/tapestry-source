<!-- 类目规则 -->
<template>
  <el-form :model="formData" :rules="rules" :inline="true" ref="ruleForm">
    <!-- 条件规则 -->
    <div>
      <div
        v-for="(item, index) in formData.conditionRule"
        :key="index"
        class="to-store-rate"
      >
        <span>当导购完成值</span>

        <!-- 高于/达到 -->
        <el-form-item
          :prop="`conditionRule.${index}.rule_if`"
          :rules="rules.rule_if"
        >
          <el-select
            v-model="item.rule_if"
            class="sm-select"
            :disabled="disabled"
          >
            <el-option
              v-for="item in condition"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 条件值 -->
        <el-form-item
          :prop="`conditionRule.${index}.rule_value`"
          :rules="rules.rule_value"
        >
          <el-input
            :disabled="disabled"
            class="bg-input"
            v-model="item.rule_value"
            oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
            placeholder="请输入大于0的正整数"
          ></el-input>
        </el-form-item>
        <span v-if="type === 16">%，</span>
        <span v-if="type === 19">元，</span>
        <span>增加</span>

        <!-- 积分 -->
        <el-form-item
          :prop="`conditionRule.${index}.integral`"
          :rules="rules.integral"
        >
          <el-input
            :disabled="disabled"
            class="bg-input"
            v-model="item.integral"
            oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
            placeholder="请输入大于0的正整数"
          ></el-input>
          <span>积分</span>
        </el-form-item>

        <!-- 按钮操作新增、删除 -->
        <template v-if="!disabled">
          <i
            class="el-icon-remove icon disabled-icon"
            v-if="formData.conditionRule.length === 1"
          />
          <i class="el-icon-remove icon" v-else @click="delTerm(index)" />

          <i
            class="el-icon-circle-plus icon disabled-icon"
            v-if="formData.conditionRule.length === 5 && !index"
          />

          <i
            class="el-icon-circle-plus icon"
            v-else-if="!index"
            @click="addTerm()"
          />
        </template>
      </div>
    </div>
  </el-form>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import { condition, ItemType } from "../../config";

@Component({})
export default class OnlineRule extends Vue {
  /** 16.奖励规则类型 线上销售业绩率（每月），17 线上销售金额（每月）*/
  @Prop() type!: number;

  /** 详情 */
  @Prop() info!: Array<ItemType>;

  /** 是否置灰 */
  @Prop({ default: false }) disabled?: boolean;

  /** todo: store_type:0默认 1奥莱 2零售 */
  @Prop({ default: 0 }) storeType?: number;

  /** 生成item */
  generateItem(): ItemType {
    return {
      rule_if: 2,
      rule_if1: 0,
      integral_type: 1,
      rule_value: "",
      integral: "",
      store_type: this.storeType
    };
  }

  /** 生成默认表单数据 */
  generateData() {
    return {
      conditionRule: [this.generateItem()]
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
  @Watch("info", { deep: true, immediate: true })
  changeInfo() {
    if (this.info && this.info.length > 0) {
      this.formData.conditionRule = this.info;
    } else {
      this.formData.conditionRule = [this.generateItem()];
    }
  }

  /** 条件枚举 */
  private condition = condition;

  /** 表单数据 */
  private formData = this.generateData();

  /** 校验规则 */
  private rules = {
    integral: [{ required: true, message: "请输入积分", trigger: "blur" }],
    rule_value: [
      { required: true, validator: this.validIntegral, trigger: "blur" }
    ],
    rule_if: [{ required: true, message: "请选择条件", trigger: "change" }],
    rule_if1: [{ required: true, message: "请选择条件", trigger: "change" }]
  };

  /** 校验规则值 */
  private validIntegral(
    rule: any,
    value: { getTime: () => number }[],
    callback: (error?: Error) => void
  ) {
    if (!value) {
      return callback(new Error("请输入大于0的正整数"));
    }

    const index = Number(rule.fullField.split(".")[1]);
    const filterItem = this.formData.conditionRule[index];

    // 判断值是否相等
    const filterRules = this.formData.conditionRule.filter(item => {
      return (
        item.rule_if === filterItem.rule_if &&
        Number(item.rule_value) === Number(filterItem.rule_value)
      );
    });
    if (filterRules.length > 1) {
      return callback(new Error("不能存在相同的规则"));
    }

    callback();
  }

  /** 供父组件获取提交数据 */
  public getData() {
    // 数据转换
    const data = this.formData.conditionRule.map(item => ({
      ...item,
      rule_value: Number(item.rule_value),
      integral: Number(item.integral)
    }));
    return data;
  }

  /** 新增项 */
  private addTerm() {
    this.formData.conditionRule.push(this.generateItem());
  }

  /** 删除项 */
  private delTerm(index: number) {
    this.formData.conditionRule.splice(index, 1);
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
