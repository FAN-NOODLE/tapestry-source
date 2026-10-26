<template>
  <div class="item" @click="handleViewDetail()">
    <div class="item-header">
      <div class="item-title">{{item.display_name}}</div>
      <div v-if="item.status === 1 && beforeEndTime<5*24*60*60*1000" class="item-countdown">
        <van-count-down :time="beforeEndTime" format="距结束DD天HH时mm分" />
      </div>
      <div class="icon-right"></div>
    </div>
    <div class="item-content">
      <div class="item-row">
        <span class="item-status" :data-status="item.status">{{statusLabel}}</span>
        <span class="item-time">执行时间：{{dateFormat(item.start_time)}}-{{dateFormat(item.end_time)}} </span>
      </div>
      <div class="item-desc">{{item.content}}</div>
    </div>
    <div class="item-footer">
      <div 
        class="item-status" 
        :class="{
          active: status === -1 || status === val.value
        }"
        v-for="(val, index) in statusList" 
        :key="index"
        @click.stop="handleViewDetail(val.value)">
        <div class="item-status-num">{{ item[val.key] || 0 }}</div>
        <div class="item-status-label">{{ val.label }}</div>
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
export default class VmDisplayTaskItem extends Vue {
  @Prop({ required: true }) item: IDisplayTaskItem
  @Prop({ required: true }) status: number
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

  get statusList() {
    return auditStatusOptions.filter(item => item.value>0)
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
   * @description 点击查看详情
   */
  handleViewDetail(value: number = -1) {
    this.$store.commit('displayTask/setDetail', {
      task_info: null,
      store_info: null
    });
    this.$router.push({
      path: 'display-task-detail-vm',
      query: {
        id: this.item.id + '',
        finish_status: String(value),
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
    padding-right: 10px;
  }
  &-countdown {
    font-weight: 400;
    .van-count-down {
      line-height: 18px;
      font-size: 12px;
      color: #70A560;
    }

    &.start .van-count-down {
      color: #C6AD7F;
    }
  }

  .icon-right {
    width: 9px;
    height: 14px;
    background: url("~@/assets/display-task/icon_right.png") no-repeat center/100%;
    margin-left: 5px;
  }

  &-content {
    padding: 10px;
    background: rgba(119, 119, 119, 0.1);
    .item-status {
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
  }

  &-row {
    height: 17px;
    font-size: 12px;
    line-height: 17px;
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
    padding: 20px 0 5px;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    
    .item-status {
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      color: #777777;
      flex: 1;
      position: relative;
      &.active {
        color: #000000;
      }
      &:not(:last-child)::after {
        content: "";
        position: absolute;
        top: 50%;
        right: 0;
        transform: translateY(-50%);
        width: 1px;
        height: 30px;
        background-color: rgba(166, 166, 166, 0.5);
      }
     .item-status-num {
        font-size: 18px;
        font-weight: 500;
        line-height: 25px;
      }
     .item-status-label {
        font-size: 14px;
        line-height: 20px;
        margin-top: 2px;
      }
    }
    
  }

  

}
</style>