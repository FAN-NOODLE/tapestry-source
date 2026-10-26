<template>
  <div class="tool">
    <div ref="btn" class="tool-btn" @click="handleClick">
      <slot></slot>
    </div>
    <div v-if="visible" class="tool-popup" @click="visible = false">
      <slot name="content"></slot>
    </div>
    <van-overlay :show="visible" :custom-style="{ backgroundColor: 'rgba(0,0,0,0)' }" @click="visible = false" />
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component({
  components: {
  },
})
export default class ToolPopup extends Vue {
  visible: boolean =false;

  handleCancel() {
    this.$emit('cancel');
  }

  handleClick() {
    this.visible = true
  }
}
</script>

<style lang="scss" scoped>
.icon-delete {
  width: 16px;
  height: 16px;
  background: url("~@/assets/store-visit/icon_delete.png") center/100% no-repeat;
}

.icon-edit {
  width: 16px;
  height: 16px;
  background: url("~@/assets/store-visit/icon_edit.png") center/16px no-repeat;
}

.tool {
  position: relative;
}

.tool-popup {
  padding: 4px 12px;
  border-radius: 8px;
  box-shadow: 4px 4px 16px 0px rgba(0, 0, 0, 0.05);
  background: var(--Grey-100, rgba(248, 249, 250, 1));
  position: absolute;
  z-index: 10;
  right: -10px;
  top: 20px;
  border: 1px solid var(--Grey-200, rgba(233, 236, 239, 1));

  // &::after {
  //   content: "";
  //   position: absolute;
  //   top: -6px;
  //   right: 10px;
  //   width: 10px;
  //   height: 6px;
  //   background: url("~@/assets/store-visit/icon_t.png") center/contain no-repeat;
  // }

  .tool-item {
    display: flex;
    align-items: center;
    font-size: 10px;
    font-weight: 600;
    height: 40px;
    // width: 80px;
    .text {
      margin-left: 4px;
      width: fit-content;
      white-space: nowrap;
    }

    &:not(:last-child) {
      border-bottom: 0.5px solid var(--Grey-200, rgba(233, 236, 239, 1))
    }
  }
}
</style>