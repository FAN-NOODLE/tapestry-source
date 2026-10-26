<template>
  <div>
    <!-- form表单 -->
    <el-form :model="formData" :rules="rules" ref="ruleForm">
      <!-- 嵌套table -->
      <el-table stripe border :data="formData.list" ref="table">
        <!-- 奖章ID -->
        <el-table-column
          label="奖章ID"
          prop="id"
          align="center"
          width="180"
          v-if="disabled"
        >
        </el-table-column>
        <el-table-column
          :label="type === 1 ? '财月' : '财季'"
          v-if="type === 1 || type === 2"
          width="100"
          prop="time_name"
        >
        </el-table-column>
        <!-- 奖章名称 -->
        <el-table-column
          label="奖章名称"
          align="center"
          :render-header="renderHeader"
          width="400"
        >
          <template slot-scope="scope">
            <el-form-item
              :prop="`list.${scope.$index}.medal_name`"
              :rules="rules.medal_name"
            >
              <el-input
                :disabled="disabled"
                v-model="scope.row.medal_name"
                placeholder="请输入奖章名称"
                maxlength="20"
                show-word-limit
              ></el-input>
            </el-form-item>
          </template>
        </el-table-column>

        <!-- 奖章icon -->
        <el-table-column
          align="center"
          label="奖章点亮图标"
          width="200"
          :render-header="renderHeader"
        >
          <template slot-scope="scope">
            <el-form-item
              :prop="`list.${scope.$index}.light_icon`"
              :rules="rules.light_icon"
            >
              <upload-image
                v-if="!disabled"
                class="upload-template"
                size="small"
                accept="image/png"
                :urls="scope.row.light_icon"
                @change="url => picSuccess(url, scope.$index, 'light_icon')"
              ></upload-image>

              <img
                v-else
                class="upload-template"
                :src="scope.row.light_icon"
                alt=""
              />
            </el-form-item>
          </template>
        </el-table-column>

        <!-- 奖章icon -->
        <el-table-column
          align="center"
          label="奖章未点亮图标"
          width="200"
          :render-header="renderHeader"
        >
          <template slot-scope="scope">
            <el-form-item
              :prop="`list.${scope.$index}.unlit_icon`"
              :rules="rules.unlit_icon"
            >
              <upload-image
                v-if="!disabled"
                class="upload-template"
                size="small"
                accept="image/png"
                :urls="scope.row.unlit_icon"
                @change="url => picSuccess(url, scope.$index, 'unlit_icon')"
              ></upload-image>

              <img
                v-else
                class="upload-template"
                :src="scope.row.unlit_icon"
                alt=""
              />
            </el-form-item>
          </template>
        </el-table-column>
      </el-table>
    </el-form>

    <!-- 提示文案 -->
    <div class="tips">
      图片格式仅支持*.png；推荐图片尺寸430*430; 图片不得超过10M
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import { CreateElement } from "vue";

/** 生成列表项 */
function generateItem(time_name: string, time_value: number) {
  return {
    /** 奖章名称 */
    medal_name: "",
    /** 奖章点亮图标 */
    light_icon: "",
    /** 奖章未点亮图标 */
    unlit_icon: "",
    /** 季度、月份 */
    time_name,
    /** 下标 */
    time_value
  };
}

@Component({
  components: { UploadImage }
})
export default class MedalStyle extends Vue {
  /** 时间类型 0无,1月,2季度,3年 */
  @Prop() type!: number;

  /** 是否置灰 */
  @Prop({ default: false }) disabled?: boolean;

  /** 详情 */
  @Prop() info?: Array<ReturnType<typeof generateItem>>;

  /** 时间类型变化，更改奖章样式 */
  @Watch("type", { immediate: true })
  changeType() {
    this.generateList();
  }

  @Watch("info", { deep: true, immediate: true })
  changeInfo() {
    if (this.info && this.info.length) {
      // 处理时间
      this.formData.list = this.info;
    }
  }

  /** 表单数据 */
  private formData = {
    list: [] as Array<ReturnType<typeof generateItem>>
  };

  /** 表单检验 */
  private rules = {
    medal_name: [
      { required: true, message: "请输入奖章名称", trigger: "change" }
    ],
    light_icon: [
      { required: true, message: "请上传奖章点亮图标", trigger: "change" }
    ],
    unlit_icon: [
      { required: true, message: "请上传奖章未点亮图标", trigger: "change" }
    ]
  };

  /** 生成默认表单数据 */
  private generateList() {
    const loopMap = {
      0: [""],
      1: [
        "P1",
        "P2",
        "P3",
        "P4",
        "P5",
        "P6",
        "P7",
        "P8",
        "P9",
        "P10",
        "P11",
        "P12"
      ],
      2: ["第1财季", "第2财季", "第3财季", "第4财季"],
      3: [""]
    };

    this.formData.list = loopMap[this.type].map((item: string, index: number) =>
      generateItem(item, index)
    );
  }

  /** 自定义头部内容 - table */
  private renderHeader(h: CreateElement, { column }: HTMLFormElement) {
    console.log(column);
    return h("div", {}, [
      h("span", { class: "required" }, "*"),
      h("span", {}, column.label)
    ]);
  }

  /** 图片上传成功 */
  private picSuccess(url: string, index: number, key: string) {
    this.$set(this.formData.list[index], key, url);
  }

  /** 供父组件获取提交数据 */
  public getData() {
    return this.formData.list;
  }
}
</script>
<style lang="scss" scoped>
/deep/.el-input {
  width: 350px !important;
}

.tips {
  color: #909399;
  font-size: 12px;
}

.upload-template {
  margin: 0 auto;
  width: 55px;
  height: 56px;
}

/deep/.required {
  color: #f67170;
  margin-right: 5px;
  display: inline-block;
}
</style>
