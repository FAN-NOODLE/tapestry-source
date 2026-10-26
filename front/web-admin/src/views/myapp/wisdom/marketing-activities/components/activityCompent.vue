<template>
  <div class="active-container">
    <layout-form>
      <template slot="content">
        <div class="header">
          <div :class="['item', active == 1 && 'active-item']">
            1、基础信息
          </div>
          <div :class="['item', active == 2 && 'active-item']">
            2、活动内容设计
          </div>
          <div :class="['item', active == 3 && 'active-item']">
            3、页面配置
          </div>
        </div>
      </template>
      <template slot="content">
        <basicInfo
          v-show="active == 1"
          ref="basicInfo"
          :activityId="activityId"
          :newActivityId="newActivityId"
          :base_info="base_info"
          @submitSuccess="submitSuccess"
        ></basicInfo>
        <activityContent
          v-show="active == 2"
          ref="activityContent"
          :activityId="activityId"
          :newActivityId="newActivityId"
          :isActivityContent="isActivityContent"
          :isMaterial="isMaterial"
          :activity_type="activity_type"
          :content="content"
          :retainCapital="retainCapital"
          @submitSuccess="submitSuccess"
        ></activityContent>
        <pageConfig
          v-show="active == 3"
          ref="pageConfig"
          :activityId="activityId"
          :newActivityId="newActivityId"
          :config="config"
          :content="content"
          :base_info_copy="base_info_copy"
          @submitSuccess="submitSuccess"
        ></pageConfig>
      </template>
      <template slot="footer">
        <div>
          <el-button size="small" v-if="active != 1" @click="previous"
            >上一步</el-button
          >
          <el-button
            type="primary"
            size="small"
            v-if="active == 3"
            @click="toSave"
            >确认保存</el-button
          >
          <el-button
            type="primary"
            size="small"
            v-if="active != 3"
            @click="nextStep"
            >下一步</el-button
          >
          <el-button size="small" @click="cancel">取 消</el-button>
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import basicInfo from "./basicInfo.vue";
import activityContent from "./activityContent.vue";
import pageConfig from "./pageConfig.vue";
import Api from "@/api/wisdom/marketing-activities";
import Bus from "@/components/bus";
@Component({
  components: {
    basicInfo,
    activityContent,
    pageConfig,
  },
})
export default class extends Vue {
  @Prop({ default: () => 1 }) statusType?: number; //类型: 1-新增 , 2-编辑 , 3-复制

  private active = 1; // 当前
  private activityId = 0;
  private newActivityId = 0;
  private isActivityContent = false; // 当前活动是否配置签到内容
  private isMaterial = false; // 当前活动书否留资
  private activity_type = 1; // 当前活动的类型 1、报名；2、签到；3、问卷
  private base_info: any = {};
  private base_info_copy: any = {};
  private content: any = "";
  private retainCapital: any = "";
  private config: any = {};

  // ----------------------------methods
  // 上一步
  private previous() {
    if (
      !this.isMaterial &&
      !this.isActivityContent &&
      this.activity_type === 2
    ) {
      this.active = 1;
    } else {
      this.active--;
    }
  }
  // 下一步
  private nextStep() {
    if (this.active === 1) {
      (this.$refs["basicInfo"] as any).submit();
    } else if (this.active === 2) {
      (this.$refs["activityContent"] as any).submit();
    }
  }
  private submitSuccess(id: number) {
    if (id !== 0) {
      if (this.activityId === 0) {
        this.activityId = id;
      } else if (this.activityId !== id) {
        this.newActivityId = id;
      }
    }
    if (this.active === 1) {
      const base_info = (this.$refs["basicInfo"] as any).activityConfigFormData;
      this.base_info_copy = {
        activity_type: base_info.active_type,
        title: base_info.active_name,
        start_time: base_info.indate[0] / 1000,
        end_time: base_info.indate[1] / 1000,
        sign_start:
          base_info.active_type === 2 ? base_info.sign_in_period[0] / 1000 : 0,
        sign_end:
          base_info.active_type === 2 ? base_info.sign_in_period[1] / 1000 : 0,
        cover_url: base_info.banner_pic,
        activity_rule: base_info.active_rule,
      };
      this.activity_type = base_info.active_type;
      if (this.activity_type === 2) {
        this.isMaterial = base_info.retain_capital === 1 ? true : false;
        this.isActivityContent = base_info.sign_in_content === 1 ? true : false;
      }
    }
    if (this.active === 2) {
      const content = {
        templateList: (this.$refs["activityContent"] as any).templateList,
      };
      this.content = JSON.stringify(content);
    }
    if (this.active === 3) {
      this.$router.push({
        name: "myapp_wisdom_marketing-activities_list",
      });
    } else {
      // 判断到签到活动，且不需要配置内容和留资
      if (
        !this.isMaterial &&
        !this.isActivityContent &&
        this.activity_type === 2
      ) {
        this.active = 3;
      } else {
        this.active += 1;
      }
    }
  }

  // 获取活动基本信息详情
  getActivityBaseDetail(steps = 1) {
    Api.activitySteps({ id: this.activityId, steps: steps })
      .then((res: any) => {
        if (res.httpCode === 200) {
          if (steps === 1) {
            // 基本信息
            this.base_info = {
              ...res.data.base_info,
              ...res.data.scope,
              ...res.data.other,
            };
            this.base_info_copy = JSON.parse(JSON.stringify(this.base_info));
            // 获取下是否留资
            if (this.base_info.is_open_remark === 1) {
              this.getActivityBaseDetail(4);
            }
            if (this.base_info.is_config_sign_content !== 0) {
              this.getActivityBaseDetail(2);
            }
          } else if (steps === 2) {
            // 活动内容
            this.content = res.data.content;
          } else if (steps === 3) {
            // 页面配置
            this.config = res.data;
          } else if (steps === 4) {
            // 留资
            this.retainCapital = res.data.remark;
          }
        }
      })
      .catch((err: any) => {
        console.log("activityBaseInfo", err);
      });
  }
  // 取消
  private cancel() {
    this.$router.go(-1);
  }
  // 确认保存
  private toSave() {
    (this.$refs["pageConfig"] as any).submit();
  }

  // 生命周期
  private mounted() {
    if (this.statusType !== 1) {
      this.activityId = Number(this.$route.query.id);
      this.getActivityBaseDetail();
      this.getActivityBaseDetail(3);
    }
  }
}
</script>
<style lang="scss" scoped>
.active-container {
  .header {
    display: flex;
    .item {
      padding: 10px 50px;
      background-color: #eeeeee;
      color: #7f7f7f;
      position: relative;
      user-select: none;
      &::after {
        content: " ";
        position: absolute;
        width: 14px;
        height: calc(100% + 2px);
        background-color: #fff;
        top: 10px;
        right: -6px;
        transform: rotate(45deg);
        z-index: 1;
      }
      &::before {
        content: " ";
        position: absolute;
        width: 14px;
        height: calc(100% + 2px);
        background-color: #fff;
        top: -10px;
        right: -6px;
        transform: rotate(-45deg);
        z-index: 1;
      }
    }
    .active-item {
      background-color: #69a3ff;
      color: #fff;
    }
  }
}
</style>
