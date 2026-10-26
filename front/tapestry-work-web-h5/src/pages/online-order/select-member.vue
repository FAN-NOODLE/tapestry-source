<template>
  <div class="container">
    <div class="header">
      <van-search 
        v-model="searchParams.keyword"
        :clearable="false"
        :left-icon="searchIcon"
        placeholder="搜索导购姓名、手机号"
        @click="handleInputClick"
        @input="handleSearch"
      ></van-search>
    </div>
    <div class="content">
      <div class="item all">
        <div class="item-title">
          <div class="flex-1">全选导购</div>
          <div class="checkbox" :class="`checked-${checkedAll}`" @click="handleCheckAll"></div>
        </div>
      </div>
      <div class="item" v-for="item, index in dataList" :key="index">
        <div class="item-title" @click="handleExpand(item)">
          <div class="triangle" :class="{active: !!item.expand}" v-if="item.third_id && item.member_list.length>0"></div>
          <div v-if="item.third_name === '离职导购'" class="flex-1">{{item.third_name}}</div>
          <div v-else class="flex-1">[{{ item.store_code }}] {{item.store_name||item.store_code}}</div>
          <div class="checkbox" :class="`checked-${item.checked}`" @click.stop="handleCheckStore(item)"></div>
        </div>
        <div v-if="item.third_id" class="item-content" :style="{height: !!item.expand ? item.member_list.length*34+20 + 'px' : '0px'}">
          <div class="item-title" v-for="(member, memberIndex) in item.member_list" :key="memberIndex">
            <div class="flex-1">
              <span>{{member.member_name}}({{ member.userid }})</span>
            </div>
            <div class="checkbox" :class="`checked-${member.checked}`" @click="handleCheckMember(member, item)"></div>
          </div>
        </div>
      </div>
    </div>
    <div class="footer">
      <van-button plain color="#000000" @click="handleReset">重置</van-button>
      <van-button color="#000000" @click="handleConfirm">确认</van-button>
    </div>
  </div>
</template>
<script lang="ts">
import { getGuideListApi, getLeaveGuideListApi } from '@/api/online-order';
import { Component, Vue } from 'vue-property-decorator';

@Component({

})
export default class OnlineOrderSelectStore extends Vue {
  searchIcon = require('@/assets/material-library/search.png')
  /** 搜索参数 */
  searchParams = {
    keyword: ''
  }

  // 全选状态
  checkedAll = 0

  guideList = []
  dataList:any[] = []
  filterParams:any = {}

  async mounted() {
    this.filterParams = JSON.parse(JSON.stringify(this.$store.getters["onlineOrder/getFilterParams"]))
    await this.getGuideList()
  }

  async getGuideList() {
    const { third_ids } = this.filterParams
    const res:any = await getGuideListApi({
      check_ids: third_ids
    })
    if(res.code == 200) {
      this.guideList = res.data.map(item => {
        return {
          checked: 0,
          ...item
        }
      })
      this.computedList(this.guideList)
    }
  }
  async getLeaveGuideList() {
    const { third_ids } = this.filterParams
    const res:any = await getLeaveGuideListApi({
      check_ids: third_ids
    })
    if(res.code == 200) {
      return res.data.map(item => {
        return {
          member_id: item
        }
      })
    } else {
      return []
    }
  }

  async computedList(list:any[]) {
    const { member_ids, leave_member_ids } = this.filterParams
    
    // 在职员工
    list.forEach(item => {
      const index = this.dataList.findIndex(val => val.third_id == item.third_id)
      item.checked = member_ids.includes(item.member_id) ? 1 : 0
      if(index == -1) {
        this.dataList.push({
          checked: 0,
          expand: !!this.searchParams.keyword ? 1: 0,
          third_id: item.third_id,
          third_name: item.third_name,
          store_name: item.store_name,
          store_code: item.store_code,
          member_list: [item]
        })
      } else {
        this.dataList[index].member_list.push(item)
      }
    })

    const leave = await this.getLeaveGuideList()
    
    this.guideList.push(...leave)
    
    if(leave.length) {
      this.dataList.push({
        checked: 0,
        expand: 0,
        third_id: '',
        third_name: '离职导购',
        member_list: leave.map(item => {
          item.checked = leave_member_ids.includes(item.member_id) ? 1 : 0
          return item
        })
      })
    }

    this.dataList = this.dataList.map(item => {
      if(item.member_list.every(member => member.checked)) {
        item.checked = 1
      } else if(item.member_list.every(member => !member.checked)) {
        item.checked = 0
      } else {
        item.checked = 2
      }
      
      return item
    })
    this.isCheckedAll()
  }

  handleSearch() {
    if(this.searchParams.keyword) {
      this.dataList = []
      this.computedList(this.guideList.filter(item => (item.member_name && item.member_name.indexOf(this.searchParams.keyword) !==-1)||item.mobile === this.searchParams.keyword))
    } else {
      this.dataList = []
      this.computedList(this.guideList.filter(item => item.member_name))
    }
  }

  handleCheckAll() {
    this.checkedAll = [1, 0][this.checkedAll] ?? 1
    this.dataList.forEach(store => {
      store.checked = this.checkedAll
      store.member_list.forEach(member => {
        member.checked = this.checkedAll
      })
    })
  }

  isCheckedAll() {
    if(this.dataList.every(item => item.checked === 1)) {
      this.checkedAll = 1
    } else if(this.dataList.every(item => item.checked === 0)) {
      this.checkedAll = 0
    } else {
      this.checkedAll = 2
    }
  }

  /**
   * 选择门店
   * @param item 
   */
  handleCheckStore(item:any) {
    item.checked = [1, 0][item.checked] ?? 1
    item.member_list.forEach(member => {
      member.checked = item.checked
    })
    console.log(item.member_list);
    
    this.isCheckedAll()
    this.getSelectResult()
  }

  /**
   * 选择导购
   */
  handleCheckMember(item, parent) {
    item.checked = [1, 0][item.checked] ?? 1

    if(parent.member_list.every(member => member.checked)) {
      parent.checked = 1
    } else if(parent.member_list.every(member => !member.checked)) {
      parent.checked = 0
    } else {
      parent.checked = 2
    }
    
    this.isCheckedAll()
    this.getSelectResult()
  }
  /**
   * 展开
   */
  handleExpand(item: any) {
    item.expand = [1, 0][item.expand]
  }


  handleReset() {
    this.$store.commit("onlineOrder/setFilterParams", {
      ...this.filterParams,
      member_ids: [],
      leave_member_ids: []
    });
    this.$store.commit("onlineOrder/setUpdate", true)
    this.$router.back()
  }

  getSelectResult() {
    const storeList = this.guideList.filter(item => item.third_id && item.checked === 1)
    this.filterParams.leave_member_ids = this.guideList.filter(item => !item.third_id && item.checked === 1).map(item => item.member_id)
    this.filterParams.member_ids = Array.from(new Set(storeList.map(item => item.member_id)))
    this.filterParams.member_name_list = Array.from(new Set(storeList.map(item => item.member_name)))
  }

  handleConfirm() {
    this.getSelectResult()
    this.$store.commit("onlineOrder/setFilterParams", {
      ...this.filterParams,
    });
    this.$store.commit("onlineOrder/setUpdate", true)
    const { from } = this.$route.query
    const { member_name_list = [], member_ids, leave_member_ids } = this.filterParams
    if(leave_member_ids.length) {
      member_name_list.push('离职导购')
    }
    // 神策数据上报 -- 订单筛选
    this.$track.trackOrderFilter({
      current_page: ['订单数据', '订单列表'][Number(from)||0],
      filter_type: '导购',
      // filter_option: [...member_ids, ...leave_member_ids].join()
      filter_option: member_name_list.join()
    })

    this.$store.commit("onlineOrder/setUpdate", true)
    this.$router.back()
  }


  handleInputClick() {
    // 神策数据上报 -- 搜索框点击
    this.$track.trackSearchClick({
      current_page: '选择导购'
    })
  }

}
</script>
<style lang="scss" scoped>
::v-deep {
  .van-search {
    padding: 0;
    height: 38px;
    .van-search__content {
      height: 100%;
      padding-left: 10px;
      background-color: #F6F7F9;
    }
    .van-cell {
      height: 100%;
      padding: 0 15px 0 0 ;
      display: flex;
      align-items: center;
      color: red;
    }
    .van-cell__value {
      height: 100%;
    }
    .van-field__body {
      height: 100%;
    
    }
    .van-field__left-icon {
      display: flex;
      align-items: center;
      justify-items: center;
      margin-right: 7px;
    }
  }
}
.checkbox {
  width: 48px;
  height: 48px;
  background: url('~@/assets/online-order/icon_c_1.png') center/16px no-repeat;
  transition: all 0.3s;
  &.checked-1 {
    background: url('~@/assets/online-order/icon_c_2.png') center/16px no-repeat;
  }
  &.checked-2 {
    background: url('~@/assets/online-order/icon_c_3.png') center/16px no-repeat;
  }
}
.container {
  width: 100vw;
  min-height: 100vh;
  background-color: #ffffff;
  padding-bottom: calc(64px + contain(safe-area-inset-bottom));
  padding-bottom: calc(64px + env(safe-area-inset-bottom));

  .header {
    padding: 10px 15px;
    background-color: #fff;
    box-shadow: 0px 6px 12px 0px rgba(222,224,227,0.5);
  }

  .content {
    padding: 5px 15px;

    .item {
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      border-bottom: 1px solid #EEEEEE;
      overflow: hidden;
      .item-title {
        display: flex;
        align-items: center;
        padding-left: 11px;
        position: relative;
        background-color: #fff;
        z-index: 1;

        .triangle {
          width: 7px;
          height: 4px;
          background: url('~@/assets/online-order/icon_triangle.png') center/100% no-repeat;
          margin-right: 6px;
          transition: all .3s;
          position: absolute;
          left: 0;
          top: 50%;
          transform: translateY(-50%) rotate(-90deg);
          &.active {
            transform: translateY(-50%) rotate(0deg);
          }
        }

        .flex-1 {
          flex: 1;
        }

        .checkbox {
          margin-right: -16px;
        }

        .tag {
          margin-left: 6px;
          padding: 0 4px;
          line-height: 17px;
          background: #F0F0F0;
        }
      }

      &.all .item-title {
        padding-left: 0;
      }

      .item-content {
        padding-left: 5px;
        height: 0;
        overflow: hidden;
        transition: all 0.3s ease-in-out;
        &.active {
          height: auto;
        }
        .item-title {
          font-weight: 400;
          font-size: 12px;
          color: #000000;
          line-height: 17px;
          z-index: 0;
          
          .flex-1 {
            display: flex;
            align-items: center;
          }
        }

        .checkbox {
          width: 34px;
          height: 34px;
          margin-right: -9px;
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