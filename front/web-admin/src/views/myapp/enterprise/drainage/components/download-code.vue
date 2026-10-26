<template>
  <div class="add-manual-label">
    <el-dialog
      title="下载联系我二维码"
      :visible.sync="showDialog"
      width="550px"
      @close="closeDialog"
    >
      <div class="code-box">
        <div class="item-code" v-for="(item, index) in typeList" :key="index">
          <div class="code-type">{{ item.name }}</div>
          <div class="type-tip">
            <div class="m-b-10">{{ item.tip1 }}</div>
            <div>{{ item.tip2 }}</div>
          </div>
          <img
            @click="download(item.width, item.height, item.name)"
            class="down-img"
            src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/code-download-icon.png"
            alt=""
          />
        </div>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/drainage-coach";
import { downloadFile } from "@/utils/download";
@Component({})
export default class extends Vue {
  @Prop({ default: "" }) currentDownLoadItem?: any;
  @Prop({ default: false }) zip?: any;
  /************data**************/
  private showDialog = false;
  private typeList: any = [
    {
      size: 1,
      name: "小尺寸",
      tip1: "适用于屏幕类，宣传册类",
      tip2: "边长8cm(258 * 258px)",
      width: 258,
      height: 258
    },
    {
      size: 2,
      name: "中尺寸",
      tip1: "适用于海报，展架等",
      tip2: "边长15cm(430 * 430px)",
      width: 430,
      height: 430
    },
    {
      size: 3,
      name: "大尺寸",
      tip1: "适用于幕布，大型广告等",
      tip2: "边长50cm(1417 * 1417px)",
      width: 1417,
      height: 1417
    }
  ]; // 下载类型列表
  /************methods**************/
  // 时间戳转日期字符
  private timestampToTime(date: number) {
    const year = new Date(date).getFullYear();
    const month = new Date(date).getMonth() + 1;
    const theDate = new Date(date).getDate();
    const hour = new Date(date).getHours();
    const minutes = new Date(date).getMinutes();
    const seconds = new Date(date).getSeconds();
    return (
      year +
      "" +
      (month >= 10 ? month : "0" + month) +
      "" +
      (theDate >= 10 ? theDate : "0" + theDate) +
      "" +
      (hour >= 10 ? hour : "0" + hour) +
      "" +
      (minutes >= 10 ? minutes : "0" + minutes) +
      "" +
      (seconds >= 10 ? seconds : "0" + seconds)
    );
  }
  // 下载二维码
  private download(width: number, height: number, typeName: string) {
    if(this.zip){
      this.downloadZip(width);
      return;
    }
    const params = {
      image_url: this.currentDownLoadItem.qr_code_url,
      width: width,
      height: height
    };
    Api.contactDownloadApi(params).then((res: any) => {
      if (res.code == 200) {
        const timestamp = this.timestampToTime(new Date().getTime());
        const name = `${this.currentDownLoadItem.name}${typeName}${timestamp}.png`;
        const a = document.createElement("a");
        a.href = res.data;
        a.download = name;
        a.click();
        this.$message.success("下载成功");
      } else {
        this.$message.error(res.message);
      }
    });
  }

  async downloadZip(w) {
    const params = {
      width: w,
      id: this.currentDownLoadItem.id
    };
    const res = await Api.contactDownloadZipApi(params)
    await downloadFile(res.data);
    this.$message.success("下载成功");
  }

  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }
  private mounted() {
    this.showDialog = true;
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 416px;
}
/deep/.el-dialog__body {
  max-height: 416px;
}
div {
  box-sizing: border-box;
}
/deep/.el-dialog__title {
  color: #33363c;
}
.code-box {
  .item-code {
    background-color: #f5f6fa;
    border-radius: 10px;
    padding: 20px 30px;
    height: 90px;
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    .code-type {
      font-size: 18px;
      color: #606266;
      font-weight: 600;
    }
    .type-tip {
      color: #606266;
      font-size: 14px;
      line-height: 20px;
      margin-left: 30px;
      width: 320px;
    }
    .down-img {
      width: 25px;
      height: 25px;
      cursor: pointer;
    }
  }
  .item-code:last-child {
    margin-bottom: 10px;
  }
}
</style>
