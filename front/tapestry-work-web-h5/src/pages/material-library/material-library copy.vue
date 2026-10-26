<template>
  <div>
    <div
      class="meterial-library-page"
      v-if="!isPreview"
    >
      <!-- 素材上传icon -->
      <div
        v-if="authAdd && !isCoach"
        class="add-icon"
        :class="[!isApplets ? '' : 'bottom']"
        @click="chooseAddType"
      ></div>

      <div class="newHeader">
        <div class="newSearch">
          <div
            class="sel-tag-div"
            @click="openCatePopup"
          >
            <span>素材标签</span>
            <img
              class="down-icon"
              src="@/assets/ico_default.png"
              alt=""
            />
          </div>
          <img
            class="search-icon"
            src="@/assets/material-library/search.png"
            alt=""
          />
          <van-search
            v-model="keyWord"
            :clearable="false"
            left-icon=""
            placeholder="请输入素材关键词"
            @input="search(active)"
          />
          <img
            v-if="keyWord !== ''"
            class="delete-icon"
            src="@/assets/material-library/delete.png"
            alt=""
            @click="deleteKeyWord"
          />
        </div>
        <div
          class="right-screen"
          @click="toSelectName"
        >
          <img
            src="@/assets/material-library/screen.png"
            alt=""
          />
        </div>
      </div>
      <!-- 产品类别浮层 -->
      <van-overlay
        :show="showCatePopup"
        @click="openCatePopup"
        z-index="98"
        :custom-style="{ backgroundColor: 'rgba(255,255,255,0)' }"
      >
      </van-overlay>
      <div
        class="cate-popup"
        v-if="showCatePopup"
      >
        <div class="direction-icon">
          <van-icon
            name="play"
            color="#ffffff"
          />
        </div>
        <ul class="cate-ul">
          <li
            v-for="(item, index) in category_list"
            :key="index"
            class="cate-li"
            :class="[category_ids.indexOf(item.id) > -1 ? 'active' : '']"
            @click="selectCategory(item.id)"
          >
            {{ item.name }}
          </li>
          <!-- <li class="cate-li active">包包</li> -->
        </ul>
        <empty-page
          v-if="!category_list.length"
          picture="noData"
          paddingTop="20"
          paddingBottom="20"
          text="暂无产品类别数据"
        ></empty-page>
        <div class="btn-view">
          <div
            class="btn-li reset"
            @click="resetCate"
          >重置</div>
          <div
            class="btn-li confirm"
            @click="confirmCate"
          >确定</div>
        </div>
      </div>
      <!-- 分组名称 -->
      <div
        class="header-groupName"
        v-if="headerGroupName !== ''"
      >
        <span class="allGroupName">
          {{ headerGroupName }}
          <img
            src="@/assets/material-library/delete-icon.png"
            @click="resetScreen"
            alt=""
          />
        </span>
      </div>
      <!--分组弹窗 -->
      <van-popup
        v-model="showPopup"
        round
        position="bottom"
        class="vPopup"
      >
        <van-sidebar
          v-model="activeSideKey"
          class="vSideBar"
        >
          <van-sidebar-item
            v-for="(item, index) in sildeListGroup"
            :key="index"
            :title="item.name"
            @click="recursionOtherGroupData(item)"
            class="vSideBarItem"
          />
        </van-sidebar>
        <div class="right-tree">
          <cyTree
            :options="groupOption"
            v-if="groupOption.length > 0"
          ></cyTree>
        </div>
        <div class="btn-view">
          <div
            class="btn-li reset"
            @click="groupCancel"
          >取消</div>
          <div
            class="btn-li confirm"
            @click="groupCheck"
          >确定</div>
        </div>
      </van-popup>
      <!-- 分类以及素材 -->
      <div class="tabContant">
        <van-tabs
          v-model="active"
          @change="handleTabChange"
          sticky
          title-inactive-color="#000000"
          background="#f4f7fb"
          color="#000000"
        >
          <!-- 图片 -->
          <van-tab
            title="图片"
            name="image"
          >
            <template #title>
              <div class="tab-text">
                图片
                <div
                  class="num-icon-box"
                  v-if="numAll[2]"
                >
                  <div class="num-text">{{ numAll[2] }}</div>
                </div>
              </div>
            </template>

            <div
              class="tips"
              v-if="active == 'image' && isApplets"
            >
              【操作提示】在预览页面长按图片可保存到本地哦～
            </div>
            <!-- todo:宫格展示 -->
            <div
              class="grid-change"
              v-if="isCoach"
            >
              <div
                class="grid-item"
                @click="changeGrid(1)"
              >
                <van-image
                  width="20"
                  height="20"
                  :src="gridType === 1 ? grid1 : grid2"
                  fit="contain"
                />
              </div>
              <div class="grid-item show-tips">
                <van-image
                  width="20"
                  height="20"
                  :src="gridType === 2 ? grid3 : grid4"
                  fit="contain"
                  @click="changeGrid(2)"
                />
                <van-image
                  v-if="showGridPopover && isFirstOpen"
                  @click="showGridPopover = false"
                  class="icon_popup_tips"
                  src="https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_popup_tips.png"
                  fit="contain"
                ></van-image>
              </div>
            </div>
            <div
              v-if="active == 'image'"
              class="outline-box"
            >
              <van-list
                v-model="isLoading"
                :finished="isFinished"
                :finished-text="displayData.length ? '没有更多了' : ''"
                @load="loadData(false, 'image')"
                :immediate-check="false"
              >
                <!-- todo：layout布局 -->
                <van-row
                  class="card-wrapper"
                  :gutter="gridType === 1 ? 15 : 10"
                >
                  <van-col
                    :span="gridType === 1 ? 12 : 8"
                    v-for="(item, index) in displayData"
                    :key="index"
                  >
                    <div
                      class="card"
                      :style="{
                        'marginBottom': gridType === 1 ? '15px' : '10px'
                      }"
                    >
                      <div
                        v-if="item.is_new == 1"
                        class="new-box"
                      ></div>
                      <div
                        v-if="!isApplets && [0, 2].includes(item.is_deal)"
                        class="sel-tag"
                        @click.stop="checkedItem(item)"
                      >
                        <van-icon :name="getIconSelectedClass(item)" />
                      </div>
                      <div
                        v-if="sentData.findIndex((f) => f.id === item.id) > -1"
                        class="sent-tag"
                        :class="gridType === 1 ? 'sent-tag-two-cols' : 'sent-tag-there-cols'"
                      >
                        已发送
                      </div>
                      <!-- todo:合并数量显示 -->
                      <div
                        class="pic-num"
                        v-if="item.merge_image && item.merge_image.length > 1"
                      >
                        <div class="num-text">{{ item.merge_image.length }}张</div>
                      </div>
                      <van-image
                        lazy-load
                        class="pic"
                        :class="gridType === 1 ? 'two-cols' : 'there-cols'"
                        @click.stop="prewImg(item)"
                        :src="item.url"
                        fit="contain"
                      />
                      <div class="title">
                        <vue-ellipsis
                          :visible-line="2"
                          :text="item.name"
                        >
                        </vue-ellipsis>
                      </div>
                      <div
                        v-if="item.is_deal === 1"
                        class="btn disabled"
                      >
                        压缩中，请稍后
                      </div>
                      <div
                        class="btn"
                        v-else-if="isApplets"
                        @click.stop="prewImg(item)"
                      >
                        预览
                      </div>
                      <!-- todo:合并素材发送 -->
                      <div
                        v-else-if="item.merge_image && item.merge_image.length > 1"
                        class="btn"
                        @click.stop="mergeImageSend(item)"
                      >
                        {{ getSendImageText(item) }}
                      </div>
                      <div
                        v-else
                        class="btn"
                        @click.stop="sendChatMessage(item)"
                      >
                        发送
                      </div>
                    </div>
                  </van-col>
                </van-row>
              </van-list>
              <empty-page
                v-if="!displayData.length && !isLoading"
                picture="noData"
                paddingTop="60"
                paddingBottom="265"
                :text="keyWord || groupVal ? '没有找到相关图片素材' : '暂无图片素材'
                  "
              ></empty-page>
            </div>
          </van-tab>
          <!-- 视频 -->
          <van-tab
            title="视频"
            name="video"
          >
            <template #title>
              <div class="tab-text">
                视频
                <div
                  class="num-icon-box"
                  v-if="numAll[5]"
                >
                  <div class="num-text">{{ numAll[5] }}</div>
                </div>
              </div>
            </template>
            <div
              class="tips"
              v-if="active == 'video' && isApplets"
            >
              【操作提示】在预览页面点击右上角…可保存到本地哦～
            </div>
            <div
              v-if="active == 'video'"
              class="outline-box"
            >
              <van-list
                v-model="isLoading"
                :finished="isFinished"
                :finished-text="displayData.length ? '没有更多了' : ''"
                @load="loadData(false, 'video')"
                :immediate-check="false"
              >
                <!-- todo：layout布局 -->
                <van-row
                  class="card-wrapper"
                  :gutter="15"
                >
                  <van-col
                    :span="12"
                    v-for="(item, index) in displayData"
                    :key="index"
                  >
                    <div
                      class="card"
                      @click="checkedItem(item)"
                    >
                      <div
                        v-if="item.is_new == 1"
                        class="new-box"
                      ></div>
                      <div
                        v-if="!isApplets && [0, 2].includes(item.is_deal)"
                        class="sel-tag"
                      >
                        <van-icon :name="getIconSelectedClass(item)" />
                      </div>

                      <div
                        v-if="sentData.findIndex((f) => f.id === item.id) > -1"
                        class="sent-tag sent-tag-two-cols'"
                      >
                        已发送
                      </div>
                      <van-image
                        @click.stop="openModel(item)"
                        class="pic two-cols"
                        :src="getCoverImage(item.cover_url)"
                        lazy-load
                        fit="contain"
                      />
                      <div class="title">
                        <vue-ellipsis
                          :visible-line="2"
                          :text="item.name"
                        >
                        </vue-ellipsis>
                      </div>
                      <div
                        v-if="item.is_deal === 1"
                        class="btn disabled"
                      >
                        压缩中，请稍后
                      </div>
                      <div
                        class="btn"
                        v-else-if="isApplets"
                        @click.stop="getDownload(item)"
                      >
                        预览
                      </div>
                      <div
                        v-else
                        class="btn"
                        @click.stop="sendChatMessage(item)"
                      >
                        发送
                      </div>
                    </div>
                  </van-col>
                </van-row>
              </van-list>
              <empty-page
                v-if="!displayData.length && !isLoading"
                picture="noData"
                paddingTop="60"
                paddingBottom="265"
                :text="keyWord || groupVal ? '没有找到相关视频素材' : '暂无视频素材'
                  "
              ></empty-page>
            </div>
          </van-tab>
          <!-- 文本 -->
          <van-tab
            title="文本"
            name="text"
          >
            <template #title>
              <div class="tab-text">
                文本
                <div
                  class="num-icon-box"
                  v-if="numAll[1]"
                >
                  <div class="num-text">{{ numAll[1] }}</div>
                </div>
              </div>
            </template>
            <div
              v-if="active == 'text'"
              class="outline-box"
            >
              <van-list
                v-model="isLoading"
                :finished="isFinished"
                :finished-text="displayData.length ? '没有更多了' : ''"
                @load="loadData(false, 'text')"
                :immediate-check="false"
              >
                <div
                  class="txt-box"
                  v-for="(item, index) in displayData"
                  :key="index"
                  @click="checkedItem(item)"
                >
                  <div
                    v-if="item.is_new == 1"
                    class="new-box"
                  ></div>
                  <p>{{ item.content }}</p>
                  <div
                    v-if="!isApplets"
                    class="sel-tag"
                  >
                    <van-icon :name="getIconSelectedClass(item)" />
                  </div>
                  <div class="btn-box clearfix">
                    <div
                      v-if="sentData.findIndex((f) => f.id === item.id) > -1"
                      class="sent-text"
                    >
                      已发送
                    </div>
                    <div
                      class="right-div"
                      v-if="isApplets"
                    >
                      <van-button
                        type="info"
                        @click.stop="copyText(item)"
                      >复制</van-button>
                    </div>
                    <div
                      class="right-div"
                      v-else
                    >
                      <van-button
                        type="info"
                        @click.stop="editText(item)"
                      >编辑</van-button>
                      <van-button
                        type="info"
                        @click.stop="sendChatMessage(item)"
                      >发送</van-button>
                    </div>
                  </div>
                </div>
              </van-list>
              <empty-page
                v-if="!displayData.length && !isLoading"
                picture="noData"
                paddingTop="60"
                paddingBottom="265"
                :text="keyWord || groupVal ? '没有找到相关文本素材' : '暂无文本素材'
                  "
              ></empty-page>
              <van-popup
                v-model="showEdit"
                closeable
                :style="{
                  width: '80%',
                }"
              >
                <div class="edit-popup">
                  <p class="title">编辑</p>
                  <div class="input-view">
                    <van-field
                      v-model="eidtRow.content"
                      rows="2"
                      class="input-field"
                      :autosize="{ maxHeight: 180 }"
                      type="textarea"
                      maxlength="500"
                      :autofocus="true"
                    />
                  </div>

                  <div
                    class="btn"
                    @click="editSend(eidtRow)"
                  >发送</div>
                </div>
              </van-popup>
            </div>
          </van-tab>
          <!-- 文件 -->
          <van-tab
            title="文件"
            name="file"
          >
            <template #title>
              <div class="tab-text">
                文件
                <div
                  class="num-icon-box"
                  v-if="numAll[6]"
                >
                  <div class="num-text">{{ numAll[6] }}</div>
                </div>
              </div>
            </template>

            <div
              class="tips"
              v-if="active == 'file' && isApplets"
            >
              【操作提示】在预览页面点击右上角...可分享给好友哦~
            </div>
            <div
              v-if="active == 'file'"
              class="outline-box"
            >
              <van-list
                v-model="isLoading"
                :finished="isFinished"
                :finished-text="displayData.length ? '没有更多了' : ''"
                @load="loadData(false, 'file')"
                :immediate-check="false"
              >
                <!-- todo：layout布局 -->
                <van-row
                  class="card-wrapper"
                  :gutter="15"
                >
                  <van-col
                    :span="12"
                    v-for="(item, index) in displayData"
                    :key="index"
                  >
                    <div
                      class="card fileCard"
                      :key="index"
                      @click="checkedItem(item)"
                    >
                      <div
                        v-if="item.is_new == 1"
                        class="new-box"
                      ></div>
                      <div
                        v-if="!isApplets"
                        class="sel-tag"
                      >
                        <van-icon :name="getIconSelectedClass(item)" />
                      </div>
                      <div
                        v-if="sentData.findIndex((f) => f.id === item.id) > -1"
                        class="sent-tag file-sent-tag"
                      >
                        已发送
                      </div>
                      <div class="filePic">
                        <van-image
                          class="fileImg two-cols"
                          :src="item.file_type | imgUrl"
                          fit="cover"
                          lazy-load
                        />
                      </div>

                      <div class="title">
                        <vue-ellipsis
                          :visible-line="2"
                          :text="item.name"
                        >
                        </vue-ellipsis>
                      </div>
                      <div
                        class="btn"
                        v-if="isApplets"
                        @click.stop="getDownload(item)"
                      >
                        预览
                      </div>
                      <div
                        v-else
                        class="btn"
                        style="margin:0 10px"
                        @click.stop="sendChatMessage(item)"
                      >
                        发送
                      </div>
                    </div>
                  </van-col>
                </van-row>
              </van-list>
              <empty-page
                v-if="!displayData.length && !isLoading"
                picture="noData"
                paddingTop="60"
                paddingBottom="265"
                :text="keyWord || groupVal ? '没有找到相关文件素材' : '暂无文件素材'
                  "
              ></empty-page>
            </div>
          </van-tab>
          <!-- 链接 -->
          <van-tab
            title="链接"
            name="link"
          >
            <template #title>
              <div class="tab-text">
                链接
                <div
                  class="num-icon-box"
                  v-if="numAll[3]"
                >
                  <div class="num-text">{{ numAll[3] }}</div>
                </div>
              </div>
            </template>
            <div
              v-if="active == 'link'"
              class="outline-box"
            >
              <van-list
                v-model="isLoading"
                :finished="isFinished"
                :finished-text="displayData.length ? '没有更多了' : ''"
                @load="loadData(false, 'link')"
                :immediate-check="false"
              >
                <div
                  class="link-box"
                  v-for="(item, index) in displayData"
                  :key="index"
                  @click="checkedItem(item)"
                >
                  <div
                    v-if="item.is_new == 1"
                    class="new-box"
                  ></div>
                  <div
                    v-if="!isApplets"
                    class="sel-tag"
                  >
                    <van-icon :name="getIconSelectedClass(item)" />
                  </div>
                  <div
                    v-if="sentData.findIndex((f) => f.id === item.id) > -1"
                    class="sent-tag link-sent-tag"
                  >
                    已发送
                  </div>
                  <van-image
                    class="link-img"
                    lazy-load
                    :src="item.cover_url"
                    alt=""
                  />

                  <div class="link-rig-box">
                    <div class="tit ellipsis">
                      {{ item.title }}
                    </div>
                    <div class="con ellipsis">
                      {{ item.description }}
                    </div>
                    <van-button
                      type="info"
                      v-if="!isApplets"
                      @click.stop="sendChatMessage(item)"
                    >发送</van-button>
                  </div>
                </div>
              </van-list>
              <empty-page
                v-if="!displayData.length && !isLoading"
                picture="noData"
                paddingTop="60"
                paddingBottom="265"
                :text="keyWord || groupVal ? '没有找到相关链接素材' : '暂无链接素材'
                  "
              ></empty-page>
            </div>
          </van-tab>
          <!-- 小程序 -->
          <van-tab
            title="小程序"
            name="miniprogram"
          >
            <template #title>
              <div class="tab-text">
                小程序
                <div
                  class="num-icon-box"
                  v-if="numAll[4]"
                >
                  <div class="num-text">{{ numAll[4] }}</div>
                </div>
              </div>
            </template>
            <div
              v-if="active == 'miniprogram'"
              class="outline-box"
            >
              <van-list
                v-model="isLoading"
                :finished="isFinished"
                :finished-text="displayData.length ? '没有更多了' : ''"
                @load="loadData(false, 'miniprogram')"
                :immediate-check="false"
              >
                <div class="miniprogram-card-wrapper">
                  <div
                    class=""
                    :class="[isApplets ? 'h-at' : '', 'miniprogram-card']"
                    v-for="item in displayData"
                    :key="item.id"
                    @click="checkedItem(item)"
                  >
                    <div
                      v-if="item.is_new == 1"
                      class="new-box"
                    ></div>
                    <div
                      v-if="!isApplets"
                      class="sel-tag"
                    >
                      <van-icon :name="getIconSelectedClass(item)" />
                    </div>
                    <div
                      v-if="sentData.findIndex((f) => f.id === item.id) > -1"
                      class="sent-tag miniprogram-sent-tag"
                    >
                      已发送
                    </div>
                    <div class="miniprogram-app-header">
                      <img
                        class="miniprogram-pic"
                        src="@/assets/material-library/mini.png"
                        alt=""
                      />
                      <p class="miniprogram-app-name">{{ item.app_name }}</p>
                    </div>
                    <div class="miniprogram-app-header-title">
                      {{ item.title }}
                    </div>
                    <div class="miniprogram-app-content">
                      <van-image
                        lazy-load
                        class="miniprogram-app-content-img"
                        :src="item.cover_url"
                      />
                    </div>
                    <div
                      v-if="!isApplets"
                      class="miniprogram-app-btn"
                      @click.stop="sendChatMessage(item)"
                    >
                      发送
                    </div>
                  </div>
                </div>
              </van-list>
              <empty-page
                v-if="!displayData.length && !isLoading"
                picture="noData"
                paddingTop="60"
                paddingBottom="265"
                :text="keyWord || groupVal
                  ? '没有找到相关小程序素材'
                  : '暂无小程序素材'
                  "
              ></empty-page>
            </div>
          </van-tab>
        </van-tabs>
      </div>
      <!-- 视频预览 -->
      <div
        class="video-content"
        v-if="showModel"
      >
        <div
          class="mask"
          @click="closeModel"
        ></div>
        <div class="content">
          <video
            :src="videoOptions.src"
            :controls="videoOptions.controls"
            class="video-js vjs-big-play-centered vjs-fluid"
            playsinline="true"
            x-webkit-airplay="allow"
            x5-playsinline
            webkit-playsinline
            object-fit="cover"
            style="width: 100%; max-height: 50vh"
            ref="videoRef"
          ></video>
        </div>
      </div>
      <!-- 组合发送 -->
      <div
        class="bottom-view"
        v-if="!isApplets"
      >
        <div class="button-group">
          <div class="num-left">
            <div>已选：{{ chooseArr.length }}</div>
            <div
              class="max-num-left"
              v-if="chooseArr.length == 9"
            >最多可选9个</div>
          </div>
          <!-- todo:清空已选按钮 -->
          <div
            class="clear-btn"
            @click="clearChoose"
            v-if="chooseArr.length > 0"
          >
            清空已选
          </div>
          <div
            class="button-right"
            :class="[!chooseArr.length ? 'dis-btn' : '']"
            @click="combinationSend"
          >
            {{ chooseArr.length > 1 ? '组合发送' : '发送' }}
          </div>
        </div>
      </div>
      <!-- 添加素材类型选择 -->
      <van-popup
        v-model="showAddModal"
        closeable
        position="bottom"
      >
        <header class="add-modal-header">请选择新增素材类型</header>
        <ul>
          <li
            class="cate"
            v-for="(cate, index) in addCategory"
            :key="index"
            @click="jumpAdd(cate.route)"
          >
            <img :src="cate.link" />
          </li>
        </ul>
      </van-popup>
    </div>
    <!-- 组合发送 -->
    <previewSent
      v-if="isPreview"
      :chooseArr="chooseArr"
      @backPage="backPage"
    ></previewSent>
    <!-- todo: 图片预览弹窗 -->
    <image-preview-popup
      v-if="isCoach && showImagePopup"
      :showPopup.sync="showImagePopup"
      :imgList="previewImgList"
      :chooseArr="chooseArr"
      :isFirstOpen="isFirstOpen"
      :showIndex="0"
      @checkedItem="checkedItem"
      @closePopup="showImagePopup = false"
    />
    <!-- toto:合并素材选择弹窗 -->
    <merge-image-send-popup
      v-if="isCoach && showMergeImageSendPopup"
      :showPopup.sync="showMergeImageSendPopup"
      :merageImgList="merageImgList"
      :chooseArr="chooseArr"
      :isFirstOpen="isFirstOpen"
      @closePopup="closeMergeImageSendPopup"
      @checkedItem="checkedItem"
      @sendChatMessage="sendChatMessage"
      @combinationSend="combinationSend"
      @clearChoose="clearChoose"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { install } from "vs-tree";
import wxUtils from "@/utils/wx";
import "vs-tree/dist/vs-tree.css";
import {
  checkAddApi,
  selectListApi,
  sourceMakeMediaApi,
  getListForStore,
  categoryList,
  getNumApi,
  getCheckFirst,
  reportTimes
} from "@/api/marketing-promotion/index";
import emptyPage from "@/components/base/empty-page.vue";
import previewSent from "@/components/material-library/preview-sent.vue";

import {
  PullRefresh,
  Icon,
  Image,
  Field,
  Sidebar,
  SidebarItem,
  Popup,
  ImagePreview,
  Lazyload,
  Col,
  Row,
  Popover,
  Overlay
} from "vant";
import wx from "weixin-js-sdk";
import cyTree from "@/components/base/cy-tree/cy-tree.vue";
// todo:判断企业
import Helper from "@/utils/helper";
// todo:图片预览弹窗
import imagePreviewPopup from "./components/image-preview-popup.vue"
// todo:合并素材选择弹窗
import mergeImageSendPopup from "./components/merge-image-send.vue"
// todo: 文本省略组件
import { VueEllipsis } from 'vue-ellipsis-component';
Vue.use(Popup);
Vue.use(install);
Vue.use(Field);
Vue.use(Sidebar);
Vue.use(SidebarItem);
Vue.use(Lazyload);
Vue.use(Col);
Vue.use(Row)
Vue.use(Popover)
interface Igr {
  name: string;
  id: number;
  children: null;
  source_lvl: number;
}
Component.registerHooks(["beforeRouteEnter"]);
@Component({
  components: {
    emptyPage,
    previewSent,
    cyTree,
    [Icon.name]: Icon,
    vanPullRefresh: PullRefresh,
    [Image.name]: Image,
    imagePreviewPopup,
    mergeImageSendPopup,
    [Popover.name]: Popover,
    VueEllipsis,
    [Overlay.name]: Overlay
  },
  filters: {
    imgUrl: (value) => {
      const wordBj = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/1.png';
      const excelBj = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/2.png';
      const pptBj = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/3.png';
      const pdfBj = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/4.png';
      const otherBj = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/file/5.png';
      switch (value) {
        case 1:
          return wordBj;
        case 2:
          return excelBj;
        case 3:
          return pptBj;
        case 4:
          return pdfBj;
        default:
          return otherBj;
      }
    },
  },
})
export default class extends Vue {
  // todo:发送按钮置灰
  private sendDisabled = false
  // todo:预览图片数组
  private previewImgList = []
  // todo:显示预览弹窗
  private showImagePopup = false
  // todo:显示合并素材选择弹窗
  private showMergeImageSendPopup = false
  // todo:合并素材选择弹窗图片数组
  private merageImgList = []
  // todo:宫格提示
  private showGridPopover = true
  // todo: 是否是首次打开
  private isFirstOpen = false
  private isCoach = Helper.isCoach()
  private defaultVideoImg = require("@/assets/video-bg.png")
  private activeSideKey = 0; // 当前popup侧边栏激活项
  /** 是否有权限上传素材 */
  private authAdd = false;
  private numAll = {};
  private isApplets: any = false;
  private keyWord = "";
  private groupVal = 0;
  private sildeListGroup: Array<Igr> = [
    { name: "全部分组", id: 0, children: null, source_lvl: 1 },
    // { name: "未分组", id: -1, children: null, source_lvl: 1 },
  ];
  private groupOption: Array<any> = [];
  private stairGroupName = ""; //左侧一级名称
  private headerGroupName = ""; //当前选中的所有分组名称

  private showPopup = false;
  private active = "image";
  //当前页
  private page = 1;
  //是否正在加载
  private isLoading = false;
  //是否加载完成
  private isFinished = false;
  private activeIndex = 0;

  //要显示的数据
  private displayData: Array<any> = [];
  private videoOptions = {
    controls: true,
    src: "", // url地址
  };
  private player = null;
  private showModel = false;
  // 编辑文本
  private showEdit = false;
  private eidtRow = {
    content: "",
  };
  private showCatePopup = false; // 产品类别弹窗是否显示
  private category_ids: Array<any> = []; // 已选中的产品类别
  private category_list: Array<any> = []; // 产品类别列表
  private isPreview = false; // 预览发送

  // 点击搜索
  private search = this.$helper.debounce((value: string) => {
    this.resetData();
    this.loadData(true, value);
  }, 500);
  private chooseArr: Array<any> = []; // 组合发送选中的数据

  //上传素材
  private showAddModal = false;
  /* 支持添加的分类 */
  private addCategory = [
    {
      link: require("@/assets/material-library/text.png"),
      route: "text-library",
    },
    {
      link: require("@/assets/material-library/image.png"),
      route: "image-library",
    },
    {
      link: require("@/assets/material-library/video.png"),
      route: "video-library",
    },
  ];
  // 记录本次已发送素材
  public get sentData(): any {
    return this.$store.getters["materialLibrary/getNowSentData"];
  }
  // todo:四宫格，九宫格切换 1:四宫格 2:九宫格
  private gridType = 1;
  private grid1 = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_4_1.png'
  private grid2 = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_4.png'
  private grid3 = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_9_1.png'
  private grid4 = 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_grade_9.png'
  changeGrid(type: number) {
    this.gridType = type;
    //重新加载数据
    this.loadData(false, this.active);
    if (this.isFirstOpen) {
      this.showGridPopover = true
    }
  }
  // todo: 获取视频图片
  getCoverImage(url: any) {
    const regex = /^(((ht|f)tps?):\/\/)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-z]{2,6}\/?/
    if (regex.test(url)) {
      return url
    } else {
      return this.defaultVideoImg
    }
  }
  // todo:合并素材发送
  mergeImageSend(item: any) {
    if (item.merge_image && item.merge_image.length > 0) {
      // 往数据里面插入merge_id字段
      item.merge_image.forEach((i: any) => {
        i.merge_id = item.merge_id
      })
    }
    this.merageImgList = item.merge_image
    this.showMergeImageSendPopup = true
  }
  // todo：关闭合并素材选择弹窗
  closeMergeImageSendPopup() {
    this.showMergeImageSendPopup = false
  }
  // todo:清空已选
  clearChoose() {
    this.chooseArr = []
    this.sendDisabled = false
    this.$store.commit("materialLibrary/setNowSentData", [])
    sessionStorage.setItem("chooseMaterial", JSON.stringify([]))
  }
  // todo:icon选中class
  getIconSelectedClass(item: any) {
    // 如果是合并图片
    if (item.merge_image && item.merge_image.length > 0) {
      const chooseArrLength = this.chooseArr.filter((i) => {
        return item.merge_image.findIndex((arrItem: any) => arrItem.id === i.id) > -1
      }).length
      if (chooseArrLength < item.merge_image.length && chooseArrLength > 0) {
        return 'https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/half.png'
      } else if (chooseArrLength === item.merge_image.length) {
        return 'checked'
      } else {
        return 'circle'
      }
    } else {
      return this.chooseArr.findIndex((i) => i.id === item.id) > -1 ? 'checked' : 'circle'
    }
  }
  // todo:获取合并素材选中数量
  getSendImageText(item: any) {
    if (item.merge_image && item.merge_image.length > 0) {
      // chooseArr中选中的素材数量
      const chooseArrLength = this.chooseArr.filter((i) => {
        return item.merge_image.findIndex((arrItem: any) => arrItem.id === i.id) > -1
      }).length
      if (chooseArrLength > 0) {
        return `选择发送(${chooseArrLength})`
      } else {
        return `选择发送`
      }
    }
  }
  // 复制
  private copyText(e: any) {
    const content = e.content;
    let oInput = document.createElement("input");
    oInput.value = content;
    document.body.appendChild(oInput);
    oInput.select(); // 选择对象;
    document.execCommand("Copy"); // 执行浏览器复制命令
    oInput.remove();
    this.$toast("复制成功");
  }
  // 下载
  private getDownload(e: any) {
    this.downLoad(e);
  }
  // 下载方法
  private downLoad(e: any) {
    console.log(e);
    // 判断文件名字是否有后缀
    if (e.name.indexOf(".") === -1) {
      // 从url中获取后缀
      let ext = e.url.split(".").pop();
      e.name = e.name + "." + ext;
      console.log(e.name);
    }
    window.wx.invoke("previewFile", {
      url: e.url, // 需要预览文件的地址(必填，可以使用相对路径)
      name: e.name, // 需要预览文件的名称(必填)
      size: e.size, // 需要预览文件的字节大小(必填)
    });
  }
  private popstate() {
    wx.closeWindow(); // 微信网页退出
  }
  //判断被选中的数据的类型
  private fileType(type: any) {
    let text: any = "";
    switch (type) {
      case 1:
        text = "text";
        break;
      case 2:
        text = "image";
        break;
      case 3:
        text = "link";
        break;
      case 4:
        text = "miniprogram";
        break;
      case 5:
        text = "video";
        break;
      default:
        text = "file";
        break;
    }
    return text;
  }

  private closeModel() {
    this.showModel = false;
    this.player.currentTime = 0;
    this.player && this.player.pause();
    this.player = null;
  }
  private openModel(item) {
    if (![0, 2].includes(item.is_deal)) {
      return;
    }
    localStorage.setItem("cacheTime", "0");
    this.videoOptions.src = item.url;
    this.showModel = true;
    setTimeout(() => {
      this.player = this.$refs.videoRef;
      this.player.play();
    }, 500);
  }
  private getWxContent() {
    const self = this;
    if (!window.wx.invoke) return;
    window.wx.invoke("getContext", {}, function (res) {
      if (res.err_msg == "getContext:ok") {
        console.log("getContext", res.entry);
        const entry = res.entry;
        // 从单聊会话的工具栏进入-打开客户画像-客户页面
        if (entry == "single_chat_tools") {
          window.wx.invoke("getCurExternalContact", {}, (res) => {
            if (res.err_msg == "getCurExternalContact:ok") {
              console.log("customer:getCurExternalContact", res.userId);
              self.$store.commit("user/setExternalId", res.userId);
              // 神策数据上报 -- 点击侧边栏
              self.$track.trackClickItem({
                feature: self.$route.path,
                customer_id: res.userId
              })
            } else {
              //错误处理
              console.log("getCurExternalContacterr:" + JSON.stringify(res));
            }
          });
        } else if (entry == "group_chat_tools") {
          // 从群聊会话的工具栏进入-打开客户画像-社群页面
          window.wx.invoke("getCurExternalChat", {}, function (res) {
            if (res.err_msg == "getCurExternalChat:ok") {
              console.log("res.chatId", res.chatId);
              self.$store.commit("common/setGroupChatId", res.chatId);
            } else {
              //错误处理
              console.log("getCurExternalChaterr:" + JSON.stringify(res));
            }
          });
        }
      } else {
        //错误处理
        console.log("getContexterr:" + JSON.stringify(res));
      }
    });
  }
  //分组选择变化
  private handleGroupChange() {
    this.resetData();
    this.loadData(true, this.active);
  }

  // 分组接口
  private resetData() {
    this.page = 1;
    this.isLoading = true;
    this.isFinished = false;
    this.displayData = [];
  }

  // 切换标签
  private handleTabChange(name: string) {
    this.active = name;
    this.resetData();
    this.loadData(true, name);
  }
  // 获取分组的选项
  private async loadGroupOptions() {
    const {
      code,
      data = [],
      message,
    }: IApiResult<Array<any>> = await getListForStore();
    if (code !== 200) {
      console.log(message);
      return;
    }
    let sildeList = this.recursionGroupData(data, 0);
    this.sildeListGroup = [...this.sildeListGroup, ...sildeList];
  }
  private toSelectName() {
    this.showPopup = !this.showPopup;
  }
  // 加载列表数据,增加节流
  private loadData = this.$helper.throttle(async (isSearch?: boolean, activeTab?: string) => {
    this.newLoadData(isSearch, activeTab);
  }, 1000);
  private async newLoadData(isSearch?: boolean, activeTab?: string) {
    //1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件
    let type = 1;
    switch (this.active) {
      case "text":
        type = 1;
        break;
      case "image":
        type = 2;
        break;
      case "link":
        type = 3;
        break;
      case "miniprogram":
        type = 4;
        break;
      case "video":
        type = 5;
        break;
      case "file":
        type = 6;
        break;
    }
    let numTabObj = JSON.parse(localStorage.getItem("numTabObj"));
    if (numTabObj) {
      let numTabObjs = JSON.parse(JSON.stringify(numTabObj));
      numTabObj[type] = 0;
      localStorage.setItem("numTabObj", JSON.stringify(numTabObj));
      this.numAll = numTabObjs;
    }
    const params = {
      type: type,
      keyword: this.keyWord,
      source_group_id: this.groupVal,
      page: this.page,
      size: 12,
      category_ids: this.category_ids,
    };
    const { code, message, data }: IApiResult = await selectListApi(params);
    this.isLoading = false;
    if (code !== 200 || data.total === 0) {
      console.log(message);
      this.isFinished = true;
      return;
    }
    // todo:去掉文件后缀
    if (data.list && data.list.length > 0) {
      data.list.forEach((item: any) => {
        if (item.name) {
          item.source_name = item.name;
          item.name = item.name.split(".")[0];
        }
      });
    }
    //防止快速切换tab时，非当前tab的接口后返回
    if (activeTab !== this.active) {
      return;
    }
    this.page++;
    if (isSearch) {
      this.displayData = data.list;
    } else {
      this.displayData = [...this.displayData, ...data.list];
    }
    this.isFinished = this.displayData.length >= data.total;
  }
  // todo:组合发送
  private combinationSend(receive?: any) {
    //判断receive是否有值，有值则是预览发送
    if (receive && receive.length > 0) {
      //将预览发送的数据赋值给chooseArr
      this.chooseArr = [...receive];
    }
    if (!this.chooseArr.length) {
      this.$toast("请选择要发送的文件");
      return;
    }
    // todo:判断选择的是否只有一个，如果只有一个直接发送
    if (this.chooseArr.length == 1) {
      this.sendChatMessage(this.chooseArr[0], true);
      return;
    }
    this.isPreview = true;
  }
  // 编辑发送
  private editSend(item: any) {
    if (!item.content) {
      this.$toast("文本不能为空");
    } else {
      this.sendChatMessage(item);
    }
  }
  // todo:点击发送
  private async sendChatMessage(item: any, isCombination?: boolean) {
    if (!item) return;
    if (this.sendDisabled && !isCombination) {
      this.$toast("需清空已选素材 才可进行单独发送")
      return
    }
    if (!isCombination) {
      this.$toast.loading({
        message: "生成中...",
        forbidClick: true,
      });
    }
    let msgtype = this.fileType(item.type);
    if (["video", "image", "file"].includes(msgtype)) {
      if (!item.media_id || Date.now() > item.media_expire_time * 1e3) {
        const params: IMakeMediaParams = {
          source_id: item.id,
          source_type: msgtype,
          source_url: item.url,
          source_name: item.source_name,
        };
        const { code, data }: IApiResult = await sourceMakeMediaApi(params);
        if (code == 200) {
          item.media_id = data.media_id;
        }
      }
    }
    let message: any = {
      msgtype,
    };
    switch (msgtype) {
      case "text":
        message.text = {
          content: item.content,
        };
        break;
      case "image":
        message.image = { mediaid: item.media_id };
        break;
      case "link":
        message = {
          msgtype: "news",
          news: {
            link: item.link_url,
            title: item.title,
            desc: item.description,
            imgUrl: item.cover_url,
          },
        };
        break;
      case "video":
        message.video = {
          mediaid: item.media_id,
        };
        break;
      case "miniprogram":
        message.miniprogram = {
          appid: item.app_id,
          title: item.title,
          imgUrl: item.cover_url,
          page: item.path,
        };
        break;
      default:
        message.file = { mediaid: item.media_id };
        break;
    }
    //发送完成之后不进入会话界面
    message.enterChat = true;
    console.log(message, 99999);
    let t = this;
    window.wx.invoke("sendChatMessage", message, function (res) {
      console.log(res);
      if (res.err_msg == "sendChatMessage:ok") {
        //发送成功
        console.log("sendChatMessage:成功回调");
        // 记录本次已发送过的文件

        let sentData = JSON.parse(JSON.stringify(t.sentData));
        sentData.push(item);
        t.$store.commit("materialLibrary/setNowSentData", sentData);
        // 神策数据上报 -- 发送素材
        t.$track.trackAssetSend({
          asset_id: item.id,
          customer_id: t.$store.getters["user/getExternalId"],
          task_id:"普通发送无法获取",
          operation_type:2
        })
        // todo: 上报次数
        reportTimes({
          ids: [item.id],
          type: 2
        })
      } else {
        if (res.err_msg == "sendChatMessage:cancel") {
          return;
        }
        if (msgtype == "miniprogram") {
          t.$toast.fail("企微官方后台未关联小程序");
        }
        if (res.err_code === 1) {
          t.$toast.fail('发送失败,请检查素材');
        }
      }
      // 关闭文本编辑弹窗
      t.showEdit = false;
      t.$toast.clear();
    });
  }
  // todo:预览图片兼容ks&sw
  private prewImg(item: any) {
    if (![0, 2].includes(item.is_deal)) {
      return;
    }
    if (this.isCoach) {
      if (item.merge_image && item.merge_image.length > 0) {
        // 往数据里面插入merge_id字段
        item.merge_image.forEach((i: any) => {
          i.merge_id = item.merge_id
        })
        this.previewImgList = item.merge_image
      } else {
        let arr = [{
          ...item
        }];
        this.previewImgList = arr
      }
      this.showImagePopup = true
    } else {
      let arr = [item.url];
      ImagePreview({
        images: arr,
        startPosition: 0,
      });
    }
  }
  // todo:右上角全选按钮
  private checkedItem(obj: any) {
    console.log(obj);
    // todo: 如果点击了选择，下面发送按钮置灰不可点击
    this.sendDisabled = true;
    if (
      (this.isApplets && this.active == "text") ||
      (this.isApplets && this.active == "link") ||
      (this.isApplets && this.active == "miniprogram")
    ) {
      return;
    }
    //图片 视频压缩中禁止选中
    if (["image", "video"].includes(this.active)) {
      if (![0, 2].includes(obj.is_deal)) {
        return;
      }
    }
    const chooseIndex = this.chooseArr.findIndex((item) => item.id === obj.id);
    const merge_image = obj.merge_image
    if (chooseIndex === -1) {
      if (this.chooseArr.length >= 9) {
        this.$toast("最多可选9个素材");
        return;
      }
      // todo:如果是合并文件存在多张
      if (merge_image && merge_image.length > 0) {
        const mergeIndex = this.chooseArr.findIndex((item) => item.merge_id === obj.merge_id);
        if (mergeIndex > -1) {
          this.deleteMergeImage(merge_image)
          return;
        }
        for (let i = 0; i < merge_image.length; i++) {
          // 需要增加type 2 是图片
          if (this.chooseArr.length >= 9) {
            this.$toast("最多可选9个素材");
            return;
          }
          merge_image[i].type = 2
          this.chooseArr.push(merge_image[i]);
        }
      } else {
        this.chooseArr.push(obj);
      }
    } else {
      // todo:删除合并素材
      if (merge_image && merge_image.length > 0) {
        this.deleteMergeImage(merge_image)
      } else {
        this.chooseArr.splice(chooseIndex, 1);
      }
      // 如果chooseArr为空值说明没选择，发送按钮恢复
      if (this.chooseArr.length == 0) {
        this.sendDisabled = false;
      }
    }
  }
  // todo:删除合并素材逻辑
  deleteMergeImage(merge_image: any) {
    // 选中数量
    const chooseArr = this.chooseArr.filter((i) => {
      return merge_image.findIndex((arrItem: any) => arrItem.id === i.id) > -1
    })
    // 如果是合并文件存在多张,删除所有合并素材
    for (let i = 0; i < chooseArr.length; i++) {
      const mergeIndex = this.chooseArr.findIndex((item) => item.id === chooseArr[i].id);
      this.chooseArr.splice(mergeIndex, 1);
    }
    // 如果chooseArr为空值说明没选择，发送按钮恢复
    if (this.chooseArr.length == 0) {
      this.sendDisabled = false;
    }
  }
  // todo:清空所有合并素材
  clearAllMergeImage() {
    this.chooseArr = this.chooseArr.filter((item) => {
      return !item.merge_id
    })
  }
  private editText(item: any) {
    this.eidtRow = JSON.parse(JSON.stringify(item));
    this.showEdit = true;
  }
  // 清空搜索
  private deleteKeyWord() {
    this.keyWord = "";
    this.resetData();
    this.loadData(true, this.active);
  }
  // 给右侧tree重组数据赋值
  private recursionOtherGroupData(option: any) {
    this.activeSideKey = option.id;
    this.stairGroupName = option.name;
    this.groupOption = [];
    if (option.id == 0 || option.id == -1) {
      this.groupVal = option.id;
      this.headerGroupName = this.screenGroupName(option);
      this.resetData();
      this.loadData(true, this.active);
      this.showPopup = false;
    } else {
      setTimeout(() => {
        if (option.children == null) {
          this.groupOption = [];
        } else {
          this.groupOption = option.children;
        }
      });
    }
    // 选择左侧，重置选择的node
    this.$store.commit("materialLibrary/setNowChecked", {});
  }
  // 点击tree
  private getClick(node: any) {
    this.groupVal = node.id;
    this.headerGroupName = this.screenGroupName(node);
    this.resetData();
    this.loadData(true, this.active);
    this.showPopup = false;
  }
  // 查询现实当前现实的各级名组
  private screenGroupName(data: any) {
    if (data.source_lvl === 1) {
      return this.stairGroupName;
    } else if (data.source_lvl === 2) {
      return `${this.stairGroupName}>${data.name}`;
    } else if (data.source_lvl === 3) {
      let secondGroupName = "";
      this.groupOption.forEach((itemTwo: any) => {
        if (itemTwo.children != null) {
          itemTwo.children.forEach((itemThree: any) => {
            if (itemThree.id == this.groupVal) {
              secondGroupName = itemTwo.name;
              return;
            }
          });
        }
      });
      return `${this.stairGroupName}>${secondGroupName}>${data.name}`;
    } else if (data.source_lvl === 4) {
      let secondGroupName = "";
      let threeGroupName = "";
      this.groupOption.forEach((itemTwo: any) => {
        if (itemTwo.children != null) {
          itemTwo.children.forEach((itemThree: any) => {
            if (itemThree.children != null) {
              itemThree.children.forEach((itemFour: any) => {
                if (itemFour.id == this.groupVal) {
                  threeGroupName = itemThree.name;
                  return;
                }
              });
              secondGroupName = itemTwo.name;
              return;
            }
          });
        }
      });
      return `${this.stairGroupName}>${secondGroupName}>${threeGroupName}>${data.name}`;
    }
  }
  // 重组分组数据
  private recursionGroupData(data: any[], level: number) {
    // sildeList左侧一级分组数据
    let sildeList = [];
    data.forEach((item: any) => {
      let newChildren = this.recursionGroupData(item.children || [], level + 1);
      if (!level) {
        newChildren.unshift({
          name: "全部",
          id: item.id,
          children: null,
          source_lvl: level + 1,
          is_show: item.is_show,
        });
      }

      sildeList.push({
        ...item,
        name: item.group_name,
        children: newChildren,
        source_lvl: level + 1,
      });
    });
    return sildeList;
  }
  // 重置筛选
  private resetScreen() {
    this.headerGroupName = "";
    this.activeSideKey = 0;
    this.groupVal = 0;
    this.groupOption = [];
    this.resetData();
    this.loadData(true, this.active);
    this.$store.commit("materialLibrary/setNowChecked", {});
  }
  // 打开产品类别弹窗
  private openCatePopup() {
    this.showCatePopup = !this.showCatePopup;
    // todo:锁定滚动
    if (this.showCatePopup) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.removeProperty('overflow')
    }
  }
  // 重置产品类别搜索
  private resetCate() {
    this.category_ids = [];
    this.resetData();
    this.loadData(true, this.active);
    this.showCatePopup = false;
    document.body.style.removeProperty('overflow')
  }
  // 产品类别搜索确定
  private confirmCate() {
    this.resetData();
    this.loadData(true, this.active);
    this.showCatePopup = false;
    document.body.style.removeProperty('overflow')
  }
  // 获取产品类别
  private async getCategoryList() {
    const res: any = await categoryList({ page: 1, size: 1999 });
    if (res.code == 200) {
      this.category_list = res.data.list;
    }
  }
  // 获取tab新数量
  private async getNumList() {
    const res: any = await getNumApi();
    if (res.code == 200) {
      let obj = {};
      res.data.forEach((e) => {
        obj[e.type] = e.count_num > 99 ? "99+" : e.count_num;
      });
      let t: any = new Date();
      let time = this.$helper.dateFormat(t / 1000, "yyyy-MM-dd");
      obj["time"] = time;
      let numTabObj = JSON.parse(localStorage.getItem("numTabObj"));
      if (numTabObj && numTabObj.time == time) {
        this.numAll = numTabObj;
      } else {
        this.numAll = obj;
        let objs = JSON.parse(JSON.stringify(obj));
        objs[1] = 0;
        localStorage.setItem("numTabObj", JSON.stringify(objs));
      }
    }
  }
  // 选择产品类别
  private selectCategory(id) {
    const index = this.category_ids.indexOf(id);
    if (index > -1) {
      this.category_ids.splice(index, 1);
    } else {
      this.category_ids.push(id);
    }
  }
  // 取消预览
  private backPage() {
    this.isPreview = false;
  }
  // 分组筛选取消
  private groupCancel() {
    this.showPopup = false;
    if(!this.headerGroupName){
      this.activeSideKey = 0;
      this.groupOption = [];
    }
  }
  // 分组筛选确定
  private groupCheck() {
    const node = this.$store.getters["materialLibrary/getNowChecked"];
    console.log("node :>> ", node);
    if (!node.id) return;
    this.getClick(node);
    // this.isPreview = false;
  }
  /**打开添加素材 类型选择弹窗 */
  private chooseAddType() {
    this.showAddModal = true;
  }
  /**跳转到对应的添加页面 */
  private jumpAdd(path) {
    //添加跳转标识
    this.$router.push({
      name: path,
    });
  }
  /** 查询是否具备添加权限 */
  private checkAddApi() {
    checkAddApi()
      .then((res: IApiResult) => {
        if (res.code === 200) {
          this.authAdd = !!res.data.auth;
        }
      })
      .catch((err) => {
        console.error("checkAddApi", err);
      });
  }
  // todo:判断是否是首次打开本页面
  private async isFirst() {
    try {
      const res: any = await getCheckFirst({});
      if (res.code == 200) {
        this.isFirstOpen = res.data === 0 ? true : false;
      }
    } catch (error) {
      console.log(error);
    }
  }
  /**生命周期 */
  private async created() {
    // 当前需要直接退出的页面
    window.history.pushState("_android_back", "", location.href);
    window.addEventListener("popstate", this.popstate, false);
    this.checkAddApi();
    this.isFirst();
  }
  private async destroyed() {
    this.player && this.player.dispose();
    this.player = null;
    // 记得清除!!!!不然单页面，其它页面也可以用
    window.removeEventListener("popstate", this.popstate, false);
  }

  private async mounted() {
    console.log(location.href, 111);
    // 非侧边栏进入 不调用sdk方法
    try {
      this.getWxContent();
    } catch (e) {
      await wxUtils.wxAgentConfig();
      this.getWxContent();
    }
    // this.getWxContent();
    this.loadGroupOptions();
    this.resetData();
    this.loadData(true, this.active);
    this.getCategoryList();
    if (sessionStorage.getItem("chooseMaterial")) {
      this.chooseArr = JSON.parse(sessionStorage.getItem("chooseMaterial"));
    }
    this.getNumList();
    if (this.$route.query.is_applets) {
      this.isApplets = this.$route.query.is_applets;
      // todo  消息进入H5隐藏右上角功能
      window.wx.invoke("hideAllNonBaseMenuItem");
    }
  }
  /** 路由生命周期 */
  private beforeRouteEnter(to, form, next) {
    console.log("beforeRouteEnter", form);
    next((vm) => {
      if (form.name === "text-library") {
        vm.active = "text";
      }
      if (form.name === "image-library") {
        vm.active = "image";
      }
      if (form.name === "video-library") {
        vm.active = "video";
      }
    });
  }
}
</script>
<style lang="scss" scoped>
.grid-tips {
  padding: 15px;
  font-size: 12px;
  font-weight: 500;
  color: #FFFFFF;
}

.meterial-library-page {
  background: #f4f7fb;
  min-height: 100vh;

  .add-icon {
    width: 36px;
    height: 36px;
    background: url("../../assets/material-library/add-icon.png") no-repeat;
    background-size: 100% 100%;
    position: fixed;
    bottom: 126px;
    right: 12px;
    z-index: 2;

    &.bottom {
      bottom: 76px;
    }
  }

  .add-modal-header {
    height: 20px;
    font-size: 14px;
    font-family: PingFangSC-Medium, PingFang SC;
    font-weight: 700;
    color: #000000;
    line-height: 20px;
    margin: 16px auto 20px;
    text-align: center;
  }

  .cate {
    padding: 4px 12px 0;
  }

  .tips {
    font-size: 12px;
    font-weight: 400;
    color: #a6a6a6;
    line-height: 17px;
    text-align: center;
    margin-top: 15px;
  }

  .grid-change {
    position: relative;
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 65px;
    margin: 15px 0 5px 20px;

    .grid-item {
      display: flex;
      align-items: center;

      &.active {
        background: #4a4a4a;
        color: #ffffff;
      }

      img {
        width: 20px;
        height: 20px;
        margin-right: 4px;
      }
    }

    .show-tips {
      position: relative;

      .icon_popup_tips {
        position: absolute;
        width: 182px;
        height: 56px;
        z-index: 999;
        left: -90%;
        top: 70%;
      }
    }
  }

  .video-content {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    z-index: 99;

    .mask {
      position: absolute;
      width: 100%;
      height: 100%;
      background: rgba($color: #000000, $alpha: 0.6);
      z-index: 2;
    }

    .content {
      position: absolute;
      width: 100%;
      top: 0;
      max-height: 50vh;
      z-index: 8;
      background: white;
      font-size: 0;
    }
  }

  .bottom-view {
    position: fixed;
    bottom: -1px;
    width: 100%;
    height: 50px;
    left: 0;
    padding: 8px;
    background-color: #ffffff;
    font-size: 14px;
    z-index: 9;
    box-shadow: 0px 1px 6px 0px rgba(0, 0, 0, 0.15);

    .button-group {
      position: relative;
      display: flex;
      align-items: center;
      width: 100%;
      height: 100%;

      .num-left {
        position: absolute;
        left: 0;

        div {
          margin: 0 4px;
          font-size: 16px;
        }

        .max-num-left {
          font-size: 10px;
          color: #E02020;
        }
      }

      .clear-btn {
        position: absolute;
        right: 180px;
        font-size: 12px;
        font-weight: 500;
        color: #000000;
        text-decoration: solid underline #000000;
      }

      .button-right {
        position: absolute;
        right: 0;
        width: 165px;
        height: 34px;
        color: #ffffff;
        line-height: 34px;
        text-align: center;
        background-color: #000;

        &.dis-btn {
          background-color: #cccccc;
        }
      }
    }
  }
}

.edit-popup {
  padding-bottom: 20px;

  .title {
    font-size: 16px;
    font-weight: 500;
    padding: 15px 0;
    border-bottom: 1px solid #f0f0f0;
    width: 90%;
    margin: 0 auto;
  }

  .input-view {
    width: 90%;
    margin: 0 auto;
    padding: 15px 0;

    // max-height: 60vh;
    // overflow: auto;
    .input-field {
      font-size: 13px;
      padding: 0;

      &::after {
        border: 0 none;
      }
    }
  }

  .btn {
    margin: 15px auto 0;
    width: 90%;
    line-height: 32px;
    font-size: 14px;
    font-weight: 500;
    color: #4a4a4a;
    background-color: transparent;
    border: 1px solid #4a4a4a;
    text-align: center;
    border-radius: 8px;

    &.disabled {
      font-size: 12px;
      border: 1px solid #CCCCCC;
      color: #CCCCCC;
    }
  }
}

.newHeader {
  padding: 15px 14px 0 14px;
  display: flex;
  align-items: center;
  background: #f4f7fb;

  .newSearch {
    width: 280px;
    height: 50px;
    background: #ffffff;
    box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
    margin-right: 13px;
    padding: 16px;
    display: flex;
    align-items: center;

    // justify-content: space-between;
    .sel-tag-div {
      font-size: 12px;
      width: 100px;
      margin-right: 8px;
      border-right: 1px solid #f0f0f0;

      .down-icon {
        width: 8px;
        margin-left: 4px;
        transform: translateY(-2px);
      }
    }

    .search-icon {
      display: block;
      width: 16px;
      height: 17px;
      margin-right: 10px;
    }

    .delete-icon {
      display: block;
      width: 12px;
      height: 12px;
    }

    .input {
      width: 85%;
      font-size: 13px;
      font-weight: 600;
      color: #000000;
      line-height: 18px;
    }
  }

  .van-search {
    width: 85%;
    background: pink;
    padding: 0px;
    font-size: 13px;
    font-weight: 600;
    color: #000000;
  }

  ::placeholder {
    font-size: 13px;
    font-weight: 100;
    color: #a1a1a1;
  }

  ::v-deep.van-search__content {
    background-color: #ffffff;
    padding-left: 0px;
    border-radius: 0;
  }

  .right-screen {
    width: 50px;
    height: 50px;
    display: flex;
    justify-content: center;
    align-items: center;
    background: #ffffff;
    box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);

    img {
      width: 16px;
      height: 16px;
    }
  }
}

.header-change {
  padding: 15px 14px;
  // display: flex;
  align-items: center;
  // height: 36px;
  // margin-bottom: 6px;
  background: #ffffff;

  ::v-deep.van-search {
    // width: 50%;
    // margin-right: 10px;
    margin-bottom: 10px;
    padding: 0px;
  }

  ::v-deep.van-search__content {
    background-color: #f4f5f8;
    border-radius: 4px;
  }

  .right-name {
    // width: 50%;
    height: 36px;
    // background-color: #f4f5f8;
    background-color: #0095f6;
    border-radius: 4px;
    color: #ffffff;
    font-size: 14px;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  ::v-deep.van-dropdown-menu {
    display: flex;
    align-items: center;
    height: 36px;
    flex: 1;
  }

  ::v-deep.van-dropdown-menu__title {
    font-size: 14px;
  }

  ::v-deep.van-ellipsis {
    color: #333333;
  }

  ::v-deep.van-dropdown-menu__bar {
    width: 100%;
    height: 100%;
    box-shadow: none;
    background: #f4f5f8;
    border: none;
    border-radius: 4px;
  }
}

.tabContant {
  .ellipsis {
    overflow: hidden;
    width: 100%;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  ::v-deep.van-tabs__line {
    width: 30px;
    height: 3px;
    background: #000000;
    border-radius: 2px;
  }

  ::v-deep.van-tab--active {
    color: black !important;
    font-weight: bold;
  }

  ::v-deep.van-tabs__nav--line.van-tabs__nav--complete {
    background-color: #f4f5f8;
  }

  // 文本
  .outline-box {
    padding: 10px 15px 50px 15px;
    overflow: hidden;

    &.p-r-0 {
      padding-right: 0;
    }

    &:hover {
      border-color: #0095f6;
    }

    .sel-tag {
      position: absolute;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 28px;
      height: 28px;
      box-shadow: -2px 2px 8px 0px rgba(0, 0, 0, 0.1);
      border-radius: 0px 0px 0px 10px;
      right: 0;
      top: -1px;
      background-color: #ffffff;
      z-index: 1;

      .van-icon {
        font-size: 16px;
        color: #4a4a4a;
      }
    }

    .sent-tag {
      position: absolute;
      background-color: rgba(54, 54, 54, 0.49);
      color: #ffffff;
      text-align: center;
      line-height: 28px;
      font-size: 14px;
      font-weight: 400;
      left: 7px;
      z-index: 1;
    }

    .sent-tag-two-cols {
      top: 130px;
      width: 150px;
    }

    .sent-tag-there-cols {
      width: 94px;
      top: 72px;
    }

    .txt-box {
      position: relative;
      background: #ffffff;
      box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
      margin-bottom: 20px;
      padding: 40px 15px 0;

      .btn-box {
        padding-bottom: 17px;

        .sent-text {
          float: left;
          font-size: 14px;
          color: #96979e;
          line-height: 32px;
        }

        .right-div {
          float: right;

          .van-button {
            height: 28px;
            width: 64px;
            background: #ffffff;
            border: 1px solid #4a4a4a;
            font-size: 13px;
            color: #4a4a4a;
            line-height: 28px;
            margin-left: 12px;
            border-radius: 0;

            &::before {
              border-color: #4a4a4a;
            }
          }
        }
      }

      p {
        margin: 0 18px 10px 0;
        font-size: 14px;
        color: #000000;
        line-height: 20px;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 4;
        -webkit-box-orient: vertical;
      }
    }
  }

  .miniprogram-card-wrapper {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;

    .miniprogram-card {
      width: 162px;
      height: 264px;
      background: #ffffff;
      box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
      padding: 10px 12px 16px;
      margin-bottom: 20px;
      position: relative;

      .miniprogram-sent-tag {
        width: 137px;
        top: inherit;
        bottom: 57px;
        left: 12px;
      }

      .miniprogram-app-header {
        display: flex;
        overflow: hidden;

        .miniprogram-pic {
          flex: 0 0 20px;
          display: block;
          width: 20px;
          height: 20px;
          object-fit: cover;
        }

        .miniprogram-app-name {
          margin: 0;
          overflow: hidden;
          flex: auto;
          height: 20px;
          line-height: 20px;
          margin-left: 5px;
          font-size: 13px;
          color: #000000;
          text-overflow: ellipsis;
          white-space: nowrap;
          font-weight: bold;
        }
      }

      .miniprogram-app-header-title {
        min-height: 36px;
        margin: 4px 0;
        font-size: 12px;
        line-height: 18px;
        color: #777777;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        word-break: break-all;
      }

      .miniprogram-app-content {
        .miniprogram-app-content-img {
          width: 100%;
          height: 132px;
          border-radius: 8px;
          display: block;
        }
      }

      .miniprogram-app-btn {
        margin-top: 16px;
        width: 100%;
        height: 32px;
        cursor: pointer;
        background: transparent;
        line-height: 32px;
        color: #4a4a4a;
        text-align: center;
        font-size: 12px;
        font-weight: 500;
        border: 1px solid #4a4a4a;

        &.disabled {
          font-size: 12px;
          border: 1px solid #CCCCCC;
          color: #CCCCCC;
        }
      }
    }

    .h-at {
      height: auto;
      padding: 14px 15px 15px;
    }
  }

  .card-wrapper {
    -webkit-overflow-scrolling: touch;

    &>* {
      transform: translateZ(0px);
    }

    .card {
      position: relative;
      margin-bottom: 15px;
      background: #ffffff;
      box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
      padding-bottom: 7px;

      .file-sent-tag {
        width: 150px;
        bottom: 86px;
      }

      .pic {
        position: relative;
        margin-bottom: 8px;
        background: #f0f0f0;

        &.video {
          background: #d8d8d8;
          font-size: 32px;
        }
      }

      .pic-num {
        position: absolute;
        top: 15px;
        left: 15px;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 17px;
        background: rgba(0, 0, 0, 0.5);
        border-radius: 9px;
        font-size: 10px;
        font-weight: 400;
        color: #FFFFFF;
        z-index: 9;
        padding: 0 5px;
      }

      .two-cols {
        width: 150px;
        height: 150px;
        margin: 7px;
      }

      .there-cols {
        width: 94px;
        height: 94px;
        margin: 6px 6px 0 7px;
      }

      .title {
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        color: #000000;
        height: 32.6px;
        overflow: hidden;
        padding: 0 11px;
        text-align: center;
      }

      .btn {
        height: 32px;
        display: flex;
        justify-content: center;
        align-items: center;
        cursor: pointer;
        border: 1px solid #000000;
        color: #000000;
        border-radius: 0;
        font-size: 13px;
        margin: 14px 7px 0px;

        &.disabled {
          font-size: 12px;
          border: 1px solid #CCCCCC;
          color: #CCCCCC;
        }
      }
    }

    .fileCard {
      padding-bottom: 15px;

      .filePic {
        padding: 25px 0 17px;

        .fileImg {
          display: block;
          width: 71px;
          height: 94px;
          margin: 0 auto;
        }
      }
    }
  }

  // 链接
  .link-box {
    width: 345px;
    padding: 15px;
    background: #ffffff;
    box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
    display: flex;
    margin-bottom: 20px;
    position: relative;

    .link-sent-tag {
      top: inherit;
      bottom: 15px;
      width: 65px;
      left: 15px;
    }

    .link-content {
      display: flex;
    }

    .link-img {
      width: 66px;
      height: 66px;
      display: block;
      margin-right: 13px;
      margin-top: 2px;
      border-radius: 8px;
    }

    .link-rig-box {
      width: 236px;
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      position: relative;

      .tit {
        margin-bottom: 5px;
        font-size: 13px;
        font-weight: bold;
        color: #000000;
        line-height: 20px;
      }

      .con {
        margin-bottom: 5px;
        line-height: 20px;
        font-size: 12px;
        color: #777777;
      }

      .van-button {
        height: 28px;
        width: 64px;
        background: #ffffff;
        border-radius: 0;
        border: 1px solid #4a4a4a;
        font-size: 13px;
        line-height: 28px;
        color: #4a4a4a;
        position: absolute;
        bottom: 0;
        right: 0;
      }
    }
  }

  // 文件
  .file-box {
    width: 345px;
    padding: 15px;
    border: 1px solid #f0f0f0;
    display: flex;
    margin-bottom: 15px;
    position: relative;

    .van-image {
      width: 39px;
      height: 44px;
    }

    .link-rig-box {
      margin-left: 15px;

      .tit {
        float: left;
        width: 151px;
        font-size: 14px;
        font-weight: 500;
        color: #4a4a4a;
        margin-right: 10px;
        line-height: 44px;
      }

      .van-button {
        float: right;
        border-radius: 0;
        border: 0 none;
        width: 82px;
        height: 28px;
        margin-right: 15px;
        background: #4a4a4a;
        font-size: 12px;
        font-weight: 500;
        margin-top: 8px;
      }
    }
  }
}

.clearfix:after {
  display: block;
  clear: both;
  content: ""; // content 属性与 :before 及 :after 伪元素配合使用，来插入生成内容。
  visibility: hidden; // 属性规定元素是否可见
  height: 0;
}

.clearfix {
  zoom: 1;
}

.header-groupName {
  padding: 0px 24px 0px 16px;
  margin: 10px 0 6px 0;
  font-size: 12px;
  color: #ffffff;
  line-height: 18px;

  .allGroupName {
    display: inline-block;
    padding: 5px 15px;
    background: #4a4a4a;
    border-radius: 24px;
    position: relative;

    img {
      width: 18px;
      height: 18px;
      position: absolute;
      right: -9px;
      top: 50%;
      margin-top: -9px;
    }
  }
}

.vPopup {
  height: 76vh;
  padding: 12px 0 0 0;
  display: flex;
  border-radius: 0;

  .btn-view {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    display: flex;
    background-color: #ffffff;
    padding: 15px;

    .btn-li {
      flex: 2;
      border: 1px solid #4a4a4a;
      color: #4a4a4a;
      font-size: 14px;
      text-align: center;
      line-height: 36px;
    }

    .confirm {
      margin-left: 15px;
      color: #ffffff;
      background-color: #4a4a4a;
    }
  }

  .right-tree {
    height: calc(76vh - 86px);
    overflow: auto;
    width: 100%;
  }

  .vSideBar {
    width: 112px;
    background: #f4f7fb;
    margin-top: -12px;
    overflow: auto;
    height: calc(76vh - 80px);

    .vSideBarItem {
      height: 60px;
      display: flex;
      align-items: center;

      .van-sidebar-item__text {
        letter-spacing: 0;
        overflow: hidden;
        display: -webkit-box;
        text-overflow: ellipsis;
        -webkit-line-clamp: 2;
        /*要显示的行数*/
        -webkit-box-orient: vertical;
        font-size: 14px;
      }
    }

    .van-sidebar-item--select {
      font-weight: 600;
      color: #000000;
      border-radius: 0;
    }

    .van-sidebar-item--select::before {
      display: none;
    }
  }

  .vsTree {
    width: calc(100% - 112px);
    max-height: 100% !important;
    padding: 0px 13px;

    ::v-deep.vs-expand-icon::after {
      width: 16px;
      height: 16px;
      background-size: 100% 100%;
      background-image: url("../../assets/material-library/show-icon.png");
      margin-top: -2px;
      margin-right: 3px;
    }

    ::v-deep.vs-tree-text {
      margin-left: 3px;
      font-size: 14px;
      color: #777777;
      line-height: 20px;
    }

    ::v-deep.vs-tree-node {
      height: 40px;
      color: #777777;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    ::v-deep.vs-tree-node:hover {
      height: 40px;
      // border-radius: 12px;
    }

    ::v-deep.selected {
      height: 40px;
      background: rgba(245, 244, 247, 1);
      // border-radius: 12px;
      font-weight: 900;
      padding-left: 5px;
    }

    ::v-deep.expand-empty {
      display: none;
    }

    ::v-deep.selected .vs-tree-text {
      color: #666668;
    }
  }
}

.cate-popup {
  background: #ffffff;
  padding: 15px 0;
  position: fixed;
  top: 90px;
  left: 15px;
  width: 343px;
  z-index: 99;
  box-shadow: 0px 6px 16px 0px rgba(73, 75, 77, 0.3);

  .direction-icon {
    position: absolute;
    top: -18px;
    left: 30px;
    transform: rotate(-90deg);
    font-size: 24px;
  }

  .cate-ul {
    display: flex;
    flex-wrap: wrap;
    padding: 0 15px 0 5px;
    max-height: 400px;
    overflow: auto;

    .cate-li {
      margin-bottom: 10px;
      margin-left: 10px;
      padding: 6px 15px;
      color: #4a4a4a;
      font-size: 10px;
      background-color: #f4f6f8;

      &.active {
        color: #ffffff;
        background-color: #4a4a4a;
      }
    }
  }

  .btn-view {
    display: flex;
    padding: 15px 15px 0;

    .btn-li {
      flex: 2;
      border: 1px solid #4a4a4a;
      color: #4a4a4a;
      font-size: 14px;
      text-align: center;
      line-height: 36px;
    }

    .confirm {
      margin-left: 15px;
      color: #ffffff;
      background-color: #4a4a4a;
    }
  }
}

.num-icon-box {
  border-radius: 50%;
  width: 14px;
  height: 14px;
  position: absolute;
  top: -1px;
  right: -14px;
  background: #000000;
  font-size: 12px;
  font-weight: 400;
  color: #ffffff;
  display: flex;
  justify-content: center;
  align-items: center;

  .num-text {
    transform: scale(0.6);
  }
}

.new-box {
  width: 39px;
  height: 39px;
  position: absolute;
  top: -13px;
  left: -13px;
  z-index: 1;
  background: url('https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/new.png') no-repeat;
  background-size: 100% 100%;
}

.tab-text {
  position: relative;
}</style>
