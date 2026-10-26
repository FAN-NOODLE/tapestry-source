<template>
  <layout-form>
    <template slot="content">
      <el-form
        :model="formData"
        :rules="rules"
        ref="ruleForm"
        label-width="120px"
        label-position="right"
      >
        <!-- 基础配置 -->
        <div class="model-title">基础配置</div>
        <el-form-item
          label="奖章规则名称"
          prop="name"
        >
          <el-input
            :disabled="pageType === 3"
            v-model="formData.name"
            maxlength="15"
            show-word-limit
          ></el-input>
        </el-form-item>

        <el-form-item
          label="奖章来源"
          required
        >
          <el-select
            v-model="formData.source"
            @change="sourceChange"
            :disabled="pageType === 3"
          >
            <el-option
              v-for="item in conditions.sourceOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item
          label="规则有效期"
          required
        >
          <el-radio-group
            v-model="formData.rule_type"
            @change="ruleTypeChange"
          >
            <el-radio
              :label="1"
              :disabled="formData.source === 1 || pageType === 3"
            >长期有效</el-radio>
            <el-radio
              :label="2"
              :disabled="formData.source === 1 || pageType === 3"
            >指定时间</el-radio>
            <el-radio
              :label="3"
              :disabled="formData.source === 2 || pageType === 3"
            >仅生效一次</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 指定时间 -->
        <el-form-item
          v-if="formData.rule_type === 2"
          label=""
          prop="time"
        >
          <!-- todo: 如果小于当前时间禁用 -->
          <el-date-picker
            :disabled="pageType === 3 || formData.time[0] < new Date().getTime()"
            v-model="formData.time[0]"
            type="datetime"
            placeholder="开始日期"
            :picker-options="{ disabledDate: disabledStartDate(new Date()) }"
          >
          </el-date-picker>
          <span style="padding: 0 10px;">至</span>
          <el-date-picker
            :disabled="pageType === 3"
            v-model="formData.time[1]"
            type="datetime"
            placeholder="结束日期"
            :picker-options="{ disabledDate: disabledEndDate(formData.time[0]) }"
          >
          </el-date-picker>
        </el-form-item>

        <!-- 奖章配置 -->
        <div class="model-title">奖章配置</div>
        <el-form-item
          label="奖章类型"
          prop="meda_type"
          required
        >
          <el-select
            v-model="formData.meda_type"
            :disabled="pageType === 3"
          >
            <el-option
              v-for="item in conditions.medalTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item
          label="奖章描述"
          prop="meda_description"
        >
          <el-input
            :disabled="pageType === 3"
            type="textarea"
            :rows="2"
            placeholder="请输入奖章描述"
            v-model="formData.meda_description"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>

        <el-form-item
          label="其它奖励"
          required
        >
          <el-radio-group
            v-model="formData.other_rewards"
            @change="otherRewardChange"
          >
            <el-radio
              :disabled="pageType === 3"
              v-for="item in conditions.otherOptions"
              :key="item.value"
              :label="item.value"
            >{{ item.label }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 其它奖励 -->
        <el-form-item label="">
          <!-- 奖励积分 -->
          <el-form-item
            class="integral"
            label-width="100px"
            label="增加"
            prop="integral"
            v-if="formData.other_rewards === 2 || formData.other_rewards === 4"
          >
            <el-input
              :disabled="pageType === 3"
              class="bg-input"
              v-model="formData.integral"
              oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
              placeholder="请输入正整数"
            ></el-input>
            <span>积分</span>
          </el-form-item>
          <!-- 仅奖励实物奖品 -->
          <el-form-item
            label-width="100px"
            label="奖品名称"
            prop="item_name"
            v-if="formData.other_rewards === 3 || formData.other_rewards === 4"
          >
            <el-input
              :disabled="pageType === 3"
              v-model="formData.item_name"
              placeholder="请输入奖品名称"
              maxlength="15"
              show-word-limit
            ></el-input>
            <div class="tips">
              备注：实物奖品在前端仅做名称展示，具体奖品需品牌在线下完成发放
            </div>
          </el-form-item>
        </el-form-item>

        <el-form-item
          label="奖章规则维度"
          required
          v-if="formData.source === 2"
        >
          <el-select
            :disabled="pageType === 3"
            v-model="formData.rule_dimension_id"
            @change="dimensionVisibleChange"
          >
            <el-option
              v-for="item in conditions.ruleDimension"
              :key="item.id"
              :label="item.rule_name"
              :value="item.id"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item
          label="奖章规则类目"
          v-if="formData.source === 2"
          prop="rule_category_id"
        >
          <el-select
            :disabled="pageType === 3"
            v-model="formData.rule_category_id"
            @change="categoryValueChange"
          >
            <el-option
              v-for="item in conditions.ruleClasses"
              :key="item.id"
              :label="item.rule_name"
              :value="item.id"
            >
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item
          label=""
          v-if="formData.source === 2"
          prop="invite_top"
        >
          <span class="invite-m">{{ description }}</span>
          <el-input
            :disabled="pageType === 3"
            class="bg-input"
            v-model="formData.invite_top"
            oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
            placeholder="请输入正整数"
          ></el-input>
          <!-- toto：提示 -->
          <div
            class="medal-tips"
            v-if="formData.rule_dimension_id === 1"
          >* 该奖章实际发放规则会区分零售排行与奥莱排行</div>
        </el-form-item>

        <!-- 避免最后一个校验不生效 -->
        <div style="margin-top: 20px"></div>

        <!-- 人群包 -->
        <el-form-item
          label="获奖人群包"
          v-if="formData.source === 1"
          required
        >
          <div class="upload-box">
            <large-file-upload
              class="upload-box-upload"
              v-model="formData.file_url"
              :onSuccess="uploadSuccess"
              :onRemove="handleRemove"
              accept=".xls,.xlsx"
            >
              <el-button
                size="mini"
                type="primary"
              >上传人群包</el-button>
            </large-file-upload>
            <el-button
              size="mini"
              type="text"
              class="upload-box-download"
              @click="download"
            >下载上传模板</el-button>
          </div>
        </el-form-item>

        <el-form-item
          label="奖章样式"
          required
        >
          <medal-style
            :disabled="pageType === 3"
            :info="formData.styles"
            :type="medalStyleOptions.type"
            ref="medalStyle"
          ></medal-style>
        </el-form-item>
      </el-form>
    </template>

    <!-- 底部按钮 -->
    <template slot="footer">
      <el-button
        type="primary"
        @click="save"
        v-if="pageType !== 3"
      >保存</el-button>
      <el-button @click="goBack">{{
        pageType === 3 ? "返回" : "取消"
      }}</el-button>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Mixins, Prop, Watch } from "vue-property-decorator";
import { sourceOptions, medalTypeOptions, otherOptions } from "../../config";
import {
  getRuleOptionsApi,
  medalAddApi,
  medalEditApi,
  downloadCrowPackApi,
  readCrowPackApi
} from "@/api/we-journey/medal";
import LargeFileUpload from "@/components/business/large-file-upload/index_o.vue";
import MedalStyle from "./medal-style.vue";
import { ApiResult } from "@/types";
import DownLoad from "@/mixins/DownLoad";

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

/** form类型 */
type formType = {
  create_time: number;
  /** 奖章来源 */
  source: number;
  /** 规则有效期 */
  rule_type: number;
  time?: Array<Date>;
  /** 规则开始时间 */
  rule_start_time: number | string;
  /** 规则结束时间 */
  rule_end_time: number | string;
  /** 奖章类型 */
  meda_type: number;
  /** 奖章描述 */
  meda_description: string;
  /** 其他奖励 */
  other_rewards: number;
  /** 积分 */
  integral: string;
  /** 奖品名称 */
  item_name: string;
  /** 奖章规则维度  */
  rule_dimension_id: string | number;
  /** 奖章规则类目 */
  rule_category_id: string | number;
  /** TOP统计多少名榜单 */
  invite_top: string;
  /** 奖章样式 */
  styles: Array<any>;
  /** 上传人群包 */
  user_ids: Array<any>;
};

/** 奖章规则类型 */
type ruleType = {
  id: number;
  /** 规则名称 */
  rule_name: string;
  /** 时间类型 0无,1月,2季度,3年 */
  time_type: number;
  /** 描述 */
  description: string;
};

@Component({
  components: { MedalStyle, LargeFileUpload }
})
export default class Medal extends Mixins(DownLoad) {
  /** 详情 */
  @Prop() info?: formType;

  /** 页面状态 1. 新增 2.编辑 3.查看 */
  @Prop() pageType?: number;

  @Watch("info")
  async changeInfo() {
    if (this.info) {
      let time: Array<Date> = [];

      // 奖章来源 = 规则统计
      if (this.info.source === 2) {
        // 处理时间
        time =
          this.info.source === 2 && this.info.rule_start_time
            ? [
              new Date((this.info.rule_start_time as number) * 1000),
              new Date((this.info.rule_end_time as number) * 1000)
            ]
            : [];

        // 获取奖章规则维度（用做回显）
        await this.getRuleOptions();
        // 获取奖章规则类目（用做回显）
        await this.getRuleOptions(this.info.rule_dimension_id as number, () => {
          const id = (this.info as formType).rule_category_id as number;
          this.formData.rule_category_id = id;

          this.categoryValueChange(id);
        });
      }

      this.formData = {
        ...this.info,
        time
      };
    }
  }

  /** 配置项 */
  private conditions = {
    sourceOptions,
    medalTypeOptions,
    otherOptions,
    ruleClasses: [] as Array<ruleType>,
    ruleDimension: [] as Array<ruleType>
  };

  /** medalStyle 子组件配置参数 */
  private medalStyleOptions = {
    type: 0
  };

  /** 表单数据 */
  private formData: formType = {
    create_time: 0,
    /** 奖章来源 */
    source: 1,
    /** 规则有效期 */
    rule_type: 3,
    time: [],
    /** 规则开始时间 */
    rule_start_time: "",
    /** 规则结束时间 */
    rule_end_time: "",
    /** 奖章类型 */
    meda_type: 1,
    /** 奖章描述 */
    meda_description: "",
    /** 其他奖励 */
    other_rewards: 1,
    /** 积分 */
    integral: "",
    /** 奖品名称 */
    item_name: "",
    /** 奖章规则维度  */
    rule_dimension_id: "",
    /** 奖章规则类目 */
    rule_category_id: "",
    /** TOP统计多少名榜单 */
    invite_top: "",
    /** 奖章样式 */
    styles: [],
    /** 上传人群包 */
    user_ids: []
  };

  /** 校验规则 */
  private rules = {
    name: [{ required: true, message: "请输入规则名称", trigger: "change" }],
    meda_description: [
      { required: true, message: "请输入奖章描述", trigger: "change" }
    ],
    integral: [
      { required: true, message: "请输入奖励积分", trigger: "change" }
    ],
    item_name: [
      { required: true, message: "请输入奖品名称", trigger: "change" }
    ],
    time: [{ required: true, validator: validateTime, trigger: "change" }],
    invite_top: [{ required: true, message: "请输入正整数", trigger: "change" }],
    // todo:请选择奖章规则类目
    rule_category_id: [{ required: true, message: "请选择奖章规则类目", trigger: "change" }]
  };

  /** 描述 */
  private description = "";

  /** 其他奖励切换 */
  private otherRewardChange() {
    this.formData.integral = "";
    this.formData.item_name = "";
    this.$nextTick(() => {
      (this.$refs["ruleForm"] as HTMLFormElement).clearValidate("integral");
      (this.$refs["ruleForm"] as HTMLFormElement).clearValidate("item_name");
    });
  }
  /** 处理时间 - 之前的时间不可选 */
  private disabledStartDate(time: Date) {
    // todo: 开始时间之前的不能选
    return (date: Date) => {
      return date.getTime() < time.getTime();
    };
  }
  /** 处理时间 - 之前的时间不可选 */
  private disabledEndDate(time: Date) {
    if (!time) return () => false;
    // todo: 开始时间之前的不能选
    // 判断time是否<当前时间
    if (time.getTime() < new Date().getTime()) {
      return (date: Date) => {
        return date.getTime() < new Date().getTime();
      };
    } else {
      return (date: Date) => {
        return date.getTime() < time.getTime();
      };
    }
  }

  /** 规则来源变化，自动更改规则有效期 */
  private async sourceChange() {
    if (this.formData.source === 1) {
      // source === 1 手动上传
      this.formData.rule_type = 3;
      this.medalStyleOptions.type = 0;
    } else {
      // 规则统计
      this.formData.rule_type = 1;
      // 获取规则下拉
      await this.getRuleOptions();
      // 获取奖章规则类目
      this.getRuleOptions(this.formData.rule_dimension_id as number, () => {
        this.categoryValueChange(this.formData.rule_category_id as number);
      });
    }
  }

  /** 规则有效期变化 - 切换有效期清空时间 */
  private ruleTypeChange() {
    // todo: 如果是长期有效
    if (this.formData.rule_type === 2) {
      if (this.pageType === 2) {
        if ((this.info as any).create_time > 0) {
          this.formData.time = [new Date((this.info as any).create_time * 1000), new Date()]
          console.log(this.formData.time)
        }
      }
    } else {
      this.formData.time = [];
    }
  }

  /** 奖章规则类目变化 - 值变化 */
  private categoryValueChange(value: number) {
    const checkedItem = this.conditions.ruleClasses.filter(
      item => item.id === value
    );
    if (!checkedItem.length) return;

    // 确定奖章样式变化
    this.medalStyleOptions.type = checkedItem[0].time_type;

    this.description = checkedItem[0].description;
  }

  /** 奖章规则维度变化 */
  private dimensionVisibleChange() {
    this.getRuleOptions(this.formData.rule_dimension_id as number, () => {
      this.categoryValueChange(this.formData.rule_category_id as number);
    });
  }

  /** 下载上传模板 - 人群包 */
  private async download() {
    try {
      const res = await downloadCrowPackApi();
      const nowTime = new Date().getTime();
      const fileName = "上传人群包模板_" + nowTime + ".xlsx";
      this.downLoadFile(res, fileName);
    } catch (e) {
      console.error(e);
    }
  }

  /** 移除人群包信息 */
  private handleRemove() {
    this.formData.user_ids = [];
  }

  /** 上传人群包 */
  private async uploadSuccess({ data }: { data: { url: string } }) {
    const res: ApiResult = await readCrowPackApi({ file_url: data.url });

    if (res.code !== 200) return this.$message.error(res.message as string);

    this.formData.user_ids = res.data;
  }

  /** 组合校验规则 */
  private composeValidate() {
    const childDom = this.$refs["medalStyle"] as HTMLFormElement;

    // 获取子组件参数
    const childData = childDom.getData();
    // 子组件校验
    const childPro = new Promise(resolve => {
      (childDom.$refs["ruleForm"] as HTMLFormElement).validate(
        (valid: boolean) => {
          this.formData.styles = childData;
          if (valid) resolve(valid);
        }
      );
    });

    // 父组件校验整合
    const parentPro = new Promise(resolve => {
      (this.$refs["ruleForm"] as HTMLFormElement).validate((valid: boolean) => {
        if (valid) resolve(valid);
      });
    });

    return Promise.all([childPro, parentPro]);
  }

  /** 保存 */
  private async save() {
    try {
      await this.composeValidate();
      const { time } = this.formData;
      const params = {
        ...this.formData,
        integral: Number(this.formData.integral),
        rule_dimension_id: Number(this.formData.rule_dimension_id),
        rule_category_id: Number(this.formData.rule_category_id),
        invite_top: Number(this.formData.invite_top),
        user_ids: this.formData.user_ids || [],
        rule_start_time:
          time && time.length ? ~~((time[0] as Date).getTime() / 1000) : 0,
        rule_end_time:
          time && time.length ? ~~((time[1] as Date).getTime() / 1000) : 0
      };

      delete params.time;

      const api = this.pageType === 1 ? medalAddApi : medalEditApi;
      const res: ApiResult = await api(params);

      if (res.code !== 200) return this.$message.error(res.message as string);

      this.$message.success("保存成功");
      this.goBack();
    } catch (e) {
      console.error("校验未通过", e);
    }
  }

  /** 取消-返回上一页 */
  private goBack() {
    this.$router.push({
      name: "myapp_wisdom_we-journey_medal-set-list"
    });
  }

  /** 获取下拉列表, 当页面初始化调取，获取奖章规则维度 更改之后调用 */
  private async getRuleOptions(pid = 0, fn?: () => void) {
    const res: ApiResult = await getRuleOptionsApi({ pid });

    if (res.code !== 200) return this.$message.error(res.message as string);

    const list = res.data.list || [];
    // 当pid === 0 奖章规则维度, 初始化获取
    if (!pid) {
      // 奖章规则维度
      this.conditions.ruleDimension = list || [];
      // 奖章规则维度id 默认赋值
      this.formData.rule_dimension_id = list[0]?.id;
    } else {
      // 奖章规则类目
      this.conditions.ruleClasses = list || [];
      // 奖章规则类目id 默认赋值
      this.formData.rule_category_id = list[0]?.id;

      // 描述 默认赋值
      this.description = list[0]?.description;
    }

    // 如果是编辑，这里需要重新赋值
    fn && fn();
  }

  /** 生命周期 */
  private created() {
    console.log("生命周期");
  }
}
</script>
<style lang="scss" scoped>
.medal-tips {
  color: red;
  margin-top: 10px;
}

/deep/.el-input {
  width: 300px;
}

/deep/.el-textarea {
  width: 300px;
}

.model-title {
  font-size: 16px;
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

.tips {
  color: #909399;
  font-size: 12px;
}

.integral {
  margin-bottom: 15px;

  span {
    margin-left: 10px;
  }
}

.invite-m {
  margin-right: 10px;
  color: #606266;
}

.upload-box {
  display: flex;
  align-items: center;

  &-download {
    color: #409eff;
    margin-left: 15px;
  }

  &-upload {
    /deep/.progress {
      width: 500px;
      margin-top: 7px;
    }
  }
}
</style>
