<template>
  <div class="add-edit-group">
    <layout-form>
      <template slot="content">
        <cy-form
          class="main-box"
          ref="cform"
          :model="ruleForm"
          label-width="140px"
          :rules="rules"
        >
          <cy-form-item label="素材分组">
            <el-cascader
              v-model="ruleForm.group_id"
              :options="groupList"
              :props="{
                label: 'group_name',
                value: 'id',
                checkStrictly: true,
                emitPath: false
              }"
              :style="{
                width: '400px'
              }"
            ></el-cascader>
          </cy-form-item>
          <cy-form-item
            label="素材标签"
            prop="category_ids"
            trigger="change"
          >
            <el-select
              style="width:350px"
              v-model="ruleForm.category_ids"
              multiple
              placeholder="请选择素材标签"
            >
              <el-option
                v-for="item in category_options"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              >
              </el-option>
            </el-select>
            <el-button
              style="font-size:20px;transform: translateY(3px);color: #000000;"
              class="m-l-20"
              type="text"
              icon="el-icon-circle-plus"
              @click="openAddCate"
            ></el-button>
          </cy-form-item>
          <cy-form-item
            label="有效期"
            prop="datePicker"
            required
            trigger="change"
          >
            <el-date-picker
              v-model="ruleForm.datePicker"
              :default-time="['00:00:00', '23:59:59']"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              clearable
              value-format="timestamp"
              @change="datePickerChange"
            >
            </el-date-picker>
          </cy-form-item>
          <div class="right-wrap">
            <div class="left-label">素材内容</div>
            <cy-form-item label="文本内容">
              <!-- 新增素材展示区 -->
              <div class="risk-content-box">
                <!-- 表情选择区域-->
                <div class="select-box-insert">
                  <add-emoji
                    @select="addEmoji($event)"
                    class="add-emoji"
                  >
                    <template slot="reference">
                      <el-button
                        class="insert-btn"
                        size="small"
                      ><img
                          class="btn-img"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/face_icon.png"
                          alt=""
                        /><span>插入表情</span></el-button>
                    </template>
                  </add-emoji>
                </div>
                <!--多行文本区域-->
                <div>
                  <el-input
                    ref="textarea"
                    class="textarea-box"
                    placeholder="请输入消息内容"
                    v-model="ruleForm.text_content"
                    :rows="4"
                    type="textarea"
                    maxlength="500"
                    show-word-limit
                  ></el-input>
                </div>
              </div>
            </cy-form-item>
            <cy-form-item
              label="素材附件"
              prop="file_type"
              trigger="change"
            >
              <div class="material">
                <div class="material-radio">
                  <el-radio-group
                    v-model="ruleForm.file_type"
                    @change="changeFileType"
                    :disabled="dialogType === 2"
                  >
                    <el-radio :label="1">图片</el-radio>
                    <el-radio :label="2">视频</el-radio>
                  </el-radio-group>

                </div>
                <div class="material-button">
                  <!-- 本地上传 素材库选择按钮 -->
                  <upload-image-button
                    v-if="ruleForm.file_type === 1"
                    class="upload-template"
                    :urls="ruleForm.sources"
                    accept="image/jpeg,image/png"
                    img-type="batch"
                    :before-upload="beforeUpload"
                    :limit="9"
                    multiple
                    :disabled="dialogType === 2"
                    @update="update"
                    @updateAI="updateAIImage"
                  >
                    <el-button
                      type="primary"
                      size="small"
                      :disabled="dialogType === 2"
                    >本地上传</el-button>
                  </upload-image-button>
                  <large-file-upload
                    ref="upload"
                    v-if="ruleForm.file_type === 2 && dialogType === 1"
                    v-model="ruleForm.urls"
                    accept="video/mp4"
                    :beforeUpload="beforeUploadFn"
                    :onSuccess="onSuccessFn"
                    :onRemove="removeFn"
                    :disabled="ruleForm.urls.length > 0 || dialogType === 2"
                  >
                    <el-button
                      type="primary"
                      size="small"
                      :disabled="ruleForm.urls.length > 0 || dialogType === 2"
                    >本地上传</el-button>
                  </large-file-upload>
                  <div class="material-button-box">

                    <el-button
                      v-if="ruleForm.file_type !== 0 && dialogType === 1"
                      size="small"
                      type="primary"
                      @click="openSelect"
                      :disabled="dialogType === 2 || (ruleForm.file_type === 2 && ruleForm.urls.length > 0)"
                    >从素材中选择</el-button>
                    <el-checkbox v-if="ruleForm.file_type === 1 && dialogType !== 2 && ruleForm.sources.length > 0" v-model="isAllAi" @change="handleAllAiChange">全选AI识图</el-checkbox>
                  </div>
                </div>
                <!-- 右侧视频选择展示区域-->
                <div
                  class="select-video-box"
                  v-if="ruleForm.file_type === 2 && ruleForm.urls.length > 0"
                >
                  <img
                    @click="
                      openVideo(ruleForm.urls[0])
                      "
                    class="video-play-icon"
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
                    alt=""
                  />
                  <video
                    width="80"
                    height="80"
                    x5-video-player-fullscreen="true"
                    x5-playsinline
                    playsinline
                    webkit-playsinline
                    muted
                    preload="auto"
                    :src="ruleForm.urls[0]"
                  ></video>
                  <div
                    v-if="dialogType === 1"
                    class="delete-video-box"
                    @click="deleteMaterial"
                  >
                    <img
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/material-delete-icon.png"
                      alt=""
                    />
                  </div>
                </div>
                <div
                  class="material-tip"
                  :style="{ marginTop: ruleForm.file_type === 1 ? '0' : '22px' }"
                >
                  <ul>
                    <li>1、最多支持9张图片 ，或者1个视频。</li>
                    <li>
                      2、选择【图片】素材：单张普通图片尺寸建议不超过1440*1080px，长图片，长边建议不超过10800px；最多支持传入9张图片，总像素不超过1555200。图片大小不超过10M。
                    </li>
                    <li>3、选择【视频】素材：最长不超过30S，最大不超过10MB。</li>
                  </ul>
                </div>
              </div>
            </cy-form-item>
            <cy-form-item :label="ruleForm.file_type === 2 ? '关联商品':'AI识图关联商品'" >
              <div v-if="ruleForm.file_type == 1" class="tip-text">勾选识图后可自动加载</div>
              <select-goods v-if="ruleForm.file_type == 2 && (!ruleForm.sources[0] || ruleForm.sources[0].relation_goods.length ==0)" :value="[]" :disabled="ruleForm.sources.length === 0" @confirm="handleSelectGoodsConfirm(ruleForm.sources[0], $event)"></select-goods>
              <!-- {{  ruleForm.sources.length }} -->
              <template v-for="item in ruleForm.sources">
                <relation-goods  :key="item.hash || item.id" v-if="item.is_ai || (item.relation_goods && item.relation_goods.length > 0)" :item="item" @confirm="handleSelectGoodsConfirm(item, $event)" @delete="handleDeleteGoods(item, $event)"></relation-goods>
              </template>
            </cy-form-item>
          </div>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button
          @click="closeDialog"
          size="small"
        >取 消</el-button>
        <el-button
          type="primary"
          @click="submitFn"
          size="small"
        >确 认</el-button>
      </template>
    </layout-form>
    <add-category
      v-if="addVisible"
      @closeVisible="closeVisible"
      @submit="addCateSubmit"
    ></add-category>
    <select-image
      v-if="showMaterialModal === 'image'"
      :groupSendType="1"
      :isFriend="true"
      @submit="getSelectImgMaterial"
    ></select-image>
    <select-video
      v-if="showMaterialModal === 'video'"
      :groupSendType="1"
      :isFriend="true"
      @submit="getSelectVideoMaterial"
    ></select-video>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="isOpenVideo"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import UploadImageButton from "@/components/business/upload-image/upload-image-button.vue";
import LargeFileUpload from "@/components/business/large-file-upload/index-new.vue";
import addCategory from "../components/add-category.vue";
import Times from "../components/datePicker";
// 品牌store
import PerModule from "@/store/modules/permission";
// 选择表情弹窗
import AddEmoji from "@/views/myapp/enterprise/components/add-emoji.vue";
// 选择图片素材弹窗
import SelectImage from "../components/pic-mterial-dialog.vue";
// 选择视频素材弹窗
import SelectVideo from "../components/video-material-dialog.vue";
// 视频预览
import OpenVideo from "../components/com-open-video.vue";
import SelectGoods from "../components/select-goods.vue";
import RelationGoods from "../components/relation-goods.vue";
import { getVideoDuration } from "@/utils/upload/utils";

interface T {
  group_id: number;
  url: string;
  id: number;
}
@Component({
  components: {
    UploadImageButton,
    addCategory,
    AddEmoji,
    SelectImage,
    SelectVideo,
    LargeFileUpload,
    OpenVideo,
    SelectGoods,
    RelationGoods
  }
})
export default class extends Vue {
  [x: string]: any;
  //------------------------------------------------------data-----------------------------------------------------------/
  dialogType: any = ""
  editRowId: any = ""
  show_dialog = true;
  testFlag = false;
  groupList: any = []; //素材列表
  ruleForm: any = {
    group_id: 0,
    name: "",
    urls: [],
    category_ids: [], //素材标签
    datePicker: Times.times,
    start_time: Times.start, //有效期开始时间
    end_time: Times.end, //有效期结束时间
    text_content: '',//文本内容
    // todo: 1. 单张展示 2. 合并展示
    show_type: 1,
    // todo:素材集合
    sources: [],
    // todo: 1. 图片 2. 视频
    file_type: 0
  };
  rules = {};
  fileObj: any = {}; //上传文件信息
  fileSuccessObj: any = {}; //上传成功信息
  // 素材标签
  category_options: any[] = [];
  addVisible = false; // 控制添加素材标签弹窗
  // todo: 当前品牌
  epName = PerModule.getEpName;
  // todo: 弹窗类型
  showMaterialModal = "";
  // todo: 是否打开视频预览弹窗
  isOpenVideo = false;
  // 是否全选AI识图
  isAllAi = false;
  //----------------------------------------------------methods-----------------------------------------------------------/
  /**
   * @description: 取消视频预览
   * @return {*}
   */
  cancelVideo() {
    this.isOpenVideo = false;
  }
  /**
   * @description: 打开视频预览弹窗
   * @return {*}
   */
  openVideo(url: string) {
    this.isOpenVideo = true;
    this.videoUrl = url;
  }
  /**
   * @description: 删除素材
   * @return {*}
   */
  deleteMaterial() {
    this.ruleForm.urls = [];
    (this.$refs.upload as any).remove();
  }
  /**
   * @description: 文件类型判断
   * @return {*}
   */
  changeFileType() {
    // 清空已经上传的数据
    this.ruleForm.urls = []
    this.ruleForm.sources = []
  }
  /**
   * @description: 获取选择的视频素材
   * @return {*}
   */
  getSelectVideoMaterial(data: any) {
    console.log(data)
    // 如果有视频素材
    if (data) {
      if (this.ruleForm.urls.length > 0) {
        this.$message.warning("视频不能超过1个");
        return;
      } else {
        this.ruleForm.urls = [data.url]
        this.ruleForm.sources = [{
          id: data.id,
          name: data.name,
          url: data.url,
          sort: 0,
          video_cover_url: data.cover_url,
          relation_goods: []
        }]
      }
    }
    this.showMaterialModal = "";
  }
  /**
   * @description: 获取选择的图片素材
   * @return {*}
   */
  getSelectImgMaterial(data: any) {
    if (data && data.length > 0) {
      // todo:判断不能超过9张
      if (this.ruleForm.urls.length + data.length > 9) {
        this.$message.warning("图片不能超过9张");
        return;
      }
      // todo:获取素材集合
      const selectUrls = data.map((item: any) => item.url);
      // todo: sources数据
      // todo:合并urls, sources
      const selectSources = data.map((item: any, index: any) => {
        return {
          id: item.id,
          name: item.name,
          url: item.url,
          sort: index,
          relation_goods: []
        }
      });
      console.log(selectSources)
      this.ruleForm.urls = [...this.ruleForm.urls, ...selectUrls];
      this.ruleForm.sources = [...this.ruleForm.sources, ...selectSources];
    }
    this.showMaterialModal = "";

    this.isAllAi = this.ruleForm.sources.every((item: any) => item.is_ai);
  }
  /**
   * @description: 打开素材库选择弹窗
   * @return {*}
   */
  openSelect() {
    if (this.ruleForm.file_type === 1) {
      this.showMaterialModal = "image";
    } else {
      this.showMaterialModal = "video";
    }
  }
  /**
   * @description: 添加表情
   * @param {*} emj
   * @return {*}
   */
  async addEmoji(emj: any) {
    const myField = (this.$refs.textarea as any).$refs.textarea;
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      // 添加不能超过500个字符
      if (this.ruleForm.text_content.length + emj.length > 500) {
        this.$message.warning("文本内容不能超过500个字符");
        return;
      }
      this.ruleForm.text_content =
        myField.value.substring(0, startPos) +
        emj +
        myField.value.substring(endPos, myField.value.length);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + emj.length, endPos + emj.length);
    } else {
      this.ruleForm.text_content += emj;
    }
  }
  // 上传成功
  update(res: any) {
    const sortImgs = this.sortImage(res)
    console.log(sortImgs)
    this.testFlag = false;
    this.ruleForm.urls = sortImgs;
    this.ruleForm.sources = sortImgs;
    // todo:判断是否全选AI识图
    this.isAllAi = sortImgs.every((item: any) => item.is_ai === 1);
  }
  // todo:增加排序字段sort
  sortImage(arr: any) {
    for (let i = 0; i < arr.length; i++) {
      arr[i].sort = i
    }
    return arr
  }
  selectNode(node: any) {
    this.ruleForm.name = node.group_name;
    this.ruleForm.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 获取素材详情
  sourceDetail() {
    Api.sourceDetail({
      source_id: this.editRowId,
      type: 7
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data[0];
        let categoryIds = [];
        if (data.category_ids.trim()) {
          categoryIds = data.category_ids
            .split(",")
            .map(item => parseInt(item));
        }
        this.ruleForm.category_ids = categoryIds;
        if (data.start_time) {
          this.ruleForm.start_time = data.start_time;
          this.ruleForm.end_time = data.end_time;
          this.ruleForm.datePicker = [
            data.start_time * 1000,
            data.end_time * 1000
          ];
        }
        // todo:merge_id为0是视频
        console.log(data.merge_image);
        
        if (data.merge_image.length == 0 && data.original_url) {
          console.log("视频");
          this.ruleForm.file_type = 2
          this.ruleForm.urls = [data.original_url];
          this.ruleForm.sources = [{
            id: data.source_id,
            name: data.name,
            url: data.original_url,
            video_cover_url: data.video_cover_url,
            relation_goods: (data.relation_goods || []).map(v => ({
              ...v,
              checked: true
            })),
          }];
        } else {
          console.log("图片")
          if (data.merge_image && data.merge_image.length > 0) {
            this.ruleForm.file_type = 1
            // 根据sort字段排序
            data.merge_image.sort((a, b) => a.sort - b.sort)
            this.ruleForm.urls = data.merge_image.map((item) => {
              return item.url
            })
            console.log(this.ruleForm.urls)
            this.ruleForm.sources = data.merge_image
          }
        }
        // todo:文本内容
        this.ruleForm.text_content = data.text_content;
        this.ruleForm.group_id = data.group_id || -1;
        if (!data.group_id) {
          // this.ruleForm.name = "全部分组";
          this.ruleForm.name = "未分组";
        } else if (data.group_id === -1) {
          this.ruleForm.name = "未分组";
        } else {
          this.ruleForm.name = data.group_name;
        }

        this.ruleForm.sources.forEach(item => {
          item.relation_goods = (item.relation_goods || []).map(v => ({
            ...v,
            checked: true
          }))
          item.is_ai = item.relation_goods.some(v => v.add_type === 2)
          this.$set(item, 'is_ai', item.is_ai)
        })
        
        this.isAllAi = this.ruleForm.sources.every(v => v.is_ai)


      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 上传视频限制
  async beforeUploadFn(file: any) {
    // const fileMaxSize = 1073741824; //视频最大上传1G
    const fileMaxSize = 1024 * 1024 * 10; //视频最大上传1G 2020.11.30改为10M
    const isMP4 = file.type === "video/mp4";
    if (!isMP4) {
      this.$message.error("请上传MP4格式的视频");
      return false;
    }
    // 校验视频时长
    const duration:number = await getVideoDuration(file);

    if (duration > 30) {
      this.$message.warning("视频时长不能超过30S,请重新选择");
      return false;
    } 

    if (file.size > fileMaxSize) {
      this.$message.warning("请上传10M大小以内的视频");
      return false;
    }
    return true;
  }
  //上传成功
  async onSuccessFn(res: any) {
    const { fileName, url, poster } = res.data;
    this.$set(this.ruleForm, "urls", [url]);
    this.$set(this.ruleForm, "sources", [{
      name: fileName,
      url: url,
      video_cover_url: poster,
      relation_goods: [],
    }]);
  }
  //移除上传的文件
  removeFn() {
    this.isShowBtn = true;
    this.ruleForm.sources = []
  }
  // 上传文件限制规范
  beforeUpload(file: any) {
    this.testFlag = true;
    console.log(file.type);
    
    if (file.type != "image/jpeg" && file.type != "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }
    const name = file.name
      .substring(file.name.lastIndexOf(".") + 1, file.name.length)
      .toLocaleLowerCase();
    console.log(file, "测试上传文件名", name);
    if (name === "jpeg") {
      this.$message.warning("企业微信不支持JPEG的图片");
      return false;
    }
    if (name === "jfif") {
      this.$message.warning("企业微信不支持JFIF格式的图片");
      return false;
    }
    if (file.size > 10 * 1024 * 1024) {
      this.$message.error("请上传10MB大小以内的图片");
      return false;
    }
  }

  // 点击完成
  submitFn() {
    console.log("this.testFlag", this.testFlag);
    if (this.testFlag) {
      this.$message.warning("素材正在上传中");
      return;
    }
    // todo:判断文本内容或者sources是否为空
    if (!this.ruleForm.text_content && this.ruleForm.sources.length === 0) {
      this.$message.warning("文本内容或者素材不能为空");
      return;
    }
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        console.log("this.dialogType", this.dialogType);
        if (this.dialogType === 1) {
          this.sourceCreate();
        } else {
          this.sourceUpdate();
        }
      }
    });
  }
  // 关闭弹框
  closeDialog() {
    // 返回上一页
    // this.$router.push({
    //   name: "myapp_enterprise_material-library_material-list",
    //   query: {
    //     type: '7'
    //   },
    // });
    this.$router.back();
  }
  // 编辑接口
  sourceUpdate() {
    const params = {
      source_type: 7,
      id: this.editRowId,
      text_content: this.ruleForm.text_content,
      group_id: this.ruleForm.group_id || 0,
      sources: JSON.parse(JSON.stringify(this.ruleForm.sources)),
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time,
      category_ids: this.ruleForm.category_ids,
      show_type: this.ruleForm.file_type === 1 ? 2 : 1
    }
    // todo:如果sources为空
    if (params.sources.length === 0) {
      params.show_type = 0
    }
    params.sources.forEach(item => {
      item.relation_goods = (item.relation_goods || []).filter(v => v.checked)
    })
    Api.friendCircleUpdate(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("编辑成功");
        this.closeDialog();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 图片上传接口
  sourceCreate() {
    const params = {
      source_type: 7,
      text_content: this.ruleForm.text_content,
      file_type: this.ruleForm.file_type,
      group_id: this.ruleForm.group_id || 0,
      sources: JSON.parse(JSON.stringify(this.ruleForm.sources)),
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time,
      category_ids: this.ruleForm.category_ids,
      show_type: this.ruleForm.file_type === 1 ? 2 : 1
    };
    // todo:如果sources为空
    if (params.sources.length === 0) {
      params.show_type = 0
    }

    params.sources.forEach(item => {
      item.relation_goods = (item.relation_goods || []).filter(v => v.checked)
    })

    Api.friendCircleCreate(params).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("添加成功");
        this.closeDialog();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 分组接口
  grouplist() {
    Api.sourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const data = res.data;
        this.groupList = data;
        const obj1 = { group_name: "未分组", id: -1 };
        (this.groupList as any).unshift(obj1);
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 打开添加素材标签弹窗
  openAddCate() {
    this.addVisible = true;
  }
  // 关闭取消弹窗
  closeVisible() {
    this.addVisible = false;
  }
  // 新增素材标签成功
  addCateSubmit() {
    this.getCategory();
  }

  // 获取数据
  async getCategory() {
    const param: any = {
      page: 1,
      size: 999
    };
    const res: any = await Api.categoryList(param);
    if (res.code == 200) {
      this.category_options = res.data.list || [];
    } else {
      this.$message.error(res.message);
    }
  }
  // 有效期选择
  datePickerChange(val) {
    if (val) {
      // 因为 parseInt 的参数类型为 string ,所以处理后的时间戳只能先转成 string 然后再取整
      this.ruleForm.start_time = parseInt((val[0] / 1000).toString());
      this.ruleForm.end_time = parseInt((val[1] / 1000).toString());
    } else {
      this.ruleForm.start_time = 0;
      this.ruleForm.end_time = 0;
    }
  }
  /**
   * 全选AI识图商品
   * @param e 
   */
  handleAllAiChange(e) {
    this.ruleForm.sources.forEach(async (item: any) => {
      if(e){
        if(item.is_ai) {
          return
        }
        this.$set(item, 'is_ai', true)
        this.$set(item, 'is_ai_loading', true)
        const res: any = await this.getTalkieImageSearch(item.url)
        console.log(res);
        item.relation_goods = item.relation_goods || []
        const upcList = (res || []).map(v => v.upc)

        item.relation_goods = [
          ...(res || []),
          ...item.relation_goods.filter(v => v.add_type === 1 && !upcList.includes(v.upc)),
        ]
        item.is_ai_loading = false
      } else {
        this.$set(item, 'is_ai', false)
        item.relation_goods = item.relation_goods.filter(v => v.add_type === 1)
      }
    })
  }
  /**
   * 获取AI识图商品
   * @param url 
   * @returns 
   */
  async getTalkieImageSearch(url: string) {
    const res:any = await Api.getTalkieImageSearchApi({
      image_url: url,
      num: 2
    })
    if(res.code == 200){
      const list = res.data || []
      if(list.length === 0){
        return []
      }
      const result:any = []
      // 如果list数据不够5个，需要重复填充5个
      let index = 0
      if (list.length <= 2) {
        list.forEach(item => {
          result.push(...item.picList.map(goods => ({
            upc: goods.upc,
            sku: goods.sku,
            color: goods.color,
            url: goods.url || '',
            name: goods.name || '',
            add_type: 2, // 1.手动添加 2.AI识别
            score: String(goods.score) || '', // 置信度 （AI识别时有）
            checked: false
          })))
        })
      } else {
        while(list.length < 5){
          const item = JSON.parse(JSON.stringify(list[index]))
          // 去除第一个元素
          item.picList.shift()
          list.push(item)
          index++
        }
        list.forEach((item: any) => {
          const goods = item.picList[0]
          result.push({
            upc: goods.upc,
            sku: goods.sku,
            color: goods.color,
            url: goods.url || '',
            name: goods.name || '',
            add_type: 2, // 1.手动添加 2.AI识别
            score: String(goods.score) || '', // 置信度 （AI识别时有）
            checked: false
          })
        })
      }
      result.sort((a, b) => Number(b.score) - Number(a.score))
      return result
    } else {
      return []
    }
    
  }
  /**
   * 更新AI识图商品1
   * @param data 
   */
  async updateAIImage(data: any) {
    const { index, is_ai } = data
    this.ruleForm.sources[index].is_ai = is_ai
    this.isAllAi = this.ruleForm.sources.every(v => v.is_ai)
    if(is_ai){
      this.$set(this.ruleForm.sources[index], 'is_ai_loading', true)
      const res: any = await this.getTalkieImageSearch(this.ruleForm.sources[index].url)
      
      this.ruleForm.sources[index].relation_goods = this.ruleForm.sources[index].relation_goods || []
      const upcList = (res || []).map(v => v.upc)
      this.ruleForm.sources[index].relation_goods = [
        ...(res || []),
        ...this.ruleForm.sources[index].relation_goods.filter(v => v.add_type === 1 && !upcList.includes(v.upc)),
      ]
      this.ruleForm.sources[index].is_ai_loading = false
    } else {
      this.ruleForm.sources[index].relation_goods = this.ruleForm.sources[index].relation_goods.filter(v => v.add_type === 1)
    }
  }
  /**
   * 删除AI识图商品
   * @param item 
   * @param index 
   */
  handleDeleteGoods(item: any, index: number) {
    item.relation_goods.splice(index, 1)
    // if(item.relation_goods.length === 0){
    //   item.is_ai = false
    // }
  }
  /**
   * 选择商品
   * @param item 
   * @param goods 
   */
  handleSelectGoodsConfirm(item: any, goods: any) {
    console.log(item);
    console.log(goods);
    this.$set(item, "relation_goods", item.relation_goods || [])
    const upcList = item.relation_goods.filter(v => v.add_type === 2).map(v => v.upc)
    const newGoods = goods.filter(v => upcList.includes(v.upc))
    console.log('重复商品', newGoods);
    if(newGoods.length > 0){
      const skuList = newGoods.map(v => v.sku+'_'+v.color)
      this.$message.warning(`商品已存在AI识图结果，将自动过滤以下选中商品：${skuList.join('、')}`)
    }
    
    item.relation_goods = [
      ...item.relation_goods.filter(v => v.add_type === 2),
      ...goods.filter(v => !upcList.includes(v.upc))
    ]
    
    item.is_ai = item.relation_goods.filter(v => v.add_type === 2).length > 0
   }

  mounted() {
    //素材分组列表
    this.grouplist();
    this.getCategory();
    const editRowParams: any = this.$route.query.editRowParams;
    const decodeParams = decodeURIComponent(editRowParams)
    const editRowParamsObj = JSON.parse(decodeParams)
    console.log(editRowParamsObj)
    const { dialogType, editRowId, groupId, name } = editRowParamsObj
    this.dialogType = Number(dialogType);
    if (this.dialogType === 2) {
      if (editRowId) {
        this.editRowId = editRowId,
          this.sourceDetail();
      }
    } else {
      this.ruleForm.group_id = groupId || -1;
      this.ruleForm.name = groupId ? name : "未分组";
      this.ruleForm.file_type = 1
      console.log("this.ruleForm.group_id", this.ruleForm.group_id);
    }
  }
}
</script>
<style lang="scss" scoped>
.tip-text {
  font-size: 12px;
  color: #909399;
  line-height: 36px;
  // margin: 4px 0;
  margin-top: 2px;
}
.material {
  width: 600px;

  &-button {
    position: relative;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    margin-top: 10px;

    &-box {
      position: absolute;
      left: 100px;
      top: 0px;
      display: flex;
      align-items: center;
      gap: 10px;
    }

    &-choose {
      position: absolute;
      left: 100px;
      top: 0px;
    }
  }

  &-tip {
    font-size: 12px;
    color: #999999;

    ul,
    li {
      list-style: none;
      margin: 0;
      padding: 0;
    }
  }
}

.right-wrap {
  position: relative;
  margin-left: 120px;
  margin-top: 80px;

  .left-label {
    position: absolute;
    left: -70px;
    top: -50px;
    color: #606266;
  }
}

.risk-content-box {
  width: 500px;
  border: solid 1px #dcdfe6;
  border-radius: 4px;

  .select-box-insert {
    height: 60px;
    margin: 0 15px;
    display: flex;
    align-items: center;

    .add-emoji {
      display: inline-flex;
      justify-content: center;
      margin-right: 10px;

      /deep/span {
        display: inline-flex;
        align-items: center;
      }
    }

    .insert-btn {
      border: none;
      border-radius: 25px;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      width: 87px;
      height: 28px;
      background-color: #f5f6fa;

      /deep/span {
        display: inline-flex;
        align-items: center;
      }

      .btn-img {
        width: 14px;
        height: 14px;
        margin-right: 5px;
      }
    }
  }

  .textarea-box {
    width: calc(100% - 30px) !important;
    margin-left: 15px;
    border-top: solid 1px #ebeef5;
    border-bottom: solid 1px #ebeef5;

    /deep/textarea {
      border: none;
      resize: none;
      height: 200px !important;
      padding: 15px 0;
      font-family: PingFangSC-Medium, PingFang SC;
      font-size: 14px;
      line-height: 20px;
      color: #606266;
      margin-bottom: 40px;
    }

    /deep/.el-input__count {
      right: 0;
    }
  }
}

.select-video-box {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 4px;
  margin-top: 12px;

  video {
    border-radius: 4px;
  }

  .video-play-icon {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    top: 0;
    margin: auto;
    width: 25px;
    height: 25px;
    z-index: 99;
    cursor: pointer;
  }

  .delete-video-box {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    cursor: pointer;
    right: -10px;
    top: -10px;
    z-index: 2;

    img {
      width: 14px;
      height: 14px;
      display: block;
    }
  }
}

</style>
