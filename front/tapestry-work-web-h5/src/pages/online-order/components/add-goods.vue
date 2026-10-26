<template>
  <div id="search">
    <van-tabs 
      v-model="activeTab"
      title-inactive-color="#A6A6A6"
      color="#000000"
      line-width="40px"
      line-height="2px"
      @change="handleTabChange"
    >
      <van-tab title="SKU CODE添加">
        <div class="box">
          <div class="input-wrap">
            <van-field 
              v-model="keyword[0]"
              ref="input"
              class="input"
              clearable
              clear-trigger="always"
              placeholder="请输入商品款_色_码，例如CW971_OPI_5B"
              @input="handleInput"
              @focus="handleFocus"
              @click="handleClick" />
              <!-- <div v-if="keyword[0]" class="input-clear" @click="handleClear"></div> -->
          </div>
          <div class="tips">所有商品将走<{{`OC226网仓`}}>查询库存并发货</div>
          <!-- 搜索结果 -->
          <div v-if="showSearchResult" class="search-result">
            <!-- 有结果 -->
            <div class="search-item" v-for="(item, index) in searchResult" :key="index" @click="handleSelectItem(item)">{{ item.product_code }}</div>
            <!-- 无结果 -->
            <div v-if="!searchResult.length" class="search-empty">
              <div class="search-empty-text">商品库未查询到对应的商品信息</div> 
              <div class="search-empty-btn" @click="handleAddPopup">手动添加</div>
            </div>
          </div>
        </div>
      </van-tab>
      <van-tab title="条形码添加">
        <div class="box">
          <van-field
            v-model="keyword[1]"
            class="input-2"
            ref="fieldInput"
            clearable
            clear-trigger="always"
            placeholder="请输入条形码查询相关数据"
            :right-icon="require('@/assets/online-order/icon_scan.png')"
            @input="handleInput2"
            @click="handleClick2"
            @click-right-icon="handleScan"
          />
          <div class="tips">所有商品将走<{{`OC226网仓`}}>查询库存并发货</div>
          <!-- 搜索结果 -->
          <div v-if="showSearchResult" class="search-result">
            <!-- 有结果 -->
            <div class="search-item" v-for="(item, index) in searchResult" :key="index" @click="handleSelectItem(item)">{{ item.product_code }}</div>
            <!-- 无结果 -->
            <div v-if="!searchResult.length" class="search-empty">
              <div class="search-empty-text">商品库未查询到对应的商品信息</div> 
              <div class="search-empty-btn" @click="handleAddPopup">手动添加</div>
            </div>
          </div>
        </div>
      </van-tab>
    </van-tabs>
    <div v-if="innerValue.length" class="list">
      <div class="tips">已添加({{ innerValue.length }})</div>
      <div class="item" v-for="(item, index) in innerValue" :key="index">
        <van-image 
          width="80px"
          height="80px"
          fit="contain"
          lazy-load
          :src="item.image||emptyImage"
        ></van-image>
		<div class="item-info">
          <div class="item-code">{{item.product_code}}</div>
          <div class="item-name">{{item.product_name}}</div>
          <div class="item-stock" :class="{error: item.quantity > item.dcl_on_hand}">
            <!-- 加载中 -->
            <div v-if="item.loading" class="flex">
              <span>网仓库存</span>
              <van-loading type="spinner" size="12" />
            </div>
            <!-- 无库存 -->
            <span v-else-if="item.dcl_on_hand === -1">商品库存未知</span>
            <span v-else-if="item.dcl_on_hand === '--'">网仓库存 --</span>
            <!-- 无库存 -->
            <span v-else-if="item.dcl_on_hand === 0">网仓无库存</span>
            <!-- 有库存 -->
            <span v-else>网仓库存 {{ item.dcl_on_hand }}</span>
          </div>
		  <div v-if="item.inventory_warning" class="inventory-warning">{{ item.inventory_warning }}</div>
        </div>
        <div class="item-delete" @click="handleDelete(index)"></div>
        <div class="item-opt">
          <!-- 减 -->
          <van-button class="item-minus" color="#000000" plain @click="handleNumChange(index, -1)">-</van-button>
          <!-- 数量 -->
          <div class="item-num">{{item.quantity}}</div>
          <!-- 加 -->
          <van-button class="item-plus" color="#000000" @click="handleNumChange(index, 1)">+</van-button>
        </div>
      </div>
    </div>
    <!-- 删除确认 -->
    <order-delete-popup v-model="showDeletePopup" :item="innerValue[current]" @confirm="handleDeleteConfirm"></order-delete-popup>
    <!-- 手动添加弹框 -->
    <add-goods-popup v-model="showAddPopup" :sku="keyword[activeTab]" @cancel="handleAddCancel" @confirm="handleAddConfirm"></add-goods-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Vue, Watch } from 'vue-property-decorator';
import OrderDeletePopup from './order-delete-popup.vue';
import AddGoodsPopup from './add-goods-popup.vue';
import { getBrandMaterialListApi } from '@/api/marketing-promotion';
import { getGoodsStockApi } from '@/api/online-order';

const that = this

@Component({
  components: {
    OrderDeletePopup,
    AddGoodsPopup
  }
})
export default class AddGoods extends Vue {
  @Model('update:value', { type: Array, default: () => [] }) value: any[]
  emptyImage = require('@/assets/material-library/v4.2.0/icon_empty.png')
  activeTab: number = 0
  // 筛选关键字
  keyword:string[] = ['', '']
  // 当前要删除的索引
  current = -1
  // 删除弹框
  showDeletePopup = false
  // 显示搜索结果
  showSearchResult = false
  // 搜索结果
  searchResult: any[] = []
  // 当前选中的商品
  selected: any = null

  innerValue: any[] = []

  showAddPopup = false

  @Watch('value', {deep:true})
  onValueChange(val: any[]) {
    this.innerValue = val
  }

  async getData() {
    this.$toast.loading()
    const res = await getBrandMaterialListApi({
      page: 1,
      size: 20,
      sku_item: this.keyword[0],
      upc_code: this.keyword[1],
      include_none_image_product: true   

    })
    this.searchResult = res.data.map(item => {
      item.image = item.brand_source_list[0]?.url
      return item
    }).reduce((prev, item:any) => {
      const { upc_list, ...other} = item
      return prev.concat(upc_list.map(val => {
        return {
          ...other,
          upc_code: val.upc_code,
          product_code: val.variant_code,
          sku_number: val.sku_number,
          style_color: other.style_code + ' ' + other.color
        }
      }))
    }, [])
    this.showSearchResult = true
    this.$toast.clear()
  }

  handleInput = this.$helper.debounce((val: any) => {
    console.log(111111, val);
    
    if(val.length > 3) {
      this.getData()
    } else {
      this.handleClose()
    }
  }, 500)
  handleInput2 = this.$helper.debounce((val: any) => {

    if(val.length > 3) {
      this.getData()
    } else {
      this.handleClose()
    }
  }, 500)

  handleFocus(e) {
    // 页面滚动到指定ID
    this.$nextTick(() => {
      document.documentElement.scrollTo({
        top: document.getElementById('search').offsetTop,
        behavior: "smooth"
      });
    })
    this.handleInput(e.target.value)
  }

  handleClick() {
    (this.$refs.input as any).focus()
  }
  handleClick2() {
    (this.$refs.fieldInput as any).focus()
  }

  /**
   * 扫码
   */
  handleScan() {
    window.wx.invoke('scanQRCode', {
      desc: 'scanQRCode desc',
      needResult: 1, // 默认为0，扫描结果由企业微信处理，1则直接返回扫描结果，
      scanType: ["qrCode", "barCode"], // 可以指定扫二维码还是条形码（一维码），默认二者都有
    }, (res) => {
      try {
        this.keyword[1] = JSON.parse(res.resultStr).scan_code.scan_result
      } catch (error) {
        this.keyword[1] = res.resultStr
      }
      this.getData()
    })
  }

  handleTabChange(e) {
    this.keyword = ['', '']
    this.showSearchResult = false
    this.searchResult = []
  }


  /**
   * 选中商品
   */
  handleSelectItem(item: any) {
    this.showSearchResult = false
    this.selected = JSON.parse(JSON.stringify(item))
    this.handleAdd()
  }

  /**
   * 添加商品
   */
  async handleAdd(skipGetStock = false) {
    if(!this.selected) return
    
    // 神策数据上报 -- 点击侧边栏
    this.$track.trackAddItems({
      add_method: ['SKU CODE添加', '条形码添加'][this.activeTab],
      commodity_name: this.selected.product_name,
      sku_id: this.selected.product_code,
      customer_id: this.$store.getters["user/getExternalId"]
    })

    const index= this.innerValue.findIndex(item => item.product_code === this.selected.product_code)
    if(index === -1) {
      this.innerValue.push({
        ...this.selected,
        loading: !skipGetStock,
        quantity: 1
      })
      this.selected = null
      // 获取库存
      if(!skipGetStock) {
        this.getStock(this.innerValue[this.innerValue.length - 1])
      }
    } else {
      this.innerValue[index].quantity +=1
      if(!skipGetStock && this.innerValue[index].dcl_on_hand == -1) {
        // this.innerValue[index] = {
        //   ...this.innerValue[index],
        //   ...this.selected,
        //   quantity: this.innerValue[index].quantity,
        //   loading: true,
        // }
        this.$set(this.innerValue, index, {
          ...this.innerValue[index],
          ...this.selected,
          quantity: this.innerValue[index].quantity,
          loading: true,
        })
        
        this.getStock(this.innerValue[index])
      }

      this.selected = null
    }

    this.$emit('update:value', this.innerValue)
  }
  /**
   * 
   * @param item 获取库存
   */
  async getStock(item) {
    // this.$set(item, 'loading', true)
    // this.$set(item, 'quantity', 1)
    try {
      const res:any = await getGoodsStockApi({
        store_code_arr: ['OC226'],
        upc_code: item.upc_code
      })
      if (res.code === 200) {
		const payload = res.data || {}
		const list = Array.isArray(payload) ? payload : payload.list || []
		if (payload.dcl_status === 'failed') {
			item.dcl_on_hand = '--'
			item.inventory_warning = '大仓库存信息可能存在延迟，请稍后重试或人工核对库存。'
		} else {
			item.dcl_on_hand = list[0] ? list[0].on_hand : 0
			item.inventory_warning = payload.inventory_status === 'partial'
				? '部分库存信息可能不完整，请稍后刷新查看。'
				: ''
		}
        item.loading = false
	  } else {
		throw new Error(res.message || 'inventory request failed')
      }
      
    } catch (error) {
      item.loading = false
	  item.dcl_on_hand = -1
	  item.inventory_warning = '库存信息暂时无法获取，请稍后再试。给您带来不便，敬请谅解。'
    }
  }
  /**
   * 修改商品数量
   * @param index 
   * @param val 
   */
  handleNumChange(index: number, val: number) {
    console.log(index, val, this.innerValue);
    
    const current = this.innerValue[index].quantity
    const max = this.innerValue[index].dcl_on_hand
    //  || current == max && val===1
    if(current === 1 && val === -1) return
    this.innerValue[index].quantity += val;

    this.$emit('update:value', this.innerValue)
  }

  /**
   * 删除
   * @param index 
   */
  handleDelete(index: number) {
    this.current = index
    this.showDeletePopup = true
  }

  /**
   * 删除确认
   */
  handleDeleteConfirm() {
    this.innerValue.splice(this.current, 1)
    this.current = -1
    this.$emit('update:value', this.innerValue)
  }

  handleClear() {
    this.keyword = ['', '']
    this.showSearchResult = false
    this.searchResult = []
  }
  handleClose() {
    this.showSearchResult = false
    this.searchResult = []
  }
  isError(text) {
    return !/^[0-9A-Za-z]+_[0-9A-Za-z/]+_[0-9A-Za-z/.]+$|^[0-9A-Za-z]+_[0-9A-Za-z/]+$/.test(text);
  }

  handleAddPopup() {
    if(this.isError(this.keyword[this.activeTab])) {
      this.showAddPopup = true
    } else {
      this.handleAddConfirm(this.keyword[this.activeTab])
    }
  }

  /**
   * 手动添加 取消
   */
  handleAddCancel() {
    this.showAddPopup = false
    this.showSearchResult = false
  }

  /**
   * 手动添加 确认
   */
  handleAddConfirm(e) {
    this.selected = {
      product_code: e,
      product_name: '',
      dcl_on_hand: -1,
      quantity: 1,
      style_color: e.split('_').slice(0, 2).join(' ')
    }
    this.handleAdd(true)
    
    this.showAddPopup = false
    this.showSearchResult = false
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  
  .van-icon-clear {
    color: #000;
  }
  .van-tabs__line {
    bottom: 20px
  }

  .van-tabs__wrap {
    padding-bottom: 5px;
  }

  .van-tabs__content {
    position: relative;

    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 50%;
      transform: translateX(-50%);
      width: 345px;
      height: 1px;
      background: #EEEEEE;
      z-index: 1;
    }
  }

  .van-button--disabled {
    opacity: 0.2;
  }
}

.box {
  padding-top: 15px;
  position: relative;

  .search-result {
    width: 315px;
    max-height: 300px;
    overflow-y: scroll;
    background: #FFFFFF;
    box-shadow: 0px 1px 5px 0px rgba(0,0,0,0.2);
    position: absolute;
    top: 60px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 11;
    .search-item {
      height: 50px;
      font-weight: 400;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      padding: 0 15px;
      line-height: 50px;

      &.active {
        background: #F7F7F7;
        font-weight: 500;
      }
    }
  }
}

.input-wrap {
  border: 1px solid #000000;
  height: 40px;
  box-sizing: border-box;
  display: flex;
  align-items: center;
}
.input {
  height: 38px;
  line-height: 38px;
  font-size: 14px;
  padding: 0 15px;
  box-sizing: border-box;
  font-weight: bold;
  width: 100%;
  margin: 0;
}
.input-2 {

  border: 1px solid #000000;
  font-size: 14px;
  padding: 0 15px;
  box-sizing: border-box;
  font-weight: bold;
  width: 100%;

  ::v-deep .van-field__control {
    height: 38px;
  }
}

.input-clear {
  width: 38px;
  height: 38px;
  background: url(~@/assets/online-order/icon_clear.png) center/16px no-repeat;
}

.btn {
  margin-top: 15px;
  height: 36px;
  font-size: 14px;
  border-radius: 0;
}

.tips {
  font-weight: 400;
  font-size: 9px;
  color: #A6A6A6;
  line-height: 13px;
  text-align: center;
  margin-top: 10px;
}

.search-empty {
  padding: 15px;

  .search-empty-text {
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: center;

    &::before {
      content: '';
      width: 20px;
      height: 20px;
      background: url(~@/assets/online-order/icon_tips_2.png) center/100% no-repeat;
    }
  }

  .search-empty-btn {
    margin: auto;
    margin-top: 15px;
    width: 76px;
    height: 27px;
    background: #000000;
    box-shadow: 0px 1px 5px 0px rgba(0,0,0,0.2);
    color: #fff;
    font-size: 12px;
    display: flex;
    align-items: center;
    justify-content: center;

  }
}



.list {
  .tips {
    text-align: left;
    font-weight: 400;
    font-size: 12px;
    color: #676767;
    line-height: 17px;
    margin-top: 15px;
    margin-bottom: 10px;
  }

  .item {
    display: flex;
    position: relative;

    &:not(:last-child) {
      margin-bottom: 15px;
    }

    .van-image {
      background: #EFEFEF;
    }

    .item-info {
      margin-left: 10px;
      flex: 1;
      overflow: hidden;

      .item-code {
        font-weight: 400;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
      }

      .item-name {
        font-weight: 400;
        font-size: 12px;
        color: #96979E;
        line-height: 17px;
        margin-top: 5px;
        height: 17px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .item-stock {
        font-weight: 400;
        font-size: 12px;
        color: #000000;
        line-height: 17px;
        margin-top: 14px;

        &.error {
          color: #F64040;
        }
      }

	  .inventory-warning {
		max-width: 220px;
		margin-top: 4px;
		color: #8a5a00;
		font-size: 11px;
		line-height: 16px;
	  }
    }

    .item-opt {
      position: absolute;
      right: 0;
      bottom: 5px;
      display: flex;
      align-items: center;
      font-size: 14px;

      .item-num {
        width: 54px;
        text-align: center;
      }
      
      .van-button {
        padding: 0;
        width: 24px;
        height: 24px;
        border: 1px solid;
      }
    }


    .item-delete {
      position: absolute;
      top: -10px;
      right: -10px;
      width: 40px;
      height: 40px;
      background: url(~@/assets/online-order/icon_clear_red.png) center/20px no-repeat;
    }
  }
}
</style>
