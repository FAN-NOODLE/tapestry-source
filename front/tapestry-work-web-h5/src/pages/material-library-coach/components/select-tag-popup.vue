<template>
  <van-popup
    v-model="visible"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    position="bottom"
  >
    <div class="popup">
      <div class="popup-content">
        <div v-if="list.length" class="list">
          <div
            v-for="(item, index) in list"
            :key="index"
            :class="['item', {
              active: item.checked
            }]"
            @click="handleSelect(item)"
          >
            {{ item.name }}
          </div>
        </div>
        <empty-page
          v-if="!list.length"
          picture="noData"
          paddingTop="200"
          paddingBottom="20"
          text="暂无素材标签"
        ></empty-page>
      </div>
      <div class="popup-footer">
        <van-button class="btn-li reset" @click="handleCancel">重置</van-button>
        <van-button class="btn-li confirm" :disabled="disabled" @click="handleConfirm">确定<span v-if="count">({{ count }})</span></van-button>
      </div>
    </div>
  </van-popup>
</template>

<script lang="ts">
import { Component, Vue, Prop,PropSync, Watch } from "vue-property-decorator";
import { Popup } from "vant";
import { categoryList } from "@/api/marketing-promotion/index";
import emptyPage from "@/components/base/empty-page.vue";
interface ITag {
  id: number;
  name: string;
  checked: boolean
}
@Component({
  components: {
    emptyPage,
    [Popup.name]: Popup,
  },
})
export default class SelectTagPopup extends Vue {
  @PropSync("show",{ type: Boolean }) visible: boolean;
  /** 选中数据 */
  @Prop({ type: Array, default: () => [] }) value: Array<number>;
  @Prop({ type: Number, default: 0 }) isCollect: number;
  /** 标签列表 */
  list: Array<ITag> = []; 

  get disabled() {
    return this.list.length === 0 || this.list.every(item => !item.checked)
  }

  get count() {
    return this.list.filter(item => item.checked).length
  }
  
  @Watch('value')
  handleValueChange(value: Array<number>) {
    this.list = this.list.map(item => {
      item.checked = value.includes(item.id);
      return item;
    })
  }

  mounted() {
    this.getCategoryList();
  }
  activated() {
    this.getCategoryList();
  }
  /** 
   * 获取标签
   */
  async getCategoryList() {
    const res: IApiResult = await categoryList({ page: 1, size: 1999, is_collect: this.isCollect });
    if (res.code == 200) {
      this.list = (res.data.list || []).map(item => {
        item.checked = this.value.includes(item.id);
        return item
      });
    }
  }
  /**
   * 选择
   * @param item 
   */
  handleSelect(item: ITag) {
    item.checked = !item.checked
  }
  /**
   * 取消
   */
  handleCancel() {
    this.visible = false;
    this.list = this.list.map(item => {
      item.checked = false;
      return item;
    })
    this.$emit('confirm', []);
  }
  /**
   * 确认
   */
  handleConfirm() {
    this.visible = false;
    const value = this.list.filter(item => item.checked).map(item => item.id);
    this.$emit('confirm', value);
  }
}
</script>

<style lang="scss" scoped>
.popup {
  max-height: 80vh;
  min-height: 40vh;
  padding: 15px 0 0;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  border-radius: 0;
  .popup-footer {
    flex-shrink: 1;
    width: 100%;
    display: flex;
    background-color: #ffffff;
    padding: 15px;
    .btn-li {
      flex: 2;
      border: 1px solid #4a4a4a;
      color: #4a4a4a;
      font-size: 14px;
      text-align: center;
      line-height: 36px;
    }
    .confirm {
      margin-left: 15px;
      color: #ffffff;
      background-color: #4a4a4a;
    }
  }
  .popup-content {
    flex-grow: 1;
    padding: 0 10px 0;
    overflow-y: auto;
    .list {
      display: flex;
      flex-wrap: wrap;
    }
    .item {
      margin-bottom: 10px;
      margin-right: 10px;
      padding: 6px 15px;
      color: #4a4a4a;
      font-size: 10px;
      background-color: #f4f6f8;
      &.active {
        color: #ffffff;
        background-color: #000000;
      }
    }
  }
}
</style>
