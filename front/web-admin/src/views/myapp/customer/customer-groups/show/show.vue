<template>
  <div class="group-show">
    <layout-form type="normal">
      <template slot="content">
        <dl class="form-layout title">
          <dt>分组名称</dt>
          <dd>
            {{ groupName }}
          </dd>
        </dl>
        <dl class="form-layout title">
          <dt>分组维度</dt>
          <dd>
            <span v-if="user_group_type === 1">标签</span>
            <span v-if="user_group_type === 2">顾客名单</span>
            <span v-if="user_group_type === 3">门店</span>
          </dd>
        </dl>
        <dl class="form-layout title" v-if="user_group_type === 2">
          <dt>名单导入详情</dt>
          <dd v-if="info.is_extract === 1">
            <el-button type="text" @click="downExport"
              >下载名单导入详情</el-button
            >
          </dd>
          <dd v-if="info.is_extract !== 1">
            匹配中
          </dd>
        </dl>
        <dl class="form-layout title store" v-if="user_group_type === 3">
          <dt>选择门店</dt>
          <dd>
            <el-table
              :data="info.Departments"
              ref="tableList"
              style="width:100%"
              :cell-style="{ height: '20px' }"
            >
              <el-table-column label="门店名称" prop="name"></el-table-column>
            </el-table>
          </dd>
        </dl>
        <dl
          class="form-layout title"
          v-if="selectInclude && selectInclude.length != 0"
        >
          <dt>包含标签</dt>
          <dd>
            <tags :tagsArr="selectInclude" type="include" isClose="no"></tags>
          </dd>
        </dl>
        <dl
          class="form-layout"
          v-if="selectExclude && selectExclude.length != 0"
        >
          <dt>排除标签</dt>
          <dd>
            <tags :tagsArr="selectExclude" type="exclude" isClose="no"></tags>
          </dd>
        </dl>
      </template>
      <template slot="footer">
        <el-button @click="goBackFn">返回</el-button>
        <el-button
          type="primary"
          @click="editFn"
          v-if="iscanEdit && info.user_group_type !== 2"
          >编辑</el-button
        >
        <el-button type="danger" @click="delFn" v-if="iscanEdit"
          >删除</el-button
        >
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import API from "@/api/customer";
import tags from "../components/tags.vue";
import CustomerModule from "@/store/modules/customer";
import { VUE_APP_API_ROOT } from "@/utils/env-config";
@Component({
  components: { tags }
})
export default class extends Vue {
  /*************************** data start ***************************/
  private groupName = ""; //分组名称
  private selectInclude: any = []; //包含标签
  private selectExclude: any = []; //排除标签
  private userCountArr: any = []; //组装数据传接口
  private curId = 0; //当前id
  private iscanEdit = true; //是否可编辑
  private user_group_type = 1;
  private info = {};

  /*************************** method start ***************************/
  private downExport() {
    const params = {
      user_group_id: this.curId
    };
    API.downListApi(params).then((res: any) => {
      const blob = new Blob([res], {
        type: "application/x-download;charset=UTF-8"
      });
      const fileName =
        "名单导入详情" + new Date().toLocaleString().substr(0, 10) + ".xlsx";
      //下载的文件名称及其后缀，后缀要和后台保持的一致

      if ("download" in document.createElement("a")) {
        // 非IE下载
        const elink = document.createElement("a");
        elink.download = fileName;
        elink.style.display = "none";
        elink.href = URL.createObjectURL(blob);
        document.body.appendChild(elink);
        elink.click();
        URL.revokeObjectURL(elink.href); // 释放URL 对象
        document.body.removeChild(elink);
      } else {
        // IE10+下载
        (navigator as any).msSaveBlob(blob, fileName);
      }
    });
  }
  //获取详情
  private getDetailFn() {
    API.groupDetailApi({ id: this.curId }).then(res => {
      const { data } = res;
      // 分组名称
      this.info = data;
      this.groupName = data.name;
      this.iscanEdit = data.can_edit;
      this.user_group_type = data.user_group_type;
      data.conditions.forEach((item: any) => {
        if (item.include_type == 1) {
          //包含
          this.selectInclude = item.tag_groups;
        } else {
          //非包含
          this.selectExclude = item.tag_groups;
        }
      });
    });
  }
  //删除
  private delFn() {
    this.$confirm("你确认要删除该分组吗？", "删除确认", {
      confirmButtonText: "确认",
      confirmButtonClass: "button-danger",
      cancelButtonText: "取消",
      type: "warning"
    })
      .then(() => {
        API.groupDeleteApi({ id: this.curId }).then(() => {
          this.$message.success("删除成功");
          this.goBackFn();
        });
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  //编辑
  private editFn() {
    CustomerModule.setGroupName("");
    CustomerModule.setGroupS([]);
    CustomerModule.setGroupU([]);
    this.$router.push({
      name: "myapp_customer_customer-groups_edit",
      query: { id: this.curId.toString() }
    });
  }
  //返回
  private goBackFn() {
    this.$router.push({
      name: "myapp_customer_customer-groups_list"
    });
  }
  /*************************** 生命周期 start ***************************/
  private mounted() {
    this.curId = Number(this.$route.query.id);
    this.getDetailFn();
  }
}
</script>
<style lang="scss" scoped>
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
.form-layout {
  display: flex;
  justify-content: center;
  align-items: center;
  padding-left: 10px;
  &.store {
    align-items: stretch;
  }
  dt {
    flex: 0 0 84px;
    margin-right: 20px;
    color: $color-text-normal;
  }
  dd {
    flex: 1;
    position: relative;
  }
}
.title {
  margin-bottom: 30px;
}
.tags {
  // margin-left: 0px;
  overflow: hidden;
  margin-bottom: 10px;
  .tags-single {
    float: left;
    overflow: hidden;
    margin-bottom: 20px;
    &-content {
      float: left;
      border-radius: 4px;
      border: 1px solid $color-second-line;
      padding: 20px 20px 0;
      dt {
        font-size: 14px;
        color: $color-text-main;
        line-height: 20px;
        margin-bottom: 20px;
      }
      dd {
        overflow: hidden;
        .tag-content-item {
          float: left;
          overflow: hidden;
          margin-bottom: 20px;
          .sub-tag,
          .sub-connect {
            float: left;
          }
          .sub-connect {
            line-height: 32px;
            padding: 0 10px;
          }
        }
      }
    }
    &-connect {
      float: left;
      padding: 0 20px;
      padding-top: 46px;
    }
  }
}
</style>
<style>
.group-show .layout-form .footer {
  margin-top: 0px !important;
}
</style>
