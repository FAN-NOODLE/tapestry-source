<template>
  <div
    class="tags"
    :style="$attrs.isClose == 'no' ? 'marginLeft:0' : 'marginLeft:94px'"
  >
    <!--重复内容，可写组件，暂留-->
    <div
      class="tags-single"
      v-for="(item, index) in $attrs.tagsArr"
      :key="item.tag_group_id"
    >
      <dl class="tags-single-content">
        <dt>{{ item.tag_group_name }}</dt>
        <dd>
          <div
            v-for="(tag, subIndex) in item.tags"
            :key="tag.tag_id"
            class="tag-content-item"
          >
            <el-tag v-if="$attrs.isClose == 'no'" class="sub-tag">
              {{ tag.name }}
            </el-tag>
            <el-tag v-else class="sub-tag" closable @close="tagCloseFn(tag)">
              {{ tag.name }}
            </el-tag>
            <div class="sub-connect" v-if="subIndex != item.tags.length - 1">
              或
            </div>
          </div>
        </dd>
      </dl>
      <div
        class="tags-single-connect"
        v-if="index != $attrs.tagsArr.length - 1"
      >
        且
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";

@Component
export default class extends Vue {
  private tagCloseFn(tag: object) {
    this.$emit("deleteFn", { item: tag, type: this.$attrs.type });
  }
}
</script>
<style lang="scss" scoped>
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
.tags {
  // margin-left: 80px;
  overflow: hidden;
  margin-bottom: 10px;
  .tags-single {
    float: left;
    overflow: hidden;
    margin-bottom: 20px;
    &-content {
      float: left;
      border: 1px solid #dcdfe6;
      padding: 20px 20px 0;
      dt {
        font-size: 14px;
        color: $color-text-main;
        line-height: 20px;
        margin-bottom: 20px;
        font-weight: bold;
      }
      dd {
        overflow: hidden;
        .tag-content-item {
          float: left;
          overflow: hidden;
          margin-bottom: 20px;
          .sub-tag,
          .sub-connect {
            float: left;
          }
          .sub-connect {
            line-height: 32px;
            padding: 0 10px;
          }
        }
      }
    }
    &-connect {
      float: left;
      padding: 0 20px;
      padding-top: 46px;
    }
  }
}
</style>
