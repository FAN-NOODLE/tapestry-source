<template>
  <div class="invite-date">
    <div class="mask" @click="closeModel()"></div>
    <div class="content">
      <div class="left-date">
        <div class="scroll">
          <div
            @click="clickItem(index)"
            :class="{
              item: true,
              active: currentIndex == index,
            }"
            v-for="(item, index) in dataList"
            :key="index"
          >
            {{ item.value }}
          </div>
        </div>
        <div class="blur-box"></div>
      </div>
      <div class="right-content">
        <div class="top-title">
          <div>选择到店时间</div>
          <div class="close-box" @click="closeModel()">
            <img
              class="icon"
              src="https://cdn-cos.coach.com.cn/uploads/2023/v3.2.0.1/close-icon.png"
            />
          </div>
        </div>
        <div class="store-invite-date">
          <div class="scroll-invite">
            <div
              @click="clickInvite(index)"
              :class="{
                'invite-item': true,
                none: item.isFinish,
                active: dateIndex == index,
              }"
              v-for="(item, index) in dataList[currentIndex].date_list"
              :key="index"
            >
              <div class="name">{{ item.name }}</div>
              <div class="invite">{{ item.inviteNum }}人已预约</div>
            </div>
          </div>
        </div>
        <div class="sure-btn" @click="sumbit()">确定</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { getDays, getKsDays } from "@/utils/invite-date";
import { Component, Vue, Prop } from "vue-property-decorator";
import { Toast } from "vant";
import Helper from "@/utils/helper";
@Component({
  components: {
    // Toast,
  },
})
export default class extends Vue {
  @Prop(Number) allDay!: number;
  // eslint-disable-next-line @typescript-eslint/ban-types
  @Prop(Object) timeObj!: Object;
  @Prop(Array) timeSlot!: Array<any>;
  dataList = [];
  currentIndex = 0;
  dateIndex = -1;

  created() {
    this.getDataList();
  }
  // 预览图片
  clickItem(index: number) {
    const currentIndex = index;
    console.log(currentIndex);
    this.currentIndex = currentIndex;
    this.dateIndex = -1;
  }
  clickInvite(index: number) {
    const dateIndex = index;
    console.log(dateIndex);
    this.dateIndex = dateIndex;
  }
  sumbit() {
    if (this.dateIndex === -1) {
      Toast("请选择时间段");
      return;
    }
    const value = this.dataList[this.currentIndex].date_list[this.dateIndex];
    // todo: 拼接时间，不改变原有数据结构
    let select_name = this.dataList[this.currentIndex].select_name
    // todo: 判断select_name是否包含-，替换成/
    if(select_name.indexOf('-') > -1) {
      select_name = select_name.replace(/-/g, '/')
    }
    console.log(select_name)
    const startHour = value.name.split(' - ')[0]
    const endHour = value.name.split(' - ')[1]
    const startTime = new Date(select_name + ' ' + startHour).getTime() / 1000
    const endTime = new Date(select_name + ' ' + endHour).getTime() / 1000
    const obj = {
      ...value,
      select_name: select_name,
      startTime,
      endTime
    };
    this.$emit("sumbit", obj);
  }
  getDataList() {
    let startDate = new Date();
    let dataList = [];
    if (Helper.isKS() || Helper.isSW()) {
      startDate = new Date();
      dataList = getKsDays(startDate, this.allDay);
    } else {
      startDate = new Date(new Date().getTime() + 1000 * 60 * 60 * 24);
      dataList = getDays(startDate, this.allDay,this.timeSlot);
    }
    const timeObj = this.timeObj;
    const keys = Object.keys(timeObj);
    dataList.forEach((item) => {
      item.date_list.forEach((cell) => {
        if (keys.includes(cell.startTime + "")) {
          cell.inviteNum = timeObj[cell.startTime];
        }
      });
    });
    this.dataList = dataList;
  }
  closeModel() {
    this.$emit("close");
  }
}
</script>
<style lang="sass" scoped>
.invite-date
  position: fixed
  top: 0
  left: 0
  right: 0
  bottom: 0
  z-index: 10
  .mask
    position: absolute
    top: 0
    left: 0
    right: 0
    bottom: 0
    background: rgba(0,0,0,0.55)
  .content
    height: 400px
    position: absolute
    bottom: 0
    left: 0
    right: 0
    background: white
    display: flex
    justify-content: flex-start
    .left-date
      width: 90px
      background: #F6F5F5
      position: relative
      padding-bottom: 51px
      flex-shrink: 0
      .blur-box
        width: 90px
        height: 51px
        background: linear-gradient(180deg, #F6F5F5 0%, #F6F5F5 100%)
        filter: blur(2px)
        position: absolute
        bottom: 0
        left: 0
        z-index: 2
    .scroll
      height: 349px
      overflow-y: auto
      .item
        width: 90px
        text-align: center
        font-size: 14px
        line-height: 46px
        font-weight: 400
        color: #777777
        &.active
          font-family: PingFangSC-Medium, PingFang SC
          font-weight: 500
          color: #4A4A4A
          background: white
    .top-title
      display: flex
      justify-content: space-between
      height: 46px
      line-height: 14px
      align-items: center
      padding-left: 20px
      font-size: 14px
      font-family: PingFangSC-Medium, PingFang SC
      font-weight: 500
      color: #4A4A4A
      position: relative
      .close-box
        position: absolute
        height: 46px
        width: 46px
        padding: 19px 13px 9px
        top: 0
        right: 0
        .icon
          width: 18px
          height: 18px
    .right-content
      width: 285px
      .store-invite-date
        padding-left: 13px
      .scroll-invite
        height: 250px
        overflow-y: auto
        .invite-item
          display: flex
          justify-content: space-between
          align-items: center
          height: 42px
          padding: 13px 15px 15px 7px
          &.none
            display: none
          &.active
            background: #F6F5F5
          .name
            font-size: 14px
            color: #4A4A4A
          .invite
            font-size: 12px
            color: #BDBEC8
      .sure-btn
        width: 84px
        line-height: 36px
        background: #ffffff
        border: 1px solid #000000
        font-size: 14px
        font-family: PingFangSC-Regular, PingFang SC
        font-weight: 400
        color: #000000
        text-align: center
        margin: 10px 0 0 185px
</style>
