<template>
  <div class="friend-circle-edit">
    <layout-form>
      <template slot="content">
        <el-form label-width="100px" label-position="left">
          <el-form-item label="任务描述:" required>{{ tempInfo.template_desc }}</el-form-item>
          <el-form-item label="任务名称:" required>{{ tempInfo.template_name }}</el-form-item>
          <el-form-item label="任务类型:" required>{{ tempInfo.task_category === 0 ? '沟通任务' : '朋友圈任务' }}</el-form-item>
          <el-form-item v-if="tempInfo.task_category === 0 " label="群发类型:" required>{{ tempInfo.group_send_type === 1 ? '企微一键群发（从群发助手发送）' : '灵活沟通转发（从任务中心发送）' }}</el-form-item>
          <el-form-item v-if="tempInfo.task_category === 0 && tempInfo.group_send_type === 2" label="完成标准:" required>顾客名单触达率达到{{ tempInfo.finish_rate }}%</el-form-item>
          <el-form-item v-if="tempInfo.task_category === 0 && tempInfo.group_send_type === 1" label="群发客户:" required>{{ tempInfo.allow_select ? '导购可编辑' : '导购不可编辑' }}</el-form-item>
          <el-form-item label="任务内容:" required>
            <div>
              <div>{{ tempInfo.content.text.content }}</div>
              <task-content-detail :content="tempInfo.content" @openVideo="openVideo"></task-content-detail>
            </div>
          </el-form-item>
          <el-form-item label="备选模版:" class="custom">
            <el-switch v-model="tempInfo.open_back_up" :active-value="1" :inactive-value="0" disabled></el-switch>
          </el-form-item>
          <div v-if="tempInfo.open_back_up" style="padding-left: 100px;">
            <el-form-item label-width="120px" label="生效时间:" required>
              <div v-if="tempInfo.template_back.end_time">
                <span class="m-r-10">在</span>
                {{ tempInfo.template_back.start_time | formatDate('yyyy-mm-dd') }}
                至
                {{ tempInfo.template_back.end_time | formatDate('yyyy-mm-dd') }}
                <span class="m-l-10">期间生效</span>
              </div>
              <div v-if="!tempInfo.template_back.end_time">
                <span class="m-r-10">在</span>
                {{ tempInfo.template_back.start_time | formatDate('yyyy-mm-dd') }}
                <span class="m-l-10">日，开始生效。</span>
              </div>
            </el-form-item>
            <el-form-item v-if="tempInfo.task_category === 0 " label-width="120px" label="任务群发类型:" required>{{ tempInfo.template_back.group_send_type === 1 ? '企微一键群发（从群发助手发送）' : '灵活沟通转发（从任务中心发送）' }}</el-form-item>
            <el-form-item label-width="120px" v-if="tempInfo.task_category === 0 && tempInfo.template_back.group_send_type === 2" label="任务完成标准:" required>顾客名单触达率达到{{ tempInfo.template_back.finish_rate }}%</el-form-item>
            <el-form-item label-width="120px" v-if="tempInfo.task_category === 0 && tempInfo.template_back.group_send_type === 1" label="群发客户:" required>{{ tempInfo.template_back.allow_select ? '导购可编辑' : '导购不可编辑' }}</el-form-item>
            <el-form-item label-width="120px" label="任务内容:" required>
              <div v-if="tempInfo.template_back.content">
                <div>{{ tempInfo.template_back.content.text.content || '--' }}</div>
                <task-content-detail :content="tempInfo.template_back.content" @openVideo="openVideo"></task-content-detail>
              </div>
              <div v-else>--</div>
              <div style="color: red; font-size: 12px;">※必若未配置备选模板的任务素材内容，则默认发送原模板任务素材内容</div>
            </el-form-item>
          </div>
        </el-form>
      </template>
      <template slot="footer">
        <el-button @click="cancel()">返回</el-button>
      </template>
    </layout-form>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="isOpenVideo"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/temp";
import openVideo from "@/views/myapp/enterprise/material-library/components/com-open-video.vue";
import TaskContentDetail from "@/views/myapp/wisdom/task-center/task-temp/components/task-content-detail.vue";
function parseResult(templateContents: any) {

  // 这里需要重新组装数据
  let newContents: any = {};
  newContents = templateContents;
  // 处理数据
  if (newContents.text.content) {
    newContents.text = {
      content: newContents.text.content
    };
  }
  if (newContents.image) {
    newContents.picture = newContents.image.imgList.map((item: any) => {
      return {
        ...item,
        url: item.pic_url
      };
    });
  } else {
    newContents.picture = []
  }
  if (newContents.video) {
    newContents.video = {
      ...newContents.video,
      url: newContents.video.video_url
    };
  }
  if (newContents.applets) {
    newContents.applet = newContents.applets;
  }
  if (newContents.link) {
    newContents.link = {
      ...newContents.link,
      cover_url: newContents.link.pic_url
    };
  }
  console.log(111111, newContents);
  
  return newContents;
}


@Component({
  components: {
    openVideo,
    TaskContentDetail
  }
})
export default class extends Vue {
  isOpenVideo = false;
  videoUrl = "";
  tempInfo: any = {
    template_name: "",
    template_desc: "",
    content: {}
  };
  // 点击取消按钮
  cancel() {
    this.$router.go(-1);
  }
  created() {
    this.getTempInfo();
  }
  // 获取任务模板详情
  getTempInfo() {
    const { id } = this.$route.query;
    Api.getTaskTempDetail({ id: Number(id) }).then((res: any) => {
      if ( res.data.template_content) {
         res.data.content = parseResult(JSON.parse( res.data.template_content));
      }
      if(res.data.template_back && res.data.template_back.template_content) {
        res.data.template_back.content = parseResult(JSON.parse(res.data.template_back.template_content));
      }
      this.tempInfo = res.data;
      

    });
  }
  // 预览视频
  openVideo(url: string) {
    this.videoUrl = url;
    this.isOpenVideo = true;
  }
  // 关闭预览视频
  cancelVideo() {
    this.isOpenVideo = false;
  }
}
</script>
<style lang="scss" scoped>
div {
  box-sizing: border-box;
}
.content {
  .font-re {
    display: flex;
    line-height: 38px;
    .title {
      min-width: 92px;
    }
    .content {
      background: none;
      margin-left: 10px;
      box-shadow: none;
      padding: 0;
    }
  }
}

.el-switch {
  font-size: 0;
}

</style>
