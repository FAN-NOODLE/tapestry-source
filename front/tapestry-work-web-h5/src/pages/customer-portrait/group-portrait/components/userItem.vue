<template>
  <div>
    <!--  && list.length > 0 -->
    <div :class="{ 'activites-list': type && type != 2 && list.length > 0 }">
      <div class="item-list" v-for="item in list" :key="item.external_userid">
        <img
          class="item-img"
          :src="
            item.avatar
              ? item.avatar
              : require('@/assets/default/default-person.png')
          "
        />
        <div class="item-right">
          <div class="item-name-box">
            <div class="item-name">{{ item.nick || item.name || "-" }}</div>
            <div class="item-form">
              <!-- 查看参加活动 -->
              <div
                v-if="type != 6"
                :class="['item-type', { underLine: type == 5 || type == 2 }]"
                @click="jumpUrl(item)"
              >
                {{
                  type == 5 || type == 2
                    ? "查看参加活动"
                    : item.user_id === ownerId
                    ? "群主"
                    : item.join_scene == 3
                    ? "扫码入群"
                    : "直接邀请入群"
                }}
              </div>
              <div class="line">|</div>
              <div class="item-form-content">
                {{
                  type == 5 || type == 2 ? item.act_num || 0 : item.time || "-"
                }}
              </div>
            </div>
          </div>
          <div
            class="item-btns"
            :style="{
              'justify-content': type == 6 ? 'center' : 'space-between',
            }"
          >
            <!-- 联系他、添加好友 -->
            <div
              v-if="item.user_id && item.user_id !== userId"
              :class="['item-btn', { underLine: type != 5 }]"
              @click="addUser(item)"
            >
              {{ item.friend ? "联系他" : "添加好友" }}
            </div>
            <div class="item-btn-else" v-else-if="item.user_id !== ownerId">
              暂无添加员工
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Inject } from "vue-property-decorator";
@Component({
  components: {},
})
export default class extends Vue {
  @Prop({ default: [], type: Array }) list!: Array<any>;
  @Prop({ default: 0, type: Number }) type!: number;
  @Prop({ default: 0, type: Number }) total!: number;
  @Prop({ default: "", type: String }) ownerId!: string;

  public get userId(): string {
    return this.$store.getters["user/getUserInfo"].user_id;
  }
  private tit = "";
  private created() {
    if (this.type) {
      this.tit = document.title;
      if (this.type == 6) this.tit = "白金卡及以下卡等未入群的管理用户";
    }
  }
  public addUser({ type, user_id, friend }) {
    // 添加好友
    if (!user_id) return this.$toast("userid不存在");
    if (friend === 1) this.openEnterpriseChat(user_id, type);
    else
      window.wx.invoke(
        "openUserProfile",
        {
          type, //1表示该userid是企业成员，2表示该userid是外部联系人
          userid: user_id, //可以是企业成员，也可以是外部联系人
        },
        (res) => {
          console.log("添加好友>>>", res);
          if (res.err_msg != "openUserProfile:ok") this.$toast(res.err_msg);
        }
      );
  }
  public openEnterpriseChat(userIds, type) {
    // 打开会话
    interface p {
      userIds?: "";
      externalUserIds?: "";
      groupName?: "";
      chatId?: "";
    }
    const params: p = {
      groupName: "", // 必填，会话名称。单聊时该参数传入空字符串""即可。
      chatId: "", //新建会话时，chatId必须为空串
    };
    if (type === 1) params.userIds = userIds;
    //参与会话的企业成员列表，格式为userid1;userid2;...，用分号隔开。
    else params.externalUserIds = userIds; // 参与会话的外部联系人列表，格式为userId1;userId2;…，用分号隔开。
    window.wx.invoke("openEnterpriseChat", params, (res) => {
      console.log("打开会话>>>", res);
      if (res.err_msg != "openEnterpriseChat:ok") this.$toast(res.err_msg);
    });
  }

  private jumpUrl(item) {
    if (this.type == 5 || this.type == 2) {
      this.$router.push({
        path: "/join-activites?user_id=" + item.user_id,
      });
    }
  }
}
</script>
<style lang="scss" scoped>
.activites-list {
  margin: 15px;
  background: #ffffff;
  box-shadow: 0px 2px 19px 0px rgba(235, 228, 217, 0.39);
  border-radius: 12px;
  border: 1px solid #ffffff;
  padding: 0 15px;
}
.item-list {
  display: flex;
  background: #ffffff;
  flex-shrink: 0;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
  .item-index {
    font-size: 12px;
    padding: 0 12px;
    margin-bottom: 10px;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  .item-img {
    margin-right: 8px;
    display: block;
    width: 44px;
    height: 44px;
    border-radius: 100%;
  }
  .item-right {
    flex: 1;
    display: flex;
    justify-content: space-between;
    align-items: center;
    .item-name-box {
      display: flex;
      justify-content: space-between;
      flex-direction: column;
      .item-name {
        overflow: hidden;
        width: 168px;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-weight: 500;
        line-height: 20px;
        font-size: 14px;
      }
    }
    .item-btns {
      flex: 1;
      display: flex;
      align-items: center;
      height: 100%;
      margin-left: 12px;
      .item-btn {
        width: 56px;
        height: 20px;
        border-radius: 4px;
        border: 1px solid #5d5d5d;
        line-height: 20px;
        text-align: center;
        font-size: 12px;
        font-weight: 400;
      }
      .item-btn-else {
        font-size: 12px;
        line-height: 20px;
        color: #a1a1a1;
        cursor: pointer;
        font-weight: bold;
      }
    }
    .item-form {
      display: flex;
      align-items: center;
      margin-top: 4px;
      font-size: 12px;
      line-height: 17px;
      color: #a1a1a1;
      min-width: 190px;
      div {
        transform: scale(0.95);
      }
      .item-type {
        font-size: 12px;
        line-height: 17px;
        &.underLine {
          color: #caa361;
          text-decoration: underline;
        }
      }
      .line {
        margin: 0 5px;
      }
      .activity {
        font-size: 14px;
      }
    }
  }
}
.item-list:last-child {
  border: 0;
}
</style>
