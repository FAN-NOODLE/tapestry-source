<template>
  <div class="com-filter">
    <div class="com-filter-item" @click="handleSelectDate">
      <div class="value">{{['全部','昨天','今天','近7天','近30天','自定义'][filterParams.date_type]}}</div>
      <div class="icon"></div>
    </div>
    <div class="com-filter-item" @click="handleSelectStore">
      <div class="value">{{filterParams.third_ids.length==0?'全部门店': `已选${filterParams.third_ids.length}家门店`}}</div>
      <div class="icon"></div>
    </div>
    <div class="com-filter-item" @click="handleSelectMember">
      <div class="value">{{memberCount==0?'全部导购': `已选${memberCount}名导购`}}</div>
      <div class="icon"></div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component({

})
export default class ComFilter extends Vue {
  @Prop({ type: Number, default: 0 }) from: number
  get filterParams() {
    return this.$store.getters["onlineOrder/getFilterParams"];
  }

  get memberCount() {
    return this.filterParams.member_ids.length + this.filterParams.leave_member_ids.length
  }
  
  handleSelectDate() {
    this.$router.push({
      path: 'online-order-select-date',
      query: {
        from: String(this.from)
      }
    })
  }

  handleSelectStore() {
    this.$router.push({
      path: 'online-order-select-store',
      query: {
        from: String(this.from)
      }
    })
  }

  handleSelectMember() {
    this.$router.push({
      path: 'online-order-select-member',
      query: {
        from: String(this.from)
      }
    })
  }
}
</script>
<style lang="scss" scoped>
.com-filter {
  height: 44px;
  display: flex;
  align-items: center;
  background-color: #fff;

  .com-filter-item {
    flex: 1;
    height: 100%;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 400;
    font-size: 14px;
    color: #666666;

    &:not(:last-child)::after {
      content: "";
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
      z-index: 1;
      width: 1px;
      height: 12px;
      background-color: #D1BFA7;
    }

    .icon {
      width: 14px;
      height: 14px;
      background: url('~@/assets/online-order/icon_tr.png') center/100% no-repeat;
    }
  }
}

</style>