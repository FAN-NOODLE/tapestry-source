<template>
  <div class="add-edit-group">
    <layout-form>
      <template slot="content">
        <!-------------------新增链接模块 -------------------------->
        <cy-form
          ref="cform"
          :rules="formRule"
          :model="formData"
          label-width="140px"
          class="add-link"
        >
          <!-- 素材分组 -->
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

          <!-- 链接标题 -->
          <cy-form-item label="链接标题" prop="title" required>
            <cy-input
              ref="inp"
              placeholder="请输入链接标题"
              v-model.trim="formData.title"
              size="medium"
              maxlength="30"
              show-word-limit
              width-size="medium"
            ></cy-input>
          </cy-form-item>

          <!-- 链接描述 -->
          <cy-form-item label="链接描述" prop="description" required>
            <cy-input
              ref="inp"
              type="textarea"
              placeholder="请输入链接描述"
              v-model.trim="formData.description"
              size="medium"
              maxlength="1000"
              show-word-limit
              width-size="medium"
            ></cy-input>
          </cy-form-item>

          <!-- 有效期 -->
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

          <!-- 封面图 -->
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

          <!-- 链接内容类型 -->
          <cy-form-item label="链接内容" prop="link_type">
            <el-radio
              v-model="formData.link_type"
              :label="1"
              :disabled="isEdit"
              @change="radioChange"
            >URL</el-radio>
            <el-radio
              :disabled="isEdit"
              v-model="formData.link_type"
              :label="2"
              @change="radioChange"
            >自定义</el-radio>
          </cy-form-item>

          <!-- URL输入 -->
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

          <!-- 关联商品 -->
          <el-form-item
            label="AI识图关联商品"
            aria-required=""
            prop="goods"
            max-width="100%"
            v-if="ep.isCoach"
          >
            <span class="tip" v-if="(!formData.relation_goods || formData.relation_goods.length == 0)" >勾选识图后可自动加载</span>

            <relation-goods 
              v-if="formData.cover_url" 
              :item="formData" 
              @confirm="handleSelectGoodsConfirm(formData, $event)" 
              @delete="handleDeleteGoods(formData, $event)"
            ></relation-goods>
          </el-form-item>
        </cy-form>
      </template>

      <!-- 底部按钮 -->
      <template slot="footer">
        <el-button @click="handleBack" size="small">取 消</el-button>
        <el-button type="primary" @click="submitFn" size="small">保 存</el-button>
      </template>
    </layout-form>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import Api from "@/api/enterprise/material-library";
// 引入 wangEditor富文本编辑器
import wangEditor from "wangeditor";
import Times from "../components/datePicker";
import SelectGoods from "../components/select-goods.vue";
import RelationGoods from "../components/relation-goods.vue";
import PerModule from "@/store/modules/permission";

/**
 * 编辑数据接口类型
 */
interface EditData {
  id: number;
  group_id: number;
}

/**
 * 创建链接素材组件
 * @description 用于创建或编辑链接素材的表单页面
 */
@Component({
  components: {
    UploadImage,
    SelectGoods,
    RelationGoods
  }
})
export default class CreateLink extends Vue {

  
  // ============================================================
  // 响应式数据
  // ============================================================

  /**
   * 素材ID（编辑时传入）
   */
  private sourceId?: number;

  /**
   * 素材分组列表（树形结构）
   */
  private groupList: any[] = [];

  /**
   * 是否为编辑模式
   */
  private isEdit = false;
  /**
   * 当前品牌标识
   */
  get ep () {
    return this.$store.getters['enterprise']
  }

  /**
   * 表单数据
   */
  private formData: any = {
    title: "",                     // 素材标题
    description: "",               // 素材描述
    group_id: 0,                   // 分组ID
    link_type: 2,                  // 链接类型 (1 URL 2 自定义)
    cover_url: "",                 // 封面图URL
    link_url: "",                  // 链接URL
    author: "",                    // 作者
    name: "",                      // 分组名称（显示用）
    content: "",                   // 正文内容
    datePicker: Times.times,       // 有效期选择器值
    start_time: Times.start,       // 有效期开始时间
    end_time: Times.end,           // 有效期结束时间
    relation_goods: []             // 关联商品
  };

  /**
   * 表单验证规则
   */
  private formRule: any = {
    link_url: [
      {
        validator: this.validateUrl,
        trigger: "blur"
      }
    ]
  };

  /**
   * 富文本编辑器实例
   */
  private editorSms: any = null;

  /**
   * 编辑器菜单配置
   */
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
    "list",
    "justify",
    "quote",
    "emoticon",
    "table",
    "splitLine",
    "undo",
    "redo"
  ];

  // ============================================================
  // 表单验证方法
  // ============================================================

  /**
   * URL格式验证
   */
  private validateUrl(rule: any, value: any, callback: any) {
    const match = /^(https?:\/\/)(([a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,})(:\d+)?(\/[^\s?#]*)?(\?[^\s#]*)?(#[^\s]*)?$/;
    const flag = match.test(value);
    flag ? callback() : callback(new Error("请输入正确的URL"));
  }

  // ============================================================
  // 事件处理方法
  // ============================================================

  /**
   * 树形选择节点点击事件
   * @param node 选中的节点数据
   */
  private selectNode(node: any) {
    this.formData.name = node.group_name;
    this.formData.group_id = node.id;
    (this.$refs.selectBlur as any).blur();
  }

  /**
   * 获取素材详情（编辑时调用）
   */
  private async getSourceDetail() {
    try {
      const res:any = await Api.sourceDetail({
        source_id: this.sourceId,
        type: 3
      });

      if (res.code === 200) {
        const data = res.data[0];
        this.formData.title = data.title;
        this.formData.description = data.description;
        this.formData.group_id = data.group_id || -1;
        this.formData.link_type = data.link_type;
        this.formData.cover_url = data.cover_url;
        this.formData.link_url = data.link_url;
        this.formData.author = data.author;
        this.formData.content = data.content;
        this.formData.relation_goods = (data.relation_goods || []).map(item => ({
          ...item,
          checked: true
        }));

        // 设置有效期
        if (data.start_time) {
          this.formData.start_time = data.start_time;
          this.formData.end_time = data.end_time;
          this.formData.datePicker = [
            data.start_time * 1000,
            data.end_time * 1000
          ];
        }

        // 设置分组名称
        if (!data.group_id || data.group_id === -1) {
          this.formData.name = "未分组";
        } else {
          this.formData.name = data.group_name;
        }

        // 如果是自定义类型，初始化编辑器
        if (this.formData.link_type === 2) {
          this.initEditor();
        }
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error("获取素材详情失败");
    }
  }

  /**
   * 获取分组列表
   */
  private async getGroupList() {
    try {
      const res: any = await Api.sourceGroupListApi({});

      if (res.code === 200) {
        const arr = [{ group_name: "未分组", id: -1 }];
        this.groupList = [...arr, ...res.data];
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error("获取分组列表失败");
    }
  }

  /**
   * 关闭弹窗
   */
  private handleBack() {
    const { from } = this.$route.query;
    if(from){
      this.$router.replace({ 
        name: 'myapp_enterprise_material-library_material-list',
        query: { 
          type: '3' 
        } 
      });
    } else {
      this.$router.back();
    }
  }

  /**
   * 提交保存
   */
  private submitFn() {
    if (this.isEdit) {
      this.updateSource();
    } else {
      this.addConfirm();
    }
  }
  

  /**
   * 新增确认
   */
  private addConfirm() {
    // 获取编辑器内容
    if (this.editorSms) {
      this.$set(this.formData, "content", this.editorSms.txt.html());
    }

    const data = JSON.parse(JSON.stringify(this.formData));
    data.group_id = data.group_id || 0;
    data.category_ids = [];

    // 根据链接类型处理数据
    if (this.formData.link_type === 1) {
      data.author = "";
      data.content = "";
    } else {
      data.link_url = "";
      data.content = this.editorSms ? this.editorSms.txt.html() : "";
    }

    data.relation_goods = data.relation_goods.filter(item => item.checked);

    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        Api.linkSaveApi(data).then((res: any) => {
          if (res.code === 200) {
            this.handleBack();
            this.$message.success("保存成功");
          } else {
            this.$message.error(res.message);
          }
        });
      }
    });
  }

  /**
   * 更新素材（编辑操作）
   */
  private updateSource() {
    const params: any = {
      source_type: 3,
      id: this.sourceId,
      category_ids: [],
      ...this.formData,
      relation_goods: this.formData.relation_goods.filter(item => item.checked)
    }
    Api.sourceUpdate(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("编辑成功");
        this.handleBack();
      } else {
        this.$message.error(res.message);
      }
    });
  }

  /**
   * 图片上传成功回调
   */
  private picSuccess(res: any) {
    if(!res){
      this.formData.relation_goods = []
    }
    
    this.$set(this.formData, "cover_url", res);
    (this.$refs["cform"] as any).validateField("cover_url");
    if (this.ep.isCoach) {
      this.updateAIImage();
    }
  }

  /**
   * 上传文件前置校验
   */
  private beforeUpload(file: any) {
    if (file.type !== "image/jpeg" && file.type !== "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }

    const ext = file.name
      .substring(file.name.lastIndexOf(".") + 1)
      .toLowerCase();

    if (ext === "jpeg") {
      this.$message.warning("企业微信不支持JPEG格式的图片");
      return false;
    }

    if (file.size > 2 * 1024 * 1024) {
      this.$message.error("请上传2MB大小以内的图片");
      return false;
    }

    return true;
  }

  /**
   * 链接类型切换处理
   */
  private radioChange() {
    if (this.formData.link_type === 2) {
      this.initEditor();
    }

    this.$nextTick(() => {
      (this.$refs["cform"] as any).clearValidate(["link_url"]);
      (this.$refs["cform"] as any).clearValidate(["content"]);
      (this.$refs["cform"] as any).clearValidate(["author"]);
    });
  }

  /**
   * 初始化富文本编辑器
   */
  private initEditor() {
    // if (this.editorSms) return;

    this.editorSms = new wangEditor("#editContent");
    this.editorSms.config.zIndex = 500;
    this.editorSms.config.menus = this.editorMenus;
    this.editorSms.create();

    if (this.formData.content) {
      this.editorSms.txt.html(this.formData.content);
    }
  }

  /**
   * 有效期选择变化回调
   */
  private datePickerChange(val: number[]) {
    if (val && val.length === 2) {
      this.formData.start_time = parseInt((val[0] / 1000).toString(), 10);
      this.formData.end_time = parseInt((val[1] / 1000).toString(), 10);
    } else {
      this.formData.start_time = 0;
      this.formData.end_time = 0;
    }
  }

  /**
   * 删除关联商品
   */
  private handleDeleteGoods(item: any, index: number) {
    item.relation_goods.splice(index, 1);
    if (item.relation_goods.length === 0) {
      item.is_ai = false;
    }
  }

  /**
   * 选择商品确认
   */
  private handleSelectGoodsConfirm(item: any, goods: any) {
    const upcList = item.relation_goods.filter(v => v.add_type === 2).map(v => v.upc);
    const newGoods = goods.filter(v => upcList.includes(v.upc))
    console.log('重复商品', newGoods);
    if(newGoods.length > 0){
      const skuList = newGoods.map(v => v.sku+'_'+v.color)
      this.$message.warning(`商品已存在AI识图结果，将自动过滤以下选中商品：${skuList.join('、')}`)
    }

    item.relation_goods = [
      ...item.relation_goods.filter((v: any) => v.add_type === 2),
      ...goods.filter((v: any) => !upcList.includes(v.upc))
    ];

    item.is_ai = item.relation_goods.filter((v: any) => v.add_type === 2).length > 0;
  }

  /**
   * AI识图接口调用
   */
  private async getTalkieImageSearch(url: string): Promise<any[]> {
    const res:any = await Api.getTalkieImageSearchApi({
      image_url: url,
      num: 2
    });

    if (res.code === 200) {
      let list = res.data || []
      if(list.length === 0){
        return []
      }
      const result: any[] = [];

      // 如果数据不够5个，重复填充
      let index = 0;
      if (list.length <= 2) {
        list.forEach(item => {
          result.push(...item.picList.map(goods => ({
            upc: goods.upc,
            sku: goods.sku,
            color: goods.color,
            url: goods.url || '',
            name: goods.name || '',
            add_type: 2, // 1.手动添加 2.AI识别
            score: String(goods.score) || '', // 置信度 （AI识别时有）1
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
      return result;
    } else {
      return [];
    }
  }

  /**
   * 更新AI识图商品
   */
  private async updateAIImage() {
    this.$set(this.formData, 'is_ai_loading', true);

    try {
      const res = await this.getTalkieImageSearch(this.formData.cover_url);
      const currentGoods = this.formData.relation_goods || [];
      const upcList = (res || []).map((v: any) => v.upc);

      this.formData.relation_goods = [
        ...(res || []),
        ...currentGoods.filter((v: any) => v.add_type === 1 && !upcList.includes(v.upc))
      ];
    } catch (error) {
      console.error('AI识图失败:', error);
    } finally {
      this.$set(this.formData, 'is_ai_loading', false);
    }
  }

  // ============================================================
  // 生命周期钩子
  // ============================================================

  /**
   * 组件挂载完成
   */
  private async mounted() {
    const { id, groupId } = this.$route.query;
    // 判断是新增还是编辑模式
    if (!id) {
      this.isEdit = false;
      this.formData.group_id = groupId || -1;
      this.formData.link_type = 1;
    } else {
      this.isEdit = true;
      this.sourceId = Number(id);
      await this.getSourceDetail();
    }

    // 获取分组列表
    await this.getGroupList();

    
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

.add-edit-group {
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
    // max-height: 465px;
    overflow-y: auto;
  }

  .upload-tip {
    color: #909399;
    font-size: 12px;
    line-height: 18px;
    margin-top: 5px;
  }
  .tip {
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
