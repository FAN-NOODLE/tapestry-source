<template>
  <div id="signPage">
    <div class="activity-detail-info">
      <div class="detail-info-body">
        <div v-for="(item, index) in moduleContent" :key="index">
          <!-- 单选 -->
          <div v-if="item.type == 1" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div
              class="custom-radio"
              v-for="(ite, inx) in item.options"
              :key="inx"
            >
              <p>
                <img src="@/assets/marketing-campaign/no-choice.png" alt="" />
                <span>{{ ite.option }}</span>
              </p>
              <div class="option-img">
                <img
                  v-for="(itm, ind) in ite.pic"
                  :key="ind"
                  :src="itm"
                  alt=""
                  :style="{
                    marginRight: indexArr.indexOf(ind) == -1 ? '10px' : '0px',
                  }"
                />
              </div>
            </div>
            <div class="custom-radio" v-if="item.isOther">
              <p>
                <img src="@/assets/marketing-campaign/no-choice.png" alt="" />
                <span>其他</span>
              </p>
            </div>
          </div>
          <!-- 单行文本输入 -->
          <div v-if="item.type == 2" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="input-style">
              <span v-if="item.min === '' || item.max === ''" class="hint-text"
                >请输入</span
              >
              <span v-else class="hint-text"
                >请输入{{ item.min }}-{{ item.max }}字符</span
              >
            </div>
          </div>
          <!-- 多选 -->
          <div v-if="item.type == 3" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div
              class="custom-radio"
              v-for="(ite, inx) in item.options"
              :key="inx"
            >
              <p>
                <img src="@/assets/marketing-campaign/no-choice.png" alt="" />
                <span>{{ ite.option }}</span>
              </p>
              <div class="option-img">
                <img
                  v-for="(itm, ind) in ite.pic"
                  :key="ind"
                  :src="itm"
                  alt=""
                  :style="{
                    marginRight: indexArr.indexOf(ind) == -1 ? '10px' : '0px',
                  }"
                />
              </div>
            </div>
            <div class="custom-radio" v-if="item.isOther">
              <p>
                <img src="@/assets/marketing-campaign/no-choice.png" alt="" />
                <span>其他</span>
              </p>
            </div>
            <p class="vice-title-p" v-if="item.min !== '' && item.max !== ''">
              选择{{ item.min }}到{{ item.max }}项
            </p>
          </div>
          <!-- 多行文本输入 -->
          <div v-if="item.type == 4" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div
              class="input-style"
              style="height: 85px; align-items: flex-start"
            >
              <span v-if="item.min === '' || item.max === ''" class="hint-text"
                >请输入</span
              >
              <span v-else class="hint-text"
                >请输入{{ item.min }}-{{ item.max }}字符</span
              >
            </div>
          </div>
          <!-- 上传图片 -->
          <div v-if="item.type == 5" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="upload-img">
              <van-icon name="plus" size="24" color="#808080" />
            </div>
            <p class="vice-title-p" v-if="item.min !== '' && item.max !== ''">
              上传{{ item.min }}到{{ item.max }}张
            </p>
          </div>
          <!-- 下拉选择 -->
          <div v-if="item.type == 6" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="input-style">
              <span class="hint-text">请选择</span>
              <van-icon name="arrow-down" size="16" color="#808080" />
            </div>
          </div>
          <!-- 时间选择 -->
          <div v-if="item.type == 7" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="input-style">
              <span class="hint-text">请选择</span>
              <van-icon name="arrow" size="16" color="#808080" />
            </div>
          </div>
          <!-- 输入姓名 -->
          <div v-if="item.type == 8" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="input-style">
              <span class="hint-text">请输入您的真实姓名</span>
            </div>
          </div>
          <!-- 性别单选 -->
          <div v-if="item.type == 9" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div
              class="custom-radio"
              v-for="(ite, inx) in item.options"
              :key="inx"
            >
              <p>
                <!-- <img
                  v-if="index == 0"
                  src="@/assets/marketing-campaign/choice.png"
                  alt=""
                /> -->
                <img src="@/assets/marketing-campaign/no-choice.png" alt="" />
                <!-- <span v-if="index == 0" style="color: #caa361">{{
                    item.option
                  }}</span> -->
                <span>{{ ite.option }}</span>
              </p>
            </div>
          </div>
          <!-- 手机号 -->
          <div v-if="item.type == 10" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="input-style">
              <span class="hint-text">请输入您的手机号</span>
            </div>
          </div>
          <!-- 邮箱 -->
          <div v-if="item.type == 11" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="input-style">
              <span class="hint-text">请输入您的邮箱地址</span>
            </div>
          </div>
          <!-- 地址 -->
          <div v-if="item.type == 12" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="select-style">
              <div class="input-style" style="width: 30%">
                <span class="hint-text">选择省</span>
                <div class="select-icon"></div>
              </div>
              <div class="input-style" style="width: 30%">
                <span class="hint-text">选择市</span>
                <div class="select-icon"></div>
              </div>
              <div class="input-style" style="width: 30%">
                <span class="hint-text">选择区</span>
                <div class="select-icon"></div>
              </div>
            </div>
            <div class="input-style" v-if="item.isAddress">
              <span class="hint-text">请输入详细地址</span>
            </div>
          </div>
          <!-- 邮编 -->
          <div v-if="item.type == 13" class="widget">
            <!-- <p class="title-p">{{ item.sequence }}.{{ item.name }}</p> -->
            <p class="title-new">
              <span v-if="item.isRequired" class="title-one">*</span>
              <span class="title-two">{{ item.sequence }}</span>
              <!-- <b class="title-b" style="margin-right: 4px">/</b> -->
              {{ item.tit }}
              <!-- <b class="title-b" style="margin-left: 4px">/</b> -->
              <!-- <b class="title-line"></b> -->
            </p>
            <p class="vice-title-p">{{ item.mark }}</p>
            <div class="input-style">
              <span class="hint-text">请输入您所在区域的邮政编码</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import { Icon } from "vant";
@Component({
  components: {
    [Icon.name]: Icon,
  },
})
export default class extends Vue {
  @Prop({ default: [] }) moduleContent: Array<any>;
  private indexArr = [2, 5, 8];
}
</script>
<style lang="scss" scoped>
::v-deep {
  .activity-detail-info {
    font-size: 16px;
    color: #33363c;
    line-height: 22px;
  }
}
#signPage {
  width: 100%;
  .activity-detail-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    .detail-info-body {
      width: 92%;
      .widget {
        padding: 15px;
        margin: 0px 0px 15px;
        background: #ffffff;
        box-shadow: 0px 9px 16px 0px rgba(0, 0, 0, 0.05);
        border-radius: 12px;
        .custom-radio {
          p {
            display: flex;
            align-items: flex-start;
            font-size: 16px;
            color: #33363c;
            line-height: 22px;
            margin: 10px 0px;
            img {
              width: 16px;
              height: 16px;
              margin-right: 11px;
              margin-top: 3px;
            }
          }
          .option-img {
            img {
              width: 98px;
              height: 98px;
              border-radius: 12px;
            }
          }
        }
        .title-p {
          font-size: 16px;
          font-weight: bold;
          color: #33363c;
          line-height: 26px;
        }
        // 新的标题样式
        .title-new {
          display: flex;
          align-items: center;
          font-size: 16px;
          font-family: PingFangSC-Medium, PingFang SC;
          color: #33363c;
          line-height: 26px;
          font-weight: bold;
          .title-one {
            font-size: 18px;
            font-family: YouSheBiaoTiHei;
            color: #fa4046;
            line-height: 23px;
            margin-right: 3px;
          }
          .title-b {
            font-size: 20px;
            color: #33363c;
            font-weight: 500;
          }
          .title-line {
            width: 60px;
            height: 0px;
            border-bottom: 1px dashed #33363c;
            margin: 0px 7px;
          }
          .title-two {
            font-size: 22px;
            // font-family: DINAlternate-Bold, DINAlternate;
            font-weight: bold;
            color: #caa361;
            line-height: 26px;
            margin-right: 3px;
          }
        }
        .vice-title-p {
          font-size: 14px;
          color: #808080;
          line-height: 24px;
        }
        .upload-img {
          width: 98px;
          height: 98px;
          background: #f6f6f9;
          border-radius: 10px;
          margin-top: 8px;
          margin-bottom: 6px;
          display: flex;
          justify-content: center;
          align-items: center;
        }
        .hint-text {
          font-size: 14px;
          color: #96979e;
          line-height: 20px;
        }
        .input-style {
          width: 315px;
          height: 44px;
          background: #f6f6f9;
          border-radius: 10px;
          margin-top: 8px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 10px;
        }
        .select-style {
          width: 100%;
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        .select-icon {
          width: 0;
          height: 0;
          border-right: 6px solid transparent;
          border-left: 6px solid transparent;
          border-top: 7px solid #96979e;
        }
      }
    }
  }
}
</style>
