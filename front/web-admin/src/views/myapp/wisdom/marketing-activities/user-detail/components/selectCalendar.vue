<template>
  <div class="lookForMonth_wrapper">
    <div class="lookForMonth_top">
      <div class="selectDate">
        <div class="upDownSelect_item icon1" @click="dateUp"></div>
        <div>{{ year }} 年 {{ month }} 月</div>
        <div class="upDownSelect_item icon2" @click="dateDown"></div>
      </div>
    </div>
    <div class="calendar">
      <div
        v-for="(item, index) in calendarData"
        class="calendar_item"
        :key="index"
        :class="{
          notCheckedIn: item.signInDay && !item.isSignIn,
          CheckedIn: item.signInDay && item.isSignIn,
          checked: item.checked,
        }"
        @click="checkTime(item)"
      >
        <p class="dateEdit">
          {{ item.label }}
        </p>
        <p v-if="item.signInDay">
          {{ item.isSignIn ? "已签到" : "未签到" }}
        </p>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
@Component({})
export default class SelectInput extends Vue {
  @Prop() signInlist!: Array<any>;
  private calendarData: any = [
    { label: "日" },
    { label: "一" },
    { label: "二" },
    { label: "三" },
    { label: "四" },
    { label: "五" },
    { label: "六" },
  ]; //日历循环渲染数据
  private year = 0; //当前日期年
  private month = 0; //当前日期月数
  private date = 0; //当前日期号数
  private day = -1; //当前星期几
  private created() {
    // 获取当前日期数据
    this.getNow();
    // 获取当前月份一号的时间戳
    this.getCalendarDate(); // 给calendarData添加当月数据
  }

  // 获取当前时间
  private getNow() {
    const now = new Date();
    this.year = +now.getFullYear();
    this.month = +now.getMonth() + 1;
    this.date = +now.getDate();
    this.day = +now.getDay();
  }

  private monthDay(month: number) {
    if ([1, 3, 5, 7, 8, 10, 12].includes(month)) {
      return 31;
    } else if ([4, 6, 9, 11].includes(month)) {
      return 30;
    } else if (month === 2) {
      //  判断当年是否为闰年
      if (
        (this.year % 4 === 0 && this.year % 100 !== 0) ||
        this.year % 400 === 0
      ) {
        return 29;
      } else {
        return 28;
      }
    }
  }

  // 给calendarData添加当月数据
  private getCalendarDate() {
    // 获取当前月份一号星期几
    const firstDay: any = new Date(
      this.year + "-" + this.month + "-" + "01"
    ).getDay();
    this.calendarData = [
      { label: "日" },
      { label: "一" },
      { label: "二" },
      { label: "三" },
      { label: "四" },
      { label: "五" },
      { label: "六" },
    ];
    const num = parseInt(firstDay);
    const nowDays: any = this.monthDay(this.month);
    const lastMonth = this.month - 1 > 0 ? this.month - 1 : 12;
    const lastDays = this.monthDay(lastMonth);
    //   循环添加上一个月数据，空数组占位
    for (let i = 0; i < num; i++) {
      this.calendarData.push({});
    }
    // 循环添加当月数据
    for (let i = 0; i < nowDays; i++) {
      // 此处添加是否为签到日
      let boll = false;
      for (let j = 0; j < this.signInlist.length; j++) {
        if (`${this.year}-${this.month}-${i + 1}` == this.signInlist[j].time) {
          boll = true;
          this.calendarData.push({
            label: i + 1,
            signInDay: true,
            isSignIn: this.signInlist[j].isSignIn,
          });
        }
      }
      if (!boll) {
        this.calendarData.push({ label: i + 1 });
      }
    }
    // 默认选中当月最后一个签到的
    for (let i = this.calendarData.length - 1; i >= 0; i--) {
      if (this.calendarData[i].signInDay && this.calendarData[i].isSignIn) {
        this.calendarData[i].checked = true;
        this.$emit(
          "checkDate",
          `${this.year}-${this.month}-${this.calendarData[i].label}`
        );
        return;
      }
    }
  }

  // 将日期调上
  private dateUp() {
    this.month--;
    if (this.month <= 0) {
      this.year--;
      this.month = 12;
    }
    this.getCalendarDate(); // 给calendarData添加当月数据
  }
  // 将日期调下
  private dateDown() {
    this.month++;
    if (this.month > 12) {
      this.year++;
      this.month = 1;
    }
    this.getCalendarDate(); // 给calendarData添加当月数据
  }

  // 选择签到的日期
  private checkTime(item: any) {
    if (item.signInDay && item.isSignIn) {
      let clickIndex = 0;
      this.calendarData.forEach((mapItem: any, index: number) => {
        if (mapItem.label == item.label) {
          mapItem.checked = true;
          clickIndex = index;
        } else {
          mapItem.checked = false;
        }
      });
      this.$set(this.calendarData, clickIndex, this.calendarData[clickIndex]);
      this.$emit("checkDate", `${this.year}-${this.month}-${item.label}`);
    }
  }
}
</script>
<style lang="scss" scoped>
.lookForMonth_wrapper {
  padding: 20px;
  width: 481px;
  background: #f4f7fa;
}
.lookForMonth_top {
  margin-bottom: 20px;
  overflow: hidden;
  .selectTeacher {
    float: left;
  }
  .selectDate {
    height: 30px;
    line-height: 30px;
    display: flex;
    justify-content: space-between;
    .upDownSelect {
      display: flex;
      flex-direction: column;
      margin-top: -2px;
      margin-left: 5px;
    }
    .upDownSelect_item {
      width: 0;
      height: 0;
      border: 7px solid transparent;
      cursor: pointer;
      margin-top: 10px;
    }
    .icon1 {
      border-right: 7px solid #666;
      margin-bottom: 5px;
      &:hover {
        border-right: 7px solid skyblue;
      }
    }
    .icon2 {
      border-left: 7px solid #666;
      &:hover {
        border-left: 7px solid skyblue;
      }
    }
  }
}
/* 日历表样式=======================================↓ */
.calendar {
  width: 481px;
  display: flex;
  flex-wrap: wrap;
  box-sizing: border-box;
  .calendar_item {
    box-sizing: border-box;
    width: 60px;
    height: 60px;
    text-align: center;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    margin: 4px;
    .status {
      margin-top: 10px;
      &.textBlue {
        color: blue;
      }
      &.textRed {
        color: brown;
      }
    }
    .el-icon-edit-outline {
      cursor: pointer;
      margin-left: 7px;
    }
  }
  .dateEdit {
    margin-bottom: 10px;
  }
  .notCheckedIn {
    border: 1px solid #b9babe;
    color: #b9babe;
  }
  .CheckedIn {
    cursor: pointer !important;
    color: #7390f5;
    border: 1px solid #7390f5;
  }
  .checked {
    background: #708df5;
    color: white;
  }
}
</style>
