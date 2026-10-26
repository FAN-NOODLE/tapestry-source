<template>
  <div>
    <el-button size="mini" type="primary" :disabled="total==0" @click="handleClick">待关联门店({{ total }})</el-button>
    <el-dialog :title="title" :visible.sync="visible" @close="handleCancel">
      <el-table :data="list" height="500px">
        <el-table-column prop="id" label="唯一ID" width="100">
        </el-table-column>
        <el-table-column prop="third_store_code" label="门店编码">
        </el-table-column>
        <el-table-column prop="store_name" label="门店名称">
        </el-table-column>
        <el-table-column prop="closed_at" label="门店状态">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.closed_at > 0" size="mini" type="danger">闭店</el-tag>
            <el-tag v-else size="mini">正常</el-tag>
          </template>
        </el-table-column>
        <!-- <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-button type="text" size="mini" @click="handleAssociate(scope.row)">关联</el-button>
          </template>
        </el-table-column> -->
      </el-table>
      <!-- <div slot="footer" class="dialog-footer">
        <el-button type="primary" @click="handleConfirm">关闭</el-button>
      </div> -->
    </el-dialog>
    <structure-select ref="select" :visible.sync="visible2" :treeList="treeList" @confirm="handleSuccess" />
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import Api from '@/api/enterprise/org-structure';
import StructureSelect from './structure-select.vue'
@Component({
  components: {
    StructureSelect
  }
})
export default class extends Vue {
  @Prop({ default: () => [] }) treeList: any[];
  
  /** 弹框控制 */
  visible = false;
  visible2 = false;
  /** 弹框标题 */
  title = '待关联门店';
  /** 待关联门店 */
  list = [];
  currentId = 0;
  /** 总条数 */
  get total() {
    return this.list.length;
  }
  created() {
    this.getData();
  }

  getData() {
    Api.getDashboardStoreList({}).then((res) => {
      this.list = res.data||[];
    });
  }

  handleClick() {
    this.visible = true;
    this.getData();
  }

  handleCancel() {
    this.visible = false;
  }

  handleConfirm() {
    this.visible = false;
  }

  handleAssociate(value: any) {
    (this.$refs.select as any).show({
      parent_id: '',
      third_id: value.wechat_department_id
    });
  }

  handleSuccess(params: any) {
    Api.createDashboardTree(params)
      .then((res: any) => {
        if(res.code === 200) {
          this.visible2 = false
          this.getData()
          this.$emit('update')
        } else {
          this.$message.error(res.message)
        }
      })
  }

}
</script>

<style lang="scss" scoped>
::v-deep {
  .el-dialog {
    max-height: none !important;
  }
}
</style>