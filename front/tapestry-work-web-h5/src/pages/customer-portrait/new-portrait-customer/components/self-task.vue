<template>
  <div class="sow">
    <div class="new-sow">
      <div class="top-box">
        <div class="title">自设任务</div>
        <div class="time-select">
          <img
            class="icon-1"
            src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/calendar-icon.png"
          />
          <div class="placeholder">
            <div v-if="info.isRelation" @click="showTimePicker = true">
              {{ dateTimeText ? dateTimeText : "选择提醒时间" }}
            </div>
            <!-- 好友显示 -->
            <div v-if="!info.isRelation" @click="tapDisableEdit">
              选择提醒时间
            </div>
          </div>
          <img
            class="icon-2"
            src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/arrow-bottom.png"
          />
        </div>
      </div>
      <div class="edit-box" @click="tapDisableEdit">
        <!-- 非好友禁用 -->
        <van-field
          class="edit"
          v-model="editValue"
          :disabled="!info.isRelation"
          :maxlength="200"
          autosize
          show-word-limit
          type="textarea"
          placeholder="请输入自设任务"
        />
      </div>
      <div
        :class="{
          btn: true,
          disabled: !dateTimeText || !editValue,
        }"
        @click="createSow"
      >
        创建
      </div>
    </div>

    <div class="sow-record">
      <div class="record-top">
        <div class="title">自设任务记录</div>
        <div class="filter-list">
          <div
            v-for="(item, index) in filterList"
            @click="tapFilter(item)"
            :key="index"
            :class="{
              'filter-item': true,
              selected: filterKey === item.key,
            }"
          >
            {{ item.label }}
          </div>
        </div>
      </div>
      <div v-if="sowList.length > 0" class="record-list">
        <van-list
          v-model="isLoading"
          @load="getList"
          :finished="!this.meta.isMore"
          finished-text="已经到底了~"
        >
          <div v-for="item in sowList" :key="item.id" class="record-item">
            <div class="top-info">
              <img
                class="icon"
                src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/clock-icon.png"
              />
              <div class="time">{{ item.remind_time }}</div>
              <div :class="{ status: true, finished: item.status === 1 }">
                {{ item.status === 1 ? "待完成" : "已完成" }}
              </div>
            </div>
            <div class="detail">{{ item.content }}</div>
            <div v-if="item.status === 1" class="btn-box">
              <div class="delete-btn" @click="openDeleteModel(item.id)">
                删除
              </div>
              <div class="confirm-btn" @click="openTipsModel(item.id)">
                确认完成
              </div>
            </div>
          </div>
        </van-list>
      </div>
      <div v-else class="empty-list">
        <img
          class="image"
          src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/no-sow-icon.png"
        />
        <div class="word">暂无自设任务记录</div>
      </div>
    </div>
    <van-popup v-model="showTimePicker" position="bottom">
      <van-datetime-picker
        v-if="info.isRelation"
        v-model="dateTime"
        type="datetime"
        title="选择提醒时间"
        :min-date="minDate"
        :max-date="maxDate"
        @confirm="timeConfirm"
        @cancel="showTimePicker = false"
      />
    </van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import dayjs from "dayjs";
import { Field, DatetimePicker, Popup, List, Dialog } from "vant";
import {
  createSowTaskApi,
  getSowTaskListApi,
  finishedSowTaskApi,
  deleteSowTaskApi,
} from "@/api/sow-api/index";
@Component({
  components: {
    [Field.name]: Field,
    [DatetimePicker.name]: DatetimePicker,
    [Popup.name]: Popup,
    [List.name]: List,
    // [Dialog.name]: Dialog,
  },
})
export default class SelfTask extends Vue {
  @Prop(Object) info!: any;
  // 自设任务信息
  editValue = "";
  filterList = [
    {
      key: 0,
      label: "全部",
    },
    {
      key: 1,
      label: "未完成",
    },
    {
      key: 2,
      label: "已完成",
    },
  ];
  // 查询字段
  filterKey = 0;
  minDate = new Date(2020, 0, 1);
  maxDate = new Date(2025, 10, 1);
  // 自设任务记录列表
  sowList = [];
  // 选择时间
  dateTime: any = new Date();
  // 选择时间文本
  dateTimeText: any = "";
  // 是否正在创建
  isCreated = false;
  // 是否正在请求列表
  isLoading = false;
  meta = {
    page: 1,
    size: 10,
    total: 0,
    isMore: true,
  };
  isOpenTipsModal = false;
  sowId = null;
  isDelete = false;
  // 显示时间picker
  showTimePicker = false;
  mounted() {
    this.getList();
    console.log("自设任务info", this.info);
  }
  /**
   * 获取自设任务记录
   */
  async getList(page = true) {
    const { meta, filterKey, isLoading, sowList } = this;
    const work_external_user_id = this.info.work_external_user_id;
    // 正在请求
    if (isLoading || !meta.isMore) return;
    this.isLoading = true;
    // 分页
    if (page) {
      this.meta.page = meta.page++;
    }

    const params = {
      page: meta.page,
      size: meta.size,
      status: filterKey,
      work_external_user_id,
    };
    try {
      const { code, data, message }: any = await getSowTaskListApi(params);
      this.isLoading = false;

      if (code !== 200) {
        this.$toast(message);
        return;
      }

      const { total } = data;
      const list = data.list || [];

      console.log(list, "list");

      const newSowList = page ? sowList.concat(list) : list;
      // 修改remind_time
      newSowList.forEach((item: any) => {
        item.remind_time = dayjs(item.remind_time * 1000).format(
          "YYYY-MM-DD HH:mm:ss"
        );
      });
      this.meta.total = total;
      this.meta.isMore = newSowList.length < total;
      this.sowList = newSowList;

      console.log(data);
    } catch (error) {
      console.log(error);
      this.isLoading = false;
      this.meta.page = Math.max(1, meta.page - 1);
    }
  }

  cdpSplitEvent() {}
  /**
   * 创建自设任务
   */
  async createSow() {
    const { dateTime, editValue, isCreated } = this;
    // 好友自设任务
    if (!this.info.isRelation) {
      this.$toast("该用户非您好友，暂无法设置自设任务");
      return;
    }

    if (isCreated) return;

    if (!dateTime) {
      this.$toast("请选择创建时间");
      return;
    }

    if (!editValue) {
      this.$toast("请填写自设任务");
      return;
    }
    let dateTimeFmt = dayjs(dateTime).format("YYYY-MM-DD HH:mm:ss");
    if (dayjs().isAfter(dateTimeFmt)) {
      this.$toast("请选择当前时间之后的日期");
      return;
    }

    this.isCreated = true;

    const work_external_user_id = this.info.work_external_user_id;

    const params = {
      work_external_user_id,
      content: editValue,
      remind_time: new Date(dateTime).getTime() / 1000,
    };

    try {
      const { code, data, message }: any = await createSowTaskApi(params);
      // this.dateTime = 0;
      this.editValue = "";
      this.isCreated = false;
      if (code !== 200) {
        this.$toast(message);
        return;
      }

      if (data) {
        this.$toast("创建成功");
        this.$emit("createSowSuccess");
        this.meta = {
          page: 1,
          size: 10,
          total: 0,
          isMore: true,
        };
        this.filterKey = 0;
        this.dateTime = null;
        this.dateTimeText = "";
        this.editValue = "";
        this.getList(false);
        this.$emit("updateCount");
        // 神策数据上报 -- 创建播种
        this.$track.trackReminderCreate({
          sow_task_id: data,
          sow_task_content: editValue,
          user_external_id: this.info.external_userid,
          // TODO 神策数据上报 入口判断 侧边栏/顾客详情
          entrance: 1
        })
      }
    } catch (error) {
      console.log(error);
      // this.dateTime = 0;
      this.editValue = "";
      this.isCreated = false;
    }
  }
  /**
   * 弹窗点击确认
   */
  confirm() {
    if (this.isDelete) {
      this.deleteSow();
    } else {
      this.confirmFinished();
    }
  }
  /**
   * 删除自设任务
   */
  async deleteSow() {
    const sowId = this.sowId;
    const work_external_user_id = this.info.work_external_user_id;

    try {
      const { data, code, message }: any = await deleteSowTaskApi({
        id: sowId,
        work_external_user_id,
      });

      if (code !== 200) {
        this.$toast(message);
        return;
      }

      if (data.status === "OK") {
        this.closeTipsModel();
        this.$toast("删除成功");
        this.meta = {
          page: 1,
          size: 10,
          total: 0,
          isMore: true,
        };
        this.filterKey = 0;
        this.sowList = [];
        this.getList(false);
        this.$emit("updateCount");
      }
    } catch (error) {
      console.log(error);
    }
  }
  /**
   * 结束自设任务
   */
  async confirmFinished() {
    const { sowId } = this;
    try {
      const { data, code, message }: any = await finishedSowTaskApi({
        id: sowId,
      });

      if (code !== 200) {
        this.$toast(message);
        return;
      }

      if (data.status === "OK") {
        this.closeTipsModel();
        this.$toast("自设任务完成");
        this.meta = {
          page: 1,
          size: 10,
          total: 0,
          isMore: true,
        };
        this.filterKey = 0;
        this.getList(false);
      }
    } catch (error) {
      console.log(error);
    }
  }
  /**
   * 时间选择
   */
  timeConfirm() {
    this.dateTimeText = dayjs(this.dateTime).format("YYYY-MM-DD HH:mm:ss");
    this.showTimePicker = false;
  }
  /**
   * 点击筛选条件
   */
  tapFilter(e: any) {
    this.filterKey = e.key;
    this.sowList = [];
    this.meta = {
      page: 1,
      size: 10,
      total: 0,
      isMore: true,
    };
    this.getList();
  }
  /**
   * 关闭弹窗
   */
  closeTipsModel() {}
  /**
   * 打开确认弹窗
   */
  openTipsModel(id) {
    Dialog.confirm({
      title: "提醒",
      message: "确认已完成该自设任务？",
    }).then(() => {
      // on confirm
      this.sowId = id;
      this.isDelete = false;
      this.confirmFinished();
    });
  }
  /**
   * 打开删除弹窗
   */
  openDeleteModel(id) {
    Dialog.confirm({
      title: "提醒",
      message: "确认删除当前自设任务？",
    }).then(() => {
      // on confirm
      this.sowId = id;
      this.isDelete = true;
      this.deleteSow();
    });
  }
  /**
   * 无法创建自设任务自设任务时点击事件
   */
  tapDisableEdit() {
    if (!this.info.isRelation) {
      this.$toast("该用户非您好友，暂无法设置自设任务");
    }
  }
}
</script>
<style lang="sass" scoped>
.sow
  padding: 15px 15px

  .new-sow
    margin-bottom: 15px
    padding: 15px 15px 25px
    background-color: #fff

    .top-box
      display: flex
      justify-content: space-between
      align-items: center

      .title
        font-size: 16px
        font-family: PingFangSC-Medium, PingFang SC
        font-weight: 500
        color: #000000

      .time-select
        width: 170px
        height: 30px
        background-color: #F6F6F9
        display: flex
        align-items: center
        padding: 0 5px

        .icon-1
          width: 15px
          height: 15px
          margin-right: 5px

        .placeholder
          flex: 1
          font-size: 12px
          font-family: PingFangSC-Regular, PingFang SC
          font-weight: 400
          color: #96979E
          margin-left: 2px

        .icon-2
          width: 14px
          height: 14px
          transform: rotate(-90deg)

    .edit-box
      margin-top: 15px
      margin-bottom: 25px
      background-color: #F6F6F9
      position: relative

      .edit
        padding: 10px
        font-size: 14px
        font-family: PingFangSC-Regular, PingFang SC
        font-weight: 400
        color: #000000
        background-color: #F6F6F9

      .word-total
        font-size: 12px
        font-family: PingFangSC-Regular, PingFang SC
        font-weight: 400
        color: #BDBEC8
        display: flex
        justify-content: flex-end
        position: absolute
        bottom: 10px
        right: 10px

    .btn
      height: 42px
      display: flex
      justify-content: center
      align-items: center
      background-color: #000000
      color: #ffffff
      font-size: 14px
      font-family: PingFangSC-Regular, PingFang SC
      font-weight: 400

      &.disabled
        background-color: #cccccc

  .sow-record
    background-color: #fff
    padding: 15px

    .record-top
      display: flex
      align-items: center
      justify-content: space-between

      .title
        font-size: 16px
        font-family: PingFangSC-Medium, PingFang SC
        font-weight: 500
        color: #000000

      .filter-list
        display: flex
        align-items: center

        .filter-item
          height: 22px
          padding: 0 10px
          background-color: #fff
          font-size: 11px
          font-family: PingFangSC-Medium, PingFang SC
          font-weight: 500
          color: #404040
          display: flex
          align-items: center
          justify-content: center
          border: 1px solid #404040
          margin-left: 10px

          &.selected
            color: #ffffff
            background-color: #000000
            border: none

    .bottom-tip, .loading-box
      width: 100%
      font-size: 14px
      color: #96979E
      text-align: center
      margin: 15px 0

    .record-item
      margin-top: 15px
      background-color: #F6F6F9
      display: flex
      flex-direction: column
      align-items: flex-end
      padding: 10px 15px
      height: 135px

      .top-info
        width: 100%
        display: flex
        align-items: center

        .icon
          width: 14px
          height: 14px
          margin-right: 10px

        .time
          flex: 1
          font-size: 12px
          font-family: PingFangSC-Medium, PingFang SC
          font-weight: 500
          color: #404040

        .status
          font-size: 12px
          font-family: PingFangSC-Medium, PingFang SC
          font-weight: 500
          color: #CFB888

          &.finished
            color: #82D17B

      .detail
        flex: 1
        width: 100%
        font-size: 14px
        font-family: PingFangSC-Regular, PingFang SC
        font-weight: 400
        color: #000000
        margin: 12px 0
        word-break: break-all
        overflow-y: scroll

      .btn-box
        display: flex
        justify-content: flex-end

        .delete-btn
          width: 64px
          height: 28px
          border: 1px solid #4A4A4A
          font-size: 12px
          font-family: PingFangSC-Regular, PingFang SC
          font-weight: 400
          color: #4A4A4A
          display: flex
          align-items: center
          justify-content: center
          margin-right: 10px

        .confirm-btn
          width: 64px
          height: 28px
          border: 1px solid transparent
          font-size: 12px
          background-color: #000
          font-family: PingFangSC-Regular, PingFang SC
          font-weight: 400
          color: #fff
          display: flex
          align-items: center
          justify-content: center

    .empty-list
      height: 233px
      display: flex
      flex-direction: column
      justify-content: center
      align-items: center

      .image
        width: 86px
        height: 86px

      .word
        font-size: 14px
        font-family: PingFangSC-Regular, PingFang SC
        font-weight: 400
        color: #4A4A4A
        margin-top: 15px
</style>
