<template>
  <layout-list v-if="isShow">
    <template slot="top-tabs">
      <el-radio-group class="tab-tit" v-model="first_active_name">
        <el-radio-button :label="1">基本信息</el-radio-button>
        <el-radio-button :label="2">活动内容</el-radio-button>
        <el-radio-button :label="3">统计与详情</el-radio-button>
      </el-radio-group>
    </template>
    <template slot="top-tabs">
      <div v-if="first_active_name === 1">
        <div class="content-info">
          <div class="info-title">基础信息</div>
          <div class="info-content">
            <p class="info-content-text">
              <span>活动类型</span>
              {{ base_info.name }}
            </p>
            <p class="info-content-text">
              <span>活动标题</span>
              {{ base_info.title }}
            </p>
            <p class="info-content-text">
              <span>活动规则</span>
              <span style="width:85%;color:#000000;">{{
                base_info.activity_rule
              }}</span>
            </p>
            <p class="info-content-text">
              <span>封面图片</span
              ><img
                :src="base_info.cover_url"
                alt=""
                @click="checkNowImage(base_info.cover_url)"
              />
            </p>
            <p class="info-content-text">
              <span>活动有效期</span>{{ base_info.period }}
            </p>
            <!-- 签到活动 -->
            <p class="info-content-text" v-if="base_info.sign_period">
              <span>签到周期</span>{{ base_info.sign_period }}
            </p>
            <p class="info-content-text" v-if="base_info.sign_type">
              <span>签到规则</span>
              {{
                base_info.sign_type === 1
                  ? "满签到"
                  : base_info.sign_type === 2
                  ? `连续签到${base_info.sign_day}天`
                  : `累计签到${base_info.sign_day}天`
              }}
            </p>
            <!-- {{ taskObj.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }}～{{
              taskObj.end_time | formatDate("yyyy-mm-dd hh:mm:ss")
            }} -->
          </div>
        </div>
        <div class="content-info">
          <div class="info-title">活动范围</div>
          <div class="info-content">
            <div class="activity-scope">
              <span>活动范围</span>
              <div class="activity-scope-content">
                <p>{{ scope.name }}</p>
                <el-tag v-for="(item, index) in scope.ba_list" :key="index"
                  ><i class="el-icon-user-solid" style="margin-right:10px;"></i
                  >{{ item.name }}</el-tag
                >
              </div>
            </div>
          </div>
        </div>
        <div class="content-info">
          <div class="info-title">其他</div>
          <div class="info-content">
            <!-- 报名活动 -->
            <p class="info-content-text" v-if="other.is_audit">
              <span>报名审核</span
              >{{ other.is_audit === 1 ? "需要" : "不要要" }}
            </p>
            <!-- 签到活动 -->
            <p class="info-content-text" v-if="base_info.type === 2">
              <span>是否配置签到内容</span>
              {{ other.is_config_sign_content_text }}
            </p>
            <p class="info-content-text" v-if="base_info.type === 2">
              <span>是否开启留资</span>
              {{ other.is_open_remark_text }}
            </p>
            <p class="info-content-text">
              <span>弹框说明</span>
              {{ other.explanation }}
            </p>
            <!-- 报名活动 -->
            <p class="info-content-text" v-if="other.is_audit === 1">
              <span>审核通过弹窗</span
              ><img
                :src="other.audit_success_url"
                alt=""
                @click="checkNowImage(other.audit_success_url)"
              />
            </p>
            <p class="info-content-text" v-if="other.is_audit === 1">
              <span>审核未通过弹窗</span
              ><img
                :src="other.audit_fail_url"
                alt=""
                @click="checkNowImage(other.audit_fail_url)"
              />
            </p>
            <!-- ------- -->
          </div>
        </div>
      </div>
      <div v-if="first_active_name === 2">
        <div class="activity-content-title">
          活动标题：{{ base_info.title }}
        </div>
        <div class="activity-content">
          <el-tabs v-model="second_active_name" @tab-click="changeTab">
            <el-tab-pane label="活动内容" name="1"></el-tab-pane>
            <el-tab-pane
              label="留资内容"
              name="2"
              v-if="base_info.type === 2"
            ></el-tab-pane>
          </el-tabs>
          <div class="content" v-if="templateList.length > 0">
            <div
              v-for="(v, i) in templateList"
              :key="i"
              class="template-content"
            >
              <div class="template-left">
                <span>【{{ v.name }}】</span>
              </div>
              <div class="template-medium">
                <div class="title">
                  {{ ++i }}、 {{ v.tit }}
                  <span class="ml-20">
                    <el-checkbox v-model="v.isRequired" disabled
                      >必填</el-checkbox
                    >
                  </span>
                </div>
                <div class="msc" v-if="v.name">{{ v.mark }}</div>
                <!-- 单选多选 -->
                <div
                  v-if="v.type == 1 || v.type == 3 || v.type == 9"
                  class="item"
                >
                  <el-checkbox-group class="checkbox-group">
                    <div v-if="v.options && v.options.length > 0">
                      <div v-for="(item, index) in v.options" :key="index">
                        <el-checkbox :label="item.option" disabled>
                        </el-checkbox>
                        <div v-if="item.pic && item.pic.length > 0">
                          <el-image fit="cover"
                            v-for="(img, pIndex) in item.pic"
                            :key="pIndex"
                            style="width: 100px; height: 100px; margin: 2px 10px"
                            :src="img"
                          >
                          </el-image>
                        </div>
                      </div>
                      <div v-if="v.isOther">
                        <el-checkbox :label="'其他'" disabled> </el-checkbox>
                      </div>
                    </div>
                  </el-checkbox-group>
                </div>
                <!-- 输入选择类 -->
                <div
                  v-if="[2, 4, 6, 7, 8, 10, 11, 12, 13].indexOf(v.type) > -1"
                >
                  <div class="input-select">请输入内容</div>
                </div>
                <!-- 图片上传 -->
                <div class="upload" v-if="v.type == 5">
                  <el-image fit="cover">
                    <div slot="error" class="image-slot">
                      <i class="el-icon-plus"></i>
                    </div>
                  </el-image>
                </div>
                <div class="msc" v-if="v.type == 5">
                  最多上传{{ v.max }}张图片
                </div>
                <div class="msc" v-if="[2, 4].indexOf(v.type) > -1">
                  <span v-if="v.min">最少输入{{ v.min }}个字符</span>&nbsp;
                  <span v-if="v.max">最多输入{{ v.max }}个字符</span>
                </div>
                <div class="msc" v-if="v.type == 3">
                  <span v-if="v.min">最少选择{{ v.min }}个项</span>&nbsp;
                  <span v-if="v.max">最多选择{{ v.max }}个项</span>
                </div>
              </div>
            </div>
          </div>
          <div
            v-if="second_active_name === '1' && templateList.length === 0"
            class="empty"
          >
            暂无配置活动内容
          </div>
          <div
            v-if="second_active_name === '2' && templateList.length === 0"
            class="empty"
          >
            暂无配置留资内容
          </div>
        </div>
      </div>
      <statistics-details
        v-if="first_active_name === 3"
        :isGoBack="isGoBack"
        :activity_info="activity_info"
        :statistical="statistical"
        :dataList="dataList"
        :total="total"
        @getHeaderDetail="getHeaderDetail(arguments)"
        @getList="getList(arguments)"
      ></statistics-details>
      <!-- 图片视频弹框 -->
      <preview
        v-if="preview_show"
        :urls="checkImgArr"
        @close="preview_show = false"
      ></preview>
    </template>
  </layout-list>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
// 选择员工
import statisticsDetails from "../activity-detail/components/statistics-details.vue";
// 图片预览
import preview from "@/components/business/preview/index.vue";
import Api from "@/api/wisdom/marketing-activities";
@Component({
  components: {
    statisticsDetails,
    preview,
  },
})
export default class extends Vue {
  private activityId = 0;
  private isShow = false;
  private first_active_name = 1;
  private second_active_name = "1";
  private base_info: any = {};
  private other: any = {};
  private scope: any = {};
  private templateList = [];
  private templateList1 = [];
  private templateList2 = [];
  private activity_info = {};
  private statistical = {};
  private dataList: Array<any> = [];
  private total = 0;
  private isGoBack = false;
  private preview_show = false;
  private checkImgArr: Array<any> = [];
  // 切换tab
  changeTab(value: any) {
    if (value.name == "1") {
      this.templateList = this.templateList1;
    } else {
      this.templateList = this.templateList2;
    }
  }
  // 获取基本信息
  private activityBaseInfo() {
    Api.activityBaseInfo({ id: this.activityId })
      .then((res: any) => {
        if (res.httpCode === 200) {
          this.base_info = res.data.base_info;
          this.other = res.data.other;
          this.scope = res.data.scope;
        }
      })
      .catch((err: any) => {
        console.log("activityBaseInfo", err);
      });
  }

  // 获取活动内容
  private getActivityContent() {
    Api.activityContent({ id: this.activityId })
      .then((res: any) => {
        if (res.httpCode === 200) {
          try {
            const data = JSON.parse(res.data);
            this.templateList1 = data.templateList || [];
            this.templateList2 = data.templateSecondList || [];
            this.templateList = this.templateList1;
          } catch (error) {
            this.templateList = [];
          }
        }
      })
      .catch((err: any) => {
        console.log("activityBaseInfo", err);
      });
  }

  // 获取统计与详情
  private getActivityStatistics(type = 1, parameter?: any) {
    const params: any = {
      id: this.activityId,
      type: type,
      ...parameter,
    };
    Api.activityStatistics(params)
      .then((res: any) => {
        if (res.httpCode === 200) {
          if (type === 1) {
            this.activity_info = res.data.activity_info;
            this.statistical = res.data.statistical;
            this.dataList = res.data.data_list.list;
            this.total = res.data.data_list.meta.total;
          } else if (type === 2) {
            this.statistical = res.data.statistical;
          } else if (type === 3) {
            this.dataList = res.data.data_list.list;
            this.total = res.data.data_list.meta.total;
          }
          if (this.isGoBack) {
            this.first_active_name = 3;
          }
          this.isShow = true;
        }
      })
      .catch((err: any) => {
        console.log("activityBaseInfo", err);
      });
  }
  // 筛选头部信息
  private getHeaderDetail(msg: any) {
    const { type, ...params } = msg[1];
    this.getActivityStatistics(2, params);
  }

  // 筛选列表
  private getList(msg: any) {
    const { type, ...params } = msg[1];
    this.getActivityStatistics(3, params);
  }

  private checkNowImage(img: string) {
    this.checkImgArr = [img];
    this.preview_show = true;
  }

  // 生命周期
  private mounted() {
    this.activityId = Number(this.$route.query.id);
    this.activityBaseInfo();
    this.getActivityContent();
    this.getActivityStatistics();
  }

  private beforeRouteEnter(to: any, form: any, next: (vm: any) => void): void {
    next((vm: any) => {
      if (form.name === "myapp_wisdom_marketing-activities_user-detail") {
        vm.isGoBack = true;
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.content-info {
  padding: 10px 0;
  .info-title {
    width: 100%;
    background: #f2f2f2;
    padding: 4px 10px;
  }
  .info-content {
    width: 100%;
    padding: 10px 18px;
    .info-content-text {
      display: flex;
      align-items: center;
      line-height: 40px;
      color: #333333;
      margin: 7px;
      span {
        display: inline-block;
        width: 130px;
        color: #7f7f7f;
      }
      img {
        width: 250px;
      }
    }
    .activity-scope {
      width: 100%;
      display: flex;
      span {
        display: inline-block;
        width: 130px;
        color: #7f7f7f;
      }
      .activity-scope-content {
        width: 100%;
        p {
          width: 95%;
          padding: 4px 10px;
          background: #ebebf0;
          margin-bottom: 15px;
        }
      }
    }
  }
}
.activity-content-title {
  font-size: 16px;
  font-weight: 900;
  line-height: 48px;
}
.tab-tit {
  margin-bottom: 20px;
}
.template-content {
  display: flex;
  margin-bottom: 15px;
  .template-left {
    width: 160px;
    color: #888 !important;
    font-size: 16px;
    padding-right: 5px;
    line-height: 36px;
  }
  .template-medium {
    flex: 1;
    padding-top: 6px;
    .title {
      font-size: 14px;
      font-weight: 600;
    }
    .item {
      .checkbox-group {
        display: flex;
        flex-direction: column;
      }
    }
  }
}
.empty {
  height: 100%;
  line-height: 404px;
  font-size: 20px;
  text-align: center;
  color: #888;
}
.el-tag {
  margin-right: 20px;
}
.msc {
  font-size: 12px;
  color: #c8c8c8;
  margin: 2px 0;
}
.ml-20 {
  margin-left: 20px;
}
.input-select {
  width: 400px;
  height: 50px;
  line-height: 50px;
  padding: 0 20px;
  background: #f8f9fc;
}
.upload {
  width: 140px;
  height: 140px;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 1px solid #9fa2ab;
}
</style>
