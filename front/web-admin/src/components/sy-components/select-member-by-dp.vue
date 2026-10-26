<template>
  <div class="select-wrap">
    <!-- 左侧组织树 -->
    <div class="tree-box">
      <div class="tree-search clearfix">
        <div class="search-ipt pull-left">
          <cy-input
            class="input-with-select"
            v-model.trim="filter_text"
            size="small"
            width-size="small"
            :placeholder="`请输入${list_types[search_type - 1].name}`"
            clearable="clearable"
          >
            <el-select
              v-if="list_types.length > 1"
              v-model="search_type"
              slot="prepend"
              @change="
                filter_text = '';
                is_show_tree = true;
              "
            >
              <el-option
                v-for="item in list_types"
                :label="item.name"
                :value="item.key"
                :key="item.key"
              ></el-option>
            </el-select>
          </cy-input>
        </div>
        <div class="search-btn pull-left m-l-20">
          <el-button size="small" type="primary" @click="changeFilter">
            查询
          </el-button>
          <el-button size="small" @click="reset">重置</el-button>
        </div>
      </div>

      <div class="tree m-t-1">
        <el-tree
          lazy
          v-show="is_show_tree"
          ref="tree"
          empty-text="暂无可执行的部门"
          :props="default_props"
          :load="loadNode"
          node-key="id"
          show-checkbox
          :expand-on-click-node="AllExpandClick"
          :default-expanded-keys="expand_key"
          @check-change="handleCheckChange"
        >
          <div slot-scope="{ node }" class="custom-tree-node">
            <span class="label">{{ node.label }}</span>
          </div>
        </el-tree>

        <!-- 搜索时隐藏树，展示列表 -->
        <div v-show="!is_show_tree" class="filter-tree">
          <div v-if="filter_obj.length">
            <div
              class="item-box"
              v-for="(item, index) in filter_obj"
              :key="index"
              @click="listClick(item)"
            >
              <span class="label">{{ item.name }}</span>
            </div>
          </div>
          <div class="default" v-if="!filter_obj.length">
            <span>没有找到数据</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧选择的数据 -->
    <div class="selected-box">
      <div class="search-ipt m-b-20" v-if="hasSelectedSearch">
        <cy-input
          class="choice-input"
          size="small"
          width-size="small"
          :placeholder="`请输入${search_types[choice_type - 1].name}`"
          v-model.trim="choice_text"
          clearable="clearable"
        >
          <el-select
            v-model="choice_type"
            slot="prepend"
            @change="choice_text = ''"
          >
            <el-option
              v-for="item in search_types"
              :label="item.name"
              :value="item.key"
              :key="item.key"
            ></el-option>
          </el-select>
        </cy-input>
      </div>
      <div class="filter-title">{{ choice_name }}</div>
      <div class="selecte-content">
        <div v-for="(item, index) in select_arr" :key="index" class="item-box">
          <template
            v-if="
              (choice_text && item.name.search(choice_text) > -1) ||
                !choice_text
            "
          >
            <span class="label">{{ item.name }}</span>
            <em
              @click="deleteStaff(item)"
              class="el-icon-circle-close font-14 color-main pull-right p-8"
            ></em>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Emit } from "vue-property-decorator";
import Api from "@/api/wisdom/Personnel";
import GuideApi from "@/api/wisdom/Guide";

@Component
export default class extends Vue {
  @Prop({ default: () => [] }) selectArr?: Array<any>; // 选择的数据
  @Prop({ default: () => [] }) filterArr?: Array<any>; // 需要过滤的数据
  @Prop({ default: true }) AllExpandClick?: boolean; // 是否允许点击父级展开下级
  @Prop({ default: "已选部门" }) choice_name?: string; // 右侧选择标题
  @Prop({ default: false }) is_show_guide?: boolean; // 是否显示导购
  @Prop({ default: true }) hasSelectedSearch?: boolean; // 是否需要已选数据的搜索
  @Prop({ default: "" }) worker_status?: string | number;

  /*** data ***/
  private expand_key: Array<any> = []; // 默认展开的key
  private search_type = 1; // 搜索类型
  private choice_type = 1;
  private default_props = {
    label: "name",
    isLeaf: "leaf",
    disabled: "disabled"
  }; // 配置参数
  private filter_text = ""; // 搜索的字段
  private choice_text = "";
  private is_show_tree = true; // 是否显示tree
  private filter_obj: Array<any> = []; // 筛选出的数据
  private select_arr: Array<any> = []; // 已选择的部门数据

  get list_types() {
    if (this.is_show_guide) {
      return [
        { name: "部门名称", key: 1 },
        { name: "导购名称", key: 2 }
      ];
    } else {
      return [{ name: "部门名称", key: 1 }];
    }
  }

  get search_types() {
    const searchType = [
      { name: `${this.is_show_guide ? "导购" : "部门"}名称`, key: 1 }
    ];
    if (this.is_show_guide) {
      searchType.push({ name: "导购ID", key: 2 });
    }
    // return [
    //   { name: `${this.is_show_guide ? "导购" : "部门"}名称`, key: 1 },
    //   { name: `${this.is_show_guide ? "导购" : "部门"}ID`, key: 2 }
    // ];
    return searchType;
  }

  /*** methods ***/
  // 选择节点
  private handleCheckChange(data: any, checked: boolean) {
    if (this.is_show_guide) {
      if (checked && !data.is_guide) this.expand_key = [data.id];
      if (checked && data.is_guide) {
        const index = this.select_arr.findIndex(v => v.id === data.id);
        if (index < 0) this.select_arr.push(data);
        console.log(111111111111, JSON.parse(JSON.stringify(this.select_arr)));
      } else {
        this.select_arr = this.select_arr.filter(item => data.id !== item.id);
      }
    } else {
      if (checked) this.expand_key = [data.id];
      if (checked) {
        const index = this.select_arr.findIndex(v => v.id === data.id);
        if (index < 0) this.select_arr.push(data);
      } else {
        this.select_arr = this.select_arr.filter(item => data.id !== item.id);
      }
    }
  }
  // 左侧搜索框搜索
  private changeFilter() {
    this.is_show_tree = !this.filter_text ? true : false;
    this.filter_obj = [];
    if (!this.filter_text) return;
    if (this.search_type === 1) {
      this.searchDp(this.filter_text);
    } else {
      this.searchGuide();
    }
  }

  // 搜索框重置
  private reset() {
    this.filter_text = "";
    this.changeFilter();
  }

  // 搜索部门
  private async searchDp(value: string) {
    const params = {
      page: 1,
      size: 500,
      filter_ids: this.filterArr,
      dp_name: value
    };
    const arr: any = await this.getDp(params);
    this.filter_obj = arr;
  }

  // 搜索导购
  private async searchGuide() {
    if (this.search_type !== 2 || !this.filter_text) return;
    const params = {
      worker_name: this.filter_text,
      worker_status: this.worker_status
    };
    GuideApi.guideListByStore(params).then(res => {
      const { list } = res.data;
      list.forEach((item: any) => {
        item.id = item.member_id;
        item.name = `${item.name}(导购ID：${item.mem_userid})`;
        item.is_guide = true;
        item.leaf = true;
      });
      this.filter_obj = list;
    });
  }

  // 搜索列表点击
  private async listClick(res: any) {
    if (this.is_show_guide) {
      if (res.is_guide) {
        if (this.select_arr.findIndex(v => v.id === res.id) === -1) {
          this.select_arr.push(res);
        }
      } else {
        const params = {
          dp_ids: [res.dp_third_id],
          worker_status: this.worker_status
        };
        const {
          data: { list }
        } = await GuideApi.guideListByStore(params);
        list.forEach((item: any) => {
          item.id = item.member_id;
          item.name = `${item.name}(导购ID：${item.mem_userid})`;
          item.is_guide = true;
          item.leaf = true;
          const data = this.select_arr.find(v => v.id === item.id);
          if (!data) this.select_arr.push(item);
        });
      }
    } else {
      const data = this.select_arr.find(v => v.id === res.id);
      if (!data) this.select_arr.push(res);
    }
    this.setChoice();
  }

  // 删除已选数据
  private deleteStaff(data: any) {
    this.select_arr = this.select_arr.filter(v => v.id !== data.id);
    (this.$refs.tree as any).setChecked(data.id, false, true);
  }

  // 关闭模态框的回调
  @Emit("close")
  closeDialog() {
    console.log("close");
  }

  // 懒加载树节点
  private async loadNode(node: any, resolve: Function) {
    if (node.level === 0) {
      const arr: any = await this.getDp({
        is_top: 1,
        filter_ids: this.filterArr
      });
      if (arr.length > 0) this.expand_key.push(arr[0].id);
      resolve(arr);
      this.setChoice();
    }
    if (node.level >= 1) {
      const params = { dp_id: node.data.id, filter_ids: this.filterArr };
      const list: any = await this.getDp(params, node);
      if (list.length > 0 && node.checked) {
        this.expand_key = [...this.expand_key, ...list.map((v: any) => v.id)];
      }
      resolve(list);
      this.setChoice();
      if (this.is_show_guide) this.getGuide(node.data.dp_third_id, node);
    }
  }

  // 获取部门
  private getDp(params: any, node?: any) {
    return new Promise(resolve => {
      Api.depListApi(params).then(res => {
        const data = res.data.list;
        data.forEach((item: any) => {
          item.id = item.dp_id;
          item.name = item.dp_name;
          if (!this.is_show_guide && !item.is_child_department) {
            item.leaf = true;
          }
          if (!this.is_show_guide && node && node.checked) {
            const data = this.select_arr.find(v => v.id === item.id);
            if (!data) this.select_arr.push(item);
          }
          if (node && node.checked) {
            (this.$refs.tree as any).setChecked(item.id, true, true);
          }
          if (node && node.checked && this.is_show_guide) {
            this.expand_key.push(item.id);
          }
        });
        resolve(data);
      });
    });
  }

  // 点击部门获取部门下的导购
  private getGuide(id: string, node: any) {
    const params = { dp_ids: [id], worker_status: this.worker_status };
    GuideApi.guideListByStore(params).then(res => {
      const { list } = res.data;
      if (list.length <= 0) return;
      list.forEach((item: any) => {
        item.id = item.member_id;
        item.name = `${item.name}(导购ID：${item.mem_userid})`;
        item.is_guide = true;
        item.leaf = true;
        if (node.checked) {
          const data = this.select_arr.find(v => v.id === item.id);
          if (!data) this.select_arr.push(item);
        }
        (this.$refs.tree as any).append(item, node);
      });
      this.setChoice();
    });
  }

  // 设置树选中
  private setChoice() {
    // let key_arr = this.select_arr.map(v => v.id);
    // this.$refs.tree.setCheckedKeys(key_arr);
    this.select_arr.forEach(item => {
      (this.$refs.tree as any).setChecked(item.id, true, true);
    });
  }

  private mounted() {
    if (this.selectArr?.length) {
      this.select_arr = this.selectArr;
      if (!this.is_show_guide) {
        this.select_arr.forEach(item => {
          item.id = item.id || item.dp_id;
          item.name = item.name || item.dp_name;
          if (!this.is_show_guide && !item.is_child_department) {
            item.leaf = true;
          }
        });
      } else {
        this.select_arr.forEach(item => {
          item.id = item.member_id;
          item.name = `${item.name}(导购ID：${item.mem_userid})`;
          item.is_guide = true;
          item.leaf = true;
        });
      }
    }
  }
}
</script>

<style lang="scss" scoped>
/deep/.el-select .el-input {
  width: 110px;
}

/deep/.el-input-group__append,
/deep/.el-input-group__prepend {
  background-color: $color-white;
}

.select-wrap {
  flex: 1;
  // width: 795px;
  height: 420px;
  display: flex;
  justify-content: space-between;

  .tree-box {
    // width: 480px;
    width: 62%;
    height: 420px;
    box-sizing: border-box;
    padding-right: 5px;
    border-right: 2px solid #e4e7ed;

    .tree-search {
      margin-left: 20px;
    }
  }

  .selected-box {
    // width: 280px;
    width: 35.5%;
    height: 100%;
    box-sizing: border-box;
  }
}

.choice-input {
  width: 148px;
}

.item-box {
  display: flex;
  justify-content: space-between;
  align-items: center;
  line-height: 26px;
  cursor: pointer;
}

.label {
  max-width: 380px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: inline-block;
}

.filter-title {
  margin-bottom: 10px;
  width: 100%;
  color: #787878;
}

.filter-tree {
  height: calc(100% - 20px);
  margin-left: 20px;
  margin-top: 20px;
}

.default {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.delete {
  cursor: pointer;
}

.tree,
.selecte-content {
  height: 380px;
  overflow: auto;

  &::-webkit-scrollbar {
    /*滚动条整体样式*/
    width: 5px;
    /*高宽分别对应横竖滚动条的尺寸*/
    height: 1px;
  }

  &::-webkit-scrollbar-thumb {
    /*滚动条里面小方块*/
    border-radius: 2px;
    // box-shadow: inset 0 0 5px rgba(0, 0, 0, .2)
    background: #e4e7ed;
  }

  &::-webkit-scrollbar-track {
    /*滚动条里面轨道*/
    // box-shadow: inset 0 0 5px rgba(0, 0, 0, .2)
    border-radius: 2px;
    background: white;
  }
}

.selecte-content {
  height: 334px;
}

.p-8 {
  padding: 8px;
}

.m-t-1 {
  margin-top: 1px;
}

.custom-tree-node {
  height: 20px;
  line-height: 20px;
  margin-left: 6px;
}
.el-select-dropdown__item.selected {
  font-weight: normal;
}
</style>
