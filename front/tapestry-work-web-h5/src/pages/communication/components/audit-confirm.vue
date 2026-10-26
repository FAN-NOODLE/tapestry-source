<template>
  <van-popup
    v-model="visible"
    position="bottom"
    @click-overlay="handleCancel">
    <div class="audit-confirm">
      <div class="audit-confirm__header">{{ ['审核不通过', '审核通过'][type] }}</div>
      <div v-if="type === 0" class="audit-confirm__content">
        <div class="tips" style="text-align: left;">确认审核不通过？</div>
        <van-field
          v-model="content"
          rows="5"
          autosize
          label=""
          type="textarea"
          maxlength="500"
          placeholder="请输入您的意见："
          show-word-limit
        />
      </div>
      <div v-else class="audit-confirm__content">
        <div class="tips">确认审核通过？</div>
      </div>
      <div class="audit-confirm__footer">
        <van-button class="btn plain" @click="handleCancel">取消</van-button>
        <van-button class="btn" @click="handleConfirm">确定</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'
import { Popup, Field } from "vant";

@Component({
  components: {
    [Popup.name]: Popup,
    [Field.name]: Field
  },
  model: {
    prop: 'visible',
    event: 'change'
  }
})
export default class AuditConfirm extends Vue {
  @Prop({ type: Boolean, default: false }) visible: boolean;
  @Prop({ type: Number, default: 0 }) type: number;
  
  content: string = '';

  get disabled() {
    return this.type === 0 && !this.content.trim();
  }

  handleCancel() {
    this.$emit('change', false);
  }

  handleConfirm() {
    this.$emit('confirm', {
      type: this.type,
      value: this.content.trim()
    });
  }
}
</script>

<style lang="scss" scoped> 
::v-deep {
  .van-field__word-limit {
    color: #A6A6A6
  }

  .van-button--disabled {
    opacity: 1;
    background: #A6A6A6 !important;
  }
}
.audit-confirm {
  width: 100vw;
  padding: 15px;
  box-sizing: border-box;
  background: #FFFFFF;
  position: relative;
  padding-bottom: calc(20px + contain(safe-area-inset-bottom));
  padding-bottom: calc(20px + env(safe-area-inset-bottom));


  &__header {
    height: 22px;
    font-weight: 500;
    font-size: 16px;
    color: #000000;
    line-height: 22px;
    text-align: center;
  }

  &__content {
    text-align: center;
    color: #000000;
    margin-top: 10px;
    
    .tips {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
    }

    .van-field {
      margin-top: 10px;
      background: rgba(166, 166, 166, 0.1);
    }
  }

  &__footer {
    display: flex;
    align-items: center;
    margin-top: 15px;

    .btn {
      flex: 1;
      height: 40px;
      background: #000000;
      color: #fff;
      font-weight: 400;
      font-size: 14px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-sizing: border-box;

      &.plain {
        background: #fff;
        color: #000;
        border: 1px solid #000;
      }

      &:nth-child(2) {
        margin-left: 10px;
      }
    }
  }

}
</style>