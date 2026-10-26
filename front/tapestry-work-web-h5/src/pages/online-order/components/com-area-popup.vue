
<template>
  <div>
    <div @click="handleShowPopup">
      <slot></slot>
    </div>
    <van-popup v-model:show="show" position="bottom">
      <van-area v-model="innerValue" title="选择省市区" :area-list="areaList" @cancel="handleCancel" @confirm="handleAreaConfirm"></van-area>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { areaList } from '@vant/area-data';
import { Component, Model, Vue, Watch } from 'vue-property-decorator';

@Component({

})
export default class ComAreaPopup extends Vue {
  @Model('change', { type: String, default: '' }) value: string;
  
  show: boolean = false;
  areaList: any = areaList;
  innerValue: string = '';

  @Watch('value')
  onValueChange(val: string){
    this.innerValue = val;
  }

  handleShowPopup() {
    this.show = true
  }

  handleAreaConfirm(values:any[]) {
    console.log(values);
    this.$emit('change', values[2].code, values);
    this.show = false
  }

  handleCancel() {
    this.$emit('cancel');
    this.show = false
  }
}
</script>
<style lang="scss" scoped>
.container {
  width: 100vw;
  min-height: 100vh;
  background: #F7F7F7;
}

</style>