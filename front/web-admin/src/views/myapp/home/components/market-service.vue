<template>
  <div class="market-service">
    <div class="container">
      <div
        class="card"
        v-for="item in market_service"
        :key="item.sort"
        @click="jumpTo(item.jump_url, item.title)"
      >
        <div class="card-left">
          <div class="image-con">
            <!-- <svg-icon class="svg-icon" :icon-class="item.icon_class"></svg-icon> -->
            <img
              v-if="item.icon == 'bargain-icon'"
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/ui/bargain-icon.png"
            />
            <img
              v-if="item.icon == 'limit-icon'"
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/ui/limit-icon.png"
            />
            <img
              v-if="item.icon == 'group-icon'"
              src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/ui/group-icon.png"
            />
          </div>
          <div>
            <div class="title">
              {{ item.title }}
            </div>
            <div class="desc">
              {{ item.desc }}
            </div>
          </div>
        </div>
        <div class="card-right">
          <i
            class="el-icon-arrow-right"
            aria-hidden="true"
          ></i>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import ShareModule from "@/shared/index";
import Permission from "@/store/modules/permission";
import { IKeysItem } from "@/types";
@Component
export default class NewIntroduction extends Vue {
  private keys = [];
  private market_service = [
    {
      title: "裂变玩法",
      desc: "支持裂变好友和裂变群成员，通过欢迎语和企微侧边栏推送让好友带来他的好友进一步提炼现有好友资源价值，同时提高客户活跃度及转化率",
      sort: 1,
      icon: "limit-icon",
      jump_url: "myapp_enterprise_drainage_fission-activities_list",
    },
    {
      title: "联系我",
      desc: "将多个员工添加至一个二维码内，客户扫码后随机分配员工添加，可设置无需经过验证自动通过好友申请",
      sort: 2,
      icon: "group-icon",
      jump_url: "myapp_enterprise_drainage_contact-me_contact-me-list",
    },
    // {
    //   title: "好友助力",
    //   desc: "邀请好友助力，付邮下单",
    //   sort: 3,
    //   icon_class: "friend-help",
    //   jump_url: "myapp_app_fissionactive_friend-helps"
    // }
    // {
    //   title: "砍价",
    //   desc: "邀请好友砍价后低价购买，裂变获客",
    //   sort: 4,
    //   icon_class: "bargain",
    //   jump_url: "myapp_app_fissionactive_bargain-new"
    // },
    // {
    //   title: "直播",
    //   desc: "可在私域小程序中完成直播带货",
    //   sort: 5,
    //   icon_class: "live",
    //   jump_url: "myapp_app_toolbox_live-streaming"
    // },
    {
      title: "营销日历",
      desc: "营销日历，创建导购需按指定日期完成的营销任务，统一导购营销素材的获取，提高导购工作效率。",
      sort: 2,
      icon: "bargain-icon",
      jump_url: "myapp_wisdom_marketing-calendar_list",
    },
  ];
  private findKey() {
    const permission = Permission.getPermission.filter(
      (item) => item.key === "myapp_app"
    );
    return permission.length > 0 && permission[0].children[0].children[0].key;
  }
  private jumpTo(url: any, title: any) {
    console.log(
      this.keys.some((p: IKeysItem) => {
        if (p.key === url) {
          return;
        }
      })
    );
    if (this.keys.some((p: IKeysItem) => p.key === url)) {
      ShareModule.getShared().pushState(url);
    } else {
      this.$message.warning(`你暂未开通${title}功能，请联系商务人员开通`);
    }
  }
  private getkeys() {
    console.log(localStorage.getItem("KEYS"));
    try {
      return JSON.parse(localStorage.getItem("KEYS") || "[]");
    } catch (error) {
      return [];
    }
  }
  mounted() {
    this.keys = this.getkeys();
  }
}
</script>
<style lang="scss" scoped>
.market-service {
  .container {
    display: flex;
    justify-content: flex-start;
    align-items: center;
    flex-wrap: wrap;
  }
  .card:nth-child(1),
  .card:nth-child(2),
  .card:nth-child(3) {
    margin-bottom: 20px;
  }
  .card {
    margin-right: 20px;
    width: 30%;
    cursor: pointer;
    padding: 20px;
    box-sizing: border-box;
    background: #f5f6fa;
    display: flex;
    justify-content: space-between;
    border-radius: 5px;
    &:hover {
      .title {
        color: #6f3112 !important;
        font-size: 14px;
      }
      .desc {
        color: #6f3112 !important;
        font-size: 12px;
      }
      .el-icon-arrow-right {
        color: #6f3112 !important;
      }
    }
    .card-left {
      display: flex;
      justify-content: flex-start;
      align-items: center;
      .title {
        color: #303133;
        font-size: 14px;
      }
      .desc {
        color: #909399;
        font-size: 12px;
        display: -webkit-box;
        overflow: hidden;
        text-overflow: ellipsis;
        word-break: break-all;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 2;
      }
      .image-con {
        width: 42px;
        height: 42px;
        background: #fff;
        margin-right: 15px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        img {
          width: 20px;
          height: 20px;
        }
      }
    }
    .card-right {
      display: flex;
      align-items: center;
      .el-icon-arrow-right {
        color: #909399;
      }
    }
  }
}
</style>
