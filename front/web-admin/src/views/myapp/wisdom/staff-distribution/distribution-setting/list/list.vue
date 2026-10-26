<template>
  <div>
    <div
      v-if="showDeployTip"
      class="content-wrap"
    >
      <div class="align-item-center">
        <img
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/ico_tip.png"
          alt=""
        />
        全员分销可在企业微信聊天工具栏(侧边栏)使用，
        <span
          class="c-blue"
          @click="checkDeploy"
        >点击查看如何配置</span>
      </div>
      <div @click="closeDeploy"><i class="el-icon-close"></i></div>
    </div>

    <layout-list>
      <template slot="top-tabs">
        <el-tabs
          v-model="activeName"
          @tab-click="handleTabClick"
        >
          <el-tab-pane
            v-for="tab in tabs"
            :key="tab.name"
            :label="tab.title"
            :name="tab.name"
          >
            <component :is="tab.name"></component>
          </el-tab-pane>
        </el-tabs>
      </template>
    </layout-list>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import BaseSetting from "./base-setting.vue";
import DistributionGoods from "./distribution-goods.vue";
@Component({
  components: {
    BaseSetting,
    DistributionGoods,
  },
})
export default class extends Vue {
  private showDeployTip = true; // 全员分销配置提示展示flag
  private tabs: Array<any> = [
    { title: "基础设置", name: "BaseSetting" },
    { title: "分销商品", name: "DistributionGoods" },
  ];
  private activeName = "BaseSetting";

  private mounted() {
    // 是否隐藏wiki提示
    this.showDeployTip = localStorage.showDistributionDeployTip
      ? JSON.parse(localStorage.showDistributionDeployTip)
      : true;
    this.activeName =
      localStorage.getItem("staff-setting-tab-active") || "BaseSetting";
  }
  private handleTabClick(active: any) {
    localStorage.setItem("staff-setting-tab-active", active.name);
  }

  // 查看配置全员分销
  private checkDeploy() {
    const protocol = "http";
    window.open(
      `${protocol}://xueyuan.vchangyi.com/bin/view/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E7%AE%A1%E7%90%86%E5%90%8E%E5%8F%B0/2.3%E5%85%A8%E5%91%98%E5%88%86%E9%94%80/7%E5%B0%86%E5%85%A8%E5%91%98%E5%88%86%E9%94%80%E9%85%8D%E7%BD%AE%E5%88%B0%E4%BC%81%E5%BE%AE%E4%BE%A7%E8%BE%B9%E6%A0%8F/`
    );
  }
  // 配置全员分销提示关闭
  private closeDeploy() {
    this.showDeployTip = false;
    localStorage.showDistributionDeployTip = this.showDeployTip;
  }
}
</script>

<style lang="scss" scoped>
/deep/.layout-list {
  > .search-content,
  > .operate-content,
  > .slide-box,
  > .table-content,
  > .footer-content {
    display: none;
  }
}
.content-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #fff;
  border-radius: 4px;
  margin-bottom: 30px;
  i {
    cursor: pointer;
    font-weight: bold;
    font-size: 18px;
    color: #c0c4cc;
  }
  span {
    cursor: pointer;
  }
  img {
    width: 14px;
    height: 14px;
    margin-right: 10px;
  }
}
.align-item-center {
  display: flex;
  align-items: center;
}
.c-blue {
  color: #708cf6;
}
</style>
