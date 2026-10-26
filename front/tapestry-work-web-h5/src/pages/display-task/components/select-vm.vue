<template>
  <div>
    <div class="select-vm__value" @click="handleShowPopup">
      <span v-if="value" class="value">{{name}}</span>
      <span v-else class="placeholder">请选择所属VM（必填）</span>
      <div class="icon"></div>
    </div>
    <van-popup
      class="popup"
      v-model:show="show"
      position="bottom"
    >
      <div ref="list" class="popup-content" >
        <div class="list">
          <div 
            class="vm-item" 
            ref="item"
            :class="{
              active: index === current,
              prev: index === current - 1,
              next: index === current + 1
            }"
            v-for="(item, index) in list" 
            :key="index"
            @click="handleSelect(index)">
            {{ item.member_name }}
          </div>
        </div>
      </div>
      <div class="popup-footer">
        <div class="btn plain" @click="handleSelectCancel">取消</div>
        <div class="btn" @click="handleSelectConfirm">确定</div>
      </div>
    </van-popup>

  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from "vue-property-decorator";
import { Popup } from 'vant';

@Component({
  components: {
    [Popup.name]: Popup
  }
})
export default class SelectVm extends Vue {
  /** 选中的VM */
  @Model('change', { type: [String, Number], default: '' }) value: string|number;
  /** VM列表 */
  @Prop({ type: Array, default: () => []}) list: IVmItem[];
  @Prop({ type: Boolean, default: false }) disabled: boolean;
  /** 是否显示弹窗 */
  show = false;
  /** 当前选中的VM下标 */
  current = -1;
  /** 选中的VM名称 */
  get name() {
    return this.list.find(item => Number(item.member_id) === Number(this.value))?.member_name
  };

  /**
   * @description 获取选中下标
   */
  getCurrent() {
    this.current = this.list.findIndex(item => Number(item.member_id) === Number(this.value));
  }

  /**
   * @description 显示弹窗
   */
  handleShowPopup() {
    if (this.disabled) return;
    this.show = true;
    this.getCurrent();
  }

  /**
   * @description 隐藏弹窗
   */
  handleHidePopup() {
    this.show = false;
  }

  /**
   * @description 选择VM
   * @param index 
   */
  handleSelect(index: number) {
    this.current = index;
  }

  /**
   * @description 确认选择
   */
  handleSelectConfirm() {
    this.$emit('change', this.list[this.current].member_id)
    this.handleHidePopup();
  }

  /**
   * @description 取消选择
   */
  handleSelectCancel() {
    this.handleHidePopup();
  }
}
</script>
<style lang="scss" scoped>
.select-vm__value {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 38px;
  background: #F6F6F6;
  font-size: 12px;
  padding: 0 15px;
  margin-top: 10px;

  .placeholder {
    color: #A1A1A1;
    flex: 1;
  }
  .value {
    flex: 1;
  }

  .icon {
    width: 7px;
    height: 4px;
    background: url('~@/assets/display-task/icon_triangle.png') center/100% no-repeat;
  }
}

.popup {
  &-content {
    max-height: 30vh;
    padding: 15px;
    overflow-y: scroll;
    transition: all 0.3s;

    .vm-item {
      padding: 11px;
      font-size: 14px;
      font-weight: 400;
      color: #A6A6A6;
      line-height: 20px;
      text-align: center;
      transition: all 0.1s;

      &.active {
        font-size: 16px;
        font-weight: 400;
        color: #000000;
        line-height: 22px;
        background-color: #F0F0F0;
        &.prev {
          padding-bottom: 8px;
        }
        &.next {
          padding-top: 8px;
        }
      }
    }
  }

  &-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 15px;
    padding-bottom: calc(15px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(15px + env(safe-area-inset-bottom)); /** iOS 11.0+ */


    .btn {
      width: 165px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-sizing: border-box;
      background-color: #000;
      color: #fff;
      font-size: 14px;

      &.plain {
        background-color: #fff;
        color: #000000;
        border: 1px solid #000000;
        margin-right: 15px;
      }
    }
  }
}
</style>