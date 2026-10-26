<template>
  <div class="personalise-list">
    <div class="container_content">
      <div class="button-group mt-20">
        <el-button type="primary" @click="createSpecialWelcome" size="small">
          新增个性化欢迎语
        </el-button>
        <p class="tips m-t-10">*此欢迎语适用于个性化条件配置下添加的企微好友，若未配置则发送默认欢迎语。</p>
      </div>
      <div class="search_content">
        <div class="line special-line">
          <search-item text="欢迎语：">
            <cy-input
              width-size="small"
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入欢迎语关键词"
              v-model="specialKeyword"
              clearable
            ></cy-input>
          </search-item>
          <div class="search_btn ">
            <el-button type="primary" @click="searchSpecial" size="small"
              >查询</el-button
            >
            <el-button @click="resetSpecial" size="small">重置</el-button>
          </div>
        </div>
      </div>
      <div class="section">
        <el-table :data="list" stripe tooltip-effect="light">
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
          <el-table-column label="条件范围" width="200">
            <template>
              第三方数据
            </template>
          </el-table-column>
          <el-table-column prop="type_name" label="条件内容" width="200">
          </el-table-column>
          <el-table-column
            prop="created_at"
            label="创建时间"
            :formatter="formatDate"
            width="200"
          >
          </el-table-column>
          <el-table-column
            prop="updated_at"
            label="编辑时间"
            :formatter="formatDate1"
            width="200"
          >
          </el-table-column>
          <el-table-column label="操作" width="160" fixed="right">
            <template slot-scope="scope">
              <el-button type="text" size="small" @click="editRowFn(scope.row)"
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
            v-if="list && list.length > 0"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="meta.page"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="meta.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="meta.total"
          >
          </el-pagination>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { dateFormat } from "../../../../customer/auto-market/util/index";
import API from "@/api/enterprise/drainage-coach";
@Component
export default class extends Vue {
  // data
  private specialKeyword = "";
  private searchLoading = false;
  private list: Array<any> = [];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // 生命周期
  private mounted() {
    console.log("hhh");
    this.getSpecialList();
  }
  private formatDate(row: any) {
    return row.created_at
      ? dateFormat(row.created_at, "yyyy-MM-dd hh:mm:ss")
      : "";
  }
  private formatDate1(row: any) {
    return row.updated_at
      ? dateFormat(row.updated_at, "yyyy-MM-dd hh:mm:ss")
      : "";
  }
  // 新增个性化欢迎语
  private createSpecialWelcome() {
    this.$router.push({
      name: "myapp_enterprise_drainage_welcome_newcreate"
    });
  }

  private editRowFn(row: any) {
    this.$router.push({
      name: "myapp_enterprise_drainage_welcome_newedit",
      query: { id: row.id, type: "1" }
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
  private async deleteRowFn(row: any) {
    const res: any = await API.drainageDeleteApi({ id: row.id });
    if (res.code === 200) {
      this.$message({
        message: "删除成功",
        type: "success",
        duration: 2000
      });
      this.getSpecialList();
    } else {
      this.$message.error(res.message);
    }
  }

  private resetSpecial() {
    this.specialKeyword = "";
    this.searchLoading = false;
    this.getSpecialList();
  }

  private searchSpecial() {
    this.searchLoading = true;
    this.getSpecialList();
  }

  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getSpecialList();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getSpecialList();
  }

  private async getSpecialList() {
    const params = {
      page: this.meta.page,
      size: this.meta.size,
      key_words: this.specialKeyword
    };
    const res: any = await API.drainageWelcomeAll(params);
    if (res.code === 200) {
      res.data.list.forEach((v: any) => {
        let arr: any = [];
        let str = "";
        v.content_text.forEach((v: any) => {
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
      this.list = res.data.list;
      this.meta.total = res.data.total;
    } else {
      this.$message.error(res.message);
    }
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
