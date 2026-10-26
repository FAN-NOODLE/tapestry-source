<template>
  <div>
    <layout-form type="normal">
      <template slot="content">
        <ul class="layout-content">
          <li v-for="item in strategyList" :key="item.id">
            <dl class="show">
              <dt>{{ item.strategy_name }}</dt>
              <dd>实时对“成交客户”且“未加企微好友”的客户发送加粉短信</dd>
              <div class="btn-show txt" v-if="item.is_finished == '1'">
                <span class="tip">当前策略未完善，请先去完善</span>
                <el-button type="primary" size="mini" @click="editFn(item.id)"
                  >完善</el-button
                >
              </div>
              <dd class="btn-show" v-if="item.is_finished == '2'">
                <div class="isSwitch">
                  <el-switch
                    v-model="item.status"
                    :active-text="item.status == 2 ? `已开启` : `已关闭`"
                    :active-value="2"
                    :inactive-value="1"
                    @change="isForbiddenFn(item.id, item.status)"
                  ></el-switch>
                </div>
                <div class="btn">
                  <el-button
                    type="primary"
                    size="mini"
                    @click="goDetail(item.id)"
                    >详情</el-button
                  >
                  <el-button
                    type="primary"
                    size="mini"
                    @click="detailFn(item.id)"
                    >查看数据</el-button
                  >
                </div>
              </dd>
            </dl>
          </li>
        </ul>
      </template>
    </layout-form>
  </div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import API from "@/api/customer";

@Component({
  components: {}
})
export default class extends Vue {
  /*************************** data start ***************************/
  private isForbidden = false; //是否禁用 1:禁用 2启用
  private strategyList: Array<any> = [];

  /*************************** method start ***************************/
  //关闭/开启禁用按钮
  private isForbiddenFn(id: any, status: any) {
    API.marketingUpdateStatusApi({ id, status }).then((res: any) => {
      const msg = status == 2 ? "启用" : "禁用";
      if (res.code === 200) {
        this.$message.success(`${msg}成功`);
      } else {
        this.$message.error(res.message);
      }
      this.getStrategyList();
    });
  }
  //完善
  private editFn(id: any) {
    this.$router.push({
      name: "myapp_customer_auto-market_insert-market-strategy_edit",
      query: { id }
    });
  }
  // 查看策略
  private goDetail(id: any) {
    this.$router.push({
      name: "myapp_customer_auto-market_insert-market-strategy_detail",
      query: { id }
    });
  }
  //查看数据
  private detailFn(id: any) {
    this.$router.push({
      name: "myapp_customer_auto-market_market-strategy_data-track",
      query: { id }
    });
  }
  // 获取列表接口
  private getStrategyList() {
    API.marketStrategyListApi({ strategy_type: 2 }).then((res: any) => {
      if (res.code === 200) {
        if (res.data.list) {
          this.strategyList = res.data.list;
        }
      } else {
        this.$message.error(res.message);
      }
    });
  }
  /*************************** 生命周期 start ***************************/
  private created() {
    this.getStrategyList();
  }
}
</script>
<style lang="scss" scoped>
.tip {
  color: $color-warn;
}
.txt {
  margin-top: 15px;
}
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
ul.layout-content {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  margin-right: -30px;
  margin-bottom: -20px;
  li {
    width: 380px;
    box-sizing: border-box;
    padding-right: 30px;
    margin-bottom: 20px;
  }
}
dl.show {
  background: $color-background;
  height: 100%;
  padding: 0 30px;
  border-radius: 5px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  dt {
    color: $color-text-main;
    font-weight: bold;
    font-size: 14px;
    line-height: 20px;
    margin: 30px 0 20px;
  }
  dd {
    font-size: 14px;
  }
  dd:nth-child(2) {
    color: $color-text-normal;
    line-height: 20px;
  }
  dd:nth-child(3) {
    color: $color-warn;
    margin: 15px 0 23px;
  }
  .btn-show {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
  }
  // align-items: center;
}
/deep/.el-button--mini {
  font-size: 14px;
}
</style>
