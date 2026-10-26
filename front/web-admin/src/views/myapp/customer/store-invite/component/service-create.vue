<template>
  <div class="contaniner">
    <el-dialog
      width="1000px"
      :title="title"
      :visible.sync="dialogVisible"
      :close-on-click-modal="false"
      @close="closeDialog"
    >
      <div class="dialog-body">
        <el-form ref="form" label-width="100px">
          <el-form-item label="服务名称：">
            <el-input
              type="text"
              placeholder="请输入服务名称"
              v-model="form.reserve_name"
              maxlength="10"
              show-word-limit
              style="width: 50%"
            />
          </el-form-item>
          <el-form-item label="服务时长：">
            <el-select
              v-model="form.reserve_time"
              placeholder="请选择时间"
              style="width: 50%"
            >
              <el-option label="1小时" :value="1"></el-option>
              <el-option label="2小时" :value="2"></el-option>
              <el-option label="3小时" :value="3"></el-option>
              <el-option label="4小时" :value="4"></el-option>
            </el-select>
          </el-form-item>

          <!-- !coach -->
          <el-form-item label="活动有效期：" v-if="epNameDetail !== 'coach'">
            <el-date-picker
              v-model="form.end_time"
              type="date"
              style="width: 50%"
              placeholder="选择日期"
            >
            </el-date-picker>
          </el-form-item>
          <el-form-item label="活动图片：">
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
          </el-form-item>

          <!-- coach -->
          <el-form-item label="活动有效期：" v-if="epNameDetail === 'coach'">
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
            ></el-date-picker>
            <!-- <el-date-picker
              prefix-icon="el-icon-date"
              :default-time="['00:00:00', '23:59:59']"
            > -->
          </el-form-item>

          <!-- coach -->
          <el-form-item label="适用范围：" v-if="epNameDetail === 'coach'">
            <el-radio-group v-model="form.service_role">
              <el-radio :label="0">全部</el-radio>
              <el-radio :label="1">指定范围</el-radio>
            </el-radio-group>
            <div v-show="form.service_role == 1" class="range">
              <el-button
                type="primary"
                size="medium"
                @click="selectPerson.visible = true"
                >选择员工</el-button
              >
              <el-table
                v-if="select_members && select_members.length > 0"
                :data="[select_members]"
                class="m-t-20"
              >
                <el-table-column label="可见部门" min-width="220">
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
                <el-table-column label="操作" width="120">
                  <template>
                    <el-button type="text" @click="selectPerson.visible = true"
                      >修改</el-button
                    >
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-form-item>
        </el-form>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialog">取 消</el-button>
        <el-button size="small" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>

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
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import Api from "@/api/store-invite/index";
import UploadImage from "@/components/business/upload-image/index.vue";
import SelectPerson from "@/components/business/select-person/index.vue";
import {
  selectPerson,
  callBackType
} from "@/components/business/select-person/type";
import { ArrayIterator } from "lodash";
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
  @Prop({ default: "新增服务" }) title?: string; // 右侧选择标题
  @Prop({ default: {} }) row: object; // 右侧选择标题
  /*** data ***/
  private dialogVisible = true; // 控制模态框显示/隐藏
  private epNameDetail = localStorage.getItem("epNameDetail");
  private form: FormItem = {
    reserve_name: "",
    reserve_time: 1,
    reserve_img: "",
    end_time: "",
    service_time: [],
    service_type: 1, // 1长期有效，0指定时间
    service_role: 0
  };
  // 上传文件限制规范
  private async beforeUpload(file: any) {
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    const _this = this;
    await new Promise(function(resolve, reject) {
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
  private mounted() {
    const row = this.row as any;
    if (!row.id) return;
    const {
      reserve_name,
      reserve_time,
      reserve_img,
      end_time,
      service_type,
      start_time,
      service_role
    } = row;

    this.form = {
      reserve_name,
      reserve_time,
      reserve_img,
      service_type,
      service_role,
      service_time: [start_time * 1000, end_time * 1000],
      end_time: end_time ? (end_time - 86399) * 1000 : "" // ks sw需要
    };
  }
  // 上传成功
  private picSuccess(res: any) {
    this.form.reserve_img = res;
  }
  // 活动有效期变化时
  private serviceTypeChange() {
    this.form.service_time = [];
  }

  //员工可见范围
  // type: 1, //员工可见范围   对应数据见【弹窗框选择员工数据select_members】

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
  // 选择员工回调;
  private handleSubmit(data: callBackType) {
    this.selectPerson.selectedData = this.select_members = data.rawData;
    console.log("SelectedData: ", data.rawData);
  }

  //处理员工可见范围数据
  private setSelectedPerson() {
    // if (scope.type == 2) {
    //   if (scope.ba_list.dp_list) {
    //     this.selectPerson.selectedData = this.select_members =
    //       scope.ba_list.dp_list;
    //   }
    //   if (scope.ba_list.user_list) {
    //     this.selectPerson.selectedData = this.select_members =
    //       scope.ba_list.user_list;
    //   }
    //   if (scope.ba_list.dp_list && scope.ba_list.user_list) {
    //     this.selectPerson.selectedData = this.select_members = scope.ba_list.dp_list.concat(
    //       scope.ba_list.user_list
    //     );
    //   }
    // }
  }

  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.$emit("close");
  }
  private submit() {
    const {
      reserve_img,
      service_time,
      end_time,
      reserve_name,
      reserve_time,
      service_type,
      service_role
    } = this.form;

    if (!reserve_name) {
      this.$message.warning("请填写服务名称");
      return;
    }
    if (!service_time && service_type === 0) {
      this.$message.warning("请选择活动有效期时间");
      return;
    }

    // this.ruleForm.start_time = parseInt((val[0] / 1000).toString());
    // this.ruleForm.end_time = parseInt((val[1] / 1000).toString());
    const params = {
      reserve_name,
      reserve_time,
      reserve_img,
      service_type,
      end_time: Number(end_time) / 1000 + 86399, // ks、sw-有值
      id: 0,
      role_rules: [],
      service_role,
      start_time: 0 // // ks、sw-无值
    };

    // 选择指定时间
    if (service_type === 0) {
      params.start_time = service_time ? (service_time[0] || 0) / 1000 : 0;
      params.end_time = service_time ? (service_time[1] || 0) / 1000 : 0;
    }

    // 选择员工范围 数据处理
    if (service_role === 1) {
      // 选择范围为指定员工
      // const dp_arr: Array<any> = this.select_members.filter(item => {
      //   if (item.is_dp) {
      //     return item.id;
      //   }
      // });
      params.role_rules = this.select_members;
    }

    let api = "updateServiceApi";
    let txt = "编辑成功";
    if (!(this.row as any).id) {
      api = "createServiceApi";
      txt = "新增成功";
      delete params.id;
    } else {
      params.id = (this.row as any).id;
    }
    Api[api](params).then(() => {
      this.$message.success(txt);
      this.$emit("submit");
    });
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 1500px !important;
  .el-form-item__label {
    color: gray;
  }
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
.date-range-box {
  height: 36px;
  margin-left: 20px;
}
</style>
