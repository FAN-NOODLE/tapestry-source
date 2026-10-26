/* eslint-disable no-useless-escape */
<template>
  <div class="select-materials-dialog">
    <el-dialog
      :title="title"
      :value="value"
      :visible.sync="selectDialogVisible"
      width="1000px"
      @close="cancel"
    >
      <!-------------------新增链接模块 -------------------------->
      <cy-form
        ref="cform"
        :rules="rules"
        :model="formData"
        label-width="90px"
        class="add-link"
      >
        <!-- <cy-form-item label="素材分组" prop="group_id">
          <cy-select
            v-model="formData.group_id"
            size="medium"
            placeholder="请选择素材分组"
            search-input
            :optionsList="groupList"
            :options="{ label: 'group_name', value: 'id' }"
          >
          </cy-select>
        </cy-form-item> -->
        <cy-form-item label="素材分组">
          <el-cascader
            v-model="formData.group_id"
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
        <cy-form-item label="链接标题" prop="title" required>
          <cy-input
            ref="inp"
            placeholder="请输入链接标题"
            v-model.trim="formData.title"
            size="medium"
            maxlength="30"
            show-word-limit
            width-size="medium"
          ></cy-input
        ></cy-form-item>
        <cy-form-item label="链接描述" prop="description" required>
          <cy-input
            ref="inp"
            placeholder="请输入链接描述"
            v-model.trim="formData.description"
            size="medium"
            maxlength="30"
            show-word-limit
            width-size="medium"
          ></cy-input>
        </cy-form-item>
        <cy-form-item
          label="有效期"
          prop="datePicker"
          required
          trigger="change"
        >
          <el-date-picker
            v-model="formData.datePicker"
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
        <cy-form-item label="封面图" required prop="cover_url" trigger="change">
          <upload-image
            class="upload-template"
            accept="image/jpeg,image/png"
            :urls="formData.cover_url"
            @change="picSuccess"
            :before-upload="beforeUpload"
          ></upload-image>
          <div class="upload-tip">支持JPG、PNG格式，图片＜2M</div>
        </cy-form-item>
        <cy-form-item label="链接内容" prop="link_type">
          <el-radio
            v-model="formData.link_type"
            :label="1"
            :disabled="isEdit"
            @change="radioChange"
            >URL
          </el-radio>
          <el-radio
            :disabled="isEdit"
            v-model="formData.link_type"
            :label="2"
            @change="radioChange"
            >自定义</el-radio
          >
        </cy-form-item>
        <cy-form-item
          label="URL"
          required
          v-if="formData.link_type == 1"
          prop="link_url"
        >
          <cy-input
            ref="inp"
            :disabled="isEdit"
            placeholder="请输入URL"
            v-model.trim="formData.link_url"
            size="medium"
          ></cy-input>
        </cy-form-item>
        <cy-form-item
          label="作者"
          required
          v-if="formData.link_type == 2"
          prop="author"
        >
          <cy-input
            ref="inp"
            placeholder="请输入作者名称"
            v-model.trim="formData.author"
            size="medium"
          ></cy-input>
        </cy-form-item>
        <cy-form-item
          label="正文"
          required
          v-if="formData.link_type == 2"
          prop="content"
        >
          <div id="editContent"></div>
        </cy-form-item>
        <el-form-item
          label="关联商品"
          aria-required=""
          prop="goods"
          max-width="100%"
        >
          <select-goods v-if="formData.cover_url && (!formData.relation_goods || formData.relation_goods.length == 0)" :value="formData.relation_goods" @confirm="handleSelectGoodsConfirm"></select-goods>
          <relation-goods v-if="formData.cover_url && formData.relation_goods && formData.relation_goods.length > 0" :item="formData" @confirm="handleSelectGoodsConfirm(formData, $event)" @delete="handleDeleteGoods(formData, $event)"></relation-goods>
          
        </el-form-item>
      </cy-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="cancel" size="small">取 消</el-button>
        <el-button type="primary" @click="submitFn" size="small"
          >保 存</el-button
        >
      </span>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import Api from "@/api/enterprise/material-library";
// 引入 wangEditor富文本编辑器
import wangEditor from "wangeditor";
import Times from "./datePicker";
import SelectGoods from "./select-goods.vue";
import RelationGoods from "./relation-goods.vue";
interface T {
  id: number;
  group_id: number;
}
@Component({
  components: {
    UploadImage,
    SelectGoods,
    RelationGoods
  }
})
export default class extends Vue {
  @Prop() value?: T;
  @Prop() title?: string;
  @Prop() groupName?: string;
  @Prop() groupId!: number;
  //------------------------------------------------------data-----------------------------------------------------------/
  private selectDialogVisible = true;
  private groupData: any = []; //左侧分组
  private linkList: any = [];
  private isActive = 0;
  private isChecked = 0;
  private defaultShow = true; //默认展示
  private selectedMaterial = {};
  // ---------新增链接data-----------------
  private groupList: any = [];
  private isEdit = false;

  private formData: any = {
    title: "", //素材标题
    description: "", //素材描述
    group_id: 0,
    link_type: 2, //(1 URL 2 自定义)
    cover_url: "", //封面图URL
    link_url: "", //链接URL
    author: "", //作者
    name: "",
    content: "", //正文内容
    datePicker: Times.times,
    start_time: Times.start, //有效期开始时间
    end_time: Times.end, //有效期结束时间
    relation_goods: [] //关联商品
  };
  private rules = {
    link_url: [
      {
        validator: this.validateself,
        trigger: "blur"
      }
    ]
  };
  private validateself(rule: any, value: any, callback: any) {
    const match = /^(https?:\/\/)(([a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,})(:\d+)?(\/[^\s?#]*)?(\?[^\s#]*)?(#[^\s]*)?$/;
    const flag = match.test(value);
    flag ? callback() : callback(new Error("请输入正确的URL"));
  }
  // no-useless-scape
  // 富文本编辑器
  private editorSms: any = null; //编辑器
  private editorMenus = [
    "head",
    "bold",
    "fontSize",
    "fontName",
    "italic",
    "underline",
    "strikeThrough",
    "indent",
    "lineHeight",
    "foreColor",
    "backColor",
    // "link",
    "list",
    "justify",
    "quote",
    "emoticon",
    // "image",
    // "video",
    "table",
    // "code",
    "splitLine",
    "undo",
    "redo"
  ]; //编辑器菜单
  //----------------------------------------------------methods-----------------------------------------------------------/
  selectNode(node: any) {
    this.formData.name = node.group_name;
    this.formData.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }
  // 获取素材详情
  private sourceDetail() {
    Api.sourceDetail({
      source_id: this.value.id,
      type: 3
    }).then((res: any) => {
      if (res.code == 200) {
        const data = res.data[0];
        this.formData.title = data.title;
        this.formData.description = data.description;
        this.formData.group_id = data.group_id || -1;
        this.formData.link_type = data.link_type;
        this.formData.cover_url = data.cover_url;
        this.formData.link_url = data.link_url;
        this.formData.author = data.author;
        this.formData.content = data.content;
        this.formData.relation_goods = data.relation_goods || [];
        if (data.start_time) {
          this.formData.start_time = data.start_time;
          this.formData.end_time = data.end_time;
          this.formData.datePicker = [
            data.start_time * 1000,
            data.end_time * 1000
          ];
        }

        if (!data.group_id) {
          this.formData.name = "未分组";
        } else if (data.group_id === -1) {
          this.formData.name = "未分组";
        } else {
          this.formData.name = data.group_name;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取分组名
  inquireName(row: any) {
    Api.inquireGroupName({ group_id: row.group_id }).then((res: any) => {
      if (res.code == 200) {
        if (this.formData.link_type === 1) {
          this.formData.author = "";
          this.formData.content = "";
        } else {
          setTimeout(() => {
            this.radioChange();
          }, 0);
          this.formData.link_url = "";
        }
        this.formData = row;
        this.formData.name = res.data.group_name;
        this.formData.group_id = row.group_id;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 获取分组
  private grouplist() {
    Api.sourceGroupListApi({}).then((res: any) => {
      if (res.code == 200) {
        const arr = [{ group_name: "未分组", id: -1 }];
        this.groupList = [...arr, ...res.data];
        if (JSON.stringify(this.value) == "{}") {
          // 新增
          this.formData.link_type = 1;
        } else {
          // 编辑
          // this.inquireName(this.value);
        }
        // 编辑
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 取消
  private cancel() {
    this.$emit("closeDialog");
  }
  // 提交保存
  private submitFn() {
    if (this.isEdit) {
      this.sourceUpdate();
    } else {
      this.addConfirm();
    }
  }
  // 确认
  private addConfirm() {
    if (this.editorSms) {
      this.$set(this.formData, "content", this.editorSms.txt.html());
    }
    const data = JSON.parse(JSON.stringify(this.formData));
    data.group_id = data.group_id || 0;
    data.category_ids = [];
    if (this.formData.link_type == 1) {
      data.author = "";
      data.content = "";
    } else {
      data.link_url = "";
      data.content = this.editorSms.txt.html();
    }

    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        Api.linkSaveApi(data).then((res: any) => {
          if (res.code == 200) {
            this.$emit("closeDialog");
            this.$message.success("保存成功");
          } else {
            this.$message.error(res.message);
          }
        });
      } else {
        return;
      }
    });
  }
  // 编辑接口
  private sourceUpdate() {
    Api.sourceUpdate({
      source_type: 3,
      id: this.value.id,
      category_ids: [],
      ...this.formData
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("编辑成功");
        this.$emit("closeDialog");
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 图片上传成功
  private picSuccess(res: any) {
    this.$set(this.formData, "cover_url", res);
    (this.$refs["cform"] as any).validateField("cover_url");
    this.updateAIImage()
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
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
    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }
  }
  private radioChange() {
    if (this.formData.link_type == 2) {
      this.editorSms = new wangEditor(`#editContent`);
      this.editorSms.config.zIndex = 500;
      this.editorSms.config.menus = this.editorMenus;
      this.editorSms.create();
      if (this.formData.content) {
        this.editorSms.txt.html(this.formData.content);
      }
    }
    this.$nextTick(() => {
      (this.$refs["cform"] as any).clearValidate(["link_url"]);
      (this.$refs["cform"] as any).clearValidate(["content"]);
      (this.$refs["cform"] as any).clearValidate(["author"]);
    });
  }
  // 有效期选择
  private datePickerChange(val) {
    if (val) {
      // 因为 parseInt 的参数类型为 string ,所以处理后的时间戳只能先转成 string 然后再取整
      this.formData.start_time = parseInt((val[0] / 1000).toString());
      this.formData.end_time = parseInt((val[1] / 1000).toString());
    } else {
      this.formData.start_time = 0;
      this.formData.end_time = 0;
    }
  }

  /**
   * 删除AI识图商品
   * @param item 
   * @param index 
   */
  handleDeleteGoods(item: any, index: number) {
    item.relation_goods.splice(index, 1)
    if(item.relation_goods.length === 0){
      item.is_ai = false
    }
  }
  /**
   * 选择商品
   * @param item 
   * @param goods 
   */
  handleSelectGoodsConfirm(item: any, goods: any) {
    console.log(item);
    console.log(goods);
    const upcList = item.relation_goods.map(v => v.upc)
    
    item.relation_goods = [
      ...item.relation_goods.filter(v => v.add_type === 2),
      ...goods.filter(v => !upcList.includes(v.upc))
    ]
    item.is_ai = item.relation_goods.filter(v => v.add_type === 2).length > 0
  }

  /**
   * 获取AI识图商品
   * @param url 
   * @returns 
   */
  async getTalkieImageSearch(url: string) {
    const res:any = await Api.getTalkieImageSearchApi({
      image_url: url,
      num: 5
    })
    if(res.code == 200){
      const list = res.data || []
      if(list.length === 0){
        return []
      }
      const result:any = []
      // 如果list数据不够5个，需要重复填充5个
      let index = 0
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
          score: String(goods.score) || '' // 置信度 （AI识别时有）
        })
      })

      result.sort((a, b) => Number(b.score) - Number(a.score))
      return result
    } else {
      return []
    }
    
  }
  /**
   * 更新AI识图商品
   * @param data 
   */
  async updateAIImage() {
    this.$set(this.formData, 'is_ai_loading', true)
    const res: any = await this.getTalkieImageSearch(this.formData.cover_url)
    this.formData.relation_goods = this.formData.relation_goods || []
    const upcList = (res || []).map(v => v.upc)
    this.formData.relation_goods = [
      ...(res || []),
      ...this.formData.relation_goods.filter(v => v.add_type === 1 && !upcList.includes(v.upc)),
    ]
    this.formData.is_ai_loading = false
  }
  
  // -------------生命周期函数--------------------
  private mounted() {
    if (JSON.stringify(this.value) == "{}") {
      this.isEdit = false;
      this.formData.group_id = this.groupId || -1;
      this.formData.name = this.groupId ? this.groupName : "未分组";
    } else {
      this.isEdit = true;
      this.sourceDetail();
    }

    this.grouplist();
  }
}
</script>
<style lang="scss" scoped>
/deep/ .el-form--label-left .el-form-item__label {
  text-align: right !important;
  padding-right: 20px;
}
.goods-image {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}
.select-materials-dialog {
  // -------------------新增链接 ui-----------------
  .header {
    display: flex;
    justify-content: space-between;
    text-align: center;
    margin-bottom: 20px;
    font-size: 18px;
    font-family: PingFangSC-Semibold, PingFang SC;
    cursor: pointer;
    .chosoe-clas {
      color: #33363c;
      font-weight: 600;
    }
    .add-clas {
      color: #33363c;
      font-weight: 600;
    }
  }
  .add-link {
    max-height: 465px;
    overflow-y: auto;
  }
  .upload-tip {
    color: #909399;
    font-size: 12px;
    line-height: 18px;
    margin-top: 5px;
  }
  #editContent {
    width: 700px;
  }
}
</style>
