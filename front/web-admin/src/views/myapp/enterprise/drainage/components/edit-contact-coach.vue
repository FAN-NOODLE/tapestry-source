<template>
  <div class="edit-contact-me">
    <layout-form>
      <template slot="content">
        <div class="content-box">
          <el-form
            class="main-box"
            ref="cform"
            :rules="rules"
            :model="info"
            label-width="120px"
          >
            <el-form-item label="二维码名称" prop="name" :rules="[{ required: true, message: '请输入二维码名称', trigger: 'blur' }]">
              <el-input
                v-model.trim="info.name"
                clearable
                width-size="medium"
                maxlength="20"
                placeholder="请输入二维码名称"
                show-word-limit
                style="width: 400px;"
              ></el-input>
            </el-form-item>
            <el-form-item
              label="渠道"
              prop="channel_id"
              required
              trigger="change"
              style="max-width:700px"
            >
              <el-select
                v-model="info.channel_id"
                width-size="medium"
                size="medium"
                search-input
                placeholder="请选择渠道"
                :popperAppendToBody="true"
                v-if="showChannelSelect"
                style="width: 400px;"
              >
                <el-option
                  :label="item.name"
                  :value="item.id"
                  :key="item.id"
                  v-for="item in channelList"
                ></el-option>
              </el-select>
              <el-button
                class="add-channel-btn"
                type="primary"
                size="small"
                @click="createFn"
                >新增渠道</el-button
              >
            </el-form-item>
            <el-form-item label="自动通过" prop="auto_friend">
              <el-switch
                v-model="info.auto_friend"
                :active-text="info.auto_friend == 1 ? '已开启' : '已关闭'"
                :active-value="1"
                :inactive-value="0"
              ></el-switch>
              <div class="switch-tip">
                开启后，客户添加时无需经过确认自动成为好友
              </div>
            </el-form-item>
            <el-form-item label="二维码类型" prop="auto_friend">
              <el-radio-group v-model="info.type" :disabled="$route.query.id">
                <el-radio :label="1">员工码+门店码</el-radio>
                <el-radio :label="2">多人共享码</el-radio>
              </el-radio-group>
              <div v-if="info.type == 1" class="switch-tip">
                根据所选的员工，每人生成一张专属的联系我二维码，同时将所有人放在一起生成一张联系我二维码
              </div>
              <div v-else class="switch-tip">
                将所选的所有员工放在一个联系我二维码中，生成一张添加二维码。
              </div>
            </el-form-item>
            <el-form-item
              label="员工"
              prop="member_ids"
              :rules="[
                {
                  ...rules.member_ids[0],
                  data: info
                }
              ]"
            >
              <el-button
                v-if="selectMember().length == 0"
                size="small"
                type="primary"
                @click="openSelectStaff"
                >选择员工</el-button
              >
              <el-table
                v-if="selectMember().length > 0"
                :data="[{staff: selectMember()}]"
                style="width: 100%"
                max-height="500"
              >
                <el-table-column prop="staff" label="员工">
                  <template slot-scope="scope">
                    <div class="staff-box">
                      <el-tag
                        class="item-staff"
                        v-for="(item, index) in scope.row.staff"
                        :key="index"
                      >
                        <img
                          class="staff-icon"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                          alt=""
                          v-if="item.is_dp"
                        />
                        <i v-else class="el-icon-user-solid"></i>
                        <span>{{ item.label }}</span>
                      </el-tag>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column fixed="right" label="操作" width="72">
                  <template>
                    <el-button
                      @click.native.prevent="openSelectStaff"
                      type="text"
                    >
                      修改
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-form-item>

            <el-form-item label="是否配置欢迎语" prop="auto_friend">
              <el-switch
                v-model="info.has_welcome"
                :active-value="1"
                :inactive-value="0"
              ></el-switch>
              <div class="switch-tip">
                配置欢迎语后，通过该联系我码添加的客户，将自动发送配置的欢迎语
              </div>
            </el-form-item>
            <el-form-item>
              <customer-welcome-form v-if="info.has_welcome" ref="welcomeFormRef" :value="info.welcome_list"></customer-welcome-form>
            </el-form-item>
          </el-form>
        </div>
      </template>
      <template slot="footer" class="dialog-footer">
        <el-button type="primary" @click="save()" size="small">保 存</el-button>
        <el-button @click="cancel()" size="small">取 消</el-button>
      </template>
    </layout-form>
    <select-member      
      :show.sync="showSelectMember"
      :value="selectMember()"
      @cancel="showSelectMember = false"
      @confirm="handleSelectMemberConfirm"
    />

    <create-or-edit-channel-dialog
      v-if="showCreateOrEditDialog"
      :dialogType="'1'"
      @channelDialog="channelDialogFn"
    ></create-or-edit-channel-dialog>
    <!-- 员工未实名弹框 -->
     <el-dialog
      title="温馨提示"
      :visible.sync="showUnverifiedUsersDialog"
      width="450px"
      center
    >
      <div>您选择的员工范围中有{{unverified_users.length}}位导购未进行企业实名认证，相关人员如下：</div>
      <div class="users">
        <div class="user" v-for="(item, index) in unverified_users" :key="index">
          {{ item.name }}({{ item.userid }})
        </div>
      </div>
      <div slot="footer" class="footer">请联系相关导购进行实名认证或从活动范围中移除未实名导购~</div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import API from "@/api/enterprise/drainage-coach";
import SelectMember from "./select-member-coach.vue";
import CustomerWelcomeForm from "./customer-welcome-form.vue";
import createOrEditChannelDialog from "../../components/create-or-eidt-channel-dialog.vue";
@Component({
  components: {
    SelectMember,
    CustomerWelcomeForm,
    createOrEditChannelDialog
  }
})
export default class extends Vue {
  private unverified_users: Array<any> = [];
  private showUnverifiedUsersDialog = false; // 控制员工未实名弹框显隐
  private is_edit = 0; // 是否是编辑
  private id = 0; // 联系我详情id
  private staffList: Array<any> = [
    // 选择的员工列表
    {
      staff: []
    }
  ];
  private channelList: Array<any> = []; // 渠道下拉列表
  private showChannelSelect = false; // 控制选择渠道下拉框显隐
  private rules = {
    name: [
      {
        required: true,
        message: "请输入渠道名称",
        trigger: "blur"
      }
    ],
    channel_id: [
      {
        required: true,
        trigger: "change",
        message: "请选择渠道"
      }
    ],
    member_ids: [
      {
        required: true,
        trigger: "change",
        validator: (rule: any, value: any, callback: any) => {
          console.log(this.info);
          
          if (value.length === 0 && rule.data.third_ids.length === 0) {
            callback(new Error("请选择员工或部门"));
          } else {
            callback();
          }
        }
      }
    ],
  };
  private infoCopy: any;
  private info: any = {
    name: "",
    channel_id: "",
    type: 1,
    auto_friend: 1,
    has_welcome: 0,
    welcome_list: [],
    member_ids: [],
    user_ids: [],
    third_ids: [],
    department_data: [],
    member_data: [],
  };
  private showSelectStaff = false; // 控制选择员工弹框显隐
  private showSelectMember = false; // 控制选择员工弹框显隐
  private showCreateOrEditDialog = false;

  createFn() {
    // console.log()
    this.showCreateOrEditDialog = true;
  }
  channelDialogFn() {
    this.showCreateOrEditDialog = false;
    this.getChannelList();
  }

  get selectMember() {
    return () => {
      const { member_data = [], department_data = [] } = this.info;
      return [
        ...(department_data || []).map(item => {
          return {
            label: item.department_name,
            value: item.third_id,
            is_dp: 1,
          }
        }),
        ...(member_data || []).map(item => {
          return {
            label: item.member_name,
            value: item.member_id,
            is_dp: 0,
          }
        }) 
      ]
    }
  }
  handleSelectMemberConfirm(e: Array<any>) {
    this.info.member_ids = [];
    this.info.third_ids = [];
    this.info.department_data = [];
    this.info.member_data = [];

    e.forEach(item => {
      if(item.is_dp) {
        this.info.third_ids.push(item.value);
        this.info.department_data.push({
          third_id: item.value,
          department_name: item.label,
        })
      } else {
        this.info.member_ids.push(item.value);
        this.info.member_data.push({
          member_id: item.value,
          member_name: item.label,
        })
      }
    })

    this.showSelectMember = false;
    (this.$refs["cform"] as any).validateField("member_ids");
  }
  // 打开选择员工弹窗
  private openSelectStaff() {
    this.showSelectMember = true;
  }
  // 打开新增渠道页面
  private goAddChannel() {
    const newLink = this.$router.resolve({
      name: "myapp_enterprise_drainage_channel_list",
      query: {
        showAddModal: "1"
      }
    });
    window.open(newLink.href, "_blank");
  }
  // 获取渠道下拉列表数据
  private getChannelList() {
    const params = {
      page: 1,
      size: 9999
    };
    API.channelsListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.channelList = res.data.list || [];
      } else {
        this.$message.error(res.message);
      }
      this.showChannelSelect = true;
    });
  }
  // 新增联系我接口
  async createContact() {
    const params = {
      ...this.info
    };
    if(params.has_welcome) {
      params.welcome_list = await (this.$refs.welcomeFormRef as any).save();
    } else {
      params.welcome_list = []
    }
    const res: any = await API.contactCreateApi(params);
    if (res.code === 200) {
      if(res.data.unverified_users && res.data.unverified_users.length > 0) {
        // this.$message.warning("请先验证员工");
        this.showUnverifiedUsersDialog = true;
        this.unverified_users = res.data.unverified_users || [];
        return;
      }
      localStorage.setItem("CONTACT_INFO", JSON.stringify({}));
      this.$message.success("保存成功");
      this.$router.go(-1);
    } else {
      this.$message.error(res.message);
    }
  }
  // 编辑联系我接口
  private async editContact() {
    const params = {
      ...this.info
    };

    if(params.has_welcome) {
      params.welcome_list = await (this.$refs.welcomeFormRef as any).save();
    } else {
      params.welcome_list = []
    }
    API.contactUpdateApi(params).then((res: any) => {
      if (res.code === 200) {

        if(res.data.unverified_users && res.data.unverified_users.length > 0) {
          // this.$message.warning("请先验证员工");
          this.showUnverifiedUsersDialog = true;
          this.unverified_users = res.data.unverified_users || [];
          return;
        }
        this.$message.success("保存成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取联系我详情
  private getDeatail() {
    const params = {
      id: this.id
    };
    API.contactDetailsApi(params).then((res: any) => {
      if (res.code === 200) {
        
        // const memberList: Array<any> = [];
        res.data.member_data = res.data.member_data || [];
        res.data.department_data = res.data.department_data || [];
        res.data.welcome_list = res.data.welcome_list || [];
        this.info.member_ids = (res.data.member_data||[]).map((item: any) => item.member_id)
        this.info.third_ids = (res.data.department_data||[]).map((item: any) => item.third_id)
        this.info = { ...this.info, ...res.data };
        this.infoCopy = JSON.parse(JSON.stringify(this.info));
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 点击保存按钮
  private save() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (!this.is_edit) {
          // 新增联系我
          this.createContact();
        } else {
          // 编辑联系我
          this.editContact();
        }
      } else {
        return;
      }
    });
  }
  // 点击取消按钮
  private cancel() {
    console.log(JSON.stringify(this.infoCopy) == JSON.stringify(this.info));
    if (JSON.stringify(this.infoCopy) == JSON.stringify(this.info)) {
      this.$router.go(-1);
    } else if (this.is_edit) {
      // 编辑时点击取消
      this.$confirm(
        "你确认要放弃编辑吗？放弃后编辑的信息将不会被保存？",
        "离开页面提示",
        {
          confirmButtonText: "确认",
          cancelButtonText: "取消",
          type: "warning"
        }
      )
        .then(() => {
          this.$router.go(-1);
        })
        .catch(() => {
          console.log("已取消");
        });
    } else {
      // 新增时点击取消,检测到有修改
      this.$confirm(
        "监测到未保存的内容，是否保存至下次新增继续编辑？",
        "离开页面提示",
        {
          confirmButtonText: "保存并离开",
          cancelButtonText: "直接离开",
          distinguishCancelAndClose: true,
          type: "warning"
        }
      )
        .then(() => {
          //保存并离开,存储到localStorage
          this.$nextTick(() => {
            localStorage.setItem(
              "CONTACT_INFO",
              JSON.stringify({
                info: this.info,
                staff: this.staffList[0].staff
              })
            );
            this.$router.go(-1);
          });
        })
        .catch(action => {
          if (action === "cancel") {
            localStorage.setItem("CONTACT_INFO", JSON.stringify({}));
            this.$router.go(-1);
          }
        });
    }
  }
  private created() {
    const { is_edit, id } = this.$route.query;
    this.is_edit = Number(is_edit) || 0;
    this.id = Number(id) || 0;
    // 新增时，拿取上次新增并已缓存的内容
    const reEdit = localStorage.getItem("CONTACT_INFO") || "{}";
    if (this.is_edit == 0 && reEdit && reEdit != "{}") {
      const obj = JSON.parse(reEdit);
      this.info = obj.info;
      this.staffList[0].staff = obj.staff;
    }
    this.infoCopy = JSON.parse(JSON.stringify(this.info));
    if (this.id) {
      this.getDeatail();
    }
    this.getChannelList();
  }
}
</script>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}
/deep/.gutter {
  width: 0px !important;
  display: none;
}
::v-deep .slot-content {
  background-color: #ffffff;
  border-top: none;
  padding-top: 0px;
}
.edit-contact-me {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  .content-box {
    display: flex;
    .main-box {
      flex: 1;
    }
    /deep/.in-item {
      min-width: 500px;
      max-width: calc(100% - 200px);
      width: auto;
    }
    .add-channel-btn {
      margin-left: 10px;
    }
    .staff-box {
      max-height: 125px;
      overflow: auto;
      display: flex;
      flex-wrap: wrap;
      gap: 10px;

      .item-staff {
        border: 1px solid #dcdfe6;
        background-color: #f7f7f7;
        padding: 0 20px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: fit-content;
        .staff-icon {
          width: 12px;
          height: 12px;
          margin-right: 5px;
          display: block;
          object-fit: contain;
        }

        .el-icon-user-solid {
          margin-right: 5px;
        }
      }
    }
    .switch-tip {
      font-size: 12px;
      color: #909399;
      line-height: 12px;
      margin-top: 2px;
    }
  }
}
.users {
  padding: 20px 0;
  color: #4D7CEC;
  line-height: 24px;
}

.footer {
  text-align: left;
}
/deep/.el-dialog__title {
  color: #733818;
}
</style>
