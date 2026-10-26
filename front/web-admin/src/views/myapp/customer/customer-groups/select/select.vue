<template>
  <div>
    <layout-form>
      <template slot="content">
        <p class="explain">
          说明：选择标签后，标签组内的标签为"或"关系，标签组互相为"且"关系，例如：当用户选择用户来源中的“商城用户”和“公众号粉丝”2个标签时，表示客户来自于商城用户或者公众号粉丝，当用
          户再选择访问活跃度中的“活跃访客”标签时，则表示客户来自于商城用户或者公众号粉丝并且是活跃访客
        </p>
        <div class="search">
          <search-item text="筛选：">
            <cy-input
              width-size="small"
              size="small"
              prefix-icon="el-icon-search"
              placeholder="请输入标签或标签组名称"
              v-model="filterKeyword"
            ></cy-input>
          </search-item>
          <el-button type="primary" size="small" @click="getLabelListFn"
            >查询</el-button
          >
          <!-- 增加重置按钮 -->
          <el-button size="small" @click="resetFn">重置</el-button>
        </div>
        <div class="select-tag" v-if="leftMenu.length">
          <div class="select-tag-left">
            <div class="title">标签组名称</div>
            <dl class="content-wrapper">
              <dd
                v-for="(item, index) in leftMenu"
                :key="item.tag_group_id"
                @click="setTagGroupFn(item, index)"
                :class="{ menuAactive: menuIndex == index + 1 }"
              >
                <span>{{ item.tag_group_name }}</span>
              </dd>
            </dl>
          </div>
          <div class="select-tag-middle">
            <div class="title">标签</div>
            <dl class="content-wrapper">
              <dd
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
                  <svg-icon
                    class="group-selected"
                    icon-class="group-selected"
                  ></svg-icon>
                  {{ item.name }}
                </span>
              </dd>
            </dl>
          </div>
          <div class="select-tag-right">
            <div class="title">已选标签（{{ selectTagsNum }}个）</div>
            <dl class="content-wrapper">
              <dd>
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
                    v-if="
                      index != selectTag.length - 1 && selectTag.length != 1
                    "
                  >
                    且
                  </div>
                </div>
              </dd>
            </dl>
          </div>
        </div>
        <div class="no-data" v-else>
          <div v-if="filterKeyword">
            <img src="@/assets/images/empty-search.png" alt="" />
            <p>
              未搜索到标签或标签组
            </p>
          </div>
          <div v-else>
            <img src="@/assets/images/empty-box.png" alt="" />
            <p>暂无标签或标签组</p>
          </div>
        </div>
        <p class="current-total" v-if="leftMenu.length">
          客户分组预计：{{ userCountNum }} 人
        </p>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="saveTagFn">保存</el-button>
        <el-button @click="goBackFn">取消</el-button>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import API from "@/api/customer";
import CustomerModule from "@/store/modules/customer";
@Component
export default class extends Vue {
  /*************************** data start ***************************/
  private filterKeyword = ""; // 标签名或标签组名称筛选
  private menuIndex = 1; // 切换标签组名，默认选中第一个
  private selectTag: any = []; //选中的标签
  private forbidTag: any = []; //禁止选的标签
  private userCountArr: any = []; //预计算人数组装数组
  private userCountNum = 0; //预计算人数
  private leftMenu: any = []; //左侧menu 默认列表
  private middleMenu: any = []; //当前选中的标签组-默认选中第menuIndex==1个
  private paramsR: any = {}; //路由参数取值
  private selectTagsNum = 0; //当前选中标签的数量
  private cancelObj: any = {}; //点击取消时，保持原值
  private forbidTagArr = []; //禁止选中的数据，需要传值给分组预计分数的接口

  /*************************** method start ***************************/
  //获取标签组、标签列表
  private getLabelListFn() {
    API.labelListApi({ keyword: this.filterKeyword, search_type: 3 }).then(
      res => {
        //合并手动标签和自动标签
        if (res.data.manual != null) {
          this.leftMenu = res.data.auto.concat(res.data.manual);
        } else {
          this.leftMenu = res.data.auto;
        }

        this.leftMenu.forEach((item: any) => {
          if (item.tags.length != 0) {
            item.tags.forEach((subItem: any) => {
              subItem.isSelected = false;
              subItem.disabled = false;
            });
          }
        });
        this.menuIndex = 1;
        if (this.leftMenu.length != 0) {
          this.middleMenu = this.leftMenu[this.menuIndex - 1].tags; //当前选中的标签组-默认选中第menuIndex==1个
          this.forbidTagFn();
          this.tagSelectFn();
        }
      }
    );
  }
  //重置
  private resetFn() {
    this.filterKeyword = "";
    this.getLabelListFn();
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

  //切换选择标签组
  private setTagGroupFn(item: object, index: number) {
    this.menuIndex = index + 1;
    this.middleMenu = this.leftMenu[this.menuIndex - 1].tags;
    this.forbidTagFn();
    this.tagSelectFn();
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
            name: item.name
          });
        } else {
          //无标签组，新增加
          this.selectTag.push({
            tag_group_id: this.leftMenu[this.menuIndex - 1].tag_group_id,
            tag_group_name: this.leftMenu[this.menuIndex - 1].tag_group_name,
            tags: [
              {
                tag_id: item.tag_id,
                name: item.name
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
              name: item.name
            }
          ]
        });
      }
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
    this.leftMenu.forEach((curItem: any) => {
      curItem.tags.forEach((subCurItem: any) => {
        if (subCurItem.tag_id == item.tag_id) {
          subCurItem.isSelected = false;
        }
      });
    });
    this.getUsercountFn();
  }
  //选中/从新建页面带来的tag
  private tagSelectFn() {
    if (this.selectTag.length != 0) {
      this.selectTag.forEach((item: any) => {
        if (
          item.tag_group_id == this.leftMenu[this.menuIndex - 1].tag_group_id
        ) {
          item.tags.forEach((subItem: any) => {
            this.leftMenu[this.menuIndex - 1].tags.forEach((curItem: any) => {
              if (subItem.tag_id == curItem.tag_id) {
                curItem.isSelected = true;
              }
            });
          });
        }
      });
    }
  }
  //取消-返回上一页
  private goBackFn() {
    this.$router.go(-1);
    if (this.paramsR.type == "include") {
      CustomerModule.setGroupS(this.cancelObj.include);
      CustomerModule.setGroupU(this.cancelObj.exclude);
    } else {
      CustomerModule.setGroupS(this.cancelObj.include);
      CustomerModule.setGroupU(this.cancelObj.exclude);
    }
  }
  //保存
  saveTagFn() {
    // this.$router.push({
    //   name: "myapp_customer_customer-groups_create",
    // });
    this.$router.go(-1);
    if (this.paramsR.type == "include") {
      CustomerModule.setGroupS(this.selectTag);
    } else {
      CustomerModule.setGroupU(this.selectTag);
    }
  }

  /*************************** 生命周期 start ***************************/
  private created() {
    this.paramsR = this.$route.params;
    this.cancelObj = JSON.parse(JSON.stringify(this.$route.params));
    //获取标签组列表
    this.getLabelListFn();
    // type=='include'||'exclude' 包含||不包含
    if (this.paramsR.type == "include") {
      this.selectTag = this.paramsR.include;
      this.forbidTag = this.paramsR.exclude;
    } else if (this.paramsR.type == "exclude") {
      this.selectTag = this.paramsR.exclude;
      this.forbidTag = this.paramsR.include;
    }
    // //初始计算带过来的tag个数
    // this.selectTag.forEach((item: any) => {
    //   this.selectTagsNum += item.tags.length;
    // });
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
}
</script>
<style lang="scss" scoped>
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
ul {
  list-style: none;
}
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
    .title {
      background: #ebebf1;
      font-size: 14px;
      color: $color-text-main;
      height: 36px;
      line-height: 36px;
      padding-left: 20px;
      font-weight: bold;
    }
    .content-wrapper {
      height: calc(100% - 36px);
      overflow-y: auto;
      overflow-x: hidden;
    }
  }
  &-left,
  &-middle {
    dd {
      margin-inline-start: 0px;
      height: 20px;
      line-height: 20px;
      // background-color: rgba(255, 255, 255, 1);
      color: $color-text-main;
      font-size: 14px;
      padding: 10px 0;
      cursor: pointer;
      span {
        display: block;
        padding-left: 20px;
      }
    }
  }
  &-left {
    dd {
      span {
        border-left: 4px solid $color-background;
      }
    }

    dd.menuAactive {
      background: $color-background;
      span {
        border-left: 4px solid $color-main !important;
        color: $color-main;
        margin-right: 20px;
      }
    }
    dd:hover {
      background: #ebebf1;
      span {
        border-left: 4px solid #ebebf1;
      }
    }
  }
  &-middle {
    border-left: 1px solid #ebebf1;
    border-right: 1px solid #ebebf1;
    dd {
      .group-selected {
        display: none;
      }
    }
    dd.unDisabled {
      // background: $color-text-normal;
      cursor: not-allowed;
      color: $color-disable;
    }
    dd.unDisabled:hover {
      // background: $color-text-normal;
      span {
        color: $color-disable;
      }
      .group-selected {
        display: none;
      }
    }
    dd.menuAactive {
      background: $color-background;
      span {
        color: $color-main;
        margin-right: 20px;
      }
      .group-selected {
        display: inline-block;
      }
    }
    dd:hover {
      background: #ebebf1;
    }
  }
  &-right {
    dd {
      padding: 20px;
    }
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
        background: $color-background;
        right: -10px;
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
.current-total {
  font-size: 14px;
  color: $color-text-main;
  line-height: 20px;
  margin-top: 20px;
}
.no-data {
  text-align: center;
  padding: 100px 0;
  p {
    color: $color-normal;
    margin-bottom: 40px;
  }
}
</style>
