<template>
  <el-dialog
    :title="{1:'下载活动码',2:'复制链接',3:'复制活动页面参数'}[type]"
    :visible.sync="visible"
    append-to-body
    @close="hide"
  >
    <div v-if="type===1" class="tip-wrap">
      <div class="title">活动码尺寸说明：</div>
      <div class="info m-t-5">【小尺寸】： 适用于屏幕类，宣传册类等，边长5cm左右（258*258px）</div>
      <div class="info">【中尺寸】：适用于海报，展架等，边长8cm左右（430*430px）</div>
      <div class="info">【大尺寸】：适用于幕布，大型广告等，边长15cm左右（1417*1417px）</div>
    </div>
    <div v-if="type===3" class="tip-wrap">
      <div class="info">此处下载内容为该投放活动的页面参数，请您复制对应渠道的参数，前往EC小程序管理后台生成对应的小程序太阳码和页面链接</div>
    </div>
    <div class="list m-t-15" v-for="item,index in detail.category_list" :key="index">
      <div>渠道{{ index+1 }}: </div>
      <el-row class="item">
        <div class="item-name">{{ item.name }}</div>
        <el-button v-if="type==1" type="primary" icon="el-icon-download" size="small" @click="download(item, '258x258')">小尺寸</el-button>
        <el-button v-if="type==1" type="primary" icon="el-icon-download" size="small" @click="download(item, '430x430')">中尺寸</el-button>
        <el-button v-if="type==1" type="primary" icon="el-icon-download" size="small" @click="download(item, '1417x1417')">大尺寸</el-button>
        <el-button v-if="type==2" type="primary" icon="el-icon-copy-document" size="small" @click="copy(item.link_url)">复制链接</el-button>
        <el-button v-if="type==3" type="primary" icon="el-icon-copy-document" size="small" @click="copy(item.query_str)">复制参数</el-button>
      </el-row>
      <el-row class="item" v-for="point,pi in item.point_list" :key="pi">
        <div class="item-name">{{ point.name }}</div>
        <el-button v-if="type==1" type="primary" icon="el-icon-download" size="small" @click="download(point, '258x258', item)">小尺寸</el-button>
        <el-button v-if="type==1" type="primary" icon="el-icon-download" size="small" @click="download(point, '430x430', item)">中尺寸</el-button>
        <el-button v-if="type==1" type="primary" icon="el-icon-download" size="small" @click="download(point, '1417x1417', item)">大尺寸</el-button>
        <el-button v-if="type==2" type="primary" icon="el-icon-copy-document" size="small" @click="copy(point.link_url)">复制链接</el-button>
        <el-button v-if="type==3" type="primary" icon="el-icon-copy-document" size="small" @click="copy(point.query_str)">复制参数</el-button>
      </el-row>
    </div>
  </el-dialog>
</template>
<script lang="ts"> 
import { getActivityDetailApi } from "@/api/enterprise/activity";
import { downloadFile } from "@/utils/download";
import dayjs from "dayjs";
import { Vue, Component, Prop, Watch, PropSync } from "vue-property-decorator";

@Component
export default class DownloadPopup extends Vue {
  visible = false;
  detail:any = {
    category_list: []
  }
  // 1- 二维码，2-链接，3-活动页面参数
  type = 1

  async show(id:any, type = 1) {
    const res = await getActivityDetailApi({id})
    this.detail = res.data
    this.visible = true;
    this.type = type
  }

  hide() {
    this.visible = false;
  }
  // 复制
  copy(data: string) {
    const url = data;
    const oInput: HTMLTextAreaElement = document.createElement("textarea");
    oInput.value = url;
    document.body.appendChild(oInput);
    oInput.select(); // 选择对象;
    document.execCommand("Copy"); // 执行浏览器复制命令
    this.$message({
      message: "复制成功",
      type: "success"
    });
    oInput.remove();
  }
  download(item, size, parent = null) {
    let filename = [this.detail.name, parent?.name, item.name, size, dayjs().format('YYYYMMDD')].filter(Boolean).join('_');
    
    downloadFile(`${item.code_url}?imageMogr2/crop/${size}`, {
      filename: `${filename}.png`
    })
    
  }
}
</script>

<style lang="scss" scoped>

.tip-wrap {
  padding: 15px;
  // border-radius: 8%;
  background: #f6f6ff;
  .title {
    font-size: 14px;
    font-weight: 500;
    color: #333333;

  }

  .info {
    font-size: 14px;
    font-weight: 400;
    color: #666666;

  }
}

.list {
  // border-radius: 8px;
  overflow: hidden;
}

.item {
  height: 60px;
  display: flex;
  align-items: center;
  background-color: #f7f7fe;
  padding: 0 20px;
  border-bottom: 1px solid #f1f2ff;

  .item-name {
    flex: 1;
  }

  &:nth-child(2) {
    background-color: #f1f2ff;
  }
}
</style>