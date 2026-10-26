<template>
  <div class="container">
    <layout-form type="normal">
      <template slot="content">
        <cy-form
          ref="form"
          :model="form"
          :rules="rules"
        >
          <cy-form-item
            label="服务名称："
            prop="reserve_name"
          >
            <el-input
              type="text"
              placeholder="请输入服务名称"
              v-model="form.reserve_name"
              maxlength="15"
              show-word-limit
            />
          </cy-form-item>
          <cy-form-item
            label="服务时长："
            prop="reserve_time"
          >
            <el-select
              v-model="form.reserve_time"
              placeholder="请选择时间"
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

          <cy-form-item
            label="活动图片："
            prop="reserve_img"
          >
            <upload-image
              class="upload-template"
              :urls="form.reserve_img"
              accept="image/jpeg,image/png"
              :before-upload="beforeUpload"
              @change="picSuccess"
            >
            </upload-image>

            <div class="tip-view">
              <p class="tip-text">1. 图片格式仅支持*.jpg；*.png</p>
              <p class="tip-text">2. 图片尺寸630*210;</p>
              <p class="tip-text">3. 图片不得超过10M</p>
            </div>
          </cy-form-item>

          <cy-form-item
            label="活动有效期："
            key="service_type"
            prop="service_time"
          >
            <el-radio-group
              v-model="form.service_type"
              @change="serviceTypeChange"
            >
              <el-radio :label="1">长期有效</el-radio>
              <el-radio :label="0">指定时间</el-radio>
            </el-radio-group>

            <el-date-picker
              v-show="form.service_type === 0"
              type="daterange"
              v-model="form.service_time"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              range-separator="至"
              class="date-range-box"
              :picker-options="pickerOptions"
              :default-time="['00:00:00', '23:59:59']"
            ></el-date-picker>
          </cy-form-item>

          <cy-form-item
            label="适用范围："
            key="service_role"
            prop="service_role"
            class="custom-width"
          >
            <el-radio-group v-model="form.service_role">
              <el-radio :label="0">全部</el-radio>
              <el-radio :label="1">指定范围</el-radio>
            </el-radio-group>
            <div
              v-show="form.service_role == 1"
              class="range"
            >
              <el-button
                v-show="select_members && select_members.length === 0"
                type="primary"
                size="medium"
                @click="selectPerson.visible = true"
              >选择门店/员工</el-button>
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
                <el-table-column
                  label="操作"
                  width="120"
                >
                  <template>
                    <el-button
                      type="text"
                      @click="selectPerson.visible = true"
                    >修改</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </cy-form-item>
          <cy-form-item
            label="时间配置项："
            v-if="epName === 'coach'"
          >
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
                  start: '00:00',
                  step: '01:00',
                  end: '23:59',
                  minTime: getStartMinTime(index),
                  maxTime: '22:59'
                }"
                @change="changeStartTime(index)"
              >
              </el-time-select>
              ~
              <el-time-select
                placeholder="结束时间"
                style="width: 150px !important;"
                v-model="item.end_time"
                :picker-options="{
                  start: '00:00',
                  step: '01:00',
                  end: '23:59',
                  minTime: item.start_time ? item.start_time : index > 0 ? timeList[index - 1].end_time : ''
                }"
                @change="changeEndTime(index)"
              >
              </el-time-select>
              <!-- 添加，删除按钮 -->
              <div class="time-btn">
                <span
                  class="el-icon-circle-plus"
                  v-if="index === 0"
                  @click="addTimeList"
                ></span>
                <span
                  class="el-icon-remove"
                  v-if="index > 0"
                  @click="delTimeList(index)"
                ></span>
              </div>
            </div>
          </cy-form-item>

        </cy-form>
      </template>
      <template slot="footer">
        <el-button
          size="small"
          @click="goList"
        >取 消</el-button>
        <el-button
          size="small"
          type="primary"
          @click="submit"
        >确 定</el-button>
      </template>
    </layout-form>

    <!-- 选择员工弹窗 -->
    <select-person
      :selectedData="selectPerson.selectedData"
      :visible.sync="selectPerson.visible"
      @submit="selectPerson.submit"
      :appendToBody="true"
    ></select-person>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Emit } from "vue-property-decorator";
import Api from "@/api/store-invite/index";
import UploadImage from "@/components/business/upload-image/index.vue";
import SelectPerson from "@/components/business/select-person/index.vue";
import {
  selectPerson,
  callBackType
} from "@/components/business/select-person/type";
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
    UploadImage,
    SelectPerson
  }
})
export default class extends Vue {
  /*** data ***/
  private rowId = Number(this.$route.query.id || "");
  private title = this.$route.query.title;
  private detailInfo = {};
  // private showCreate = this.$route.query.showCreate;
  private form: FormItem = {
    reserve_name: "",
    reserve_time: 1,
    reserve_img: "",
    end_time: "",
    service_time: [],
    service_type: 1,
    service_role: 0
  };
  //todo: 选择时间范围默认9:00-12:00,12:00-22:00
  private timeList = [];
  // todo: 当前品牌
  private epName = PerModule.getEpName;
  // 获取开始时间的最小值
  private getStartMinTime(index) {
    if (index > 0) {
      // 获取上一个时间范围的结束时间 - 1小时
      const endTime = this.timeList[index - 1].end_time;
      if (endTime) {
        let hours: any = endTime.split(':')[0];
        let minutes: any = endTime.split(':')[1];
        if (minutes == '00') {
          hours = hours - 1;
          minutes = '59';
        } else {
          minutes = minutes - 1;
        }
        return `${hours}:${minutes}`;
      } else {
        return '';
      }
    } else {
      return '';
    }
  }
  // 开始时间change
  private changeStartTime(index) {
    if (this.timeList[index].start_time && this.timeList[index].end_time) {
      if (this.timeList[index].start_time >= this.timeList[index].end_time) {
        this.$message.error("开始时间不能大于或等于结束时间");
        this.timeList[index].start_time = "";
      }
    }
  }
  // 结束时间change
  private changeEndTime(index) {
    if (this.timeList[index].start_time && this.timeList[index].end_time) {
      if (this.timeList[index].start_time >= this.timeList[index].end_time) {
        this.$message.error("结束时间不能小于或等于开始时间");
        this.timeList[index].end_time = "";
      }
    }
  }
  // 添加时间范围
  private addTimeList() {
    // 判断是否已经有时间范围
    if (this.timeList.length > 0) {
      // 判断最后一个时间范围是否已经填写
      if (this.timeList[this.timeList.length - 1].start_time && this.timeList[this.timeList.length - 1].end_time) {
        this.timeList.push({
          start_time: '',
          end_time: ''
        });
      } else {
        this.$message.error('请填写完整时间范围');
      }
    } else {
      this.timeList.push({
        start_time: '',
        end_time: ''
      });
    }
  }
  // 删除时间范围
  private delTimeList(index) {
    this.timeList.splice(index, 1);
  }
  private pickerOptions = {
    disabledDate(time) {
      return time.getTime() < Date.now() - 8.64e7;
    }
  };
  private validateServiceTime = (rule, value, callback) => {
    console.log("validateServiceTime", value);
    if (this.form.service_type == 0) {
      if (value.length > 0) {
        callback();
      } else {
        callback(new Error("请选择活动有效期"));
      }
    } else {
      callback();
    }
  };
  //校验规则
  private rules = {
    reserve_name: [
      { required: true, message: "请填写服务名称", trigger: "blur" }
    ],
    reserve_time: [
      { required: true, message: "请选择服务时长", trigger: "change" }
    ],
    end_time: [
      { required: true, message: "请选择活动有效期", trigger: "change" }
    ],
    service_time: [
      { required: true, validator: this.validateServiceTime, trigger: "change" }
    ]
  };
  // 上传文件限制规范
  private async beforeUpload(file: any) {
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return Promise.reject();
    }
    if (file.size > 10 * 1024 * 1024) {
      this.$message.error("图片大小不超过10M");
      return Promise.reject();
    }
    const _this = this;
    await new Promise(function (resolve, reject) {
      const img = new Image();
      img.src = window.URL.createObjectURL(file);
      img.onload = () => {
        console.log(img.width, img.height, "====");
        const valid = img.width === 630 && img.height === 210;
        const proportion = img.width / img.height === 630 / 210;
        valid || proportion ? resolve(true) : reject();
      };
    }).then(
      () => {
        return file;
      },
      () => {
        (_this as any).$message.error(`上传的图片尺寸不符合, 只能是630*210!`);
        return Promise.reject();
      }
    );
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
          service_time: [start_time * 1000, end_time * 1000],
        };
        if (time_slot && time_slot.length > 0) {
          this.timeList = time_slot
        } else {
          //默认选择时间范围默认10:00-13:00,13:00-22:00
          this.timeList = [{
            start_time: '10:00',
            end_time: '13:00'
          }, {
            start_time: '13:00',
            end_time: '22:00'
          }]
        }
        //指定人员
        if (service_role === 1) {
          const role_rules = res.data.role_rules.map(item => {
            return {
              ...item,
              mem_subscribe_status: item.mem_subscribe_status == 1
            };
          });
          this.selectPerson.selectedData = this.select_members = role_rules;
        }
      }
    });
  }

  // 上传成功
  private picSuccess(res: any) {
    this.form.reserve_img = res;
  }

  // 活动有效期变化时
  private serviceTypeChange() {
    this.form.service_time = [];
    (this.$refs.form as HTMLFormElement).clearValidate("service_time");
  }

  //------------------  弹窗框选择员工 --------------------
  // 已选中员工数据
  private select_members: Array<any> = [];
  // 选择员工依赖
  private selectPerson: selectPerson = {
    visible: false,
    rawData: [], // 原始数据
    selectedData: [], // 选中ids && 员工可见范围
    submit: this.handleSubmit
  };
  // 选择员工回调
  private handleSubmit(data: callBackType) {
    console.log("SelectedData: ", data.rawData);
    this.selectPerson.selectedData = this.select_members = data.rawData;
  }

  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.$emit("close");
  }
  // 取消操作
  private goList() {
    this.$router.push({
      name: "myapp_customer_store-invite_service-setting_list"
    });
  }
  /**
   * 保存
   */
  private submit() {
    const {
      reserve_img,
      service_time,
      reserve_name,
      reserve_time,
      service_type,
      service_role,
      end_time
    } = this.form;

    (this.$refs.form as HTMLFormElement).validate(valid => {
      if (valid) {
        //活动图片必传
        if (!this.form.reserve_img) {
          this.$message.error("请上传活动图片");
          return;
        }
        // todo: 时间配置项验证
        if (this.timeList.length > 0 && this.epName === 'coach') {
          for (let i = 0; i < this.timeList.length; i++) {
            if (!this.timeList[i].start_time || !this.timeList[i].end_time) {
              this.$message.error("请填写完整时间范围");
              return;
            }
          }
        }
        // todo: 时间配置项不能超过15个
        if (this.timeList.length >= 15 && this.epName === 'coach') {
          this.$message.error("时间配置项不能超过15个");
          return;
        }
        const params = {
          reserve_name,
          reserve_time,
          reserve_img,
          service_type,
          end_time: Number(end_time) / 1000 + 86399, // ks、sw-有值
          id: 0,
          role_rules: [],
          service_role,
          start_time: 0,// // ks、sw-无值
          time_slot: this.epName === 'coach' ? this.timeList : [] // ks、sw-有值
        };

        // 选择指定时间
        if (service_type === 0) {
          params.start_time = service_time ? (service_time[0] || 0) / 1000 : 0;
          params.end_time = service_time ? (service_time[1] || 0) / 1000 : 0;
        }
        // 选择员工范围 数据处理
        if (service_role === 1) {
          console.log("人员", this.select_members);
          let role_rules = [];
          role_rules = this.select_members.map(item => {
            return {
              id: item.id,
              third_id: item.third_id || 0,
              name: item.name,
              is_dp: item.is_dp,
              store_type: item.store_type || 0,
              mem_subscribe_status: item.mem_subscribe_status ? 1 : 0,
              user_id: item.user_id
            };
          });
          params.role_rules = role_rules;
        }

        let api = "updateServiceApi";
        let txt = "编辑成功";
        if (!this.rowId) {
          api = "createServiceApi";
          txt = "新增成功";
          delete params.id;
        } else {
          params.id = this.rowId;
        }
        Api[api](params).then(res => {
          if (res.code !== 200) return this.$message.error(res.message);
          this.$message.success(txt);
          this.$router.push({
            name: "myapp_customer_store-invite_service-setting_list"
          });
        });
      } else {
        console.log("error submit!!");
        return false;
      }
    });
  }
  private mounted() {
    if (this.rowId) {
      this.serviceDetail(this.rowId);
    } else {
      this.timeList = [{
        start_time: '10:00',
        end_time: '13:00'
      }, {
        start_time: '13:00',
        end_time: '22:00'
      }]
    }
  }
}
</script>

<style lang="scss" scoped>.custom-width {
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
  width: 100%;
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

.time-area {
  margin-bottom: 10px;

  .time-btn {
    display: inline-block;

    span {
      font-size: 20px;
    }

    .el-icon-circle-plus {
      padding-left: 10px;
    }

    .el-icon-remove {
      padding-left: 10px;
    }
  }
}</style>
