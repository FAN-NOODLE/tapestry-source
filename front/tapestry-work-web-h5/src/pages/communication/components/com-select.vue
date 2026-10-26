<template>
  <div class="com-select">
    <div class="com-select__value" @click="handleToggle">
      <div>{{valueName}}</div>
      <div :class="['icon-down', {active: show}]"></div>
    </div>
    <div v-if="show" class="com-select__content">
      <div class="com-select__item" v-for="(item, index) in options" :key="index" @click="handleSelect(item)">
        <div :class="['label', {selected: item.value === value}]">{{ item.label }}</div>
      </div>
    </div>
    <div v-if="show" class="com-select__mark" @click="handleToggle"></div>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

interface IOption {
  value: string | number;
  label: string;
}

@Component({
  name: 'SelectType',
  model: {
    prop: 'value',
    event: 'change',
  },
})
export default class ComSelect extends Vue {
  @Prop({ type: [Number, String], default: null }) value!: string|string;
  @Prop({ type: Array, default: () => [] }) options!: Array<IOption>;
  @Prop({ type: String, default: 'value' }) rowKey: string;

  show = false

  get valueName() {
    return this.options.find((item) => item[this.rowKey] === this.value).label;
  }

  handleToggle() {
    this.show = !this.show
  }

  handleSelect(item: IOption) {
    this.handleToggle()
    this.$emit('change', item[this.rowKey])
  }
}
</script>
<style lang="scss" scoped>
.icon-down {
  width: 8px;
  height: 4px;
  position: relative;
  transform: rotate(0);
  transition: transform 0.3s;
  &.active {
    transform: rotate(180deg);
  }
  &::after {
    content: "";
    // 三角形；
    display: block;
    width: 0;
    height: 0;
    border-left: 4px solid transparent;
    border-right: 4px solid transparent;
    border-top: 4px solid #000;
    border-bottom: 4px solid transparent;
  }
}
.com-select {
  width: 100%;
  position: relative;

  &__mark {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    width: 100vw;
    z-index: 1;
  }

  &__value {
    width: 100%;
    font-size: 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  &__content {
    font-size: 14px;
    position: absolute;
    bottom: -15px;
    right: 0;
    transform: translateY(100%);
    z-index: 2;
    width: 150px;
    background: #FFFFFF;
    box-shadow: 0px 1px 5px 0px rgba(0,0,0,0.2);
    padding: 20px;
    box-sizing: border-box;
    
  }

  &__item {
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 400;
    font-size: 14px;
    color: rgba(74,74,74,0.6);
    line-height: 20px;

    &:not(:last-child) {
      margin-bottom: 20px;
    }

    .label {
      position: relative;
      
      &.selected {
        color: #000;
        font-weight: 500;
        &::after {
          content: '';
          position: absolute;
          left: -10px;
          top: 50%;
          transform: translate(-100%, -50%);
          width: 20px;
          height: 20px;
          background: url('~@/assets/communication/icon_select.png') center/100% no-repeat;
        }
      }
    }
  }
}
</style>