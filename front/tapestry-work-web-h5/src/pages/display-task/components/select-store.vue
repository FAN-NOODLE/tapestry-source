<template>
  <van-popup v-model="show" class="popup" position="bottom" get-container="body">
    <div class="popup-close" @click="handleCancel"></div>
    <div class="popup-header">
      <div class="popup-title">选择门店</div>
      <div class="popup-search">
        <van-search
          v-model="searchParams.name"
          clearable
          clear-trigger="always"
          placeholder="请输入门店编码/城市名称"
          :left-icon="searchIcon"
          @search="handleSearch"
        />
        <div v-if="searchParams.name" class="search-clear" @click="handleClear"></div>
      </div>
    </div>
    <div class="tips">仅展示已通过审核的门店</div>
    <div class="popup-content">
      <div 
        class="city-item" 
        v-for="(item, index) in list" 
        :key="index"
        :class="{opened: item.opened}">
        <div class="city-item__name" @click="handleToggleOpened(item)">{{ item.store_city }}</div>
        <div class="city-item__store">
          <div 
            class="store-item"
            :class="{active: store.store_id === selectedStoreId}"
            v-for="(store, storeIndex) in item.store_info" 
            :key="storeIndex"
            @click="handleSelect(store)">
            <span class="store-item__text">{{ store.store_name }}</span>
          </div>
        </div>
      </div>
      <empty-page v-if="list.length === 0" text="暂无门店"></empty-page>
    </div>
    <div class="popup-footer">
      <div class="btn plain" @click="handleCancel">取消</div>
      <div class="btn" @click="handleConfirm">确定</div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Popup } from 'vant';
import EmptyPage from "@/components/base/empty-page.vue";
import { getStoreListApi } from '@/api/display-task';

interface IStoreItem {
  store_id: number
  store_name: string
}

interface IStoreCityItem {
  store_city: string
  opened: boolean
  store_info: IStoreItem[]
}
@Component({
  components: {
    [Popup.name]: Popup,
    EmptyPage
  }
})
export default class SelectStore extends Vue {
  /** 是否显示 */
  show = false
  /** 搜索字段 */
  searchParams = {
    name: '',
    id: '',
    store_id: ''
  }
  /** 门店列表 */
  list: IStoreCityItem[] = []
  /** 选中的门店id */
  selectedStoreId: number = null
  /** 搜索图标 */
  searchIcon = require('@/assets/display-task/icon_search.png')

  showPopup(value: any) {
    this.show = true
    this.searchParams.id = value.id
    this.searchParams.store_id = value.store_id
    this.fetchStoreList()
  }

  hidePopup() {
    this.show = false
    this.searchParams = {
      name: '',
      id: '',
      store_id: ''
    }
    this.selectedStoreId = null
  }

  handleCancel () {
    this.hidePopup()
  }

  handleConfirm () {
    if(!this.selectedStoreId) {
      this.$toast.fail('请选择门店')
      return
    }
    this.$emit('confirm', {
      id: this.searchParams.id,
      store_id: this.selectedStoreId
    })

    // this.hidePopup()
  }

  /**
   * 展开收起
   * @param item 
   */
  handleToggleOpened (item: IStoreCityItem) {
    item.opened =!item.opened
  }

  /**
   * @description 选择门店
   * @param store 
   */
  handleSelect(store: IStoreItem) {
    this.selectedStoreId = store.store_id
  }

  /**
   * @description 搜索门店
   */
  handleSearch() {
    this.fetchStoreList()
  }

  /**
   * @description 清空搜索
   */
  handleClear() {
    this.searchParams.name = ''
    this.fetchStoreList()
  }

  /** 
   * @description 获取门店列表
   */
  async fetchStoreList() {
    const params = this.searchParams
    const { code, message, data }:IApiResult = await getStoreListApi(params)
    if (code !== 200) {
      this.$toast.fail(message)
      return 
    }
    this.list = data.map((item, index) => {
      return {
        store_city: item.store_city,
        opened: index === 0,
        store_info: item.store_info
      }
    })
  }
}
</script>
<style lang="scss" scoped>
.popup {
  &-close {
    position: absolute;
    width: 18px;
    height: 18px;
    right: 15px;
    top: 15px;
    background: url("~@/assets/display-task/icon_popup_close.png") center/100% no-repeat;
  }

  &-header {
    padding: 15px;
  }

  &-title {
    text-align: center;
    font-size: 16px;
    font-weight: 500;
    color: #000000;
    line-height: 22px;
  }

  &-search {
    margin-top: 20px;
    height: 38px;
    background: #F6F6F6;
    position: relative;
    .search-clear {
      position: absolute;
      right: 7px;
      top: 50%;
      transform: translateY(-50%);
      width: 24px;
      height: 24px;
      background: url("~@/assets/display-task/icon_search_clear.png") center/100% no-repeat;
    }
    ::v-deep .van-search {
      padding: 0;
      height: 100%;
      .van-search__content {
        height: 100%;
        padding-left: 10px;
      }
      .van-cell {
        height: 100%;
        padding: 0 15px 0 0 ;
        display: flex;
        align-items: center;
        font-size: 12px;
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
  .tips {
    font-size: 11px;
    font-weight: 400;
    color: #777777;
    line-height: 16px;
    padding: 0 20px;
  }
  &-content {
    padding: 0 15px;
    height: 50vh;
    overflow-y: scroll;

    

    .city-item {
      padding: 15px 0;

      &:not(:last-child) {
        border-bottom: 1px solid #F2F2F2;
      }

      &__name {
        font-size: 14px;
        font-weight: 500;
        color: #000000;
        line-height: 20px;
        padding-left: 22px;
        position: relative;
        &::before {
          content: "";
          transition: all 0.3s;
          width: 7px;
          height: 4px;
          position: absolute;
          left: 5px;
          top: 50%;
          transform: translateY(-50%);
          background: url("~@/assets/display-task/icon_triangle.png") center/100% no-repeat;
        }
      }


      &__store {
        padding: 0 5px;
        display: none;
        flex-wrap: wrap;

        .store-item {
          margin-top: 10px;
          margin-right: 10px;
          min-width: 84px;
          height: 24px;
          border-radius: 12px;
          border: 1px solid #A6A6A6;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 11px;
          color: #A6A6A6;
          transition: all 0.3s;
          
          &__text {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            padding:  0 10px;
          }

          &.active {
            background: #000000;
            color: #ffffff;
          }
        }
      }

      &.opened {
        .city-item__name::before {
          transform: translateY(-50%) rotate(180deg);
        }
        .city-item__store {
          display: flex;
        }
      }
    }
  }

  &-footer {
    padding: 15px;
    display: flex;
    align-items: center;
    padding-bottom: calc(15px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(15px + env(safe-area-inset-bottom)); /** iOS 11.0+ */

    .btn {
      flex: 1;
      height: 40px;
      background: #000000;
      color: #ffffff;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      &.plain {
        color: #000;
        background: #fff;
        border: 1px solid #000;
        margin-right: 15px;
      }
    }
  }

}
</style>