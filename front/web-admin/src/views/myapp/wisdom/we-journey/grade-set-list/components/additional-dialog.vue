<!-- 附加条件 -->
<template>
  <div class="additional-dialog">
    <el-form
      :model="formData"
      :rules="rules"
      ref="ruleForm"
      inline
      label-width="100px"
      label-position="left"
    >
      <div
        v-for="(item, index) in formData.relation_rule"
        :key="index"
        class="rule"
      >
        <span class="title">能力指标</span>
        <el-form-item
          :prop="`relation_rule.${index}.incentive_config_type`"
          :rules="rules.incentive_config_type"
        >
          <el-select
            v-model="item.incentive_config_type"
            @visible-change="categoryVisibleChange"
            @change="categoryChange(item)"
          >
            <el-option
              v-for="item in conditions.abilityOptions"
              :key="item.incentive_config_type"
              :label="item.incentive_config_name"
              :value="item.incentive_config_type"
              :disabled="item.disabled"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <!-- 高于/达到 -->
        <el-form-item
          :prop="`relation_rule.${index}.rule_condition`"
          :rules="rules.rule_condition"
        >
          <el-select v-model="item.rule_condition" class="sm-select">
            <el-option
              v-for="item in conditions.condition"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 等级 -->
        <el-form-item
          :prop="`relation_rule.${index}.rule_give`"
          :rules="rules.rule_give"
        >
          <el-select
            v-model="item.rule_give"
            class="sm-select"
            @visible-change="giveVisibleChange($event, index)"
            @change="giveChange(item)"
          >
            <el-option
              v-for="item in conditions.condition4"
              :key="item.rule_give"
              :label="item.rule_give"
              :value="item.rule_give"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 按钮操作新增、删除 -->
        <i class="el-icon-remove icon" @click="delTerm(index)" />

        <i
          class="el-icon-circle-plus icon disabled-icon"
          v-if="
            formData.relation_rule.length ===
              conditions.abilityOptions.length && !index
          "
        />
        <i
          class="el-icon-circle-plus icon"
          v-else-if="
            formData.relation_rule.length !==
              conditions.abilityOptions.length && !index
          "
          @click="addTerm()"
        />
      </div>
    </el-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import { condition, abilityOptions } from "../../config";
import { getOptionApi } from "@/api/we-journey/grade";
import { ApiResult } from "@/types";

/** 等级 */
type giveType = {
  rule_give: string;
  rule_value: number;
  id: number;
};

/** 下拉类型 */
type optionsType = {
  /** 维度type */
  incentive_config_type: number;
  /** 维度名称 */
  incentive_config_name: string;
  /** 是否已选 */
  disabled: boolean;
  /** 规则类目 */
  incentive_relation_rules: Array<giveType>;
};

/** 生成item */
function generateItem() {
  return {
    incentive_config_type: 0,
    rule_condition: 3,
    rule_give: "",
    incentive_relation_rule_id: 0
  };
}

/** 生成默认表单数据 */
function generateData() {
  return {
    /** 规则数组数据 */
    relation_rule: [generateItem()]
  };
}

const newCondition = condition.splice(1, 1);

@Component({})
export default class AdditionalDialog extends Vue {
  /** list */
  @Prop() list!: Array<ReturnType<typeof generateItem>>;

  /** 表单数据 */
  private formData = generateData();

  /** private  */
  private conditions = {
    abilityOptions: [] as Array<optionsType>,
    condition: newCondition,
    condition4: [] as Array<giveType>,
    ruleTitle: ""
  };

  /** 校验规则 */
  private rules = {
    incentive_config_type: [
      { required: true, message: "请选择能力指标", trigger: "change" }
    ],
    rule_condition: [{ required: true, message: "请选择", trigger: "change" }],
    rule_give: [
      { required: true, message: "请选择能力指标等级", trigger: "change" }
    ]
  };

  /** 类目下拉框出现时触发 */
  private categoryVisibleChange(e: boolean) {
    // 如果是隐藏, 直接返回
    if (!e) return;

    // 已选中的类目ids
    const checkedIds = this.formData.relation_rule.map(item =>
      Number(item.incentive_config_type)
    );

    // 过滤出能选择的类目
    this.conditions.abilityOptions = this.conditions.abilityOptions.map(
      item => {
        return {
          ...item,
          disabled: checkedIds.includes(item.incentive_config_type)
        };
      }
    );
  }

  /** 等级下拉框出现时触发 */
  private giveVisibleChange(isShow: boolean, index: number) {
    // 如果是隐藏, 直接返回
    if (!isShow) return;

    // 已选中的类目id
    const ruleId = this.formData.relation_rule[index].incentive_config_type;

    // 从规则类目中过滤出等级
    const checkedItem = this.conditions.abilityOptions.filter(
      item => item.incentive_config_type === ruleId
    );
    this.conditions.condition4 = checkedItem[0].incentive_relation_rules;
  }

  /** 类目变化更换默认等级 */
  private categoryChange(ruleItem: ReturnType<typeof generateItem>) {
    ruleItem.rule_give = "";
  }

  /** 给等级赋值 */
  private giveChange(ruleItem: ReturnType<typeof generateItem>) {
    const checkedItem = this.conditions.condition4.filter(
      item => item.rule_give === ruleItem.rule_give
    );
    ruleItem.incentive_relation_rule_id = checkedItem[0].id;
  }

  /** 组合下拉框筛选 */
  private async getOptions() {
    try {
      const res: ApiResult = await getOptionApi();

      if (res.code !== 200) return this.$message.error(res.message as string);

      // 过滤掉等级为空的数组
      const newDate = res.data.filter(
        (item: optionsType) =>
          item.incentive_relation_rules && item.incentive_relation_rules.length
      );

      // 为元素添加名称、disabled
      this.conditions.abilityOptions = newDate.map((item: optionsType) => ({
        ...item,
        incentive_config_name:
          abilityOptions[item.incentive_config_type - 1].label,
        disabled: false
      }));

      // 新增需要设置默认选中
      if (!(this.list && this.list.length)) {
        // 初始化选中
        const checkedOptions = this.conditions.abilityOptions[0];

        this.defaultSetValue(this.formData.relation_rule[0], checkedOptions);
      }
    } catch (e) {
      console.error(e, "数据解析失败");
    }
  }

  /** 附加条件默认赋值 */
  private defaultSetValue(
    relation_rule: ReturnType<typeof generateItem>,
    selectedList: optionsType
  ) {
    try {
      const { incentive_config_type, incentive_relation_rules } = selectedList;
      const { rule_give, id } = incentive_relation_rules[0] || {};

      relation_rule.incentive_config_type = incentive_config_type;
      relation_rule.rule_give = rule_give;
      relation_rule.incentive_relation_rule_id = id;
    } catch (e) {
      console.error(e, "附加条件赋值失败");
    }
  }

  /** 新增项 */
  private addTerm() {
    this.formData.relation_rule.push(generateItem());

    // 更新列表
    this.categoryVisibleChange(true);

    // 删选可选的
    const selectedList = this.conditions.abilityOptions.filter(
      item => !item.disabled
    );

    // 需要赋值的元素
    const relation_rule = this.formData.relation_rule[
      this.formData.relation_rule.length - 1
    ];

    this.defaultSetValue(relation_rule, selectedList[0]);
  }

  /** 删除项 */
  private delTerm(index: number) {
    this.formData.relation_rule.splice(index, 1);
    // 通知父组件关闭附件条件
    if (!this.formData.relation_rule.length) {
      this.$emit("close");
    }
  }

  /** 供父组件获取提交数据 */
  public getData() {
    return this.formData.relation_rule;
  }

  /** 生命周期 */
  private created() {
    // 查询下拉列表
    this.getOptions();

    if (this.list && this.list.length) {
      return (this.formData.relation_rule = this.list);
    }

    // 清空原有数据
    this.formData = generateData();
  }
}
</script>
<style lang="scss" scoped>
.additional-dialog {
  margin-left: 95px;
  .title {
    margin-right: 10px;
  }
}
.icon {
  padding-top: 10px;
  cursor: pointer;
  font-size: 20px;
  vertical-align: middle;
  margin-left: 10px;
}

.disabled-icon {
  color: #909399;
  cursor: not-allowed;
}

.rule {
  display: flex;
  margin-bottom: 10px;
  /deep/.el-form-item {
    display: inline-block;
  }
}
.sm-select {
  margin: 0 5px;
  /deep/.el-input {
    width: 100px !important;
  }
}

.bg-input {
  width: 170px !important;
  margin: 0 10px;
}
</style>
