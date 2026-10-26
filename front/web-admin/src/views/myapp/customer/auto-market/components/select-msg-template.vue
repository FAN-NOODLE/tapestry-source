<template>
  <el-dialog
    title="选择短信模板"
    :visible.sync="dialogVisible"
    class="msg-table-hidden-radio-label"
    :close-on-click-modal="false"
    @close="closeDialogFn"
  >
    <span style="display:inline-block">
      <cy-input
        placeholder="请输入模板名称"
        v-model="searchPageText"
        size="small"
        width-size="small"
      >
      </cy-input>
    </span>

    <el-button
      type="primary"
      size="small"
      class="m-b-20 marginL20"
      @click="onSearchFn"
    >查询</el-button>
    <el-button
      size="small"
      class="m-b-20"
      @click="resetFn"
    >重置</el-button>
    <el-table
      :data="list"
      ref="table"
      stripe
    >
      <el-table-column
        label=""
        width="55"
      >
        <template slot-scope="{ row }">
          <el-radio
            :label="row.id"
            v-model="templateRadio"
            @change.native="getMsgRow(row)"
          >&nbsp;</el-radio>
        </template>
      </el-table-column>
      <el-table-column
        prop="template_name"
        label="模板名称"
        min-width="310"
        :show-overflow-tooltip="true"
      ></el-table-column>
      <el-table-column
        prop="message_str"
        label="短信内容"
        min-width="310"
        :formatter="formatTempFn"
        :show-overflow-tooltip="true"
      >
        <template slot-scope="{ row }">
          <el-popover
            placement="top-start"
            width="600"
            trigger="hover"
            :disabled="row.message_str && row.message_str.length < 20"
          >
            <div>{{ row.message_str }}</div>

            <span
              slot="reference"
              v-if="row.message_str && row.message_str.length < 20"
            >{{ row.message_str }}</span>
            <span
              slot="reference"
              v-if="row.message_str && row.message_str.length > 20"
            >{{ row.message_str.substr(0, 20) + "..." }}</span>
          </el-popover>
        </template>
      </el-table-column>
      <div
        slot="empty"
        class="table-empty-label"
      >
        <img
          class="empty-img"
          :src="searchPageText
              ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
              : require('@/assets/images/empty-box.png')
            "
          alt=""
        />
        <p>{{ searchPageText ? "没有找到该短信模板" : "暂无短信模板" }}</p>
      </div>
    </el-table>
    <div
      class="clearfix"
      style="margin-top:20px;"
    >
      <el-pagination
        class="pull-right"
        @current-change="handleCurrentChangeFn"
        :current-page="currentPage"
        :page-size="5"
        layout="total, prev, pager, next, jumper"
        :total="total"
      >
      </el-pagination>
    </div>

    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button
        @click="closeDialogFn"
        size="small"
      >取 消</el-button>
      <el-button
        type="primary"
        class="marginL20"
        @click="submitFn"
        size="small"
      >确 认
      </el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/customer";
@Component({})
export default class extends Vue {
  @Prop() msgName!: any;
  /************data**************/
  private templateRadio = "";
  private total = 0;
  private list = [];
  private dialogVisible = true;
  private searchPageText = "";
  private currentPage = 1;
  private selectedRow = null;
  private choosed = null;
  created() {
    this.choosed = this.msgName;
    this.queryTableFn();
  }
  queryTableFn() {
    const params = {
      page: this.currentPage,
      size: 5,
      search_type: 2,
      search_text: this.searchPageText || "",
      status: 5 //4审核成功，5已赋值
    };
    API.smsTemplateListApi(params).then((res: any) => {
      if (res.code === 200) {
        const data = res.data.list;
        this.total = res.data.total || 0;
        this.list = data;
        if (this.choosed) {
          this.list.forEach(v => {
            if (v["template_name"] == this.choosed) {
              this.templateRadio = v["id"];
              this.selectedRow = v;
            }
          });
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  formatTempFn(row: any, column: any, cellValue: any) {
    const data: any = cellValue.replace(/(<[^>]+>)/g, function(substr: string) {
      const reg: any = new RegExp("(%[^%]+%)");
      return reg.exec(substr) ? reg.exec(substr)[0] : "";
    });
    return data;
  }
  async submitFn() {
    this.$emit("addMsgDialog", this.selectedRow);
  }
  closeDialogFn() {
    this.$emit("addMsgDialog");
  }
  onSearchFn() {
    this.currentPage = 1;
    this.queryTableFn();
  }
  resetFn() {
    this.searchPageText = "";
    this.currentPage = 1;
    this.queryTableFn();
  }
  getMsgRow(row: any) {
    this.selectedRow = row;
    this.choosed = row.template_name;
  }
  handleCurrentChangeFn(val: number) {
    this.currentPage = val;
    this.queryTableFn();
    //获取当前页数据后根据selectList里当前page存的值进行渲染
  }
}
</script>
<style lang="scss">
.msg-table-hidden-radio-label {
  .el-radio__label {
    display: none;
  }
}
</style>
<style lang="scss" scoped>
.empty-img {
  margin-top: 50px;
  width: 124px;
}

.marginL20 {
  margin-left: 20px;
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
}</style>
