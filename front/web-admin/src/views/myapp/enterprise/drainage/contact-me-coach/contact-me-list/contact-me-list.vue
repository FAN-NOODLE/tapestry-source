<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button
          type="primary"
          @click="goEdit()"
          size="small"
        >新增活动码</el-button>
      </template>
      <template slot="search">
        <search-item text="活动名称：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入活动名称"
            clearable
            class="custom-name"
            v-model="searchParams.name"
          ></cy-input>
        </search-item>
        <search-item text="渠道：">
          <cy-select
            v-model="searchParams.channel_id"
            placeholder="请选择渠道"
            search-input
            :optionsList="channelList"
          >
          </cy-select>
        </search-item>
        <search-item text="二维码类型：">
          <cy-select
            v-model="searchParams.type"
            placeholder="请选择二维码类型"
            :optionsList="codeTypeList"
          >
          </cy-select>
        </search-item>
        <search-item text="创建时间：">
          <el-date-picker
            size="medium"
            v-model="createTime"
            type="datetimerange"
            store_type-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="timestamp"
            :default-time="['00:00:00', '23:59:59']"
          >
          </el-date-picker>
        </search-item>
        <search-item>
          <el-button
            type="primary"
            @click="search()"
            size="small"
          >查询</el-button>
          <el-button
            @click="reset()"
            size="small"
          >重置</el-button>
        </search-item>
      </template>
      <template slot="table">
        <el-table
          :data="contactList"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            label="活动名称"
            min-width="150"
            align="center"
            prop="name"
            show-overflow-tooltip
          >
          </el-table-column>
          <el-table-column
            label="渠道"
            min-width="150"
            align="center"
            prop="channel_name"
          >
          </el-table-column>
          <el-table-column
            label="渠道标签"
            min-width="300"
            align="center"
            prop="tags"
          >
            <template slot-scope="scope">
              <!-- 显示前三个，超出的省略，鼠标上去tooltip显示所有 -->
              <div v-if="scope.row.tags.length == 0">--</div>
              <div v-else-if="scope.row.tags.length <=3" class="tags">
                <el-tag v-for="(item, index) in scope.row.tags" :key="index" type="blue">{{ item }}</el-tag>
              </div>
              <el-popover v-else placement="top" width="200px" trigger="hover">
                <div class="tags" slot="reference">
                  <el-tag v-for="(item, index) in scope.row.tags.slice(0,3)" :key="index" type="blue">{{ item }}</el-tag>
                  <span>...</span>
                </div>
                <div class="tags">
                  <el-tag v-for="(item, index) in scope.row.tags" :key="index" type="blue">{{ item }}</el-tag>
                </div>
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            label="二维码类型"
            min-width="150"
            align="center"
            prop="type"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.type === 1">员工码+门店码</span>
              <span v-else>多人共享码</span>
            </template>
          </el-table-column>
          <el-table-column
            label="二维码数量"
            min-width="150"
            align="center"
            prop="code_number"
          >
          </el-table-column>
          <el-table-column
            prop="add_customer_number"
            label="新增好友数"
            min-width="150"
            align="center"
          ></el-table-column>
          <el-table-column
            prop="is_config_welcome"
            label="是否配置欢迎语"
            min-width="150"
            align="center"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.has_welcome === 1" style="color: green;">已配置</span>
              <span v-else>未配置</span>
            </template>
          </el-table-column>
          <el-table-column
            label="创建时间"
            min-width="200"
            prop="created_at"
            align="center"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.created_at > 0">{{
                scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            width="280"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                type="text"
                @click="downLoad(scope.row)"
              >
                下载
              </el-button>
              <el-button
                type="text"
                @click="goDetail(scope.row.id)"
              >
                查看二维码
              </el-button>
              <el-button
                type="text"
                @click="goEdit(scope.row.id)"
              >
                编辑
              </el-button>
              <el-button
                type="text"
                @click="deleteContact(scope.row.id)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              class="empty-img"
              :src="searchLoading
                  ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                  : require('@/assets/images/empty-box.png')
                "
              alt=""
            />
            <div>{{ searchLoading ? "没有找到相关联系我" : "暂无联系我" }}</div>
          </div>
        </el-table>
      </template>
      <template
        slot="footer"
        v-if="meta.total > 0"
      >
        <el-pagination
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
    <!-- 下载二维码弹框-->
    <download-code
      v-if="showDownLoad"
      :currentDownLoadItem="currentDownLoadItem"
      @closeDialog="showDownLoad = false"
      :zip="true"
    ></download-code>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/drainage-coach";
import downloadCode from "../../components/download-code.vue";
import { downloadFile } from "@/utils/download";
const defaultSearchParams = {
  name: "", // 二维码
  channel_id: 0, // 渠道
  create_time: {
    start_time: 0, // 创建时间-开始
    end_time: 0 // 创建时间-结束
  }, // 创建时间-开始
  type: 0 // 二维码类型
};

@Component({
  components: {
    downloadCode
  }
})
export default class extends Vue {
  /*************************************** data ****************************************/
  private searchParams = JSON.parse(JSON.stringify(defaultSearchParams));
  private contactList: Array<any> = []; // 联系我列表数组
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private searchLoading = false; // 是否正在搜索
  private showDownLoad = false; // 控制下载二维码的弹框显隐
  private currentDownLoadItem: any; // 当前下载二维码的url
  private channelList = [];
  private codeTypeList = [
    {
      value: 0,
      label: "全部"
    },
    {
      value: 1,
      label: "员工码+门店码"
    },
    {
      value: 2,
      label: "多人共享码"
    }
  ];

  /** 创建时间 读取 */
  get createTime() {
    const { create_time } = this.searchParams
    return [create_time.start_time*1000, create_time.end_time*1000].filter(Boolean)
  }

  /** 创建时间 写入 */
  set createTime(value: number[]) {
    const [start_time, end_time] = value || []
    this.searchParams.create_time.start_time = parseInt(start_time/1000+'');
    this.searchParams.create_time.end_time = parseInt(end_time/1000+'');
  }  
  /*************************************** method ****************************************/
  get memberList() {
    return (row) => {
      return [
        ...(row.member_data||[]).map((item) => {
          return {
            name: item.member_name,
            is_dp: false
          };
        }),
        ...(row.department_data||[]).map((item) => {
          return {
            name: item.department_name,
            is_dp: true
          };
        })
      ];
    };
  }
  // 去发布、编辑联系我页面
  private goEdit(id?: string) {
    this.$router.push({
      name: id
          ? "myapp_enterprise_drainage_contact-me-coach_edit-contact-me"
          : "myapp_enterprise_drainage_contact-me-coach_add-contact-me",
      query: { id }
    });
  }
  private goDetail(id?: string) {
    this.$router.push({
      name:  "myapp_enterprise_drainage_contact-me-coach_detail",
      query: { id }
    });
  }
  // 获取渠道下拉列表数据
  private getChannelList() {
    const params = {
      page: 1,
      size: 9999
    };
    Api.channelsListApi(params).then((res: any) => {
      if (res.code == 200) {
        if (res.data.list) {
          this.channelList = [
            {
              value: 0,
              label: "全部"
            },
            ...res.data.list.map((item) => {
              return {
                value: item.id,
                label: item.name
              };
            })
          ];
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 联系我列表接口
  private getContactList() {
    const params = {
      ...this.searchParams,
      page: this.meta.page,
      size: this.meta.size
    };
    Api.contactListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.contactList = res.data.list;
        this.meta.total = res.data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 搜索查询
  private search() {
    this.meta.page = 1;
    this.searchLoading = true;
    this.getContactList();
  }
  // 重置
  private reset() {
    this.searchParams = JSON.parse(JSON.stringify(defaultSearchParams));
    this.meta.page = 1;
    this.searchLoading = false; //是否正在查询重置
    this.getChannelList();
    this.getContactList();
  }
  // 内部绑定函数（获取每页数据量）
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getContactList();
  }
  // 内部绑定函数（获取当前第几页）
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getContactList();
  }
  
  // 点击下载按钮
  private downLoad(item: any) {
    this.currentDownLoadItem = item;
    this.showDownLoad = true;
    // downloadFile(item.zip_url);
  }
  // 删除联系我接口
  private deleteContact(id: number) {
    this.$confirm("你确认要删除该联系我吗？", "删除联系我", {
      confirmButtonText: "确定删除",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        Api.contactDeleteApi({
          id: id
        }).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              type: "success",
              message: "删除成功!"
            });
            this.reset();
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
  private mounted() {
    this.getChannelList();
    this.getContactList();
  }
}
</script>
<style lang="scss" scoped>.code-img {
  width: 50px;
  height: 50px;
}

/deep/.el-dropdown-menu__item.color-red {
  color: $color-danger;
}

/deep/.el-dropdown-menu__item.color-edit {
  color: $color-text-button;
}

/deep/.el-tag.el-tag--blue {
  color: #333 !important;
  background-color: #EFF6FF !important;
  border-color: #EFF6FF !important;
}

.tags {
  display: flex;
  flex-wrap: nowrap;
  gap: 5px;
  justify-content: center;
}

.table-empty-label {
  display: flex;
  height: 340px;
  flex-direction: column;

  img {
    width: 124px;
    margin: 80px auto 0;
  }

  div {
    line-height: 20px;
    margin-top: 14px;
  }
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.tag {
  max-width: 100px;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.staff-icon {
  width: 14px;
  height: 14px;
  margin-right: 5px;
  vertical-align: middle;
  
}
</style>
