<template>
  <div class="card-box">
    <layout-form type="normal">
      <template slot="content">
        <div class="create-wrap">
          <div class="sms-show">
            <p class="time">下午 6:30</p>
            <div
              class="content-msg pull-left"
              v-html="msgTextStr"
              ref="editWordLenDom"
            ></div>
          </div>
          <div class="sms-set">
            <cy-form ref="cTemplate" :rules="rules" :model="ruleForm">
              <cy-form-item
                label="模版名称"
                prop="template_short_name"
                required
              >
                <cy-input
                  width-size="medium"
                  placeholder="请输入模板名称"
                  v-model="ruleForm.template_short_name"
                  maxlength="30"
                  show-word-limit
                  trims
                  @keyup.native="checkSpace()"
                ></cy-input>
              </cy-form-item>
              <cy-form-item label="短信签名" prop="sign_id" trims required>
                <cy-select
                  width-size="medium"
                  size="small"
                  v-model="ruleForm.sign_id"
                  :optionsList="templateStatus"
                  :options="{ label: 'sign_name', value: 'id' }"
                  @change="changeSmsFn"
                ></cy-select>
              </cy-form-item>
              <cy-form-item
                label="短信内容"
                prop="isCheck"
                required
                class="new-class"
              >
                <div class="sms-set-msg">
                  <dl>
                    <dt>插入参数<span @click="tipsFn">如何插入参数？</span></dt>
                    <dd class="tag-wrap">
                      <el-button
                        class="tag-item"
                        size="mini"
                        v-for="(item, index) in paramArr"
                        :key="index"
                        :type="item.params_type == 1 ? 'warning' : 'primary'"
                        @click="setParamFn(item)"
                        :disabled="item.isForbidden ? true : false"
                        plain
                        >{{ item.name }}</el-button
                      >
                    </dd>
                    <dd class="sms-content">
                      <div id="editorSmsInput" ref="getSmsSign"></div>
                    </dd>
                  </dl>
                  <div class="sms-set-msg-tips clearfix">
                    <div class="pull-left">
                      已占用{{ curWordLen + 4 }}个字，预计计费为{{
                        Math.ceil((curWordLen + 4) / 70)
                      }}条短信
                    </div>
                    <div class="pull-right">
                      <span>*</span> 计费规则：1条短信70个字
                    </div>
                  </div>
                </div>
                <dl class="tips-list">
                  <!-- <dd>
                    1.营销类短信，请在短信内容后面增加 “回T退订”。
                  </dd> -->
                  <dd>
                    1.计费说明：单条短信 (签名+正文)
                    最大长度为70个字，超过70个字即为长消息，将按67字/条分隔成多条计费。例如，短信长度为150字，则按照67字/67字/16字
                    分隔成3条计费。
                  </dd>
                  <dd>
                    2.汉字、字母、数字、标点符号 (不区分全角/半角)
                    以及空格等，都按1个字计算。
                  </dd>
                  <dd>3.短信模版内容不能含有【】{}字符。</dd>
                  <!-- <dd>
                    4.本系统插入短链自动在链接前后各加入1个空格。
                  </dd> -->
                </dl>
              </cy-form-item>
            </cy-form>
          </div>
        </div>
      </template>
      <template slot="footer">
        <!-- <el-button type="primary" @click="saveFn(1, 'cTemplate')" size="small"
          >保存并提审</el-button
        > -->
        <el-button type="primary" @click="saveFn(2, 'cTemplate')" size="small"
          >保存</el-button
        >
        <el-button @click="cancelFn" size="small">取消</el-button>
      </template>
    </layout-form>
    <!-- 返回框弹框 -->
    <back-pop v-if="isPop" @goBack="goBackFn"></back-pop>
    <el-dialog
      title="如何插入参数"
      :visible.sync="isTips"
      width="420px"
      @close="tipsCloseFn"
    >
      <dl class="tipsText">
        <dd>
          短信模板中可放置动态参数。动态参数为黄色文字，发送短信时动态参数的实际值，会自动用实际值替换动态参数
        </dd>
      </dl>
    </el-dialog>
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
        <el-button @click="goCloseFn" size="small">取消</el-button>
        <el-button type="primary" @click="saveGoBackFn" size="small"
          >确定</el-button
        >
      </span>
    </el-dialog>
    <div ref="messageTemText" style="display:none" v-html="msgTextHtml"></div>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import API from "@/api/customer";
import backPop from "../../components/back-pop.vue";
// 引入 wangEditor
import wangEditor from "wangeditor";

@Component({
  components: { backPop }
})
export default class extends Vue {
  /*************************** data start ***************************/
  private editorSms: any = null; //编辑器
  private isCreate: any = "1"; //创建1或编辑2
  private rules = {};
  //表单
  private ruleForm = {
    template_short_name: "", // 20201221 短信模版优化 显示字段取消enterprise_id
    template_name: "",
    sign_id: 117567463708480,
    isCheck: "no"
  };
  private msgTextStr = ""; //显示的短信内容
  private msgTextHtml = ""; //组装的短信内容
  //参数
  private paramArr: any = [];
  private msgArr: any = [];
  private curId: any; //当前id
  //模板状态
  private templateStatus = [];
  private curWordLen = 0; //短信长度
  private editWordLen = 0; //编辑文字长度
  private menuWordLen = 0; //参数长度
  private isPop = false; //返回弹出框
  private showStr = ""; //短信模板字符串
  private isTips = false; //是否显示弹框
  private oldObj: any = {}; //存储旧的编辑数据
  private isTipsCancel = false; //取消弹框提示
  private templateContentObj: any = {}; //设置的参数
  private messageStr = "";
  private flag = true;
  /*************************** method start ***************************/
  //不能输入空格
  private checkSpace() {
    this.ruleForm.template_short_name = this.ruleForm.template_short_name.replace(
      /^\s+|\s+$/gm,
      ""
    );
  }
  //select change事件
  private changeSmsFn() {
    this.curWordLen = this.editWordLen + this.menuWordLen;
  }
  //获取参数列表
  private getSmsParamsListFn() {
    API.smsParamsListApi({}).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        this.paramArr = data.list;
        this.paramArr.forEach((item: any) => {
          item.isForbidden = false;
          item.strHtml = "";
          item.type = "dynamic";
        });
        if (this.isCreate == "2") {
          //编辑
          this.getEditInfoFn();
        } else {
          //新建
          const reEdit = localStorage.getItem("SMS_TEMPLATE");
          if (reEdit && reEdit != "{}") {
            const obj = JSON.parse(reEdit);
            this.ruleForm.template_short_name = obj.template_short_name; //TODO
            this.ruleForm.sign_id = obj.template_sign_id;
            this.showStr = obj.message_template.split("%").join("");
            this.showInfoFn();
          }
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //设置参数方法
  private setParamFn(item: any) {
    this.menuWordLen += item.name.length;
    // if (this.editWordLen == 0 && this.menuWordLen > 0) {
    //   this.menuWordLen = this.menuWordLen - 1;
    // }
    // this.curWordLen = this.editWordLen + this.menuWordLen + this.smsSignLen;
    const colorSet = item.params_type == 1 ? "status-orange" : "status-blue";
    const str =
      '<input type="button" class="' +
      colorSet +
      '" value="%' +
      item.name +
      '%">';
    this.editorSms.cmd.do("insertHTML", str);
    item.isForbidden = true;
    item.strHtml = str;
    // this.$set(this.templateContentObj, item.key, "");
  }
  //保存前组装信息方法
  private packageInfoFn() {
    if (this.flag) {
      this.messageStr = this.msgTextHtml;
    }
    this.paramArr.forEach((item: any) => {
      if (item.strHtml && this.msgTextHtml.indexOf(item.strHtml) != -1) {
        const newStr = this.msgTextHtml.replace(
          item.strHtml,
          "%" + item.key + "%"
        );
        this.msgTextHtml = newStr;
      }
    });
  }
  //保存
  private saveFn(isExamine: number, name: string) {
    if (!this.ruleForm.template_short_name) {
      this.$message.warning("请完善信息");
    }
    (this.$refs[name] as any).validate((valid: boolean) => {
      if (valid) {
        this.packageInfoFn();
        this.$nextTick(() => {
          const params: any = {
            template_name: this.ruleForm.template_short_name,
            template_sign_id: this.ruleForm.sign_id,
            message_template: (this.$refs.messageTemText as HTMLElement)
              .innerText,
            template_content: JSON.stringify(this.templateContentObj),
            is_examine: isExamine,
            message_str: this.messageStr,
            text_content: (this.$refs.messageTemText as HTMLElement).innerText
          };
          let api = "";
          if (this.isCreate == "2") {
            api = "smsTemplateEditApi";
            params.id = this.curId;
          } else {
            api = "smsTemplateCreateApi";
          }
          API[api](params).then((res: any) => {
            if (res.code === 200) {
              localStorage.setItem("SMS_TEMPLATE", JSON.stringify({}));
              this.$router.push({
                name: "myapp_customer_auto-market_sms-mgt_list"
              });
            } else {
              this.flag = false;
              this.$message.error(res.message);
            }
          });
        });
      }
    });
  }
  //取消
  private cancelFn() {
    if (
      (this.ruleForm.template_short_name ||
        this.ruleForm.sign_id != 117567463708480 ||
        this.msgTextHtml) &&
      this.isCreate == 1
    ) {
      this.isPop = true;
    } else if (
      (this.oldObj.template_short_name != this.ruleForm.template_short_name ||
        this.oldObj.sign_id != this.ruleForm.sign_id ||
        this.oldObj.message_str != this.msgTextHtml) &&
      this.isCreate == 2
    ) {
      this.isTipsCancel = true;
    } else {
      this.directBcak();
    }
  }
  //关闭和取消
  private goCloseFn() {
    this.isTipsCancel = false;
  }
  //确定
  private saveGoBackFn() {
    this.isTipsCancel = true;
    this.directBcak();
  }
  //返回按钮
  private goBackFn(b: string) {
    if (b == "save") {
      this.packageInfoFn();
      //保存并离开,存储到localStorage
      this.$nextTick(() => {
        localStorage.setItem(
          "SMS_TEMPLATE",
          JSON.stringify({
            template_name: this.ruleForm.template_short_name,
            template_sign_id: this.ruleForm.sign_id,
            message_template: (this.$refs.messageTemText as HTMLElement)
              .innerText
          })
        );
        this.directBcak();
      });
    } else if (b == "cancel") {
      localStorage.setItem("SMS_TEMPLATE", JSON.stringify({}));
      this.directBcak();
    }
    this.isPop = false;
  }
  //直接离开
  private directBcak() {
    this.$router.push({
      name: "myapp_customer_auto-market_sms-mgt_list"
    });
  }
  //获取编辑详情
  private getEditInfoFn() {
    API.smsTemplateInfoApi({ id: this.curId }).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        this.ruleForm.template_short_name = data.template_short_name;
        this.ruleForm.template_name = data.template_name;
        this.ruleForm.sign_id = data.template_sign_id;
        this.showStr = data.message_template.split("%").join("");
        this.oldObj = {
          template_short_name: data.template_short_name,
          template_name: data.template_name,
          sign_id: data.template_sign_id,
          message_str: data.message_str
        };
        if (data.template_content) {
          this.templateContentObj = data.template_content;
        }
        // const lidiaArr = data.message_template.split("%").filter(Boolean);
        // lidiaArr.map((item: string) => {
        //   this.$set(this.templateContentObj, item, "");
        // });
        this.showInfoFn();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //编辑展示详情信息
  private showInfoFn() {
    this.paramArr.forEach((item: any) => {
      const colorSet = item.params_type == 1 ? "status-orange" : "status-blue";
      const getVal = Object.values(this.templateContentObj)[
        Object.keys(this.templateContentObj).indexOf(item.key)
      ];
      const getShowVal = getVal ? getVal : item.name;
      const str =
        '<input type="button" class="' +
        colorSet +
        '" attr="' +
        item.key +
        '" value="%' +
        getShowVal +
        '%">';
      if (this.showStr.indexOf(item.key) != -1) {
        const newStr = this.showStr.replace(item.key, str);
        item.strHtml = str;
        this.showStr = newStr;
        item.isForbidden = true;
        this.menuWordLen += item.name.length;
      }
    });
    this.msgTextStr = this.showStr.split("%").join("");
    this.editorSms.txt.html(this.showStr);
  }
  //短信签名列表
  private smsSignListFn() {
    const params: any = {
      status: 3,
      page: 1,
      size: 1000
    };

    API.smsSignListApi(params).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        this.templateStatus = data.list;
        // this.smsSignName = "【" + data.list[0].sign_name + "】";
        this.curWordLen = this.editWordLen + this.menuWordLen;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //参数提示弹框打开
  private tipsFn() {
    this.isTips = true;
  }
  //参数弹框关闭
  private tipsCloseFn() {
    this.isTips = false;
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    this.isCreate = this.$route.query.type; //1新建  2编辑
    this.curId = Number(this.$route.query.id);
    //获取参数列表
    this.getSmsParamsListFn();
    //获取短信签名列表
    this.smsSignListFn();
  }
  private mounted() {
    this.editorSms = new wangEditor(`#editorSmsInput`);
    // 配置菜单栏，删减菜单，调整顺序
    this.editorSms.config.menus = [];
    // 配置全屏功能按钮是否展示
    this.editorSms.config.showFullScreen = false;
    //编辑区域高度
    this.editorSms.config.height = 210;
    //自定义placeholder提示文字
    this.editorSms.config.placeholder = "";
    //配置 z-index
    this.editorSms.config.zIndex = 1;

    const _this = this;
    // 配置 onchange 回调函数，将数据同步
    this.editorSms.config.onchange = (newHtml: any) => {
      const reg = RegExp(/[\[\]\{\}\【\】]/);
      const bool = newHtml.match(reg);
      newHtml = newHtml.replace(/[\[\]\{\}\【\】]/g, "");
      if (bool) {
        _this.editorSms.txt.html(newHtml);
        return false;
      }
      _this.paramArr.forEach((item: any) => {
        if (
          item.strHtml &&
          item.isForbidden &&
          newHtml.indexOf(item.strHtml) == -1
        ) {
          item.isForbidden = false;
          item.strHtml = "";
          this.menuWordLen -= item.name.length;
        }
      });

      //显示给页面的短信内容
      this.msgTextStr = newHtml.split("%").join("");
      this.msgTextStr = this.msgTextStr.replace(/<br>/g, "");
      setTimeout(() => {
        this.editWordLen = (this.$refs
          .editWordLenDom as HTMLElement).innerText.length;
        //给接口传数据的短信内容
        this.msgTextHtml = newHtml;
        //计算文本长度
        this.curWordLen = this.editWordLen + this.menuWordLen;
      }, 0);
    };
    // 创建编辑器
    this.editorSms.create();
    document.onkeydown = function(e) {
      const ev: any = typeof event != "undefined" ? window.event : e;
      if (ev.keyCode == 13) {
        return false;
      }
    };
  }
}
</script>
<style lang="scss" scoped>
.card-box {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
}
li {
  list-style: none;
}
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
.create-wrap {
  display: flex;
  .sms-show {
    padding: 100px 0;
    flex: 0 0 375px;
    margin-right: 30px;
    background: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/phone.png") no-repeat 100% top;
    .time {
      text-align: center;
      line-height: 60px;
      color: $color-text-normal;
      font-size: 14px;
    }
    .content-msg {
      padding: 10px 20px;
      font-size: 12px;
      background: $color-background;
      margin: 0px 14px;
      line-height: 18px;
      border-radius: 36px;
      color: $color-text-normal;
      max-width: 255px;
      min-height: 18px;
      word-wrap: break-word;
      white-space: normal;
      word-break: break-all;
    }
    .content-msg::after {
      content: "\9000\8ba2\56de\0054";
    }
    .content-msg::before {
      content: attr(attr);
    }
  }
  .sms-set {
    flex: 1;
    border-left: 1px solid $color-third-line;
    &-msg {
      border-radius: 4px;
      border: 1px solid #dcdfe6;
      padding: 0 15px;
      dt {
        height: 50px;
        line-height: 50px;
        span {
          color: $color-main;
          margin-left: 10px;
          cursor: pointer;
        }
      }
      dd.tag-wrap {
        padding-bottom: 15px;
        border-bottom: 1px solid $color-third-line;
        line-height: 34px;
      }
      .sms-content {
        padding: 30px 0;
        color: $color-text-normal;
      }
      &-tips {
        background: #ebebf1;
        margin: 0 -15px;
        padding: 0 15px;
        .pull-left {
          color: $color-text-normal;
        }
        .pull-right {
          color: $color-text-second;
          span {
            color: #f67270;
          }
        }
      }
      .tag-item {
        // margin-right: 5px;
        // margin-bottom: 10px;
        cursor: pointer;
      }
      /deep/.el-tag {
        border-radius: 3px;
      }
    }
    .tips-list {
      padding-top: 20px;
      // margin-bottom: -40px;
      dd {
        line-height: 20px;
        color: $color-text-second;
        margin-bottom: 10px;
      }
    }
  }
}
.send-btn {
  margin-left: 20px;
}

/deep/.in-item {
  width: auto;
}
//button 字体
/deep/.el-button--small {
  font-size: 14px;
}
/deep/.el-button--mini {
  padding: 0 8px !important;
  height: 24px !important;
}
/deep/.el-button.is-disabled:not(.el-button--text) {
  border: 1px solid !important;
}
.tipsText {
  dd {
    line-height: 20px;
    margin: 5px 0;
  }
  padding-bottom: 20px;
}
.new-class {
  max-width: none !important;
}
</style>
<style lang="scss">
/* 覆盖编辑器样式 */
.w-e-text-container {
  border: none !important;
}
.w-e-toolbar {
  display: none !important;
}
/* //覆盖编辑器样式 */
.w-e-text {
  margin: 0;
}
#editorSmsInput {
  margin-bottom: 20px;
}
// #editorSmsInput::before {
//   content: attr(attr);
// }
// #editorSmsInput::after {
//   content: "\9000\8ba2\56de\0054";
//   position: absolute;
//   left: 27px;
// }
#editorSmsInput p {
  word-wrap: break-word;
  white-space: normal;
  word-break: break-all;
}
.create-wrap .sms-set-msg .sms-content {
  padding: 0;
  color: $color-text-normal;
}
#editorSmsInput,
.content-msg {
  input[type="button"] {
    border: none;
    background: none;
    outline: none;
    margin: 0 1px;
    padding: 0;
  }
  .status-blue {
    color: #708cf6 !important;
  }
  .status-orange {
    color: #f6ab3c !important;
  }
}

.content-msg {
  input[type="button"] {
    font-size: 12px !important;
    color: $color-text-normal !important;
    margin: 0;
  }
}
.new-class {
  max-width: none !important;
}
</style>
