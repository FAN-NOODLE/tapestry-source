<template>
  <el-dialog
    title="选择标签"
    :visible.sync="visible"
    append-to-body
  >
    <el-row>
      <el-input v-model="keyword" style="width: 400px;" placeholder="搜索标签名称" prefix-icon="el-icon-search"></el-input>
    </el-row>
    <div class="m-t-20" style="max-height: 60vh;overflow-y: scroll;">
      <el-collapse v-model="active">
        <el-collapse-item v-for="item,index in showGroupList" :key="index" :title="item.name" :name="index">
          <template slot="title">
            <div class="title">
              <div>{{ item.name }}</div>
              <div class="check-all" @click.stop="handleCheckAllChange(item)">
                <el-checkbox v-model="item.checked">全选</el-checkbox>
                <div class="check-all-mask"></div>
              </div>
            </div>
          </template>
          <div v-if="item.children" class="list">
            <el-tag v-for="tag in item.children" :key="tag.id" :effect="isChecked(tag) ? 'dark': 'plain'" @click="handleSelect(tag, item)">{{ tag.name }}</el-tag>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
    <div
      slot="footer"
      class="dialog-footer"
    >
      <span
        class="footer_color m-r-20"
        v-show="selection.length > 0"
      >
        <span>已选标签:</span>
        <span>{{ selection.length }}</span><span v-if="limit>0">/{{limit}}</span>
      </span>

      <el-button
        @click="handleCancel"
        size="small"
      >取消</el-button>
      <el-button
        type="primary"
        @click="handleConfirm()"
        size="small"
      >确认</el-button>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch, PropSync } from "vue-property-decorator";
import API from "@/api/customer/index";
@Component({})
export default class extends Vue {
  @Prop({ default: () => [] }) updateSelectedTags?: any;
  @Prop({ default: () => null }) type?: any;
  /************data**************/
  private tagGroupList: Array<any> = [];
  private currentIndex = 0;
  private activeIndex = 0;
  private selectedTags: Array<any> = [];
  private checkAllTags = false;
  private isIndeterminate = false;
  private checkAllMap = {};
  private searchText = "";
  private checkboxDisabled = false;
  private searchLoading = false;
  private idMapData = {};
  private groupListLength = 0;
  
  @Watch("searchText")
  changeSearchLoadingStatus(value: string) {
    if (value) {
      this.searchLoading = true;
    } else {
      this.searchLoading = false;
    }
  }

  // ====================================================================================================================

  @PropSync('show', { type: Boolean, default: false }) visible: boolean
  @Prop({ default: () => [] }) value?: any[];
  @Prop({ default: () => 10 }) limit?: number;

  // 搜索
  keyword = "";

  // 展开项
  active:number[] = [];
  // 标签组
  groupList: any[] = [];
  // 选中项
  selection: any[] = [];

  /** 
   * 查询
   */
  get showGroupList() {
    // 同过keyword过滤 name
    const list = JSON.parse(JSON.stringify(this.groupList))
    return list.map(item => {
      // 过滤子元素
      if(item.children) {
        item.children = item.children.filter(child => {
          return !this.keyword || child.name.indexOf(this.keyword) !== -1
        })
      } else {
        item.children = []
      }
      return item
    }).filter(item => {
      // 过滤父元素
      return !this.keyword || item.name.indexOf(this.keyword) !== -1 || item.children.length > 0
    });
  }

  get isChecked() {
    return item => {
      const ids = this.selection.map(v => v.id)
      return ids.includes(item.id);
    }
  }

  @Watch('show')
  async handleShow(val: boolean) { 
    if(val) {
      this.selection = this.value
      this.keyword = ''
    }
  }



  /**
   * 获取标签分组列表
   * @param index 
   */
  async getTagsList() {
    const res: any = await API.labelListApi({
      searchText: this.searchText
    });
    if (res.code == 200) {
      console.log(res);
      this.groupList = res.data.manual.map((item: any) => {
        return {
          name: item.tag_group_name,
          id: item.tag_group_id,
          ...item
        }
      })
      this.active = this.groupList.map((_, i) => i)
      this.getNewTagInfoList()
    } else {
      this.$message.error(res.message);
    }
  }

  /**
   * 获取标签列表
   */
  async getNewTagInfoList() {
    const groudIds = this.groupList.map(item => item.id);
    const params: any = {
      tag_group_ids: groudIds,
      page: 1,
      size: 9999,
      tag_keywords: this.searchText,
      order: 0,
      ignore_tag_count: true
    };
    const res: any = await API.tagInfoApi(params);
    if (res.code == 200) {
      
      this.groupList.forEach((item: any) => {
        this.$set(item, "children", res.data.list.filter(v => v.tag_group_id === item.id).map(v => {
          return {
            ...v,
            id: v.tag_id,
            name: v.tag_name,
          }
        }))
      })
    } else {
      this.$message.error(res.message);
    }
  }


  /**
   * 选择
   */
  handleSelect(item, parent) {
    const ids = this.selection.map(v => v.id)
    
    if (ids.includes(item.id)) {
      this.selection = this.selection.filter(v => v.id !== item.id);
    } else {
      if(this.limit > 0 && this.selection.length >= this.limit) {
        this.$message.error(`最多只能选择 ${this.limit} 条数据`);
      } else {
        this.selection.push(item);
      }
    }
    this.handleSelectChange(parent)
  }

  handleSelectChange(item) {

    const ids = this.selection.map(v => v.id)
    const isCheckedAll = item.children.every(v => ids.includes(v.id))
    
    this.$set(item, 'checked', isCheckedAll)
  }


  /**
   * 全选
   * @param e 
   */
  handleCheckAllChange(item) {
    if(!item.checked) {
      // 全选
      // 判断数据是否超出已选数据是否 超出限制
      const ids = this.selection.map(v => v.id);
      const selection = [
        ...this.selection,
        ...item.children.filter(v => !ids.includes(v.id))
      ]
      if(this.limit > 0 && selection.length > this.limit) {
        this.$message.error(`最多只能选择 ${this.limit} 条数据`);

      } else {
        this.selection = selection;
      }
    } else {
      // 取消全选
      const ids = item.children.map(v => v.id);
      this.selection = this.selection.filter(v => !ids.includes(v.id))
    }

    this.handleSelectChange(item)
  }
  
  /**
   * 取消
   */
  handleCancel() {
    this.$emit('cancel')
  }

  /**
   * 确定
   */
  handleConfirm() {
    this.$emit('confirm', this.selection)
  }

  // ====================================================================================================================

  /************methods**************/
  //确认

  public get getCurrentList(): any {
    if (this.searchText) {
      const arr: any[] = this.flatten(this.tagGroupList);
      const arr1: any[] = [];
      const arr2: any[] = [];
      arr.forEach((item: any) => {
        arr1.push(item.list);
      });
      let tempArr: any[] = this.flatten(arr1);
      tempArr = tempArr.filter(
        (v: any) => v.name.indexOf(this.searchText) > -1
      );
      tempArr.forEach((item: any) => {
        arr2.push(item.tag_id);
      });
      this.checkAllTags = arr2.every((val: any) =>
        this.selectedTags.includes(val)
      );
      return tempArr;
    } else {
      const arr1: any[] = [];
      this.tagGroupList[this.currentIndex].list.forEach((item: any) => {
        arr1.push(item.tag_id);
      });
      this.checkAllTags = arr1.every((val: any) =>
        this.selectedTags.includes(val)
      );
      return this.tagGroupList[this.currentIndex].list;
    }
  }

  // 数组扁平化
  flatten(arr: any[]) {
    let res: any[] = [];
    arr.forEach((item: any) => {
      if (Array.isArray(item)) {
        res = res.concat(this.flatten(item));
      } else {
        res.push(item);
      }
    });
    return res;
  }
  addTag() {
    const tags: any = this.packageSelectedTags();
    if (tags.ids.length > 10) {
      this.$message.warning("最多只能选择10个标签");
    } else {
      this.$emit("selectedTags", tags);
    }
  }
  //关闭
  closeFn() {
    this.$emit("selectedTags");
  }
  packageSelectedTags() {
    const list: any = { ids: [], showList: [] };
    list.ids = this.selectedTags.map((v: any) => {
      list.showList.push(this.idMapData[v]?.name);
      return this.idMapData[v].id;
    });
    return list;
  }
  setStatus() {
    const curList: Array<string> = this.tagGroupList[
      this.currentIndex
    ].list.map((v: any) => v.tag_id);
    this.checkboxDisabled = curList.length > 0 ? false : true;
    let checkAll: any = false;
    if (curList.length > 0 && this.selectedTags.length >= curList.length) {
      let status: any = true;
      curList.forEach((v: string) => {
        if (!this.selectedTags.includes(v)) {
          status = false;
        }
      });
      checkAll = status;
    }
    this.checkAllTags = checkAll;
    this.checkAllMap[this.currentIndex] = checkAll;
    this.isIndeterminate = !checkAll;
  }
  changeIndeterminateStatus() {
    const curList: Array<string> = this.tagGroupList[
      this.currentIndex
    ].list.map((v: any) => v.tag_id);
    let status: any = false;
    this.selectedTags.forEach((v: string) => {
      if (curList.includes(v)) {
        status = true;
      }
    });
    this.isIndeterminate = status;
    if (this.checkAllTags) {
      this.isIndeterminate = false;
    }
  }
  handleCheckedTagsChange(value: any) {
    this.setStatus();
    if (value.length == 0) {
      this.isIndeterminate = false;
    }
  }
  updateContent() {
    this.checkAllTags = this.checkAllMap[this.currentIndex] || false;
    this.setStatus();
    this.changeIndeterminateStatus();
  }
  packageData(data: any) {
    let taglist: Array<any> = [];
    taglist = data.map((v: any) => {
      const ele: object = {
        tag_group_id: v.tag_group_id,
        tag_group_name: v.tag_group_name,
        list: v.tags.map((g: any) => {
          const key: string = v.tag_group_id + "," + g.tag_id;
          const value: any = {
            id: v.tag_group_id + "," + g.tag_id + "," + v.work_tag_id + "," + g.work_tag_id,
            name: g.tag_name,
          };
          this.idMapData[key] = value;
          return {
            name: g.tag_name,
            tag_id: v.tag_group_id + "," + g.tag_id,
          };
        }),
      };
      return ele;
    });
    return taglist;
  }
  // todo: 任务标签详情列表
  async getTagInfoList(currentIndexData: any) {
    const params: any = {
      tag_group_id: currentIndexData.tag_group_id,
      page: 1,
      size: 9999,
      tag_keywords: this.searchText,
      order: 0
    };
    const res: any = await API.tagInfoApi(params);
    if (res.code == 200) {
      // 将获取来数据放入currentIndexData
      currentIndexData.tags = res.data.list;
      // 将currentIndexData变成数组
      const currentIndexDataArr: any = [currentIndexData];
      const dealData: any = this.packageData(currentIndexDataArr);
      this.tagGroupList = [];
      this.tagGroupList = dealData;
      this.currentIndex = 0;
      this.$nextTick(() => {
        try {
          this.updateContent();
        } catch (error) {
          console.log(error);
        }
      });
    } else {
      this.$message.error(res.message);
    }
  }

  refresh() {
    this.getTagsList();
  }
  created() {
    this.getTagsList();
    //this.isIndeterminate = true;
    if (this.updateSelectedTags.length > 0) {
      this.selectedTags = this.updateSelectedTags.map((v: any) => {
        const arr: Array<string> = v.split(",").slice(0, 2);
        return arr.join(",");
      });
    }
  }
}
</script>
<style lang="scss" scoped>
.list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;

  .el-tag {
    cursor: pointer;
  }
}

.title {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 20px;
}

.check-all {
  position: relative;

  .check-all-mask {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 1;
  }
}
</style>
