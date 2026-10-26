<template>
  <div class="detail">
    <layout-list>
      <template slot="table">
        <el-form label-width="90px">
          <el-form-item label="标签组名称">
            <div>{{ details.group_name }}</div>
          </el-form-item>
          <el-form-item label="标签名称">
            <div>{{ details.tag_name }}</div>
          </el-form-item>
          <el-form-item label="标签人数">
            <div>{{ meta.total }}</div>
          </el-form-item>
          <el-form-item label="客户名单">
            <el-button type="text" @click="exportList">导出名单</el-button>
          </el-form-item>
          <el-form-item>
            <el-table
              :data="list"
              ref="tableList"
              style="width:80%"
              :cell-style="{ height: '20px' }"
            >
              <el-table-column
                label="客户昵称"
                min-width="100"
                prop="name"
              ></el-table-column>
              <el-table-column
                label="客户手机号"
                width="180"
                prop="mobile"
              ></el-table-column>
              <el-table-column
                label="union ID"
                min-width="180"
                prop="union_id"
              ></el-table-column>
            </el-table>
          </el-form-item>
        </el-form>
      </template>
      <template slot="footer">
        <el-pagination
          v-if="list.length > 0"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        >
        </el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="tsx" type="text/tsx">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import Api from "@/api/customer";
@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  private list = [1, 2, 3];
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private details = {
    tag_id: 0,
    tag_name: "",
    group_name: "",
    tag_type: 0
  };
  //*****************************方法************************************ */
  private mounted() {
    Object.keys(this.details).map(key => {
      this.details[key] = this.$route.query[key] || "";
    });
    console.log(this.details);
    this.getList();
  }
  // 时间格式转换
  private formateDate(time: any) {
    const year = time.getFullYear();
    let month = time.getMonth() + 1;
    let day = time.getDate();

    month = month < 10 ? "0" + month : month;
    day = day < 10 ? "0" + day : day;

    return year + "-" + month + "-" + day;
  }
  // 导出客户列表
  private exportList() {
    console.log("导出页面=====");
    Api.exportTagDetailApi({ tag_id: Number(this.details.tag_id),tag_type:Number(this.details.tag_type) }).then(
      (res: any) => {
        // if (res.code !== 200) {
        //   this.$message.error("导出文件失败");
        //   return;
        // }
        const blob = new Blob([res], {
          type: "application/x-xls;charset=utf-8"
        });
        console.log(blob);
        let fileName = "";
        const nowTime = new Date();
        fileName = "客户标签_" + this.formateDate(nowTime) + ".xlsx";

        //下载的文件名称及其后缀，后缀要和后台保持的一致

        if ("download" in document.createElement("a")) {
          // 非IE下载
          const elink = document.createElement("a");
          elink.download = fileName;
          elink.style.display = "none";
          elink.href = URL.createObjectURL(blob);
          document.body.appendChild(elink);
          elink.click();
          URL.revokeObjectURL(elink.href); // 释放URL 对象
          document.body.removeChild(elink);
        } else {
          // IE10+下载
          (navigator as any).msSaveBlob(blob, fileName);
        }
      }
    );
  }
  // 获取标签用户列表
  private getList() {
    const params = {
      tag_id: Number(this.details.tag_id),
      tag_type: Number(this.details.tag_type),
      page: this.meta.page,
      size: this.meta.size
    };
    Api.getTagCustomersApi(params).then(res => {
      const { total, user_list } = res.data;
      this.list = user_list || [];
      this.meta.total = total;
    });
  }
  //内部绑定函数（获取每页数据量）
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getList();
  }
  //内部绑定函数（获取当前第几页）
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getList();
  }
}
</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
/deep/.footer-content {
  padding-right: 20%;
}
</style>
