# 表单详情布局组件

## 组件参数说明

|    参数    | 类型          | 说明          | 可选值    | 默认值  |  是否必填   |
|---------- |---------- |----------|------------  |-------- | ------- |
|   type  | String  | 常规布局和左右布局   | 'normal'/'mobile'    | 'nromal'  |  否   |

## 组件插槽说明

> type为normal的插槽

|   插槽名称  |                         说明                           |
|-----------|--------------------------------------------------------|
|   content  |             主内容区域，用于插槽表单/详情等内容             |
|   footer   |    底部区域，用于放置保存取消等操作按钮，当主区域未超出一屏幕会紧跟主区域下方，若超出一屏幕，会定位在浏览器             |

> type为mobile的插槽

|   插槽名称  |                         说明                           |
|-----------|--------------------------------------------------------|
|   left  |             布局左侧区域,一般放置一些移动端预览图             |
|   right  |             表单内容区域,一般放置表单/详情信息等内容             |
|   footer   |    底部区域，用于放置保存取消等操作按钮，当主区域未超出一屏幕会紧跟主区域下方，若超出一屏幕，会定位在浏览器             |

#### 效果图

> type为normal的效果图
<f-image src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/layout-form-1.png" width="500px" />

> type为mobile的效果图
<f-image src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/layout-form-2.png" width="500px" />

#### 代码参考

> type为normal的代码示例

```html
<template>
  <div>
    <layout-form>
      <template slot="content">
        <cy-form ref="cform" :rules="rules" :model="ruleForm">
          <cy-form-item label="测试form-one" prop="name" number required>
            <cy-input
              v-model="ruleForm.name"
              clearable
              width-size="medium"
            ></cy-input>
          </cy-form-item>
          <cy-form-item
            label="收费方式"
            prop="fee"
            desc="这是对兴趣的描述"
            required
            trigger="change"
          >
            <el-switch
              v-model="ruleForm.fee"
              active-text="按月付费"
              inactive-text="按年付费"
            ></el-switch>
          </cy-form-item>
          <cy-form-item
            label="兴趣"
            prop="faverate"
            desc="这是对兴趣的描述"
            required
            trigger="change"
          >
            <el-checkbox-group v-model="ruleForm.faverate">
              <el-checkbox label="美食/餐厅线上活动" name="type"></el-checkbox>
              <el-checkbox label="地推活动" name="type"></el-checkbox>
              <el-checkbox label="线下主题活动" name="type"></el-checkbox>
              <el-checkbox label="单纯品牌曝光" name="type"></el-checkbox>
            </el-checkbox-group>
          </cy-form-item>
          <cy-form-item
            label="兴趣"
            prop="sex"
            desc="性别咯"
            required
            trigger="change"
          >
            <el-radio-group disabled v-model="ruleForm.sex">
              <el-radio label="男" value="1"></el-radio>
              <el-radio label="女" value="2"></el-radio>
            </el-radio-group>
          </cy-form-item>
          <cy-form-item label="上传">
            <upload-image
              class="upload-template"
              :urls="ruleForm.avator"
              disabled
            ></upload-image>
          </cy-form-item>
        </cy-form>
      </template>
      <template slot="footer">
        <el-button type="primary">保存</el-button>
        <el-button>取消</el-button>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
@Component({
  components: {
    UploadImage
  }
})
export default class extends Vue {
  private text = "";
  private ruleForm = {
    name: "",
    phone: "",
    email: "",
    faverate: ["美食/餐厅线上活动"],
    fee: false,
    sex: "",
    idcard: "",
    action: "",
    date: "",
    dateVal: [],
    time: "",
    timeVal: [new Date(2016, 9, 10, 8, 40), new Date(2016, 9, 10, 9, 40)],
    select: "",
    avator: ""
  };
  private options = [
    {
      value: "选项1",
      label: "黄金糕"
    },
    {
      value: "选项2",
      label: "双皮奶"
    },
    {
      value: "选项3",
      label: "蚵仔煎"
    },
    {
      value: "选项4",
      label: "龙须面"
    },
    {
      value: "选项5",
      label: "北京烤鸭"
    }
  ];
  private rules = {
    name: [
      {
        validator: this.validateself,
        trigger: "blur"
      }
    ]
  };
  private validateself(rule: any, value: any, callback: any) {
    value.length > 5 ? callback(new Error("不允许")) : callback();
  }
}
</script>
```

> type为mobile的代码示例

```html
<template>
  <div class="home-page">
    <layout-form type="mobile">
      <template slot="left">
        <div class="mobile-area">
          <img class="phone" src="@/assets/images/ui/bargain.png" alt="" />
        </div>
      </template>
      <template slot="right">
        <span>表单内容</span>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="toast1">保存</el-button>
        <el-button @click="toast2">取消</el-button>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
@Component
export default class extends Vue {
  private toast1() {
    this.$message.success("保存");
  }
  private toast2() {
    this.$message.error("取消");
  }
}
</script>
<style lang="sass" scoped>
.phone {
  width: 375px;
  height: 798px;
}
</style>


```
