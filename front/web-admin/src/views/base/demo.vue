<template>
  <div class="container">
    <div class="box">
      <div class="title">文件提取</div>
      <el-input type="textarea" v-model="formData.content" :rows="5" placeholder="请输入文件提取密钥"></el-input>
      <div class="btn-box">
        <el-button type="primary" @click="init">下载</el-button>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { downloadFile } from "@/utils/download";
import { Vue, Component } from "vue-property-decorator";
import COS from 'cos-js-sdk-v5';
import { getPrivateByKeyApi } from "@/api/enterprise/sensitive-word";
import { resolveCosDownloadObject } from "./cos-download-object";
import md5 from "js-md5";
const uatcdn = 'https://cdn-cos-uat.coach.com.cn/'
const prodcdn = 'https://cdn-cos.coach.com.cn/'

@Component({ name: "404" })
export default class extends Vue {
  formData = {
    content: ''
  }
  formRules = {
    content: [
      { required: true, message: '请输入文件链接' },
    ]
  }
  async init() {
    const secret = this.formData.content.trim();
    const key = this.$route.query.secret;

    if(secret !== md5(key)) {
      return this.$message.error('链接无效');
    }

    const res: any = await getPrivateByKeyApi({key});
    const { url, credential } = res.data
    const client = new COS({
      getAuthorization: async function (options, callback) {
        if (res.code !== 200) throw new Error('获取临时密钥失败');
        callback({
          TmpSecretId: credential?.secret_id,
          TmpSecretKey: credential?.secret_key,
          XCosSecurityToken: credential?.session_token,
          StartTime: credential?.start_time,
          ExpiredTime: credential?.expired_time,
        });
      },
    })
    const objectInfo = resolveCosDownloadObject(url, credential, [prodcdn, uatcdn]);
    client.getObjectUrl({
      Bucket: objectInfo.bucket /* 必须 */,
      Region: credential?.region /* 存储桶所在地域，必须字段 */,
      Key: objectInfo.key, // 存储在桶里的对象键（例如1.jpg，a/b/test.txt），支持中文，必须字段
    },
    function (err, data) {
      console.log(err || data.Url);
      if(!err) {
        const filename = url.substring(url.lastIndexOf('/') + 1)
        downloadFile(data.Url, {
          filename
        })
      }
    })
  }
}
</script>
<style lang="scss" scoped>
.container {
  width: 100vw;
  height: 100vh;

  display: flex;
  align-items: center;
  justify-content: center;
  background: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/login/login-bac.png") center/cover no-repeat;

  .box {
    width: 600px;
    box-shadow: 0 0 10px #ccc;
    border-radius: 10px;
    background: #fff;
    padding: 0 50px;

    .title {
      padding: 20px 0;
      font-size: 16px;
      font-weight: bold;
    }

    .btn-box {
      padding: 20px 0;
      text-align: right;
    }

  }
}
</style>
