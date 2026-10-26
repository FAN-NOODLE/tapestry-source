<template>
  <van-popup
    v-model="show"
    @close="groupCancel"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    round
    position="bottom"
    class="material-tag"
  >
    <div class="cate-wrap" v-if="categoryList.length">
      <ul class="cate-ul">
        <li
          v-for="(item, index) in categoryList"
          :key="index"
          class="cate-li"
          :class="[categoryIds.indexOf(item.id) > -1 ? 'active' : '']"
          @click="selectCategory(item)"
        >
          {{ item.name }}
        </li>
      </ul>
    </div>
    <empty-page
      v-if="!categoryList.length"
      picture="noData"
      paddingTop="20"
      paddingBottom="20"
      text="暂无素材标签"
    ></empty-page>
    <div class="btn-view">
      <div class="btn-li reset" @click="groupCancel">取消</div>
      <div class="btn-li confirm" @click="groupCheck">确定</div>
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
}
@Component({
  components: {
    emptyPage,
    [Popup.name]: Popup,
  },
})
export default class MaterialTag extends Vue {
  @PropSync("showPopup",{ type: Boolean }) show;
  @Prop({
    type: Array,
    default: () => {
      return [];
    },
  })
  ids;
  @Prop({
    type: Array,
    default: () => {
      return [];
    },
  })
  names;

  @Watch("showPopup")
  showChange(e) {
    if (e) {
      this.categoryIds = JSON.parse(JSON.stringify(this.ids));
      this.categoryName = JSON.parse(JSON.stringify(this.names));
    }
  }
  private categoryIds: Array<number> = []; // 已选中的标签id
  private categoryName: Array<string> = []; // 已选中的标签名称
  private categoryList: Array<ITag> = []; // 标签列表
  /**获取标签 */
  private async getCategoryList() {
    //TODO: 待优化
    const res: IApiResult = await categoryList({ page: 1, size: 1999 });
    if (res.code == 200) {
      this.categoryList = res.data.list || [];
    }
  }
  /**选择标签 */
  private selectCategory(tag) {
    const index = this.categoryIds.indexOf(tag.id);
    if (index > -1) {
      this.categoryIds.splice(index, 1);
      this.categoryName.splice(index, 1);
    } else {
      this.categoryIds.push(tag.id);
      this.categoryName.push(tag.name);
    }
  }
  /** 分组筛选取消 */
  private groupCancel() {
    this.$emit("cancel");
  }
  /** 分组筛选确定 */
  private groupCheck() {
    const categoryIds = JSON.parse(JSON.stringify(this.categoryIds));
    const categoryName = JSON.parse(JSON.stringify(this.categoryName));
    this.$emit("confirm", { names: categoryName, ids: categoryIds });
  }
  /**生命周期 */
  mounted() {
    this.getCategoryList();
  }
}
</script>

<style lang="scss" scoped>
.material-tag {
  max-height: 453px;
  padding: 15px 0 0;
  display: flex;
  justify-content: space-between;
  flex-direction: column;
  border-radius: 0;
  .btn-view {
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
  .cate-wrap {
    flex-grow: 1;
    padding: 0 10px 0;
    overflow-y: auto;
    .cate-ul {
      display: flex;
      flex-wrap: wrap;
    }
    .cate-li {
      text-align: center;
      min-width: 78px;
      margin: 0 5px 15px;
      padding: 6px 18px;
      color: #606266;
      font-size: 14px;
      background-color: #f5f5f5;
      &.active {
        color: #ffffff;
        background-color: #000000;
      }
    }
  }
}
</style>
