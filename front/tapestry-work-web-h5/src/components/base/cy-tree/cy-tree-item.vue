<template>
  <li class="tree-item" :class="{ 'disable-item': !data.is_show }">
    <div
      @click="toggle"
      class="left-span p-l-30"
      :class="{ 'select-name': onNodeId === data.id }"
    >
      <span v-if="hasChild()" class="open-icon">
        <span :class="[open ? 'open-icon-icon' : '']"
          ><van-icon
            name="play"
            :color="onNodeId === data.id ? '#000000' : '#777777'"
        /></span>
      </span>
      <span class="name-text">{{ data.name }}</span>
      <span v-if="data.is_show" class="radio-icon" @click.stop="onChecked">
        <van-icon v-if="onNodeId !== data.id" name="circle" />
        <van-icon v-else name="checked" color="#000000" />
      </span>
    </div>
    <ul v-show="open" v-if="hasChild()">
      <cy-tree-item
        v-for="(item, index) in data.children"
        :data="item"
        :key="index"
      ></cy-tree-item>
    </ul>
  </li>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch, Emit } from "vue-property-decorator";
import { Icon } from "vant";
interface T {
  children: any;
  name: string;
  id: number;
  is_show: any;
}
@Component({
  name: "cy-tree-item",
  components: {
    [Icon.name]: Icon,
  },
})
export default class extends Vue {
  @Prop() data: T;
  // @Prop() onNodeId: number;
  private open = false;
  private itemId = 0;
  // 记录本次单选已选择的数据
  public get onNodeId(): any {
    return this.$store.getters["materialLibrary/getNowChecked"].id || 0;
  }
  private hasChild() {
    return this.data.children && this.data.children.length;
  }
  private toggle() {
    if (this.hasChild()) {
      this.open = !this.open;
    }
  }
  private onChecked() {
    if (!this.data.is_show) {
      return;
    }
    this.$store.commit("materialLibrary/setNowChecked", this.data);
    this.$emit("onChecked", this.data);
  }
}
</script>
<style lang="scss" scoped>
.tree-item {
  // line-height: 40px;
  font-size: 14px;
  color: #777777;
  padding-left: 10px;
  &.disable-item {
    color: #cccccc;
  }
  .left-span {
    padding: 10px 0;
    position: relative;
    padding-left: 12px;
    &.p-l-30 {
      padding-left: 30px;
    }
    &.select-name {
      font-weight: 600;
      color: #000000;
      background-color: #f4f7fb;
    }
    .open-icon {
      color: #666666;
      position: absolute;
      left: 8px;
      top: 50%;
      transform: translateY(-50%);

      .open-icon-icon {
        display: block;
        transform: rotate(90deg);
        .van-icon {
          font-size: 14px;
        }
      }
    }
    .name-text {
      display: inline-block;
      width: 88%;
    }
    .radio-icon {
      display: block;
      color: #666666;
      position: absolute;
      right: 8px;
      top: 50%;
      transform: translateY(-39%);
      .van-icon {
        font-size: 16px;
      }
    }
  }
}
</style>
