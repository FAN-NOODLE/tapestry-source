<template>
  <div ref="select" class="select">
    <div class="select-column">
      <van-picker
        :defaultIndex="current"
        :show-toolbar="false"
        valueKey="label"
        itemHeight="48"
        :columns="columns"
        @change="handleChange"
      />
    </div>
  </div>
</template>
<script lang="ts">
import dayjs from 'dayjs';
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';


@Component({
  components: {}
})
export default class SelectYear extends Vue {
  @Model('update:value', { type: Number, default: new Date().getFullYear() }) value: number;
  @Prop({ type: Array, default: () => [] }) defaultColumns: any[];
  
  innerValue: number = new Date().getFullYear();

  current = 0
  
  @Watch('value', { immediate: true })
  onValueChange(val) {
    this.innerValue = val;
    this.current = this.columns.findIndex(item => item.value === val)
  }

  get columns() {
    if(this.defaultColumns.length == 0) return []
    return this.defaultColumns[0].children
  }

  mounted() {
    
  }
  /**
   * 变化
   */
  handleChange(e, column, index) {
    this.current = index;
    this.$emit('update:value', column.value)
  }

}
</script>
<style lang="scss" scoped>
.van-picker {
  height: 100%;
  background-color: transparent;
}
::v-deep .van-picker__mask {
  background-image: 
  linear-gradient(180deg, rgba(255, 255, 255, 0), rgba(255, 255, 255, 0)), 
  linear-gradient(0deg, rgba(255, 255, 255, 0), rgba(255, 255, 255, 0));
}

::v-deep .van-picker-column__item {
  font-size: 20px;
  color: rgba(0, 0, 0, 1);
  font-weight: 700;
  opacity: 0.2;
}

::v-deep .van-picker-column__item--selected {
  opacity: 1;
}

.select {
  width: 100%;
  height: 304px;
  display: flex;
  overflow: hidden;

  .select-column {
    flex: 1;
    height: 100%;
    overflow: hidden;
    position: relative;
  }
}
</style>