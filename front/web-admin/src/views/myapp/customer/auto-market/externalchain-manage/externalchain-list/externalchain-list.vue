<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button
          type="primary"
          size="small"
          @click="lookItemInfoFn('add', 0)"
        >新增外链</el-button>
      </template>
      <template slot="search">
        <div
          class="search-project"
          :class="{ 'advanced-search-project': advancedSearch }"
        >
          <search-item text="一级渠道：">
            <cy-select
              width-size="small"
              size="small"
              placeholder="请选择一级渠道"
              @change="getSecondChannelData"
              v-model="searchBox.one_level_id"
              :optionsList="searchBox_opt.levelOneChannels"
              :options="{ label: 'name', value: 'id' }"
              search-input
            >
            </cy-select>
          </search-item>
          <search-item text="二级渠道：">
            <cy-select
              width-size="small"
              size="small"
              placeholder="请选择二级渠道"
              v-model="searchBox.two_level_id"
              :optionsList="searchBox_opt.levelSecondChannels"
              :options="{ label: 'name', value: 'id' }"
              search-input
            >
            </cy-select>
          </search-item>
          <search-item text="页面名称：">
            <cy-select
              width-size="small"
              size="small"
              placeholder="请选择页面名称："
              v-model="searchBox.page_name"
              :optionsList="searchBox_opt.pageNames"
              :options="{ label: 'label', value: 'value' }"
            >
            </cy-select>
          </search-item>
          <search-item text="创建人：">
            <cy-input
              width-size="small"
              size="small"
              placeholder="请输入创建人"
              v-model="searchBox.operator"
              clearable
              prefix-icon="el-icon-search"
            ></cy-input>
          </search-item>
          <search-item text="创建时间：">
            <el-date-picker
              class="cy-el-date-picker"
              v-model="searchBox.createTime"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
            >
            </el-date-picker>
          </search-item>
        </div>
        <div class="search-bottom">
          <div class="btn-list">
            <el-button
              size="small"
              type="primary"
              @click="searchDataFn"
            >查询</el-button>
            <el-button
              size="small"
              @click="resetFn"
            >重置</el-button>
          </div>
          <div
            class="senior-search"
            :class="{ 'general-search': advancedSearch }"
            @click="highSearch"
          >
            高级搜索
            <i></i>
          </div>
        </div>
      </template>
      <template slot="table">
        <el-table
          :data="channelList"
          stripe
          style="padding-right:20px;box-sizing:border-box;"
        >
          <el-table-column
            prop="one_level_name"
            label="一级渠道名称"
            min-width="180"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top-start"
                :content="scope.row.one_level_name"
                trigger="hover"
              >
                <div
                  class="channel-name"
                  slot="reference"
                >
                  {{ scope.row.one_level_name }}
                </div>
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            prop="two_level_name"
            label="二级渠道名称"
            min-width="180"
          >
            <template slot-scope="scope">
              <el-popover
                placement="top-start"
                :content="scope.row.two_level_name"
                trigger="hover"
              >
                <div
                  class="channel-name"
                  slot="reference"
                >
                  {{ scope.row.two_level_name }}
                </div>
              </el-popover>
            </template>
          </el-table-column>
          <el-table-column
            label="页面名称"
            min-width="130"
          >
            <template slot-scope="scope">
              {{ nowPageNames(scope.row.url_type) }}
            </template>
          </el-table-column>
          <el-table-column
            prop="created_at"
            label="创建时间"
            min-width="180"
          >
            <template slot-scope="scope">
              {{ scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
            </template>
          </el-table-column>
          <el-table-column
            prop="operator"
            label="创建人"
            min-width="130"
          >
          </el-table-column>
          <el-table-column
            label="操作"
            width="260"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                type="text"
                :disabled="scope.row.short_url === ''"
                @click="copyLink(scope.row.short_url)"
              >复制通用短外链</el-button>
              <el-button
                type="text"
                @click="lookItemInfoFn(scope.row.id, 2)"
              >查看</el-button>
              <el-dropdown
                class="m-l-20"
                @command="handleCommand"
              >
                <el-button
                  size="small"
                  type="text"
                >更多<i
                    aria-hidden="true"
                    class="el-icon-arrow-down el-icon--right"
                  ></i></el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item
                    class="color-edit"
                    :command="`edit_${scope.row.id}`"
                  >编辑</el-dropdown-item>
                  <el-dropdown-item
                    class="color-red"
                    :command="`delete_${scope.row.id}`"
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
            <p>{{ searchLoading ? "查询无结果" : "暂无数据" }}</p>
          </div>
        </el-table>
      </template>
      <template
        slot="footer"
        v-if="channelList.length > 0"
      >
        <el-pagination
          @size-change="handleSizeChangeFn"
          @current-change="handleCurrentChangeFn"
          :current-page="searchBox.current_page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="currentSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
        >
        </el-pagination>
      </template>
    </layout-list>
  </div>
</template>

<script lang="ts">
import LayoutList from "@/components/public/layout-list/index.vue";
import API from "@/api/customer";
export default {
  name: "externlchainList",
  data() {
    return {
      searchLoading: false, // 搜索结果表示无结果
      currentSize: 10, // 默认的一页数据
      advancedSearch: false, // 默认不是高级搜索
      searchBox_opt: {
        levelOneChannels: [],
        levelSecondChannels: [],
        pageNames: [
          {
            label: "全部",
            value: ""
          },
          {
            label: "首页",
            value: 1
          },
          {
            label: "商品详情页",
            value: 2
          },
          {
            label: "活动页",
            value: 3
          },
          {
            label: "直播页",
            value: 4
          },
          {
            label: "活动",
            value: 5
          }
        ]
      },
      searchBox: {
        one_level_id: "",
        two_level_id: "",
        page_name: "",
        createTime: [],
        operator: "",
        current_page: 1,
        page_num: 10
      },
      channelList: [],
      total: 0
    };
  },
  // 监听
  watch: {},
  // 组件实例创建完成
  created() {
    (this as any).getLevelOneChannelData();
    (this as any).getDataFn();
  },
  // vue方法
  methods: {
    /**
     * 获取一级渠道的数据
     *
     * @method getLevelOneChannelData
     */
    async getLevelOneChannelData(): Promise<void> {
      const nowEnterpriseId = Number(
        (this as any).$store.state.user.USERINFO.enterprise_id || ""
      );
      let list: any = [];
      const data: object = {
        enterprise_id: nowEnterpriseId,
        pid: 0
      };
      try {
        const res: any = await API.geChannelName(data);
        if (res.code === 200) {
          list = res.data.list;
        }
      } catch (error) {
        console.log(error);
        console.warn("获取渠道数据失败！");
      }
      (this as any).searchBox_opt.levelOneChannels = list;
    },
    /**
     * 获取二级渠道的数据
     *
     * @method getSecondChannelData
     * @param 渠道ID
     */
    async getSecondChannelData(channelId: number) {
      const nowEnterpriseId = Number(localStorage.getItem("ENTERPRISE_ID"));
      let list: any = [];
      const data: object = {
        enterprise_id: nowEnterpriseId,
        pid: channelId
      };
      try {
        const res: any = await API.geChannelName(data);
        if (res.code === 200) {
          list = res.data.list;
        }
      } catch (error) {
        console.log(error);
        console.warn("获取渠道数据失败！");
      }
      (this as any).searchBox_opt.levelSecondChannels = list;
    },
    /**
     * 高级搜索
     *
     * @method highSearch
     */
    highSearch(): void {
      (this as any).advancedSearch = !(this as any).advancedSearch;
    },
    /**
     * 查看外链的详情/新增一个外链
     *
     * @method lookItemInfoFn
     * @param {item 如果是新增那么这个是0，否则是当前对象的相关id}
     * @param {type 0表示新增 1表示编辑 2表示查看}
     */
    lookItemInfoFn(item: any, type: any) {
      const nowPath =
        type === 0
          ? "myapp_customer_auto-market_externalchain-manage_externalchain-set"
          : type === 1
          ? "myapp_customer_auto-market_externalchain-manage_externalchain-edit"
          : "myapp_customer_auto-market_externalchain-manage_externalchain-look";
      const nowParams = type === 0 ? "" : item;
      (this as any).$router.push({
        name: nowPath,
        query: {
          id: nowParams,
          type: type
        }
      });
    },
    nowPageNames(data: number) {
      console.log(data, 999);
      switch (data) {
        case 1:
          return "首页";
        case 2:
          return "商品详情页";
        case 3:
          return "活动页";
        case 4:
          return "直播";
        case 5:
          return "活动";
        default:
          break;
      }
    },
    /**
     * 复制链接
     *
     * @method copyLink
     * @param {data 要复制的那条数据的信息}
     */
    copyLink(data: string) {
      const url: string = data;
      if (url === "") return false;
      const oInput: HTMLTextAreaElement = document.createElement("textarea");
      oInput.value = url;
      document.body.appendChild(oInput);
      oInput.select(); // 选择对象;
      document.execCommand("Copy"); // 执行浏览器复制命令
      (this as any).$message({
        message: "短外链复制成功",
        type: "success",
        duration: 2000
      });
      oInput.remove();
    },
    handleCommand(command: string) {
      const [key, id] = command.split("_");
      console.log(command, 99);
      switch (key) {
        case "edit":
          (this as any).lookItemInfoFn(id, 1);
          break;
        case "delete":
          (this as any).deletelink(id);
          break;
        default:
          break;
      }
    },
    /**
     * 删除短链接
     *
     * @method deletelink
     * @param {id 此链接的唯一标识}
     */
    async deletelink(id: number | string) {
      (this as any)
        .$confirm("是否需要删除该短链", "删除确认", {
          confirmButtonText: "确认删除",
          confirmButtonClass: "button-danger",
          cancelButtonText: "我再想想",
          type: "warning"
        })
        .then(async () => {
          const data: object = {
            id: Number(id)
          };
          try {
            const res: any = await API.deleteChannelData(data);
            if (res.code === 200) {
              (this as any).$message({
                message: "删除成功",
                type: "success",
                duration: 2000
              });
              (this as any).getDataFn();
            }
          } catch (error) {
            console.log(error);
            console.warn("删除数据失败");
          }
        })
        .catch(() => {
          (this as any).$message({
            type: "info",
            message: "已取消删除"
          });
        });
    },
    /**
     * 点击搜索
     * @method searchDataFn
     */
    searchDataFn(): void {
      (this as any).searchBox.current_page = 1;
      (this as any).getDataFn();
      (this as any).searchLoading = true;
    },
    /**
     * 获取搜索数据的方法
     *
     * @method getDataFn
     */
    async getDataFn(): Promise<void> {
      let list: any = [];
      let total: any = 0;
      try {
        const data: object = {
          enterprise_id:
            (this as any).$store.state.user.USERINFO.enterprise_id || "",
          one_level_id: Number((this as any).searchBox.one_level_id),
          two_level_id: Number((this as any).searchBox.two_level_id),
          url_type: Number((this as any).searchBox.page_name),
          create_start:
            Date.parse((this as any).searchBox.createTime[0]) / 1000 || 0,
          create_end:
            Date.parse((this as any).searchBox.createTime[1]) / 1000 || 0,
          operator: (this as any).searchBox.operator,
          current_page: (this as any).searchBox.current_page,
          page_num: (this as any).searchBox.page_num
        };
        const res = await API.searchChannelData(data);
        list = res.data.list;
        total = res.data.total;
      } catch (error) {
        console.log(error);
        console.warn("获取渠道数据失败！");
        list = [];
      }
      (this as any).channelList = list;
      (this as any).total = total;
    },
    /**
     * 重置方法
     * @method resetFn
     */
    resetFn(): void {
      (this as any).searchBox = {
        one_level_id: "",
        two_level_id: "",
        page_name: "",
        createTime: [],
        operator: "",
        current_page: 1,
        page_num: 10
      };
      (this as any).searchBox_opt.levelSecondChannels = [];
      (this as any).getDataFn();
      (this as any).searchLoading = false;
    },
    /**
     * 修改当页数据的方法handleSizeChangeFn
     * @param { val number 此时选择的一页的数据}
     * @method handleSizeChangeFn
     */
    handleSizeChangeFn(val: number) {
      (this as any).searchBox.page_num = val;
      (this as any).getDataFn();
    },
    /**
     * 修改当页数据的方法handleCurrentChangeFn
     * @param { val number 当前的页数}
     * @method handleCurrentChangeFn
     */
    handleCurrentChangeFn(val: number) {
      (this as any).searchBox.current_page = val;
      (this as any).getDataFn();
    }
  },
  components: {
    LayoutList
  }
};
</script>
<style lang="scss">
#crs-scrm .el-range-editor {
  height: 36px !important;
}
</style>
<style lang="scss" scoped>
.channel-name {
  width: 160px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.search-project {
  height: auto;
  transition: all 0.6s;
  -moz-transition: all 0.6s;
  /* Firefox 4 */
  -webkit-transition: all 0.6s;
  /* Safari 和 Chrome */
  -o-transition: all 0.6s;
  /* Opera */
}

.advanced-search-project {
  height: 53px;
  transition: all 0.6s;
  -moz-transition: all 0.6s;
  /* Firefox 4 */
  -webkit-transition: all 0.6s;
  /* Safari 和 Chrome */
  -o-transition: all 0.6s;
  /* Opera */
  overflow: hidden;
}

.search-bottom {
  overflow: hidden;

  .senior-search {
    float: left;
    cursor: pointer;
    line-height: 36px;
    color: #6F3112;
    width: 78px;
    height: 36px;
    position: relative;

    i {
      display: block;
      right: 7px;
      width: 10px;
      top: 0;
      height: 36px;
      position: absolute;
      background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAMAAADXqc3KAAAARVBMVEX19vpwjfX09fqDnPXz9PqzwvdyjvXX3vjk6fmPpfZxjfWuvve6yPd4k/V+mPXS2/h5lPW/y/iarvbr7vm7yPfH0fjH0viknlVUAAAAUklEQVR4XtWOOQ6AIAAEWe771v8/1VCQoKE2Ot3ONkNehtK9V4ypjeZCAlLwp9cG1jkLo+/+DIiJkBQRjtV7oLYxWgX8cvRcZlXJ/UtVg1n1Py5txgOvnHa61gAAAABJRU5ErkJggg==") no-repeat right center;
      background-size: 13px auto;
    }
  }

  .general-search {
    i {
      display: block;
      right: 7px;
      top: 0;
      width: 10px;
      height: 36px;
      position: absolute;
      transform: rotate(180deg);
      -ms-transform: rotate(180deg);
      /* IE 9 */
      -moz-transform: rotate(180deg);
      /* Firefox */
      -webkit-transform: rotate(180deg);
      /* Safari 和 Chrome */
      -o-transform: rotate(180deg);
      /* Opera */
      background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAMAAADXqc3KAAAARVBMVEX19vpwjfX09fqDnPXz9PqzwvdyjvXX3vjk6fmPpfZxjfWuvve6yPd4k/V+mPXS2/h5lPW/y/iarvbr7vm7yPfH0fjH0viknlVUAAAAUklEQVR4XtWOOQ6AIAAEWe771v8/1VCQoKE2Ot3ONkNehtK9V4ypjeZCAlLwp9cG1jkLo+/+DIiJkBQRjtV7oLYxWgX8cvRcZlXJ/UtVg1n1Py5txgOvnHa61gAAAABJRU5ErkJggg==") no-repeat right center;
      background-size: 13px auto;
    }
  }

  .btn-list {
    float: left;
    width: 180px;
  }
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
