<template>
  <div>
    <div v-if="!showSendPreview" class="container">
      <div v-if="!searchStatus" class="goods-material">
        <div class="goods-material__header">
          <div class="title">官网商品素材查询</div>
        </div>
        <div class="goods-material__content">
          <div class="tab-wrap">
            <van-tabs 
              v-model="tabs.active"
              title-inactive-color="#000000"
              color="#000000"
              line-width="20px"
              line-height="2px"
              @change="handleTabChange">
              <van-tab v-for="(item, index) in tabs.list" :title="item.title" :key="index"></van-tab>
            </van-tabs>
          </div>
          
          <!-- 商品信息查询 -->
          <div v-if="tabs.active===0" class="search-wrap">
            <div class="search-type" @click="handleTypeChange">{{types.list[types.active].title}}</div>
            <input ref="input1" class="search-input" type="text" :placeholder="`请输入${types.list[types.active].placeholder}查询`" contenteditable v-model="searchValue" @click="handleInputClick" />
            <div v-if="searchValue.length" class="icon-clear" @click="handleClear"></div>
          </div>
          <!-- 条形码查询 -->
          <div v-if="tabs.active===1" class="search-wrap">
            <input class="search-input" type="text" placeholder="请输入条形码查询相关数据" contenteditable v-model="searchValue" />
            <div v-if="searchValue.length" class="icon-clear" @click="handleClear"></div>
            <div class="scan" @click="handleScan"></div>
          </div>
          <div :class="['btn', {disabled: !searchValue.length}]" @click="handleSearch">立即查询</div>
        </div>
      </div>
      <!-- 搜索结果 -->
      <div v-else class="goods-material__result">
        <div class="header">
          <!-- 商品信息查询 -->
          <div v-if="tabs.active===0" class="search-wrap">
            <div class="search-type" @click="handleTypeChange">{{types.list[types.active].title}}</div>
            <input class="search-input" type="text" :placeholder="`请输入${types.list[types.active].placeholder}查询`" contenteditable v-model="searchValue" />
            <div v-if="searchValue.length" class="icon-clear" @click="handleClear"></div>
            <div :class="['search-btn', {disabled: !searchValue.length}]" @click="handleSearch">搜索</div>
          </div>
          <!-- 条形码查询 -->
          <div v-if="tabs.active===1" class="search-wrap">
            <input class="search-input" type="text" placeholder="请输入条形码查询相关数据" contenteditable v-model="searchValue" />
            <div v-if="searchValue.length" class="icon-clear" @click="handleClear"></div>
            <div class="scan" @click="handleScan"></div>
            <div :class="['search-btn', {disabled: !searchValue.length}]" @click="handleSearch">搜索</div>
          </div>
        </div>
        <select-grid v-model="grid" @change="handleGridChange"></select-grid>
        <van-list 
          v-model="isLoading"
          :finished="isFinished"
          :finished-text="list.length ? '没有更多了' : ''"
          @load="getData" 
        >
          <van-row :gutter="gutter">
            <van-col v-for="(item, index) in list" :key="index" :span="span">
              <material-item 
                :item="item"
                :selection="selection"
                :small="grid === 1"
                @select="handleSelect"
                @preview="handlePreview"
                @selectImage="handleSelectImage"
                @poster="handlePoster($event, item)"
              ></material-item>
            </van-col>
          </van-row>
        </van-list>
        <empty-page
          v-if="!list.length && !isLoading"
          picture="noData"
          paddingTop="60"
          paddingBottom="265"
          :text="'暂无素材'"
        ></empty-page>
  
        <!-- footer -->
        <div class="footer">
          <div class="num-left">
            <div>已选：{{ selection.length }}</div>
            <div
              class="max-num-left"
              v-if="selection.length == maxlength"
            >最多可选{{maxlength}}个</div>
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
      </div>
      <!-- 选择筛选类型 -->
      <select-search-popup 
        :show.sync="types.show" 
        :list="types.list" 
        :active="types.active" 
        @confirm="handleTypeChangeConfirm"
      ></select-search-popup>
  
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
        :show.sync="showPreviewPopup"
        :show-index="previewIndex"
        :item="previewItem"
        :selection="selection"
        :isFirstOpen="false"
        @select="handleSelect"
      ></preview-image-popup>
    </div>
    <!-- 发送预览 -->
    <preview-sent 
      v-if="showSendPreview"
      :chooseArr="selection"
      :is_official="2"
      @backPage="handleSendCancel"
      @send="handleSendAll"
    ></preview-sent>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import {
  Col,
  Row,
  List,
} from "vant";
import MaterialItem from './components/material-item.vue';
import SelectGrid from './components/select-grid.vue';
import EmptyPage from "@/components/base/empty-page.vue";
import SelectMergeImagePopup from "./components/select-merge-image-popup.vue"
import PreviewImagePopup from "./components/preview-image-popup.vue"
import PreviewSent from "./components/preview-sent.vue";
import SelectSearchPopup from "./components/select-search-popup.vue";

import { getBrandCollectMaterialListApi, reportTimes, sourceMakeMediaApi } from '@/api/marketing-promotion';
import { getCosCompressUrl } from '@/utils/upload/utils';

@Component({
  name: 'GoodsMaterial',
  components: {
    [Row.name]: Row,
    [Col.name]: Col,
    [List.name]: List,
    MaterialItem,
    SelectGrid,
    EmptyPage,
    PreviewImagePopup,
    PreviewSent,
    SelectMergeImagePopup,
    SelectSearchPopup
  }
})
export default class GoodsMaterial extends Vue {
  /** 查询方式切换 */
  tabs = {
    /** 下标  */
    active: 0,
    /** 列表 */
    list: [
      {
        title: '商品信息查询',
        value: 1
      },
      {
        title: '条形码查询',
        value: 'upc_code'
      }
    ]
  }
  types = {
    list: [ 
      {
        title: 'Skuitem',
        value: 'sku_item',
        placeholder: '商品Skuitem'
      }, 
      {
        title: '商品名称',
        value: 'product_name',
        placeholder: '商品名称'
      } 
    ],
    active: 0,
    show: false
  }
  /** 搜索值 */
  searchValue = ''
  /** 搜索状态 */
  searchStatus = false
  /** 布局参数 */
  grid = 0
  span = 12
  gutter = 15

  /** 分页标记 */
  last_code = ''
  /** 加载状态 */
  isLoading = false
  /** 是否加载完成 */
  isFinished = false
  
  /** 列表数据 */
  list = []
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

  get searchParams() {
    let key = this.tabs.list[this.tabs.active].value
    if(key === 1) {
      key = this.types.list[this.types.active].value
    }
    return {
      [key]: this.searchValue
    }
  }
  
  activated() {
    this.selection = this.$store.getters["materialLibrary/getSelection"]
  }

  @Watch('selection', {immediate: true, deep: true})
  handleSelectionChange(val) {
    console.log(val);
    
    this.$store.commit("materialLibrary/setSelection", val);
  }
  
  
  handleTabChange() {
    this.searchValue = ''
  }
  /**
   * 搜索类型 切换
   */
  handleTypeChange() {
    this.types.show = true
  }
  /**
   * 搜索类型 确定
   * @param value 
   */
  handleTypeChangeConfirm(value) {
    this.types.active = value
    this.searchValue = ''
  }
  /**
   * 搜索
   */
  handleSearch() {
    if(!this.searchValue.trim()) return
    this.list = []
    this.last_code = ''
    this.searchStatus = true
    this.isLoading = true
    this.getData()
    // 神策数据上报 -- 查询官网商品素材
    let key = this.tabs.list[this.tabs.active].value
    if(key === 1) {
      key = this.types.list[this.types.active].value
    }
    this.$track.trackAssetInquiry({
      inquiry_type: key
    })
  }
  /**
   * 清空搜索
   */
  handleClear() {
    this.searchValue = ''
  }
  /**
   * 点击获取焦点
   */
  handleInputClick() {
    (this.$refs.input1 as any).focus()
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
      console.log(res);
      
      try {
        this.searchValue = JSON.parse(res.resultStr).scan_code.scan_result
      } catch (error) {
        this.searchValue = res.resultStr
      }
      this.handleSearch()
    })
  }

  /** 
   * 布局切换
   */
  handleGridChange(value) {
    this.gutter = [15, 10][value]
    this.span = [12, 8][value]
  }
  /**
   * 获取数据
   */
  async getData() {
    const { last_code, searchParams, list } = this
    const params: any = {
      last_code,
      size: 15,
      ...searchParams
    };

    const { code, message, data = [] }: IApiResult = await getBrandCollectMaterialListApi(params);
    this.isLoading = false;
    if (code !== 200) {
      this.$toast.fail(message);
      this.isFinished = true;
      return;
    }

    this.list = [...list, ...data.map((item: any) => {
      item.merge_image = (item.brand_source_list || []).map(image => {
        return {
          ...image,
          id: image.resource_key,
          product_name: item.product_name,
          product_code: item.product_code,
          upc_code: item.upc_code
        }
      })
      // const video = item.merge_image.find(image => image.type === 5)
      // console.log(video);
      
      // if(video) {
      //   video.poster = ''
      // }
      // const imageList = item.merge_image.filter(image => image.type === 2)
      // item.merge_image = video ? [video, ...imageList] : imageList
      const currentImage = item.merge_image.find(image => image.is_main === 1) || item.merge_image[0]
      return {
        ...item,
        url: currentImage ? currentImage.url: '',
        name: currentImage ? currentImage.name : item.product_code,
        _item_type: currentImage? currentImage.type : '',
        type: 8
      };
    })];
    this.last_code = data.length ? data[data.length - 1].id : '';
    this.isFinished = data.length < params.size;
  }

  handlePoster(value: string, item: any) {
    item.merge_image[0].poster = value
    this.$set(item, 'merge_image[0].poster', value)
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
        const params: any = {
          resource_key: item.id,
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
          // 神策数据上报 -- 发送素材
          this.$track.trackAssetSend({
            asset_id: item.id,
            customer_id: this.$store.getters["user/getExternalId"],
            task_id:"普通发送无法获取",
            operation_type: 2,
            is_official: 2,
            asset_type_modify: item.type,
            official_url: item.url,
            official_upc: item.upc_code,
            official_product: item.product_name,
            official_name: item.product_code,
            current_page: this.$route.path,
            entry_point: 2,
            library_type: 2,
            is_top: false,
            is_favorite: !!item.is_collect
          })
          // todo: 上报次数
          reportTimes({
            resource_keys: [item.id],
            type: 2 // 1: 下载次数 2: 转发次数 3: 复制次数（文本）
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
        this.$toast.clear();
      });
    })
  }
}
</script>
<style lang="scss" scoped>
::v-deep .van-tabs--line .van-tabs__wrap {
  height: 30px;
}
::v-deep .van-tab {
  font-size: 12px;
}

::v-deep .van-field__body {
  height: 100%;

  input {
    height: 100%;
  }
}
.icon-clear {
  width: 35px;
  height: 35px;
  background: url("~@/assets/material-library/v4.2.0/icon_clear.png") center/20px no-repeat;
}

.placeholder {
  color: #96979E;
}

.goods-material {
  
  &__header {
    background: url('~@/assets/material-library/v4.2.0/bg_1.png') center/100% no-repeat;
    width: 375px;
    height: 240px;
    overflow: hidden;

    .title {
      font-weight: 500;
      font-size: 16px;
      color: #000000;
      line-height: 22px;
      margin-top: 96px;
      text-align: center;
    }
  }

  &__content {
    width: 345px;
    height: 203px;
    background: #FFFFFF;
    margin: auto;
    margin-top: -70px;
    z-index: 1;
    padding: 0 15px;
    box-sizing: border-box;

    .tab-wrap {
      height: 46px;
      position: relative;
      padding: 5px 15px 0;
      box-sizing: border-box;

      &::after {
        content: "";
        width: 315px;
        height: 1px;
        background-color: #EEEEEE;
        position: absolute;
        bottom: 0;
        left: 0;
      }
    }

    .search-wrap {
      width: 285px;
      height: 36px;
      border: 1px solid #000000;
      margin: auto;
      margin-top: 24px;
      box-sizing: border-box;
      display: flex;
      align-items: center;

      .search-type {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100px;
        height: 100%;
        position: relative;
        font-size: 13px;
        white-space: nowrap;

        &::after {
          content: "";
          width: 8px;
          height: 4px;
          background: url("~@/assets/material-library/v4.2.0/icon_down.png") center/100% no-repeat;
          margin-left: 5px;
        }

        &::before {
          content: "";
          width: 1px;
          height: 20px;
          background-color: #CCCCCC;
          position: absolute;
          right: 0;
          top: 50%;
          transform: translateY(-50%);
        }
      }

      .search-input {
        flex: 1;
        height: 100%;
        padding: 0 0 0 18px;
        font-size: 12px;
      }

      .scan {
        width: 40px;
        height: 40px;
        background: url("~@/assets/material-library/v4.2.0/icon_scan.png") center/18px no-repeat;
          
      }
    }

    .btn {
      width: 285px;
      height: 36px;
      margin: auto;
      margin-top: 30px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      background: #000000;
      font-size: 14px;

      &.disabled {
        background: #CCCCCC;
      }
    }
  }

  &__footer {
    padding: 16px 15px;
    display: flex;
    align-items: center;
    box-shadow: 0px 1px 6px 0px rgba(0,0,0,0.15);
    box-sizing: border-box;
    
    .number {
      font-size: 16px;
      font-weight: 400;
      color: #000000;
      line-height: 22px   ;
      flex: 1;
      position: relative;
      .tips {
        font-size: 10px;
        font-weight: 400;
        color: #E02020;
        line-height: 14px
      }
    }
    .btn-clear {
      font-size: 12px;
      font-weight: bold;
      margin-right: 20px;
      text-decoration: underline;
    }
    .btn {
      width: 330px;
      height: 34px;
      background: #000;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      font-weight: bold;
      color: #fff;

      &.disabled {
        background-color: #ccc
      }
    }
  }
}
.goods-material__result {
  background: #f4f7fb;
  .header {
    padding-top: 15px;
    position: sticky;
    z-index: 10;
    top: 50px;
    background: #f4f7fb;
  }

  .search-wrap {
    width: 345px;
    height: 50px;
    background: #FFFFFF;
    box-shadow: 0px 6px 12px 0px rgba(222,224,227,0.5);
    display: flex;
    align-items: center;
    margin: auto;
    padding: 0 9px 0 0;

    .search-type {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 90px;
      height: 100%;
      position: relative;
      font-size: 13px;

      &::after {
        content: "";
        width: 8px;
        height: 4px;
        background: url("~@/assets/material-library/v4.2.0/icon_down.png") center/100% no-repeat;
        margin-left: 5px;
      }

      &::before {
        content: "";
        width: 1px;
        height: 20px;
        background-color: #CCCCCC;
        position: absolute;
        right: 0;
        top: 50%;
        transform: translateY(-50%);
      }
    }

    .search-input {
      flex: 1;
      height: 100%;
      padding: 0 0 0 18px;
      font-size: 12px;
      user-select: auto;
      -webkit-user-select: auto;
    }

    .scan {
      width: 35px;
      height: 35px;
      background: url("~@/assets/material-library/v4.2.0/icon_scan.png") center/18px no-repeat;
        
    }
    
    .search-btn {
      width: 60px;
      height: 32px;
      background: #000000;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      font-size: 14px;
      &.disabled {
        background-color: #ccc
      }
    }
  }
  .van-list {
    padding: 15px;
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