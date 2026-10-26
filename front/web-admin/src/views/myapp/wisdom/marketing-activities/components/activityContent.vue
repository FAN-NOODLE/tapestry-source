<template>
  <div>
    <div class="active-content get-info">
      <el-radio-group
        class="tab-tit"
        v-model="content_type"
        v-if="activity_type === 2"
      >
        <el-radio-button
          :label="1"
          :disabled="activity_type === 2 && !isActivityContent"
          >活动内容</el-radio-button
        >
        <el-radio-button
          :label="2"
          :disabled="activity_type === 2 && !isMaterial"
          >留资内容</el-radio-button
        >
      </el-radio-group>
      <div class="active-question" v-if="content_type === 1">
        <draggable
          v-if="templateList.length > 0"
          v-model="templateList"
          :options="{
            animation: 150,
            scroll: true,
            scrollSensitivity: 200,
            handle: '.move',
          }"
        >
          <div v-for="(v, i) in templateList" :key="i" class="template-content">
            <div class="template-left">
              <span>【{{ v.name }}】</span>
              <span class="fr">{{ i + 1 }}.</span>
            </div>
            <div class="template-medium">
              <!-- 单选 -->
              <div v-if="v.type == 1">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10 tit">选择列表</div>
                <div class="m-t-10">
                  <div
                    class="choose-item"
                    v-for="(item, index) in v.options"
                    :key="index"
                  >
                    <span class="tit choose-num">选项{{ index + 1 }}</span>
                    <el-input
                      class="del-icon-input"
                      v-model="item.option"
                      size="medium"
                      placeholder="请输入选项"
                      maxlength="64"
                      show-word-limit
                    >
                      <el-tooltip
                        class="item"
                        effect="dark"
                        slot="suffix"
                        content="删除"
                        placement="bottom"
                      >
                        <i
                          v-if="index > 1"
                          class="el-icon-delete-solid el-input__icon"
                          @click="delOption(i, index)"
                        >
                        </i>
                      </el-tooltip>
                    </el-input>
                    <div class="m-t-10">
                      <el-button
                        v-show="item.pic && item.pic.length == 0"
                        type="text"
                        @click="
                          triggerUpload(`image_upload_${i}_${index}`, i, index)
                        "
                        >添加照片</el-button
                      >
                      <upload-image
                        v-show="!(item.pic && item.pic.length == 0)"
                        :ref="`image_upload_${i}_${index}`"
                        :urls="item.pic"
                        img-type="batch"
                        accept="image/jpeg,image/png"
                        :before-upload="beforeUpload"
                        @change="
                          (res) => {
                            return picSuccess(res, i, index);
                          }
                        "
                        :limit="9"
                        multiple
                      ></upload-image>
                    </div>
                  </div>
                  <div class="m-t-10">
                    <el-button size="mini" @click="addChoose(i)"
                      >添加选项</el-button
                    >
                    <el-checkbox
                      v-model="v.isOther"
                      class="m-l-10 m-r-10"
                    ></el-checkbox>
                    <span class="tit">添加'其他'选项</span>
                  </div>
                </div>
              </div>
              <!-- 单行文本 -->
              <div v-if="v.type == 2">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10">
                  最少填写
                  <el-input
                    v-model="v.min"
                    size="mini"
                    class="input-size-small"
                    @input="minInput($event, i)"
                  ></el-input>
                  个字符， 最多填写
                  <el-input
                    v-model="v.max"
                    size="mini"
                    class="input-size-small"
                    @input="maxInput($event, i)"
                  ></el-input>
                  个字符
                </div>
              </div>
              <!-- 多项选择 -->
              <div v-if="v.type == 3">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10 tit">选择列表</div>
                <div class="m-t-10">
                  <div
                    class="choose-item"
                    v-for="(item, index) in v.options"
                    :key="index"
                  >
                    <span class="tit choose-num">选项{{ index + 1 }}</span>
                    <el-input
                      class="del-icon-input"
                      v-model="item.option"
                      size="medium"
                      placeholder="请输入选项"
                      maxlength="64"
                      show-word-limit
                    >
                      <el-tooltip
                        class="item"
                        effect="dark"
                        slot="suffix"
                        content="删除"
                        placement="bottom"
                      >
                        <i
                          v-if="index > 2"
                          class="el-icon-delete-solid el-input__icon"
                          @click="delOption(i, index)"
                        >
                        </i>
                      </el-tooltip>
                    </el-input>
                    <div class="m-t-10">
                      <el-button
                        v-show="item.pic && item.pic.length == 0"
                        type="text"
                        @click="
                          triggerUpload(`image_upload_${i}_${index}`, i, index)
                        "
                        >添加照片</el-button
                      >
                      <upload-image
                        v-show="!(item.pic && item.pic.length == 0)"
                        :ref="`image_upload_${i}_${index}`"
                        :urls="item.pic"
                        img-type="batch"
                        accept="image/jpeg,image/png"
                        :before-upload="beforeUpload"
                        @change="
                          (res) => {
                            return picSuccess(res, i, index);
                          }
                        "
                        :limit="9"
                        multiple
                      ></upload-image>
                    </div>
                  </div>
                  <div class="m-t-10">
                    <el-button size="mini" @click="addChoose(i)"
                      >添加选项</el-button
                    >
                    <el-checkbox
                      v-model="v.isOther"
                      class="m-l-10 m-r-10"
                    ></el-checkbox>
                    <span class="tit">添加'其他'选项</span>
                  </div>
                  <div class="m-t-10">
                    最少选择
                    <el-input
                      v-model="v.min"
                      size="mini"
                      class="input-size-small"
                      @input="limitSelect($event, i, v.options.length, 0)"
                    ></el-input>
                    个选项， 最多选择
                    <el-input
                      v-model="v.max"
                      size="mini"
                      class="input-size-small"
                      @input="limitSelect($event, i, v.options.length, 1)"
                    ></el-input>
                    个选项
                  </div>
                </div>
              </div>
              <!-- 多行文本 -->
              <div v-if="v.type == 4">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10">
                  最少填写
                  <el-input
                    v-model="v.min"
                    size="mini"
                    class="input-size-small"
                    @input="minInput($event, i)"
                  ></el-input>
                  个字符， 最多填写
                  <el-input
                    v-model="v.max"
                    size="mini"
                    class="input-size-small"
                    @input="maxInput($event, i)"
                  ></el-input>
                  个字符
                </div>
              </div>
              <!-- 上传图片 -->
              <div v-if="v.type == 5">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10">
                  最多上传
                  <cy-select
                    width-size="medium"
                    size="small"
                    class="input-size-small"
                    placeholder=""
                    v-model="v.max"
                    :optionsList="numList"
                    :options="{ label: 'val', value: 'val' }"
                  ></cy-select>
                  张图片
                </div>
              </div>
              <!-- 下拉选择 -->
              <div v-if="v.type == 6">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  size="medium"
                  placeholder="请输入说明文字"
                  maxlength="200"
                  show-word-limit
                ></el-input>
                <div class="m-t-10 tit">选择列表</div>
                <div class="m-t-10">
                  <div
                    class="choose-item m-b-10"
                    v-for="(item, index) in v.options"
                    :key="index"
                  >
                    <span class="tit choose-num">选项{{ index + 1 }}</span>
                    <el-input
                      v-model="item.option"
                      size="medium"
                      placeholder="请输入选项"
                      maxlength="64"
                      show-word-limit
                    >
                      <el-tooltip
                        class="item"
                        effect="dark"
                        slot="suffix"
                        content="删除"
                        placement="bottom"
                      >
                        <i
                          v-if="index > 1"
                          class="el-icon-delete-solid el-input__icon"
                          @click="delOption(i, index)"
                        >
                        </i>
                      </el-tooltip>
                    </el-input>
                  </div>
                  <div class="m-t-10">
                    <el-button size="mini" @click="addOption(i)"
                      >添加选项</el-button
                    >
                    <el-checkbox
                      v-model="v.checked"
                      class="m-l-10 m-r-10"
                    ></el-checkbox>
                    <span class="tit">添加'其他'选项</span>
                  </div>
                </div>
              </div>
              <!-- 日期时间 -->
              <div v-if="v.type == 7">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10">
                  <span class="m-r-10">日期格式</span>
                  <el-radio-group v-model="v.radio">
                    <el-radio :label="1">时:分</el-radio>
                    <el-radio :label="2">年-月-日</el-radio>
                    <el-radio :label="3">年-月-日 时:分</el-radio>
                  </el-radio-group>
                </div>
              </div>
              <!-- 姓名 -->
              <div v-if="v.type == 8">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
              <!-- 性别 -->
              <div v-if="v.type == 9">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  size="medium"
                  placeholder="请输入说明文字"
                  maxlength="200"
                  show-word-limit
                ></el-input>
                <div class="m-t-10 tit">选择列表</div>
                <div class="m-t-10">
                  <div
                    class="choose-item m-b-10"
                    v-for="(item, index) in v.options"
                    :key="index"
                  >
                    <span class="tit choose-num">选项{{ index + 1 }}</span>
                    <el-input
                      class="del-icon-input"
                      v-model="item.option"
                      size="medium"
                      placeholder="请输入选项"
                      maxlength="64"
                      show-word-limit
                    >
                      <el-tooltip
                        class="item"
                        effect="dark"
                        slot="suffix"
                        content="删除"
                        placement="bottom"
                      >
                        <i
                          v-if="index == 2"
                          class="el-icon-delete-solid el-input__icon"
                          @click="delOption(i, index)"
                        ></i>
                      </el-tooltip>
                    </el-input>
                  </div>
                </div>
              </div>
              <!-- 手机号 -->
              <div v-if="v.type == 10">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
              <!-- 邮箱 -->
              <div v-if="v.type == 11">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
              <!-- 地址 -->
              <div v-if="v.type == 12">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10">
                  <el-checkbox
                    v-model="v.isAddress"
                    class="m-l-10 m-r-10"
                  ></el-checkbox>
                  <span class="tit">需填写详细地址</span>
                </div>
              </div>
              <!-- 邮编 -->
              <div v-if="v.type == 13">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
            </div>
            <div class="template-right">
              <el-checkbox
                v-model="v.isRequired"
                class="m-l-10 m-r-10"
              ></el-checkbox>
              <span>必填</span>

              <el-tooltip
                class="item"
                effect="dark"
                slot="suffix"
                content="移动"
                placement="bottom"
              >
                <i class="el-icon-rank icon-btn move"></i>
              </el-tooltip>
              <el-tooltip
                class="item"
                effect="dark"
                slot="suffix"
                content="复制"
                placement="bottom"
              >
                <i
                  class="el-icon-document-copy icon-btn"
                  @click="copyItem(v, i)"
                ></i>
              </el-tooltip>
              <el-tooltip
                class="item"
                effect="dark"
                slot="suffix"
                content="删除"
                placement="bottom"
              >
                <i
                  class="el-icon-delete-solid icon-btn"
                  @click="delItem(i)"
                ></i>
              </el-tooltip>
            </div>
          </div>
        </draggable>
        <div v-else class="empty">点击右边控件添加题目</div>
      </div>
      <div class="active-question" v-if="content_type === 2">
        <draggable
          v-if="templateSecondList.length > 0"
          v-model="templateSecondList"
          :options="{
            animation: 150,
            scroll: true,
            scrollSensitivity: 200,
            handle: '.move',
          }"
        >
          <div
            v-for="(v, i) in templateSecondList"
            :key="i"
            class="template-content"
          >
            <div class="template-left">
              <span>【{{ v.name }}】</span>
              <span class="fr">{{ i + 1 }}.</span>
            </div>
            <div class="template-medium">
              <!-- 姓名 -->
              <div v-if="v.type == 8">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
              <!-- 性别 -->
              <div v-if="v.type == 9">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10 tit">选择列表</div>
                <div class="m-t-10">
                  <div
                    class="choose-item m-b-10"
                    v-for="(item, index) in v.options"
                    :key="index"
                  >
                    <span class="tit choose-num">选项{{ index + 1 }}</span>
                    <el-input
                      class="del-icon-input"
                      v-model="item.option"
                      size="medium"
                      placeholder="请输入选项"
                      maxlength="64"
                      show-word-limit
                    >
                      <el-tooltip
                        class="item"
                        effect="dark"
                        slot="suffix"
                        content="删除"
                        placement="bottom"
                      >
                        <i
                          v-if="index == 2"
                          class="el-icon-delete-solid el-input__icon"
                          @click="delOption(i, index)"
                        ></i>
                      </el-tooltip>
                    </el-input>
                  </div>
                </div>
              </div>
              <!-- 手机号 -->
              <div v-if="v.type == 10">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
              <!-- 邮箱 -->
              <div v-if="v.type == 11">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
              <!-- 地址 -->
              <div v-if="v.type == 12">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
                <div class="m-t-10">
                  <el-checkbox
                    v-model="v.isAddress"
                    class="m-l-10 m-r-10"
                  ></el-checkbox>
                  <span class="tit">需填写详细地址</span>
                </div>
              </div>
              <!-- 邮编 -->
              <div v-if="v.type == 13">
                <el-input
                  v-model="v.tit"
                  size="medium"
                  placeholder="请输入标题"
                  maxlength="64"
                  show-word-limit
                ></el-input>
                <el-input
                  v-model="v.mark"
                  class="m-t-10"
                  type="textarea"
                  :rows="6"
                  maxlength="200"
                  resize="none"
                  show-word-limit
                  placeholder="请输入说明文字"
                ></el-input>
              </div>
            </div>
            <div class="template-right">
              <el-checkbox
                v-model="v.isRequired"
                class="m-l-10 m-r-10"
              ></el-checkbox>
              <span>必填</span>

              <el-tooltip
                class="item"
                effect="dark"
                slot="suffix"
                content="移动"
                placement="bottom"
              >
                <i class="el-icon-rank icon-btn move"></i>
              </el-tooltip>
              <el-tooltip
                class="item"
                effect="dark"
                slot="suffix"
                content="复制"
                placement="bottom"
              >
                <i
                  class="el-icon-document-copy icon-btn"
                  @click="copyItem(v, i)"
                ></i>
              </el-tooltip>
              <el-tooltip
                class="item"
                effect="dark"
                slot="suffix"
                content="删除"
                placement="bottom"
              >
                <i
                  class="el-icon-delete-solid icon-btn"
                  @click="delItem(i)"
                ></i>
              </el-tooltip>
            </div>
          </div>
        </draggable>
        <div v-else class="empty">点击右边控件添加题目</div>
      </div>
      <!-- 右侧控件选择 -->
      <div class="active-components">
        <div
          v-for="(item, index) in components"
          :key="index"
          class="active-type"
        >
          <div class="components-name">{{ item.name }}</div>
          <div class="component-content">
            <div
              v-for="(v, i) in item.components"
              :key="i"
              class="component-item"
              @click="addComponents(v)"
            >
              {{ v.name }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import draggable from "vuedraggable";
import UploadImage from "@/components/business/upload-image/index.vue";
import Api from "@/api/wisdom/marketing-activities";
@Component({
  components: {
    draggable,
    UploadImage,
  },
})
export default class extends Vue {
  @Prop() info: any;
  @Prop({ default: () => 0 }) activityId?: number;
  @Prop({ default: () => 0 }) newActivityId?: number;
  @Prop({ default: () => false }) isActivityContent?: boolean;
  @Prop({ default: () => false }) isMaterial?: boolean;
  @Prop({ default: () => 1 }) activity_type?: number;
  @Prop() content: any;
  @Prop() retainCapital: any;
  @Watch("isMaterial")
  isMaterialChange(data: any) {
    if (this.activity_type === 2) {
      if (!this.isActivityContent && data) {
        this.content_type = 2;
      }
    }
  }
  @Watch("content")
  contentChange(data: any) {
    try {
      this.templateList = JSON.parse(data).templateList || [];
    } catch (d) {
      this.templateList = [];
    }
    console.log("this.templateList", this.templateList);
  }
  @Watch("retainCapital")
  retainCapitalChange(data: any) {
    try {
      this.templateSecondList = JSON.parse(data).templateSecondList || [];
    } catch (d) {
      this.templateSecondList = [];
    }
    console.log("this.templateSecondList", this.templateSecondList);
  }
  private statusType = Number(this.$route.query.status_type);
  private components = [
    {
      name: "基础组件",
      components: [
        {
          tit: "单项选择",
          name: "单项选择",
          type: 1,
          isRequired: false,
          isOther: false,
          options: [
            { option: "选项1", pic: [] },
            { option: "选项1", pic: [] },
          ],
          isBase: 1,
        },
        {
          tit: "单行文本",
          name: "单行文本",
          type: 2,
          min: "",
          max: "",
          isRequired: false,
          isBase: 1,
        },
        {
          tit: "多项选择",
          name: "多项选择",
          type: 3,
          options: [
            { option: "选项1", pic: [] },
            { option: "选项2", pic: [] },
            { option: "选项3", pic: [] },
          ],
          min: "",
          max: "",
          isRequired: false,
          isBase: 1,
        },
        {
          tit: "多行文本",
          name: "多行文本",
          type: 4,
          min: "",
          max: "",
          isRequired: false,
          isBase: 1,
        },
        {
          tit: "上传图片",
          name: "上传图片",
          type: 5,
          max: "1",
          isRequired: false,
          isBase: 1,
        },
        {
          tit: "下拉选择",
          name: "下拉选择",
          type: 6,
          options: [
            { option: "选项1", pic: [] },
            { option: "选项2", pic: [] },
          ],
          isRequired: false,
          isBase: 1,
        },
        {
          tit: "日期时间",
          name: "日期时间",
          type: 7,
          radio: 2,
          isRequired: false,
          isBase: 1,
        },
      ],
    },
    {
      name: "联系人控件",
      components: [
        { tit: "姓名", name: "姓名", type: 8 },
        {
          tit: "性别",
          name: "性别",
          type: 9,
          options: [{ option: "男" }, { option: "女" }, { option: "其他" }],
        },
        { tit: "手机号", name: "手机号", type: 10 },
        { tit: "邮箱", name: "邮箱", type: 11 },
        {
          tit: "地址",
          name: "地址",
          type: 12,
          isAddress: false,
        },
        { tit: "邮编", name: "邮编", type: 13 },
      ],
    },
  ];
  private numList = [
    { val: 1 },
    { val: 2 },
    { val: 3 },
    { val: 4 },
    { val: 5 },
    { val: 6 },
    { val: 7 },
    { val: 8 },
    { val: 9 },
  ];
  private content_type = 1;
  private templateList: Array<any> = [];
  private templateSecondList: Array<any> = [];
  // ------------------------------------- methods ----
  // 添加组件
  private addComponents(row: object) {
    const obj = JSON.parse(JSON.stringify(row));
    if (this.activity_type === 2) {
      // 签到
      if (this.isActivityContent && this.isMaterial) {
        // 1.既有内容又有留资产
        if (this.content_type === 1) {
          this.templateList.push(obj);
        } else {
          if (obj.isBase === 1) {
            this.$message.info("留资只允许选择联系人组件");
            return;
          }
          this.templateSecondList.push(obj);
        }
      } else if (this.isActivityContent && !this.isMaterial) {
        // 2. 只有内容
        this.templateList.push(obj);
      } else if (!this.isActivityContent && this.isMaterial) {
        // 3.只有留资
        if (obj.isBase === 1) {
          this.$message.info("留资只允许选择联系人组件");
          return;
        }
        this.content_type = 2;
        this.templateSecondList.push(obj);
      }
    } else {
      this.templateList.push(obj);
    }
  }
  // 复制
  private copyItem(v: any, i: number) {
    if (this.content_type === 1) {
      this.templateList.splice(i + 1, 0, JSON.parse(JSON.stringify(v)));
    } else {
      this.templateSecondList.splice(i + 1, 0, JSON.parse(JSON.stringify(v)));
    }
  }
  // 删除
  private delItem(i: number) {
    if (this.content_type === 1) {
      this.templateList.splice(i, 1);
    } else {
      this.templateSecondList.splice(i, 1);
    }
  }
  // 上传图片
  private triggerUpload(ref: string, i: number, index: number) {
    this.$nextTick(() => {
      (this.$refs as any)[
        ref
      ][0].$children[0].$children[0].$children[0].$children[0].$refs[
        "upload-inner"
      ].handleClick();
    });
  }
  private beforeUpload(file: any) {
    console.log("file", file);
  }
  private picSuccess(res: any, i: number, index: number) {
    this.$set(this.templateList[i].options[index], "pic", res);
  }

  private minInput(val: string, i: number) {
    const inputValue = String(val);
    let str: any = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    if (Number(str) > 10000) {
      str = 10000;
      this.$message.warning("最多为10000个字符");
    }
    this.$nextTick(() => {
      this.templateList[i].min = str;
    });
  }
  private maxInput(val: string, i: number) {
    const inputValue = String(val);
    let str: any = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    if (Number(str) > 10000) {
      str = 10000;
      this.$message.warning("最多为10000个字符");
    }
    this.$nextTick(() => {
      this.templateList[i].max = str;
    });
  }
  private limitSelect(val: string, i: number, maxLength: number, type: number) {
    const inputValue = String(val);
    let str: any = inputValue.replace(/\D/g, "").replace(/^0{1,}/g, "");
    if (Number(str) > maxLength) {
      str = maxLength;
      this.$message.warning("最多不能超过选项的长度");
    }
    this.$nextTick(() => {
      if (type === 0) {
        this.templateList[i].min = str;
      } else {
        this.templateList[i].max = str;
      }
    });
  }
  // 添加选项
  private addChoose(i: number) {
    if (this.content_type === 1) {
      this.templateList[i].options.push({ pic: [], option: "" });
    } else {
      this.templateSecondList[i].options.push({ pic: [], option: "" });
    }
  }
  // 添加下拉选项
  private addOption(i: number) {
    if (this.content_type === 1) {
      this.templateList[i].options.push({});
    } else {
      this.templateSecondList[i].options.push({});
    }
  }
  private delOption(i: number, index: number) {
    if (this.content_type === 1) {
      this.templateList[i].options.splice(index, 1);
    } else {
      this.templateSecondList[i].options.splice(index, 1);
    }
  }

  // 验证自定义的表单内容
  private verify() {
    // 验证活动内容配置
    for (let i = 0; i < this.templateList.length; i++) {
      // 标题验证
      if (this.templateList[i].tit === "") {
        this.$message.warning("活动内容标题不能为空");
        return false;
      }
      // 同时设置最大最小字符集的验证
      if (this.templateList[i].max !== "" && this.templateList[i].min !== "") {
        if (
          Number(this.templateList[i].max) < Number(this.templateList[i].min)
        ) {
          this.$message.warning(
            "活动内容设置字符限制时，最大字符应大于最小字符"
          );
          return false;
        }
      }
      // 1. 单项选择 && 3.多项选择 && 6. 下拉选择 && 9.性别
      if ([1, 3, 6, 9].indexOf(this.templateList[i].type) > -1) {
        for (
          let index = 0;
          index < this.templateList[i].options.length;
          index++
        ) {
          if (this.templateList[i].options[index].option === "") {
            this.$message.warning("活动内容单项选择的每一项都不能为空");
            return false;
          }
        }
      }
    }
    // 验证留资内容
    for (let i = 0; i < this.templateSecondList.length; i++) {
      // 标题验证
      if (this.templateSecondList[i].tit === "") {
        this.$message.warning("留资内容标题不能为空");
        return false;
      }
      if (this.templateSecondList[i].type === 9) {
        for (
          let index = 0;
          index < this.templateSecondList[i].options.length;
          index++
        ) {
          if (this.templateSecondList[i].options[index].option === "") {
            this.$message.warning("留资内容性别选项不能为空");
            return false;
          }
        }
      }
    }
    return true;
  }
  private submit() {
    if (
      this.activity_type !== 2 &&
      this.templateList &&
      this.templateList.length <= 0
    ) {
      this.$message.warning("请配置活动内容");
      return;
    }
    if (
      this.activity_type === 2 &&
      this.isActivityContent &&
      this.templateList &&
      this.templateList.length <= 0
    ) {
      this.$message.warning("请配置活动内容");
      return;
    }
    if (
      this.activity_type === 2 &&
      this.isMaterial &&
      this.templateSecondList &&
      this.templateSecondList.length <= 0
    ) {
      this.$message.warning("请配置留资内容");
      return;
    }
    if (this.verify()) {
      const params: any = {
        base_info: {
          activity_type: this.activity_type,
        },
      };
      const content: any = {};
      const remark: any = {};
      if (this.activity_type !== 2) {
        content.templateList = this.templateList;
      } else {
        if (this.isActivityContent) content.templateList = this.templateList;
        if (this.isMaterial)
          remark.templateSecondList = this.templateSecondList;
      }
      if (this.newActivityId !== 0 && this.newActivityId !== this.activityId) {
        params.id = this.newActivityId;
      } else {
        params.id = this.activityId;
      }
      if (this.statusType === 1 || this.statusType === 3) {
        if (this.activity_type === 2) {
          if (this.isMaterial && this.isActivityContent) {
            // 内容 + 留资
            const params1 = {
              ...params,
              steps: 2,
              content: JSON.stringify(content),
            };
            const params2 = {
              ...params,
              steps: 4,
              remark: JSON.stringify(remark),
            };
            Promise.all([
              Api.activityAddSubmit(params1),
              Api.activityAddSubmit(params2),
            ])
              .then((res: any) => {
                if (res[0].httpCode === 200 && res[1].httpCode === 200) {
                  this.$emit("submitSuccess", res[0].data.id || 0);
                } else {
                  this.$message.warning("保存失败，请稍后重试！");
                }
              })
              .catch((err: any) => {
                console.log("activityUpdata", err);
              });
          } else if (this.isMaterial) {
            // 留资
            const tempParams = {
              ...params,
              steps: 4,
              remark: JSON.stringify(remark),
            };
            Api.activityAddSubmit(tempParams)
              .then((res: any) => {
                if (res.httpCode === 200) {
                  this.$emit("submitSuccess", res.data.id || 0);
                } else {
                  this.$message.warning("保存失败，请稍后重试！");
                }
              })
              .catch((err: any) => {
                console.log("err", err);
              });
          } else if (this.isActivityContent) {
            const tempParams = {
              ...params,
              steps: 2,
              content: JSON.stringify(content),
            };
            Api.activityAddSubmit(tempParams)
              .then((res: any) => {
                if (res.httpCode === 200) {
                  this.$emit("submitSuccess", res.data.id || 0);
                } else {
                  this.$message.warning("保存失败，请稍后重试！");
                }
              })
              .catch((err: any) => {
                console.log("err", err);
              });
          }
        } else {
          const tempParams = {
            ...params,
            steps: 2,
            content: JSON.stringify(content),
          };
          Api.activityAddSubmit(tempParams)
            .then((res: any) => {
              if (res.httpCode === 200) {
                this.$emit("submitSuccess", res.data.id || 0);
              } else {
                this.$message.warning("保存失败，请稍后重试！");
              }
            })
            .catch((err: any) => {
              console.log("err", err);
            });
        }
      } else if (this.statusType === 2) {
        if (this.activity_type === 2) {
          // 添加ID走编辑逻辑接口
          if (this.isMaterial && this.isActivityContent) {
            // 内容 + 留资
            const params1 = {
              ...params,
              content: JSON.stringify({ templateList: this.templateList }),
              steps: 2,
            };
            const params2 = {
              ...params,
              remark: JSON.stringify({
                templateSecondList: this.templateSecondList,
              }),
              steps: 4,
            };
            Promise.all([
              Api.activityUpdata(params1),
              Api.activityUpdata(params2),
            ])
              .then((res: any) => {
                if (res[0].httpCode === 200 && res[1].httpCode === 200) {
                  this.$emit("submitSuccess", 0);
                } else {
                  this.$message.warning("保存失败，请稍后重试！");
                }
              })
              .catch((err: any) => {
                console.log("activityUpdata", err);
              });
          } else if (this.isActivityContent) {
            const tempParams = {
              ...params,
              content: JSON.stringify({ templateList: this.templateList }),
              steps: 2,
            };
            Api.activityUpdata(tempParams)
              .then((res: any) => {
                if (res.httpCode === 200) {
                  this.$emit("submitSuccess", res.data.id || 0);
                } else {
                  this.$message.warning("保存失败，请稍后重试！");
                }
              })
              .catch((err: any) => {
                console.log("activityOneSubmit", err);
              });
          } else if (this.isMaterial) {
            const tempParams = {
              ...params,
              content: JSON.stringify({ templateList: this.templateList }),
              steps: 4,
            };
            Api.activityUpdata(tempParams)
              .then((res: any) => {
                if (res.httpCode === 200) {
                  this.$emit("submitSuccess", res.data.id || 0);
                } else {
                  this.$message.warning("保存失败，请稍后重试！");
                }
              })
              .catch((err: any) => {
                console.log("activityOneSubmit", err);
              });
          }
        } else {
          const tempParams = {
            ...params,
            content: JSON.stringify({ templateList: this.templateList }),
            steps: 2,
          };
          Api.activityUpdata(tempParams)
            .then((res: any) => {
              if (res.httpCode === 200) {
                this.$emit("submitSuccess", res.data.id || 0);
              } else {
                this.$message.warning("保存失败，请稍后重试！");
              }
            })
            .catch((err: any) => {
              console.log("activityOneSubmit", err);
            });
        }
      }
    }
  }
}
</script>
<style lang="scss" scoped>
.active-content {
  display: flex;
  margin-top: 30px;
  border-top: 1px solid #ccc;
  padding: 0 0 30px 30px;
  position: relative;
  .tab-tit {
    position: absolute;
    top: -36px;
  }
  .active-question {
    padding-top: 30px;
    flex: 1;
    border-right: 1px solid #ccc;
    max-height: 450px;
    overflow: auto;
  }
  .active-components {
    width: 300px;
    padding-top: 30px;
    .active-type:nth-child(1) {
      margin-bottom: 30px;
    }
    .components-name {
      padding-left: 15px;
      border-left: 5px solid #2976ea;
    }
    .component-content {
      display: flex;
      width: 100%;
      flex-wrap: wrap;
      margin-top: 20px;
      box-sizing: border-box;
      padding: 0 30px;
      .component-item {
        width: 50%;
        text-align: center;
        padding: 10px;
        border: 1px solid #ccc;
        box-sizing: border-box;
        cursor: pointer;
      }
      .component-item:hover {
        color: #708df5;
        border: 1px solid #708df5;
      }
    }
  }
  .template-content {
    display: flex;
    margin-bottom: 15px;
    .template-left {
      width: 160px;
      color: #888 !important;
      font-size: 16px;
      padding-right: 5px;
      line-height: 36px;
    }
    .template-medium {
      flex: 1;
    }
    .template-right {
      width: 320px;
      font-size: 16px;
      color: #888;
      line-height: 36px;
      .icon-btn {
        margin: 0 10px;
        cursor: pointer;
      }
      .move {
        cursor: move;
      }
    }
  }
  .fr {
    float: right;
  }
  .tit {
    color: #888 !important;
    font-size: 16px;
  }
  .choose-item {
    position: relative;
    .choose-num {
      position: absolute;
      left: -50px;
      top: 6px;
    }
  }
  .input-size-small {
    display: contents;
    /deep/ input {
      width: 70px !important;
      padding-right: 20px !important;
    }
  }
  .m-r-10 {
    margin-right: 10px;
  }
}
.get-info {
  margin-top: 66px;
}
.empty {
  height: 100%;
  line-height: 404px;
  font-size: 20px;
  text-align: center;
  color: #888;
}
/deep/ .el-input__inner {
  padding-right: 50px;
}
.del-icon-input {
  /deep/ .el-input__inner {
    padding-right: 70px;
  }
}
</style>
