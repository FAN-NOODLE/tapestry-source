<template>
  <div class="role-block">
    <layout-form title_name="查看角色">
      <template slot="content">
        <cy-form
          label-width="104px"
          :rules="role_rules"
        >
          <cy-form-item
            label="创建时间"
            prop=""
          >
            {{ role.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </cy-form-item>
          <cy-form-item
            label="更新时间"
            prop=""
          >
            {{ role.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </cy-form-item>
          <cy-form-item
            label="创建人"
            prop=""
          >
            {{ role.admin_name }}
          </cy-form-item>
          <cy-form-item
            label="关联员工数量"
            :span="8"
            prop=""
          >
            <span v-if="role.admin_relation_count === 0">{{ role.admin_relation_count }}
            </span>
            <el-button
              v-else
              size="middle"
              type="text"
              @click="newPage"
            >{{ role.admin_relation_count }}
            </el-button>
          </cy-form-item>
          <cy-form-item
            label="角色名称"
            prop=""
          >{{ role.name }}</cy-form-item>
          <cy-form-item
            label="角色适用范围"
            prop=""
          >{{ roleAppName(role.app_code) }}</cy-form-item>
          <cy-form-item
            label="权限配置"
            widthSize="100%"
            prop=""
          >
            <el-collapse
              class="mar-b-20"
              v-model="active_name"
            >
              <el-collapse-item name="1">
                <template slot="title">
                  <el-checkbox
                    class="m-l-20 m-r-20"
                    disabled="disabled"
                    v-model="check_alls"
                  >全部选择</el-checkbox>
                </template>
                <el-checkbox-group
                  v-model="checked_shop_powers"
                  disabled="disabled"
                >
                  <template v-for="item in power_list">
                    <el-checkbox
                      class="parent-power"
                      v-if="item.type === 1"
                      :label="item.id"
                      :key="item.id"
                    >
                      {{ item.name }}
                    </el-checkbox>
                    <template v-if="item.children">
                      <template v-for="items in item.children">
                        <el-checkbox
                          class="no-child"
                          v-if="items.type === 1 || items.type === 4"
                          :label="items.id"
                          :key="items.id"
                        >{{ items.name }}
                        </el-checkbox>
                        <template v-if="items.children">
                          <template v-for="itemss in items.children">
                            <el-checkbox
                              class="last-child"
                              v-if="itemss.type === 1 || itemss.type === 4"
                              :label="itemss.id"
                              :key="itemss.id"
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
        <div class="mar-l-60">
          <el-button
            size="small"
            @click="goBack"
          >返回 </el-button>
          <el-button
            type="primary"
            size="small"
            @click="edit"
            v-if="!default_roles.includes(role.mark)"
          >编辑
          </el-button>
          <el-button
            size="small"
            @click="deleteRole"
            v-if="
              !default_roles.includes(role.mark) &&
                !console_default_roles.includes(role.mark)
            "
          >删除
          </el-button>
        </div>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import Api from "@/api/setting/Role";
import { Vue, Component } from "vue-property-decorator";
import { roleAppName, WEB_ADMIN_APP_CODE } from "../app-options";
import BaseData from "../base-data";
@Component({
  components: {},
  mixins: [BaseData],
})
export default class extends Vue {
  private id = "";
  private identity_type = "";
  private scope = 1;
  private type = this.$route.query.type;
  private role = {
    name: "",
    describe: "",
    app_code: WEB_ADMIN_APP_CODE,
  };
  private role_rules = {};
  private power_list = []; // 权限列表
  private checked_crs_powers = [];
  private checked_shop_powers = []; // 选中的权限列表
  private check_alls = false; //是否全部选择
  private crs_check_all = false; //是否选中全部crs
  private shop_check_all = false; // 是否全部选中shop
  private active_name = ["1"]; //默认展开
  private role_id = null; //角色编辑
  private serve_type = null; //当前的服务类型
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
  private goBack() {
    this.$router.go(-1);
  }
  // 删除弹框提示
  private deleteRole() {
    this.$confirm("确认要删除该角色吗? 删除后将无法恢复", "删除角色", {
      confirmButtonText: "确认删除",
      cancelButtonText: "我再想想",
      type: "warning",
    })
      .then(() => {
        const obj = {
          role_id: this.role_id,
        };
        Api.deleteRole(obj).then(() => {
          this.$message({
            type: "success",
            message: "删除成功",
          });
          this.$router.push({
            name: "myapp_setting_accounts_roles_list",
          });
        });
      })
      .catch();
  }
  private edit() {
    this.$router.push({
      name: "myapp_setting_accounts_roles_role-edit",
      query: {
        id: this.$route.query.id,
        type: "2",
        bread: "编辑角色",
      },
    });
  }
  private newPage() {
    const route_url = this.$router.resolve({
      name: "myapp_setting_accounts_accounts-news_list",
      query: {
        role: this.role.name,
      },
    });
    window.open(route_url.href, "_blank");
  }
  //获取角色详情
  private roleAppName = roleAppName;
  private async getRoleInfo() {
    const { data } = await Api.roleInfo({ role_id: this.role_id });
    this.role = {
      ...data,
      app_code: data.app_code || WEB_ADMIN_APP_CODE,
    };
    this.checked_shop_powers = data.permissions_id_list;
  }
  private async getConfigurablePermissions() {
    const { data } = await Api.configurablePermissions(this.role.app_code);
    this.power_list = data;
  }
  created() {
    const { id = "" } = this.$route.query;
    id && ((this as any).role_id = Number(id));
  }
  async mounted() {
    await this.getRoleInfo();
    await this.getConfigurablePermissions();
  }
}
</script>
<style lang="scss" scoped>
.role-block {
  /deep/ {
    .label {
      padding-left: 0 !important;
    }
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
    .el-form-item {
      margin-bottom: 9px !important;
    }
    .el-form-item__content {
      margin-left: 108px !important;
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
</style>
