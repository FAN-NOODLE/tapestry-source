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
              disabled
            />
          </el-form-item>
          <el-form-item label="服务时长：">
            <el-select
              v-model="form.reserve_time"
              placeholder="请选择时间"
              style="width: 50%"
              disabled
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
              disabled
            >
            </el-date-picker>
          </el-form-item>
          <el-form-item label="活动图片：">
            <upload-image
              class="upload-template"
              :urls="form.reserve_img"
              accept="image/jpeg,image/png"
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
            <el-radio-group v-model="form.service_type" disabled>
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
              disabled
            ></el-date-picker>
            <!-- <el-date-picker
              prefix-icon="el-icon-date"
              :default-time="['00:00:00', '23:59:59']"
            > -->
          </el-form-item>

          <!-- coach -->
          <el-form-item label="适用范围：" v-if="epNameDetail === 'coach'">
            <el-radio-group v-model="form.service_role" disabled>
              <el-radio :label="0">全部</el-radio>
              <el-radio :label="1">指定范围</el-radio>
            </el-radio-group>
            <div v-show="form.service_role == 1" class="range">
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
              </el-table>
            </div>
          </el-form-item>
        </el-form>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="closeDialog">关 闭</el-button>
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

  // 关闭模态框的回调
  @Emit("close")
  private closeDialog() {
    this.$emit("close");
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
