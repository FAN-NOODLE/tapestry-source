<template>
  <div class="add-edit-group">
    <el-dialog
      :title="dialogType === 1 ? '上传图片' : '编辑图片'"
      :visible.sync="show_dialog"
      width="1000px"
      @close="closeDialog"
    >
      <cy-form
        class="main-box"
        ref="cyForm"
        :model="ruleForm"
        label-width="120px"
        :rules="rules"
      >
        <!-- 素材分组选择 -->
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

        <!-- 素材标签选择 -->
        <cy-form-item label="素材标签" prop="category_ids" trigger="change">
          <div class="flex">
            <el-select
              style="width:400px"
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

        <!-- 图片上传 -->
        <cy-form-item label="本地图片" required prop="sources" trigger="change">
          <upload-image
            class="upload-template"
            :urls="ruleForm.sources"
            accept="image/jpeg,image/png"
            img-type="batch"
            :before-upload="beforeUpload"
            :limit="ruleForm.show_type === 1 ? 20 : 9"
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
        <cy-form-item label="图片展示样式" prop="show_type" trigger="change" v-if="epName === 'coach'">
          <el-radio-group v-model="ruleForm.show_type" :disabled="dialogType === 2">
            <el-radio :label="1">单张展示(最多20张)</el-radio>
            <el-radio :label="2">合并展示(最多9张)</el-radio>
          </el-radio-group>
        </cy-form-item>
        <cy-form-item label="关联商品" >
          <span class="tip-text">勾选识图后可自动加载</span>
          <template v-for="item in ruleForm.sources">
            <relation-goods  :key="item.hash || item.id" v-if="item.is_ai || (item.relation_goods && item.relation_goods.length > 0)" :item="item" @confirm="handleSelectGoodsConfirm(item, $event)" @delete="handleDeleteGoods(item, $event)"></relation-goods>
          </template>
        </cy-form-item>
      </cy-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeDialog" size="small">取 消</el-button>
        <el-button type="primary" @click="submitFn" size="small"
          >确 认</el-button
        >
      </span>
    </el-dialog>
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
import addCategory from "./add-category.vue";
import Times from "./datePicker";
import PerModule from "@/store/modules/permission";
import SelectGoods from "./select-goods.vue";
import RelationGoods from "./relation-goods.vue";
interface T {
  group_id: number;
  url: string;
  id: number;
}
@Component({
  components: {
    UploadImage,
    addCategory,
    SelectGoods,
    RelationGoods
  }
})
export default class extends Vue {
  @Prop() dialogType!: number;
  @Prop() groupId!: number;
  @Prop() name!: string;
  @Prop() editTextDate?: T;
  //------------------------------------------------------data-----------------------------------------------------------/
  private show_dialog = true;
  private testFlag = false;
  private groupList: any = []; //素材列表
  private ruleForm: any = {
    group_id: 0,
    name: "",
    category_ids: [], //素材标签
    datePicker: Times.times,
    start_time: Times.start, //有效期开始时间
    end_time: Times.end, //有效期结束时间
    // todo: 1. 单张展示 2. 合并展示
    show_type: 1,
    // todo:素材集合
    sources: []
  };
  private rules = {};
  // 素材标签
  private category_options: any[] = [];
  private addVisible = false; // 控制添加素材标签弹窗
  // todo: 当前品牌
  private epName = PerModule.getEpName;
  //----------------------------------------------------methods-----------------------------------------------------------/
  // 上传成功
  private update(res: any) {
    const sortImgs = this.sortImage(res)
    console.log(sortImgs)
    this.testFlag = false;
    this.ruleForm.sources = sortImgs;
  }
  // todo:增加排序字段sort
  sortImage(arr: any){
    for(let i = 0;i < arr.length;i++){
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
  private sourceDetail() {
    Api.sourceDetail({
      source_id: (this.editTextDate as T).id,
      type: 2
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
        // todo:合并图片展示
        if(data.merge_image && data.merge_image.length >0 && this.epName === 'coach'){
          // 根据sort字段排序
          data.merge_image.sort((a,b) => a.sort - b.sort)
          this.ruleForm.sources = data.merge_image
        }else{
          this.ruleForm.sources =  [{
            id:data.source_id,
            name:data.name,
            url:data.url,
            is_ai: (data.relation_goods || []).filter((item: any) => item.add_type == 2).length > 0,
            relation_goods: data.relation_goods || []

          }]
          console.log(this.ruleForm.sources)
        }
        // todo:如果有merage_id 1. 单张展示 2. 合并展示
        this.ruleForm.show_type = data.merge_id ? 2:1;
        this.ruleForm.group_id = data.group_id || -1;
        if (!data.group_id) {
          // this.ruleForm.name = "全部分组";
          this.ruleForm.name = "未分组";
        } else if (data.group_id === -1) {
          this.ruleForm.name = "未分组";
        } else {
          this.ruleForm.name = data.group_name;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    this.testFlag = true;
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
    if (file.size > 10 * 1024 * 1024) {
      this.$message.error("请上传10MB大小以内的图片");
      return false;
    }
  }

  // 点击完成
  private submitFn() {
    console.log("this.testFlag", this.testFlag);
    if (this.testFlag) {
      this.$message.warning("图片正在上传中");
      return;
    }
    // todo: 判断单张上传source不能超过20张，合并上传source不能超过9张
    if(this.ruleForm.show_type === 1 && this.ruleForm.sources.length > 20){
      this.$message.warning("单张上传不能超过20张");
      return;
    }
    if(this.ruleForm.show_type === 2 && this.ruleForm.sources.length > 9){
      this.$message.warning("合并上传不能超过9张");
      return;
    }
    // todo：合并素材至少上传2张
    if(this.ruleForm.show_type === 2 && this.ruleForm.sources.length < 2){
      this.$message.warning("合并上传最少上传2张");
      return;
    }
    (this.$refs["cform"] as any).validate((valid: any) => {
      if (valid) {
        if (this.dialogType === 1) {
          this.imageAdd();
        } else {
          this.sourceUpdate();
        }
      }
    });
  }
  // 关闭弹框
  private closeDialog() {
    this.$emit("closeDialog");
  }
  // 编辑接口
  private sourceUpdate() {
    Api.sourceUpdate({
      source_type: 2,
      id: (this.editTextDate as T).id,
      group_id: this.ruleForm.group_id || 0,
      sources: this.ruleForm.sources,
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time,
      category_ids: this.ruleForm.category_ids,
      show_type: this.epName === 'coach' ? this.ruleForm.show_type : 1
    }).then((res: any) => {
      if (res.code == 200) {
        this.$message.success("编辑成功");
        this.closeDialog();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  // 图片上传接口
  imageAdd() {
    const params = {
      group_id: this.ruleForm.group_id || 0,
      images: this.ruleForm.sources,
      start_time: this.ruleForm.start_time,
      end_time: this.ruleForm.end_time,
      category_ids: this.ruleForm.category_ids,
      show_type: this.ruleForm.show_type
    };
    Api.imageAddApi(params).then((res: any) => {
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
  private openAddCate() {
    this.addVisible = true;
  }
  // 关闭取消弹窗
  private closeVisible() {
    this.addVisible = false;
  }
  // 新增素材标签成功
  private addCateSubmit() {
    this.getCategory();
  }

  // 获取数据
  async getCategory() {
    const param: object = {
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
  private datePickerChange(val) {
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
      item.is_ai = e;
      if(e){
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
  async updateAIImage(data: any) {
    const { index, is_ai } = data
    
    this.ruleForm.sources[index].is_ai = is_ai
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
  private mounted() {
    //素材分组列表
    this.grouplist();
    this.getCategory();
    if (this.dialogType === 2) {
      this.sourceDetail();
    } else {
      this.ruleForm.group_id = this.groupId || -1;
      this.ruleForm.name = this.groupId ? this.name : "未分组";
    }
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-dialog {
  max-height: 700px !important;
}
.main-box{
  max-height: 500px;
  overflow-y: auto;
}

.tip-view {
  font-size: 12px;
  color: $color-danger;
  line-height: 18px;
  margin-top: 10px;
}
.flex {
  display: flex;
  align-items: center;
  width: 500px;
}

.tip-view {
  margin-top: 10px;
}
.tip-text {
  font-size: 12px;
  color: #909399;
  margin: 4px 0;
}
</style>
