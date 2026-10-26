<template>
  <div class="content-block">
    <el-form ref="form" label-position="left" label-width="150px">
      <el-form-item label="企业名称">
        <span> {{ formData.enterprise_name }}</span>
      </el-form-item>
      <el-form-item label="企业简称">
        <span> {{ formData.short_name }}</span>
      </el-form-item>
      <el-form-item label="企业logo">
        <el-image fit="cover" :src="formData.logo" class="img-src"> </el-image>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="goEdit" size="small">编辑</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/setting/public";
interface IDataForm {
  enterprise_name?: string; // 企业名称
  short_name?: string; // 企业简介
  logo?: string; // 企业logo
}
@Component
export default class extends Vue {
  private formData: IDataForm = {
    enterprise_name: "",
    short_name: "",
    logo: ""
  };
  // 编辑商户信息
  private goEdit() {
    this.$router.push({
      name: "myapp_setting_merchant-info_edit"
    });
  }
  // 商户信息接口调用
  private getDetail() {
    Api.merchantDetail().then((res: any) => {
      const data: IDataForm = res.data;
      this.formData = { ...data };
    });
  }
  mounted() {
    this.getDetail();
  }
}
</script>
<style lang="scss" scoped>
.content-block {
  background-color: $color-white;
  padding: 30px;
  border-radius: 10px;
  .img-src {
    width: 100px;
    height: 100px;
  }
}
</style>
