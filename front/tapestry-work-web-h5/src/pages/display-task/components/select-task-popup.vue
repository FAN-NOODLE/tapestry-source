<template>
  <div class="task-popup">
    <div class="task-name" @click="show=true">
      <span class="task-name__text" 
        :style="{
          paddingRight: renderValue.count>1?'60px':'0'
        }"
        >
        <div class="ellipsis">{{renderValue.name}}</div>
        <div v-if="renderValue.count>1" class="task-num">等{{renderValue.count}}个任务</div>
      </span>
    </div>
    <van-popup v-model="show" class="popup" position="top" get-container="body" @click-overlay="handleCancel">
      <div class="popup-header">
        <div v-if="renderValue.name" class="popup-value">
          <div class="flex-1">
            <div class="ellipsis">{{ renderValue.name }}</div>
          </div>
          <div v-if="renderValue.count>1">等{{ renderValue.count }}个任务</div>
        </div>
        <div v-else class="popup-value placeholder">请选择任务</div>
      </div>
      <!-- 搜索 -->
      <div class="popup-search">
        <van-search v-if="show" v-model="keyword" placeholder="请输入任务名称快捷搜索" :left-icon="searchIcon" :clearable="false" />
        <div v-if="keyword" class="search-clear" @click="handleClear"></div>
      </div>
      <div v-if="renderData.length" class="popup-content">
        <select-item v-for="(item, index) in renderData" :key="index" :data="item" nameKey="display_name" @select="handleSelect"></select-item>
      </div>
      <empty-page
        v-else
        picture="noSearchData"
        paddingTop="60"
        paddingBottom="0"
        text="暂无搜索内容"
      ></empty-page>
      <div class="popup-footer">
        <div class="btn plain" @click="handleReset">重置</div>
        <div class="btn" @click="handleConfirm">确定</div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch, PropSync } from 'vue-property-decorator'
import { Popup } from 'vant';
import SelectItem from './select-item.vue'
import { TreeNode } from '@/utils/tree'
import EmptyPage from "@/components/base/empty-page.vue";

interface ITreeNode extends TreeNode {
  isChecked: boolean
  isLeaf: boolean
}
export interface ITaskItem extends ITreeNode {
  /** 任务名称 */
  display_name: string
  /** 任务ID */
  id: number
}
@Component({
  components: {
    [Popup.name]: Popup,
    SelectItem,
    EmptyPage
  }
})
export default class SelectTaskPopup extends Vue {
  @Prop({ type: Array, default: () => [] }) defaultSelected: number[]
  @Prop({ type: Array, default: () => [] }) selected: number[]
  @Prop({ type: String, default: 'id' }) rowKey: string
  @Prop({ type: Array, default: () => [] }) list: ITaskItem 
  // 显示弹框
  show: boolean = false
  // 数据
  data: ITaskItem[] = []
  /** 搜索图标 */
  searchIcon = require('@/assets/display-task/icon_search.png')
  /** 关键字搜索 */
  keyword = ''
  
  /** 渲染数据 */
  get renderData() {
    return this.data.filter(item => item.display_name.indexOf(this.keyword) > -1)
  }

  /** 渲染值 */
  get renderValue() {
    const list = this.data.filter(item => item.isChecked)

    return {
      name: list.map(item => item.display_name).join(','),
      count: list.length
    }
  }

  @Watch('show') 
  onShowChange() {
    this.init()
  }
  @Watch('list') 
  onDataChange() {
    this.init()
  }

  mounted() {
    this.init()
  }

  /**
   * 初始化
   */
  init(selected?: number[]) {
    this.data = this.list.map(item => {
      const status = (selected||this.selected).includes(item[this.rowKey])
      return {
        ...item,
        isLeaf: true,
        isChecked: status,
        level: 0
      }
    })
  }

  /**
   * 选择
   * @param item 
   */
  handleSelect(item: ITreeNode) {
    item.isChecked = !item.isChecked
  }

  /**
   * 清空搜索
   */
  handleClear() {
    this.keyword = ''
  }

  /**
   * 重置
   */
  handleReset() {
    this.init(this.defaultSelected)
  }

  /**
   * 确认
   */
  handleConfirm() {
    const ids = this.getInfo().ids
    if(ids.length === 0) {
      this.$toast.fail('请先选择任务哦！')
      return
    }
    this.$emit('confirm', ids)
    this.show = false
  }

  /**
   * 取消
   */
  handleCancel() {
    this.show = false
    this.$emit('cancel')
  }

  /** 
   * 获取数据
   */
  getInfo() {
    const ids = this.data
      .filter(item => item.isChecked)
      .map(item => item.id)

    return {
      ids,
      task_name: this.renderValue.name,
      task_num: this.renderValue.count
    }
  }
}
</script>
<style lang="scss" scoped>
  .task-popup {
    flex: 1;
    display: flex;
    overflow: hidden;
    .task-name {
      flex: 1;
      overflow: hidden;
      height: 38px;
      background: #F6F6F6;
      position: relative;
      padding-left: 15px;
      padding-right: 26px;
      display: flex;
      align-items: center;

      &::after {
        content: '';
        display: block;
        position: absolute;
        width: 15px;
        height: 15px;
        right: 8px;
        top: 50%;
        transform: translateY(-50%);
        background: url('~@/assets/display-task/icon_triangle_2.png') center/100% no-repeat;
    
      }

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
  }
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
    display: flex;
    align-items: center;
    padding-right: 40px;

    .flex-1 {
      flex: 1;
      overflow: hidden;
    }

    .ellipsis {
      display: block;
      width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }


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

    &.placeholder {
     color: #A6A6A6;
    }
  }
  .popup-search {
    padding: 15px;
    margin-top: -15px;
    position: relative;
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
  .popup-title {
    font-size: 14px;
    font-weight: 500;
    color: #000000;
    line-height: 20px;
    margin-bottom: 15px;
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
    .select-item {
      padding: 15px 0;
      margin-bottom: 0;
    }
    .select-item__icon {
      display: none;
    }
    .select-item__name {
      font-size: 14px !important;
    }
  }
</style>