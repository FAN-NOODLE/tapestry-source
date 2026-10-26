<template>
  <div 
    v-if="show(data)"
    :class="[
      'select-item',
      'level-' + data.level,
      {
        'has-leaf': !data.isLeaf,
        'expand': data.isExpand
      }
    ]">
    <div class="select-item__content" @click="handleSelect()">
      <div class="select-item__icon" @click.stop="handleExpand()"></div>
      <div class="select-item__name">{{ data[nameKey] }}</div>
      <div 
        :class="[
          'select-item__check',
          {
            'checked': data.isChecked,
            'half': data.isHalf
          }
        ]"
      ></div>
    </div>
    <div v-if="data.isExpand && data.children && data.children.length" class="select-item__collapse">
      <select-item v-for="(item, index) in data.children" :key="index" :data="item" :nameKey="nameKey" :show="show" @expand="handleExpand" @select="handleSelect"></select-item>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'
import { TreeNode } from '@/utils/tree'
@Component({
})
export default class  extends Vue {
  @Prop({ type: Object, default: () => ({}) }) data: TreeNode;
  @Prop({ type: String, default:'name' }) nameKey: string;
  @Prop({ type: Function, default: (item: Object) => true }) show: (item: TreeNode) => boolean;

  handleSelect(value?: TreeNode) {
    this.$emit('select', value||this.data)
  }

  handleExpand(value?: TreeNode) {
    this.$emit('expand', value||this.data)
  }
}
</script>

<style lang="scss" scoped>
.select-item {
  margin-bottom: 15px;
  
  &__content {
    display: flex;
    align-items: center;
  }

  &__icon {
    width: 17px;
    height: 20px;
    visibility: hidden;
    &::after {
      content: '';
      display: block;
      width: 17px;
      height: 20px;
      background: url('~@/assets/display-task/icon_triangle.png') center/7px 4px no-repeat;
      transform: rotate(-90deg);
    }
  }

  &__name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 12px;
    line-height: 20px;
    font-weight: 400;
  }

  &__check {
    width: 16px;
    height: 16px;
    background: url('~@/assets/display-task/icon_check.png') center/100% no-repeat;
    margin-left: 5px;

    &.checked {
      background: url('~@/assets/display-task/icon_check_checked.png') center/100% no-repeat;
    }

    &.half {
      background: url('~@/assets/display-task/icon_check_half.png') center/100% no-repeat;
    }
  }

  &__collapse {
    margin-top: 15px;
  }

  &.has-leaf {
    padding-left: 10px;

    &>.select-item__content .select-item__icon {
      visibility: visible;
    }
    &>.select-item__content .select-item__name {
      font-size: 14px;
    }
  }

  &.expand {
    &>.select-item__content .select-item__icon::after {
      transform: rotate(0deg);
    }
  }

  &.level-0 {
    padding-left: 0;
    &:not(:last-child) {
      border-bottom: 1px solid #F2F2F2;
    }
  }
}
</style>