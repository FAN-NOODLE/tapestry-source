<template>
  <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
    <div class="popup">
      <div class="popup-header">
        <div v-if="showBack" class="icon-left" @click="handleBack"></div>
        <div class="popup-title">{{getName(title)}} {{ $t('storeVisit.tips.31') }}</div>
        <div class="popup-close" @click.stop="handleCancel"></div>
      </div>
      <div class="popup-content">
        <div class="list">
          <div class="item" v-for="item,index in list" :key="index" @click="handleClick(item)">
            <div class="item-name">{{ getName(item.name) }}</div>
            <div class="item-label">{{ $t('storeVisit.tips.26') }}:</div>
            <div class="item-value">{{ item.total_count }}</div>
            <div class="icon-right"></div>
          </div>
        </div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import cityNameMap from '../city'

@Component({

})
export default class SelectStoreListPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: Object, default: () => {} }) defaultData: any;


  visible: boolean = false;
  
  list: any[] = []
  title: string = '';
  showBack: boolean = false;


  get lang () {
    return this.$store.getters["lang/getLang"]
  }
  get getName() {
    return (name) => {
      if(this.lang === 'en-US') {
        if(name.indexOf(' Total') > -1) {
          return cityNameMap[name.replace(' Total', '')] + ' Total'
        }
        return cityNameMap[name] || name
      }
      
      return name
    }
  }
  
  
  @Watch('show', { immediate: true })
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.title = this.defaultData.name
      const children = this.defaultData._children || this.defaultData.children || []
      
      // 老版本
      // this.list = children.reduce((prev, cur) => {
      //   console.log(cur);
        
      //   if(cur.children) {
      //     prev.push(...cur.children)
      //   } else {
      //     prev.push(cur)
      //   }
      //   return prev;
      // }, [])

      // 新版本， 如果是四大直辖市 直接获取门店
      const name = this.defaultData.name.replace(' Total', '')
      if(['北京市', '天津市', '上海市', '重庆市'].includes(name)) {
        this.list = children.reduce((prev, cur) => {
          console.log(cur);
          
          if(cur.children) {
            prev.push(...cur.children)
          } else {
            prev.push(cur)
          }
          return prev;
        }, [])
      } else {
        this.list = children
      }

    }
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
  }

  handleClick(item) {
    // this.$emit('update:show', false);
    console.log(item);
    if(item.level === 3) {
      this.showBack = true
      this.title = item.name
      this.list = item.children
    } else {
      this.$emit('confirm', item);
    }
  }

  handleBack() {
    this.showBack = false
    this.title = this.defaultData.name
    this.list = this.defaultData._children || this.defaultData.children || []
  }
}
</script>
<style lang="scss" scoped>

.icon-left {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_left.png) center/100% no-repeat;
}
.icon-right {
  width: 24px;
  height: 24px;
  background: url('~@/assets/store-visit/icon_right_black.png') center/100% no-repeat
}

.popup {
  position: relative;
  width: 100vw;
  background: #FFFFFF;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    display: flex;
    align-items: center;
    
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

  .popup-content {
    background: rgba(248, 249, 250, 1);
    max-height: 450px;
    overflow-y: auto;
    padding-bottom: calc(env(safe-area-inset-bottom));
    padding-bottom: calc(constant(safe-area-inset-bottom));

    .list {
      padding: 0 24px;

      .item {
        height: 64px;
        display: flex;
        align-items: center;

        &:not(:last-child) {
          border-bottom: 0.5px solid var(--Grey-200, rgba(233, 236, 239, 1))
        }

        .item-name {
          font-weight: 400;
          font-size: 14px;
          color: #000000;
          flex: 1;
          transition: all 0.3s;
        }
        .item-label {
          font-weight: 400;
          font-size: 14px;
          color: rgba(153, 153, 153, 0.8);
          margin-right: 8px;
        }
        .item-value {
          font-weight: 400;
          font-size: 14px;
          color: rgba(0, 0, 0, 1);
          margin-right: 8px;
        }
      }
    }
  }
}
</style>