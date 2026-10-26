<template>
  <div>
    <el-form ref="form" :model="formInfo" :rules="formRules" label-width="140px">
      <el-form-item label="活动名称" prop="name">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入活动名称"
          v-model="formInfo.name"
          maxlength="30"
          show-word-limit
          :disabled="readonly"
        ></cy-input>
      </el-form-item>
      <el-form-item label="活动对外名称" prop="alias_name">
        <cy-input
          size="medium"
          width-size="medium"
          placeholder="请输入活动对外名称"
          v-model="formInfo.alias_name"
          maxlength="30"
          show-word-limit
          :disabled="readonly"
        ></cy-input>
        <div class="tips">非必填，用于前端活动页展示，未配置则不展示活动名</div>
      </el-form-item>
      <el-form-item label="活动banner图" prop="banner_url">
        <div style="display: flex; align-items: flex-end;">

          <upload-image
            class="upload-template"
            :urls="formInfo.banner_url"
            accept="image/jpeg,image/png"
            :limit="1"
            :before-upload="handleUploadBefore"
            @change="handleUploadChange"
          >
          </upload-image>
          <el-button class="m-l-10" type="text" size="small" @click="handlPreview">预览效果图</el-button>
        </div>
        <div class="tips">
          <span>非必填，未上传则展示默认图</span>
          <span class="tips-link">（查看默认图）
            <el-image class="img" :src="defaultImage" :preview-src-list="[defaultImage]"></el-image>
          </span>
          <span>；若需上传图片素材，可参考默认图比例，上传长高比为1.3:1的图片(推荐尺寸1125*837px)，支持png、jpg格式</span>
        </div>
      </el-form-item>
      <el-form-item label="活动有效期" prop="time_type">
        <el-radio-group v-model="formInfo.time_type">
          <el-radio :label="0">长期有效</el-radio>
          <el-radio :label="1">指定时间</el-radio>
        </el-radio-group>
        <el-form-item v-if="formInfo.time_type === 1" prop="start_time" :rules="formRules.start_time">
          <el-date-picker
            type="datetimerange"
            v-model="expirationDate"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            range-separator="至"
            :picker-options="pickerOptions"
            :default-time="['00:00:00', '23:59:59']"
          ></el-date-picker>
        </el-form-item>
      </el-form-item>
      <el-form-item label="投放渠道" prop="category_list">
        <el-select v-model="formInfo.category_list" placeholder="请选择投放渠道" multiple style="width:400px" value-key="id">
          <el-option
            v-for="item in JSON.parse(JSON.stringify(categoryList))"
            :key="item.id"
            :label="item.name"
            :value="item"
          ></el-option>
        </el-select>
        <div v-if="formInfo.category_list.length" class="m-t-10">关联点位配置<span class="tips">（可根据运营需求，在投放渠道下增加细分点位，一个投放点位可对应一个活动码/活动链接）</span></div>
        <div v-if="formInfo.category_list.length" class="list">
          <div class="item" v-for="item,index in formInfo.category_list" :key="item.id">
            <div class="item-title">
              <span>渠道：{{ item.name }}</span>
              <el-button class="m-l-10" type="primary" size="small" @click="handleAddPoint(item)">添加点位</el-button>
            </div>
            <div v-if="item.points?.length" class="point-list">
              <el-form-item label-width="auto" label="" v-for="(point, pointIndex) in item.points" :key="pointIndex" :prop="`category_list[${index}].points[${pointIndex}].name`" :rules="{
                ...formRules.point_name[0],
                data: item.points
              }">
                <el-row class="point-item">
                  <cy-input v-model="point.name" placeholder="请输入" maxlength="30" size="medium" width-size="medium" show-word-limit></cy-input>
                  <!-- <el-input v-model="point.name" placeholder="请输入" maxlength="30" show-word-limit></el-input> -->
                  <el-button
                    class="m-l-5"
                    type="primary"
                    icon="el-icon-delete"
                    circle
                    @click="handleRemovePoint(item, pointIndex)"
                  ></el-button>
                </el-row>
              </el-form-item>
            </div>
          </div>
        </div>
      </el-form-item>
      <el-form-item label="是否为客户打标签">
        <el-row>
          <el-checkbox v-model="formInfo.add_name_tag">活动名称标签</el-checkbox>
          <el-checkbox v-model="formInfo.add_category_tag" >渠道标签</el-checkbox>
        </el-row>
        <div class="tips">勾选后，从不同渠道添加的客户将自动打上活动名称标签及渠道标签</div>
      </el-form-item>
      <el-form-item label="自定义标签" prop="">
        <el-button v-if="!formInfo.tag_list.length" type="primary" @click="showSelectTag=true">选择标签</el-button>
        <div class="tips">通过该活动添加的客户，将自动打上所选标签</div>
        <el-row class="m-t-10" v-if="formInfo.tag_list.length">
          <el-table
              ref="multipleTable"
              :data="[{list: formInfo.tag_list}]"
            >
              <el-table-column label="标签">
                <template slot-scope="scope">
                  <div class="flex">
                    <el-tag
                      v-for="(item, index) in scope.row.list"
                      :key="index"
                      >{{ item.name }}</el-tag
                    >
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="72" fixed="right">
                <template>
                  <el-button @click="showSelectTag = true" type="text" size="small"
                    >修改</el-button
                  >
                </template>
              </el-table-column>
            </el-table>
        </el-row>
      </el-form-item>
      <el-form-item 
        label="门店范围" 
        prop="scope_area_list" 
        :rules="{
          ...formRules.scope_area_list[0],
          data: scopeList()
        }">
        <el-button type="primary" @click="showSelectMember=true">选择门店范围</el-button>
        <el-row class="m-t-10" v-if="scopeList().length>0">
          <el-table :data="[{list:scopeList()}]">
            <el-table-column prop="name" label="员工">
              <template slot-scope="scope">
                <div class="flex">
                  <el-tag class="" v-for="(item, index) in scope.row.list" :key="index">
                    <div class="flex-center">
                      <img
                        class="img m-r-5"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                        alt=""
                        v-if="item.type !==5"
                      />
                      <i v-else class="el-icon-user-solid m-r-5"></i>
                      <span >{{ item.name }}</span>
                    </div>
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column v-if="!readonly" label="操作" width="120">
              <el-button type="text" @click="showSelectMember=true">修改</el-button>
            </el-table-column>
          </el-table>
        </el-row>
      </el-form-item>
      <el-form-item label="自动通过" prop="">
        <el-switch v-model="formInfo.skip_verify" ></el-switch>
        <div class="tips">开启后，客户添加时无需经过确认自动成为好友</div>
      </el-form-item>
      <el-form-item label="是否配置欢迎语" prop="has_welcome">
        <el-switch v-model="formInfo.has_welcome" ></el-switch>
        <div class="tips">配置欢迎语后，通过该活动渠道码添加的客户，将自动发送配置的欢迎语</div>
      </el-form-item>
      <customer-welcome-form v-if="formInfo.has_welcome" ref="welcomeFormRef" :value="formInfo.welcome_list"></customer-welcome-form>
    </el-form>
    <!-- 导入人群包 -->
    <select-member
      :show.sync="showSelectMember"
      :value="scopeList()"
      @cancel="showSelectMember = false"
      @confirm="handleSelectMemberConfirm"
    ></select-member>
    <select-tag
      :show.sync="showSelectTag"
      :value="formInfo.tag_list"
      @cancel="showSelectTag = false"
      @confirm="handleSelectTagConfirm"
    ></select-tag>
    <!-- 图片预览 -->
    <el-dialog width="400px" title="预览效果图" :visible.sync="showPreview" @opened="handleOpened">
      <div style="text-align: center;">
        <div class="tips m-b-10">请用手机扫描下方二维码预览banner图效果</div>
        <canvas class="canvas" id="canvas"></canvas>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from 'vue-property-decorator'
import UploadImage from "@/components/business/upload-image/index.vue";
import CustomerWelcomeForm from "./customer-welcome-form.vue";
// import SelectMember from "@/components/business/select-member-copy/index.vue";
var QRCode = require('qrcode');
// https://docs.apipost.net/docs/detail/4c842b180c9a000?locale=zh-cn&target_id=8383b45746001

// 264239
// 选择员工
import SelectMember from "./select-member.vue";
import SelectTag from './select-tag.vue';
import { getCategoryListApi } from '@/api/enterprise/activity';


export interface ICategoryItem {
  id: number;
  name?:string;
  points: {
    id?: number;
    name: string;
  }[];
}

export interface IWelcomeItem {
  id: number;
  type: 0 | 1; // 0 文本，1 图片
  content: string[]; // 文本段落/图片 URL
}
export interface IFormInfo {
  enterprise_id: number;
  member_id: number;
  name: string;
  alias_name: string;
  banner_url: string;
  time_type: 0 | 1; // 1 固定时间，0 长期有效
  start_time: number; // Unix 时间戳（秒）
  end_time: number;   // Unix 时间戳（秒）
  skip_verify: boolean; // 0 需审核，1 无需审核
  add_name_tag: boolean;
  add_category_tag: boolean;
  tag_ids?: number[];
  tag_list: {
    id: number;
    name: string;
  }[];
  scope_member_ids?: number[];
  scope_member_list: {
    id: number;
    name: string;
  }[];
  scope_area_ids?: number[];
  scope_area_list: {
    id: number;
    name: string;
  }[];
  category_list: ICategoryItem[];
  has_welcome: boolean;
  welcome_list: IWelcomeItem[];
}

export const defaultFormInfo: IFormInfo = {
  enterprise_id: 0,
  member_id: 0,
  name: '',
  alias_name: '',
  banner_url: '',
  time_type: 0,
  start_time: 0,
  end_time: 0,
  skip_verify: true,
  add_name_tag: false,
  add_category_tag: true,
  tag_ids: [],
  scope_member_ids: [],
  scope_area_ids: [],
  category_list: [],
  welcome_list: [],
  tag_list: [],
  scope_member_list: [],
  scope_area_list: [],
  has_welcome: false
}
@Component({
  components: {
    UploadImage,
    SelectMember,
    SelectTag,
    CustomerWelcomeForm
  }
})
export default class EditForm extends Vue {
  /** 是否只读 */
  @Prop({ type: Boolean, default: false }) readonly: boolean;
  /** 表单数据 */
  @Prop({ type: Object, default: () => defaultFormInfo }) data: IFormInfo;
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))
  /** 表单校验规则 */
  formRules = {
    name: [{ required: true, 
      validator: (rule: any, value: any, callback: any) => {
        if((value || '').trim().length === 0) {
          callback(new Error('请输入名称，活动名称不能为空'));
        } else {
          callback();
        }
      },
      trigger: 'blur' }],
    point_name: [{ required: true, 
      validator: (rule: any, value: any, callback: any) => {
        if((value || '').trim().length === 0) {
          callback(new Error('请输入渠道点位名称，渠道点位名称不能为空'));
        } else if(rule.data.filter((item:any) => item.name === (value || '').trim()).length > 1) {
          callback(new Error('渠道点位名称已存在'));
        } else {
          callback();
        }
      },
      trigger: 'blur' }],
    time_type: [{ required: true, message: '请选择活动有效期', trigger: 'change' }],
    start_time: [
      { required: true, 
        validator: (rule: any, value: any, callback: any) => {
  
          if (!value) {
            return callback(new Error('请选择活动有效期'));
          }
          return callback();
      }  , message: '请选择活动有效期', trigger: 'change' }
    ],
    category_list: [
      { required: true, message: '请选择投放渠道', trigger: 'change' }
    ],
    scope_area_list: [{required: true,  validator: (rule: any, value: any, callback: any) => {
      if (rule.data.length === 0) {
        return callback(new Error('请选择门店范围'));
      }
      callback();
    } , trigger: 'change' }],
  }
  
  /** 选择员工 */
  showSelectMember = false
  /** 选择标签 */
  showSelectTag = false

  pickerOptions = {
    disabledDate(time) {
      return time.getTime() < new Date().setHours(0,0,0,0);
    }
  }

  categoryList:any[] = []

  defaultImage = require('@/assets/images/default.png')

  showPreview = false

  get selectMember() {
    return []
  }


  /** 执行时间 读取 */
  get expirationDate() {
    const { start_time, end_time } = this.formInfo
    return [start_time*1000, end_time*1000].filter(Boolean)
  }
  /** 执行时间 写入 */
  set expirationDate(value: number[]) {
    const [start_time, end_time] = value || []
    this.formInfo.start_time = parseInt(start_time/1000+'');
    this.formInfo.end_time = parseInt(end_time/1000+'');
  }

  get scopeList() {
    return () => {
      return [...this.formInfo.scope_area_list, ...this.formInfo.scope_member_list]
    }
  }


  @Watch('data')
  onChange(value: IFormInfo) {
    if(value) {
      console.log(11, value);
      
      this.formInfo = value
    }
  }

  mounted() {
    if(document.getElementsByName('v-modal')[0]) {
      document.getElementsByName('v-modal')[0].style.display='none'
    }
    this.getCategory()
  }
  handlPreview() {
    this.showPreview = true
  }
  handleOpened() {
    this.$nextTick(() => {
      const url = 'https://tapestry-web-uat.coach.com.cn/work/457a684cb1477448d1798ef36a510102/user/preview?url=' + encodeURIComponent(this.formInfo.banner_url)
      
      QRCode.toCanvas(document.getElementById('canvas'), url, function (error) {
        if (error) console.error(error)
        console.log('success!');
      })
    })
  }
  
  getCategory() {
    getCategoryListApi({}).then((res:any) => {
      if (res.code == 200) {
        this.categoryList = res.data
      } else {
        this.$message.error(res.message);
      }
    })
  }

  /**
   * 上传前校验
   * @param file 
   */
  async handleUploadBefore(file: File) {
    if (file.size > 10 * 1024 * 1024) {
      this.$message.error("请上传大小不超过10MB的图片");
      throw new Error("请上传大小不超过10MB的图片");
    }
    return true
  }

  /**
   * 上传成功
   * @param e 
   */
  handleUploadChange(e: string) {
    this.formInfo.banner_url = e
  }

  /**
   * 添加点位
   */
  handleAddPoint(item:ICategoryItem) {
    if(!item.points) this.$set(item, 'points', [])
    item.points.push({id: 0, name: ''})
  }

  /**
   * 删除点位
   */
  handleRemovePoint(item:ICategoryItem, index: number) {
    item.points.splice(index, 1)
  }


  /**
   * 选择员工确认
   */
  handleSelectMemberConfirm(value:any[]) {
    console.log(value);
    
    this.showSelectMember = false;
    this.formInfo.scope_member_list = value.filter(item => item.type === 5)
    this.formInfo.scope_member_ids = this.formInfo.scope_member_list.map(item => item.id)
    this.formInfo.scope_area_list = value.filter(item => item.type !== 5)
    this.formInfo.scope_area_ids = this.formInfo.scope_area_list.map(item => item.id)
  }
  /**
   * 选择标签
   */
  handleSelectTagConfirm(value) {
    console.log(value)
    this.showSelectTag = false;
    this.formInfo.tag_list = value;
    this.formInfo.tag_ids = value.map(item => item.id);
  }

  /**
   * 表单校验
   */
  validate(callback: (data: IFormInfo) => void) {
    (this.$refs.form as any).validate(async (valid: boolean) => {
      if (valid) {
        
        const params = JSON.parse(JSON.stringify(this.formInfo));
        params.name = params.name.trim();
        params.alias_name = params.alias_name.trim();
        params.skip_verify = params.skip_verify ? 1 : 0
        params.has_welcome = params.has_welcome ? 1 : 0
        params.add_name_tag = params.add_name_tag ? 1 : 0
        params.add_category_tag = params.add_category_tag ? 1 : 0

        params.category_list = params.category_list.map((item: any) => {
          return {
            id: item.id,
            points: (item.points||[]).map((point: any) => {
              return point.name.trim()
            }).filter(Boolean)
          }
        })
        if(params.has_welcome) {
          params.welcome_list = await (this.$refs.welcomeFormRef as any).save();
        } else {
          params.welcome_list = []
        }
        

        
        delete params.scope_area_list
        delete params.scope_member_list
        delete params.tag_list
        callback(params)
      }
    })
  }
}
</script>
<style lang="scss" scoped> 
.canvas {
  margin: auto;
}
.tips {
  font-size:12px;
  color: #999;
}
.tips-link {
  color: #409eff;
  position: relative;
  .img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 0;
  }
}
.list {
  margin-top: 15px;
}
.item {
  line-height: normal;
  // border: 1px solid #ebeef5;
  // padding: 10px 15px;
  // border-radius: 6px;
  margin-bottom: 10px;

  .item-title {
    background: #f6f7ff;
    padding: 10px 15px;
    // border-bottom: 1px solid #ebeef5;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .point-list {
    border: 1px solid #ebeef5;
    padding: 10px 15px;
    // border-radius: 6px;
    display: grid;
    // 三列
    grid-template-columns: repeat(3, 1fr);
    grid-gap: 15px;
  }

  .point-item {
    // width: 400px;
    display: flex;
    width: 100%;
    flex: 1;
  }



  ::v-deep .cy-input-medium {
    width: auto !important;
    flex: 1;
  }
}

.flex {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.flex-center {
  display: flex;
  align-items: center;
  overflow: hidden;

  .img {
    width: 14px;
    display: block;
  }
}
</style>