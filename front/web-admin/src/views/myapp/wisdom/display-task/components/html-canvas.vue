<template>
  <div class="fixed">
    <div ref="container" v-if="detail" class="container">
      <!-- 任务详情 -->
      <div class="detail">
        <div class="flex-center">
          <div class="title">门店：{{detail.store_name}}_{{ detail.store_code }}</div>
          <div v-if="detail.finish_status === 3" class="audit-status" :data-status="detail.finish_status">已通过</div>
        </div>
        <div class="tips">
          <span>提交时间：{{ detail.submit_time | formatDate('yyyy年mm月dd日 hh:mm:ss') }}</span>
          <span v-if="detail.finish_status === 3" class="grade">{{ detail.grade }}</span>
        </div>
      </div>
      <!-- VM留言 -->
      <div v-if="detail.review_remark" class="remark">
        <div class="remark-title">VM留言：</div>
        <div class="remark-content">{{ detail.review_remark }}</div>
      </div>
      <!-- 陈列内容 -->
      <div class="content">
        <div class="title">陈列内容</div>
        <div ref="sectionRef" class="section" v-for="(item, index) in detail.template_info" :key="index">
          <div class="section-title">{{ item.template_name }}:</div>
          <div v-if="item.image_list.length>0" class="image-list">
            <div class="image-item" v-for="(image, index) in item.image_list" :key="index">
              <!-- 图片 -->
              <div class="image-item__inner">
                <el-image class="image-item__img" :src="allImageMap.get(image.url)" />
                <div v-if="image._status === 1" class="image-item__fail"></div>
              </div>
              <!-- 评论 -->
              <div v-for="(comment, commentIndex) in image.comment" :key="commentIndex" class="comment-item">
                <el-image class="comment-item__avatar" :src="allImageMap.get(comment.avatar)"/>
                <div class="comment-item__body">
                  <div class="comment-item__info">
                    <span class="comment-item__name">{{comment.member_name}}</span>
                    <span v-if="comment.store_code && comment.is_vm !== 1" class="comment-item__store">{{comment.store_code}}</span>
                    <span v-if="comment.is_self == 1" class="comment-item__tag" :class="{gray: comment.is_vm != 1}">我</span>
                    <span v-else-if="comment.is_vm" class="comment-item__tag">VM</span>
                    <span class="comment-item__time">{{comment.comment_time | formatDate('yyyy/mm/dd hh:mm:ss')}}</span>
                  </div>
                  <div class="comment-item__content">
                    <div class="tr"></div>
                    <div v-if="comment.content.text" class="comment-item__text">
                      <pre class="comment-item__pre"><span v-if="comment.reply_name" class="prefix">回复 <span class="prefix-name">{{ comment.reply_name }}</span>：</span>{{comment.content.text}}</pre>
                    </div>
                    <div v-if="comment.content.image_url" class="comment-item__image-list">
                      <img 
                        class="image-item"
                        v-for="(image, imageIndex) in comment.content.image_url" 
                        :key="imageIndex" 
                        :src="allImageMap.get(image.url)"/>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div v-if="item.image_list.length === 0" class="is-empty">未上传内容</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from '@/api/wisdom/displayTask';
import { IStoreInfo } from "./audit-form.vue";
import html2canvas, { Options } from '@whowasme/html2canvas';
import { download } from '@/utils/download';
import api from '@/api/wisdom/displayTask'
let canvas = null;
/**
 * @description: 图片转base64
 * @param url 图片地址
 * @returns
 */
function imageToBase64(url: string): Promise<string> {
  return new Promise((resolve) => {
    const img = new Image();
    img.crossOrigin = 'anonymous';
    img.src = url + '?v=' + Date.now();
    img.onload = function () {
      if(!canvas) {
        canvas = document.createElement('canvas');
        console.log('我创建了一个canvas', canvas);
      }
      canvas.width = img.width;
      canvas.height = img.height;
      const ctx = canvas.getContext('2d')
      if(ctx) {
        // 清空画布
        ctx.fillRect(0, 0, img.width, img.height);
        // 绘制图片
        ctx.drawImage(img, 0, 0);
        // 转换为 base64
        resolve(canvas.toDataURL());
        // 清除画布
        ctx.clearRect(0, 0, img.width, img.height);
      }
    };
  });
}
/**
 * @description 替换图片链接
 */
async function replaceImageUrl(dom: HTMLElement, callback: Function) {
  
  // 获取页面上所有的图片元素
  const images = dom.getElementsByTagName('img');
  const imageMap = new Map()

  async function inner(image) {
    // 获取图片元素的 src 属性
    const oldSrc = image.getAttribute('src') || '';

    // 对 oldSrc 进行替换，得到新的图片地址
    // 如果图片是远程地址，则使用 imageToBase64 方法将其转换为 base64
    const isUrl = /^http(s?):\/\//i.test(oldSrc);
    
    let newSrc = oldSrc;
    if(isUrl) {
      if(imageMap.has(oldSrc)) {
        newSrc = imageMap.get(oldSrc)
      } else {
        newSrc = await imageToBase64(oldSrc)
        imageMap.set(oldSrc, newSrc)
      }
    }

    // 将图片元素的 src 属性替换为新地址
    image.setAttribute('src', newSrc);
  }

  for (let i = 0; i < images.length; i++) {
    await inner(images[i]);
    callback({
      current: i+1,
      total: images.length,
      message: `正在处理第${i+1}/${images.length}张图片`
    })
  }

  // return Promise.all([...images].map((image) => inner(image)));
}

@Component({
  components: {}
})
export default class HtmlCanvas extends Vue {
  /** 表单数据 */
  detail: IStoreInfo = null
  defaultAvatar = 'https://cdn-cos.coach.com.cn/uploads/common/icon_logo.png'
  taskDetail: any = null
  loading: any = null
  allImageMap = new Map()

  async getTaskDetail(params: any) {
    const res: any = await  Api.getDisplayTaskDetailApi({
      id: Number(params.id)
    })
    if (res.code == 200) {
      this.taskDetail = res.data
    } else {
      this.$message.error(res.message);
    }
  }
  /**
   * @description 获取任务详情
   */
  async getDetail(params: any) {
    const res: any = await Api.getAllotInfoApi( { 
      id: Number(params.id),
      store_id: Number(params.store_id)
    })

    if(res.code === 200) {
      // 清空数据
      this.allImageMap.clear()
      // 格式化图片数据
      const template_info = await Promise.all(this.formatData(res.data.image_list)) 
      this.loading = this.$loading({
        lock: true,
        text: '正在生成图片 0%',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)',
      });
      await this.loadAssets()
      console.log(template_info);
      
      // 赋值
      this.detail = {
        ...res.data,
        template_info
      }
    }
  }

  /**
   * 格式化图片数据
   * @param image_list
   */
   formatData(image_list) {
    return this.taskDetail.template_info.map(async (item) => {
      const imageSet = new Map()
      const imageList = await Promise.all(
        image_list
        .filter(image => image.identifier === item.template_name)
        .map(async value => {
          this.allImageMap.set(value.url, '')
          value._status = value.status
          value._image_list = []
          value.comment = await Promise.all(value.comment.map(async comment => {
            const images = comment.content.image_url?.map(url => {
              // 记录评论图片链接
              this.allImageMap.set(url, '')
              return {
                url: url as unknown as string,
              }
            })||[]
            // 头像链接替换
            if(comment.avatar && comment.avatar.indexOf('https://tapestry') === -1) {
              if(imageSet.has(comment.avatar)) {
                comment.avatar = imageSet.get(comment.avatar)
              } else {
                const originAvatar = comment.avatar
                comment.avatar = (await api.getCosUrlApi({ url: originAvatar }) as any).data.cos_avatar
                imageSet.set(originAvatar, comment.avatar)
              }
            } else {
              comment.avatar = this.defaultAvatar
            }
            // 记录头像链接
            this.allImageMap.set(comment.avatar, '')
            return {
              ...comment,
              _reply_content: '',
              _reply_name: '',
              _reply_show: false,
              content: images.length ? {image_url: images} : {text: comment.content.text}
            }
          }))
          return value
        }))

        console.log(imageList);
        
      return {
        ...item,
        image_list: imageList
      }
    })
  }

  /**
   * 加载资源
   */
  async loadAssets() {
    const len = this.allImageMap.size
    const list = [...this.allImageMap.keys()]
    for (let i = 0; i < len; i++) {
      const url = list[i];
      const value = await imageToBase64(url)
      this.allImageMap.set(url, value) 
      const percent = parseInt(((i+1)/(len+1))*100 + '');
      console.log('加载成功', url);
      this.loading.setText(`正在生成长图 ${percent}%`)
    }
  }

  // TODO 导出长图
  async save(params: any) {
    // this.loading = this.$loading({
    //   lock: true,
    //   text: '正在生成图片 0%',
    //   spinner: 'el-icon-loading',
    //   background: 'rgba(0, 0, 0, 0.7)',
    // });
    // 获取数据并格式化数据
    await this.getTaskDetail(params)
    await this.getDetail(params)

    this.$nextTick(() => {
      setTimeout(async() => {
        const element = this.$refs.container as HTMLElement
        const scale = 1
        // 方案一 一次绘制所有图片
        canvas.width = element.clientWidth * scale;
        canvas.height = element.clientHeight * scale;
        console.log(canvas.width, canvas.height);
        // 获取浏览器是IE还是谷歌或者火狐
        
        
        await html2canvas(element, { scale, canvas });
        this.loading.setText(`正在生成长图 99%`)
        // const base64Data = canvas.toDataURL('image/jpg', 0.9)
        // console.log(base64Data);
        canvas.toBlob((blob) => {
          this.loading.setText(`正在生成长图 100%`)
          download(blob, 'image/jpeg', this.detail.id + '.jpeg')
          this.loading.close()
          // 清除数据
          this.allImageMap.clear()
          this.detail = null
        }, 'image/jpeg', 0.95)

        // 方案二 分块绘制图片
        // const sections = this.$refs.sectionRef as HTMLElement[]
        

      }, 1000);
    })
  }
}
</script>
<style lang="scss" scoped>

.fixed {
  position: fixed;
  left: -800px;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}
.container {
  width: 357 * 2px;
  background-color: #F7F7F7;
  
  .detail {
    padding: 15 * 2px;
    background-color: #fff;
    .title {
      font-size: 16 * 2px;
      font-weight: 600;
      color: #000000;
      line-height: 22 * 2px;
      display: flex;
      padding-right: 15 * 2px;
      flex: 1;
    }
    .audit-status {
      font-weight: 400;
      font-size: 10 * 2px;
      width: 60 * 2px;
      height: 15 * 2px;
      border-radius: 15 * 2px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 5 * 2px;
      &[data-status="0"] {
        color: #000;
        background: rgba(0,0,0,0.1);
      }
      &[data-status="1"] {
        color: #000;
        background: #D1BFA7;
      }
      &[data-status="2"] {
        color: #E02020;
        background: #FFD0D0;
      }
      &[data-status="3"] {
        color: #000;
        background: #BBD1A7;
      }
    }

    .tips {
      font-size: 12 * 2px;
      font-weight: 400;
      color: #A6A6A6;
      line-height: 17 * 2px;
      margin-top: 10 * 2px;
      position: relative;

      .grade {
        position: absolute;
        right: 0;
        top: 40%;
        transform: translateY(-50%);
        font-size: 22 * 2px;
        font-weight: 600;
        color: #70A560;
        line-height: 30 * 2px;
        display: flex;
        align-items: center;
        justify-content: center;
        &::before {
          content: "评分:";
          font-size: 12 * 2px;
          font-weight: 400;
          margin-right: 2 * 2px;
          padding-top: 5 * 2px;
          white-space: nowrap;
        }
      }
    }

   .line {
      width: 345 * 2px;
      height: 37 * 2px;
      background: rgba(119, 119, 119, 0.1);
      justify-content: flex-start;
      padding: 0 15 * 2px;
      margin-top: 10 * 2px;

      .status {
        font-size: 12 * 2px;
        font-weight: bold;
        margin-right: 10 * 2px;
        &[data-status="0"] {
          color: #C6AD7F;
        }
        &[data-status="1"] {
          color: #70A560;
        }
        &[data-status="2"] {
          color: #A6A6A6;
        }
      }
      .time {
        font-size: 12 * 2px;
        font-weight: 400;
        color: #000000;
        line-height: 17 * 2px;
      }
    }

    .desc {
      font-size: 12 * 2px;
      font-weight: 400;
      color: #777777;
      line-height: 17 * 2px;
      margin-top: 10 * 2px;
      text-align: justify;
      word-break: break-all;
    }

    .link {
      margin-top: 15 * 2px;

      &-item {
        border: 2px solid #A6A6A6;
        box-sizing: border-box;
        padding: 20px 30px;

        &:not(:last-child) {
          margin-bottom: 15 * 2px;
        }
      }

      &-text {
        font-size: 12 * 2px;
        font-weight: 400;
        color: #000000;
        line-height: 17 * 2px;
        flex: 1;
        padding-right: 15 * 2px;
        word-break: break-all;
      }

      &-btn {
        width: 62 * 2px;
        height: 24 * 2px;
        background: #000000;
        font-size: 11 * 2px;
        font-weight: 400;
        color: #FFFFFF;
        line-height: 17 * 2px;
      }
    }

    .countdown {
      margin-top: 15 * 2px;
      .van-count-down {
        line-height: 18 * 2px;
        font-size: 12 * 2px;
        color: #70A560;
      }

      &.start .van-count-down {
        color: #C6AD7F;
      }
    }
  }

  .remark {
    background: rgba(0,0,0,0.1);
    box-shadow: 0 6px 34px 0 rgba(54,68,58,0.06);
    padding: 30px 40px;
    margin-top: 15 * 2px;

    &-title {
      font-size: 14 * 2px;
      font-weight: 500;
      color: #000000;
      line-height: 20 * 2px;
    }
    &-content {
      padding: 20px 20px 10px 20px;
      font-size: 12 * 2px;
      font-weight: 400;
      color: #000000;
      line-height: 17 * 2px;
      text-align: justify;
      word-break: break-all;
    }
  }

  .content {
    background-color: #fff;
    padding: 30px 40px;
    margin-top: 15 * 2px;
    .title {
      font-size: 14 * 2px;
      font-weight: 500;
    }

    .section {
      margin-top: 15 * 2px;

      .section-title {
        font-size: 12 * 2px;
        font-weight: 400;
        color: #000000;
        line-height: 17 * 2px;
      }

      .image-item {
        margin-top: 15 * 2px;

        .comment-item {
          margin-top: 15 * 2px;
          margin-bottom: 0;
        }

        &__inner {
          position: relative;
        }

        &__img {
          width: 100%;
          background-color: #F7F7F7;
          display: block;
        }

        &__fail {
          position: absolute;
          left: 0;
          top: 0;
          width: 55 * 2px;
          height: 55 * 2px;
          background: url("~@/assets/images/icon_fail_big.png") center/100% no-repeat;
        }
      }

      .is-empty {
        font-size: 12 * 2px;
        font-weight: 400;
        color: #A6A6A6;
        line-height: 17 * 2px;
        margin-top: 15 * 2px;
      }
    }
  }

}

.prefix {
  pointer-events: none; /* 防止伪元素阻碍用户与textarea交互 */
  font-size: 12 * 2px;

  &-name {
    color: #777777;
    font-weight: 500;
  }
}
.comment-item {
  font-size: 14 * 2px;
  color: #000000;
  display: flex;
  align-items: flex-start;
  margin-bottom: 15 * 2px;

  &__avatar {
    width: 35 * 2px;
    height: 35 * 2px;
    margin-right: 20 * 2px;
    border-radius: 50%;
  }

  &__body {
    flex: 1;
    // overflow: hidden;
  }

  &__info {
    display: flex;
    align-items: center;
    white-space: nowrap;
  }

  &__name {
    font-weight: 500;
    line-height: 20 * 2px;
    margin-right: 4 * 2px;
  }

  &__store {
    color: #777777;
    line-height: 17 * 2px;
    margin-right: 4 * 2px;
  }

  &__tag {
    font-size: 10 * 2px;
    padding: 5px 20px;
    border-radius: 17 * 2px;
    background: #C6AD7F;
    color: #fff;
    margin-right: 4 * 2px;

    &.gray {
      background: #A6A6A6;
    }
  }

  &__time {
    font-size: 10 * 2px;
    color: #A6A6A6;
    flex: 1;
    text-align: right;
  }

  &__content {
    margin-top: 6 * 2px;
    background: #F6F5F5;
    box-shadow: -4px 6px 16px 0 rgba(0,0,0,0.1);
    padding: 15 * 2px;
    position: relative;
    z-index: 1;
    .tr {
      width: 26 * 2px;
      height: 31 * 2px;
      position: absolute;
      left: -44px;
      top: 5 * 2px;
      background: url('~@/assets/images/icon_tr.png') no-repeat right center/30px 50px;
      z-index: -1;
    }
  }

  &__pre {
    font-size: 12 * 2px;
    color: #777777;
    white-space: pre-wrap;
    margin: 0;
    line-height: 17 * 2px;
  }

  &__image-list {
    display: flex;
    flex-wrap: wrap;
    .image-item {
      width: 100 * 2px;
      height: 100 * 2px;
      margin-right: 15 * 2px;
      position: relative;
      margin-bottom: 15 * 2px;
      margin-top: 0;
      object-fit: cover;

      &:nth-child(3n) {
        margin-right: 0;
      }
    }

    .image-item:nth-last-child(1) {
      margin-bottom: 0;
    }
    .image-item:nth-last-child(2) {
      margin-bottom: 0;
    }
  }
}

</style>