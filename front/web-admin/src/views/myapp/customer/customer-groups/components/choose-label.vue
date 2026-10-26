<template>
  <div class="select-materials-dialog">
    <el-dialog
      :visible.sync="selectDialogVisible"
      width="1000px"
      title="添加分组"
      @close="cancel"
    >
      <div class="header">
        <select-input
          ref="tagSelect"
          v-model="search.tagGroupSearch"
          :list="tagList"
          @selectChange="changed"
          :placeholder="search.type == 1 ? '请输入标签组' : '请输入标签'"
        ></select-input>
        <div class="head-foot">
          <el-button size="small" type="primary" @click="groupSearch"
            >查询</el-button
          >
          <el-button size="small" plain @click="reset">重置</el-button>
        </div>
      </div>
      <div class="select-tag">
        <div class="select-tag-left">
          <h4 class="title">标签组</h4>
          <ul class="content-wrapper" v-if="leftMenu.length">
            <li
              v-for="(item, index) in leftMenu"
              :key="item.tag_group_id"
              @click="setTagGroupFn(item, index)"
              :class="{ menuAactive: menuIndex == index + 1 }"
            >
              <span>{{ item.tag_group_name }}</span>
            </li>
          </ul>
          <div class="no-data" v-else>
            <div v-if="search.tagGroupSearch">
              <img src="@/assets/images/empty-search.png" alt="" />
              <p>
                未找到相关标签组
              </p>
            </div>
            <div v-else>
              <img src="@/assets/images/empty-box.png" alt="" />
              <p>暂无标签组</p>
            </div>
          </div>
        </div>
        <div class="line"></div>
        <div class="select-tag-middle">
          <h4 class="title">标签</h4>
          <ul class="content-wrapper" v-if="middleMenu.length">
            <li
              v-for="(item, index) in middleMenu"
              :key="item.tag_id"
              @click="setTagCheckFn(item, index)"
              :class="
                item.disabled
                  ? 'unDisabled'
                  : item.isSelected
                  ? 'menuAactive'
                  : ''
              "
            >
              <span>
                {{ item.tag_name }}
                <i class="el-icon-check group-selected"></i>
              </span>
            </li>
          </ul>
          <div class="no-data" v-else>
            <div v-if="search.tagGroupSearch">
              <img src="@/assets/images/empty-search.png" alt="" />
              <p>
                未找到相关标签
              </p>
            </div>
            <div v-else>
              <img src="@/assets/images/empty-box.png" alt="" />
              <p>暂无标签</p>
            </div>
          </div>
        </div>
        <div class="line"></div>
        <div class="select-tag-right">
          <h4 class="title">已选标签</h4>
          <ul class="content-wrapper">
            <li>
              <div
                class="selected-tag-groups"
                v-for="(item, index) in selectTag"
                :key="item.tag_group_id"
              >
                <div class="content">
                  <h3>{{ item.tag_group_name }}</h3>
                  <div class="tag-content">
                    <div
                      v-for="(subItem, subIndex) in item.tags"
                      :key="subItem.tag_id"
                      class="tag-content-item"
                    >
                      <el-tag
                        class="tag"
                        closable
                        @close="tagCloseFn(subItem)"
                        >{{ subItem.name }}</el-tag
                      >
                      <div
                        class="connect"
                        v-if="subIndex != item.tags.length - 1"
                      >
                        或
                      </div>
                    </div>
                  </div>
                </div>
                <!--边框-->
                <ul class="border" v-if="selectTag.length != 1">
                  <li v-if="index != 0"></li>
                  <li v-if="index != selectTag.length - 1"></li>
                </ul>
                <div
                  class="position-key"
                  v-if="index != selectTag.length - 1 && selectTag.length != 1"
                >
                  且
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>
      <div class="fotter">
        <span v-if="leftMenu.length">客户分组预计：{{ userCountNum }} 人</span>
        <el-button @click="cancel">取 消</el-button>
        <el-button type="primary" @click="submitPar">确 认</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import API from "@/api/customer";
import CustomerModule from "@/store/modules/customer";
import SelectInput from "@/components/business/select-input/index.vue";
@Component({
  components: {
    SelectInput
  }
})
export default class extends Vue {
  @Prop() toDialog!: any;
  private tagList = [
    { key: 1, name: "标签组" },
    { key: 2, name: "标签" }
  ];
  private search = {
    tagGroupSearch: "",
    type: 1
  };
  private selectDialogVisible = true;
  private menuIndex = 1; // 切换标签组名，默认选中第一个
  private selectTag: any = []; //选中的标签
  private forbidTag: any = []; //禁止选的标签
  private userCountArr: any = []; //预计算人数组装数组
  private userCountNum = 0; //预计算人数
  private leftMenu: any = []; //左侧menu 默认列表
  private middleMenu: any = []; //当前选中的标签组-默认选中第menuIndex==1个
  private paramsR: any = {}; //路由参数取值
  private selectTagsNum = 0; //当前选中标签的数量
  private storageSelect: Array<any> = [];
  // private cancelObj: any = {}; //点击取消时，保持原值
  private forbidTagArr = []; //禁止选中的数据，需要传值给分组预计分数的接口
  //------------------------------------------------------data-----------------------------------------------------------/

  //----------------------------------------------------methods-----------------------------------------------------------/
  private changed(val: number) {
    this.search.type = val;
  }
  private groupSearch() {
    this.getLabelListFn();
  }
  private reset() {
    (this.$refs["tagSelect"] as any).reset();
    this.search.tagGroupSearch = "";
    this.search.type = 1;
    this.getLabelListFn();
  }
  private created() {
    this.paramsR = this.toDialog;
    //获取标签组列表
    this.getLabelListFn();
    // type=='include'||'exclude' 包含||不包含
    if (this.paramsR.type == "include") {
      this.selectTag = JSON.parse(JSON.stringify(this.paramsR.include));
      this.forbidTag = JSON.parse(JSON.stringify(this.paramsR.exclude));
    } else if (this.paramsR.type == "exclude") {
      this.selectTag = JSON.parse(JSON.stringify(this.paramsR.exclude));
      this.forbidTag = JSON.parse(JSON.stringify(this.paramsR.include));
    }
    if (this.forbidTag.length != 0) {
      this.forbidTagArr = JSON.parse(JSON.stringify(this.forbidTag));
      this.forbidTagArr.forEach((item: any) => {
        delete item.tag_group_name;
        const temp_arr: any = [];
        item.tags.forEach((subItem: any) => {
          temp_arr.push(subItem.tag_id);
          item.tag_ids = temp_arr;
        });
        delete item.tags;
      });
    }
    this.getUsercountFn();
  }
  //获取标签组、标签列表
  private getLabelListFn() {
    API.labelListApi({
      keyword: this.search.tagGroupSearch,
      search_type: this.search.type
    }).then((res: any) => {
      //合并手动标签和自动标签
      if (res.data.manual != null && res.data.auto != null) {
        this.leftMenu = res.data.auto.concat(res.data.manual);
      } else if (res.data.manual == null && res.data.auto == null) {
        this.leftMenu = [];
      } else {
        this.leftMenu = res.data.manual || res.data.auto;
      }
      this.leftMenu.forEach((item: any) => {
        if (item.tags && item.tags !== null) {
          if (item.tags.length > 0) {
            item.tags.forEach((subItem: any) => {
              subItem.isSelected = false;
              subItem.disabled = false;
            });
          }
        }
      });
      this.menuIndex = 1;
      if (this.leftMenu.length > 0) {
        this.getTagListFn(this.leftMenu[0].tag_group_id);
      } else {
        this.middleMenu = [];
      }
    });
  }
  // 获取标签列表的接口
  private getTagListFn(tag_group_id: number) {
    const params = {
      tag_group_id: tag_group_id,
      order: 1,
      page: 1,
      size: 9999
    };
    API.tagInfoApi(params).then((res: any) => {
      if (res.code === 200) {
        // 获取子标签列表
        this.middleMenu = res.data.list;
        this.forbidTagFn();
        this.tagSelectFn();
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //已经选过的标签置灰，不可选
  private forbidTagFn() {
    if (this.forbidTag.length != 0) {
      this.forbidTag.forEach((item: any) => {
        if (
          this.leftMenu[this.menuIndex - 1].tag_group_id == item.tag_group_id
        ) {
          item.tags.forEach((subItem: any) => {
            this.middleMenu.forEach((curItem: any) => {
              if (subItem.tag_id == curItem.tag_id) {
                curItem.disabled = true;
              }
            });
          });
        }
      });
    }
  }
  //选择标签
  private setTagCheckFn(item: any) {
    //判断当前标签名不可选
    if (item.disabled) {
      return false;
    }
    item.isSelected = !item.isSelected;
    if (item.isSelected) {
      //选中
      //先判断改标签组有没有已经被选过
      let curIndex = 0;
      let bool = true;
      if (this.selectTag.length != 0) {
        for (let i = 0; i < this.selectTag.length; i++) {
          const curItem = this.selectTag[i];
          bool = false;
          if (
            curItem.tag_group_id ==
            this.leftMenu[this.menuIndex - 1].tag_group_id
          ) {
            bool = true;
            curIndex = i;
            break;
          }
        }
        if (bool) {
          //已经有标签组
          this.selectTag[curIndex].tags.push({
            tag_id: item.tag_id,
            name: item.tag_name
          });
        } else {
          //无标签组，新增加
          this.selectTag.push({
            tag_group_id: this.leftMenu[this.menuIndex - 1].tag_group_id,
            tag_group_name: this.leftMenu[this.menuIndex - 1].tag_group_name,
            tags: [
              {
                tag_id: item.tag_id,
                name: item.tag_name
              }
            ]
          });
        }
      } else {
        this.selectTag.push({
          tag_group_id: this.leftMenu[this.menuIndex - 1].tag_group_id,
          tag_group_name: this.leftMenu[this.menuIndex - 1].tag_group_name,
          tags: [
            {
              tag_id: item.tag_id,
              name: item.tag_name
            }
          ]
        });
      }
      this.storageSelect = this.selectTag;
    } else {
      //取消选中
      this.tagCloseFn(item);
    }
    this.getUsercountFn();
  }
  // 组装数据，计算符合当前条件的客户数
  private getUsercountFn() {
    this.userCountArr = [];
    this.userCountArr = JSON.parse(JSON.stringify(this.selectTag));
    this.userCountArr.forEach((item: any) => {
      delete item.tag_group_name;
      const temp_arr: any = [];
      item.tags.forEach((subItem: any) => {
        temp_arr.push(subItem.tag_id);
        item.tag_ids = temp_arr;
      });
      delete item.tags;
    });
    // type=='include'||'exclude' 包含||不包含  1:包含。2：排除
    const params = {
      conditions: [
        {
          tag_groups: this.userCountArr,
          include_type: this.paramsR.type == "include" ? 1 : 2
        }
      ]
    };
    if (this.forbidTagArr.length != 0) {
      params.conditions.push({
        tag_groups: this.forbidTagArr,
        include_type: this.paramsR.type == "include" ? 2 : 1
      });
    }

    API.groupUsercountApi(params).then((res: any) => {
      if (res.code == 200) {
        this.userCountNum = res.data.user_count;
      } else {
        this.$message.error(res.message);
      }
    });
    //计算当前选中的标签数
    this.selectTagsNum = 0;
    this.userCountArr.forEach((item: any) => {
      this.selectTagsNum += item.tag_ids.length;
    });
  }
  //取消选中/tag删除
  private tagCloseFn(item: any) {
    this.selectTag.forEach((curItem: any, cueIndex: number) => {
      curItem.tags.forEach((subCurItem: any, subCurIndex: number) => {
        if (subCurItem.tag_id == item.tag_id) {
          curItem.tags.splice(subCurIndex, 1);
        }
      });
      if (curItem.tags.length == 0) {
        this.selectTag.splice(cueIndex, 1);
      }
    });
    this.middleMenu.forEach((curItem: any) => {
      if (curItem.tag_id == item.tag_id) {
        curItem.isSelected = false;
      }
    });
    this.getUsercountFn();
  }
  //切换选择标签组
  private setTagGroupFn(item: any, index: number) {
    this.menuIndex = index + 1;
    this.getTagListFn(item.tag_group_id);
  }
  //选中/从新建页面带来的tag
  private tagSelectFn() {
    if (this.selectTag.length != 0) {
      this.selectTag.forEach((item: any) => {
        item.tags.forEach((subItem: any) => {
          this.middleMenu.forEach((curItem: any) => {
            if (subItem.tag_id == curItem.tag_id) {
              curItem.isSelected = true;
            }
          });
        });
      });
    }
  }
  //取消-返回上一页
  private cancel() {
    if (this.paramsR.type == "include") {
      CustomerModule.setGroupS(this.toDialog.include);
      CustomerModule.setGroupU(this.toDialog.exclude);
    } else {
      CustomerModule.setGroupS(this.toDialog.include);
      CustomerModule.setGroupU(this.toDialog.exclude);
    }
    this.$emit("closeDialog");
  }
  //保存
  private submitPar() {
    if (this.selectTag.length == 0) {
      this.$message.error("请选择标签");
      return;
    }
    if (this.paramsR.type == "include") {
      console.log(this.selectTag, "1111q");
      CustomerModule.setGroupS(this.selectTag);
    } else {
      CustomerModule.setGroupU(this.selectTag);
    }
    this.$emit("closeDialog");
  }
}
</script>
<style lang="scss" scoped>
.header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  .head-foot {
    margin-left: 20px;
  }
}
.no-data {
  text-align: center;
  margin-top: 100px;
}
ul,
li {
  list-style: none;
  padding: 0;
}
.select-tag {
  display: flex;
  height: 400px;
  .search-div {
    text-align: center;
    margin-top: 100px;
  }
  .line {
    width: 1px;
    height: 100%;
    background: #ebeef5;
    margin: 0 20px;
  }
  h4 {
    font-size: 14px;
    font-weight: 500;
    color: #303133;
    line-height: 20px;
    margin-bottom: 10px;
  }
  .tags {
    .group-selected {
      margin-left: 10px;
      color: $color-text-button;
    }
  }
  // .tag-group {}
  .explain {
    font-size: 14px;
    color: $color-text-second;
    line-height: 20px;
    margin-bottom: 20px;
  }
  .search {
    padding: 20px 20px 5px;
    background: $color-background;
    margin-bottom: 20px;
    ::v-deep.el-button--small {
      font-size: 14px;
    }
  }
  .select-tag {
    display: flex;
    background: $color-background;
    height: calc(100vh - 460px);
    &-left,
    &-middle,
    &-right {
      flex: 1;
      .content-wrapper {
        height: 350px;
        overflow-y: scroll;
        overflow-x: hidden;
      }
    }
    &-left,
    &-middle {
      li {
        margin-inline-start: 0px;
        height: 20px;
        line-height: 20px;
        color: $color-text-main;
        font-size: 14px;
        padding: 10px 0;
        cursor: pointer;
      }
    }
    &-left {
      li.menuAactive {
        // background: #f5f7fa;
        span {
          color: $color-main;
          margin-right: 20px;
        }
      }
      li:hover {
        background: $color-background;
      }
    }
    &-middle {
      li {
        .group-selected {
          display: none;
        }
      }
      li.unDisabled {
        // background: $color-text-normal;
        cursor: not-allowed;
        color: $color-disable;
      }
      li.unDisabled:hover {
        // background: $color-text-normal;
        span {
          color: $color-disable;
        }
        .group-selected {
          display: none;
        }
      }
      li.menuAactive {
        // background: #f5f7fa;
        span {
          color: $color-main;
          margin-right: 20px;
        }
        .group-selected {
          display: inline-block;
          margin-left: 10px;
          font-size: 16px;
          color: $color-text-button;
        }
      }
      li:hover {
        background: $color-background;
      }
    }
    &-right {
      .selected-tag-groups {
        display: flex;
        margin-bottom: 20px;
        position: relative;
        .content {
          flex: 1;
          border-radius: 4px;
          border: 1px solid $color-second-line;
          padding: 0 20px;
          background: none;
          h3 {
            padding: 20px 0;
            line-height: 20px;
            font-size: 14px;
            color: $color-text-main;
            font-weight: bold;
          }
          .tag-content {
            overflow: hidden;
            &-item {
              float: left;
              margin-bottom: 20px;
              display: flex;
              align-items: center;
              .connect {
                width: 14px;
                padding: 0 10px;
              }
            }
          }
        }
        .border {
          flex: 0 0 22px;
          display: flex;
          flex-direction: column;
          li {
            height: 50%;
            width: 14px;
            border-right: 1px solid $color-second-line;
          }
          li:nth-child(1) {
            border-bottom: 1px solid $color-second-line;
          }
          // height: 50px;
        }
        .position-key {
          position: absolute;
          width: 20px;
          height: 26px;
          line-height: 26px;
          text-align: center;
          // background: $color-background;
          font-size: 12px;
          right: -2px;
          bottom: -20px;
          z-index: 1;
        }
      }
      .selected-tag-groups:first-child {
        .border {
          justify-content: flex-end;
          li {
            // background: red;
            border-top: 1px solid $color-second-line;
            border-bottom: none;
          }
        }
      }
    }
  }
}
.fotter {
  text-align: right;
  margin-top: 20px;
  span {
    margin-right: 20px;
  }
}
</style>
