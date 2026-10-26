<template>
  <div class="function-introduce clearfix">
    <h4 class="title">公告</h4>
    <div class="announcement-list m-b-10">
      <div
        class="announcement-item clearfix"
        v-for="(item, index) in announceList"
        :key="index"
        @click="selectNews(item)"
      >
        <div class="item-title">{{ item.title }}</div>
        <div class="item-time">
          {{ item.released_at | formatDate("mm-dd") }}
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Bus from "@/components/bus";
import { AnnounceItem } from "./interface";
@Component({})
export default class AnnouncementList extends Vue {
  @Prop({ type: Array }) announceList!: Array<AnnounceItem>;
  private announcement = {};
  private selectNews(data: AnnounceItem) {
    const { content, released_at, title } = data;
    Object.assign(this.announcement, { content, released_at, title });
    Bus.$emit("showAnnounceMent", this.announcement);
  }
}
</script>
<style lang="scss" scoped>
.function-introduce {
  border-bottom: 1px solid $color-third-line;
  .title {
    margin: 20px 0 5px 0;
    font-size: 16px;
    color: $color-text-main;
    font-weight: 500;
    line-height: 26px;
    text-align: left;
  }
}
.announcement-list {
  width: 189px;
  height: 160px;
  overflow-y: auto;
  .announcement-item {
    text-align: left;
    cursor: pointer;
    margin-bottom: 5px;
    &:hover {
      .item-title {
        color: $color-main;
      }
    }
    .item-title {
      float: left;
      width: 134px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      color: $color-text-normal;
      line-height: 20px;
      font-size: 14px;
    }
    .item-time {
      float: right;
      color: $color-text-second;
      font-size: 14px;
      line-height: 20px;
      margin-right: 5px;
    }
  }
}
</style>
