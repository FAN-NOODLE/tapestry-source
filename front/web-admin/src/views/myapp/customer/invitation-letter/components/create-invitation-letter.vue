<template>
  <div>
    <layout-form type="normal">
      <template slot="content">
        <cy-form ref="cform" :rules="rules" :model="cardForm">
          <p class="headline">基础信息</p>
          <cy-form-item label="模板名称" prop="title_name" required>
            <cy-input
              maxlength="15"
              show-word-limit
              v-model="cardForm.title_name"
              clearable
              width-size="medium"
              placeholder="请输入模板名称"
            ></cy-input>
          </cy-form-item>
          <cy-form-item label="关联预约服务" required class="custom-width">
            <el-button
              v-show="!serviceData.length"
              type="primary"
              class="select-material-btn"
              size="small"
              @click="handleToSelectService"
              ><span>选择预约服务</span></el-button
            >
            <el-table
              v-show="serviceData.length"
              :data="serviceData"
              stripe
              row-key="sid"
              :row-class-name="serviceData.length > 1 ? 'pointer' : ''"
            >
              <el-table-column
                label="服务名称"
                prop="reserve_name"
                min-width="130"
              >
              </el-table-column>
              <el-table-column label="时长" prop="reserve_time">
                <template slot-scope="scope">
                  <span>{{ scope.row.reserve_time }}小时</span>
                </template>
              </el-table-column>
              <el-table-column label="有效期" prop="limit_time" min-width="240">
                <template slot-scope="{ row }">
                  <span v-if="row.service_type === 1">长期有效</span>
                  <span v-else>
                    {{ row.start_time | formatDate("yyyy-mm-dd") }}
                    ~
                    {{ row.end_time | formatDate("yyyy-mm-dd") }}
                  </span>
                </template></el-table-column
              >
              <el-table-column label="操作" fixed="right">
                <template>
                  <el-button
                    @click="handleToSelectService"
                    type="text"
                    size="small"
                    >修改</el-button
                  >
                </template>
              </el-table-column>
            </el-table>
          </cy-form-item>
          <cy-form-item label="服务适用范围" required class="custom-width">
            <p v-show="!serviceData.length" class="mark">
              请先选择【关联服务预约】
            </p>
            <div v-show="serviceData.length" class="range-div">
              <header class="range-header">适用范围</header>
              <el-tag
                v-show="rangeData.length"
                class="range-data"
                v-for="(item, index) in rangeData"
                :key="index"
              >
                <!-- icon -> 部门 -->
                <img
                  class="range-icon"
                  v-if="item.is_dp"
                  src="@/assets/images/group-ic.png"
                  alt=""
                />
                <!-- icon -> 正常员工 -->
                <img
                  class="range-icon"
                  v-else
                  src="@/assets/images/staff_icon.png"
                  alt=""
                />
                {{ item.name }}</el-tag
              >
              <p v-show="!rangeData.length" class="range-p">暂无适用范围</p>
            </div>
          </cy-form-item>
          <cy-form-item label="邀请函有效期" required>
            <template>
              <p v-show="!serviceData.length" class="mark">
                请先选择【关联服务预约】
              </p>
              <cy-form
                class="limitTimeForm"
                ref="limitTimeForm"
                :rules="rules"
                :model="cardForm"
                v-show="serviceData.length"
              >
                <cy-form-item label="" prop="status_type">
                  <template>
                    <el-radio-group
                      :disabled="statusTypeDisabled"
                      v-model="cardForm.status_type"
                      @change="handleStatusTypeChange"
                    >
                      <el-radio :label="0">长期有效</el-radio>
                      <el-radio :label="1">指定时间</el-radio>
                    </el-radio-group>
                  </template>
                </cy-form-item>
                <cy-form-item
                  label="有效期"
                  prop="created_time"
                  v-show="cardForm.status_type === 1"
                  required
                  tigger="change"
                >
                  <template>
                    <el-date-picker
                      v-model="cardForm.created_time"
                      format="yyyy-MM-dd"
                      type="daterange"
                      range-separator="至"
                      start-placeholder="开始日期"
                      end-placeholder="结束日期"
                      :picker-options="pickerOptions"
                    >
                    </el-date-picker>
                  </template>
                </cy-form-item>
              </cy-form>
            </template>
          </cy-form-item>
          <cy-form-item
            style="max-width:1000px"
            label="背景图"
            prop="background_url"
            required
            trigger="change"
          >
            <upload-image
              class="upload-template"
              :urls="cardForm.background_url"
              accept="image/jpeg,image/png"
              img-type="batch"
              :before-upload="beforeUpload"
              @update="update"
              :limit="limitNum"
              multiple
            >
            </upload-image>
            <div class="upload-tip">
              可上传多张图片，建议上传尺寸为690*440的图片，支持JPG、PNG格式，图片小于1MB
            </div>
            <div class="upload-tip" style="color:#ff0000">
              注: 因邀请函样式设计，图片底部有
              59px距离可能会被遮挡，建议保存时可点击[预览] 提前查看效果
            </div>
          </cy-form-item>
          <cy-form-item label="文案" prop="content" required>
            <cy-input
              v-model="cardForm.content"
              clearable
              width-size="medium"
              type="textarea"
              maxlength="150"
              show-word-limit
              :autosize="{ minRows: 5, maxRows: 10 }"
              placeholder="请输入文案"
            ></cy-input>
          </cy-form-item>

          <p class="headline">分享内容配置</p>
          <cy-form-item
            label="分享卡片标题"
            prop="share_title"
            required
            trigger="change"
          >
            <cy-input
              v-model="cardForm.share_title"
              width-size="medium"
              placeholder="请输入分享标题"
              maxlength="15"
              show-word-limit
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            style="max-width:1000px"
            label="分享卡片图片"
            prop="share_url"
            required
            trigger="change"
          >
            <upload-image
              class="upload-template"
              :urls="cardForm.share_url"
              accept="image/jpeg,image/png"
              img-type="batch"
              :before-upload="beforeUploadSharePic"
              @update="updateSharePic"
              :limit="1"
              multiple
            >
            </upload-image>
            <span class="upload-tip">
              建议尺寸500*400；支持JPG、PNG格式，图片小于2MB</span
            >
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button @click="goCancelFn">取消</el-button>
        <el-button
          :disabled="cardForm.content === '' || !cardForm.background_url.length"
          type="primary"
          @click="saveClickFn('preview')"
          >预览</el-button
        >
        <el-button type="primary" @click="saveClickFn">保存</el-button>
      </template>
    </layout-form>
    <!-- 创建离开页面弹框 -->
    <el-dialog
      title="离开页面提示"
      :visible.sync="isCancelDialog"
      width="420px"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        检测到未保存的内容，是否保存到下次新增继续编辑？</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goBackFn">直接离开</el-button>
        <el-button type="primary" @click="saveGoBackFn">保存并离开</el-button>
      </span>
    </el-dialog>
    <!-- 编辑离开页面弹框 -->
    <el-dialog
      title="离开页面提示"
      :visible.sync="isTipsCancel"
      width="420px"
      @close="goCloseFn"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        你确定要放弃编辑吗？放弃后编辑的信息将不会被保存</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goCloseFn">取消</el-button>
        <el-button type="primary" @click="editGoBackFn">确定</el-button>
      </span>
    </el-dialog>

    <SelectService
      title="选择预约服务"
      :selected="selectData"
      @close="showSelectService = false"
      @submit="selectConfirm"
      v-if="showSelectService"
    />

    <el-dialog
      title="预览二维码"
      :visible.sync="showPreview"
      width="30%"
      :before-close="handleClose"
    >
      <div class="previewTips">
        <img :src="previewCode" />
        <p>请使用微信扫一扫，扫码预览哦~</p>
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="showPreview = false">关闭</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/customer/invitation-letter";
import UploadImage from "@/components/business/upload-image/index.vue";
import SelectService from "./select-service.vue";

@Component({
  components: { UploadImage, SelectService }
})
export default class extends Vue {
  @Prop() type?: string;

  /*************************** data start ***************************/
  /** 选择预约服务弹窗是否显示 */
  private showSelectService = false;

  /** 取消按钮是否弹框 */
  private isCancelDialog = false;
  /** 编辑取消按钮是否弹框 */
  private isTipsCancel = false;
  /** 当前id */
  private inviteId = 0;
  private cardForm = {
    title_name: "",
    reserve_service_id: "",
    status_type: 0,
    start_time: "",
    end_time: "",
    background_url: [],
    content: "",
    share_title: "",
    share_url: [],
    created_time: []
  };
  private serviceData = [];
  private selectData = {};
  private preObj = {};
  private rules = {};
  private rangeData = [];

  private isUploading = false;
  private limitNum = 10;

  /** 预览窗口是否显示 */
  private showPreview = false;
  /** 预览二维码 */
  private previewCode = "";
  /** 邀请函有效期是否可切换 */
  private statusTypeDisabled = false;
  /** 禁用时间 */
  private pickerOptions = {
    // 只能选择今天之后的日期
    disabledDate: (time: any) => {
      return time.getTime() < Date.now() - 8.64e7;
    }
  };

  /** methods */
  private created() {
    if (this.type == "edit") {
      this.inviteId = Number(this.$route.query.id);
      this.getDetailFn();
    }
  }
  // 选择预约服务
  handleStatusTypeChange(val: any) {
    if (val === 1 && this.type === "create") {
      this.pickerOptions = {
        // 只能选择今天之后的日期
        disabledDate: (time: any) => {
          return time.getTime() < Date.now() - 8.64e7;
        }
      };
    } else {
      this.cardForm.created_time = [];
    }
    if (val === 0) {
      this.cardForm.created_time = ["", ""];
    }
  }
  /** 创建页面进来，拿sessionStorage值 */
  private createShow() {
    this.cardForm = JSON.parse(
      sessionStorage.getItem("create_card_obj") || JSON.stringify(this.cardForm)
    );
  }

  /** 获取详情 */
  private getDetailFn() {
    API.invitationDataInfo({ id: this.inviteId }).then((res: any) => {
      if (res.code === 200) {
        const background_url = JSON.parse(res.data.background_url).map(
          (item: any, index: number) => {
            return {
              name: "name_" + index,
              url: item
            };
          }
        );
        const share_url = JSON.parse(res.data.share_url).map(
          (item: any, index: number) => {
            return {
              name: "name_" + index,
              url: item
            };
          }
        );

        const {
          title_name,
          status_type,
          content,
          share_title,
          reserve_service_id,
          start_time,
          end_time
        } = res.data;

        const params = {
          reserve_service_id,
          title_name,
          status_type,
          content,
          share_title,
          background_url,
          share_url,
          created_time: [start_time * 1000, end_time * 1000]
        };

        // 页面提交数据
        this.cardForm = Object.assign(this.cardForm, params);

        this.getServiceInfo(false);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /** 选择预约服务、修改关联预约服务 */
  private handleToSelectService() {
    this.showSelectService = true;
  }

  /** 确认选择预约服务 */
  private selectConfirm(data: any) {
    // 关联预约服务ID
    this.cardForm.reserve_service_id = data.id;

    // 关闭弹窗
    this.showSelectService = false;
    this.getServiceInfo(true);
  }

  /** 获取已选服务适用范围 */
  private getServiceInfo(isChange) {
    const params = {
      id: this.cardForm.reserve_service_id
    };
    API.serviceInfo(params).then((res: any) => {
      if (res.code === 200) {
        const { data } = res;
        if (data.service_role === 0) {
          this.rangeData = [
            {
              name: "全部"
            }
          ];
        } else {
          const dpList = data.role_rules.filter((item: any) => {
            return item.is_dp;
          });
          const employeesList = data.role_rules.filter((item: any) => {
            return !item.is_dp;
          });
          this.rangeData = dpList.concat(employeesList);
        }

        this.serviceData = [data];
        this.selectData = data;

        // 邀请函有效期
        // data.service_type  1--长期有效  0--指定时间
        if (isChange) {
          this.cardForm.status_type = data.service_type === 0 ? 1 : 0;
          this.statusTypeDisabled = data.service_type === 0;

          if (this.statusTypeDisabled) {
            this.cardForm.created_time = [
              data.start_time * 1000,
              data.end_time * 1000
            ];

            this.pickerOptions = {
              disabledDate: (time: any) => {
                return (
                  time.getTime() < data.start_time * 1000 ||
                  time.getTime() > data.end_time * 1000
                );
              }
            };
          } else {
            this.cardForm.created_time = ["", ""];
            this.pickerOptions = {
              // 只能选择今天之后的日期
              disabledDate: (time: any) => {
                return time.getTime() < Date.now() - 8.64e7;
              }
            };
          }
        }
        // 详情备份，用来点击取消按钮进行数据对比
        this.preObj = Object.assign({}, this.cardForm);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /** 保存 */
  async saveClickFn(operateType) {
    const cFormValidate = await (this.$refs[
      "cform"
    ] as HTMLFormElement).validate();
    if (!this.serviceData.length) {
      this.$message.warning("请选择预约服务");
      return;
    }

    const limitTimeFormValidate = await (this.$refs[
      "limitTimeForm"
    ] as HTMLFormElement).validate();
    if (limitTimeFormValidate && cFormValidate) {
      const background_url = this.cardForm.background_url.map(
        (item: { url: string }) => item.url
      );
      const share_url = this.cardForm.share_url.map(
        (item: { url: string }) => item.url
      );

      const {
        reserve_service_id,
        title_name,
        status_type,
        content,
        share_title,
        created_time
      } = this.cardForm;

      const params = {
        reserve_service_id,
        title_name,
        status_type,
        background_url,
        content,
        share_title,
        share_url,
        start_time: Number(created_time[0] / 1000),
        end_time: Number(created_time[1] / 1000),

        id: this.inviteId
      };
      if (this.type !== "edit") {
        delete params.id;
      }

      API.invitationDataAdd(params).then((res: any) => {
        if (res.code === 200) {
          if (operateType === "preview") {
            this.inviteId = res.data.id;
            this.type = "edit";
            // 预览
            API.invitationPreview({
              path: "pagesOrder/pages/invite-check/invite-check",
              params: `id=${res.data.id}&from=code`,
              apiVersion: "old"
            }).then((res: any) => {
              if (res.code === 200) {
                this.previewCode = res.data.code_url;
                this.showPreview = true;
              } else {
                this.$message.error(res.message);
              }
            });
          } else {
            this.$message.success("保存成功");
            this.routeBackFn();
          }
        } else {
          this.$message.error(res.message);
        }
      });
    }
  }

  /** 关闭预览二维码弹窗 */
  private handleClose() {
    this.showPreview = false;
  }

  /** 判断数据是否更改 */
  private isObjectValueEqual(a: any, b: any) {
    //取对象a和b的属性名
    const aProps = Object.getOwnPropertyNames(a);
    const bProps = Object.getOwnPropertyNames(b);
    //判断属性名的length是否一致
    if (aProps.length != bProps.length) {
      return false;
    }
    //循环取出属性名，再判断属性值是否一致
    for (let i = 0; i < aProps.length; i++) {
      const propName = aProps[i];
      if (a[propName] !== b[propName]) {
        return false;
      }
    }
    return true;
  }

  /** 取消 */
  private goCancelFn() {
    if (this.type === "create") {
      this.createCancelFn();
    } else if (this.type === "edit") {
      this.editCancelFn();
    }
  }

  /** 取消创建 */
  private createCancelFn() {
    if (!this.isObjectValueEqual(this.cardForm, this.preObj)) {
      this.isCancelDialog = true;
    } else {
      this.goBackFn();
    }
  }

  /** 取消编辑 */
  private editCancelFn() {
    if (!this.isObjectValueEqual(this.cardForm, this.preObj)) {
      this.isTipsCancel = true;
    } else {
      this.editRouteBackFn();
    }
  }

  /** 返回上一页 */
  private goBackFn() {
    this.isCancelDialog = false;
    this.routeBackFn();
  }

  /** 保存返回上一页 */
  private saveGoBackFn() {
    this.isCancelDialog = false;
    // 新建的值保存sessionStorage
    sessionStorage.setItem("create_card_obj", JSON.stringify(this.cardForm));
    this.routeBackFn();
  }

  /** 路由跳转 */
  private routeBackFn() {
    this.$router.push({
      name: "myapp_customer_invitation-letter_list"
    });
  }

  /** 路由跳转 */
  private editRouteBackFn() {
    this.$router.push({
      name: "myapp_customer_invitation-letter_list",
      query: { id: this.inviteId.toString() }
    });
  }

  /** 编辑关闭和取消 */
  private goCloseFn() {
    this.isTipsCancel = false;
  }

  /** 编辑确定 */
  private editGoBackFn() {
    this.isTipsCancel = true;
    this.editRouteBackFn();
  }

  /** 背景图上传文件限制规范 */
  private beforeUpload(file: any) {
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    if (name === "jpeg") {
      this.$message.warning("企业微信不支持JPEG的图片");
      return false;
    }
    if (file.size > 1024 * 1024) {
      this.$message.error("请上传1Mb大小以内的图片");
      return false;
    }
    this.isUploading = true;
  }

  /** 背景图上传完成 */
  private update(res: any) {
    this.cardForm.background_url = res;
    this.isUploading = false;
  }

  /** 分享卡片图片上传文件限制规范 */
  private beforeUploadSharePic(file: any) {
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    if (name === "jpeg") {
      this.$message.warning("企业微信不支持JPEG的图片");
      return false;
    }
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2Mb大小以内的图片");
      return false;
    }
    this.isUploading = true;
  }

  /** 分享卡片图片上传完成 */
  private updateSharePic(res: any) {
    this.cardForm.share_url = res;
    this.isUploading = false;
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-input__count {
  background: transparent;
  bottom: -8px;
  right: 8px;
}
.custom-width {
  max-width: 800px !important;
}
.upload-tip {
  font-size: 12px;
  color: #666666;
}

.headline {
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 30px;
}

.mark {
  font-size: 14px;
  font-weight: normal;
  color: #999;
}

.limitTimeForm {
  /deep/ .el-form-item__content {
    margin-left: 0 !important;
  }
}

.range-header {
  height: 36px;
  background-color: #ebebeb;
  padding-left: 20px;
  color: #606266;
  font-weight: 600;
}

.range-data {
  margin-left: 20px;
  margin-top: 17px;
}

.range-div {
  padding-bottom: 17px;
  border-bottom: 1px solid #ebebeb;
}

.range-p {
  margin-left: 20px;
  margin-top: 17px;
}

.range-icon {
  margin: auto;
}

.previewTips {
  width: 100%;
  text-align: center;

  img {
    margin-bottom: 15px;
  }
}
</style>
