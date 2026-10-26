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
        <div class="text">门店销量排行榜</div>
      </div>
      <div class="head">
        <van-row class="row" gutter="10">
          <van-col span="4">
            <div class="cell">排名</div>
          </van-col>
          <van-col span="12">
            <div class="cell">门店名称</div>
          </van-col>
          <van-col span="4">
            <div class="cell">发货商品</div>
          </van-col>
          <van-col span="4">
            <div class="cell">发货订单</div>
          </van-col>
        </van-row>
      </div>
      <van-list class="body">
        <div class="row" v-for="(item, index) in storeList" :key="index">
          <van-row gutter="10">
            <van-col span="4">
              <div class="cell" v-if="index<3">
                <img class="rank-icon" v-if="index==0" src="@/assets/online-order/icon_rank_1.png" alt=""></img>
                <img class="rank-icon" v-if="index==1" src="@/assets/online-order/icon_rank_2.png" alt=""></img>
                <img class="rank-icon" v-if="index==2" src="@/assets/online-order/icon_rank_3.png" alt=""></img>
              </div>
              <div class="cell" v-else>
                <div class="rank-num">{{index+1}}</div>
              </div>
            </van-col>
            <van-col span="12">
              <div class="cell">{{ item.store_name }}</div>
            </van-col>
            <van-col span="4">
              <div class="cell">{{item.sku_num}}</div>
            </van-col>
            <van-col span="4">
              <div class="cell">{{item.num}}</div>
            </van-col>
          </van-row>
        </div>
        <emptyPage v-if="storeList.length==0"></emptyPage>
      </van-list>
    </div>
    <div class="section">
      <div class="title">
        <div class="text">SA销量TOP10</div>
      </div>
      <div class="list">
        <div class="row" v-for="(item, index) in saList" :key="index">
          <van-row gutter="10">
            <van-col span="14">
              <div class="user">
                <div class="user-avatar">
                  <van-image
                    width="50px"
                    height="50px"
                    fit="cover"
                    round
                    lazy-load
                    :src="item.avatar"
                  />
                  <img v-if="index===0" src="@/assets/online-order/icon_crown_1.png" class="user-crown">
                  <img v-if="index===1" src="@/assets/online-order/icon_crown_2.png" class="user-crown">
                  <img v-if="index===2" src="@/assets/online-order/icon_crown_3.png" class="user-crown">
                  <img v-if="index===0" src="@/assets/online-order/icon_no_1.png" class="user-no">
                  <img v-if="index===1" src="@/assets/online-order/icon_no_2.png" class="user-no">
                  <img v-if="index===2" src="@/assets/online-order/icon_no_3.png" class="user-no">
                </div>
                <div class="flex-1">
                  <div class="user-name">{{item.member_name}}</div>
                  <div class="user-store">{{item.store_name}}</div>
                </div>
              </div>
            </van-col>
            <van-col span="5">
              <div class="cell">
                <div class="value">{{item.sku_num}}</div>
                <div class="label">发货商品</div>
              </div>
            </van-col>
            <van-col span="5">
              <div class="cell">
                <div class="value">{{item.num}}</div>
                <div class="label">发货订单</div>
              </div>
            </van-col>
          </van-row>
        </div>
        <emptyPage v-if="saList.length==0"></emptyPage>
      </div>
    </div>
    <van-calendar 
      ref="calendar"
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
import { getMonthListApi, getSaRankListApi, getStoreRankListApi } from '@/api/online-order';
import dayjs from 'dayjs';
import { Component, Vue } from 'vue-property-decorator';
import emptyPage from "@/components/base/empty-page.vue";

@Component({
  components: {
    emptyPage
  }
})
export default class SaleRank extends Vue {

  active = 0

  list: any[] = []
  saList:any[] = [] 
  storeList: any[] = []

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

      console.log(data[current.year]);
      

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
    this.getSaRankList()
    this.getStoreRankList()
  }

  handleTabChange(index: number) {
    if(index === 0) {
      this.active = index
      this.getSaRankList()
      this.getStoreRankList()
      // 神策数据上报 -- 订单筛选
      this.$track.trackOrderFilter({
        current_page: '销售排行榜',
        filter_type: '榜单类型',
        filter_option: '月榜'
      })
    }else if(index === 1) {
      this.active = index
      this.getSaRankList()
      this.getStoreRankList()
      // 神策数据上报 -- 订单筛选
      this.$track.trackOrderFilter({
        current_page: '销售排行榜',
        filter_type: '榜单类型',
        filter_option: '年榜'
      })
    } else {
      this.show = true;
      (this.$refs.calendar as any).getInitialDate()
    }
  }

  async getMonthList() {
    
    const res = await getMonthListApi()
    
    this.list = res.data
  }

  async getSaRankList() {
    const res = await getSaRankListApi(this.searchDate)
    
    this.saList = res.data
  }

  async getStoreRankList() {
    const res = await getStoreRankListApi(this.searchDate)
    this.storeList = res.data
  }

  handleDateConfirm(value: any) {
    console.log(value);
    this.active = 2;
    this.show = false;
    this.customDate = [
      parseInt(value[0].getTime() / 1000 + ''),
      parseInt(new Date(value[1]).setHours(23,59,59) / 1000 + '')
    ]
    this.getSaRankList()
    this.getStoreRankList()

    // 神策数据上报 -- 订单筛选
    this.$track.trackOrderFilter({
      current_page: '销售排行榜',
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
      }
  
      &::before,
      &::after {
        content: "";
        width: 93px;
        height: 1px;
        background: #F0F0F0;
      }
    }

    .head {
      background-color: #FAFAFA;
      .cell {
        font-weight: 400;
        font-size: 11px;
        color: #676767;
        padding: 10px 0;
        line-height: 16px;
        text-align: center;
      }
    }
    .body {
      height: auto;
      padding-bottom: 15px;
      .cell {
        font-weight: 400;
        font-size: 12px;
        color: #404040;
        padding: 15px 0;
        line-height: 17px;
        text-align: center;
      }

      .row {
        position: relative;

        &::after {
          content: "";
          position: absolute;
          bottom: 0;
          left: 50%;
          transform: translateX(-50%);
          width: 315px;
          height: 1px;
          background-color: #F2F2F2;
        }
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
    }

    .list {
      padding: 0 15px;
      margin-top: -15px;
      .row {
        padding: 15px 0;
        position: relative;
        &:not(:last-child)::after {
          content: "";
          position: absolute;
          bottom: 0;
          left: 50%;
          transform: translateX(-50%);
          width: 315px;
          height: 1px;
          background-color: #F2F2F2;
        }
      }

      .user {
        display: flex;
        align-items: center;

        .user-avatar {
          position: relative;

          .van-image {
            display: block;
            position: relative;
            z-index: 1;
          }

          .user-crown {
            width: 18px;
            height: 16px;
            position: absolute;
            top: -5px;
            left: 0;
          }
          .user-no {
            width: 40px;
            height: 16px;
            position: absolute;
            bottom: 0;
            left: 50%;
            transform: translateX(-50%);
            z-index: 2;
          }
        }

        .flex-1 {
          flex: 1;
          margin-left: 8px;
          overflow: hidden;
        }

        .user-name {
          font-weight: 600;
          font-size: 16px;
          color: #000000;
          line-height: 22px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .user-store {
          font-weight: 400;
          font-size: 11px;
          color: #A6A6A6;
          line-height: 16px;
          margin-top: 5px;
        }
      }

      .cell {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;

        .value {
          font-weight: 400;
          font-size: 16px;
          color: #000000;
          line-height: 22px
        }

        .label {
          font-weight: 400;
          font-size: 11px;
          color: #A6A6A6;
          line-height: 16px;
          margin-top: 5px;
        }
      }
    }
  }
  
}


</style>