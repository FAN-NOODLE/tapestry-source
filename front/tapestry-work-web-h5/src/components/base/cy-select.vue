<template>
  <div>
    <div class="cy-select">
      <div
        style="
          height: 32px;
          padding: 5px 8px;
          display: flex;
          font-size: 14px;
          justify-content: space-between;
          background: #f4f5f8;
        "
        @click.stop="toggleSelect"
        :class="{ select: selectIcon == 'arrow-up' }"
      >
        <span
          style="color: #a6a6a6"
          v-if="placeholder && selectedObj == null"
          >{{ placeholder }}</span
        >
        <span v-else>{{ selectedObj ? selectedObj.label : "" }}</span>
        <van-icon :name="selectIcon" style="top: 20%" />
      </div>
      <ul
        class="panel"
        v-if="selectIcon == 'arrow-up'"
        @click.stop="choiceLi($event)"
      >
        <li
          v-for="(item, index) in options"
          :key="index"
          :data-ind="index"
          :class="{
            active: item.value == (selectedObj ? selectedObj.value : null),
          }"
        >
          {{ item.label }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch, Emit } from "vue-property-decorator";
import { Icon } from "vant";
@Component({
  name: "cy-select",
  components: {
    [Icon.name]: Icon,
  },
})
export default class extends Vue {
  @Prop({ default: [] }) options: Array<any>;
  @Prop({ default: "" }) placeholder: any;
  private selectIcon = "arrow-down";
  private selectIconList = {
    "arrow-down": "arrow-up",
    "arrow-up": "arrow-down",
  };
  @Watch("value")
  setValue(val: string | any): void {
    this.selectedObj = this.options.filter((v: any) => {
      return v.value == val;
    })[0];
  }
  private value: any = "";
  private selectedObj: any = null;
  @Emit("select")
  private choiceLi(e) {
    this.selectedObj = this.options[e.target.dataset.ind];
    this.value = this.selectedObj.value;
    this.toggleSelect();
    return this.value;
  }
  private toggleSelect() {
    this.selectIcon = this.selectIconList[this.selectIcon];
  }
  private closePanel() {
    this.selectIcon = "arrow-down";
  }
  private mounted() {
    document.body.addEventListener("click", this.closePanel);
  }
  private destoryed() {
    document.body.removeEventListener("click", this.closePanel);
  }
}
</script>
<style lang="scss" scoped>
.cy-select {
  position: relative;
  width: 50%;
  .select:after {
    content: "";
    position: absolute;
    width: 0;
    height: 0;
    border: 8px solid transparent;
    border-bottom: 8px solid #fff;
    bottom: -13px;
    right: 20px;
  }
  .panel {
    width: 100%;
    background: #fff;
    font-size: 14px;
    position: absolute;
    top: 42px;
    height: 210px;
    overflow: auto;
    li:hover {
      background: #f4f5f8;
      color: #0095f6;
    }
    li:nth-child(even) {
      background: #f4f5f8;
    }
    .active {
      color: #0095f6;
    }
    li {
      padding: 10px 20px;
    }
  }
}
</style>
