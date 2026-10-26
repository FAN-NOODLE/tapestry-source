<template>
  <el-table
    :data="list"
    style="width: 100%"
    max-height="500"
    stripe
  >
    <el-table-column label="标签名称" prop="name"></el-table-column>
    <el-table-column label="同步总人数" prop="value"></el-table-column>
    <el-table-column label="标签人数" prop="value">
      <template slot-scope="{ row }">
        <el-button type="text" @click.native.prevent="goCustomerList(row)">
          {{ row.value }}
        </el-button>
      </template>
    </el-table-column>
    <el-table-column label="同步时间">
      <!-- <template slot-scope="scope">
        <span>{{scope.row.value_update_at | formatDate("yyyy-mm-dd hh:mm:ss")}}</span>
      </template> -->
    </el-table-column>
    <el-table-column label="未匹配名单">
      <template slot-scope="scope">
        <el-button type="text" @click.native.prevent="handleDownLoad(scope.row.tag_id)">
          下载
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>
<script lang="ts">
import { Vue, Prop, Component } from "vue-property-decorator";
// import API from "@/api/customer";
@Component({
  components: {}
})
export default class extends Vue {
  @Prop() private outsideLabelList: any;

  /************data**************/

  private list: any = [];

  /************methods**************/

  // 跳转客户列表
  private goCustomerList(row: any) {
    this.$router.replace({
      name: "myapp_customer_customer-list",
      query: {
        tag_id: row.tag_id
      }
    });
  }

  // 下载报表
  private handleDownLoad(id: number) {
    console.log("点击下载", id);
  }

  /************生命周期**************/

  private created() {
    this.list = this.outsideLabelList[0].tags;
  }
}
</script>
