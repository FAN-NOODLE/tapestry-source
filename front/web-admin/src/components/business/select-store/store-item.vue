<template>
  <div class="store-item">
    <div class="store-item__icon">
      <img
        class="group-img"
        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
        alt=""
      />
    </div>
    <div class="store-item__name">{{ item.name||item.dp_store_name||item.store_name }}</div>
    <div v-if="clearable" class="store-item__delete" @click="handleDelete">
      <i class="el-icon-circle-close"></i>
    </div>
  </div>
</template>
<script lang="ts">

import { Component, Vue, Prop } from "vue-property-decorator";
@Component({
  name: 'StoreItem'
})
export default class StoreItem extends Vue {
  @Prop({default: ()=>({}), type: Object }) item!: any;
  @Prop({default: false, type: Boolean }) clearable!: boolean;

  handleDelete() {
    this.$emit('delete');
  }
}
</script>
<style scoped lang="scss">
  .store-item {
    display: flex;
    align-items: center;
    height: 30px;

    &:not(:last-child) {
      border-bottom: 1px solid #f7f7f7;
    }

    
    &__icon {
      width: 14px;
      height: 14px;
      .group-img {
        width: 100%;
        height: 100%;
        object-fit: scale-down;
        display: block;
      }
    }

    &__name {
      flex: 1;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      padding: 0 10px;
    }

    &__delete {
      width: 30px;
      height: 30px;
      display: flex;
      justify-content: center;
      align-items: center;
      cursor: pointer;
      font-size: 16px;
    }
  }
</style>
