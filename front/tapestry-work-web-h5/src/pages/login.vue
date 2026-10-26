<template>
  <div class="content">
    <div class="form-item">
      <van-field v-model="userid" placeholder="请输入userID" />
    </div>
    <div class="form-item">
      <van-field v-model="enterprise_id" placeholder="请输入企业ID" />
    </div>
    <van-button class="add-btn" type="primary" size="large" @click="getToken">
      获取TOKEN</van-button
    >
    <div class="acount">
      <span v-for="(item, index) in ['13992092140', 'ChengBin', '18792889956', 'youyou', 'Yoyoyou', '18092066590', 'yanjiaxin']" :key="index" @click="userid=item">{{ item }}</span>
    </div>
    <!-- <div class="red" v-if="isToken">获取成功，已存储token</div> -->
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Button, Toast, Field } from "vant";
import { mockApi } from "@/api/common/index";
import { sensorsInit, sensorsLogin, setSensorsProfile } from "@/utils/sensors";
@Component({
  components: {
    [Button.name]: Button,
    [Field.name]: Field,
  },
})
export default class extends Vue {
  // 18792889956 - 3363603667487424
  // ChengBin - 3363603667487424
  // 13992092140 - 3363603667487424
  private userid = "18792889956";
  private enterprise_id = "3363603667487424";//uat coach信息
  private isToken = false;
  private async getToken() {
    this.$toast.loading({ message: "加载中...", forbidClick: true });
    let params = {
      userid: this.userid,
      enterprise_id: this.enterprise_id,
    };
    const data = await mockApi(params);
    if (data.data.token) {
      this.$store.commit("user/setToken", data.data.token);

      // // 神策数据上报 -- 初始化
      // const sensors = sensorsInit()
      // // 神策数据上报 -- 初始化并挂载全局
      // Vue.prototype.sensors = sensors
      // 神策数据上报 -- 模拟登陆
      sensorsLogin(data.data.userid)
      // 神策数据上报 -- 设置用户属性
      setSensorsProfile({
        store_name: data.data.store_name,
        store_code: data.data.store_code || data.data.store_name
      })
      this.$toast.success("模拟登陆成功");
      const pageUrl = sessionStorage.getItem("pageUrl");
      if (pageUrl) {
        sessionStorage.removeItem("pageUrl");
        window.location.href = pageUrl;
      }

    } else {
      this.$toast.fail("失败");
    }
    setTimeout(() => {
      this.$toast.clear();
    }, 1500);
    this.isToken = true;
  }
  private async mounted() {
    Toast.clear();
    // this.getToken()
  }
}
</script>
<style lang="scss" scoped>
.content {
  padding: 20px;
  text-align: center;
}
.form-item {
  margin-bottom: 20px;
}
.red {
  font-size: 16px;
  margin-top: 60px;
  color: orange;
}
.acount {
  font-size: 16px;
  display: flex;
  margin-top: 20px;
  flex-wrap: wrap;

  span {
    margin-right: 20px;
    margin-bottom: 20px;
  }
}
</style>
