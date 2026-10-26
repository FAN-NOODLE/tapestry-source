<template>
  <div class="setStratage">
    <layout-form>
      <template slot="content">
        <el-form
          ref="cform"
          :rules="strategyrules"
          :model="strategyForm"
          label-width="90px"
        >
          <p class="headline ">基础信息</p>
          <el-form-item
            label="策略名称"
            prop="strategyName"
            string
            required
            trims
          >
            <cy-input
              v-model.trim="strategyForm.strategyName"
              maxlength="30"
              minlength="1"
              placeholder="请输入策略名称"
              :show-word-limit="true"
              width-size="medium"
            ></cy-input>
          </el-form-item>
          <el-form-item
            label="客户分组"
            prop="strategyItem"
            required
            trigger="change"
          >
            <cy-select
              width-size="small"
              size="small"
              placeholder="请选择客户分组"
              v-model="strategyForm.strategyItem"
              @change="setStrategyItemInfoFn"
              :optionsList="strategyItems"
              :options="{ label: 'label', value: 'value' }"
              search-input
              :disabled="strategy_status == 3 ? true : false"
            >
            </cy-select>
            <span class="spanText">
              <span> 预计触达人数 ：{{ strategyItemInfo.total }}人</span>
              <!-- <span class="main">
                <span>
                  （ {{ strategyItemInfo.phone }}人有手机号；{{
                    strategyItemInfo.wechat
                  }}人是企微好友；{{ strategyItemInfo.visitor }}人是商城访客；{{
                    strategyItemInfo.Accounts
                  }}人是公众号粉丝 ）</span
                ></span
              > -->
            </span>
          </el-form-item>
          <el-form-item label="添加渠道" required>
            <el-dropdown @command="channelChangeFn">
              <el-button
                type="primary"
                :style="{
                  opacity:
                    selectedChannels.length == 1 || strategy_status == 3
                      ? 0.5
                      : 1
                }"
              >
                选择渠道<i
                  aria-hidden="true"
                  class="el-icon-arrow-down el-icon--right"
                ></i>
              </el-button>
              <el-dropdown-menu
                v-if="selectedChannels.length != 1 && strategy_status != 3"
                slot="dropdown"
              >
                <el-dropdown-item
                  v-for="(item, index) in channels"
                  :key="index"
                  :command="item"
                  >{{ item.name }}</el-dropdown-item
                >
              </el-dropdown-menu>
            </el-dropdown>

            <el-row
              v-for="(item, index) in selectedChannels"
              :key="index"
              class="channel-container"
            >
              <!-- 企微群发 -->
              <el-card
                class="box-card other-card"
                v-if="item.list[0].dialog == 'wechartMsgDialog'"
              >
                <div slot="header" class="clearfix">
                  <span class="title">{{ item.title }}</span>
                  <a
                    v-if="strategy_status != 3"
                    class="del-btn el-icon-close"
                    @click.prevent="channelsDelThisFn(index)"
                  >
                  </a>
                  <span>
                    <el-button
                      type="text"
                      size="small"
                      class="enterprisemicro-button"
                      @click.prevent="channelReplaceThisFun(item)"
                      >编辑</el-button
                    >
                  </span>
                </div>
                <div class="box-body">
                  <div
                    class="clearfix"
                    v-for="(ele, ind) in item.list"
                    :key="ind"
                  >
                    <!-- 区分进入渠道，创建/编辑进行展示， -->
                    <span
                      class="pull-left cloumn-name"
                      v-if="
                        strategy_id &&
                          item.title == '企微群发' &&
                          item.list.length === 1 &&
                          ele.originData.content.text != ''
                      "
                    >
                      {{
                        ele.originData.content.text
                          ? ele.originData.content.text
                          : ""
                      }}
                    </span>
                    <span
                      class="pull-left cloumn-name m-b-20"
                      v-if="
                        strategy_id &&
                          item.title == '企微群发' &&
                          item.list.length === 1
                      "
                    >
                      {{
                        ele.originData.content
                          ? ele.originData.content.applet.app_name ||
                            ele.originData.content.link.title ||
                            ele.originData.content.picture.name
                          : ""
                      }}
                    </span>
                    <span
                      class="pull-left cloumn-name"
                      v-else
                      v-show="ele.originData.name !== ''"
                    >
                      {{ ele.originData.name }}
                    </span>
                    <span
                      class="pull-left cloumn-content clearfix enterprisemicro-button"
                    >
                      <span
                        class="pull-left cloumn-types content-type"
                        v-if="
                          strategy_id &&
                            item.title == '企微群发' &&
                            item.list.length === 1
                        "
                      >
                        {{
                          ele.originData.content.text &&
                          ele.originData.content.text != ""
                            ? "文本"
                            : ""
                        }}</span
                      >
                      <span
                        class="pull-left cloumn-types"
                        v-if="
                          strategy_id &&
                            item.title == '企微群发' &&
                            item.list.length === 1
                        "
                      >
                        {{
                          ele.originData.content.picture.type === 2
                            ? "图片"
                            : ""
                        }}
                        {{
                          ele.originData.content.link.type === 3 ? "链接" : ""
                        }}
                        {{
                          ele.originData.content.applet.type === 4
                            ? "小程序"
                            : ""
                        }}
                      </span>
                      <span
                        class="pull-left cloumn-types"
                        v-else
                        v-show="ele.originData.name !== ''"
                        :class="
                          ele.contentType == '图片' ||
                          ele.contentType == '链接' ||
                          ele.contentType == '小程序'
                            ? 'm-b-20'
                            : ''
                        "
                      >
                        {{ ele.originData.name == "" ? "" : ele.contentType }}
                      </span>
                    </span>
                  </div>
                </div>
              </el-card>
              <!-- 短信 -->
              <el-card
                class="box-card other-card"
                v-if="item.list[0].dialog == 'msgDialog'"
              >
                <div slot="header" class="clearfix">
                  <span class="title">{{ item.title }}</span>
                  <a
                    v-if="strategy_status != 3"
                    class="del-btn el-icon-close"
                    @click.prevent="channelsDelThisFn(index)"
                  >
                  </a>
                </div>
                <div class="box-body">
                  <div
                    class="clearfix"
                    v-for="(ele, ind) in item.list"
                    :key="ind"
                  >
                    <span class="pull-left cloumn-name">
                      {{ ele ? ele.content : "" }}
                    </span>
                    <span class="pull-left cloumn-content clearfix">
                      <span class="pull-left cloumn-type">
                        {{ ele ? ele.contentType : "" }}</span
                      >
                      <span class="pull-left cloumn-btn clearfix">
                        <el-button
                          type="text"
                          size="small"
                          @click.prevent="
                            channelReplaceThisFn(
                              index,
                              ind,
                              ele ? ele.content : ''
                            )
                          "
                          >更换</el-button
                        >
                        <el-button
                          type="text"
                          size="small"
                          @click.prevent="channelDelThisFn(index, ind)"
                          >删除</el-button
                        >
                      </span>
                    </span>
                  </div>
                </div>
              </el-card>
              <!-- 自动发券 -->
              <el-card
                v-if="item.list[0].dialog == 'autoCouponDialog'"
                class="box-card"
                style="width: 100%;"
              >
                <div slot="header" class="clearfix">
                  <span class="title">{{ item.title }}</span>
                  <a
                    v-if="strategy_status != 3"
                    class="del-btn el-icon-close"
                    @click.prevent="channelsDelThisFn(index)"
                  >
                  </a>
                  <span>
                    <el-button
                      type="text"
                      size="small"
                      class="enterprisemicro-button"
                      @click.prevent="couponReplaceThisFun(item)"
                      >编辑</el-button
                    >
                  </span>
                </div>
                <div class="box-body" style="max-height: 300px;">
                  <el-table
                    :data="item.list[0].content.coupon_list"
                    max-height="310px"
                    ref="table"
                  >
                    <el-table-column
                      prop="coupon_stock_id"
                      label="优惠券批次号"
                      min-width="200"
                      :show-overflow-tooltip="true"
                    ></el-table-column>
                    <el-table-column
                      prop="coupon_name"
                      label="优惠券名称"
                      min-width="150"
                    >
                      <template slot-scope="{ row }">
                        <div v-if="row.coupon_name.length > 10">
                          <el-tooltip
                            effect="light"
                            :content="row.coupon_name"
                            placement="top-start"
                          >
                            <span>{{
                              row.coupon_name.slice(0, 10) + "..."
                            }}</span>
                          </el-tooltip>
                        </div>
                        <div v-else>
                          {{ row.coupon_name }}
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column
                      v-if="
                        showChannel &&
                          item.list[0]['content']['customer_channel'] == 1
                      "
                      label="客户渠道"
                      min-width="120"
                    >
                      <span>不限</span>
                    </el-table-column>
                    <el-table-column
                      v-if="
                        showChannel &&
                          item.list[0]['content']['customer_channel'] == 2
                      "
                      label="场景值渠道"
                      min-width="220"
                    >
                      <template slot-scope="{ row }">
                        <div v-if="false">
                          {{ row.coupon_name }}
                        </div>
                        <div
                          v-show="item.list[0].content.customer_channel == 2"
                        >
                          <div
                            v-show="sceneChannel && sceneChannel.length > 2"
                            class="user-LableButtonList"
                          >
                            <el-tooltip
                              effect="light"
                              :content="sceneChannel[0]"
                              placement="top-start"
                            >
                              <div class="user-lablebutton">
                                <span>{{ sceneChannel[0] }}</span>
                              </div>
                            </el-tooltip>
                            <el-tooltip
                              effect="light"
                              :content="sceneChannel[1]"
                              placement="top-start"
                            >
                              <div class="user-lablebutton">
                                <span>{{ sceneChannel[1] }}</span>
                              </div>
                            </el-tooltip>
                            <el-popover
                              v-if="showChannel"
                              placement="top-start"
                              trigger="hover"
                              popper-class="popover-Tip"
                            >
                              <el-tag
                                v-for="(item, i) in sceneChannel"
                                :key="i"
                                :disable-transitions="false"
                                class="m-b-5 m-r-5"
                              >
                                <div class="user-groupbutton">
                                  {{ item }}
                                </div>
                              </el-tag>
                              <el-button
                                size="small"
                                type="text"
                                slot="reference"
                                icon="el-icon-more"
                              ></el-button>
                            </el-popover>
                          </div>
                          <div
                            class="user-LableButtonList"
                            v-show="sceneChannel && sceneChannel.length <= 2"
                          >
                            <div
                              class="user-lablebutton"
                              v-for="(itemTags, i) in sceneChannel"
                              :key="i"
                            >
                              <el-tooltip
                                effect="light"
                                :content="itemTags"
                                placement="top-start"
                              >
                                <span>{{ itemTags }}</span>
                              </el-tooltip>
                            </div>
                          </div>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column
                      v-if="
                        showChannel &&
                          item.list[0]['content']['customer_channel'] == 3
                      "
                      prop="scope"
                      label="一级渠道"
                      min-width="150"
                      :show-overflow-tooltip="true"
                    >
                      <template>
                        <div>
                          <div v-if="channelName > 10">
                            <el-tooltip
                              effect="light"
                              :content="channelName"
                              placement="top-start"
                            >
                              <span>{{
                                channelName.slice(0, 10) + "..."
                              }}</span>
                            </el-tooltip>
                          </div>
                          <div v-else>
                            {{ channelName }}
                          </div>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column
                      v-if="
                        showChannel &&
                          item.list[0]['content']['customer_channel'] == 3
                      "
                      label="二级渠道名称"
                      min-width="270"
                    >
                      <template slot-scope="{ row }">
                        <div v-if="false">
                          {{ row.coupon_name }}
                        </div>
                        <div
                          v-show="item.list[0].content.customer_channel == 3"
                        >
                          <div
                            v-show="channelNames && channelNames.length > 2"
                            class="user-LableButtonList"
                          >
                            <el-tooltip
                              effect="light"
                              :content="channelNames[0]"
                              placement="top-start"
                            >
                              <div class="user-lablebutton">
                                <span>{{ channelNames[0] }}</span>
                              </div>
                            </el-tooltip>
                            <el-tooltip
                              effect="light"
                              :content="channelNames[1]"
                              placement="top-start"
                            >
                              <div class="user-lablebutton">
                                <span>{{ channelNames[1] }}</span>
                              </div>
                            </el-tooltip>
                            <el-popover
                              v-if="showChannel"
                              placement="top-start"
                              trigger="hover"
                              popper-class="popover-Tip"
                            >
                              <el-tag
                                v-for="(item, i) in channelNames"
                                :key="i"
                                :disable-transitions="false"
                                class="m-b-5 m-r-5"
                              >
                                <div class="user-groupbutton">
                                  {{ item }}
                                </div>
                              </el-tag>
                              <el-button
                                size="small"
                                type="text"
                                slot="reference"
                                icon="el-icon-more"
                              ></el-button>
                            </el-popover>
                          </div>
                          <div
                            class="user-LableButtonList"
                            v-show="channelNames && channelNames.length <= 2"
                          >
                            <div
                              class="user-lablebutton"
                              v-for="(itemTags, i) in channelNames"
                              :key="i"
                            >
                              <el-tooltip
                                effect="light"
                                :content="itemTags"
                                placement="top-start"
                              >
                                <span>{{ itemTags }}</span>
                              </el-tooltip>
                            </div>
                          </div>
                        </div>
                      </template>
                    </el-table-column>
                    <el-table-column
                      label="是否查券"
                      min-width="150"
                      :show-overflow-tooltip="true"
                    >
                      <template>
                        {{
                          item.list[0].content.query_coupon == 1 ? "是" : "否"
                        }}
                      </template>
                    </el-table-column>
                    <el-table-column
                      label="券状态"
                      min-width="150"
                      :show-overflow-tooltip="true"
                    >
                      <template slot-scope="{ row }">
                        <el-tag
                          :type="
                            row.start_time > parseInt(+new Date() / 1000)
                              ? ''
                              : row.end_time < parseInt(+new Date() / 1000)
                              ? 'danger'
                              : 'success'
                          "
                        >
                          {{
                            row.start_time > parseInt(+new Date() / 1000)
                              ? "未开始"
                              : parseInt(+new Date() / 1000) > row.end_time
                              ? "已失效"
                              : "已开始"
                          }}
                        </el-tag>
                      </template>
                    </el-table-column>
                    <el-table-column
                      prop="template_name"
                      label="有效期"
                      min-width="400"
                    >
                      <template slot-scope="{ row }">
                        {{ row.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
                        -
                        {{ row.end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
                      </template>
                    </el-table-column>
                  </el-table>
                </div>
              </el-card>
            </el-row>
          </el-form-item>
          <el-form-item
            label="执行方式"
            prop="runMode"
            required
            trigger="change"
            :style="{
              marginBottom:
                (selectedChannels[0] && selectedChannels[0].id != 3) ||
                !isShortTermExecution
                  ? ''
                  : 0
            }"
          >
            <el-radio-group v-model="strategyForm.runMode" class="runMode">
              <el-col
                v-if="
                  (selectedChannels[0] && selectedChannels[0].id != 3) ||
                    !isShortTermExecution
                "
                :span="24"
                :class="strategyForm.condition == 1 ? '' : 'm-t-20'"
                style="padding-top:12px"
              >
                <el-radio :disabled="strategy_status == 3" :label="2"
                  >单次执行</el-radio
                >
                <div
                  class="clearfix marginT20"
                  v-if="strategyForm.runMode == 2"
                >
                  <span class="padding-l-r-20 pull-left lineHeight40 "
                    >计划执行时间</span
                  >
                  <span class="pull-left ">
                    <el-radio-group
                      v-model="strategyForm.onceRunMode"
                      @change="onceRunModeFn"
                      class="runMode"
                    >
                      <el-col :span="24" class="block">
                        <el-radio :label="1">立即执行</el-radio>
                      </el-col>
                      <el-col :span="24" class="block margin-b-0">
                        <el-radio :label="2">定时执行</el-radio>
                        <div
                          class="inlineBlock"
                          style="position: relative;top:-11px"
                        >
                          <el-form
                            v-if="strategyForm.onceRunMode != 1"
                            :model="strategyForm.short"
                            label-width="90px"
                            ref="afrom"
                            :inline="true"
                          >
                            <el-form-item
                              prop="planTime"
                              :rules="[
                                {
                                  required:
                                    strategyForm.onceRunMode == 1
                                      ? false
                                      : true,
                                  message: '计划执行时间不能为空',
                                  trigger: ['blur', 'change']
                                }
                              ]"
                            >
                              <el-date-picker
                                v-model="strategyForm.short.planTime"
                                size="small"
                                type="datetime"
                                placeholder="选择日期时间"
                                format="yyyy-MM-dd HH:mm"
                                :picker-options="startTimeOptSetting"
                              >
                              </el-date-picker>
                            </el-form-item>
                          </el-form>
                        </div>
                      </el-col>
                    </el-radio-group>
                  </span>
                </div>
              </el-col>
            </el-radio-group>
            <div class="tip">企微群发每月可发送四次</div>
          </el-form-item>
          <el-form-item label="触达规则" prop="strategyPolicy" required>
            <el-radio-group
              v-model="strategyForm.strategyPolicy"
              class="runMode"
            >
              <el-col :span="24" class="block" style="margin-bottom:0;">
                <el-radio :label="2">通知所有渠道</el-radio>
              </el-col>
            </el-radio-group>
          </el-form-item>
          <div class="splitLine"></div>
          <!-- <p class="headline">
            异常提醒
            <span class="mark"
              >长期执行的策略将于24小时后到期或短信条数不足100条或任意可发券的张数不足100张的情况通过短信提醒</span
            >
          </p> -->
          <!-- <el-form-item
            label="提醒手机号"
            prop="phone"
            style="margin-bottom: 0;"
          >
            <cy-input
              v-model.trim="strategyForm.phone"
              maxlength="11"
              minlength="0"
              placeholder="请输入手机号码"
              width-size="medium"
            ></cy-input>
          </el-form-item> -->
        </el-form>
      </template>
      <template slot="footer">
        <el-button type="primary" @click="savePageFn" size="small"
          >保存</el-button
        >
        <el-button @click="cancelPageFn" size="small">取消</el-button>
      </template>
    </layout-form>
    <select-shop
      v-if="shopListDialog"
      :hasSelectArray="selectedshopList"
      @addShopsDialog="getSelectShopsFn"
    ></select-shop>
    <select-video
      v-if="videoListDialog"
      :hasSelectArray="selectedVideoList"
      @addVideosDialog="getSelectVideosFn"
    ></select-video>
    <select-active-stretch
      v-if="activeDialog"
      :hasSelectArray="selectedActiveList"
      :listType="Number(strategyForm.long.activityType)"
      @addActiveDialog="getSelectActiveFn"
    ></select-active-stretch>

    <select-active-page
      v-if="activePageDialog"
      :hasSelectArray="selectedPageList"
      @addPagesDialog="getSelectPagesFn"
    ></select-active-page>
    <select-msg-template
      v-if="msgDialog"
      :msgName="msgName"
      @addMsgDialog="getSelectMsgFn"
    ></select-msg-template>
    <select-wechat-msg-template
      v-if="WeChatWorkMsgDialog"
      :msgToDialog="msgToDialog"
      @addWeChatMsgDialog="getSelectWeChatWorkMsgFn"
    ></select-wechat-msg-template>
    <select-auto-coupon-template
      v-if="autoCouponDialog"
      :isused="isused"
      :sceneList="sceneList"
      :chooseInfo="chooseInfo"
      @addAutoCouponDialog="getSelectAutoCouponFn"
    ></select-auto-coupon-template>
    <replace-wechat-msg-template
      v-if="wechartMsgDialog"
      :msgToDialog="msgToDialog"
      @replaceWeChatMsgDialog="replaceWeChatWorkMsgFn"
    ></replace-wechat-msg-template>
    <back-pop v-if="isPop" @goBack="goBackFn"></back-pop>
    <edit-back-pop v-if="isEditPop" @goEditBack="goEditBackFn"></edit-back-pop>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Watch } from "vue-property-decorator";
import UploadImage from "@/components/business/upload-image/index.vue";
import API from "@/api/customer";
import draggable from "vuedraggable";
import selectShop from "./select-shop.vue";
import selectVideo from "./select-video.vue";
import selectedShops from "./selected-shops.vue";
import selectedVideos from "./selected-videos.vue";
import selectActivePage from "./select-active-page.vue";
import selectedActivePage from "./selected-active-page.vue";
import selectMsgTemplate from "./select-msg-template.vue";
import selectWechatMsgTemplate from "./select-wechat-msg-template.vue";
import selectAutoCouponTemplate from "./select-auto-coupon-template.vue";
import replaceWechatMsgTemplate from "./replace-wechat-msg-template.vue";
import selectedActiveStretch from "./selected-active-stretch.vue";
import selectActiveStretch from "./select-active-stretch.vue";
import backPop from "./back-pop.vue";
import editBackPop from "./edit-back-pop.vue";
import CustomerModule from "@/store/modules/customer";
import { limitMinDate, sendTimeDate } from "../util/index";
@Component({
  components: {
    UploadImage,
    selectShop,
    selectedShops,
    selectActivePage,
    selectedActivePage,
    draggable,
    selectMsgTemplate,
    selectWechatMsgTemplate,
    selectAutoCouponTemplate,
    replaceWechatMsgTemplate,
    selectVideo,
    selectedVideos,
    selectedActiveStretch,
    selectActiveStretch,
    backPop,
    editBackPop
  }
})
export default class extends Vue {
  get CUSTOMERGROUPS() {
    return CustomerModule.getCustomerGroups;
  }
  private strategyItemInfo = {
    total: 0,
    phone: 0,
    wechat: 0,
    visitor: 0,
    Accounts: 0
  };
  private isPop = false;
  private isEditPop = false;
  private shopListDialog = false;
  private currentPage = 1;
  private isPageOpenList = true;
  private isActiveOpenList = true;
  private activeDialog = false;
  private strategy_id = 0;
  private strategy_status = 0;
  private isLongSrategy = false; // 是否是编辑长期执行策略
  private startTimeOptSetting = {
    disabledDate(time: Date) {
      const now: number = limitMinDate();
      return time.getTime() < now;
    }
  };
  private pickerOptionsStart = {
    disabledDate: (time: Date) => {
      const endDateVal = Number(this.strategyForm.long.runTimeRange[1]);
      if (endDateVal) {
        return (
          time.getTime() > new Date(endDateVal).getTime() ||
          time.getTime() < Date.now() - 3600 * 1000 * 24
        );
      } else {
        return time.getTime() < Date.now() - 3600 * 1000 * 24;
      }
    }
  };
  private pickerOptionsEnd = {
    disabledDate: (time: Date) => {
      const beginDateVal = Number(this.strategyForm.long.runTimeRange[0]);
      if (beginDateVal) {
        return (
          time.getTime() <
            new Date(beginDateVal).getTime() - 3600 * 1000 * 24 + 1 ||
          time.getTime() < Date.now() - 3600 * 1000 * 24
        );
      } else {
        return time.getTime() < Date.now() - 3600 * 1000 * 24;
      }
    }
  };
  private conditionList = [
    {
      value: 1,
      label: "不限",
      isCoupon: true
    },
    {
      value: 2,
      label: "访问",
      isCoupon: true
    },
    {
      value: 3,
      label: "加购",
      isCoupon: false
    },
    {
      value: 4,
      label: "下单",
      isCoupon: false
    },
    {
      value: 5,
      label: "支付",
      isCoupon: false
    }
  ];
  private timeTypes = [
    {
      value: 1,
      label: "秒"
    },
    {
      value: 2,
      label: "分钟"
    },
    {
      value: 3,
      label: "小时"
    },
    {
      value: 4,
      label: "天"
    }
  ];
  private relationTypes = [
    {
      value: 1,
      label: "首页"
    },
    {
      value: 2,
      label: "商品详情页"
    },
    {
      value: 3,
      label: "直播"
    },
    {
      value: 4,
      label: "活动页"
    }
    // {
    //   value: 5,
    //   label: "活动"
    // }
  ];
  private activityTypes: Array<any> = [];
  private channels = [
    // {
    //   value: 1,
    //   name: "短信",
    //   hasChecked: false,
    //   dialog: "msgDialog"
    // },
    {
      value: 2,
      name: "企微群发",
      hasChecked: false,
      dialog: "WeChatWorkMsgDialog"
    }
    // {
    //   value: 3,
    //   name: "自动发券",
    //   hasChecked: false,
    //   dialog: "autoCouponDialog"
    // }
    // {
    //   value: 4,
    //   name: "公众号",
    //   hasChecked: false,
    //   dialog: ""
    // }
  ];
  private activityNames = [];
  private assignTypes = {
    "2": [
      {
        value: 1,
        label: "全部商品"
      },
      {
        value: 2,
        label: "指定商品"
      }
    ],
    "3": [
      {
        value: 1,
        label: "全部直播"
      },
      {
        value: 2,
        label: "指定直播"
      }
    ],
    "4": [
      {
        value: 1,
        label: "全部活动页"
      },
      {
        value: 2,
        label: "指定活动页"
      }
    ],
    "5": [
      {
        value: 1,
        label: "全部活动"
      },
      {
        value: 2,
        label: "指定活动"
      }
    ]
  };
  private strategyForm: any = {
    strategyItem: "",
    strategyName: "",
    runMode: 2,
    long: {
      runTimeRange: [],
      timeInterval: "",
      relationType: "",
      activityType: 99,
      timeType: 1
    },
    short: {
      planTime: ""
    },
    condition: 1,
    assignType: 1,
    pageType: 1,
    videoType: 1,
    activeType: 1,
    onceRunMode: 1,
    strategyPolicy: 2,
    phone: ""
  };
  private strategyItems: Array<any> = [];
  private validateMap = {
    strategyName(rule: any, value: any, callback: Function) {
      if (!value) {
        callback(new Error("策略名称不允许为空！"));
      } else {
        callback();
      }
    },
    strategyItem(rule: any, value: any, callback: Function) {
      if (!value) {
        callback(new Error("客户分组不允许为空！"));
      } else {
        callback();
      }
    },
    runTimeRange(rule: any, value: any, callback: Function) {
      if (!value || value.length < 1) {
        callback(new Error("计划执行时间不允许为空！"));
      } else {
        callback();
      }
    },
    planTime(rule: any, value: any, callback: Function) {
      if (!value) {
        callback(new Error("执行时间不允许为空！"));
      } else {
        callback();
      }
    },
    phone(rule: any, value: any, callback: Function) {
      const reg = /^1[0-9]{10}$/g;
      if (!value) {
        callback();
      } else if (!reg.test(value)) {
        callback(new Error("请输入11位正确的手机号！"));
      } else {
        callback();
      }
    }
  };
  private strategyrules = {
    strategyName: {
      validator: this.validateMap.strategyName,
      trigger: ["blur", "change", "focus"]
    },
    strategyItem: {
      validator: this.validateMap.strategyItem,
      trigger: ["blur", "change"]
    },
    long: {
      runTimeRange: {
        validator: this.validateMap.runTimeRange,
        trigger: ["blur", "change", "focus"]
      }
    },
    short: {
      planTime: {
        validator: this.validateMap.planTime,
        trigger: ["blur", "change", "focus"]
      }
    },
    phone: {
      validator: this.validateMap.phone,
      trigger: ["blur"]
    }
  };
  private sortPolicyList: Array<any> = [];
  private sortPolicyLevels = ["1", "2", "3", "4"];
  private selectedShopNum = 0;
  private selectedVideoNum = 0;
  private selectedPageNum = 0;
  private selectedActiveNum = 0;
  private channelNameById = {
    "1": "短信",
    "2": "企微群发",
    "3": "自动发券",
    "4": "公众号"
  };
  private oldStrategy: any = null;
  private selectedVideoList: Array<any> = [];
  private selectedshopList: Array<any> = [];
  private selectedPageList: Array<any> = [];
  private selectedActiveList: Array<any> = [];
  private isLongTermExecution = false; // 长期执行禁用
  private isShortTermExecution = false; // 短期执行禁用
  private isShopOpenList = true;
  private isVideoOpenList = true;
  private activePageDialog = false;
  private videoListDialog = false;
  private msgDialog = false;
  private WeChatWorkMsgDialog = false;
  private autoCouponDialog = false;
  private wechartMsgDialog = false;
  private editable = true;
  private msgToDialog = null;
  private channelReplaceIndex = -1;
  private channelReplaceFIndex = -1;
  private delayedDragging = false;
  private isDragging = false;
  private isEnterpriseMicro = false;
  private selectedChannels: Array<any> = []; //{title:"",list:[{content:"",contentType:""}]}
  private selectedChanneles: object = {};
  private channelName = "";
  private channelNames: Array<any> = [];
  private isused = [];
  private msgName = "";
  public get dragOptions() {
    return {
      animation: 0,
      group: "description",
      disabled: !this.editable,
      ghostClass: "ghost"
    };
  }
  private newChannels: Array<any> = [];
  private sceneList: Array<any> = [];
  private chooseInfo: any = {};
  private showChannel = true;
  private showTrigger = true;
  private sceneChannel: Array<any> = [];
  @Watch("isDragging")
  delayedDraggingFn(newVal: boolean) {
    if (newVal) {
      this.delayedDragging = true;
      return;
    }
    this.$nextTick(() => {
      this.delayedDragging = false;
    });
  }
  @Watch("selectedChannels", { deep: true, immediate: true })
  channelRefCustomGroup(newVal: Array<string>) {
    console.log("出发了newChannelswatch");
    if (
      newVal[0] &&
      newVal[0]["id"] == 3 &&
      newVal[0]["list"][0]["content"]["customer_channel"] == 3
    ) {
      this.showTrigger = false;
      API.channelListApi({
        enterprise_id: this.$store.state.user.USERINFO.enterprise_id || "",
        pid: 0
      }).then((res: any) => {
        if (res.code == 200) {
          res.data.list.forEach((v: any) => {
            if (
              v.indice ==
              Object.keys(
                newVal[0]["list"][0]["content"]["customer_channel_value"]
              )[0]
            ) {
              this.showChannel = false;
              this.channelName = v.name;
              API.channelListApi({
                enterprise_id:
                  this.$store.state.user.USERINFO.enterprise_id || "",
                pid: v.id
              }).then(val => {
                if (res.code == 200) {
                  this.channelNames = [];
                  val.data["list"].forEach((g: any) => {
                    if (
                      (Object.values(
                        newVal[0]["list"][0]["content"][
                          "customer_channel_value"
                        ]
                      )[0] as any).indexOf(g["indice"]) != -1
                    ) {
                      this.channelNames.push(g.name);
                    }
                  });
                  this.showChannel = true;
                }
              });
            }
          });
        }
      });
    } else if (
      newVal[0] &&
      newVal[0]["id"] == 3 &&
      newVal[0]["list"][0]["content"]["customer_channel"] == 2
    ) {
      this.showChannel = false;
      this.sceneChannel = [];
      const scene = Object.values(
        newVal[0]["list"][0]["content"]["customer_channel_value"]
      )[0];
      console.log(scene, "scene");
      this.sceneList.forEach(v => {
        if ((scene as any).indexOf(v.val) != -1) {
          this.sceneChannel.push(v.name);
        }
      });
      this.showChannel = true;
    }
    this.setStrategyItemInfoFn();
  }
  onMoveFn(context: any) {
    const { relatedContext, draggedContext } = context;
    const relatedElement = relatedContext.element;
    const draggedElement = draggedContext.element;
    return (!relatedElement || !relatedElement.fixed) && !draggedElement.fixed;
  }
  onceRunModeFn() {
    this.strategyForm.short.planTime = "";
  }
  getItemIndexFn(arr: Array<object>, index: number) {
    const ind: number = arr.findIndex((v: any) => {
      return v ? v.value == index : false;
    });
    return ind;
  }
  clearSelectedActiveFn() {
    this.selectedActiveNum = 0;
    this.selectedActiveList = [];
  }
  /**
   * @description: 策列人群
   * @策列人群:list[{group_id:0,group_name:"",info:{total:0,phone:0,wechat:0,visitor:0,Accounts:0}}]
   * @param {type}
   * @return {type}
   */
  async initStrategyItemsFn() {
    const res: any = await API.groupListApi({});
    const strategyItems: Array<any> = [];
    if (res.code == 200) {
      const customerGroup: any = res.data.custom;
      customerGroup.forEach((v: any) => {
        // v.info = {
        //   total: v.member_count,
        //   phone: v.mobile_user_count,
        //   wechat: v.qw_user_count,
        //   visitor: v.shop_user_count,
        //   Accounts: v.chat_official_account_count
        // };
        strategyItems.push({
          label: v.group_name,
          value: v.group_id,
          num: v.value
          // info: v.info || this.strategyItemInfo
        });
      });
      this.strategyItems = strategyItems;
      if (this.CUSTOMERGROUPS.id && this.$route.query.type == "1") {
        this.strategyForm.strategyItem = this.CUSTOMERGROUPS.id;
        this.setStrategyItemInfoFn(this.CUSTOMERGROUPS.id);
      }
      this.setPageDataFn();
    } else {
      this.$message.error(res.message);
    }
  }
  /*
   * @description: 活动列表初始化
   * @param {type}
   * @return {type}
   */
  async initActiveTypeListFn() {
    const res: any = await API.marketingSelectorTextApi();
    const list: Array<any> = [];
    if (res.code == 200) {
      const activityTypes: any = res.data.ActivityTypes;
      if (Object.keys({ activityTypes }).length > 0) {
        for (const key in activityTypes) {
          const element: string | null = activityTypes[key];
          list.push({
            label: element,
            value: Number(key)
          });
        }
      }
      this.activityTypes = list.reverse();
      this.strategyForm.long.activityType = 99;
    }
  }

  initPageOptionFn() {
    //初始化策略人群
    this.initStrategyItemsFn();
    this.initActiveTypeListFn();
  }
  // 获取场景值列表
  getSceneList() {
    API.sceneListApi({}).then((res: any) => {
      if (res.code == 200) {
        this.sceneList = res.data;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  cancelPageFn() {
    if (this.strategy_id) {
      this.isEditPop = true;
    } else {
      this.isPop = true;
    }
  }
  goBackFn(b: string) {
    if (b == "save") {
      this.savePageLocalFn();
    } else if (b == "cancel") {
      const localData: any = window.localStorage.getItem("STRATAGE_LOCAL_DATA");
      if (localData) {
        window.localStorage.removeItem("STRATAGE_LOCAL_DATA");
      }
      this.$router.push({
        name: "myapp_customer_auto-market_market-strategy_market-list",
        params: JSON.parse(JSON.stringify(this.$route.query))
      });
    }

    this.isPop = false;
  }
  savePageLocalFn() {
    const stratageLocalData: object = {
      selectedChannels: this.selectedChannels,
      sortPolicyList: this.sortPolicyList,
      strategyForm: this.strategyForm
    };
    window.localStorage.setItem(
      "STRATAGE_LOCAL_DATA",
      JSON.stringify(stratageLocalData)
    );
    this.$router.go(-1);
  }
  goEditBackFn(b: string) {
    if (b == "close") {
      this.isEditPop = false;
    } else {
      this.$router.go(-1);
    }
  }
  async savePageFn() {
    let flag = false;
    if (this.selectedChannels[0] && this.selectedChannels[0]["id"] == 3) {
      flag = this.selectedChannels[0]["list"][0]["content"]["coupon_list"].some(
        (v: any) => Math.round(+new Date() / 1000) > v.end_time
      );
    }
    const save = () => {
      if (this.strategyForm.runMode == 1) {
        if (!this.strategyForm.long.runTimeRange[0]) {
          this.$message.warning("计划执行开始时间不允许为空");
          return;
        }
        if (!this.strategyForm.long.runTimeRange[1]) {
          this.$message.warning("计划执行结束时间不允许为空");
          return;
        }
      }
      const runtime = sendTimeDate(this.strategyForm.long.runTimeRange[0]);
      const plantime = sendTimeDate(this.strategyForm.short.planTime);
      const ciurrectTime = sendTimeDate(new Date().getTime());
      const dt = ciurrectTime.split(" ");
      const dl = runtime.split(" ");
      const dd = plantime.split(" ");
      if (runtime && !this.isLongSrategy) {
        if (dl[0] == dt[0] && dt[1] > dl[1]) {
          this.$message.warning("计划执行时间不可早于当前时间");
          return;
        }
      }
      if (plantime && !this.isLongSrategy) {
        if (dd[0] == dt[0] && dt[1] > dd[1]) {
          this.$message.warning("计划执行时间不可早于当前时间");
          return;
        }
      }
      window.localStorage.setItem("STRATAGE_LOCAL_DATA", "");
      this.summitPageFn();
    };
    if (flag) {
      this.$confirm(
        "优惠券列表中含有已失效的优惠券, 请确认是否保存?",
        "保存提示",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning"
        }
      )
        .then(() => {
          save();
        })
        .catch(() => {
          return;
        });
    } else {
      save();
    }
    //不确定有没有提交，先留着
    // const params: object = this.packagePageDataFn(false);
    // const res: any = await API.marketStrategySaveApi(params);
    // if (res.code == 200) {
    //   this.$message.success("保存成功！");
    //   this.$router.go(-1);
    // } else {
    //   this.$message.error(res.message);
    // }
  }
  //保存页面
  async summitPageFn() {
    (this.$refs["cform"] as any).validate(async (valid: boolean) => {
      let bresult: any = true;
      let cresult: any = true;

      this.$refs["bfrom"] &&
        (this.$refs["bfrom"] as any).validate((kvalid: boolean) => {
          bresult = kvalid;
        });
      this.$refs["afrom"] &&
        (this.$refs["afrom"] as any).validate((kvalid: boolean) => {
          cresult = kvalid;
        });

      if (valid && bresult && cresult) {
        const params: object = this.packagePageDataFn(true);
        console.log(params, "params");
        let res: any;
        if (!params) return;
        if (this.strategy_id) {
          res = await API.marketStrategyRecordEditSaveApi(params);
        } else {
          res = await API.marketStrategySaveApi(params);
        }
        if (res.code == 200) {
          this.$message.success("保存成功！");
          this.$router.go(-1);
        } else {
          this.$message.error(res.message);
        }
      } else {
        this.$message.warning("请仔细检查填写数据！");
        return false;
      }
    });
  }
  //组装relation给后台
  formatRelationFn(data: any, type: number) {
    let arr: any = [];
    if (type == 2) {
      arr = data.map((v: any) => {
        return { relation_id: v.id };
      });
    }
    return arr;
  }

  //设置页面初始化业务数据
  async setPageDataFn() {
    if (!this.strategy_id) return;
    const params: object = {
      id: this.strategy_id
    };
    const res: any = await API.marketStrategyDetailApi(params);
    if (res.code == 200) {
      this.oldStrategy = res.data;
      this.setInitDataFn(res.data);
    } else {
      this.$message.error(res.message);
    }
  }

  getSortPolicyByIdFn(type: number) {
    const result: any = {
      name: this.channelNameById[type],
      value: type,
      fixed: false
    };
    return result;
  }
  getChannelFn(obj: any) {
    console.log(obj, "--=======obj");
    const result: any = {
      title: this.channelNameById[obj.channel_type],
      id: obj.channel_type,
      list: []
    };
    if (obj.channel_type == 1) {
      result.list = [
        {
          content: obj.channel_detail ? obj.channel_detail.name : "",
          contentType: "",
          type: "",
          dialog: "msgDialog",
          originData: obj
        }
      ];
      this.channels[0].hasChecked = true;
    } else if (obj.channel_type == 2) {
      result.list = [
        {
          content: obj.content.text ? obj.content.text : "",
          contentType: "",
          type: "",
          dialog: "wechartMsgDialog",
          originData: obj
        }
      ];
      const index = this.channels.findIndex(item => {
        return item.dialog === "WeChatWorkMsgDialog";
      });
      this.channels[index].hasChecked = true;
      //这块接口需重新定义
    } else if (obj.channel_type == 3) {
      this.isused = obj.content.coupon_id;
      result.list = [
        {
          content: obj.content,
          dialog: "autoCouponDialog",
          originData: obj
        }
      ];
    }
    return result;
  }
  formatDetailRelationsFn(data: any = [], type: number) {
    const relation: any = [];
    if (type == 2) {
      data.forEach((v: any) => {
        relation.push({
          id: v.relation_id,
          spu_name: v.relation_detail.name,
          spu_source: {
            url: v.relation_detail.img_url
          }
        });
      });
    } else if (type == 3) {
      //暂时不做
      data.forEach((v: any) => {
        relation.push({
          id: v.relation_id,
          live_title: v.relation_detail.name
        });
      });
    } else if (type == 4) {
      data.forEach((v: any) => {
        relation.push({
          id: v.relation_id,
          template_name: v.relation_detail.name
        });
      });
    } else if (type == 5) {
      data.forEach((v: any) => {
        relation.push({
          id: v.relation_id,
          name: v.relation_detail.name
        });
      });
    }

    return relation;
  }
  setInitDataFn(data: any) {
    if (!data || Object.keys(data).length < 1) return;
    const form = this.strategyForm;
    form.strategyName = data.strategy_name;
    form.strategyItem = data.customer_group_id;
    // this.strategyItems.forEach((item: any) => {
    //   if (item.value == data.customer_group_id) {
    //     this.strategyItemInfo = item.info;
    //   }
    // });
    if (data.strategy_method == 1) {
      form.runMode = 1;
      this.isLongSrategy = true;
      const start: any = data.start_time
        ? new Date(data.start_time * 1000)
        : "";
      const end: any = data.start_time ? new Date(data.end_time * 1000) : "";
      form.long.runTimeRange = [start, end];
      form.condition = data.condition;
      if (form.condition != 1) {
        form.long.timeInterval = data.time_interval;
        form.long.timeType = data.time_type;
      }
      if (form.condition == 2) {
        form.long.relationType = data.relation_type;
        //访问页面；1=首页，2=商品，3=直播，4=活动页，5=活动 assignType videoType pageType activeType
        if (form.long.relationType == 2) {
          form.assignType = data.assign_type;
          this.selectedshopList = this.formatDetailRelationsFn(
            data.relations,
            data.relation_type
          );
          this.selectedShopNum = this.selectedshopList.length;
        } else if (form.long.relationType == 3) {
          form.videoType = data.assign_type;
          this.selectedVideoList = this.formatDetailRelationsFn(
            data.relations,
            data.relation_type
          );
          this.selectedVideoNum = this.selectedVideoList.length;
        } else if (form.long.relationType == 4) {
          form.pageType = data.assign_type;
          this.selectedPageList = this.formatDetailRelationsFn(
            data.relations,
            data.relation_type
          );
          this.selectedPageNum = this.selectedPageList.length;
        } else if (form.long.relationType == 5) {
          form.activityType = data.belong_type;
          form.activeType = data.assign_type;
          this.selectedActiveList = this.formatDetailRelationsFn(
            data.relations,
            data.relation_type
          );
          this.selectedActiveNum = this.selectedActiveList.length;
        }
      }
    } else {
      form.runMode = 2;
      if (data.strategy_method == 2) {
        form.onceRunMode = 1;
      } else if (data.strategy_method == 3) {
        form.onceRunMode = 2;
        form.short.planTime = data.start_time
          ? new Date(data.start_time * 1000)
          : "";
      }
    }
    form.phone = data.mobile;
    if (data.channels.length > 0) {
      try {
        this.selectedChannels = data.channels.map((item: any) =>
          this.getChannelFn(item)
        );
      } catch (error) {
        console.error(error, "===++++++");
      }
      console.log(this.selectedChannels, "=========");
      data.channels.forEach((v: any) => {
        this.sortPolicyList[String(v.order - 1)] = this.getSortPolicyByIdFn(
          v.channel_type
        );
      });
    }
    form.strategyPolicy = data.strategy_policy;
  }
  //组装参数方法type:false(保存)，true:(提交)
  packagePageDataFn(type = false) {
    const obj: any = this.strategyForm;
    const params: any = {
      strategy_name: obj.strategyName,
      customer_group_id: obj.strategyItem,
      strategy_type: 1,
      trigger_count: Number(this.strategyItemInfo.total)
    };
    if (this.strategy_id) {
      params.id = this.strategy_id;
    }
    //执行方式判断strategy_method1，2，3
    if (obj.runMode == 1) {
      params.strategy_method = 1;
      params.start_time = new Date(obj.long.runTimeRange[0]).getTime() / 1000;
      params.end_time = new Date(obj.long.runTimeRange[1]).getTime() / 1000;
      //触发条件;1=不限，2=访问，3=加购，4=下单，5=支付
      params.condition = obj.condition;
      if (obj.condition != 1) {
        params.time_interval = Number(obj.long.timeInterval);
        params.time_type = Number(obj.long.timeType);
      }
      if (obj.condition == 2) {
        params.relation_type = obj.long.relationType;
        //访问页面；1=首页，2=商品，3=直播，4=活动页，5=活动 assignType videoType pageType activeType
        if (params.relation_type == 2) {
          params.assign_type = obj.assignType;
          if (params.assign_type == 1) {
            params.relations = [];
          } else {
            params.relations = this.formatRelationFn(
              this.selectedshopList,
              params.assign_type
            );
            if (type && params.relations.length < 1) {
              this.$message.warning("请选择商品！");
              return;
            }
          }
        } else if (params.relation_type == 3) {
          params.assign_type = obj.videoType;
          if (params.assign_type == 1) {
            params.relations = [];
          } else {
            params.relations = this.formatRelationFn(
              this.selectedVideoList,
              params.assign_type
            );
            if (type && params.relations.length < 1) {
              this.$message.warning("请选择直播项！");
              return;
            }
          }
        } else if (params.relation_type == 4) {
          params.assign_type = obj.pageType;
          if (params.assign_type == 1) {
            params.relations = [];
          } else {
            params.relations = this.formatRelationFn(
              this.selectedPageList,
              params.assign_type
            );
            if (type && params.relations.length < 1) {
              this.$message.warning("请选择活动页！");
              return;
            }
          }
        } else if (params.relation_type == 5) {
          params.assign_type = obj.activeType;
          if (params.assign_type == 1) {
            params.relations = [];
          } else {
            params.belong_type = obj.activityType;
            params.relations = this.formatRelationFn(
              this.selectedActiveList,
              params.assign_type
            );
            if (type && params.relations.length < 1) {
              this.$message.warning("请选择活动！");
              return;
            }
          }
        }
      }
    } else {
      if (obj.onceRunMode == 1) {
        params.strategy_method = 2;
        params.start_time = 0;
        params.end_time = 0;
      } else {
        params.strategy_method = 3;
        params.start_time = new Date(obj.short.planTime).getTime() / 1000;
        params.end_time = new Date(obj.short.planTime).getTime() / 1000;
      }
    }
    //添加渠道
    const channels: Array<any> = [];
    const orderObj: object = {
      "1": -1,
      "2": -1,
      "3": -1,
      "4": -1
    };
    params.strategy_policy = obj.strategyPolicy;
    this.sortPolicyList.forEach((v: any, i: number) => {
      orderObj[String(v.value)] = i + 1;
    });
    if (this.selectedChannels.length < 1) {
      this.$message.warning("请添加渠道！");
      return;
    }
    this.selectedChannels.forEach((e: any) => {
      if (e.list.length > 0) {
        if (e.id == 1) {
          const obj: any = e.list[0].originData;
          channels.push({
            channel_type: e.id,
            junior_id: obj.junior_id ? obj.junior_id : obj.id,
            order: orderObj[e.id],
            content: {}
          });
        } else if (e.id == 2) {
          if (e.list.length === 2) {
            const obj: any = e.list[1].originData;
            channels.push({
              channel_type: e.id,
              junior_id: 0,
              order: orderObj[e.id],
              content:
                {
                  type: obj.data.type - 1,
                  video: {},
                  text: e.list[0].originData.name || "",
                  link: obj.data.type - 1 == 2 ? obj.data : {},
                  picture: obj.data.type - 1 == 1 ? obj.data : {},
                  applet: obj.data.type - 1 == 3 ? obj.data : {}
                } || {}
            });
          } else {
            channels.push({
              channel_type: e.list[0].originData.channel_type,
              junior_id: 0,
              order: e.list[0].originData.order,
              content: e.list[0].originData.content
            });
          }
        } else if (e.id == 3) {
          channels.push({
            channel_type: e.id,
            junior_id: 0,
            order: orderObj[e.id],
            content: e.list[0].content
          });
        }
      }
    });
    params.channels = channels;
    params.mobile = obj.phone;
    if (this.strategy_id) {
      return Object.assign(this.oldStrategy, params);
    }
    return params;
  }
  setStrategyItemInfoFn(value: number | string = 0) {
    if (!(value || this.strategyForm.strategyItem)) return;
    const id = value || this.strategyForm.strategyItem;
    this.strategyItems.forEach(item => {
      if (item.value == id) {
        this.strategyItemInfo.total = item.num;
      }
    });
  }
  /**
   * @description:删除卡片，并释放已选类型的状态，更新触达规则
   * @param {index} 当前卡片在selectedChannels的下表
   * @return {type}
   */
  channelsDelThisFn(index: number) {
    const id: number = this.selectedChannels[index].id;
    const ind: number = this.getItemIndexFn(this.channels, id);
    if (ind > -1) {
      this.channels[ind].hasChecked = false;
    }
    const i: number = this.getItemIndexFn(this.sortPolicyList, id);
    if (i > -1) {
      this.sortPolicyList.splice(i, 1);
    }
    this.chooseInfo = {};
    this.selectedChannels.splice(index, 1);
    this.msgToDialog = null;
    this.isLongTermExecution = false;
    this.isShortTermExecution = false;
    this.showTrigger = true;
  }
  /**
   * @description: 卡片中多行中会出现删除按钮，这是该按钮的逻辑
   * @param {fInd} 当前this.selectedChannels的item的下标
   * @param {ind} 当前this.selectedChannels中list的item的下标
   * @return {type}
   */
  channelDelThisFn(fInd: number, ind: number) {
    this.selectedChannels[fInd].list.splice(ind, 1);
    if (this.selectedChannels[fInd].list.length < 1) {
      this.channelsDelThisFn(fInd);
    }
  }
  channelReplaceThisFn(fInd: number, ind: number, content: string) {
    //this.channels
    //selectedChannels
    this.msgName = content;
    const item: any = this.selectedChannels[fInd].list[ind];
    const dialog: string = item.dialog;
    this[dialog] = true;
    this.msgToDialog = item;
    this.channelReplaceIndex = ind;
    this.channelReplaceFIndex = fInd;
  }
  channelReplaceThisFun(item: any) {
    // const item: any = this.selectedChannels[2];
    // const dialog: string = item.dialog;
    console.log(item);
    this.WeChatWorkMsgDialog = true;
    this.msgToDialog = item;
    this.channelReplaceFIndex = 2;
  }
  couponReplaceThisFun(item: any) {
    this.chooseInfo = item["list"][0]["content"];
    this.autoCouponDialog = true;
    this.channelReplaceFIndex = 3;
  }
  channelChangeFn(item: any) {
    const dialog: string = item.dialog;
    if (dialog) {
      this[dialog] = true;
    }
  }
  conditionChangeFn(value: number) {
    this.strategyForm.long.timeInterval = "";
    switch (value) {
      case 2:
      case 5:
        this.strategyForm.long.timeType = 1;
        break;
      case 3:
        this.strategyForm.long.timeType = 3;
        break;
      case 4:
        this.strategyForm.long.timeType = 2;
        break;
      default:
        this.strategyForm.long.timeType = 1;
        break;
    }
  }
  delSelectShopsFn(result: any) {
    this.selectedshopList = result;
    this.selectedShopNum = result.length;
  }
  delSelectVideosFn(result: any) {
    this.selectedVideoList = result;
    this.selectedVideoNum = result.length;
  }
  delSelectPagesFn(result: any) {
    this.selectedPageList = result;
    this.selectedPageNum = result.length;
  }
  delSelectActiveFn(result: any) {
    this.selectedActiveList = result;
    this.selectedActiveNum = result.length;
  }

  openShopsFn() {
    this.shopListDialog = true;
  }
  openVideosFn() {
    this.videoListDialog = true;
  }
  openActivePagesFn() {
    this.activePageDialog = true;
  }
  openActiveFn() {
    this.activeDialog = true;
  }
  getSelectVideosFn(rows: any) {
    this.videoListDialog = !this.videoListDialog;
    if (rows) {
      this.selectedVideoList = rows;
      this.selectedVideoNum = this.selectedVideoList.length;
    }
  }
  getSelectActiveFn(rows: any) {
    this.activeDialog = !this.activeDialog;
    if (rows) {
      this.selectedActiveList = rows;
      this.selectedActiveNum = this.selectedActiveList.length;
    }
  }
  ArrHasItemFn(arr: any, item: any) {
    const has: number = arr.findIndex((v: any) => {
      return v.id == item.id;
    });
    if (has > -1) {
      return true;
    }
    return false;
  }
  getSelectShopsFn(rows: any) {
    this.shopListDialog = !this.shopListDialog;
    if (rows) {
      this.selectedshopList = rows;
      this.selectedShopNum = this.selectedshopList.length;
    }
  }
  getSelectPagesFn(rows: any) {
    this.activePageDialog = !this.activePageDialog;
    if (rows) {
      this.selectedPageList = rows;
      this.selectedPageNum = this.selectedPageList.length;
    }
  }
  getSelectMsgFn(row: any) {
    this.msgDialog = !this.msgDialog;
    const id = 1;
    if (row) {
      const item: any = {
        title: this.channelNameById[String(id)],
        id: id,
        list: [
          {
            content: row.template_name,
            contentType: "",
            type: "",
            dialog: "msgDialog",
            originData: row
          }
        ]
      };

      const ind: number = this.selectedChannels.findIndex((item: any) => {
        return item.id == 1;
      });

      if (ind > -1) {
        this.selectedChannels[ind].list = item.list;
      } else {
        this.selectedChannels.push(item);
        this.channels[0].hasChecked = true;
        this.sortPolicyList.push({
          name: this.channelNameById[String(id)],
          value: id,
          fixed: false
        });
      }
    }
  }
  // 企微群发弹窗
  getSelectWeChatWorkMsgFn(row: Array<any>) {
    this.WeChatWorkMsgDialog = !this.WeChatWorkMsgDialog;
    const id = 2;

    if (row) {
      this.isEnterpriseMicro = true;
      this.channelReplaceIndex = -1;
      this.channelReplaceFIndex = -1;
      //type 1：文本，2：图片，3：链接，4：小程序卡片
      if (row[1].data) {
        switch (row[1].data.type) {
          case 2:
            row[1].typeName = "图片";
            break;
          case 3:
            row[1].typeName = "链接";
            break;
          case 4:
            row[1].typeName = "小程序";
            break;
        }
      }
      console.log(row);
      const item: any = {
        title: this.channelNameById[String(id)],
        id: id,
        list: row.map((item: any) => {
          return {
            // content: row[1].data && row[1].data !== {} ? otherName : "",
            // contentText: textName,
            contentType: item.typeName,
            type: item.type,
            dialog: "wechartMsgDialog",
            originData: item
          };
        })
      };
      const ind: number = this.selectedChannels.findIndex((item: any) => {
        return item.id == id;
      });
      if (this.strategy_id && ind > -1) {
        this.selectedChannels[ind].list = item.list;
      } else {
        if (ind > -1) {
          this.selectedChannels[ind].list = item.list;
        } else {
          this.selectedChannels.push(item);
          this.channels[1].hasChecked = true;
          this.sortPolicyList.push({
            name: this.channelNameById[String(id)],
            value: id,
            fixed: false
          });
        }
      }
      this.strategyForm.runMode = 2;
      this.msgToDialog = item;
    }
  }
  // 自动发券弹窗
  getSelectAutoCouponFn(row: any) {
    this.autoCouponDialog = !this.autoCouponDialog;
    const id = 3;
    if (row) {
      this.isShortTermExecution = true;
      const item: any = {
        title: this.channelNameById[String(id)],
        id,
        list: [
          {
            content: {
              query_coupon: row.queryCoupon ? 1 : 2,
              customer_channel: Number(row.radio),
              customer_channel_value:
                row.radio == "2"
                  ? row.sceneList
                  : row.radio == "3"
                  ? row.channelList
                  : {},
              coupon_id: row.coupon,
              coupon_list: row.couponList
            },
            dialog: "autoCouponDialog",
            originData: row
          }
        ]
      };
      this.showTrigger = Number(row.radio) != 3;
      const ind: number = this.selectedChannels.findIndex((item: any) => {
        return item.id == id;
      });
      if (this.strategy_id && ind > -1) {
        this.selectedChannels[ind].list = item.list;
      } else {
        if (ind > -1) {
          this.selectedChannels[ind].list = item.list;
        } else {
          this.selectedChannels.push(item);
          this.channels[1].hasChecked = true;
          this.sortPolicyList.push({
            name: this.channelNameById[String(id)],
            value: id,
            fixed: false
          });
        }
      }
      this.strategyForm.runMode = 1;
      this.strategyForm.condition = 1;
      this.strategyForm.long.timeInterval = "";
    }
  }
  replaceWeChatWorkMsgFn(row: any) {
    this.wechartMsgDialog = !this.wechartMsgDialog;
    if (row) {
      //type 1：文本，2：图片，3：链接，4：小程序卡片
      const item: any = {
        content: row.name,
        contentType: row.typeName,
        type: row.type,
        dialog: "wechartMsgDialog",
        originData: row
      };
      this.selectedChannels[this.channelReplaceFIndex].list[
        this.channelReplaceIndex
      ] = item;
    }
  }
  //生命周期函数
  created() {
    const id: string | (string | null)[] = this.$route.query.id;
    const status: string | (string | null)[] = this.$route.query.status;

    if (id) {
      this.strategy_id = Number(id);
      this.strategy_status = Number(status);
      this.$route.query.way.indexOf("企微群发") != -1
        ? (this.isEnterpriseMicro = true) && (this.isLongTermExecution = true)
        : (this.isEnterpriseMicro = false) &&
          (this.isLongTermExecution = false);
      this.$route.query.way.indexOf("自动发券") != -1
        ? (this.isShortTermExecution = true)
        : (this.isShortTermExecution = false);
    }
    this.initPageOptionFn();
    let localData: any = window.localStorage.getItem("STRATAGE_LOCAL_DATA");
    if (localData && id === undefined) {
      localData = JSON.parse(localData);
      const { selectedChannels, sortPolicyList, strategyForm } = localData;
      this.selectedChannels = selectedChannels;
      this.sortPolicyList = sortPolicyList;
      this.strategyForm = strategyForm;
      this.selectedChannels[0].id === 2
        ? (this.isEnterpriseMicro = true) && (this.isLongTermExecution = true)
        : (this.isEnterpriseMicro = false) &&
          (this.isLongTermExecution = false);
      this.selectedChannels[0].id === 2
        ? (this.channels[1].hasChecked = true)
        : (this.channels[1].hasChecked = false);
      this.selectedChannels[0].id === 3
        ? (this.isShortTermExecution = true)
        : (this.isShortTermExecution = false);
    }
    // this.getSceneList();
  }
}
</script>
<style lang="scss">
.setStratage {
  .el-dialog {
    max-height: inherit !important;
  }
}

.popover-Tip {
  height: 100px;
  min-width: 300px;
  max-width: 400px;
  overflow: auto;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}
</style>
<style lang="scss" scoped>
// 公共颜色
@media screen and (max-width: 1440px) {
  .setStratage {
    .container_top_20 {
      margin-top: 20px;
    }
  }
}
.setStratage {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  .marginT20 {
    margin-top: 20px;
  }
  .channel-container {
    .other-card {
      width: 500px;
      position: relative;
      overflow: inherit;
      .del-btn {
        position: absolute;
        right: -10px;
        top: -10px;
      }
    }
    margin: 20px 0;
    .title {
      font-size: 16px;
      font-weight: bold;
    }
    .del-btn {
      font-size: 16px;
      position: absolute;
      right: -10px;
      top: -10px;
      background: #999999a3;
      color: #fff;
      border-radius: 50%;
      padding: 5px;
      cursor: pointer;
    }
    .del-btn:hover {
      background: #9999998a;
    }
    .box-body {
      line-height: 36px;
      // height: 36px;
      .cloumn-name {
        height: 36px;
        width: 290px;
        min-height: 2px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
      .cloumn-type {
        height: 36px;
        width: 50px;
        min-height: 2px;
        text-align: center;
        color: #999999f7;
      }
      .content-type {
        position: absolute;
        top: 100px;
        right: 50px;
      }
      .cloumn-btn {
        width: 110px;
        text-align: right;
      }
    }
  }
  .splitLine {
    height: 0;
    border-top: 1px solid #ebeef5;
    margin-bottom: 20px;
  }
  .marginT0 {
    margin-top: 0;
  }
  .margin-b-0 {
    margin-bottom: 0 !important;
  }
  .marginB10 {
    margin-bottom: 10px;
  }
  .flip-list-move {
    transition: transform 0.5s;
  }

  .no-move {
    transition: transform 0s;
  }

  .ghost {
    opacity: 0.5;
    background: #c8ebfb;
  }

  .list-group {
    min-height: 20px;
  }

  .list-group-item {
    cursor: move;
  }

  .list-group-item i {
    cursor: pointer;
  }
  .drag-container {
    width: 350px;
    margin-left: 20px;
    .drag-header {
      background: #ebebf1;
      font-size: 14px;
      color: #303133;
      height: 36px;
      line-height: 36px;
      padding-left: 10px;
      font-weight: bold;
      padding-right: 5px;
      div {
        padding: 0 5px;
      }
    }
    .drag-body {
      ul {
        padding: 0;
        margin: 0;
        list-style: none;
      }
      .cell {
        padding: 10px 15px;
        border-bottom: 1px solid #dcdfe6;
        height: 35px;
        .drag_icon {
          height: 14px;
          background: url("https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/customer/drag_icon.png") no-repeat
            right center;
          background-size: 14px;
          margin-right: 4px;
        }
      }
      .cell.paddingL20 {
        padding-left: 25px;
      }
    }
  }
  .lineHeight40 {
    line-height: 40px;
  }
  .paddingT0 {
    padding-top: 0;
  }
  .padding-l-r-20 {
    text-align: right;
    vertical-align: middle;
    font-size: 14px;
    color: #606266;
    line-height: 40px;
    -webkit-box-sizing: border-box;
    box-sizing: border-box;
    width: 106px;
    margin-right: 20px;
    display: inline-block;
  }
  .margin-t-b-20 {
    margin: 20px 0;
  }
  .margin-B15 {
    margin-bottom: 15px;
  }
  .shop-box-container {
    padding-left: 50px;
    width: 700px;
  }
  .inlineBlock {
    display: inline-block;
  }
  // .p-l-10 {
  //   padding-left: 10px;
  // }
  .paddingL20 {
    padding-left: 20px;
  }
  .marginL50 {
    margin-left: 50px;
  }
  .paddingR10 {
    padding-right: 10px;
  }
  // .m-l-10 {
  //   margin-left: 10px;
  // }
  .runMode.el-radio-group {
    display: block;
    font-size: inherit;
    .block {
      margin: 20px 0;
      margin-top: 12px;
    }
  }
  .conditionInfo {
    margin-left: 75px;
    .timeType .el-select.el-select--small {
      width: 80px !important;
    }
  }
  .headline {
    font-size: 16px;
    font-weight: 700;
    margin-bottom: 30px;
    .mark {
      font-size: 14px;
      font-weight: normal;
      color: #999;
    }
  }
  .spanText {
    display: inline-block;
    padding-left: 20px;
    line-height: 16px;
    .main {
      color: #999;
    }
  }
}
.tip {
  color: #999;
  font-size: 12px;
  margin-top: 5px;
}
.enterprisemicro-button {
  float: right;
  margin-right: 30px;
  margin-top: 3px;
}
.cloumn-types {
  height: 36px;
  min-height: 2px;
  text-align: center;
  color: #999999f7;
}
.home-page {
  background-color: $color-white;
  border-radius: 10px;
  padding: 30px;
}

/deep/.cell {
  .user-LableButtonList {
    display: flex;
    .user-groupbutton {
      font-size: 12px;
      width: 52pt;
      height: 24pt;
      line-height: 24pt;
      background: #fef8ea;
      color: #f6ab3c;
      border-radius: 3pt;
      box-sizing: border-box;
      border: 1px solid #fce5c4;
      overflow: hidden;
      cursor: pointer;
      margin-right: 5px;
      text-align: center;
    }
  }
  .user-LableButtonLists {
    display: flex;
  }

  .user-LableButtonList {
    display: flex;
    margin-left: 5px;
    .user-lablebutton {
      font-size: 12px;
      height: 24pt;
      line-height: 24pt;
      background: #f0f3fe;
      color: #6f3112;
      border-radius: 3pt;
      box-sizing: border-box;
      border: 1px solid #d4dcfc;
      overflow: hidden;
      cursor: pointer;
      margin-right: 5px;
      text-align: center;
      width: 99px;
      display: inline-block;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      vertical-align: middle;
    }
  }
}
</style>
