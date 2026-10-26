<template>
  <div class="container">
    <div v-if="type === 0" class="header">
      <div class="task-name">
        <span class="task-name__text" 
          :style="{
            paddingRight: searchParams.task_num>1?'60px':'0'
          }"
          >
          <div class="ellipsis">{{searchParams.task_name}}</div>
          <div v-if="searchParams.task_num>1" class="task-num">等{{searchParams.task_num}}个任务</div>
        </span>
      </div>
      <div class="store-name">
        <span v-if="searchParams.is_total">全部区域/全部门店</span>
        <span v-else>{{searchParams.area_num}}个区域/{{searchParams.store_num}}家门店</span>
      </div>
    </div>
    <div class="title" v-else> 
      <span>{{typeLabel(type)}}：</span>
      <span>{{typeName}}</span>
    </div>
    <!-- 区域门店 -->
    <template v-if="type < 4">
      <!-- 搜索 -->
      <div class="search">
        <van-search v-model="keyword" :placeholder="'请输入'+(typeLabel(currentType))+'快捷搜索'" :left-icon="searchIcon" :clearable="false" />
        <div v-if="keyword" class="search-clear" @click="handleClear"></div>
      </div>
      <!-- 列表 -->
      <div class="table-wrap">
        <com-table 
          :columns="columns" 
          :data="renderList"
          @row-click="handleRowClick"
          >
          <template v-slot:operation="{row}">
            <div v-if="row.store_num>0 || row.task_num>0" class="icon-right"></div>
          </template>
          <template v-slot:empty>
            <empty-page
              picture="noSearchData"
              paddingTop="60"
              paddingBottom="60"
              text="暂无数据"
            ></empty-page>
          </template>
        </com-table>
      </div>
    </template>
    <!-- 任务 -->
    <template v-else>
      <div v-if="list.length" class="task-list">
        <div 
          class="task-item" 
          v-for="(item, index) in list" 
          :key="index" 
          @click="handleViewDetail(item)">
          <display-task-item 
            :item="item" 
            :status="searchParams.finish_status"
            disabled>
          </display-task-item>
        </div>
        <div class="load-more">～没有更多了～</div>
      </div>
      <empty-page v-else></empty-page>
    </template>
    <div class="footer">
      <div class="btn" @click="handleBack(back)">返回首页</div>
      <div class="btn" v-if="type !== 0" @click="handleBack(-1)">返回上一级</div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import EmptyPage from "@/components/base/empty-page.vue";
import ComTable from './components/com-table.vue';
import { getDashboardListApi, getStoreTaskListApi } from '@/api/display-task'
import { statusName } from './const'
import DisplayTaskItem from "./components/display-task-item.vue";

interface IStoreItem {
  /** 节点ID */
  id: number;
  /** 区域名称 */
  name: string;
  /** 节点层及 */
  type: number;
  /** 门店数据 */
  store_num: number;
  /** 任务数 */
  task_num: number;
  /** 子ID */
  children_third_ids: number[];
}

@Component({
  name: 'DisplayStoreList',
  components: {
    EmptyPage,
    ComTable,
    DisplayTaskItem
  }
})
export default class DisplayStoreList extends Vue {
  /** 返回层级 */
  back = -1
  /** 搜索关键字 */
  keyword = '';
  /** 搜索图标 */
  searchIcon = require('@/assets/display-task/icon_search.png')
  /** 数据列表 */
  list: IStoreItem[] = []
  /** 穿透类型 */
  type = 0
  typeName = '';

  currentType = 0

  get typeLabel() {
    return value => ['', '渠道名称', '区域名称', '区域名称', '门店名称', '门店名称'][value]
  }

  /** 列数据 */
  get columns() {
    const name = statusName[this.searchParams.finish_status]
    const type = this.list[0]?.type || 0
    return [
      {
        label: this.typeLabel(type),
        prop: 'name',
      },
      type === 4 ? {
        label: name.replace('门店','任务'),
        prop: 'task_num',
        width: 150,
        sortable: true,
        sortType: '',
        align: 'center'
      }: {
        label: name,
        prop: 'store_num',
        width: 150,
        sortable: true,
        sortType: '',
        align: 'center'
      },
      {
        label: '',
        prop: 'operation',
        width: 40
      },
    ]
  }

  get searchParams() {
    const params = {
      finish_status: Number(this.$route.query.finish_status),
      ...this.$store.getters["displayTask/getSearchParams"]
    }
    if(this.$route.query.id) {
      params.target_area_id = Number(this.$route.query.id)
    }
    return params
  }

  get renderList() {
    return this.list.filter((item:any) => {
      return item.name.indexOf(this.keyword) > -1
    })
  }

  @Watch('$route.query')
  onRouteChange() {
    this.keyword = ''
    this.list = []
    this.init()
  }

  mounted() {
    this.init()
  }


  async init() {
    const { type, name, back} = this.$route.query
    this.back = Number(back) || -1
    this.type = Number(type) || 0;
    this.typeName = name as string
    if(this.type === 4) {
      await this.getStoreTaskList()
    } else {
      await this.getList()
    }
  }

  async getList() {
    const params = {
      ...this.searchParams
    }
    const res: any = await getDashboardListApi(params)
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    this.list = res.data
    const total = res.data.map(item => item.store_num).reduce((prev, cur) => prev + cur)
    document.title = statusName[params.finish_status] + ': ' + total
    console.log(res.data[0].type);
    
    if(res.data[0]) {
      this.currentType = res.data[0].type
    }
  }

  async getStoreTaskList(id?: number) {

    const params = {
      ids: this.searchParams.ids,
      store_id: id || this.searchParams.target_area_id,
      finish_status: this.searchParams.finish_status,
    }
    const res: any = await getStoreTaskListApi(params)
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    if(res.data.length === 1) {
      this.handleViewDetail(res.data[0])
    } else {
      this.list = res.data
      const total = res.data.length
      document.title = statusName[params.finish_status].replace('门店', '任务') + ': ' + total
    }
  }

  handleClear() {
    this.keyword = ''
  }

  handleRowClick(row: IStoreItem) {
    const id = row.type === 4 ? row.children_third_ids[0]: row.id
    // 如果下钻的门店数量或者任务数量都是0 则不能点击
    if(row.store_num===0 && row.task_num===0) return
    // 门店到任务，如果只有一个任务就直接到任务详情
    if(row.type === 4 && row.task_num === 1) {
      this.getStoreTaskList(id)
    } else {
      this.$router.push({
        name: 'display-store-list',
        query: {
          id: id + '',
          name: row.name,
          type: row.type + '',
          back: this.back - 1 + '',
          finish_status: this.searchParams.finish_status + ''
        }
      })
    }
  }

  handleViewDetail(item: any) {
    this.$store.commit('displayTask/setDetail', {
      task_info: null,
      store_info: null
    });
    this.$store.commit('displayTask/setUpdate', true);
    this.$router.push({
      path: 'display-task-detail',
      query: {
        id: item.id + '',
        store_id: item.store_id + '',
        from: '1',
        manager: '1' // 从数据管理进入
      }
    });
  }

  handleBack(type: number) {
    this.$router.go(type)
  }
}
</script>
<style lang="scss" scoped>
  .container {
    min-height: 100vh;
    background-color: #f7f7f7;
    font-size: 12px;
    padding-bottom: calc(56px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(56px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    
    .header {
      padding: 10px 15px;
      background: #FFFFFF;
      box-shadow: 0px 6px 12px 0px rgba(222,224,227,0.5);
      display: flex;
      align-items: center;
      .task-name {
        flex: 1;
        overflow: hidden;
        height: 38px;
        background: #F6F6F6;
        position: relative;
        padding: 10px 15px;
        display: flex;
        align-items: center;

        &__text {
          overflow: hidden;
          position: relative;

          .ellipsis {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .task-num {
            position: absolute;
            right: 0;
            top: 50%;
            transform: translateY(-50%);
          }
        }
      }

      .store-name {
        padding-left: 10px;
      }
    }

    .title {
      padding: 15px 20px 0;
      font-size: 14px;
      font-weight: 500;
      color: #000000;
      line-height: 20px;
    }

    .search {
      padding: 15px;
      margin-top: 3px;
      position: relative;
      .van-search {
        height: 38px;
        padding: 0;
        border: 1px solid #A6A6A6;
        background-color: transparent;
      }

      .van-search__content {
        height: 100%;
        background-color: transparent;
      }
      .van-cell {
        height: 100%;
        padding: 0 15px 0 0 ;
        display: flex;
        align-items: center;
        font-size: 12px;
      }
      ::v-deep {
        .van-cell__value {
          height: 100%;
        }
        .van-field__body {
          height: 100%;
        }
        input.van-field__control {
          height: 100%;
        }
      }
      .van-field__left-icon {
        display: flex;
        align-items: center;
        justify-items: center;
        margin-right: 7px;
      }
      .search-clear {
        position: absolute;
        right: 20px;
        top: 50%;
        transform: translateY(-50%);
        width: 20px;
        height: 20px;
        background: url("~@/assets/display-task/icon_search_clear_2.png") center/100% no-repeat;
      }
    }

    .table-wrap {
      margin-top: -15px;
      padding: 0 15px;
      ::v-deep {
        .table {
          background-color: #fff;
        }
        .table__header {
          background-color: #f7f7f7;

          .table__col {
            padding: 15px 0;
            color: #A6A6A6;
            font-weight: normal;
            text-align: center;
          }
        }


        .table__body {
          .table__col {
            padding: 15px 0;
          }
          .table__row {
            &:not(:last-child) {
              border-bottom: 1px solid #F2F2F2;
            }
          }
        }

        // .empty {
        //   background-color: #f7f7f7;
        // }
        .default {
          height: auto;
          // background-color: #f7f7f7;
        }
      } 
    }

    .task-list {
      padding: 15px;

      .task-item {
        &:not(:last-child) {
          margin-bottom: 15px;
        }
      }
    }
  }
  .icon-right {
    width: 11px;
    height: 17px;
    background: url('~@/assets/display-task/icon_right_black.png') center/100% no-repeat;
  }
  .load-more {
    text-align: center;
    color: #A6A6A6;
  }

  .footer {
    display: flex;
    align-items: center;
    justify-content: space-around;
    position: fixed;
    bottom: 0;
    width: 100vw;
    background-color: #fff;
    padding: 10px 20px;
    padding-bottom: calc(10px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(10px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    z-index: 100;
    .btn {
      flex: 1;
      height: 36px;
      border: 1px solid #000000;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;

      &:nth-child(2) {
        margin-left: 14px;
      }
    }
  }
</style>