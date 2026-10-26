<template>
  <div>
    <!-- 选择门店下人人员弹窗-->
    <selectMember
      v-if="showSelectStaff"
      :rootId="{
        id: storeInfo.wechat_department_id,
        name: storeInfo.third_store_code
      }"
      :hasSelectArr="this.staffList[0].staff"
      @submit="getMember"
      dialogTitle="选择员工"
      @close="showSelectStaff = false"
    />

    <el-dialog
      :title="is_edit ? '编辑门店' : '新增门店'"
      :visible="true"
      @close="closeDialog"
      :close-on-click-modal="false"
    >
      <div class="dialog-body">
        <cy-form
          class="main-box"
          ref="cform"
          :rules="rules"
          :model="formData"
          label-width="140px"
        >
          <cy-form-item
            label="部门ID"
            prop="wechat_department_id"
            v-if="epNameDetail === 'coach'"
            class="item-mb"
          >
            <div v-if="is_edit" class="txt-line-height">
              {{ storeInfo.wechat_department_id }}
            </div>
            <cy-input
              v-else
              placeholder="请输入部门ID"
              v-model.trim="formData.wechat_department_id"
              size="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item label="门店编码" prop="store_code">
            <div v-if="is_edit" class="txt-line-height">
              {{ storeInfo.third_store_code }}
            </div>
            <cy-input
              v-else
              placeholder="请输入门店编码"
              v-model.trim="formData.store_code"
              size="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item v-if="epName !== 'coach'" label="门店对外名称" prop="store_name">
            <cy-input
              :maxlength="32"
              placeholder="请输入门店对外名称"
              v-model.trim="formData.store_name"
              size="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item v-if="epName === 'coach'" label="门店中文名称" prop="store_name">
            <cy-input
              :maxlength="32"
              placeholder="请输入门店中文名称"
              v-model.trim="formData.store_name"
              size="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item label="门店属性" prop="store_type">
            <cy-select v-model="formData.store_type" size="medium">
              <cy-option label="奥莱" :value="1"></cy-option>
              <cy-option label="零售" :value="2"></cy-option>
            </cy-select>
          </cy-form-item>
          <cy-form-item label="门店地址" prop="address">
            <cy-input
              :maxlength="40"
              show-word-limit
              placeholder="请输入地址"
              v-model.trim="formData.address"
              size="medium"
            ></cy-input>
          </cy-form-item>
          <!-- todo:coach库存专属字段 -->
          <div v-if="epNameDetail == 'coach'">
            <cy-form-item v-if="epName !== 'coach'" label="库存查询门店名" prop="inventory_name">
              <cy-input
                placeholder="请输入库存查询门店名"
                v-model.trim="formData.inventory_name"
                size="medium"
              ></cy-input>
            </cy-form-item>
            <cy-form-item v-if="epName === 'coach'" label="门店英文名称" prop="inventory_name">
              <cy-input
                placeholder="请输入门店英文名称"
                v-model.trim="formData.inventory_name"
                size="medium"
              ></cy-input>
            </cy-form-item>
            <cy-form-item label="门店所属城市" prop="city">
              <cy-input
                placeholder="请输入门店所属城市"
                v-model.trim="formData.city"
                size="medium"
              ></cy-input>
            </cy-form-item>
            <cy-form-item label="门店电话" prop="phone">
              <cy-input
                placeholder="请输入门店电话"
                v-model.trim="formData.phone"
                size="medium"
              ></cy-input>
            </cy-form-item>
          </div>
          <!-- coach 定制线上订单功能 -->
          <template v-if="epNameDetail === 'coach'">
            <cy-form-item label="门店邮箱" prop="email">
              <cy-input
                placeholder="请输入地址"
                v-model.trim="formData.email"
                size="medium"
              ></cy-input>
            </cy-form-item>
            <cy-form-item
              key="popup_store_status"
              label="Popup Store"
              prop="popup_store_status"
            >
              <el-radio v-model="formData.popup_store_status" :label="1"
                >支持</el-radio
              >
              <el-radio v-model="formData.popup_store_status" :label="0"
                >不支持</el-radio
              >
            </cy-form-item>
            <template v-if="formData.popup_store_status === 1">
              <cy-form-item
                key="sa_share_order_url"
                label="SA分享订单链接"
                prop="sa_share_order_url"
              >
                <cy-input
                  :maxlength="150"
                  show-word-limit
                  v-model.trim="formData.sa_share_order_url"
                  size="medium"
                ></cy-input>
              </cy-form-item>
              <cy-form-item
                label="店长查看门店订单链接"
                prop="sa_check_order_url"
                key="sa_check_order_url"
              >
                <cy-input
                  :maxlength="150"
                  show-word-limit
                  v-model.trim="formData.sa_check_order_url"
                  size="medium"
                ></cy-input>
              </cy-form-item>
              <cy-form-item
                label="查看门店全部订单权限"
                prop="sa_check_role_status"
                key="sa_check_role_status"
              >
                <el-radio v-model="formData.sa_check_role_status" :label="1"
                  >仅店长</el-radio
                >
                <el-radio v-model="formData.sa_check_role_status" :label="2"
                  >店长和指定店员</el-radio
                >
                <el-button
                  v-show="this.formData.sa_check_role_status == 2"
                  type="primary"
                  @click="toggleMemberDialog"
                  >选择店员</el-button
                >
                <el-table
                  v-if="
                    this.staffList[0].staff.length > 0 &&
                      formData.sa_check_role_status == 2
                  "
                  :data="staffList"
                  class="staff"
                  max-height="500"
                >
                  <el-table-column prop="staff" label="群主">
                    <template slot-scope="scope">
                      <div
                        class="item-staff"
                        v-for="(item, index) in scope.row.staff"
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
                  <el-table-column fixed="right" label="操作" width="72">
                    <template>
                      <el-button
                        @click.native.prevent="toggleMemberDialog"
                        type="text"
                      >
                        修改
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </cy-form-item>
              <cy-form-item
                label="签名密钥"
                key="popup_store_keys"
                prop="popup_store_keys"
              >
                <cy-input
                  :maxlength="150"
                  show-word-limit
                  v-model.trim="formData.popup_store_keys"
                  size="medium"
                ></cy-input>
              </cy-form-item>
            </template>
          </template>
        </cy-form>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button type="primary" @click="confirm">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/business";
// 选择员工弹窗
import SelectMember from "./select-member/index.vue";

@Component({
  components: {
    SelectMember
  }
})
export default class extends Vue {
  @Prop() storeInfo!: {};
  /*** data ***/
  /** 选择员工弹窗显隐 */
  private showSelectStaff = false;
  /** 页面展示员工信息 */
  private staffList: Array<any> = [
    {
      staff: []
    }
  ];
  private formData = {
    store_type: "",
    store_name: "",
    store_code: "",
    address: "",
    wechat_department_id: "",
    popup_store_status: 0, //是否启用popup store 0 否 1 是
    sa_share_order_url: "", //SA分享订单链接
    sa_check_order_url: "", //店长查看门店订单链接
    popup_store_keys: "", //popup store 秘钥
    sa_check_role_status: 1, //查看门店全部订单权限 1 仅店长 2店长和指定店员
    sa_check_role: [], //查看门店全部订单权限的导购企微userid
    /** 邮箱 */
    email: "",
    /** 门店库存名称 */
    inventory_name: "",
    /** 门店所属城市 */
    city: "",
    /** 门店电话 */
    phone: ""
  };
  private epNameDetail = localStorage.getItem("epNameDetail");
  private rules = {
    store_name: [
      { required: true, message: "请输入门店名称", trigger: ["blur", "change"] }
    ],
    store_type: [
      { required: true, message: "请输入门店类型", trigger: ["blur", "change"] }
    ],
    email: [{ required: false, message: "请输入门店邮箱", trigger: "blur" }],
    address: [
      { required: true, message: "请输入门店地址", trigger: ["blur", "change"] }
    ],
    sa_share_order_url: [
      {
        required: true,
        message: "请填写SA分享订单链接",
        trigger: ["blur", "change"]
      }
    ],
    sa_check_order_url: [
      {
        required: true,
        message: "请填写店长查看门店链接",
        trigger: ["blur", "change"]
      }
    ],
    popup_store_keys: [
      { required: true, message: "请填写签名密钥", trigger: ["blur", "change"] }
    ],
    //电话验证包含座机或者手机
    phone: [
      {
        required: true,
        message: "请输入门店电话",
        trigger: ["blur", "change"]
      },
      {
        validator: (rule: any, value: string, callback: any) => {
          if (!value) {
            callback();
          } else {
            //只能输入数字
            const reg = /^[0-9]*$/;
            if (!reg.test(value)) {
              callback(new Error("请输入正确的电话号码"));
            } else {
              callback();
            }
          }
        },
        trigger: "blur"
      }
    ],
  };
  private is_edit = false;
  private epName = ''
  // *******************************方法*******************************
  private mounted() {
    this.is_edit = !!(this.storeInfo["id"] || 0);
    Object.keys(this.formData).forEach(key => {
      if (
        ![
          "popup_store_status",
          "sa_check_role_status",
          "sa_check_role"
        ].includes(key)
      ) {
        this.formData[key] = this.storeInfo[key] || "";
      }
    });

    //编辑
    if (this.is_edit) {
      this.formData.popup_store_status =
        this.storeInfo["popup_store_status"] == 0 ? 0 : 1;
      this.formData.sa_check_role_status =
        this.storeInfo["sa_check_role_status"] || 1;

      // 同时之前有选择员工，则反选
      if (this.storeInfo["sa_check_role"].length > 0) {
        this.formData.sa_check_role = this.storeInfo["sa_check_role"].split(
          ","
        );
        this.staffList[0].staff =
          this.storeInfo["sa_check_role_arr"].map(item => {
            return {
              name: item.name,
              user_id: item.userid
            };
          }) || [];
        console.log("mounted", this.staffList[0].staff);
      }
    }
    this.getCompanyConfig()
  }
  /**
   * 取消操作
   */
  private closeDialog() {
    this.$emit("closeDialog");
  }
  /**
   * 打开/关闭指定店员
   */
  private toggleMemberDialog() {
    this.showSelectStaff = true;
  }

  getCompanyConfig() {
    Api.getEpConfig().then(res => {
      // // 企业名称：coach;ks;sw
      const { name } = res.data;
      this.epName = name;
    });
  }

  /****选人组件回调 */

  private getMember(memberList: any) {
    this.showSelectStaff = false;
    if (memberList.length > 0) {
      const list: Array<any> = [];
      memberList.map((item: any) => {
        const value = {
          name: item.name,
          user_id: item.userid,
          is_dp: false
        };
        list.push(value);
      });
      this.staffList[0].staff = list.map(item => {
        return {
          user_id: item.user_id,
          name: item.name
        };
      });
      this.formData.sa_check_role = list.map((item: any) => item.user_id);
      console.log(
        "getMember",
        this.staffList[0].staff,
        this.formData.sa_check_role
      );
    } else {
      this.staffList[0].staff = [];
      this.formData.sa_check_role = [];
    }
  }
  /****
   * 表单校验
   */
  private confirm() {
    (this.$refs["cform"] as HTMLFormElement).validate((valid: boolean) => {
      if (valid) {
        console.log(this.formData);
        //开启popup store 且指定员工店长，但是没选员工
        if (
          this.formData.popup_store_status == 1 &&
          this.formData.sa_check_role_status == 2 &&
          this.formData.sa_check_role.length == 0
        ) {
          this.$message.warning("请选择店员");
          return;
        }
        const params = {
          id: this.storeInfo["id"] || "",
          ...this.formData
        };
        console.log(params);
        //未开启popupstore 或者 非2店长和指定店员，不传递sa_check_role
        if (
          !(
            this.formData.popup_store_status == 1 &&
            this.formData.sa_check_role_status == 2
          )
        ) {
          delete params.sa_check_role;
        }
        let text = "更新成功";
        if (!this.is_edit) {
          delete params.id;
          text = "新增成功";
        }
        Api.editStoreApi(params).then(() => {
          this.$message.success(text);
          this.$emit("success");
        });
      } else {
        return;
      }
    });
  }
}
</script>

<style lang="scss" scoped>
.select-btn {
  margin: 20px 0;
}
.staff {
  width: 100%;
  margin-top: 10px;
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
}

.c-blue {
  color: #708df5;
  font-size: 12px;
  cursor: pointer;
}
/deep/ .el-table__body-wrapper {
  max-height: 300px !important;
  overflow-y: scroll;
}
/deep/ .table-empty-label img {
  position: relative;
  top: 40px;
}
.txt-line-height {
  line-height: 40px;
}
.item-mb {
  margin-bottom: 15px !important;
}
::v-deep {
  .el-dialog__body {
    max-height: 457px;
    overflow: auto;
  }
}
</style>
