<template>
  <div class="signPage">
    <div class="activity-detail-info">
      <div class="detail-info-body">
        <div class="info-content">
          <div
            class="content-list"
            v-for="(item, index) in signInfoList"
            :key="index"
          >
            <div class="top-list">
              <img
                class="top-list-img"
                :src="index === 0 ? noBg0 : index === 1 ? achieveBg0 : awaitBg0"
                alt=""
              />
              <div class="top-line" v-if="indexArr.indexOf(index) == -1"></div>

            </div>
            <div class="bottom-list">
              <div
                class="list-img"
                :style="{
                  width: '58px',
                  height: '58px',
                  borderRadius: '29px',
                  background:
                    index === 0 ? '#f6f6f9' : 'rgba(202, 163, 97, 0.05)',
                }"
              >
                <template v-if="index == 0">
                  <h3 style="color: #e6e6e6">
                    {{ item.split("-")[0] }}
                  </h3>
                  <h6 style="color: #a7a7a7">
                    {{ item.split("-")[1] }}
                  </h6>
                </template>

                <template v-if="index >= 1">
                  <h3 style="color: #ede3d2">
                    {{ item.split("-")[0] }}
                  </h3>
                  <h6 style="color: #caa361">
                    {{ item.split("-")[1] }}
                  </h6>
                </template>
              </div>
              <p class="sign-text" v-if="index === 0" style="color: #808080">
                未签到
              </p>
              <p class="sign-text" v-if="index === 1" style="color: #caa361">
                已签到
              </p>
              <p class="sign-text" v-if="index > 1" style="color: #caa361">
                立即签到
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { dateFormat } from "../../../customer/auto-market/util/index";
@Component
export default class extends Vue {
  @Prop({ default: 0 }) startTime?: number;
  @Prop({ default: 0 }) endTime?: number;
  private noBg0 = "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/no-icon.png";
  private awaitBg0 = "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/await-icon.png";
  private achieveBg0 = "https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/marketing-campaign/sign-icon.png";
  private signInfoList: Array<any> = [];
  private indexArr: Array<any> = [];
  private mounted() {
    console.log("this.startTime", this.startTime);
    console.log("this.endTime", this.endTime);

    const startTime = dateFormat(this.startTime, "yyyy-MM-dd");
    const endTime = dateFormat(this.endTime, "yyyy-MM-dd");
    const timeLead = Date.parse(endTime) - Date.parse(startTime);
    const dateLead = timeLead / 86400 / 1000;
    for (let i = 0; i <= dateLead; i++) {
      const time = Date.parse(startTime) / 1000 + i * 86400;
      const newTime = dateFormat(time, "MM-dd");
      this.signInfoList = [...this.signInfoList, newTime];
    }
    for (let i = 1; i <= 20; i++) {
      const number = i * 4 - 1;
      this.indexArr = [...this.indexArr, number];
    }
    this.indexArr.push(this.signInfoList.length - 1);
    console.log("this.signInfoList", this.signInfoList);
  }
}
</script>
<style lang="scss" scoped>
.signPage {
  background: #ffffff;
  padding-top: 6px;
  .activity-detail-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    .detail-info-body {
      width: 90%;
      .info-content {
        width: 100%;
        display: flex;
        flex-wrap: wrap;
        .content-list {
          width: 25%;
          position: relative;
          .top-list {
            position: absolute;
            left: 26px;
            display: flex;
            align-items: center;
            .top-list-img {
              width: 24px;
              height: 24px;
            }
            .top-line {
              width: 50px;
              border-bottom: 2px dotted #1da095;
              margin: 0px 2px;
            }
          }
          .bottom-list {
            .list-img {
              display: flex;
              justify-content: center;
              align-items: center;
              font-size: 16px;
              color: #ffffff;
              line-height: 25px;
              margin: 42px 9px 12px;
              position: relative;
              h3 {
                position: absolute;
                font-size: 44px;
                font-family: DINAlternate-Bold, DINAlternate !important;
                font-weight: bold;
              }
              h6 {
                position: absolute;
                bottom: 7px;
                right: 0px;
                z-index: 99;
                font-size: 22px;
                font-family: DINAlternate-Bold, DINAlternate !important;
                font-weight: bold;
              }
            }
            .sign-text {
              font-size: 14px;
              line-height: 21px;
              text-align: center;
              margin-bottom: 28px;
            }
          }
        }
      }
    }
  }
}
</style>
