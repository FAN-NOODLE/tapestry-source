<template>
  <div class="select">
    <div class="select" @click="handleClick">
      <slot :text="innerValue">{{ getName(innerValue)  }}</slot>
    </div>
    <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
      <div class="popup">
        <div class="popup-header">
          <div class="popup-title">{{ $t('storeVisit.popup.title.15') }}</div>
          <div class="popup-close" @click.stop="handleCancel"></div>
        </div>
        <div class="popup-tabs">
          <div class="tabs">
            <div class="tab" :class="{active: firstIndex === index}" v-for="item, index in tree" :key="index" @click="handleTabClick(index)">{{ item.name }}</div>
          </div>
        </div>
        <div class="popup-content">
          <div class="list">
            <div class="item" v-for="item,index in areaList" :key="index">
              <div class="item-header" :class="{active: item.key == innerValue}" @click="handleItemClick($event, item, index)">
                <div class="item-name">{{ getName(item.name) }}</div>
                <div v-if="item.children && item.children.length" :class="[item.isExpand?'icon-minus':'icon-plus']"></div>
                <div v-else class="icon-check"></div>
              </div>
              <div class="item-content">
                <div class="item-child" :class="{active: child.key == innerValue}" v-for="(child, subIndex) in cityList(item.children,item.name)" :key="subIndex" @click="handleItemClick($event, child, subIndex)">
                  <div class="item-name">{{ getName(child.name) }}</div>
                  <div class="icon-check"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </van-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import cityNameMap from '../city'


@Component({
  components: {}
})
export default class SelectDistrictPopup extends Vue {
  @Model('change', { type: String, default: '' }) value: string;
  @Prop({ type: Array, default: () => [] }) list: any[];

  
  visible: boolean = false;

  innerValue: string = '';

  tree: any[] = []
  // 默认选中的索引
  firstIndex = 0;

  get showText() {
    return ''
  }

  get areaList () {
    if(!this.tree[this.firstIndex]) return []
    return this.tree[this.firstIndex].children
  }

  get cityList() {
    return (list, name) => {
      if(['北京市', '天津市', '上海市', '重庆市'].includes(name)) {
        return [list[0]]
      } else {
        return list
      }
    }
  }



  get lang () {
    return this.$store.getters["lang/getLang"]
  }
  get getName() {
    return (name) => {
      if(this.lang === 'en-US') {
        // 通过下划线分组
        const arr = name.split('_')
        // console.log(arr);
        
        return arr.map(item => {
          if(item.indexOf(' Total') > -1) {
            return cityNameMap[item.replace(' Total', '')] + ' Total'
          }
          return cityNameMap[item] || item

        }).join('_')
      }
      
      return name
    }
  }


  @Watch('value', { immediate: true })
  onValueChange(val) {
    this.innerValue = val;
  }
  @Watch('list', { immediate: true })
  onTreeChange(val) {
    this.tree = val
  }

  handleClick() {
    this.visible = true;
    this.$nextTick(() => {
      this.setDefaultExpand();
    })
  }

  handleTabClick(index) {
    this.firstIndex = index;
    this.innerValue = ''
    this.$nextTick(() => {
      this.setDefaultExpand();
    })
  }
  /**
   * 取消
   */
  handleCancel() {
    this.visible = false;
    this.$emit('cancel');
  }

  /**
   * 确定
   */
  handleConfirm(value) {
    this.visible = false;
    this.$emit('change', value.key);
    this.$emit('confirm', value);
  }

  handleItemClick(e, item, index) {
    // this.handleConfirm(item);
    if(item.level == 2 && item.children && item.children.length > 0) {
      this.handleToggleExpand(e, item, index)
    } else {
      this.handleConfirm(item);
    }
  }

  handleToggleExpand(event, item, index) {
    event.stopPropagation();
    const isExpand = this.areaList[index].isExpand;
    // 关闭其他展开的节点
    this.areaList.forEach((item, i) => {
      if (i !== index) {
        this.$set(item, 'isExpand', false);
      } else {
        this.$set(item, 'isExpand', !isExpand);
      }
    })
    // 展开收起内容
    document.querySelectorAll('.item-content').forEach((item, i) => {
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
    console.log(this.value)
    const [a, b, c] = this.value.split('_');
    
    // this.firstIndex = this.tree.findIndex(item => item.name === a);
    const secondIndex = this.tree[this.firstIndex].children.findIndex(item => item.name === b.replace(' Total', ''));
    console.log(secondIndex)
    
    this.tree[this.firstIndex].children.forEach((item, index) => {
      if(index === secondIndex) {
        this.$set(item, 'isExpand', true);
        document.querySelectorAll('.item-content').forEach((item, i) => {
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

}
</script>
<style lang="scss" scoped>

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

.icon-check {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_check_blue.png) center/24px no-repeat;
}
.icon-down {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_down_circle.png) center/24px no-repeat;
}
.select {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.popup {
  position: relative;
  width: 100vw;
  height: 80vh;
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
  

  .popup-tabs {
    padding: 0 12px 12px 12px;
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
  }

  .popup-content {
    flex: 1;
    overflow-y: auto;
    .item {
      &:not(:last-child) {
        border-bottom: 1px solid rgba(248, 249, 250, 1);
      }
    }

    .item-header {
      height: 68px;
      padding: 0 24px;
      display: flex;
      align-items: center;
      justify-content: space-between;

      .item-name {
        font-size: 14px;
        color: rgba(0, 0, 0, 1);
      }
      .icon-check {
        opacity: 0;
        transition: all 0.3s;
      }

      &.active {
        .item-name {
          font-weight: 700;
          color: rgba(13, 110, 253, 1);
        }
        .icon-check {
          opacity: 1;
        }
      }
    }

    .item-content {
      padding: 0 24px;
      overflow: hidden;
      background-color: rgba(248, 249, 250, 1);
      transition: all 0.3s;

      .item-child {
        height: 64px;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .item-name {
          font-size: 14px;
          color: rgba(0, 0, 0, 1);
          transition: all 0.3s;
        }
        .icon-check {
          opacity: 0;
          transition: all 0.3s;
        }

        &.active {
          .item-name {
            font-weight: 700;
            color: rgba(13, 110, 253, 1);
          }
          .icon-check {
            opacity: 1;
          }
        }
      }

    }
  }
}
</style>