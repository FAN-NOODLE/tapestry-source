<template>
  <div class="layout-list">
    <layout-card>
      <template slot="content">
        <!-- <el-dropdown @command="handleCommand" size="small">
          <el-button type="primary" size="small" class="font-14">
            接入公众号
            <i class="el-icon-arrow-down el-icon--right"></i>
          </el-button>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="unsame"
              >非同一主体公众号接入</el-dropdown-item
            >
            <el-dropdown-item command="same"
              >同一主体公众号接入</el-dropdown-item
            >
          </el-dropdown-menu>
        </el-dropdown> -->
        <div
          class="overflow-hidden "
          :class="officialAccountList.length > 0 ? '' : 'height-500'"
        >
          <div v-if="officialAccountList.length > 0">
            <div
              v-for="(item, index) in officialAccountList"
              :key="index"
              class="card"
            >
              <div class="left flex-s-c">
                <div class="img-container">
                  <img :src="item.head_img" alt="" />
                </div>
                <div class="desc m-l-10">{{ item.nick_name }}</div>
              </div>
              <div class="right">
                <el-button
                  size="small"
                  class="font-14"
                  @click="unbindOfficial(item.authorizer_appid)"
                  >解除绑定</el-button
                >
              </div>
            </div>
          </div>
          <div class="m-t-130" v-else>
            <!-- <empty-content :description="unbindDesc"></empty-content> -->
            <div class="empty-content">
              <img src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/default/default_noList.png" alt="" />
              <!-- <p>{{ description }}</p> -->
              <p>当前还未绑定公众号</p>
              <el-dropdown @command="handleCommand" size="small" class="m-t-30">
                <el-button type="primary" size="small" class="font-14">
                  接入公众号
                  <i class="el-icon-arrow-down el-icon--right"></i>
                </el-button>
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item command="unsame"
                    >非同一主体公众号接入</el-dropdown-item
                  >
                  <el-dropdown-item command="same"
                    >同一主体公众号接入</el-dropdown-item
                  >
                </el-dropdown-menu>
              </el-dropdown>
            </div>
          </div>
        </div>
        <el-dialog
          title="非同一主体公众号接入"
          :visible="isSameDialog"
          width="420px"
          @close="closeUnsame"
        >
          <div class="flex-s">
            <p>
              <i class="el-icon-warning" aria-hidden="true"></i>
            </p>
            <p class="m-l-10">
              公众号与商城小程序主体不一致时，绑定后将无法针对粉丝在商城内的交互行为使用消息策略，确认是否绑定？
            </p>
          </div>
          <span slot="footer" class="dialog-footer">
            <el-button size="small" @click="closeUnsame">取消</el-button>
            <el-button size="small" type="primary" @click="jumpToWXAuth"
              >确认</el-button
            >
          </span>
        </el-dialog>
        <el-dialog
          title="解除绑定"
          :visible="unbindDialog"
          width="420px"
          @close="closeUnbindDialog"
        >
          <div class="flex-s">
            <p>
              <i class="el-icon-warning" aria-hidden="true"></i>
            </p>
            <p class="m-l-10">
              你确定要解除绑定公众号吗？解绑后将不可恢复！
            </p>
          </div>
          <span slot="footer" class="dialog-footer">
            <el-button size="small" @click="closeUnbindDialog">取消</el-button>
            <el-button size="small" type="primary" @click="unbindSave"
              >确认</el-button
            >
          </span>
        </el-dialog>
      </template>
    </layout-card>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import LayoutCard from "./components/layout-card/index.vue";
import EmptyContent from "./components/empty-content/index.vue";
import Api from "@/api/setting/setting";
import LoginAboat from "@/utils/LoginManager";

@Component({
  components: {
    LayoutCard,
    EmptyContent
  }
})
export default class extends Vue {
  private officialAccountList: Array<any> = [];
  private currentAuthId = "";
  private auth_url = "";
  private commnadType = 0; //0代表同一主体 1代表非同一主体
  private isSameDialog = false;
  private unbindDialog = false;
  private unbindDesc = "当前还未绑定公众号";
  private back_url = "";
  private handleCommand(command: string) {
    if (command === "unsame") {
      this.commnadType = 1;
      this.isSameDialog = !this.isSameDialog;
    } else {
      this.commnadType = 0;
      this.jumpToWXAuth();
    }
  }
  private getOfficialAccount() {
    const enterprise_id = LoginAboat.getLoginData().enterprise_id;
    const params = {
      enterprise_id: enterprise_id
    };
    Api.getOfficialAccount(params).then(res => {
      if (res.data) {
        this.officialAccountList = res.data.data ? res.data.data : [];
        this.auth_url = `${res.data.auth_url}${res.data.enterprise_id}`;
      }
    });
  }
  private closeUnsame(): void {
    this.isSameDialog = false;
  }
  private closeUnbindDialog(): void {
    this.unbindDialog = false;
  }
  private unbindOfficial(id: string) {
    this.unbindDialog = !this.unbindDialog;
    this.currentAuthId = id;
  }
  private jumpToWXAuth() {
    const currentUrl = window.location.href;
    window.location.href = `${this.auth_url}/${this.commnadType}?&ref=${currentUrl}`;
  }
  private unbindSave() {
    const params = {
      authorizer_appid: this.currentAuthId
    };
    Api.unbindOfficialAccount(params).then((res: any) => {
      if (res.code === 200) {
        this.$message.success("公众号解绑成功！");
        this.getOfficialAccount();
        this.unbindDialog = false;
      } else {
        this.$message.error("公众号解绑失败！");
        this.unbindDialog = false;
      }
    });
  }
  created() {
    //获取公众号列表
    this.getOfficialAccount();
  }
}
</script>
<style lang="scss" scoped>
div {
  user-select: none;
}
.card {
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f5f6fa;
  margin-bottom: 20px;
  border-radius: 4px;
  .left {
    .img-container {
      width: 50px;
      height: 50px;
      border-radius: 4px;
      overflow: hidden;
      img {
        width: 100%;
        height: 100%;
        border-radius: 4px;
      }
    }
    .desc {
      color: #606266;
      font-size: 14px;
    }
  }
}
.card:last-child {
  margin-bottom: 0px;
}
.el-icon-warning {
  font-size: 20px;
  color: #f6ad42;
}
.flex-s {
  display: flex;
  justify-content: flex-start;
}
.flex-s-c {
  display: flex;
  justify-content: flex-start;
  align-items: center;
}
.m-t-130 {
  margin-top: 130px;
}
.height-500 {
  height: 500px;
}
.overflow-hidden {
  overflow: hidden;
  overflow-y: scroll;
}

.empty-content {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  img {
    width: 124px;
    height: 124px;
  }
  p {
    margin-top: 15px;
    color: #909399;
    font-size: 14px;
  }
}
</style>
