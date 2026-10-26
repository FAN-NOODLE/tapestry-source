<template>
  <layout-list>
    <template slot="table">
      <div v-for="(item, index) in config" :key="index">
        <!-- 模块标题 -->
        <div class="model-title">{{ item.name }}</div>

        <!-- 模块内容 -->
        <div class="model-container">
          <div
            class="box"
            v-for="(model, key) in item.list"
            :key="key"
            @click="goModel(model.url)"
          >
            <div class="box-icon">
              <img :src="model.icon" />
            </div>
            <div class="box-content">
              <div>{{ model.name }}</div>
              <div>{{ model.describe }}</div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { weJourneyConfig } from "../config";

@Component({})
export default class WeJourney extends Vue {
  /** 默认数据 */
  private config = weJourneyConfig;

  /** 跳转到对应页面 */
  private goModel(path: string) {
    this.$router.push(path);
  }
}
</script>
<style lang="scss" scoped>
.model-title {
  font-size: 16px;
  font-weight: bold;
  position: relative;
  margin-left: 10px;
  margin-bottom: 15px;
  &::after {
    position: absolute;
    content: "";
    left: -15px;
    top: 2px;
    height: 15px;
    width: 5px;
    background: #000000;
  }
}
.model-container {
  display: grid;
  margin-bottom: 15px;
  grid-template-columns: repeat(3, 1fr);
  grid-gap: 30px;
  .box {
    background: #f9f9f9;
    height: 120px;
    padding: 30px;
    cursor: pointer;
    display: flex;
    align-items: center;
    box-sizing: border-box;
    &-icon {
      width: 47px;
      height: 47px;
      margin-right: 5px;
      align-self: center;
      img {
        width: 100%;
        height: 100%;
        border-radius: 50%;
      }
    }
    &-content {
      flex: 1;
      font-size: 12px;
      div:first-child {
        font-weight: bold;
        font-size: 16px;
        margin-bottom: 15px;
      }
    }
  }
}
</style>
