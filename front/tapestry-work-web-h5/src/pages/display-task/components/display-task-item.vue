<template>
  <div class="item">
    <div class="item-header">
      <div class="item-title">{{item.display_name}}</div>
      <div class="item-audit-status" :data-status="item.finish_status">{{auditStatusLabel}}</div>
    </div>
    <div class="item-content">
      <div class="item-row">
        <span class="item-status" :data-status="item.status">{{statusLabel}}</span>
        <span class="item-time">执行时间：{{dateFormat(item.start_time)}}-{{dateFormat(item.end_time)}} </span>
      </div>
      <div class="item-desc">{{item.content}}</div>
    </div>
    <div v-if="!disabled" class="item-footer">
      <!-- 倒计时 -->
      <div class="flex">
        <div v-if="item.status === 0 && beforeStartTime<5*24*60*60*1000" class="item-countdown start">
          <van-count-down :time="beforeStartTime" format="距开始DD天HH时mm分" />
        </div>
        <div v-if="item.status === 1 && beforeEndTime<5*24*60*60*1000" class="item-countdown">
          <van-count-down :time="beforeEndTime" format="距结束DD天HH时mm分" />
        </div>
        <div v-if="item.finish_status === 3" class="item-score">{{item.grade}}</div>
      </div>
      <div class="btn-group">
        <!-- 进行中/已结束 -->
        <div v-if="item.status !== 0" class="btn" @click="handleViewOtherStore">他店陈列图</div>
        <!-- 进行中 & 待提交 -->
        <div v-if="item.status === 1 && item.finish_status === 0" class="btn primary" @click="handleViewDetail">去上传</div>
        <!-- 进行中/已结束 & 未通过 -->
        <div v-else-if="item.status !== 0 && item.finish_status === 2" class="btn primary" @click="handleViewDetail">重新上传</div>
        <!-- 其他 -->
        <div v-else class="btn" @click="handleViewDetail">查看</div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import { CountDown } from 'vant';
import { statusOptions, auditStatusOptions } from '../const' 
@Component({
  components: {
    [CountDown.name]: CountDown
  }
})
export default class DisplayTaskItem extends Vue {
  @Prop({ required: true }) item: IDisplayTaskItem
  @Prop({ default: false, type: Boolean }) disabled: boolean
  /** 活动状态 */
  get statusLabel() {
    return statusOptions.find((item) => item.value === this.item.status).label
  }
  /** 审核状态 */
  get auditStatusLabel() {
    return auditStatusOptions.find((item) => item.value === this.item.finish_status).label
  }
  /** 剩余多少时间开始 */
  get beforeStartTime() {
    return this.item.start_time * 1000 - Date.now()
  }
  /** 剩余多少时间结束 */
  get beforeEndTime() {
    return this.item.end_time * 1000 - Date.now()
  }
  /**
   * @description 格式化时间
   * @param time 
   * @param format 
   */
  dateFormat(time: number, format = "yyyy/MM/dd") {
    return this.$helper.dateFormat(time, format);
  }

  /**
   * @description 点击查看其他店铺
   */
  handleViewOtherStore() {
    this.$emit('view')
  }

  /**
   * @description 点击查看详情
   */
  handleViewDetail() {
    this.$store.commit('displayTask/setDetail', {
      task_info: null,
      store_info: null
    });
    this.$store.commit('displayTask/setUpdate', true);
    this.$router.push({
      path: 'display-task-detail',
      query: {
        id: this.item.id+'',
        store_id: this.item.store_id + '',
        from: '1'
      }
    });
  }
}
</script>
<style scoped lang="scss">
.item {
  width: 345px;
  background: #FFFFFF;
  box-shadow: 0px 3px 17px 0px rgba(54,68,58,0.06);
  padding: 15px;
  margin: auto;
  line-height: normal;

  &:not(:last-child) {
    margin-bottom: 15px;
  }

  &-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;

  }
  &-title {
    font-size: 14px;
    color: #000;
    flex: 1;
    font-weight: bold;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &-audit-status {
    font-weight: 400;
    font-size: 10px;
    width: 60px;
    height: 15px;
    border-radius: 15px;
    display: flex;
    align-items: center;
    justify-content: center;
    // margin-right: 5px;
    &[data-status="0"] {
      color: #000;
      background: rgba(0,0,0,0.1);
    }
    &[data-status="1"], &[data-status="4"] {
      color: #000;
      background: #D1BFA7;
    }
    &[data-status="2"] {
      color: #E02020;
      background: #FFD0D0;
    }
    &[data-status="3"] {
      color: #000;
      background: #BBD1A7;
    }
  }

  &-content {
    padding: 10px;
    background: rgba(119, 119, 119, 0.1);
  }

  &-row {
    height: 17px;
    font-size: 12px;
    line-height: 17px;
  }

  &-status {
    &[data-status="0"] {
      color: #C6AD7F;
    }
    &[data-status="1"] {
      color: #70A560;
    }
    &[data-status="2"] {
      color: #A6A6A6;
    }
  }

  &-time {
    font-weight: 400;
    color: #000000;
    margin-left: 10px;
  }

  &-desc {
    font-size: 12px;
    font-weight: 400;
    color: #7F7F7F;
    line-height: 17px;
    margin-top: 10px;
    // 超出两行省略
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }


  &-footer {
    margin-top: 10px;
    .flex {
      justify-content: flex-end;
      align-items: center;
      margin-bottom: 10px;
    }
  }

  &-countdown {
    font-weight: 400;
    flex: 1;
    .van-count-down {
      line-height: 18px;
      font-size: 12px;
      color: #70A560;
    }

    &.start .van-count-down {
      color: #C6AD7F;
    }
  }

  &-score {
    color: #70A560;
    font-size: 22px;
    font-weight: bold;
    height: 30px;
    line-height: 30px;
    &::before {
      content: "评分:";
      font-weight: 400;
      font-size: 12px;
      margin-right: 2px;
    }
  }

  .btn-group {

    display: flex;
    align-items: center;

    .btn {
      font-size: 14px;
      flex: 1;
      height: 36px;
      background: #FFFFFF;
      border: 1px solid #000000;
      box-sizing: border-box;
      display: flex;
      align-items: center;
      justify-content: center;

      &:nth-child(2) {
        margin-left: 15px;
      }

      &.primary {
        background-color: #000000;
        color: #fff;
      }
    }
  }
}
</style>