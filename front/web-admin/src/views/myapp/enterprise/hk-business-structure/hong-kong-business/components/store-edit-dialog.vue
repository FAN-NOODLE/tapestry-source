<template>
  <el-dialog
    title="编辑门店"
    :visible.sync="visible"
    width="640px"
    custom-class="hk-store-edit-dialog"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="form"
      :model="formData"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="企微部门 ID">
        <el-input :value="formData.wechat_department_id" disabled></el-input>
      </el-form-item>
      <el-form-item label="门店编码">
        <el-input :value="formData.store_code" disabled></el-input>
      </el-form-item>
      <el-form-item label="门店名称" prop="store_name">
        <el-input v-model.trim="formData.store_name" maxlength="100"></el-input>
      </el-form-item>
      <el-form-item label="区域" prop="area">
        <el-select v-model="formData.area" placeholder="请选择区域">
          <el-option label="香港" value="hk"></el-option>
          <el-option label="澳门" value="mo"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="门店类型">
        <el-select v-model="formData.store_type" placeholder="请选择">
          <el-option label="未设置" :value="0"></el-option>
          <el-option label="奥莱" :value="1"></el-option>
          <el-option label="零售" :value="2"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="英文名">
        <el-input v-model.trim="formData.inventory_name" maxlength="100"></el-input>
      </el-form-item>
      <el-form-item label="所属城市">
        <el-input v-model.trim="formData.city" maxlength="100"></el-input>
      </el-form-item>
      <el-form-item label="门店地址">
        <el-input v-model.trim="formData.address" maxlength="200"></el-input>
      </el-form-item>
      <el-form-item label="联系电话">
        <el-input v-model.trim="formData.phone" maxlength="30"></el-input>
      </el-form-item>
      <el-form-item label="联系邮箱" prop="email">
        <el-input v-model.trim="formData.email" maxlength="100"></el-input>
      </el-form-item>
    </el-form>
    <template slot="footer">
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">
        保存
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import Api from "@/api/enterprise/hong-kong-business";

@Component
export default class StoreEditDialog extends Vue {
  @Prop({ default: false }) value!: boolean;
  @Prop({ default: () => ({}) }) store!: any;

  visible = false;
  saving = false;
  formData: any = {};
  rules = {
    store_name: [
      { required: true, message: "请输入门店名称", trigger: "blur" }
    ],
    area: [
      { required: true, message: "请选择区域", trigger: "change" }
    ],
    email: [
      { type: "email", message: "请输入正确的邮箱地址", trigger: "blur" }
    ]
  };

  @Watch("value", { immediate: true })
  handleVisibleChange(value: boolean) {
    this.visible = value;
  }

  @Watch("store", { immediate: true, deep: true })
  handleStoreChange(value: any) {
    this.formData = {
      id: value.id,
      wechat_department_id: value.wechat_department_id,
      store_code: value.store_code || "",
      store_name: value.store_name || "",
      area: value.area || "",
      store_type: value.store_type || 0,
      inventory_name: value.inventory_name || "",
      address: value.address || "",
      city: value.city || "",
      phone: value.phone || "",
      email: value.email || ""
    };
  }

  handleClose() {
    this.$emit("input", false);
  }

  handleSave() {
    (this.$refs.form as any).validate(async (valid: boolean) => {
      if (!valid) return;
      this.saving = true;
      try {
        const {
          id,
          store_name,
          area,
          store_type,
          inventory_name,
          address,
          city,
          phone,
          email
        } = this.formData;
        const res: any = await Api.saveStore({
          id,
          store_name,
          area,
          store_type,
          inventory_name,
          address,
          city,
          phone,
          email
        });
        if (res.code !== 200) {
          this.$message.error(res.message || "保存失败");
          return;
        }
        this.$message.success("保存成功");
        this.$emit("success");
        this.handleClose();
      } finally {
        this.saving = false;
      }
    });
  }
}
</script>

<style lang="scss" scoped>
.el-select {
  width: 100%;
}

/deep/ .hk-store-edit-dialog {
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 80px);

  .el-dialog__header,
  .el-dialog__footer {
    flex-shrink: 0;
  }

  .el-dialog__body {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
  }
}
</style>
