<!-- 新增积分规则弹窗 -->
<template>
  <el-dialog :visible.sync="syncVisible" width="1000px" :title="title" center>
    <div class="atlas-body">
      <el-form
        :model="formData"
        :rules="rules"
        ref="ruleForm"
        label-width="100px"
        label-position="left"
      >
        <el-form-item label="能力图谱" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :label="2">关闭</el-radio>
            <el-radio :label="1">开启</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 开启时候才展示 -->
        <div v-if="formData.status === 1">
          <el-form-item label="能力指标" prop="ability_index">
            <el-select v-model="formData.ability_index">
              <el-option
                v-for="item in conditions.condition5"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
            <div class="tips">{{ tip }}</div>
          </el-form-item>
          <!-- todo: 指标渠道：通用，区分零售和奥莱 -->
          <el-form-item label="指标渠道" prop="is_diff_store">
            <el-radio-group
              v-model="formData.is_diff_store"
              @change="changeDiffStore"
            >
              <el-radio :label="0">通用规则</el-radio>
              <el-radio :label="1">区分零售和奥莱</el-radio>
            </el-radio-group>
          </el-form-item>
          <!-- todo:通用规则 -->
          <el-form-item
            label="指标规则"
            required
            v-if="formData.is_diff_store === 0"
          >
            <RelationRules
              ref="relationRules"
              :type="type"
              :rulesArray="formData.relation_rule"
              :storeType="0"
              @changeGive="changeGive"
              @addTerm="addTerm"
              @delTerm="delTerm"
            ></RelationRules>
          </el-form-item>
          <!-- todo:奥莱和零售 -->
          <el-form-item
            label="指标规则"
            required
            v-if="formData.is_diff_store === 1"
          >
            <div class="rules-content">
              <div class="rules-title">零售渠道</div>
              <div class="rules-item">
                <RelationRules
                  ref="lsRelationRules"
                  :type="type"
                  :rulesArray="formData.ls_relation_rule"
                  :storeType="2"
                  @changeGive="changeGive"
                  @addTerm="addTerm"
                  @delTerm="delTerm"
                ></RelationRules>
              </div>
              <div class="rules-title">奥莱渠道</div>
              <div class="rules-item">
                <RelationRules
                  ref="alRelationRules"
                  :type="type"
                  :rulesArray="formData.al_relation_rule"
                  :storeType="1"
                  @changeGive="changeGive"
                  @addTerm="addTerm"
                  @delTerm="delTerm"
                ></RelationRules>
              </div>
            </div>
          </el-form-item>
        </div>
      </el-form>
    </div>
    <span slot="footer" class="dialog-footer">
      <el-button @click="syncVisible = false">取 消</el-button>
      <el-button type="primary" @click="save">确 定</el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component, PropSync, Prop, Watch } from "vue-property-decorator";
import { condition3, condition4, condition5, atlasTips } from "../../config";
import { getRelationSaveApi } from "@/api/we-journey/Integral";
import { ApiResult } from "@/types";
import RelationRules from "./relation-rules.vue";
/** form 类型 */
type fromType = {
  /** 能力指标条件 */
  ability_index: number;
  /** 开启状态 */
  status: number;
  /** 规则数组数据 */
  relation_rule: Array<{
    rule_condition: string;
    rule_value: string;
    rule_give: string | number;
  }>;
  // todo:奥莱渠道
  al_relation_rule: Array<{
    rule_condition: string;
    rule_value: string;
    rule_give: string | number;
  }>;
  // todo:零售渠道
  ls_relation_rule: Array<{
    rule_condition: string;
    rule_value: string;
    rule_give: string | number;
  }>;
  /** 指标渠道 */
  is_diff_store: number;
};

/** 生成item */
function generateItem(storeType = 0) {
  return {
    rule_condition: "",
    rule_value: "",
    rule_give: "",
    store_type: storeType
  };
}

/** 生成默认表单数据 */
function generateData(): fromType {
  return {
    /** 能力指标条件 */
    ability_index: 0,
    /** 开启状态 */
    status: 2,
    /** 规则数组数据 */
    relation_rule: [generateItem(0)],
    // todo:奥莱渠道
    al_relation_rule: [generateItem(1)],
    // todo:零售渠道
    ls_relation_rule: [generateItem(2)],
    /** 指标渠道 */
    is_diff_store: 1
  };
}

/** 默认标题 */
const defaultText = [
  { title: "", ruleTitle: "" },
  { title: "邀约能力指标配置", ruleTitle: "当导购邀约目标率" },
  { title: "任务能力指标配置", ruleTitle: "当导购任务完成率" },
  { title: "客户互动能力指标配置", ruleTitle: "当导购总互动积分排行" },
  { title: "培训能力指标配置", ruleTitle: "当导购总培训积分排行" },
  { title: "线上能力指标配置", ruleTitle: "当导购线上目标率" }
];

@Component({
  components: {
    RelationRules
  }
})
export default class AtlasDialog extends Vue {
  /** 判断父属性显示影藏，双向绑定 */
  @PropSync("visible", { type: Boolean }) syncVisible!: boolean;

  /** 维度id */
  @Prop() type!: number;

  /** 详情 */
  @Prop() info?: any;

  @Watch("syncVisible")
  changeVisible(show: boolean) {
    if (!show) return;
    // 清空原有数据
    this.formData = generateData();
    // 初始化下拉框数据
    this.initOptions();

    // 编辑情况下赋值
    if (JSON.stringify(this.info) !== "{}") {
      //todo:需要将relation_rule中store_type为1的数据放到al_relation_rule中，store_type为2的数据放到ls_relation_rule中
      const relation_rule = this.info.relation_rule.filter(
        (item: { store_type: number }) => item.store_type === 0
      );
      const al_relation_rule = this.info.relation_rule.filter(
        (item: { store_type: number }) => item.store_type === 1
      );
      const ls_relation_rule = this.info.relation_rule.filter(
        (item: { store_type: number }) => item.store_type === 2
      );
      this.formData = {
        ...this.info,
        relation_rule:
          relation_rule.length === 0 ? [generateItem(0)] : relation_rule,
        al_relation_rule:
          al_relation_rule.length === 0 ? [generateItem(1)] : al_relation_rule,
        ls_relation_rule:
          ls_relation_rule.length === 0 ? [generateItem(2)] : ls_relation_rule
      } as fromType;
    }
  }

  /** 弹窗title */
  private title = "";

  /** 表单数据 */
  private formData = generateData();

  /** private  */
  private conditions = {
    condition3,
    condition4,
    condition5,
    ruleTitle: ""
  };

  /** 校验规则 */
  private rules = {
    status: [{ required: true, message: "请选择开启状态", trigger: "change" }],
    ability_index: [
      { required: true, message: "请选择能力指标", trigger: "change" }
    ],
    rule_condition: [{ required: true, message: "请选择", trigger: "change" }],
    rule_value: [
      { required: true, message: "请输入大于0的正整数", trigger: "blur" }
    ],
    rule_give: [
      { required: true, message: "请选择能力指标等级", trigger: "change" }
    ],
    is_diff_store: [
      { required: true, message: "请选择指标渠道", trigger: "change" }
    ]
  };

  get tip() {
    return atlasTips[this.type - 1];
  }
  /** 能力指标变化, 指定置灰项 */
  private changeGive(receiveObj: any) {
    const { open, storeType } = receiveObj;
    if (!open) return;
    if (storeType === 0) {
      const ids = this.formData.relation_rule.map(item => item.rule_give);
      this.conditions.condition4.forEach(item => {
        item.disabled = ids.includes(item.label);
      });
    } else if (storeType === 1) {
      const al_ids = this.formData.al_relation_rule.map(item => item.rule_give);
      this.conditions.condition4.forEach(item => {
        item.disabled = al_ids.includes(item.label);
      });
    } else {
      const ls_ids = this.formData.ls_relation_rule.map(item => item.rule_give);
      this.conditions.condition4.forEach(item => {
        item.disabled = ls_ids.includes(item.label);
      });
    }
  }

  /** 初始化配置 */
  private initOptions() {
    // 根据状态对能力指标单独处理 - 下拉框option
    this.conditions.condition5 = condition5.filter(
      item => item.value === this.type
    );
    // 设置默认选中
    this.formData.ability_index = this.conditions.condition5[0].value;

    // 根据状态对指标规则做处理 - 下拉框option
    switch (this.type) {
      // 任务维度
      case 2:
      case 3:
      case 4:
        this.conditions.condition3 = condition3.filter(
          item => item.value !== 1
        );
        break;
      // 客户互动维度、 培训维度

      // this.conditions.condition3 = condition3.filter(
      //   item => item.value === 3 || item.value === 4
      // );
      // break;
      default:
        this.conditions.condition3 = condition3;
    }

    // 根据状态对指标规则文案做处理 - 文案
    this.conditions.ruleTitle = defaultText[this.type].ruleTitle;

    // 设置dialog标题
    this.title = defaultText[this.type].title;
  }

  /** 新增项 */
  private addTerm(storeType: number) {
    if (storeType === 0) {
      this.formData.relation_rule.push(generateItem(0));
    } else if (storeType === 1) {
      this.formData.al_relation_rule.push(generateItem(1));
    } else {
      this.formData.ls_relation_rule.push(generateItem(2));
    }
  }

  /** 删除项 */
  private delTerm(receiveObj: any) {
    const { index, storeType } = receiveObj;
    if (storeType === 0) {
      this.formData.relation_rule.splice(index, 1);
    } else if (storeType === 1) {
      this.formData.al_relation_rule.splice(index, 1);
    } else {
      this.formData.ls_relation_rule.splice(index, 1);
    }
  }
  /**
   * @description: 子组件表单验证方法
   */
  private allValidate() {
    //  子组件表单验证集合
    const childValidate = [];
    // 如果通用规则存在
    if (this.formData.is_diff_store === 0) {
      // 通用规则验证
      const childDom = this.$refs["relationRules"] as HTMLFormElement;
      const childForm = childDom.$refs["ruleForm"] as HTMLFormElement;
      childValidate.push(
        new Promise<void>((resolve, reject) => {
          childForm.validate((valid: boolean) => {
            if (!valid) {
              reject(valid);
            } else {
              resolve();
            }
          });
        })
      );
    } else {
      // 奥莱规则验证
      const alChildDom = this.$refs["alRelationRules"] as HTMLFormElement;
      const alChildForm = alChildDom.$refs["ruleForm"] as HTMLFormElement;
      childValidate.push(
        new Promise<void>((resolve, reject) => {
          alChildForm.validate((valid: boolean) => {
            if (!valid) {
              reject();
            } else {
              resolve();
            }
          });
        })
      );
      // 零售规则验证
      const lsChildDom = this.$refs["lsRelationRules"] as HTMLFormElement;
      const lsChildForm = lsChildDom.$refs["ruleForm"] as HTMLFormElement;
      childValidate.push(
        new Promise<void>((resolve, reject) => {
          lsChildForm.validate((valid: boolean) => {
            if (!valid) {
              reject();
            } else {
              resolve();
            }
          });
        })
      );
    }
    // 父元素验证
    const parentDom = this.$refs["ruleForm"] as HTMLFormElement;
    childValidate.push(
      new Promise<void>((resolve, reject) => {
        parentDom.validate((valid: boolean) => {
          if (!valid) {
            reject();
          } else {
            resolve();
          }
        });
      })
    );
    return Promise.all(childValidate);
  }
  /** 保存 */
  private async save() {
    try {
      // 验证父子组件所有表单
      await this.allValidate();
      // 整合提交数据,将al_relation_rule和ls_relation_rule中的数据放到relation_rule中
      const data: any = {
        ...this.formData,
        incentive_config_type: this.type,
        time_rule: this.conditions.condition5[0].time_rule
      };
      if (data.is_diff_store === 1) {
        data.relation_rule = [
          ...data.al_relation_rule,
          ...data.ls_relation_rule
        ];
      }
      // 删除多余字段
      delete data.al_relation_rule;
      delete data.ls_relation_rule;
      delete data.incentive_relation_rules;
      console.log(data);

      const res: ApiResult = await getRelationSaveApi(data);

      if (res.code !== 200) return this.$message.error(res.message as string);

      this.syncVisible = false;
      this.$emit("success");
    } catch (error) {
      console.log("验证失败");
    }
  }
  /**
   * @description: 渠道变化
   */
  changeDiffStore() {
    this.formData.relation_rule = [generateItem(0)];
    this.formData.al_relation_rule = [generateItem(1)];
    this.formData.ls_relation_rule = [generateItem(2)];
  }
}
</script>
<style lang="scss" scoped>
.tips {
  color: #909399;
  font-size: 12px;
}

.atlas-body {
  max-height: 400px;
  overflow-y: scroll;
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
.rules-content {
  .rules-title {
    color: #708df5;
  }
}
</style>
