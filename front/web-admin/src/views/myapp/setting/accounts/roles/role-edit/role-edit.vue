<template>
  <div class="role-block">
    <layout-form>
      <template slot="content">
        <cy-form
          :model="role"
          ref="role"
          :rules="role_rules"
          label-width="114px"
        >
          <cy-form-item
            label="角色名称"
            prop="name"
            required
          >
            <cy-input
              widthSize="medium"
              v-model.trim="role.name"
              placeholder="请输入角色名称"
              show-word-limit="show-word-limit"
              maxlength="20"
            >
            </cy-input>
          </cy-form-item>
          <cy-form-item
            label="角色适用范围"
            required
            prop="app_code"
          >
            <el-radio-group
              v-model="role.app_code"
              @change="handleAppCodeChange"
            >
              <el-radio
                v-for="item in app_options"
                :key="item.app_code"
                :label="item.app_code"
              >{{ item.name }}</el-radio>
            </el-radio-group>
            <div class="role-app-tip">每个角色仅可关联一个系统。</div>
          </cy-form-item>
          <cy-form-item
            label="权限配置"
            required
            prop=""
            widthSize="100%"
          >
            <el-collapse
              class="m-b-20"
              v-model="active_name"
            >
              <el-collapse-item name="1">
                <template slot="title">
                  <el-checkbox
                    class="m-l-20 m-r-20"
                    v-model="check_alls"
                    @change="ChangeAlls"
                  >全部选择</el-checkbox>
                </template>
                <el-checkbox-group v-model="checked_shop_powers">
                  <template v-for="item in power_list">
                    <el-checkbox
                      class="parent-power"
                      v-if="item.type === 1"
                      :label="item.id"
                      :key="item.id"
                      @change="isSelected(item)"
                    >{{ item.name }}
                    </el-checkbox>
                    <template v-if="item.children">
                      <template v-for="items in item.children">
                        <el-checkbox
                          class="no-child"
                          v-if="items.type === 1 || items.type === 4"
                          :label="items.id"
                          :key="items.id"
                          @change="isSelected(items, item)"
                        >{{ items.name }}</el-checkbox>
                        <template v-if="items.children">
                          <template v-for="itemss in items.children">
                            <el-checkbox
                              class="last-child"
                              v-if="itemss.type === 1 || itemss.type === 4"
                              :label="itemss.id"
                              :key="itemss.id"
                              @change="isSelected(itemss, items, item)"
                            >{{ itemss.name }}
                            </el-checkbox>
                          </template>
                        </template>
                      </template>
                    </template>
                  </template>
                </el-checkbox-group>
              </el-collapse-item>
            </el-collapse>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <div>
          <el-button
            type="primary"
            size="small"
            @click="saveData"
          >
            保存
          </el-button>
          <el-button
            size="small"
            @click="goBack"
          >取消</el-button>
        </div>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import Api from "@/api/setting/Role";
import { Vue, Component } from "vue-property-decorator";
import {
  ACCOUNT_SYSTEM_APP_CODES,
  ROLE_APP_NAMES,
  WEB_ADMIN_APP_CODE,
} from "../app-options";
@Component({
  components: {},
})
export default class extends Vue {
  private checkName = (rule: any, value: any, callback: any) => {
    // eslint-disable-next-line no-useless-escape
    const reg = /^((?!\?|\*|\<|\>|\/|\:|\||\"|\\).)*$/;
    if (!reg.test(value)) {
      return callback(new Error("角色名称不支持特殊字符，请重新输入"));
    }
    callback();
  };
  private title_name =
    this.$route.query.type == "1"
      ? "创建角色"
      : this.$route.query.type == "2"
      ? "编辑角色"
      : ""; //主标题
  private type = this.$route.query.type;
  private role = {
    name: "",
    describe: "",
    app_code: WEB_ADMIN_APP_CODE,
  };
  private roleString = "";
  private checked_shop_powers_string = "";
  private role_rules = {
    name: [{ validator: this.checkName }],
    app_code: [
      { required: true, message: "请选择角色适用范围", trigger: "change" },
    ],
  };
  private app_options: Array<{ app_code: string; name: string }> = [];
  private power_list = []; // 权限列表
  private checked_crs_powers = [];
  private checked_shop_powers: Array<any> = []; // 选中的权限列表
  private check_alls = false; //是否全部选择
  private crs_check_all = false; //是否选中全部crs
  private shop_check_all = false; // 是否全部选中shop
  private role_edit = 1; //判断是编辑还是修改
  private active_name = ["1"]; //默认展开
  private role_id = null; //角色编辑
  private application = []; //后端返回的开通的应用
  private crs_powers = []; //crs的所有权限列表
  private shop_powers = []; //shop的所有权限列表
  private modal_show = false; //是否打开选择门店弹框
  private store_list = []; //选中的门店列表
  private see_more = false; // 查看/隐藏更多选择的门店
  private tags_height = 0; //获取tags的高度
  private store_ids = []; //选中的门店id
  private is_nature = false; //是否开启自然流量
  private is_store = false; //是否开启选择门店
  // 返回
  goBack() {
    const arrayTemp =
      this.checked_shop_powers.length !== 0
        ? JSON.stringify(
            this.checked_shop_powers.sort(function (a, b) {
              return a - b;
            })
          )
        : "";
    if (
      this.roleString == JSON.stringify(this.role) &&
      this.checked_shop_powers_string == arrayTemp
    ) {
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
          type: "warning",
        }
      ).then(() => {
        this.$router.go(-1);
      });
    }
  }
  //父节点选中时选中他的子节点,未选中时去掉所有的子节点
  private addRemoveChild(arr: any, data: any, type: any) {
    if (data.constructor === Array) {
      for (let i = data.length; i--; ) {
        if (type === 1) arr.indexOf(data[i].id) > -1 || arr.push(data[i].id);
        if (type === 2)
          arr.indexOf(data[i].id) == -1 ||
            arr.splice(arr.indexOf(data[i].id), 1);
        if (data[i].children) {
          this.addRemoveChild(arr, data[i].children, type);
        }
      }
    } else if (data.children) {
      for (let i = data.children.length; i--; ) {
        if (type === 1)
          arr.indexOf(data.children[i].id) > -1 ||
            arr.push(data.children[i].id);
        if (type === 2)
          arr.indexOf(data.children[i].id) == -1 ||
            arr.splice(arr.indexOf(data.children[i].id), 1);
        if (data.children[i].children) {
          this.addRemoveChild(arr, data.children[i].children, type);
        }
      }
    }
  }

  //顶部全部选中,crs,商城选择框 判断他们对应的子节点是否选中,如果选中就选中父级
  private isAllCheck(data: any, arr: any) {
    const indexArr = [];
    let state = false;
    for (let i = data.length; i--; ) {
      const index = arr.includes(data[i].id);
      indexArr.push(index);
    }
    state = indexArr.includes(false) ? false : true;
    return state;
  }

  //子级全部选中将父级的key加入数组,反之去掉父级的key
  private isAllSelect(arr: any, data: any) {
    const indexArr = [];
    for (let i = data.children.length; i--; ) {
      const index = arr.includes(data.children[i].id);
      indexArr.push(index);
    }
    if (indexArr.includes(true)) {
      arr.indexOf(data.id) > -1 || arr.push(data.id);
    } else {
      arr.indexOf(data.id) == -1 || arr.splice(arr.indexOf(data.id), 1);
    }
  }

  //全部选中按钮勾选
  private ChangeAlls(val: any) {
    if (val) {
      this.addRemoveChild(this.checked_shop_powers, this.power_list, 1);
      this.crs_check_all = true;
      this.shop_check_all = true;
    } else {
      this.addRemoveChild(this.checked_shop_powers, this.power_list, 2);
      this.crs_check_all = false;
      this.shop_check_all = false;
    }
  }
  //节点勾选
  private isSelected(dataChilds: any, dataChild: any, data: any) {
    if ((this.checked_shop_powers as any).includes(dataChilds.id)) {
      this.addRemoveChild(this.checked_shop_powers, dataChilds, 1);
    } else {
      this.addRemoveChild(this.checked_shop_powers, dataChilds, 2);
    }
    if (dataChild) this.isAllSelect(this.checked_shop_powers, dataChild);
    if (data) this.isAllSelect(this.checked_shop_powers, data);
    this.check_alls = this.isAllCheck(
      this.power_list,
      this.checked_shop_powers
    );
    this.crs_check_all = this.isAllCheck(
      this.crs_powers,
      this.checked_shop_powers
    );
    this.shop_check_all = this.isAllCheck(
      this.shop_powers,
      this.checked_shop_powers
    );
  }

  // 新增/编辑角色
  private editRole() {
    let api = "addRole";
    const permissions = [...this.checked_shop_powers];
    if (this.role.app_code === WEB_ADMIN_APP_CODE) {
      permissions.push(1); // WeClient 管理后台保留首页权限
    }
    const checked_list = [...new Set(permissions)];
    const params = {
      name: this.role.name,
      type: 1,
      app_code: this.role.app_code,
      permission_id_list: checked_list,
    };
    if (this.role_edit == 2) {
      api = "editRole";
      params["role_id"] = this.role_id;
    }
    Api[api](params).then(() => {
      this.$router.go(-1);
      this.$message({
        showClose: true,
        message: this.role_edit === 2 ? "编辑角色成功" : "新增角色成功",
        type: "success",
      });
    });
  }

  // 表单提交
  private saveData() {
    (this.$refs["role"] as HTMLFormElement).validate((valid: any) => {
      if (valid) {
        if (this.checked_shop_powers.length <= 0) {
          this.$message({
            showClose: true,
            message: "请选择角色权限范围",
            type: "error",
          });
          return;
        }
        this.editRole();
      }
    });
  }
  //获取角色详情
  private async getRoleInfo() {
    const { data } = await Api.roleInfo({ role_id: this.role_id });
    this.role = {
      ...data,
      app_code: data.app_code || WEB_ADMIN_APP_CODE,
    };
    this.checked_shop_powers = data.permissions_id_list;
    this.roleString = JSON.stringify(this.role);
    this.checked_shop_powers_string = JSON.stringify(
      this.checked_shop_powers
    );
  }
  private async getConfigurablePermissions() {
    const { data } = await Api.configurablePermissions(this.role.app_code);
    this.power_list = data;
  }
  private async getRoleApplications() {
    const { data } = await Api.roleApplications();
    this.app_options = (data.list || [])
      .map((item: { id: number }) => {
        const app_code = ACCOUNT_SYSTEM_APP_CODES[item.id];
        return app_code
          ? { app_code, name: ROLE_APP_NAMES[app_code] }
          : null;
      })
      .filter(Boolean);
    if (this.app_options.length === 0) {
      this.app_options = [
        { app_code: WEB_ADMIN_APP_CODE, name: ROLE_APP_NAMES[WEB_ADMIN_APP_CODE] },
      ];
    }
  }
  private async handleAppCodeChange() {
    this.checked_shop_powers = [];
    this.check_alls = false;
    await this.getConfigurablePermissions();
  }
  created() {
    const { type, id } = this.$route.query;
    type && ((this as any).role_edit = type);
    id && ((this as any).role_id = Number(id));
  }
  async mounted() {
    await this.getRoleApplications();
    if (this.role_edit == 2) {
      await this.getRoleInfo();
    }
    await this.getConfigurablePermissions();
    this.roleString = JSON.stringify(this.role);
  }
}
</script>
<style lang="scss" scoped>
.role-block {
  /deep/ {
    .el-collapse {
      border: 1px solid #dcdfe6;
    }
    .el-collapse-item__header {
      border-bottom: 1px solid #dcdfe6;
      height: 55px;
      line-height: 55px;
    }
    .el-collapse-item__content {
      padding-bottom: 20px;
      margin-top: -20px;
    }
    .el-collapse-item__wrap {
      border-top: 1px solid #dcdfe6;
      border-bottom: 1px solid #dcdfe6;
    }
    .el-checkbox {
      + {
        .el-checkbox {
          margin-left: 0px;
        }
      }
    }
    .el-textarea__inner {
      padding-bottom: 20px;
    }
  }
  .no-padding {
    margin-bottom: 0px;
  }
  .parent-power {
    margin-top: 20px;
    background: rgba(232, 232, 232, 0.88);
    width: calc(100% - 40px);
    padding: 15px 20px;
  }
  .no-child {
    display: block;
    margin-left: 30px !important;
    margin-top: 20px;
  }
  .last-child {
    margin-left: 50px !important;
    margin-top: 20px;
  }
  .natural {
    margin-right: 10px;
  }
  .hover-span {
    font-size: 16px;
    cursor: pointer;
  }
  .has-stor {
    display: inline-block;
    margin-right: 0px;
  }
  .store-col {
    color: #0486fe;
  }
  .search-more {
    cursor: pointer;
    vertical-align: top;
  }
  .search-word {
    font-size: 14px;
  }
  .search-more-img {
    margin-left: 7px;
    width: 10px;
    height: 10px;
    vertical-align: middle;
  }
  .more {
    color: #409eff !important;
  }
  .no-more {
    margin-top: -15px;
  }
  .show-tag {
    display: inline-block;
  }
  .hide-tag {
    display: inline-block;
    height: 40px;
    overflow: hidden;
  }
  .number-v {
    width: calc(100% - 10px);
    text-align: right;
    margin-top: -33px;
    padding: 0 10px 0 0;
  }
  .role-app-tip {
    color: #909399;
    font-size: 12px;
    line-height: 20px;
  }
  .main-tag {
    background-color: rgba(64, 158, 255, 0.1);
    display: inline-block;
    padding: 0 10px;
    height: 32px;
    line-height: 30px;
    font-size: 12px;
    color: #409eff;
    border-radius: 4px;
    -webkit-box-sizing: border-box;
    box-sizing: border-box;
    border: 1px solid rgba(64, 158, 255, 0.2);
    white-space: nowrap;
    margin-right: 10px;
    /deep/ {
      .el-icon-close {
        border-radius: 50%;
        text-align: center;
        position: relative;
        cursor: pointer;
        font-size: 12px;
        height: 16px;
        width: 16px;
        line-height: 16px;
        vertical-align: middle;
        top: -1px;
        right: -5px;
        color: #409eff;
        &:hover {
          background-color: #409eff;
          color: #fff;
        }
      }
    }
    .tag-title {
      max-width: 95px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      display: inline-block;
      vertical-align: middle;
      cursor: pointer;
    }
  }
}
/deep/ {
  // .el-form-item__label {
  //   width: 110px !important;
  // }
  .el-button {
    + {
      .el-button {
        margin-left: 20px;
      }
    }
  }
  .el-form-item {
    margin-bottom: 9px !important;
  }
  .el-form-item__content {
    // margin-left: 118px !important;
  }
  .mar-l-60 {
    margin-left: 30px !important;
  }
}
</style>
