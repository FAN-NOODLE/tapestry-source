<template>
  <van-row gutter="10">
    <van-col span="24">
      <div class="copy-input">
        <textarea
          ref="textarea"
          class="textarea"
          v-model="keywords"
          type="textarea"
          contenteditable 
          autofocus
          placeholder="请粘贴或输入文本，点击“识别”按钮。智能拆分姓名、电话和地址（请输入带有完整省市区的地址信息）"
          @click="handleClick('textarea')"
        ></textarea>
        <div class="btn" @click="handlePasteIdentify">{{keywords?'识别':'识别'}}</div>
      </div>
    </van-col>
    <van-col span="12">
      <div class="input" :class="{error: isError(innerValue.name)}">
        <div class="label">姓名</div>
        <van-field
          class="value"
          v-model="innerValue.name"
          placeholder="请输入姓名"
          ref="name"
          @click="handleClick('name')"
        />
      </div>
    </van-col>
    <van-col span="12">
      <div class="input" :class="{error: isError(innerValue.phone)}">
        <div class="label">手机号</div>
        <van-field
          class="value"
          v-model="innerValue.phone"
          placeholder="请输入手机号"
          ref="phone"
          @click="handleClick('phone')"
        />
      </div>
    </van-col>
    <van-col span="24">
      <div class="input" :class="{error: isError(innerValue.countyCode)}">
        <div class="label">省市区</div>
        <com-area-popup v-model="innerValue.countyCode" @change="handleAreaChange">
          <van-field
            class="value"
            v-model="province"
            readonly
            :right-icon="require('@/assets/online-order/icon_right_2.png')"
            placeholder="请选择省市区"
          />
        </com-area-popup>
      </div>
    </van-col>
    <van-col span="24">
      <div class="input" :class="{error: isError(innerValue.address)}">
        <div class="label">详细地址</div>
        <van-field
          class="value"
          v-model="innerValue.address"
          clearable
          clear-trigger="always"
          placeholder="请输入，例如**街**号"
          ref="address"
          @click="handleClick('address')"
        />
      </div>
    </van-col>
    <van-col v-if="showClear" span="24">
      <div class="flex-end">
        <div class="btn" @click="handleClear">清空</div>
      </div>
    </van-col>
  </van-row>
</template>
<script lang="ts">
import smart from 'address-smart-parse'
import { Component, Model, Vue, Watch } from 'vue-property-decorator';
import ComAreaPopup from './com-area-popup.vue';
import wx from 'weixin-js-sdk'
export interface IAddressInfo {
  zipCode: string;
  province: string;
  provinceCode: string;
  city: string;
  cityCode: string;
  county: string;
  countyCode: string;
  street: string;
  streetCode: string;
  address: string;
  name: string;
  phone: string;
  idCard: string;
}

@Component({
  components: {
    ComAreaPopup
  }
})
export default class AddAddress extends Vue {
  @Model('change', { type: Object, default: () => ({}) }) value: IAddressInfo;

  // 粘贴板内容
  keywords = ''
  // 是否已识别
  isIdentify = false

  innerValue: IAddressInfo = {
    name: '',
    phone: '',
    county: '',
    countyCode: '',
    address: '',
    zipCode: '',
    province: '',
    provinceCode: '',
    city: '',
    cityCode: '',
    street: '',
    streetCode: '',
    idCard: ''
  }

  get showClear() {
    return this.keywords || Object.values(this.innerValue).some(v => !!v)
  }

  get province() {
    const { province, city, county } = this.innerValue
    return (province||'') + (city||'') + (county||'')
  }

  get isError() {
    return (value:string) => {
      if(!this.isIdentify) return false

      if(!value) return true

      return false
    }
  }
  
  @Watch('value')
  onValueChange(val: IAddressInfo) {
    this.innerValue = val
  }
  handleClick(key) {
    (this.$refs[key] as any).focus()
  }
  getClipboardContent() {
    return new Promise((resolve, reject) => {
      console.log(window.wx.getClipboardData);
      window.wx.invoke('getClipboardData', (res) => {
        console.log(res);
        
      })

      navigator.clipboard.readText()
        .then(text => {
            //my code to handle paste
            console.log(111, text);
            resolve(text)
        })
        .catch(err => {
            console.error('Failed to read clipboard contents: ', err);
            reject(err)
        });
    })
  }
  /**
   * 粘贴识别
   */
  async handlePasteIdentify() {
    // if(!this.keywords) {
    //   this.keywords = (await this.getClipboardContent()) as string;
    // }

    if(!this.keywords) {
      this.$toast('内容不能为空，请粘贴或手动输入');
      return;
    }
    // 数据格式化
    const text = this.keywords
      // .replaceAll(' ', '')
      .replace(/(收件人|联系人)：/, '')
      .replace(/(手机号|手机号码)：/, '')
      .replace(/(所在地区|地区)：/, '')
      .replace(/详细地址：/, '')
      .replace(/邮政编码：/, '')
      .replaceAll('\n', ' ')

    
    // 数据解析
    const res = smart(text);
    // 地址合并

    res.address = (res.street || '') + res.address
    Object.keys(res).map(key => {
      res[key] = res[key]||''
    })
    console.log(res);
    
    // emit
    this.$emit('change', res);
    this.isIdentify = true
  }

  /**
   * 省市区选择
   * @param value 
   * @param values 
   */
  handleAreaChange(value, values) {
    console.log(value, values);
    this.innerValue.province = values[0].name || ''
    this.innerValue.provinceCode = values[0].code || ''
    this.innerValue.city = values[1].name || ''
    this.innerValue.cityCode = values[1].code || ''
    this.innerValue.county = values[2].name || ''
    this.innerValue.countyCode = values[2].code || ''
    this.$emit('change', this.innerValue);
  }

  handleClear() {
    this.innerValue = {
      name: '',
      phone: '',
      county: '',
      countyCode: '',
      address: '',
      zipCode: '',
      province: '',
      provinceCode: '',
      city: '',
      cityCode: '',
      street: '',
      streetCode: '',
      idCard: ''
    }
    this.keywords = ''
    this.$emit('change', this.innerValue);
  
  }
}
</script>
<style lang="scss" scoped>
::v-deep .van-icon-clear {
  color: #000;
}
.copy-input {
  background-color: #F6F5F5;
  overflow: hidden;
  margin-bottom: 15px;

  .textarea {
    padding: 6px 10px;
    font-size: 12px;
    line-height: 17px;
    background-color: #F6F5F5;
    display: block;
    width: 100%;
    user-select: all;

    ::placeholder {
      color: #A6A6A6;
    }
  }

  .btn {
    font-weight: 400;
    font-size: 12px;
    color: #000000;
    line-height: 17px;
    text-decoration: underline;
    float: right;
    padding: 0 10px 10px;
    overflow: hidden;
  }

}

.input {
  margin-bottom: 15px;

  .label {
    font-weight: 400;
    font-size: 14px;
    color: #000000;
    line-height: 20px;
  }

  .value {
    padding: 15px 0;
    border-bottom: 1px solid #EEEEEE;
    font-size: 14px;

    ::placeholder {
      color: #A6A6A6;
    }

    
  }

  &.error {
    .value {
      color: #F64040;
      border-color: #F64040;

      ::placeholder {
        color: #F64040 !important;
      }
    }
  }
  
}

.flex-end {
  display: flex;
  justify-content: flex-end;

  .btn {
    font-weight: 400;
    font-size: 14px;
    color: #000000;
    line-height: 17px;
    text-decoration: underline;
  }
}
</style>