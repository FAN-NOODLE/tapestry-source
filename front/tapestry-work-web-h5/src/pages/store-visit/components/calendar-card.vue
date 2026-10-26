<template>
  <div class="calendar-wrap">
    <div class="calendar-opt">
      <div class="year">{{currentYear}}.</div>
      <div class="flex">
        <SelectDatePopup type="year-month" :value="defaultValue" @confirm="handleYearMonthChange">
          <div class="month">
            <div class="month-text">{{showMonth}}</div>
            <div class="icon-down-circle"></div>
          </div>
        </SelectDatePopup>
        <slot></slot>
      </div>
    </div>
    <div class="calendar">
      <div class="calendar-header">
        <div class="row">
          <div class="col" v-for="(item, index) in week" :key="index">{{ item }}</div>
        </div>
      </div>
      <div class="calendar-content" @touchstart="handleTouchStart" @touchend="handleTouchEnd">
        <div class="row">
          <div class="col" v-for="(item, index) in days" :key="index">
            <div class="day" :class="{
              'disabled': disabled(item),
              'current': item.unixDate === currentDate,
              'today': item.isToday
            }" @click="handleClick(item)">
              <div class="day-text">{{ item.dayText }}</div>
              <div v-if="schedules.includes(item.unixDate)" class="day-dot"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
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
export default class CalendarCard extends Vue {
  @Model('update:current', { type: String, default: '' }) current: string
  @Prop({ type: Array, default: () => [] }) schedules: string[]
  @Prop({ type: Number, default: new Date().setHours(0, 0, 0, 0) }) minDate: number
  
  week = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];
  days = []
  // 当前选中的日期 年-月-日
  currentDate = this.getCurrentDate();
  // 今天
  today = this.getCurrentDate()

  // 当前年份
  currentYear = new Date().getFullYear();
  // 当前月份
  currentMonth = new Date().getMonth() + 1;

  startX = 0;

  get showMonth() {
    return months[this.currentMonth].abbreviation
  }

  get disabled() {
    return (item) => {
      return item.unix < this.minDate
    }
  }

  get defaultValue() {
    return `${this.currentYear}-${this.currentMonth}-1`
  }

  @Watch('current', {
    immediate: true,
    deep: true
  })
  async onCurrentChange(val: string) {
    this.currentDate = val || this.getCurrentDate();
    
  }

  mounted() {
    this.init();
  }
  init() {
    const date = new Date();
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
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
        day: '',
        unix: 0,
        unixDate: '',
      });
    }
    for (let i = 1; i <= days; i++) {
      const unixDate = dayjs(new Date(year, month - 1, i)).format('YYYY-MM-DD')
      arr.push({
        day: i,
        dayText: String(i),
        unix: new Date(year, month - 1, i).getTime(),
        unixDate: unixDate,
        isToday: unixDate === this.today,
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
    // if(this.disabled(item)) {
    //   return;
    // }
    // 神策数据上报 -- 日历页点击
    this.$track.trackCalendarPageClick({
      button_type: '选择日期',
      button_name: item.dayText,
    })
    this.currentDate = item.unixDate;
    this.$emit('update:current', item.unixDate);
    this.$emit("change", item);


  }
  /**
   * 切换年月月份
   */
  handleYearMonthChange(date: Date) {
    console.log(date);
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    this.currentYear = year;
    this.currentMonth = month;
    this.days = this.getDaysInMonth(year, month);

    // 神策数据上报 -- 日历页点击
    this.$track.trackCalendarPageClick({
      button_type: '切换月份',
      button_name: this.showMonth,
    })
    
    this.$emit('change-month', {
      year,
      month,
    })

  }
  handleTouchStart(e) {
    this.startX = e.changedTouches[0].pageX;
  }
  handleTouchEnd(e) { 
    const endX = e.changedTouches[0].pageX;
    if (Math.abs(endX - this.startX) > 50) {
      if (endX > this.startX) {
        this.handleMonthChange(-1);
      } else {
        this.handleMonthChange(1);
      }
    }
  }

  handleMonthChange(value) {
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
    this.$emit('change-month', {
      year:  this.currentYear,
      month: this.currentMonth,
    })
  }
}
</script>

<style lang="scss" scoped>
.calendar-wrap {
  width: 100%;
}
.calendar-opt {
  padding: 16px 24px 0;
  .year {
    font-size: 12px;
    font-weight: 700;
    color: rgba(153, 153, 153, 1);
  }
  .flex {
    display: flex;
    ::v-deep .select {
      flex: 1;
      justify-content: flex-start;
      display: flex;
    }

    .month {
      display: flex;
      align-items: center;

      .month-text {
        font-size: 24px;
        font-weight: 700;
        color: rgba(0, 0, 0, 1);
      }

      .icon-down-circle {
        margin-left: 4px;
        width: 24px;
        height: 24px;
        background: url('~@/assets/store-visit/icon_down_circle.png') no-repeat center center/100%;
      }
    }
  }
}
.calendar {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  padding: 0 12px;
  .row {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
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
    .row .col {
      font-weight: 400;
      font-size: 12px;
      height: 40px;
    }
  }

  .calendar-content {

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
      
      &.disabled {
        color: rgba(0, 0, 0, 0.3);
      }


      &.today {
        .day-text {
          color: #0D6EFD;
        }
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