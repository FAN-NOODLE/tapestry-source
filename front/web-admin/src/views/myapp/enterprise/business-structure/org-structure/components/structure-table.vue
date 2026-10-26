<template>
  <div class="structure-table">
    <div class="header">
      <div class="title">{{ title }}</div>
      <associate-store :treeList="treeList" @update="handleUpdate"></associate-store>
    </div>
    <div class="content">
      <el-table :data="list" stripe>
        <el-table-column prop="id" label="唯一ID" width="100">
        </el-table-column>
        <el-table-column prop="store_code" label="门店编码" min-width="100">
        </el-table-column>
        <el-table-column prop="store_name" label="门店名称" min-width="140" show-overflow-tooltip>
        </el-table-column>
        <el-table-column prop="path_name" label="所属分区" min-width="200" show-overflow-tooltip>
        </el-table-column>
        <el-table-column prop="closed_at" label="门店状态" width="100">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.closed_at > 0" size="mini" type="danger">闭店</el-tag>
            <el-tag v-else size="mini">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="最近更新时间" width="200">
          <template slot-scope="scope">
            <span>{{
              scope.row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss")
            }}</span>
          </template>
        </el-table-column>
        <!-- <el-table-column prop="updated_at" label="操作" width="100" fixed="right">
          <template slot-scope="scope">
            <el-button type="text" size="mini" @click="handleRemove(scope.row)">移动</el-button>
          </template>
        </el-table-column> -->
      </el-table>
    </div>
    <div class="footer">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
      >
      </el-pagination>
    </div>

    <structure-select ref="select" :treeList="treeList" :visible.sync="visible" @confirm="handleSuccess" />
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import AssociateStore from './associate-store.vue'
import Api from '@/api/enterprise/org-structure';
import StructureSelect from './structure-select.vue'

@Component({
  components: {
    AssociateStore,
    StructureSelect
  }
})
export default class extends Vue {
  @Prop({ default: "全部门店" }) title: string;
  @Prop({ default: () => [] }) data: any[];
  @Prop({ default: () => [] }) treeList: any[];
  visible = false
  /** 分页参数 */
  meta = {
    page: 1,
    size: 10,
  };

  get total() {
    return this.data.length
  }

  /** 列表数据 */
  get list() {
    const { page, size } = this.meta;
    return this.data.slice((page - 1) * size, page * size);
  }

  @Watch('total')
  onTotalChange() {
    this.meta.page = 1;
  }

  /**
   * 分页
   * @param val 
   */
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
  }

  /**
   * 分页
   * @param val 
   */
  handleCurrentChange(val: number) {
    this.meta.page = val;
  }

  handleRemove(value: any) {
    (this.$refs.select as any).show({
      parent_id: value.parent_id,
      id: value.id
    });
  }

  handleSuccess(params: any) {
    Api.createDashboardTree(params)
      .then((res: any) => {
        if(res.code === 200) {
          this.visible = false
          this.$emit("update");
        } else {
          this.$message.error(res.message)
        }
      })
  }

  handleUpdate() {
    this.$emit("update")
  }
}
</script>

<style lang="scss" scoped>
.structure-table {
  flex: 1;
  overflow: hidden;
  
  .header {
    display: flex;
    align-items: center;

    .title {
      font-weight: bold;
      flex: 1;
    }
  }

  .content {
    margin-top: 10px;
  }

  .footer {
    text-align: right;
    margin-top: 10px;
  }
}
</style>