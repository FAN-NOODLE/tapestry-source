<template>
  <div class="container">
    <layout-form type="normal">
      <template slot="content">
        <cy-form ref="form" :rules="rules">
          <cy-form-item label="服务名称：">
            <el-input
              type="text"
              placeholder="请输入服务名称"
              v-model="form.reserve_name"
              maxlength="10"
              show-word-limit
              disabled
            />
          </cy-form-item>
          <cy-form-item label="服务时长：">
            <el-select
              v-model="form.reserve_time"
              placeholder="请选择时间"
              disabled
            >
              <el-option
                label="1小时"
                :value="1"
              ></el-option>
              <el-option
                label="2小时"
                :value="2"
              ></el-option>
              <el-option
                label="3小时"
                :value="3"
              ></el-option>
              <el-option
                label="4小时"
                :value="4"
              ></el-option>
            </el-select>
          </cy-form-item>

          <cy-form-item label="活动图片：">
            <upload-image
              class="upload-template"
              :urls="form.reserve_img"
              accept="image/jpeg,image/png"
            >
            </upload-image>
          </cy-form-item>

          <cy-form-item
            label="活动有效期："
            key="service_type"
          >
            <span v-if="form.service_type === 1">长期有效</span>
            <span v-if="form.service_type === 0">
              {{ this.form.service_time[0] | formatDate("yyyy-mm-dd") }}-{{
                this.form.service_time[1] | formatDate("yyyy-mm-dd")
              }}</span>
          </cy-form-item>

          <cy-form-item
            label="适用范围："
            class="custom-width"
            key="service_role"
          >
            <div
              v-if="form.service_role == 1"
              class="range"
            >
              <el-table
                v-if="select_members && select_members.length > 0"
                :data="[select_members]"
                class="m-t-20"
              >
                <el-table-column
                  label="可见门店/员工"
                  min-width="220"
                >
                  <template slot-scope="{ row }">
                    <div
                      class="item-staff"
                      v-for="(item, index) in row"
                      :key="index"
                    >
                      <img
                        v-if="item.is_dp"
                        class="staff-icon"
                        src="@/assets/images/group-icon.png"
                        alt=""
                      />
                      <img
                        v-else
                        class="staff-icon"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/staff_icon.png"
                        alt=""
                      />
                      <span>{{ item.name }}</span>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
            <div v-else>全部</div>
          </cy-form-item>
          <cy-form-item label="时间配置项：" v-if="epName === 'coach'">
            <div
              class="time-area"
              v-for="(item, index) in timeList"
              :key="index"
            >
              <!-- 时间范围选择 -->
              <el-time-select
                placeholder="开始时间"
                style="width: 150px !important;"
                v-model="item.start_time"
                :picker-options="{
                  start: '08:00',
                  step: '01:00',
                  end: '23:59',
                }"
                disabled
              >
              </el-time-select>
              ~
              <el-time-select
                placeholder="结束时间"
                style="width: 150px !important;"
                v-model="item.end_time"
                :picker-options="{
                  start: '08:00',
                  step: '01:00',
                  end: '23:59',
                }"
                disabled
              >
              </el-time-select>
            </div>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button
          size="small"
          type="primary"
          @click="goList"
        >返 回</el-button>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Emit } from "vue-property-decorator";
import Api from "@/api/store-invite/index";
import UploadImage from "@/components/business/upload-image/index.vue";
import PerModule from "@/store/modules/permission";

interface FormItem {
  reserve_name: string;
  reserve_time: number;
  reserve_img: string;
  end_time: string | number;
  service_time: Array<any>;
  service_type: number;
  service_role: number;
}

@Component({
  components: {
    UploadImage
  }
})
export default class extends Vue {
  private rowId = Number(this.$route.query.id || "");
  private title = this.$route.query.title;
  private detailInfo = {};
  private epNameDetail = localStorage.getItem("epNameDetail");
  private form: FormItem = {
    reserve_name: "",
    reserve_time: 1,
    reserve_img: "",
    end_time: "",
    service_time: [],
    service_type: 0,
    service_role: 0
  };
  private timeList = []
  private rules = {}
  // todo: 当前品牌
  private epName = PerModule.getEpName;
  private mounted() {
    if (this.rowId) {
      this.serviceDetail(this.rowId);
    }
  }
  // 获取服务详情
  private serviceDetail(id) {
    Api.serviceDetail({ id }).then(res => {
      if (res.data) {
        this.detailInfo = res.data;
        const {
          reserve_name,
          reserve_time,
          reserve_img,
          end_time,
          service_type,
          service_role,
          start_time,
          time_slot
        } = res.data;
        this.form = {
          reserve_name,
          reserve_time,
          reserve_img,
          service_type,
          service_role,
          end_time: end_time ? (end_time - 86399) * 1000 : "", // ks sw需要
          service_time: [start_time, end_time]
        };
        if (time_slot && time_slot.length > 0) {
          this.timeList = time_slot.map(item => {
            return {
              start_time: item.start_time,
              end_time: item.end_time
            }
          })
        }
        //指定人员
        if (service_role === 1) {
          const role_rules = res.data.role_rules.map(item => {
            return {
              ...item,
              mem_subscribe_status: item.mem_subscribe_status == 1
            };
          });
          this.select_members = role_rules;
        }
      }
    });
  }

  //------------------  弹窗框选择员工 --------------------
  // 已选中员工数据
  private select_members: Array<any> = [];
  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.$emit("close");
  }
  private goList() {
    this.$router.push({
      name: "myapp_customer_store-invite_service-setting_list"
    });
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 700px !important;

  .cy-form-item__label {
    color: gray;
  }
}

.custom-width {
  max-width: 800px !important;
}

.tip-view {
  margin-top: 20px;

  .tip-text {
    font-size: 12px;
    color: $color-danger;
    line-height: 18px;
  }
}

.range {
  width: 90%;
  // padding-left: 110px;
  margin-top: 15px;

  .text-button {
    margin-left: 20px;
    color: #708df5;
    cursor: pointer;
  }
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
.time-area{
  margin-bottom: 10px;
  .time-btn{
    display: inline-block;
    span{
      font-size: 20px;
    }
    .el-icon-circle-plus{
      padding-left: 10px;
    }
    .el-icon-remove{
      padding-left: 10px;
    }
  }
}
</style>
