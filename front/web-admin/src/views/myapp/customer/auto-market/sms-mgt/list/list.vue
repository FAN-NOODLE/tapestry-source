<template>
  <div class="sms-layout">
    <layout-form type="normal">
      <template slot="content">
        <div>
          <div class="msg-number">
            <!-- <dl>
              <dt>可发送短信条数</dt>
              <dd>{{ availableNum }}</dd>
            </dl> -->
            <dl>
              <dt>已发送条数</dt>
              <dd>{{ usedNum }}</dd>
            </dl>
          </div>
          <el-tabs v-model="tabSms" class="tab-sms" @tab-click="tabSmsFn">
            <el-tab-pane label="短信模板" name="template">
              <el-button
                class="create-btn"
                type="primary"
                size="small"
                @click="createTemplateFn"
                >创建短信模板</el-button
              >
              <p class="time-show">
                数据更新时间：{{
                  updateTime | formatDate("yyyy-mm-dd hh:mm:ss")
                }}
              </p>
              <div class="search-nav">
                <search-item text="筛选：">
                  <select-input
                    ref="select"
                    v-model="curTemplateName"
                    :list="templateType"
                    @selectChange="selectChangedFn"
                    :placeholder="'请输入' + curTemplatePla"
                  ></select-input>
                </search-item>
                <!-- <search-item text="模板状态：">
                  <cy-select
                    width-size="small"
                    size="small"
                    v-model="curTemplateStatus"
                  >
                    <cy-option
                      v-for="(item, index) in templateStatus"
                      :key="index"
                      :label="item"
                      :value="index"
                    ></cy-option>
                  </cy-select>
                </search-item> -->
                <el-button
                  type="primary"
                  class="m-b-20 none-m-r"
                  @click="search"
                  size="small"
                  >查询</el-button
                >
                <el-button @click="resetFn" size="small">重置</el-button>
              </div>
              <ul class="layout-content" v-if="smsList.length">
                <li v-for="item in smsList" :key="item.id">
                  <dl class="show">
                    <dt>{{ item.template_name }}</dt>
                    <!-- <dd>ID：{{ item.template_id }}</dd> -->
                    <!-- 1：未审核、2：审核中、3：审核失败、4：审核成功、5：已赋值 -->
                    <!-- <dd>
                      状态：
                      <span :class="auditState[item.status - 1].type">
                        {{ auditState[item.status - 1].title }}
                      </span>

                      <el-popover
                        placement="top-start"
                        width="200"
                        trigger="hover"
                        :content="item.fail_reason"
                        v-if="item.status == 3"
                      >
                        <svg-icon
                          slot="reference"
                          class="help-icon"
                          icon-class="group-help"
                        ></svg-icon>
                      </el-popover>
                    </dd> -->
                    <dd v-html="item.message_str"></dd>
                    <dd class="btn-show clearfix">
                      <!-- 未审核 -->
                      <!-- <div class="pull-right" v-if="item.status == 1">
                        <el-button
                          type="primary"
                          size="mini"
                          @click="submitReviewFn(item.id)"
                          >提交审核</el-button
                        >
                      </div> -->
                      <!-- 审核通过 -->
                      <!-- <div class="pull-right" v-if="item.status == 4">
                        <el-button
                          type="primary"
                          size="mini"
                          @click="setParamFn(item.id)"
                          >完善参数</el-button
                        >
                      </div> -->
                      <!-- 已赋值(赋值通过) -->
                      <div class="pull-right" v-if="item.status == 5">
                        <!-- <el-button
                          type="primary"
                          size="mini"
                          @click="setParamFn(item.id)"
                          >修改参数</el-button
                        > -->
                      </div>
                      <a
                        class="pull-right status-blue"
                        href="javascript:;"
                        @click="editFn(item.id)"
                        >编辑</a
                      >
                      <a
                        class="pull-right status-red"
                        href="javascript:;"
                        @click="deleteMsgFn(item.id)"
                        >删除</a
                      >
                    </dd>
                    <span class="tips" v-if="item.type == 1">默认</span>
                  </dl>
                </li>
              </ul>
              <div class="no-data" v-else>
                <div v-if="searchLoading">
                  <img src="@/assets/images/empty-search.png" alt="" />
                  <p>
                    未搜索到短信模板
                  </p>
                </div>
                <div v-else>
                  <img src="@/assets/images/empty-box.png" alt="" />
                  <p>暂无短信模板</p>
                </div>
              </div>
              <div class="paging" v-if="smsList.length">
                <el-pagination
                  @size-change="sizeChangeFn"
                  @current-change="currentChangeFn"
                  :current-page="paging.page"
                  :page-sizes="[9, 18, 27, 36]"
                  :page-size="paging.size"
                  layout="total, sizes, prev, pager, next, jumper"
                  :total="paging.total"
                >
                </el-pagination>
              </div>
            </el-tab-pane>
            <el-tab-pane label="短信签名" name="sign">
              <sms-sign :isTabApi="isTabApi" v-if="isTabApi"></sms-sign>
            </el-tab-pane>
          </el-tabs>
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import SelectInput from "@/components/business/select-input/index.vue";
import API from "@/api/customer";
import smsSign from "../../components/sms-sign.vue";
@Component({
  components: { SelectInput, smsSign }
})
export default class extends Vue {
  /*************************** data start ***************************/
  private updateTime = 0; //更新时间
  private smsList = []; //短信列表
  private availableNum = 0; //可发送短信条数
  private usedNum = 0; //已发送条数
  private curTemplateName = ""; //选择的模板的name
  private curTemplateType = 2; //选择模板的key
  private curTemplatePla = "模板名称"; //选择模板的key
  private tabSms = "template";
  //模板类型
  private templateType = [
    { key: 2, name: "模板名称" },
    { key: 1, name: "模板ID" },
    { key: 3, name: "短信内容" }
  ];
  //模板状态
  private templateStatus = [
    "全部",
    "未审核",
    "审核中",
    "审核失败",
    "审核通过",
    "已赋值"
  ];
  private curTemplateStatus = 0; //当前模板状态
  //分页信息
  private paging = {
    page: 1,
    size: 9,
    total: 0
  };
  private auditState = [
    { title: "未审核", type: "status-blue" },
    { title: "审核中", type: "status-orange" },
    { title: "审核失败", type: "status-red" },
    { title: "审核通过", type: "status-green" },
    { title: "已赋值(审核通过)", type: "status-green" }
  ];
  private isTabApi = false; //是否切换tab请求接口
  private searchLoading = false; //控制搜索缺省图
  /*************************** method start ***************************/
  //创建短信模板
  private createTemplateFn() {
    this.$router.push({
      path: "/myapp/customer/auto-market/sms-mgt/create",
      query: {
        type: "1"
      }
    });
  }
  //企业短信数据统计
  private smsTemplateTotalFn() {
    API.smsTemplateTotalApi({}).then((res: any) => {
      if (res.code == 200) {
        const data = res.data;
        this.availableNum = data.available_num; //可发送短信条数
        this.usedNum = data.used_num; //已发送条数
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //模板
  private selectChangedFn(val: any) {
    this.curTemplateType = val;
    this.templateType.forEach(item => {
      if (item.key == val) {
        this.curTemplatePla = item.name;
      }
    });
  }
  //查询
  private search() {
    this.paging.page = 1;
    this.smsTemplateListFn();
  }
  //短信模板列表
  private smsTemplateListFn() {
    const params = {
      page: this.paging.page, //页码
      size: this.paging.size, //每页数量
      search_type: this.curTemplateType, //搜索类型(1:模板ID、2:模板名称、3:模板内容)
      search_text: this.curTemplateName, //搜索内容文本
      status: this.curTemplateStatus //模板状态（0：全部、1：未审核、2：审核中、3：审核失败、4：审核成功、5：已赋值）
    };
    API.smsTemplateListApi(params).then((res: any) => {
      if (
        this.curTemplateStatus ||
        this.curTemplateName ||
        this.curTemplateType
      ) {
        this.searchLoading = true;
      } else {
        this.searchLoading = false;
      }

      if (res.code == 200) {
        const data = res.data;
        this.smsList = data.list;
        this.updateTime = data.updateTime;
        this.paging.total = data.total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 查询
  private onsearch() {
    this.paging.page = 1;
    this.smsTemplateListFn();
  }
  //重置
  private resetFn() {
    (this.$refs["select"] as any).reset();
    this.curTemplatePla = "模板名称";
    this.curTemplateType = 2;
    this.curTemplateName = "";
    this.curTemplateStatus = 0;
    this.paging.page = 1;
    this.smsTemplateListFn();
  }
  //改变size
  private sizeChangeFn(val: number) {
    this.paging.page = 1;
    this.paging.size = val;
    this.smsTemplateListFn();
  }
  //改变当前页
  private currentChangeFn(val: number) {
    this.paging.page = val;
    this.smsTemplateListFn();
  }
  //删除
  private deleteMsgFn(id: number) {
    this.$confirm("是否需要删除该短信模板", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        const params = {
          id: id
        };
        //删除子标签接口
        API.smsTemplateDeleteApi(params).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              message: "删除成功",
              type: "success",
              duration: 2000
            });
            this.smsTemplateListFn();
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
  //提交审核
  private submitReviewFn(id: number) {
    API.smsTemplateArraignmentApi({ id: id }).then((res: any) => {
      if (res.code === 200) {
        //列表接口
        this.smsTemplateListFn();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //编辑
  private editFn(id: any) {
    this.$router.push({
      path: "/myapp/customer/auto-market/sms-mgt/create",
      query: {
        type: "2",
        id: id
      }
    });
  }
  //设置参数、编辑参数
  private setParamFn(id: any) {
    this.$router.push({
      path: "/myapp/customer/auto-market/sms-mgt/set-param",
      query: {
        id: id
      }
    });
  }
  //切换tab
  private tabSmsFn() {
    if (this.tabSms == "template") {
      this.paging.page = 1;
      this.isTabApi = false;
      this.smsTemplateListFn();
    } else {
      this.isTabApi = true;
    }
  }

  /*************************** 生命周期 start ***************************/
  private created() {
    // //企业短信数据统计
    this.smsTemplateTotalFn();
    //短信列表
    this.smsTemplateListFn();
  }
}
</script>
<style lang="scss" scoped>
li {
  list-style: none;
}
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
a {
  text-decoration: none;
}
.none-m-r {
  margin-right: 0 !important;
}
.help-icon {
  cursor: pointer;
}
//下拉输入框行显示
.select-input {
  display: inline-block;
}
.msg-number {
  overflow: hidden;
  dl {
    width: 216px;
    margin-right: 30px;
    text-align: center;
    color: $color-text-main;
    padding: 30px;
    background: $color-background;
    float: left;
    font-weight: bold;
    dt {
      font-size: 14px;
      line-height: 20px;
      margin-bottom: 20px;
    }
    dd {
      font-size: 26px;
      line-height: 38px;
    }
  }
}
.tab-sms {
  margin-top: 10px;
  .no-data {
    text-align: center;
    margin: 20px 0;
    img {
      margin-top: 30px;
    }
    p {
      color: $color-normal;
      margin-bottom: 40px;
    }
  }
}
.create-btn {
  margin-bottom: 20px;
}
.time-show {
  color: $color-text-normal;
  font-size: 14px;
  line-height: 20px;
}
.search-nav {
  background: $color-background;
  padding-top: 20px;
  padding-left: 20px;
  margin: 20px 0;
}
ul.layout-content {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  margin-right: -30px;
  li {
    width: 33.333%;
    box-sizing: border-box;
    padding-right: 30px;
    margin-bottom: 20px;
  }
}
dl.show {
  background: $color-background;
  height: 100%;
  padding: 0 30px;
  border-radius: 5px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
  dt {
    color: $color-text-main;
    font-weight: bold;
    font-size: 14px;
    line-height: 20px;
    margin: 30px 0 20px;
  }
  dd {
    font-size: 14px;
  }
  dd:nth-child(2) {
    color: $color-normal;
    line-height: 20px;
  }
  dd:nth-child(3) {
    margin-top: 10px;
    line-height: 20px;
  }
  dd:nth-child(4) {
    color: $color-text-normal;
    margin: 10px 0 38px;
    min-height: 40px;
    //两行省略号显示
    text-overflow: -o-ellipsis-lastline;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }
  .btn-show {
    margin-bottom: 30px;
    div,
    a {
      line-height: 32px;
      font-size: 14px;
    }
    a {
      margin-right: 20px;
    }
  }
  .status-green {
    color: #48ba9b;
  }
  .status-red {
    color: #f56c6c;
  }
  .status-blue {
    color: $color-text-button;
  }
  .status-orange {
    color: #f6ab3c;
  }
  span.tips {
    position: absolute;
    top: 10px;
    right: 10px;
    z-index: 1;
    color: $color-main;
    font-size: 12px;
    line-height: 18px;
  }
}
.paging {
  text-align: right;
}
/deep/.el-button--mini,
.el-button--small {
  font-size: 14px;
}
</style>
<style lang="scss">
.layout-content {
  input[type="button"] {
    border: none;
    background: none;
    outline: none;
    margin: 0 1px;
    padding: 0;
    color: $color-text-normal;
  }
}
.sms-layout {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  .beyound-container {
    margin-bottom: 0px !important;
  }
}
</style>
