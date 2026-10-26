<template>
  <div class="content">
    <choose-range
      :show="is_show_select_range"
      ref="chooseRange"
      :defaultList="showList"
      :on-update="onRangeUpdate"
      :on-close="
        () => {
          is_show_select_range = false;
        }
      "
    ></choose-range>
    <choose-role ref="select" :value="role_ids" @change="handleRoleChange"></choose-role>
    <layout-form :title_name="title_name">
      <template slot="content">
        <cy-form
          class="demo-ruleForm"
          ref="base"
          :model="formData"
          label-width="104px"
          size="medium"
          :rules="base_rules"
        >
          <cy-form-item label="账号名称" prop="account" required>
            <cy-input
              type="text"
              widthSize="medium"
              v-model.trim="formData.account"
              placeholder="请输入账号名称"
              :disabled="edit_type == 1"
              maxlength="20"
              show-word-limit="show-word-limit"
            ></cy-input>
          </cy-form-item>
          <cy-form-item label="账号手机号" prop="mobile" required>
            <cy-input
              v-if="edit_type == 1"
              type="text"
              widthSize="medium"
              v-model.trim="formData.newMobile"
              :disabled="edit_type == 1"
              placeholder="请输入账号手机号"
              maxlength="11"
              show-word-limit="show-word-limit"
            >
              <template slot="prepend">{{area}}</template>
            </cy-input>
            <cy-input
              v-else
              type="text"
              widthSize="medium"
              v-model.trim="formData.mobile"
              :disabled="edit_type == 1"
              placeholder="请输入账号手机号"
              maxlength="11"
              show-word-limit="show-word-limit"
            >
              <template slot="prepend">{{area}}</template>
            </cy-input>
          </cy-form-item>
          <cy-form-item label="账号姓名" prop="username" required>
            <cy-input
              type="text"
              widthSize="medium"
              v-model.trim="formData.username"
              placeholder="请输入真实姓名"
            ></cy-input>
          </cy-form-item>
          <div v-if="edit_type != 1">
            <cy-form-item label="账号密码" required>
              <p class="passDsc">
                创建时密码默认为账号+手机号后4位,登录时需修改密码
              </p>
            </cy-form-item>
          </div>
          <cy-form-item label="账号可用系统" class="p-l-0">
            <el-checkbox-group v-model="formData.admin_available_system">
              <el-checkbox
                :label="item.id"
                v-for="item in checkList"
                :key="item.id"
                >{{ item.name }}</el-checkbox
              >
            </el-checkbox-group>
          </cy-form-item>
          <cy-form-item label="账号角色" required prop="">
            <el-button type="primary" @click="handleSelectRole()" size="small">
              选择角色
            </el-button>
          </cy-form-item>
          <cy-form-item widthSize="100%">
            <el-table
              class="pdb"
              v-if="role_ids && role_ids.length"
              :data="[1]"
              :key="indexValue"
              :header-cell-style="{
                background: '#ebebf1',
                fontSize: '14px',
                color: '#606266',
                fontWeight: '700',
                padding: '7px 0 7px 0'
              }"
            >
              <el-table-column label="角色" min-width="320" prop="name"
                ><template>
                  <el-tag
                    v-for="(tag, index) in role_ids"
                    :key="index"
                    @close="handleDeleteRole(index, tag)"
                    closable
                    :type="tag.type"
                  >
                    {{ tag.name }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" min-width="120">
                <!-- <template slot-scope="scope">
                  <el-button
                    class="operation-error"
                    @click.native.prevent="handleDeleteRole(-1)"
                    type="text"
                    size="small"
                    :disabled="scope.row.status === 1"
                  >删除</el-button>
                </template> -->
              </el-table-column>
            </el-table>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <div>
          <el-button type="primary" size="small" @click="submitForm()">
            保存
          </el-button>
          <el-button size="small" @click="goBack()">取消</el-button>
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script>
import Api from "@/api/setting/Account";
import ChooseRange from "../ChooseRange";
import ChooseRole from "../ChooseRole";
import BaseData from "./base-data";
let count = 0;
const mobileRegexp = /(1[3456789]\d{9}|\d{8})/;
const areaRegexp = /(\+(?:1|44|81|91|33|49|61|7|82|65|86|852|853|886))/
const areaMobileRegexp = new RegExp(`(?:${areaRegexp.source})?${mobileRegexp.source}`)

export default {
  name: "AccountsEdit",
  components: {
    ChooseRange,
    ChooseRole
  },
  mixins: [BaseData],
  data() {
    const validatePass = (rule, value, callback) => {
      const reg = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{15,}$/;
      if (!value) {
        return callback(new Error("请输入登录密码"));
      }
      if (!reg.test(value)) {
        return callback(
          new Error(
            "密码至少包含15个字符，一个大写字母，一个小写字母和一个数字"
          )
        );
      }
      callback();
    };
    const validUsename = (rule, value, callback) => {
      const reg = /^[0-9a-zA-Z]*$/g;
      if (!value) {
        return callback(new Error("请输入账号名称"));
      }
      if (!reg.test(value)) {
        return callback(new Error("只允许输入字母+数字"));
      }
      callback();
    };
    return {
      id: this.$route.query.id,
      edit_type: this.$route.query.isEdit,
      title_name: "",
      showList: [], //门店已选择数据处理后的展示集合
      is_show_select_range: false, //是否选择账号范围
      is_show_select_role: false, //是否选择角色
      checkList: [],
      formData: {
        account: "", //名称
        mobile: "", //手机号
        email: "", //邮箱
        password: "", //密码
        username: "", // 真实姓名
        admin_available_system: [], //可用系统
        range_type: 1, //是否指定范围
        range: [] //范围
        // role_ids: [] //角色
      },
      role_ids: [],
      formDataString: "",
      base_rules: {
        account: [
          {
            validator: validUsename,
            trigger: ["blur", "change"]
          }
        ],
        mobile: [
          {
            validator: this.validateself,
            trigger: "blur"
          }
        ],
        password: [
          {
            validator: validatePass,
            trigger: ["blur", "change"]
          }
        ]
      },
      indexValue: Math.random()
    };
  },
  methods: {
    //删除范围
    validateself(rule, value, callback) {
      // eslint-disable-next-line no-useless-escape
      console.log(value);
      if (!areaMobileRegexp.test(value)) {
        return callback(new Error("手机号格式有误"));
      }
      callback();
    },
    handleDeleteFw(item, index) {
      if (item) {
        this.showList.splice(index, 1);
      } else if (item.text == "全部门店") {
        this.showList.forEach((i, dex) => {
          if (i.type != 3) {
            this.showList.splice(dex, 1);
          }
        });
      }
    },
    //删除角色
    handleDeleteRole(index, tag) {
      console.log(index, tag, "tag======>");
      this.indexValue = Math.random();
      if (index > -1) {
        if (this.inner_no_choice.includes(tag.mark)) {
          this.$alert("删除该角色后SA将无法使用WeClient Pro,请勿删除", "提示", {
            confirmButtonText: "确定",
            type: "warning",
            callback: action => {
              return;
            }
          });
        } else {
          this.role_ids.splice(index, 1);
          this.$forceUpdate();
        }
      } else {
        this.role_ids = [];
      }
      this.$forceUpdate();
    },
    //选择账号范围
    handleSelectRange() {
      this.$refs.chooseRange.init(this.showList);
      this.is_show_select_range = true;
    },
    // 递归查询当前节点下的门店数量
    findShopCount(data) {
      if (!data) {
        return 0;
      } else {
        data.forEach(item => {
          if (item.children && item.children.length > 0) {
            this.findShopCount(item.children);
          } else if (item.type == 2) {
            count += 1;
          }
        });
        return count;
      }
    },
    //账号范围回调
    onRangeUpdate(way, e) {
      const self = this;
      if (way == 1) {
        this.showList = e.map(item => {
          count = 0;
          return { ...item, ...{ count: self.findShopCount([item]) } };
        });
      }
      if (way == 2) {
        let hasMall = false;
        this.showList.forEach(item => {
          if (item.id == e.id) {
            hasMall = true;
          }
        });
        if (!hasMall) {
          this.showList.push(e);
        }
      }
    },
    //选择角色
    handleSelectRole() {
      this.$refs.select.handleShow()
    },
    //选择角色回调
    handleRoleChange(e) {
      this.indexValue = Math.random();
      this.role_ids = e;
    },
    // 获取账号详情
    getDetail() {
      Api.accountDetail({ id: Number(this.id) }).then(res => {
        this.formData = res.data;
        const reg = /^(\+(?:1|44|81|86|91|33|49|61|7|82|852|65|853|886))?(\d+)$/
        const [, area, mobile] = res.data.mobile.match(reg)|| [];
        
        const pat = /(\d{3})\d*(\d{4})/;
        this.formData.newMobile = mobile.replace(pat, "$1****$2");
        this.area = area||'+86'
        this.showList = res.data.range.map(item => {
          return {
            label: item.name,
            type: item.type,
            pid: item.pid,
            id: item.id,
            count: item.count
          };
        });
        this.role_ids = res.data.new_role_list;
        this.formDataString = JSON.stringify(this.formData);
      });
    },
    changeRange() {
      const tempData = [];
      this.showList.forEach(item => {
        tempData.push({
          object_id: item.id,
          pid: item.pid,
          type: item.type,
          count: item.countW
        });
      });
      return tempData;
    },
    // 添加账号
    addDep() {
      Api.accountAdd({
        account: this.formData.account, //名称
        mobile: `${this.formData.mobile}`, //手机号
        username: this.formData.username, //真实姓名
        email: this.formData.email, // 账号邮箱
        password: this.formData.password, // 账号密码
        admin_available_system: this.formData.admin_available_system, //可用系统
        range_type: this.formData.range_type, //是否指定范围
        range: this.changeRange(), //范围
        role_ids: this.role_ids.map(res => {
          return res.id;
        }) //角色
      }).then(() => {
        // 清除缓存
        localStorage.removeItem(this.SESSION_NAME);
        this.$router.go(-1);
        this.$message({
          showClose: true,
          message: "添加账号成功",
          type: "success"
        });
      });
    },

    // 编辑账号
    editDep() {
      Api.accountEdit({
        id: this.id,
        account: this.formData.account, //名称
        mobile: `${this.formData.mobile}`, //手机号
        username: this.formData.username, //名称
        email: this.formData.email, // 账号邮箱
        password: this.formData.password, // 账号密码
        admin_available_system: this.formData.admin_available_system, //可用系统
        range_type: this.formData.range_type, //是否指定范围
        range: this.changeRange(), //范围
        role_ids: this.role_ids.map(res => {
          return res.id;
        }) //角色
      }).then(() => {
        // 清除缓存
        localStorage.removeItem(this.SESSION_NAME);
        this.$router.go(-1);
        this.$message({
          showClose: true,
          message: "编辑账号成功",
          type: "success"
        });
      });
    },

    // 表单提交
    submitForm() {
      this.$refs["base"]["validate"](valid => {
        if (valid) {
          if (this.role_ids.length == 0) {
            this.$message({
              showClose: true,
              message: "请选择账号角色",
              type: "warning"
            });
            return;
          }
          this.edit_type === 1 ? this.editDep() : this.addDep();
        } else {
          return false;
        }
      });
    },
    //获取区域可用系统
    async accountSystem() {
      const res = await Api.accountSystem();
      this.checkList = res.data ? res.data.list : [];
    },
    goBack() {
      if (this.formDataString == JSON.stringify(this.formData)) {
        this.$router.go(-1);
      } else {
        this.$confirm(
          "你确认要放弃编辑吗？放弃后编辑的信息将不会被保存",
          "离开页面提示",
          {
            distinguishCancelAndClose: true,
            confirmButtonText: "确认",
            confirmButtonClass: "button-danger",
            cancelButtonText: "取消",
            type: "warning"
          }
        ).then(() => {
          this.$router.go(-1);
        });
      }
    }
  },
  async mounted() {
    await this.accountSystem();
    this.edit_type = Number(this.$route.query.is_edit);
    if (this.edit_type === 1) {
      this.title_name = "编辑账号";
      this.base_rules.password = {
        required: true,
        trigger: ["blur", "change"],
        message: "请输入登录密码"
      };
      this.getDetail();
    } else {
      this.title_name = "新增账号";
      this.setAndGetSesssion();
      this.formData.admin_available_system = this.checkList[0]
        ? [this.checkList[0].id]
        : [];
      this.formDataString = JSON.stringify(this.formData);
    }
  },
  beforeMount() {
    this.area = "+86";
  }
};
</script>
<style lang="scss" scoped>
.content {
  .p-l-0 {
    /deep/.label {
      padding-left: 0 !important;
    }
  }
  .pdb {
    margin-bottom: 20px;
  }
  .el-tag {
    margin-right: 10px;
    margin-bottom: 5px;
  }
  .operation-error {
    color: red;
  }
}
.passDsc {
  color: #dc1e13;
  line-height: 40px;
}
</style>
