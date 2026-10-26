<template>
  <div class="select-materials-dialog">
    <!-- 新增小程序 -->
    <div class="add-box" v-if="addSmallBox">
      <cy-form
        class="main-box"
        ref="cform"
        :rules="rules"
        :model="ruleForm"
        label-width="90px"
      >
        <cy-form-item label="素材分组" prop="libraryGroup" trigger="change">
          <cy-select
            v-model="ruleForm.libraryGroup"
            width-size="medium"
            placeholder="请选择素材分组"
            search-input
            :optionsList="newGroupData"
            :options="{ label: 'group_name', value: 'id' }"
          >
          </cy-select>
        </cy-form-item>
        <cy-form-item label="小程序" prop="miniPro" required>
          <cy-select
            v-model="ruleForm.miniPro"
            width-size="medium"
            placeholder="请选择小程序"
            search-input
            :optionsList="smallList"
            :options="{ label: 'app_name', value: 'id' }"
          >
          </cy-select>
          <el-button
            type="text"
            style="margin-left:20px;"
            @click="addSmall"
            size="small"
            >添加小程序</el-button
          >
          <span class="shugang">|</span>
          <el-button type="text" @click="editSmall" size="small"
            >编辑小程序</el-button
          >
        </cy-form-item>
        <cy-form-item label="卡片标题" prop="title" max="28" required>
          <cy-input
            v-model="ruleForm.title"
            placeholder="输入卡片标题"
            show-word-limit
            maxlength="28"
            show-overflow-tooltip
            clearable
            width-size="medium"
          ></cy-input>
        </cy-form-item>
        <cy-form-item label="页面路径" prop="pageurl" required>
          <cy-input
            v-model="ruleForm.pageurl"
            placeholder="请输入小程序跳转页面路径"
            clearable
            width-size="medium"
          ></cy-input>
          <el-button
            type="text"
            style="margin-left:20px;"
            @click="getPageUrl"
            size="small"
            >如何获取页面路径？</el-button
          >
        </cy-form-item>
        <cy-form-item label="上传图片" required prop="urls" trigger="change">
          <upload-image
            class="upload-template"
            accept="image/jpeg,image/png"
            :urls="ruleForm.urls"
            :before-upload="beforeUpload"
            @change="update"
          >
          </upload-image>
          <span style="color: #707277;">支持JPG、PNG格式，图片小于2MB</span>
        </cy-form-item>
      </cy-form>
    </div>
    <!-- 添加小程序 -->
    <div class="push-box" v-show="pushSmallBox">
      <cy-form
        class="addmini-box"
        ref="addFormRef"
        :rules="rulesEdit"
        :model="addruleForm"
        label-width="130px"
      >
        <cy-form-item label="小程序名称" prop="mininame" required>
          <cy-input
            v-model="addruleForm.mininame"
            placeholder="请输入小程序名称"
            maxlength="30"
            show-word-limit
            clearable
            width-size="medium"
          ></cy-input>
        </cy-form-item>
        <cy-form-item label="小程序APP ID" prop="appid" required>
          <cy-input
            v-model="addruleForm.appid"
            placeholder="请输入小程序APP ID"
            maxlength="30"
            show-word-limit
            clearable
            width-size="medium"
          ></cy-input>
          <div>在微信公众平台>基础配置模块复制APP ID</div>
        </cy-form-item>
        <cy-form-item label="小程序原始ID" prop="originid" required>
          <cy-input
            v-model="addruleForm.originid"
            placeholder="请输入小程序原始ID"
            clearable
            width-size="medium"
          ></cy-input>
          <div>在微信公众平台>设置>公众号设置模块复制原始ID</div>
        </cy-form-item>
      </cy-form>
    </div>
    <!-- 编辑小程序 -->
    <div class="edit-box" v-if="editSmallArgin">
      <el-table :data="smallList" height="510" stripe tooltip-effect="light">
        <el-table-column
          label="小程序名称"
          prop="app_name"
          min-width="550"
          show-overflow-tooltip
        >
        </el-table-column>
        <el-table-column label="小程序类型" prop="type" width="150">
          <template slot-scope="scope">
            <span>{{ scope.row.type == 1 ? "手动添加" : "授权" }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="112">
          <template slot-scope="scope" v-if="scope.row.type == 1">
            <el-button
              type="text"
              class="m-r-20"
              @click="editItem(scope.row)"
              size="small"
              >编辑</el-button
            >
            <el-popconfirm
              title="你确定删除该小程序吗?"
              @onConfirm="delItem(scope.row.id)"
              icon="el-icon-info"
              cancel-button-type="Primary"
              confirm-button-text="确认"
            >
              <el-button
                type="text"
                slot="reference"
                style="color:#f56c6c;"
                size="small"
                >删除</el-button
              >
            </el-popconfirm>
          </template>
        </el-table-column>
        <div slot="empty">
          <img
            class="empty-img"
            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png"
            alt=""
          />
          <div class="empty-box">
            <div class="desc">暂无小程序</div>
          </div>
        </div>
      </el-table>
      <div style="text-align: right;" v-if="editmeta.total > 0">
        <el-pagination
          @size-change="edithandleSizeChange"
          @current-change="edithandleCurrentChange"
          :current-page.sync="editmeta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="10"
          layout="prev, pager, next"
          :total="editmeta.total"
        >
        </el-pagination>
      </div>
    </div>

    <div v-show="pushSmallBox" class="addfoot">
      <el-button @click="goBack" size="small">返 回</el-button>
      <el-button size="small" type="primary" @click="addchildSure"
        >保 存</el-button
      >
    </div>

    <div v-if="addSmallBox" class="addfoot">
      <el-button @click="goBack" size="small">返 回</el-button>

      <el-button type="primary" @click="childSure" size="small">{{
        editObj && editObj.type == "create" ? "保 存" : "确 认"
      }}</el-button>
    </div>
  </div>
</template>
<script lang="ts">
import { Mixins, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import UploadImage from "@/components/business/upload-image/index.vue";
import UploadValidate from "@/mixins/UploadValidate";
import UploadDrag from "@/components/business/upload-drag/index.vue";
interface FormIn {
  libraryGroup: number | undefined;
  miniPro: string;
  title: string;
  pageurl: string;
  urls: string;
}
@Component({
  components: {
    UploadImage,
    UploadDrag
  }
})
export default class extends Mixins(UploadValidate) {
  @Prop() appletsDialog!: boolean;
  @Prop() editObj!: any;
  @Prop({ default: 0, type: Number }) chooseNum?: number;
  //------------------------------------------------------data-----------------------------------------------------------/
  // 标题显示
  private addOredit = 1;
  private addBox = false; //新增小程序
  private pushBox = false; //添加小程序
  private editBox = false; //编辑小程序
  private editArginBox = false; //编辑里面的编辑
  // 内容显示
  private chooseSmallBox = true; //选择小程序素材
  private addSmallBox = false; //新增小程序
  private pushSmallBox = false; //添加,编辑小程序编辑
  private editSmallArgin = false; //编辑小程序

  private tableList: any = []; // 表格数据
  private smallList: any = []; // 小程序数据
  private selectDialogVisible = true;
  private emptyFlag = false;
  private ruleForm: FormIn = {
    libraryGroup: 0,
    miniPro: "",
    title: "",
    pageurl: "",
    urls: ""
  };
  // 搜索区域
  private search = {
    keyWord: ""
  };
  // 是否全选
  private checked = false;
  // 列表移动分组数据
  private groupData: Array<any> = [];
  private originData: any[] = [
    { group_name: "全部分组", id: 0 },
    { group_name: "未分组", id: -1 }
  ];
  private newGroupData: Array<any> = [];
  private newOriginData: any[] = [{ group_name: "未分组", id: -1 }]; //左侧分组
  private isActive = 0;
  //分页
  private meta = {
    size: 10,
    page: 1,
    total: 0
  };
  private editmeta = {
    size: 10,
    page: 1,
    total: 0
  };
  private groupList: any = []; //下拉列表
  private addruleForm = {
    mininame: "",
    appid: "",
    originid: ""
  };
  private rules = {};
  private rulesEdit = {};
  private miniId = 0; //编辑ID
  // 小程序列表选是否中
  private isChecked: any = [];
  // 小程序列表选中数据
  private selectedMaterial: any = [];
  // --------------------------------------------------------------------面包屑点击返回---------------------------------------------------------------------
  /*
  chooseSmallBox//选择小程序素材
  addSmallBox //新增小程序
  pushSmallBox //添加,编辑小程序编辑
  editSmallArgin //编辑小程序
*/
  // 点击选择小程序
  private choosebacksmal() {
    this.addBox = false;
    this.chooseSmallBox = true;
    this.pushSmallBox = false;
    this.addSmallBox = false;
    this.editSmallArgin = false;
  }
  // 新增小程序
  private addbacksmal() {
    this.editmeta.size = 999;
    this.miniList();
    this.ruleForm.libraryGroup = -1;
    this.ruleForm.miniPro = "";
    this.ruleForm.title = "";
    this.ruleForm.pageurl = "";
    this.ruleForm.urls = "";
    this.pushBox = false;
    this.editBox = false;
    this.pushSmallBox = false;
    this.chooseSmallBox = false;
    this.addSmallBox = true;
    this.editSmallArgin = false;
    this.editArginBox = false;
    // console.log(11111);
  }
  // 添加小程序
  private pushbacksmal() {
    this.editBox = false;
    this.pushSmallBox = true;
    this.chooseSmallBox = false;
    this.addSmallBox = false;
    this.editSmallArgin = false;
  }
  // 编辑小程序
  private editArginbacksmal() {
    //  console.log("444444444444444");
    this.editBox = true;
    this.editArginBox = false;
    this.editSmallArgin = true;
    this.addSmallBox = false;
    this.chooseSmallBox = false;
    this.pushSmallBox = false;
  }
  // -----------------------------------------------------------------------------------页面点击进入下一页--------------------------------------------------------------------------
  // 点击新增小程序
  private addTextmat() {
    if (this.ruleForm.libraryGroup == 0 || this.ruleForm.libraryGroup == -1) {
      this.ruleForm.libraryGroup = undefined;
    }
    this.editmeta.size = 999;
    this.miniList();
    this.ruleForm.miniPro = "";
    this.ruleForm.title = "";
    this.ruleForm.pageurl = "";
    this.ruleForm.urls = "";
    this.addBox = true;
    this.pushSmallBox = false;
    this.addSmallBox = true;
    this.chooseSmallBox = false;
    this.editSmallArgin = false;
    this.editArginBox = false;
  }
  // 点击添加小程序
  private addSmall() {
    this.addruleForm.mininame = "";
    this.addruleForm.appid = "";
    this.addruleForm.originid = "";
    this.pushBox = true;
    this.pushSmallBox = true;
    this.addSmallBox = false;
    this.chooseSmallBox = false;
    this.editSmallArgin = false;
  }
  // 点击编辑小程序
  private editSmall() {
    this.editmeta.size = 10;
    this.editmeta.page = 1;
    this.miniList();
    this.editBox = true;
    this.editSmallArgin = true;
    this.addSmallBox = false;
    this.chooseSmallBox = false;
    this.pushSmallBox = false;
  }
  // 点击编辑小程序编辑
  private editItem(data: any) {
    this.miniId = data.id;
    this.addruleForm.mininame = data.app_name;
    this.addruleForm.appid = data.app_id;
    this.addruleForm.originid = data.origin_id;
    this.editArginBox = true;
    this.addOredit = 2;
    this.pushSmallBox = true;
    this.addSmallBox = false;
    this.chooseSmallBox = false;
    this.editSmallArgin = false;
  }
  //----------------------------------------------------选择小程序methods-----------------------------------------------------------/
  // 分组接口
  grouplist() {
    Api.sourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const data = res.data;
        this.groupData = [...this.originData, ...data];
        this.newGroupData = [...this.newOriginData, ...data];
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // -----------------------------------------------------------------------------------新增小程序methods--------------------------------------------------------------------------

  // 获取页面路径跳转
  private getPageUrl() {
    window.open(
      "https://jingyan.baidu.com/article/f3ad7d0f4c39aa09c3345bf1.html",
      "_blank"
    );
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    if (name !== "jpg" && name !== "png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }
  }
  private update(res: Array<string>) {
    // this.ruleForm.urls = res;
    this.$set(this.ruleForm, "urls", res);
    (this.$refs["cform"] as any).validateField("urls");
  }
  // 新增小程序点击确定
  private childSure() {
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        this.addMiniProgram();
      }
    });
  }
  addMiniProgram() {
    const params: any = {
      group_id: this.ruleForm.libraryGroup || -1,
      title: this.ruleForm.title,
      mini_program_id: this.ruleForm.miniPro,
      path: this.ruleForm.pageurl,
      cover_url: this.ruleForm.urls
    };
    if (this.editObj && this.editObj.type == "edit" && this.appletsDialog) {
      params.id = this.editObj.id;
    }
    Api.miniProgramSaveApi(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("添加成功");
        //列表的新增过来的新增小程序
        if (this.appletsDialog) {
          this.$emit("submit", { save: true });
        } else {
          this.$emit("submit", res.data);
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //编辑详情
  private editMiniProgramFn() {
    Api.miniProgramDetailApi({ id: this.editObj.id }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data;
        this.ruleForm = {
          libraryGroup: data.group_id,
          miniPro: data.mini_program_id,
          title: data.title,
          pageurl: data.path,
          urls: data.cover_url
        };
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // -------------------------------------------------------------------------------------添加小程序methods----------------------------------------------------------------------------
  private addchildSure() {
    (this.$refs["addFormRef"] as any).validate((valid: any) => {
      if (valid) {
        if (this.editArginBox) {
          this.editArginBox = false;
          this.editmeta.size = 10;
          this.editmeta.page = 1;
          this.miniProgramEdit();
          this.editBox = true;
          this.editSmallArgin = true;
          this.addSmallBox = false;
          this.chooseSmallBox = false;
          this.pushSmallBox = false;
        } else {
          setTimeout(() => {
            this.miniProgramAdd();
            this.addBox = true;
            this.pushSmallBox = false;
            this.addSmallBox = true;
            this.chooseSmallBox = false;
            this.editSmallArgin = false;
          }, 300);
        }
      }
    });
  }
  private miniProgramAdd() {
    const params = {
      app_name: this.addruleForm.mininame,
      app_id: this.addruleForm.appid,
      origin_id: this.addruleForm.originid
    };
    Api.miniProgramAddApi(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("添加成功");
        this.miniList();
        setTimeout(() => {
          this.ruleForm.miniPro = res.data.id;
        }, 500);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private miniList() {
    const params = {
      page: this.editmeta.page,
      size: this.editmeta.size
    };
    Api.miniListApi(params).then((res: any) => {
      if (res.code == 200) {
        this.smallList = res.data.list;
        this.editmeta.total = res.data.total;
        if (this.appletsDialog && this.smallList.length != 0) {
          //默认选中第一个
          this.ruleForm.miniPro = this.smallList[0].id;
        }
        //调用编辑
        if (this.editObj && this.editObj.type == "edit" && this.appletsDialog) {
          this.editMiniProgramFn();
        } else if (
          this.editObj &&
          this.editObj.type == "create" &&
          this.appletsDialog &&
          this.editObj.id != 0 &&
          this.editObj.id != -1
        ) {
          this.ruleForm.libraryGroup = this.editObj.id;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // -------------------------------------------------------------编辑小程序编辑------------------------------------------------------------
  // 删除小程序
  private miniProgramDel(id: any) {
    Api.miniProgramDelApi({ id: id }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("删除成功");
        this.editmeta.size = 10;
        this.editmeta.page = 1;
        this.miniList();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 编辑小程序
  private miniProgramEdit() {
    const params = {
      app_name: this.addruleForm.mininame,
      app_id: this.addruleForm.appid,
      origin_id: this.addruleForm.originid,
      id: this.miniId
    };
    Api.miniProgramEditApi(params).then((res: any) => {
      if (res.code == 200) {
        console.log("编辑成功");
        this.editmeta.size = 10;
        this.editmeta.page = 1;
        this.miniList();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  private goBack() {
    this.$emit("submit");
  }
  // 点击删除小程序按钮
  private delItem(id: any) {
    this.miniProgramDel(id);
  }
  // 分页
  edithandleSizeChange(val: number) {
    this.editmeta.size = val;
    this.editmeta.page = 1;
    this.miniList();
  }
  edithandleCurrentChange(val: number) {
    this.editmeta.page = val;
    this.miniList();
  }
  mounted() {
    if (this.appletsDialog) {
      this.ruleForm = {
        libraryGroup: undefined,
        miniPro: "",
        title: "",
        pageurl: "",
        urls: ""
      };
      this.chooseSmallBox = false;
      this.addSmallBox = true;
      this.addBox = true;
      this.miniList();
    }
    this.grouplist();
  }
  created() {
    this.addTextmat();
  }
}
</script>
<style lang="scss" scoped>
.chosoe-clas {
  color: #33363c;
  font-weight: 600;
}
/deep/.el-dialog {
  max-height: 660px;
}

.fotter {
  margin-top: 20px;
  text-align: right;
}

.foot {
  margin-top: 20px;
  text-align: right;
  div {
    margin-bottom: 20px;
  }
}
.select-materials-dialog {
  .shugang {
    color: #bbbbbb;
    margin: 0 20px;
  }
  .in-item {
    width: 100% !important;
  }
  .header {
    display: flex;
    justify-content: space-between;
    text-align: center;
    margin-bottom: 20px;
    font-size: 16px;
    cursor: pointer;
  }
  /deep/.el-dialog__header {
    display: none;
  }
  .search {
    margin-bottom: 5px;
    position: relative;
    .addClass {
      position: absolute;
      right: 0;
    }
  }
  .content {
    display: flex;
    .menu-left {
      margin-right: 10px;
      background: $color-background;
      width: 200px;
      height: 380px;
      overflow-y: scroll;
      .sub-wrap {
        dd {
          display: flex;
          align-items: center;
          justify-content: space-between;
          height: 20px;
          line-height: 20px;
          span {
            padding-left: 20px;
          }
          img {
            width: 3px;
            height: 17px;
          }
        }
        dd.check-color span {
          border-left: 4px solid #708df5;
          color: #708df5;
        }
        img {
          display: inline-block;
          width: 20px;
          height: 20px;
          margin-right: 10px;
          border-radius: 50%;
        }
      }
      dd {
        margin-inline-start: 0px;
        height: 20px;
        line-height: 20px;
        color: $color-text-main;
        font-size: 14px;
        padding: 10px 0;
        cursor: pointer;
      }
      dd:hover {
        background: #ebebf1;
      }
    }
    .right-content {
      flex: 1;
      width: 731px;
      .right-cont-top {
        height: 363px;
        clear: both;
        overflow-y: scroll;
      }
      .cont-item:nth-child(4n) {
        margin-right: 0px;
      }
      .cont-item {
        width: 145px;
        // height: 191px;
        background: #ffffff;
        border-radius: 4px;
        margin: 0 20px 20px 0;
        padding: 8px;
        border: 1px solid #dcdfe6;
        float: left;
        cursor: pointer;
        .small-tit {
          margin-bottom: 5px;
          display: flex;
          .minilogo {
            display: inline-block;
            margin-right: 5px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
          }
          .smallName {
            color: #707277;
          }
        }
        span {
          width: 100%;
          display: inline-block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .miniimg {
          display: block;
          width: 145px;
          height: 116px;
          border-radius: 4px;
        }
      }
      .selected-item {
        border: 1px solid #708df5;
        position: relative;
      }
      .item-select {
        position: absolute;
        right: 0;
        top: 0;
        width: 14px;
        height: 14px;
      }
    }
  }
}
.tit-name {
  text-overflow: -o-ellipsis-lastline;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  margin-bottom: 10px;
  height: 40px;
}
.add-box,
.push-box {
  height: 528px;
  overflow: auto;
}
.addfoot {
  text-align: right;
  margin-top: 20px;
}
</style>
