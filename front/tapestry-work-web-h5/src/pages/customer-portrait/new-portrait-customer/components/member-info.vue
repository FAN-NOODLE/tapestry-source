<template>
  <div class="member-info">
    <!-- 客户标签 -->
    <div v-if="isSW && info.crm_tag && info.crm_tag.length>0" class="customer-tag">
      <div class="title-box">
        <div class="title">偏好标签</div>
      </div>
      <div class="tag-list">
        <div
          v-for="(item, index) in info.crm_tag"
          :key="index"
          class="tag-item"
          :style="{
            'padding-right': '10px',
          }"
        >
          {{ item }}
        </div>
      </div>
    </div>
    <!-- 客户标签 -->
    <div class="customer-tag">
      <div class="title-box" @click="openToast(0)">
        <div class="title">客户标签</div>
        <img
          class="icon"
          src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/popover-icon.png"
        />
      </div>
      <div class="tag-list">
        <div
          v-for="(item, index) in info.tagList"
          :key="index"
          class="tag-item"
          :style="{
            'padding-right': item.third_tag_description ? '5px' : '10px',
          }"
        >
          {{ item.name }}
          <img
            v-if="item.third_tag_description"
            class="icon"
            @click="openToast(2, item.third_tag_description)"
            src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/popover-icon-white.png"
          />
        </div>
      </div>
    </div>

    <!-- 客户印象 -->
    <div class="customer-impression" v-if="!isKS">
      <div class="title-box" @click="openToast(1)">
        <div class="title">客户印象</div>
        <div class="icon-box">
          <img
            class="icon"
            src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/popover-icon.png"
          />
        </div>
      </div>
      <div class="impression-list">
        <div
          v-for="(item, index) in info.user_impression_tags"
          :key="item.group_id"
          class="item-box"
        >
          <div class="impression-item">
            <div class="top-box">
              <div class="title">{{ item.group_name }}</div>
              <img
                v-if="item.impression_tags_list.length > 0"
                @click="handleExpand(item)"
                :class="['image', item.expand ? 'icon-expand' : '']"
                src="https://cdn-cos.coach.com.cn/uploads/2022/v3.0.2/weclient/arrow-bottom.png"
              />
            </div>
            <div v-if="item.expand" class="impression-detail-list">
              <div
                v-for="tag in item.impression_tags_list"
                @click="
                  tagsImpressionOperate(tag, item.impression_tags_list, index)
                "
                :key="tag.tag_id"
                :class="['tag', tag.is_have === 1 ? 'selected' : '']"
              >
                {{ tag.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 提示文案弹窗 -->
    <van-popup v-model="toastOption.show" closeable class="popup-content">
      {{ toastOption.content }}
    </van-popup>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { Popup } from "vant";
import { impressionTagsOperateApi } from "@/api/customer-portrait";
import { InfoType } from "../index.vue";
import Helper from "@/utils/helper";
/** 提示文案 */
const DECS = [
  "客户标签仅提供推荐建议，导购无法编辑客户标签。",
  `
    为顾客添加或移除印象标签
    1、点击客户印象标签可以为顾客添加或移除印象标签
    2、已添加的标签以深色底色展示
    3、所有导购都可以为顾客添加或移除印象标签，页面展示最后一次编辑
    `,
];

@Component({
  components: {
    [Popup.name]: Popup,
  },
})
export default class MemberInfo extends Vue {
  /** 顾客详情 */
  @Prop() info!: InfoType;

  /** 外部联系人id */
  @Prop() openId!: string;

  /** 提示弹窗信息 */
  private toastOption = {
    show: false,
    content: "",
  };
  // 是否是KS
  private isKS = Helper.isKS();
  // 是否是sw
  private isSW = Helper.isSW();
  /** 打开提示弹窗 */
  private openToast(index = 0, content = "") {
    this.toastOption.content = DECS[index] || content;
    this.toastOption.show = true;
  }

  /** 展开收起 */
  private handleExpand(item) {
    item.expand = !item.expand;
  }

  /** 选中和取消标签 */
  private async tagsImpressionOperate(tag, tagArr, index) {
    const { isRelation, customer_code, user_id, work_external_user_id } = this.info;

    /** 判断是否好友 */
    if (!isRelation) {
      return this.$toast.fail("抱歉，非本人好友暂不可操作");
    }

    /** 判断是否会员 */
    if (!Helper.isSW() && !customer_code) {
      this.$toast.fail("顾客还未注册会员，无法使用标签功能");
      return;
    }

    let add_tag_ids = [];
    const del_tag_ids = [];

    tagArr.forEach((item) => {
      if (item.is_have) {
        del_tag_ids.push(item.tag_id);
      }
    });
    if (!tag.is_have) {
      add_tag_ids = [tag.tag_id];
    }

    // 设置加载状态为进行中
    this.$toast.loading({
      message: "加载中...",
      forbidClick: true,
    });
    // 组合数据
    const params = {
      work_external_user_id,
      external_userid: this.openId,
      enterprise_user_id: user_id,
      add_tag_ids: add_tag_ids,
      del_tag_ids: del_tag_ids,
    };

    try {
      const { code, message }: IApiResult = await impressionTagsOperateApi(
        params
      );

      this.$toast.clear();
      if (code !== 200) {
        return this.$toast.fail(message);
      }

      let newTags = tagArr.map((item) => {
        return {
          ...item,
          is_have: item.tag_id === tag.tag_id && !tag.is_have ? 1 : 0,
        };
      });

      this.$set(this.info.user_impression_tags, index, {
        ...this.info.user_impression_tags[index],
        impression_tags_list: newTags,
      });
    } catch (error) {
      this.$toast.clear();
      console.log(error);
    }
  }
}
</script>

<style lang="scss" scoped>
.member-info {
  //   padding: 15px 15px 0 15px;

  .popup-content {
    width: 240px;
    padding: 30px 20px;
    font-size: 14px;

    ::v-deep .van-popup__close-icon--top-right {
      position: absolute;
      top: 10px;
      right: 10px;
    }
  }
}

.customer-tag {
  padding: 15px 10px 15px 15px;
  background-color: #fff;

  .title-box {
    margin-bottom: 10px;
    display: flex;
    align-items: center;

    .title {
      font-size: 16px;
      font-family: PingFangSC-Medium, PingFang SC;
      font-weight: 500;
      color: #33363c;
      margin-right: 7px;
    }

    .icon {
      width: 20px;
      height: 20px;
    }
  }

  .tag-list {
    display: flex;
    flex-wrap: wrap;

    .tag-item {
      height: 30px;
      background: #000000;
      font-size: 14px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      color: #ffffff;
      display: flex;
      align-items: center;
      padding-left: 10px;
      padding-right: 5px;
      margin-right: 10px;
      margin-bottom: 10px;

      .icon {
        width: 20px;
        height: 20px;
        margin-left: 8px;
      }
    }
  }
}

.customer-impression {
  background-color: #fff;
  padding: 15px 15px;
  margin-top: 10px;

  .title-box {
    display: flex;
    align-items: center;
    margin-bottom: 15px;

    .title {
      font-size: 16px;
      font-family: PingFangSC-Medium, PingFang SC;
      font-weight: 500;
      color: #33363c;
      margin-right: 7px;
    }

    .icon-box {
      flex: 1;
      display: flex;
      align-items: center;
      .icon {
        width: 20px;
        height: 20px;
      }
    }

    .arrow {
      width: 14px;
      height: 14px;
      transform: rotate(180deg);
    }
  }

  .impression-list {
    background-color: #f3f4f6;
    position: relative;

    .item-box:not(:last-child)::after {
      content: "";
      display: block;
      height: 1px;
      background-image: linear-gradient(to right, #cccccc 55%, #f3f4f6 55%);
      background-size: 8px 100%;
      background-repeat: repeat-x;
      width: calc(100% - 30px);
      margin-left: 15px;
    }
  }
}

.impression-item {
  .top-box {
    padding: 15px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .title {
      font-size: 16px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      color: #33363c;
    }

    .image {
      width: 14px;
      height: 14px;
      transform: rotate(180deg);

      &.icon-expand {
        transform: rotate(0);
      }
    }
  }

  .impression-detail-list {
    padding: 0 5px 4px 15px;
    display: flex;
    flex-wrap: wrap;

    .tag {
      background-color: #fff;
      border: 1px solid #000000;
      color: #000000;
      padding: 5px 10px;
      font-size: 14px;
      font-family: PingFangSC-Regular, PingFang SC;
      font-weight: 400;
      margin-right: 10px;
      margin-bottom: 10px;

      &.selected {
        background-color: #000000;
        color: #ffffff;
      }
    }
  }
}
</style>
