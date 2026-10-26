<template>
  <div class="corporate-item" @click="handleView">
    <van-image class="corporate-item__cover" lazy-load :src="item.cover || defaultCover" fit="contain"></van-image>
    <div class="corporate-item__content">
      <div class="title">{{ item.title }}</div>
      <div class="remark">{{ item.ro_name }}</div>
      <div class="time">{{ formatTime(item.send_time) }}</div>
      <div class="num">
        <div class="num-item">
          <div class="num-item__icon"></div>
          <div class="num-item__value">{{item.view_people_num||0}}</div>
        </div>
        <div class="num-item">
          <div class="num-item__icon"></div>
          <div class="num-item__value">{{item.view_num||0}}</div>
        </div>
      </div>
    </div>
    <div class="is-new" v-if="!item.is_read"></div>
    <!-- 类型 -->
    <div class="type">{{ typeName }}</div>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator';
import { Image } from "vant";
import { defaultCover, typeOptions } from '../const';
export interface ICorporateItem {
  /** 传讯ID */
  id?: number;
  /** 传讯标题 */
  title: string;
  /** 封面图 */
  cover: string;
  /** 创建时间 */
  send_time: number;
  /** 传讯状态 */
  status: number;
  /** 发布人 */
  ro_name: string;
  /** 是否已读 */
  is_read: boolean;
  /** 阅读人数 */
  view_people_num: number;
  /** 阅读次数 */
  view_num: number;
  /** 类型 */
  type: number;
}

@Component({
  components: {
    [Image.name]: Image,
  },
  model: {
    prop: 'item',
    event: 'update'
  }
})
export default class CorporateItem extends Vue {
  /** 数据 */
  @Prop({ required: true }) item!: ICorporateItem;

  /** 默认封面图 */
  defaultCover = defaultCover

  activated() {
    const id = Number(localStorage.getItem('communicationId'))
    if(id && id === this.item.id) {
      this.$emit('update', {...this.item, is_read: true})
      localStorage.removeItem('communicationId')
    }
  }

  get typeName() {
    return typeOptions.find(item => item.value === this.item.type)?.label
  }

  /**
   * 时间格式化
   * @param time 
   */
  formatTime(time: number) {
    return this.$helper.dateFormat(time, "yyyy/MM/dd hh:mm");
  }

  handleView() {
    localStorage.setItem('communicationId', String(this.item.id))
    this.$router.push({
      name: 'corporate-detail',
      query: {
        id: this.item.id + '',
        from: '1'
      }
    })
  }
}
</script>
<style lang="scss" scoped>
.corporate-item {
  width: 345px;
  min-height: 124px;
  background: #FFFFFF;
  box-shadow: 0px 1px 3px 0px rgba(0,0,0,0.1);
  margin: auto;
  margin-bottom: 15px;
  background-color: #fff;
  display: flex;
  padding: 8px;
  box-sizing: border-box;
  position: relative;

  &__cover {
    width: 108px;
    height: 108px;
    background: rgba(0, 0, 0, 0.2);
    // box-shadow: 0px 1px 3px 0px rgba(0,0,0,0.1);

    .van-image {
      width: 100%;
      height: 100%;
    }
  }

  &__content {
    flex: 1;
    overflow: hidden;
    padding-left: 20px;
    position: relative;
    padding-bottom: 30px;

    .title {
      font-weight: 600;
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      text-align: left;
      font-style: normal;
      padding-top: 7px;
      // 超出两行显示
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
      overflow: hidden;

    }

    .remark {
      height: 17px;
      font-weight: 400;
      font-size: 12px;
      color: #000000;
      line-height: 17px;
      text-align: left;
      font-style: normal;
      margin-top: 4px;
    }

    .time {
      height: 17px;
      font-weight: 400;
      font-size: 12px;
      color: #A6A6A6;
      line-height: 17px;
      text-align: left;
      font-style: normal;
      margin-top: 4px;
    }

    .num {
      display: flex;
      align-items: center;
      position: absolute;
      bottom: 0;
      left: 14px;

      .num-item {
        display: flex;
        align-items: center;
        font-size: 14px;
        color: #777777;
        margin-right: 18px;

        &__icon {
          width: 24px;
          height: 24px;
          background: url("~@/assets/communication/icon_user.png") center/cover no-repeat;
        }
      }
      .num-item:nth-child(2) .num-item__icon {
        background-image: url("~@/assets/communication/icon_view.png");
      }
    }
  }

  .is-new {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 33px;
    height: 33px;
    overflow: hidden;

    &::before {
      content: "";
      position: absolute;
      bottom: 0;
      right: 0;
      width: 100%;
      height: 100%;
      transform: translate(50%, 50%)rotate(-45deg);
      background-color: #000;
    }

    &::after {
      content: "new";
      position: absolute;
      bottom: 0;
      right: 0;
      width: 100%;
      height: 100%;
      transform: translate(50%, 50%)rotate(-45deg);
      color: #fff;
      font-size: 8px;
      line-height: 11px;
      z-index: 1;
      display: flex;
      justify-content: center;
      align-items: flex-start;
    }
  }

  .type {
    position: absolute;
    top: 8px;
    left: 8px;
    height: 18px;
    line-height: 18px;
    background: #000;
    border-radius: 0px 0px 9px 0px;
    text-align: center;
    font-size: 10px;
    color: #fff;
    padding: 0 4px;
  }
}
</style>