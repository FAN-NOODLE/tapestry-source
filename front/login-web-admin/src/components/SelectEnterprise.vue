<template lang="pug">
  el-dialog(title="选择企业" :visible.sync="show" :before-close="handleClose" :close-on-click-modal="false" width="830px")
    .dialog-body
      .list-box
        .list
          .item(v-for="(item,index) in tableView" :key="item.id" :class="[item.checked?'checked':'']" @click="handleSelect(item,index)")
            .item-info
                img.logo-img(v-if="item.logo" :src="item.logo")
                .logo(v-if="!item.logo") LOGO
                .title-box
                  .title {{item.shop_name}}
                  .enter-btn(v-if="!has_expire(item)") 进入
                  //- 过期
                  img.has-expire(v-if="has_expire(item)" src="@/assets/has_expire.png")
                  .ex-btn(v-if="has_expire(item)") 续费
            p.mt10.character-hiding-o 企业名称： {{item.name}}
      .pagination-box(v-show="meta.total > 6")
        el-pagination(
          :page-size="meta.size"
          :total="meta.total"
          layout="prev, pager, next"
          :current-page="meta.page"
          @current-change="handleCurrentChange")
</template>
<script>
export default {
  name: "ComSelectEnterprise",
  components: {},
  props: {
    //控制弹窗 显示
    show: {
      type: Boolean,
      default: false
    },
    listData: {
      type: Array,
      default: null
    },
    onClose: {
      type: Function,
      default: () => {}
    },
    onSelect: {
      type: Function,
      default: () => {}
    }
  },
  computed: {},
  data() {
    return {
      meta: {
        total: null, // 总数
        page: 1, // 当前页
        size: 6 // 每页数量
      },
      tableView: []
    };
  },
  methods: {
    has_expire(item) {
      var timestamp = new Date().getTime() / 1000;
      return timestamp > item.account_expire_at && item.version > 1;
    },
    handleClose() {
      this.onClose();
    },
    handleSelect(item, index) {
      // 跳转立即续费文档
      if (this.has_expire(item)) {
        window.open(
          "http://xueyuan.vchangyi.com/bin/view/CRS%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E7%AE%A1%E7%90%86%E5%90%8E%E5%8F%B0/%E6%99%BA%E6%85%A7%E5%AF%BC%E8%B4%AD%E4%BB%98%E8%B4%B9%E8%AE%A2%E9%98%85/",
          "_blank"
        );
      } else {
        var timestamp = new Date().getTime() / 1000;
        if (timestamp > item.account_expire_at && item.version > 1) {
          return false;
        }
        this.tableView.forEach(item => {
          item.checked = false;
        });
        this.tableView[index].checked = true;
        this.onSelect(item);
      }
    },
    // 选择页码
    handleCurrentChange(val) {
      this.meta.page = val;

      console.log(val);
      //为什么必须不是1 呢 因为点击1是返回初始数据
      if (val != 1) {
        val = (val - 1) * this.meta.size; //比如点击是2，并且size是10的时候 那么数据下标是 10开始就是10
      } else {
        val = 0;
      }
      // 当点击2的时候如果第二页是最后一页进入else全部展示
      if (this.meta.total > val + this.meta.size) {
        this.tableView = this.listData.slice(val, val + this.meta.size);
      } else {
        this.tableView = this.listData.slice(val, this.listData.length);
      }
    },
    // 处理列表，分页
    pageDeal() {
      this.tableView = [];
      this.meta.total = this.listData.length;
      //进入if 因为他有多页码 else 只有1页
      if (this.listData.length >= this.meta.size) {
        for (let i = 0; i < this.meta.size; i++) {
          this.tableView.push(this.listData[i]);
        }
      } else {
        this.tableView = this.listData;
      }
      console.log(this.tableView);
    }
  },
  mounted() {
    this.pageDeal();
  }
};
</script>
<style lang="sass" scoped>
.mt10
  margin-top: 10px
/deep/.scrollbar-wrapper
  overflow-x: hidden !important
  margin-top: 0px
/deep/.el-dialog__header
  padding: 30px 30px 10px 30px
/deep/.el-dialog__body
  padding: 0 30px 30px 30px
.dialog-body
    height: 440px
/deep/.el-dialog__title
  font-weight: 700
.list-box
  height: 412px
.list
  overflow: hidden
  display: flex
  justify-content: space-between
  flex-wrap: wrap
  flex-direction: row
  padding: 2px 2px 2px 2px
  box-sizing: border-box
  .item
    position: relative
    overflow: hidden
    width: 44.2%
    font-size: 12px
    height: 82px
    border-radius: 4px
    margin-bottom: 20px
    cursor: pointer
    padding: 20px 15px 20px 20px
    box-shadow: 0 0 0 1px #ebeef5
    font-size: 16px
    .item-info
      display: flex
      line-height: 45px
      .logo-img
        min-width: 50px
        min-height: 50px
        max-width: 50px
        max-height: 50px
        border-radius: 50%
      .logo
        box-shadow: 0 0 0 1px #C0C4CC
        border-radius: 50%
        min-width: 50px
        min-height: 50px
        max-width: 50px
        max-height: 50px
        color: #999
        font-size: 12px
        font-weight: 700
        line-height: 50px
        text-align: center
      .title-box
        display: flex
        justify-content: space-between
        width: 100%
        align-items: center
        height: 50px
        .title
          color: #606266
          font-size: 16px
          margin-left: 10px
          font-weight: 500
          width: 144px
          line-height: 22px
          display: -webkit-box
          overflow: hidden
          text-overflow: ellipsis
          word-break: break-all
          -webkit-box-orient: vertical
          -webkit-line-clamp: 2
        .enter-btn
          visibility: hidden
          padding: 0 12px
          background-color: #708DF5
          color: #fff
          border-radius: 2px
          height: 27px
          line-height: 27px
          font-size: 14px
        .has-expire
          width: 55px
          height: 47px
        .ex-btn
          visibility: hidden
          padding: 0 12px
          background-color: #708DF5
          color: #fff
          border-radius: 2px
          height: 27px
          line-height: 27px
          font-size: 14px
    .character-hiding-o
      color: #909399
      font-size: 16px
  .item
    &:hover
      box-shadow: 0 0 0 1px #708DF5
      .item-info .title-box .enter-btn, .item-info .title-box .ex-btn
        visibility: visible
  .checked
      box-shadow: 0 0 0 1px #708DF5 !important
.pagination-box
  width: 100%
  text-align: right
  padding-top: 13px
</style>
