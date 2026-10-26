<template>
  <div class="child-task">
    <layout-form>
      <template slot="content">
        <cy-form
          ref="cform"
          label-width="90px"
          :rules="planRules"
          :model="planInfo"
        >
          <!-- 计划名称 -->
          <cy-form-item
            label="子任务类型"
            required
          >
            {{ planInfo.task_category === 0 ? ep.isCoach?"企微沟通任务":"沟通任务" : "朋友圈任务" }}
          </cy-form-item>
          <!-- 计划名称 -->
          <cy-form-item
            label="子任务名称"
            prop="childTaskName"
            required
          >
            <cy-input
              size="small"
              width-size="medium"
              placeholder="请输入子任务名称"
              v-model="planInfo.childTaskName"
              @input="accountInput"
              class="input-with-select"
              maxlength="20"
              show-word-limit
            ></cy-input>
          </cy-form-item>
          <!-- 子任务描述 -->
          <cy-form-item
            label="子任务描述"
            prop="text"
            required
          >
            <cy-input
              class="textarea-box"
              placeholder="请输入子任务描述内容"
              v-model="planInfo.text"
              :rows="4"
              type="textarea"
              maxlength="200"
              show-word-limit
            ></cy-input>
          </cy-form-item>
          <!-- 任务群发类型 企业ks/sw -->
          <cy-form-item
            label="任务群发类型"
            required
            v-if="planInfo.task_category === 1"
          >
            <el-radio-group
              v-model="planInfo.group_send_type"
              @change="changeGroupTask"
              style="width: 608px;"
            >
              <el-radio :label="1">企微一键群发(从客户朋友圈发送)</el-radio>
              <el-radio :label="0">{{ep.isCoach? '标签任务群发': '标签任务群发（从任务中心发送）'}}</el-radio>
            </el-radio-group>
            <div class="message-info">
              <span
                v-if="planInfo.task_category === 0">※每位客户/客户群每天可接收每个SA的1条群发消息，标签任务群发和企微一键群发共享额度；超过接收上限的客户/客户群将无法再收到群发消息。</span>
              <span v-if="planInfo.task_category === 1">※每位客户（除外部企微联系人）一天可以看到每个导购发送的最多三条任务群发的朋友圈内容。</span>
            </div>
          </cy-form-item>
          <cy-form-item
            label="任务群发类型"
            required
            v-if="ep.isCoach && planInfo.task_category===0"
          >
            <el-radio-group
              v-model="planInfo.group_send_type"
              @change="changeGroupTask"
              style="width: 608px;"
            >
              <el-radio :label="1">企微一键群发（从群发助手发送）</el-radio>
              <el-radio :label="2">灵活沟通转发（从任务中心发送）</el-radio>
            </el-radio-group>
            <div class="message-info">※每位客户/每个客户群每月最多可接收条数为当月天数，企微一键群发和灵活沟通转发共享额度，超过接收上限的客户/客户群将无法再收到群发消息。（当前为每位客人每天接收1条群发消息。可从企业微信-群发助手-设置群发规则中调整其它使用规则）</div>
          </cy-form-item>
          <cy-form-item
            label="群发顾客"
            required
            v-if="ep.isCoach && planInfo.task_category===0 && planInfo.group_send_type === 1"
          >
            <el-radio-group
              v-model="planInfo.allow_select"
              style="width: 608px;"
            >
              <el-radio :label="0">导购不可编辑</el-radio>
              <el-radio :label="1">导购可编辑</el-radio>
            </el-radio-group>
          </cy-form-item>
          <!-- 任务完成标准 -->
          <el-form-item
            v-if="ep.isCoach && planInfo.group_send_type === 2"
            label="任务完成标准"
            prop="finish_rate"
            :rules="planRules.finish_rate"
          >
            <el-row class="flex">
              <span class="m-r-5">顾客名单触达率达到</span>
              <el-input v-model="planInfo.finish_rate" style="width:100px"></el-input>
              <span class="m-l-5">%</span>
            </el-row>
          </el-form-item>
          <!-- 子任务有效期 -->
          <cy-form-item v-if="ep.isCoach" label="子任务有效期" required style="max-width: 800px;">
            <el-radio-group v-model="planInfo.task_expire_type">
              <el-radio v-for="(item, index) in taskExpireTypeList" :key="index" :label="item.value">{{ item.label }} <span class="tips">{{ item.tips }}</span></el-radio>
            </el-radio-group>
            <cy-form-item class="picker" v-if="planInfo.task_expire_type === 1" prop="task_start_time">
              <el-date-picker 
                class="m-t-10"
                v-model="taskExpireTime"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                value-format="timestamp"
                :default-time="defaultTime"
                :picker-options="pickerOptions"
              ></el-date-picker>
            </cy-form-item>
          </cy-form-item>
          <!-- 条件维度 编辑的时候不可修改 -->
          <cy-form-item
            label="条件维度"
            required
          >
            <cy-select
              class="input-with-select"
              v-model="planInfo.dimension"
              placeholder="请选择子任务维度"
              @change="changeDimension"
              :optionsList="dimension_list"
              :options="{ label: 'name', value: 'value' }"
              :disabled="ep.isCoach && task_id && $route.query.status != '2'"
            >
            </cy-select>
          </cy-form-item>
          <!-- 选择员工 当 dimension 1或者2 展示 -->
          <cy-form-item
            label="员工范围"
            prop="choiceStaff"
            class="with-auto"
            v-if="getDimension('choose_staff')"
          >
            <el-radio-group v-model="chooseMemberType">
              <el-radio :label="1">选择员工</el-radio>
              <el-radio :label="2">导入员工</el-radio>
            </el-radio-group>
            <div class="m-t-10 choose-member">
              <el-button
                v-if="chooseMemberType == 1"
                size="small"
                type="primary"
                @click="showChoiceStaff = true"
              >选择员工</el-button>
              <el-button
                v-else
                size="small"
                type="primary"
                @click="staffImportDialog = true"
              >上传员工</el-button>
            </div>
            <div v-if="planInfo.range && planInfo.range.length > 0" style="width: 800px;">
              <!-- 标题 -->
              <div class="title-item flex">
                <div>员工</div>
                <div>操作</div>
              </div>
              <!-- 可见范围显示区 -->
              <div class="range-content flex flex-align-center">
                <div class="left">
                  <div
                    v-for="(rangeItem, rangeIndex) in planInfo.range"
                    :key="rangeIndex"
                    class="flex item flex-align-center"
                  >
                    <img
                      class="group-img"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                      alt=""
                      v-if="rangeItem.is_dp"
                    />
                    <i v-else class="el-icon-user-solid"></i>
                    <div class="item-name font-re">
                      {{ rangeItem.label }}
                    </div>
                  </div>
                </div>
                <!-- 修改按钮 -->
                <el-button
                  class="right font-re pointer-cursor"
                  type="text"
                  @click="chooseMemberType == 1 ? (showChoiceStaff = true) : (staffImportDialog = true)"
                >
                  {{ chooseMemberType == 1 ? "修改" : "重新导入" }}
                </el-button>
              </div>
            </div>
          </cy-form-item>
          <!-- dimension 3.4的时候需要下载模板和导入 -->
          <template v-if="getDimension('down_url')">
            <cy-form-item v-if="task_id && $route.query.status != '2'" label="上传名单" required> 
              <div class="link" @click="downLoadFile(planInfo.file_name, planInfo.file_url)">
                <i class="el-icon-document"></i>
                {{ planInfo.file_name }}</div>
              <div class="message-info m-t-5">名单已在解析中，暂不支持编辑，如需调整可重新复制该子任务，上传新的人群包数据</div>
            </cy-form-item>
            <div v-else>
              <cy-form-item label="下载模板">
                <el-button
                  type="text"
                  @click="temDown(getDimension('down_url'))"
                >{{
                  getDimension("down_text")
                }}</el-button>
              </cy-form-item>
              <cy-form-item
                label="上传名单"
                required
              >
                <el-upload
                  ref="upload"
                  type="file"
                  :limit="1"
                  :on-exceed="excedFun"
                  :data="exceldata"
                  :action="uploadUrl"
                  :before-upload="beforeUpload"
                  :headers="headers"
                  :on-success="onSuccessFun"
                  :file-list="file_list"
                  :on-remove="onRemoveFn"
                  drag
                >
                  <i class="el-icon-upload"></i>
                  <div class="el-upload__text">
                    将文件拖到此处，或<em>点击上传</em>
                    <div>只能上传xlsx文件</div>
                  </div>
                </el-upload>
              </cy-form-item>
            </div>
          </template>
          <!-- 生成任务标签: ks/sw(朋友圈任务 标签群发任务) -->
          <template v-if="planInfo.group_send_type === 0 && getDimension('down_url')">
            <cy-form-item
              label="生成任务标签"
              required
            >
              <el-radio-group v-model="planInfo.is_add_tag" :disabled="ep.isCoach && task_id && $route.query.status != '2'">
                <el-radio :label="1">是</el-radio>
                <el-radio :label="0">否</el-radio>
              </el-radio-group>
            </cy-form-item>
            <!-- 生成任务标签form组件 -->
            <taskDefaultValue
              v-if="planInfo.is_add_tag === 1"
              :tagName="planInfo.childTaskName"
              :isDimissionChange="isDimissionChange"
              :tagExpireForm="tagExpireForm"
              ref="tagValueForm"
              @save="defaultValueSave"
            ></taskDefaultValue>
          </template>
          <!-- dimension 1 2时有条件内容 -->
          <template v-if="getDimension('has_condition')">
            <cy-form-item
              label="条件内容"
              required
            >
              <cy-select
                class="input-with-select"
                v-model="planInfo.condition_type"
                placeholder="请选择条件内容"
                @change="changeConditionType"
                :optionsList="condition_list"
                :options="{ label: 'name', value: 'value' }"
              >
              </cy-select>
            </cy-form-item>
          </template>
          <cy-form-item
            required
            label="任务发送时间"
            prop="birthday"
            v-if="planInfo.condition_type == 8"
          >
            <el-input
              v-model="planInfo.birthday"
              @input="inputNumber($event)"
            >
              <template slot="prepend">提前</template>
              <template slot="append">天</template>
            </el-input>
          </cy-form-item>
          <!-- 任务周期 -->
          <cy-form-item
            label="任务计划状态"
            required
            v-if="planInfo.dimension === 2"
          >
            <cy-select
              class="cy-select"
              :disabled="planInfo.condition_type == 8"
              v-model="planInfo.task_type"
              placeholder="请选择任务计划状态"
            >
              <cy-option
                label="单次"
                :value="1"
              ></cy-option>
              <cy-option
                label="持续"
                :value="2"
              ></cy-option>
            </cy-select>
            <p class="status-desc">
              单次任务：即在任务有效期内，当匹配任务条件后，仅执行一次的任务;
            </p>
            <p class="status-desc">
              持续任务：即在任务有效期内，当匹配任务条件后，有效期内每天执行一次的任务;
            </p>
          </cy-form-item>
          <!-- 任务对象 -->
          <cy-form-item
            label="任务对象"
            required
            v-if="planInfo.dimension === 2 || planInfo.dimension === 3"
          >
            <el-radio
              v-model="planInfo.send_type"
              :label="1"
              :disabled="!ep.isCoach && planInfo.condition_type !== 9 || ep.isCoach && task_id && $route.query.status != '2' && planInfo.dimension !== 2"
            >关联导购</el-radio>
            <el-radio
              v-model="planInfo.send_type"
              :label="2"
              :disabled="!ep.isCoach && planInfo.condition_type === 9 || ep.isCoach && task_id && $route.query.status != '2' && planInfo.dimension !== 2"
            >管理导购</el-radio>
            <!-- 提示文案 -->
            <div class="tips m-t-5">关联导购：添加了客户的所有SA；</div>
            <div class="tips m-t-10">管理导购：添加了客户的主SA；</div>
          </cy-form-item>
          <!-- 条件设置 -->
          <div
            class="header-condition-title"
            v-if="!getDimension('down_url')"
          >
            条件设置
          </div>
          <div
            class="risk-content-wrap"
            v-for="(item, index) in planInfo.termArr"
            :key="index + 'only'"
          >
            <div
              class="condition-title"
              v-if="!getDimension('down_url')"
            >
              条件{{ index + 1 }}
              <i
                class="el-icon-delete delete-icon"
                v-if="planInfo.termArr && planInfo.termArr.length > 1 && index > 0"
                @click.stop="deleteRiskEvn(index)"
              ></i>
            </div>
            <!-- 子任务条件 -->
            <cy-form-item
              v-if="planInfo.dimension === 1"
              :label="getlabelName()"
              :prop="`termArr.${index}.value`"
              required
            >
              <cy-select
                class="input-with-select w-120"
                v-model="item.term"
              >
                <!-- <cy-option label="小于等于" value="3"></cy-option> -->
                <cy-option
                  label="大于等于"
                  value="5"
                ></cy-option>
              </cy-select>
              <cy-input
                class="w-180"
                v-model="item.value"
                @input="changeNumberInput($event, index)"
              />
              <span class="text-position">人</span>
              <span class="m-l-40">推送如下子任务内容：</span>
            </cy-form-item>
            <!-- 选择标签 -->
            <cy-form-item
              v-if="planInfo.dimension === 2 && planInfo.condition_type != 9"
              label="客户标签"
              required
            >
              <el-button
                class="range-choose-btn"
                size="small"
                @click="selectTagsIndex(index)"
              >选择标签</el-button>
              <!-- 标题 -->
              <div
                v-if="item.customerLabel && item.customerLabel.length > 0"
                class="title-item flex"
              >
                <div class="flex">
                  <div class="font-me">
                    已选标签
                  </div>
                </div>
                <div>操作</div>
              </div>
              <!-- 可见范围显示区 -->
              <div
                v-if="item.customerLabel && item.customerLabel.length > 0"
                class="range-content flex flex-align-center"
              >
                <div class="left">
                  <div
                    v-for="(items, i) in item.customerLabel"
                    :key="i"
                    class="flex item flex-align-center"
                  >
                    <div class="item-name font-re">
                      {{ items }}
                    </div>
                  </div>
                  <div
                    v-for="(items, i) in item.deleteTag"
                    :key="i"
                    class="flex item flex-align-center is-delete"
                  >
                    <div class="item-name font-re">
                      {{ items }}
                    </div>
                  </div>
                </div>
                <!-- 修改按钮 -->
                <div
                  class="right font-re pointer-cursor"
                  @click="selectTagsIndex(index)"
                >
                  修改
                </div>
              </div>
              <div
                class="tip"
                v-if="item.deleteTag && item.deleteTag.length > 0"
              >
                标红标签为：客户标签中已删除标签，已删除标签下次修改时不展示
              </div>
            </cy-form-item>
            <!-- 任务内容 -->
            <el-form-item
              class="risk-form-item"
              label="子任务内容"
              trigger="blur"
              :prop="`termArr.${index}`"
              :rules="planRules.content"
              :required="planInfo.group_send_type === 1"
            >
              <!-- 新增素材展示区 -->
              <div class="risk-content-box">
                <!-- 表情选择区域-->
                <div class="select-box-insert">
                  <add-emoji
                    @select="addEmoji($event, item, index)"
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
                  <el-button
                    class="insert-btn"
                    size="small"
                    @click="opentextDialog(index)"
                  ><img
                      class="btn-img"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/text_icon.png"
                      alt=""
                    /><span>添加文本</span></el-button>
                </div>
                <!--多行文本区域-->
                <div>
                  <el-input
                    :ref="'textarea' + index"
                    class="textarea-box"
                    placeholder="请输入营销话术内容"
                    v-model="item.text.content"
                    :rows="4"
                    type="textarea"
                    maxlength="500"
                    show-word-limit
                  ></el-input>
                </div>
                <!-- 选择素材 -->
                <meterial-select
                  :item="item"
                  :index="index"
                  :planInfo="planInfo"
                  :showMaterialModal="showMaterialModal"
                  :textDialog="textDialog"
                  @openSelectMaterial="openSelectMaterial"
                  @deleteMaterial="deleteMaterial"
                  @getSelectMaterial="getSelectMaterial"
                  @textData="textData"
                  @getSelectCard="getSelectCard"
                ></meterial-select>
                <!-- 选择素材提示 -->
                <div
                  class="material-tip"
                  v-if="planInfo.task_category === 1"
                >
                  <ul>
                    <li>最多支持9个图片 类型，或者1个视频，或者1个链接。类型只能三选一</li>
                    <li>
                      选择【图片】素材：像素不超过1555200。图片大小不超过10M。最多支持传入9个；（普通图片尺寸建议不超过1440*1080px，长图片，长边建议不超过10800px）
                    </li>
                    <li>选择【视频】素材：最长不超过30S，最大不超过10MB。</li>
                  </ul>
                </div>
                <div
                  class="material-tip"
                  v-else-if="ep.isCoach"
                >
                  <ul>
                    <li>【本地上传】可支持上传图片及视频素材，最大不超过10M</li>
                  </ul>
                </div>
              </div>
            </el-form-item>
          </div>
          <!-- <el-button
            class="range-choose-btn margin-btn"
            v-if="planInfo.termArr.length < 999 && !getDimension('down_url') && planInfo.condition_type != 8"
            size="small"
            @click="addTerm"
            >添加条件
          </el-button> -->
        </cy-form>
      </template>
      <template slot="footer">
        <div>
          <el-button
            class="range-choose-btn"
            @click="save"
          >保存</el-button>
          <el-button @click="cancel">取消</el-button>
        </div>
      </template>
    </layout-form>
    <!-- 可见范围弹窗  -->
    <!-- <selectMember
      v-if="showChoiceStaff"
      :hasSelectArr="planInfo.range"
      @close="showChoiceStaff = false"
      @submit="submitDialog"
    /> -->
    <!-- <select-member      
      v-if="showChoiceStaff"
      :show.sync="showChoiceStaff"
      :value="planInfo.range"
      rowKey="user_id"
      @cancel="showChoiceStaff = false"
      @confirm="submitDialog"
    /> -->
    <select-member      
      :show.sync="showChoiceStaff"
      :value="planInfo.range"
      @cancel="showChoiceStaff = false"
      @confirm="submitDialog"
    />
    <!-- 客户标签弹窗  -->
    <select-tag-dialog
      v-if="tagDialog"
      :updateSelectedTags="selectedTags"
      @selectedTags="selectedTagsFn"
    ></select-tag-dialog>
    <!-- 上传员工人群包弹窗 -->
    <upload-file
      v-if="staffImportDialog"
      @close="staffImportDialog = false"
      @submit="submitDialog"
    ></upload-file>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/wisdom/planTask";
import UserModule from "@/store/modules/user";
// 选择员工
import SelectMember from "./select-member.vue";
// 上传员工人群包弹窗
import uploadFile from "./upload-file.vue";
import { VUE_APP_API_ROOT_RE, VUE_APP_API_VERSION } from "@/utils/env-config";
import LoginAboat from "@/utils/LoginManager";
import PerModule from "@/store/modules/permission";
// 选择表情弹窗 ../components/add-emoji.vue
import AddEmoji from "@/views/myapp/wisdom/marketing-calendar/components/add-emoji.vue";
// 选中标签
import selectTagDialog from "../../../../enterprise/components/select-tag-dialog.vue";
// 企业任务类型
import { tag_task, group_task } from "./dimension_data";
// 调整默认值弹窗
import taskDefaultValue from "./task-defalut-value.vue";
// 素材选择框拆离
import meterialSelect from "./material-select.vue";
import { ApiResult } from "@/types";
import settingApis from "@/api/setting/setting";
import { formatDate } from '@/utils/date'
import LocalUpload from '@/mixins/LocalUpload';
import { downloadFile } from "@/utils/download";
@Component({
  components: {
    SelectMember,
    uploadFile,
    AddEmoji,
    selectTagDialog,
    taskDefaultValue,
    meterialSelect,
  },
  mixins: [LocalUpload],
})
export default class extends Vue {
  
  private planInfo: any = {
    // 任务群发类型
    is_add_tag: 0,
    // 任务计划类型
    task_category: 0,
    tag_name: "",
    group_send_type: 1, // 任务群发类型（沟通任务） 1: 企微一键群发 2: 灵活沟通转发  任务群发类型（朋友圈任务） 1: 企微任务群发 0: 标签任务群发
    allow_select: 0, // 企微沟通任务 - 导购是否可以编辑 0 不可编辑 1 可编辑
    tag_expire_time: 0,
    tag_expire_type: 3,
    birthday: "", // 提前多少天
    consumer_file_id: 0,
    childTaskName: "",
    text: "",
    range: [], // 可见范围
    dimension: 3,
    condition_type: 1,
    task_type: 1,
    send_type: 1,
    termArr: [
      {
        term: "5",
        value: "",
        deleteTag: [],
        deleteTagId: [],
        customerLabel: [], //已选标签内容
        customerLabelIds: [], //已选标签的ID
        text: {
          content: "", // 任务内容
        },
        content: {
          // 素材内容
          videoId: "0",
          videoUrl: "",
          videoFlag: false,
          appletId: "0",
          applet: {},
          pictureId: [],
          picture: [],
          video: {},
          link: {},
          card: {},
        },
        hasMaterial: false, // 是否有已选素材
        currentMaterialType: "", // 当前选择素材的类型
      },
    ],
    // 任务完成标准
    finish_rate: 100,
    // 子任务的有效期类型，0为根据主任务的有效期一致，1指定时间有效期
    task_expire_type: 0,
    // 子任务的开始时间
    task_start_time: 0,
    // 子任务的结束时间
    task_end_time: 0,
  };
  private headers = {
    "Enterprise-Hash": LoginAboat.getRealHash(),
    "Auth-Token": LoginAboat.getToken(),
  };
  private chooseMemberType = 1; // 选择员工的方式 1 组件选 2 上传选择
  private staffImportDialog = false;
  private validatePlanName = (rule: any, value: any, callback: Function) => {
    if (value == "") {
      callback(new Error("请输入子任务名称"));
    }
    if (value.length > 20) {
      callback(new Error("子任务名称不能超过20个字"));
    } else {
      callback();
    }
  };
  private validateBirthday = (rule: any, value: any, callback: Function) => {
    if (this.planInfo.condition_type == 8 && value == "") {
      callback(new Error("请输入任务发送时间"));
    } else {
      callback();
    }
  };
  private planRules = {
    childTaskName: [{ validator: this.validateBirthday, trigger: ["blur", "change"] }],
    birthday: [{ validator: this.validatePlanName, trigger: ["blur", "change"] }],
    finish_rate: [{ 
      required: true,
      validator: (rule, value, callback) => {
        if (!/^[1-9][0-9]?$|100/.test(value) || Number(value)>100) {
          callback(new Error("请输入1-100的整数"));
        } else {
          callback();
        }
      }
    }],
    task_start_time: [{ 
      info: this.planInfo,
      validator: (rule, value, callback) => {
        const { st, et } = this.$route.query as any;
        if(!value) {
          callback(new Error("请选择任务有效期"));
        } else if(value < Number(st) || rule.info.task_end_time > Number(et)) {
          callback(new Error("仅可选择主任务有效期范围内的时间"));
        } else {
          callback();
        }
      }
    }],
    content: [{ validator: (rule, value, callback) => {
      console.log(1111, value);
      if(this.planInfo.group_send_type != 1) {
        return callback();
      }
      if (value.text.content == "" && value.content.picture.length == 0 && !value.content.link.link_url && !value.content.video.url && !value.content.applet.app_name) {
        callback(new Error("请填写任务内容"));
      } else {
        callback();
      }
    }, message: "请输入子任务内容", trigger: ["blur", "change"] }],
  };
  // 数据抽离
  private dimension_list = group_task;
  private isDimissionChange = false;
  private tagExpireForm: any = {
    name: "",
    radio: 1,
    date: 0,
    days: 10,
  };
  private condition_list: any = [];
  private status = 0; // 页面状态 0-新增 1-编辑 2-复制
  private plan_id = 0; // 任务计划ID
  private task_id = 0; // 子任务的ID
  private loading = true;
  private childTerm = "群人数";
  private childTermText = "人";
  
  private taskExpireTypeList = [
    {
      label: "主任务时间有效",
      value: 0,
    },
    {
      label: "按指定时间有效",
      value: 1,
      tips: '（*仅可选择主任务有效期范围内的时间）'
    },
  ]
  private showChoiceStaff = false; //选择员工的弹窗显隐
  private termIndex = 0; // 条件的下标
  private showMaterialModal = ""; // 控制当前选择素材的弹框显隐
  private textDialog = false; // 添加文本弹窗
  private tagDialog = false; // 选择客户标签
  private selectedTags: Array<any> = [];
  private delet_task_conditions: Array<number> = []; //记录编辑状态删除的条件
  private exceldata: any = {}; // 上传时附带的额外参数
  private uploadUrl = `${VUE_APP_API_VERSION}task/csvFile/upload`;
  private upFileId = 0;
  private file_list = [];
  private showDefalueValue = false;
  get ep() {
    return this.$store.getters['enterprise']
  }
  // 计算属性
  get USER() {
    return UserModule.getUserInfo;
  }

  get defaultTime() {
    const { st, et } = this.$route.query as any;
    console.log([formatDate(st*1000, 'hh:mm:ss'), formatDate(et*1000, 'hh:mm:ss')]);
    
    return [formatDate(st*1000, 'hh:mm:ss'), formatDate(et*1000, 'hh:mm:ss')]
  }

  // 子任务有效期不能大于主任务有效期
  pickerOptions = {
    disabledDate:(time: Date)=> {
      const currentStart = new Date(time).setHours(23, 59, 59);
      const currentEnd = new Date(time).setHours(0, 0, 0);
      
      const { st, et } = this.$route.query as any;
      return currentStart < Math.floor(st *1000) || currentEnd > Math.floor(et*1000);
    },
  };

  // 关联任务有效期
  get taskExpireTime() {
    const { task_start_time, task_end_time } = this.planInfo
    if(!task_start_time) {
      return []
    }
    return [task_start_time*1000, task_end_time*1000]
  }
  set taskExpireTime(value) {
    const [st, et] = value || []
    this.planInfo.task_start_time = st ? Math.floor(st/1000) : 0
    this.planInfo.task_end_time = et ? Math.floor(et/1000) : 0
  }

  /**
   * *********** 任务类型 task_category **********
   * 0 企微沟通任务
   * 1 朋友圈任务 
   * 
   * *********** 任务群发类型 group_send_type **********
   * 1 企微一键群发
   * 2 标签任务群发
   * 
   * *********** 条件维度 dimension **********
   * 1 群发
   * 2 用户画像
   * 3 导入顾客名单
   * 4 导入顾客&导购匹配关系
   * 
   * ********** 总结 **********
   * 任务未开始 & (企微沟通任务 | 朋友圈任务 & 用户画像)
   * 子任务可以编辑
   * 条件维度不可修改
   *    企微沟通任务 & 导购顾客名单
   *    上传名单不可修改
   *    任务对象不可修改
   *    企微沟通任务 & 导入顾客&导购匹配关系
   *    上传名单不可修改
   */
  

  /* 生命周期 */
  private async mounted() {
    await this.getEpConfig()
    // 增加任务类型
    const { status, plan_id, task_id, type } = this.$route.query;
    this.status = Number(status);
    this.plan_id = Number(plan_id);
    this.task_id = Number(task_id);
    this.planInfo.task_category = Number(type);
    // Coach&ks默认为管理导购，可修改为关联导购；SW默认为关联导购，可修改为管理导购
    this.planInfo.send_type = this.ep.isSW ? 1 : 2;
    // coach 新建的企微沟通任务不需要任务群发类型
    if(this.ep.isCoach && Number(type) === 0 && !task_id) {
      this.planInfo.group_send_type = 1;
    } else {
      this.planInfo.group_send_type = 1;
    }
    this.init();
  }
  /** 获取企业身份，区分ks、coach、sw */
  private async getEpConfig() {
    const res: ApiResult = await settingApis.getEpConfigApi();
    if (res.code !== 200) return this.$message.error(res.message as string);

    // 企业名称：coach;ks;sw
    const { name } = res.data;
    // 缓存企业名称，给详情页面用
    PerModule.setEqName(name);
  }
  // 调整默认值确定
  defaultValueSave(defaultValue: any) {
    console.log(defaultValue);
    this.planInfo.tag_name = defaultValue.name;
    this.planInfo.tag_expire_type = defaultValue.tag_expire_type;
    this.planInfo.tag_expire_time = defaultValue.tag_expire_time;
    this.showDefalueValue = false;
  }
  // 寻找dimension_list对应的key值
  getDimension(key: any) {
    //找到planInfo中的dimension对应的值
    const dimension = this.planInfo.dimension;
    const dimension_list = this.dimension_list;
    const obj = dimension_list.find((item: any) => {
      return item.value === dimension;
    });
    return obj ? obj[key] : "";
  }
  // 任务群发类型change
  changeGroupTask(val: number) {
    if(this.ep.isCoach && this.task_id && this.$route.query.status != '2') {
      return
    }
    // 清空dimension
    this.dimension_list = [];
    this.planInfo.dimension = 3;
    if (val === 1 || val === 2) {
      this.dimension_list = group_task;
      this.isDimissionChange = false;
      this.planInfo.is_add_tag = 0;
    } else {
      // 如果是朋友圈任务删除dimension_list社群数据维度
      if (this.planInfo.task_category === 1) {
        // 如果是朋友圈任务删除dimension_list社群数据维度
        this.dimension_list = tag_task.filter((item: any) => {
          return item.value !== 1;
        });
      } else {
        this.dimension_list = tag_task;
      }
      this.tagExpireForm = {
        name: this.planInfo.childTaskName,
        radio: 1,
        date: 0,
        days: 10,
      };
      this.isDimissionChange = true;
      this.planInfo.is_add_tag = 1;
    }
    // 清空上传文件
    this.file_list = [];
    this.planInfo.consumer_file_id = 0;
  }
  private getlabelName() {
    const obj = this.condition_list.find((item: { value: any }) => {
      return item.value === this.planInfo.condition_type;
    });
    return obj ? obj.name : "";
  }
  private getSelectCard(item: any) {
    if (item) {
      // this.selectCard = item;
      this.planInfo.termArr[this.termIndex].content[this.showMaterialModal] = item;
      this.planInfo.termArr[this.termIndex].hasMaterial = true;
    }
    this.showMaterialModal = "";
  }
  // 模板下载
  private temDown(url: string) {
    window.location.href = `${VUE_APP_API_ROOT_RE}${url}?task_type=${this.planInfo.dimension}`;
  }
  // 改变任务标签
  private changeConditionType(val: number) {
    if (this.planInfo.condition_type === 8) {
      this.planInfo.task_type = 2;
    }
    // 如果选择全部顾客，则隐藏客户标签，关联导购
    if (val === 9) {
      this.planInfo.send_type = 1;
    }
  }
  // 上传文件个数超过范围
  private excedFun() {
    this.$message.error("上传失败，每次仅支持上传一个文件！");
  }
  // 上传文件限制规范
  private beforeUpload(file: any) {
    this.exceldata.task_type = this.planInfo.dimension;
    const testmsg = file.name.substring(file.name.lastIndexOf(".") + 1);
    const extension = testmsg === "xlsx";
    if (!extension) {
      this.$message.error("请上传 xlsx 类型的文件!");
      return false;
    }
  }
  // 上传成功
  private onSuccessFun(res: any) {
    if (res.code == 200) {
      this.planInfo.consumer_file_id = res.data.id;
    } else {
      this.file_list = [];
      this.$message.error(res.message);
    }
  }
  // 移除上传
  private onRemoveFn() {
    this.file_list = [];
    this.planInfo.consumer_file_id = 0;
  }
  private init() {
    // 0-新增 1-编辑 2-复制
    if (this.status == 0) {
      this.loading = false;
      this.changeGroupTask(this.planInfo.group_send_type)
    } else {
      // 查看任务详情
      this.getChildTaskDetail();
    }
  }
  private accountInput(val: string) {
    // const str = val.replace(/[^\a-\z\A-\Z0-9\u4E00-\u9FA5]/g, "");
    this.$nextTick(() => {
      this.planInfo.childTaskName = val;
    });
  }
  private inputNumber(val: any) {
    const inputValue = String(val);
    const str = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    this.$nextTick(() => {
      this.planInfo.birthday = str;
    });
  }
  private changeNumberInput(val: any, index: number) {
    const inputValue = String(val);
    let str: any = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    if (str > 500) {
      str = 500;
      this.$message.warning("最大人数为500人");
    }
    this.$nextTick(() => {
      this.planInfo.termArr[index].value = str;
    });
  }
  // 可见范围
  private submitDialog(list: any) {
    console.log(list);
    
    this.planInfo.range = list;
    this.showChoiceStaff = false;
    this.staffImportDialog = false;
  }

  // 改变子任务维度
  private changeDimension(val: number, condition_type?: number) {
    // Coach&ks默认为管理导购，可修改为关联导购；SW默认为关联导购，可修改为管理导购
    this.planInfo.send_type = this.ep.isSW ? 1 : 2;
    this.planInfo.condition_type = "";
    this.planInfo.dimension = val;
    // 读取任务下面条件内容condition
    const dimension_list = this.dimension_list;
    const obj = dimension_list.find((item: any) => {
      return item.value === val;
    });
    const condition = obj ? obj.condition : [];
    if (condition && condition.length > 0) {
      // 复制一份condition
      this.condition_list = JSON.parse(JSON.stringify(condition));
      // 如果是沟通任务，则删除掉全部顾客标签
      if (this.planInfo.task_category === 0 || !this.ep.isCoach) {
        this.condition_list.splice(1, 1)
      }
      this.planInfo.condition_type = condition_type ? condition_type : condition[0].value;
    }
    // 清空上传文件
    if (val === 3 || val === 4) {
      this.file_list = [];
      this.planInfo.consumer_file_id = 0;
      this.isDimissionChange = true;
      if (this.planInfo.group_send_type === 0) {
        this.planInfo.is_add_tag = 1;
      }
    } else {
      // 清空任务标签
      this.planInfo.is_add_tag = 0;
      this.isDimissionChange = false;
    }
  }
  // 获取任务详情
  private getChildTaskDetail() {
    Api.childTaskDetail({ task_id: this.task_id }).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        // todo：任务群发类型
        this.planInfo = Object.assign({}, this.planInfo, data);
        this.tagExpireForm = {
          name: data.tag_name,
          radio: data.tag_expire_type,
          date: data.tag_expire_time * 1000,
          days: data.tag_expire_time,
        };
        // 获取任务维度
        this.changeGroupTask(data.group_send_type);
        this.changeDimension(data.dimension);
        // 获取任务条件
        this.changeConditionType(data.condition_type);
        this.planInfo.condition_type = data.condition_type;
        this.file_list = [];
        this.planInfo.consumer_file_id = 0;
        this.planInfo.childTaskName = this.status === 2 ? (data.name += "-副本") : data.name;
        this.planInfo.text = data.remark;
        this.planInfo.range = (data.plan_task_specify_crowd || []).map(item => {
          return {
            ...item,
            is_dp: item.third_id ? 1 : 0,
            value: item.user_id || item.third_id,
            label: item.user_name
          }
        });
        
        if (this.planInfo.dimension === 2 || this.planInfo.dimension === 3) {
          this.planInfo.task_type = data.task_type;
          this.planInfo.send_type = data.send_type;
        }
        if (data.condition_type == 8) {
          this.planInfo.birthday = data.birthday_pre_time + "";
        }
        if (data.plan_task_conditions && data.plan_task_conditions.length > 1) {
          this.planInfo.termArr = [];
          for (let i = 0; i < data.plan_task_conditions.length; i++) {
            this.planInfo.termArr.push({
              term: "5",
              value: "",
              deleteTag: [],
              deleteTagId: [],
              customerLabel: [], //已选标签name
              customerLabelIds: [], // 已选标签ID
              text: {
                content: "", // 任务内容
              },
              content: {
                // 素材内容
                videoId: "0",
                videoUrl: "",
                videoFlag: false,
                appletId: "0",
                applet: {},
                pictureId: [],
                picture: [],
                video: {},
                link: {},
                card: {},
              },
              hasMaterial: false, // 是否有已选素材
              currentMaterialType: "", // 当前选择素材的类型
            });
            if (this.status === 1) {
              this.planInfo.termArr[i].id = data.plan_task_conditions[i].id;
            }
          }
        } else {
          // 复制，编辑相当于编辑需要用到id
          if (this.status === 1 || this.status === 2) {
            this.planInfo.termArr[0].id = data.plan_task_conditions[0].id;
          }
        }
        const termArr = this.planInfo.termArr;
        // 处理标签或子任务条件
        data.plan_task_conditions.forEach((item: any, index: number) => {
          console.log(this.planInfo.dimension, "number");
          if (this.planInfo.dimension === 1) {
            termArr[index].value = item.conditions.logical.number;
            const expression = item.conditions.logical.expression;
            switch (expression) {
              case "ge":
                termArr[index].term = "5";
                break;
              case "le":
                termArr[index].term = "3";
                break;
              default:
                termArr[index].term = "5";
            }
          } else {
            item.conditions.tag &&
              item.conditions.tag.forEach((v: any) => {
                if (v.is_delete) {
                  termArr[index].deleteTag.push(v.tag_name);
                  termArr[index].deleteTagId.push(`${v.group_id},${v.tag_id}`);
                } else {
                  termArr[index].customerLabel.push(v.tag_name);
                  termArr[index].customerLabelIds.push(`${v.group_id},${v.tag_id},${v.qw_group_id},${v.qw_tag_id}`);
                }
              });
          }
          termArr[index].text.content = item.content.text?.content;
          // 包含图片
          if ("image" in item.content) {
            termArr[index].currentMaterialType = "picture";
            termArr[index].hasMaterial = true;
            termArr[index].content.picture = item.content.image.imgList;
            termArr[index].content.picture.map((item: any) => {
              return (item.id = item.media_id), (item.url = item.pic_url);
            });
          }
          // 包含视频
          if ("video" in item.content) {
            termArr[index].currentMaterialType = "video";
            termArr[index].hasMaterial = true;
            termArr[index].content.video = {
              id: item.content.video.media_id,
              url: item.content.video.video_url,
            };
          }
          // 包含小程序
          if ("applets" in item.content) {
            termArr[index].currentMaterialType = "applet";
            termArr[index].hasMaterial = true;
            termArr[index].content.applet = {
              ...item.content.applets,
            };
          }
          // 包含链接
          if ("link" in item.content) {
            termArr[index].currentMaterialType = "link";
            termArr[index].hasMaterial = true;
            termArr[index].content.link = {
              ...item.content.link,
              cover_url: item.content.link.pic_url,
            };
          }
          if ("card" in item.content) {
            termArr[index].currentMaterialType = "card";
            termArr[index].hasMaterial = true;
            termArr[index].content.card = {
              ...item.content.card,
            };
          }
        });

        this.planInfo.task_start_time = data.task_start_time;
        this.planInfo.task_end_time = data.task_end_time;
        this.planInfo.consumer_file_id = data.consumer_file_id;
        this.planInfo.is_add_tag = data.is_add_tag;
      }
    });
  }
  //添加条件
  private addTerm() {
    this.planInfo.termArr.push({
      term: "5",
      value: "",
      customerLabel: [], //已选标签name
      customerLabelIds: [], // 已选标签ID
      deleteTag: [],
      deleteTagId: [],
      text: {
        content: "", // 任务内容
      },
      content: {
        // 素材内容
        videoId: "0",
        videoUrl: "",
        videoFlag: false,
        appletId: "0",
        applet: {},
        pictureId: [],
        picture: [],
        video: {},
        link: {},
        card: {},
      },
      hasMaterial: false, // 是否有已选素材
      currentMaterialType: "", // 当前选择素材的类型
    });
  }
  // 删除条件
  deleteRiskEvn(index: number) {
    this.$confirm("你确认删除该条件吗？", "删除确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
    })
      .then(() => {
        if (this.status === 1) {
          this.delet_task_conditions.push(this.planInfo.termArr[index].id);
        }
        this.planInfo.termArr.splice(index, 1);
      })
      .catch((err: any) => {
        console.log(err);
      });
  }
  // 选择新增标签的下标
  selectTagsIndex(index: number) {
    this.termIndex = index;
    this.selectedTags = this.planInfo.termArr[index].customerLabelIds;
    this.tagDialog = true;
  }
  // 标签选择回调
  selectedTagsFn(list: any) {
    this.tagDialog = false;
    if (list) {
      this.planInfo.termArr[this.termIndex].customerLabelIds = list.ids;
      this.planInfo.termArr[this.termIndex].customerLabel = list.showList;
      this.planInfo.termArr[this.termIndex].deleteTag = [];
      this.planInfo.termArr[this.termIndex].deleteTagId = [];
    }
  }
  /*
    添加表情
    emj: 选中的表情
    item: 当前任务详情
    index: 当前第几个任务
   */
  private async addEmoji(emj: any, item: any, index: number) {
    const myField = (this.$refs[`textarea${index}`] as any)[0].$refs.textarea;
    // 多行文本内容长度
    if ((this.$refs[`textarea${index}`] as any)[0].value) {
      const textLength = (this.$refs[`textarea${index}`] as any)[0].value.length;
      if (textLength >= 499) return;
    }
    if (myField.selectionStart || myField.selectionStart === 0) {
      const startPos = myField.selectionStart;
      const endPos = myField.selectionEnd;
      (item as any).text.content =
        myField.value.substring(0, startPos) + emj + myField.value.substring(endPos, myField.value.length);

      this.$set(this.planInfo, `termArr.${index}.text.content`, (item as any).text.content);
      await this.$nextTick();
      myField.focus();
      myField.setSelectionRange(endPos + emj.length, endPos + emj.length);
    } else {
      (item as any).text.content += emj;
    }
  }
  // 添加对应素材
  private openSelectMaterial(type: string, index: number) {
    if(type === 'local-image') {
      this.handleUpload('image');
    } else if(type === 'local-video') {
      this.handleUpload('video');
    } else {
      this.termIndex = index;
      this.planInfo.termArr[this.termIndex].currentMaterialType = type;
      this.showMaterialModal = type;
    }
  }
  async handleUpload(type: string) {
    
    try {
      // 上传多个媒体文件
      const results = await (this as any).uploadMedias({
        multiple: type === 'image',
        maxCount: type === 'image' ? 9-this.planInfo.termArr[this.termIndex].content.picture.length : 1,
        maxSize: 10 * 1024 * 1024, // 10MB
        acceptTypes: type === 'image' ? ['image/jpeg','image/png'] : ['video/mp4', 'video/m4v'],
        validate: this.planInfo.task_category !== 0,
        compress: this.planInfo.task_category !== 0,
      });
      console.log('多文件上传结果:', results);
      
      if(type === 'image') {
        this.planInfo.termArr[this.termIndex].currentMaterialType = "picture";
        this.planInfo.termArr[this.termIndex].content.picture = [
          ...this.planInfo.termArr[this.termIndex].content.picture, 
          ...results.map(item => ({
            url: item.url,
            type: 2,
            id: 0
          }))
        ];
      } else {
        this.planInfo.termArr[this.termIndex].currentMaterialType = "video";
        this.planInfo.termArr[this.termIndex].content.video = {
          url:  results[0].url,
          type: 5,
          id: 0,
          duration: results[0].duration,
          video_image_url: results[0].cover_url,
        };
      }
      this.planInfo.termArr[this.termIndex].hasMaterial = true;
      console.log(this.planInfo.termArr[this.termIndex].content);
      

    } catch (error) {
      console.error('上传失败:', error);
      // this.$message.error(error.message);
    }
  }
  // 打开文本弹窗
  opentextDialog(index: number) {
    this.textDialog = true;
    this.termIndex = index;
  }
  // 文本弹窗的关闭
  private async textData(obj: any) {
    this.textDialog = false;
    if (obj) {
      const myField = (this.$refs[`textarea${this.termIndex}`] as any)[0].$refs.textarea;
      // 多行文本内容长度
      const textLength = (this.$refs[`textarea${this.termIndex}`] as any)[0].value.length;
      if (textLength + obj.content.length > 500) {
        this.$message.warning("文本内容不足");
        return;
      }
      if (myField.selectionStart || myField.selectionStart === 0) {
        const startPos = myField.selectionStart;
        const endPos = myField.selectionEnd;
        (this.planInfo.termArr[this.termIndex] as any).text.content =
          myField.value.substring(0, startPos) + obj.content + myField.value.substring(endPos, myField.value.length);

        this.$set(
          this.planInfo,
          `termArr.${this.termIndex}.text.content`,
          (this.planInfo.termArr[this.termIndex] as any).text.content
        );
        await this.$nextTick();
        myField.focus();
        myField.setSelectionRange(endPos + obj.content.length, endPos + obj.content.length);
      } else {
        (this.planInfo.termArr[this.termIndex] as any).text.content += obj.content;
      }
    }
  }
  // 点击素材的叉叉按钮
  private deleteMaterial(type: string, index: number, fIndex?: number) {
    switch (type) {
      case "picture":
        this.planInfo.termArr[index].content.picture.splice(fIndex, 1);
        if (this.planInfo.termArr[index].content.picture.length <= 0) {
          this.planInfo.termArr[index].currentMaterialType = "";
          this.planInfo.termArr[index].hasMaterial = false;
        }
        break;
      case "video":
      case "applet":
      case "link":
      case "card":
        this.planInfo.termArr[index].content[type] = {};
        this.planInfo.termArr[index].currentMaterialType = "";
        this.planInfo.termArr[index].hasMaterial = false;
        break;
    }
  }
  // 选择素材弹窗后回调
  private getSelectMaterial(material?: any) {
    console.log(material);
    
    // 选择图片点击取消按钮回调处理
    if (
      this.planInfo.termArr[this.termIndex].currentMaterialType == "picture" &&
      (!material || material.length == 0) &&
      this.planInfo.termArr[this.termIndex].content.picture.length > 0
    ) {
      this.showMaterialModal = "";
      return;
    }
    // 若素材返回为空，则直接关闭弹窗
    if (!material || Object.keys(material).length == 0 || material.length == 0) {
      this.planInfo.termArr[this.termIndex].currentMaterialType = "";
      this.showMaterialModal = "";
      return;
    }
    switch (this.showMaterialModal) {
      case "picture":
        this.planInfo.termArr[this.termIndex].content.picture.push(...material);
        this.planInfo.termArr[this.termIndex].hasMaterial = true;
        this.showMaterialModal = "";
        break;
      case "video":
      case "applet":
      case "link":
        this.planInfo.termArr[this.termIndex].content[this.showMaterialModal] = material;
        this.planInfo.termArr[this.termIndex].hasMaterial = true;
        this.showMaterialModal = "";
        break;
    }
    console.log(1111, this.planInfo.termArr[this.termIndex].content);
    

  }
  private cancel() {
    this.$confirm("你确认退出该任务吗？", "取消确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
    })
      .then(() => {
        this.$router.go(-1);
      })
      .catch((err: any) => {
        console.log(err);
      });
  }
  private save() {
    //验证生成任务标签form表单
    if (this.planInfo.is_add_tag === 1 && this.planInfo.group_send_type === 0 && this.getDimension('down_url')) {
      const tagForm = this.$refs.tagValueForm as any;
      tagForm.save();
    }
    const temArr = this.planInfo.termArr;
    // 用户画像维度，判断标签是否选择
    // 当为客户画像和不为顾客生日的时候，或者是全部顾客的时候，需要选择标签
    if (this.planInfo.dimension === 2 && this.planInfo.condition_type != 9) {
      for (let i = 0; i < temArr.length; i++) {
        // 定位到当前条件
        if (temArr[i].customerLabel && temArr[i].customerLabel.length <= 0) {
          this.$message.error("请先选择当前的标签");
          (this.$refs as any)["textarea" + i][0].$el.scrollIntoViewIfNeeded();
          return;
        }
      }
    }
    console.log(this.planInfo, "====");
    if ((this.planInfo.dimension === 3 || this.planInfo.dimension === 4) && !this.planInfo.consumer_file_id) {
      this.$message.error("请先上传文件");
      return;
    }
    (this.$refs.cform as any).validate((valid: boolean) => {
      if (!valid) {
        for (let i = 0; i < temArr.length; i++) {
          // 定位到当前任务
          if (!temArr[i].text) {
            (this.$refs as any)["textarea" + i][0].$el.scrollIntoViewIfNeeded();
            break;
          }
        }
      } else {
        const dimensionType = this.planInfo.dimension; // 子任务维度
        // 处理任务内容入参
        const params: any = {
          plan_id: this.plan_id,
          name: this.planInfo.childTaskName,
          // 增加参数
          task_category: Number(this.planInfo.task_category),
          group_send_type: Number(this.planInfo.group_send_type),
          dimension: Number(this.planInfo.dimension),
          condition_type: this.planInfo.condition_type || 1,
          remark: this.planInfo.text,
          plan_task_specify_crowd: this.planInfo.range.map((v: any) => {
            if(v.is_dp) {
              return {
                user_id: '',
                third_id: v.value,
                user_name: v.label,
              };
            }
            return {
              user_id: v.value,
              third_id: 0,
              user_name: v.label,
            };
          }),
          plan_task_conditions: [],
          allow_select: this.planInfo.allow_select
        };
        // 增加参数
        if (this.planInfo.is_add_tag === 1) {
          params.is_add_tag = this.planInfo.is_add_tag;
          params.tag_name = this.planInfo.tag_name;
          params.tag_expire_type = this.planInfo.tag_expire_type;
          params.tag_expire_time = this.planInfo.tag_expire_time;
        }
        if (dimensionType === 2 || dimensionType === 3) {
          params.task_type = this.planInfo.task_type;
          params.send_type = this.planInfo.send_type;
        }
        if (dimensionType === 3 || dimensionType === 4) {
          params.task_type = 1;
          params.consumer_file_id = this.planInfo.consumer_file_id;
        }
        // 任务的生日时间
        if (this.planInfo.condition_type == 8) {
          params.birthday_pre_time = parseInt(this.planInfo.birthday);
        }
        if (this.status === 1) {
          params.task_id = this.task_id;
        }
        if (this.delet_task_conditions && this.delet_task_conditions.length > 0) {
          params.delet_task_conditions = this.delet_task_conditions;
        }
        // 新增参数
        params.finish_rate = Number(this.planInfo.finish_rate)
        // 子任务有效期
        params.task_expire_type = this.planInfo.task_expire_type
        if(this.planInfo.task_expire_type === 1) {
          params.task_start_time = this.planInfo.task_start_time
          params.task_end_time = this.planInfo.task_end_time
        }



        temArr.forEach((item: any) => {
          let conditionsValue: any = [];
          if (dimensionType === 1) {
            conditionsValue = {
              number: Number(item.value),
              expression: item.term === "5" ? "ge" : "le",
            };
          } else {
            const conditions = item.customerLabelIds.map((v: any) => {
              const labelArr: Array<any> = v.split(",");
              const req: any = {
                group_id: Number(labelArr[0]),
                tag_id: Number(labelArr[1]),
                qw_group_id: "",
                qw_tag_id: labelArr[3],
              };
              return req;
            });
            item.customerLabel.forEach((v: any, index: number) => {
              conditionsValue.push({
                tag_name: v,
                tag_id: conditions[index].tag_id,
                group_id: conditions[index].group_id,
                qw_group_id: "",
                qw_tag_id: conditions[index].qw_tag_id,
              });
            });
          }

          let tempObj: any = {};
          switch (item.currentMaterialType) {
            case "card":
              tempObj = {
                card: item.content.card,
              };
              break;
            case "picture":
              tempObj = {
                image: {
                  imgList: [],
                },
              };
              item.content.picture.forEach((imgItem: any) => {
                tempObj.image.imgList.push({
                  pic_url: imgItem.url,
                  media_id: String(imgItem.id),
                });
              });
              break;
            case "video":
              tempObj = {
                video: {
                  media_id: String(item.content.video.id),
                  video_url: item.content.video.url,
                  duration: item.content.video.duration,
                  cover_url: item.content.video.cover_url,
                },
              };
              break;
            case "applet":
              tempObj = {
                applets: {
                  ...item.content.applet,
                },
              };
              break;
            case "link":
              tempObj = {
                link: {
                  ...item.content.link,
                  pic_url: item.content.link.cover_url,
                },
              };
              break;
          }
          let conditionsKey: any = "";
          dimensionType === 1 ? (conditionsKey = "logical") : (conditionsKey = "tag");
          params.plan_task_conditions.push({
            conditions: {
              [conditionsKey]: conditionsValue,
            },
            content: {
              text: item.text,
              ...tempObj,
            },
            id: item.id,
          });
        });
        if (this.status === 0) {
          this.addChildTask(params);
        } else if (this.status === 1) {
          this.editChildTask(params);
        } else if (this.status === 2) {
          this.copyChildTask(params);
        }
      }
    });
  }
  // 新增子任务
  private addChildTask(params: any) {
    Api.addChildTask(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("新增成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 编辑保存子任务
  private editChildTask(params: any) {
    Api.editChildTask(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("修改成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  // 复制保存子任务
  private copyChildTask(params: any) {
    Api.addChildTask(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("复制成功");
        this.$router.go(-1);
      } else {
        this.$message.error(res.message);
      }
    });
  }

  downLoadFile(filename: string, fileUrl: string) {
    downloadFile(fileUrl, {
      filename,
    });
  }
}
</script>
<style lang="scss" scoped>
.link {
  padding: 5px 10px;
  background: #f6f7f8;
  border-radius: 4px;
  cursor: pointer;
}
.material-tip {
  padding: 0 20px 0 30px;
  font-size: 12px;
  color: #999999;

  ul {
    margin: 0;
    padding: 0;
  }
}

.message-info {
  width: 608px;
  font-size: 12px;
  color: red;
  line-height: 20px;
}

.tip {
  font-size: 12px;
  color: red;
  margin-top: 8px;
}

.main-color {
  margin-left: 10px;
  color: #6f3112;
}

.header-condition-title {
  margin-left: 120px;
  padding-left: 20px;
  height: 50px;
  line-height: 50px;
  border: 1px solid #c0c4cc;
}

.risk-content-wrap {
  margin-left: 120px;

  /deep/ .in-item {
    width: auto;
    margin-right: 120px;
  }

  .condition-title {
    height: 50px;
    line-height: 50px;
    padding-left: 20px;
    margin-bottom: 20px;
    background: #e7e7e7;
    border-radius: 4px;
  }
}

.range-choose-btn {
  border-color: $color-main;
  color: $color-main;
}

.margin-btn {
  margin: 20px 0 20px 120px;
}

/deep/ .el-form-item__label {
  width: 120px !important;
}

/deep/ .el-form-item__content {
  margin-left: 120px !important;
}
/deep/ .picker .el-form-item__content {
  margin-left: 0 !important;
}
.picker {
  margin-bottom: 0 !important;
}

.textarea-box {
  width: calc(100% - 10px) !important;

  /deep/ textarea {
    height: 180px !important;
    min-height: 180px !important;
    max-height: 180px !important;
  }
}

.w-120 {
  width: 120px !important;
}

.w-180 {
  width: 180px !important;

  /deep/ .el-input__inner {
    padding-right: 40px !important;
  }
}

.text-position {
  position: relative;
  left: -28px;
}

.risk-content-wrap {
  /deep/ .el-form-item__content {
    width: 100%;
  }
}

.risk-content-box {
  width: 600px;
  border: solid 1px #dcdfe6;
  border-radius: 4px;
  min-height: 425px;

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

.delete-icon {
  float: right;
  line-height: 50px;
  margin-right: 20px;
  cursor: pointer;
}

.title-item {
  width: 100%;
  margin-top: 10px;
  box-sizing: border-box;
  line-height: 20px;
  height: 36px;
  background: #ebebf1;
  padding: 8px 20px;
  justify-content: space-between;
}

.flex {
  display: flex;
}

.flex-align-center {
  align-items: center;
}

.flex-wrap {
  flex-wrap: wrap;
}

.range-content {
  padding: 0 20px 5px;
  justify-content: space-between;
  margin-top: 15px;
  border-bottom: 1px solid #ddd;

  .left {
    display: flex;
    flex-wrap: wrap;
    max-height: 132px;
    overflow: auto;

    .item {
      background-color: #e8e8e9;
      border: #d1d2d2 solid 1px;
      padding: 8px 10px;
      margin: 0 10px 10px 0;

      &.is-delete {
        background: white;
        border: red solid 1px;
        color: red;
      }

      img {
        width: 14px;
        height: 14px;
        margin-right: 4px;
      }
    }

    .item-name {
      font-size: 12px;
      line-height: 12px;
    }
  }

  .right {
    color: #6f3112;
    font-size: 14px;
    min-width: 60px;
    text-align: right;
  }
}

.pointer-cursor {
  cursor: pointer;
}

.footer {
  padding-bottom: 0;

  .range-choose-btn {
    color: #fff;
    background: $color-main;
  }
}

.tips {
  font-size: 12px;
  color: $color-text-second;
  height: 12px;
  line-height: 12px;
  overflow: hidden;
}

.status-desc {
  width: 150%;
  line-height: 10px;
  color: #999999;
  margin-top: 20px;
}</style>
