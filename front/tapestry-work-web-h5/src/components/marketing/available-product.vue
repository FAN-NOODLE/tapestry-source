<template>
  <van-action-sheet v-model="show" title="部分可用商品" @close="close">
    <div class="outline">
      <van-list
        v-model="isLoading"
        :finished="isFinished"
        finished-text="没有更多了"
        :immediate-check="false"
      >
        <div class="content" v-for="item in couponGoodlist" :key="item.spu_id">
          <van-image :src="item.image" />
          <div class="tit moreLlipsis">{{ item.spu_name }}</div>
        </div>
      </van-list>
    </div>
  </van-action-sheet>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { ActionSheet, Image, List } from "vant";
import { wxCouponGoodlistApi } from "@/api/marketing-promotion/index";
@Component({
  components: {
    vanActionSheet: ActionSheet,
    [Image.name]: Image,
    vanList: List,
  },
})
export default class extends Vue {
  @Prop(Number) couponId!: number;
  private show = true;
  private couponGoodlist = [];
  private isLoading = false;
  private isFinished = false;
  private mounted() {
    this.goodsList(this.couponId);
  }

  private close() {
    this.$emit("close");
  }

  private async goodsList(coupon_id) {
    this.isLoading = true;
    const { code, message, data }: IApiResult = await wxCouponGoodlistApi(
      coupon_id
    );
    this.isLoading = false;
    this.isFinished = true;
    if (code !== 200) {
      console.log(message);
      return;
    }
    this.couponGoodlist = data.list;
  }
}
</script>
<style lang="less" scoped>
.van-action-sheet__header {
  border-bottom: 1px solid #ebebeb;
}
// 两行显示
.moreLlipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}
.outline {
  padding: 15px;
  background: white;
  height: 50vh;
}
.content {
  width: 345px;
  height: 100px;
  background: #f4f5f8;
  border-radius: 2px;
  padding: 10px;
  display: flex;
  margin-bottom: 15px;
  .van-image {
    width: 76px;
    height: 80px;
    margin-right: 13px;
  }
  .tit {
    font-size: 12px;
    color: #333333;
  }
}
</style>
