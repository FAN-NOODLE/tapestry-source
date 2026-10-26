<template>
  <div class="display-data-list">
    <!-- 按钮 -->
    <el-row class="buttons">
      <el-button type="primary" @click="handleExport">导出列表数据</el-button>
      <el-button type="primary" @click="handleSetRemind">设置自动提醒</el-button>
      <el-button type="primary" :disabled="total===0" @click="handleRemind">一键提醒</el-button>
    </el-row>
    <!-- 筛选 -->
    <el-row class="search">
      <search-item text="标题：">
        <cy-input
          size="medium"
          placeholder="请输入标题名称"
          v-model="searchParams.display_name"
          clearable
          @change="handleNameChange"
        ></cy-input>
      </search-item>
      <search-item text="门店信息：">
        <cy-input
          size="medium"
          placeholder="请输入门店名称或编码"
          v-model="searchParams.store_name"
          clearable
        ></cy-input>
      </search-item>
      <search-item text="所属VM：">
        <cy-select
          class="search-select"
          size="medium"
          v-model="searchParams.owner_member_id"
          placeholder="请输入所属VM"
          :optionsList="vmList"
        >
          <cy-option
            v-for="item in vmList"
            :key="item.member_id"
            :label="item.member_name"
            :value="item.member_id"
          ></cy-option>
        </cy-select>
      </search-item>
      <search-item text="完成情况：">
        <cy-select
          class="search-select"
          size="medium"
          v-model="searchParams.finish_status"
          placeholder="请选择完成情况"
        >
          <cy-option
            v-for="item in statusOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          ></cy-option>
        </cy-select>
      </search-item>
      <search-item text="最近提交时间：">
        <el-date-picker
          class="cy-el-date-picker"
          v-model="submitTime"
          type="datetimerange"
          store_type-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          value-format="timestamp"
          :default-time="['00:00:00', '23:59:59']"
        >
        </el-date-picker>
      </search-item>
      <el-button size="medium" type="primary" @click="reload">查询</el-button>
      <el-button size="medium" plain @click="handleReset">重置</el-button>
    </el-row>
    <el-row class="panel" :gutter="20">
      <el-col v-for="(item, index) in statusOptions" :key="index" :span="4">
        <div class="panel-item">
          <div class="panel-item-label">{{item.label}}</div>
          <div class="panel-item-value">{{countInfo[item.key]}}</div>
        </div>
      </el-col>
    </el-row>
    <!-- 列表 -->
    <el-row class="table">
      <el-table :data="list" stripe>
        <el-table-column prop="display_name" label="标题" min-width="100"></el-table-column>
        <el-table-column prop="store_code" label="门店ID" min-width="100"></el-table-column>
        <el-table-column prop="store_name" label="门店名称" min-width="100"></el-table-column>
        <el-table-column prop="store_type" label="门店类型" min-width="100">
          <template slot-scope="{row}">
            <span>{{ ['--', '奥莱', '零售'][row.store_type] }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="store_city" label="所属城市" min-width="100">
          <template slot-scope="{row}">
            <span>{{ row.store_city || '--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="finish_status" label="完成情况" min-width="100">
          <template slot-scope="{row}">
            <el-tag size="small" :type="['', 'warning', 'danger', 'success', 'blue'][row.finish_status]">{{ finishStatus(row.finish_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="grade" label="评分" min-width="100">
          <template slot-scope="{row}">
            <span>{{ row.finish_status===3? row.grade:'--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="submit_name" label="最近提交导购" width="150">
          <template slot-scope="{row}">
            <span>{{ row.submit_name || '--' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="owner_member_name" label="所属VM" min-width="100">
          <template slot-scope="{row}">
            <span>{{ row.finish_status > 0 ? row.owner_member_name : '--' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="最近提交时间" width="200">
          <template slot-scope="{row}">
            <span v-if="row.submit_time">{{ row.submit_time | formatDate('yyyy/mm/dd hh:mm:ss') }}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column label="最近审批时间" width="200">
          <template slot-scope="{row}">
            <span v-if="row.review_time">{{ row.review_time | formatDate('yyyy/mm/dd hh:mm:ss') }}</span>
            <span v-else>--</span>
          </template>
        </el-table-column>
        <el-table-column prop="review_name" label="最近审批人" min-width="150">
          <template slot-scope="{row}">
            <span>{{ row.review_name || '--' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="160">
          <template slot-scope="{row}">
            <el-button type="text" size="small" @click="handleView(row)" :disabled="row.finish_status === 0">查看</el-button>
            <el-button v-show="row.finish_status !== 3" type="text" size="small" @click="handleAudit(row)" :disabled="row.finish_status !== 1 && row.finish_status !== 4">审核</el-button>
            <el-button v-show="row.finish_status === 3" type="text" size="small" @click="handleExportLong(row)">导出长图</el-button>
          </template>
        </el-table-column>
      </el-table> 
    </el-row>
    <!-- 分页 -->
    <el-row class="footer">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        class="page-right font-weight-400"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      ></el-pagination>
    </el-row>
    <!-- 设置自动提醒 弹框 -->
    <el-dialog 
      title="设置自动提醒" 
      :visible.sync="visible">
      <el-form :model="formInfo" :rules="formRules" label-width="110px">
        <el-form-item label="自动提醒设置：">
          <div>
            <span>执行时间结束前</span>
            <el-input 
              class="input-width" 
              v-model="formInfo.value"
              oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
              placeholder="请输入正整数"></el-input>
            <span>天，自动提醒未提交及未通过门店。</span>
          </div>
          <div class="tips">注：默认提醒时间为早上10:00</div>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="handleRemindCancel">取消</el-button>
        <el-button type="primary" @click="handleSetRemind">确定</el-button>
      </div>
    </el-dialog>
    <html-canvas ref="canvas"></html-canvas>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import Api from '@/api/wisdom/displayTask';
import HtmlCanvas from "./html-canvas.vue";
/** 查询条件 */
interface ISearchParams {
  display_task_id?: number;
  /** 任务标题 */
  display_name: string;
  /** 门店名称 */
  store_name: string;
  /** 所属VM */
  owner_member_id: number;
  /** 完成情况 */
  finish_status: number;
  /** 最近提交时间 开始 */
  submit_start_time: number;
  /** 最近提交时间 结束 */
  submit_end_time: number;

}

/** 下拉选项 */
interface ISelectItem {
  label: string;
  value: number;
  key?: string;
}
/** 分页 */
interface IMate {
  /** 页数 */
  page: number;
  /** 每页数量 */
  size: number;
  /** 总数 */
  total: number;
}
interface IItem {
  /** 主键 */
  allot_id: number;
  /** 任务标题 */
  display_name: string;
  /** 任务ID */
  display_task_id: number;
  /** 完成情况 */
  finish_status: number;
  /** 评分 */
  grade: number;
  /** 所属VM 导购id */
  owner_member_id: number;
  /** 所属VM 导购名称 */
  owner_member_name: string;
  /** 最近审批人 */
  review_name: string;
  /** 最近审批时间 */
  review_time: number;
  /** 门店城市 */
  store_city: string;
  /** 门店编码 */
  store_code: string;
  /** 门店ID */
  store_id: number;
  /** 门店名称 */
  store_name: string;
  /** 门店类型 */
  store_type: number;
  /** 最近提交导购id */
  submit_member_id: number;
  /** 最近提交导购名称 */
  submit_name: string;
  /** 最近提交时间 */
  submit_time: number;
}

interface IVmItem {
  member_name: string;
  member_id: number;
  user_id: string;
}
             
const defaultSearchParams: ISearchParams = {
  display_task_id: 0,
  display_name: '',
  store_name: '',
  owner_member_id: 0,
  finish_status: -1,
  submit_start_time: 0,
  submit_end_time: 0
}

@Component({
  components: {
    HtmlCanvas
  }
})
export default class DisplayDataList extends Vue {
  @Prop({ default: '' }) title: string;
  @Prop({ default: '' }) currentId: string;
  @Watch('title')
  titleChange(value: string) {
    if(value) {
      this.searchParams = { ...defaultSearchParams}
    }
    this.searchParams.display_name = value;
    this.searchParams.display_task_id = Number(this.currentId)
    this.reload()
  }
  /** 状态选项 */
  statusOptions: ISelectItem[] = [
    { label: '全部', value: -1, key: 'all_num' },
    { label: '未提交', value: 0, key: 'wait_num' },
    { label: '待审核', value: 1, key: 'check_num' },
    { label: '待复审', value: 4, key: 'double_check_num' },
    { label: '未通过', value: 2, key: 'fail_num' },
    { label: '已通过', value: 3, key: 'success_num' },
  ]
  /** 筛选数据 */ 
  searchParams: ISearchParams = { ...defaultSearchParams }
  /** 分页 */
  meta: IMate = {
    page: 1,
    size: 10,
    total: 0
  }
  /** 列表数据 */
  list: IItem[] = []

  /** 统计数据 */
  countInfo = {
    all_num: 0,
    wait_num: 0,
    check_num: 0,
    double_check_num: 0,
    fail_num: 0,
    success_num: 0,
  }

  visible = false
  formInfo = {
    value: ''
  }
  formRules = {
    value: [
      { required: true, message: '请输入正整数', trigger: 'blur' }
    ]
  }
  total = 0

  vmList: IVmItem[] = []

  get finishStatus() {
    return (value: number) => {
      return this.statusOptions.find(item => item.value === value).label
    }
  }

  /** 最近提交时间 读取 */
  get submitTime() {
    const { submit_start_time, submit_end_time } = this.searchParams
    return [submit_start_time*1000, submit_end_time*1000].filter(Boolean)
  }
  /** 最近提交时间 写入 */
  set submitTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.submit_start_time = parseInt(start_time/1000+'');
    this.searchParams.submit_end_time = parseInt(end_time/1000+'');
  }

  handleNameChange() {
    this.searchParams.display_task_id = null
  }

  /**
   * @description 查询
   */
  reload() {
    this.meta.page = 1
    this.fetchData()
    this.getAllotNum()
    this.getNoticeNum()
  }

  /**
   * @description 重置
   */
  handleReset() {
    this.searchParams = { ...defaultSearchParams }
    this.reload()
  }
  /**
   * @description 分页
   * @param size 
   */
  handleSizeChange(size: number) {
    this.meta.page = 1
    this.meta.size = size
    this.fetchData()
  }
  /**
   * @description 下一页
   * @param page 
   */
  handleCurrentChange(page: number) {
    this.meta.page = page
    this.fetchData()
  }

  /**
   * @description 查看
   * @param row 
   */
  handleView(row: IItem) {
    this.$router.push({
      name: "myapp_wisdom_display-task_audit-detail",
      query: { 
        id: String(row.display_task_id),
        store_id: String(row.store_id)
       }
    });
  }

  /**
   * @description 审核
   * @param row 
   */
  handleAudit(row: IItem) {
    this.$router.push({
      name: "myapp_wisdom_display-task_audit",
      query: { 
        id: String(row.display_task_id),
        store_id: String(row.store_id)
      }
    });
  }

  /**
   * @description 导出
   */
  handleExport() {
    console.log('导出');
    Api.exportAllotListApi(this.searchParams).then((res: any) => {
      if(res.code ===200){
        this.$confirm("文件生成中，稍后将自动下载，详情记录请前往【下载中心】查看。", "文件下载提示", {
          confirmButtonText: "去下载中心",
          cancelButtonText: "留在当前页",
          type: "success"
        })
          .then(() => {
            this.$router.push({
              path: "/myapp/setting/download-center/list"
            });
          })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消前往下载中心"
            });
          });
      }
    })
  }

  /**
   * @description 导出长图
   */
  handleExportLong(row: IItem) {
    console.log('导出长图');
    (this.$refs.canvas as any).save({ 
      id: String(row.display_task_id),
      store_id: String(row.store_id)
    })
  }

  getNoticeSetting() {
    Api.getNoticeSettingApi({}).then(res => {
      this.formInfo.value = res.data.value
    })
  }

  handleRemindCancel() {
    this.getNoticeSetting()
    this.visible = false
  }

  /**
   * @description 设置提醒
   */
  handleSetRemind() {
    if(this.visible === false) {
      this.visible = true
      return 
    }
    Api.setNoticeApi({
      value: Number(this.formInfo.value)
    }).then((res: any) => {
      if(res.code === 200) {
        this.$message.success('设置成功');
        this.visible = false
      }
    })
  }

  /**
   * @description 一键提醒
   */
  handleRemind() {
    const title = '一键提醒未提交门店';
    const message = '确认一键提醒？确定后将立即发送消息通知，给当前未提交及未通过的门店？'
    this.$confirm(message, title, {
      confirmButtonText: "确认",
      cancelButtonText: "取消"
    }).then(() => {
      
      Api.sendNoticeApi(this.searchParams).then((res: any) => {
        if(res.code === 200) {
          this.$message.success('一键提醒成功');
        }
      })
    })
  }

  /**
   * @description 获取列表数据
   */
  fetchData() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      ...this.searchParams
    }
    Api.getAllotListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  getAllotNum() {
    Api.getAllotNumApi(this.searchParams).then((res: any) => {
      if (res.code == 200) {
        this.countInfo = res.data
      } else {
        this.$message.error(res.message);
      }
    });
  }
  getNoticeNum() {
    Api.getNoticeNumApi(this.searchParams).then((res: any) => {
      if (res.code == 200) {
        this.total = res.data.num
      } else {
        this.$message.error(res.message);
      }
    });
  }
  getVmList() {
    Api.getVmListApi().then((res: any) => {
      if (res.code == 200) {
        this.vmList = [{
          member_id: 0,
          member_name: '全部',
          user_id: ''
        },...res.data]
      } else {
        this.$message.error(res.message);
      }
    });
  }

  mounted() {
    this.fetchData()
    this.getAllotNum()
    this.getNoticeSetting()
    this.getNoticeNum()
    this.getVmList()
  }
  activated() {
    this.fetchData()
    this.getAllotNum()
    this.getNoticeSetting()
    this.getNoticeNum()
    this.getVmList()
  }
}
</script>
<style lang="scss" scoped>
/deep/.cy-select.search-select .el-select .el-input--medium, 
/deep/.cy-select.search-select .el-select .el-input--medium .el-input__inner, 
/deep/.cy-select.search-select .el-select .el-input--small, 
/deep/.cy-select.search-select .el-select .el-input--small .el-input__inner {
  height: 40px !important;
  line-height: 40px !important;
}
/deep/.el-tag.el-tag--blue {
  color: #409EFF !important;
  background-color: rgb(236, 245, 255) !important;
  border-color: rgb(160, 207, 255) !important;
}
.display-data-list {
  .buttons {
    margin-bottom: 20px;
  }
  .search {
    padding: 20px 20px 0px 20px;
    margin-bottom: 20px;
    background-color: #f7f7f7;
    border-radius: 4px;
  }

  .panel {
    padding: 20px 20px 0px 20px;
    margin-bottom: 20px;
    background-color: #f7f7f7;
    border-radius: 4px;
    &-item {
      margin-bottom: 20px;
      height: 100px;
      background-color: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      border-radius: 4px;

      &-label {
        color: #666;
      }

      &-value {
        margin-top: 20px;
        font-size: 30px;
        font-weight: bold;
      }
    }
  }

  .footer {
    text-align: right;
    margin-top: 20px;
  }
}
.input-width {
  width: 200px;
  margin-left: 5px;
  margin-right: 5px;
}
.tips {
  color: #999;
}
</style>