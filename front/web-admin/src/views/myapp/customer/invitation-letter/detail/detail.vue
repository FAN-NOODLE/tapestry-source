<template>
  <div>
    <layout-form type="normal">
      <template slot="content">
        <cy-form
          ref="cform"
          :model="cardForm"
        >
          <p class="headline">基础信息</p>
          <cy-form-item
            label="模板名称"
            prop="title_name"
            required
          >
            <cy-input
              v-model="cardForm.title_name"
              clearable
              width-size="medium"
              placeholder="请输入模板名称"
              disabled
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="关联预约服务"
            required
            class="custom-width"
          >
            <el-table
              :data="serviceData"
              stripe
            >
              <el-table-column
                label="服务名称"
                prop="reserve_name"
                min-width="130"
              >
              </el-table-column>
              <el-table-column
                label="时长"
                prop="reserve_time"
              >
                <template slot-scope="scope">
                  <span>{{ scope.row.reserve_time }}小时</span>
                </template>
              </el-table-column>
              <el-table-column
                label="有效期"
                prop="limit_time"
                min-width="240"
              >
                <template slot-scope="{ row }">
                  <span v-if="row.service_type === 1">长期有效</span>
                  <span v-else>
                    {{ row.start_time | formatDate("yyyy-mm-dd") }}
                    ~
                    {{ row.end_time | formatDate("yyyy-mm-dd") }}
                  </span>
                </template></el-table-column>
            </el-table>
          </cy-form-item>
          <cy-form-item
            label="服务适用范围"
            required
            class="custom-width"
          >
            <p
              v-show="!serviceData.length"
              class="mark"
            >
              请先选择【关联服务预约】
            </p>
            <div
              v-show="serviceData.length"
              class="range-div"
            >
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
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                  alt=""
                />
                <!-- icon -> 正常员工 -->
                <img
                  class="range-icon"
                  v-else
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/staff_icon.png"
                  alt=""
                />{{ item.name }}</el-tag>
              <p
                v-show="!rangeData.length"
                class="range-p"
              >暂无适用范围</p>
            </div>
          </cy-form-item>
          <cy-form-item
            label="邀请函有效期"
            required
          >
            <template>
              <cy-form
                class="limitTimeForm"
                ref="limitTimeForm"
                :rules="rules"
                :model="cardForm"
              >
                <cy-form-item
                  label=""
                  prop="status_type"
                >
                  <template>
                    <el-radio-group
                      v-model="cardForm.status_type"
                      disabled
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
                >
                  <template>
                    <el-date-picker
                      v-model="cardForm.created_time"
                      type="datetimerange"
                      prefix-icon="el-icon-date"
                      format="yyyy-MM-dd HH:mm:ss"
                      value-format="timestamp"
                      range-separator="至"
                      start-placeholder="开始日期"
                      end-placeholder="结束日期"
                      :default-time="['00:00:00', '23:59:59']"
                      style="height:36px"
                      disabled
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
              :limit="limitNum"
              multiple
              disabled
            >
            </upload-image>
            <span class="upload-tip">可上传多张图片，建议上传尺寸为750*1624的图片，支持JPG、PNG格式，图片小于1MB</span>
          </cy-form-item>
          <cy-form-item
            label="文案"
            prop="content"
            required
          >
            <cy-input
              v-model="cardForm.content"
              clearable
              width-size="medium"
              type="textarea"
              maxlength="150"
              show-word-limit
              :autosize="{ minRows: 3, maxRows: 10 }"
              placeholder="请输入文案"
              disabled
            ></cy-input>
          </cy-form-item>

          <p class="headline">分享内容配置</p>
          <cy-form-item
            label="分享卡片标题"
            prop="share_title"
            required
          >
            <cy-input
              v-model="cardForm.share_title"
              width-size="medium"
              placeholder="请输入分享标题"
              maxlength="15"
              show-word-limit
              disabled
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
              :limit="1"
              multiple
              disabled
            >
            </upload-image>
            <span class="upload-tip">
              建议尺寸500*400；支持JPG、PNG格式，图片小于2MB</span>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button @click="goBack">返回</el-button>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import API from "@/api/customer/invitation-letter";

@Component
export default class extends Vue {
  /** data START */
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
    created_time: [],
  };

  private serviceData = [];
  private rangeData = [];

  private rules = {};

  private limitNum = 10;

  /** data END */

  /** methods START */
  private created() {
    this.inviteId = Number(this.$route.query.id);
    this.getDetail();
  }

  /** 获取详情 */
  private getDetail() {
    API.invitationDataInfo({ id: this.inviteId }).then((res: any) => {
      if (res.code === 200) {
        const background_url = JSON.parse(res.data.background_url).map(
          (item: any, index: number) => {
            return {
              name: "name_" + index,
              url: item,
            };
          }
        );
        const share_url = JSON.parse(res.data.share_url).map(
          (item: any, index: number) => {
            return {
              name: "name_" + index,
              url: item,
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
          end_time,
        } = res.data;

        const params = {
          reserve_service_id,
          title_name,
          status_type,
          content,
          share_title,
          background_url,
          share_url,
          created_time: [start_time * 1000, end_time * 1000],
        };

        // 页面提交数据
        this.cardForm = Object.assign(this.cardForm, params);

        this.getServiceInfo();
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /** 获取已选服务适用范围 */
  private getServiceInfo() {
    const params = {
      id: this.cardForm.reserve_service_id,
    };
    API.serviceInfo(params).then((res: any) => {
      if (res.code === 200) {
        const { data } = res;

        if (data.service_role === 0) {
          this.rangeData = [
            {
              name: "全部",
            },
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

        // 邀请函有效期
        // data.service_type  1--长期有效  0--指定时间
        this.cardForm.status_type = data.service_type === 0 ? 1 : 0;
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /** 返回数据列表 */
  private goBack() {
    this.$router.push({
      name: "myapp_customer_invitation-letter_list",
    });
  }

  /** methods END */
}
</script>

<style lang="scss" scoped>.custom-width {
  max-width: 800px !important;
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
}</style>
