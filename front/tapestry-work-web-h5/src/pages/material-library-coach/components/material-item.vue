<template>
  <div :class="['item', {
    'item-image': [2, 5, 8].includes(item.type),
    'item-text': item.type === 1,
    'item-file': item.type === 6,
    'item-link': item.type === 3,
    'item-miniprogram': item.type === 4,
    'small': small
  }]" @click="handleSelect">
    <!-- 新标 -->
    <div v-if="item.is_new === 1" class="item-new"></div>
    <!-- 选择框 -->
    <div v-if="!disabledSelect" class="item-checkbox">
      <!-- <i class="check-box" :class="status === 2 ? 'checked-box': ''"></i> -->
      <van-icon :name="icon" />
    </div>

	<div v-if="item.member_name || item.member_avatar" class="item-uploader">
		<img :src="item.member_avatar || defaultAvatar" alt="上传者头像" />
		<span>{{ item.member_name || 'Coach KOS team' }}</span>
	</div>

    <!-- 顶部位置 -->
    <div class="top-box">
      <!-- 0否，1是 -->
      <div class="top-up" v-if="item.is_top === 1 && page === 'material'">置顶</div>
      <!-- 数量 -->
      <div v-if="isMergeImage" class="num">{{item.merge_image.length}}张</div>
    </div>

    <!-- 文本 -->
    <div v-if="item.type === 1" class="item-content">
      <div class="content">
        <pre>{{ item.content }}<span v-if="isApplets" class="copy-icon" @click.stop="handleCopy"></span></pre>
      </div>

      <div class="btn-box-text">
        <div class="btn-collect" :class="[{disabled: disabledSend}]" @click.stop="handleCollect" v-if="item.is_collect === 1"><i class="btn-collected-img"></i><span style="opacity: 0;">收藏</span></div>
        <div class="btn-collect":class="[{disabled: disabledSend}]" @click.stop="handleCollect" v-else><i class="btn-collect-img"></i>收藏</div>
        <div class="btn-text" :class="[{disabled: disabledSend}]" @click.stop="handleEdit"><i class="btn-edit-img"></i>编辑</div>
        <div class="btn-text" :class="[{disabled: disabledSend}]" @click.stop="handleSendMessage"><i class="btn-send-img"></i>发送</div>
      </div>
    </div>
    <!-- 图片  -->
    <div v-if="item.type === 2" class="item-content">
      <div v-if="hasSend" class="has-send">已发送</div>
      <!-- <div v-if="isMergeImage" class="num">{{item.merge_image.length}}张</div> -->
      <van-image :src="cosUrl(item.url)" lazy-load fit="contain" @click.stop="handlePreview"></van-image>
      <div class="name">
        <span class="ellipsis">{{ item.name.split(".")[0] }}</span>
      </div>

      <!-- 按钮 -->
      <div class="btn-box" :class="{'btn-merge': isMergeImage, 'btn-small': small}">
        <!-- 左 -->
        <div class="btn-left">
          <div class="btn-collect" @click.stop="handleCollect" v-if="item.is_collect === 1"><i class="btn-collected-img"></i><span v-if="!small" style="opacity: 0;">收藏</span></div>
          <div class="btn-collect" @click.stop="handleCollect" v-else><i class="btn-collect-img"></i><span v-if="!small">收藏</span></div>
        </div>
        <!-- 右 -->
        <div class="btn-right">
          <div v-if="disabled" class="btn disabled">压缩中，请稍后</div>
          <div v-else-if="isApplets" class="btn" @click.stop="handlePreview">预览</div>
          <div v-else-if="isMergeImage" class="btn" @click.stop="handleSelectImage"><i class="btn-send-img"></i>{{ selectCount ? '发送' : '选择发送' }}<span v-if="selectCount">({{selectCount}})</span></div>
          <div v-else :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage"><i class="btn-send-img"></i>发送</div>
        </div>
      </div>
      <!-- 按钮 end -->
    </div>
    <!-- 链接 -->
    <div v-if="item.type === 3" class="item-content">
      <div v-if="hasSend" class="has-send">已发送</div>
      <van-image :src="item.cover_url" lazy-load fit="contain" @click.stop="handlePreviewLink(item)"></van-image>
      <div class="item-link__info">
        <div class="title">{{ item.title }}</div>
        <div class="desc">{{ item.description }}</div>
        <!-- <div v-if="!isApplets" :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage">发送</div> -->
        <!-- 按钮 -->
        <div class="btn-box">
          <!-- 左 -->
          <div class="btn-left">
            <div class="btn-collect" @click.stop="handleCollect" v-if="item.is_collect === 1"><i class="btn-collected-img"></i><span style="opacity: 0;">收藏</span></div>
            <div class="btn-collect" @click.stop="handleCollect" v-else><i class="btn-collect-img"></i>收藏</div>
          </div>
          <!-- 右 -->
          <div class="btn-right">
            <div v-if="!isApplets" :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage"><i class="btn-send-img"></i>发送</div>
          </div>
        </div>
        <!-- 按钮 end -->
      </div>
    </div>
    <!-- 小程序 -->
    <div v-if="item.type === 4" class="item-content">
      <div class="app-name">
        <img
          class="icon"
          src="@/assets/material-library/mini.png"
          alt=""
        />
        <span>{{ item.app_name }}</span>
      </div>
      <div class="title">{{ item.title }}</div>
      <div v-if="hasSend" class="has-send">已发送</div>
      <van-image :src="item.cover_url" lazy-load></van-image>
      <!-- <div v-if="!isApplets" :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage">发送</div> -->
      <!-- 按钮 -->
      <div class="btn-box">
        <!-- 左 -->
        <div class="btn-left">
          <div class="btn-collect" @click.stop="handleCollect" v-if="item.is_collect === 1"><i class="btn-collected-img"></i><span style="opacity: 0;">收藏</span></div>
          <div class="btn-collect" @click.stop="handleCollect" v-else><i class="btn-collect-img"></i>收藏</div>
        </div>
        <!-- 右 -->
        <div class="btn-right">
          <div v-if="!isApplets" :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage"><i class="btn-send-img"></i>发送</div>
        </div>
      </div>
      <!-- 按钮 end -->
    </div>
    <!-- 视频 -->
    <div v-if="item.type === 5" class="item-content">
      <div v-if="hasSend" class="has-send">已发送</div>
      <div class="video">
        <van-image :src="coverImage(item.cover_url)" lazy-load fit="contain" @click.stop="handlePreviewVideo"></van-image>
        <div @click.stop="handlePreviewVideo" v-if="item.cover_url" class="mark"></div>
      </div>
      <div class="name">
        <span class="ellipsis">{{ item.name.split(".")[0] }}</span>
      </div>
      <!-- <div v-if="disabled" class="btn disabled">压缩中，请稍后</div>
      <div v-else-if="isApplets" class="btn" @click.stop="handlePreviewFile">预览</div>
      <div v-else :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage">发送</div> -->

      <!-- 按钮 -->
      <div class="btn-box">
        <!-- 左 -->
        <div class="btn-left">
          <div class="btn-collect" @click.stop="handleCollect" v-if="item.is_collect === 1"><i class="btn-collected-img"></i><span style="opacity: 0;">收藏</span></div>
          <div class="btn-collect" @click.stop="handleCollect" v-else><i class="btn-collect-img"></i>收藏</div>
        </div>
        <!-- 右 -->
        <div class="btn-right">
          <div v-if="disabled" class="btn disabled">压缩中，请稍后</div>
          <div v-else-if="isApplets" class="btn" @click.stop="handlePreviewFile">预览</div>
          <div v-else :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage"><i class="btn-send-img"></i>发送</div>
        </div>
      </div>
      <!-- 按钮 end -->
    </div>
    <!-- 文件 -->
    <div v-if="item.type === 6" class="item-content">
      <div v-if="hasSend" class="has-send">已发送</div>
      <van-image :src="fileTypeImage" lazy-load fit="contain"></van-image>
      <div class="name">
        <span class="ellipsis">{{ item.name.split(".")[0] }}</span>
      </div>
      <!-- <div v-if="isApplets" class="btn" @click.stop="handlePreviewFile">预览</div>
      <div v-else :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage">发送</div> -->
      <!-- 按钮 -->
      <div class="btn-box">
        <!-- 左 -->
        <div class="btn-left">
          <div class="btn-collect" @click.stop="handleCollect" v-if="item.is_collect === 1"><i class="btn-collected-img"></i><span style="opacity: 0;">收藏</span></div>
          <div class="btn-collect" @click.stop="handleCollect" v-else><i class="btn-collect-img"></i>收藏</div>
        </div>
        <!-- 右 -->
        <div class="btn-right">
          <div v-if="isApplets" class="btn" @click.stop="handlePreviewFile">预览</div>
          <div v-else :class="['btn', {disabled: disabledSend}]" @click.stop="handleSendMessage"><i class="btn-send-img"></i>发送</div>
        </div>
      </div>
      <!-- 按钮 end -->
    </div>
    <!-- 商品 -->
    <div v-if="item.type === 8" class="item-content">
      <div v-if="hasSend" class="has-send">已发送</div>
      <!-- <div v-if="item.merge_image.length>1" class="num">{{item.merge_image.length}}张</div> -->
      <div v-if="item._item_type === 5" class="video">
        <video ref="video" class="video-content" 
          :controls="false" 
          :src="item.url"
          playsinline="true"
          x-webkit-airplay="allow"
          x5-playsinline
          webkit-playsinline
          object-fit="contain"
          :poster="item.poster"
        ></video>
        <div @click.stop="handleViewDetail" class="mark"></div>
      </div>
      <van-image v-else-if="item.url" :src="item.url" lazy-load fit="contain" @click.stop="handleViewDetail"></van-image>
      <van-image v-else :src="emptyImage" lazy-load fit="contain"></van-image>
      <div class="name">
        <span class="ellipsis">{{ item.name.split(".")[0] }}</span>
      </div>

      <!-- 按钮 -->
      <div class="btn-box" :class="{'btn-merge': isMergeImage, 'btn-small': small}">
        <!-- 左 -->
        <div class="btn-left">
          <div class="btn-collect" @click.stop="handleCollect" v-if="item.is_collect === 1"><i class="btn-collected-img"></i><span v-if="!small" style="opacity: 0;">收藏</span></div>
          <div class="btn-collect" @click.stop="handleCollect" v-else><i class="btn-collect-img"></i><span v-if="!small">收藏</span></div>
        </div>
        <!-- 右 -->
        <div class="btn-right">
          <div v-if="disabled" class="btn disabled">压缩中，请稍后</div>
          <div v-else-if="isApplets" class="btn" @click.stop="handlePreview">预览</div>
          <div v-else-if="item.merge_image.length>0" class="btn" @click.stop="handleViewDetail"><i class="btn-send-img"></i>{{ selectCount ? '发送' : '选择发送' }}<span v-if="selectCount">({{selectCount}})</span></div>
          <div v-else class="btn" @click.stop="handleViewDetail">查看详情</div>
        </div>
      </div>
      <!-- 按钮 end -->
    </div>
    <!-- 预览链接 -->
    <div class="link-wrapper" v-if="showPreview">
      <div class="link-header">链接预览
        <div class="link-close" @click.stop="showPreview = false"></div>
      </div>
      <iframe class="link-content" :src="iframeUrl" frameborder="0"></iframe>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, PropSync } from 'vue-property-decorator'
import { Icon, Image, ImagePreview } from 'vant';
import Helper from "@/utils/helper";
import { getVideoPoster } from '@/utils/tools';

import {setCollectApi } from '@/api/marketing-promotion';

export interface IMaterialItem {
  /** 素材id */
  id: number;
  /** 素材类型 1文本 2图片 3链接 4小程序 5视频 6文件 7朋友圈素材 8商品 */
  type: 1|2|3|4|5|6|7|8;
  _item_type: 1|2|3|4|5|6|7|8;
  /** 素材地址 */
  url: string;
  /** 素材名称 */
  name: string;
  /** 素材大小 */
  size: number;
  /** 压缩处理：0:处理成功 1:处理中 2:处理失败 */
  is_deal: number;
  /** 是否新标 */
  is_new: number;
  /** 素材组合id */
  merge_id: number;
  /** 素材组合图片 */
  merge_image: Array<any>;
  /** 文本素材内容 */
  content: string;
  /** 视频素材封面图 */
  cover_url: string;
  /** 文件素材-文件类型 1:word 2:excel 3:ppt 4:pdf 5:其他 */
  file_type: 1|2|3|4|5;
  /** 小程序素材标题 */
  app_name: string;
  /** 小程序素材分享标题 / 链接素材标题*/
  title: string;
  /** 链接素材描述 */
  description: string;
  /** 商品唯一标识 */
  upc_code: string;
  /** 视频封面 */
  poster?: string,
  /** 收藏 */
  is_collect: number;
  /** 置顶 */
  is_top: number;
  product_code: string;
  product_name: string;
  /** 收藏数量 */
  collect_num: number;
  member_name?: string;
  member_avatar?: string;
}

@Component({
  components: {
    [Icon.name]: Icon,
    [Image.name]: Image
  },
})
export default class MaterialItem extends Vue {
  @Prop({ type: Object, default: () => ({})}) item: IMaterialItem;
  @Prop({ type: Boolean, default: false}) small: boolean;
  @Prop({ type: Number, default: 9}) maxlength: number;
  @Prop({ type: Array, default: () => ([]) }) selection: IMaterialItem[]
  @Prop({ type: String, default: ''}) page: string;
  videoPoster = ''
  emptyImage = require('@/assets/material-library/v4.2.0/icon_empty.png')
	defaultAvatar = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_logo.png'
  iframeUrl = ''
  showPreview = false
  /**
   * 是否来自消息通知
   */
  get isApplets() {
    return this.$route.query.is_applets;
  }


  get cosUrl() {
    return (url) => {
      return url + '?imageMogr2/thumbnail/!30p'
    }
 }

  /** 选中状态 */
  get status() {
    const selectedIds = this.selection.map(item => item.id)
    
    if(this.isMergeImage) {
      const checkedLength = this.item.merge_image.filter(item => selectedIds.includes(item.id))
      if(checkedLength.length === this.item.merge_image.length) {
        return 2
      } else if(checkedLength.length > 0) {
        return 1
      } else {
        return 0
      }
    }
    if (selectedIds.includes(this.item.id)) {
      return 2
    }
    return 0
  }
  /** 图标 */
  get icon() {
    return {
      0: 'circle',
      1: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/half.png',
      2: 'checked'
    }[this.status]
  }

  /** 判断企业 */
  get isCoach() {
    return Helper.isCoach()
  }

  /** 是否是组合图片 */
  get isMergeImage() {
    const { merge_image, type } = this.item
    if(type === 8) {
      return merge_image.length > 0
    }
    return merge_image && merge_image.length > 1
  }

  /** 
   * 已发送
   */
  get hasSend() {
    return false
    // const list = this.$store.getters["materialLibrary/getNowSentData"]
    // return list.findIndex(item => item.id === this.item.id) > -1
  }

  /**
   * 视频封面图片
   */
  get coverImage() {
    return (url: string) => {
      const regex = /^(((ht|f)tps?):\/\/)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-z]{2,6}\/?/
      if (regex.test(url)) {
        return url
      } else {
        return require("@/assets/video-bg.png")
      }
    }
  }

  get fileTypeImage() {
    const { file_type } = this.item
    return {
      1: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/1.png',
      2: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/2.png',
      3: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/3.png',
      4: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/4.png',
      5: 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/5.png'
    }[file_type]
  }

  /**
   * 是否禁用
   */
  get disabled() {
    const { type, is_deal } = this.item
    return is_deal === 1 && [2, 5].includes(type)
  }

  get disabledSelect() {
    const { type, is_deal, merge_image } = this.item
    // 消息进入禁止选择
    if( this.isApplets ) return true
    // 小红书素材不支持多选
    if(type === 3) return true
    // 商品素材 没有图片禁止选择
    if(type === 8) {
      return merge_image.length === 0
    }
    // 图片视频没有处理完成 禁止选择
    return is_deal === 1 && [2, 5].includes(type)
  }

  /** 组合素材选中数量 */
  get selectCount() {
    const ids = this.selection.map(item => item.id)
    return this.item.merge_image.filter(item => ids.includes(item.id)).length
  }

  /** 是否可以发送 */
  get disabledSend() {
    return this.selection.length > 0
  }

   mounted() {
    this.$nextTick(async() => {
      try {
        const videoRef = this.$refs.video as any
        if(videoRef) {
          const res:any = await getVideoPoster(videoRef)
          this.videoPoster = res.base64
          this.$emit('poster', res.base64)
        }
        
      } catch (error) {
        console.log(error);
        
      }
    })
  }

  /**
   * 选择
   */
  handleSelect() {
    if(this.disabledSelect) return
    this.$emit("select", {
      item: this.item, 
      status: this.status,
      isMergeImage: this.item.type===8 ? true : this.isMergeImage
    })
  }
  
  /**
   * 选择发送
   */
  handleSelectImage() {
    this.$emit('selectImage', this.item)
  }
  /** 收藏 */
  async handleCollect() {
    let item = this.item
    const trackList = []
    const params:any = {
      source_info: [
        {
          type: item.type
        }
      ],
      is_collect: item.is_collect ? 2 : 1
    }
  
    switch (item.type) {
      case 1:
        params.source_info[0].source_id = [item.id]
        trackList.push({
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 1,
          asset_type_modify: item.type
        })
        break
      case 2:
        if (item.merge_id) {
          const images = item.merge_image.filter(v => {
            if (params.is_collect === 1) {
              return v.is_collect === 0
            } else {
              return v.is_collect === 1
            }
          })
          trackList.push(...images.map(v => ({
            is_top: !!item.is_top,
            asset_id: v.id,
            is_official: 1,
            asset_type_modify: item.type
          })))
          params.source_info[0].source_id = images.map(v => v.id)
        } else {
          params.source_info[0].source_id = [item.id]
          trackList.push({
            is_top: !!item.is_top,
            asset_id: item.id,
            is_official: 1,
            asset_type_modify: item.type
          })
        }

        params.count = params.source_info[0].source_id.length < item.merge_image.length ? 0 : 1
        break
      case 3:
        params.source_info[0].source_id = [item.id]
        trackList.push({
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 1,
          asset_type_modify: item.type
        })
        break
      case 4:
        params.source_info[0].source_id = [item.id]
        trackList.push({
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 1,
          asset_type_modify: item.type
        })
        break
      case 5:
        params.source_info[0].source_id = [item.id]
        trackList.push({
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 1,
          asset_type_modify: item.type
        })
        break
      case 6:
        params.source_info[0].source_id = [item.id]
        trackList.push({
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 1,
          asset_type_modify: item.type
        })
        break
      case 7:
        params.source_info[0].merge_id = [item.merge_id]
        trackList.push({
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 1,
          asset_type_modify: item.type
        })
        break
      case 8:
        params.source_info[0].product_code = [item.product_code]
        trackList.push({
          is_top: !!item.is_top,
          asset_id: item.id,
          is_official: 2,
          asset_type_modify: item.type,
          official_url: item.url,
          official_upc: item.upc_code,
          official_product: item.product_name,
          official_name: item.product_code
        })
        break

      default:
        break
    }
      
    const res:IApiResult = await setCollectApi(params)
    if (res.code == 200) {
      item.is_collect = params.is_collect === 2 ? 0 : 1
      item.collect_num = res.data.collect_num
      // 组合素材
      if (item.type === 2 && item.merge_id) {
        item.merge_image.forEach(image => {
          const v = res.data.source_list[0]
          image.is_collect = params.is_collect === 2 ? 0 : 1
          image.collect_num = v.collect_num
        })
      }
      this.$toast(item.is_collect === 1 ? "已收藏": "已取消收藏")

      // 取消收藏
      if(this.page === 'collect' && item.is_collect === 0) {
        this.$emit('cancelCollect', item)
      }
      trackList.forEach(v => {
        // 神策数据上报 -- 收藏素材
        this.$track.trackAssetFavorite({
          current_page: this.$route.path,
          entry_point: 2,
          library_type: 2,
          operation_position: 1, // 固定值
          operation_type: params.is_collect,
          official_url: '非官网商品素材',
          official_upc: '非官网商品素材',
          official_product: '非官网商品素材',
          official_name: '非官网商品素材',
          ...v
        })
      })
    }
  }
  /**
   * 发送消息
   */
  handleSendMessage() {
    if (this.disabledSend) {
      this.$toast("需清空已选素材 才可进行单独发送")
      return
    }
    this.$emit('send', this.item)
  }
  /**
   * 预览图片
   */
  handlePreview() {
    const { url } = this.item;
    // 处理中 和 处理失败 的图片不能预览
    if (this.disabled) {
      return;
    }
    if(this.isCoach) {
      if(this.isMergeImage) {
        this.item.merge_image = this.item.merge_image.map((item) => {
          return {
            type: this.item.type,
            ...item
          }
        })
      }
      this.$emit('preview', {
        item: this.item,
        index: 0
      })
    } else {
      ImagePreview({
        images: [url],
        startPosition: 0,
      });
    }
  }
  /**
   * 预览视频
   */
  handlePreviewVideo() {
    this.$emit('preview-video', this.item)
    // 神策数据上报 -- 查看素材
    this.$track.trackAssetView({
      source_id: this.item.id,
      task_id: '素材库无任务ID',
      is_official: 1,
      asset_type_modify: this.item.type,
      official_url: '非官网商品素材',
      official_upc: '非官网商品素材',
      official_product: '非官网商品素材',
      official_name: '非官网商品素材',
      current_page: this.$route.path,
      entry_point: 2,
      library_type: 2,
      is_top: !!this.item.is_top,
      is_favorite: !!this.item.is_collect
    })
  }
  /**
   * 预览文件
   */
  handlePreviewFile() {
    const e = this.item;
    // 判断文件名字是否有后缀
    if (e.name.indexOf(".") === -1) {
      // 从url中获取后缀
      let ext = e.url.split(".").pop();
      e.name = e.name + "." + ext;
      console.log(e.name);
    }
    window.wx.invoke("previewFile", {
      url: e.url, // 需要预览文件的地址(必填，可以使用相对路径)
      name: e.name, // 需要预览文件的名称(必填)
      size: e.size, // 需要预览文件的字节大小(必填)
    });
  }

  /** 
   * 编辑文本素材
   */
  handleEdit() {
    this.$emit('edit', this.item)
  }

  /**
   * 复制文本素材
   */
  handleCopy() {
    const { content } = this.item
    let oInput = document.createElement("input");
    oInput.value = content;
    document.body.appendChild(oInput);
    oInput.select(); // 选择对象;
    document.execCommand("Copy"); // 执行浏览器复制命令
    oInput.remove();
    this.$toast("复制成功");
  }

  handleViewDetail() {
    this.$store.commit('materialLibrary/setDetail', this.item)
    this.$router.push({
      name: "material-detail",
      query: {
        upc_code: this.item.upc_code,
        path: this.$route.path,
      },
    });
  }

  /**
   * 预览链接
   */
  handlePreviewLink(item) {
    console.log(item);
    
    window.open(item.link_url)
    // this.iframeUrl = item.link_url
    // this.showPreview = true
  }
}
</script>
<style lang="scss" scoped>
.link-wrapper {
  width: 100vw;
  height: 100vh;
  position: fixed;
  z-index: 1000;
  top: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  .link-header {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    font-weight: bold;
    border-bottom: 1px solid #E5E5E5;

    .link-close {
      position: absolute;
      right: 0;
      top: 0;
      cursor: pointer;
      width: 40px;
      height: 40px;
      // 写一个x
      &::before { 
        content: '';
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        width: 15px;
        height: 2px;
        background-color: #000000;
        transform: translate(-50%, -50%) rotate(45deg);
      }
      &::after {
        content: '';
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%) rotate(-45deg);
        width: 15px;
        height: 2px;
        background-color: #000000;
      }
    }
  }

  .link-content { 
    width: 100%;
    height: calc(100vh - 40px);
    overflow-y: auto;
  }
}
.item {
  position: relative;
  margin-bottom: 15px;
  background-color: #fff;
  box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);

	.item-uploader {
		display: flex;
		align-items: center;
		min-width: 0;
		padding: 10px 12px 0;
		font-size: 12px;
		color: #666;

		img {
			width: 24px;
			height: 24px;
			flex: 0 0 24px;
			border-radius: 50%;
			margin-right: 6px;
		}

		span {
			overflow: hidden;
			text-overflow: ellipsis;
			white-space: nowrap;
		}
	}
    

  .item-new {
    width: 39px;
    height: 39px;
    position: absolute;
    top: -13px;
    left: -13px;
    z-index: 1;
    background: url('https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/new.png') no-repeat;
    background-size: 100% 100%;
  }

  .item-checkbox {
    position: absolute;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    box-shadow: -2px 2px 8px 0px rgba(0, 0, 0, 0.1);
    border-radius: 0px 0px 0px 10px;
    right: 0;
    top: -1px;
    background-color: #ffffff;
    z-index: 1;

    .van-icon {
      font-size: 16px;
      color: #4a4a4a;
    }
  }

  // .item-checkbox {
  //   position: absolute;
  //   display: flex;
  //   align-items: center;
  //   justify-content: center;
  //   width: 28px;
  //   height: 28px;
  //   box-shadow: -2px 2px 8px 0px rgba(0, 0, 0, 0.1);
  //   border-radius: 0px 0px 0px 10px;
  //   right: 0;
  //   top: -1px;
  //   background-color: #ffffff;
  //   z-index: 1;

  //   .check-box{
  //     display: inline-block;
  //     width: 15px;
  //     height: 15px;
  //     margin-right: 5px;
  //     background: url('~@/assets/material-library/v5.1.5/m0529-check.png') no-repeat center center/15px 15px;
  //   }
  //   .checked-box{
  //     background: url('~@/assets/material-library/v5.1.5/m0529-checked.png') no-repeat center center/15px 15px;
  //   }
  // }


  // .btn {
  //   height: 32px;
  //   display: flex;
  //   justify-content: center;
  //   align-items: center;
  //   cursor: pointer;
  //   border: 1px solid #000000;
  //   color: #000000;
  //   border-radius: 0;
  //   font-size: 13px;

  //   &.disabled {
  //     font-size: 12px;
  //     border: 1px solid #CCCCCC;
  //     color: #CCCCCC;
  //   }
  // }

  .btn-collect-img{
    display: inline-block;
    width: 12px;
    height: 12px;
    background: url('~@/assets/material-library/v5.1.5/m0529-collect01.png') no-repeat center center/100% 100%;
    margin-right: 5px;
  }
  .btn-collected-img{
    display: inline-block;
    width: 12px;
    height: 12px;
    background: url('~@/assets/material-library/v5.1.5/m0529-collect02.png') no-repeat center center/100% 100%;
    margin-right: 5px;
  }
  .btn-send-img{
    display: inline-block;
    width: 12px;
    height: 12px;
    background: url('~@/assets/material-library/v5.1.5/m0529-send.png') no-repeat center center/12px 12px;
    margin-right: 5px;
  }
  .btn-edit-img{
    display: inline-block;
    width: 12px;
    height: 12px;
    background: url('~@/assets/material-library/v5.1.5/m0529-edit.png') no-repeat center center/100% 100%;
    margin-right: 5px;
  }

  .btn-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 12px;
    border-top: 1px dashed #EEEEEE;
    padding-top: 10px;
    .btn-left{
      width: 50%;
      display: flex;
      justify-content: center;
      .btn-collect{
        display: flex;
        align-items: center;
        flex: 1;
      }
    }
    .btn-right {
      width: 50%;
      display: flex;
      justify-content: center;
      .btn {
        display: flex;
        align-items: center;
        cursor: pointer;
        text-align: right;
        color: #000000;
        font-size: 12px;

        &.disabled {
          font-size: 12px;
          opacity: 0.4;
        }
      }
    }
  }

  .btn-merge {
    .btn-left {
      width: 30%;
    }
    .btn-right {
      width: 70%;
      justify-content: flex-end;
    }
  }

  // 按钮 - 三列
  .btn-small {
    .btn-left { width: auto;}
    .btn-right { width: auto;}
  }

  // 顶部按钮
  .top-box {
    position: absolute;
    left: 14px;
    top: 14px;
    z-index: 1;
    display: flex;
    align-items: center;
    .top-up {
      width: 31px;
      height: 14px;
      background: #000000;
      border-radius: 9px;
      line-height: 14px;
      text-align: center;
      font-size: 10px;
      color: #fff;
    }
    .num {
      display: flex;
      align-items: center;
      justify-content: center;
      min-width: 32px;
      height: 14px;
      background: rgba(0, 0, 0, 0.5);
      border-radius: 9px;
      font-size: 10px;
      font-weight: 400;
      color: #FFFFFF;
      z-index: 9;
      padding: 0 5px;
    }
    .top-up + .num {
      margin-left: 5px;
    }
  }

  &.item-image {
    width: calc(145px + 20px);
    padding: 10px;

    .has-send {
      position: absolute;
      background-color: rgba(54, 54, 54, 0.49);
      color: #ffffff;
      text-align: center;
      line-height: 28px;
      font-size: 14px;
      font-weight: 400;
      left: 50%;
      transform: translateX(-50%);
      z-index: 1;
      top: 130px;
      width: 145px;
    }
    
    // .num {
    //   position: absolute;
    //   top: 15px;
    //   left: 15px;
    //   display: flex;
    //   align-items: center;
    //   justify-content: center;
    //   min-width: 32px;
    //   height: 17px;
    //   background: rgba(0, 0, 0, 0.5);
    //   border-radius: 9px;
    //   font-size: 10px;
    //   font-weight: 400;
    //   color: #FFFFFF;
    //   z-index: 9;
    //   padding: 0 5px;
    // }

    .van-image {
      width: 145px;
      height: 145px;
      background: #f0f0f0;
      display: block;
    }
    
    .video {
      width: 145px;
      height: 145px;
      position: relative;

      .video-content {
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

    
    .name {
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      color: #000000;
      padding: 5px 10px;
      height: 44px;
      
      .ellipsis {
        line-height: 17px;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: normal;
        word-break: break-all;
        text-align: center;
      }
    }

    &.small {
      width: calc(96px + 14px);
      margin-bottom: 10px;
      padding: 7px;

      .item-new {
        width: 34px;
        height: 34px;
        top: -12px;
        left: -12px;
      }

      .num {
        top: 10px;
        left: 10px;
      }

      .van-image {
        width: 96px;
        height: 96px
      }
      .video {
        width: 96px;
        height: 96px
      }

      .has-send {
        width: 94px;
        top: 72px
      }

      .name {
        padding: 3px 7px;
      }

      .btn {
        height: 28px;
        line-height: 28px;
        font-size: 12px;
      }
    }
  }

  &.item-file {
    width: calc(145px + 20px);
    padding: 15px 10px;

    .has-send {
      position: absolute;
      background-color: rgba(54, 54, 54, 0.49);
      color: #ffffff;
      text-align: center;
      line-height: 28px;
      font-size: 14px;
      font-weight: 400;
      left: 50%;
      transform: translateX(-50%);
      z-index: 1;
      width: 145px;
      bottom: 96px;
    }
    
    // .num {
    //   position: absolute;
    //   top: 15px;
    //   left: 15px;
    //   display: flex;
    //   align-items: center;
    //   justify-content: center;
    //   min-width: 32px;
    //   height: 17px;
    //   background: rgba(0, 0, 0, 0.5);
    //   border-radius: 9px;
    //   font-size: 10px;
    //   font-weight: 400;
    //   color: #FFFFFF;
    //   z-index: 9;
    //   padding: 0 5px;
    // }
    
    .van-image {
      width: 70px;
      height: 80px;
      margin: 10px auto;
      display: block;
    }

    .name {
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 12px;
      color: #000000;
      padding: 5px 10px;
      height: 44px;
      
      .ellipsis {
        line-height: 17px;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: normal;
        word-break: break-all;
        text-align: center;
      }
    }
  }

  &.item-text {
    padding: 15px;

    .content {
      padding: 25px 10px 10px 0;
      font-size: 14px;
      color: #000000;
      pre {
        line-height: 20px;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 4;
        -webkit-box-orient: vertical;
        white-space: pre-wrap;
        word-break: break-all;
      }
    }

    .btn-group {
      display: flex;
      justify-content: flex-end;
    }

    .btn {
      height: 28px;
      line-height: 28px;
      width: 64px;
      margin-left: 12px;
    }

    .has-send {
      float: left;
      font-size: 14px;
      color: #96979e;
      line-height: 32px;
    }
    .btn-box-text{
      display: flex;
      align-items: center;
      justify-content:space-between;
      font-size: 12px;
      border-top: 1px dashed #EEEEEE;
      padding: 10px 0 0;
      .btn-collect{
        width: 33.33%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .btn-text {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 33.33%;
        color: #000000;
        font-size: 13px;
      }
    }
    .copy-icon{
      display: inline-block;
      width: 12px;
      height: 12px;
      background: url('~@/assets/material-library/v5.1.5/m0529-copy.png') no-repeat center center/100% 100%;
      margin-left: 5px;
    }
  }

  &.item-miniprogram {
    width: calc(145px + 20px);
    padding: 10px;

    .app-name {
      display: flex;
      align-items: center;
      margin: 0;
      overflow: hidden;
      flex: auto;
      height: 20px;
      line-height: 20px;
      font-size: 13px;
      color: #000000;
      text-overflow: ellipsis;
      white-space: nowrap;
      font-weight: bold;
      .icon {
        flex: 0 0 20px;
        display: block;
        width: 20px;
        height: 20px;
        object-fit: cover;
        margin-right: 5px;
      }
    }

    .title {
      height: 36px;
      margin: 4px 0;
      font-size: 12px;
      line-height: 18px;
      color: #777777;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      word-break: break-all;
    }

    .has-send {
      position: absolute;
      background-color: rgba(54, 54, 54, 0.49);
      color: #ffffff;
      text-align: center;
      line-height: 28px;
      font-size: 14px;
      font-weight: 400;
      left: 50%;
      transform: translateX(-50%);
      z-index: 1;
      width: 145px;
      bottom: 50px;
    }

    .van-image {
      width: 145px;
      height: 145px;
      display: block;
    }
    // .btn {
    //   margin-top: 10px;
    // }
    .btn-box{
      border: none;
      padding-top: 10px;
    }

    .top-box{
      top: 80px;
      left: 16px;
    }
  }

  &.item-link {
    padding: 15px;
    .item-content {
      display: flex;
    }

    .van-image {
      width: 97px;
      height: 97px;
      display: block;
      background-color: #f0f0f0;
    }

    .item-link__info {
      flex: 1;
      overflow: hidden;
      margin-left: 10px;
      position: relative;

      .title {
        margin-bottom: 5px;
        font-size: 13px;
        font-weight: bold;
        color: #000000;
        line-height: 20px;
      }

      .desc {
        display: -webkit-box;
        overflow: hidden;
        text-overflow: ellipsis;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 2;
        margin-bottom: 5px;
        line-height: 20px;
        font-size: 12px;
        color: #777777;
      }
      .btn-box{
        border: none;
        padding-top: 0;
      }

      // .btn {
      //   position: absolute;
      //   right: 0;
      //   bottom: 0;
      //   height: 28px;
      //   line-height: 28px;
      //   width: 64px;
      // }
    }

    .has-send {
      position: absolute;
      background-color: rgba(54, 54, 54, 0.49);
      color: #ffffff;
      text-align: center;
      line-height: 28px;
      font-size: 14px;
      font-weight: 400;
      z-index: 1;
      bottom: 15px;
      width: 97px;
      left: 15px;
    }
    .btn-box{
      .btn-left, .btn-right, .btn-collect{
        width: auto;
      }
    }
    .top-box{
      top: 20px;
      left: 20px;
    }
  }

}
</style>
