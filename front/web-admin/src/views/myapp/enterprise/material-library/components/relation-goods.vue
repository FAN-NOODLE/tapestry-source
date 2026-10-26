<template>
  <el-row class="card" :key="item.hash || item.id">
    <el-row class="card-header">
      <div class="flex">
        <el-image :src="item.cover_url || item.video_cover_url || item.url" style="width: 40px; height: 40px;"></el-image>
        <span>{{item.name}}</span>
      </div>
      <select-goods v-if="!item.is_ai_loading" :value="item.relation_goods.filter(goods => goods.add_type === 1)" @confirm="handleSelectGoodsConfirm"></select-goods>
    </el-row>
    <el-table :data="item.relation_goods" 
      v-loading="item.is_ai_loading" 
      element-loading-text="AI识图加载中，请稍后..."
      element-loading-spinner="el-icon-loading"
      element-loading-background="rgba(255, 255, 255, 1)"
      empty-text="暂未识别到商品，您可手动添加商品"
      >
      <el-table-column label="商品名称" prop="name" min-width="100px">
        <template slot-scope="scope">
          <span v-if="!scope.row.name">--</span>
          <el-tooltip v-else :content="scope.row.name" placement="top">
            <div class="title">{{ scope.row.name }}</div>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="款号" prop="sku" />
      <el-table-column label="颜色" prop="color" />
      <el-table-column label="图片" prop="url">
        <template slot-scope="scope">
          <el-image :src="scope.row.url" class="goods-image" :preview-src-list="[scope.row.url]">
            <div slot="error" class="image-slot">
              <i class="el-icon-picture-outline"></i>
            </div>
          </el-image>
        </template>
      </el-table-column>
      <el-table-column label="置信度" prop="score">
        <template slot-scope="scope">
          {{ Number(scope.row.score) ? `${(scope.row.score * 100).toFixed(2)}%`  : '/' }}
        </template>
      </el-table-column>
      <el-table-column label="来源" prop="add_type">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.add_type === 1" type="success">手动添加</el-tag>
          <el-tag v-else type="warning">AI识图</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" prop="operation">
        <template slot-scope="scope">
          <!-- <el-button type="text" size="small" @click="handleDeleteGoods(scope.$index)">删除</el-button> -->
          <el-checkbox v-model="scope.row.checked"></el-checkbox>
        </template>
      </el-table-column>
      
    </el-table>
  </el-row>
</template>
<script lang="ts"> 
import { Component, Vue, Prop } from 'vue-property-decorator';
import SelectGoods from "./select-goods.vue";

@Component({
  components: {
    SelectGoods
  }
})
export default class extends Vue {
  //------------------------------------------------------prop-----------------------------------------------------------/
  @Prop({ default: () => {} }) item: any;
  
  //------------------------------------------------------method-----------------------------------------------------------/
  handleSelectGoodsConfirm(goods: any[]) {
    this.$emit('confirm', goods);
  }
  handleDeleteGoods(index: number) {
    this.$emit('delete', index);
  }
}
</script>
<style lang="scss" scoped> 

.card {
  width: 800px;
  margin-bottom: 10px;
  .card-header {
    padding: 10px;
    display: flex;
    align-items: center;
    background-color: #f2ebda;

    .flex {
      display: flex;
      align-items: center;
      gap: 10px;
      flex: 1;
    }
  }

  .goods-image {
    width: 60px;
    height: 60px;

    /deep/.image-slot {
      width: 60px;
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: #f5f7fa;
    }
  }

  .el-table {
    min-height: 200px;
  }
}

.title {
  height: 40px;
  line-height: 20px;
  // 两行省略
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>