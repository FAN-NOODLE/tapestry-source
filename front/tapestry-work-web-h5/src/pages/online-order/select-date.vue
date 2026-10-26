<template>
  <div class="container">
    <div class="content">
      <div class="list">
        <div class="item" :class="{ active: item.value===filterParams.date_type }" v-for="item, index in list" :key="index" @click="handleChange(item)">{{ item.text }}</div>
      </div>
      <div v-if="filterParams.date_type===5" class="date" @click="handleSelectDate">
        <div class="date-value">
          <span v-if="filterParams.start_time">{{$helper.dateFormat(filterParams.start_time, "yyyy-MM-dd")}}</span>
          <span v-else class="placeholder">开始时间</span>
        </div>
        <div>-</div>
        <div class="date-value">
          <span v-if="filterParams.start_time">{{$helper.dateFormat(filterParams.end_time, "yyyy-MM-dd")}}</span>
          <span v-else class="placeholder">结束时间</span>
        </div>
      </div>
    </div>
    <div class="footer">
      <van-button plain color="#000000" @click="handleReset">重置</van-button>
      <van-button color="#000000" @click="handleConfirm">确认</van-button>
    </div>
    <van-calendar 
      v-model:show="show" 
      type="range" 
      :min-date="minDate"
      :max-date="maxDate"
      allow-same-day
      :default-date="[new Date(filterParams.start_time * 1000), new Date(filterParams.end_time * 1000)]"
      @confirm="handleDateConfirm"
    />
  </div>
</template>
<script lang="ts">
import dayjs from 'dayjs';
import { Component, Vue } from 'vue-property-decorator';

@Component({

})
export default class OnlineOrderSelectDate extends Vue {
  filterParams:any = {}
  list = [
    {
      text: '昨天',
      value: 1
    },
    {
      text: '今天',
      value: 2
    },
    {
      text: '近7天',
      value: 3
    },
    {
      text: '近30天',
      value: 4
    },
    {
      text: '自定义时间',
      value: 5
    },
  ]
  active = 2
  show = false

  get minDate() {
    return new Date(2024, 0, 1)
  }
  
  get maxDate() {
    return new Date()
  }

  mounted() {
    this.filterParams = JSON.parse(JSON.stringify(this.$store.getters["onlineOrder/getFilterParams"]))
  }

  handleChange(item:any) {
    this.filterParams.date_type = item.value
    this.updateDate(item.value)
  }
  updateDate(type: number) {
    switch (type) {
      case 1:
        // 昨天 00:00:00 - 23:59:59 时间戳
        this.filterParams.start_time = dayjs(dayjs().subtract(1, 'day').format('YYYY-MM-DD') + ' 00:00:00').unix()
        this.filterParams.end_time = dayjs(dayjs().subtract(1, 'day').format('YYYY-MM-DD') + ' 23:59:59').unix()
        break;
      
      case 2:
        // 今天 00:00:00 - 当前时间 时间戳
        this.filterParams.start_time = dayjs(dayjs().format('YYYY-MM-DD') + ' 00:00:00').unix()
        this.filterParams.end_time = dayjs().unix()
        break;

      case 3:
        // 近7天 00:00:00 - 当前时间 时间戳
        this.filterParams.start_time = dayjs(dayjs().subtract(6, 'day').format('YYYY-MM-DD') + ' 00:00:00').unix()
        this.filterParams.end_time = dayjs().unix()
        break;

      case 4:
        // 近30天 00:00:00 - 当前时间 时间戳
        this.filterParams.start_time = dayjs(dayjs().subtract(29, 'day').format('YYYY-MM-DD') + ' 00:00:00').unix()
        this.filterParams.end_time = dayjs().unix()
        break;
      
      default:
        break;
    }

  }
  handleSelectDate() {
    this.show = true
  }
  handleDateConfirm(value) {
    console.log(value);
    this.show = false
    this.filterParams.start_time = dayjs(dayjs(value[0]).format('YYYY-MM-DD') + ' 00:00:00').unix()
    this.filterParams.end_time = dayjs(dayjs(value[1]).format('YYYY-MM-DD') + ' 23:59:59').unix()
  }
  handleReset() {
    this.$store.commit("onlineOrder/setFilterParams", {
      ...this.filterParams,
      date_type: 0,
      start_time: 0,
      end_time: 0
    });
    this.$store.commit("onlineOrder/setUpdate", true)
    this.$router.back()
  }
  handleConfirm() {
    this.$store.commit("onlineOrder/setFilterParams", JSON.parse(JSON.stringify(this.filterParams)))
    this.$store.commit("onlineOrder/setUpdate", true)
    const { from } = this.$route.query
    const { date_type, start_time, end_time } = this.filterParams
    const params = {
      current_page: ['订单数据', '订单列表'][Number(from)||0],
      filter_type: '时间段',
      filter_option: ''
    }

    if(date_type < 5) {
      params.filter_option = date_type == 0 ? '全部' : this.list.find(item => item.value==date_type).text
    } else {
      params.filter_option = [dayjs(start_time*1000).format('YYYY-MM-DD'),dayjs(end_time*1000).format('YYYY-MM-DD')].join(' - ')
    }


    // 神策数据上报 -- 订单筛选
    this.$track.trackOrderFilter(params)

    this.$store.commit("onlineOrder/setUpdate", true)
    this.$router.back()
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .van-icon {
    position: absolute;
    font-size: 20px;
  }
  .van-button--danger {
    background-color: #000000;
    border-color: #000000;
  }

  .van-calendar__day--middle {
    color: #000000;
  }

  .van-calendar__day--end, .van-calendar__day--start, .van-calendar__day--start-end, .van-calendar__day--multiple-middle, .van-calendar__day--multiple-selected {
    background-color: #000000;
  }
}

.container {
  width: 100vw;
  min-height: 100vh;
  background: #F7F7F7;

  .content {
    padding: 0 15px;

    .item {
      padding: 15px 0;
      font-weight: 400;
      font-size: 14px;
      color: #666666;
      line-height: 20px;
      text-align: center;

      &.active {
        font-weight: bold;
        position: relative;

        &::after {
          content: '';
          position: absolute;
          width: 16px;
          height: 16px;
          background: url('~@/assets/online-order/icon_check_2.png') center/100% no-repeat;
          right: 0;
          top: 50%;
          transform: translateY(-50%);
        }
      }
    }

    .date {
      display: flex;
      align-items: center;
      padding: 10px 15px;
      border: 1px solid #ccc;
      font-size: 14px;
      color: #000000;
      line-height: 20px;

      .date-value {
        flex: 1;
        text-align: center;

        .placeholder {
          color: #999;
        }
      }
    }
  }

  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 10px 15px;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding-bottom: calc(10px + contain(safe-area-inset-bottom));
    padding-bottom: calc(10px + env(safe-area-inset-bottom));
    box-shadow: 0px 1px 6px 0px rgba(0,0,0,0.15);

    .van-button {
      flex: 1;

      &:nth-child(2) {
        margin-left: 15px;
      }
    }
  }
}

</style>