<template>
  <div class="signIn-detail">
    <div class="cy-title">
      基本信息
    </div>
    <div class="tabs-firstrow">
      <el-image fit="cover"
        :src="userInfo.avatar
            ? userInfo.avatar
            : 'https://oss-crs.vchangyi.com/wechat/default-avatar-icon.png'
          "
        class="user-avatar m-b-20"
      >
      </el-image>
      <div class="message-grid">
        <div class="message-item">
          <span>{{ userInfo.name || "-" }}</span>
          <span>签到天数:
            <span class="c-blue">{{ userInfo.day || "-" }}</span></span>
        </div>
        <div class="mt-30">手机号码: {{ userInfo.phone || "-" }}</div>
      </div>
    </div>
    <div class="cy-title">
      签到信息
    </div>
    <el-table
      :data="signInList"
      stripe
      tooltip-effect="light"
    >
      <el-table-column
        label="签到时间"
        min-width="140"
      >
        <template slot-scope="scope">
          <span>
            {{ scope.row.sign_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
          </span>
        </template>
      </el-table-column>
      <el-table-column
        prop="sign_content"
        label="签到内容"
        min-width="140"
      >
      </el-table-column>
      <el-table-column
        label="签到图片"
        width="140"
      >
        <template slot-scope="scope">
          <el-image fit="cover"
            style="width: 50px; height: 50px"
            :src="scope.row.sign_pic"
            @click="checkNowImage(scope.row.sign_pic)"
          ></el-image>
        </template>
      </el-table-column>
      <div
        slot="empty"
        class="table-empty-label-page"
      >
        <img
          class="empty-img"
          :src="require('@/assets/images/empty-box.png')"
          alt=""
        />
        <div>暂无签到信息</div>
      </div>
    </el-table>
    <div class="table_footer">
      <el-pagination
        v-if="signInList && signInList.length > 0"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="meta.page"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="meta.size"
        layout="total, sizes, prev, pager, next, jumper"
        :total="meta.total"
      >
      </el-pagination>
    </div>
    <div class="cy-title">
      收货信息
    </div>
    <el-table
      :data="list"
      stripe
      tooltip-effect="light"
    >
      <el-table-column
        prop="receiver_name"
        label="收货人"
        width="120"
      >
      </el-table-column>
      <el-table-column
        prop="receiver_phone"
        label="收货人手机号"
        width="160"
      >
      </el-table-column>
      <el-table-column
        prop="receiver_address"
        label="收货地址"
        min-width="240"
      >
      </el-table-column>
      <div
        slot="empty"
        class="table-empty-label-page"
      >
        <img
          class="empty-img"
          :src="require('@/assets/images/empty-box.png')"
          alt=""
        />
        <div>暂无收货信息</div>
      </div>
    </el-table>
    <div class="mt-30">
      <el-button
        size="medium"
        @click="goBack"
      >返回</el-button>
    </div>
    <!--图片预览 -->
    <preview
      v-if="preview_show"
      :urls="checkImgArr"
      @close="preview_show = false"
    ></preview>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
// 图片预览
import preview from "@/components/business/preview/index.vue";
import Api from "@/api/wisdom/researchActivities";
@Component({
  components: {
    preview
  }
})
export default class extends Vue {
  // data
  private activity_id = 0;
  private union_id: any = "";
  private userInfo: any = {};
  private preview_show = false;
  private checkImgArr: Array<any> = [];
  private signInList: Array<any> = [];
  private list: Array<any> = [];
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  // 生命周期
  private mounted() {
    this.activity_id = Number(this.$route.query.activity_id);
    this.union_id = this.$route.query.union_id;
    this.getDateil();
  }
  private getDateil() {
    const params = {
      page: {
        page_num: this.meta.page,
        page_size: this.meta.size
      },
      id: this.activity_id,
      union_id: this.union_id
    };
    Api.signInDetail(params).then((res: any) => {
      if (res.code === 200) {
        const data = res.data;
        this.userInfo = {
          avatar: data.avatar,
          name: data.name,
          phone: data.phone,
          day: data.day
        };
        this.signInList = data.sign_in_info_detail;
        this.meta.total = data.count;
        this.list = [
          {
            receiver_name: data.shr_name,
            receiver_phone: data.shr_phone,
            receiver_address: data.shr_address
          }
        ];
      }
    });
  }
  // 每页
  private handleSizeChange(val: any) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getDateil();
  }

  // 当前页
  private handleCurrentChange(val: any) {
    this.meta.page = val;
    this.getDateil();
  }
  // 查看图片
  private checkNowImage(item: string) {
    this.checkImgArr = [item];
    this.preview_show = true;
  }
  private goBack() {
    this.$router.go(-1);
  }
}
</script>
<style lang="scss" scoped>
.signIn-detail {
  background-color: #fff;
  padding: 20px 30px;
  border-radius: 10px;
}

.cy-title {
  background: #f2f2f2;
  height: 40px;
  line-height: 40px;
  padding: 0 10px;
  border-radius: 2px;
  margin: 10px 0 30px 0;
}

.tabs-firstrow {
  display: flex;
  padding: 0 20px;

  .user-avatar {
    width: 70px;
    height: 70px;
    border-radius: 50%;
    margin-right: 10px;
  }

  .message-grid {
    margin-left: 10px;
    width: 100%;

    .message-item {
      display: flex;
      font-size: 18px;
      width: 100%;
      justify-content: space-between;
    }
  }
}

.mt-30 {
  margin-top: 30px;
}

.table-empty-label-page {
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

.table_footer {
  display: flex;
  margin-top: 20px;
  justify-content: flex-end;
}</style>
