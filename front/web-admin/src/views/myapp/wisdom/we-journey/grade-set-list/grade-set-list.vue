<template>
  <layout-form>
    <template slot="content">
      <!-- 标题 -->
      <div class="title">等级规则设置</div>

      <!-- 提示文案 -->
      <div class="tips">备注：等级icon建议上传132*132的png格式图片</div>

      <!-- form表单 -->
      <el-form :model="formData" :rules="rules" ref="ruleForm">
        <!-- 嵌套table -->
        <el-table stripe border :data="formData.list" ref="table" row-key="id">
          <el-table-column label="序号" width="80" align="center">
            <template slot-scope="scope">{{ scope.$index + 1 }}</template>
          </el-table-column>

          <!-- 等级icon -->
          <el-table-column align="center" label="等级icon" width="100">
            <template slot-scope="scope">
              <el-form-item
                :prop="`list.${scope.$index}.level_icon`"
                :rules="rules.level_icon"
              >
                <upload-image
                  class="upload-template"
                  size="small"
                  accept="image/jpeg,image/png"
                  :urls="scope.row.level_icon"
                  @change="url => picSuccess(url, scope.$index)"
                ></upload-image>
              </el-form-item>
            </template>
          </el-table-column>

          <!-- 等级名称 -->
          <el-table-column label="等级名称" align="center">
            <template slot-scope="scope">
              <el-form-item
                :prop="`list.${scope.$index}.level_name`"
                :rules="rules.level_name"
              >
                <el-input
                  v-model="scope.row.level_name"
                  placeholder="请输入等级名称"
                  maxlength="20"
                  show-word-limit
                ></el-input>
              </el-form-item>
            </template>
          </el-table-column>

          <!-- 最大分值-->
          <el-table-column label="最大分值" align="center" width="160">
            <template slot-scope="scope">
              <span v-if="scope.$index === formData.list.length - 1"
                >大于等于{{ scope.row.min_value }}</span
              >
              <el-form-item
                v-else
                :prop="`list.${scope.$index}.max_value`"
                :rules="
                  scope.$index === formData.list.length - 1
                    ? ''
                    : rules.max_value
                "
              >
                <el-input
                  v-model="scope.row.max_value"
                  placeholder="请输入正整数"
                  oninput="value=value.replace(/^(0+)|[^\d]+/g,'')"
                  @change="value => maxValueChange(value, scope.$index)"
                ></el-input>
              </el-form-item>
            </template>
          </el-table-column>

          <!-- 积分范围 -->
          <el-table-column align="center" label="积分范围" width="150">
            <template slot-scope="scope">
              <span
                :class="{
                  'error-formula':
                    scope.$index !== formData.list.length - 1 &&
                    Number(scope.row.max_value) <= Number(scope.row.min_value)
                }"
              >
                <span v-if="scope.row.min_value">
                  {{ scope.row.min_value }}</span
                >
                <span v-if="scope.row.min_value && scope.row.max_value"
                  >≤ X &lt;</span
                >
                <span v-else-if="scope.row.min_value">≤ X</span>
                <span v-else-if="scope.row.max_value">X &lt;</span>
                <span v-if="scope.row.max_value">{{
                  scope.row.max_value
                }}</span>
              </span>
            </template>
          </el-table-column>

          <!-- 操作 -->
          <el-table-column align="center" label="操作" width="150">
            <template slot-scope="scope">
              <div v-if="scope.$index === formData.list.length - 1">
                <el-tooltip
                  v-if="scope.$index"
                  effect="dark"
                  content="设置附加条件"
                  placement="top"
                >
                  <i
                    class="el-icon-s-operation i-btn special-btn"
                    @click="handleExpand(scope.row)"
                  ></i>
                </el-tooltip>
                <div class="tips tips-m-0">最大等级不设上</div>
                <div class="tips tips-m-0">限，由系统自动生成</div>
              </div>
              <span v-else>
                <el-tooltip
                  effect="dark"
                  content="添加下一等级"
                  placement="top"
                >
                  <i
                    class="el-icon-circle-plus-outline i-btn"
                    @click="handleAdd(scope.$index)"
                  ></i>
                </el-tooltip>

                <el-tooltip
                  effect="dark"
                  content="删除"
                  placement="top"
                  v-if="formData.list.length !== 2"
                >
                  <i
                    class="el-icon-delete i-btn"
                    @click="handleDel(scope.$index)"
                  ></i>
                </el-tooltip>

                <el-tooltip
                  v-if="scope.$index"
                  effect="dark"
                  content="设置附加条件"
                  placement="top"
                >
                  <i
                    class="el-icon-s-operation i-btn"
                    @click="handleExpand(scope.row)"
                  ></i>
                </el-tooltip>
              </span>
            </template>
          </el-table-column>

          <el-table-column type="expand" width="1">
            <template slot-scope="props">
              <el-form-item
                required
                label-width="100"
                label="附加条件："
                class="rules-container"
              >
                <additional-dialog
                  :list="props.row.level_rule_attach_arr"
                  :ref="`integral${props.$index}`"
                  @close="closeExpand(props.row)"
                ></additional-dialog>
              </el-form-item>
            </template>
          </el-table-column>
        </el-table>
      </el-form>
    </template>
    <template slot="footer">
      <el-button type="primary" @click="save">保存</el-button>
      <el-button @click="goBack">取消</el-button>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import AdditionalDialog from "./components/additional-dialog.vue";
import { getLevelDetailApi, getLevelSaveApi } from "@/api/we-journey/grade";
import { ApiResult } from "@/types";

/** 生成列表项 */
function generateItem(
  name = "",
  icon = "",
  min_value = "",
  max_value = "",
  id = ""
) {
  return {
    /** 等级名称 */
    level_name: name,
    /** 等级icon */
    level_icon: icon,
    /** 最小分值 */
    min_value,
    /** 最小分值 */
    max_value,
    /** 是否打开附加条件 */
    open: false,
    /** 附加条件 */
    level_rule_attach_arr: [],
    /** 是否删除 */
    is_delete: 0,
    id: id || `__${new Date().getTime()}`
  };
}

/** 列表项类型 */
type ItemType = ReturnType<typeof generateItem>;

/** 生成默认表单数据 */
function generateList() {
  return [
    generateItem(
      "LV1",
      "https://cdn-cos.coach.com.cn/uploads/2023/v3.3/level-1-new.png",
      "",
      "100",
      "__1"
    ),
    generateItem(
      "LV2",
      "https://cdn-cos.coach.com.cn/uploads/2023/v3.3/level-2.png",
      "100",
      "200",
      "__2"
    ),
    generateItem(
      "LV3",
      "https://cdn-cos.coach.com.cn/uploads/2023/v3.3/level-3.png",
      "200",
      "300",
      "__3"
    ),
    generateItem(
      "LV4",
      "https://cdn-cos.coach.com.cn/uploads/2023/v3.3/level-4-new.png",
      "300",
      "400",
      "__4"
    ),
    generateItem(
      "LV5",
      "https://cdn-cos.coach.com.cn/uploads/2023/v3.3/level-5.png",
      "400",
      "",
      "__5"
    )
  ];
}

@Component({
  components: { UploadImage, AdditionalDialog }
})
export default class GradeSetList extends Vue {
  /** 表单数据 */
  private formData = {
    list: generateList()
  };

  /** 缓存附加条件已经展开的列表 */
  private cacheList = [] as ReturnType<typeof generateList>;

  /** 缓存删除的元素, 为了添加删除标识 */
  private deleteList = [] as any;

  /** 表单检验 */
  private rules = {
    level_name: [
      { required: true, message: "请输入等级名称", trigger: "blur" }
    ],
    max_value: [
      { required: true, message: "请输入最大分值", trigger: "change" }
    ],
    level_icon: [
      { required: true, message: "请上传等级icon", trigger: "change" }
    ]
  };

  /** 打开展开行 */
  private handleExpand(row: ItemType) {
    if (row.open) return;
    row.open = true;

    const table = this.$refs.table as HTMLFormElement;
    table.toggleRowExpansion(row);
    this.cacheList.push(row);
  }

  /** 关闭展开行 */
  private closeExpand(row: ItemType) {
    row.open = false;
    const table = this.$refs.table as HTMLFormElement;
    table.toggleRowExpansion(row);

    this.cacheList = this.cacheList.filter(item => item !== row);
  }

  /** 添加 */
  private handleAdd(index: number) {
    const min_value = this.formData.list[index].max_value;
    this.formData.list.splice(index + 1, 0, generateItem("", "", ""));

    // 为新增元素添加min_value
    this.formData.list[index + 1].min_value = min_value;

    if (index + 2 === this.formData.list.length) return;

    /** 计算差值 */
    const nextItem =
      this.formData.list[index + 2] && this.formData.list[index + 2];

    const calcValue = nextItem.max_value
      ? (Number(nextItem.max_value) - Number(nextItem.min_value)) / 2
      : 100;

    const max_value = Number(min_value) + ~~calcValue + "";

    // 赋值最大值。取平均值的一半
    this.formData.list[index + 1].max_value = max_value;
    this.formData.list[index + 2].min_value = max_value;
  }

  /** 删除 */
  private handleDel(index: number) {
    // 当元素删除后，最小值自动递增到下一级， 必须是删除前获取元素
    const min_value = index ? this.formData.list[index - 1].max_value : "";

    // 缓存删除元素
    // 判断是否自定义id __开头
    const reg = /__/;
    const deleteItem = this.formData.list[index];
    // 只有原本存在的元素才需要缓存，新增的元素不需要缓存
    if (!reg.test(deleteItem.id)) {
      this.deleteList.push({
        ...this.formData.list[index],
        is_delete: 1,
        max_value: Number(deleteItem.max_value),
        min_value: Number(deleteItem.min_value) || 0
      });
    }

    // 删除元素
    this.formData.list.splice(index, 1);

    this.maxValueChange(min_value, index - 1);
  }

  /** 最大值变化，往下一级添加最小值 */
  private maxValueChange(value: string, index: number) {
    const formItem = this.formData.list[index + 1];
    if (formItem) {
      formItem.min_value = value;
    }
  }

  /** 图片上传成功 */
  private picSuccess(url: string, index: number) {
    this.$set(this.formData.list[index], "level_icon", url);
  }

  /** 组合校验规则 */
  private composeValidate() {
    // 子组件校验整合
    const childPro = [];
    // 从列表中循环获取子组件实例
    this.formData.list.forEach((item, index) => {
      const promise = new Promise(resolve => {
        // 获取组件实例
        const childDom = this.$refs[`integral${index}`] as HTMLFormElement;

        // 没有值说明当前item 未开启附加条件
        if (!childDom) resolve(true);

        // 获取子组件参数
        const childData = childDom.getData();
        // 子组件校验
        (childDom.$refs["ruleForm"] as HTMLFormElement).validate(
          (valid: boolean) => {
            if (valid) {
              // 添加规则
              item.level_rule_attach_arr = childData;
              resolve(childData);
            }
          }
        );
      });
      childPro.push(promise);
    });
    // 父组件校验整合
    const parentPro = new Promise(resolve => {
      (this.$refs["ruleForm"] as HTMLFormElement).validate((valid: boolean) => {
        if (valid) resolve(valid);
      });
    });
    childPro.push(parentPro);

    return Promise.all(childPro);
  }

  /** 保存 */
  private async save() {
    try {
      // 判断子组件有没有填写
      await this.composeValidate();

      // 判断是否自定义id __开头
      const reg = /__/;

      // 判断积分范围是否正确
      let list = this.formData.list.map((item, index) => ({
        ...item,
        max_value: Number(item.max_value),
        min_value: Number(item.min_value) || 0,
        level_type: index + 1,
        is_delete: 0,
        id: reg.test(item.id) ? 0 : item.id
      }));

      const isTrue = list.every((item, index) => {
        return item.max_value > item.min_value || index === list.length - 1;
      });

      if (!isTrue) return this.$message.warning("积分范围不正确");

      // 添加删除元素
      list = list.concat(this.deleteList as any);

      console.log("删除元素", this.deleteList);

      console.log("校验通过了", list);
      const res: ApiResult = await getLevelSaveApi({ level_rules: list });

      if (res.code !== 200) return this.$message.error(res.message as string);

      this.$message.success("保存成功");

      // 获取列表之前先关闭折叠行
      this.cacheList.forEach(item => {
        this.closeExpand(item);
      });
      this.getList();
    } catch (e) {
      console.error("校验未通过", e);
    }
  }

  /** 取消 */
  private goBack() {
    this.$router.push({
      name: "myapp_wisdom_we-journey"
    });
  }

  /** 获取列表 */
  private async getList() {
    const res: ApiResult = await getLevelDetailApi({});

    if (res.code !== 200) return this.$message.error(res.message as string);

    // 如果没有值，采用默认数据
    if (!(res.data && res.data.length)) return;

    // 有值,采用接口数据
    this.formData.list = res.data;

    this.deleteList = [];

    // 手动打开附加条件
    res.data.forEach((item: ItemType) => {
      // 默认关闭附件条件
      item.open = false;
      // 判断是否有附加条件
      if (item.level_rule_attach_arr && item.level_rule_attach_arr.length) {
        // 下次渲染打开附件条件
        this.$nextTick(() => {
          this.handleExpand(item);
        });
      }
    });
  }

  /** 生命周期 */
  private created() {
    this.getList();
  }
}
</script>
<style lang="scss" scoped>
.rules-container {
  margin-left: 50px;
}
.title {
  margin-bottom: 20px;
  font-weight: bold;
  font-size: 16px;
  text-align: center;
}

.special-btn {
  margin: 0 0 10px 0 !important;
}

.tips {
  font-size: 12px;
  color: $color-text-second;
  height: 12px;
  line-height: 12px;
  margin-bottom: 15px;
}

.tips-m-0 {
  margin-bottom: 0;
}
.upload-template {
  margin: 0 auto;
  width: 56px;
  height: 56px;
}

.i-btn {
  cursor: pointer;
  margin-left: 15px;
}

/deep/.el-form-item {
  margin-bottom: 0 !important;
}

/deep/.el-table__expand-column {
  display: none;
}

.error-formula {
  color: red;
}
</style>
