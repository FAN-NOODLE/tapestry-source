<template>
  <div class="select-store">
    <div class="store-name" @click="show = true" >
      <span v-if="isTotal">全部区域/全部门店</span>
      <span v-else>{{checked.area}}个区域/{{checked.store}}家门店</span>
    </div>
    <van-popup v-model="show" class="popup" position="top" get-container="body" @click-overlay="handleCancel">
      <div class="popup-header">
        <div class="popup-value">
          <span v-if="isTotal">全部区域/全部门店</span>
          <span v-else>{{checked.area}}个区域/{{checked.store}}家门店</span>
        </div>
      </div>
      <!-- 搜索 -->
      <div class="popup-search">
        <van-search v-if="show" v-model="keyword" placeholder="请输入门店或区域名称快捷搜索" :left-icon="searchIcon" :clearable="false" />
        <div v-if="keyword" class="search-clear" @click="handleClear"></div>
      </div>
      <div v-if="data.length && !isEmpty" class="popup-content">
        <div class="popup-title">全部门店</div>
        <div ref="listRef" class="list">
          <template v-for="(item, index) in data">
            <select-item :key="index" :data="item" :show="showItem" @expand="handleExpand" @select="handleSelect"></select-item>
          </template>
        </div>
      </div>
      <empty-page
        v-else
        picture="noSearchData"
        paddingTop="60"
        paddingBottom="100"
        text="暂无搜索内容"
      ></empty-page>
      <div class="popup-footer">
        <div class="btn plain" @click="handleReset">重置</div>
        <div class="btn" @click="handleConfirm">确定 <span v-if="checked.store">（{{checked.store}}家门店）</span></div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch, PropSync } from 'vue-property-decorator'
import { Popup } from 'vant';
import SelectItem from './select-item.vue'
import { TreeNode,traverseTreeList, bfs, dfs, hasMatchingChild } from '@/utils/tree'
import EmptyPage from "@/components/base/empty-page.vue";

interface ITreeNode extends TreeNode {
  isChecked: boolean
  isExpand: boolean
  isHalf: boolean
  isLeaf: boolean
}
export interface IStoreItem extends ITreeNode {
  /** 节点ID */
  id: number
  /** 节点名称 */
  name: string
  /** 节点级别 1 渠道 2大区 3分区 4门店 */
  type: string
  /** 父节点ID */
  parent_id: number
  /** 父节点名称 */
  parent_name: string;
  /** 门店编码 */
  store_code: string;
  /** 门店名称 */
  store_name: string;
  /** 子节点 */
  children?: IStoreItem[]
}
@Component({
  components: {
    [Popup.name]: Popup,
    SelectItem,
    EmptyPage
  }
})
export default class SelectStorePopup extends Vue {
  @Prop({ type: Array, default: () => [] }) defaultSelected: number[]
  @Prop({ type: Array, default: () => [] }) selected: number[]
  @Prop({ type: String, default: 'id' }) rowKey: string
  @Prop({ type: Array, default: () => [] }) list: IStoreItem[]
  
  show: boolean = false
  // 数据
  data: IStoreItem[] = []
  /** 搜索图标 */
  searchIcon = require('@/assets/display-task/icon_search.png')
  /** 搜索关键字 */
  keyword = ''

  /** 选中数量 */
  checked = {
    store: 0,
    storeIds: [],
    area: 0,
    areaIds: []
  }

  isEmpty = false

  get isTotal() {
    return this.data.every(item => item.isChecked)
  }

  @Watch('keyword')
  onKeywordChange() {
    this.isEmpty = false
    this.$nextTick(() => {
      const listDom: HTMLElement|undefined = this.$refs.listRef as HTMLElement
      if(listDom) {
        this.isEmpty =  listDom.children.length === 0
      }
    })
  }

  @Watch('data', {deep: true})
  handleDataChange(value) {
    const storeIds = new Set()
    const areaIds = new Set()
    value.forEach(item => {
      bfs(item, (node) => {
        // 门店
        if (node.type === 4 && node.isChecked) {
          storeIds.add(node.id)
        }
        // 区域
        if (node.type === 3 && (node.isChecked|| node.isHalf)) {
          areaIds.add(node.id)
        }
      })
    });
    this.checked.store = storeIds.size
    this.checked.storeIds = Array.from(storeIds)
    this.checked.area = areaIds.size
    this.checked.areaIds = Array.from(areaIds)
  }

  @Watch('show') 
  onShowChange() {
    this.init()
  }

  @Watch('list', {deep: true})
  onListChange() {
    this.init()
  }
  
  mounted() {
    this.init()
  }

  init(selected?: number[]) {
    const data = traverseTreeList<IStoreItem>(JSON.parse(JSON.stringify(this.list)),(item) => {
      const status = (selected||this.selected).includes(item[this.rowKey])
      item.isChecked = status
      if(status) {
        dfs(item, (node) => {
          node.isChecked = status
          node.isHalf = false
        })
      }
      return item
    })
    
    // 展开选中门店
    this.data = data.map(item => {
      dfs(item, (node) => {
        this.updateParent(node as ITreeNode)
      })
      return item
    })
  }

  updateParent(node:ITreeNode) {
    if (node.parent) {
      node.parent.isChecked = node.parent.children.every(child => child.isChecked);
      
      if(!node.parent.isChecked)  {
        node.parent.isHalf = node.parent.children.some(child => child.isChecked || child.isHalf);
      } else {
        node.parent.isHalf = false
      }

      if (node.parent.isLeaf) {
        node.parent.isHalf = false
      } else {
        node.parent.isExpand = node.parent.isExpand || node.parent.isHalf || node.parent.isChecked
      }

      this.updateParent(node.parent as ITreeNode);
    }
  }

  handleSelect(item: ITreeNode) {
    const status = !item.isChecked
    if(item.isLeaf) {
      item.isChecked = status
    } else {
      bfs(item, (node) => {
        node.isChecked = status
        node.isHalf = false
      })
    }
    // 重新判断选中状态
    this.updateParent(item)
  }
  handleExpand(item: ITreeNode) {
    item.isExpand = !item.isExpand
  }

  handleReset() {
    this.init(this.defaultSelected)
  }

  handleConfirm() {
    if(!this.checked.store) {
      this.$toast.fail('请至少选择一个门店或区域')
      return
    }
    
    this.$emit('confirm', this.getInfo().area_ids)
    this.show = false
  }

  handleCancel() {
    this.show = false
    this.$emit('cancel')
  }

  showItem(item: IStoreItem) {
    if(this.keyword.trim()) {
      // 判断字节点包含
      const value = hasMatchingChild(item, node => {
        return node.name.indexOf(this.keyword) > -1
      })
      if(value) {
        item.isExpand = true
      }
      return value
    }
    return true
  }

  handleClear() {
    this.keyword = ''
  }

  /** 
   * 获取数据
   */
  getInfo() {
    const ids = []
    this.data.forEach(item => {
      bfs(item, (node) => {
        if (node.isChecked) {
          ids.push(node.id)
        }
      })
    })

    return {
      area_ids: ids,
      area_num: this.checked.area,
      store_num: this.checked.store,
      is_total: this.isTotal
    }
  }
}
</script>
<style lang="scss" scoped>
  .popup-header {
    padding: 10px 15px;
  }
  .popup-value {
    width: 345px;
    height: 38px;
    line-height: 38px;
    background: #F6F6F6;
    position: relative;
    font-size: 12px;
    font-weight: 400;
    color: rgba(0,0,0,0.85);
    padding: 0 15px;


    &::after {
      content: '';
      display: block;
      position: absolute;
      width: 15px;
      height: 15px;
      right: 8px;
      top: 50%;
      transform: translateY(-50%);
      background: url('~@/assets/display-task/icon_triangle_1.png') center/100% no-repeat;
    }
  }
  .popup-search {
    padding: 15px;
    position: relative;
    margin-top: -15px;
    .van-search {
      height: 38px;
      padding: 0;
      border: 1px solid #A6A6A6;
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
  .popup-title {
    font-size: 14px;
    font-weight: 500;
    color: #000000;
    line-height: 20px;
  }
  .popup-content {
    padding: 0 15px;
    max-height: 60vh;
    overflow-y: scroll;
  }

  .popup-footer {
    padding: 15px;
    display: flex;
    align-items: center;
    .btn {
      height: 36px;
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      font-weight: 500;
      color: #FFFFFF;
      line-height: 20px;
      background-color: #000;
      border: 1px solid #000000;
      box-sizing: border-box;

      &.plain {
        background-color: #fff;
        color: #000000;
        margin-right: 15px;
      }
    }
  }

  ::v-deep {
    .default {
      height: auto;
    }
    .select-item {
      padding-top: 15px;
      padding-bottom: 15px;
      margin-bottom: 0;
    }
    .select-item__collapse .select-item {
      padding-top: 0;
      padding-bottom: 0;
      &:not(:last-child) {
        margin-bottom: 15px;
      }
    }
  }
  .store-name {
    padding-left: 10px;
    padding-right: 18px;
    position: relative;

    &::after {
      content: '';
      display: block;
      position: absolute;
      width: 15px;
      height: 15px;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
      background: url('~@/assets/display-task/icon_triangle_2.png') center/100% no-repeat;
  
    }
  }
</style>