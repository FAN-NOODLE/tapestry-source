<template>
  <div>
    <layout-form type="normal">
      <template slot="content">
        <div class="create-wrap">
          <div class="sms-show">
            <p class="time">下午 6:30</p>
            <div
              class="content-msg pull-left"
              v-html="msgTextStr"
              :attr="smsSignName"
            ></div>
          </div>
          <div class="sms-set">
            <cy-form ref="cTemplate" :rules="rules" :model="ruleForm">
              <cy-form-item label="模版名称" required>
                <div class="sms-set-name">
                  {{ ruleForm.template_short_name }}
                </div>
              </cy-form-item>
              <cy-form-item label="短信签名" required>
                <cy-select
                  width-size="medium"
                  :disabled="occupancy == 1"
                  size="small"
                  v-model="ruleForm.sign_id"
                  :optionsList="templateStatus"
                  :options="{ label: 'sign_name', value: 'id' }"
                  @change="changeSmsFn"
                ></cy-select>
              </cy-form-item>
              <cy-form-item label="短信内容" class="new-class" required>
                <div class="show-msg">
                  <h3>点击下方内容区变量可进行赋值</h3>
                  <dl class="sms-set-msg">
                    <dd
                      class="sms-content"
                      ref="setParamDom"
                      v-html="msgStrHtml"
                    ></dd>
                  </dl>
                  <div class="sms-set-msg-tips clearfix">
                    <div class="pull-left">
                      已占用{{ totalCount }}个字，预计计费为
                      {{ Math.ceil(totalCount / 70) }}条短信
                    </div>
                    <div class="pull-right">
                      <span>*</span> 计费规则：1条短信70个字
                    </div>
                  </div>
                </div>
                <dl class="tips-list">
                  <dd>
                    1.您可在短信中插入动态参数和固定参数，如：%用户昵称%，%为英文格式符号，固定参数可在审核通过后赋值
                  </dd>
                  <dd>2.短信内容中不可包括【】、[]、测试、test字符。</dd>
                  <dd>
                    3.您填写的链接首尾都需要包含一个空格，且以http://或https://格式开头。
                  </dd>
                  <dd>
                    4.短信模板审核：请在内容创建之前给到相关人员进行人工审核。
                  </dd>
                  <dd>
                    5.由于变量和签名等原因，实际发送短信字数与短信模板字数会有不一致，实际短信计费条数以运营商返回为准，请在大批量发送前做好测试。
                  </dd>
                </dl>
              </cy-form-item>
              <cy-form-item label="发送预览" prop="telNumber" phone
                ><cy-input
                  width-size="small"
                  placeholder="请输入预览手机号码"
                  v-model="ruleForm.telNumber"
                  maxlength="11"
                  minlength="0"
                  show-word-limit
                  @input="_checkInt(ruleForm, 'telNumber', ruleForm.telNumber)"
                ></cy-input>
                <el-button
                  class="send-btn"
                  type="primary"
                  size="small"
                  @click="sendSmsFn"
                  >测试发送</el-button
                >
              </cy-form-item>
            </cy-form>
          </div>
        </div>
      </template>
      <template slot="footer">
        <!-- <el-button
          :disabled="template_sign_id == ruleForm.sign_id"
          type="primary"
          @click="saveFn('change')"
          >保存并提审</el-button
        > -->
        <el-button type="primary" @click="saveFn('none')">保存</el-button>
        <el-button @click="cancelFn">取消</el-button>
      </template>
    </layout-form>
    <set-param
      v-if="isShowSetMsg"
      :curPopObj="curPopObj"
      @closeDialog="setMsgFn"
    ></set-param>
    <el-dialog
      title="离开页面提示"
      :visible.sync="isTips"
      width="420px"
      @close="goCloseFn"
    >
      <span
        ><i aria-hidden="true" class="el-icon-warning"></i>
        你确定要放弃编辑吗？放弃后编辑的信息将不会被保存</span
      >
      <span slot="footer" class="dialog-footer">
        <el-button @click="goCloseFn">取消</el-button>
        <el-button type="primary" @click="saveGoBackFn">确定</el-button>
      </span>
    </el-dialog>
    <div ref="messageTemText" style="display:none" v-html="msgTextHtml"></div>
  </div>
</template>
<script lang="ts">
import { Component, Mixins } from "vue-property-decorator";
import API from "@/api/customer";
import setParam from "../../components/set-param.vue";
import CheckNum from "@/mixins/NumberInput";

@Component({
  components: { setParam }
})
export default class extends Mixins(CheckNum) {
  /*************************** data start ***************************/
  private editorSms: any = null; //编辑器
  private rules = {};
  //表单
  private ruleForm = {
    template_short_name: "",
    template_name: "",
    sign_id: "畅移科技",
    telNumber: ""
  };
  private template_sign_id = ""; // 接口返回的短信签名
  private msgTextStr = ""; //显示的短信内容
  private msgTextHtml = ""; //组装的短信内容
  private msgStrHtml = ""; //右侧显示的短信内容
  private occupancy = 0; // 短信模版是否被占用
  //参数
  private paramArr: any = [];
  private msgArr: any = [];
  private isShowSetMsg = false; //是否显示设置弹框
  private curId = 0; //当前id
  private curPopObj: any = {}; //当前需要改参数
  private templateContentObj: any = {}; //设置的参数
  private templateStatus = []; //模板状态
  private curWordLen = 4; //短信长度
  private signCount = 0; //短信签名长度
  private totalCount = 0; //整体短信长度
  private smsSignName = ""; //短信签名
  private oldObj: any = {}; //取消弹框存储旧数据
  private isTips = false; //是否提示框
  private selectDataType = ""; //日期类型
  /*************************** method start ***************************/
  //select change事件
  private changeSmsFn() {
    this.templateStatus.forEach((item: any) => {
      if (item.id == this.ruleForm.sign_id) {
        this.smsSignName = "【" + item.sign_name + "】";
      }
    });
    this.signCount = this.smsSignName.length;
    this.totalCount = this.curWordLen + this.signCount;
  }
  //获取参数列表
  private getSmsParamsListFn() {
    API.smsParamsListApi({}).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        this.paramArr = data.list;
        this.paramArr.forEach((item: any) => {
          item.strHtml = "";
          item.type = "dynamic";
        });
        this.getEditInfoFn();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //点击事件
  private domClickFn() {
    setTimeout(() => {
      const objs = (this.$refs.setParamDom as HTMLElement).getElementsByTagName(
        "input"
      );
      let attr: any = "";
      for (let i = 0; i < objs.length; i++) {
        objs[i].addEventListener("click", () => {
          attr = objs[i].getAttribute("attr");
          if (
            attr == "user_name" ||
            attr == "address_name" ||
            attr == "order_id"
          ) {
            return false;
          } else {
            if (objs[i].value.indexOf("至") != -1) {
              this.selectDataType = "2";
            } else {
              this.selectDataType = "1";
            }
            this.curPopObj = {
              type: objs[i].getAttribute("attr"),
              name: objs[i].value,
              selectData: this.selectDataType
            };
            this.paramArr.forEach((item: any) => {
              if (objs[i].getAttribute("attr") == item.key) {
                this.curPopObj.title = item.name;
              }
            });
            this.isShowSetMsg = true;
          }
        });
      }
    }, 0);
  }
  //设置参数弹框
  private setMsgFn(obj?: any) {
    if (obj) {
      this.selectDataType = obj.selectData; //日期类型
      //计算字数
      const count =
        obj.val.length - this.curPopObj.name.split("%").join("").length;
      this.curWordLen += count;
      this.totalCount = this.curWordLen + this.signCount;
      if (this.msgStrHtml.indexOf(obj.type) != -1) {
        const newStr = this.msgStrHtml.replace(obj.type, "%" + obj.val + "%");
        this.msgStrHtml = newStr;
        this.msgTextHtml = newStr;
        this.msgTextStr = newStr.split("%").join("");
        this.templateContentObj[obj.num] = obj.val;
        this.paramArr.forEach((item: any) => {
          const colorSet =
            item.params_type == 1 ? "status-orange" : "status-blue";
          if (item.key == obj.num) {
            item.strHtml =
              '<input type="button" class="' +
              colorSet +
              '" attr="' +
              obj.num +
              '" value="%' +
              obj.val +
              '%">';
          }
        });
      }
      this.domClickFn();
    }
    this.isShowSetMsg = !this.isShowSetMsg;
  }
  //保存
  private saveFn(type: string) {
    const messageStr = this.msgTextHtml;
    this.paramArr.forEach((item: any) => {
      if (item.strHtml && this.msgTextHtml.indexOf(item.strHtml) != -1) {
        const newStr = this.msgTextHtml.replace(
          item.strHtml,
          "%" + item.key + "%"
        );
        this.msgTextHtml = newStr;
      }
    });
    let isAllParam = true;
    console.log("1111111111", this.templateContentObj);
    for (const key in this.templateContentObj) {
      if (
        !this.templateContentObj[key] &&
        key !== "user_name" &&
        key !== "address_name" &&
        key !== "order_id"
      ) {
        console.log(this.templateContentObj[key], "ffffff", key);
        isAllParam = false;
        break;
      }
    }
    if (!isAllParam) {
      this.$message.error("请完善所有需要赋值的变量");
      return false;
    }
    setTimeout(() => {
      const params: any = {
        id: this.curId,
        template_name: this.ruleForm.template_name,
        template_sign_id: this.ruleForm.sign_id,
        template_content: JSON.stringify(this.templateContentObj),
        message_template: (this.$refs.messageTemText as HTMLElement).innerText,
        message_str: messageStr,
        is_examine: type == "change" ? 1 : 0
      };
      API.smsTemplateEditApi(params).then((res: any) => {
        if (res.code === 200) {
          this.$message.success("保存成功");
          this.$router.push({
            name: "myapp_customer_auto-market_sms-mgt_list"
          });
        } else {
          this.$message.error(res.message);
        }
      });
    }, 0);
  }
  //获取编辑详情
  private getEditInfoFn() {
    API.smsTemplateInfoApi({ id: this.curId }).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        this.oldObj = {
          message_str: data.message_str,
          sign_id: data.template_sign_id
        };

        this.ruleForm.template_short_name = data.template_short_name;
        this.ruleForm.template_name = data.template_name;
        this.ruleForm.sign_id = data.template_sign_id;
        this.template_sign_id = data.template_sign_id;
        this.occupancy = JSON.parse(data.occupancy);
        this.templateStatus.forEach((item: any) => {
          if (item.id == data.template_sign_id) {
            this.smsSignName = "【" + item.sign_name + "】";
          }
        });
        this.signCount = this.smsSignName.length;
        let showStr = data.message_template.split("%").join("");
        if (data.template_content) {
          this.templateContentObj = JSON.parse(data.template_content);
        }
        let msgCount = 0; //标签字数
        this.paramArr.forEach((item: any) => {
          const colorSet =
            item.params_type == 1 ? "status-orange" : "status-blue";
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
          if (showStr.indexOf(item.key) != -1) {
            const newStr = showStr.replace(item.key, str);
            item.strHtml = str;
            showStr = newStr;
            if (!data.template_content) {
              this.templateContentObj[item.key] = "";
            }
            msgCount += getShowVal.length;
          }
        });
        this.msgTextStr = showStr.split("%").join("");
        this.msgStrHtml = showStr;
        this.msgTextHtml = showStr;
        this.$nextTick(function() {
          this.curWordLen +=
            (this.$refs.messageTemText as HTMLElement).innerText.length +
            msgCount;
          this.totalCount = this.curWordLen + this.signCount;
        });
        this.domClickFn();
      } else {
        this.$message.error(res.message);
      }
    });
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
        //获取参数列表
        this.getSmsParamsListFn();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //发送短信
  private sendSmsFn() {
    const params: any = {
      phone: this.ruleForm.telNumber,
      template_id: this.curId
    };
    API.smsSendApi(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("发送成功！");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //取消
  private cancelFn() {
    if (
      this.oldObj.sign_id != this.ruleForm.sign_id ||
      this.oldObj.message_str != this.msgTextHtml
    ) {
      this.isTips = true;
    } else {
      this.$router.push({
        name: "myapp_customer_auto-market_sms-mgt_list"
      });
    }
  }
  //关闭和取消
  private goCloseFn() {
    this.isTips = false;
  }
  //确定
  private saveGoBackFn() {
    this.isTips = true;
    this.$router.push({
      name: "myapp_customer_auto-market_sms-mgt_list"
    });
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    this.curId = Number(this.$route.query.id);
    //获取短信签名列表
    this.smsSignListFn();
  }
}
</script>
<style lang="scss" scoped>
.new-class {
  max-width: none !important;
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
      background: #ebebf1;
      margin: 0px 14px;
      line-height: 18px;
      border-radius: 36px;
      color: $color-text-normal;
      max-width: 255px;
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
    .show-msg {
      h3 {
        color: $color-text-normal;
        line-height: 40px;
        font-size: 14px;
      }
    }
    &-name {
      line-height: 40px;
      color: $color-text-normal;
    }
    &-msg {
      border-radius: 4px 4px 0 0;
      border: 1px solid #dcdfe6;
      border-bottom: none;
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
      // .sms-content::after {
      //   content: "\9000\8ba2\56de\0054";
      // }
      // .sms-content::before {
      //   content: attr(attr);
      // }
      &-tips {
        background: #ebebf1;
        padding: 0 15px;
        border-radius: 0 0 4px 4px;
        border: 1px solid #dcdfe6;
        border-top: none;
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
  padding: 0 8px;
  height: 24px;
}
</style>
<style lang="scss">
.new-class {
  max-width: none !important;
}
.create-wrap .sms-set-msg .sms-content,
.create-wrap .content-msg {
  // padding: 0;
  // color: $color-text-normal;
  input[type="button"] {
    border: none;
    background: none;
    outline: none;
    margin: 0 1px;
    cursor: pointer;
    padding: 0;
  }
  .status-blue {
    color: #708cf6 !important;
  }
  .status-orange {
    color: #f6ab3c !important;
  }
}
.create-wrap .content-msg {
  input[type="button"] {
    font-size: 12px !important;
    color: $color-text-normal !important;
    margin: 0;
  }
}
</style>
