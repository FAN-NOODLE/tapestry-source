<template>
  <div>
    <div>
      <el-button type="primary" @click="handleOpenClick"  size="small"
        >添加商品</el-button
      >
    </div>
    <el-dialog
      id="crs-scrm"
      class="select-goods-dialog"
      :title="title"
      :visible.sync="visible"
      width="1000px"
      @close="handleClose"
      @closed="handleClosed"
      @open="handleOpen"
      append-to-body
    >
      <div class="select-goods-container">
        <!-- 左侧商品列表 -->
        <div class="goods-list-section">
          <div class="section-header">
            <el-input
              v-model="searchKeyword"
              :placeholder="`请输入${searchType === 'sku_color' ? '商品款号' : '商品名称'}`"
              style="width: 400px"
              clearable
              @keyup.enter.native="handleSearch"
            >
              <template slot="prepend">
                <el-select v-model="searchType" placeholder="请选择" style="width: 120px;" @change="searchKeyword = ''">
                  <el-option label="商品款号" value="sku_color"></el-option>
                  <el-option label="商品名称" value="product_name"></el-option>
                </el-select>
              </template>
            </el-input>
            <el-button class="m-l-10" type="primary" size="small" :disabled="!searchKeyword" @click="handleSearch">立即查询</el-button>
          </div>
          <!-- <div v-if="goodsList.length === 0 && searchKeyword" class="no-data">未搜索到商品</div> -->
          <div v-if="showContent" class="goods-list">
            <el-table 
            ref="goodsTable" 
            :data="goodsList" 
            style="width: 100%" 
            row-key="upc_code" 
            height="500" 
            stripe
            @select="handleSelect"
            @select-all="handleSelectAll"
            >
              <el-table-column label="选择" type="selection" width="55" reserve-selection></el-table-column>
              <el-table-column label="商品名称" prop="product_name" min-width="200">
                <template slot-scope="scope">
                  <div class="item-title">{{ scope.row.product_name || '--' }}</div>
                </template>
              </el-table-column>
              <el-table-column label="款号" prop="sku_item"  />
              <el-table-column label="颜色" prop="color"  />
              <el-table-column label="图片" width="120">
                <template slot-scope="scope">
                  <el-image :src="scope.row.brand_source_list[0]?.url" class="goods-image" :preview-src-list="scope.row.brand_source_list.map(item => item.url)">
                    <div slot="error" class="image-slot">
                      <i class="el-icon-picture-outline"></i>
                    </div>
                  </el-image>
                </template>
              </el-table-column>
              <!-- 缺省 -->
              
              <div slot="empty">
                <img
                  class="empty-img"
                  v-if="!showSearchIcon"
                  src="@/assets/images/empty-box.png"
                  alt=""
                />
                <div class="empty-box" v-if="!showSearchIcon">
                  <div class="desc">请输入查询条件，检索后展示商品列表</div>
                </div>
                <img
                  class="empty-img"
                  v-if="showSearchIcon"
                  src="@/assets/images/empty-search.png"
                  alt=""
                />
                <div class="empty-box" v-if="showSearchIcon">
                  <div class="desc">暂无匹配数据，请调整查询条件后重试</div>
                </div>
              </div>
            </el-table>
            <!-- <div class="mask"></div> -->
          </div>
        </div>

        <!-- 右侧已选列表 -->
        <div v-if="showContent" class="selected-list-section">
          <div class="section-header">
            <div>已选商品 ({{ selectedGoodsList.length }}个)</div>
            <el-button
              type="text"
              @click="handleDeleteAll"
              :disabled="selectedGoodsList.length === 0"
            >
              <i class="el-icon-delete"></i>
              清空已选
            </el-button>
          </div>
          <div class="selected-scroll">
            <div class="selected-list">
              <div
                v-for="(goods, index) in selectedGoodsList"
                :key="goods.upc"
                class="selected-item"
              >
                <div class="selected-info">
                  <div class="selected-details">
                    <div>
                      <i class="el-icon-goods"></i>
                    </div>
                    <div class="selected-name">{{ goods.sku }}_{{ goods.color }}</div>
                  </div>
                </div>
                <el-button
                  type="text"
                  size="small"
                  @click="handleDelete(index)"
                >
                  <i class="el-icon-delete"></i>
                </el-button>
              </div>
              <div v-if="selectedGoodsList.length === 0" class="empty-selected">
                暂无已选商品
              </div>
            </div>
          </div>
        </div>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleConfirm" :disabled="!showContent">
          确定
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop, Emit, Watch } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";

interface Goods {
  upc_code: string;
  [key: string]: any;
}


@Component
export default class SelectGoods extends Vue {
  @Prop({ default: "添加关联商品" }) title!: string;
  @Prop({ default: () => [] }) value!: Goods[];
  @Prop({ default: false }) disabled!: boolean;

  visible = false;
  searchType = "sku_color";
  loading = false;
  last_code:string = null
  hasMore = false
  maxSize = 10
  showContent = false
  showSearchIcon = false
  private searchKeyword = "";
  private goodsList: Goods[] = [];
  private selectedGoodsList: any[] = [];


  mounted() {}

  handleOpenClick() {
    if(this.disabled) {
      this.$message.warning("请先上传素材");
      return
    }

    this.visible = true
  }

  async initGoodsList() {
    if(this.loading) return
    this.loading = true;
    const params = {
      size: 20,
      last_code: this.last_code,
      [this.searchType]: this.searchKeyword,
      include_none_image_product: true
    }
    // 模拟从接口获取商品列表
    const res = await Api.searchGoodsListApi(params);
    const list = res.data || []
    this.goodsList = this.last_code ? [...this.goodsList, ...list] : [...list];
    const last_code = this.goodsList[this.goodsList.length - 1]?.id
    this.loading = false;
    this.hasMore = list.length === params.size;
    
    this.last_code = last_code,
    this.showContent = true
    this.showSearchIcon = true

    


    // 监听 el-table__body-wrapper is-scrolling-none 滚动触底
    this.$nextTick(() => {
      if(this.$refs.goodsTable) {
        this.goodsList.forEach(item => {
          if(this.selectedGoodsList.some(goods => goods.upc === item.upc_code)){
            (this.$refs.goodsTable as any).toggleRowSelection(item, true);
          } else {
            (this.$refs.goodsTable as any).toggleRowSelection(item, false);
          }
        })
        const bodyWrapper = (this.$refs.goodsTable as any).$el.querySelector('.el-table__body-wrapper');
        // 不要重复添加滚动监听
        if (bodyWrapper && !bodyWrapper.getAttribute('scroll-listener')) {
          bodyWrapper.setAttribute('scroll-listener', 'true');
          bodyWrapper.addEventListener('scroll', (e) => {
            
            const { scrollTop, scrollHeight, clientHeight } = e.target;
            if (scrollTop + clientHeight + 10 >= scrollHeight) {
              this.loadMore();
              console.log('触底了');
            }
          });
        }
      }
    })
  }

  handleSearch() {
    this.last_code = null
    this.initGoodsList();
  }
  /**
   * @description 加载更多
   */
  loadMore() {
    if (this.hasMore) {
      this.initGoodsList();
    }
  }

  handleSelect(selection, row:any) {
    console.log(selection, row);
    // 判断当前项是否在已中中选中了，如果在，就从已选中中删除除，否则，就添加到已选中中去
    if(selection.includes(row)){
      if(this.selectedGoodsList.length >= this.maxSize) {
        (this.$refs.goodsTable as any).toggleRowSelection(row, false)
        this.$message({
          message: `已选择商品已达最大数量${this.maxSize}`,
          type: 'warning'
        })
        return
      }
      this.selectedGoodsList.push(this.transform(row))
    } else {
      this.selectedGoodsList = this.selectedGoodsList.filter(goods => goods.upc !== row.upc_code)
    }
  }
  handleSelectAll(selection) {
    console.log(selection);
    if(selection.length === 0) {
      const ids = this.goodsList.map(item => item.upc_code)
      this.selectedGoodsList = this.selectedGoodsList.filter(goods => !ids.includes(goods.upc))
      return
    }
    
    const selectedIds = this.selectedGoodsList.map(item => item.upc)

    this.goodsList.forEach(item => {
      if(!selectedIds.includes(item.upc_code)) {
        this.selectedGoodsList.push(this.transform(item))
      }
    })
}

  transform(item:any) {
    return {
        upc: item.upc_code,
        sku: item.sku_item,
        color: item.color,
        url: item.brand_source_list[0]?.url||'',
        name: item.product_name|| '',
        add_type: 1,
        score: '',
        checked: true,
      }
  }

  private handleDelete(index) {
    const [row] = this.selectedGoodsList.splice(index, 1);
    const deletedGoods = this.goodsList.find(item => item.upc_code === row.upc);
    (this.$refs.goodsTable as any).toggleRowSelection(deletedGoods, false)

  }

  private handleDeleteAll() {
    this.selectedGoodsList = [];
    (this.$refs.goodsTable as any).clearSelection()
  }

  handleOpen() {
    this.selectedGoodsList = [...this.value];
    this.showContent = this.value.length > 0
  }

  handleClosed() {
    this.goodsList = []
    this.selectedGoodsList = []
    this.searchKeyword = ""
    this.last_code = null
    this.showContent = false
    this.showSearchIcon = false
  }


  @Emit('close')
  private handleClose() {
    this.visible = false;
    return;
  }

  @Emit('confirm')
  private handleConfirm() {
    if(this.selectedGoodsList.length > this.maxSize) {
      this.$message({
        message: `选择商品数量不能超过${this.maxSize}`,
        type: 'warning'
      })
      return
    }
    this.visible = false;
    return this.selectedGoodsList;
  }
}
</script>

<style lang="scss" scoped>
.empty-box {
  line-height: 40px;
}


.goods-image {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}
// .image-slot {
//   display: flex;
//   align-items: center;
//   justify-content: center;
//   font-size: 24px;
//   color: #909399;
//   background-color: #f5f7fa;
//   width: 50px;
//   height: 50px;
// }
.select-goods-container {
  display: flex;

  .goods-list-section {
    flex: 1;
    overflow: hidden;
    border-right: 1px solid #ebeef5;
    margin-right: 10px;
    padding-right: 10px;

    .goods-list { 
      margin-top: 10px;
      position: relative;
      .mask {
        position: absolute;
        top: 0;
        left: 0;
        width: 50px;
        height: 36px;
        background-color: #ebebeb;
        z-index: 5;

      }
    }
    .pagination {
      margin-top: 10px;
      text-align: right;
    }
  }
  .selected-list-section {
    width: 220px;
    height: 550px;
    display: flex;
    flex-direction: column;

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px;
      border-bottom: 1px solid #ebeef5;
    
    }

    .selected-scroll {
      flex: 1;
      overflow: hidden;
    }

    .selected-list {
      height: 100%;
      overflow: auto;
      // padding: 10px;

      .selected-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 10px 0;
        border-bottom: 1px solid #ebeef5;
      }

      .selected-details {
        display: flex;
        align-items: center;
      }
      .selected-name {
        margin-left: 5px;
      }
      .el-icon-delete {
        margin-left: 5px;
      }
    }

    .empty-selected{
      text-align: center;
      color: #909399;
      text-align: center;
      padding: 100px 0;
    }


  }
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 100px 0;
  width: 100%;
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 10px;
}
</style>

<style lang="scss">
#crs-scrm.select-goods-dialog .el-dialog {
  max-height: none !important;
}
</style>