<template>
  <div class="welcome_list">
    <div class="container_content">
      <div class="button-group">
        <el-button type="primary" @click="createWelcomeFn" size="small">
          新增客户欢迎语
        </el-button>
        <!-- <p class="tips m-t-10">*当个性化欢迎语没有配置欢迎语时，默认展示此欢迎语，此欢迎语适用于联系我关联的渠道添加的企微好友</p> -->
        <p class="tips m-t-10">*此欢迎语适用于联系我关联的渠道添加的企微好友，若未配置则发送个性化欢迎语。</p>
      </div>
      <div class="search_content">
        <div class="line">
          <search-item text="欢迎语：">
            <cy-input
              width-size="small"
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入欢迎语关键词"
              v-model="searchGroup.keyword"
              clearable
            ></cy-input>
          </search-item>
          <search-item text="发送渠道：">
            <cy-select
              width-size="small"
              size="small"
              v-model="searchGroup.channel"
              placeholder="请选择渠道"
              :optionsList="channelList"
              :options="{ label: 'channel', value: 'id' }"
              search-input
            >
            </cy-select>
          </search-item>
        </div>
        <div class="search_btn">
          <el-button type="primary" @click="searchData" size="small"
            >查询</el-button
          >
          <el-button @click="resetFn" size="small">重置</el-button>
        </div>
      </div>
      <div class="section">
        <el-table :data="welcomeTable.table" stripe tooltip-effect="light">
          <el-table-column
            label="欢迎语"
            min-width="380"
            width="400"
            fixed="left"
          >
            <template slot-scope="scope">
              <div class="aligns-center">
                <el-popover
                  placement="top-start"
                  width="250"
                  trigger="hover"
                  :disabled="scope.row.text.length <= 20"
                >
                  <div
                    style="overflow:hidden"
                    v-html="
                      scope.row.text
                        .replace(/\n/g, '<br/>')
                        .replace(/\s/g, '&nbsp;')
                    "
                  ></div>
                  <span
                    v-if="scope.row.text.length <= 20"
                    class="m-l-10 goods-name el-popover-ref"
                    slot="reference"
                    >{{ scope.row.text }}</span
                  >
                  <span
                    slot="reference"
                    v-else
                    class="m-l-10 goods-name el-popover-ref"
                    >{{ scope.row.text.substr(0, 20) + "..." }}</span
                  >
                </el-popover>
                <span>
                  {{ scope.row.types }}
                </span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="channel"
            label="发送渠道"
            show-overflow-tooltip
            width="250"
          >
          </el-table-column>
          <el-table-column
            prop="created_at"
            label="创建时间"
            :formatter="formatDate"
            min-width="200"
          >
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template slot-scope="scope">
              <el-button
                type="text"
                size="small"
                @click="editRowFn(scope.row, 2)"
                >编辑</el-button
              >
              <el-button
                type="text"
                size="small"
                slot="reference"
                style="color:#f56c6c;"
                @click="delRow(scope.row)"
                >删除</el-button
              >
            </template>
          </el-table-column>
          <div slot="empty" class="table-empty-label-page">
            <img
              class="empty-img"
              :src="
                searchLoading
                  ? require('@/assets/images/empty-search.png')
                  : require('@/assets/images/empty-box.png')
              "
              alt=""
            />
            <div>
              {{ searchLoading ? "没有找到相关欢迎语" : "暂无客户欢迎语" }}
            </div>
          </div>
        </el-table>
        <div class="table_footer">
          <el-pagination
            @size-change="handleSizeChangeFn"
            @current-change="handlePageChangeFn"
            v-if="welcomeTable.total > 0"
            :current-page="welcomeTable.currentPage"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="welcomeTable.currentSize"
            layout="total, sizes, prev, pager, next, jumper"
            :total="welcomeTable.total"
          >
          </el-pagination>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import { dateFormat } from "../../../../customer/auto-market/util/index";
import API from "@/api/enterprise/drainage-coach";
@Component({
  components: {
    LayoutList
  }
})
export default class extends Vue {
  private defalutWelList: any = {
    table: []
  };
  private welcomeTable = {
    table: [],
    currentSize: 10,
    currentPage: 1,
    total: 0
  };
  handleSizeChangeFn(val: number) {
    this.welcomeTable.currentPage = 1;
    this.welcomeTable.currentSize = val;
    this.getTableDataFn();
  }
  handlePageChangeFn(val: number) {
    this.welcomeTable.currentPage = val;
    this.getTableDataFn();
  }
  formatDate(row: any) {
    return row.created_at
      ? dateFormat(row.created_at, "yyyy-MM-dd hh:mm:ss")
      : "";
  }
  private searchGroup = {
    channel: 0,
    keyword: ""
  };
  private channelList: any = [
    {
      channel: "全部",
      id: 0
    }
  ];
  private searchLoading = false;
  async searchData() {
    this.searchLoading = true;
    this.welcomeTable.currentPage = 1;
    this.getTableDataFn();
  }
  formatUserName(data: any) {
    if (data && data.length > 0) {
      return data.map((v: any) => {
        const item = v;
        if (item.user_names) {
          item.user_names = item.user_names.replace(/,/g, ";");
        }
        return item;
      });
    }
    return [];
  }
  async getTableDataFn() {
    const params: object = {
      page: this.welcomeTable.currentPage,
      size: this.welcomeTable.currentSize,
      content: this.searchGroup.keyword || undefined,
      channel_id: this.searchGroup.channel
      // member_id: this.searchGroup.staff
    };
    const res: any = await API.welcomeListApi(params);
    if (res.code == 200) {
      this.welcomeTable.total = res.data.total;
      res.data.list.forEach((v: any) => {
        let arr: any = [];
        let str = "";
        v.content.forEach((v: any) => {
          if ("text" in v) {
            str = v.text.content;
          }
          if ("image" in v) {
            arr.push("[图片]");
          } else if ("link" in v) {
            arr.push("[链接]");
          } else if ("miniprogram" in v) {
            arr.push("[小程序]");
          } else if ("video" in v) {
            arr.push("[视频]");
          } else if ("file" in v) {
            arr.push("[文件]");
          }
        });
        arr = [...new Set(arr)].join(" ");
        v.types = arr;
        v.text = str;
      });
      this.welcomeTable.table = this.formatUserName(res.data.list);
    } else {
      this.$message.error(res.message);
    }
  }
  async getDetaultTableDataFn() {
    const res: any = await API.welcomeDefaultListApi({});
    if (res.code == 200) {
      let arr: any = [];
      let str = "";
      res.data.content.forEach((v: any) => {
        if ("text" in v) {
          str = v.text.content;
        }
        if ("image" in v) {
          arr.push("[图片]");
        } else if ("link" in v) {
          arr.push("[链接]");
        } else if ("miniprogram" in v) {
          arr.push("[小程序]");
        } else if ("video" in v) {
          arr.push("[视频]");
        } else if ("file" in v) {
          arr.push("[文件]");
        }
      });
      arr = [...new Set(arr)].join(" ");
      res.data.types = arr;
      res.data.text = str;
      this.defalutWelList.table.push(res.data);
    } else {
      this.$message.error(res.message);
    }
  }
  resetFn() {
    this.searchGroup = {
      channel: 0,
      keyword: ""
    };
    this.welcomeTable.currentPage = 1;
    this.welcomeTable.currentSize = 10;
    this.searchLoading = false;
    this.getTableDataFn();
  }
  editRowFn(row: any, type: any) {
    this.$router.push({
      name: "myapp_enterprise_drainage_welcome_edit",
      query: { id: row.id, type: type }
    });
  }
  private delRow(row) {
    this.$confirm("确认删除该欢迎语吗", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        this.deleteRowFn(row);
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  async deleteRowFn(row: any) {
    const res: any = await API.welcomeDeleteApi({ id: row.id });
    if (res.code === 200 && res.data.result) {
      this.$message({
        message: "删除成功",
        type: "success",
        duration: 2000
      });
      this.getTableDataFn();
    } else {
      this.$message.error(res.message);
    }
  }
  createWelcomeFn() {
    this.$router.push({
      name: "myapp_enterprise_drainage_welcome_create"
    });
  }
  //获取渠道列表
  async getChannelList() {
    const params: object = {
      add_channel: false
    };
    const res: any = await API.welcomeChannelListApi(params);
    if (res.code == 200) {
      const obj: object = { id: 0, channel: "全部" };
      const channelList: Array<any> = res.data;
      channelList.unshift(obj);
      this.channelList = channelList;
    } else {
      this.$message.error(res.message);
    }
  }

  mounted() {
    this.getChannelList();
    this.getDetaultTableDataFn();
    this.getTableDataFn();
  }
}
</script>
<style>
.el-popconfirm__main {
  margin: 10px 0;
}
</style>
<style lang="scss" scoped>
.tips {
  color: #999;
}
.el-popover-ref {
  pointer-events: auto;
}
.welcome_container {
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
  .container_header {
    margin-bottom: 30px;
    h3 {
      margin-bottom: 10px;
      font-weight: bold;
    }
    p {
      margin-bottom: 10px;
      color: #999;
    }
    .aligns-center {
      display: flex;
      align-items: center;
      .goods-name {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
    }
  }
}
.container_content {
  background-color: #fff;
  padding-bottom: 10px;
  border-radius: 10px;
  .button-group {
    margin-bottom: 20px;
  }
  .search_content {
    background-color: #f5f6fa;
    padding: 20px;
    margin-bottom: 20px;
  }
  .section {
    .marginR20 {
      margin-right: 20px;
    }
    .aligns-center {
      display: flex;
      align-items: center;
      .goods-name {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }
    }

    .deleteBtn.el-button--text {
      color: #f67170;
    }
    .deleteBtn.el-button--text:hover {
      color: #f56c6c;
    }
    .table_footer {
      margin: 20px 0;
      text-align: right;
    }
  }
}
.table-empty-label-page {
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
/deep/ .el-tabs__nav {
  .is-active {
    color: #fff;
    background-color: #708df5;
  }
}
.special-line {
  display: inline-flex;
}
.mt-20 {
  margin-top: 20px;
}
</style>
