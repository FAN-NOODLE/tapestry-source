<template>
  <div>
    <div v-if="detail && !showSendPreview" class="container">
      <div class="content">
        <div class="info">
          <div class="row title">
            <div class="col">
              <div class="label">品名：</div>
              <div class="value">{{detail.product_name||'--'}}</div>
            </div>
          </div>
          <div class="row">
            <div v-if="detail.store_code" class="col">
              <div class="label">门店code：</div>
              <div class="value">{{detail.store_code||'--'}}</div>
            </div>
            <div class="col">
              <div class="label">UPC：</div>
              <div class="value">{{detail.upc_code||'--'}}</div>
            </div>
            <div class="col">
              <div class="label">颜色：</div>
              <div class="value">{{detail.color||'--'}}</div>
            </div>
            <div class="col">
              <div class="label">吊牌价：</div>
              <div class="value">{{detail.sku_tag_price?`¥${detail.sku_tag_price}`:'--'}}</div>
            </div>
            <div class="col">
              <div class="label">款式：</div>
              <div class="value">{{detail.style_code||'--'}}</div>
            </div>
            <div v-if="detail.store_code" class="col">
              <div class="label">零售价：</div>
              <div class="value">{{detail.store_price?`¥${detail.store_price}`:'--'}}</div>
            </div>
          </div>
          <div class="row">
            <div class="col">
              <div class="label">产品分类：</div>
              <div class="value">{{detail.department||'--'}}</div>
            </div>
          </div>
          <div class="row">
            <div class="col">
              <div class="label">廓形说明：</div>
              <div class="value">{{detail.silhouette_desc||'--'}}</div>
            </div>
          </div>
        </div>
        <div v-if="detail.short_description_zh || detail.long_description_zh" ref="desc" :class="['desc',{expanded: isExpand}]">
          <div v-if="detail.short_description_zh" class="row">
            <div class="label">短描述：</div>
            <div class="value">
              <pre v-html="detail.short_description_zh"></pre>
            </div>
          </div>
          <div v-if="detail.long_description_zh" class="row">
            <div class="label">长描述：</div>
            <div class="value">
              <pre v-html="detail.long_description_zh"></pre>
            </div>
          </div>
        </div>
        <div v-if="showExpand" :class="['expand', {expanded: isExpand}]" @click="handleExpand">{{isExpand?'收起':'展开全部'}}</div>
        <div class="btn" @click="handleSearch">商品库存查询</div>
      </div>
      <div v-if="detail.merge_image.length" class="list">
        <div class="flex">
          <div class="select" @click="handleSelectAll">
            <div :class="['icon-select', icon]"></div>
            <div>全选</div>
          </div>
        </div>
        <van-row gutter="20">
          <van-col span="8" v-for="(item, index) in detail.merge_image" :key="index">
            <div class="item">
              <div :class="['item__select', {active: isChecked(item)}]" @click="handleSelectSingle(item)"></div>
              <div class="item__image" @click="handlePreviewImage(index)">
                <van-image v-if="item.type==2" class="image" :src="item.url" fit="contain"></van-image>
                <video ref="video" v-else class="image" :controls="false" :src="item.url" :poster="item.poster"></video>
                <div v-if="item.type==5" class="mark"></div>
              </div>
              <div class="item__title">
                <div class="text">{{item.name.split(".")[0]}}</div>
              </div>
            </div>
          </van-col>
        </van-row>
      </div>
      <!-- 底部操作 -->
      <div v-if="detail.merge_image.length" class="footer">
        <div class="number">已选：{{selection.length}}
          <div v-if="selection.length === maxlength" class="tips">最多可选{{maxlength}}个</div>
        </div>
        <div v-if="selection.length>0" class="btn-clear" @click="handleClearSelection">清空已选</div>
        <div :class="['btn', { disabled: selection.length === 0}]" @click="handleSend">{{ selection.length > 1 ? '组合发送' : '发送' }}</div>
      </div>
      <!-- 预览图片 -->
      <preview-image-popup 
        v-if="showPreviewPopup"
        :show.sync="showPreviewPopup"
        :show-index="previewIndex"
        :item="detail"
        :selection="selection"
        :isFirstOpen="false"
        @select="handleSelect"
      ></preview-image-popup>
    </div>

    <empty-page
      v-if="!isLoading && !detail"
      picture="noData"
      paddingTop="60"
      paddingBottom="265"
      :text="'未查询到该商品'"
    ></empty-page>

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
import { Component, Vue, Watch } from "vue-property-decorator";

import EmptyPage from "@/components/base/empty-page.vue";
import wxUtils from "@/utils/wx";
import { getBrandCollectMaterialListApi, reportTimes, sourceMakeMediaApi } from '@/api/marketing-promotion';
import { API } from '@/utils/env-config';
import PreviewImagePopup from "./components/preview-image-popup.vue"
import PreviewSent from "./components/preview-sent.vue";
import { getVideoPoster } from "@/utils/tools";
import { getCosCompressUrl } from "@/utils/upload/utils";

@Component({
  components: {
    PreviewImagePopup,
    PreviewSent,
    EmptyPage
  },
})
export default class GoodsDetail extends Vue {
  isLoading = true
  /** 详情 */
  detail = this.$store.getters['materialLibrary/getDetail']
  /** 选中项 */
  selection: any[] = this.$store.getters["materialLibrary/getSelection"]
  maxlength = 9

  /** 展开收起 */
  showExpand = false
  isExpand = true

  /** 查看预览图 */
  showPreviewPopup = false
  previewIndex = 0

  /** 发送预览 */
  showSendPreview = false

  get status() {
    const selectedIds = this.selection.map(item => item.id)
    
    const checkedLength = this.detail.merge_image.filter(item => selectedIds.includes(item.id))
    if(checkedLength.length === this.detail.merge_image.length) {
      return 2
    } else if(checkedLength.length > 0) {
      return 1
    } else {
      return 0
    }
  }

  /** 图标 */
  get icon() {
    return {
      0: '',
      1: 'active--half',
      2: 'active'
    }[this.status]
  }


  get isChecked() {
    return (item)=> {
      const ids = this.selection.map(item => item.id)
      return ids.includes(item.id)
    }
  }

  @Watch('selection', {immediate: true, deep: true})
  handleSelectionChange(val) {
    this.$store.commit("materialLibrary/setSelection", val);
  }

  created() {
    // this.getDetail()
  }
  async mounted() {

    this.$nextTick(() => {
      const videoRef = (this.$refs.video||[]) as HTMLVideoElement[]
      videoRef.forEach(async element => {
        const res:any = await getVideoPoster(element)
        element.poster = res.base64
      });
    })
    try {
      await wxUtils.wxAgentConfig();
      this.getWxContent();
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
  /**
   * 获取详情
   */
  async getDetail() {
    const res: any = await getBrandCollectMaterialListApi({
      last_code: '',
      upc_code: this.$route.query.upc_code
    })
    if(res.code === 200 && res.data.length) {
      this.detail = res.data.map((item: any) => {
        item.merge_image = (item.brand_source_list || []).map(image => {
          return {
            ...image,
            id: image.resource_key,
            is_collect: item.is_collect
          }
        })
        const currentImage = item.merge_image.find(image => image.is_main === 1) || item.merge_image[0]
        return {
          ...item,
          url: currentImage ? currentImage.url: '',
          name: currentImage ? currentImage.name : item.product_code,
          type: 8
        };
      })[0]

      this.isLoading = false
      this.setShowExpand()
    } else {
      // this.$toast('未查询到该商品')
      this.isLoading = false
    }
  }
  setShowExpand() {
    this.$nextTick(() => {
      const dom = this.$refs.desc as HTMLElement;
      if(dom) {
        this.showExpand = dom.clientHeight>200
      }
    })
  }
  handleExpand() {
    this.isExpand = !this.isExpand
  }

  /**
   * 查看库存详情
   */
  handleSearch() {
    const params = {
      searchValue: this.detail.upc_code,
      active: 1
    }
    // 转成字符串
    const paramsStr = encodeURIComponent(JSON.stringify(params))
    window.wx.invoke('launchMiniprogram', {
      "appid" : API.appid, // 需跳转的小程序appid
			"path" : `pagesStock/pages/list/list.html?from=3&searchParams=${paramsStr}`, // 所需跳转的小程序内页面路径及参数。非必填
    }, res => {
			if(res.err_msg == "launchMiniprogram:ok") {
				// 神策数据上报 -- 库存查询
        this.$track.trackStockInquiry({
          entry_point: 3,
          inquiry_type: 2
        })
			} else {
				// 错误处理
        this.$toast.fail('跳转失败')
			}
    });  
  }

  handleSelectAll() {
    this.handleSelect({
      status: this.status,
      item: this.detail,
      isMergeImage: true
    })
  }

  handleSelectSingle(item:any) {
    const ids = this.selection.map(val => val.id)
    // 判断是否全部选中
    this.handleSelect({
      status: ids.includes(item.id) ? 2: 0,
      item: {type: 2, ...item},
      isMergeImage: false
    })
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

  handlePreviewImage(index) {
    this.showPreviewPopup = true
    this.previewIndex = index
  }

  handleClearSelection() {
    this.selection = []
  }


  /**
   * 发送预览
   */
  handleSend() {
    if(this.selection.length==0 ) return
    if(this.selection.length==1) {
      this.sendChatMessage(this.selection[0])
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
    console.log(111111, this.$route.query.path);
    
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
            current_page: this.$route.query.path||this.$route.path,
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
.icon-select {
  position: relative;
  width: 16px;
  height: 16px;
  background-color: #fff;
  border-radius: 0px 0px 0px 10px;
  z-index: 2;
  &::after {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background-color: #fff;
    border: 1px solid #4A4A4A;
    box-sizing: border-box;
    z-index: 2;
  }
  &::before {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_select.png") no-repeat center center/100%;
    z-index: 1;
  }
  &.active {
    &::before {
      z-index: 2;
    }
    &::after {
      z-index: 1
    }
  }
  
  &.active--half {
    &::before {
      background: #fff;
      width: 8px;
      height: 3px;
      border-radius: 0;
      z-index: 2;
    }
    &::after {
      z-index: 1;
      background: #000
    }
  }
}


::v-deep li {
  &::marker {
    display: none !important;
  }
}

.container {
  min-height: 100vh;
  background-color: #F4F7FB;
  padding-bottom: calc(50px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
  padding-bottom: calc(50px + env(safe-area-inset-bottom)); /** iOS 11.0+ */


  .content {
    width: 375px;
    background: #FFFFFF;
    box-shadow: 0px 3px 17px 0px rgba(54,68,58,0.06);
    box-sizing: border-box;
    padding: 10px 15px 15px;

    .info {
      width: 345px;
      background: rgba(119, 119, 119, 0.1);
      box-sizing: border-box;
      padding: 15px;
      font-size: 12px;
      color: #4A4A4A;
      line-height: 17px;

      .title {
        line-height: 20px;
        font-size: 14px;
        .value {
          color: #000000;
        }
      }
      .row {
        display: flex;
        flex-wrap: wrap;

        .col {
          width: 50%;
          display: flex;
          .label {
            white-space: nowrap;
            padding-right: 5px;
          }
        }
        &:not(:last-child) .col {
          margin-bottom: 10px;
        }

        &.title .col {
          width: 100%;
        }
      }
    }

    .desc {
      font-size: 12px;
      color: #777777;
      line-height: 17px;
      transition: all 0.5s;
      overflow: hidden;
      height: 0;
      .row {
        display: flex;
        &:not(:last-child) {
          margin-bottom: 15px;
        }
        .label {
          width: 53px;
        }
        .value {
          flex: 1;
          text-align: justify;
          pre {
            white-space: pre-wrap;
          }
        }
      }

      &.expanded {
        height: auto;
        margin-top: 15px;
      }
    }

    .expand {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 15px 0 10px;
      font-weight: 400;
      font-size: 14px;
      color: rgba(0,0,0,0.5);
      line-height: 20px;

      &::after {
        content: "";
        width: 14px;
        height: 9px;
        background: url("~@/assets/material-library/v4.2.0/icon_top.png") no-repeat center center/100%;
        margin-left: 5px;
        transform: rotate(180deg);
        transition: all 0.3s;
      }
      
      &.expanded::after {
        transform: rotate(0);
      }
    }

    .btn {
      width: 285px;
      height: 36px;
      background: #000000;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      color: #FFFFFF;
      margin: auto;
      margin-top: 15px;
    }

    .expand + .btn {
      margin-top: 0;
    }
  }

  .select {
    display: flex;
    align-items: center;
    margin-bottom: 15px;
    font-size: 14px;
  
    .icon-select {
      margin-right: 5px;
      .van-icon {
        font-size: 16px;
        color: #4a4a4a;
      }
    }
  }

  .list {
    margin-top: 15px;
    padding: 0 15px;

    .item {
      width: 108px;
      padding: 7px 7px 10px 7px;
      box-sizing: border-box;
      margin-bottom: 15px;
      position: relative;
      background-color: #fff;
    
      &__select {
        position: absolute;
        top: 0;
        right: 0;
        width: 28px;
        height: 28px;
        background-color: #fff;
        box-shadow: -2px 2px 8px 0px rgba(0,0,0,0.1);
        border-radius: 0px 0px 0px 10px;
        z-index: 2;
        &::after {
          content: "";
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background-color: #fff;
          border: 1px solid #4A4A4A;
          box-sizing: border-box;
          z-index: 2;
        }
        &::before {
          content: "";
          position: absolute;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_select.png") no-repeat center center/100%;
    
          z-index: 1;
        }
    
        &.active {
          &::before {
            z-index: 2;
          }
          &::after {
            z-index: 1;
          }
        }
      }
    
      &__image {
        width: 94px;
        height: 94px;
        background: #F0F0F0;
        position: relative;
        .image {
          width: 100%;
          height: 100%;
        }
        .mark {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          background: rgba(0, 0, 0, 0.5) url("~@/assets/material-library/v4.2.0/icon_play.png") center/30px no-repeat;
        }
      }
      
      &__title {
        margin-top: 7px;
        font-size: 12px;
        line-height: 17px;
        font-weight: bold;
        display: flex;
        align-items: center;
        justify-content: center;
        text-align: center;
        height: 34px;
    
        .text {
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          word-break: break-all;
        }
      }
    
      &__btn {
        width: 94px;
        height: 32px;
        border: 1px solid #000;
        font-size: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-top: 10px;
      }
    }
  }
    
  .footer {
    width: 100vw;
    padding: 8px 15px;
    display: flex;
    align-items: center;
    box-shadow: 0px 1px 6px 0px rgba(0,0,0,0.15);
    box-sizing: border-box;
    position: fixed;
    bottom: 0;
    padding-bottom: calc(8px + constant(safe-area-inset-bottom)); /** iOS 11.0+ */
    padding-bottom: calc(8px + env(safe-area-inset-bottom)); /** iOS 11.0+ */
    background-color: #fff;
    z-index: 10;
    
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
      width: 165px;
      height: 34px;
      background: #000;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      font-weight: bold;
      color: #fff;
  
      &.disabled {
        background-color: #ccc
      }
    }
  }
}
</style>
