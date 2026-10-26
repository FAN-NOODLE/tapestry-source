<template>
  <div class="communication-item">
    <div class="communication-item__content">
      <div class="title">{{ item.title }}</div>
      <div class="date">
        <span>{{ formatTime(item.update_at) }}</span>
        <span class="name">{{ item.submit_name }}</span>
      </div>
      <div class="text">
        <span>{{ item.department_name }}</span>
        <span> | </span>
        <span :class="{
          type: item.type === 3
        }">{{ typeName }}</span>
      </div>
    </div>
    <div v-if="type === 'ro'" class="communication-item__footer">
      <div :class="['btn plain', {disabled: item.status === 5}]" @click="handleInvalid">作废</div>
      <div v-if="item.status === 1" class="btn" @click="handleView">去审核</div>
      <div v-else class="btn plain" @click="handleView">查看</div>
    </div>
    <div v-if="type === 'ft'" class="communication-item__footer">
      <div v-if="item.status === 1" class="btn plain" @click="handleRecall">撤回</div>
      <div class="btn plain" @click="handleView">查看</div>
    </div>
    <div class="status" :data-status="item.status">{{ statusName }}</div>
    <div v-if="!item.is_read" class="is-new"></div>
    <com-message 
      v-model="showRecallPopup" 
      title="撤回提醒" 
      content="确定撤回当前内容？撤回后需重新提交RO审核。"
      @cancel="showRecallPopup = false"
      @confirm="handleRecallConfirm"></com-message>
    <com-message
      v-model="showInvalidPopup" 
      title="作废提醒" 
      content="确定作废当前内容，作废后该内容将不可修改使用！"
      @cancel="showInvalidPopup = false"
      @confirm="handleInvalidConfirm"></com-message>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator';
import { Image } from "vant";
import { typeOptions, statusOptions } from '../const'
import ComMessage from './com-message.vue'
import { cancelApi, recallApi } from '@/api/communication';

export interface ICommunicationItem {
  /** 传讯ID */
  id?: number;
  /** 传讯标题 */
  title: string;
  /** 提交时间 */
  submit_time: number;
  /** 传讯状态 */
  status: number;
  /** 发布人 */
  submit_name: string;
  /** 最近更新时间 */
  update_at: number;
  /** 所属部门 */
  department_name: string;
  /** 传讯类型 1-周传讯；2-转货传讯；3-紧急传讯 */
  type: number;
  /** 是否新的 */
  is_read: boolean;
}

@Component({
  components: {
    [Image.name]: Image,
    ComMessage
  }
})
export default class CommunicationItem extends Vue {
  /** 数据 */
  @Prop({ required: true }) item!: ICommunicationItem;
  @Prop({ type: String, default: 'ro' }) type!: string;

  showRecallPopup = false
  showInvalidPopup = false

  get typeName() {
    return typeOptions.find(item => item.value === this.item.type).label
  }

  get statusName() {
    return statusOptions.find(item => item.value === this.item.status).label
  }
  activated() {
    const id = Number(localStorage.getItem('communicationId'))
    if(id && id === this.item.id) {
      this.$emit('update', {...this.item, is_read: true})
      localStorage.removeItem('communicationId')
    }
  }

  formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy/MM/dd hh:mm");
  }

  /** 撤回 */
  handleRecall() {
    this.showRecallPopup = true
  }
  /** 撤回确认 */
  async handleRecallConfirm() {
    
    try {
      const {code, message, data}:IApiResult = await recallApi({
        id: this.item.id
      })
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.showRecallPopup = false
      this.item.status = 2
      // 神策数据上报 -- 撤回传讯
      this.$track.trackWithdrawTransmittal({
        transmittal_title: this.item.title,
        transmittal_publisher: this.item.submit_name,
        release_time: this.item.update_at,
      })
    } catch (error) {
      console.log(error);
    }
  }

  /** 作废 */
  handleInvalid() {
    if(this.item.status === 5) return
    this.showInvalidPopup = true
  }

  /** 作废确认 */
  async handleInvalidConfirm() {
    try {
      const {code, message, data}:IApiResult = await cancelApi({
        id: this.item.id
      })
      if (code!== 200) {
        this.$toast.fail(message);
        return;
      }
      this.showInvalidPopup = false
      this.item.status = 5
      // 神策数据上报 -- 作废传讯
      this.$track.trackVoidTransmittal({
        transmittal_title: this.item.title,
        transmittal_publisher: this.item.submit_name,
        release_time: this.item.update_at,
      })
    } catch (error) {
      console.log(error);
    }
  }
  
  /** 查看 */
  handleView() {
    localStorage.setItem('communicationId', String(this.item.id))
    this.$router.push({
      path: 'communication-detail/'+this.type,
      query: {
        id: String(this.item.id),
        from: '1'
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.communication-item {
  width: 345px;
  padding: 15px;
  position: relative;
  background-color: #fff;
  box-shadow: 0px 1px 3px 0px rgba(0,0,0,0.1);
  margin-bottom: 15px;

  &__content {
    .title {
      width: 290px;
      font-weight: 600;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      text-align: left;
      font-style: normal;
    }

    .date {
      height: 14px;
      font-weight: 300;
      font-size: 10px;
      color: #A6A6A6;
      line-height: 14px;
      text-align: left;
      font-style: normal;
      margin-top: 10px;
    }

    .name {
      font-size: 11px;
      font-weight: 500;
      padding-left: 5px;
    }

    .text {
      font-weight: 400;
      font-size: 10px;
      color: #000000;
      line-height: 14px;
      text-align: left;
      font-style: normal;
      margin-top: 5px;
    }
    
    .type {
      color: #FF8F52;
    }
  }
  &__footer {
    margin-top: 10px;
    display: flex;
    justify-content: space-between;
    .btn {
      flex: 1;
      height: 30px;
      box-shadow: 0px 1px 3px 0px rgba(0,0,0,0.1);
      border: 1px solid #000000;
      background: #000000;
      color: #fff;
      font-size: 12px;
      display: flex;
      align-items: center;
      justify-content: center;

      &.plain {
        background: #fff;
        color: #000;
      }

      &.disabled {
        border: 1px solid #A6A6A6;
        color: #A6A6A6;
      }

      &:nth-child(2) {
        margin-left: 10px;
      }
    }
  }

  .status {
    width: 40px;
    height: 21px;
    box-shadow: 0px 1px 3px 0px rgba(0,0,0,0.1);
    border-radius: 100px 0px 0px 100px;
    position: absolute;
    top: 10px;
    right: -7px;
    font-size: 10px;
    color: #000;
    display: flex;
    align-items: center;
    justify-content: center;

    &::after {
      content: '';
      width: 7px;
      height: 6px;
      position: absolute;
      bottom: -6px;
      right: 0;
    }

    &[data-status="0"] {
      background: #A6A6A6;
      &::after {
        background: linear-gradient(to right bottom, #787777, #787777 50%, transparent 51%, transparent);
      }
    }
    &[data-status="1"] {
      background: #C7AD7F;
      &::after {
        background: linear-gradient(to right bottom, #8E713E, #8E713E 50%, transparent 51%, transparent);
      }
    }
    &[data-status="3"] {
      background: #E6B8B8;

      &::after {
        background: linear-gradient(to right bottom, #B27575, #B27575 50%, transparent 51%, transparent);
      }
    }
    &[data-status="4"] {
      background: #BBD1A7;

      &::after {
        background: linear-gradient(to right bottom, #7A9F59, #7A9F59 50%, transparent 51%, transparent);
      }
    }
    &[data-status="2"] {
      background: #E1E1E1;

      &::after {
        background: linear-gradient(to right bottom, #C0C0C0, #C0C0C0 50%, transparent 51%, transparent);
      }
    }
    &[data-status="5"] {
      background: #E1E1E1;
      opacity: 0.5;

      &::after {
        background: linear-gradient(to right bottom, #C0C0C0, #C0C0C0 50%, transparent 51%, transparent);
      }
    }
  }

  .is-new {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 33px;
    height: 33px;
    overflow: hidden;

    &::before {
      content: "";
      position: absolute;
      bottom: 0;
      right: 0;
      width: 100%;
      height: 100%;
      transform: translate(50%, 50%)rotate(-45deg);
      background-color: #000;
    }

    &::after {
      content: "new";
      position: absolute;
      bottom: 0;
      right: 0;
      width: 100%;
      height: 100%;
      transform: translate(50%, 50%)rotate(-45deg);
      color: #fff;
      font-size: 8px;
      line-height: 11px;
      z-index: 1;
      display: flex;
      justify-content: center;
      align-items: flex-start;
    }
  }
}
</style>