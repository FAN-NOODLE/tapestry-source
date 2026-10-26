<template>
  <div class="guide_set">
    <layout-form>
      <template slot="content">
        <el-form
          :model="form"
          label-position="left"
          :rules="rules"
          ref="form"
          label-width="120px"
          class="demo-ruleForm"
          :validate-on-rule-change="false"
        >
          <cy-form-item
            label="绑定时效"
            prop="duration_time"
            desc="若设置绑定时效有效期，则时长不可超过365天"
            class="cy-item-w-100 align"
          >
            <el-radio-group
              v-model="form.binding_type"
              size="small"
              :disabled="!is_edit"
              @change="changeRadio(1)"
            >
              <el-radio
                v-for="item in time_list"
                :key="item.key"
                :label="item.key"
                >{{ item.name }}</el-radio
              >
            </el-radio-group>
            <cy-input
              v-model="form.duration_time"
              width-size="mini"
              :disabled="!is_edit || form.binding_type !== 2"
              class="search-input"
              @input="_checkPositive(form, 'duration_time', form.duration_time)"
            >
              <div slot="append">天</div>
            </cy-input>
          </cy-form-item>

          <cy-form-item
            size="medium"
            label-width="120px"
            label="是否允许他人抢客"
            prop="protected_time"
            class="protected-time cy-item-w-100 align"
            desc="保护期，当顾客与导购绑定后，x天内不允许被抢客，最小可设置为0天，且保护期不可大于绑定时间有效期"
          >
            <el-radio-group
              v-model="form.change_owner"
              size="small"
              :disabled="!is_edit"
              @change="changeRadio(2)"
            >
              <el-radio
                v-for="item in allow_list"
                :key="item.key"
                :label="item.key"
                >{{ item.name }}</el-radio
              >
            </el-radio-group>
            <div class="item-inline">
              <cy-input
                v-model="form.protected_time"
                width-size="mini"
                :disabled="!is_edit || form.change_owner !== 2"
                class="search-input"
                @input="
                  _checkInt(form, 'protected_time', form.protected_time + '')
                "
              >
                <div slot="append">天</div>
              </cy-input>
              <span style="color:#606266">内不允许抢客 </span>
            </div>
          </cy-form-item>

          <cy-form-item
            label="导购购买商品"
            desc="选中计算业绩后，所有导购和自己绑定客户关系，自己购买商品会为自己计算业绩"
            class="cy-item-w-100"
          >
            <el-radio-group
              v-model="form.commission_type"
              size="small"
              :disabled="!is_edit"
            >
              <el-radio
                v-for="item in money_list"
                :key="item.key"
                :label="item.key"
                >{{ item.name }}</el-radio
              >
            </el-radio-group>
            <div class="remind-info"></div>
          </cy-form-item>

          <cy-form-item
            label="裂变分享"
            desc=" 选中计算业绩后，客户将导购分享的内容再次分享给其他客户，其他客户购买为导购计算业绩"
            class="cy-item-w-100"
          >
            <el-radio-group
              v-model="form.fission_type"
              size="small"
              :disabled="!is_edit"
            >
              <el-radio
                v-for="item in set_list"
                :key="item.key"
                :label="item.key"
                >{{ item.name }}</el-radio
              >
            </el-radio-group>
          </cy-form-item>
        </el-form>
      </template>
      <template slot="footer">
        <el-button v-if="!is_edit" type="primary" size="small" @click="toEdit"
          >编辑
        </el-button>

        <el-button v-if="is_edit" type="primary" size="small" @click="saveData"
          >保存</el-button
        >
        <el-button v-if="is_edit" size="small" @click="cancel">取消</el-button>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Mixins } from "vue-property-decorator";
import Api from "@/api/wisdom/Guide";
import NumberInput from "@/mixins/NumberInput";
// import personApi from "@/api/wisdom/Personnel";

@Component
export default class extends Mixins(NumberInput) {
  private is_edit = false; // 保存取消按钮是否展示
  private form = {
    binding_type: 1, // 绑定时效类型
    duration_time: "", // 有效期天数
    change_owner: 1, // 是否允许他人抢客类型
    protected_time: "", // 保护期天数
    commission_type: 2, // 导购购买商品是否计算业绩
    fission_type: 1 // 裂变分享是否计算业绩
  };
  private rules = {
    duration_time: [{ validator: this.checkDuration }],
    protected_time: [{ validator: this.checkProtect }]
  };
  private time_list = [
    { name: "永久有效", key: 1 },
    { name: "有效期", key: 2 }
  ];
  private allow_list = [
    { name: "不允许", key: 1 },
    { name: "允许,保护期", key: 2 }
  ];
  private money_list = [
    { name: "计算业绩", key: 1 },
    { name: "不计算业绩", key: 2 }
  ];
  private set_list = [
    { name: "计算业绩", key: 2 },
    { name: "不计算业绩", key: 1 }
  ];
  private toEdit() {
    this.is_edit = true;
    console.log(this.is_edit);
  }
  private changeRadio(type: any) {
    switch (type) {
      case 1:
        this.form.duration_time = "";
        break;
      case 2:
        this.form.protected_time = "";
        break;
    }
  }
  //校验
  private checkDuration(rule: any, value: any, callback: any) {
    if (this.form.binding_type === 2 && value === "") {
      return callback(new Error(`请输入有效期`));
    }
    if (value > 365) {
      return callback(new Error(`有效期天数不能超过365天`));
    }
    callback();
  }
  private checkProtect(rule: any, value: any, callback: any) {
    if (this.form.change_owner === 2 && value === "") {
      return callback(new Error("请输入保护期"));
    }
    if (
      this.form.binding_type === 2 &&
      this.form.duration_time &&
      Number(value) > Number(this.form.duration_time)
    ) {
      return callback(new Error("保护期不可大于绑定时效有效期"));
    }
    callback();
  }
  // 获取设置详情
  private getDetail() {
    Api.achievemenDetail({}).then(res => {
      const { data } = res;
      for (const key in data) {
        // if (this.form.hasOwnProperty(key)) {
        this.form[key] = data[key];
        // }
      }
      if (!this.form.change_owner) this.form.change_owner = 1;
    });
  }
  // 保存前校验
  private saveData() {
    (this.$refs["form"] as any).validate((valid: any) => {
      if (valid) {
        const params = { ...this.form };
        this.saveApi(params);
      }
    });
  }
  // 保存数据接口
  private saveApi(params: any) {
    params.duration_time = parseInt(params.duration_time || 0);
    params.protected_time = parseInt(params.protected_time || 0);
    Api["achievemenSave"](params).then(() => {
      this.$message.success("保存成功");
      this.is_edit = false;
    });
  }
  cancel() {
    this.getDetail();
    this.is_edit = false;
    // this.reload();
  }
  // 生命周期
  private created() {
    console.log("created");
  }
  private mounted() {
    this.getDetail();
    console.log("mounted");
  }
  private destroyed() {
    console.log("destroyed");
  }
}
</script>
<style lang="scss" scoped>
.demo-ruleForm {
  /deep/ .label {
    max-width: unset;
    padding-left: 0;
  }
  // /deep/ .align .el-form-item__content {
  //   display: flex;
  //   align-items: center;
  // }
}
.item-inline {
  display: inline-block;
}
.select-box {
  display: inline-block;
  max-width: 231px;
  line-height: 1.5;
  margin: 3px 10px 3px 0;
  padding: 3px 11px 3px 7px;
  border: 1px solid #e4e6e9;
  border-radius: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  vertical-align: middle;
  img {
    width: 20px;
    vertical-align: top;
    &.del-icon {
      cursor: pointer;
    }
  }
  .label {
    margin: 0 5px;
  }
}
.update {
  vertical-align: middle;
}
.search-input {
  width: 180px;
  margin: 0 10px;
}
.remind-info {
  font-size: 12px;
  color: #5d5d5d;
  line-height: 20px;
}
.line-info {
  display: inline-block;
  vertical-align: text-top;
}
.button-wrap {
  margin-left: 120px;
}
.cy-item-w-100 {
  width: 100%;
}
/deep/.el-radio-group {
  margin-bottom: 4px;
}
.demo-ruleForm /deep/ .desc {
  margin-left: 120px !important;
}
// .duration-time{
//   /deep/.el-form-item__error;
//   left: 184px !important
//   }
// .protected-time
//   /deep/.el-form-item__error
//   left: 200px !important
</style>
