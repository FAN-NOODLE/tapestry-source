<template>
  <div class="welcome">
    <layout-form>
      <template slot="content">
        <p class="headline ">基础信息</p>
        <div class="item-container">
          <p class="line">
            <span class="line-key">策略名称</span>
            <span class="line-value">{{ strategy.strategyName }}</span>
          </p>
          <p class="line">
            <span class="line-key">客户分组</span>
            <span class="line-value">
              <span class="like-input">{{ strategy.strategyItem }}</span>
              <span class="like-input-tips">
                <span>
                  预计触达人数 ：{{ strategy.strategyItemInfo.total }}</span
                >
                <!-- <span class="main">
                  <span>
                    （ {{ strategy.strategyItemInfo.phone }}有手机号；{{
                      strategy.strategyItemInfo.wechat
                    }}是企微好友；{{
                      strategy.strategyItemInfo.visitor
                    }}是商城访客；{{
                      strategy.strategyItemInfo.Accounts
                    }}是公众号粉丝 ）</span
                  ></span
                > -->
              </span></span
            >
          </p>
          <p class="line flex-row">
            <span class="line-key">触达渠道</span>
            <span class="line-object">
              <div
                class="channel-container"
                v-for="(item, index) in selectedChannels"
                :key="index"
              >
                <el-card v-if="item.id == 3" class="coupon-card">
                  <span class="title">{{ item.title }}</span>
                  <div class="box-body" style="max-height: 300px;">
                    <div>
                      <el-table
                        :data="item.list[0].content.coupon_list"
                        max-height="310px"
                        ref="table"
                      >
                        <el-table-column
                          prop="coupon_stock_id"
                          label="优惠券批次号"
                          width="200"
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
                              v-show="
                                item.list[0].content.customer_channel == 2
                              "
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
                                v-show="
                                  sceneChannel && sceneChannel.length <= 2
                                "
                              >
                                <div
                                  v-for="(itemTags, i) in sceneChannel"
                                  :key="i"
                                >
                                  <el-tooltip
                                    effect="light"
                                    :content="itemTags"
                                    placement="top-start"
                                  >
                                    <div class="user-lablebutton">
                                      <span>{{ itemTags }}</span>
                                    </div>
                                  </el-tooltip>
                                </div>
                              </div>
                            </div>
                          </template>
                        </el-table-column>
                        <el-table-column
                          v-if="
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
                              v-show="
                                item.list[0].content.customer_channel == 3
                              "
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
                                <div>
                                  <el-tooltip
                                    effect="light"
                                    :content="channelNames[1]"
                                    placement="top-start"
                                  >
                                    <div class="user-lablebutton">
                                      <span>{{ channelNames[1] }}</span>
                                    </div>
                                  </el-tooltip>
                                </div>
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
                                v-show="
                                  channelNames && channelNames.length <= 2
                                "
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
                          width="150"
                          :show-overflow-tooltip="true"
                        >
                          <template>
                            {{
                              item.list[0].content.query_coupon == 1
                                ? "是"
                                : "否"
                            }}
                          </template>
                        </el-table-column>
                        <el-table-column
                          label="券状态"
                          width="150"
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
                            {{
                              row.start_time | formatDate("yyyy-mm-dd hh:mm:ss")
                            }}
                            -
                            {{
                              row.end_time | formatDate("yyyy-mm-dd hh:mm:ss")
                            }}
                          </template>
                        </el-table-column>
                      </el-table>
                    </div>
                  </div>
                </el-card>
                <el-card v-else class="box-card m-t-20 other-card">
                  <span class="title">{{ item.title }}</span>
                  <div class="box-body">
                    <div
                      class="clearfix"
                      v-for="(ele, ind) in item.list"
                      :key="ind"
                    >
                      <span class="pull-left cloumn-content">
                        <el-popover
                          placement="top-start"
                          width="290"
                          trigger="hover"
                          :disabled="ele.content.length <= 20"
                        >
                          <div>
                            {{ ele ? ele.content : "" }}
                          </div>
                          <span
                            slot="reference"
                            v-if="ele.content.length <= 20"
                            class="el-popover-ref"
                            >{{ ele.content }}</span
                          >
                          <span slot="reference" v-else>{{
                            ele.content.substr(0, 20) + "..."
                          }}</span>
                        </el-popover>
                        <span
                          class="pull-left cloumn-content"
                          v-if="item.title == '企微群发'"
                        >
                          {{
                            ele.originData.content
                              ? ele.originData.content.applet.app_name ||
                                ele.originData.content.link.title ||
                                ele.originData.content.picture.name
                              : ""
                          }}
                        </span>
                      </span>
                      <span
                        class="pull-left cloumn-type clearfix"
                        v-if="ele.content !== ''"
                      >
                        {{ ele ? ele.contentType : "" }}
                      </span>

                      <span
                        class="pull-left cloumn-type clearfix"
                        v-if="item.title == '企微群发'"
                      >
                        {{ ele ? ele.type : "" }}
                      </span>
                    </div>
                  </div>
                </el-card>
              </div>
            </span>
          </p>
          <div class="line flex-row">
            <span class="line-key">执行方式</span>
            <div class="line-object">
              <!-- 单次执行-定时 -->
              <span class="container" v-if="strategy.strategy_method == 3">
                <span class="model">单次执行 </span>
                <div>
                  <span class="value like-input"
                    >定时执行： {{ strategy.short.planTime }}
                  </span>
                </div>
              </span>
              <!-- 立即执行 -->
              <span class="container" v-if="strategy.strategy_method == 2">
                <div>
                  <span class="model">单次执行 </span>
                </div>
                <div>
                  <span class="key like-input"> 立即执行</span>
                </div>
              </span>
              <!-- 长期执行   v-if="strategy.strategy_method == 1"-->
              <span class="block" v-if="strategy.strategy_method == 1">
                <span class="row">
                  长期执行
                </span>
                <div class="row mode-content">
                  <div class="row">
                    <span class="key">计划执行时间</span>
                    <span class="value">{{
                      strategy.long.runTimeRange.join("  -  ")
                    }}</span>
                  </div>
                  <div class="row">
                    <span class="key">触发条件</span>
                    <span class="value" v-if="strategy.condition == 1"
                      >不限</span
                    >
                    <!-- testData -->
                    <span
                      class="value"
                      v-if="
                        strategy.condition == 2 &&
                          selectedChannels[0] &&
                          selectedChannels[0].id != 3
                      "
                      >{{
                        `访问后 ${strategy.long.timeInterval} ${
                          timeTypesObj[String(strategy.long.timeType)]
                        } 触达客户`
                      }}</span
                    >
                    <span
                      class="value"
                      v-if="
                        strategy.condition == 3 &&
                          selectedChannels[0] &&
                          selectedChannels[0].id != 3
                      "
                      >{{
                        `加购后 ${strategy.long.timeInterval} ${
                          timeTypesObj[String(strategy.long.timeType)]
                        } 未成交触达客户`
                      }}</span
                    >
                    <span
                      class="value"
                      v-if="
                        strategy.condition == 4 &&
                          selectedChannels[0] &&
                          selectedChannels[0].id != 3
                      "
                      >{{
                        `下单后 ${strategy.long.timeInterval} ${
                          timeTypesObj[String(strategy.long.timeType)]
                        } 未成交触达客户`
                      }}</span
                    >
                    <span
                      class="value"
                      v-if="
                        strategy.condition == 5 &&
                          selectedChannels[0] &&
                          selectedChannels[0].id != 3
                      "
                      >{{
                        `支付后 ${strategy.long.timeInterval} ${
                          timeTypesObj[String(strategy.long.timeType)]
                        } 内触达客户`
                      }}</span
                    >
                  </div>
                  <div class="row" v-if="strategy.condition == 2">
                    <span class="key">访问页面</span>
                    <span class="value" v-if="strategy.long.relationType == 1">
                      <span class="value">{{
                        assignTypes[strategy.long.relationType][0]
                      }}</span>
                    </span>
                    <span class="value" v-if="strategy.long.relationType == 2">
                      <span class="value">{{
                        assignTypes[strategy.long.relationType][
                          strategy.assignType - 1
                        ]
                      }}</span>

                      <span class="value" v-if="strategy.assignType == 2">
                        {{ `已选 ${selectedshopList.length} 个商品 ` }}
                        <el-button
                          type="text"
                          v-if="selectedshopList.length > 0"
                          @click="isOpen = !isOpen"
                          >{{ isOpen ? " 收起" : " 展开" }}
                          <i
                            aria-hidden="true"
                            class="el-icon--right"
                            :class="
                              isOpen ? 'el-icon-arrow-up' : 'el-icon-arrow-down'
                            "
                          ></i>
                        </el-button>
                      </span>
                    </span>
                    <span class="value" v-if="strategy.long.relationType == 3">
                      <span class="value">{{
                        assignTypes[strategy.long.relationType][
                          strategy.videoType - 1
                        ]
                      }}</span>
                      <span class="value" v-if="strategy.videoType == 2">
                        {{ `已选 ${selectedVideoList.length} 个直播` }}
                        <el-button
                          type="text"
                          v-if="selectedVideoList.length > 0"
                          @click="isOpen = !isOpen"
                          >{{ isOpen ? " 收起" : " 展开" }}
                          <i
                            aria-hidden="true"
                            class="el-icon--right"
                            :class="
                              isOpen ? 'el-icon-arrow-up' : 'el-icon-arrow-down'
                            "
                          ></i>
                        </el-button>
                      </span>
                    </span>
                    <span class="value" v-if="strategy.long.relationType == 4">
                      <span class="value">{{
                        assignTypes[strategy.long.relationType][
                          strategy.pageType - 1
                        ]
                      }}</span>
                      <span class="value" v-if="strategy.pageType == 2">
                        {{ `已选 ${selectedPageList.length} 个活动页` }}
                        <el-button
                          type="text"
                          v-if="selectedPageList.length > 0"
                          @click="isOpen = !isOpen"
                          >{{ isOpen ? " 收起" : " 展开" }}
                          <i
                            aria-hidden="true"
                            class="el-icon--right"
                            :class="
                              isOpen ? 'el-icon-arrow-up' : 'el-icon-arrow-down'
                            "
                          ></i>
                        </el-button>
                      </span>
                    </span>
                  </div>
                  <!-- 商品列表 -->
                  <div
                    class="row list"
                    v-if="strategy.assignType == 2 && isOpen"
                  >
                    <selected-shops
                      :list="selectedshopList"
                      :operationBtn="false"
                    ></selected-shops>
                  </div>
                  <!-- 直播列表 -->
                  <div
                    class="row list"
                    v-if="strategy.videoType == 2 && isOpen"
                  >
                    <selected-videos
                      :list="selectedVideoList"
                      :operationBtn="false"
                    ></selected-videos>
                  </div>
                  <!-- 活动列表 -->
                  <div class="row list" v-if="strategy.pageType == 2 && isOpen">
                    <selected-active-page
                      :list="selectedPageList"
                      :operationBtn="false"
                    ></selected-active-page>
                  </div>
                </div>
              </span>
            </div>
          </div>
        </div>
        <div class="item-container">
          <p class="line flex-row">
            <span class="line-key">触达规则</span>
            <span class="line-value" v-if="strategy.strategyPolicy == 2"
              >通知所有渠道</span
            >
            <span class="line-object" v-else></span>
          </p>
        </div>
        <!-- <div class="splitLine"></div> -->
        <!-- <p class="headline">
          异常提醒
          <span class="mark">短信条数不足100条时通过短信提醒</span>
        </p> -->
      </template>
      <template slot="footer">
        <el-button @click="goBack()" plain>返回</el-button>
        <el-button
          type="primary"
          plain
          v-if="status == 1 || status == 3"
          @click="itemStartFn()"
          >开始</el-button
        >
        <el-button
          v-if="
            status == 2 &&
              strategy.strategy_method != 2 &&
              strategy.strategy_method != 3
          "
          @click="itemPauseFn()"
          type="primary"
          plain
          >暂停</el-button
        >
        <el-button
          v-if="status == 1 || status == 3"
          @click="editItem()"
          type="primary"
          plain
          >编辑</el-button
        >
        <el-button v-if="status == 1" @click="delItem()" type="danger" plain
          >删除</el-button
        >
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Watch } from "vue-property-decorator";
import selectedShops from "../../components/selected-shops.vue";
import selectedVideos from "../../components/selected-videos.vue";
import selectedActivePage from "../../components/selected-active-page.vue";
import API from "@/api/customer";
import { dateFormat } from "../../util/index";
@Component({
  components: {
    selectedShops,
    selectedVideos,
    selectedActivePage
  }
})
export default class extends Vue {
  private status = 0;
  private itemId = 0;
  private isOpen = true;
  private strategy: any = {
    strategyName: "",
    strategyItem: "",
    strategyItemInfo: {
      total: 0,
      phone: 0,
      wechat: 0,
      visitor: 0,
      Accounts: 0
    },

    strategy_method: 0,
    long: {
      runTimeRange: [],
      timeInterval: "",
      relationType: "",
      activityType: 0,
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
    strategyPolicy: 2,
    phone: "未填写"
  };
  private selectedVideoList: Array<any> = [];
  private selectedshopList: Array<any> = [];
  private selectedPageList: Array<any> = [];
  private selectedActiveList: Array<any> = [];
  private selectedChannels: Array<any> = [];
  private selectedChannel = { title: "" };
  private info: any = {
    content: {
      message: "", // 文本
      text: "", // 图片/链接/小程序名称
      fileData: null, // 原始数据
      fileType: 0 //文本：1/ 图片：2/链接：3/小程序：4
    }
  };
  private channelNameById = {
    "1": "短信",
    "2": "企微群发",
    "3": "自动发券",
    "4": "公众号"
  };
  private timeTypesObj = {
    "1": "秒",
    "2": "分钟",
    "3": "小时",
    "4": "天"
  };
  private assignTypes = {
    "1": ["首页"],
    "2": ["全部商品", "指定商品"],
    "3": ["全部直播", "指定直播"],
    "4": ["全部活动页", "指定活动页"],
    "5": ["全部活动", "指定活动"]
  };
  private sortPolicyList: Array<any> = [];
  private sortPolicyLevels = ["1", "2", "3", "4"];
  private showChannel = true;
  private channelName = "";
  private channelNames: Array<any> = [];
  private sceneChannel = [];
  private sceneList = [];
  @Watch("selectedChannels", { deep: true, immediate: true })
  channelRefCustomGroup(newVal: Array<string>) {
    if (
      newVal[0] &&
      newVal[0]["id"] == 3 &&
      newVal[0]["list"][0]["content"]["customer_channel"] == 3
    ) {
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
                      this.channelNames.push(g["name"]);
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
      // API.sceneListApi({}).then((res: any) => {
      //   if (res.code == 200) {
      //     this.sceneList = res.data;
      //     this.showChannel = false;
      //     this.sceneChannel = [];
      //     const scene = Object.values(
      //       newVal[0]["list"][0]["content"]["customer_channel_value"]
      //     )[0];
      //     console.log(scene, "scene");
      //     this.sceneList.forEach(v => {
      //       if ((scene as any).indexOf(v["val"]) != -1) {
      //         this.sceneChannel.push(v["name"]);
      //       }
      //     });
      //     this.showChannel = true;
      //   } else {
      //     this.$message.error(res.message);
      //   }
      // });
    }
  }
  itemStartFn() {
    this.$confirm("是否需要开始该营销策略", "开始确认", {
      confirmButtonText: "确认开始",
      confirmButtonClass: "button-primary",
      cancelButtonText: "我再想想",
      type: "info"
    })
      .then(async () => {
        const params: object = {
          id: Number(this.itemId),
          strategy_status: 2
        };
        const res: any = await API.marketingUpdateStatusApi(params);
        if (res.code === 200) {
          this.$message({
            message: "操作成功",
            type: "success",
            duration: 2000
          });
          this.setPageDataFn();
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消操作"
        });
      });
  }
  itemPauseFn() {
    this.$confirm("是否需要暂停该营销策略", "暂停确认", {
      confirmButtonText: "确认暂停",
      confirmButtonClass: "button-primary",
      cancelButtonText: "我再想想",
      type: "info"
    })
      .then(async () => {
        const params: object = {
          id: this.itemId,
          strategy_status: 3
        };
        const res: any = await API.marketingUpdateStatusApi(params);
        if (res.code === 200) {
          this.$message({
            message: "操作成功",
            type: "success",
            duration: 2000
          });
          this.setPageDataFn();
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消操作"
        });
      });
  }
  delItem() {
    this.$confirm("是否需要删除该营销策略", "删除确认", {
      confirmButtonText: "确认删除",
      confirmButtonClass: "button-danger",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(async () => {
        const params: object = {
          id: this.itemId
        };
        const res: any = await API.marketStrategyRecordDelApi(params);
        if (res.code === 200) {
          this.$message({
            message: "删除成功",
            type: "success",
            duration: 2000
          });
          this.$router.push({
            name: "myapp_customer_auto-market_market-strategy_market-list"
          });
        } else {
          this.$message.error(res.message);
        }
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  editItem() {
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_edit-stratage",
      query: {
        id: String(this.itemId),
        status: String(this.status),
        way: this.$route.query.way,
        page: this.$route.query.page,
        size: this.$route.query.size,
        param: this.$route.query.param
      }
    });
  }
  goBack() {
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_market-list",
      params: JSON.parse(JSON.stringify(this.$route.query))
    });
  }
  getChannelFn(obj: any) {
    const result: any = {
      title: this.channelNameById[obj.channel_type],
      id: obj.channel_type,
      list: []
    };
    if (obj.channel_type == 1) {
      result.list = [
        {
          content: obj.channel_detail ? obj.channel_detail.name : "",
          contentType: "文本",
          type: "",
          dialog: "msgDialog",
          originData: obj
        }
      ];
    } else if (obj.channel_type == 2) {
      let typeName = "";
      switch (obj.content.type) {
        case 1:
          typeName = "图片";
          break;
        case 2:
          typeName = "链接";
          break;
        case 3:
          typeName = "小程序";
      }
      result.list = [
        {
          content: obj.content.text ? obj.content.text : "",
          contentType: "文本",
          type: typeName ? typeName : "",
          dialog: "WeChatWorkMsgDialog",
          originData: obj
        }
      ];
      //这块接口需重新定义
    } else if (obj.channel_type == 3) {
      result.list = [
        {
          content: obj.content,
          dialog: "autoCouponDialog",
          originData: obj,
          contentType: "自动发券",
          type: ""
        }
      ];
    }
    return result;
  }
  getSortPolicyByIdFn(type: number) {
    const result: any = {
      name: this.channelNameById[type],
      value: type,
      fixed: false
    };
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
    const form = this.strategy;
    form.strategyName = data.strategy_name;
    this.status = data.strategy_status;
    form.strategyItem = data.customer_group_name;
    this.setStrategyItemInfoFn(data.customer_group_id);
    if (data.strategy_method == 1) {
      form.strategy_method = 1;
      const start: any = data.start_time
        ? dateFormat(data.start_time, "yyyy-MM-dd hh:mm")
        : "";
      const end: any = data.start_time
        ? dateFormat(data.end_time, "yyyy-MM-dd hh:mm")
        : "";
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
        } else if (form.long.relationType == 3) {
          form.videoType = data.assign_type;
          this.selectedVideoList = this.formatDetailRelationsFn(
            data.relations,
            data.relation_type
          );
        } else if (form.long.relationType == 4) {
          form.pageType = data.assign_type;
          this.selectedPageList = this.formatDetailRelationsFn(
            data.relations,
            data.relation_type
          );
        } else if (form.long.relationType == 5) {
          form.activityType = data.belong_type;
          form.activeType = data.assign_type;
          this.selectedActiveList = this.formatDetailRelationsFn(
            data.relations,
            data.relation_type
          );
        }
      }
    } else if (data.strategy_method == 2) {
      form.strategy_method = 2;
    } else if (data.strategy_method == 3) {
      form.strategy_method = 3;
      form.short.planTime = data.start_time
        ? dateFormat(data.start_time, "yyyy-MM-dd hh:mm")
        : "";
    }
    form.phone = data.mobile || "未填写";
    if (data.channels.length > 0) {
      this.selectedChannels = data.channels.map((item: any) =>
        this.getChannelFn(item)
      );
      data.channels.forEach((v: any) => {
        this.sortPolicyList[String(v.order - 1)] = this.getSortPolicyByIdFn(
          v.channel_type
        );
      });
    }
    form.strategyPolicy = data.strategy_policy;
  }
  //设置页面初始化业务数据
  async setPageDataFn() {
    if (!this.itemId) return;
    const params: object = {
      id: this.itemId
    };
    const res: any = await API.marketStrategyDetailApi(params);
    if (res.code == 200) {
      this.setInitDataFn(res.data);
    } else {
      this.$message.error(res.message);
    }
  }
  async setStrategyItemInfoFn(value: number | string) {
    const res: any = await API.groupListApi({});
    if (res.code == 200) {
      const customerGroup: any = res.data.custom;
      customerGroup.forEach((v: any) => {
        if (v.group_id == value) {
          this.strategy.strategyItemInfo.total = v.value;
        }
      });
    } else {
      this.$message.error(res.message);
    }
  }
  created() {
    this.itemId = Number(this.$route.query.id);
    this.setPageDataFn();
  }
}
</script>
<style lang="scss" scoped>
.el-button--primary.is-plain,
.el-button--danger.is-plain {
  background: transparent;
}
.el-button--primary.is-plain:hover,
.el-button--primary.is-plain:focus {
  background: $color-main;
  border-color: $color-main;
  color: #ffffff;
}
.el-button--danger.is-plain:hover,
.el-button--danger.is-plain:focus {
  background: #f67170;
  border-color: #f67170;
  color: #ffffff;
}
.welcome {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  min-width: 910px;
  .item-container {
    .line {
      color: #606266;
      line-height: 40px;
      font-size: 14px;
      margin-bottom: 20px;
      display: flex;
      .line-key {
        vertical-align: middle;
        font-size: 14px;
        box-sizing: border-box;
        margin-right: 20px;
        display: inline-block;
        min-width: 56px;
        width: 70px;
        text-align: right;
      }
      .line-value {
        display: inline-block;

        .like-input {
          vertical-align: middle;
          margin-right: 20px;
          background: $color-background;
          display: inline-block;
          min-width: 200px;
          padding: 0 10px;
          height: 40px;
          color: 000000;
        }
        .like-input-tips {
          color: #606266;
          display: inline-block;
          .main {
            color: #999;
          }
        }
      }
      .line-object {
        .block {
          .list {
            margin: 10px 30px;
          }
          .row {
            display: block;

            .key {
              vertical-align: middle;
              font-size: 14px;
              box-sizing: border-box;
              margin-right: 20px;
              display: inline-block;
              width: 84px;
              text-align: right;
            }
            .value {
              vertical-align: middle;
              font-size: 14px;
              box-sizing: border-box;
              margin-right: 20px;
              display: inline-block;
            }
          }
        }
        .container {
          display: inline-block;
          color: #606266;
          line-height: 40px;
          font-size: 14px;
          vertical-align: middle;
          .model {
            display: inline-block;
            margin-right: 20px;
          }
          .key {
            display: inline-block;
            margin-right: 20px;
          }
          .value {
            display: inline-block;
            margin-right: 20px;
          }
        }
      }
    }
    .flex-row {
      flex-flow: row;
      .line-key {
        flex: 0 0 70px;
      }
      li {
        list-style: none;
      }
      .line-object {
        width: calc(100% - 90px);
        //flex: 1;
        .channel-container {
          .coupon-card {
            width: 100%;
          }
          .other-card {
            width: 380px;
          }
          .title {
            font-weight: bold;
          }
          .cloumn-content {
            width: 290px;
          }
          .cloumn-type {
            text-align: right;
            width: 45px;
          }
        }
        .drag_container {
          width: 450px;
          .drag_header {
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
          .drag_body {
            .cell {
              border-bottom: 1px solid #dcdfe6;
              height: 40px;
              line-height: 40px;
              padding-left: 10px;
            }
            .paddingL20 {
              padding-left: 25px;
            }
          }
        }
        .mode-content {
          border: 1px solid rgb(211, 212, 214);
          padding: 10px 20px;
          margin-top: 10px;
          width: 650px;
        }
        .mark {
          font-size: 14px;
          font-weight: normal;
          color: #999;
        }
      }
    }
  }
  .splitLine {
    height: 0;
    border-top: 1px solid #ebeef5;
    margin-bottom: 20px;
  }
  .headline {
    font-size: 16px;
    font-weight: 700;
    margin-bottom: 30px;
    .mark {
      font-size: 14px;
      font-weight: normal;
      color: #999;
      padding-left: 20px;
    }
  }
}

.like-input {
  vertical-align: middle;
  margin-right: 20px;
  background: $color-background;
  display: inline-block;
  min-width: 200px;
  padding: 0 10px;
  height: 40px;
  color: #000000;
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
<style lang="scss"></style>
