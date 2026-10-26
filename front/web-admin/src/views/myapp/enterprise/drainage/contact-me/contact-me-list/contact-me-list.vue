<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button
          type="primary"
          @click="goEdit(0)"
          size="small"
        >新增联系我</el-button>
      </template>
      <template slot="search">
        <search-item text="二维码：">
          <cy-input
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入二维码名称"
            clearable
            class="custom-name"
            v-model="search.name"
          ></cy-input>
        </search-item>
        <search-item text="所属渠道：">
          <cy-select
            v-model="search.channel.value"
            placeholder="全部"
            search-input
            :optionsList="search.channel.list"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </search-item>
        <div>
          <el-button
            type="primary"
            @click="clickSearch()"
            size="small"
          >查询</el-button>
          <el-button
            @click="reset()"
            size="small"
          >重置</el-button>
        </div>
      </template>
      <template slot="table">
        <el-table
          :data="contactList"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            label="二维码名称"
            width="150"
            align="center"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top-start"
                width="300"
                trigger="hover"
                :disabled="scope.row.code_name.length <= 12"
              >
                <div>{{ scope.row.code_name || "--" }}</div>
                <span
                  slot="reference"
                  v-if="scope.row.code_name.length <= 12"
                >{{ scope.row.code_name || "--" }}</span>
                <span
                  slot="reference"
                  v-if="scope.row.code_name.length > 12"
                >{{
                  scope.row.code_name.substr(0, 12) + "..."
                }}</span>
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            label="二维码"
            width="100"
            align="center"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top-start"
                title=""
                trigger="hover"
              >
                <img
                  slot="reference"
                  class="code-img"
                  :src="scope.row.qr_code_url"
                  alt=""
                />
                <img
                  :src="scope.row.qr_code_url"
                  alt=""
                  style="width: 150px;height: 150px"
                />
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            prop="channel_name"
            label="所属渠道"
            min-width="120"
            align="center"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top-start"
                width="300"
                trigger="hover"
                :disabled="scope.row.channel_name.length <= 12"
              >
                <div>{{ scope.row.channel_name || "--" }}</div>
                <span
                  slot="reference"
                  v-if="scope.row.channel_name.length <= 12"
                >{{ scope.row.channel_name || "--" }}</span>
                <span
                  slot="reference"
                  v-if="scope.row.channel_name.length > 12"
                >{{ scope.row.channel_name.substr(0, 12) + "..." }}</span>
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            prop="member_count"
            label="部门/员工"
            min-width="320"
            align="center"
          >
        
            <template slot-scope="scope">
              <!-- <el-tag v-for="(item, index) in scope.row.category_names" :key="index">{{item}}</el-tag> -->
              <el-tooltip placement="top" effect="light" v-if="memberList(scope.row).length > 3">
                <div class="tags" style="cursor: pointer;">
                  <!-- 显示前2个 -->
                  <el-tag class="tag" v-for="item,index in memberList(scope.row).slice(0,3)" :key="index">
                    <!-- icon -> 部门 -->
                    <img
                      class="staff-icon"
                      v-if="item.is_dp"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                      alt=""
                    />
                    <!-- icon -> 人 -->
                    <img
                      class="staff-icon"
                      v-else
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/staff_icon.png"
                      alt=""
                    />
                    {{item.name}}
                  </el-tag>
                  <span>...</span>
                </div>
                <div class="tags" slot="content">
                  <el-tag v-for="item,index in memberList(scope.row)" :key="index">
                    <!-- icon -> 部门 -->
                    <img
                      class="staff-icon"
                      v-if="item.is_dp"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                      alt=""
                    />
                    <!-- icon -> 人 -->
                    <img
                      class="staff-icon"
                      v-else
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/staff_icon.png"
                      alt=""
                    />
                    {{item.name}}
                  </el-tag>
                </div>
              </el-tooltip>
              <div v-else-if="memberList(scope.row).length <= 3 && memberList(scope.row).length > 0" class="tags">
                <el-tag v-for="item,index in memberList(scope.row)" :key="index">
                  <!-- icon -> 部门 -->
                  <img
                    class="staff-icon"
                    v-if="item.is_dp"
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                    alt=""
                  />
                  <!-- icon -> 人 -->
                  <img
                    class="staff-icon"
                    v-else
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/staff_icon.png"
                    alt=""
                  />
                  {{item.name}}</el-tag>
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="sum_add_member_number"
            label="新增好友数"
            min-width="150"
            align="center"
          ></el-table-column>
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
            width="150"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                type="text"
                @click="downLoad(scope.row)"
              >
                下载
              </el-button>
              <el-dropdown
                class="m-l-20"
                @command="handleCommand"
              >
                <el-button type="text">更多<i
                    aria-hidden="true"
                    class="el-icon-arrow-down el-icon--right"
                  ></i></el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item
                    :command="'edit_' + scope.row.id"
                    class="color-edit"
                  >编辑</el-dropdown-item>
                  <el-dropdown-item
                    class="color-red"
                    :command="'delete_' + scope.row.id"
                  >删除</el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
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
    ></download-code>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/drainage";
import downloadCode from "../../components/download-code.vue";
@Component({
  components: {
    downloadCode
  }
})
export default class extends Vue {
  /*************************************** data ****************************************/
  private search = {
    name: "", // 二维码
    channel: {
      list: [],
      value: "全部"
    }
  };
  private contactList: Array<any> = []; // 联系我列表数组
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private searchLoading = false; // 是否正在搜索
  private showDownLoad = false; // 控制下载二维码的弹框显隐
  private currentDownLoadItem: any; // 当前下载二维码的url
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
  private goEdit(is_edit?: string, id?: string) {
    this.$router.push({
      name: `${
        is_edit
          ? "myapp_enterprise_drainage_contact-me_edit-contact-me"
          : "myapp_enterprise_drainage_contact-me_add-contact-me"
      }`,
      query: {
        is_edit,
        id
      }
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
        const obj = { id: 0, name: "全部" };
        if (res.data.list) {
          res.data.list.unshift(obj);
          this.search.channel.list = res.data.list;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 联系我列表接口
  private getContactList() {
    const params = {
      channel_id:
        (this.search.channel.value !== "全部" && this.search.channel.value) ||
        undefined,
      name: this.search.name || undefined,
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
  private clickSearch() {
    this.meta.page = 1;
    this.searchLoading = true;
    this.getContactList();
  }
  // 重置
  private reset() {
    this.search.name = "";
    this.search.channel.value = "全部";
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
  // 点击操作-更多按钮
  private handleCommand(command: string) {
    if (command.indexOf("edit") > -1) {
      const id = command.replace("edit_", "");
      this.goEdit("1", id);
    }
    if (command.indexOf("delete") > -1) {
      const id = Number(command.replace("delete_", ""));
      this.deleteContact(id);
    }
  }
  // 点击下载按钮
  private downLoad(item: any) {
    this.currentDownLoadItem = item;
    this.showDownLoad = true;
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
          ids: [id]
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
