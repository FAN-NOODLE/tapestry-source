<template>
  <div class="card">
      <div class="item-title" @click="detailItem(item.link_url)">
        {{ item.title }}
      </div>
      <div class="item-info" @click="detailItem(item.link_url)">
        <span class="item-txt">
          {{ item.description }}
        </span>
        <img class="item-img" :src="item.cover_url" alt="" />
      </div>
      <!-- 商品标签 -->
      <div v-if="!item.sku_list || item.sku_list.length === 0" class="tags h"></div>
      <div v-else-if="item.sku_list.length <=3" class="tags h">
        <el-tag v-for="(item, index) in item.sku_list" :key="index" type="gray">{{ item }}</el-tag>
      </div>
      <el-popover v-else placement="top" width="200px" trigger="hover">
        <div class="tags h" slot="reference">
          <el-tag v-for="(item, index) in item.sku_list.slice(0,3)" :key="index" type="gray">{{ item }}</el-tag>
          <span class="more">...</span>
        </div>
        <div class="tags w">
          <el-tag v-for="(item, index) in item.sku_list" :key="index" type="gray">{{ item }}</el-tag>
        </div>
      </el-popover>
      <div class="item-info">
        <el-checkbox :label="item.id" :key="index"></el-checkbox>
        <span class="date-time">{{
          item.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
        }}</span>
      </div>
      <div class="btns">
        <span v-if="item.is_kos_sync === 1" class="source-origin">CH同步</span>
        <el-button
          type="text"
          size="small"
          :style="{ color: item.source_status === 0 ? '#F56C6C' : '' }"
          @click="updateStatus(item)"
          >{{ item.source_status === 0 ? "禁用" : "启用" }}</el-button
        >
        <el-button
          size="small"
          type="text"
          class="unload"
          @click="editItem(item)"
          >编辑</el-button
        >
        <el-button
          style="color:#303133;"
          size="small"
          type="text"
          @click="delItem(item.id)"
          >删除</el-button
        >
        <el-button
          v-if="ep.isCoach && filterObj.source_group_id !== -2"
          type="text"
          size="small"
          :disabled="item.is_top === 0 && !item.can_top"
          @click="handleTop(item)"
          >{{ item.is_top === 0 ? "置顶" : "取消置顶" }}</el-button
        >
      </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Emit } from "vue-property-decorator";
import PerModule from "@/store/modules/permission";

@Component
export default class LinkMaterialItem extends Vue {
  @Prop() item!: any;
  @Prop() index!: number;
  @Prop() checkedMaterials!: any[];
  @Prop() filterObj!: any;

  get ep() {
    return PerModule.enterprise
  }

  @Emit('detail')
  detailItem(url: string) {
    return url;
  }

  @Emit('edit')
  editItem(item: any) {
    return item;
  }

  @Emit('delete')
  delItem(id: number) {
    return id;
  }

  @Emit('update-status')
  updateStatus(item: any) {
    return item;
  }

  @Emit('top')
  handleTop(item: any) {
    return item;
  }

  @Emit('checked-change')
  handleCheckedMaterialsChange(value: any) {
    return value;
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-tag.el-tag--gray {
  color: #666 !important;
  background-color: #F6F7F8 !important;
  border-color: #F6F7F8 !important;
}
.tags {
  margin-top: 5px;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  &.h {
    height: 32px;
  }

  &.w {
    width: 300px;
  }

  .more {
    color: #666;
    font-size: 12px;
  }
}
.card {
  float: left;
  width: 290px;
  min-height: 186px;
  background: #ffffff;
  border-radius: 4px;
  border: 1px solid #dcdfe6;
  margin-right: 20px;
  margin-bottom: 20px;
  padding: 10px;
  box-sizing: border-box;
  .item-title {
    width: 270px;
    height: 40px;
    font-size: 14px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    color: #303133;
    line-height: 20px;
    text-overflow: -o-ellipsis-lastline;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2; //行数需设置
    line-clamp: 2;
    -webkit-box-orient: vertical;
    &:hover {
      cursor: pointer;
    }
  }
  .item-info {
    display: flex;
    justify-content: space-between;
    margin-top: 5px;
    margin-bottom: 10px;
    color: #909399;
    &:hover {
      cursor: pointer;
    }
    .item-txt {
      width: 220px;
      height: 34px;
      font-size: 14px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      line-height: 17px;
      text-overflow: -o-ellipsis-lastline;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2; //行数需设置
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }
    .item-img {
      width: 50px;
      height: 50px;
      border-radius: 4px;
      img {
        width: 100%;
        height: 100%;
        border-radius: 4px;
        object-fit: cover;
      }
    }
    .date-time {
      font-size: 12px;
    }
    /deep/.el-checkbox__label {
      font-size: 0;
    }
  }
  .btns {
    border-top: 1px solid $color-second-line;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    .source-origin {
      margin-right: auto;
      color: #909399;
      font-size: 12px;
    }
    .del-btn {
      color: #f67270;
    }
  }
}
</style>
