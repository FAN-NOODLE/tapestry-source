<template>
  <div class="child-task-detail">
    <layout-list>
      <template slot="top-tabs">
        <el-tabs
          v-model="first_active_name"
          @tab-click="firstHandleClick"
        >
          <el-tab-pane
            label="子任务详情"
            name="1"
          ></el-tab-pane>
          <el-tab-pane
            label="子任务完成详情"
            name="2"
          ></el-tab-pane>
        </el-tabs>
      </template>
      <template slot="search">
        <div class="content-item m-b-20 flex-wrap flex-start">
          <span>子任务名称：{{ planInfo.childTaskName }}</span>
          <span v-if="planInfo.dimension !== '1'">任务模式：{{ planInfo.task_type === 1 ? "单次" : "持续" }}</span>
          <span>条件维度：{{ dimension_list[planInfo.dimension] }}</span>
          <span v-if="planInfo.dimension < '3'">条件内容：{{ conditionName }}</span>
          <span v-if="planInfo.condition_type === 8">任务发送时间：生日提前{{ planInfo.birthday_pre_time }}天</span>
          <span v-if="planInfo.dimension == '2'">任务周期：{{ taskObj.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }}～{{
            taskObj.end_time | formatDate("yyyy-mm-dd hh:mm:ss")
          }}</span>
          <span v-if="ep.isCoach">子任务有效期：{{ (planInfo.task_start_time||0) | formatDate("yyyy-mm-dd hh:mm:ss") }}-{{ (planInfo.task_end_time||0) | formatDate("yyyy-mm-dd hh:mm:ss") }}</span>
        </div>
        <!-- 描述 -->
        <div v-if="ep.isCoach" class="content-item m-b-20">子任务描述：{{ planInfo.remark }}</div>
        <div
          class="content-item flex-start m-b-20"
          style="position: relative;"
          :style="{
            paddingBottom: (planInfo.range.length > 50  ? '40px' : '10px')
          }"
          v-if="planInfo.range.length > 0 || showDownStaff()"
        >
          <!-- todo:展开，折叠 -->
          <div class="range-content-collapse">
            <el-button
              type="text"
              :icon="isCollapse ? 'el-icon-caret-top' : 'el-icon-caret-bottom'"
              @click="isCollapse = !isCollapse"
              v-if="planInfo.range.length > 50"
            >
              {{ isCollapse ? '展开' : '折叠' }}
            </el-button>
          </div>
          <!-- 标题 -->
          <div class="title-item">员工：</div>
          <!-- 可见范围显示区 -->
          <div
            class="range-content"
            v-if="!showDownStaff()"
            :style="[{ height: (planInfo.range.length > 50 && isCollapse ? '120px' : 'auto') },
             { overflow: (isCollapse ? 'hidden' : 'visible') }]"
          >
            <div
              v-for="( rangeItem, rangeIndex ) in  planInfo.range "
              :key="rangeIndex"
              class="range-label"
            >
              <img
                class="member-icon"
                src="@/assets/images/member_normal.png"
                alt=""
              />
              <div class="item-name font-re">
                {{ rangeItem.name }}
              </div>
            </div>
          </div>
          <!-- 导入客户名单 -->
          <div v-else>
            <el-button
              type="text"
              style="color: #6f3112"
              @click="downGuide()"
            >{{
              planInfo.is_extract == 1 ? showDownStaff() : "名单匹配中..."
            }}</el-button>
          </div>
        </div>

        <div class="content-item flex-start m-b-20">
          <div>时间：</div>
          <el-date-picker
            v-model="searchTime"
            type="daterange"
            range-separator="-"
            :clearable="false"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="timestamp"
            :picker-options="pickerOptions"
            :default-time="['00:00:00', '23:59:59']"
            @blur="changeSelector()"
          >
          </el-date-picker>
        </div>
        <div
          class="top-box"
          v-if="planInfo.dimension === '1'"
          style="width:100%;"
        >
          <div
            v-for="( item, index ) in  topTwoList "
            :key="index"
            :class="['content-box width-20', index % 3 == 2 ? 'mar-r-line' : '']"
          >
            <p>
              <span class="title-name">{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
            </p>
            <p class="sum">{{ subTaskData[item.key] }}</p>
          </div>
        </div>
        <div
          class="top-box"
          style="width:100%;"
          v-else
        >
          <div
            v-for="( item, index ) in  topOneList "
            :key="index"
            :class="['content-box  width-20', index % 3 == 2 ? 'mar-r-line' : '']"
          >
            <p>
              <span class="title-name">{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
            </p>
            <p class="sum">{{ subPortraitData[item.key] }}</p>
          </div>
        </div>
        <div v-if="!ep.isCoach" class="data-info m-b-30">
          <div class="time">
            以下数据更新截止时间：{{ updateTime | formatDate("yyyy-mm-dd") }}
            23:59:59
          </div>
        </div>
        <div
          class="top-box flex-nowarp"
          style="width:100%;"
        >
          <div
            v-for="( item, index ) in  dataList "
            :key="index"
            :class="['content-box', 'm-b-0', index % 5 == 4 ? 'mar-r-0' : '']"
            style="width:20%;"
          >
            <p>
              <span class="title-name">{{ item.name }}</span>
              <el-popover
                :ref="`popoverA${index}`"
                placement="top"
                width="270"
                trigger="hover"
                visible-arrow="true"
                :content="item.desc"
              >
                <svg-icon
                  slot="reference"
                  class="help-icon"
                  icon-class="group-help"
                ></svg-icon>
              </el-popover>
            </p>
            <p
              class="sum"
              v-if="planInfo.dimension === '1'"
            >
              <template v-if="item.key !== 'trade_transform'">
                <div>
                  {{ (subTaskData[item.key] || 0) | formatNumber }}
                </div>
              </template>
              <template v-else>
                <div>
                  {{ subTaskData[item.key] || 0 }}
                </div>
              </template>
            </p>
            <p
              class="sum"
              v-if="planInfo.dimension !== '1'"
            >
              <template v-if="item.key !== 'trade_transform'">
                <div>
                  {{ (subPortraitData[item.key] || 0) | formatNumber }}
                </div>
              </template>
              <template v-else>
                <div>
                  {{ subPortraitData[item.key] || 0 }}
                </div>
              </template>
            </p>
          </div>
        </div>
      </template>
    </layout-list>
    <layout-list class="margin-top-30">
      <template slot="top-tabs">
        <el-tabs
          v-model="active_name"
          @tab-click="handleClick"
        >
          <el-tab-pane
            v-for="( item, index ) in  conditionList "
            :key="index"
            :label="item.label"
            :name="item.name"
          ></el-tab-pane>
        </el-tabs>
      </template>
      <template slot="search">
        <div
          v-if="first_active_name === '1'"
          class="child-task-detail"
        >
          <div class="child-yask-item">
            <div
              class="font-re"
              v-if="planInfo.dimension === '1'"
            >
              <div class="title">
                条件规则:
              </div>
              <div class="content">
                {{ conditionName }}{{ planInfo.termArr[active_name - 1].term === "3" ? "小于等于" : "大于等于"
                }}{{ planInfo.termArr[active_name - 1].value }}人
              </div>
            </div>
            <div
              class="font-re"
              v-if="planInfo.dimension == 2"
            >
              <div
                class="content work-break"
                v-if="planInfo.condition_type !== 8 && planInfo.condition_type !== 9"
              >
                标签顾客：
                <span
                  class="label-item"
                  v-for="( items, i ) in  planInfo.termArr[active_name - 1].customerLabel "
                  :key="i"
                >
                  {{ items }}
                </span>
              </div>
            </div>
            <div class="font-re">
              <div class="title">
                子任务内容：
              </div>
              <div class="content">
                {{ planInfo.termArr[active_name - 1].text.content }}
              </div>
            </div>
            <div class="material-container">
              <div
                class="select-picture-box"
                v-if="planInfo.termArr[active_name - 1].currentMaterialType == 'picture'"
              >
                <div
                  class="item-picture"
                  v-for="( pItem, pindex ) in  planInfo.termArr[active_name - 1].content.picture "
                  :key="pindex"
                >
                  <img
                    class="item-img"
                    :src="pItem.url"
                    alt=""
                  />
                </div>
              </div>
              <div
                class="select-video-box"
                v-if="planInfo.termArr[active_name - 1].hasMaterial &&
                  planInfo.termArr[active_name - 1].content.video &&
                  planInfo.termArr[active_name - 1].content.video.url
                  "
              >
                <img
                  @click="openVideo(planInfo.termArr[active_name - 1].content.video.url)"
                  class="video-play-icon"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/play_icon.png"
                  alt=""
                />
                <video
                  width="80"
                  height="80"
                  x5-video-player-fullscreen="true"
                  x5-playsinline
                  playsinline
                  webkit-playsinline
                  muted
                  preload="auto"
                  :src="planInfo.termArr[active_name - 1].content.video.url"
                ></video>
              </div>
              <!-- 卡片 -->
              <div
                class="select-applet-box"
                v-if="planInfo.termArr[active_name - 1].currentMaterialType == 'card'"
              >
                <div class="select-applet-name">
                  <img
                    class="mini-img"
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                    alt=""
                  /><span>贺卡</span>
                </div>
                <div class="select-applet-title">
                  {{ planInfo.termArr[active_name - 1].content.card.template_name }}
                </div>
                <img
                  class="select-applet-img"
                  :src="getCardPictureUrl(planInfo.termArr[active_name - 1].content.card)"
                  alt=""
                />
              </div>
              <!-- 小程序 -->
              <div
                class="select-applet-box"
                v-if="planInfo.termArr[active_name - 1].currentMaterialType == 'applet'"
              >
                <div class="select-applet-name">
                  <img
                    class="mini-img"
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/enterprise/mini-applet-icon.png"
                    alt=""
                  /><span>{{
                    planInfo.termArr[active_name - 1].content.applet.app_name
                  }}</span>
                </div>
                <div class="select-applet-title">
                  {{ planInfo.termArr[active_name - 1].content.applet.title }}
                </div>
                <img
                  class="select-applet-img"
                  :src="planInfo.termArr[active_name - 1].content.applet.cover_url"
                  alt=""
                />
              </div>
              <div
                class="select-link-box"
                v-if="planInfo.termArr[active_name - 1].currentMaterialType == 'link'"
              >
                <div class="select-link-title">
                  {{ planInfo.termArr[active_name - 1].content.link.title }}
                </div>
                <div class="select-link-bottom">
                  <span class="select-link-description">
                    {{ planInfo.termArr[active_name - 1].content.link.description }}
                  </span>
                  <img
                    class="select-applet-img"
                    :src="planInfo.termArr[active_name - 1].content.link.cover_url"
                    alt=""
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-if="first_active_name === '2'">
          <div class="table-container">
            <div
              class="font-new"
              v-if="planInfo.dimension === '1'"
              style="display:flex;"
            >
              <div class="title">
                条件规则:
              </div>
              <div class="content">
                {{ conditionName }}{{ planInfo.termArr[active_name - 1].term === "3" ? "小于等于" : "大于等于"
                }}{{ planInfo.termArr[active_name - 1].value }}人
              </div>
            </div>
            <div
              class="font-new"
              v-if="planInfo.dimension === '2'"
            >
              <div
                class="content work-break"
                v-if="planInfo.condition_type !== 8"
              >
                客户标签：
                <span
                  class="label-item"
                  v-for="( items, i ) in  planInfo.termArr[active_name - 1].customerLabel "
                  :key="i"
                >
                  {{ items }}
                </span>
              </div>
            </div>
            <div class="header-search">
              <search-item
                text="社群名称："
                v-if="planInfo.dimension === '1'"
              >
                <cy-input
                  size="small"
                  width-size="small"
                  placeholder="请输入社群名称"
                  v-model.trim="search.name"
                  prefix-icon="el-icon-search"
                  clearable
                ></cy-input>
              </search-item>
              <search-item
                text="导购名称："
                v-else
              >
                <cy-input
                  size="small"
                  width-size="small"
                  placeholder="请输入导购名称"
                  v-model.trim="search.name"
                  prefix-icon="el-icon-search"
                  clearable
                ></cy-input>
              </search-item>
              <search-item text="所在部门：">
                <el-select
                  class="search-select"
                  placeholder="全部"
                  ref="selectBlur"
                  v-model="search.department.name"
                >
                  <el-option
                    :value="search.department.value"
                    :label="search.department.name"
                  >
                    <el-tree
                      class="treeStyle"
                      :props="search.department.data"
                      :load="loadNode"
                      node-key="dp_id"
                      :expand-on-click-node="true"
                      :check-strictly="true"
                      :check-on-click-node="true"
                      @node-click="handleNodeClick"
                      lazy
                      ref="tree"
                      @current-change="blurFn"
                    ></el-tree>
                  </el-option>
                </el-select>
              </search-item>
              <search-item
                text="触发时间："
                v-if="planInfo.dimension === '1'"
              >
                <el-date-picker
                  v-model="search.taskTime"
                  type="daterange"
                  range-separator="-"
                  start-placeholder="开始时间"
                  end-placeholder="结束时间"
                  :picker-options="taskPickerOptions"
                >
                </el-date-picker>
              </search-item>
              <search-item text="完成时间：">
                <el-date-picker
                  v-model="search.finishTime"
                  type="daterange"
                  range-separator="-"
                  start-placeholder="开始时间"
                  end-placeholder="结束时间"
                  :picker-options="finishPickerOptions"
                >
                </el-date-picker>
              </search-item>
              <search-item text="任务状态：">
                <cy-select
                  v-model="search.state"
                  size="medium"
                  :optionsList="stateList"
                  :options="{ label: 'label', value: 'value' }"
                >
                </cy-select>
              </search-item>
              <el-button
                size="small"
                type="primary"
                @click="searchList"
              >
                查询
              </el-button>
              <el-button
                class="mar-l-20 mar-b-20"
                size="small"
                @click="reset"
                style="margin-left: 10px;"
              >
                重置
              </el-button>
              <el-button
                class="mar-l-20 mar-b-20"
                size="small"
                @click="AssociationExportFile"
                style="margin-left: 10px;"
              >
                任务数据导出
              </el-button>
              <el-button
                v-if="planInfo.dimension - 0 !== 1"
                class="mar-l-20 mar-b-20"
                size="small"
                @click="PortrayalExportFile"
                style="margin-left: 10px;"
              >
                数据明细导出
              </el-button>
            </div>
            <div class="table-content">
              <!-- 社群数据维度列表 -->
              <el-table
                v-if="planInfo.dimension === '1'"
                stripe
                :data="list"
                class="mar-t-20"
              >
                <el-table-column
                  label="任务触发时间"
                  min-width="180"
                >
                  <template slot-scope="scope">
                    <span>
                      {{ scope.row.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column
                  prop="group_chat_name"
                  label="社群名称"
                  min-width="120"
                ></el-table-column>
                <el-table-column
                  prop="user_name"
                  min-width="130"
                  label="社群群主"
                ></el-table-column>
                <el-table-column
                  prop="belonging_area"
                  min-width="130"
                  label="区域"
                ></el-table-column>
                <el-table-column
                  prop="department_name"
                  min-width="180"
                  label="所属部门"
                ></el-table-column>
                <el-table-column
                  label="任务完成状态"
                  min-width="160"
                >
                  <template slot-scope="scope">
                    <span>
                      {{ scope.row.is_finished === 1 ? "已完成" : "未完成" }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column
                  label="完成时间"
                  min-width="180"
                >
                  <template slot-scope="scope">
                    <span>
                      {{ scope.row.completed_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column min-width="200">
                  <template slot="header">
                    <span class="tag-header">触达客户数（去重）</span>
                    <el-tooltip
                      placement="top-start"
                      effect="light"
                    >
                      <div slot="content">
                        <div v-if="planInfo.task_category === 0">
                          <b>【 触达客户数(去重) 】：</b>此处的触达客户代表客户状态被变更为【已联系】的客户数
                        </div>
                        <div v-else>
                          <b>【 触达客户数(去重) 】：此处的触达客户代表可以看到导购发送朋友圈的客户数</b>
                        </div>
                      </div>
                      <svg-icon
                        class="help-icon help-icon-m"
                        icon-class="group-help"
                      ></svg-icon>
                    </el-tooltip>
                  </template>
                  <template slot-scope="scope">
                    <span>
                      {{ (scope.row.total_reach_customer_number || 0) | formatNumber }}
                    </span>
                  </template>
                </el-table-column>
                <!-- <el-table-column min-width="130" label="下单客户数">
                  <template slot-scope="scope">
                    <span>
                      {{
                        (scope.row.total_order_customer_number || 0)
                          | formatNumber
                      }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column min-width="130" label="订单总数">
                  <template slot-scope="scope">
                    <span>
                      {{ (scope.row.total_order_number || 0) | formatNumber }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column min-width="150" label="交易总额（元）">
                  <template slot-scope="scope">
                    <span>
                      {{
                        (scope.row.total_order_fee / 100 || 0) | formatNumber
                      }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column min-width="170" label="交易转化率（%）">
                  <template slot-scope="scope">
                    <span
                      v-if="
                        scope.row.total_order_customer_number > 0 &&
                          scope.row.total_reach_customer_number > 0
                      "
                    >
                      {{
                        Math.floor(
                          (scope.row.total_order_customer_number /
                            scope.row.total_reach_customer_number) *
                            1000
                        ) / 10
                      }}
                    </span>
                    <span v-else>0%</span>
                  </template>
                </el-table-column> -->
                <div
                  slot="empty"
                  class="table-empty-label"
                >
                  <img
                    src="@/assets/images/empty-box.png"
                    alt=""
                  />
                  <p>暂无数据</p>
                </div>
              </el-table>
              <!-- 用户画像维度列表 -->
              <el-table
                v-else
                stripe
                :data="list"
                class="mar-t-20"
              >
                <el-table-column
                  prop="name"
                  width="120"
                  label="导购姓名"
                ></el-table-column>
                <el-table-column
                  prop="department_name"
                  min-width="130"
                  label="所属部门"
                ></el-table-column>
                <el-table-column
                  label="任务状态"
                  width="120"
                >
                  <template slot-scope="scope">
                    <span>
                      {{ scope.row.is_finished === 1 ? "已完成" : "未完成" }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column
                  label="完成时间"
                  width="220"
                >
                  <template slot-scope="scope"><span>{{ scope.row.end_time | formatDate("yyyy-mm-dd") }}</span></template>
                </el-table-column>
                <el-table-column
                  prop="customer_num"
                  label="符合条件客户数(人)"
                  width="200"
                ></el-table-column>
                <el-table-column min-width="200">
                  <template slot="header">
                    <span class="tag-header">触达客户数（去重）</span>
                    <el-tooltip
                      placement="top-start"
                      effect="light"
                    >
                      <div slot="content">
                        <div v-if="planInfo.task_category === 0">
                          <b>【 触达客户数(去重) 】：</b>此处的触达客户代表客户状态被变更为【已联系】的客户数
                        </div>
                        <div v-else>
                          <b>【 触达客户数(去重) 】：此处的触达客户代表可以看到导购发送朋友圈的客户数</b>
                        </div>
                      </div>
                      <svg-icon
                        class="help-icon help-icon-m"
                        icon-class="group-help"
                      ></svg-icon>
                    </el-tooltip>
                  </template>
                  <template slot-scope="scope">
                    <span>
                      {{ (scope.row.total_reach_customer_number || 0) | formatNumber }}
                    </span>
                  </template>
                </el-table-column>
                <!--  && !(planInfo.task_category==0&&planInfo.group_send_type!=2)  -->
                <el-table-column v-if="ep.isCoach" min-width="130" label="顾客触达率">
                  <template slot="header">
                    <span class="tag-header">顾客触达率</span>
                    <el-tooltip
                      placement="top-start"
                      effect="light"
                    >
                      <div slot="content">
                        <b>顾客触达率=触达客户数（去重）/符合条件客户数（人）*100% </b>
                      </div>
                      <svg-icon
                        class="help-icon help-icon-m"
                        icon-class="group-help"
                      ></svg-icon>
                    </el-tooltip>
                  </template>
                  <template slot-scope="scope">
                    <span>{{scope.row.finish_rate}}</span>
                  </template>
                </el-table-column>
                <!-- <el-table-column min-width="130" label="下单客户数">
                  <template slot-scope="scope">
                    <span>
                      {{
                        (scope.row.total_order_customer_number || 0)
                          | formatNumber
                      }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column min-width="130" label="订单总数">
                  <template slot-scope="scope">
                    <span>
                      {{ (scope.row.total_order_number || 0) | formatNumber }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column min-width="140" label="交易总额（元）">
                  <template slot-scope="scope">
                    <span>
                      {{
                        (scope.row.total_order_fee / 100 || 0) | formatNumber
                      }}
                    </span>
                  </template>
                </el-table-column>
                <el-table-column min-width="160" label="交易转化率（%）">
                  <template slot-scope="scope">
                    <span
                      v-if="
                        scope.row.total_order_customer_number > 0 &&
                          scope.row.total_reach_customer_number > 0
                      "
                    >
                      {{
                        Math.floor(
                          (scope.row.total_order_customer_number /
                            scope.row.total_reach_customer_number) *
                            1000
                        ) / 10
                      }}
                    </span>
                    <span v-else>0%</span>
                  </template>
                </el-table-column> -->
                <el-table-column
                  prop="zip"
                  label="操作"
                  width="100"
                  fixed="right"
                >
                  <template slot-scope="scope">
                    <span
                      class="c-blue"
                      @click="viewEvn(scope.row)"
                    >查看 </span>
                  </template>
                </el-table-column>
                <div
                  slot="empty"
                  class="table-empty-label"
                >
                  <img
                    src="@/assets/images/empty-box.png"
                    alt=""
                  />
                  <p>暂无数据</p>
                </div>
              </el-table>
            </div>
          </div>
        </div>
      </template>
      <template slot="footer">
        <el-pagination
          v-if="list && list.length > 0 && page_init && first_active_name === '2'"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          class="page-right font-weight-400"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        ></el-pagination>
        <el-button
          class="range-choose-btn"
          @click="cancel"
        >返回</el-button>
      </template>
    </layout-list>
    <!-- 预览视频弹框 -->
    <open-video
      v-if="isOpenVideo"
      :video_url="videoUrl"
      @closeModel="cancelVideo"
    ></open-video>
    <!-- 标签成员弹窗 -->
    <label-member-dialog
      v-if="isOpenDialog"
      :selectedTerm="new_id"
      :selectedBA="selectedBA"
      @closeDialogBA="closeDialogBA"
    ></label-member-dialog>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
// 视频预览
import OpenVideo from "@/views/myapp/wisdom/marketing-calendar/components/com-open-video.vue";
import personApi from "@/api/setting/setting";
import Api from "@/api/wisdom/planTask";
import labelMemberDialog from "../components/labelMemberDialog.vue";
import { downFilePost } from "@/utils/tools";
import PerModule from "@/store/modules/permission";
import { downloadFile } from "@/utils/download";
import { getPrivateFileUrlApi } from "@/api/enterprise/sensitive-word";
@Component({
  components: {
    OpenVideo,
    labelMemberDialog,
  },
})
export default class extends Vue {
  get ep() {
    return PerModule.enterprise
  }
  get updateTime() {
    return (this.searchTime[1] - 60 * 60 * 24 * 1000) / 1000;
  }
  // todo: 折叠
  private isCollapse = true;
  private dimension_list = ["", "社群数据", "用户画像", "导入顾客名单", "导入顾客&导购匹配关系"];
  private page_init = false;
  private first_active_name = "1";
  private conditionName = "";
  private active_name = "1";
  private active_id = -1;
  private taskObj: any = {};
  private new_id = 0;
  private planInfo: any = {
    task_category: '',
    /** 下载导购名单匹配结果 */
    exportUrl: "",
    range: [],
    childTaskName: "",
    text: "",
    dimension: "1",
    condition: "2",
    is_extract: 0,
    consumer_file_id: 0,
    termArr: [
      {
        term: "5",
        value: "",
        customerLabel: [], //已选标签内容
        customerLabelIds: [], //已选标签的ID
        text: {
          content: "", // 任务内容
        },
        content: {
          // 素材内容
          videoId: "0",
          videoUrl: "",
          videoFlag: false,
          appletId: "0",
          applet: {},
          pictureId: [],
          picture: [],
          video: {},
          link: {},
          card: {},
        },
        hasMaterial: false, // 是否有已选素材
        currentMaterialType: "", // 当前选择素材的类型
      },
    ],
  };
  private searchTime: any = "";
  private topStartTime: any = "";
  private topEndTime: any = "";
  private topOneList = [
    {
      name: "触达导购总人数（去重）",
      desc: "当前子任务 触达的导购人数（去重）",
      key: "touch_ba_count",
    },
    {
      name: "匹配客户总人数（不去重）",
      desc: "统计触达的客户人数（不去重）",
      key: "customer_count",
    },
    {
      name: "匹配客户总人数 (去重）",
      desc: "统计触达的客户人数（去重）",
      key: "customer_distict_count",
    },
    {
      name: "任务总件数",
      desc: "所有导购接收到的所有子任务数量",
      key: "task_count",
    },
    {
      name: "已完成任务件数",
      desc: "所有导购已完成的 所有子任务数量",
      key: "task_complete_count",
    },
    {
      name: "任务完成率",
      desc: "已完成任务数 / 任务总数 * 100%",
      key: "task_complete_rate",
    },
  ];
  private topTwoList = [
    {
      name: "触达群主人数（去重）",
      desc: "统计触达的群主人数（去重）",
      key: "touch_master_count",
    },
    {
      name: "匹配社群数",
      desc: "查询周期内匹配任务条件所能触达的社群总数",
      key: "group_chat_count",
    },
    {
      name: "新增任务数（件）",
      desc: "查询周期内匹配任务条件的新增任务总件数",
      key: "new_task_count",
    },
    {
      name: "任务总数（件）",
      desc: "截止查询日累计匹配的任务件数之和",
      key: "task_count",
    },
    {
      name: "已完成任务（件）",
      desc: "统计每日已完成的任务件数之和",
      key: "task_complete_count",
    },
    {
      name: "任务完成率",
      desc: "已完成任务 / 任务总数 * 100%",
      key: "task_complete_rate",
    },
  ];
  get dataList() {
    if(this.ep.isCoach) {
      const { task_category } = this.planInfo
      return [
        {
          name: "累计触达客户数（去重）",
          desc: ['当前任务，导购匹配客户数中；客户联系状态变更为【已联系】的客户数。', '当前任务，导购匹配客户数中；可以看到导购发送朋友圈的客户数。'][task_category],
          key: "total_reach_customer_number",
        },
      ]
    }
    return [
      {
        name: "累计触达客户数（去重）",
        desc: "当前任务中，导购已完成任务后所触达的客户数（去重）",
        key: "total_reach_customer_number",
      },
    ]
  }
  private subTaskData = {
    group_chat_count: 0, //匹配社群数
    task_complete_count: 0, //匹配社群数
    task_complete_rate: 0, //匹配社群数
    new_task_count: 0, //新增任务件数
    task_count: 0, //累计任务总件数
    touch_master_count: 0, //匹配社群数
    total_reach_customer_number: 0, // 触达客户数
    total_order_customer_number: 0, // 累计下单客户数
    total_order_number: 0, // 累计订单数
    total_order_fee: 0, // 累计交易总额
  };
  private subPortraitData = {
    customer_count: 0, //匹配客户数
    customer_distict_count: 0, //匹配客户数去重
    task_complete_count: 0, //任务完成数
    task_complete_rate: 0, //任务完成率
    task_count: 0, //任务总件数
    touch_ba_count: 0, //触达ba数
    total_reach_customer_number: 0, // 触达客户数
    total_order_customer_number: 0, // 累计下单客户数
    total_order_number: 0, // 累计订单数
    total_order_fee: 0, // 累计交易总额
  };
  private isOpenVideo = false; //是否打开视频
  private videoUrl = ""; // 当前播放视频的url
  private conditionList: Array<any> = [];
  private stateList: Array<any> = [
    {
      label: "全部",
      value: -1,
    },
    {
      label: "已完成",
      value: 1,
    },
    {
      label: "未完成",
      value: 2,
    },
  ];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10, // 每页数量
  };
  private search: any = {
    name: "", // 导购名称
    department: {
      value: "", // 部门id
      list: [], // 部门数据
      name: "",
      isShowTree: false,
      data: {
        label: "name",
        children: [],
        isLeaf: "leaf",
      },
    },
    taskTime: 0, // 任务时间
    finishTime: 0, // 完成时间
    state: -1,
  };
  private list: Array<any> = [];
  private isOpenDialog = false;
  private selectedBA = "";
  private task_id = 0;
  private taskPickerOptions = {
    onPick: (obj: any) => {
      this.search.taskTime = new Date(obj.minDate).getTime();
    },
    disabledDate: (time: any) => {
      if (this.search.taskTime) {
        const day1 = 90 * 24 * 3600 * 1000;
        const maxTime = this.search.taskTime + day1;
        const minTime = this.search.taskTime - day1;
        return time.getTime() >= maxTime || time.getTime() <= minTime;
      }
    },
  };
  private finishPickerOptions = {
    onPick: (obj: any) => {
      this.search.finishTime = new Date(obj.minDate).getTime();
    },
    disabledDate: (time: any) => {
      if (this.search.finishTime) {
        const day1 = 90 * 24 * 3600 * 1000;
        const maxTime = this.search.finishTime + day1;
        const minTime = this.search.finishTime - day1;
        return time.getTime() >= maxTime || time.getTime() <= minTime;
      }
    },
  };
  private pickerOptions = {};

  // 生命周期
  private mounted() {
    const { task_id, taskObj } = this.$route.query;
    this.taskObj = taskObj;
    this.taskObj = JSON.parse(this.taskObj);
    this.searchTime = [this.taskObj.created_at * 1000, Date.now()];
    console.log(this.searchTime);
    
    const nowTime = new Date();
    this.topStartTime = this.taskObj.created_at;
    this.topEndTime = parseInt(Date.now()/1000 + '');
    this.task_id = Number(task_id);
    this.getChildTaskDetail(Number(task_id));
    this.search.department.name = "全部";
    // 部门
    this.fetchStoreList();
    this.pickerOptions = {
      disabledDate: (time: any) => {
        return (
          time.getTime() / 1000 < this.taskObj.created_at - 24 * 3600 ||
          time.getTime() / 1000 > new Date().valueOf() / 1000
        );
      },
    };
  }

  private showDownStaff() {
    let text = "";
    if (this.planInfo.dimension - 0 === 3 || this.planInfo.dimension - 0 === 4) {
      text = "下载导购名单匹配结果";
    }
    return text;
  }
  private async downGuide() {
    if (this.planInfo.is_extract === 1) {
      // 接口请求方式，替换为url下载地址
      // window.open(this.planInfo.exportUrl, "_blank");
      const res: any = await getPrivateFileUrlApi({ url: this.planInfo.exportUrl });
      downloadFile(res.data.url, {
        filename: `${this.planInfo.childTaskName}_匹配结果.csv`,
      });
    }
  }
  // 根据时间查询
  private changeSelector() {
    console.log(this.searchTime);
    this.topStartTime = parseInt(this.searchTime[0] /1000 + '');
    this.topEndTime = parseInt(this.searchTime[1] /1000 + '');
    this.getTotalData(Number(this.planInfo.dimension));
  }
  private getTotalData(type: number) {
    // 获取用户画像数据 type 2
    let api = "subTaskPortrayal",
      key = "subPortraitData";
    if (type === 1) {
      api = "subTaskAssociation";
      key = "subTaskData";
    }
    const params = {
      plan_id: this.taskObj.id, //主任务id
      task_id: this.task_id, //子任务id
      start_time: this.topStartTime, //开始时间
      end_time: this.topEndTime, //结束时间
    };
    Api[api](params).then((res: any) => {
      const { code, data } = res;
      if (code === 200) {
        data.trade_transform =
          data.total_order_customer_number > 0 && data.total_reach_customer_number > 0
            ? `${Math.floor((data.total_order_customer_number / data.total_reach_customer_number) * 1000) / 10}%`
            : "0%";
        data.total_order_fee = data.total_order_fee / 100;
        this[key] = data;
        console.log(this[key], `======优化后${key}`);
      }
    });
  }

  private hasPictureMaterial(image: any) {
    return Array.isArray(image?.imgList) && image.imgList.some((item: any) => item?.pic_url || item?.url);
  }

  private hasVideoMaterial(video: any) {
    return Boolean(video?.video_url || video?.url);
  }

  private hasAppletMaterial(applet: any) {
    return Boolean(applet?.app_id || applet?.app_name || applet?.title || applet?.cover_url || applet?.path);
  }

  private hasLinkMaterial(link: any) {
    return Boolean(link?.link_url || link?.url || link?.title || link?.description || link?.pic_url || link?.cover_url);
  }

  private hasCardMaterial(card: any) {
    const pictureUrl = this.getCardPictureUrl(card);
    return Boolean(card?.template_name || card?.template_content || pictureUrl);
  }

  private getCardPictureUrl(card: any) {
    const pictureUrl = card?.template_picture_url;
    if (Array.isArray(pictureUrl)) {
      return pictureUrl[0] || "";
    }
    return pictureUrl || "";
  }

  // 详情切换
  private firstHandleClick() {
    if (this.first_active_name === "2") {
      this.handleClick();
    }
  }
  // 条件tab切换
  private handleClick() {
    if (this.first_active_name === "2") {
      this.conditionList.forEach((item: any) => {
        if (item.name === this.active_name) {
          this.active_id = item.id;
        }
      });
      this.meta = {
        total: 0, // 总数
        page: 1, // 当前页
        size: 10, // 每页数量
      };
      this.reset();
    }
  }
  // 选择部门
  private handleNodeClick(node: any) {
    this.search.department.name = node.name;
    this.search.department.value = node.id;
    (this.$refs.selectBlur as any).blur();
    this.getConditionComplete();
  }
  // 关闭筛选
  private blurFn() {
    (this.$refs.selectBlur as any).blur();
  }
  //动态加载节点
  private loadNode(node: any, resolve: any) {
    // 一级节点处理
    if (node.level === 0) {
      this.setChecked(this.search.department.list);
      return resolve(this.search.department.list);
    }
    if (node.level >= 1) {
      this.getIndex(node, resolve);
    }
  }
  //判断部门是否有子节点
  private setChecked(arr: any) {
    arr.forEach((item: any) => {
      if (!item.children_count) item.leaf = true;
    });
  }

  // 加载叶子节点数据函数
  private getIndex(node: any, resolve: any) {
    const params = { pid: node.data.third_id };
    personApi.depatListApi(params).then((res: any) => {
      this.setChecked(res.data);
      resolve(res.data);
    });
  }

  // 获取部门
  private fetchStoreList() {
    const _this = this;
    personApi.depatListApi({}).then(
      (res) => {
        const obj: any = { id: 0, name: "全部" };
        res.data.unshift(obj);
        _this.search.department.list = res.data;
      },
      (err) => {
        if (err.error === "ERROR") {
          _this.$message({
            showClose: true,
            message: err.message,
            type: "error",
          });
        }
      }
    );
  }
  // 查询子任务的详情
  private getChildTaskDetail(task_id: number) {
    Api.childTaskDetail({ task_id: task_id }).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        // todo:增加参数
        this.planInfo.remark = data.remark;
        this.planInfo.group_send_type = data.group_send_type;
        this.planInfo.task_category = data.task_category;
        this.planInfo.exportUrl = data.down_load_url;
        this.planInfo.childTaskName = data.name;
        this.planInfo.text = data.remark;
        this.planInfo.dimension = data.dimension + "";
        this.planInfo.birthday_pre_time = data.birthday_pre_time;
        this.planInfo.condition_type = data.condition_type;
        this.planInfo.range = data.plan_task_specify_crowd || [];
        this.planInfo.consumer_file_id = data.consumer_file_id;
        this.planInfo.is_extract = data.is_extract;
        this.planInfo.tags = data.tags;
        if (this.planInfo.range) {
          this.planInfo.range.forEach((item: any) => {
            item.name = item.user_name;
          });
        }
        switch (this.planInfo.condition_type) {
          case 1:
            this.conditionName = "标签";
            break;
          case 2:
            this.conditionName = "群人数";
            break;
          case 3:
            this.conditionName = "群内新增人数";
            break;
          case 7:
            this.conditionName = "退群人数";
            break;
          case 8:
            this.conditionName = "顾客生日";
            break;
          case 9:
            this.conditionName = "全部客户";
            break;
          default:
            break;
        }

        this.planInfo.task_type = data.task_type;
        this.getTotalData(Number(this.planInfo.dimension));
        if (data.plan_task_conditions && data.plan_task_conditions.length >= 1) {
          this.planInfo.termArr = [];
          for (let i = 0; i < data.plan_task_conditions.length; i++) {
            this.planInfo.termArr.push({
              term: "5",
              value: "",
              customerLabel: [], //已选标签name
              customerLabelIds: [], // 已选标签ID
              text: {
                content: "", // 任务内容
              },
              content: {
                // 素材内容
                videoId: "0",
                videoUrl: "",
                videoFlag: false,
                appletId: "0",
                applet: {},
                pictureId: [],
                picture: [],
                video: {},
                link: {},
                card: {},
              },
              hasMaterial: false, // 是否有已选素材
              currentMaterialType: "", // 当前选择素材的类型
            });
            this.conditionList.push({
              label: `子任务条件${i + 1}`,
              name: `${i + 1}`,
              id: data.plan_task_conditions[i].id,
            });
          }
        }
        const termArr = this.planInfo.termArr;
        // 处理标签
        data.plan_task_conditions.forEach((item: any, index: number) => {
          if (this.planInfo.dimension === "1") {
            termArr[index].value = item.conditions.logical.number;
            const expression = item.conditions.logical.expression;
            switch (expression) {
              case "ge":
                termArr[index].term = "5";
                break;
              case "le":
                termArr[index].term = "3";
                break;
              default:
                termArr[index].term = "5";
            }
          } else {
            this.planInfo.tags.forEach((v: any) => {
              termArr[index].customerLabel.push(v);
              // termArr[index].customerLabelIds.push(`${v.group_id},${v.tag_id}`);
            });
          }
          termArr[index].text.content = item.content.text?.content;
          // 包含图片
          if (this.hasPictureMaterial(item.content.image)) {
            termArr[index].currentMaterialType = "picture";
            termArr[index].hasMaterial = true;
            termArr[index].content.picture = item.content.image.imgList;
            termArr[index].content.picture.map((item: any) => {
              return (item.id = item.media_id), (item.url = item.pic_url);
            });
          }
          // 包含视频
          if (this.hasVideoMaterial(item.content.video)) {
            termArr[index].currentMaterialType = "video";
            termArr[index].hasMaterial = true;
            termArr[index].content.video = {
              id: item.content.video.media_id,
              url: item.content.video.video_url,
            };
          }
          // 包含小程序
          if (this.hasAppletMaterial(item.content.applets)) {
            termArr[index].currentMaterialType = "applet";
            termArr[index].hasMaterial = true;
            termArr[index].content.applet = {
              ...item.content.applets,
            };
          }
          // 包含链接
          if (this.hasLinkMaterial(item.content.link)) {
            termArr[index].currentMaterialType = "link";
            termArr[index].hasMaterial = true;
            termArr[index].content.link = {
              ...item.content.link,
              cover_url: item.content.link.pic_url,
            };
          }
          if (this.hasCardMaterial(item.content.card)) {
            termArr[index].currentMaterialType = "card";
            termArr[index].hasMaterial = true;
            termArr[index].content.card = {
              ...item.content.card,
            };
          }
        });

        this.planInfo.task_expire_type = data.task_expire_type
        this.planInfo.task_start_time = data.task_start_time;
        this.planInfo.task_end_time = data.task_end_time;
      }
    });
  }
  // 预览视频
  private openVideo(url: string) {
    this.videoUrl = url;
    this.isOpenVideo = true;
  }
  // 关闭预览视频
  private cancelVideo() {
    this.isOpenVideo = false;
  }
  private handleSizeChange(val: number) {
    this.meta.size = val;
    this.getConditionComplete();
  }
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getConditionComplete();
  }
  private searchList() {
    this.meta.page = 1;
    this.getConditionComplete();
  }
  // 列表和导出的参数
  private getListParams(isExport = false) {
    const obj = {
      id: this.active_id,
      ba_name: this.search.name, // 导购名称
      department_id: this.search.department.value, // 部门编号
      is_finished: this.search.state || -1,
      end_time: Date.parse(`${this.searchTime[1]} 23:59:59`) / 1000 || 0,
      ...(this.search.finishTime && {
        completion_start_time: this.search.finishTime[0] / 1000,
        completion_end_time: this.search.finishTime[1] / 1000,
      }),
      ...(this.search.taskTime &&
        Number(this.planInfo.dimension) === 1 && {
        plan_start_time: this.search.taskTime[0] / 1000,
        plan_end_time: this.search.taskTime[1] / 1000,
      }),
      ...(!isExport && {
        page: this.meta.page,
        size: this.meta.size,
      }),
    };
    return obj;
  }
  // 查询
  private getConditionComplete() {
    const params: any = this.getListParams();
    const apis = ["", "subTaskList", "subPortrayalList", "subPortrayalList", "subPortrayalList"];
    const api = apis[Number(this.planInfo.dimension)];
    api &&
      Api[api](params).then((res: any) => {
        if (res.code === 200) {
          this.page_init = true;
          this.list = res.data.list;
          this.meta.total = res.data.total;
        }
      });
  }
  private async PortrayalExportFile() {
    try {
      const params = this.getListParams(true);
      const task_type = this.planInfo.task_type;
      const res: any = await Api.subTaskPortrayalExport({ ...params, task_type });
      if(res.code ===200){
        this.$confirm("文件生成中，稍后将自动下载，详情记录请前往【下载中心】查看。", "文件下载提示", {
          confirmButtonText: "去下载中心",
          cancelButtonText: "留在当前页",
          type: "success"
        })
          .then(() => {
            this.$router.push({
              path: "/myapp/setting/download-center/list"
            });
          })
          .catch(() => {
            this.$message({
              type: "info",
              message: "已取消前往下载中心"
            });
          });
      }
    } catch (e) {
      console.error(e);
    }
  }
  // 社群任务列表导出
  private AssociationExportFile() {
    const params = this.getListParams(true);
    const arr = [
      {
        api: "subTaskAssociationExport",
        file: "社群任务列表导出",
      },
      {
        api: "subTaskPortrayalListExport",
        file: "用户画像列表数据导出",
      },
      {
        api: "subTaskPortrayalListExport",
        file: "导入顾客列表导出",
      },
      {
        api: "subTaskPortrayalListExport",
        file: "导入顾客&导购匹配关系列表导出",
      },
    ];
    const obj = arr[this.planInfo.dimension - 1];
    downFilePost(Api[obj.api], params, obj.file);
  }
  // 重置
  private reset() {
    this.search.name = "";
    this.search.department.name = "";
    this.search.department.value = 0;
    this.search.taskTime = 0;
    this.search.finishTime = 0;
    this.search.state = "";
    this.getConditionComplete();
  }

  // 查看
  private viewEvn(item: any) {
    this.isOpenDialog = true;
    this.selectedBA = item.user_id;
    this.new_id = item.id;
  }

  // 关闭弹窗
  private closeDialogBA() {
    this.isOpenDialog = false;
  }

  // 返回
  private cancel() {
    this.$confirm("你确认要离开吗？", "取消确认", {
      confirmButtonText: "确认",
      cancelButtonText: "取消",
    }).then(() => {
      this.$router.go(-1);
    });
  }
  private getSearchTime(nowTime: any) {
    const year = nowTime.getFullYear();
    let month = nowTime.getMonth() + 1;
    let day = nowTime.getDate();
    if (month < 10) {
      month = "0" + month;
    }
    if (day < 10) {
      day = "0" + day;
    }
    const nowDate = `${year}-${month}-${day}`;
    return nowDate;
  }
}
</script>
<style lang="scss" scoped>
.child-task-detail {
  .margin-top-30 {
    margin-top: 30px;
  }

  /deep/.search-content {
    background: white;
    padding: 0;
    margin-bottom: 0;
  }
}

.mar-b-20 {
  margin-bottom: 20px;
}

.content-item {
  font-size: 14px;
  display: flex;
  justify-content: space-between;
  line-height: 30px;
  color: #303133;
  background: #f7f7f7;
  padding: 20px 30px;
  font-family: PingFangSC-Medium, PingFang SC;

  &.m-b-20 {
    margin-bottom: 20px;
  }

  &.flex-start {
    justify-content: flex-start;
    align-items: center;
  }

  &.flex-wrap {
    flex-wrap: wrap;

    span {
      margin-right: 50px;
    }
  }

  .title-item {
    line-height: 20px;
    font-size: 14px;
    width: 42px;
  }

  .range-content {
    position: relative;
    width: calc(100% - 42px);
    display: flex;
    justify-content: flex-start;
    align-items: center;
    flex-wrap: wrap;

    .range-label {
      background-color: #e8e8e9;
      border: #d1d2d2 solid 1px;
      display: flex;
      align-items: center;
      padding: 4px 12px;
      font-size: 14px;
      line-height: 20px;
      border-radius: 2px;
      margin-right: 10px;
      margin-bottom: 10px;

      .member-icon {
        width: 16px;
        height: 16px;
        margin-right: 4px;
      }
    }

    &-collapse {
      position: absolute;
      bottom: 5px;
      width: 100%;
      text-align: center;
    }
  }
}

.top-box {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  flex-wrap: wrap;

  &.flex-nowarp {
    flex-wrap: nowrap;
  }

  &.width-60 {
    width: calc((100% - 120px) / 5 * 3 + 90px);

    .content-box {
      width: calc(100% / 3);
      margin-right: 30px;
    }
  }

  .content-box {
    height: 137px;
    background: #f7f7f7;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    margin-bottom: 30px;
    margin-right: 30px;

    &.mar-r-0 {
      margin-right: 0;
    }

    &.m-b-0 {
      margin-bottom: 0;
    }

    &.mar-r-line {
      margin-right: calc((100% - 120px) / 5 * 2 - 30px);
    }

    &.width-20 {
      width: calc((100% - 120px) / 5);
    }

    p {
      text-align: center;
      line-height: 26px;
      font-size: 17px;
    }
  }
}

.c-blue {
  cursor: pointer;
  color: $color-main;
}

.table-empty-label {
  display: flex;
  height: 340px;
  flex-direction: column;

  img {
    width: 124px;
    margin: 80px auto 0;
  }

  div {
    line-height: 20px;
    margin-top: 14px;
  }
}

.child-yask-item {
  width: 55%;
  margin-top: 10px;
  padding-bottom: 40px;

  .font-re {
    display: flex;
    line-height: 38px;

    .title {
      min-width: 92px;
    }
  }

  .content {
    margin-left: 10px;
  }

  .material-container {
    margin: 20px 0 0 20px;
  }
}

.el-select-dropdown__item.selected {
  font-weight: normal !important;
  height: auto !important;
}

.el-scrollbar {
  .el-scrollbar__view {
    .el-select-dropdown__item {
      font-weight: normal !important;
      height: auto;
    }
  }
}

.font-new {
  padding: 5px 5px 5px 15px;
  background: #f6f6f9;
  margin-bottom: 15px;
  font-size: 14px;
  line-height: 38px;
  display: flex;

  .title {
    // min-width: 100px;
    text-align: end;
    margin-right: 15px;
  }
}

.right {
  float: right;
}

.label-item {
  background: #f6f6f9;
  border-radius: 3px;
  border: 1px solid #404040;
  padding: 2px 10px;
  margin-left: 20px;
  color: #404040;
  font-size: 14px;
  min-width: 60px;
  text-align: right;
}

.select-picture-box {
  display: flex;
  flex-wrap: wrap;
  margin-top: 12px;

  .item-picture {
    position: relative;
    width: 80px;
    height: 80px;
    margin-right: 10px;
    margin-bottom: 10px;

    .item-img {
      width: 100%;
      height: 100%;
      border-radius: 4px;
    }

    .delete-img-box {
      width: 20px;
      height: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      position: absolute;
      cursor: pointer;
      right: -10px;
      top: -10px;
      z-index: 2;

      img {
        width: 14px;
        height: 14px;
        display: block;
      }
    }
  }

  .upload-card {
    background-color: $color-white;
    border: 1px solid $color-second-line;
    border-radius: 4px;
    position: relative;
    width: 80px;
    height: 80px;

    &:hover {
      border: 1px solid $color-main;
    }

    .center {
      position: absolute;
      bottom: 0;
      top: 0;
      left: 0;
      right: 0;
      margin: auto;
      text-align: center;
      height: 33px;

      .avatar-uploader-icon {
        font-size: 14px;
        display: block;
        margin-bottom: 5px;
      }

      .tips {
        font-size: 12px;
        color: $color-text-second;
        height: 12px;
        line-height: 12px;
        overflow: hidden;
      }
    }
  }
}

.select-video-box {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 4px;
  margin-top: 12px;

  video {
    border-radius: 4px;
  }

  .video-play-icon {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    top: 0;
    margin: auto;
    width: 25px;
    height: 25px;
    z-index: 99;
    cursor: pointer;
  }

  .delete-video-box {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    cursor: pointer;
    right: -10px;
    top: -10px;
    z-index: 2;

    img {
      width: 14px;
      height: 14px;
      display: block;
    }
  }
}

.select-applet-box {
  width: 165px;
  height: 211px;
  border: solid 1px #dcdfe6;
  border-radius: 4px;
  position: relative;
  padding: 10px;

  .select-applet-name {
    display: flex;
    font-size: 12px;
    color: #707277;
    line-height: 20px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-bottom: 4px;

    .mini-img {
      width: 20px;
      height: 20px;
      margin-right: 5px;
    }
  }

  .select-applet-title {
    font-size: 14px;
    color: #303133;
    line-height: 20px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .select-applet-img {
    width: 145px;
    height: 116px;
    margin-top: 10px;
    border-radius: 4px;
  }

  .delete-applet-box {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    cursor: pointer;
    right: -10px;
    top: -10px;
    z-index: 2;

    img {
      width: 14px;
      height: 14px;
      display: block;
    }
  }
}

.select-link-box {
  width: 225px;
  height: 83px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  position: relative;
  padding: 10px;

  .select-link-bottom {
    display: flex;
    margin-top: 3px;

    .select-link-description {
      width: 155px;
      font-size: 12px;
      color: #909399;
      line-height: 17px;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }

    img {
      width: 40px;
      height: 40px;
      display: block;
      border-radius: 4px;
      margin-left: 10px;
    }
  }

  .delete-link-box {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    cursor: pointer;
    right: -10px;
    top: -10px;
    z-index: 2;

    img {
      width: 14px;
      height: 14px;
      display: block;
    }
  }
}

.work-break {
  word-break: keep-all;
}

/deep/ .layout-list {
  padding-top: 10px;
}

.help-icon {
  margin-left: 5px;
  cursor: pointer;
  width: 13px;
}

.sum {
  font-size: 24px !important;
  font-weight: 600;
  margin-top: 10px;
}

.title-name {
  font-size: 14px !important;
}

.flex {
  display: flex;
}

.data-info {
  .time {
    color: $color-text-normal;
    font-size: 14px;
    line-height: 20px;
  }
}
</style>
