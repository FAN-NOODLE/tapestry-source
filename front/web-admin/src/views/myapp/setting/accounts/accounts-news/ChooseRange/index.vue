<template>
  <el-dialog
    title="选择账号范围"
    :visible.sync="show"
    :before-close="handleClose"
    :close-on-click-modal="false"
    width="900px"
  >
    <div class="dialog-body">
      <el-tabs v-model="active" tabPosition="left" style="height: 375px">
        <el-tab-pane label="门店" name="1">
          <el-row>
            <el-col class="left" :span="12">
              <cy-input
                prefix-icon="el-icon-search"
                placeholder="请输入门店名称"
                clearable="clearable"
                v-model="filterText"
              ></cy-input>
              <el-scrollbar
                wrap-class="scrollbar-wrapper"
                style="height:340px;width: 100% !important;"
              >
                <div class="box-all" v-if="listData.length">
                  <el-checkbox
                    v-model="check_all"
                    @change="handleCheckAllChange"
                    >全部</el-checkbox
                  >
                </div>
                <el-tree
                  :data="listData"
                  @check="handleCheckNode"
                  show-checkbox="show-checkbox"
                  node-key="id"
                  default-expand-all="default-expand-all"
                  :filter-node-method="filterNode"
                  ref="tree"
                  empty-text="暂无门店"
                ></el-tree>
              </el-scrollbar>
            </el-col>
            <el-col :span="10">
              <div class="title">已选门店范围</div>
              <el-scrollbar
                wrap-class="scrollbar-wrapper"
                style="height:320px;width: 100% !important;padding-left: 30px;"
              >
                <div class="content" v-if="!check_all">
                  <div class="info" v-for="item in showList" :key="item.id">
                    <div class="info-item">
                      <span>{{ item.text }}</span
                      ><span
                        class="btn-delete el-icon-error"
                        @click="handleDetele(item.item)"
                      ></span>
                    </div>
                  </div>
                </div>
                <div class="content" v-if="check_all">
                  <div class="info-item">
                    <span>{{ showList.length ? showList[0].text : "" }}</span
                    ><span
                      class="btn-delete el-icon-error"
                      @click="handleDetele()"
                    ></span>
                  </div>
                </div>
              </el-scrollbar>
            </el-col>
          </el-row>
        </el-tab-pane>
        <el-tab-pane label="商城" name="2">
          <el-row>
            <el-col class="left" :span="12">
              <cy-input
                prefix-icon="el-icon-search"
                placeholder="请输入商城名称"
                clearable="clearable"
                v-model="storeSelected"
                @input="handleInput"
                style="padding-bottom: 20px"
              ></cy-input>
              <el-checkbox-group
                v-model="selectMallList"
                @change="handleChange"
              >
                <el-scrollbar
                  wrap-class="scrollbar-wrapper"
                  style="height:380px;width: 100% !important;"
                >
                  <el-checkbox
                    :label="item.id"
                    v-for="item in mallList"
                    :key="item.id"
                  >
                    {{ item.name }}
                  </el-checkbox>
                </el-scrollbar>
              </el-checkbox-group>
            </el-col>
            <el-col :span="12">
              <div class="title">已选商城</div>
              <el-scrollbar
                wrap-class="scrollbar-wrapper"
                style="height:320px;width: 100% !important;padding-left: 30px;"
              >
                <div class="content" v-if="selectMallList">
                  <div
                    class="info"
                    v-for="item in selectedMallName"
                    :key="item.id"
                  >
                    <div class="info-item">
                      <span>{{ item.name }}</span
                      ><span
                        class="btn-delete el-icon-error"
                        @click="handleDeteleMall(item.id)"
                      ></span>
                    </div>
                  </div>
                </div>
              </el-scrollbar>
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </div>
    <div class="btn-position">
      <el-button size="small" @click="handleClose()">取消</el-button>
      <el-button type="primary" size="small" @click="handleSave()"
        >确定
      </el-button>
    </div>
  </el-dialog>
</template>
<script>
let TEMP_COUNT = 0;
let TEMP_ITEM = null; //递归保存临时选中项
import Api from "@/api/setting/Account";
export default {
  name: "ChooseRange",
  components: {},
  props: {
    //控制弹窗 显示
    show: {
      type: Boolean,
      default: false
    },
    onClose: {
      type: Function,
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      default: () => {}
    },
    onUpdate: {
      type: Function,
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      default: () => {}
    },
    defaultList: {
      type: Array,
      default: null
    }
  },
  watch: {
    filterText(val) {
      this.$refs.tree.filter(val);
    },
    // 展示所需list加工
    selectData(val) {
      this.showList = [];
      let showItem = null;
      const tempVal = JSON.parse(JSON.stringify(val));
      if (this.check_all) {
        showItem = this.resultLable(this.listData);
        if (showItem) {
          this.showList.push(showItem);
        }
      } else {
        tempVal.forEach(item => {
          delete item.children;
          showItem = this.resultLable(item);
          if (showItem) {
            this.showList.push(showItem);
          }
        });
      }
    }
  },
  data() {
    return {
      active: "1",
      filterText: "",
      check_all: false,
      selectIdsList: [], //门店已选择id集合
      selectData: [], //门店已选择数据集合
      showList: [], //门店已选择数据处理后的展示集合
      selectMallList: [], //商城已选择
      listDataLenth: 0, //总数据长度
      mallList: [],
      listData: [], //区域门店数据
      storeSelected: "",
      selectedMallName: []
    };
  },
  methods: {
    init(val) {
      if (!this.$refs.tree) {
        setTimeout(() => {
          this.init(val);
        }, 300);
      } else {
        this.$refs.tree.setCheckedNodes(val);
        if (!val || !val.length) {
          this.check_all = false;
          this.showList = [];
          this.selectMallList = [];
        } else {
          let isAll = false;
          let hasAre = false;
          let hasMall = false;
          val.forEach(item => {
            if (item.type == 3) {
              hasMall = true;
              this.selectMallList = [item.id];
            }
            if (item.type == 2) {
              hasAre = true;
            }
            if (item.type == 0) {
              isAll = true;
            }
          });
          if (!hasAre && !isAll) {
            console.log("没有区域");
            this.check_all = false;
            this.showList = [];
          }
          if (!hasMall) {
            console.log("没有商城");
            this.selectMallList = [];
          }
          if (isAll) {
            console.log("全部门店");
            this.check_all = true;
            this.handleCheckAllChange();
          }
        }
        this.selectIdsList = this.$refs.tree.getCheckedKeys();
        this.selectData = this.$refs.tree.getCheckedNodes();
      }
    },
    handleDeteleMall(id) {
      const flag = this.selectMallList.indexOf(id);
      this.selectMallList.splice(flag, 1);
      this.selectedMallName.forEach((item, key) => {
        if (item.id === id) {
          this.selectedMallName.splice(key, 1);
        }
      });
    },
    handleChange(val) {
      const arr = [];
      val.forEach(id => {
        this.mallList.forEach(item => {
          if (id === item.id) {
            arr.push({ id: item.id, name: item.name });
          }
        });
      });
      this.selectedMallName = arr;
    },
    handleInput(val) {
      const tempData = JSON.parse(localStorage.getItem("USERINFO"));
      if (val != "") {
        if (this.mallList.length === 0) {
          this.mallList.push({
            name: tempData.name,
            id: tempData.enterprise_id
          });
        }
        this.mallList = this.mallList.filter(item => {
          return item.name.indexOf(val) > -1;
        });
      } else {
        this.mallList = [];
        this.mallList.push({
          name: tempData.name,
          id: tempData.enterprise_id
        });
      }
    },
    //保存门店选择
    handleSave() {
      if (this.check_all) {
        this.onUpdate(1, [{ type: 0, label: "全部", id: 0 }]);
      } else {
        this.onUpdate(
          1,
          this.showList.map(res => {
            return res.item;
          })
        );
      }
      if (this.selectMallList && this.selectMallList.length) {
        this.onUpdate(2, {
          id: this.mallList[0].id,
          label: this.mallList[0].name,
          type: 3
        });
      }
      this.handleClose();
    },
    //删除已选项目
    handleDetele(item) {
      if (item) {
        this.deleteAllIdsByItem(item);
      } else {
        this.selectIdsList = [];
      }
      this.$refs.tree.setCheckedKeys(this.selectIdsList);
      this.selectData = this.$refs.tree.getCheckedNodes();
      this.checkIsAll();
    },
    // 递归删除当前选项下所有节点
    deleteAllIdsByItem(data) {
      this.selectIdsList.forEach((i, index) => {
        if (i == data.id) {
          this.selectIdsList.splice(index, 1);
        }
      });
      data.children &&
        data.children.length &&
        data.children.forEach(i => {
          this.deleteAllIdsByItem(i);
        });
    },
    //TODO:获取数据列表
    getData() {
      Api.areaStoreList().then(res => {
        this.listData = res.data.list || [];
        this.getListLength(this.listData);
      });
      const tempData = JSON.parse(localStorage.getItem("USERINFO"));
      this.mallList = [
        {
          name: tempData.name,
          id: tempData.enterprise_id
        }
      ];
    },
    // 全选或反选
    handleCheckAllChange() {
      if (this.check_all) {
        this.$refs.tree.setCheckedNodes(this.listData);
      } else {
        this.$refs.tree.setCheckedKeys([]);
      }
      this.selectIdsList = this.$refs.tree.getCheckedKeys();
      this.selectData = this.$refs.tree.getCheckedNodes();
    },
    //右侧展示列表加工
    resultLable(item) {
      TEMP_COUNT = 0;
      if (this.check_all) {
        this.findShopCountAll(item);
        return {
          text: `全部门店（共${TEMP_COUNT}家门店）`,
          item: this.listData
        };
      } else if (!item.pid || !this.selectIdsList.includes(item.pid)) {
        if (item.type == 2) {
          return {
            text: item.label,
            item: item
          };
        } else {
          this.findTarget(this.listData, item.id);
          this.findShopCount(TEMP_ITEM);
          return {
            text: `${TEMP_ITEM.label}（共${TEMP_COUNT}家门店）`,
            item: TEMP_ITEM
          };
        }
        // }
      }
    },
    // 递归查询当前对应的选项
    findTarget(source, targetId) {
      source.forEach(item => {
        if (item.id == targetId) {
          TEMP_ITEM = item;
          return item;
        } else if (item.children && item.children.length) {
          return this.findTarget(item.children, targetId);
        }
      });
    },
    // 递归查询所有的门店数量
    findShopCountAll(data) {
      if (data && data.length) {
        data.forEach(item => {
          if (item.type === 2) {
            TEMP_COUNT += 1;
          }
          this.findShopCountAll(item.children);
        });
      }
    },
    // 递归查询当前节点下的门店数量
    findShopCount(data) {
      if (data.type === 2) {
        TEMP_COUNT += 1;
      }
      data.children &&
        data.children.length &&
        data.children.forEach(i => {
          this.findShopCount(i);
        });
    },
    // 获取所有内容的长度-->全选判断使用
    getListLength(data) {
      if (data && data.length) {
        this.listDataLenth += data.length;
        data.forEach(item => {
          this.getListLength(item.children);
        });
      }
    },
    // 选择监听
    handleCheckNode() {
      this.selectIdsList = this.$refs.tree.getCheckedKeys();
      this.selectData = this.$refs.tree.getCheckedNodes();
      this.checkIsAll();
    },
    // 刷新全选box
    checkIsAll() {
      if (this.listDataLenth == this.selectIdsList.length) {
        this.check_all = true;
      } else {
        this.check_all = false;
      }
    },
    filterNode(value, data) {
      if (!value) return true;
      return data.label.indexOf(value) !== -1;
    },
    handleClose() {
      this.onClose();
    }
  },
  mounted() {
    this.getData();
  }
};
</script>
<style lang="scss" scoped>
/deep/ {
  .el-tabs--left .el-tabs__nav-wrap.is-left::after {
    height: 100%;
    width: 1px;
    bottom: auto;
    top: 0;
    background-color: #ebeef5;
  }
  .el-tabs__item {
    color: #303133 !important;
    &.is-active {
      color: #414244 !important;
      font-weight: bold !important;
    }
  }
  .el-dialog__body {
    padding-top: 0;
  }
  .el-input__inner {
    height: 35px !important;
  }
  .scrollbar-wrapper {
    overflow-x: hidden !important;
  }
  .el-tabs__header {
    margin: 0 0 20px !important;
  }
  .left {
    padding-right: 30px !important;
  }
  .el-tree {
    padding: 0 0 20px 0 !important;
  }
}
.left {
  padding: 0 30px;
  border-right: 1px solid #ebeef5;
}
.el-checkbox {
  display: block;
  margin-bottom: 10px;
}
.box-all {
  margin-top: 20px;
}
.content {
  .info-item {
    display: flex;
    justify-content: space-between;
    width: 90%;
    margin-top: 10px;
    .btn-delete {
      cursor: pointer;
    }
  }
}
.btn-box {
  background-color: #fff;
  position: absolute;
  bottom: 0;
  right: 10px;
}
.el-button {
  + {
    .el-button {
      margin-left: 20px !important;
    }
  }
}
.cancel {
  &:hover {
    color: #798df5 !important;
    background: #f0f3fe !important;
    border: 1px solid #d4dcfc !important;
  }
}
.title {
  padding-left: 30px !important;
  padding-top: 10px;
}
.scrollbar-wrapper.el-scrollbar__wrap {
  margin-left: 8px !important;
}
.btn-position {
  display: flex;
  justify-content: flex-end;
}
</style>
