<template>
  <div>
    <el-button
      class="create-btn"
      type="primary"
      size="small"
      @click="createSignFn"
      >创建短信签名</el-button
    >
    <p class="time-show">
      数据更新时间：{{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}
    </p>
    <div class="search-nav">
      <search-item text="签名：">
        <select-input
          ref="select"
          v-model="curTemplateName"
          :list="templateType"
          @selectChange="selectChangedFn"
          :placeholder="'请输入' + curTemplatePla"
        ></select-input>
      </search-item>
      <search-item text="状态：">
        <cy-select width-size="small" size="small" v-model="curTemplateStatus">
          <cy-option
            v-for="(item, index) in templateStatus"
            :key="index"
            :label="item"
            :value="index"
          ></cy-option>
        </cy-select>
      </search-item>
      <el-button
        type="primary"
        size="small"
        class="m-b-20"
        @click="getSmsSignFn"
        >查询</el-button
      >
      <el-button size="small" class="m-b-20" @click="resetFn">重置</el-button>
    </div>
    <!--table-->
    <el-table :data="tableData" stripe>
      <el-table-column prop="id" label="签名ID"> </el-table-column>
      <el-table-column prop="sign_name" label="签名"> </el-table-column>
      <!-- <el-table-column prop="province" label="状态">
        <template slot-scope="scope">
          <el-tag
            :type="
              scope.row.third_sign_status == 1
                ? 'warning'
                : scope.row.third_sign_status == 2
                ? 'danger'
                : 'success'
            "
            >{{ auditStatus[scope.row.third_sign_status - 1] }}</el-tag
          >
        </template>
      </el-table-column> -->
      <el-table-column label="操作" width="100">
        <template slot-scope="scope">
          <el-button
            style="color:#F56C6C"
            @click="deleteFn(scope.row.id)"
            type="text"
            size="small"
            >删除</el-button
          >
        </template>
      </el-table-column>
      <div slot="empty" class="table-empty">
        <div v-if="searchLoading">
          <img src="@/assets/images/empty-search.png" alt="" />
          <div>
            未搜索到短信签名
          </div>
        </div>
        <div v-else>
          <img src="@/assets/images/empty-box.png" alt="" />
          <p>暂无短信签名</p>
        </div>
      </div>
    </el-table>
    <div class="paging" v-if="tableData.length">
      <el-pagination
        @size-change="sizeChangeFn"
        @current-change="currentChangeFn"
        :current-page="paging.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="paging.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="paging.total"
      >
      </el-pagination>
    </div>
    <!--创建短信签名-->
    <set-sign v-if="isSetSign" @closeSign="createSignFn"></set-sign>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from "vue-property-decorator";
import API from "@/api/customer";
import setSign from "./set-sign.vue";
import SelectInput from "@/components/business/select-input/index.vue";
@Component({
  components: { SelectInput, setSign }
})
export default class extends Vue {
  @Prop() isTabApi?: boolean;
  /************data**************/
  private updateTime = 0; //更新时间
  private curTemplateName = ""; //选择的模板的name
  private curTemplatePla = "签名"; //选择模板的key
  private curTemplateType = 1; //选择模板的key
  //模板类型
  private templateType = [
    { key: 1, name: "签名" },
    { key: 2, name: "签名ID" }
  ];
  //模板状态
  private templateStatus = ["全部", "审核中", "审核失败", "审核成功"];
  private auditStatus = ["审核中", "审核失败", "审核成功"];
  private curTemplateStatus = 0; //当前模板状态
  private tableData = [];
  private isSetSign = false; //是否创建短信签名
  //分页信息
  private paging = {
    page: 1,
    size: 10,
    total: 0
  };
  private searchLoading = false; //控制搜索缺省图
  /************watch**************/
  @Watch("isTabApi")
  onChangeIsTabApi() {
    if (this.isTabApi) {
      this.getSmsSignFn();
      this.paging.page = 1;
    }
  }
  /************methods**************/
  //创建短信模板
  private createSignFn(type?: string) {
    this.isSetSign = !this.isSetSign;
    if (type) {
      this.getSmsSignFn();
    }
  }
  //模板
  private selectChangedFn(val: any) {
    this.curTemplateType = val;
    this.templateType.forEach(item => {
      if (item.key == val) {
        this.curTemplatePla = item.name;
      }
    });
  }
  //删除
  private deleteFn(id: number) {
    this.$confirm("确认删除该签名吗", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        const params = {
          id: id
        };
        //删除子标签接口
        API.smsSignDeleteApi(params).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              message: "删除成功",
              type: "success",
              duration: 2000
            });
            this.getSmsSignFn();
          } else {
            this.$message.error(res.message);
          }
        });
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  //列表
  private getSmsSignFn() {
    const params: any = {
      search_type: this.curTemplateType,
      search_text: this.curTemplateName,
      status: this.curTemplateStatus,
      page: this.paging.page,
      size: this.paging.size
    };
    API.smsSignListApi(params).then((res: any) => {
      if (res.code === 200) {
        //判断缺省
        if (
          this.curTemplateStatus ||
          this.curTemplateName ||
          this.curTemplateType
        ) {
          this.searchLoading = true;
        } else {
          this.searchLoading = false;
        }
        const data = res.data;
        this.tableData = data.list;
        this.updateTime = data.updateTime;
        this.paging.total = data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //重置
  private resetFn() {
    (this.$refs["select"] as any).reset();
    this.curTemplatePla = "签名";
    this.curTemplateType = 1;
    this.curTemplateName = "";
    this.curTemplateStatus = 0;
    this.getSmsSignFn();
  }
  //改变size
  private sizeChangeFn(val: number) {
    this.paging.page = 1;
    this.paging.size = val;
    this.getSmsSignFn();
  }
  //改变当前页
  private currentChangeFn(val: number) {
    this.paging.page = val;
    this.getSmsSignFn();
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    //列表
    this.getSmsSignFn();
  }
}
</script>
<style lang="scss" scoped>
li {
  list-style: none;
}
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
a {
  text-decoration: none;
}
//下拉输入框行显示
.select-input {
  display: inline-block;
}
.create-btn {
  margin-bottom: 20px;
}
.time-show {
  color: $color-text-normal;
  font-size: 14px;
  line-height: 20px;
}
.search-nav {
  background: $color-background;
  padding-top: 20px;
  padding-left: 20px;
  margin: 20px 0 30px;
}
/deep/.el-button--mini,
.el-button--small {
  font-size: 14px;
}
.paging {
  text-align: right;
  padding-top: 20px;
}
.table-empty {
  height: 340px;
  text-align: center;
  margin: 80px 0 0 0;
  img {
    margin-top: 30px;
  }
  p {
    margin-top: -40px;
  }
}
</style>
