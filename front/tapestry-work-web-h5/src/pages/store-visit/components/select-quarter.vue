<template>
  <div ref="select" class="select">
    <div class="select-column">
      <van-picker
        :defaultIndex="current[0]"
        :show-toolbar="false"
        valueKey="label"
        itemHeight="48"
        :columns="columns"
        @change="handleChangeYear"
      />
    </div>
    <div class="select-column">
      <van-picker
        :defaultIndex="current[1]"
        :show-toolbar="false"
        valueKey="label"
        itemHeight="48"
        :columns="quarters"
        @change="handleChangeMonth"
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
export default class SelectQuarter extends Vue {
  @Model('update:value', { type: Array, default: () => {} }) value: number[];
  @Prop({ type: Array, default: () => [] }) defaultColumns: any[];
  
  visible: boolean = false;
  innerValue: any[] = [];

  current = [1, 1]

  @Watch('value', { immediate: true })
  onValueChange(val) {
    this.innerValue = val;
    this.current = [
      this.columns.findIndex(item => item.value === val[0]),
      this.quarters.findIndex(item => item.value === val[1])
    ];
  }

  get columns() {
    if(this.defaultColumns.length == 0) return []
    return this.defaultColumns[0].children
  }

  get quarters () {
    if(this.defaultColumns.length == 0) return []
    return this.defaultColumns[1].children.filter(item => {
      return item.year == this.innerValue[0]
    })
  }

  mounted() {
  }
  
  handleChangeYear(e, column, index) {
    this.$set(this.current, 0, index);
    this.$set(this.innerValue, 0, column.value);
    this.$emit('update:value', this.innerValue)
  }

  handleChangeMonth(e, column, index) {
    this.$set(this.current, 1, index);
    this.$set(this.innerValue, 1, column.value);
    this.$emit('update:value', this.innerValue)
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
    &:not(:last-child) {
      border-right: 1px solid var(--Grey-200, rgba(233, 236, 239, 1))
    }
  }

  .select-list {
    padding:  12px 8px;
    position: absolute;
    top: 50%;
    left: 0;
    width: 0;
    width: 100%;
    box-sizing: border-box;
    transition: all 0.3s ease;

    .select-item { 
      height: 48px;
      line-height: 48px;
      text-align: center;
      font-size: 20px;
      color: rgba(0, 0, 0, 1);
      font-weight: 700;
      opacity: 0.2;

      &.active {
        opacity: 1;
      }

      &:not(:last-child) {
        margin-bottom: 10px;
      }
    }
  }
}
</style>