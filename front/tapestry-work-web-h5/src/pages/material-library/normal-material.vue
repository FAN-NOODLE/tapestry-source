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
              placeholder="请输入素材关键词"
              @input="handleSearch"
            />
            <div v-if="searchParams.keyword" class="delete" @click="handleClear"></div>
          </div>
        </div>
        <div class="filter" @click="handleSelectGroup"></div>
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
          <select-grid v-if="tab.value==2 && isCoach" v-model="tab.grid" :showFirstTips="isFirstOpen" @change="handleGridChange"></select-grid>
          <van-list 
            v-model="tab.isLoading"
            :finished="tab.isFinished"
            :finished-text="tab.list.length ? '没有更多了' : ''"
            @load="getData" 
          >
            <van-row :gutter="tab.gutter">
              <van-col v-for="(item, index) in tab.list" :key="index" :span="tab.span">
                <material-item 
                  :item="item"
                  :selection="selection"
                  :small="tab.grid === 1"
                  @select="handleSelect"
                  @send="handleSend"
                  @preview="handlePreview"
                  @preview-video="handlePreviewVideo"
                  @selectImage="handleSelectImage"
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
      <div class="footer" v-if="!isApplets">
        <div class="num-left">
          <div>已选：{{ selection.length }}</div>
          <div
            class="max-num-left"
            v-if="selection.length == 9"
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
        <div
          class="button-right"
          :class="[!selection.length ? 'dis-btn' : '']"
          @click="handleSend()"
        >
          {{ selection.length > 1 ? '组合发送' : '发送' }}
        </div>
      </div>
      <!-- 选择属性 -->
      <select-tag-popup 
        :show.sync="showTagPopup" 
        :value="searchParams.category_ids"
        @confirm="handleSelectTagConfirm"
      ></select-tag-popup>
      <!-- 选择分组 -->
      <select-group-popup 
        :show.sync="showGroupPopup"
        :value="searchParams.group_value"
        text="全部分组"
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
        v-if="previewItem"
        :show.sync="showPreviewPopup"
        :show-index="previewIndex"
        :item="previewItem"
        :selection="selection"
        :isFirstOpen="isFirstOpen"
        @select="handleSelect"
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
      <add-popup v-if="hasAuthAdd && isKS"></add-popup>
    </div>
    <!-- 发送预览 -->
    <preview-sent 
      v-if="showSendPreview"
      :chooseArr="selection"
      @backPage="handleSendCancel"
      @send="handleSendAll"
    ></preview-sent>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator'
import {
  Col,
  Row,
  List,
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
import { checkAddApi, getCheckFirst, getNumApi, reportTimes, selectListApi, sourceMakeMediaApi } from '@/api/marketing-promotion';
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
    AddPopup
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
      list: [],
      grid: 0,
      tips: '【操作提示】在预览页面长按图片可保存到本地哦～'
    },
    {
      label: '视频',
      value: 5,
      span: 12,
      gutter: 15,
      page: 1,
      isLoading: false,
      isFinished: false,
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
      list: [],
      tips: '【操作提示】在预览页面点击右上角...可分享给好友哦~'
    },
    {
      label: '链接',
      value: 3,
      span: 24,
      gutter: 0,
      page: 1,
      isLoading: false,
      isFinished: false,
      list: []
    },
    {
      label: '小程序',
      value: 4,
      span: 12,
      gutter: 15,
      page: 1,
      isLoading: false,
      isFinished: false,
      list: []
    }
  ]
  /** 当前tab */
  active = 0

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
  showGroupPopup = true
  
  /** 角标数据 */
  numAll = {}

  /** 选中项 */
  selection: any[] = []
  maxlength = 9

  /** 选择的组合素材 */
  selectItem: any = null
  showMergeImagePopup = false

  /** 预览图片 */
  previewItem: any = null
  previewIndex: number = 0
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
  get isKS() {
    return Helper.isKS()
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
      size: 15,
      keyword: searchParams.keyword,
      category_ids: searchParams.category_ids
    };

    if(searchParams.group_value.length && searchParams.group_value[0] !== 0) {
      params.source_group_id = searchParams.group_value[searchParams.group_value.length -1]
    }

    const { code, message, data }: IApiResult = await selectListApi(params);
    this.tabs[active].isLoading = false;
    if (code !== 200) {
      this.$toast.fail(message);
      this.tabs[active].isFinished = true;
      return;
    }

    this.tabs[active].list = [...list, ...data.list];
    this.tabs[active].page = page + 1;
    this.tabs[active].isFinished = data.total <= this.tabs[active].list.length;

    this.updateCatchData()
  }
 
  /**
   * 搜索
   */
  handleSearch() {
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

  /**
   * 选择素材
   * @param e 
   */
  handleSelect(e) {
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
        const ids = item.merge_image.map(val => val.id)
        // 如果是选中状态，取消选中
        this.selection = this.selection.filter(val => !ids.includes(val.id))
      } else {
        const ids = this.selection.map(val => val.id)
        // 如果是取消状态，选中
        const list = item.merge_image.filter(val => !ids.includes(val.id)).map(val => {
          return {
            type: item.type,
            ...val,
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
        const index = selection.findIndex(val => val.id === item.id)
        this.selection.splice(index, 1)
      } else {
        this.selection.push(item)
      }
    }
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
    this.selection = []
  }

  /**
   * 预览图片
   * @param item 
   */
  handlePreview({item, index}) {
    this.previewItem = item
    this.previewIndex = index
    this.showPreviewPopup = true
  }

  handlePreviewVideo(item) {
    this.previewItem = item
    this.showPreviewVideoPopup = true
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
    if(!item && this.selection.length==0 ) return
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
        const params: IMakeMediaParams = {
          source_id: item.id,
          source_type: msgtype,
          source_url: msgtype === 'image' ? getCosCompressUrl({ url: item.url, maxWidth: 1080, maxHeight: 1440 }) : item.url,
          source_name: item.name,
        };
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
            is_official: 1,
            asset_type_modify: item.type,
            official_url: '非官网商品素材',
            official_upc: '非官网商品素材',
            official_product: '非官网商品素材',
            official_name: '非官网商品素材',
            current_page: this.$route.path,
            entry_point: 2,
            library_type: 2,
            is_top: !!item.is_top,
            is_favorite: !!item.is_collect
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

}
</script>
<style lang="scss" scoped>

::v-deep .van-tabs__wrap {
  height: 40px !important;
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
      margin-right: 10px;
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
      width: 165px;
      height: 34px;
      color: #ffffff;
      line-height: 34px;
      text-align: center;
      background-color: #000;

      &.dis-btn {
        background-color: #cccccc;
      }
    }
  }

}
</style>