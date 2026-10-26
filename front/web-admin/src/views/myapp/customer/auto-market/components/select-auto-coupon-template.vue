<template>
  <el-dialog
    title="自动发券"
    :visible.sync="dialogVisible"
    class="auto-table-hidden-radio-label"
    width="1000px"
    :close-on-click-modal="false"
    @close="closeDialogFn"
  >
    <el-form label-width="80px">
      <el-form-item
        label="是否查券"
        required
      >
        <el-switch
          v-model="queryCoupon"
          active-color="#13ce66"
          inactive-color="#ff4949"
          active-text="是"
          inactive-text="否"
        >
        </el-switch>
      </el-form-item>
      <el-form-item
        label="客户渠道"
        required
      >
        <el-radio
          v-model="radio"
          label="1"
        >不限</el-radio>
        <el-radio
          v-model="radio"
          label="2"
        >场景值渠道</el-radio>
        <el-radio
          v-model="radio"
          label="3"
        >带参渠道</el-radio>
        <div
          id="scene"
          v-show="radio == 2"
        >
          <cy-select
            v-model="scene"
            multiple
            collapse-tags
            search-input
            placeholder="请选择场景值渠道"
            placeholder-search="搜索场景值名称"
            :options-list="sceneList"
            :options="{ label: 'name', value: 'val' }"
          >
          </cy-select>
        </div>
        <div
          id="channel"
          v-show="radio == 3"
        >
          <span class="marginR5">一级渠道</span>
          <cy-select
            v-if="radio == 3"
            v-model="channel1"
            search-input
            size="small"
            placeholder="请选择一级渠道"
            placeholder-search="搜索渠道名称"
            :options-list="channelList1"
            :options="{ label: 'name', value: 'id' }"
            @change="
              channelSelect($event);
            channel2 = [];
            "
          >
          </cy-select>
          <span
            v-show="channel1"
            class="marginL20 marginR5"
          >二级渠道</span>
          <cy-select
            v-show="channel1"
            v-model="channel2"
            multiple
            search-input
            collapse-tags
            widthSize="medium"
            clearable
            filterable
            placeholder="请选择二级渠道"
            placeholder-search="搜索渠道名称"
            :options-list="channelList2"
            :options="{ label: 'name', value: 'indice' }"
          >
          </cy-select>
        </div>
      </el-form-item>
    </el-form>
    <span>
      <cy-input
        size="small"
        width-size="small"
        :placeholder="`请搜索${typeSelect == 1 ? '批次号' : '名称'}`"
        v-model.trim="couponSearch"
        class="input-with-select"
        clearable
      >
        <el-select
          v-model="typeSelect"
          slot="prepend"
        >
          <el-option
            label="优惠券批次号"
            value="1"
          ></el-option>
          <el-option
            label="优惠券名称"
            value="2"
          ></el-option>
        </el-select>
      </cy-input>
      <el-date-picker
        v-model="couponDate"
        type="datetimerange"
        format="yyyy-MM-dd HH:mm:ss"
        value-format="timestamp"
        range-separator="至"
        start-placeholder="开始时间"
        end-placeholder="结束时间"
        :default-time="['00:00:00', '23:59:59']"
        class="date-picker"
      >
      </el-date-picker>
    </span>

    <el-button
      type="primary"
      size="small"
      class="m-b-20 marginL20"
      @click="onSearchFn"
    >查询</el-button>
    <el-button
      size="small"
      class="m-b-20"
      @click="resetFn"
    >重置</el-button>
    <el-table
      ref="table"
      :data="list"
      stripe
      @select="selectionChange"
    >
      <el-table-column
        type="selection"
        width="55"
        :selectable="checkSelectable"
      ></el-table-column>
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
          <div v-if="row.coupon_name.length > 5">
            <el-tooltip
              effect="light"
              :content="row.coupon_name"
              placement="top-start"
            >
              <span>{{ row.coupon_name.slice(0, 5) + "..." }}</span>
            </el-tooltip>
          </div>
          <div v-else>
            {{ row.coupon_name }}
          </div>
        </template>
      </el-table-column>
      <el-table-column
        prop="scope"
        label="优惠券类别"
        min-width="150"
        :show-overflow-tooltip="true"
      >
        <template slot-scope="{ row }">
          <div>
            {{ row.scope == "GLOBAL" ? "全场券" : "单品券" }}
          </div>
        </template>
      </el-table-column>
      <el-table-column
        prop="condition"
        label="优惠券条件"
        min-width="170"
      >
        <template slot-scope="{ row }">
          <div v-if="row.condition.length > 8">
            <el-tooltip
              effect="light"
              :content="row.condition"
              placement="top-start"
            >
              <span>{{ row.condition.slice(0, 8) + "..." }}</span>
            </el-tooltip>
          </div>
          <div v-else>
            {{ row.condition }}
          </div>
        </template>
      </el-table-column>
      <el-table-column
        prop="is_overlay"
        label="是否可叠加"
        min-width="150"
        :show-overflow-tooltip="true"
      >
        <template slot-scope="{ row }">
          {{ row.is_overlay == 1 ? "是" : "否" }}
        </template>
      </el-table-column>
      <el-table-column
        prop="template_name"
        label="有效期"
        min-width="400"
      >
        <template slot-scope="{ row }">
          {{ row.start_time | formatDate("yyyy-mm-dd hh:mm:ss") }} -
          {{ row.end_time | formatDate("yyyy-mm-dd hh:mm:ss") }}
        </template>
      </el-table-column>
      <div
        slot="empty"
        class="table-empty-label"
      >
        <img
          class="empty-img"
          :src="couponSearch
              ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
              : require('@/assets/images/empty-box.png')
            "
          alt=""
        />
        <p>{{ couponSearch ? "没有找到该优惠券" : "暂无优惠券" }}</p>
      </div>
    </el-table>
    <div
      class="clearfix"
      style="margin-top:20px;"
    >
      <el-pagination
        class="pull-right"
        @current-change="handleCurrentChangeFn"
        :current-page="currentPage"
        :page-size="5"
        layout="prev, pager, next, jumper"
        :total="total"
      >
      </el-pagination>
      <span
        class="pull-right"
        style="line-height:32px;"
      >已选：{{ this.checkCoupon.length }}/10</span>
    </div>

    <span
      slot="footer"
      class="dialog-footer"
    >
      <el-button
        @click="closeDialogFn"
        size="small"
      >取 消</el-button>
      <el-button
        type="primary"
        class="marginL20"
        @click="submitFn"
        size="small"
      >确 认
      </el-button>
    </span>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import DoubleSelect from "@/components/public/double-select/index.vue";
import API from "@/api/customer";
@Component({
  components: {
    DoubleSelect
  }
})
export default class extends Vue {
  @Prop() sceneList!: any;
  @Prop() chooseInfo!: any;
  @Prop() isused!: any;
  /************data**************/
  private templateRadio = "";
  private total = 0;
  private list = [];
  private dialogVisible = true;
  private currentPage = 1;
  private queryCoupon = true;
  private checkCoupon = [];
  private radio = "1";
  private channel1 = "";
  private channel2 = [];
  private channelList1: Array<any> = [];
  private channelList2: Array<any> = [];
  private scene = [];
  private couponSearch = "";
  private typeSelect = "1";
  private couponDate = [];
  created() {
    this.queryTableFn();
    this.queryChannelFn();
  }
  private mounted() {
    this.getFormInfo(this.chooseInfo);
  }
  private getFormInfo(item: object) {
    if (item["coupon_id"]) {
      this.checkCoupon = JSON.parse(JSON.stringify(item["coupon_list"]));
      this.queryCoupon = item["query_coupon"] == 1;
      this.radio = String(item["customer_channel"]);
      if (this.radio == "2") {
        this.scene = item["customer_channel_value"]["scene"];
      } else if (this.radio == "3") {
        for (const key in item["customer_channel_value"]) {
          API.channelListApi({
            enterprise_id: this.$store.state.user.USERINFO.enterprise_id || "",
            pid: 0
          }).then((res: any) => {
            if (res.code == 200) {
              res.data.list.forEach((v: any) => {
                if (v.name.length > 20) {
                  v.name = v.name.slice(0, 20) + "...";
                }
              });
              this.channelList1 = res.data.list;
              this.channel1 = this.channelList1.find(v => v.indice == key).id;
              this.channelSelect(this.channel1);
            }
          });
          this.$nextTick(() => {
            this.channel2 = item["customer_channel_value"][key];
          });
        }
      }
    }
  }
  checkSelectable(row: any) {
    return row.is_used == 0 || this.isused.indexOf(row.coupon_id) != -1;
  }
  channelSelect(row: any) {
    API.channelListApi({
      enterprise_id: this.$store.state.user.USERINFO.enterprise_id || "",
      pid: row
    }).then((res: any) => {
      if (res.code == 200) {
        res.data.list.forEach((v: any) => {
          if (v.name.length > 20) {
            v.name = v.name.slice(0, 20) + "...";
          }
        });
        this.channelList2 = res.data.list;
        console.log(this.channelList2, "chanaa");
      }
    });
  }
  queryChannelFn() {
    API.channelListApi({
      enterprise_id: this.$store.state.user.USERINFO.enterprise_id || "",
      pid: 0
    }).then((res: any) => {
      if (res.code == 200) {
        res.data.list.forEach((v: any) => {
          if (v.name.length > 20) {
            v.name = v.name.slice(0, 20) + "...";
          }
        });
        this.channelList1 = res.data.list;
      }
    });
  }
  queryTableFn() {
    const params = {
      enterprise_id: this.$store.state.user.USERINFO.enterprise_id || "",
      coupon_name: this.typeSelect == "2" ? this.couponSearch : "",
      stock_id: this.typeSelect == "1" ? this.couponSearch : "",
      start_time: this.couponDate[0] / 1000,
      end_time: this.couponDate[1] / 1000,
      page: this.currentPage,
      size: 5
    };
    API.couponListApi(params).then((res: any) => {
      if (res.code === 200) {
        const data = res.data.list;
        this.total = res.data.total || 0;
        this.list = data || [];
        this.$nextTick(() => {
          this.list.forEach((v, i) => {
            this.checkCoupon.forEach(t => {
              if (v["coupon_id"] == t["coupon_id"]) {
                (this.$refs["table"] as any).toggleRowSelection(
                  this.list[i],
                  true
                );
              }
            });
          });
        });
      } else {
        this.$message.error(res.message);
      }
    });
  }
  //选择优惠券
  selectionChange(e: any, s: any) {
    this.checkCoupon.forEach((v: any, i: number) => {
      if (v.coupon_id == s.coupon_id) {
        this.checkCoupon.splice(i, 1);
      }
    });
    if (this.checkCoupon.length > 9) {
      (this.$refs.table as any).toggleRowSelection(s, false);
      this.$message.warning("最多选择10张优惠券");
      return;
    }
    this.checkCoupon = this.checkCoupon.concat(e);
    const result: any = [];
    const obj: any = [];
    this.checkCoupon.forEach((v: any) => {
      if (obj.indexOf(v.coupon_id) == -1) {
        result.push(v);
        obj.push(v.coupon_id);
      }
    });
    console.log(result);
    this.checkCoupon = result;
  }
  async submitFn() {
    if (this.radio == "2" && this.scene.length == 0) {
      this.$message.warning("请选择场景值渠道");
      return;
    } else if (this.radio == "3" && this.channel2.length == 0) {
      this.$message.warning("请选择带参渠道");
      return;
    }
    if (this.checkCoupon.length == 0) {
      this.$message.warning("请选择优惠券");
      return;
    }
    const obj = {
      queryCoupon: this.queryCoupon,
      radio: this.radio,
      scene: this.scene,
      channel1: this.channel1,
      channel2: this.channel2,
      couponList: this.checkCoupon,
      coupon: [],
      sceneList: {},
      channelList: {}
    };
    let channel1 = "";
    this.channelList1.forEach(v => {
      if (v.id == this.channel1) {
        channel1 = v.indice;
      }
    });
    obj.couponList.forEach(v => {
      obj.coupon.push(v["coupon_id"]);
    });
    if (this.radio == "2") {
      obj.sceneList["scene"] = this.scene;
    } else if (this.radio == "3") {
      obj.channelList[channel1] = this.channel2;
    }
    this.$emit("addAutoCouponDialog", obj);
  }
  closeDialogFn() {
    this.$emit("addAutoCouponDialog");
  }
  onSearchFn() {
    this.currentPage = 1;
    this.queryTableFn();
  }
  resetFn() {
    this.couponSearch = "";
    this.currentPage = 1;
    this.queryTableFn();
  }
  handleCurrentChangeFn(val: number) {
    this.currentPage = val;
    this.queryTableFn();
    //获取当前页数据后根据selectList里当前page存的值进行渲染
  }
}
</script>
<style lang="scss" scoped>/deep/ .el-table__header-wrapper {
  .el-checkbox__inner {
    display: none;
  }
}

/deep/ .el-tag--mini:nth-child(1) {
  .el-select__tags-text {
    width: 87px;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: middle;
  }
}

/deep/ .el-input-group__prepend {
  width: 80px;
}

.date-picker {
  height: 36px;
  margin-left: 20px;
  vertical-align: middle;
}

.empty-img {
  margin-top: 50px;
  width: 124px;
}

.marginL20 {
  margin-left: 20px;
}

.marginR5 {
  margin-right: 5px;
}

.table-empty-label {
  display: flex;
  height: 340px;
  flex-direction: column;

  img {
    width: 124px;
    margin: 80px auto 0;
  }

  p {
    line-height: 20px;
    margin-top: 14px;
  }
}</style>
