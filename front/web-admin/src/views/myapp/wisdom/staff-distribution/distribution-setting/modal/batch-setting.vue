<template>
  <el-dialog :visible.sync="dialogVisible" width="450px" @close="close">
    <template slot="title">
      <div class="dialog-title">
        <span class="title-text">批量设置</span>
        <span class="select-num">已选商品{{ selectSkuNum }}</span>
      </div>
    </template>
    <el-form ref="form" :model="form" :rules="rules" label-width="82px">
      <el-form-item
        v-if="display.includes('promotion_commission_amount')"
        label="推广佣金"
        prop="promotion_commission_amount"
      >
        <el-input
          v-model="form.promotion_commission_amount"
          placeholder="请输入推广佣金"
          @keyup.native="limitInput(form, 'promotion_commission_amount')"
          size="medium"
        >
          <template slot="append">%</template>
        </el-input>
      </el-form-item>
      <el-form-item
        v-if="display.includes('service_commission_amount')"
        label="服务佣金"
        prop="service_commission_amount"
      >
        <el-input
          v-model="form.service_commission_amount"
          placeholder="请输入服务佣金"
          @keyup.native="limitInput(form, 'service_commission_amount')"
          size="medium"
        >
          <template slot="append">%</template>
        </el-input>
      </el-form-item>
      <el-form-item
        v-if="display.includes('date_range')"
        label="分销期限"
        prop="date_range"
      >
        <el-date-picker
          v-model="form.date_range"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="timestamp"
          :picker-options="pickerOption"
          size="medium"
          style="width: 308px !important;"
        >
        </el-date-picker>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button size="mini" @click="close">取 消</el-button>
      <el-button size="mini" type="primary" @click="save">确 认</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import { handleDisabledDate } from "../date-picker";
@Component
export default class BatchSetting extends Vue {
  @Prop(Boolean) visible!: boolean;
  @Prop({ default: 0 }) selectSkuNum?: number;
  @Prop({ default: () => [] }) display?: Array<any>;

  private dialogVisible = true; // 控制模态框显示/隐藏
  private pickerOption = {
    disabledDate: handleDisabledDate()
  };
  private form = {};
  private rules = {};

  // 限制小数点位数2与值范围0~100
  private limitInput(row: any, key: string) {
    let value = row[key];
    value = value > 100 ? 100 : value;
    value = value < 0 ? 0 : value;
    value = value + "";
    if (value.length > 1 && value[0] === "0" && value[1] !== ".") {
      value = value.substring(1);
    }
    row[key] = value
      .replace(/[^\d.]/g, "")
      .replace(/^(\d+)\.(\d\d).*$/, "$1.$2");
  }

  @Emit()
  private close() {
    this.dialogVisible = false;
  }

  private save() {
    (this.$refs["form"] as any).validate((valid: boolean) => {
      if (valid) {
        this.submit();
      }
    });
  }

  @Emit()
  private submit() {
    return this.form;
  }

  onDisplayChange() {
    this?.display?.forEach(item => {
      if (item === "date_range") {
        this.$set(this.form, item, []);
        this.rules[item] = [
          { required: true, message: "请输入分销期限", trigger: "blur" }
        ];
      } else {
        this.$set(this.form, item, 0);
        this.rules[item] = [
          {
            required: true,
            message:
              item === "promotion_commission_amount"
                ? "请输入推广佣金比例"
                : "请输入服务佣金比例",
            trigger: "blur"
          }
        ];
      }
    });
  }

  private mounted() {
    this.onDisplayChange();
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog__header {
  font-weight: 400;
  .dialog-title {
    display: flex;
    align-items: center;
    .title-text {
      line-height: 26px;
      font-size: 18px;
      font-weight: 500;
    }
    .select-num {
      font-size: 14px;
      font-weight: 400 !important;
      margin-left: 30px;
    }
  }
}
/deep/.el-dialog__body {
  padding-bottom: 0 !important;
}
</style>
