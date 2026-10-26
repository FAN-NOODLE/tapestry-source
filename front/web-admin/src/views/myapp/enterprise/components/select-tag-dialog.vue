<template>
  <el-dialog
    title="选择标签"
    :visible="true"
    @close="closeFn"
    append-to-body
  >
    <div class="select-tag-dialog">
      <div
        class="aside"
        v-if="!searchLoading"
      >
        <ul>
          <li
            v-for="(item, index) in groupList"
            :key="index"
            :class="activeIndex == index ? 'active' : 'normal'"
            @click="getTagsList(index)"
          >
            {{ item.tag_group_name }}
          </li>
        </ul>
      </div>
      <div class="content">
        <div class="btns">
          <div class="left">
            <el-checkbox
              :indeterminate="isIndeterminate"
              v-model="checkAllTags"
              :disabled="checkboxDisabled"
              @change="handleCheckAllChange"
            >全选</el-checkbox>
            <span class="marginLR10"></span>
            <span class="split"></span>
            <el-button
              type="text"
              class="refresh"
              @click="refresh"
            >刷新</el-button>
          </div>
          <div class="right">
            <cy-input
              width-size="small"
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入标签名称"
              clearable
              v-model.trim="searchText"
            ></cy-input>
          </div>
        </div>
        <div class="line"></div>
        <div
          class="list"
          v-if="tagGroupList[currentIndex] && getCurrentList.length > 0"
        >
          <el-checkbox-group
            v-model="selectedTags"
            @change="handleCheckedTagsChange"
          >
            <div
              class="item"
              v-for="(item, index) in getCurrentList"
              :key="index"
            >
              <el-checkbox
                :key="index"
                :label="item.tag_id"
              >{{
                item.name
              }}</el-checkbox>
            </div>
          </el-checkbox-group>
        </div>
        <div
          class="table-empty-label"
          v-else
        >
          <img
            :src="searchLoading
                ? require('@/assets/images/empty-search.png')
                : require('@/assets/images/empty-box.png')
              "
            alt=""
          />
          <div>{{ searchLoading ? "没有找到该标签" : "暂无标签" }}</div>
        </div>
      </div>
    </div>
    <span
      slot="footer"
      class="dialog-footer"
    >
      <span
        class="footer_color m-r-20"
        v-show="this.selectedTags.length > 0"
      >
        <span>已选标签:</span>
        <span>{{ this.selectedTags.length }}</span><span>/10</span>
      </span>

      <el-button
        @click="closeFn"
        size="small"
      >取消</el-button>
      <el-button
        type="primary"
        @click="addTag"
        size="small"
      >确认</el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component, Prop, Watch } from "vue-property-decorator";
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
  private groupList = [];
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
  handleCheckAllChange(val: boolean) {
    let ids: Array<string> = [];
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
      ids = arr2;
    } else {
      this.checkAllMap[this.currentIndex] = val;
      ids = this.tagGroupList[this.currentIndex].list.map((v: any) => v.tag_id);
    }
    if (val) {
      this.selectedTags = Array.from(new Set(this.selectedTags.concat(ids)));
    } else {
      const arr: Array<string> = [];
      this.selectedTags.forEach((v: string) => {
        if (!ids.includes(v)) {
          arr.push(v);
        }
      });
      this.selectedTags = arr;
    }
    this.isIndeterminate = false;
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
  async getTagsList(index: number) {
    this.activeIndex = index;
    const params: any = {};
    if (this.type) {
      params.type = this.type;
    }
    const res: any = await API.labelListApi(params);
    if (res.code == 200) {
      this.groupList = res.data.manual;
      this.groupListLength = res.data.manual.length;
      if (this.updateSelectedTags.length > 0 && this.activeIndex == 0) {
        // 循环调用 getTagInfoList
        const groudIdArr: any = [];
        for (let i = 0; i < this.groupListLength; i++) {
          groudIdArr.push(res.data.manual[i].tag_group_id)
        }
        await this.getNewTagInfoList(groudIdArr)
      }
      // 取出currentIndex数据
      const currentIndexData: any = res.data.manual[index];
      await this.getTagInfoList(currentIndexData);
    } else {
      this.$message.error(res.message);
    }
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
  // todo: 任务标签详情列表
  async getNewTagInfoList(groudIdArr: any) {
    const params: any = {
      tag_group_ids: groudIdArr,
      page: 1,
      size: 9999,
      tag_keywords: this.searchText,
      order: 0
    };
    const res: any = await API.tagInfoApi(params);
    if (res.code == 200) {
      const tagAllList: any = res.data.list
      for(let i = 0; i< tagAllList.length;i++){
        tagAllList[i].tags = res.data.list
        // 每个对象里面取出tag_group_id相同的
        
      }
      this.packageData(tagAllList);
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
    this.getTagsList(this.currentIndex);
  }
  created() {
    this.getTagsList(0);
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
<style lang="scss">
.el-dialog__header {
  font-family: PingFangSC-Semibold, PingFang SC;
  font-weight: 600;
}

.select-tag-dialog {
  .content {
    .list::-webkit-scrollbar {
      width: 6px;
    }

    .list::-webkit-scrollbar-thumb {
      background-color: $color-background;
      border-radius: 2px;
    }
  }

  .aside::-webkit-scrollbar {
    width: 0px;
  }

  .aside::-webkit-scrollbar:hover {
    width: 6px;
  }

  .list::-webkit-scrollbar-thumb {
    background-color: $color-background;
    border-radius: 2px;
  }
}
</style>
<style lang="scss" scoped>
::v-deep.el-button--small {
  padding: 11px 20px;
}

.marginLR10 {
  margin: 0 10px;
}

.select-tag-dialog {
  display: flex;
  height: 400px;

  .aside {
    flex: none;
    width: 200px;
    height: 400px;
    overflow-y: auto;
    background: #f7f7f7;

    .tag-title {
      font-weight: 700;
      height: 40px;
      font-size: 16px;
      // text-align: center;
      line-height: 30px;
      padding-left: 11px;
    }

    ul {
      list-style: none;
      padding: 0;
      margin: 0;

      li {
        padding: 11px;
        padding-left: 20px;
      }

      li:hover {
        background: #ebebeb !important;
        cursor: pointer;
      }

      li.normal {
        background: #f5f6fa;
      }

      li.active {
        color: #000000;
        font-weight: bold;
        position: relative;

        &::after {
          content: "";
          width: 5px;
          height: 14px;
          background: #000000;
          border-radius: 0px 2px 2px 0px;
          position: absolute;
          top: 50%;
          left: 0;
          margin-top: -7px;
        }
      }
    }
  }

  .content {
    flex: auto;
    flex-direction: column;
    padding: 0 20px;
    height: 400px;

    .btns {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .left {
        .newTag {
          margin-right: 10px;
        }

        .split {
          width: 1px;
          height: 100%;
          border-left: 2px solid #99999947;
        }

        .refresh {
          margin-left: 10px;
        }
      }
    }

    .line {
      margin: 15px 0;
    }

    .list {
      height: 344px;
      overflow-y: auto;

      .item {
        display: inline-block;
        margin-bottom: 10px;
        margin-right: 10px;
        padding: 10px 15px 5px 10px;
        border: 1px solid #f4f4f5;
      }
    }

    .table-empty-label {
      display: flex;
      height: 340px;
      flex-direction: column;

      img {
        width: 124px;
        margin: 80px auto 0;
      }

      div {
        text-align: center;
        line-height: 20px;
        margin-top: 14px;
      }
    }
  }
}

.footer_color {
  color: #777;
  display: inline-block;
  margin-right: 10px;
}

/deep/ .cy-input-small {
  width: 220px !important;
}</style>
