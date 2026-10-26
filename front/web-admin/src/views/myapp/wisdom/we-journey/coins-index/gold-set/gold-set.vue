<template>
  <layout-form>
    <template slot="content">
      <el-form
        :model="formData"
        :rules="rules"
        ref="ruleForm"
        label-width="100px"
        label-position="left"
        inline
      >
        <!-- 积分与金币比例 -->
        <div class="gold-title">积分与金币比例</div>
        <div class="gold-desp">备注：导购获取积分的同时获取金币，金币用于在兑礼中心兑换礼品。可以设置积分与金币的比例
          （兑换后金币会做四舍五入保留一位小数展示在管理后台的金币明细中，小程序前端金币总数四舍五入保留整数）</div>
        <el-form-item
          prop="ratio_points"
          label-width="0"
        >
          <el-input
            class="mini-input"
            v-model="formData.ratio_points"
          ></el-input>
          积分=
        </el-form-item>
        <el-form-item
          prop="ratio_coins"
          label-width="0"
        >
          <el-input
            class="mini-input"
            v-model="formData.ratio_coins"
          ></el-input>
          金币
        </el-form-item>
        <div class="right-edit">
          <el-button
            type="text"
            icon="el-icon-document"
            @click="showRecord"
          >修改记录</el-button>
        </div>
        <!-- 兑换中心配置 -->
        <div class="gold-title">兑换中心配置</div>
        <div class="gold-desp">备注：按钮开启，Wejourney展示兑礼中心入口；按钮关闭，Wejourney将不展示兑礼中心入口</div>
        <el-form-item
          prop="mall_is_open"
          label="是否开启："
        >
          <el-radio-group v-model="formData.mall_is_open">
            <el-radio label="1">开启</el-radio>
            <el-radio label="0">关闭</el-radio>
          </el-radio-group>
        </el-form-item>
        <!-- 金币过期配置 -->
        <div class="gold-title">金币过期配置</div>
        <div class="gold-desp">备注：关闭过期设置，导购金币会持续累积；开启过期设置并设置过期时间，到达设定时间会将导购账户中的剩余金币清零</div>
        <el-form-item
          prop="expiration"
          label="过期设置："
        >
          <el-radio-group v-model="formData.coins_expire_is_open">
            <el-radio label="1">开启</el-radio>
            <el-radio label="0">关闭</el-radio>
          </el-radio-group>
        </el-form-item>
        <!-- 过期时间设置 -->
        <div>
          <el-form-item
            prop="coins_expire_time"
            label="过期时间："
            v-if="formData.coins_expire_is_open === '1'"
          >
            <el-date-picker
              v-model="formData.coins_expire_time"
              type="datetime"
              placeholder="选择日期"
              value-format="yyyy-MM-dd HH:mm:ss"
              style="width: 100%;"
              :picker-options="{
                disabledDate: (time) => {
                  return time.getTime() < Date.now() - 8.64e7
                },
              }"
            ></el-date-picker>
          </el-form-item>
        </div>
        <!-- 金币过期记录 -->
        <div class="gold-title">金币过期记录</div>
        <el-form-item label-width="0">
          <!-- table 序号，过期时间 -->
          <el-table
            :data="list"
            stripe
            border
            style="width: 381px;"
            max-height="400"
          >
            <el-table-column
              prop="id"
              label="记录ID"
              width="180"
            ></el-table-column>
            <el-table-column
              prop="expired_at"
              label="过期时间"
            >
              <template slot-scope="scope">
                {{ scope.row.expired_at | formatDate("yyyy-mm-dd hh:mm:ss") }}
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
      <gold-update-record
        v-if="goldUpdateRecordVisible"
        :goldUpdateRecordVisible="goldUpdateRecordVisible"
        :coins_settings_edit_record="coins_settings_edit_record"
        @closeGoldUpdateRecord="showRecord"
      ></gold-update-record>
    </template>

    <!-- 底部按钮 -->
    <template slot="footer">
      <el-button
        type="primary"
        @click="save"
      >保存</el-button>
      <el-button @click="goBack">返回</el-button>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator"
import goldUpdateRecord from "../components/gold-update-record/index.vue"
import { saveGoldConfig, getGoldConfig } from "@/api/we-journey/gold"
@Component({
  components: {
    goldUpdateRecord,
  },
})
export default class GoldSet extends Vue {
  [x: string]: any
  formData: any = {
    ratio_coins: "",
    ratio_points: "",
    mall_is_open: "0",
    coins_expire_is_open: "0",
    coins_expire_time: "",
  }
  validateExpirationTime(rule: any, value: any, callback: any) {
    if (!value) {
      callback(new Error("请选择过期时间"))
    } else {
      // 选择的时间不能小于当前时间
      const now = new Date().getTime()
      const selectTime = new Date(value).getTime()
      if (selectTime < now) {
        callback(new Error("过期时间需大于当前时间"))
      } else {
        callback()
      }
    }
  }
  validateRatioCoins(rule: any, value: any, callback: any) {
    if (!value) {
      callback(new Error("请输入金币数量"))
    } else {
      // 只能输入大于0的整数
      const reg = /^[1-9]\d*$/g
      if (!reg.test(value)) {
        callback(new Error("金币数量只能输入整数"))
      }
      // 不能超过4位数
      if (value.length > 4) {
        callback(new Error("金币数量不能超过4位数"))
      } else {
        callback()
      }
    }
  }
  validateRatioPoints(rule: any, value: any, callback: any) {
    if (!value) {
      callback(new Error("请输入积分数量"))
    } else {
      // 只能输入大于0的整数
      const reg = /^[1-9]\d*$/g
      if (!reg.test(value)) {
        callback(new Error("积分数量只能输入整数"))
      }
      // 不能超过4位数
      if (value.length > 4) {
        callback(new Error("积分数量不能超过4位数"))
      } else {
        callback()
      }
    }
  }
  rules: any = {
    ratio_coins: [
      { required: true, trigger: "blur",validator: this.validateRatioCoins },
    ],
    ratio_points: [
      { required: true, trigger: "blur",validator: this.validateRatioPoints },
    ],
    mall_is_open: [
      { required: true, message: "请选择是否开启", trigger: "blur" },
    ],
    coins_expire_time: [
      { required: true, trigger: "change", validator: this.validateExpirationTime },
    ],
  }
  goldUpdateRecordVisible = false
  created() {
    this.getGoldConfig()
  }
  /**
   * @description: 获取导购金币配置
   */
  getGoldConfig() {
    getGoldConfig({}).then((res: any) => {
      if (res.code === 200) {
        this.formData.ratio_coins = res.data.ratio_coins
        this.formData.ratio_points = res.data.ratio_points
        this.formData.mall_is_open = String(res.data.mall_is_open)
        this.formData.coins_expire_is_open = String(res.data.coins_expire_is_open)
        this.formData.coins_expire_time = res.data.coins_expire_time ? new Date(res.data.coins_expire_time * 1000) : ""
        this.list = res.data.coins_expire_record_list
        this.coins_settings_edit_record = res.data.coins_settings_edit_record
      }
    })
  }
  /**
   * @description: 显示金币修改记录弹窗
   */
  showRecord() {
    this.goldUpdateRecordVisible = !this.goldUpdateRecordVisible
  }
  /**
   * @description: 保存
   */
  save() {
    (this.$refs.ruleForm as any).validate((valid: any) => {
      if (valid) {
        //设置的过期时间要比当前时间晚5分钟
        if (this.formData.coins_expire_is_open === "1") {
          const now = new Date().getTime()
          const selectTime = new Date(this.formData.coins_expire_time).getTime()
          if (selectTime < now + 300000) {
            this.$message.error("过期时间需大于当前时间5分钟")
            return false
          }
        }
        const params = {
          ratio_coins: Number(this.formData.ratio_coins),
          ratio_points: Number(this.formData.ratio_points),
          mall_is_open: Number(this.formData.mall_is_open),
          coins_expire_is_open: Number(this.formData.coins_expire_is_open),
          coins_expire_time: new Date(this.formData.coins_expire_time).getTime() / 1000
        }
        if (params.coins_expire_is_open === 0) {
          delete params.coins_expire_time
        }
        saveGoldConfig(params).then((res: any) => {
          if (res.code === 200) {
            this.$message.success("保存成功")
            this.$router.go(-1)
          } else {
            this.$message.error(res.message)
          }
        })
      } else {
        return false
      }
    })
  }
  /**
   * @description: 返回
   * @return {*}
   */
  goBack() {
    this.$router.go(-1)
  }
}
</script>
<style lang="scss" scoped> .gold-title {
   font-size: 16px;
   font-weight: bold;
   position: relative;
   margin-left: 10px;
   margin-bottom: 15px;

   &::after {
     position: absolute;
     content: "";
     left: -15px;
     top: 2px;
     height: 15px;
     width: 5px;
     background: #000000;
   }
 }

 .gold-desp {
   font-size: 12px;
   color: #999999;
   margin-bottom: 15px;
 }

 .mini-input {
   width: 100px;
 }

 .right-edit {
   position: relative;
   top: 10px;
   left: 200px;
   display: inline-block;
 }
</style>
