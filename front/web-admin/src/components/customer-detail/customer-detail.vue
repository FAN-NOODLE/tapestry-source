<template>
  <div>
    <layout-list>
      <!-- 顶部标签页 -->
      <template slot="top-tabs">
        <el-tabs v-model="activeName" @tab-click="resetFn">
          <el-tab-pane label="客户详情" name="first"></el-tab-pane>
          <el-tab-pane label="客户足迹" name="third"></el-tab-pane>
        </el-tabs>
      </template>
      <!-- 用户基础信息部分 -->
      <template slot="tabs">
        <div v-if="activeName == 'first'">
          <p class="tags-h">基础信息</p>
          <div class="tabs-firstrow">
            <el-image fit="cover"
              :src="
                lists.avatar
                  ? lists.avatar
                  : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
              "
              class="user-avatar m-b-20"
            >
            </el-image>
            <div class="message-grid">
              <span class="w-6">客户姓名: {{ lists.user_name || "-" }}</span>
              <span>生日: {{ lists.birthday || "-" }}</span>
              <span></span>
              <span v-if="lists.mobile"
                >手机号:
                {{
                  lists.mobile.substr(0, 3) +
                    "****" +
                    lists.mobile.substr(7, 11)
                }}</span
              >
              <span v-else>手机号:-</span>
            </div>
          </div>
          <div>
            <span class="user-remark">备注：</span>
            {{ lists.remark || "-" }}
          </div>
          <hr />
          <div class="customer-label">
            <span class="m-r-10">客户标签</span>
            <el-col :span="20">
              <el-tag
                v-for="(tag, index) in lists.tag_list"
                :key="index + 'tag'"
                class="m-b-10 m-r-10 customer-labelbutton"
              >
                <el-popover
                  placement="top-start"
                  trigger="hover"
                  :disabled="tag && tag.length < 8"
                >
                  <div>{{ tag }}</div>
                  <span slot="reference" v-if="tag && tag.length < 8">{{
                    tag
                  }}</span>
                  <span slot="reference" v-if="tag && tag.length >= 8">{{
                    tag.substr(0, 8) + "..."
                  }}</span>
                </el-popover>
              </el-tag>
            </el-col>
          </div>
          <div class="customer-label">
            <span class="m-r-10">印象标签</span>
            <el-col :span="20">
              <el-tag
                v-for="(tag, index) in lists.impression_tag"
                :key="index + 'tag'"
                class="m-b-10 m-r-10 customer-labelbutton"
              >
                <el-popover
                  placement="top-start"
                  trigger="hover"
                  :disabled="tag && tag.length < 8"
                >
                  <div>{{ tag }}</div>
                  <span slot="reference" v-if="tag && tag.length < 8">{{
                    tag
                  }}</span>
                  <span slot="reference" v-if="tag && tag.length >= 8">{{
                    tag.substr(0, 8) + "..."
                  }}</span>
                </el-popover>
              </el-tag>
            </el-col>
          </div>
          <div class="m-t-40 m-b-20">
            <!-- gender 1,男 2,女 0,未知 空为- -->
            <span class="m-r-50">性别: {{ lists.sex }}</span>
            <span class="m-r-50">年龄: {{ lists.age }}</span>
          </div>
          <hr />
          <!-- <h3 class="tags-h">收货信息</h3> -->
        </div>
      </template>
      <!-- 客户足迹 -->
      <template slot="table">
        <div v-if="activeName == 'third'" class="customer-footer">
          <el-timeline>
            <el-timeline-item
              v-for="(item, index) in customerFooter"
              :key="index"
              :timestamp="item.created_at | formatDate('yyyy-mm-dd hh:mm:ss')"
              placement="bottom"
            >
              <h3>{{ item.event_txt ? item.event_txt : "--" }}</h3>
              <span>{{ item.event_desc ? item.event_desc : "--" }}</span>
            </el-timeline-item>
          </el-timeline>
          <div v-if="customerFooter.length <= 0" class="customer-Footer">
            <el-image fit="cover" src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png" alt="" />
            <p>暂无客户足迹</p>
          </div>
        </div>
      </template>
      <!-- 底部表格 -->
      <template slot="footer">
        <div v-if="activeName == 'first'" class="first-table">
          <h3 class="tags-h">企业微信</h3>
          <!-- 企微好友表 -->
          <div
            :key="Math.random()"
          >
            <!-- 表一所属员工 -->
            <div>
              <el-table
                :data="firstList"
                max-height="440"
                tooltip-effect="light"
                ref="tableList"
              >
                <el-table-column
                  label="所属员工"
                  prop="member_name"
                  min-width="120"
                >
                  <template slot-scope="scope">
                    <span>{{
                      scope.row.member_name ? scope.row.member_name : "--"
                    }}</span>
                  </template>
                </el-table-column>
                <el-table-column
                  label="员工所在部门"
                  prop="department_list"
                  min-width="170"
                  show-overflow-tooltip
                >
                  <template slot-scope="scope">
                    <span
                      v-if="scope.row.department_list"
                      class="user-department"
                    >
                      <!-- <span
                        v-for="(item, index) in scope.row.department_list"
                        :key="index"
                        >{{ item.department_name }};
                      </span> -->
                      <span>{{ scope.row.dep_name }}</span>
                    </span>
                    <span v-else>--</span>
                  </template>
                </el-table-column>
                <el-table-column label="所属渠道" prop="state" min-width="150">
                  <template slot-scope="scope">
                    <span>{{ scope.row.state ? scope.row.state : "--" }}</span>
                  </template>
                </el-table-column>
                <el-table-column
                  label="添加好友时间"
                  prop="create_time"
                  min-width="220"
                >
                  <template slot-scope="scope">
                    <span>{{
                      scope.row.create_time | formatDate("yyyy-mm-dd hh:mm:ss")
                    }}</span>
                  </template>
                </el-table-column>
                <el-table-column
                  label="添加群数"
                  prop="belong_to_group_chat_count"
                  min-width="100"
                >
                  <template slot-scope="scope">
                    <div>
                      <el-button
                        type="text"
                        v-if="scope.row.belong_to_group_chat_count == 0"
                        >{{ scope.row.belong_to_group_chat_count }}</el-button
                      >
                      <el-button
                        type="text"
                        @click="
                          addGroupNumber(scope.row.belong_to_group_chat_list)
                        "
                        v-else
                        >{{ scope.row.belong_to_group_chat_count }}</el-button
                      >
                    </div>
                  </template>
                </el-table-column>
                <el-table-column
                  label="备注"
                  prop="remark_list"
                  min-width="200"
                >
                  <template slot-scope="scope">
                    <div v-if="scope.row.remark_list">
                      <div class="relations-inside">
                        <div v-if="scope.row.remark_list.length > 2">
                          <el-tag
                            :disable-transitions="false"
                            class="m-b-5 m-r-5"
                          >
                            <el-popover
                              placement="top-start"
                              trigger="hover"
                              :disabled="
                                scope.row.remark_list[0] &&
                                  scope.row.remark_list[0].length <= 10
                              "
                            >
                              <div>{{ scope.row.remark_list[0] }}</div>

                              <span
                                slot="reference"
                                v-if="
                                  scope.row.remark_list[0] &&
                                    scope.row.remark_list[0].length <= 10
                                "
                                >{{ scope.row.remark_list[0] }}</span
                              >
                              <span
                                slot="reference"
                                v-if="
                                  scope.row.remark_list[0] &&
                                    scope.row.remark_list[0].length > 10
                                "
                                >{{
                                  scope.row.remark_list[0].substr(0, 10) + "..."
                                }}</span
                              >
                            </el-popover>
                          </el-tag>
                          <el-tag
                            :disable-transitions="false"
                            class="m-b-5 m-r-5"
                          >
                            <el-popover
                              placement="top-start"
                              trigger="hover"
                              :disabled="
                                scope.row.remark_list[1] &&
                                  scope.row.remark_list[1].length <= 10
                              "
                            >
                              <div>{{ scope.row.remark_list[1] }}</div>

                              <span
                                slot="reference"
                                v-if="
                                  scope.row.remark_list[1] &&
                                    scope.row.remark_list[1].length <= 10
                                "
                                >{{ scope.row.remark_list[1] }}</span
                              >
                              <span
                                slot="reference"
                                v-if="
                                  scope.row.remark_list[1] &&
                                    scope.row.remark_list[1].length > 10
                                "
                                >{{
                                  scope.row.remark_list[1].substr(0, 10) + "..."
                                }}</span
                              >
                            </el-popover>
                          </el-tag>
                          <el-popover
                            placement="top-start"
                            width="350"
                            trigger="click"
                            popper-class="popover-Tip"
                          >
                            <el-tag
                              v-for="(tagRemarks, index) in scope.row
                                .remark_list"
                              :key="'tagRemarks' + index"
                              class="m-b-5 m-r-5"
                            >
                              <el-popover
                                placement="top-start"
                                trigger="hover"
                                :disabled="
                                  tagRemarks && tagRemarks.length <= 20
                                "
                              >
                                <div>{{ tagRemarks }}</div>

                                <span
                                  slot="reference"
                                  v-if="tagRemarks && tagRemarks.length <= 20"
                                  >{{ tagRemarks }}</span
                                >
                                <span
                                  slot="reference"
                                  v-if="tagRemarks && tagRemarks.length > 20"
                                  >{{ tagRemarks.substr(0, 20) + "..." }}</span
                                >
                              </el-popover>
                            </el-tag>
                            <el-button
                              size="small"
                              type="text"
                              slot="reference"
                              icon="el-icon-more"
                            ></el-button>
                          </el-popover>
                        </div>
                        <div
                          v-else-if="
                            scope.row.remark_list.length > 0 &&
                              scope.row.remark_list.length <= 2
                          "
                        >
                          <el-tag
                            v-for="(tagRemarks, index) in scope.row.remark_list"
                            :key="'tagRemarks' + index"
                            :disable-transitions="false"
                            class="m-b-5 m-r-5"
                          >
                            <el-popover
                              placement="top-start"
                              trigger="hover"
                              :disabled="tagRemarks && tagRemarks.length <= 10"
                            >
                              <div v-if="tagRemarks.length !== 0">
                                {{ tagRemarks == " " ? "--" : tagRemarks }}
                              </div>

                              <span
                                slot="reference"
                                v-if="tagRemarks && tagRemarks.length <= 10"
                                >{{ tagRemarks }}</span
                              >
                              <span
                                slot="reference"
                                v-if="tagRemarks && tagRemarks.length > 10"
                                >{{ tagRemarks.substr(0, 10) + "..." }}</span
                              >
                            </el-popover>
                          </el-tag>
                        </div>
                        <div v-else-if="scope.row.remark_list.length == 0">
                          --
                        </div>
                      </div>
                    </div>
                  </template>
                </el-table-column>
                <div slot="empty" class="table-empty-firstlabel">
                  <el-image fit="cover"
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png"
                    alt=""
                  />
                  <p>暂无数据</p>
                </div>
              </el-table>
            </div>
          </div>
          <layout-form>
            <template slot="footer">
              <el-button
                @click="routerReturn"
                class="Refresh m-t-20"
                size="small"
                >返回
              </el-button>
            </template>
          </layout-form>
        </div>
        <div v-if="activeName == 'second'">
          <el-pagination
            class="m-t-20 m-b-30"
            v-if="list.length > 0"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="meta.page"
            :page-sizes="[10, 20, 30, 40]"
            :page-size="meta.size"
            layout="total, sizes, prev, pager, next, jumper"
            :total="meta.total"
          >
          </el-pagination>
          <layout-form>
            <template slot="footer">
              <el-button
                @click="routerReturn"
                class="Refresh m-t-20"
                size="small"
                >返回
              </el-button>
            </template>
          </layout-form>
        </div>
        <div v-if="activeName == 'third'">
          <layout-form>
            <template slot="footer">
              <el-button
                @click="routerReturn"
                class="Refresh m-t-20"
                size="small"
                >返回
              </el-button>
            </template>
          </layout-form>
        </div>
      </template>
    </layout-list>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutList from "@/components/public/layout-list/index.vue";
import LayoutForm from "@/components/public/layout-form/index.vue";
import API from "@/api/customer/index";
// vuex

@Component({
  components: {
    LayoutList,
    LayoutForm
  }
})
export default class extends Vue {
  /********************************************data开始**********************************/
  private activeName = "first"; //top-tab标签页默认显示位置
  private userId = 0; //用户ID
  //表格数据源
  private firstList: Array<any> = [
    {
      member_id: 0, // 员工ID
      member_name: "", // 员工名称
      department_list: [{ department_name: "" }], //所属部门
      field_10: "", //
      create_time: 0, // 添加好友时间
      state: "", // 所属渠道
      belong_to_group_chat_count: 0, // 添加群数
      // 群列表
      belong_to_group_chat_list: [
        {
          member_id: 0, //员工ID
          group_chat_id: 0, // 群ID
          group_chat_name: "", // 群名称
          user_id: 0, // 用户ID
          join_time: 0 //加群时间
        }
      ]
    }
  ]; 
  //表格数据源（tabs为first时）
  private list: Array<any> = [];
  //表格数据源（tabs为second时）
  private lists: any = [];
  //客户足迹
  private customerFooter = [];
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private metas = {
    size: 10,
    page: 1,
    total: 0
  };
  private footprintMetas = {
    size: 999,
    page: 1,
    total: 0
  };
  /********************************************data结束**********************************/
  /***********************************************methods开始****************************/
  // 标签页切换
  private resetFn(value: any) {
    if (value["name"] === "first") {
      //客户详情基本信息 + 收货地址
      this.customerDetail();
      //客户详情所属员工列表
      this.customerBelongToMembers();
    }  else {
      //客户详情足迹
      this.customerFootprintList();
    }
  }
  //储存用户ID
  private saveUserId() {
    this.userId = Number(this.$route.query.user_id);
    sessionStorage.setItem("userIdDate", this.$route.query.user_id as string);
    if (sessionStorage.getItem("userIdDate")) {
      const userId = Number(sessionStorage.getItem("userIdDate"));
      this.userId = userId;
    }
  }
  //客户详情基本信息 + 收货地址API接口
  private customerDetail() {
    const _this = this;
    const params = {
      user_id: _this.userId,
      page: _this.meta.page,
      size: _this.meta.size
    };
    API.customerDetailListApi(params).then((res: any) => {
      if (res.code === 200) {
        if (res.data) {
          _this.lists = res.data;
          // 处理性别
          switch (_this.lists.gender) {
            case 1:
              _this.lists.sex = "男";
              break;
            case 2:
              _this.lists.sex = "女";
              break;
            case 0:
              _this.lists.sex = "未知";
              break;
            default:
              _this.lists.sex = "-";
          }
          // 处理年龄
          if (_this.lists.birthday) {
            const birthday = new Date(_this.lists.birthday.replace(/-/g, "\/"));
            const date = new Date();
            _this.lists.age =
              date.getFullYear() -
              birthday.getFullYear() -
              (date.getMonth() < birthday.getMonth() ||
              (date.getMonth() == birthday.getMonth() &&
                date.getDate() < birthday.getDate())
                ? 1
                : 0);
          } else {
            _this.lists.age = "-";
          }
          if (_this.lists.address_list && _this.lists.address_list.length) {
            _this.lists.address_list.map((item: any) => {
              item.consignee_mobile = _this.changeStr(
                item.consignee_mobile,
                4,
                7
              );
              return item;
            });
          }
        }
      }
    });
  }
  private changeStr(
    str: string,
    start_index: number,
    end_index: number,
    changeStr = "****"
  ) {
    return (
      str.substring(0, start_index - 1) + changeStr + str.substring(end_index)
    );
  }
  //客户详情足迹API接口
  private customerFootprintList() {
    const params = {
      user_id: this.userId,
      page: this.footprintMetas.page,
      size: this.footprintMetas.size
    };
    API.customerDetailFootprintListApi(params).then((res: any) => {
      if (res.code === 200) {
        this.footprintMetas.total = res.data.total;
        if (res.data.list && res.data.list !== 0) {
          this.customerFooter = res.data.list;
        }
      }
    });
  }
  //客户详情所属员工列表API接口
  private customerBelongToMembers() {
    const params = {
      user_id: this.userId,
      page: this.meta.page,
      // size: this.meta.size
      size: 999
    };
    API.customerDetailBelongToMembersApi(params).then((res: any) => {
      if (res.code === 200) {
        if (res.data.list && res.data.list.length !== 0) {
          this.firstList = res.data.list;
          this.firstList
            .map(v => v.department_list)
            .forEach((item, index) => {
              this.firstList[index].dep_name = item
                .map((v1: any) => v1.department_name)
                .join("；");
            });
        } else {
          this.firstList = [];
        }
      }
    });
  }
  //内部绑定函数（获取每页数据量）
  handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
  }
  //内部绑定函数（获取当前第几页）
  handleCurrentChange(val: number) {
    this.meta.page = val;
    this.metas.page = val;
  }
  //返回上一层
  private routerReturn() {
    this.$router.go(-1);
  }
  /**************************************************methods结束****************************/
  // 生命周期
  private created() {
    //获取储存用户ID
    this.saveUserId();
    //客户详情基本信息 + 收货地址
    this.customerDetail();
  }
  private mounted() {
    //客户详情所属员工列表
    this.customerBelongToMembers();
  }
}
</script>
<style lang="scss" scoped>
// 公共颜色
i {
  font-style: normal;
}
hr {
  background-color: #ebeef5;
  height: 1px;
  border: none;
}
.tags-h {
  color: #292a30;
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 20px;
}
.tabs-firstrow {
  display: flex;
  .user-avatar {
    width: 70px;
    height: 70px;
    border-radius: 50%;
    margin-right: 10px;
  }
  .user-nickName {
    color: #292a30;
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 10px;
  }
  .sex {
    margin-right: 30px;
  }
  .user-info {
    span {
      color: #606266;
    }
  }
  .user-note {
    margin-top: 10px;
    margin-bottom: 20px;
    display: flex;
    .user-remark {
      position: relative;
      top: 4px;
      white-space: nowrap;
      line-height: 32px;
    }
  }
}
.user-department {
  cursor: pointer;
}
h3 {
  margin-top: 20px;
  color: #292a30;
}
.customer-label {
  width: 90%;
  margin-top: 30px;
  display: flex;
  span {
    white-space: nowrap;
    color: #606266;
    padding-top: 6px;
  }
  .customer-labelbutton {
    white-space: normal;
    font-size: 12px;
    background: rgba(0,0,0,.03);
    color: #000000;
    box-sizing: border-box;
    border: 1px solid #000000;
    span {
      color: #000000;
    }
  }
}
.table-empty-label {
  text-align: center;
  height: 340px;
  img {
    width: 124px;
    margin-top: 80px;
  }
  p {
    margin-top: -10px;
    line-height: 20px;
  }
}
.table-empty-firstlabel {
  text-align: center;
  height: 340px;
  img {
    width: 124px;
    margin-top: 80px;
  }
  p {
    margin-top: -10px;
    line-height: 20px;
  }
}
.table-empty-firsttwolabel {
  text-align: center;
  height: 340px;
  img {
    width: 124px;
    margin-top: 50px;
  }
  p {
    margin-top: -10px;
    line-height: 20px;
  }
}
.customer-Footer {
  margin: 0 auto;
  text-align: center;
  height: 340px;
  img {
    width: 124px;
    margin-top: 50px;
  }
  p {
    margin-top: -10px;
    line-height: 20px;
  }
}


/deep/.container {
  display: none;
}
.slot-content {
  padding-top: 0;
  .Refresh {
    margin-top: 0;
  }
}
.customer-footer {
  h3 {
    margin-top: 0;
  }
}
/deep/.layout-form .content {
  padding: 0 !important;
  margin-top: 20px;
}
.layout-list .table-content {
  margin-bottom: 0px !important;
}
#crs-scrm .el-timeline-item {
  padding-bottom: 20px !important;
}
.customer-footer {
  margin-top: 30px;
}
/* 为了解决表格线粗的问题 */
.el-table__fixed:before {
  content: "";
  position: absolute;
  left: 0;
  bottom: 0;
  width: 100%;
  height: 1px;
  background-color: transparent;
  z-index: 4;
}
.el-table:before {
  content: "";
  position: absolute;
  background-color: transparent;
  z-index: 1;
}
.popover-Tip {
  height: 100px;
  min-width: 300px;
  max-width: 450px;
  overflow: auto;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}
.message-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  grid-column-gap: 40px;
  margin-left: 30px;
}
.w-6 {
  font-weight: 600;
}
</style>