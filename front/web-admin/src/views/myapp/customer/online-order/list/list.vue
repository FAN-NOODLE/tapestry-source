<template>
  <layout-list>
    <template slot="buttons">
      <el-button size="small" plain @click="handler(1)">导出订单信息</el-button>
      <el-button size="small" type="primary" @click="handler(2)">导入订单物流号</el-button>
      <el-button size="small" type="primary" @click="handler(3)">批量上传转货单</el-button>
    </template>
    <template slot="search">
      <el-form inline size="medium">
        <el-form-item label="门店">
          <el-input
            v-model="inputType.value"
            :placeholder="placeholder"
            clearable
            @keyup.enter.native="handleSearch"
          >
            <el-select
              slot="prepend"
              v-model="inputType.key"
              class="input-with-select"
            >
              <el-option
                v-for="item in inputType.options"
                :value="item.value"
                :label="item.label"
                :key="item.value"
              ></el-option>
            </el-select>
          </el-input>
        </el-form-item>
        <el-form-item label="销售单号">
          <el-input
            size="medium"
            placeholder="请输入销售单号"
            v-model="searchParams.order_sn"
            clearable
          ></el-input>
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select
            v-model="searchParams.status"
            placeholder="请选择订单状态"
          >
            <el-option label="全部" :value="0" ></el-option>
            <el-option label="待填写" :value="1" ></el-option>
            <el-option label="待发货" :value="2" ></el-option>
            <el-option label="已发货" :value="3" ></el-option>
            <el-option label="已拒单" :value="4" ></el-option>
            <el-option label="已作废" :value="5" ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="转货单状态">
          <el-select
            v-model="searchParams.transfer_status"
            placeholder="请选择转货单状态"
          >
            <el-option label="全部" :value="0" ></el-option>
            <el-option label="待上传" :value="1" ></el-option>
            <el-option label="已上传" :value="2" ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            class="cy-el-date-picker"
            v-model="createTime"
            type="datetimerange"
            store_type-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="timestamp"
            :default-time="['00:00:00', '23:59:59']"
          >
          </el-date-picker>
        </el-form-item>
        <el-button class="m-b-20" size="small" type="primary" @click="handleSearch">查询</el-button>
        <el-button class="m-b-20" size="small" plain @click="handleReset">重置</el-button>
      </el-form>
    </template>
    <template slot="table">
      <el-table :data="list" stripe>
        <el-table-column prop="store_code" label="门店代码" width="100"></el-table-column>
        <el-table-column prop="store_name" label="门店名称" min-width="100"></el-table-column>
        <el-table-column prop="order_sn" label="销售单号" min-width="100"></el-table-column>
        <el-table-column prop="total_quantity" label="商品数量" min-width="100"></el-table-column>
        <el-table-column prop="status" label="订单状态" min-width="100">
          <template slot-scope="scope">
            <el-tag :type="['', 'primary', 'warning', 'success', 'danger', 'info'][scope.row.status]">{{['', '待填写', '待发货', '已发货', '已拒单', '已作废'][scope.row.status]}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="transfer_status" label="转货单状态" min-width="120">
          <template slot-scope="scope">
            <el-tag :type="['', 'info','success'][scope.row.transfer_status]">{{['', '待上传','已上传'][scope.row.transfer_status]}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="创建时间" width="180">
          <template slot-scope="scope">
            {{ scope.row.created_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180">
          <template slot-scope="scope">
            {{ scope.row.updated_at | formatDate('yyyy/mm/dd hh:mm:ss') }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template slot-scope="scope">
            <el-button type="text" @click="handler(4, scope.row)">查看</el-button>
            <el-button type="text" :disabled="scope.row.status !== 2" @click="handler(5, scope.row)">发货</el-button>
            <el-button type="text" :disabled="scope.row.status !== 2" @click="handler(6, scope.row)">拒单</el-button>
          </template>
        </el-table-column>
      </el-table>
      <!-- 批量导入 -->
      <batch-import
        v-model="showBatchImport"
        @download="handleDownload"
        @confirm="handleBatchImportConfirm"
      />
      <!-- 批量上传 -->
      <batch-update v-model="showBatchUpdate" :params="{status: 2}" @view-order="handleViewOrder" />
      <!-- 订单列表 -->
      <store-order-list v-model="showOrderList" :store-info="currentStore"></store-order-list>
    </template>
    <template slot="footer">
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
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import Api from '@/api/customer/online-order';
import BatchImport from '../components/batch-import.vue';
import BatchUpdate from '../components/batch-update.vue';
import StoreOrderList from '../components/store-order-list.vue';
import { download, downloadFile } from "@/utils/download";
import dayjs from "dayjs";
import { cosUploader } from "@/utils/upload";
/**
 * 1.导出订单信息;  
 * 2.导入订单物流号;  
 * 3.批量上传转货单;  
 * 4.查看;  
 * 5.发货
 * 6.拒单
 */
type IOperationType = 1 | 2 | 3 | 4 | 5 | 6;
/**
 * 列表项
 */
interface IItem {
  /** 主键ID */
  id: number
  /** 导购ID */
  member_id: number;
  /** 导购名称 */
  member_name: string;
  /** 销售单号 */
  order_sn: number;
  /** 收货地址 */
  receiver_address: number;
  /** 收获人名称 */
  receiver_name: number;
  /** 收获手机号 */
  receiver_phone: number;
  /** 发货时间 */
  shipped_at: number;
  /** 订单状态，0待发货，1已发货，2已据单 */
  status: number;
  /** 门店编码 */
  store_code: number;
  /** 门店名称 */
  store_name: number;
  /** 运单号 */
  tracking_number: string;
  /** 转货单信息 */
  transfer_file: {
    name: string;
    url: string;
  };
  /** 转货单状态 1待上传，2已上传  */
  transfer_status: number;
  /** 商品数量 */
  total_quantity: number
}

interface ISearchParams {
  /** 门店code */
  store_code: string;
  /** 门店名称 */
  store_name: string;
  /** 销售单号 */
  order_sn: string;
  /** 订单状态 0全部，'待填写', '待发货', '已发货', '已拒单', '已作废' */
  status: 0|1|2|3|4|5;
  /** 转货单状态 0全部，1待上传，2已上传*/
  transfer_status: 0|1|2;
  /** 创建时间 开始 */
  start_time: number;
  /** 创建时间 结束 */
  end_time: number;
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

/** 默认筛选条件 */
const defaultSearchParams: ISearchParams = {
  store_code: "",
  store_name: "",
  order_sn: "",
  status: 0,
  transfer_status: 0,
  start_time: 0,
  end_time: 0
}

@Component({
  name: 'OnlineOrderList',
  components: {
    BatchImport,
    BatchUpdate,
    StoreOrderList
  }
})
export default class OnlineOrderList extends Vue {
  /** 搜索参数 */
  searchParams: ISearchParams = {...defaultSearchParams}
  /** 分页 */
  meta: IMate = {
    page: 1,
    size: 10,
    total: 0
  }
  /** 列表数据 */
  list: IItem[] = []

  inputType = {
    key: 'store_code',
    value: '',
    options: [
      { value: 'store_code', label: '门店代码'},
      { value: 'store_name', label: '门店名称'},
    ]
  }
  showBatchImport = false
  showBatchUpdate = false
  showOrderList = false
  currentStore:any = {}
  get placeholder() {
    const item = this.inputType.options.find(item => item.value === this.inputType.key)
    return `请输入${item.label}`
  }

  /** 创建时间 读取 */
  get createTime() {
    const { start_time, end_time } = this.searchParams
    return [start_time*1000, end_time*1000].filter(Boolean)
  }

  /** 创建时间 写入 */
  set createTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.start_time = parseInt(start_time/1000+'');
    this.searchParams.end_time = parseInt(end_time/1000+'');
  }
  get params() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      ...this.searchParams,
    }
    params[this.inputType.key] = this.inputType.value
    return params
  }
  mounted() {
    this.fetchData();
  }
  activated() {
    this.fetchData();
  }
  /**
   * 获取列表
   */
  fetchData() {
    Api.getOnlineOrderListApi(this.params).then((res: any) => {
      if (res.code == 200) {
        this.list = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /**
   * 操作
   * @param key 
   * @param row 
   */
  handler(key: IOperationType, row?: IItem) {
    switch (key) {
      case 1:
        this.handleExport()
        break;
      case 2:
        this.showBatchImport = true
        break;
      case 3:
        this.showBatchUpdate = true
        break;
      case 4:
        this.$router.push({
          name: "myapp_customer_online-order_detail",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 5:
        this.$router.push({
          name: "myapp_customer_online-order_deliver",
          query: {
            id: String(row.id)
          }
        });
        break;
      case 6:
        this.handleReject(row)
        break;
    
      default:
        break;
    }
  }

  /**
   * 拒单
   */
  handleReject(row: any) {
    this.$confirm(`销售单号:${row.order_sn} 商品数量:${row.total_quantity}`, '请确认是否全部拒单', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      Api.refundOrderApi({
        order_id: row.id,
      }).then((res: any) => {
        if(res.code == 200) {
          this.$message.success('操作成功')
          this.refresh()
        }
      })
    })
  }

  /**
   * @description 查询
   */
  handleSearch() {
    this.meta.page = 1
    this.fetchData()
  }
  
  /**
   * @description 重置
   */
  handleReset() {
    this.searchParams = { ...defaultSearchParams }
    this.inputType.value = ''
    this.inputType.key = 'store_code'
    this.handleSearch()
  }
  /**
   * 每页数量变化
   * @param val 
   */
  handleSizeChange(val: number) {
    this.meta.size = val;
    this.meta.page = 1;
    this.fetchData();
  }

  /**
   * 当前页
   * @param val 
   */
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.fetchData();
  }

  /**
   * 删除后刷新数据
   */
  refresh() {
    const { total, page, size } = this.meta
    if((total - 1) % size === 0 && (total - 1) / size < page) {
      this.meta.page = Math.max(page-1, 1)
    }
    this.fetchData()
  }

  

  /**
   * @description 导出数据
   */
  async handleExport() {
    const res:any = await Api.exportOrderListApi(this.params)
    if(res.code == 200) {
      cosUploader.getObjectUrl(res.data.url).then(res => {
        downloadFile(res, {
          filename: `网仓订单表_${dayjs().format('YYYYMMDD')}.xlsx`
        })
      })
    } else {
      this.$message.error('导出失败')
    }
  }

  async handleDownload() {
    const res:any = await Api.exportOrderListApi({
      page: 1,
      size: 10,
      status: 2,
      export_type: 'import'
    })
    if(res.code == 200) {
      cosUploader.getObjectUrl(res.data.url).then(res => {
        downloadFile(res, {
          filename: `网仓订单表_${dayjs().format('YYYYMMDD')}.xlsx`
        })
      })
    } else {
      this.$message.error('下载失败')
    }
  }


  /**
   * 模版导入确认
   */
  handleBatchImportConfirm() {
    this.showBatchImport = false
  }

  handleViewOrder(row: any) {
    this.currentStore = row
    setTimeout(() => {
      this.showOrderList = true
    }, 10);
  }

}
</script>

<style lang="scss" scoped>
::v-deep.layout-list .search-content {
  padding-bottom: 0 !important;
}

.el-form-item {
  margin-bottom: 15px;
}

.el-input,
.el-select,
.el-date-editor {
  width: 280px;
}

.input-with-select {
  width: 110px;
}
  
</style>