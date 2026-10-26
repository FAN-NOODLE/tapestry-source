<template>
  <div class="stall-detail">
    <layout-form>
      <template slot="content">
        <div class="row">
          <span class="title">员工头像</span>
          <span class="right img">
            <img :src="detailList.avatar ? detailList.avatar : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'" alt="" />
          </span>
        </div>
        <div class="row">
          <span class="title">员工姓名</span>
          <span class="right">{{ detailList.name }}</span>
        </div>
        <div class="row">
          <span class="title">员工性别</span>
          <span class="right">{{ detailList.gender_text }}</span>
        </div>
        <div class="row">
          <span class="title">所属部门</span>
          <span class="right">{{ detailList.forms }}</span>
        </div>
        <div class="row">
          <span class="title">手机号码</span>
          <span class="right">{{ detailList.mobile }}</span>
        </div>
        <div class="row">
          <span class="title">加入状态</span>
          <span class="right">{{ detailList.status_text }}</span>
        </div>
        <div class="row">
          <span class="title">加入时间</span>
          <span class="right" v-if="detailList.created_at">{{
            detailList.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</span>
          <span class="right" v-else>-</span>
        </div>
      </template>
      <template slot="footer">
        <el-button size="small" @click="goBack()">
          返回
        </el-button>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/enterprise";
@Component({})
export default class extends Vue {
  // 数据变量
  private detailList: Array<any> = [];

  // 详情数据
  private memberDetail(id: any) {
    Api.memberDetailApi({ id: id }).then((res: any) => {
      if (res.code == 200) {
        this.detailList = res.data;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 返回上一级
  private goBack() {
    // this.$router.push({
    //   name: "myapp_enterprise_organize-structure_staff-list",
    //   query: {
    //     page: this.$route.query.page
    //   }
    // });
    this.$router.go(-1);
  }
  private mounted() {
    const id = Number(this.$route.query.id);
    this.memberDetail(id);
  }
}
</script>
<style lang="scss" scoped>
.row {
  margin-bottom: 30px;
  align-items: center;
  display: flex;
  color: $color-text-main;
  .title {
    margin-right: 20px;
    color: $color-text-normal;
  }
}
img {
  width: 30px;
  height: 30px;
  border-radius: 50%;
}
</style>
