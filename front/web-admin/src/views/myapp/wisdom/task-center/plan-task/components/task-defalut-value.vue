<template>
  <cy-form
    ref="defaultValueForm"
    :rules="rules"
    :model="defalueValueForm"
  >
    <!-- 标签名称 -->
    <cy-form-item
      label="标签名称"
      required
      prop="name"
    >
      <el-input
        placeholder="请输入标签名称"
        v-model="defalueValueForm.name"
        maxlength="30"
        clearable
      ></el-input>
    </cy-form-item>
    <!-- 标签有效期radio -->
    <cy-form-item
      label="标签有效期"
      required
    >
      <el-radio-group
        v-model="defalueValueForm.radio"
        @change="changeRadio"
      >
        <el-radio :label="1">任务结束后删除</el-radio>
        <el-radio :label="2">自定义删除时间</el-radio>
        <el-radio :label="3">长期有效</el-radio>
      </el-radio-group>
    </cy-form-item>
    <!-- 任务有效期结束XX天后，自动删除 -->
    <cy-form-item
      v-if="defalueValueForm.radio == 1"
      prop="days"
    >
      任务有效期结束
      <el-input
        style="width: 100px;"
        v-model="defalueValueForm.days"
        clearable
      ></el-input>
      天后，自动删除
    </cy-form-item>
    <!-- 有效期 -->
    <cy-form-item
      v-if="defalueValueForm.radio == 2"
      prop="date"
    >
      <el-date-picker
        v-model="defalueValueForm.date"
        type="datetime"
        placeholder="选择日期"
        style="width: 100%;"
        :picker-options="{
          disabledDate: (time) => {
            // 选择当天之后 不包含当天
            return time.getTime() < Date.now() - 8.64e6;
          }
        }"
      ></el-date-picker>
    </cy-form-item>
  </cy-form>
</template>

<script lang="ts">
import { Vue, Component, Prop, Watch } from "vue-property-decorator";
@Component({
  components: {}
})
export default class extends Vue {
  @Prop({ default: "" }) tagName!: string;
  @Prop({ default: false }) isDimissionChange!: boolean;
  @Prop() tagExpireForm!: object;
  @Watch("tagName", { immediate: true })
  onTagNameChange(val: any) {
    this.defalueValueForm.name = val;
  }
  @Watch("isDimissionChange", { immediate: true })
  onIsDimissionChange(val: any) {
    if (val) {
      // todo:重置表单
      this.defalueValueForm.radio = 1;
      this.defalueValueForm.date = "";
      this.defalueValueForm.days = 10;
      this.defalueValueForm.name = this.tagName;
    }
  }
  @Watch("tagExpireForm", { immediate: true, deep: true })
  onTagExpireFormChange(val: any) {
    this.defalueValueForm = val
  }
  // 控制弹框的关闭隐藏
  visibleDialog = true;
  defalueValueForm = {
    name: "",
    radio: 1,
    date: "",
    days: 10
  };
  rules = {
    // 标签名称验证
    name: [
      {
        message: "请输入标签名称",
        trigger: "blur"
      }
    ],
    // 天数验证只能是数字
    days: [
      {
        validator: (rule: any, value: any, callback: any) => {
          if (!value) {
            callback(new Error("请输入天数"));
          } else {
            // 只能是大于0数字
            const reg = /^[1-9]\d*$/;
            if (!reg.test(value)) {
              callback(new Error("请输入大于0的整数"));
            // 不能大于1000
            } else if(value > 1000){
              callback(new Error("请输入小于1000的整数"));
            }else{
              callback();
            }
          }
        },
        trigger: "blur"
      }
    ],
    date: [
      {
        validator: (rule: any, value: any, callback: any) => {
          if (!value) {
            callback(new Error("请选择日期"));
            // 只能选择大于当前日期+1天的日期
          } else if (new Date(value).getTime() < new Date().setHours(0, 0, 0, 0) + 8.64e7) {
            callback(new Error("请选择大于今天的日期"));
          } else {
            callback();
          }
        },
        trigger: "blur"
      }
    ]
  };
  // radio改变
  changeRadio(val: any) {
    if (val == 2) {
      this.defalueValueForm.date = "";
    } else if (val == 1) {
      this.defalueValueForm.days = 10;
    }
  }
  // 保存
  save() {
    console.log("保存");
    // 表单验证
    (this.$refs.defaultValueForm as any).validate((valid: any) => {
      if (valid) {
        console.log("验证通过");
        const { radio, date, days } = this.defalueValueForm;
        const params = {
          name: this.defalueValueForm.name,
          tag_expire_type: radio,
          tag_expire_time: 0
        };
        if (radio == 1) {
          params.tag_expire_time = Number(days);
        } else if (radio == 2) {
          params.tag_expire_time = new Date(date).getTime() / 1000;
        }
        this.$emit("save", params);
      } else {
        console.log("验证失败");
        return false;
      }
    });
  }
}
</script>
<style lang="scss" scoped></style>
