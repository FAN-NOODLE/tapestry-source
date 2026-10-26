<template>
  <div>
    <layout-form>
      <template slot="content">
        <cy-form
          ref="cform"
          :rules="activityConfigFormRules"
          :model="activityConfigFormData"
          label-width="130px"
        >
          <div class="header-title">基础信息</div>
          <cy-form-item
            label="活动类型"
            prop="active_type"
            required
          >
            <el-radio-group
              v-model="activityConfigFormData.active_type"
              :disabled="statusType !== 1 || activityId !== 0"
            >
              <el-radio :label="1">报名活动</el-radio>
              <el-radio :label="3">问卷活动</el-radio>
              <el-radio :label="2">签到活动</el-radio>
            </el-radio-group>
          </cy-form-item>
          <cy-form-item
            label="活动标题"
            prop="active_name"
            required
          >
            <cy-input
              :trim="true"
              placeholder="请输入活动标题"
              :maxlength="20"
              :showWordLimit="true"
              v-model="activityConfigFormData.active_name"
              widthSize="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="活动规则"
            prop="active_rule"
            required
          >
            <el-input
              type="textarea"
              :rows="6"
              :maxlength="200"
              resize="none"
              show-word-limit
              placeholder="请输入活动规则"
              v-model="activityConfigFormData.active_rule"
            >
            </el-input>
            <div class="tip">活动规则说明文本，规则内容将在前端展示。</div>
          </cy-form-item>
          <cy-form-item
            label="封面图片"
            required
            style="width: 600px"
          >
            <upload-image
              :urls="activityConfigFormData.banner_pic"
              accept="image/jpeg,image/png"
              :before-upload="beforeUpload"
              @change="bannerPicSuccess"
            ></upload-image>
            <div class="tips">
              建议上传尺寸为630*316px的图片，大小不超过5000k, 建议上传 &lt; 200k
            </div>
          </cy-form-item>
          <cy-form-item
            label="活动有效期"
            prop="indate"
            required
          >
            <el-date-picker
              v-model="activityConfigFormData.indate"
              type="datetimerange"
              :clearable="false"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              :default-time="['00:00:00', '23:59:59']"
              @change="selectIndate"
              :picker-options="pickerOptionsStart"
            >
            </el-date-picker>
          </cy-form-item>
          <div v-if="activityConfigFormData.active_type == 2">
            <cy-form-item
              label="签到周期"
              prop="sign_in_period"
              required
            >
              <el-date-picker
                v-model="activityConfigFormData.sign_in_period"
                type="datetimerange"
                :clearable="false"
                range-separator="-"
                :default-time="['00:00:00', '23:59:59']"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                @change="selectSignIn"
                :picker-options="pickerOptionsSignIn"
              >
              </el-date-picker>
            </cy-form-item>
            <cy-form-item
              label="签到规则"
              required
            >
              <el-radio-group v-model="activityConfigFormData.signinType">
                <el-radio :label="1">满签</el-radio>
                <el-radio :label="2">连续签到</el-radio>
                <el-radio :label="3">累计签到</el-radio>
              </el-radio-group>
              <div
                class="tips-signin"
                v-if="activityConfigFormData.signinType == 1"
              >
                连续完成签到周期的每一天签到，即可完成签到任务
              </div>
              <div
                class="tips-signin"
                v-if="activityConfigFormData.signinType == 2"
              >
                连续签到<el-input
                  v-model="activityConfigFormData.days"
                  @input="dayInput"
                  class="num-input"
                  size="medium"
                ></el-input>天，即可完成签到任务
              </div>
              <div
                class="tips-signin"
                v-if="activityConfigFormData.signinType == 3"
              >
                累计签到<el-input
                  prop="days"
                  v-model="activityConfigFormData.days"
                  @input="dayInput"
                  class="num-input"
                  size="medium"
                ></el-input>天，即可完成签到任务
              </div>
            </cy-form-item>
          </div>
          <div class="header-title">活动范围</div>
          <cy-form-item label="范围类型">
            <span>员工/导购</span>
          </cy-form-item>
          <cy-form-item
            label="活动范围"
            prop="active_type"
            required
            style="width: 100%;"
          >
            <el-radio-group v-model="activityConfigFormData.range">
              <el-radio :label="1">全部</el-radio>
              <el-radio :label="2">指定范围</el-radio>
            </el-radio-group>
            <div
              v-show="activityConfigFormData.range == 2"
              class="range"
            >
              <el-button
                type="primary"
                size="medium"
                @click="showChoiceStaff = true"
              >选择可见范围</el-button>
              <el-table
                v-if="select_members && select_members.length > 0"
                :data="[select_members]"
                class="m-t-20"
              >
                <el-table-column
                  label="可见范围"
                  min-width="220"
                >
                  <template slot-scope="{ row }">
                    <div
                      class="item-staff"
                      v-for="(item, index) in row"
                      :key="index"
                    >
                      <img
                        class="staff-icon"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/staff_icon.png"
                        alt=""
                      />
                      <span>{{ item.name }}</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column
                  label="操作"
                  width="120"
                >
                  <template>
                    <el-button
                      type="text"
                      @click="showChoiceStaff = true"
                    >修改</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </cy-form-item>
          <div class="header-title">其它</div>
          <div v-if="activityConfigFormData.active_type === 1">
            <cy-form-item
              label="报名审核"
              required
            >
              <el-radio-group v-model="activityConfigFormData.registration_review">
                <el-radio :label="1">需要</el-radio>
                <el-radio :label="0">不需要</el-radio>
              </el-radio-group>
            </cy-form-item>
            <cy-form-item
              label="弹窗说明"
              prop="dialog_mesg"
              required
            >
              <el-input
                type="textarea"
                :rows="6"
                :maxlength="50"
                resize="none"
                show-word-limit
                placeholder="请输入弹窗说明，最多可输入50字符。如：
提交成功：
工作人员会尽快审核您的报名信息"
                v-model="activityConfigFormData.dialog_mesg"
              >
              </el-input>
              <div class="tip">
                用于用户报名活动内容成功提交后，前端弹框提示。
              </div>
            </cy-form-item>
            <cy-form-item
              label="审核通过弹窗"
              required
              style="width: 600px"
              v-if="activityConfigFormData.registration_review === 1"
            >
              <upload-image
                :urls="activityConfigFormData.arrow_pic"
                accept="image/jpeg,image/png"
                :before-upload="beforeUpload"
                @change="arrowPicSuccess"
              ></upload-image>
              <div class="tips">
                建议上传尺寸为590*560px的图片，大小不超过5000k, 建议上传 &lt;
                200k
              </div>
              <div class="tip">
                用于用户报名内容审核通过后，前端弹窗提醒
              </div>
            </cy-form-item>
            <cy-form-item
              label="审核未通过弹窗"
              required
              style="width: 600px"
              v-if="activityConfigFormData.registration_review === 1"
            >
              <upload-image
                :urls="activityConfigFormData.fail_pic"
                accept="image/jpeg,image/png"
                :before-upload="beforeUpload"
                @change="failPicSuccess"
              ></upload-image>
              <div class="tips">
                建议上传尺寸为590*560px的图片，大小不超过5000k, 建议上传 &lt;
                200k
              </div>
              <div class="tip">
                用于用户报名内容审核未通过后，前端弹窗提醒
              </div>
            </cy-form-item>
          </div>
          <cy-form-item
            label="弹窗说明"
            prop="dialog_mesg"
            required
            v-if="activityConfigFormData.active_type === 3"
          >
            <el-input
              type="textarea"
              :rows="6"
              :maxlength="50"
              resize="none"
              show-word-limit
              placeholder="请输入弹窗说明，最多可输入50字符。如：
恭喜您获得奖品：
工作人员会在7～15个工作日为您派发奖品"
              v-model="activityConfigFormData.dialog_mesg"
            >
            </el-input>
            <div class="tip">
              用于用户问卷活动内容成功提交后，前端弹框提示。
            </div>
          </cy-form-item>
          <div v-if="activityConfigFormData.active_type === 2">
            <cy-form-item
              label="是否配置签到内容"
              required
            >
              <el-radio-group v-model="activityConfigFormData.sign_in_content">
                <el-radio :label="1">是</el-radio>
                <el-radio :label="0">否</el-radio>
              </el-radio-group>
            </cy-form-item>
            <cy-form-item
              label="是否开启留资"
              required
            >
              <el-radio-group v-model="activityConfigFormData.retain_capital">
                <el-radio :label="1">开启</el-radio>
                <el-radio :label="0">关闭</el-radio>
              </el-radio-group>
            </cy-form-item>
            <cy-form-item
              label="弹窗说明"
              prop="dialog_mesg"
              required
            >
              <el-input
                type="textarea"
                :rows="6"
                :maxlength="50"
                resize="none"
                show-word-limit
                placeholder="请输入弹窗说明，最多可输入50字符。如：
恭喜您获得奖品：
工作人员会在7～15个工作日为您派发奖品"
                v-model="activityConfigFormData.dialog_mesg"
              >
              </el-input>
              <div class="tip">
                用于用户完成签到规则任务后，前端弹框提示。
              </div>
            </cy-form-item>
          </div>
        </cy-form>
      </template>
    </layout-form>
    <!-- 可见范围弹窗  -->
    <selectMember
      v-if="showChoiceStaff"
      :hasSelectArr="select_members"
      @close="showChoiceStaff = false"
      @submit="submitDialog"
    />
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import SelectMember from "@/components/business/select-member-new/index.vue";
import UserModule from "@/store/modules/user";
import Api from "@/api/wisdom/marketing-activities";
@Component({
  components: {
    UploadImage,
    SelectMember,
  },
})
export default class extends Vue {
  @Prop({ default: () => 0 }) activityId?: number;
  @Prop({ default: () => 0 }) newActivityId?: number;
  @Prop() base_info: any;
  @Watch("base_info")
  baseInfoChange(data: any) {
    this.activityConfigFormData = {
      active_type: data.activity_type,
      active_name: data.title,
      active_rule: data.activity_rule,
      banner_pic: data.cover_url,
      indate: [data.start_time * 1000, data.end_time * 1000],
      sign_in_period: [data.sign_start * 1000, data.sign_end * 1000],
      range: data.type,
      registration_review: data.is_audit,
      signinType: data.sign_type,
      days: data.sign_day,
      dialog_mesg: data.explanation,
      fail_pic: data.audit_fail_url,
      arrow_pic: data.audit_success_url,
      sign_in_content: data.is_config_sign_content,
      retain_capital: data.is_open_remark,
    };
    this.select_members = data.ba_list;
  }
  private statusType = Number(this.$route.query.status_type);
  private activityConfigFormData: any = {
    active_type: 1, // 活动类型
    active_name: "", // 活动标题
    active_rule: "", // 活动规则
    banner_pic: "", // 封面图片
    indate: [], // 活动的有效期
    sign_in_period: [],
    range: 1, // 活动范围
    registration_review: 1, // 是否需要报名审核
    signinType: 1,
    days: "",
    dialog_mesg: "", // 弹窗说明
    fail_pic: "", // 审核通过图片
    arrow_pic: "", // 审核未通过图片
    sign_in_content: 1, // 是否配置签到内容
    retain_capital: 1, // 是否留资
  };
  private rangeTabName = 0;
  private rangeType = "";
  private select_members = [];
  private showChoiceStaff = false; // 选择导购弹窗
  private activityConfigFormRules = {};
  private pickerOptionsStart = {
    disabledDate: (time: Date) => {
      return time.getTime() < Date.now() - 3600 * 1000 * 24;
    },
  };
  private pickerOptionsSignIn: any = {};
  private selectIndate() {
    if (this.activityConfigFormData.active_type === 2) {
      this.activityConfigFormData.sign_in_period = [];
      if (this.activityConfigFormData.indate.length === 0) {
        this.pickerOptionsSignIn = {};
      } else {
        this.pickerOptionsSignIn = {
          disabledDate: (time: Date) => {
            const startTime = this.activityConfigFormData.indate[0] / 1000;
            const endTime = this.activityConfigFormData.indate[1] / 1000;
            return (
              time.getTime() < startTime * 1000 - (24 * 60 * 60 * 1000 - 1) ||
              time.getTime() > endTime * 1000
            );
          },
        };
      }
    }
  }
  private selectSignIn() {
    if (this.activityConfigFormData.indate.length === 0) {
      this.$message.warning("请先选择活动有效期");
      this.activityConfigFormData.sign_in_period = [];
    }
    this.activityConfigFormData.days = "";
  }

  private dayInput(val: string) {
    const inputValue = String(val);
    let str: any = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    if (this.activityConfigFormData.sign_in_period.length === 0) {
      this.$message.warning("请先选择签到周期");
      str = "";
    } else {
      const days =
        ((this.activityConfigFormData.sign_in_period[1] / 1000 -
          this.activityConfigFormData.sign_in_period[0] / 1000) *
          1000) /
          (24 * 60 * 60 * 1000) +
        1;
      if (Number(str) > Math.floor(days)) {
        str = Math.floor(days);
        this.$message.warning("最大天数不能超过周期值");
      }
    }
    this.$nextTick(() => {
      this.activityConfigFormData.days = str;
    });
  }

  // ------------------------------------------------methods
  // 上传文件限制规范
  private beforeUpload(file: any) {
    const fileName = file.name
      .substring(file.name.lastIndexOf(".") + 1)
      .toLowerCase();
    if (
      (file.type != "image/jpeg" && file.type != "image/png") ||
      (fileName !== "png" && fileName !== "jpg")
    ) {
      this.$message.error("请上传JPG、PNG格式的图片");
      return new Promise((res, rej) => {
        return rej();
      });
    }
    if (file.size > 5000 * 1024) {
      this.$message.error("请上传5000k大小以内的图片");
      return new Promise((res, rej) => {
        return rej();
      });
    }
  }
  // 封面图片
  private bannerPicSuccess(res: any) {
    this.activityConfigFormData.banner_pic = res;
  }
  // 审核通过图片
  private arrowPicSuccess(res: any) {
    this.activityConfigFormData.arrow_pic = res;
  }
  // 审核未通过
  private failPicSuccess(res: any) {
    this.activityConfigFormData.fail_pic = res;
  }
  // 选择员工
  private submitDialog(list: any) {
    this.select_members = list;
    this.showChoiceStaff = false;
  }
  private submit() {
    if (!this.activityConfigFormData.banner_pic) {
      this.$message.warning("请上传封面图片");
      return false;
    }
    // 活动有效期开始时间校验
    if (this.activityConfigFormData.indate[0] < new Date().getTime()) {
      this.$message.warning("活动开始时间应大于当前时间");
      return;
    }
    // 活动有效期结束时间校验
    if (this.activityConfigFormData.indate[1] < new Date().getTime()) {
      this.$message.warning("活动结束时间应大于当前时间");
      return;
    }
    // 活动有效期区间校验
    if (
      this.activityConfigFormData.indate[0] ===
      this.activityConfigFormData.indate[1]
    ) {
      this.$message.warning("活动开始时间不应等于结束时间");
      return;
    }
    // 签到周期的验证
    if (this.activityConfigFormData.active_type === 2) {
      if (this.activityConfigFormData.sign_in_period.length === 0) {
        this.$message.warning("请选择签到时间范围");
        return;
      }
      if (
        this.activityConfigFormData.sign_in_period[0] <
          this.activityConfigFormData.indate[0] ||
        this.activityConfigFormData.sign_in_period[1] >
          this.activityConfigFormData.indate[1]
      ) {
        this.$message.warning("签到时间应在活动有效期内");
        return;
      }
    }
    // 报名活动
    if (this.activityConfigFormData.active_type === 1) {
      if (this.activityConfigFormData.registration_review === 1) {
        if (!this.activityConfigFormData.fail_pic) {
          this.$message.warning("请上传审核未通过图片");
          return false;
        }
        if (!this.activityConfigFormData.arrow_pic) {
          this.$message.warning("请上传审核通过图片");
          return false;
        }
      }
    }
    (this.$refs.cform as any).validate((valid: boolean) => {
      if (valid) {
        const params: any = {
          /* 所有活动公共参数 */
          base_info: {
            activity_type: this.activityConfigFormData.active_type, // 活动类型
            title: this.activityConfigFormData.active_name, // 活动标题
            activity_rule: this.activityConfigFormData.active_rule, // 活动规则
            cover_url: this.activityConfigFormData.banner_pic, // 封面图片
            start_time: this.activityConfigFormData.indate[0] / 1000, // 活动有效期开始时间
            end_time: this.activityConfigFormData.indate[1] / 1000, // 活动有效期结束时间
          },
          scope: {
            type: this.activityConfigFormData.range, // 活动范围
            ba_list: [], // 指定范围
          },
          other: {
            explanation: this.activityConfigFormData.dialog_mesg, // 弹窗说明
          },
        };
        // 活动范围被指定
        if (this.activityConfigFormData.range === 2) {
          if (this.select_members.length === 0) {
            this.$message.warning("请选择指定范围");
            return;
          }
          params.scope.ba_list = this.select_members;
        }
        if (this.activityConfigFormData.active_type === 1) {
          /* 报名活动补漏参数 */
          params.other.is_audit =
            this.activityConfigFormData.registration_review; // 报名活动时的审核
          params.other.audit_success_url =
            this.activityConfigFormData.arrow_pic; // 审核通过的弹窗图片
          params.other.audit_fail_url = this.activityConfigFormData.fail_pic; // 审核未通过的弹窗图片
        } else if (this.activityConfigFormData.active_type === 2) {
          /* 签到活动补漏参数 */
          // 签到周期开始时间
          params.base_info.sign_start =
            this.activityConfigFormData.sign_in_period[0] / 1000;
          // 签到周期结束时间
          params.base_info.sign_end =
            this.activityConfigFormData.sign_in_period[1] / 1000;
          // 签到规则
          params.base_info.sign_type = this.activityConfigFormData.signinType;
          // 天数
          if (this.activityConfigFormData.signinType !== 1) {
            if (this.activityConfigFormData.days === "") {
              this.$message.warning("请输入天数");
              return;
            }
            params.base_info.sign_day = this.activityConfigFormData.days;
          }
          // 是否配置签到内容
          params.other.is_config_sign_content =
            this.activityConfigFormData.sign_in_content;
          // 是否开启留资
          params.other.is_open_remark =
            this.activityConfigFormData.retain_capital;
        }
        if (
          this.statusType === 2 || // 编辑
          (this.statusType === 1 && this.activityId !== 0) || // 新增时回到第一步
          (this.statusType === 3 && this.newActivityId !== 0) // 复制时回到第一步
        ) {
          // 添加ID走编辑逻辑接口
          params.id =
            this.newActivityId === 0 ? this.activityId : this.newActivityId;
          params.steps = 1;
          Api.activityUpdata(params)
            .then((res: any) => {
              if (res.httpCode === 200) {
                this.$emit("submitSuccess", res.data.id || 0);
              } else {
                this.$message.warning("保存失败，请稍后重试！");
              }
            })
            .catch((err: any) => {
              console.log("activityAddSubmit", err);
            });
        } else if (this.statusType === 1 || this.statusType === 3) {
          params.steps = 1;
          Api.activityAddSubmit(params)
            .then((res: any) => {
              if (res.httpCode === 200) {
                this.$emit("submitSuccess", res.data.id || 0);
              } else {
                this.$message.warning("保存失败，请稍后重试！");
              }
            })
            .catch((err: any) => {
              console.log("activityOneSubmit", err);
            });
        }
      }
    });
  }
  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }
}
</script>
<style lang="scss" scoped>
.header-title {
  background-color: #eeeeee;
  padding: 5px 15px;
  font-size: 17px;
  margin-bottom: 25px;
}
::v-deep .layout-form .content {
  padding: 30px 0 !important;
}
.tips {
  font-size: 13px;
  color: #7f7f7f;
}
.tips-signin {
  font-size: 16px;
  color: #7f7f7f;
  margin-top: 10px;
}
.tip {
  font-size: 13px;
  color: #c0c4cc;
}
.range {
  margin-top: 15px;
}

.item-staff {
  display: inline-flex;
  align-items: center;
  height: 28px;
  background-color: #f0f3fe;
  border: solid 1px #d6dffd;
  padding: 0 10px;
  border-radius: 3px;
  margin-right: 10px;
  margin-bottom: 10px;
  min-width: 87px;
  justify-content: center;
  .staff-icon {
    width: 14px;
    height: 14px;
    margin-right: 5px;
  }
}
.num-input {
  width: 80px;
  margin: 0 3px;
}
/deep/ .el-form-item__label {
  width: 150px !important;
}
/deep/.label {
  max-width: 150px !important;
}
/deep/ .el-form-item__content {
  margin-left: 150px !important;
}
</style>
