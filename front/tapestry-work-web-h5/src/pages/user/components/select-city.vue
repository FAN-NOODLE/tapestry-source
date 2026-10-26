<template>
  <div class="select-city-wrap">
    <div class="select-city" @click="visible = !visible">
      <div v-if="value" class="value">{{ value }}</div>
      <div v-else class="placeholder">更多城市</div>
      <!-- <div v-if="value" class="icon-close-circle" @click.stop="handleClear()"></div> -->
      <!-- <div class="icon-down"></div> -->
    </div>
    <div v-if="visible" class="select-popup">
      <div class="head" @click.stop="handleClear()">重置筛选</div>
      <div ref="bodyRef" class="body">
        <div v-for="item,index in options" :key="index" :ref="'id_'+item.letter" class="item">
          <div class="item-head">{{ item.letter }}</div>
          <div class="item-body">
            <div class="item-city" v-for="city, ci in item.list" :key="ci" @click="handleSelect(city)">{{city}}</div>
          </div>
        </div>
      </div>
      <div class="letter-list">
        <div v-for="item,index in options" :key="index" class="letter-item" @click="handleScrollTo(item)">{{ item.letter }}</div>
      </div>
    </div>
    <div v-if="visible" class="select-mask" @click="handleCancel"></div>
  </div>
</template>
<script lang="ts"> 
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import wxUtil from '@/utils/wx';
import { getQrcodeByUnionidApi } from "@/api/user";


@Component
export default class SelectCity extends Vue {
  @Prop({ type: Array, default: () => [] }) options: any
  @Prop({ type: String, default: '' }) value: string

  visible = false

  @Watch('visible') 
  async handleVisibleChange(val) {
    if (val) {
      document.body.classList.add('no-scroll')
    } else {
      document.body.classList.remove('no-scroll')
    }
  }

  handleScrollTo(item) {
    const dom = this.$refs.bodyRef as HTMLElement
    const itemDom = this.$refs[`id_${item.letter}`][0] as HTMLElement
    const h = document.body.getElementsByClassName('item-head')[0].clientHeight
    // 让dom 的内容滚动到指定位置
    dom.scrollTo({
      top: itemDom.offsetTop - 50,
      behavior: "smooth"
    });
  }

  handleSelect(item) {
    this.visible = false
    this.$emit('change', item)
  }
  handleClear() {
    this.visible = false
    this.$emit('change', '')
  }
  handleCancel() {
    this.visible = false
  }
}
</script>

<style>
.no-scroll {
  overflow: hidden;
  height: 100%;          /* 兼容 iOS 10/11 */
  position: fixed;       /* 防止部分安卓机仍可滑 */
  width: 100%;
}
</style>

<style lang="scss" scoped>

.icon-address {
  width: 16px;
  height: 22px;
  background: url('~@/assets/user/icon_address.png') center/contain no-repeat;

}
.icon-down {
  width: 10px;
  height: 10px;
  background: url('~@/assets/user/icon_down.png') center/contain no-repeat;
}
.icon-search {
  width: 14px;
  height: 14px;
  background: url('~@/assets/user/icon_search.png') center/contain no-repeat;
}

.icon-close-circle {
  width: 24px;
  height: 24px;
  background: url("~@/assets/store-visit/icon_close_circle.png") center/14px no-repeat;
}
.select-city {
  width: 90px;
  height: 30px;
  // border: 1px solid #C5C5C5;
  font-size: 12px;
  display: flex;
  align-items: center;
  padding: 0 12px 0 10px;
  background-color: #000;
  color: #fff;
  text-align: center;
  line-height: 30px;

  .value {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    height: 100%;
  }


  .placeholder {
    // color: #949494;
    flex: 1;
  }
}

.select-city-wrap {
  position: relative;

  .select-mask {
    position: fixed;
    top: 112px;
    left: 0;
    // transform: translateX(-50%);
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.5);
    z-index: 2;

  }

  .select-popup {
    position: fixed;
    top: 112px;
    left: 20px;
    width: 335px;
    height: calc(100vh - 150px);
    background: #fff;
    z-index: 3;
    display: flex;
    flex-direction: column;
    padding-bottom: 10px;


    .head {
      height: 37px;
      // background-color: #F0F0F0;
      padding: 0 10px;
      font-size: 14px;
      line-height: 37px;
      text-decoration: underline;
    }
    .body {
      width: 100%;
      flex: 1;
      margin-top: 8px;
      overflow-y: scroll;

      .item-head {
        height: 37px;
        background-color: #F0F0F0;
        padding: 0 10px;
        position: sticky;
        top: 0;
        z-index: 1;
        font-size: 14px;
        line-height: 37px;
      }
      .item-city {
        height: 37px;
        padding: 0 10px;
        font-size: 14px;
        line-height: 37px;
      }
    }

    .letter-list {
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
      z-index: 2;
      .letter-item { 
        padding: 0 10px;
        height: 23px;
        line-height: 23px;
        font-size: 14px;
      }
    }
  }
}

</style>