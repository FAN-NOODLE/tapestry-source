<template>
  <!-- 侧边菜单栏主要为菜单手风琴部分 -->
  <div class="aside-bar">
    <div class="crs-logo"></div>
    <!-- 一级菜单栏 -->
    <el-menu
      class="crs-menu"
      active-text-color="#F5F6FA"
      text-color="#000000"
      :default-active="ROOT"
    >
      <div
        v-for="(route, $index) in list"
        :key="route.path"
      >
        <!-- 展开的手风琴 -->
        <el-popover
          v-model="route.model"
          :visible-arrow="false"
          v-if="route.children && route.children.length > 0"
          placement="right"
          width="158"
          trigger="hover"
          popper-class="menu-popper"
          :open-delay="50"
          @show="showSecondMenu(route.children, $index)"
        >
          <!-- 二级菜单 -->
          <el-menu
            class="crs-menu-second"
            background-color="#000000"
            text-color="#FFFFFF"
            active-text-color="#FFFFFF"
            :ref="`window${$index}`"
            :default-active="LEAF"
          >
            <template v-for="_route in route.children">
              <!-- 有子集，可展开查看三级菜单栏 -->
              <el-submenu
                v-if="_route.children &&
                  _route.children.length > 0 &&
                  _route.children.some(
                    item => item.type === 1 || item.type === 4
                  )
                  "
                :key="_route.path"
                :index="_route.key"
              >
                <template slot="title">
                  <span
                    class="crs-menu-title"
                    v-text="_route.name"
                  ></span>
                </template>
                <el-menu-item-group>
                  <el-menu-item
                    :class="item.key == `${ROOT}_${_route.path}_${item.path}`
                      ? 'crs-menu-child group-child is-active-self'
                      : 'crs-menu-child group-child'
                      "
                    v-show="item.type === 1 || item.type === 4"
                    v-for="item in _route.children"
                    :key="item.path"
                    :index="item.key"
                    @click="selectMenu(item, 3)"
                  >
                    <span class="crs-menu-text">{{ item.name }}</span>
                  </el-menu-item>
                </el-menu-item-group>
              </el-submenu>
              <!-- 二级菜单栏无可展开的菜单项 -->
              <el-menu-item
                v-else-if="_route.type === 1 || _route.type === 4"
                :class="_route.key == `${ROOT}_${_route.path}`
                  ? 'crs-menu-child crs-menu-hasnot-children is-active-self'
                  : 'crs-menu-child crs-menu-hasnot-children is-active-self'
                  "
                :key="_route.path"
                :index="_route.key"
                @click="selectMenu(_route, 2)"
              >
                <span> {{ _route.name }} </span>
              </el-menu-item>
            </template>
          </el-menu>
          <!-- 一级菜单栏有子集的菜单项 -->
          <el-menu-item
            @click="selectMenu(route, 1)"
            v-if="route.type === 1 || route.type === 4"
            slot="reference"
            :index="route.key"
            :class="{
              'crs-menu-item': true,
              'has-children': true,
              'is-active-self': ROOT == route.key,
              'continue-hover': route.model
            }"
          >
            <svg-icon
              class="menu-icon"
              :icon-class="ROOT == route.key ? route.icon + '-selected' : route.icon
                "
            ></svg-icon>
            <span>{{ route.name }}</span>
          </el-menu-item>
        </el-popover>
        <!-- 一级菜单栏上的子菜单项，没有子集 -->
        <el-menu-item
          v-else-if="route.type === 1 || route.type === 4"
          class="crs-menu-item"
          :index="route.key"
          @click="selectMenu(route, 1)"
          :class="ROOT == route.key ? 'crs-menu-item  is-active' : 'crs-menu-item'
            "
        >
          <svg-icon
            class="menu-icon"
            :icon-class="ROOT == route.key ? route.icon + '-selected' : route.icon
              "
          ></svg-icon>
          <span class="menu-text">{{ route.name }}</span>
        </el-menu-item>
      </div>
    </el-menu>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import MenuModule from "@/store/modules/menu";
import PerModule from "@/store/modules/permission";
import { INodeItem } from "@/types";
@Component
export default class AsiderBar extends Vue {
  get MENU() {
    return PerModule.getPermission;
  }
  get ROOT() {
    return MenuModule.getRoot;
  }
  get LEAF() {
    return MenuModule.getLeaf;
  }
  mounted() {
    this.list = this.MENU.map(item => {
      return {
        ...{ model: false },
        ...item
      };
    });
  }
  private list: Array<any> = [];
  private selectMenu(item: any) {
    const route: any = this.getFirstChildMenu(item);
    this.$router.push({
      name: route.key,
    });
  }
  /**
   * @description 找到菜单的第一个子叶节点
   * @param node 菜单项
   * */
  private getFirstChildMenu(node: any): any {
    if (
      node.children &&
      node.children.length > 0 &&
      node.children[0].type !== "button"
    ) {
      const tmp: any = node.children[0];
      return this.getFirstChildMenu(tmp);
    } else {
      return node;
    }
  }
  /**
   * @description 二级菜单每次出现默认展开三级菜单
   * @param menu 二级菜单list
   * @param index 一级菜单索引
   * */
  private showSecondMenu(menu: Array<INodeItem>, index: number) {
    menu.forEach(_route => {
      if (
        _route.children &&
        _route.children.length > 0 &&
        _route.children.some(item => item.type !== 2)
      ) {
        this.$refs[`window${index}`][0].open(_route.key);
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.aside-bar {
  width: 100%;
  height: 100%;
  overflow: hidden;
  overflow-y: auto;
  text-align: center;

  .crs-logo {
    width: 40px;
    height: 40px;
    background-color: $color-main;
    background-image: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/logo.png");
    background-size: cover;
    border-radius: 40px;
    margin: 30px auto;
    text-align: center;
    line-height: 40px;
    box-shadow: 0px 0px 4px 0px rgba(0, 0, 0, 0.12);
  }

  .crs-menu {
    .menu-icon {
      margin-left: 4px;
      margin-right: 10px;
    }

    .menu-text {
      display: inline-block;
      vertical-align: middle;
      line-height: 14px;
    }

    margin-left: 10px;
    border-right: none !important;

    .crs-menu-item {
      text-align: left;
      height: 44px !important;
      line-height: 44px !important;
      margin-bottom: 10px;
      padding-left: 10px !important;

      &:hover {
        background-color: $root-menu-hover !important;
      }
    }

    .menu-poper {
      background: #000000;
    }

    // 一级菜单栏的持续hover效果
    .continue-hover {
      background-color: $root-menu-hover !important;
    }

    // 修改element默认选中样式
    /deep/.is-active.has-children {
      background: $color-dark !important;

      &:hover {
        background-color: $root-menu-hover !important;
      }
    }

    // 修改element默认选中样式
    /deep/.is-active {
      background: $color-main !important;

      &:hover {
        background-color: $color-main !important;
      }
    }

    // 自定义默认选中样式
    .has-children.is-active-self {
      background: $color-main !important;

      &:hover {
        background-color: $color-main !important;
      }
    }
  }
}

/*
* 二级菜单
*/
.crs-menu-second {
  border-right: none !important;
  background: $color-main !important;

  .crs-menu-title {
    color: #ffffff;
    font-weight: 550;
  }

  /deep/ {
    .el-submenu__title {
      height: 44px !important;
      line-height: 44px !important;
    }

    .el-menu-item {
      padding-left: 10px !important;
      height: 44px !important;
      line-height: 44px !important;
      border-radius: 0 !important;
    }

    .el-submenu__title {
      padding-left: 10px !important;
    }

    .el-menu-item-group__title {
      padding: 0 !important;
    }

    .el-submenu>.el-submenu__title .el-submenu__icon-arrow {
      font-weight: 550;
      line-height: 16px;
      font-size: 14px !important;
      position: absolute;
      right: 10px;
    }
  }

  .crs-menu-child {
    color: #ffffff !important;
    margin-bottom: 3px;

    &:hover {
      color: #ffffff !important;
      background-color: #232323 !important;
    }
  }

  /*
  * 二级菜单无子集样式
  */
  // .crs-menu-hasnot-children {
  //   font-weight: 550;
  //   color: $color-text-main !important;
  // }
  /*
  * 选中样式
  */
  .is-active {
    background-color: #ffffff !important;
    color: $color-main !important;

    &:hover {
      background-color: #ffffff !important;
    }
  }

  .is-active.crs-menu-child.is-active-self {
    background: #ffffff !important;
    color: $color-main !important;

    &:hover {
      background-color: $color-background !important;
    }
  }

  .is-active.crs-menu-child.is-active-self.group-child {
    color: $color-text-main !important;
    padding-left: 10px !important;
    background: $color-background !important;

    &:hover {
      background-color: $color-background !important;
      color: $color-main !important;
    }
  }
}

/*
* 三级菜单项
*/
.group-child {
  // max-width: 149px !important;
  padding-left: 10px !important;
  min-width: 0 !important;
}

/deep/.el-submenu__title {
  position: relative;

  &:hover {
    color: #ffffff !important;
    background-color: #232323 !important;
  }
}
</style>
