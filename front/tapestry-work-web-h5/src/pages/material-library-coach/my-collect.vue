<template>
  <div>
    <div class="container" v-if="!showSendPreview">
      <div class="header">
        <div class="search">
          <div class="search-tag" @click="handleSelectTag">属性
            <div v-if="searchParams.category_ids.length" class="count">{{ searchParams.category_ids.length }}</div>
          </div>
          <div class="search-inner">
            <van-search
              v-model="searchParams.keyword"
              :clearable="false"
              :left-icon="searchIcon"
              placeholder="请输入素材名称/商品名称/款号/款色"
              @input="handleSearch"
            />
            <div v-if="searchParams.keyword" class="delete" @click="handleClear"></div>
          </div>
        </div>
      </div>
      <div v-if="searchParams.group_name" class="group-name">
        <div class="text">
          {{ searchParams.group_name }}
          <div class="delete" @click="handleClearGroup"></div>
        </div>
      </div>
      <van-tabs
        v-model="active"
        @change="handleTabChange"
        title-inactive-color="#000000"
        background="#f4f7fb"
        color="#000000"
        line-width="30px"
        line-height="3px"
        sticky
      >
        <van-tab v-for="(tab, tabIndex) in tabs" :title="tab.label" :key="tabIndex">
          <template #title>
            <div class="tab-text">
              {{tab.label}}
              <div
                class="num-icon-box"
                v-if="numAll[2]"
              >
                <div class="num-text">{{ numAll[2] }}</div>
              </div>
            </div>
          </template>
          <div v-if="isApplets && tab.tips" class="tips">{{ tab.tips }}</div>
          <!-- <select-grid v-if="tab.value==2 && isCoach" v-model="tab.grid" :showFirstTips="isFirstOpen" @change="handleGridChange"></select-grid> -->

          <!-- 多选 -->
          <div class="select-box">
            <div class="select-left">素材数：{{ tab.list.length }}</div>
            <div v-if="tab.value !== 3" class="select-right" @click="selectAll"><i class="check-box" :class="tabs[active].isAllChecked ? 'checked-box': ''"></i>全选</div>
          </div>
          <!-- 多选 -->

          <van-list 
            v-model="tab.isLoading"
            :finished="tab.isFinished"
            :finished-text="tab.list.length ? '没有更多了' : ''"
            @load="getData" 
          >
            <van-row :gutter="tab.gutter">
              <van-col v-for="(item, index) in tab.list" :key="index" :span="tab.span">
                <material-item 
                  page="collect"
                  :item="item"
                  :selection="selection"
                  :small="tab.grid === 1"
                  @select="handleSelect"
                  @send="handleSend"
                  @preview="handlePreview($event, index)"
                  @preview-video="handlePreviewVideo"
                  @selectImage="handleSelectImage"
                  @cancelCollect="handleCancelCollect"
                  @edit="handleEditText"
                ></material-item>
              </van-col>
            </van-row>
          </van-list>
          <empty-page
            v-if="!tab.list.length && !tab.isLoading"
            picture="noData"
            paddingTop="60"
            paddingBottom="265"
            :text="emptyText(tab)"
          ></empty-page>
        </van-tab>
      </van-tabs>
      <!-- footer -->
      <div class="footer" v-if="!isApplets && tabs[active].value !== 3">
        <div class="num-left">
          <div>已选：{{ selection.length }}</div>
          <div
            class="max-num-left"
            v-if="selection.length >= 9"
          >最多可选9个</div>
        </div>
        <!-- todo:清空已选按钮 -->
        <div
          class="clear-btn"
          @click="handleClearSelection"
          v-if="selection.length > 0"
        >
          清空已选
        </div>

        <!-- 取消收藏 -->
        <div
          class="button-right"
          :class="[!selection.length ? 'dis-btn' : '']"
          @click="handleCollect()">取消收藏</div>
        <!-- 发送 -->
        <div
          class="button-right"
          :class="[(!selection.length || selection.length > 9) ? 'dis-btn' : '']"
          @click="handleSend()"
        >
          {{ selection.length > 1 ? '组合发送' : '发送' }}
        </div>
      </div>
      <!-- 选择属性 -->
      <select-tag-popup 
        :isCollect="1"
        :show.sync="showTagPopup" 
        :value="searchParams.category_ids"
        @confirm="handleSelectTagConfirm"
      ></select-tag-popup>
      <!-- 选择分组 -->
      <select-group-popup 
        :show.sync="showGroupPopup"
        :value="searchParams.group_value"
        @cancel="handleSelectGroupCancel"
        @confirm="handleSelectGroupConfirm"
      ></select-group-popup>
      <!-- 组合图片选择 -->
      <select-merge-image-popup
        :show.sync="showMergeImagePopup"
        :item="selectItem"
        :selection="selection"
        @select="handleSelect"
        @clear="handleClearSelection"
        @send="handleSend"
        @preview="handlePreview"
      ></select-merge-image-popup>
      <!-- 预览图片 -->
      <preview-image-popup 
        v-if="showPreviewPopup"
        ref="previewImagePopupRef"
        :show.sync="showPreviewPopup"
        :show-index="previewIndex"
        :item="previewItem"
        :selection="selection"
        :isFirstOpen="isFirstOpen"
        pageType="collect"
        @select="handleSelect"
        @close="handleClose"
        @refresh="handlerefresh"
      ></preview-image-popup>
      <!-- 编辑文本发送 -->
      <edit-text-popup 
        :show.sync="showEditPopup"
        :item="selectItem"
        @send="handleSend"
      ></edit-text-popup>
      <!-- 预览视频 -->
      <preview-video-popup :show.sync="showPreviewVideoPopup" :item="previewItem"></preview-video-popup>
      <!-- 添加素材 -->
      <add-popup v-if="hasAuthAdd && !isCoach"></add-popup>
    </div>
    <!-- 发送预览 -->
    <preview-sent 
      v-if="showSendPreview"
      :chooseArr="selection"
      @backPage="handleSendCancel"
      @send="handleSendAll"
    ></preview-sent>

    <!-- 悬浮球 -->
    <draggable-ball :width="100" :height="34" @click="backMaterialPage">
      <template v-slot:icon>
        <van-icon :name="require('@/assets/material-library/v5.1.5/m0529-fanhui.png')"></van-icon>
      </template>
      <template v-slot:text>
        <span>返回素材库</span>
      </template>
    </draggable-ball>

    <!-- 删除确认 -->
    <confirm-popup ref="deletePopup" title="" @confirm="handleDeleteConfirm">
      <div>取消收藏内容</div>
    </confirm-popup>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch} from 'vue-property-decorator'
import {
  Col,
  Row,
  List,
  Dialog,
} from "vant";
import Helper from "@/utils/helper";
import MaterialItem from './components/material-item.vue';
import SelectTagPopup from './components/select-tag-popup.vue';
import SelectGroupPopup from './components/select-group-popup.vue';
import SelectGrid from './components/select-grid.vue';
import EmptyPage from "@/components/base/empty-page.vue";
import SelectMergeImagePopup from "./components/select-merge-image-popup.vue"
import PreviewImagePopup from "./components/preview-image-popup.vue"
import PreviewSent from "./components/preview-sent.vue";
import EditTextPopup from "./components/edit-text-popup.vue";
import PreviewVideoPopup from "./components/preview-video-popup.vue";
import AddPopup from "./components/add-popup.vue";
import DraggableBall from '@/components/material-library/draggable-ball.vue';
import ConfirmPopup from "./components/confirm-popup.vue";
import { checkAddApi, getCheckFirst, getNumApi, reportTimes, getCollectListApi, sourceMakeMediaApi,setCollectApi } from '@/api/marketing-promotion';

import wx from "weixin-js-sdk";
import wxUtils from "@/utils/wx";
import { getCosCompressUrl } from '@/utils/upload/utils';


@Component({
  components: {
    [Row.name]: Row,
    [Col.name]: Col,
    [List.name]: List,
    EmptyPage,
    SelectTagPopup,
    SelectGroupPopup,
    SelectMergeImagePopup,
    MaterialItem,
    SelectGrid,
    PreviewImagePopup,
    PreviewSent,
    EditTextPopup,
    PreviewVideoPopup,
    AddPopup,
    DraggableBall,
    ConfirmPopup
  },
  beforeRouteEnter(to, form, next) {
    next((vm:any) => {
      if (form.name === "text-library") {
        vm.active = 2;
      }
      if (form.name === "image-library") {
        vm.active = 0;
      }
      if (form.name === "video-library") {
        vm.active = 1;
      }
    });
  }
})
export default class NormalMaterial extends Vue {
  searchIcon = require('@/assets/material-library/search.png')
  /** tabs */
  tabs = [
    {
      label: '图片',
      value: 2,
      span: 12,
      gutter: 15,
      page: 1,
      isLoading: false,
      isFinished: false,
      isAllChecked: false, // 是否全选
      list: [],
      grid: 0,
      tips: '【操作提示】在预览页面长按图片可保存到本地哦～'
    },
    {
      label: '小红书',
      value: 3,
      span: 24,
      gutter: 0,
      page: 1,
      isLoading: false,
      isFinished: false,
      isAllChecked: false, // 是否全选
      list: []
    },
    {
      label: '视频',
      value: 5,
      span: 12,
      gutter: 15,
      page: 1,
      isLoading: false,
      isFinished: false,
      isAllChecked: false, // 是否全选
      list: [],
      tips: '【操作提示】在预览页面点击右上角…可保存到本地哦～'
    },
    {
      label: '文本',
      value: 1,
      span: 24,
      gutter: 0,
      page: 1,
      isLoading: false,
      isFinished: false,
      isAllChecked: false, // 是否全选
      list: []
    },
    {
      label: '官网商品',
      value: 8,
      span: 12,
      gutter: 0,
      page: 1,
      isLoading: false,
      isFinished: false,
      isAllChecked: false, // 是否全选
      list: []
    },
    {
      label: '文件',
      value: 6,
      span: 12,
      gutter: 15,
      page: 1,
      isLoading: false,
      isFinished: false,
      isAllChecked: false, // 是否全选
      list: [],
      tips: '【操作提示】在预览页面点击右上角...可分享给好友哦~'
    },
    {
      label: '小程序',
      value: 4,
      span: 12,
      gutter: 15,
      page: 1,
      isLoading: false,
      isFinished: false,
      isAllChecked: false, // 是否全选
      list: []
    }
  ]
  /** 当前tab */
  active = 0
  // 全部选中
  isAllChecked = false

  /** 筛选条件 */
  searchParams = {
    /** 关键字 */
    keyword: '',
    /** 分组ID */
    group_value: [],
    group_name: '',
    /** 分类ID */
    category_ids: []
  }

  /** 属性弹框 */
  showTagPopup = false
  /** 分组弹框 */
  showGroupPopup = false
  
  /** 角标数据 */
  numAll = {}

  /** 选中项 */
  selection: any[] = []
  maxlength = 999

  /** 选择的组合素材 */
  selectItem: any = null
  showMergeImagePopup = false

  /** 预览图片 */
  previewItem: any = null
  previewIndex: number = 0
  current = 0
  showPreviewPopup = false
  showPreviewVideoPopup = false

  /** 发送预览 */
  showSendPreview = false

  /** 编辑文本发送 */
  showEditPopup = false

  /** 是否是首次打开 */
  isFirstOpen = false
  /** 是否可以新增 */
  hasAuthAdd = false
  /** 添加 */
  showAddPopup = false

  /** 判断企业 */
  get isCoach() {
    return Helper.isCoach()
  }

  get sentData(): any {
    return this.$store.getters["materialLibrary/getNowSentData"];
  }

  /** 搜索状态 */
  get isSearch() {
    const { keyword, group_value, category_ids  } = this.searchParams
    if(keyword != '' || group_value.length > 0 || category_ids.length>0) return true
    return false
  }

  /** 缺省文案 */
  get emptyText() {
    return (item) => {
      if(this.isSearch) {
        return `没有找到相关${item.label}素材`
      }
      return `暂无${item.label}素材`
    }
  }
  
  /**
   * 是否来自消息通知
   */
  get isApplets() {
    return this.$route.query.is_applets;
  }

  activated() {
    this.selection = this.$store.getters["materialLibrary/getSelection"]
    
  }
  deactivated() {
    // 记得清除!!!!不然单页面，其它页面也可以用
    this.showSendPreview = false
  }


  @Watch('selection', {immediate: true, deep: true})
  handleSelectionChange(val) {
    console.log(val);
    
    this.$store.commit("materialLibrary/setSelection", val);
  }
  

  created() {
    this.checkAddApi()
    this.isFirst()
    this.getNumList()
  }


  /** 查询是否具备添加权限 */
  checkAddApi() {
    checkAddApi()
      .then((res: IApiResult) => {
        if (res.code === 200) {
          this.hasAuthAdd = !!res.data.auth;
        }
      })
      .catch((err) => {
        console.error("checkAddApi", err);
      });
  }
  // todo:判断是否是首次打开本页面
  async isFirst() {
    try {
      const res: any = await getCheckFirst({});
      if (res.code == 200) {
        this.isFirstOpen = res.data === 0
      }
    } catch (error) {
      console.log(error);
    }
  }
  /**
   * 获取数字角标
   */
  async getNumList() {
    const res: any = await getNumApi();
    if (res.code == 200) {
      const data = res.data.reduce((prev, item) => {
        prev[item.type] = item.count_num > 99 ? "99+" : item.count_num
        return prev;
      }, {
        time: this.$helper.dateFormat(Date.now() / 1000, "yyyy-MM-dd")
      });
      const catchData = JSON.parse(localStorage.getItem("numTabObj"));
      if (catchData && catchData.time == data.time) {
        this.numAll = catchData;
      } else {
        this.numAll = data;
        const key = this.tabs[this.active].value
        localStorage.setItem("numTabObj", JSON.stringify({...data, [key]: 0}));
      }
    }
  }

  updateCatchData() {
    const key = this.tabs[this.active].value

    const catchData = JSON.parse(localStorage.getItem("numTabObj"));
    if (catchData) {
      this.numAll = {
        ...catchData,
        [key]: 0
      }
      localStorage.setItem("numTabObj", JSON.stringify(this.numAll));
    }
  }

  /**
   * tab切换
   */
  handleTabChange() {
    this.handleSearch()
    console.log('selection', this.selection)
  }

  /** 
   * 布局切换
   */
  handleGridChange(value) {
    const { active } = this
    this.tabs[active].gutter = [15, 10][value]
    this.tabs[active].span = [12, 8][value]
  }

  /** 
   * @description 获取数据
   */
  async getData() {
    const { tabs, active, searchParams } = this
    const { value, page, list } = tabs[active]
    const params: any = {
      type: value,
      page: page,
      size: 99999,
      keyword: searchParams.keyword,
      category_ids: searchParams.category_ids
    };

    if(searchParams.group_value.length && searchParams.group_value[0] !== 0) {
      params.source_group_id = searchParams.group_value[searchParams.group_value.length -1]
    }

    const { code, message, data }: IApiResult = await getCollectListApi(params);
    this.tabs[active].isLoading = false;
    if (code !== 200) {
      this.$toast.fail(message);
      this.tabs[active].isFinished = true;
      return;
    }
    const type = this.tabs[this.active].value
    if(type === 8) {
      data.list = data.list.map((item) => {
        item.merge_image = (item.brand_source_list || []).map(image => {
          return {
            ...image,
            id: image.resource_key,
            product_name: item.product_name,
            product_code: item.product_code,
            upc_code: item.upc_code,
            is_collect: item.is_collect,
            _type: 8,
            _item_id: item.id,
          }
        })
        const currentImage = item.merge_image.find(image => image.is_main === 1) || item.merge_image[0]
        return {
          ...item,
          url: currentImage ? currentImage.url: '',
          name: currentImage ? currentImage.name : item.product_code,
          _item_type: currentImage? currentImage.type : '',
          type: 8
        };
      }) 
    }

    if(type === 2) {
      data.list.forEach(item => {
        item.merge_image.forEach(image => {
          image.is_top = item.is_top
        })
      });
    }
    this.tabs[active].list = [...list, ...data.list];
    this.tabs[active].page = page + 1;
    this.tabs[active].isFinished = data.total <= this.tabs[active].list.length;

    this.updateCatchData()
  }
 
  /**
   * 搜索
   */
  handleSearch = this.$helper.debounce(() => {
    this.handleSearchFun()
  }, 500)
  handleSearchFun() {
    this.tabs[this.active].list = []
    this.tabs[this.active].page = 1
    this.tabs[this.active].isLoading = true
    this.getData()
  }
  /**
   * 清空
   */
  handleClear() {
    this.searchParams.keyword = ''
    this.handleSearch()
  }

  /**
   * 选择分组
   */
  handleSelectGroup() {
    this.showGroupPopup = true
  }
  /**
   * 选择分组
   */
  handleSelectGroupCancel() {
    this.showGroupPopup = false
  }
  /**
   * 选择分组
   */
  handleClearGroup() {
    this.searchParams.group_value = []
    this.searchParams.group_name = ''
    this.handleSearch()
  }
  /**
   * 选择分组 确定
   */
  handleSelectGroupConfirm({value, valueName}) {
    this.searchParams.group_value = value
    this.searchParams.group_name = valueName
    this.showGroupPopup = false
    this.handleSearch()
  }

  /**
   * 选择属性
   */
  handleSelectTag() {
    this.showTagPopup = true
  }
  /**
   * 选择属性 确认
   */
  handleSelectTagConfirm(value) {
    this.searchParams.category_ids = value
    this.handleSearch()
  }

  /** ---------------------- 选择操作 --------------------------- */
  get isMergeImage() {
    return (item) => {

      const { merge_image, type } = item
      if(type === 8) {
        return merge_image.length > 0
      }
      return merge_image && merge_image.length > 1
    }
  }
  // 全选
  selectAll() {
    const tab = this.tabs[this.active];
    const isAllChecked = tab.isAllChecked;
    const { maxlength } = this;

    // 选中项集合
    let newSelection = [...this.selection];

    if (!isAllChecked) {
      // ✅ 全选当前 tab 下的内容
      for (const item of tab.list) {
        if (this.isMergeImage(item)) {
          const ids = newSelection.map(i => i.id);
          const mergedItems = item.merge_image
            .filter(val => !ids.includes(val.id))
            .map(val => ({ 
              _type: item.type, 
              type: item.type, 
              ...val 
            }));
          newSelection = [...newSelection, ...mergedItems];
        } else {
          const exists = newSelection.some(i => i.id === item.id);
          if (!exists) {
            newSelection.push(item);
          }
        }
      }

      // 限制数量
      newSelection = newSelection.slice(0, maxlength);

    } else {
      // ✅ 取消选中当前 tab 下的内容
      const currentIdsToRemove = tab.list.flatMap(item => {
        if (item.merge_image && item.merge_image.length > 1) {
          return item.merge_image.map(i => i.id);
        } else {
          return [item.id];
        }
      });

      newSelection = newSelection.filter(i => !currentIdsToRemove.includes(i.id));
    }

    this.selection = newSelection;
    tab.isAllChecked = !isAllChecked;
  }
  /**
   * 选择素材
   * @param e 
   */
  handleSelect(e) {
    console.log('e', e);
    const { status, item, isMergeImage } = e
    const { selection, maxlength } = this
    if(isMergeImage) {
      // 剩余可选择的数量
      const count = maxlength - selection.length
      if(status !== 2 && count === 0) {
        this.$toast("最多只能选择" + maxlength + "个素材")
        return
      }

      if(status === 2) {
        // （多选，取消选中）
        const ids = item.merge_image.map(val => val.id)
        // 如果是选中状态，取消选中
        this.selection = this.selection.filter(val => !ids.includes(val.id))
      } else {
        // （多选，选中）
        const ids = this.selection.map(val => val.id)
        // 如果是取消状态，选中
        const list = item.merge_image.filter(val => !ids.includes(val.id)).map(val => {
          console.log('type', item.type)
          return {
            type: item.type,
            ...val,
            _type: item.type,
            _item: item
          }
        })
        this.selection = [...this.selection, ...list].splice(0, maxlength)
      }
      
    } else {
      if(status == 0 && selection.length === maxlength) {
        this.$toast("最多只能选择" + maxlength + "个素材")
        return
      }
      if(status == 2) {
        // （单个，取消选中）
        this.selection = this.selection.filter(val => val.id !== item.id)
      } else {
        // （单个，选中）
        this.selection.push(item)

        console.log('list2', this.selection)
      }
    }

    // 全选状态，取消
    this.tabs[this.active].isAllChecked = this.checkIsAllChecked(this.tabs[this.active], this.selection);
  }
  // 全选状态
  checkIsAllChecked(tab, selection) {
    if(tab.list.length === 0) return false
    const idsInSelection = selection.map(i => i.id);
    const tabIds = tab.list.flatMap(item =>
      item.merge_image && item.merge_image.length > 1
        ? item.merge_image.map(i => i.id)
        : [item.id]
    );

    return tabIds.every(id => idsInSelection.includes(id));
  }
  /** 取消所有全选 */
  clearAllChecked() {
    // 取消全选
    this.tabs = this.tabs.map(item => {
      return {
        ...item,
        isAllChecked: false
      }
    })
  }
  // 取消收藏
  handleCancelCollect(item) {
    console.log('cancel item', item)

    this.tabs[this.active].list = this.tabs[this.active].list.filter(val => val.id!== item.id)
    // 如果是选中状态，取消选中
    this.selection = this.selection.filter(val => val.id !== item.id)

  }
  /** 
   * 选择组合素材
   */
  handleSelectImage(item) {
    this.selectItem = item
    this.showMergeImagePopup = true
  }
  /** 
   * 清空选择
   */
  handleClearSelection() {
    this.selection = [];
    this.clearAllChecked();
  }
 /** ---------------------- 选择操作 end --------------------------- */
  /**
   * 预览图片
   * @param item 
   */
  handlePreview({item, index}, current=0) {
    this.previewItem = item
    this.previewIndex = index
    this.current = current
    this.showPreviewPopup = true
  }

  handlePreviewVideo(item) {
    this.previewItem = item
    this.showPreviewVideoPopup = true
  }

  /** 图片关闭 */
  handleClose() {
    let type = this.tabs[this.active].value
    if(type !== 8) {
      // this.handleSearch()
    }
  }
  handlerefresh({index, item, isMergeImage}) {
    
    if(item.type === 2 && item.merge_image.length) {
      item.merge_image = item.merge_image.filter(v=> v.is_collect)
      if(item.merge_image.length === 0) {
        this.showPreviewPopup = false

      } else {
        if(index === item.merge_image.length) {
          index = item.merge_image.length-1;
        }
        item.is_collect = 1;
        (this.$refs.previewImagePopupRef as any).handleSwipeTo(index);
      }
    } else {
      this.showPreviewPopup = false
      this.handleSearch()
    }
    this.handleSearch()
    console.log('this.previewItem', this.previewItem)
  }

  /**
   * 编辑文本
   */
  handleEditText(item) {
    this.selectItem = item
    this.showEditPopup = true
  }

  /**
   * 发送预览
   */
  handleSend(item?:any) {
    console.log('item', item, this.selection)
    if((!item && this.selection.length==0) || this.selection.length > 9 ) return
    if(item || this.selection.length==1) {
      this.sendChatMessage(item||this.selection[0])
    } else {
      this.showSendPreview = true
    }
  }

  /**
   * 取消发送
   */
  handleSendCancel() {
    this.showSendPreview = false
  }

  /**
   * 发送所有选中数据
   */
  async handleSendAll() {
    for (let i = 0; i < this.selection.length; i++) {
      const item = this.selection[i]
      await this.sendChatMessage(item, i==this.selection.length-1)
    }
  }

  /**
   * 创建消息
   * @param item 
   */
  async createMessage(item: any) {
    console.log(item);
    
    const msgtype = {
      1: 'text',
      2: 'image',
      3: 'news',
      4: 'miniprogram',
      5: 'video',
      6: 'file',
    }[item.type]
    const message: any = {
      msgtype
    }
    // 获取mediaid
    if (["video", "image", "file"].includes(msgtype)) {
      if (!item.media_id || Date.now() > item.media_expire_time * 1e3) {
        const params: any = {
          source_type: msgtype,
          source_url: msgtype === 'image' ? getCosCompressUrl({ url: item.url, maxWidth: 1080, maxHeight: 1440 }) : item.url,
          source_name: item.name,
        };
        if(item._type === 8) {
          params.resource_key = item.id
        } else {
          params.source_id = item.id
        }

        const { code, data }: IApiResult = await sourceMakeMediaApi(params);
        if (code == 200) {
          item.media_id = data.media_id;
        }
      }
    }

    switch (msgtype) {
      case "text":
        message.text = {
          content: item.content,
        };
        break;
      case "image":
        message.image = { mediaid: item.media_id };
        break;
      case "news":
        message.news = {
          link: item.link_url,
          title: item.title,
          desc: item.description,
          imgUrl: getCosCompressUrl({ url: item.cover_url, maxWidth: 1080, maxHeight: 1440 }),
        }
        break;
      case "video":
        message.video = {
          mediaid: item.media_id,
        };
        break;
      case "miniprogram":
        message.miniprogram = {
          appid: item.app_id,
          title: item.title,
          imgUrl: getCosCompressUrl({ url: item.cover_url, maxWidth: 1080, maxHeight: 1440 }),
          page: item.path,
        };
        break;
      default:
        message.file = { mediaid: item.media_id };
        break;
    }
    return message
  }

  /**
   * 发送消息
   * @param item 
   */
  async sendChatMessage(item: any, enterChat = true) {
    console.log(item);
    
    this.$toast.loading({
      message: "生成中...",
      forbidClick: true,
    });
    const message = await this.createMessage(item)
    message.enterChat = enterChat;
    return new Promise((resolve, reject) => {
      window.wx.invoke("sendChatMessage", message, (res) => {
        if (res.err_msg == "sendChatMessage:ok") {
          //发送成功
          console.log("sendChatMessage:成功回调");
          // 记录本次已发送过的文件
          this.$store.commit("materialLibrary/setNowSentData", [...this.sentData, item]);
          // 神策数据上报 -- 发送素材
          this.$track.trackAssetSend({
            asset_id: item.id,
            customer_id: this.$store.getters["user/getExternalId"],
            task_id:"普通发送无法获取",
            operation_type:2,
            is_official: item._type==8?2:1,
            asset_type_modify: item.type,
            official_url: item._type==8 ? item.url: '非官网商品素材',
            official_upc: item._type==8 ? item.upc_code: '非官网商品素材',
            official_product: item._type==8 ? item.product_code: '非官网商品素材',
            official_name: item._type==8 ? item.product_name: '非官网商品素材',
            current_page: this.$route.path,
            entry_point: 2,
            library_type: 2,
            is_top: item._type==2&&item._item.merge_image.length>1 ? !!item._item.is_top : !!item.is_top,
            is_favorite: item._type==8? !!item._item.is_collect : !!item.is_collect
          })
          // todo: 上报次数
          reportTimes({
            ids: [item.id],
            type: 2
          })
          resolve(res)
        } else {
          if (res.err_msg == "sendChatMessage:cancel") {
            return;
          }
          if (message.msgtype == "miniprogram") {
            this.$toast.fail("企微官方后台未关联小程序");
          }
          if (res.err_code === 1) {
            this.$toast.fail('发送失败,请检查素材');
          }
          reject(res)
        }
        // 关闭文本编辑弹窗
        this.showEditPopup = false;
        this.$toast.clear();
      });
    })
  }

  handleAdd() {
    this.showAddPopup = true
  }

 /** ---------------------- 收藏 --------------------------- */

  // 取消收藏
  handleCollect() {
    if(this.selection.length == 0) return

    (this.$refs.deletePopup as any).showPopup()

    // Dialog.confirm({
    //   // title: '提示',
    //   message: '取消收藏内容',
    // })
    //   .then(() => {
    //     // 点击确认后的逻辑
    //     console.log('用户确认')
    //     this.setSollectFn();
    //   })
    //   .catch(() => {
    //     // 点击取消或关闭对话框的逻辑
    //     console.log('用户取消')
    //   })
  }
  /** 设置收藏 */
  handleDeleteConfirm() {
    let tempObj = {};
    const selection = JSON.parse(JSON.stringify(this.selection));
    
    selection.forEach(item => {
      let type = item._type || item.type;
      // 素材类别，8是官网素材，1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件 7 朋友圈组合
      // console.log(item);
      console.log('_type', item._type);
      // console.log('type', item.type);
      
      
      if([8].includes(item._type)) {
        // 官网素材
        if(!tempObj['type' + type]) {
          // 没有属性，先创建
          tempObj['type' + type] = {
            type: item._type,
            product_code: []
          }
        }
        // 有属性了，直接 push
        if(tempObj['type' + type].product_code.indexOf(item.product_code) == -1) {
          tempObj['type' + type].product_code.push(item.product_code)
        }
      } else {
        if(!tempObj['type' + type]) {
          // 没有属性，先创建
          tempObj['type' + type] = {
            type: item._type || item.type,
            source_id: []
          }
        }
        // 有属性了，直接 push
        tempObj['type' + type].source_id.push(item.id)
      }
    })

    let params: any = {
      source_info: Object.values(tempObj),
      is_collect: 2, // 1收藏，2取消收藏
    };
    console.log(params);
    // return 
    setCollectApi(params).then((res: IApiResult) => {
      console.log(res);
      if (res.code == 200) {
        this.handleSearch();
        this.$toast("已取消收藏")
        this.selection = [];

        // 取消全选
        this.clearAllChecked();
        console.log(selection);
        
        const trackList = selection.map(item => {
          if (item._type === 8) {
            return {
              is_top: false,
              asset_id: item._item_id,
              is_official: 2,
              asset_type_modify: item._type,
              official_url: item.url,
              official_upc: item.upc_code,
              official_product: item.product_name,
              official_name: item.product_code
            }
          }
          return {
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item._type || item.type
          }
        })
        

        Object.values(trackList.reduce((acc, cur) => {
          acc[cur.asset_id] = cur
          return acc
        }, {})).forEach((item:any) => {
          this.$track.trackAssetFavorite({
            current_page: this.$route.path,
            entry_point: 2,
            library_type: 2,
            operation_position: 1, // 固定值
            operation_type: 2,
            official_url: '非官网商品素材',
            official_upc: '非官网商品素材',
            official_product: '非官网商品素材',
            official_name: '非官网商品素材',
            ...item
          })
        })
      } 
    })
  }

  /** 去素材*/
  backMaterialPage() {
    // 神策数据上报 -- 点击收藏入口
    this.$track.trackMyFavoriteEntranceClick({
      operation_type: '返回素材列表',
      entry_point: '1V1聊天窗口',
      library_type: '素材库'
    })
    this.showSendPreview = false
    this.$router.back()
  }

   // 监听路由变化（刷新逻辑）
  @Watch('$route')
  onRouteChange() {
    if (this.$route.name === 'my-collect') {
      this.handleSearch();
      this.selection = [];
      this.clearAllChecked();
    }
  }

  async mounted() {
    // 刷新接口
    this.handleSearch()

    try {
      wxUtils.wxAgentConfig();
      // this.getWxContent();
    } catch (error) {
      console.log(error);
    }
  }

  
  getWxContent() {
    const self = this;
    if (!window.wx.invoke) return;
    window.wx.invoke("getContext", {}, function (res) {
      if (res.err_msg == "getContext:ok") {
        console.log("getContext", res.entry);
        const entry = res.entry;
        // 从单聊会话的工具栏进入-打开客户画像-客户页面
        if (entry == "single_chat_tools") {
          window.wx.invoke("getCurExternalContact", {}, (res) => {
            if (res.err_msg == "getCurExternalContact:ok") {
              console.log("customer:getCurExternalContact", res.userId);
              self.$store.commit("user/setExternalId", res.userId);
              // 神策数据上报 -- 点击侧边栏
              self.$track.trackClickItem({
                feature: self.$route.path,
                customer_id: res.userId
              })
            } else {
              //错误处理
              console.log("getCurExternalContacterr:" + JSON.stringify(res));
            }
          });
        } else if (entry == "group_chat_tools") {
          // 从群聊会话的工具栏进入-打开客户画像-社群页面
          window.wx.invoke("getCurExternalChat", {}, function (res) {
            if (res.err_msg == "getCurExternalChat:ok") {
              console.log("res.chatId", res.chatId);
              self.$store.commit("common/setGroupChatId", res.chatId);
            } else {
              //错误处理
              console.log("getCurExternalChaterr:" + JSON.stringify(res));
            }
          });
        }
      } else {
        //错误处理
        console.log("getContexterr:" + JSON.stringify(res));
      }
    });
  }
}
</script>
<style lang="scss" scoped>
::v-deep .van-tabs__wrap {
  height: 40px !important;
}
::v-deep .van-tab__pane {
  min-height: calc(100vh - 120px);
  padding-bottom: 100px;
}
::v-deep .van-search {
  padding: 0;
  height: 100%;
  .van-search__content {
    height: 100%;
    padding-left: 10px;
    background-color: transparent;
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
.van-list {
  padding: 15px;
}
.container {
  background: #f4f7fb;

  .header {
    display: flex;
    padding: 15px 15px 0;

    .search {
      flex: 1;
      background-color: #fff;
      box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
      height: 50px;
      padding: 0 15px;
      display: flex;
      .search-tag {
        height: 100%;
        display: flex;
        align-items: center;
        font-weight: 400;
        font-size: 13px;
        color: #000000;
        position: relative;
        padding-right: 10px;

        &::after {
          content: "";
          width: 8px;
          height: 4px;
          background: url('~@/assets/ico_default.png') center/100% no-repeat;
          margin-left: 5px;
        }

        &::before {
          content: "";
          width: 1px;
          height: 16px;
          background-color: #EEEEEE;
          position: absolute;
          right: 0;
          top: 50%;
          transform: translateY(-50%);
        }

        .count {
          position: absolute;
          top: 10px;
          right: 10px;
          font-size: 10px;
          width: 12px;
          height: 12px;
          border-radius: 50%;
          background-color: #000;
          color: #fff;
          display: flex;
          align-items: center;
          justify-content: center;
        }
      }

      .search-inner {
        flex: 1;
        height: 100%;
        position: relative;

        .delete {
          width: 30px;
          height: 50px;
          position: absolute;
          right: -15px;
          top: 50%;
          transform: translateY(-50%);
          background: url("~@/assets/material-library/delete.png") center/10px no-repeat;
        }
      }
    }

    .filter {
      width: 50px;
      height: 50px;
      background: #ffffff url('~@/assets/material-library/screen.png') center/16px no-repeat;
      box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
    }
  }
  .group-name {
    padding: 0px 24px 0px 16px;
    margin: 10px 0 6px 0;
    font-size: 12px;
    color: #ffffff;
    line-height: 18px;
    
    .text {
      display: inline-block;
      padding: 5px 15px;
      background: #4a4a4a;
      border-radius: 24px;
      position: relative;
      
    }
    .delete {
      width: 18px;
      height: 18px;
      position: absolute;
      right: -9px;
      top: 50%;
      margin-top: -9px;
      background: url('~@/assets/material-library/delete-icon.png') center/100% no-repeat;
    }
  }
  .tab-text {
    position: relative;
    .num-icon-box {
      border-radius: 50%;
      width: 14px;
      height: 14px;
      position: absolute;
      top: -1px;
      right: -14px;
      background: #000000;
      font-size: 12px;
      font-weight: 400;
      color: #ffffff;
      display: flex;
      justify-content: center;
      align-items: center;
    
      .num-text {
        transform: scale(0.6);
      }
    }
  }
  
  .tips {
    font-size: 12px;
    font-weight: 400;
    color: #a6a6a6;
    line-height: 17px;
    text-align: center;
    margin-top: 15px;
  }

  .footer {
    position: fixed;
    bottom: -1px;
    width: 100%;
    height: 50px;
    left: 0;
    padding: 8px;
    background-color: #ffffff;
    font-size: 14px;
    z-index: 9;
    box-shadow: 0px 1px 6px 0px rgba(0, 0, 0, 0.15);
    display: flex;
    align-items: center;
    .num-left {
      flex: 1;
      div {
        margin: 0 4px;
        font-size: 16px;
      }

      .max-num-left {
        font-size: 10px;
        color: #E02020;
      }
    }

    .clear-btn {
      font-size: 12px;
      font-weight: 500;
      color: #000000;
      text-decoration: solid underline #000000;
      margin-right: 15px;
    }

    .button-right {
      width: 88px;
      height: 34px;
      color: #ffffff;
      line-height: 34px;
      text-align: center;
      background-color: #000;

      &.dis-btn {
        background-color: #cccccc;
      }
    }
    .button-right+.button-right{
      margin-left: 10px;
    }
  }

}
// 全选
.select-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px 0;
  font-size: 12px;
  color: #000000;
  line-height: 18px;
  .select-right{
    display: flex;
    align-items: center;
    .check-box{
      display: inline-block;
      width: 15px;
      height: 15px;
      margin-right: 5px;
      background: url('~@/assets/material-library/v5.1.5/m0529-check.png') no-repeat center center/15px 15px;
    }
    .checked-box{
      background: url('~@/assets/material-library/v5.1.5/m0529-checked.png') no-repeat center center/15px 15px;
    }
  }
}

</style>
