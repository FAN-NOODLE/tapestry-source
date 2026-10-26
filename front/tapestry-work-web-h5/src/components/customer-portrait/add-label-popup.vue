<!--添加企业标签弹窗-->
<template>
  <div class="label-popup">
    <van-popup
      v-model="showLabelPopup"
      class="popup-box"
      position="bottom"
      closeable
      @click-overlay="closePopup"
      @click-close-icon="closePopup"
    >
      <div class="popup-title">添加企业标签</div>
      <div class="label-list-box">
        <div class="item-group" v-for="(item, index) in allTag" :key="index">
          <div class="item-group-name">{{ item.tag_group_name }}</div>
          <div
            class="item-label"
            v-for="(itemLabel, indexLabel) in item.tags"
            :key="indexLabel"
            @click="() => changeTagStatus(item.tag_group_id, itemLabel.qw_id)"
          >
            <img
              class="item-check"
              :src="
                itemLabel.is_choice
                  ? require('@/assets/customer-portrait/is_check.png')
                  : require('@/assets/customer-portrait/no_check.png')
              "
              alt=""
            />
            <div class="label-name">{{ itemLabel.name }}</div>
          </div>
        </div>
      </div>
      <div class="border-line"></div>
      <div class="save-box">
        <div class="save-btn" @click="setTag">保存</div>
      </div>
    </van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Emit, Prop } from "vue-property-decorator";
import { Popup } from "vant";
import { tagListApi, setTagApi } from "@/api/customer-portrait";
@Component({
  components: {
    [Popup.name]: Popup,
  },
})
export default class extends Vue {
  @Prop({ default: 0, type: Number }) relationId!: number;
  private showLabelPopup = true;
  private groupList = [1, 2, 3, 4, 5, 6]; // 添加的员工列表
  private allTag: any = []; // 企业标签库
  private addTag: Array<string> = []; // 新增的标签IDS（不包含现有的）
  private removeTag: Array<string> = []; // 删除的标签IDS
  private initialSelectTag: Array<string> = []; // 原有的已选中的标签
  private currentSelectTag: Array<string> = []; // 最终的已选中的标签
  @Emit("close")
  private closePopup(type?: string) {
    return type;
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    this.tagList();
  }
  /*************************** method start ***************************/
  // 标签列表接口
  private async tagList() {
    const { code, message = "", data }: IApiResult = await tagListApi(
      this.relationId
    );
    if (code === 200) {
      let allTag = data.manual;
      allTag.forEach((item: any) => {
        item.tags.forEach((itemTag: any) => {
          if (itemTag.is_choice) {
            this.initialSelectTag.push(itemTag.qw_id);
          }
        });
      });
      console.log("initialSelectTag", this.initialSelectTag);
      this.currentSelectTag =
        this.initialSelectTag.length > 0
          ? JSON.parse(JSON.stringify(this.initialSelectTag))
          : [];
      this.allTag = allTag;
    } else {
      this.$toast.fail(message);
    }
  }
  // 设置客户标签列表接口
  private async setTag() {
    // 过滤出新增的
    this.addTag = this.currentSelectTag.filter(
      (item: string) => this.initialSelectTag.indexOf(item) == -1
    );
    // 过滤出删除的
    this.removeTag = this.initialSelectTag.filter(
      (item: string) => this.currentSelectTag.indexOf(item) == -1
    );
    console.log(this.addTag, this.removeTag);
    // 无修改，直接关闭弹框
    if (this.addTag.length == 0 && this.removeTag.length == 0) {
      this.closePopup();
      return;
    }
    const params = {
      relation_id: this.relationId,
      add_tag: this.addTag,
      remove_tag: this.removeTag,
    };
    const { code, message = "" }: IApiResult = await setTagApi(params);
    if (code === 200) {
      this.$toast.success({
        duration: 1000,
        message: "保存成功",
        className: "toast-one-row",
      });
      this.closePopup("refresh");
    } else {
      this.$toast.fail(message);
    }
  }
  // 点击选择标签弹框选择标签
  private changeTagStatus(tag_group_id, qw_id) {
    let copy_data = JSON.parse(JSON.stringify(this.allTag));
    copy_data.map((item: any) => {
      if (item.tag_group_id === tag_group_id) {
        item.tags.forEach((itemTag: any) => {
          if (itemTag.qw_id === qw_id) {
            itemTag.is_choice = !itemTag.is_choice;
            itemTag.is_choice
              ? this.currentSelectTag.push(qw_id)
              : this.currentSelectTag.splice(
                  this.currentSelectTag.findIndex((v) => v === qw_id),
                  1
                );
          }
        });
      }
    });
    this.allTag = copy_data;
  }
}
</script>
<style lang="scss" scoped>
.label-popup {
  .popup-box {
    overflow-x: hidden;
    border-radius: 4px 4px 0 0;
    ::v-deep .van-popup__close-icon {
      color: #808080;
      font-size: 14px;
      z-index: 2;
    }
    .popup-title {
      height: 50px;
      line-height: 50px;
      width: 100%;
      text-align: center;
      color: #272727;
      font-size: 16px;
      font-weight: 600;
      border-bottom: solid 1px #ebebeb;
      position: sticky;
      top: 0;
      z-index: 2;
    }
    .label-list-box {
      padding: 15px 0 0 15px;
      height: 280px;
      overflow-y: auto;
      z-index: 1;
      .item-group {
        margin-bottom: 10px;
        .item-group-name {
          font-size: 14px;
          line-height: 20px;
          color: #272727;
          font-weight: 600;
          margin-bottom: 9px;
        }
        .item-label {
          display: inline-flex;
          border: solid 1px #b9bec4;
          text-align: center;
          border-radius: 4px;
          height: 28px;
          align-items: center;
          justify-content: center;
          padding: 0 15px;
          margin-right: 20px;
          margin-bottom: 10px;
          .item-check {
            width: 12px;
            height: 12px;
          }
          .label-name {
            font-size: 12px;
            color: #313131;
            margin-left: 8px;
            max-width: 288px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }
      }
    }
    .border-line {
      width: calc(100vw - 15px);
      height: 1px;
      background: #ebebeb;
      position: relative;
      left: 15px;
    }
    .save-box {
      height: 70px;
      position: relative;
      .save-btn {
        width: 100px;
        height: 28px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #ffffff;
        background: #4A4A4A;
        font-size: 14px;
        position: absolute;
        right: 15px;
        bottom: 20px;
      }
    }
  }
}
</style>
