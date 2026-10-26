<template>
  <van-popup
    v-model="show"
    @close="groupCancel"
    @opened="init"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    round
    position="bottom"
    class="material-tag"
  >
    <div class="wrap">
      <div class="header">
        <span
          :class="[currentPick === 'start' ? 'on' : '']"
          @click="setTime('start')"
          >{{
            this.timeRange.start ? this.timeRange.start : "请选择开始时间"
          }}</span
        >
        <i>-</i>
        <span
          :class="[currentPick === 'end' ? 'on' : '']"
          @click="setTime('end')"
          >{{
            this.timeRange.end ? this.timeRange.end : "请选择结束时间"
          }}</span
        >
      </div>
      <div class="title">
        <span>年</span>
        <span>月</span>
        <span>日</span>
      </div>
      <van-datetime-picker
        ref="time-range"
        v-model="currentDate"
        :show-toolbar="false"
        type="date"
        :visible-item-count="5"
        :min-date="minDate"
        :max-date="maxDate"
        @change="pickTime"
      />
    </div>
    <div class="btn-view">
      <div class="btn-li reset" @click="groupCancel">取消</div>
      <div
        class="btn-li confirm"
        :class="[btnDisabled ? 'empty' : '']"
        @click="groupCheck"
      >
        确定
      </div>
    </div>
  </van-popup>
</template>

<script lang="ts">
import { Component, Vue, PropSync, Watch } from "vue-property-decorator";
import { Popup, DatetimePicker } from "vant";

@Component({
  components: {
    [Popup.name]: Popup,
    [DatetimePicker.name]: DatetimePicker,
  },
})
export default class MaterialTag extends Vue {
  @PropSync("showPopup", { type: Boolean }) show;

  private currentPick = "start";
  private currentDate = new Date();
  private minDate = new Date();
  private maxDate = new Date();
  private timeRange = {
    start: "",
    end: "",
  };
  /** 确认按钮禁用 */
  private get btnDisabled() {
    return !this.timeRange.start || !this.timeRange.end;
  }
  /** 分组筛选取消 */
  private groupCancel() {
    this.$emit("cancel");
  }
  /** 分组筛选确定 */
  private groupCheck() {
    if (this.timeRange.start && this.timeRange.end) {
      this.$emit("confirm", { ...this.timeRange });
    }
  }
  /** 选择时间变化时 */
  private pickTime(e) {
    console.log("pickTime", e.getValues());
    const timeArr = e.getValues();
    let time = `${timeArr[0]}-${timeArr[1]}-${timeArr[2]}`;
    if (this.currentPick === "start") {
      //选择开始时间
      time += " 00:00:00";
      if (this.timeRange.end) {
        //有结束时间时，进行对比，合法的赋值
        const start = time.replace(/-/g, "/");
        const startTime = new Date(start).getTime();

        const end = this.timeRange.end.replace(/-/g, "/");
        const endTime = new Date(end).getTime();
        if (startTime <= endTime) {
          this.timeRange[this.currentPick] = time;
        } else {
          this.$toast("开始时间必须早于结束时间");
        }
      } else {
        this.timeRange[this.currentPick] = time;
      }
    } else {
      time += " 23:59:59";
      //选择结束时间
      if (this.timeRange.start) {
        //有开始时间时，进行对比，合法的赋值
        const start = this.timeRange.start.replace(/-/g, "/");
        const startTime = new Date(start).getTime();

        const end = time.replace(/-/g, "/");
        const endTime = new Date(end).getTime();
        if (endTime >= startTime) {
          this.timeRange[this.currentPick] = time;
        } else {
          this.$toast("结束时间必须晚于开始时间");
        }
      } else {
        this.timeRange[this.currentPick] = time;
      }
    }
  }
  private setTime(type) {
    this.currentPick = type;
    //首次切换到结束时间，把当前开始时间复赋值进来
    if (type === "end" && !this.timeRange.end) {
      const timeArr = (this.$refs["time-range"] as HTMLFormElement)
        .getPicker()
        .getValues();
      let time = `${timeArr[0]}-${timeArr[1]}-${timeArr[2]}`;
      time += " 23:59:59";
      this.timeRange.end = time;
    }
  }
  /** 初始化日历开始结束时间*/
  private init() {
    this.timeRange.start = "";
    this.timeRange.end = "";
    this.currentPick = "start";

    const min = new Date().getTime() / 1000;
    const maxDate = min + 3600 * 24 * 365 * 12;
    const minString = this.$helper.dateFormat(min, "yyyy/MM/dd");
    const maxString = this.$helper.dateFormat(maxDate, "yyyy/MM/dd");
    console.log("time-range", minString, maxString);
    this.minDate = new Date(minString);
    this.maxDate = new Date(maxString);
    //等待组件初始化
    (this.$refs["time-range"] as HTMLFormElement)
      .getPicker()
      .setIndexes([0, 0, 0]);

    //默认起始为当前时间
    this.timeRange.start = minString.replace(/\//g, "-") + " 00:00:00";
  }
}
</script>

<style lang="scss" scoped>
.material-tag {
  height: 453px;
  padding: 15px 0 0;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  border-radius: 0;
  .wrap {
    flex-grow: 1;
    padding: 0 10px 0;
  }
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #000000;
    margin-bottom: 20px;
    height: 36px;
    line-height: 36px;
    span {
      width: 158px;
      border: 1px solid #eeeeee;
      font-size: 14px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      color: #999999;
      text-align: center;
      &.on {
        color: #000;
      }
    }
    i {
      margin: 0 10px;
    }
  }
  .title {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #96979e;
    text-align: center;
  }
  .btn-view {
    flex-shrink: 1;
    width: 100%;
    display: flex;
    background-color: #ffffff;
    padding: 15px;
    .btn-li {
      flex: 2;
      border: 1px solid #4a4a4a;
      color: #4a4a4a;
      font-size: 14px;
      text-align: center;
      line-height: 36px;
    }
    .confirm {
      margin-left: 15px;
      color: #ffffff;
      background-color: #4a4a4a;
    }
    .empty {
      color: #ffffff;
      background-color: #cccccc;
      border-color: #cccccc;
    }
  }
}
</style>
