<template>
  <van-popup v-model="visible" position="bottom" round  @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-header">
        <div class="range">
          <div class="range-item" :class="{ 'active': current == 0 }" @click="current=0">
            <div class="label">{{ $t('storeVisit.popup.title.21') }}</div>
            <div class="value">{{showDate(0)}}</div>
          </div>
          <div class="range-item" :class="{ 'active': current == 1 }" @click="current=1">
            <div class="label">{{ $t('storeVisit.popup.title.22') }}</div>
            <div class="value">{{showDate(1)}}</div>
          </div>
        </div>
        <div class="btn" @click="handleConfirm">{{ $t('storeVisit.popup.btn.20') }}</div>
      </div>
      <div class="calendar">
        <div class="calendar-header">
          <div class="month">
            <div class="icon-left-circle" @click="handleMonthChange(-1)"></div>
            <div class="month-text">{{showMonth}}</div>
            <div class="icon-right-circle" @click="handleMonthChange(1)"></div>
          </div>
          <div class="row">
            <div class="col" v-for="(item, index) in week" :key="index">{{ item }}</div>
          </div>
        </div>
        <div class="calendar-content">
          <div class="row">
            <div class="col" v-for="(item, index) in days" :key="index">
              <div class="day" :class="{
                'disabled': disabled(item),
                'current': item.unixDate === currentDate,
                'today': item.isToday
              }" @click="handleClick(item)">
                <div class="day-text">{{ item.day }}</div>
                <div v-if="schedules.includes(item.unixDate)" class="day-dot"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import SelectDatePopup from './select-date-popup.vue'
import dayjs from 'dayjs';

const months = {
  1: {
    full: "January",
    abbreviation: "Jan."
  },
  2: {
    full: "February",
    abbreviation: "Feb."
  },
  3: {
    full: "March",
    abbreviation: "Mar."
  },
  4: {
    full: "April",
    abbreviation: "Apr."
  },
  5: {
    full: "May",
    abbreviation: "May"
  },
  6: {
    full: "June",
    abbreviation: "June"
  },
  7: {
    full: "July",
    abbreviation: "July"
  },
  8: {
    full: "August",
    abbreviation: "Aug."
  },
  9: {
    full: "September",
    abbreviation: "Sept."
  },
  10: {
    full: "October",
    abbreviation: "Oct."
  },
  11: {
    full: "November",
    abbreviation: "Nov."
  },
  12: {
    full: "December",
    abbreviation: "Dec."
  }
};


@Component({
  components: {
    SelectDatePopup
  }
})
export default class CalendarPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: Array, default: () => [] }) schedules: string[]
  @Prop({ type: Array, default: () => ['', ''] }) defaultValue: string[]
  @Prop({ type: Number, default: 0 }) minDate: number
  @Prop({ type: Number, default: 0 }) maxDate: number
  visible = false
  week = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];
  days = []
  // 今天
  today = this.getCurrentDate()
  // 当前年份
  currentYear = new Date().getFullYear();
  // 当前月份
  currentMonth = new Date().getMonth() + 1;

  innerValue = ''
  innerRangeValue = ['', '']
  current = 0

  get showMonth() {
    return [this.currentYear, this.currentMonth].join('.')
  }
  get currentDate() {
    return this.innerRangeValue[this.current] || ''
  }
  get showDate() {
    return (index) => this.innerRangeValue[index] ? this.innerRangeValue[index].replaceAll('-', '.') : '--.--.--'
  }

  get disabled() {
    return (item) => {
      if(this.current === 1) {
        //
        if(this.innerRangeValue[0]) {
          return item.unix < dayjs(this.innerRangeValue[0]).unix()
        }
        return this.maxDate && item.unix < dayjs(this.maxDate).unix()
      } else {
        if(this.innerRangeValue[1]) {
          return item.unix > dayjs(this.innerRangeValue[1]).unix()
        }

        return this.minDate && item.unix < dayjs(this.minDate).unix()
      }
      
    }
  }

  @Watch('show', {
    immediate: true,
    deep: true
  })
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.current = 0;
      this.innerRangeValue = [...this.defaultValue]
    }
  }

  mounted() {
    this.init();
  }
  init() {
    const date = new Date();
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    this.innerRangeValue = this.defaultValue || ['', '']
    this.currentYear = year;
    this.currentMonth = month;
    this.days = this.getDaysInMonth(year, month);
  }
  /**
   * 获取当前月份的天数
   */
  getDaysInMonth(year: number, month: number) {
    const date = new Date(year, month, 0);
    const days = date.getDate();
    const firstDay = new Date(year, month - 1, 1).getDay();
    const arr = [];
    for (let i = 0; i < firstDay; i++) {
      arr.push({
        day: null,
        unix: 0,
        unixDate: null,
      });
    }
    for (let i = 1; i <= days; i++) {
      const unixDate = dayjs(new Date(year, month - 1, i)).format('YYYY-MM-DD')
      arr.push({
        day: i,
        unix: dayjs(new Date(year, month - 1, i)).unix(),
        unixDate: unixDate,
        isToday: unixDate === this.today
      });
    }
    return arr;
  }
  /**
   * 获取当前日期
   * @param d 
   */
  getCurrentDate() {
    const date = new Date();
    return dayjs(date).format('YYYY-MM-DD');
  }
  
  /**
   * 点击日期
   * @param item 
   */
  handleClick(item) {
    if(this.disabled(item)) {
      return;
    }
    this.innerRangeValue[this.current] = item.unixDate;
    this.$set(this.innerRangeValue, this.current, item.unixDate)
  }
  /**
   * 切换年月月份
   */
  handleMonthChange(value: number) {
    if(value < 0) {
      const month = this.currentMonth + value;
      if(month < 1) {
        this.currentMonth = 12;
        this.currentYear = this.currentYear - 1;
      } else {
        this.currentMonth = month;
      }
    } else {
      const month = this.currentMonth + value;
      if(month > 12) {
        this.currentMonth = 1;
        this.currentYear = this.currentYear + 1;
      } else {
        this.currentMonth = month;
      }

    }
    this.days = this.getDaysInMonth(this.currentYear, this.currentMonth);
  }

  handleConfirm() {
    if(!this.innerRangeValue[1]) return
    this.$emit('update:show', false)
    this.$emit('confirm', this.innerRangeValue);
  }

  handleCancel() {
    this.$emit('update:show', false)
  }
}
</script>

<style lang="scss" scoped>
.popup {

  .popup-header {
    padding: 16px;
    display: flex;
    align-items: center;

    .range {
      height: 32px;
      flex: 1;
      background-color: rgba(248, 249, 250, 1);
      border-radius: 8px;
      box-sizing: border-box;
      padding: 4px;
      display: flex;
      align-items: center;

      .range-item {
        flex: 1;
        height: 24px;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 12px;
        transition: all 0.3s;

        &:nth-child(1) {
          margin-right: 4px;
        }

        .label {
          font-size: 10px;
          font-weight: 600;
          color: rgba(153, 153, 153, 0.3);
        }
        .value {
          font-size: 12px;
          color: rgba(153, 153, 153, 1);
        }

        &.active {
          background-color: rgba(255, 255, 255, 1);
          .label {
            color: rgba(153, 153, 153, 1);
          }
          .value {
            color: rgba(0, 0, 0, 1);
          }
        }
      }
    }

    .btn {
      margin-left: 16px;
      width: 32px;
      height: 24px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: 600;
      color: rgba(13, 110, 253, 1);
    }
  }
}
.calendar {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  padding: 0;
  .row {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    padding:  0 20px;
    .col { 
      width: 14.28571428571429%;
      height: 48px;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 14px;
      font-weight: 700;
    }
  }

  .calendar-header {
    color: rgba(0, 0, 0, 0.3);
    padding: 16px 0;
    .row .col {
      font-weight: 400;
      font-size: 12px;
      height: 24px;
    }

    .month {
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 8px;
      padding: 0 20px;

      .month-text {
        font-size: 20px;
        font-weight: 700;
        color: rgba(0, 0, 0, 1);
      }

      .icon-right-circle { 
        width: 32px;
        height: 32px;
        background: url('~@/assets/store-visit/icon_right_circle.png') center/20px no-repeat;
        margin-left: -8px;
      }
      .icon-left-circle { 
        width: 32px;
        height: 32px;
        background: url('~@/assets/store-visit/icon_left_circle.png') center/20px no-repeat;
        margin-right: -8px;
      }
    }
  }

  .calendar-content {
    background-color: rgba(248, 249, 250, 1);
    height: 348px;
    padding: 20px 0;
    .day {
      color: rgba(0, 0, 0, 1);
      width: 100%;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
      position: relative;

      .day-text {
        width: 32px;
        height: 32px;
        border-radius: 50%;
        display: flex;
        justify-content: center;
        align-items: center;
        transition: all 0.3s;
      }

      .day-dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background-color: rgba(173, 181, 189, 1);
        position: absolute;
        bottom: 0;
        left: 50%;
        transform: translateX(-50%);
      }
      &.today {
        color: rgba(13, 110, 253, 1);
      }
      
      &.disabled {
        color: rgba(0, 0, 0, 0.3);
      }

      &.current {

        .day-text {
          color: rgba(255, 255, 255, 1);
          background-color: rgba(33, 37, 41, 1);
        }

        .day-dot {
          background-color: rgba(13, 110, 253, 1);
        }
      }
    }
  }
}
</style>