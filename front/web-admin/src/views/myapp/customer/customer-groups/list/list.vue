<template>
  <div class="layout-list">
    <p class="time-show">
      数据更新时间：{{ updateTime | formatDate("yyyy-mm-dd hh:mm:ss") }}
      <span @click="refreshFn">刷新</span>
    </p>
    <dl class="layout-content">
      <!-- <dt>
        <span>忠诚度模型AIPL（针对所有商城客户）</span>
        <el-popover
          placement="top-start"
          width="200"
          trigger="hover"
          :content="aiplStr"
        >
          <svg-icon
            slot="reference"
            class="help-icon"
            icon-class="group-help"
          ></svg-icon>
        </el-popover>
      </dt>
      <dd>
        <ul>
          <li v-for="item in aiplArr" :key="item.group_id">
            <div class="show">
              <h5>
                {{ item.group_name }}({{ item.group_mark }})
                <el-popover
                  placement="top-start"
                  trigger="hover"
                  :content="item.group_tips"
                >
                  <svg-icon
                    slot="reference"
                    class="help-icon-sub"
                    icon-class="group-help"
                  ></svg-icon>
                </el-popover>
              </h5>
              <div class="number" @click="goCustomerList(item)">
                {{ (item.value || 0) | formatNumber }}
              </div>
              <el-button type="primary" size="mini" @click="createBtnFn(item)"
                >创建营销策略</el-button
              >
            </div>
          </li>
        </ul>
      </dd>
      <dt>
        <span>RFM模型分层（仅针对成交客户分层）</span>
        <svg-icon
          class="help-second-icon"
          icon-class="group-help"
          @click="getRfmBoxFn"
        ></svg-icon>
      </dt>
      <dd>
        <ul>
          <li v-for="item in rfmArr" :key="item.group_id">
            <div class="show">
              <h5>
                {{ item.group_name }}
                <el-popover
                  placement="top-start"
                  trigger="hover"
                  width="200"
                  :content="item.group_tips"
                >
                  <svg-icon
                    slot="reference"
                    class="help-icon-sub"
                    icon-class="group-help"
                  ></svg-icon>
                </el-popover>
              </h5>
              <div class="number" @click="goCustomerList(item)">
                {{ (item.value || 0) | formatNumber }}
              </div>
              <el-button type="primary" size="mini" @click="createBtnFn(item)"
                >创建营销策略</el-button
              >
            </div>
          </li>
        </ul>
      </dd> -->
      <dt><span>自定义分组</span></dt>
      <dd>
        <ul>
          <li>
            <div class="show create" @click="setCreateFn">
              <i class="el-icon-plus plus-icon" aria-hidden="true"></i>
              <p>添加用户分组</p>
            </div>
          </li>
          <li v-for="item in customArr" :key="item.group_id">
            <div class="show">
              <!-- <span class="tips" v-if="item.new_agge_result == 2">计算中…</span> -->
              <h5>{{ item.group_name }}</h5>
              <div class="number" @click="goCustomerList(item)">
                {{ (item.value || 0) | formatNumber }}
              </div>
              <div>
                <el-button type="primary" size="mini" @click="createBtnFn(item)"
                  >创建营销策略</el-button
                >
                <el-button size="mini" @click="showBtnFn(item.group_id)"
                  >查 看</el-button
                >
              </div>
            </div>
          </li>
        </ul>
      </dd>
    </dl>
    <!--RFM模型分层-->
    <rfm-model :isShow="rfmTips" @isShow="closeRfmBoxFn"></rfm-model>
  </div>
</template>
<script lang="ts">
import { Vue, Component, Emit } from "vue-property-decorator";
import rfmModel from "../components/rfm-model.vue";
import API from "@/api/customer";
import CustomerModule from "@/store/modules/customer";

@Component({
  components: { rfmModel }
})
export default class extends Vue {
  /*************************** data start ***************************/
  // 忠诚度模型字段解释
  private aiplStr =
    "针对商城所有用户，从陌生到体验，到成交转化，最终变成一个忠实粉丝的过程，我们需要通过产品优化、品牌升级、活动运营等方式，让更多用户逐渐向上层运动；";
  private aiplTipsArr: any = [
    { name: "认知用户", tips: "仅访问落地页的客户", mark: "A" },
    { name: "兴趣用户", tips: "访问页面>=2但未成交客户", mark: "I" },
    { name: "购买用户", tips: "成交1次的客户", mark: "P" },
    { name: "忠诚用户", tips: "成交>=2次的客户", mark: "L" }
  ];
  // RFM模型分层字段解释
  private rfmTipsArr = [
    {
      name: "重要价值客户",
      tips: "客户最近消费时间近、消费频次和消费金额都很高。R=1，F=1，M=1"
    },
    {
      name: "重要发展客户",
      tips:
        "最近消费时间较远，但消费频次和金额都很高，说明这是个有一段时间没来的忠实客户，主动和他保持联系。R=0，F=1，M=1"
    },
    {
      name: "重要深耕客户",
      tips:
        "最近消费时间较近、消费金额高，但频次不高，忠诚度不高，很有潜力的用户，必须重点发展。R=1，F=0，M=1"
    },
    {
      name: "重要挽留客户",
      tips:
        "最近消费时间较远、消费频次不高，但消费金额高的用户，可能是将要流失或者已经要流失的用户。R=0，F=0，M=1"
    },
    { name: "一般价值客户", tips: "忠诚客户，但累计消费金额低。R=1，F=1，M=0" },
    {
      name: "一般保持客户",
      tips: "新客户，消费金额低，消费频次少。R=1，F=0，M=0"
    },
    {
      name: "一般发展客户",
      tips: "已流失的曾经消费频次高，但是消费金额低的客户。R=0，F=1，M=0"
    },
    {
      name: "一般挽留客户",
      tips: "已流失的曾经消费频次低，消费金额也低的客户。R=0，F=0，M=0"
    }
  ];
  private rfmTips = false; //RFM模型分层tips
  private aiplArr: any = []; //忠诚度分组列表
  private rfmArr: any = []; //rfm分组列表
  private customArr: any = []; //自定义分组列表
  private updateTime = 0; //数据更新时间
  @Emit("isShow")

  /*************************** method start ***************************/
  //获取标签列表
  private getListFn() {
    API.groupListApi({}).then(res => {
      const list = res.data;
      //忠诚度分组列表
      this.aiplArr = list.aipl || [];
      this.aiplArr.forEach((item: any) => {
        this.aiplTipsArr.forEach((subItem: any) => {
          if (subItem.name == item.group_name) {
            item.group_tips = subItem.tips;
            item.group_mark = subItem.mark;
          }
        });
      });

      //rfm分组列表
      this.rfmArr = list.rfm || [];
      this.rfmArr.forEach((item: any) => {
        this.rfmTipsArr.forEach((subItem: any) => {
          if (subItem.name == item.group_name) {
            item.group_tips = subItem.tips;
          }
        });
      });

      //自定义分组列表
      this.customArr = list.custom;
    });
  }
  private updata() {
    const date = new Date().getTime();
    //返回格式化后的日期
    this.updateTime = date / 1000;
  }
  //刷新
  private refreshFn() {
    this.getListFn();
    this.updata();
  }
  //添加用户分组
  private setCreateFn() {
    CustomerModule.setGroupName("");
    CustomerModule.setGroupS([]);
    CustomerModule.setGroupU([]);
    this.$router.push({
      name: "myapp_customer_customer-groups_create"
    });
  }
  //显示RFM模型分层tips
  private getRfmBoxFn() {
    this.rfmTips = true;
  }
  //关闭RFM模型分层tips
  private closeRfmBoxFn(val: boolean) {
    this.rfmTips = val;
  }
  //创建营销策略
  private createBtnFn(item: any) {
    CustomerModule.setCustomerGroups({
      id: item.group_id,
      val: item.value ? item.value : 0
    });
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_set-stratage",
      query: { type: "1" }
    });
  }
  //查看
  private showBtnFn(id: any) {
    this.$router.push({
      name: "myapp_customer_customer-groups_show",
      query: { id: id }
    });
    // this.$router.push({
    //   name: "myapp_customer_customer-groups_create"
    // });
  }
  private goCustomerList(row: any) {
    console.log(row.group_id, "row");
    this.$router.replace({
      name: "myapp_customer_customer-list",
      query: {
        group_name: row.group_name,
        group_id: row.group_id
      }
    });
  }

  /*************************** 生命周期 start ***************************/
  private created() {
    //获取列表
    this.getListFn();
    this.updata();
  }
}
</script>
<style lang="scss" scoped>
li {
  list-style: none;
}
dl,
dd,
ul,
li {
  margin: 0;
  padding: 0;
}
.help-icon {
  cursor: pointer;
}
.help-second-icon {
  @extend .help-icon;
  margin-top: 1px;
}
.help-icon-sub {
  @extend .help-icon;
  margin-top: -3px;
}
.layout-list {
  box-shadow: 0px 0px 20px 1px rgba(0, 0, 0, 0.08);
  background-color: $color-white;
  padding: 30px;
  color: $color-text-main;
  .time-show {
    color: $color-text-normal;
    font-size: 14px;
    line-height: 20px;
    span {
      font-size: 14px;
      color: $color-text-button;
      line-height: 20px;
      margin-left: 10px;
      cursor: pointer;
    }
  }
  .layout-content {
    margin-right: -30px;
    margin-bottom: -20px;
    dt {
      font-weight: bold;
      font-size: 18px;
      display: flex;
      align-items: center;
      span {
        display: inline-block;
        height: 18px;
        line-height: 18px;
        margin-right: 5px;
      }
    }
    dt:first-child {
      padding: 20px 0;
    }
    dt:not(:first-child) {
      padding-top: 10px;
      padding-bottom: 20px;
    }
    dd {
      ul {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        li {
          width: 25%;
          height: 188px;
          box-sizing: border-box;
          padding-right: 30px;
          margin-bottom: 20px;
        }
      }
    }
    div.show {
      background: $color-background;
      height: 100%;
      border-radius: 5px;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      position: relative;
      ::v-deep.el-button--mini {
        font-size: 14px;
      }
      span.tips {
        position: absolute;
        top: 10px;
        right: 10px;
        z-index: 1;
        color: $color-main;
        font-size: 12px;
        line-height: 18px;
      }
      h5 {
        line-height: 20px;
        font-size: 14px;
        font-weight: bold;
      }
      div.number {
        font-size: 26px;
        line-height: 36px;
        padding: 20px 0;
      }
    }
    div.create {
      background: #fff !important;
      border: 1px solid $color-second-line;
      cursor: pointer;
      p {
        font-size: 12px;
        line-height: 18px;
        color: $color-text-second;
        padding: 5px 0;
      }
      .plus-icon {
        font-size: 18px;
        color: $color-normal;
      }
    }
    div.create:hover {
      border: 1px solid $color-main;
      p {
        color: $color-main;
      }
      .plus-icon {
        color: $color-main;
      }
    }
  }
}
</style>
