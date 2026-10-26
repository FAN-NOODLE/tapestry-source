<template>
  <layout-form>
    <template slot="content">
      <el-form
        :model="formData"
        :rules="rules"
        ref="ruleForm"
        label-width="100px"
        label-position="left"
      >
        <el-form-item label="规则名称" prop="rule_name">
          <el-input
            v-model="formData.rule_name"
            maxlength="15"
            :disabled="pageType === 3"
            show-word-limit
          ></el-input>
        </el-form-item>

        <div class="model-title" v-if="typeCategory.length !== 1">规则类目</div>

        <div v-for="(item, index) in formData.rule_detail" :key="index">
          <div class="category" v-if="typeCategory.length !== 1">
            <div>
              类目{{
                ["一", "二", "三", "四", "五", "六", "七", "八", "九"][index]
              }}：
            </div>

            <i
              class="el-icon-delete del-icon"
              @click="delCategory(index)"
              v-if="pageType !== 3 && index"
            />
          </div>
          <!-- 规则类目 -->
          <el-form-item
            label="规则类目"
            :prop="`rule_detail.${index}.rule_category_id`"
            :rules="rules.rule_category_id"
          >
            <el-select
              :disabled="pageType === 3"
              v-model="item.rule_category_id"
              placeholder="请选择规则类目"
              @visible-change="categoryChange"
            >
              <el-option
                v-for="item in typeCategory"
                :key="item.value"
                :label="item.label"
                :value="item.value"
                :disabled="item.disabled"
              >
              </el-option>
            </el-select>
          </el-form-item>
          <!-- todo:规则渠道 -->
          <el-form-item
            label="规则渠道"
            :prop="`rule_detail.${index}.is_diff_store`"
            :rules="rules.is_diff_store"
          >
            <el-radio-group
              v-model="item.is_diff_store"
              :disabled="pageType === 3"
              @change="radioChange(index)"
            >
              <el-radio :label="0">通用规则</el-radio>
              <el-radio :label="1">区分零售和奥莱</el-radio>
            </el-radio-group>
          </el-form-item>
          <!-- todo:通用奖励规则 -->
          <el-form-item
            class="is-required"
            v-if="item.is_diff_store === 0"
            label="奖励规则"
          >
            <!-- todo: store_type:0默认 1奥莱 2零售 -->
            <component
              :is="currentCom"
              :storeType="0"
              :info="item.detail_value"
              :type="item.rule_category_id"
              :disabled="pageType === 3"
              :ref="`integral${index}`"
            ></component>
          </el-form-item>
          <!-- todo：区分零售和奥莱奖励规则 -->
          <el-form-item
            class="is-required"
            v-if="item.is_diff_store === 1"
            label="奖励规则"
          >
            <div class="rules-content">
              <div class="rules-title">零售渠道</div>
              <div class="rules-item">
                <!-- todo: store_type:0默认 1奥莱 2零售 -->
                <component
                  :is="currentCom"
                  :storeType="2"
                  :info="item.detail_value_ls"
                  :type="item.rule_category_id"
                  :disabled="pageType === 3"
                  :ref="`ls_integral${index}`"
                ></component>
              </div>
              <div class="rules-title">奥莱渠道</div>
              <div class="rules-item">
                <!-- todo: store_type:0默认 1奥莱 2零售 -->
                <component
                  :is="currentCom"
                  :storeType="1"
                  :info="item.detail_value_al"
                  :type="item.rule_category_id"
                  :disabled="pageType === 3"
                  :ref="`al_integral${index}`"
                ></component>
              </div>
            </div>
          </el-form-item>
          <!-- 规则有效期 -->
          <el-form-item
            label="规则有效期"
            :prop="`rule_detail.${index}.validity_type`"
            :rules="rules.validity_type"
          >
            <el-radio-group
              v-model="item.validity_type"
              :disabled="pageType === 3"
            >
              <el-radio :label="1">长期有效</el-radio>
              <el-radio :label="2">指定时间</el-radio>
            </el-radio-group>
          </el-form-item>

          <!-- 指定时间 -->
          <el-form-item
            v-if="item.validity_type === 2"
            label=""
            :prop="`rule_detail.${index}.time`"
            :rules="rules.time"
          >
            <el-date-picker
              :disabled="pageType === 3"
              v-model="item.time"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              :default-time="['00:00:00', '23:59:59']"
              :picker-options="{ disabledDate }"
            >
            </el-date-picker>
          </el-form-item>
          <!-- 是否有上限 -->
          <el-form-item
            v-if="showAstrict(item.rule_category_id)"
            label="是否有上限"
            :prop="`rule_detail.${index}.is_astrict`"
            :rules="rules.is_astrict"
          >
            <el-radio-group
              v-model="item.is_astrict"
              :disabled="pageType === 3"
            >
              <el-radio :label="1">无</el-radio>
              <el-radio :label="2">有</el-radio>
            </el-radio-group>
          </el-form-item>

          <!-- 积分 -->
          <div class="inline-form" v-if="item.is_astrict === 2">
            <!-- todo：每天，每周，每月能，每年下拉 -->
            <span>导购</span>
            <el-form-item
              label=""
              :prop="`rule_detail.${index}.astrict_type`"
              :rules="rules.astrict_type"
            >
              <el-select
                class="small-select"
                :disabled="pageType === 3"
                v-model="item.astrict_type"
                placeholder="请选择"
              >
                <el-option
                  v-for="item in integralType"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
            <span>该类型积分最多可增加</span>
            <el-form-item
              label=""
              :prop="`rule_detail.${index}.max_integral`"
              :rules="rules.max_integral"
            >
              <el-input
                :disabled="pageType === 3"
                class="bg-input"
                v-model="item.max_integral"
                oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
                placeholder="请输入大于0的正整数"
              ></el-input>
            </el-form-item>
            <span>积分</span>
          </div>
        </div>
      </el-form>

      <!-- 添加类目 -->
      <el-button
        v-if="
          pageType !== 3 && formData.rule_detail.length !== typeCategory.length
        "
        type="primary"
        size="mini"
        @click="addCategory"
        >添加类目</el-button
      >
    </template>
    <template slot="footer">
      <el-button type="primary" @click="save" v-if="pageType !== 3"
        >保存</el-button
      >
      <el-button @click="goBack">{{
        pageType === 3 ? "返回" : "取消"
      }}</el-button>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import IntegralRule from "./integral-rule.vue";
import TaskRule from "./task-rule.vue";
import CustomerRule from "./customer-rule.vue";
import TrainRule from "./train-rule.vue";
import OnlineRule from "./onLine-rule.vue";
import { TYPE_CATEGORY, categoryType, integralType } from "../../config";
import {
  getIntegralAddApi,
  getIntegralEditApi
} from "@/api/we-journey/Integral";
import { ApiResult } from "@/types";

type DetailType = {
  /** 类目规则id */
  rule_category_id: number;
  /** '1.长期 2.指定时间' */
  validity_type: number;
  /** 开始时间 */
  start_time: number;
  /** 结束时间 */
  end_time: number;
  /** 指定时间 可选 */
  time?: Array<Date>;
  /** 奖励规则 */
  rule_value: Array<any>;
  /** todo:零售规则数组 */
  rule_value_ls?: any;
  detail_value_ls?: Array<any>;
  /** todo:奥莱规则数组 */
  rule_value_al?: any;
  detail_value_al?: Array<any>;
  /** 通用规则 */
  detail_value?: Array<any>;
  /** 是否有上限 */
  is_astrict?: number;
  /** 最多积分 */
  max_integral?: number | string;
  /** 规则渠道 */
  is_diff_store?: number;
  /** 限制类型 */
  astrict_type?: number;
};

/** 生成默认表单数据 */
function generateData(): DetailType {
  return {
    rule_category_id: 0,
    validity_type: 1,
    rule_value: [],
    rule_value_ls: [],
    rule_value_al: [],
    time: [],
    start_time: 0,
    end_time: 0,
    is_astrict: 1,
    astrict_type: 1,
    is_diff_store: 0
  };
}

/** 校验时间 */
function validateTime(
  rule: unknown,
  value: { getTime: () => number }[],
  callback: (error?: Error) => void
) {
  if (!value || !value.length) {
    return callback(new Error("请选择时间"));
  }
  callback();
}

/** 需要重新承接，防止内部改动 */
const typeCateGory = TYPE_CATEGORY;

@Component({
  components: { IntegralRule, TaskRule, CustomerRule, TrainRule, OnlineRule }
})
export default class Integral extends Vue {
  /** 类型 1.邀约维度 2.任务维度 3.客户互动维度 4.培训维度 5.线上能力 */
  @Prop() type!: number;

  /** 详情 */
  @Prop() info?: { rule_name: string; rule_detail: Array<DetailType> };

  /** 页面状态 1. 新增 2.编辑 3.查看 */
  @Prop() pageType?: number;

  @Watch("info")
  changeInfo() {
    if (this.info) {
      this.generateInitOptions();
      // 处理时间
      this.formData = {
        ...this.info,
        rule_detail: this.info.rule_detail.map(item => ({
          ...item,
          time:
            item.validity_type === 2
              ? [
                  new Date(item.start_time * 1000),
                  new Date(item.end_time * 1000)
                ]
              : [],
          // todo:这里需要将store_type为2的数据放到detail_value_ls中，store_type为1的数据放到detail_value_al中
          detail_value_ls: item.detail_value?.filter(
            (item: any) => item.store_type === 2
          ),
          detail_value_al: item.detail_value?.filter(
            (item: any) => item.store_type === 1
          )
        }))
      };
      // 先清空校验
      this.$nextTick(() => {
        (this.$refs["ruleForm"] as HTMLFormElement).clearValidate();
      });
    }
  }

  /** 监听企业信息变化 */
  @Watch("$store.state.permission.EQ_NAME")
  change(value: string) {
    if (!value) return;

    // 防止还未初始化配置
    this.generateInitOptions();
    // 【KS】&【SW】本期下拉前端先隐藏 邀约成功率（每周）（邀约成功数/总邀约数）、邀约成功（每次）、邀约占比（每周）（邀约金额/个人销售金额）
    if (value !== "coach" && this.type === 1) {
      // 由于编辑和新建调用的时机不同，导致重新赋值的原有函数不同
      typeCateGory[0] = this.typeCategory = this.typeCategory.filter(item =>
        [1, 2, 6].includes(item.value)
      );
    }
  }

  /** 表单数据 */
  private formData = {
    rule_name: "",
    rule_detail: [generateData()]
  };

  /** 当前规则组件 */
  private currentCom = "";

  /** 规则类目 options - 原始值不变 */
  private typeCategory = [] as Array<categoryType>;

  /** todo: 每天，每周，每月，每年 */
  private integralType = integralType as Array<any>;
  /** 校验规则 */
  private rules = {
    rule_name: [{ required: true, message: "请输入规则名称", trigger: "blur" }],
    rule_category_id: [
      { required: true, message: "请选择规则类目", trigger: "change" }
    ],
    validity_type: [
      { required: true, message: "请选择规则有效期", trigger: "change" }
    ],
    time: [{ required: true, validator: validateTime, trigger: "change" }],
    is_astrict: [
      { required: true, message: "请选择是否有上限", trigger: "change" }
    ],
    is_diff_store: [
      { required: true, message: "请选择规则渠道", trigger: "change" }
    ],
    max_integral: [{ required: true, message: "请输入积分", trigger: "blur" }],
    astrict_type: [
      { required: true, message: "请选择上限类型", trigger: "change" }
    ]
  };

  /** 是否显示 是否有上限 */
  private get showAstrict() {
    /** 缓存类目ids */
    const ids = typeCateGory
      .flat()
      .filter(item => item.isAstrict)
      .map(item => item.value);

    return (id: number) => {
      return ids.includes(id);
    };
  }

  /** 处理时间 - 之前的时间不可选 */
  private disabledDate(time: Date) {
    return time.getTime() < Date.now() - 8.64e7;
  }

  /** 类目下拉框出现时触发 */
  private categoryChange(e: boolean) {
    // 如果是隐藏, 直接返回
    if (!e) return;

    // 已选中的类目ids
    const checkedIds = this.formData.rule_detail.map(item =>
      Number(item.rule_category_id)
    );

    // 过滤出能选择的类目, 由于培训导致类目找不到，暂时隐藏
    // const leftArr: Array<categoryType> = [];
    // const rightArr: Array<categoryType> = [];
    this.typeCategory.forEach(item => {
      const disabled = checkedIds.includes(item.value);

      item.disabled = disabled;
      // disabled ? rightArr.push(item) : leftArr.push(item);
    });

    // this.typeCategory = [...leftArr, ...rightArr];
  }

  /** 获取初始化配置 */
  private generateInitOptions() {
    this.typeCategory = typeCateGory[this.type - 1];
    // 规则类目默认选中第一个
    this.formData.rule_detail[0].rule_category_id = this.typeCategory[0].value;

    // 维度类型
    const typeMap = [
      "IntegralRule",
      "TaskRule",
      "CustomerRule",
      "TrainRule",
      "OnlineRule"
    ];
    this.currentCom = typeMap[Number(this.type) - 1];
  }

  /** 添加类目 - button */
  private addCategory() {
    // 更新最新可选类目
    this.categoryChange(true);
    this.formData.rule_detail.push(generateData());

    // 过滤出未选中的元素
    const notSelectedArr = this.typeCategory.filter(item => !item.disabled);

    // 默认选中第一项
    this.formData.rule_detail[
      this.formData.rule_detail.length - 1
    ].rule_category_id = notSelectedArr[0].value;
  }

  /** 删除类目 - button */
  private delCategory(index: number) {
    this.formData.rule_detail.splice(index, 1);
  }

  /** 组合校验规则 */
  private composeValidate() {
    // 子组件校验整合
    const childPro = [];
    // 从列表中循环获取子组件实例
    this.formData.rule_detail.forEach((item, index) => {
      // todo:通用规则表单校验
      if (item.is_diff_store === 0) {
        const promise = new Promise(resolve => {
          const childDom = this.$refs[`integral${index}`][0] as HTMLFormElement;
          // 获取子组件参数
          const childData = childDom.getData();
          // 子组件校验
          (childDom.$refs["ruleForm"] as HTMLFormElement).validate(
            (valid: boolean) => {
              if (valid) {
                // 处理指定时间
                if (item.validity_type === 2 && item.time && item.time.length) {
                  item.start_time = ~~((item.time[0] as Date).getTime() / 1000);
                  item.end_time = ~~((item.time[1] as Date).getTime() / 1000);
                  // delete item.time;
                }
                // 添加规则
                item.rule_value = childData;
                // 处理最大积分
                item.max_integral = Number(item.max_integral) || 0;
                resolve(childData);
              }
            }
          );
        });
        childPro.push(promise);
      } else {
        // todo:零售规则表单校验
        const ls_promise = new Promise(resolve => {
          const ls_childDom = this.$refs[
            `ls_integral${index}`
          ][0] as HTMLFormElement;
          // 获取子组件参数
          const ls_childData = ls_childDom.getData();
          // 子组件校验
          (ls_childDom.$refs["ruleForm"] as HTMLFormElement).validate(
            (valid: boolean) => {
              if (valid) {
                // 处理指定时间
                if (item.validity_type === 2 && item.time && item.time.length) {
                  item.start_time = ~~((item.time[0] as Date).getTime() / 1000);
                  item.end_time = ~~((item.time[1] as Date).getTime() / 1000);
                  // delete item.time;
                }
                // 添加规则
                item.rule_value_ls = ls_childData;
                // 处理最大积分
                item.max_integral = Number(item.max_integral) || 0;
                resolve(ls_childData);
              }
            }
          );
        });
        childPro.push(ls_promise);
        // todo:奥莱规则表单校验
        const al_promise = new Promise(resolve => {
          const al_childDom = this.$refs[
            `al_integral${index}`
          ][0] as HTMLFormElement;
          // 获取子组件参数
          const al_childData = al_childDom.getData();
          // 子组件校验
          (al_childDom.$refs["ruleForm"] as HTMLFormElement).validate(
            (valid: boolean) => {
              if (valid) {
                // 处理指定时间
                if (item.validity_type === 2 && item.time && item.time.length) {
                  item.start_time = ~~((item.time[0] as Date).getTime() / 1000);
                  item.end_time = ~~((item.time[1] as Date).getTime() / 1000);
                  // delete item.time;
                }
                // 添加规则
                item.rule_value_al = al_childData;
                // 处理最大积分
                item.max_integral = Number(item.max_integral) || 0;
                resolve(al_childData);
              }
            }
          );
        });
        childPro.push(al_promise);
      }
    });

    // 父组件校验整合
    const parentPro = new Promise(resolve => {
      (this.$refs["ruleForm"] as HTMLFormElement).validate((valid: boolean) => {
        if (valid) resolve(valid);
      });
    });
    childPro.push(parentPro);

    return Promise.all(childPro);
  }

  /** 保存 */
  private async save() {
    try {
      await this.composeValidate();
      const params = {
        type: this.type,
        ...this.formData
      };
      // todo:删除parmas rule_detail中的detail_value
      for (let i = 0; i < params.rule_detail.length; i++) {
        // 将rule_value_al数组，rule_value_ls数组合并到rule_value数组中
        if (params.rule_detail[i].is_diff_store === 1) {
          params.rule_detail[i].rule_value = [
            ...params.rule_detail[i].rule_value_al,
            ...params.rule_detail[i].rule_value_ls
          ];
          delete params.rule_detail[i].detail_value_al;
          delete params.rule_detail[i].detail_value_ls;
          delete params.rule_detail[i].rule_value_al;
          delete params.rule_detail[i].rule_value_ls;
        }
        delete params.rule_detail[i].detail_value;
      }
      const api = this.pageType === 1 ? getIntegralAddApi : getIntegralEditApi;
      const res: ApiResult = await api(params);

      if (res.code !== 200) return this.$message.error(res.message as string);

      this.$message.success("操作成功");

      this.goBack();
    } catch (e) {
      console.error("校验未通过", e);
    }
  }

  /** 取消 */
  private goBack() {
    this.$router.push({
      name: "myapp_wisdom_we-journey_integral-set-list"
    });
  }
  /** 规则渠道变更事件 */
  radioChange(index: number) {
    delete this.formData.rule_detail[index].detail_value_al;
    delete this.formData.rule_detail[index].detail_value_ls;
    delete this.formData.rule_detail[index].detail_value;
    this.$forceUpdate();
  }
  /** 生命周期 */
  private created() {
    this.pageType === 1 && this.generateInitOptions();
  }
}
</script>
<style lang="scss" scoped>
.inline-form {
  display: flex;
  align-items: center;
  .el-form-item {
    /deep/.el-form-item__content {
      margin-left: 0 !important;
    }
  }
  span {
    display: inline-block;
    height: 40px;
  }
  .small-select {
    width: 100px;
    margin: 0 10px;
    /deep/.el-input {
      width: 100px;
    }
  }
}

/deep/.el-input {
  width: 300px;
}

.bg-input {
  width: 170px !important;
  margin: 0 10px;
}

.model-title {
  font-size: 15px;
  font-weight: bold;
  position: relative;
  margin-left: 10px;
  margin-bottom: 15px;
  &::after {
    position: absolute;
    content: "";
    left: -15px;
    top: 2px;
    height: 15px;
    width: 5px;
    background: #000000;
  }
}

.category {
  background: #f9f9f9;
  padding: 10px;
  margin: 15px 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  .del-icon {
    color: red;
    cursor: pointer;
    font-size: 18px;
  }
}
.rules-content {
  .rules-title {
    color: #708df5;
  }
}
</style>
