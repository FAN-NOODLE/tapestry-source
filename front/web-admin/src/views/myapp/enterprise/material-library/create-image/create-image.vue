<template>
  <div class="add-edit-group">
    <layout-form>
      <template slot="content">
        <!-- 表单组件 -->
        <cy-form
          class="main-box"
          ref="cform"
          :model="formData"
          label-width="140px"
          :rules="formRule"
        >
          <!-- 素材分组选择 -->
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

          <!-- 素材标签选择 -->
          <cy-form-item label="素材标签" prop="category_ids" trigger="change">
            <div class="flex">
              <el-select
                style="width:400px"
                v-model="formData.category_ids"
                multiple
                placeholder="请选择素材标签"
              >
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                >
                </el-option>
              </el-select>
              <el-button
                class="m-l-20"
                type="text"
                icon="el-icon-circle-plus"
                size="small"
                @click="openAddCate"
              ></el-button>
            </div>
          </cy-form-item>

          <!-- 有效期选择 -->
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

          <!-- 图片上传 -->
          <cy-form-item label="本地图片" required prop="sources" trigger="change">
            <el-checkbox v-if="dialogType === 1 && formData.sources.length > 0" v-model="isAllAi" @change="handleAllAiChange">全选AI识图</el-checkbox>
            <upload-image
              class="upload-template"
              :urls="formData.sources"
              accept="image/jpeg,image/png"
              img-type="batch"
              :before-upload="beforeUpload"
              :limit="formData.show_type === 1 ? 20 : 9"
              multiple
              :disabled="dialogType === 2"
              @update="update"
              @updateAI="updateAIImage"
            >
            </upload-image>
            <div class="tip-view">
              <p class="tip-text">1. 支持上传JPG,PNG格式图片</p>
              <p class="tip-text">2. 最大可上传10MB</p>
            </div>
          </cy-form-item>

          <!-- 图片展示样式（仅coach品牌显示） -->
          <cy-form-item label="图片展示样式" prop="show_type" trigger="change" v-if="ep.isCoach">
            <el-radio-group v-model="formData.show_type" :disabled="dialogType === 2">
              <el-radio :label="1">单张展示(最多20张)</el-radio>
              <el-radio :label="2">合并展示(最多9张)</el-radio>
            </el-radio-group>
          </cy-form-item>
          <cy-form-item v-if="ep.isCoach" label="AI识图关联商品" >
            <div class="tip">勾选识图后可自动加载</div>
            <template v-for="item in formData.sources">
              <relation-goods  :key="item.hash || item.id" v-if="item.is_ai || (item.relation_goods && item.relation_goods.length > 0)" :item="item" @confirm="handleSelectGoodsConfirm(item, $event)" @delete="handleDeleteGoods(item, $event)"></relation-goods>
            </template>
          </cy-form-item>
        </cy-form>
      </template>

      <!-- 底部按钮 -->
      <template slot="footer">
        <el-button
          @click="handleBack"
          size="small"
        >取 消</el-button>
        <el-button
          type="primary"
          @click="submitFn"
          size="small"
        >确 认</el-button>
      </template>
    </layout-form>

    <!-- 添加素材标签弹窗 -->
    <add-category
      v-if="addVisible"
      @closeVisible="closeVisible"
      @submit="addCateSubmit"
    ></add-category>
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import Api from "@/api/enterprise/material-library";
import UploadImage from "@/components/business/upload-image/index.vue";
import AddCategory from "../components/add-category.vue";
import Times from "../components/datePicker";
import PerModule from "@/store/modules/permission";
import SelectGoods from "../components/select-goods.vue";
import RelationGoods from "../components/relation-goods.vue";


/**
 * 创建图片素材组件
 * @description 用于创建或编辑图片素材的表单页面
 */
@Component({
  components: {
    UploadImage,
    AddCategory,
    SelectGoods,
    RelationGoods
  }
})
export default class CreateImage extends Vue {

  // ============================================================
  // 响应式数据
  // ============================================================
  /**
   * 操作类型：1-新增，2-编辑
   */
  private dialogType = 1;
  /**
   * 上传状态标记，用于防止上传过程中提交表单
   */
  private uploadFlag = false;

  /**
   * 素材分组列表（树形结构）
   */
  private groupList: any[] = [];

  /**
   * 表单数据
   */
  private formData: any = {
    group_id: 0,                           // 分组ID
    name: "",                              // 分组名称（显示用）
    category_ids: [],                      // 素材标签ID数组
    datePicker: Times.times,               // 有效期选择器值
    start_time: Times.start,               // 有效期开始时间戳
    end_time: Times.end,                   // 有效期结束时间戳
    show_type: 1,                          // 展示类型：1-单张展示，2-合并展示
    sources: []                            // 图片素材数组
  };

  /**
   * 表单验证规则
   */
  private formRule: any = {};

  /**
   * 素材标签选项列表
   */
  private categoryOptions: any[] = [];

  /**
   * 控制添加素材标签弹窗显示
   */
  private addVisible = false;
  /**
   * 是否全选AI识图
   */
  isAllAi = false;

  /**
   * 当前品牌标识
   */
  get ep () {
    return this.$store.getters['enterprise']
  }

  get disabledCheck() {
    return () => {
      console.log(this.formData.sources);
      
      return this.formData.sources.some(v => v.is_ai_loading)
    }
  }

  // ============================================================
  // 事件处理方法
  // ============================================================

  /**
   * 图片上传成功回调
   * @param res 上传返回的图片数据数组
   */
  private update(res: any) {
    // 为图片添加排序字段
    const sortedImgs = this.sortImages(res);
    this.uploadFlag = false;
    this.formData.sources = sortedImgs;
    this.isAllAi = this.formData.sources.every(v => v.is_ai)
  }

  /**
   * 为图片数组添加排序字段
   * @param arr 图片数据数组
   */
  private sortImages(arr: any[]): any[] {
    return arr.map((item, index) => ({
      ...item,
      sort: index
    }));
  }

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
  private async getSourceDetail(id: number) {
    try {
      const res: any = await Api.sourceDetail({
        source_id: id,
        type: 2
      });

      if (res.code === 200) {
        const data = res.data[0];
        
        // 处理素材标签
        let categoryIds: number[] = [];
        if (data.category_ids && data.category_ids.trim()) {
          categoryIds = data.category_ids
            .split(",")
            .map((item: string) => parseInt(item, 10));
        }
        this.formData.category_ids = categoryIds;

        // 处理有效期
        if (data.start_time) {
          this.formData.start_time = data.start_time;
          this.formData.end_time = data.end_time;
          this.formData.datePicker = [
            data.start_time * 1000,
            data.end_time * 1000
          ];
        }

        // 处理图片数据（区分合并展示和单张展示）
        if (data.merge_image && data.merge_image.length > 0 && this.ep.isCoach) {
          // 合并展示：按sort字段排序
          data.merge_image.sort((a: any, b: any) => a.sort - b.sort);
          this.formData.sources = data.merge_image.map(item => {
            return {
              ...item,
              relation_goods: (item.relation_goods || []).map(v => ({
                ...v,
                checked: true
              }))
            }
          })
        } else {
          // 单张展示
          this.formData.sources = [{
            id: data.source_id,
            name: data.name,
            url: data.url,
            is_ai: (data.relation_goods || []).filter((item: any) => item.add_type == 2).length > 0,
            relation_goods: (data.relation_goods || []).map(v => ({
              ...v,
              checked: true
            }))
          }];
        }

        // 设置展示类型
        this.formData.show_type = data.merge_id ? 2 : 1;
        this.formData.group_id = data.group_id || -1;
        this.formData.id = id;

        console.log(this.formData);

        this.formData.sources.forEach(item => {
          item.relation_goods = item.relation_goods || []
          item.is_ai = item.relation_goods.some(v => v.add_type === 2)
          this.$set(item, 'is_ai', item.is_ai)
        })
        
        this.isAllAi = this.formData.sources.every(v => v.is_ai)


        
        
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error("获取素材详情失败");
    }
  }

  /**
   * 上传文件前置校验
   * @param file 上传的文件对象
   */
  private beforeUpload(file: any) {
    this.uploadFlag = true;
    console.log(file);
    

    // 校验文件类型
    if (file.type !== "image/jpeg" && file.type !== "image/png") {
      this.$message.error("请上传JPG、PNG格式的图片");
      return false;
    }

    // 校验文件扩展名（企业微信不支持jpeg）
    const ext = file.name
      .substring(file.name.lastIndexOf(".") + 1)
      .toLowerCase();
    if (ext === "jpeg") {
      this.$message.warning("企业微信不支持JPEG格式的图片");
      return false;
    }
    if (ext === "jfif") {
      this.$message.warning("企业微信不支持JFIF格式的图片");
      return false;
    }

    // 校验文件大小（最大10MB）
    if (file.size > 10 * 1024 * 1024) {
      this.$message.error("请上传10MB大小以内的图片");
      return false;
    }

    return true;
  }

  /**
   * 表单提交处理
   */
  private async submitFn() {
    // 检查是否正在上传
    if (this.uploadFlag) {
      this.$message.warning("图片正在上传中，请稍后");
      return;
    }

    // 校验单张展示数量（最多20张）
    if (this.formData.show_type === 1 && this.formData.sources.length > 20) {
      this.$message.warning("单张展示最多上传20张图片");
      return;
    }

    // 校验合并展示数量（最多9张，最少2张）
    if (this.formData.show_type === 2 && this.formData.sources.length > 9) {
      this.$message.warning("合并展示最多上传9张图片");
      return;
    }
    if (this.formData.show_type === 2 && this.formData.sources.length < 2) {
      this.$message.warning("合并展示至少上传2张图片");
      return;
    }

    // 表单验证
    const valid = await (this.$refs["cform"] as any).validate();
    if (!valid) return;

    // 根据操作类型执行不同逻辑
    if (this.dialogType === 1) {
      await this.addImage();
    } else {
      await this.updateSource();
    }
  }


  /**
   * 更新素材（编辑操作）
   */
  private async updateSource() {
    try {
      const params:any = {
        source_type: 2,
        id: this.formData.id,
        group_id: this.formData.group_id || 0,
        sources: JSON.parse(JSON.stringify(this.formData.sources)),
        start_time: this.formData.start_time,
        end_time: this.formData.end_time,
        category_ids: this.formData.category_ids,
        show_type: this.ep.isCoach ? this.formData.show_type : 1
      }
      params.sources.forEach(item => {
        item.relation_goods = (item.relation_goods || []).filter(v => v.checked)
      })
      const res: any = await Api.sourceUpdate(params);

      if (res.code === 200) {
        this.$message.success("编辑成功");
        this.handleBack();
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error("编辑失败");
    }
  }

  /**
   * 添加图片素材（新增操作）
   */
  private async addImage() {
    try {
      const params:any = {
        group_id: this.formData.group_id || 0,
        images: JSON.parse(JSON.stringify(this.formData.sources)),
        start_time: this.formData.start_time,
        end_time: this.formData.end_time,
        category_ids: this.formData.category_ids,
        show_type: this.formData.show_type
      };
      params.images.forEach(item => {
        item.relation_goods = (item.relation_goods || []).filter(v => v.checked)
      })

      const res:any = await Api.imageAddApi(params);

      if (res.code === 200) {
        this.$message.success("添加成功");
        this.handleBack();
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error("添加失败");
    }
  }

  /**
   * 获取素材分组列表
   */
  private async getGroupList() {
    try {
      const res: any = await Api.sourceGroupListApi({});

      if (res.code === 200) {
        const data = res.data;
        this.groupList = data;
        
        // 在列表开头添加"未分组"选项
        this.groupList.unshift({ group_name: "未分组", id: -1 });
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error("获取分组列表失败");
    }
  }

  /**
   * 打开添加素材标签弹窗
   */
  private openAddCate() {
    this.addVisible = true;
  }

  /**
   * 关闭添加素材标签弹窗
   */
  private closeVisible() {
    this.addVisible = false;
  }

  /**
   * 新增素材标签成功后的回调
   */
  private addCateSubmit() {
    this.getCategoryList();
    this.closeVisible();
  }

  /**
   * 处理返回操作
   */
  private handleBack() {
    const { from } = this.$route.query;
    if(from){
      this.$router.replace({ 
        name: 'myapp_enterprise_material-library_material-list',
        query: { 
          type: '2' 
        } 
      });
    } else {
      this.$router.back();
    }
  }

  /**
   * 获取素材标签列表
   */
  private async getCategoryList() {
    try {
      const params = {
        page: 1,
        size: 999
      };

      const res: any = await Api.categoryList(params);

      if (res.code === 200) {
        this.categoryOptions = res.data.list || [];
      } else {
        this.$message.error(res.message);
      }
    } catch (error) {
      this.$message.error("获取标签列表失败");
    }
  }

  /**
   * 有效期选择变化回调
   * @param val 选择的日期范围数组
   */
  private datePickerChange(val: number[]) {
    if (val && val.length === 2) {
      // 转换为秒级时间戳
      this.formData.start_time = parseInt((val[0] / 1000).toString(), 10);
      this.formData.end_time = parseInt((val[1] / 1000).toString(), 10);
    } else {
      this.formData.start_time = 0;
      this.formData.end_time = 0;
    }
  }

  /**
   * 全选AI识图商品
   * @param e 
   */
  handleAllAiChange(e) {
    this.formData.sources.forEach(async (item: any) => {
      if(e){
        // 如果已经选中，直接返回
        if(item.is_ai){
          return
        }
        // 如果没有选中，需要获取AI识图商品
        this.$set(item, 'is_ai', e)
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
        this.$set(item, 'is_ai', e)
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
   * 更新AI识图商品
   * @param data 
   */
  async updateAIImage(data: any) {
    const { index, is_ai } = data
    
    this.$set(this.formData.sources[index], 'is_ai', is_ai)
    this.isAllAi = this.formData.sources.every(v => v.is_ai)
    if(is_ai){
      this.$set(this.formData.sources[index], 'is_ai_loading', true)
      const res: any = await this.getTalkieImageSearch(this.formData.sources[index].url)
      this.formData.sources[index].relation_goods = this.formData.sources[index].relation_goods || []
      const upcList = (res || []).map(v => v.upc)
      this.formData.sources[index].relation_goods = [
        ...(res || []),
        ...this.formData.sources[index].relation_goods.filter(v => v.add_type === 1 && !upcList.includes(v.upc)),
      ]
      this.formData.sources[index].is_ai_loading = false
    } else {
      this.formData.sources[index].relation_goods = this.formData.sources[index].relation_goods.filter(v => v.add_type === 1)
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

  // ============================================================
  // 生命周期钩子
  // ============================================================

  /**
   * 组件挂载完成
   */
  private async mounted() {
    // 获取分组列表和标签列表
    await Promise.all([
      this.getGroupList(),
      this.getCategoryList()
    ]);
    const { id, groupId } = this.$route.query;
    // 如果是编辑模式，获取素材详情
    if (id) {
      this.dialogType = 2;
      await this.getSourceDetail(Number(id));
    } else {
      this.dialogType = 1;
      // 新增模式：设置默认分组
      this.formData.group_id = Number(groupId) || -1;
    }
  }
}
</script>

<style lang="scss" scoped>
.flex {
  display: flex;
  align-items: center;
  width: 500px;
}

.tip-view {
  font-size: 12px;
  color: $color-danger;
  line-height: 18px;
  margin-top: 10px;
}
.tip {
  font-size: 12px;
  color: #909399;
  line-height: 36px;
  margin-top: 2px;
}
.upload-template {
  width: 600px;
}
</style>