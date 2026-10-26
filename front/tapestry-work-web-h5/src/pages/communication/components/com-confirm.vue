<template>
  <van-popup 
    v-model="visible"
    @click-overlay="handleCancel"
    position="bottom">
    <div class="com-confirm">
      <div class="com-confirm__close" @click="handleCancel"></div>
      <div class="com-confirm__header">下载内容</div>
      <div class="com-confirm__tips">请输入您的邮箱，用于接收下载文件</div>
      <div class="com-confirm__content">
        <div class="input">
          <input ref="input" type="email" v-model="email" placeholder="请输入邮箱" />
          <van-icon class="delete" name="clear" v-if="email" @click="email = ''" />
        </div>
        <div class="select">
          <div class="select-value" @click="handleShow">{{emailSuffix?emailSuffix:'请选择'}}</div>
          <div :class="['select-box', {show: show}]">
            <div class="select-list">
              <div :class="['select-item', {active: item === emailSuffix}]" v-for="(item, index) in ['@coach.com', '@tapestry.com']" :key="index" @click="handleSelect(item)">{{item}}</div>
            </div>
          </div>
        </div>
      </div>
      <div class="com-confirm__footer">
        <van-button class="btn" :disabled="!isEmail" @click="handleConfirm">确认发送</van-button>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from 'vue-property-decorator'
import { Popup, Icon } from "vant";
let that =null
@Component({
  components: {
    [Popup.name]: Popup,
    [Icon.name]: Icon,
  },
  model: {
    prop: 'visible',
    event: 'change'
  }
})
export default class ComConfirm extends Vue {
  @Prop({ type: Boolean, default: false }) visible: boolean;
  /** 邮箱 */
  email = ''
  /** 邮箱后缀 */
  emailSuffix = ''
  /** 显示 */
  show = false

  /** 是否是正确的邮箱 @coach.com、@tapestry.com */
  get isEmail() {
    if (!this.email.trim()) return false;
    return /^[a-zA-Z0-9_.-]+(@vchangyi.com|@coach.com|@tapestry.com)$/.test(this.email+this.emailSuffix);
  }

  @Watch('visible', { immediate: true })
  onVisibleChange(val: boolean) {
    if(val) {
      this.$nextTick(() => {
        (this.$refs.input as HTMLElement).focus();
        // 让页面滚动到最底部
        window.scrollTo(0, document.body.scrollHeight);
        
      })
    }
  }

  mounted(){
    that = this
    const value = localStorage.getItem('email') || '';
    if(value) {
      this.email = value.split('@')[0];
      this.emailSuffix = '@' + (value.split('@')[1]||'');
    }
  }

  handleShow() {
    this.show = !this.show
    
  }

  handleSelect(item: string) {
    this.emailSuffix = item
    this.handleShow()
  }
  handleCancel() {
    this.$emit('change', false);
    this.$emit('cancel');
  }

  handleConfirm = this.$helper.debounce(() => {
    
    (that.$refs.input as HTMLElement).blur();
    const value = that.email+that.emailSuffix
    console.log(111, value);
    that.$emit('confirm', value);
    localStorage.setItem('email', value)
  }, 500)
  
}
</script>

<style lang="scss" scoped> 
::v-deep {
  .van-button--disabled {
    opacity: 1;
    background: #A6A6A6 !important;
  }
}
.com-confirm {
  width: 100vw;
  padding: 20px;
  box-sizing: border-box;
  background: #FFFFFF;
  position: relative;
  padding-bottom: calc(40px + contain(safe-area-inset-bottom));
  padding-bottom: calc(40px + env(safe-area-inset-bottom));


  &__close {
    position: absolute;
    top: 0;
    right: 0;
    width: 48px;
    height: 48px;
    background: url('~@/assets/communication/icon_popup_close.png') center/18px no-repeat;
  }

  &__header {
    height: 22px;
    font-weight: 500;
    font-size: 16px;
    color: #000000;
    line-height: 22px;
  }

  &__tips {
    margin-top: 15px;
    font-weight: 400;
    font-size: 14px;
    color: #000000;
    line-height: 20px;
  }

  &__content {
    margin-top: 15px;
    font-weight: 400;
    font-size: 14px;
    color: #000000;
    line-height: 20px;
    margin-bottom: 50px;
    border: 1px solid #A6A6A6;
    display: flex;
    
    .select {
      width: 130px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      position: relative;

      .select-value {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        &::after{
          content: "";
          display: block;
          width: 0;
          height: 0;
          border-left: 6px solid transparent;
          border-right: 6px solid transparent;
          border-top: 6px solid #000;
          border-bottom: 6px solid transparent;
          margin-top: 6px;
          margin-left: 5px;
        }
      }

      .select-box {
        position: absolute;
        right: 0;
        top: 55px;
        width: 150px;
        height: 100px;
        background: #FFFFFF;
        box-shadow: 0px 1px 5px 0px rgba(0,0,0,0.2);
        z-index: 1;
        display: none;
        &.show {
          display: block;
        }
        &::after {
          content: "";
          position: absolute;
          width: 10px;
          height: 10px;
          background-color: #fff;
          top: -5px;
          left: 50%;
          transform: translateX(-50%) rotate(45deg);
          box-shadow: 0px 1px 5px 0px rgba(0,0,0,0.2);
        }

        .select-list {
          background-color: #fff;
          position: relative;
          z-index: 1;
          height: 100%;
          overflow-y: auto;
          &::-webkit-scrollbar {
            display: none;
          }

        }

        .select-item {
          padding: 15px;
          text-align: center;
          background: #FFFFFF;

          &.active {
            font-weight: 500;
            background-color: rgba($color: #000000, $alpha: 0.2);
          }
        }
      }
    }

    .input {
      flex: 1;
      position: relative;

      &::after {
        content: "";
        width: 1px;
        height: 23px;
        background-color: #979797;
        position: absolute;
        right: 0;
        top: 50%;
        transform: translateY(-50%);
      }
    }
  }

  &__footer {
    display: flex;
    align-items: center;

    .btn {
      width: 100%;
      margin: auto;
      height: 44px;
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
.input {
  width: 100%;
  height: 40px;
  position: relative;

  .delete {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
  }
  input {
    height: 100%;
    width: 100%;
    padding: 0 20px;
    box-sizing: border-box;
    font-size: 14px;

    &::placeholder {
      color: #A6A6A6;
    }
  }
}
.tips {
  font-weight: 400;
  font-size: 10px;
  color: rgba($color: #000000, $alpha: 0.5);
  line-height: 14px;
  margin-top: 10px;
  text-align: left;
}
</style>