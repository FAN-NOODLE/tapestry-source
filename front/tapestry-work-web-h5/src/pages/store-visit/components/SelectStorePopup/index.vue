<template>
  <div>
    <van-popup v-model="visible" position="bottom" round @close="handleClose" @click-overlay="handleClose">
      <div class="popup">
        <div class="popup-header">
          <div class="popup-title">{{ $t('storeVisit.popup.title.2') }}</div>
          <div class="popup-close" @click.stop="handleClose"></div>
        </div>
        <div class="popup-search">
          <div class="input" @click="handleInputClick">
            <div class="icon-search"></div>
            <input class="input-inner" ref="input" type="text" :placeholder="String($t('storeVisit.placeholder.1'))" v-model="searchValue" @input="handleInput" />
            <div class="icon-close" v-if="searchValue.length" @click="handleClear"></div>
          </div>
        </div>
        <div v-if="!showSearchResult" class="popup-tabs">
          <div class="tabs">
            <div class="tab" :class="{active: tabIndex[0] === index}" v-for="item, index in storeTree" :key="index" @click="handleTabClick(index)">{{ item.name }}</div>
          </div>
        </div>
        <div v-if="!showSearchResult && showStoreList" class="popup-nav" @click="handleNavClick">
          <div class="icon-left"></div>
          <div class="nav-item" v-for="item, index in navList" :key="index">{{ item.label }}</div>
        </div>
        <div v-if="!showSearchResult" class="popup-content">
          <div class="wrap" :class="{active: showStoreList}">
            <div class="tree">
              <div class="tree-item" v-for="(item, index) in list" :key="index">
                <div class="tree-item-head" @click="handleToggleExpand(index, $event)">
                  <div class="tree-item-title" :class="{active: item.isExpand}">{{ item.name }}</div>
                  <div :class="[item.isExpand?'icon-minus':'icon-plus']"></div>
                </div>
                <div v-if="item" class="tree-item-body">
                  <div class="tree-sub-item" v-for="sub, subIndex in item.children" :key="subIndex" @click="handleClick(index, subIndex)">
                    <div class="tree-sub-item-title">{{ sub.name }}</div>
                    <div class="icon icon-right"></div>
                  </div>
                </div>
              </div>
            </div>
            <div class="list">
              <div class="store-item" :class="{active: isChecked(item.third_id)}" v-for="(item, index) in storeList" :key="index" @click="handleSelect(item)">
                <div class="flex">
                  <div class="flex-1">
                    <div class="store-name">{{ item.store_name }}</div>
                    <div class="store-code">{{ item.store_code }}</div>
                  </div>
                  <div class="icon-check"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-if="!showSearchResult && !immediate" class="popup-footer">
          <div v-if="selection.length" class="popup-selection">
            <div class="label">
              <div>{{ $t('storeVisit.popup.tips.5') }}{{ selection.length }}{{ $t('storeVisit.popup.tips.6') }}：</div>
              <div class="all" @click="handleShowStore">{{ $t('storeVisit.popup.btn.5') }}</div>
            </div>
            <div class="scroll-x">
              <div class="select-list">
                <div class="select-item" v-for="(item, index) in selection" :key="index">
                  <div class="name">{{ item.store_name }}</div>
                  <div class="icon-close-circle" @click="handleRemoveSelection(index)"></div>
                </div>
              </div>
            </div>
          </div>
          <div class="flex">
            <van-button color="#000000" plain @click="handleCancel">{{ $t('storeVisit.popup.btn.3') }}</van-button>
            <van-button color="#000000" :disabled="!selection.length" @click="handleConfirm">{{ $t('storeVisit.popup.btn.4') }}</van-button>
          </div>
        </div>
        <div v-if="showSearchResult" class="popup-content search-result">
          <div v-if="searchResult.length" class="result">
            <div class="list">
              <div class="item" v-for="item,index in searchResult" :key="index" @click="handleSelectItem(item)">
                <div class="item-name">{{ item.store_name }}</div>
                <div class="item-tag">{{ item.store_code}}</div>
              </div>
            </div>
          </div>
          <div v-else class="empty">{{ $t('storeVisit.empty.2') }}</div>
        </div>
      </div>
    </van-popup>
    <!-- 查看已选人员 -->
    <van-popup v-model="showSelectStore" round position="center">
      <div class="popup popup-user">
        <div class="popup-header">
          <div class="popup-title">{{ $t('storeVisit.popup.tips.5') }}{{ selection.length }}{{ $t('storeVisit.popup.tips.6') }}</div>
          <div class="popup-close" @click.stop="handleCloseStore"></div>
        </div>
        <scrollbar class="popup-content" max-height="400">
          
          <div class="user-list">
            <div class="user-item" v-for="(item, index) in selection" :key="index">
              <div class="flex-1">
                <div class="user-name">{{item.store_name}}</div>
              </div>
              <div class="icon-close-circle-2" @click="handleDeleteStore(index)"></div>
            </div>
          </div>
        </scrollbar>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';

import { getStoreTreeApi } from '@/api/display-task';
import { TreeNode, traverseTreeList } from '@/utils/tree';
import { nextTick } from 'vue';
import { searchStoreApi } from '@/api/store-visit';
import Scrollbar from '../scrollbar.vue'

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


let that =null
@Component({
  components: {
    Scrollbar
  }
})
export default class SelectStorePopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: Boolean, default: true }) multiple: boolean;
  @Prop({ type: Number, default: 0 }) limit: number;  // multiple 为true时有效
  @Prop({ type: Boolean, default: false }) immediate: boolean;
  @Prop({ type: String, default: 'third_id' }) rowKey: string;
  @Prop({  type: Array, default: () => [] }) defaultValue: number[];
  @Prop({  type: Array, default: () => [] }) storeTree: IStoreItem[];
  
  visible: boolean = false;
  showAddEmail = false;
  
  showSearchResult = false;
  showSelectStore = false;
  searchValue: string = '';
  searchResult: any[] = [];

  // 门店列表
  // storeTree: IStoreItem[] = [];
  // 选中的门店
  selection: IStoreItem[] = [];
  // 导航
  tabIndex: number[] = [0,-1,-1];

  @Watch('show', { immediate: false })
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      
      nextTick(() => {
        this.setDefaultExpand();
      })
      // 清除选项
      this.selection = []
      // 获取选中项
      traverseTreeList(this.storeTree, (node:any) => {
        if(this.defaultValue.includes(node[this.rowKey])) {
          this.selection.push(node)
        }
      });
      if(this.defaultValue.length == 0) {
        this.tabIndex = [0,-1,-1];
      }
    }
  }

  get list() {
    if(this.storeTree[this.tabIndex[0]]) {
      return this.storeTree[this.tabIndex[0]].children;
    }
    return [];
  }

  get storeList() {
    const [one, two, three] = this.tabIndex;
    if(this.storeTree[one]) {
      if(this.storeTree[one].children[two]) {
        if(this.storeTree[one].children[two].children[three]) {
          return this.storeTree[one].children[two].children[three].children;
        }
      }
    }
    return [];
  }

  get navList() {
    const [one, two, three] = this.tabIndex;
    const result = []
    if(this.storeTree[one]) {
      // result.push({
      //   label: this.storeTree[one].name,
      //   value: one,
      // })
      if(this.storeTree[one].children[two]) {
        result.push({
          label: this.storeTree[one].children[two].name,
          value: two,
        })
        if(this.storeTree[one].children[two].children[three]) {
          result.push({
            label: this.storeTree[one].children[two].children[three].name,
            value: three,
          })
        }
      }
    }
    return result;
  }

  get showStoreList() {
    return this.navList.length>0;
  }

  get isChecked() {
    return (id) => {
      return this.selection.some(item => item.third_id === id);
    }
  }

  mounted() {
    that = this
    // this.getStoreTree()
  }
  destroyed() {
  }
  deactivated() {
    this.visible = false
  }

  async getStoreTree() {
    const res: any = await getStoreTreeApi({
      has_children: true,
      ignore_close_store: true
    })
    if (res.code !== 200) {
      this.$toast.fail(res.message);
      return
    }
    
    this.storeTree = res.data.filter((item, index) => index<2)

  }

  handleClose() {
    this.$emit('update:show', false);
    this.$emit('close');
  }

  /**
   * 取消
   */
  handleCancel() {
    this.visible = false
    this.$emit('cancel');
    this.showSearchResult = false;
    this.searchValue = '';
    this.searchResult = [];
  }
  /**
   * 确认事件
   */
  handleConfirm() {
    this.visible = false
    this.$emit('confirm', this.selection);
    this.showSearchResult = false;
    this.searchValue = '';
    this.searchResult = [];
  }

  /**
   * 点击获取焦点
   */
  handleInputClick() {
    (this.$refs.input as any).focus()
  }

  /**
   * 搜索
   * @param e 
   */
  handleInput = this.$helper.debounce((e) => {
    console.log(e.target.value);
    this.showSearchResult = true;
    this.searchStore(e.target.value)
  }, 500, false)

  /**
   * 搜索成员
   * @param keyword 
   */
  async searchStore(keyword) {
    if (keyword.length === 0) {
      return;
    }
    const res:any = await searchStoreApi({
      store_str: keyword,
      ignore_close_store: true
    })
    if (res.code === 200) {
      this.searchResult = res.data.map((item) => {
        return {
          third_id: item.wechat_department_id,
          store_name: item.store_name,
          store_code: item.third_store_code,
        }
      });
      this.showSearchResult = true;
    }
  }

  handleClear() {
    this.searchValue = '';
    this.showSearchResult = false;
    this.searchResult = [];
  }

  handleTabClick(index) {
    this.$set(this, 'tabIndex', [index, -1, -1]);
    nextTick(() => {
      this.setDefaultExpand();
    })
  }
  handleNavClick() {
    this.$set(this, 'tabIndex', [this.tabIndex[0], -1, -1]);
  }

  /**
   * 选择区域
   * @param index 
   */
  handleClick(index, subIndex) {
    this.$set(this, 'tabIndex', [this.tabIndex[0], index, subIndex]);
  }

  handleToggleExpand(index, event) {
    event.stopPropagation();
    const isExpand = this.list[index].isExpand;
    // 关闭其他展开的节点
    this.list.forEach((item, i) => {
      if (i !== index) {
        this.$set(this.storeTree[this.tabIndex[0]].children[i], 'isExpand', false);
      } else {
        this.$set(this.storeTree[this.tabIndex[0]].children[i], 'isExpand', !isExpand);
      }
    })
    // 展开收起内容
    document.querySelectorAll('.tree-item-body').forEach((item, i) => {
      if (i !== index || isExpand) {
        item.setAttribute('style', 'height: 0px');
      } else {
        const totalHeight = [...item.children].reduce((acc, child) => {
          return acc + child.clientHeight;
        }, 0);
        item.setAttribute('style', 'height: ' + totalHeight + 'px');
      }
    })
    
  }

  setDefaultExpand() {
    const [one] = this.tabIndex;
    if(!this.storeTree[one]) return;
    this.storeTree[one].children.forEach((item, index) => {
      if(index === 0) {
        this.$set(this.storeTree[one].children[index], 'isExpand', true);
        document.querySelectorAll('.tree-item-body').forEach((item, i) => {
          if (i !== index) {
            item.setAttribute('style', 'height: 0px');
          } else {
            const totalHeight = [...item.children].reduce((acc, child) => {
              return acc + child.clientHeight;
            }, 0);
            item.setAttribute('style', 'height: ' + totalHeight + 'px');
          }
        })
      }
    })
  }

  /**
   * 选择门店
   * @param item 
   */
  handleSelect(item) {
    // 单选
    if(!this.multiple || this.limit === 1) {
      this.selection = [item];
      this.$emit('select', {
        track: {
          filter_option: this.navList.map(item => item.label).join('-'),
          seek_method: '筛选',
        },
        item
      })
      if(this.immediate) {
        this.handleConfirm()
      }
      return
    }
    if(this.isChecked(item.third_id)) {
      // 如果以选中，则取消选中
      this.selection = this.selection.filter(i => i.third_id !== item.third_id);
    } else if(this.limit > 1 && this.selection.length >= this.limit) {
      this.$toast(`最多可同时发起${this.limit}家门店的巡店预约!`);
    } else {
      this.selection.push(item);
      this.$emit('select', {
        track: {
          filter_option: this.navList.map(item => item.label).join('-'),
          seek_method: '筛选',
        },
        item
      })
    }
  }

  handleRemoveSelection(index) {
    this.selection.splice(index, 1);
  }

  /**
   * 选择搜索 门店
   * @param item 
   */
  handleSelectItem(item) {
    // 单选
    if(!this.multiple || this.limit === 1) {
      this.selection = [item];
      this.$emit('select', {
        track: {
          search_word: this.searchValue,
          seek_method: '搜索',
        },
        item
      })
    } else if(this.limit > 1 && this.selection.length >= this.limit) {
      this.$toast(`最多可同时发起${this.limit}家门店的巡店预约!`);
    } else if(!this.isChecked(item.third_id)) {
      this.selection.push(item);
      
      this.$emit('select', {
        track: {
          search_word: this.searchValue,
          seek_method: '搜索',
        },
        item
      })
    }
    this.showSearchResult = false
    this.searchValue = ''
    this.searchResult = []

    function getParentList(node) {
        return node.parent ? [...getParentList(node.parent), node] : [node];
    }

    // 确定 item 在 storeTree 中的位置，以及路径
    traverseTreeList(this.storeTree, (node:any) => {
      if(node.third_id === item.third_id) {
        // console.log(node, getParentList(node));
        const parentList = getParentList(node)
        // 获取 parentListIndex
        const index = []
        index[0] = this.storeTree.findIndex(item => item.id === parentList[0].id)
        index[1] = parentList[0] ? parentList[0].children.findIndex(item => item.id === parentList[1].id) : -1
        index[2] = parentList[1] ? parentList[1].children.findIndex(item => item.id === parentList[2].id) : -1
        // console.log(index);
        this.tabIndex = index
        this.$nextTick(() => {  
          this.setDefaultExpand()
        })
      }
    });

    if(this.immediate) {
      this.handleConfirm()
    }
  }

  handleShowStore() {
    this.showSelectStore = true;
  }
  handleCloseStore() {
    this.showSelectStore = false;
  }
  handleDeleteStore(index) {
    this.selection.splice(index, 1);
  }
}
</script>
<style lang="scss" scoped>

.icon-close-circle-2 {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_close_circle_2.png) center/16px no-repeat;
}

.icon-plus {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_plus.png) center/100% no-repeat;
}
.icon-minus {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_minus.png) center/100% no-repeat;
}

.icon-right {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_right_black.png) center/100% no-repeat;
}
.icon-left {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_left.png) center/100% no-repeat;
}
.icon-check {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_check.png) center/24px no-repeat;
}
.icon-close-circle {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_close_circle.png) center/16px no-repeat;
}

.tree {
  width: 100%;
  .tree-item {
    box-sizing: border-box;
    overflow: hidden;
    position: relative;

    &:not(:last-child)::after {
      content: "";
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      height: 1PX;
      background: rgba(248, 249, 250, 1);
      z-index: 1;
      width: 100%;
    }

    .tree-item-head {
      height: 68px;
      display: flex;
      align-items: center;
      background-color: #FFFFFF;
      padding: 0 24px;
      position: relative;
      z-index: 1;

      .tree-item-title {
        font-size: 14px;
        color: rgba(0, 0, 0, 1);
        flex: 1;

        &.active {
          font-weight: 600;
        }
      }

    }
    .tree-item-body {
      padding: 0 24px;
      background-color: rgba(248, 249, 250, 1);
      height: 0;
      overflow: hidden;
      transition: all 0.3s;

      .tree-sub-item {
        display: flex;
        align-items: center;
        height: 60px;
        box-sizing: border-box;
        position: relative;

        &:not(:last-child)::after {
          content: "";
          position: absolute;
          left: 0;
          right: 0;
          bottom: 0;
          height: 1PX;
          background: rgba(233, 236, 239, 1);
          z-index: 1;
          width: 100%;
        }

        .tree-sub-item-title {
          font-size: 14px;
          color: rgba(0, 0, 0, 0.8);
          flex: 1;
        }
      }
    }
  }
}


.tabs {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: rgba(248, 249, 250, 1);
  padding: 2px;

  .tab {
    flex: 1;
    height: 20px;
    line-height: 20px;
    text-align: center;
    font-size: 12px;
    color: rgba(0, 0, 0, 0.3);
    border-radius: 3px;
    transition: all 0.2s;

    &.active {
      font-weight: bold;
      color: rgba(0, 0, 0, 1);
      background-color: rgba(255, 255, 255, 1);
    }

    &:nth-child(1) {
      margin-right: 2px;
    }
  }
}


.popup {
  position: relative;
  width: 100vw;
  height: 90vh;
  background: #FFFFFF;
  display: flex;
  flex-direction: column;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    
    .popup-title {
      font-weight: 600;
      font-size: 16px;
      line-height: 24px;
      color: #000000;
    }

    .popup-close {
      width: 72px;
      height: 64px;
      position: absolute;
      right: 0;
      top: 0;
      background: url(~@/assets/store-visit/icon_popup_close.png) center/24px no-repeat;
    }
  }

  .popup-search {
    padding: 12px;
    
    .input {
      height: 32px;
      border-radius: 16px;
      background-color: rgba(248, 249, 250, 1);
      padding: 4px 8px;
      box-sizing: border-box;
      display: flex;
      align-items: center;

      .icon-search {
        width: 24px;
        height: 24px;
        background: url(~@/assets/store-visit/icon_search.png) center/24px no-repeat;
      }

      .input-inner {
        flex: 1;
        margin-left: 8px;
        font-size: 12px;
        background-color: transparent;
      }

      .input-inner::placeholder {
        color: rgba(0, 0, 0, 0.2);
      }

      .icon-close {
        width: 24px;
        height: 24px;
        background: url(~@/assets/store-visit/icon_close_circle.png) center/24px no-repeat;

      }
    }
  }

  .popup-nav {
    padding: 12px;
    display: flex;
    align-items: center;

    .icon-left {
      margin-right: 4px;
    }

    .nav-item {
      height: 20px;
      line-height: 20px;
      font-size: 10px;
      font-weight: 600;
      padding: 0 12px;
      background-color: rgba(233, 236, 239, 1);
      color: #000000;
      border-radius: 3px;
      margin-right: 11px;
      position: relative;

      &:not(:last-child)::after {
        content: ";";
        width: 11px;
        text-align: center;
        position: absolute;
        right: 0;
        top: 50%;
        transform: translate(100%, -50%);
        color: rgba(153, 153, 153, 1);
      }
    }
  }

  .popup-tabs {
    padding: 0 12px 12px 12px;


  }

  .popup-content {
    flex: 1;
    background-color: rgba(248, 249, 250, 1);
    display: flex;
    overflow: hidden;

    .wrap {
      width: 200vw;
      height: 100%;
      display: flex;
      transition: all 0.3s;

      .tree {
        width: 100vw;
        height: 100%;
        overflow-y: auto;

        &::-webkit-scrollbar {
          display: none;
        }

      }
      .list {
        width: 100vw;
        height: 100%;
        padding: 0 24px;
        overflow-y: auto;
        &::-webkit-scrollbar {
          display: none;
        }

        .store-item {
          height: 64px;
          position: relative;
          .flex {
            display: flex;
            align-items: center;
            height: 100%;
          }
          .flex-1 {
            display: flex;
            align-items: center;
            flex: 1;
          }

          .store-name {
            font-size: 14px;
            color: rgba(0, 0, 0, 0.8);
          }
          .store-code {
            font-size: 12px;
            color: rgba(0, 0, 0, 0.5);
            margin-left: 8px;
          }

          .icon-check {
            display: none;
          }

          &:not(:last-child) {
            border-bottom: 1px solid rgba(233, 236, 239, 1);
          }

          &.active {
            font-weight: 600;

            .icon-check {
              display: block;
            }
          }

        }
      }
      &.active {
        transform: translateX(-50%);
      }
    }
    

    &.search-result {
      padding: 12px;
      height: auto;
      background-color: rgba(255, 255, 255, 1);

      .result {
        height: 100%;
        border-radius: 12px;
        background-color: rgba(248, 249, 250, 1);
        padding: 8px 8px 8px 40px;
        position: relative;
        width: 100%;

        .list {
          height: 100%;
          overflow-y: auto;
          &::-webkit-scrollbar {
            display: none;
          }
        }

        &::before {
          content: "";

          position: absolute;
          left: 8px;
          top: 8px;
          width: 24px;
          height: 24px;
          background: url(~@/assets/store-visit/icon_search_black.png) center/24px no-repeat;
        }

        .item {
          height: 24px;
          display: flex;
          align-items: center;

          &:not(:last-child) {
            margin-bottom: 12px;
          }
          
          .item-name {
            font-size: 12px;
            color: rgba(0, 0, 0, 1);
          }

          .item-tag {
            height: 16px;
            line-height: 16px;
            font-size: 10px;
            padding: 0 4px;
            margin-left: 4px;
            color: rgba(153, 153, 153, 1);
          }
        }
      }

      .empty {
        font-size: 12px;
        color: rgba(153, 153, 153, 1);
        text-align: center;
        padding-top: 50px;
        width: 100%;
      }
    }
  }


  .popup-footer {
    z-index: 2;
    background-color: #fff;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);
    border-top: 1PX solid rgba(233, 236, 239, 1);


    .popup-selection {
      padding: 16px 0;
      .label {
        font-size: 12px;
        color: rgba(0, 0, 0, 1);
        font-weight: 600;
        line-height: 16px;
        padding: 0 24px;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .all {
          font-size: 12px;
          font-weight: bold;
          color: #0D6EFD;
        }
      }

      .scroll-x {
        margin-top: 4px;
        overflow-x: scroll;

        &::-webkit-scrollbar {
          display: none;
        }
      }

      .select-list {
        display: flex;
        flex-wrap: nowrap;
        width: fit-content;
        padding: 0 24px;

        .select-item {
          height: 24px;
          line-height: 24px;
          font-size: 12px;
          color: rgba(0, 0, 0, 0.8);
          border-radius: 3px;
          display: flex;
          align-items: center;

          .name {
            white-space: nowrap;
          }

          .icon-close-circle {
            margin-left: 4px;
          }

          &:not(:last-child) {
            margin-right: 8px;

            &::after {
              content: ";";
              margin-left: 8px;
              font-size: 10px;
              color: rgba(153, 153, 153, 1);
            }
          }
        }
      }
    }

    .flex {
      padding: 12px 24px;
      display: flex;
      align-items: center;
    }

    .van-button {
      flex: 1;
      height: 40px;
      border-radius: 4px;
      font-weight: 600;

      &:nth-child(2) {
        margin-left: 8px;
      }

      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }
    }
  }
}
.popup-user {
  width: 311px;
  height: auto;

  .popup-content {
    width: 100%;
    height: auto;
    max-height: 400px;
    background: #fff;
    
    .user-list {
      box-sizing: border-box;
      width: 100%;
      padding: 16px 24px 24px;
      height: fit-content;
    }
    

    .user-item {
      width: 100%;
      // height: 32px;
      background-color: #F8F9FA;
      border-radius: 4px;
      padding: 8px 32px 8px 8px;
      line-height: 17px;
      position: relative;
      .flex-1 {
        flex: 1;
        display: flex;
        align-items: center;
      }

      .user-name {
        font-size: 12px;
        color: rgba(0, 0, 0, 1);
      }

      .user-dept {
        font-weight: 700;
        font-size: 10px;
        color: rgba(0, 0, 0, 1);
        background-color: rgba(206, 212, 218, 1);
        line-height: 16px;
        padding: 0 4px;
        border-radius: 2px;
        margin-left: 4px;

        &.ext {
          color: rgba(153, 153, 153, 1);
          background-color: rgba(255, 255, 255, 1);
        }
      }

      .icon-close-circle-2 {
        width: 32px;
        height: 32px;
        position: absolute;
        right: 0;
        top: 50%;
        transform: translateY(-50%);
      }

      &:not(:last-child) {
        margin-bottom: 4px;
      }
    }
  }
}
</style>