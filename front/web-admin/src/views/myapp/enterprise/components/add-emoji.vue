<template>
  <el-popover placement="bottom-start" width="400" trigger="click">
    <div class="emj-boxs">
      <el-tabs tab-position="top" style="height: 200px;">
        <el-tab-pane
          v-for="(item, index) in emoji_list"
          :key="index"
          :label="item.type"
        >
          <div class="emj-box">
            <span
              class="emjs"
              v-for="(v, i) in item.list"
              :key="i"
              @click="addEmj(v)"
              >{{ v }}</span
            >
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
    <span class="text-head-button" slot="reference">
      <slot name="reference">
        <img
          class="image"
          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/face_icon.png"
          alt=""
        />
      </slot>
    </span>
  </el-popover>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
@Component({})
export default class extends Vue {
  private emoji_list: any = [];
  private addEmj(emj: any) {
    this.$emit("select", emj);
  }
  private emoji = require("./emoji.json");
  private mounted() {
    for (const key in this.emoji) {
      const arr: any = [];
      for (const em of this.emoji[key]) {
        arr.push(em);
      }
      const item = {
        type: key,
        list: arr
      };
      this.emoji_list.push(item);
    }
  }
}
</script>
<style lang="scss" scoped>
.emj-boxs {
  width: 400px;
  height: 300px;
}
.emj-box {
  height: 250px;
  overflow: auto;
  font-size: 16px;
  user-select: none;
}
.emjs {
  display: inline-block;
  width: 28px;
  height: 28px;
  padding: 3px;
  text-align: center;
  cursor: pointer;
}
.text-head-button {
  font-size: 14px;
  font-family: PingFangSC-Regular, PingFang SC;
  font-weight: 400;
  color: rgba(96, 98, 102, 1);
  cursor: pointer;
}
.image {
  width: 14px;
  cursor: pointer;
  height: 14px;
}
</style>
