<template>
  <div class="layout-list">
    <p class="time-show">
      数据更新时间：{{ currentDate | formatDate("yyyy-mm-dd hh:mm:ss") }}
      <span @click="refreshFn">刷新</span>
    </p>
    <dl class="search-nav-self">
      <dd>筛选：</dd>
      <dd>
        <select-input
          ref="select"
          v-model="labelName"
          :list="list"
          @selectChange="selectChangedFn"
          :placeholder="labelId == 1 ? '请输入标签组名称' : '请输入标签名称'"
        ></select-input>
      </dd>
      <dd>
        <el-button type="primary" size="small" @click="getLabelListFn"
          >查询</el-button
        >
      </dd>
      <dd>
        <el-button size="small" @click="setResetFn">重置</el-button>
      </dd>
    </dl>
    <div class="label-content">
      <el-tabs v-model="activeTab" @tab-click="tabClickFn">
        <el-tab-pane label="客户标签" name="manual">
          <!--预留客户标签组件-->
          <manual-label
            :group="manualArr"
            :keywords="labelName"
            :labelId="labelId"
            :activeTab="activeTab"
            :authStatus="auth_status"
            v-if="
              activeTab == 'manual' &&
                (!hasValueSearch || (hasValueSearch && manualArr.length > 0))
            "
            @tagSonCase="tagSonCaseFn"
          ></manual-label>
          <div class="no-data" v-if="!manualArr.length">
            <div v-if="hasValueSearch">
              <img src="https://cdn-cos.coach.com.cn/uploads/2023/3.5.1/empty-search.png" alt="" />
              <p v-if="labelId == 1">未搜索到标签组</p>
              <p v-else>未搜索到标签</p>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="印象标签" name="effect">
          <!--预留客户标签组件-->
          <manual-label
            :group="manualArr"
            :keywords="labelName"
            :labelId="labelId"
            :auth_status="false"
            :activeTab="activeTab"
            v-if="
              activeTab == 'effect' &&
                (!hasValueSearch || (hasValueSearch && manualArr.length > 0))
            "
            @tagSonCase="tagSonCaseFn"
          ></manual-label>
          <div class="no-data" v-if="!manualArr.length">
            <div v-if="hasValueSearch">
              <img src="https://cdn-cos.coach.com.cn/uploads/2023/3.5.1/empty-search.png" alt="" />
              <p v-if="labelId == 1">未搜索到标签组</p>
              <p v-else>未搜索到标签</p>
            </div>
          </div>
        </el-tab-pane>
        <!-- 暂不使用 -->
        <el-tab-pane label="外部标签同步任务" name="outside" v-if="false">
          <outside-label
            :outsideLabelList="manualArr"
            v-if="activeTab == 'outside'"
          />
          <div class="no-data" v-if="!manualArr.length">
            <div v-if="hasValueSearch">
              <img src="https://cdn-cos.coach.com.cn/uploads/2023/3.5.1/empty-search.png" alt="" />
              <p v-if="labelId == 1">未搜索到标签组</p>
              <p v-else>未搜索到标签</p>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import ManualLabel from "../components/manual-label.vue";
import SelectInput from "@/components/business/select-input/index.vue";
import OutsideLabel from "../components/outside-label.vue";
import API from "@/api/customer";
@Component({
  components: {
    ManualLabel,
    SelectInput,
    OutsideLabel,
  }
})
export default class extends Vue {
  /*************************** data start ***************************/
  //标签下拉数组
  private list = [
    { key: 1, name: "标签组名称" },
    { key: 2, name: "标签名称" }
  ];
  private auth_status = false;
  private labelName = ""; //选择的标签/标签组名称的key
  private labelId = 1; //选择的标签/标签组
  private activeTab = "manual"; //标签选项卡 automatic==自动标签、manual==客户标签，默认automatic
  private manualArr: any = []; //客户标签
  private tagSonCase: any = ""; //子组件的实例
  private currentDate = 0; //当前更新的时间
  private hasValueSearch = false; // 用来判断是否是带值搜索
  
  private searchLoading = false; //缺省页面判断
  private showAddLable = false; //新增标签弹框

  /*************************** method start ***************************/
  //获取标签列表
  private async getLabelListFn() {
    const params = {
      search_type: this.labelId, //搜索类型（1：标签组名称，2：标签名称）
      keyword: this.labelName.trim(),
      ...(this.activeTab == "effect" && {
        page: 1,
        size: 1000
      })
    };
    const api =
      this.activeTab === "effect" ? "impressionGroupListApi" : "labelListApi";
    try {
      const res = await API[api](params);
      if (this.labelId || this.labelName.trim()) {
        this.searchLoading = true;
      } else {
        this.searchLoading = false;
      }
      this.hasValueSearch = params.keyword.length > 0; // 是否是带值搜索
      if (this.activeTab === "effect") {
        this.manualArr = res.data.list.map(item => {
          item.tag_group_name = item.group_name;
          item.tag_group_id = item.id;
          return item;
        });
        return;
      }
      const list = res.data;
      this.manualArr = list.manual || [];
    } catch (error) {
      console.error(error);
    }
  }
  //下拉
  private selectChangedFn(val: any) {
    console.log(val);
    
    this.labelId = val;
  }
  //切换自动/手动选项卡
  private tabClickFn() {
    // todo: 重置
    this.setResetFn()
    this.manualArr = [];
  }
  //重置
  private setResetFn() {
    this.labelName = "";
    this.labelId = 1;
    (this.$refs["select"] as any).reset();
    //重置调用标签列表接口
    this.getLabelListFn();
  }
  //刷新 父组件调用子组件方法
  private refreshFn() {
    this.updata();
    this.tagSonCase.getSubListFn();
  }
  private updata() {
    const date = new Date().getTime();
    //返回格式化后的日期
    this.currentDate = date / 1000;
  }
  tagSonCaseFn(obj: any) {
    this.tagSonCase = obj._this;
    // this.currentDate = obj.data;
  }
  private getAuthStatus() {
    API.getAuthApi().then((res: any) => {
      this.auth_status = res.data.auth;
    });
  }

  /*************************** 生命周期 start ***************************/
  private created() {
    //获取标签列表
    this.getLabelListFn();
    this.updata();
    this.getAuthStatus();
    const model: any = this.$route.query.model || "";
    if (model) {
      this.activeTab = model;
    }
  }
}
</script>
<style lang="scss" scoped>
dl,
dd,
dt {
  padding: 0;
  margin: 0;
}
.layout-list {
  background-color: $color-white;
  padding: 30px;
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  // border-radius: 10px;
  .time-show {
    color: $color-text-normal;
    font-size: 14px;
    line-height: 20px;
    padding-bottom: 20px;
    span {
      font-size: 14px;
      color: $color-main;
      line-height: 20px;
      margin-left: 10px;
      cursor: pointer;
    }
  }
  .search-nav-self {
    padding: 20px;
    background: $color-background;
    color: $color-text-main;
    font-size: 14px;
    display: flex;
    align-items: center;
    margin-bottom: 10px;
    dd {
      margin: 0;
      margin-right: 20px;
      ::v-deep.el-button--small {
        font-size: 14px;
      }
    }
    dd:first-child {
      margin-right: 0;
    }
  }
  .label-content {
    .export-btn {
      margin-bottom: 20px;
    }
    .no-data {
      height: 320px;
      text-align: center;
      margin-top: 80px;
      p {
        color: $color-normal;
        margin-top: 14px;
      }
    }
  }
}
/deep/ {
  .select-input {
    width: 340px !important;
  }
  .select-input .select {
    width: 150px !important;
  }
}
// /deep/ .el-tabs__active-bar {
//   width: 60px !important;
// }
</style>
