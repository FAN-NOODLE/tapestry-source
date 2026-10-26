<template>
  <div class="content">
    <layout-form title_name="查看账号">
      <template
        slot="content"
        v-if="detail"
      >
        <cy-form
          ref="cform"
          :rules="rules"
          :model="ruleForm"
          label-width="104px"
        >
          <cy-form-item label="创建时间">
            {{ detail.created_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </cy-form-item>
          <cy-form-item label="更新时间">
            {{ detail.updated_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </cy-form-item>
          <cy-form-item label="账号状态">
            {{ detail.status == 1 ? "启用中" : "已禁用" }}
          </cy-form-item>
          <cy-form-item label="账号名称">
            {{ detail.account }}
          </cy-form-item>
          <cy-form-item label="账号手机号">
            {{ detail.mobile }}
          </cy-form-item>
          <cy-form-item
            label="账号可用系统"
            v-if="
              detail.admin_available_system &&
                detail.admin_available_system.length > 0
            "
          >
            <div class="info">
              <div
                class="info-item"
                v-for="(item, index) in detail.admin_available_system"
                :key="index"
              >
                <span v-if="item == 1">管理后台</span><span v-if="item == 2">智慧导购小程序</span><span v-if="item == 3">天擎CDP</span>
              </div>
            </div>
          </cy-form-item>
          <cy-form-item
            label="账号角色"
            widthSize="100%"
          >
            <div class="info">
              <el-table
                class="pdb"
                :data="[{name: detail.new_role_list}]"
                size="mini"
                :header-cell-style="{
                  background: '#ebebf1',
                  fontSize: '14px',
                  color: '#606266',
                  fontWeight: '700',
                  padding: '7px 0 7px 0'
                }"
              >
                <el-table-column
                  label="角色"
                  min-width="320"
                  prop="name"
                ><template>
                    <div v-if="detail.new_role_list && detail.new_role_list.length">
                      <el-tag
                        v-for="tag in detail.new_role_list"
                        :key="tag.name"
                      >
                        {{ tag.name }}
                      </el-tag>
                    </div>
                    <div v-else>-</div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button
          size="small"
          @click="goBack"
        >返回 </el-button>
        <el-button
          type="primary"
          @click="handleEdit()"
          size="small"
          v-if="!detail.disable_edit"
        >编辑</el-button>
        <el-button
          @click="handleOpenOrClose(true)"
          size="small"
          v-if="detail.status === 2"
        >启用</el-button>
        <el-button
          @click="handleOpenOrClose()"
          v-if="!detail.disable_edit && detail.status === 1"
        >禁用</el-button>
        <el-button
          @click="handleDelete()"
          v-if="!detail.disable_edit && detail.status === 2"
        >删除</el-button>
      </template>
    </layout-form>
  </div>
</template>
<script>
import Api from "@/api/setting/Account";
// import LayoutForm from "@/components/layout/LayoutForm";
import { disOrOpenAccount, deleteAccount } from "./../utils";
import BaseData from "../base-data";
export default {
  name: "AccountsEdit",
  mixins: [BaseData],
  components: {},
  data() {
    return {
      id: this.$route.query.id, //当前账号id
      showList: [], //账号范围
      detail: {},
      ruleForm: {},
      rules: {}
    };
  },
  methods: {
    // 返回上一页
    goBack() {
      this.$router.go(-1);
    },
    // 获取账号详情
    getDetail() {
      Api.accountDetail({ id: this.id }).then(res => {
        this.detail = res.data;
        const pat = /(\d{3})\d*(\d{4})/;
        this.detail.mobile = this.detail.mobile.replace(pat, "$1****$2");
        const index = this.detail.new_role_list.findIndex(v =>
          this.no_choice.includes(v.mark)
        );
        if (index > -1) this.detail.disable_edit = true;
      });
    },
    // 编辑
    handleEdit() {
      this.$router.push({
        name: "myapp_setting_accounts_accounts-news_accounts-edit",
        query: {
          id: this.id,
          is_edit: 1,
          bread: "编辑账号"
        }
      });
    },
    // 启用/禁用
    handleOpenOrClose(e) {
      disOrOpenAccount(this.id, e).then(() => {
        this.getDetail();
      });
    },
    // 删除
    handleDelete() {
      deleteAccount(this.id).then(() => {
        this.goBack();
      });
    }
  },
  mounted() {
    this.getDetail();
  }
};
</script>
<style lang="scss" scoped>
.content {
  /deep/.label {
    padding-left: 0;
  }
}
.item {
  display: flex;
  width: 100%;
  margin: 20px 0;
  font-size: 14px;
}
.item .title {
  margin-right: 10px;
  min-width: 90px;
  max-width: 90px;
  text-align: left;
}
.item .info {
  width: 100%;
  display: flex;
}
.item .info .info-item {
  margin-right: 20px;
}

/deep/ .content-button {
  margin-left: 0 !important;
}

.el-tag {
  margin-right: 10px;
  margin-bottom: 5px;
  background: #f0f3fe !important;
  border: #d4dcfc 1px solid !important;
}

/deep/ .el-table__header {
  height: 20px;
}

/deep/ .content-block {
  padding: 10px 30px 30px 30px !important;
}
</style>
