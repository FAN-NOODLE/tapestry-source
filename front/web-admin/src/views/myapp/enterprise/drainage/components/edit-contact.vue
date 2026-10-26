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
            label-width="100px"
          >
            <el-form-item label="二维码名称" prop="code_name" required>
              <el-input
                v-model.trim="info.code_name"
                clearable
                width-size="medium"
                maxlength="30"
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
                @click="goAddChannel"
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
            <el-form-item
              label="员工"
              class="m-b-0"
              prop="member_ids"
              :rules="[
                {
                  ...rules.member_ids[0],
                  data: info
                }
              ]"
            >
              <el-button
                v-if="this.staffList[0].staff.length == 0"
                size="small"
                type="primary"
                @click="openSelectStaff"
                >选择员工</el-button
              >
              <el-table
                v-if="this.staffList[0].staff.length > 0"
                :data="staffList"
                style="width: 100%"
                max-height="500"
              >
                <el-table-column prop="staff" label="员工">
                  <template slot-scope="scope">
                    <div class="staff-box">
                      <div
                        class="item-staff"
                        v-for="(item, index) in scope.row.staff"
                        :key="index"
                      >
                      
                        <!-- icon -> 部门 -->
                        <img
                          class="staff-icon"
                          v-if="item.is_dp"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-ic.png"
                          alt=""
                        />
                        <!-- icon -> 人 -->
                        <img
                          class="staff-icon"
                          v-else
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/staff_icon.png"
                          alt=""
                        />
                        <span>{{ item.name }}</span>
                      </div>
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
          </el-form>
        </div>
      </template>
      <template slot="footer" class="dialog-footer">
        <el-button type="primary" @click="save()" size="small">保 存</el-button>
        <el-button @click="cancel()" size="small">取 消</el-button>
      </template>
    </layout-form>
    <!-- 选择员工弹窗-->
    <selectMember
      v-if="showSelectStaff"
      :hasSelectArr="this.staffList[0].staff"
      :limiteNum="100"
      dialogTitle="选择员工"
      @submit="getMember"
      @close="showSelectStaff = false"
    />
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import API from "@/api/enterprise/drainage";
import SelectMember from "./select-member.vue";
@Component({
  components: {
    SelectMember
  }
})
export default class extends Vue {
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
    code_name: [
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
    code_name: "",
    channel_id: "",
    auto_friend: 1,
    member_ids: [],
    user_ids: [],
    third_ids: []
  };
  private showSelectStaff = false; // 控制选择员工弹框显隐
  // 选人组件回调
  private getMember(memberList: Array<any>) {
    this.showSelectStaff = false;
    console.log("组件传递回来的数据-> ", memberList);
    if (memberList.length > 0) {
      // if (this.staffList[0].staff.length == 0) {
      //   this.$message.success("添加成功");
      // } else {
      //   this.$message.success("更新成功");
      // }
      const list: Array<any> = [];
      memberList.map((item: any) => {
        const value = {
          name: item.name,
          user_id: item.user_id,
          member_id: item.member_id,
          third_id: item.third_id,
          is_dp: item.is_dp,
          value: item.value
        };
        list.push(value);
      });
      this.staffList[0].staff = list;
      this.info.member_ids = list.map((item: any) => item.member_id).filter(Boolean);
      this.info.user_ids = list.map((item: any) => item.user_id).filter(Boolean);
      this.info.third_ids = list.map((item: any) => item.third_id).filter(Boolean);


    } else {
      this.staffList[0].staff = [];
      this.info.member_ids = [];
      this.info.user_ids = [];
      this.info.third_ids = [];
    }
    (this.$refs["cform"] as any).validateField("member_ids");
  }
  // 打开选择员工弹窗
  private openSelectStaff() {
    this.showSelectStaff = true;
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
  private createContact() {
    const params = {
      ...this.info
    };
    API.contactCreateApi(params).then((res: any) => {
      if (res.code === 200) {
        localStorage.setItem("CONTACT_INFO", JSON.stringify({}));
        this.$message.success("保存成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 编辑联系我接口
  private editContact() {
    const params = {
      code_name: this.info.code_name,
      channel_id: this.info.channel_id,
      auto_friend: this.info.auto_friend,
      user_ids: this.info.user_ids,
      member_ids: this.info.member_ids,
      third_ids: this.info.third_ids,
      id: this.id
    };
    API.contactUpdateApi(params).then((res: any) => {
      if (res.code === 200) {
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
        
        const memberList: Array<any> = [];
        res.data.member_data?.map((item: any) => {
          const value = {
            name: item.member_name,
            user_id: item.user_id,
            member_id: item.member_id,
            value: item.user_id,
            third_id: 0,
            is_dp: false
          };
          memberList.push(value);
        });
        res.data.department_data?.map((item: any) => {
          const value = {
            name: item.department_name,
            user_id: 0,
            member_id: 0,
            third_id: item.third_id,
            value: item.third_id,
            is_dp: true
          };
          memberList.push(value);
        })


        this.staffList[0].staff = memberList;
        this.info.member_ids = memberList.map((item: any) => item.member_id).filter(Boolean);
        this.info.user_ids = memberList.map((item: any) => item.user_id).filter(Boolean);
        this.info.third_ids = memberList.map((item: any) => item.third_id).filter(Boolean);
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
      .item-staff {
        display: inline-flex;
        align-items: center;
        height: 36px;
        border: 1px solid #dcdfe6;
        background-color: #f7f7f7;
        padding: 0 20px;
        margin-right: 10px;
        margin-bottom: 10px;
        // min-width: 87px;
        justify-content: center;
        .staff-icon {
          width: 14px;
          height: 14px;
          margin-right: 10px;
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
</style>
