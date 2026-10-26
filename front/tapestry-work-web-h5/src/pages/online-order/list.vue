<template>
  <div class="container">
    <div class="header">
      <!-- office端 筛选条件 -->
      <com-filter v-if="role===5" :from="1"></com-filter>
      <div class="search">
        <van-search 
          v-model="searchParams.order_sn"
          :clearable="false"
          :left-icon="searchIcon"
          placeholder="请输入销售单号或客户名称"
          @click="handleInputClick"
          @input="handleSearch"
        ></van-search>
        <div v-if="role===3" class="date" @click="handleSelectDate">
          <div class="date-icon"></div>
          <div class="date-text">{{['全部','昨天','今天','近7天','近30天','自定义'][filterParams.date_type]}}</div>
        </div>
      </div>
      <div class="tabs-wrap">
        <van-tabs 
          v-model="tabs.active"
          title-inactive-color="#A6A6A6"
          color="#000000"
          line-width="40px"
          line-height="2px"
          :ellipsis="false"
          @change="handleTabChange">
          <van-tab v-for="(tab, index) in tabList(true)" :title="tab.title" :key="index">
            <template #title>
              <div class="tab-text">
                {{tab.title}}({{tipsNum[tab.key]||0}})
              </div>
            </template>
          </van-tab>
        </van-tabs>
        <div class="drop" :class="{active: showDrop}" @click="showDrop = !showDrop"></div>
        <div v-if="showDrop" class="drop-menu">
          <div class="drop-menu-item" v-for="(tab, index) in tabList(false)" :key="index" @click="handleSelectTab(tab)">{{ tab.title }}({{tipsNum[tab.key]||0}})</div>
        </div>
      </div>
    </div>
    <div class="content">
      <van-list
        v-model="loading"
        :finished="finished"
        :finished-text="finishedText"
        :immediate-check="false"
        @load="fetchData"
      >
        <div v-if="tabList(true)[tabs.active].value === 1" class="checkbox">请提醒客人务必在今日 23 时 59 分 59 秒之前完成收件地址的填写，若逾期未完成，订单将自动作废。</div>
        <order-item v-for="(item, index) in list" :key="index" :item="item"></order-item>
      </van-list>
      <van-empty
        v-if="list.length == 0 && !loading"
        image-size="66"
        :image="emptyIcon"
        :description="emptyText"
      >
      </van-empty>
    </div>
    <div class="footer">
      <div class="btn plain" @click="showExportPopup=true">导出当前列表Excel</div>

    </div>
    <!-- 下载弹框 -->
    <com-confirm v-model="showExportPopup" @confirm="handleExportConfirm"></com-confirm>
    <!-- 提交成功 -->
    <order-success-popup v-model="showSuccessPopup" :show-footer="false" :content="['下载内容稍后将发送到指定邮箱']"></order-success-popup>
  </div>
</template>
<script lang="ts">
import { exportEmailApi, getOrderListApi, getOrderNumApi } from '@/api/online-order';
import { Component, Vue } from 'vue-property-decorator';
import OrderItem, { IOrderItem } from './components/order-item.vue';
import ComFilter from './components/com-filter.vue';
import ComConfirm from '../communication/components/com-confirm.vue';
import OrderSuccessPopup from './components/order-success-popup.vue';
import { getUserAuthApi } from '@/api/customer-portrait';
import emptyPage from "@/components/base/empty-page.vue";
import wx from 'weixin-js-sdk'
import dayjs from 'dayjs';
@Component({
  components: {
    emptyPage,
    OrderSuccessPopup,
    ComConfirm,
    ComFilter,
    OrderItem
  }
})
export default class OnlineOrderList extends Vue {
  role = 3
  searchIcon = require('@/assets/material-library/search.png')
  /** 查询方式切换 */
  tabs = {
    /** 下标  */
    active: 0,
    /** 列表 */
    list: [
      {
        title: '全部订单',
        value: 0,
        key: 'all_num',
        show: true
      },
      {
        title: '待填写',
        value: 1,
        key: 'wait_num',
        show: true
      },
      {
        title: '待发货',
        value: 2,
        key: 'pending_num',
        show: true
      },
      {
        title: '已发货',
        value: 3,
        key: 'send_num',
        show: true
      },
      {
        title: '已拒单',
        value: 4,
        key: 'refuse_num',
        show: false
      },
      {
        title: '已作废',
        value: 5,
        key: 'cancel_num',
        show: false
      },
    ]
  }
  /** 是否正在加载数据  */
  loading = false;
  /** 是否没有更多数据 */
  finished = false;
  /** 搜索参数 */
  searchParams = {
    order_sn: ''
  }
  /** 分页 */
  meta: IMeta = {
    page: 1,
    size: 10,
  }
  /** 列表数据 */
  list: IOrderItem[] = []
  /** 导出弹框 */
  showExportPopup = false
  /** 导出成功弹框 */
  showSuccessPopup = false

  tipsNum = {
    all_num: 0,
    refuse_num: 0,
    send_num: 0,
    wait_num: 0,
  }

  showDrop = false

  get finishedText() {
    return this.list.length > 0 ? '～～没有更多了～～':''
  }

  /** 空页面提示 */
  get emptyText() {
    if(this.searchParams.order_sn) {
      return '未找到相关订单'
    }
    return `暂无订单信息`
  }
  /** 空页面图标 */
  get emptyIcon() {
    return require('@/assets/online-order/empty.png')
  }

  get filterParams() {
    return this.$store.getters["onlineOrder/getFilterParams"];
  }

  get isUpdate() {
    return this.$store.getters["onlineOrder/getUpdate"];
  }

  get tabList() {
    return (type: boolean) => {
      return this.tabs.list.filter(item => item.show === type)
    }
  }

  async mounted() {
    await this.getUserAuth()

    const index = this.tabs.list.findIndex(item => item.value ===  Number(this.$route.query.status || 0))
    if(index> 3) {
      this.tabs.list[3].show = false
      this.tabs.list[index].show = true
      this.tabs.active = 3
    } else {
      this.tabs.active = index
    }
    console.log(11111, this.tabs);
    
    this.reload()
    this.getOrderNum()
    this.$store.commit("onlineOrder/setUpdate", false)
    if(this.$route.query.from != '1') {
      // 神策数据上报 -- 进入网仓
      this.$track.trackEnterOnlineWarehouse({
        OnlineWarehouse_entrance: '企业微信消息',
        message_type: Number(this.$route.query.status) == 1 ? '您有客人未提交收货信息订单' :'订单状态更新',
      })
    }

  }
  
  async activated() {
    if(this.isUpdate) {
      this.reload()
      this.getOrderNum()
    }
    this.$store.commit("onlineOrder/setUpdate", false)
  }

  popstate() {
    wx.closeWindow(); // 微信网页退出
  }

  /** 
   * @description 获取列表数据
   */
  async fetchData() {
    try {
      const { start_time, end_time, member_ids, third_ids, leave_member_ids } = this.filterParams
      const params = {
        page: this.meta.page,
        size: this.meta.size,
        status: this.tabList(true)[this.tabs.active].value,
        start_time,
        end_time,
        member_ids,
        third_ids,
        leave_member_ids,
        ...this.searchParams,
      }
      const {code, message, data}: IApiResult = await getOrderListApi(params)
      if (code !== 200) {
        this.finished = true;
        this.loading = false;
        this.$toast.fail(message);
        return;
      }
      this.loading = false;
      this.finished = data.list.length === 0;
      this.list = this.meta.page == 1 ? data.list : this.list.concat(data.list)
      this.meta.page = this.meta.page + 1

    } catch (error) {
      console.log(error);
      this.finished = true;
      this.loading = false;
      this.list = [{status:1}, {status:2}, {status:3}]
    }
  }

  async getOrderNum() {
    const { start_time, end_time, member_ids,third_ids, leave_member_ids } = this.filterParams
    const res = await getOrderNumApi({ 
      start_time, 
      end_time,
      member_ids,
      third_ids,
      leave_member_ids,
      ...this.searchParams,
    })
    this.tipsNum = res.data
  }

  async getUserAuth() {
    const res = await getUserAuthApi({})
    this.role = res.data.role
    this.$store.commit('user/setRole', res.data.role)
    if(res.data.role == 5) {
      this.tabs.list = [
        {
          title: '全部订单',
          value: 0,
          key: 'all_num',
          show: true
        },
        {
          title: '待发货',
          value: 2,
          key: 'pending_num',
          show: true
        },
        {
          title: '已发货',
          value: 3,
          key: 'send_num',
          show: true
        },
        {
          title: '已拒单',
          value: 4,
          key: 'refuse_num',
          show: true
        },
        {
          title: '待填写',
          value: 1,
          key: 'wait_num',
          show: false
        },
        {
          title: '已作废',
          value: 5,
          key: 'cancel_num',
          show: false
        },
      ]
    }
  }

  /**
   * @description 刷新
   */
  reload() {
    if(this.loading) return
    this.loading = true;
    this.showDrop = false
    this.meta.page = 1;
    this.finished = false
    this.list = [];
    this.fetchData()
  }

  /**
   * @description 切换tab
   */
  handleTabChange(index) {
    // 神策数据上报 -- 订单筛选
    this.$track.trackOrderFilter({
      current_page: '订单列表',
      filter_type: '订单状态',
      // filter_option: [...member_ids, ...leave_member_ids].join()
      filter_option: this.tabList(true)[this.tabs.active].title
    })
    this.reload()
  }

  /**
   * @description 搜索
   */
  handleSearch() {
    this.reload();
    this.getOrderNum()
  }

  /**
   * @description 新建
   */
  handleCreate() {
    this.$router.push({
      path: 'online-order-create'
    })
  }

  /**
   * @description 导出
   * @param e 
   */
  async handleExportConfirm(e) {
    this.$toast.loading()
    this.showExportPopup = false
    const { start_time, end_time, member_ids, third_ids, third_name_list, member_name_list,leave_member_ids } = this.filterParams
    const params = {
      page: 1,
      size: 10,
      status: this.tabList(true)[this.tabs.active].value,
      start_time,
      end_time,
      member_ids,
      third_ids,
      email: e,
      ...this.searchParams,
    }
    const res:any =await exportEmailApi(params)
    this.$toast.clear()
    if(res.code == 200){
      this.showSuccessPopup = true
      
      this.$track.trackExportOrderList({
        period: [dayjs(start_time*1000).format('YYYY-MM-DD'),dayjs(end_time*1000).format('YYYY-MM-DD')].join(' - '),
        store_name: third_name_list.join(),
        sales_name: member_name_list.join(),
        // sales_name: [...member_ids, ...leave_member_ids].join(),
        original_words: this.searchParams.order_sn,
        order_status: this.tabList(true)[this.tabs.active].title,
      })
    } else {
      this.$toast('下载失败，请重新下载')
    }
  }

  handleSelectDate() {
    this.$router.push({
      path: 'online-order-select-date',
    })
  }

  handleSelectTab(item: any) {
    // const index = this.tabs.list.findIndex(item => item.value === tab.value)
    this.tabList(true)[3].show = false
    item.show = true
    this.tabs.active = 3
    this.reload();
    // 神策数据上报 -- 订单筛选
    this.$track.trackOrderFilter({
      current_page: '订单列表',
      filter_type: '订单状态',
      filter_option: this.tabList(true)[this.tabs.active].title
    })
  }

  handleInputClick() {
    // 神策数据上报 -- 搜索框点击
    this.$track.trackSearchClick({
      current_page: '订单列表'
    })
  }
  
}
</script>
<style lang="scss" scoped>
.btn {
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000000;
  color: #fff;
  font-size: 14px;

  &.plain {
    background: #fff;
    color: #000;
    border: 1px solid #000;
  }
}
::v-deep {
  .van-empty {
    padding-top: 100px;
  }
  .van-tabs__line {
    bottom: 22px
  }

  .van-search {
    padding: 0;
    height: 38px;
    .van-search__content {
      height: 100%;
      padding-left: 10px;
      background-color: #F6F6F6;
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

  .checkbox {
      position: relative;
      font-weight: 400;
      font-size: 12px;
      color: #676767;
      line-height: 17px;
      // padding-left: 22px;
      margin-bottom: 10px;
    }
}
.container {
  width: 100vw;
  min-height: 100vh;
  background: #F7F7F7;

  padding-bottom: calc(64px + contain(safe-area-inset-bottom));
  padding-bottom: calc(64px + env(safe-area-inset-bottom));

  .header {
    background-color: #fff;
    position: sticky;
    top: 0;
    z-index: 1;

    .search {
      padding: 10px 0 10px 15px;
      display: flex;

      .van-search {
        flex: 1;
        margin-right: 15px;
      }

      .date {
        width: 62px;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        margin-left: -15px;

        .date-icon {
          width: 20px;
          height: 20px;
          background: url('~@/assets/online-order/icon_date.png') center/100% no-repeat;
        }

        .date-text {
          font-size: 13px;
          color: #000000;
        }
      }
    }

    .com-filter~.search {
      padding-top: 0;
    }

    .tab-text {
      position: relative;
      .num-icon-box {
        position: absolute;
        top: 0;
        right: 2px;
        .num-text {
          border-radius: 6px;
          min-width: 12px;
          height: 12px;
          padding: 0 2px;
          position: absolute;
          top: 50%;
          left: 0;
          transform: translateY(-50%);
          background: #000000;
          font-size: 8px;
          font-weight: 3400;
          color: #ffffff;
          display: flex;
          justify-content: center;
          align-items: center;
        }
      }
    } 
    
    .tabs-wrap {
      position: relative;

      padding-right: 30px;


      .drop {
        position: absolute;
        height: 44px;
        width: 30px;
        right: 5px;
        top: 0;
        background: #fff url('~@/assets/online-order/icon_triangle.png') center/10px no-repeat;
        transition: all 0.3;

        &.active {
          transform: rotate(180deg);
        }
      }
      .drop-menu {
        position: absolute;
        right: 5px;
        top: 54px;
        width: 150px;
        background: #FFFFFF;
        box-shadow: 0px 1px 5px 0px rgba(0,0,0,0.2);

        &::after {
          content: "";
          width: 10px;
          height: 10px;
          box-shadow: 0px 1px 5px 0px rgba(0,0,0,0.2);
          position: absolute;
          top: 0;
          left: 70%;
          transform: translate(-50%, -50%) rotate(45deg);
          background-color: #fff;
        }
      }
      .drop-menu-item {
        height: 50px;
        line-height: 50px;
        text-align: center;
        font-weight: 400;
        font-size: 14px;
        color: #A6A6A6;
        background-color: #fff;
        position: relative;
        z-index: 1;
      }
    }
  }

  .content {
    padding: 15px;
    .btn {
      width: 165px;
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

    .btn {
      flex: 1;
      &:not(:first-child) {
        margin-left: 15px
      }
    }
  }
}

</style>