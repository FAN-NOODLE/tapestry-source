<template>
  <div class="audit">
    <div class="audit-header">
      <van-search
        v-model="searchParams.title"
        clearable
        clear-trigger="always"
        placeholder="请输入标题名称搜索"
        :left-icon="searchIcon"
        @input="handleSearch"
        @click="handleInputClick"
      />
      <com-select v-model="searchParams.type" :options="typeOptions" @change="handleTypeChange"></com-select>
      <div class="line"></div>
      <com-select v-model="searchParams.status" :options="statusOptions" @change="handleStatusChange"></com-select>
    </div>
    <div class="audit-content">
      <van-list
        v-model="loading"
        :finished="finished"
        :finished-text="list.length > 0 ? '没有更多了~':''"
        @load="fetchData"
      >
        <communication-item v-for="(item, index) in list" :key="index" :item="item" type="ft" @update="handleUpdate($event, index)"></communication-item>
      </van-list>
      <empty-page
        v-if="list.length == 0 && !loading"
        :picture="emptyIcon"
        paddingTop="60"
        paddingBottom="0"
        :text="emptyText"
      ></empty-page>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import CommunicationItem, {ICommunicationItem} from './components/communication-item.vue'
import EmptyPage from "@/components/base/empty-page.vue";
import ComSelect from './components/com-select.vue'
import { statusOptions, typeOptions } from './const';
import { getContentListApi } from '@/api/communication';

@Component({
  components: {
    CommunicationItem,
    EmptyPage,
    ComSelect,
  }
})
export default class CommunicationMyList extends Vue {
  typeOptions = [
    {
      label: '全部类型',
      value: -1
    },
    ...typeOptions
  ]
  statusOptions = [{
      label: '全部状态',
      value: -1
    },
    ...statusOptions
  ]
  /** 搜索图标 */
  searchIcon = require('@/assets/communication/icon_search.png')
  /** 是否正在加载数据  */
  loading = true;
  /** 是否没有更多数据 */
  finished = false;
  /** 搜索参数 */
  searchParams = {
    title: '',
    type: -1,
    status: -1
  }
  /** 分页 */
  meta: IMeta = {
    page: 1,
    size: 10,
    total: 0,
  }
  /** 列表数据 */
  list: ICommunicationItem[] = []

  /** 空页面提示 */
  get emptyText() {
    if(this.searchParams.title) {
      return '未找到相关传讯'
    }
    return `暂无传讯`
  }
  /** 空页面图标 */
  get emptyIcon() {
    if(this.searchParams.title) {
      return 'noSearchData'
    }
    return 'noDisplayTask'
  }

  mounted() {
    this.fetchData()
  }

  /** 
   * @description 获取列表数据
   */
  async fetchData() {
    this.loading = true;
    try {
      const params = {
        page: this.meta.page,
        size: this.meta.size,
        ...this.searchParams
      }
      const {code, message, data}: IApiResult = await getContentListApi(params)
      if (code !== 200) {
        this.finished = true;
        this.loading = false;
        this.$toast.fail(message);
        return;
      }
      this.loading = false;
      this.finished = data.list.length === 0;
      this.list = this.list.concat(data.list)
      this.meta.page = this.meta.page + 1

    } catch (error) {
      console.log(error);
      this.finished = true;
      this.loading = false;
    }
  }

  handleInputClick() {
    this.$track.trackSearchClick({
      tab_name: '我的传讯'
    })
  }

  handleTypeChange(e) {
    console.log(e);
    const current = this.typeOptions.find(item => item.value == e)
    this.$track.trackFilter({
      tab_name: '我的传讯',
      filter_type_sv: '传讯类型',
      filter_option: current.label
    })
    this.handleSearch()
  }
  handleStatusChange(e) {
    const current = this.statusOptions.find(item => item.value == e)
    this.$track.trackFilter({
      tab_name: '我的传讯',
      filter_type_sv: '传讯状态',
      filter_option: current.label
    })
    this.handleSearch()
  }

  /**
   * @description 搜索
   */
  handleSearch() {
    this.meta.page = 1;
    this.list = [];
    this.finished = false;
    this.fetchData();
  }
  handleUpdate(item, index) {
    this.$set(this.list, index, item)
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .van-search {
    background-color: transparent;
    padding: 0;
    height: 38px;
    border: 1px solid #A6A6A6;
    width: 170px;
    margin-right: 15px;
  }

  .van-cell {
    font-size: 12px;
  }

  .van-field__left-icon {
    display: flex;
    align-items: center;
  }
}
.audit-header {
  padding: 15px 15px;
  display: flex;
  align-items: center;
  background-color: #f7f7f7;
  z-index: 1;
  position: sticky;
  top: 0;

  .com-select {
    flex: 1;
  }

  .line {
    width: 1px;
    height: 21px;
    background: #D8D8D8;
    margin: 0 15px;
  }
}
.audit-content {
  padding: 0 15px;
}
</style>