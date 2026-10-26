# 列表页面布局组件

## 组件插槽说明

|   插槽名称  |                         说明                           |
|-----------|--------------------------------------------------------|
|   buttons  |         位于整个布局最上方，一般用于放置按钮如跳转新增        |
|   top-tabs  |        位于搜索栏上方的tab切换，一般用于切换整个搜索栏和table     |
|   search  |          插槽用于放置搜索栏          |
|   tabs  |            位于搜索栏下方的tab切换，一般用于切换table         |
|   operate  |         批量操作按钮等部分信息放置于此                  |
|   table  |           列表内容              |
|   footer  |          底部内容放置分页器       |

#### 效果图

<f-image src="https://img-crs.vchangyi.com/组件说明样式资源/管理后台/ts/layout-list.png" width="500px" />

#### 代码示例

```html
<template>
  <div class="supernatant">
    <layout-list>
      <template slot="buttons">
        <el-button type="primary" size="small">新增按钮</el-button>
      </template>
      <template slot="top-tabs">
        <el-tabs v-model="activeName">
          <el-tab-pane label="标题" name="first"></el-tab-pane>
          <el-tab-pane label="标题一" name="second"></el-tab-pane>
          <el-tab-pane label="标题二" name="third"></el-tab-pane>
          <el-tab-pane label="标题三" name="fourth"></el-tab-pane>
        </el-tabs>
      </template>
      <template slot="search">
        <search-item text="用户：">
          <el-radio-group v-model="radio">
            <el-radio-button label="全部"></el-radio-button>
            <el-radio-button label="老用户"></el-radio-button>
            <el-radio-button label="新用户"></el-radio-button>
          </el-radio-group>
        </search-item>
        <search-item text="词条：">
          <cy-input
            width-size="small"
            size="small"
            prefix-icon="el-icon-search"
            placeholder="请输入搜索词条"
            v-model="keyword"
          ></cy-input>
        </search-item>
        <search-item text="页面：">
          <cy-select width-size="small" size="small" v-model="select">
            <cy-option label="全部" value="1"></cy-option>
            <cy-option label="商城" value="2"></cy-option>
            <cy-option label="门店" value="3"></cy-option>
          </cy-select>
        </search-item>
        <search-item text="选择日期：">
          <el-date-picker
            placeholder="选择日期"
            v-model="date"
            size="medium"
          ></el-date-picker>
        </search-item>
      </template>
      <template slot="tabs">
        <el-tabs v-model="activeName">
          <el-tab-pane label="标题" name="first"></el-tab-pane>
          <el-tab-pane label="标题一" name="second"></el-tab-pane>
          <el-tab-pane label="标题二" name="third"></el-tab-pane>
          <el-tab-pane label="标题三" name="fourth"></el-tab-pane>
        </el-tabs>
      </template>
      <template slot="operate">
        <div>
          <span>已选： 5</span>
          <el-button class="m-l-20" size="small">批量操作按钮</el-button>
        </div>
      </template>
      <template slot="table">
        <el-table :data="list">
          <el-table-column prop="date" label="日期" width="250" fixed="left">
          </el-table-column>
          <el-table-column prop="name" label="姓名" min-width="220">
          </el-table-column>
          <el-table-column prop="province" label="省份" min-width="220">
          </el-table-column>
          <el-table-column prop="city" label="市区" min-width="220">
          </el-table-column>
          <el-table-column prop="address" label="地址" min-width="300">
          </el-table-column>
          <el-table-column
            prop="zip"
            label="邮编"
            min-width="120"
            fixed="right"
          >
          </el-table-column>
        </el-table>
      </template>
      <template slot="footer">
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[100, 200, 300, 400]"
          :page-size="100"
          layout="total, sizes, prev, pager, next, jumper"
          :total="400"
        >
        </el-pagination>
      </template>
    </layout-list>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
@Component
export default class extends Vue {
  private keyword = "";
  private activeName = "first";
  private radio = "全部";
  private select = "1";
  private date = "";
  private item = {
    date: "2016-05-03",
    name: "王小虎",
    province: "上海",
    city: "普陀区",
    address: "上海市普陀区金沙江路 1518 弄",
    zip: 200333
  };
  private list: Array<any> = [];
  currentPage = 4;
  handleSizeChange(val: string) {
    console.log(`每页 ${val} 条`);
  }
  handleCurrentChange(val: string) {
    console.log(`当前页: ${val}`);
  }
  mounted() {
    for (let i = 0; i < 12; i++) {
      this.list.push(this.item);
    }
  }
}
</script>

```