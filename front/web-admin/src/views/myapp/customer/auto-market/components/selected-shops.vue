<template>
  <el-table ref="multipleTable" :data="list" stripe max-height="500">
    <el-table-column prop="id" label="SPU ID"></el-table-column>
    <el-table-column prop="goods" label="商品">
      <template slot-scope="scope">
        <div class="aligns-center">
          <!--图片-->
          <el-image fit="cover"
            style="width: 50px; height: 50px"
            :src="scope.row.spu_source ? scope.row.spu_source.url : ''"
          
          ></el-image>
          <el-tooltip
            :disabled="scope.row.spu_name.length > 33 ? false : true"
            :content="scope.row.spu_name"
            placement="bottom"
            effect="light"
          >
            <span class="m-l-10 goods-name">{{ scope.row.spu_name }}</span>
          </el-tooltip>
        </div>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="100" v-if="operationBtn">
      <template slot-scope="scope">
        <el-button @click="delFn(scope.$index)" type="text" size="small"
          >删除</el-button
        >
      </template>
    </el-table-column>
  </el-table>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
@Component({})
export default class extends Vue {
  /************data**************/
  @Prop() list!: Array<any>;
  @Prop({ default: true }) operationBtn!: boolean;
  delFn(ind: number) {
    this.list.splice(ind, 1);
    this.$emit("getShopChange", this.list);
  }
}
</script>
<style lang="scss">
.el-table--scrollable-y {
  .el-table__body-wrapper::-webkit-scrollbar {
    width: 6px !important;
  }
}
</style>
<style lang="scss" scoped>
.el-range-editor {
  width: 280px !important;
}
.input-with-select {
  .el-select {
    width: 100px !important;
  }
}
.el-button--text {
  color: #f67170;
}
.el-button--text:hover {
  color: #f56c6c;
}
.marginL10 {
  margin-left: 10px;
}
/deep/.gutter {
  width: 0px;
  display: none;
}
.table-body {
  display: flex;
  flex-direction: column;
  /deep/.el-select .el-input {
    width: 110px;
  }
  /deep/.el-table__body-wrapper {
    height: 304px !important;
  }
  .table-empty-label {
    display: flex;
    height: 340px;
    flex-direction: column;
    img {
      width: 124px;
      margin: 80px auto 0;
    }
    p {
      line-height: 20px;
      margin-top: 14px;
    }
  }
}
.aligns-center {
  display: flex;
  align-items: center;
  .goods-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }
}
</style>
