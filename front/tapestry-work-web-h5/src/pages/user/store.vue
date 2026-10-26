<template>
  <div class="container">
    <div class="header">
      <div class="tips">
        <div v-if="$route.query.other=='1'" class="back" @click="$router.back()"></div>
        <span>请从以下门店选择您的专属服务顾问</span>
      </div>
      <div class="flex">
        <div class="search">
          <div class="icon-search"></div>
          <input class="value" v-model="searchForm.storeName" type="text" placeholder="搜索附近的门店" @change="handleSearch">
        </div>
        <select-city :value="searchForm.city" :options="cityList" @change="handleCityChange"></select-city>
      </div>

      <div v-if="!searchForm.lat && !loading" class="address" @click="handleGetLocation">
        <div class="icon-address"></div>
        <div class="">推荐离我最近的门店</div>
      </div>
    </div>
    <div class="content">
      <div v-for="item,index in showStoreList" class="item" :key="index">
        <div class="item-flex">
          <div class="item-name">{{ item.store_name }}</div>
          <div class="itme-btn" @click="hanldeAdd(item)">添加顾问</div>
        </div>
        <div class="item-flex">
          <div class="flex-1">
            <div class="item-address" @click="copy(item.store_address)">{{ item.store_address }}</div>
          </div>
          <div class="item-distance">{{ distance(item.distance) }}</div>
        </div>
      </div>
      <div class="more">已经到底了哦～</div>
    </div>
  </div>
</template>
<script lang="ts"> 
import { Component, Vue, Watch } from "vue-property-decorator";
import wxUtil from '@/utils/wx';
import { getAllStoreListApi, getQrcodeByUnionidApi } from "@/api/user";
import SelectCity from "./components/select-city.vue";

@Component({
  components: {
    SelectCity
  }
})
export default class StoreList extends Vue {
  storeList:any[] = []
  cityList:any[] = []
  searchForm = {
    city: '',
    storeName: '',
    unionid: '',
    state: '',
    lat: 0,
    lng: 0
  }
  loading: boolean = true

  get showStoreList() {
    const res = this.storeList.filter(item => {
      if(this.searchForm.city && this.searchForm.storeName) {
        return item.city === this.searchForm.city && item.store_name.includes(this.searchForm.storeName)
      } else if(this.searchForm.city && !this.searchForm.storeName) {
        return item.city === this.searchForm.city
      } else if(!this.searchForm.city && this.searchForm.storeName) {
        return item.store_name.includes(this.searchForm.storeName)
      } else {
        return true
      }
    })
    return res
  }

  get distance() {
    return (value) => {
      // 门店距离小于1km的展示单位m，>50km的不展示具体距离，展示>50km
      if(value === 0) {
        return ''
      } else if(value < 1000) {
        return `${value}m`
      } else if(value < 50000) {
        return `${(value / 1000).toFixed(1)}km`
      } else {
        return '>50km'
      }
    }
  }
  async mounted() {
    console.log(this.$router);
    
    // const host = window.location.origin
    // const hash = window.location.pathname.split("/")[2];
    // const IDENTIFIER = window.location.pathname.split("/")[1] || "";
    // const url = `${host}/${IDENTIFIER}/${hash}${this.$route.fullPath}`
    // await wxUtil.wxConfig()
    
    this.searchForm.unionid = (this.$route.query.unionid) as string
    this.searchForm.state = (this.$route.query.scene) as string
    this.searchForm.lng = Number(this.$route.query.lng || 0)
    this.searchForm.lat = Number(this.$route.query.lat || 0)
    if(localStorage.getItem('lbs') !== '2') {
      await this.handleGetLocation()
    }
    await this.getStoreList()
    this.loading = false
  }


  async getData() {
    const res = await getQrcodeByUnionidApi(this.searchForm)
    console.log(res);
    this.storeList = res.data.store_list
  }

  async getStoreList() {
    const res = await getAllStoreListApi(this.searchForm)
    this.storeList = res.data.store_list
    
    this.cityList = Object.entries(res.data.store_list.reduce((prev, cur) => {
      if(!prev[cur.city_pin_yin]) {
        prev[cur.city_pin_yin] = []
      }

      prev[cur.city_pin_yin].push(cur.city)
      // 去重
      prev[cur.city_pin_yin] = Array.from(new Set(prev[cur.city_pin_yin]))
      return prev
    }, {})).map(([key, value]) => {
      return {
        letter: key,
        list: value
      }
    })
    
    this.cityList.sort((a, b) => {
      return a.letter.charCodeAt(0) - b.letter.charCodeAt(0)
    })
    
    
    this.$track.trackStoreListView({
      previous_action: this.$route.query.other ? '选择其他服务顾问' : '未添加过投放范围内的导购，且没有绑定对应的管理柜台',
      is_lbs: !!this.searchForm.lat,
      store_quantity: this.storeList.length,
      filter_option: this.searchForm.city,
      search_key_word: '',
    })
    
  }
  
  async handleGetLocation() {
    try {
      this.$toast.loading({
        mask: true,
        duration: 0,
        message: '正在获取位置信息...'
      })
      // const host = window.location.origin
      // const hash = window.location.pathname.split("/")[2];
      // const IDENTIFIER = window.location.pathname.split("/")[1] || "";
      // const url = `${host}/${IDENTIFIER}/${hash}${this.$route.fullPath}`
    
      const res:any = await wxUtil.getLocation()
      this.searchForm.lat = Number(res.latitude)
      this.searchForm.lng = Number(res.longitude)
      this.getStoreList()
      this.$track.trackLBSAuth({
        button_type: '授权',
      })
      localStorage.setItem('lbs', '1')
      

    } catch (error) {
      this.$track.trackLBSAuth({
        button_type: '拒绝',
      })
      localStorage.setItem('lbs', '2')
    } finally {
      this.$toast.clear()
    }
  }

  handleCityChange(e) {
    this.searchForm.city = e


    this.$track.trackStoreListView({
      previous_action: this.$route.query.other ? '选择其他服务顾问' : '未添加过投放范围内的导购，且没有绑定对应的管理柜台',
      is_lbs: !!this.searchForm.lat,
      store_quantity: this.showStoreList.length,
      filter_option: this.searchForm.city,
      search_key_word: this.searchForm.storeName,
    })
  }

  handleSearch() {


    this.$track.trackStoreListView({
      previous_action: this.$route.query.other ? '选择其他服务顾问' : '未添加过投放范围内的导购，且没有绑定对应的管理柜台',
      is_lbs: !!this.searchForm.lat,
      store_quantity: this.showStoreList.length,
      filter_option: this.searchForm.city,
      search_key_word: this.searchForm.storeName,
    })
  }

  hanldeAdd(item) {
    this.$track.trackAddSAButtonClick({
      previous_action: this.$route.query.other ? '选择其他服务顾问' : '未添加过投放范围内的导购，且没有绑定对应的管理柜台',
      is_lbs: !!this.searchForm.lat,
      store_quantity: this.showStoreList.length,
      filter_option: this.searchForm.city,
      search_key_word: this.searchForm.storeName,
      store_code: item.store_code,
      store_name: item.store_name,
    })
    this.$router.push({
      path: '/user/qrcode',
      query: {
        unionid: this.searchForm.unionid,
        scene: this.searchForm.state,
        third_id: item.third_id
      }
    })
    // const url = window.location.href.replace('/user/store', '/user/qrcode').split('?')[0]
    // window.location.href = url + `?unionid=${this.searchForm.unionid}&scene=${this.searchForm.state}&third_id=${item.third_id}&t=${new Date().getTime()}`
      

  }


  copy(msg) {
    this.$helper.copyMessage(msg);
  }
}
</script>

<style lang="scss" scoped>


.icon-address {
  width: 14px;
  height: 14px;
  background: url('~@/assets/user/icon_address.png') center/contain no-repeat;

}
.icon-down {
  width: 10px;
  height: 10px;
  background: url('~@/assets/user/icon_down.png') center/contain no-repeat;
}
.icon-search {
  width: 14px;
  height: 14px;
  background: url('~@/assets/user/icon_search.png') center/contain no-repeat;
}
.container {
  background-color: #F7F7F7;
  min-height: 100vh;

  .header {
    background-color: #fff;
    padding: 12px 20px 20px 16px;
    position: sticky;
    top: 0;
    z-index: 1;
    .flex {
      margin-top: 16px;
      padding: 0 12px 0 10px;
      display: flex;
      align-items: flex-end;
      .select-city {
        width: 126px;
        height: 37px;
        border: 1px solid #C5C5C5;
        font-size: 14px;
        display: flex;
        align-items: center;
        padding: 0 10px;
        box-sizing: border-box;

        .icon-search {
          margin-right: 12px;
        }

        .value {
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          color: #000;
          height: 100%;
        }


        .placeholder {
          color: #949494;
          flex: 1;
        }
      }
      .search {
        flex: 1;
        margin-right: 16px;
        height: 38px;
        // border-bottom: 1px solid #C5C5C5;
        font-size: 14px;
        display: flex;
        align-items: center;
        padding: 0;
        box-sizing: border-box;
        position: relative;

        &::after { 
          content: "";
          width: 100%;
          height: 1PX;
          background: #C5C5C5;
          position: absolute;
          left: 0;
          bottom: 0;
          opacity: 0.7;
        } 

        .value {
          flex: 1;
          overflow: hidden;
          color: #000;
          height: 100%;
          &::placeholder {
            color: #949494;
          }
        }


        .icon-search {
          margin-right: 12px;
        }

      }
    }

    .tips {
      font-size: 14px;
      color: #000000;
      line-height: 20px;
      padding: 0 0 16px 0;
      text-align: center;
      // border-bottom: 0.5px solid #C5C5C5;
      position: relative;
      margin: 0 10px;
      &::after { 
        content: "";
        width: 100%;
        height: 1PX;
        background: #C5C5C5;
        position: absolute;
        left: 0;
        bottom: 0;
        opacity: 0.7;
      } 


      .back {
        position: absolute;
        left: 0;
        top: 0;
        width: 20px;
        height: 20px;
        background: url('~@/assets/user/icon_left.png') center/14px 14px no-repeat;
      }
    }
    
    .address {
      margin-top: 12px;
      height: 22px;
      line-height: 14px;
      display: flex;
      align-items: center;
      font-size: 14px;
      text-decoration: underline;
      color: #000000;
      padding: 4px 10px 0;

      .icon-address {
        margin-right: 12px;
      }
    }
  }

  .content {
    padding: 12px 20px;

    .item {
      margin-bottom: 12px;
      background: #fff url(~@/assets/user/bg_1.png) center -1px/100% no-repeat;
      box-shadow: 1px 1px 6px 0px rgba(0, 0, 0, 0.102);
      padding: 16px 12px;
      .item-flex {
        display: flex;
        .item-name {
          height: 30px;
          font-size: 16px;
          color: #000000;
          line-height: 22px;
          flex: 1;
          flex-wrap: wrap;
          padding: 4px 12px 4px 0;
          font-weight: bold;
        }
        .itme-btn {
          background-color: #000000;
          color: #FFFFFF;
          font-size: 12px;
          height: 30px;
          display: flex;
          width: 90px;
          align-items: center;
          justify-content: center;
          // margin-right: 2px;
        }
      }

      .item-flex:nth-child(2) {
        margin-top: 34px;
      }

      .flex-1 {
        flex: 1;
        overflow: hidden;
        display: flex;
      }

      .item-address {
        font-size: 12px;
        line-height: 17px;
        padding-right: 20px;
        position: relative;

        &::after {
          content: '';
          position: absolute;
          top: 0;
          right: 0;
          width: 14px;
          height: 17px;
          background: url('~@/assets/user/icon_copy.png') center/contain no-repeat
        }
      }
      .item-distance {
        width: 63px;
        text-align: right;
        font-size: 12px;
        line-height: 17px;
      }
    }
  }
}

.more {
  text-align: center;
  font-size: 14px;
  color: #C5C5C5;
}
</style>