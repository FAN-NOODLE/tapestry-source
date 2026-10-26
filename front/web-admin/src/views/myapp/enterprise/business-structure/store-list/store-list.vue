<template>
  <div class="store-list">
    <layout-list>
      <template slot="search">
        <div class="title-mb">
          数据同步时间: {{ updateDate | formatDate("yyyy-mm-dd hh:mm:ss") }}
        </div>
        <search-item text="门店信息：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入门店名称或编码"
            clearable
            v-model="query.keyword"
          ></cy-input>
        </search-item>
        <search-item text="门店类型：">
          <cy-select
            v-model="query.type"
            placeholder="全部"
            :optionsList="type_list"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <search-item text="门店状态：">
          <cy-select
            v-model="query.close_status"
            placeholder="选择门店状态"
            :optionsList="close_status_list"
          >
          </cy-select>
        </search-item>
        <search-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="reset">重置</el-button>
          <el-button
            type="primary"
            @click="addStore"
            v-if="epNameDetail !== 'coach'"
            >新增门店</el-button
          >
          <el-button
            type="primary"
            @click="updateStore"
            v-if="epNameDetail === 'coach'"
            >获取企微门店信息</el-button
          >
          <el-button
            type="primary"
            @click="uploadStore"
            v-if="epNameDetail === 'coach'"
            >上传门店信息</el-button
          >
        </search-item>
      </template>
      <template slot="table">
        <el-table
          tooltip-effect="light"
          :data="store_list"
          stripe
          max-height="600"
          style="width:100%;"
          row-key="id"
        >
          <el-table-column
            label="部门ID"
            prop="wechat_department_id"
            key="wechat_department_id"
            v-if="epNameDetail === 'coach'"
            width="100"
          ></el-table-column>

          <el-table-column
            label="门店编号(部门名称)"
            prop="third_store_code"
            width="180"
            key="third_store_code"
          ></el-table-column>

          <el-table-column
            v-if="epNameDetail === 'coach'"
            prop="store_name"
            label="门店中文名称"
            key="store_name"
            min-width="200"
          ></el-table-column>
          <el-table-column
            v-if="epNameDetail !== 'coach'"
            prop="store_name"
            label="门店对外名称"
            key="store_name"
            min-width="200"
          ></el-table-column>
          <el-table-column
            prop="inventory_name"
            label="门店英文名称"
            key="inventory_name"
            min-width="200"
            v-if="epNameDetail === 'coach'"
          >
            <template slot-scope="{ row }">
              {{ row.inventory_name || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            v-if="epNameDetail === 'coach'" 
            prop="closed_at" 
            label="门店状态" 
            width="100">
            <template slot-scope="{row}">
              <el-tag v-if="row.closed_at > 0" size="mini" type="danger">闭店</el-tag>
              <el-tag v-else size="mini">正常</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="门店属性" width="150" key="store_type">
            <template slot-scope="{ row }">
              {{ type_list[row.store_type || 0].text }}
            </template>
          </el-table-column>
          <el-table-column
            prop="city"
            label="门店所属城市"
            key="city"
            min-width="200"
            v-if="epNameDetail === 'coach'"
          >
            <template slot-scope="{ row }">
              {{ row.city || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            label="门店邮箱"
            width="150"
            key="email"
            v-if="epNameDetail === 'coach'"
          >
            <template slot-scope="{ row }">
              {{ row.email || "--" }}
            </template>
          </el-table-column>
          <el-table-column label="门店地址" width="150" key="address">
            <template slot-scope="{ row }">
              {{ row.address || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            key="created_at"
            :label="epNameDetail !== 'coach' ? '创建时间' : '同步时间'"
            width="185"
          >
            <template slot-scope="{ row }">
              {{ row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="修改时间" width="185" key="updated_at">
            <template slot-scope="{ row }">
              {{ row.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column label="操作人" width="180" key="operate_name">
            <template slot-scope="{ row }">
              {{ row.operate_name || "--" }}
            </template>
          </el-table-column>
          <el-table-column
            label="已同步到企微"
            width="130"
            key="sync_status"
            v-if="epNameDetail !== 'coach'"
          >
            <template slot-scope="{ row }">
              {{ row.sync_status === 1 ? "是" : "否" }}
            </template>
          </el-table-column>
          <el-table-column
            label="上传同步时间"
            width="185"
            key="sync_at"
            v-if="epNameDetail !== 'coach'"
          >
            <template slot-scope="{ row }">
              {{ row.sync_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template slot-scope="{ row }">
              <el-button type="text" @click="addStore(row)">编辑</el-button>
              <el-button v-if="epNameDetail !== 'coach'" type="text" style="color: red" @click="delStore(row)"
                >删除</el-button
              >
              <el-dropdown v-if="epNameDetail === 'coach'" @command="handleCommand($event, row)">
                <el-button class="m-l-10" type="text">
                  <i class="el-icon-more"></i>
                </el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item v-if="row.closed_at > 0" :command="0">重启门店</el-dropdown-item>
                  <el-dropdown-item v-else :command="1">闭店</el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
            </template>
          </el-table-column>
          <template slot="empty" v-if="store_list.length == 0">
            <img
              class="empty-img m-t-40"
              src="@/assets/images/empty-box.png"
              alt=""
            />
            <div class="empty-box">
              <div class="desc">暂无门店</div>
            </div>
          </template>
        </el-table>
      </template>
      <template slot="footer" v-if="store_list.length !== 0">
        <el-pagination
          class="float-right m-t-20"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page.sync="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        >
        </el-pagination>
      </template>
    </layout-list>
    <edit-store
      v-if="store_model"
      :storeInfo="edit_store"
      @closeDialog="closeModel"
      @success="successModel"
    ></edit-store>
    <upload-store
      v-if="uploadStoreModel"
      @close="closeUploadModel"
      @success="sureUploadModel"
      @updateStoreList="getStoreList"
    ></upload-store>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import EditStore from "./../component/edit-store.vue";
import UploadStore from "./../component/upload-store.vue";
import Api from "@/api/enterprise/business";
@Component({
  components: {
    LayoutList,
    EditStore,
    UploadStore
  }
})
export default class extends Vue {
  // ---------------------------------------------------------------此页面所需变量---------------------------------------------------------------
  // 分页配置
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private query = {
    keyword: "",
    type: "",
    close_status: 0
  };
  private type_list = [
    {
      name: "全部",
      id: "",
      text: "--"
    },
    {
      name: "奥莱",
      text: "奥莱",
      id: 1
    },
    {
      name: "零售",
      text: "零售",
      id: 2
    }
  ];
  private store_list = [];
  private store_model = false;
  private edit_store = {};
  private epNameDetail = "";
  private uploadStoreModel = false;
  private updateDate = "";
  private close_status_list= [
    { value: 0, label: '全部' },
    { value: 1, label: '正常' },
    { value: 2, label: '闭店' },
  ]
  /************************************ 方法 ********************************************* */
  private getStoreList() {
    const params = {
      store_type: this.query.type || 0,
      code_or_name: this.query.keyword,
      close_status: this.query.close_status,
      page: this.meta.page,
      size: this.meta.size
    };
    Api.getStoreListApi(params).then(res => {
      const { total, page, size, list } = res.data;
      this.store_list = list;
      this.meta = {
        total,
        page,
        size
      };
    });
  }

  // 获取企业
  getEpConfig() {
    Api.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epNameDetail = name;
      // 缓存企业名称，给详情页面用
      localStorage.setItem("epNameDetail", name);
    });
  }

  // 搜索
  private search() {
    this.meta.page = 1;
    this.getStoreList();
  }
  private reset() {
    this.query = {
      keyword: "",
      type: "",
      close_status: 0
    };
    this.search();
  }
  private addStore(row = {}) {
    this.store_model = true;
    this.edit_store = row;
  }
  // 更新门店数据：1；获取更新时间：0
  private updateStoreData(is_sync) {
    const params = {
      is_sync
    };
    Api.getStoreSyncData(params).then(res => {
      if (res.data) {
        if (params.is_sync === 0) {
          this.updateDate = res.data.last_sync_time;
        } else {
          this.$message.success("更新成功");
        }
      }
    });
  }
  // 同步最新的门店数据
  private updateStore() {
    this.updateStoreData(1);
  }
  //上传门店信息
  private uploadStore() {
    this.uploadStoreModel = true;
  }
  private closeUploadModel() {
    this.uploadStoreModel = false;
  }
  private sureUploadModel() {
    this.closeUploadModel();
  }
  private closeModel() {
    this.store_model = false;
  }
  private successModel() {
    this.closeModel();
    this.reset();
  }
  private delStore(row) {
    this.$confirm("是否确认删除该门店", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        console.log("确认删除========");
        this.delSure(row);
      })
      .catch(() => {
        this.$message.info("已取消删除");
      });
  }
  private delSure(row) {
    const params = { department_id: row.id };
    Api.delStoreApi(params).then(() => {
      this.$message.success("删除成功");
      this.getStoreList();
    });
  }
  // 分页
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getStoreList();
  }
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getStoreList();
  }

  /**
   * 闭店/开店
   * @param key 
   * @param row 
   */
  handleCommand(key: number, row: any) {
    const message = [
      '确定开启该门店，开启后该门店将出现在业务组织架构中。',
      '确定关闭该门店，关闭后该门店将不在出现在业务组织架构中？'
    ][key]
    const title = [
      '开店提醒',
      '闭店提醒'
    ][key]
    this.$confirm(message, title, {
      confirmButtonText: "确认",
      confirmButtonClass: "button-danger",
      cancelButtonText: "取消",
      type: "warning"
    })
      .then(() => {
        Api.closeStoreApi({
          id: row.id,
          is_closed: key
        }).then(() => {
          this.$message.success("操作成功");
          this.getStoreList();
        })
      }).catch(() => {
        console.log('取消操作');
        
      })
  }

  // ---------------------------------------------------------------生命周期---------------------------------------------------------------

  // 实例挂载完
  mounted() {
    console.log("进入页面");
    this.getStoreList();
    this.getEpConfig();
    this.updateStoreData(0);
  }
}
</script>
<style lang="scss" scoped>
.title-mb {
  margin-bottom: 25px;
}
</style>
