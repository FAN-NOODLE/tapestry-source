<template>
  <empty-page
    picture="noData"
    paddingTop="60"
    v-if="steps.length == 0 && showEmpty"
    paddingBottom="60"
    text="暂无足迹"
  ></empty-page>
  <div class="container" v-else>
    <van-steps direction="vertical" :active="0" active-color="#000">
      <van-step v-for="(step, index) in steps" :key="index">
        <h3 class="step-content">
          <span class="type">{{ step.event_txt }}</span
          ><span class="desc">{{ step.event_desc }}</span>
        </h3>
        <p class="step-time">{{ formatTime(step.created_at) }}</p>
      </van-step>
    </van-steps>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Step, Steps } from "vant";
import { customerFootprintListApi } from "@/api/customer-portrait";
import emptyPage from "@/components/base/empty-page.vue";
@Component({
  components: {
    emptyPage,
    [Step.name]: Step,
    [Steps.name]: Steps,
  },
})
export default class extends Vue {
  private steps = [];
  private showEmpty = false;
  private mounted() {
    if (this.$store.getters["user/getCustomerId"]) {
      this.getfootList();
    }
  }
  private formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy-MM-dd hh:mm:ss");
  }
  async getfootList(): Promise<void> {
    const { code, data, message }: IApiResult = await customerFootprintListApi({
      user_id: this.$store.getters["user/getCustomerId"],
      page: 1,
      size: 99999,
    });
    if (code == 200) {
      data.list.length > 0 && (this.steps = data.list);
      this.showEmpty = true;
    } else {
      this.$toast.fail(message);
    }
  }
}
</script>
<style lang="scss" scoped>
div.container {
  background-color: #ffffff;
  height: calc(100vh - 46px);
  overflow: auto;
  .step-content {
    font-size: 12px;
    color: #4a4a4a;
    line-height: 17px;
    .type {
      font-size: 12px;
      display: inline-block;
      margin-right: 5px;
      line-height: 17px;
    }
    .desc {
      font-size: 12px;
      display: inline-block;
      line-height: 17px;
    }
  }
  .step-time {
    display: inline-block;
    color: #96979e;
    line-height: 17px;
    font-size: 12px;
    padding-top: 10px;
  }
}
</style>
