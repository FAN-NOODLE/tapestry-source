<template>
  <div class="content">
    <div class="filter">
      <div class="date" :class="{active: active==2}" @click="handleTabChange(2)">{{value}}</div>
      <div class="tab">
        <div class="tab-item" :class="{active: active==0}" @click="handleTabChange(0)">月榜</div>
        <div class="tab-item" :class="{active: active==1}" @click="handleTabChange(1)">年榜</div>
      </div>
    </div>
    <div class="section">
      <div class="title">
        <div class="text">已发货商品款_色 TOP100</div>
      </div>
      <van-list class="list">
        <div class="item" v-for="(item, index) in goodsList" :key="index">
          <div class="item-index">
            <div class="cell" v-if="index<3">
              <img class="rank-icon" v-if="index==0" src="@/assets/online-order/icon_rank_1.png" alt=""></img>
              <img class="rank-icon" v-if="index==1" src="@/assets/online-order/icon_rank_2.png" alt=""></img>
              <img class="rank-icon" v-if="index==2" src="@/assets/online-order/icon_rank_3.png" alt=""></img>
            </div>
            <div class="cell" v-else>
              <div class="rank-num">{{index+1}}</div>
            </div>
          </div>
          <van-image class="item-image" :src="item.image_url || imageEmpty"></van-image>
          <div class="item-content">
            <div class="item-name">
              <span class="item-d">{{ item.dept||'未知' }}</span>
              <span class="item-n">{{ item.sku_info.replaceAll(' ', '_') }}</span>
              <div class="item-c">{{ item.num }}件</div>
            </div>
            <div class="item-desc">{{ item.sku_name }}</div>
          </div>
        </div>
        <emptyPage v-if="goodsList.length==0"></emptyPage>
      </van-list>
    </div>
    <van-calendar 
      v-model:show="show" 
      type="range" 
      :min-date="minDate"
      :max-date="maxDate"
      allow-same-day
      :default-date="[new Date(searchDate.start_time*1000), new Date(searchDate.end_time*1000)]"
      @confirm="handleDateConfirm"
    />
  </div>
</template>
<script lang="ts">
import { getGoodsRankApi, getMonthListApi, getSaRankListApi, getStoreRankListApi } from '@/api/online-order';
import dayjs from 'dayjs';
import { Component, Vue } from 'vue-property-decorator';
import emptyPage from "@/components/base/empty-page.vue";

@Component({
  components: {
    emptyPage
  }
})
export default class SaleRank extends Vue {
  imageEmpty = require('@/assets/material-library/v4.2.0/icon_empty.png')

  active = 0

  list: any[] = []
  saList:any[] = [] 
  goodsList: any[] = []

  customDate = []
  show = false
  get minDate() {
    return new Date(2024, 0, 1)
  }
  
  get maxDate() {
    return new Date()
  }

  get value() {
    if(this.active === 0) {
      const current = this.list.find(item => item.now == 1)
      if(!current) return ''
      return `${dayjs(current.start_time*1000).format('MM月DD日')}-${dayjs(current.end_time*1000).format('MM月DD日')}`
    } else if (this.active === 1) {
      const current = this.list.find(item => item.now == 1)
      // const year = 
      const data:any = {}
      this.list.forEach(item => {
        if(data[item.year]) {
          data[item.year].push(item)
        } else {
          data[item.year] = [item]
        }
      })

      console.log(data[current.year]);
      return data[current.year][0].year
    } else {
      return `${dayjs(this.customDate[0]*1000).format('MM月DD日')}-${dayjs(this.customDate[1]*1000).format('MM月DD日')}`
    }
  }

  get searchDate() {
    if(this.active === 0) {
      const current = this.list.find(item => item.now == 1)
      if(!current) return ''
      return {
        start_time: current.start_time,
        end_time: current.end_time
      }
    } else if (this.active === 1) {
      const current = this.list.find(item => item.now == 1)
      // const year = 
      const data:any = {}
      this.list.forEach(item => {
        if(data[item.year]) {
          data[item.year].push(item)
        } else {
          data[item.year] = [item]
        }
      })

      return {
        start_time: data[current.year][0].start_time,
        end_time: data[current.year][11].end_time
      }
    } else {
      return {
        start_time: this.customDate[0],
        end_time: this.customDate[1]
      }
    }
  }

  async mounted() {
    await this.getMonthList()
    this.getGoodsRank()
  }

  handleTabChange(index: number) {
    if(index === 0) {
      this.active = index
      this.getGoodsRank()
      // 神策数据上报 -- 订单筛选
      this.$track.trackOrderFilter({
        current_page: '商品排行榜',
        filter_type: '榜单类型',
        filter_option: '月榜'
      })
    }else if(index === 1) {
      this.active = index
      this.getGoodsRank()
      // 神策数据上报 -- 订单筛选
      this.$track.trackOrderFilter({
        current_page: '商品排行榜',
        filter_type: '榜单类型',
        filter_option: '年榜'
      })
    } else {
      this.show = true
    }
  }

  async getMonthList() {
    
    const res = await getMonthListApi()
    
    this.list = res.data
  }

  async getGoodsRank() {
    const res = await getGoodsRankApi(this.searchDate)
    
    this.goodsList = res.data||[]
  }

  handleDateConfirm(value: any) {
    console.log(value);
    this.active = 2;
    this.show = false;
    this.customDate = [
      parseInt(value[0].getTime() / 1000 + ''),
      parseInt(new Date(value[1]).setHours(23,59,59) / 1000 + '')
    ]
    this.getGoodsRank()

    // 神策数据上报 -- 订单筛选
    this.$track.trackOrderFilter({
      current_page: '商品排行榜',
      filter_type: '时间段',
      filter_option: [dayjs(value[0]).format('YYYY-MM-DD'), dayjs(value[1]).format('YYYY-MM-DD')].join(' - ')
    })
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

  .default {
    height: auto;
    padding: 50px 0 !important;
  }
}
.content {
  padding: 0 15px 15px;
  
  .filter {
    display: flex;
    justify-content: space-between;
    position: sticky;
    top: 44px;
    z-index: 3;
    background-color: #f7f7f7;
    padding: 15px 0;

    .date {
      // width: 162px;
      height: 30px;
      border: 1px solid #A6A6A6;
      font-weight: 400;
      font-size: 14px;
      color: #A6A6A6;
      display: flex;
      align-items: center;
      justify-content: space-between;
      text-align: center;
      padding: 0 10px;

      &.active {
        border-color: #000000;
        color: #000000;
      }

      &::before {
        content: "";
        width: 9px;
        height: 14px;
        background: url('~@/assets/online-order/icon_right.png') center/100% no-repeat;
        margin-right: 10px;
        transform: rotate(180deg);
      }
      &::after {
        content: "";
        width: 9px;
        height: 14px;
        background: url('~@/assets/online-order/icon_right.png') center/100% no-repeat;
        margin-left: 10px;
      }
    }
    .tab {
      display: flex;

      .tab-item {
        width: 47px;
        height: 30px;
        border: 1px solid #A6A6A6;
        font-weight: 400;
        font-size: 14px;
        color: #A6A6A6;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-left: 5px;

        &.active {
          background-color: #000000;
          border-color: #000000;
          color: #fff;
        }
      }
    }
  }

  .section {
    background-color: #fff;
    margin-bottom: 15px;
  
    .title {
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: center;
  
      .text {
        font-size: 14px;
        color: #000000;
        padding: 0 15px;
        white-space: nowrap;
      }
  
      &::before,
      &::after {
        content: "";
        width: 93px;
        height: 1px;
        background: #F0F0F0;
      }
    }

    .list {
      padding: 0 15px;
      margin-top: -15px;

      .item {
        display: flex;
        align-items: center;
        padding: 15px 0;
        border-bottom: 1px solid #F2F2F2;

        .item-index {
          width: 40px;
          height: 40px;
          display: flex;
          align-items: center;
          justify-content: center;
        }

        .rank-icon {
          width: 17px;
          height: 17px;
          display: block;
          margin: auto;
        }

        .rank-num {
          width: 17px;
          height: 17px;
          line-height: 17px;
          text-align: center;
          margin: auto;
          border-radius: 50%;
          background-color: #F0F0F0;
          color: #404040;
          font-size: 11px;
        }

        .item-image {
          width: 60px;
          height: 60px;
        }

        .item-content {
          flex: 1;
          overflow: hidden;
          margin-left: 10px;
          font-size: 14px;
        }

        .item-name {
          // display: flex;
          // align-items: flex-start;
          font-size: 14px;
          color: #000;
          position: relative;
          padding-right: 50px;
          white-space: wrap;
          word-break: break-all;

          .item-d {
            padding: 0 7px;
            background-color: #000;
            color: #fff;
            font-size: 12px;
            margin-right: 5px;
          }

          .item-n {
            margin-right: 5px;
          }

          .item-c {
            font-weight: bold;
            position: absolute;
            top: 0;
            right: 0;
          }

        }

        .item-desc {
          font-size: 12px;
          color: #96979E;
          margin-top: 8px;
        }
      }
    }
  }
  
}


</style>