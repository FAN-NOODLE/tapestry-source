<template>
  <div>
    <layout-form type="normal">
      <template slot="content">
        <dl class="form-layout key-required">
          <dt>分组名称</dt>
          <dd>
            <cy-input
              width-size="medium"
              placeholder="请输入分组名称"
              v-model="groupName"
              maxlength="10"
              show-word-limit
              @blur="checkFn"
              @keyup.native="trimFn"
              @input="isValueFn"
            ></cy-input>
            <transition name="error">
              <p class="error-tips" v-if="isGroupName">
                请输入分组名称
              </p>
            </transition>
          </dd>
        </dl>
        <dl class="form-layout key-required">
          <dt>分组维度</dt>
          <dd>
            <template>
              <el-radio v-model="formData.user_group_type" :label="1" :disabled="type === 'edit'"
                >标签</el-radio
              >
              <el-radio v-model="formData.user_group_type" :label="2" :disabled="type === 'edit'"
                >顾客名单</el-radio
              >
              <el-radio v-model="formData.user_group_type" :label="3" :disabled="type === 'edit'"
                >门店</el-radio
              >
            </template>
          </dd>
        </dl>
        <!-- 标签 原有逻辑 -->
        <template v-if="formData.user_group_type === 1">
          <dl class="form-layout key-required">
            <dt>包含标签</dt>
            <dd>
              <el-button
                type="primary"
                size="small"
                @click="selectTagFn('include')"
                >添加</el-button
              >
              <span>同一标签组内的标签为“或”关系，标签组互相为“且”关系</span>
              <transition name="error">
                <p class="error-tips" v-if="isInclude">
                  请选择包含标签
                </p>
              </transition>
            </dd>
          </dl>
          <tags
            :tagsArr="selectInclude"
            type="include"
            @deleteFn="tagCloseFn"
            v-if="selectInclude.length"
          ></tags>
          <dl class="form-layout">
            <dt>排除标签</dt>
            <dd>
              <el-button
                type="primary"
                size="small"
                @click="selectTagFn('exclude')"
                >添加</el-button
              >
              <span>同一标签组内的标签为“或”关系，标签组互相为“且”关系</span>
            </dd>
          </dl>
          <tags
            :tagsArr="selectExclude"
            type="exclude"
            @deleteFn="tagCloseFn"
            v-if="selectExclude.length"
          ></tags>
        </template>
        <!-- 顾客名单 -->
        <template v-if="formData.user_group_type === 2">
          <dl class="form-layout">
            <dt>下载模板</dt>
            <dd>
              <el-button type="text" @click="temDown"
                >顾客名单模板下载</el-button
              >
            </dd>
          </dl>
          <dl class="form-layout key-required">
            <dt>上传文件</dt>
            <dd>
              <el-upload
                ref="upload"
                type="file"
                :limit="1"
                :on-exceed="excedFun"
                :data="exceldata"
                :action="uploadUrl"
                :before-upload="beforeUpload"
                :headers="headers"
                :file-list="file_list"
                :on-success="onSuccessFun"
                :on-remove="onRemoveFn"
                drag
              >
                <i class="el-icon-upload"></i>
                <div class="el-upload__text">
                  将文件拖到此处，或<em>点击上传</em>
                  <div>只能上传xls、xlsx文件</div>
                </div>
              </el-upload>
            </dd>
          </dl>
        </template>
        <!-- 门店 -->
        <template v-if="formData.user_group_type === 3">
          <dl class="form-layout key-required">
            <dt>选择门店</dt>
            <dd>
              <el-button type="text" @click="selectStore">选择门店</el-button>
            </dd>
          </dl>
        </template>
        <template v-if="formData.user_group_type === 3 && storeList.length > 0">
          <dl>
            <dt></dt>
            <dd>
              <el-table
                :data="storeList"
                ref="tableList"
                style="width:100%"
                :cell-style="{ height: '20px' }"
              >
                <el-table-column
                  label="门店名称"
                  prop="dp_name"
                ></el-table-column>
              </el-table>
            </dd>
          </dl>
        </template>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="saveClickFn">保存</el-button>
        <el-button @click="goCancelFn">取消</el-button>
      </template>
    </layout-form>
    <!-- 创建离开页面弹框 -->
    <el-dialog
      title="离开页面提示"
      :visible.sync="isCancelDialog"
      width="420px"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        检测到未保存的内容，是否保存到下次新增继续编辑？</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goBackFn">直接离开</el-button>
        <el-button type="primary" @click="saveGoBackFn">保存并离开</el-button>
      </span>
    </el-dialog>
    <!-- 编辑离开页面弹框 -->
    <el-dialog
      title="离开页面提示"
      :visible.sync="isTipsCancel"
      width="420px"
      @close="goCloseFn"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        你确定要放弃编辑吗？放弃后编辑的信息将不会被保存</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goCloseFn">取消</el-button>
        <el-button type="primary" @click="editGoBackFn">确定</el-button>
      </span>
    </el-dialog>
    <select-store
      v-if="showStoreModel"
      @submit="submitStore"
      @close="closeStore"
      choice_name="已选择的门店"
      :showCate="false"
      :selectArr="storeList"
    ></select-store>
    <choose-label
      v-if="chooseLabelDialog"
      :toDialog="toDialog"
      @closeDialog="closeDialog"
    ></choose-label>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from "vue-property-decorator";
import API from "@/api/customer";
import tags from "../components/tags.vue";
import chooseLabel from "../components/choose-label.vue";
import CustomerModule from "@/store/modules/customer";
import SelectStore from "@/components/business/select-store/index.vue";
import { VUE_APP_API_ROOT, VUE_APP_API_ROOT_RE } from "@/utils/env-config";
import LoginAboat from "@/utils/LoginManager";
@Component({
  components: { tags, chooseLabel, SelectStore }
})
export default class extends Vue {
  @Prop() type?: string;
  //从vuex取包含标签
  get GROUPS() {
    return CustomerModule.getGroupS;
  }
  //从vuex取排除标签
  get GROUPU() {
    return CustomerModule.getGroupU;
  }
  //从vuex取分组名称
  get GROUPNAME() {
    return CustomerModule.GROUPNAME;
  }
  @Watch("GROUPS", { deep: true })
  @Watch("GROUPU", { deep: true })
  private onChangeGroupuType() {
    console.log("watchFn :>> ");
    this.watchFn();
  }
  /*************************** data start ***************************/
  private groupName = ""; //分组名称
  private selectInclude: any = []; //包含标签
  private selectExclude: any = []; //排除标签
  private userCountArr: any = []; //组装数据传接口
  private conditionsArr: any = []; //组装数据，传给后台的入参
  private isCancelDialog = false; //取消按钮是否弹框
  private isTipsCancel = false; //编辑取消按钮是否弹框
  private formData: any = {
    user_group_type: 1,
    file_id: ""
  };
  private file_list = [];
  private showStoreModel = false;
  private preObj: any = {
    include: [],
    exclude: [],
    name: ""
  }; //以前保存的值
  private resetObj: any = {
    include: [],
    exclude: [],
    name: ""
  }; //重置的空值
  private isGroupName = false; //是否输入数组名称
  private isInclude = false; //是否选择包含标签
  private msgTips: any; // 处理提示框问题
  private curId = 0; //当前id
  private chooseLabelDialog = false; //弹框
  private toDialog = {};
  private export_model = false;
  private exceldata: any = {}; // 上传时附带的额外参数
  private uploadUrl = `${VUE_APP_API_ROOT}user/group-user-import`;
  private headers = {
    "Auth-Token": LoginAboat.getToken()
  };
  private storeList = [];

  /*************************** method start ***************************/
  private created() {
    if (this.type == "create") {
      this.createShow();
    } else if (this.type == "edit") {
      this.curId = Number(this.$route.query.id);
      this.getDetailFn();
    }
  }
  private delStore(index) {
    this.storeList.splice(index, 1);
  }
  private selectStore() {
    this.showStoreModel = true;
  }
  private closeStore() {
    this.showStoreModel = false;
  }
  private submitStore(arr) {
    this.storeList = arr;
    this.showStoreModel = false;
  }
  // 模板下载
  private temDown() {
    window.location.href = `${VUE_APP_API_ROOT_RE}ba-performance/group-download`;
  }
  // 上传文件个数超过范围
  private excedFun() {
    this.$message.error("上传失败，每次仅支持上传一个文件！");
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    // this.exceldata.file_type = "mars_code";
    const testmsg = file.name.substring(file.name.lastIndexOf(".") + 1);
    const extension = testmsg === "xlsx" || testmsg === "xls";
    if (!extension) {
      this.$message.error("请上传 xls、xlsx 类型的文件!");
      return false;
    }
  }
  // 上传成功
  private onSuccessFun(res: any) {
    if (res.code == 200) {
      this.formData.file_id = res.data.id;
    } else {
      this.file_list = [];
      this.$message.error(res.message);
    }
  }
  // 移除上传
  private onRemoveFn() {
    this.file_list = [];
    this.formData.file_id = 0;
  }
  private watchFn() {
    if (this.type == "create") {
      this.createShow();
    } else if (this.type == "edit") {
      this.preObj = JSON.parse(
        sessionStorage.getItem("create_obj") || JSON.stringify(this.preObj)
      );
      this.selectInclude =
        this.GROUPS.length != 0 ? this.GROUPS : this.selectInclude;
      this.selectExclude =
        this.GROUPU.length != 0 ? this.GROUPU : this.selectExclude;
    }

    // 包含标签提示语是否显示
    this.isInclude = !this.selectInclude.length;
  }
  //input 不输入空格
  private trimFn() {
    this.groupName = this.groupName.replace(/^\s+|\s+$/gm, "");
  }
  //input 判断是否有值
  private isValueFn() {
    this.isGroupName = this.groupName ? false : true;
  }
  // 组装数据
  private getUsercountFn(arr: any, type: string) {
    this.userCountArr = JSON.parse(JSON.stringify(arr));
    this.userCountArr.forEach((item: any) => {
      delete item.tag_group_name;
      const temp_arr: any = [];
      item.tags.forEach((subItem: any) => {
        temp_arr.push(subItem.tag_id);
        item.tag_ids = temp_arr;
      });
      delete item.tags;
    });
    this.conditionsArr.push({
      tag_groups: this.userCountArr,
      include_type: type == "include" ? 1 : 2
    });
  }
  //input 失焦方法
  private checkFn() {
    CustomerModule.setGroupName(this.groupName);
  }
  //保存
  private saveClickFn() {
    if (!this.groupName) {
      this.isGroupName = true;
    }
    this.conditionsArr = [];
    const params: any = {
      name: this.groupName,
      file_id: 0,
      department_ids: [],
      user_group_type: this.formData.user_group_type
    };
    if (this.formData.user_group_type === 1) {
      // 包含标签提示语是否显示
      this.isInclude = !this.selectInclude.length;
      if (!this.groupName || this.selectInclude.length == 0) {
        this.msgTips ? this.msgTips.close() : "";
        this.msgTips = this.$message.warning("分组信息未完善！");
        return false;
      }
      this.getUsercountFn(this.selectInclude, "include");
      this.getUsercountFn(this.selectExclude, "exclude");
    }
    if (this.formData.user_group_type === 2) {
      this.conditionsArr = [];
      if (!this.groupName || !this.formData.file_id) {
        this.msgTips ? this.msgTips.close() : "";
        this.msgTips = this.$message.warning("分组信息未完善！");
        return false;
      }
      params.file_id = this.formData.file_id;
    }
    if (this.formData.user_group_type === 3) {
      this.conditionsArr = [];
      if (!this.groupName || this.storeList.length === 0) {
        this.msgTips ? this.msgTips.close() : "";
        this.msgTips = this.$message.warning("分组信息未完善！");
        return false;
      }
      params.department_ids = this.storeList.map(item => {
        return item.department_id;
      });
    }
    params.conditions = this.conditionsArr;
    let api = "";
    if (this.curId) {
      params.id = this.curId;
      api = "groupEditApi";
    } else {
      api = "groupCreateApi";
    }
    API[api](params).then((res: any) => {
      if (res.code == 200) {
        sessionStorage.setItem("create_obj", JSON.stringify(this.resetObj));
        if (this.curId) {
          this.$message.success("分组编辑完成");
        } else {
          this.$message.success("分组创建完成");
        }
        this.$router.push({ name: "myapp_customer_customer-groups_list" });
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private selectTagFn(type: string) {
    // type=='include'||'exclude' 包含||不包含
    this.toDialog = {
      type: type,
      include: this.selectInclude,
      exclude: this.selectExclude
    };
    this.chooseLabelDialog = true;
  }
  private closeDialog() {
    this.chooseLabelDialog = false;
  }
  //取消
  private goCancelFn() {
    if (this.type === "create") {
      this.createCancelFn();
    } else if (this.type === "edit") {
      this.editCancelFn();
    }
  }
  //创建取消
  private createCancelFn() {
    if (
      this.groupName ||
      this.selectInclude.length != 0 ||
      this.selectExclude.length != 0
    ) {
      this.isCancelDialog = true;
    } else {
      this.goBackFn();
    }
  }
  //编辑取消
  private editCancelFn() {
    if (
      this.groupName != this.preObj.name ||
      JSON.stringify(this.selectInclude) !=
        JSON.stringify(this.preObj.include) ||
      JSON.stringify(this.selectExclude) != JSON.stringify(this.preObj.exclude)
    ) {
      this.isTipsCancel = true;
    } else {
      this.editRouteBackFn();
    }
  }
  //返回上一页
  private goBackFn() {
    this.isCancelDialog = false;
    sessionStorage.setItem("create_obj", JSON.stringify(this.resetObj));
    this.routeBackFn();
  }
  //保存返回上一页
  private saveGoBackFn() {
    this.isCancelDialog = false;
    const local_arr = {
      include: this.selectInclude,
      exclude: this.selectExclude,
      name: this.groupName
    };
    //新建的值保存sessionStorage
    sessionStorage.setItem("create_obj", JSON.stringify(local_arr));
    this.routeBackFn();
  }
  //路由跳转
  private routeBackFn() {
    this.$router.push({
      name: "myapp_customer_customer-groups_list"
    });
  }
  //路由跳转
  private editRouteBackFn() {
    this.$router.push({
      name: "myapp_customer_customer-groups_show",
      query: { id: this.curId.toString() }
    });
  }
  //取消选中/tag删除
  private tagCloseFn(obj: any) {
    let arrTags: any = [];
    if (obj.type == "include") {
      arrTags = this.selectInclude;
    } else {
      arrTags = this.selectExclude;
    }
    arrTags.forEach((curItem: any, cueIndex: number) => {
      curItem.tags.forEach((subCurItem: any, subCurIndex: number) => {
        if (subCurItem.tag_id == obj.item.tag_id) {
          curItem.tags.splice(subCurIndex, 1);
        }
      });
      if (curItem.tags.length == 0) {
        arrTags.splice(cueIndex, 1);
      }
    });
  }
  //创建页面进来，拿sessionStorage值
  private createShow() {
    // type=='include'||'exclude' 包含||不包含
    this.preObj = JSON.parse(
      sessionStorage.getItem("create_obj") || JSON.stringify(this.preObj)
    );
    this.selectInclude =
      this.GROUPS.length != 0 ? this.GROUPS : this.preObj.include;
    this.selectExclude =
      this.GROUPU.length != 0 ? this.GROUPU : this.preObj.exclude;
    this.groupName = this.GROUPNAME
      ? JSON.parse(JSON.stringify(this.GROUPNAME))
      : this.preObj.name;
  }
  //编辑关闭和取消
  private goCloseFn() {
    this.isTipsCancel = false;
  }
  //编辑确定
  private editGoBackFn() {
    this.isTipsCancel = true;
    this.editRouteBackFn();
  }

  //获取详情
  private getDetailFn() {
    API.groupDetailApi({ id: this.curId }).then(res => {
      const { data } = res;
      // 分组名称
      this.groupName = this.GROUPNAME
        ? JSON.parse(JSON.stringify(this.GROUPNAME))
        : data.name;
      this.preObj.name = this.GROUPNAME
        ? JSON.parse(JSON.stringify(this.GROUPNAME))
        : data.name;
      this.formData = {
        user_group_type: data.user_group_type,
        file_id: ""
      };
      this.storeList = data.Departments?.map(item => {
        return {
          department_id: item.id,
          dp_name: item.name
        };
      });
      data.conditions.forEach((item: any) => {
        if (item.include_type == 1) {
          //包含
          this.selectInclude =
            this.GROUPS.length != 0 ? this.GROUPS : item.tag_groups;
          this.preObj.include = JSON.parse(JSON.stringify(this.selectInclude));
        } else {
          //非包含
          this.selectExclude =
            this.GROUPU.length != 0 ? this.GROUPU : item.tag_groups;
          this.preObj.exclude = JSON.parse(JSON.stringify(this.selectExclude));
        }
      });
    });
  }
  /*************************** 生命周期 start ***************************/
}
</script>
<style lang="scss" scoped>
/deep/.el-upload-list {
  width: 500px;
}
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
.form-layout {
  display: flex;
  margin-bottom: 30px;
  justify-content: center;
  align-items: center;
  padding-left: 10px;
  dt {
    flex: 0 0 64px;
    margin-right: 20px;
    color: $color-text-normal;
    span {
      color: #f56c6c;
      display: inline-block;
      line-height: 20px;
    }
  }
  dd {
    flex: 1;
    position: relative;
    span {
      color: $color-text-second;
      font-size: 14px;
      margin-left: 10px;
    }
    .error-tips {
      color: #f56c6c;
      font-size: 12px;
      line-height: 20px;
      position: absolute;
      top: 100%;
      left: 0;
      z-index: 1;
      overflow: hidden;
    }
  }
}
.key-required {
  position: relative;
  dt:after {
    content: "*";
    color: red;
    position: absolute;
    width: 10px;
    height: 10px;
    line-height: 10px;
    left: 0;
    top: 50%;
    margin-top: -5px;
    z-index: 1;
  }
}
.tags {
  margin-left: 95px;
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

.error-leave-active,
.error-enter-active {
  transition: all 0.5s ease;
}
.error-leave-active,
.error-enter {
  height: 0px !important;
}
.error-leave,
.error-enter-active {
  height: 20px;
}
</style>
