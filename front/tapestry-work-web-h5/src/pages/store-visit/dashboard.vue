<template>
  <div class="container">
    <div class="header">
      <div class="row">
        <select-district-popup v-model="defaultArea" :list="list" @confirm="handleDistrictConfirm">
          <template #default="{text}">
            <div class="value">{{getName(defaultArea)}}</div>
            <div class="icon-down"></div>
          </template>
        </select-district-popup>
        <div class="switch">
          <div class="label">{{ $t('storeVisit.tips.29') }}</div>
          <van-switch v-model="searchForm.data_type" :active-value="1" :inactive-value="2" @change="handleDataTypeChange"></van-switch>
        </div>
      </div>
      <div class="tips">{{ $t('storeVisit.tips.28') }}：{{$helper.dateFormat(Date.now()/1000, "yyyy.MM.dd")}}</div>
      <div class="date-wrap">
        <select-daterange-popup :defaultValue="defaultValue" :defaultCurrent="defaultCurrent" :columns="columns" @confirm="handleDaterangeConfirm">
          <div class="date">
            <div class="label">{{ $t('storeVisit.tips.30') }}:</div>
            <div class="value">{{defaultTime}}</div>
            <div class="icon-caret-down"></div>
          </div>
        </select-daterange-popup>
        <div class="btn" @click="handleExport">{{ $t('storeVisit.btn.20') }}</div>
      </div>
    </div>
    <!-- 地图图例 -->
    <div class="legend">
      <div class="list">
        <div class="item" v-for="(item, index) in colors" :key="index">
          <div class="item-color" :style="{'background-color':item.value}"></div>
          <div class="item-text">{{ item.label }}</div>
        </div>
      </div>
    </div>
    <div class="content">
      <!-- 地图 -->
      <div ref="map" style="width: 100%; height: 100%"></div>
      <!-- 数据面板 -->
      <div v-if="dataList.length" class="dashboard">
        <div v-if="level == 1" class="list">
          <div class="item" v-for="item,index in dataList" :key="index">
            <div class="item-name">{{ item.name }}</div>
            <div class="flex">
              <div class="item-label">{{ $t('storeVisit.tips.26') }}</div>
              <div class="item-value">{{ item.total_count }}</div>
            </div>
          </div>
        </div>
        <div v-else class="data" @click="handleSelectStore">
          <div class="data-head">
            <div class="data-name">{{ getName(currentData.name) }}</div>
            <div class="data-more">
              <div class="text">{{ $t('storeVisit.tips.31') }}</div>
              <div class="icon-rigth"></div>
            </div>
          </div>
          <div class="data-content">
            <div class="data-label">{{ $t('storeVisit.tips.26') }}</div>
            <div class="data-value">{{ currentData.total_count }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="footer">
      <div class="tabbar">
        <div class="tabbar-item" @click="handleNav">{{ $t('storeVisit.tabs.1') }}</div>
        <div class="tabbar-item active">{{ $t('storeVisit.tabs.2') }}</div>
      </div>
    </div>
    <!-- 导出弹框 -->
    <com-export-popup v-model="showExportPopup" @confirm="handleExportConfirm"></com-export-popup>
    <select-store-list-popup v-model="showSelectStoreListPopup" :default-data="currentData" @confirm="handleSelectStoreConfirm"></select-store-list-popup>
  </div>
</template>

<script>
import { Component, Prop, Vue } from "vue-property-decorator";
import axios from 'axios';
import { exportDashboardListApi, getDashboardTreeApi } from "@/api/store-visit";
import { API } from "@/utils/env-config";
import SelectDistrictPopup from './components/select-district-popup.vue'
import SelectDaterangePopup from './components/select-daterange-popup.vue'
import SelectStoreListPopup from './components/select-store-list-popup.vue'
import ComExportPopup from './components/com-export-popup.vue'
import { findTreeNode } from '@/utils/tree';
import { getMonthListApi } from "@/api/online-order";
import dayjs from "dayjs";
import cityNameMap from './city'

function getContrastColor(hexColor) {
  // 将HEX颜色转换为RGB
  const hex = hexColor.replace("#", "");
  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);

  // 计算亮度
  const brightness = 0.299 * r + 0.587 * g + 0.114 * b;

  // 判断深浅
  return brightness < 128 ? "#ffffff" : "#000000";
}

function m(list, parent = null) {
  for (let i = 0; i < list.length; i++) {
    const item = list[i];
    if(parent) {
      const [a, b] = parent.key.split('_')
      if(b && item.name.indexOf(b)  !== -1) {
        item.key = `${a}_${item.name}`
      } else {

        item.key = `${parent.key}_${item.name}`
      }
    } else {
      item.key = item.name
    }
    // 遍历子元素
    if (item.children && item.children.length) {
      if(item.level === 1) {
        item.children = [{  name: 'China', children: [], dataList: item.dataList, level: 1, _children: item.children}, ...item.children]
      } else if(item.level === 2) {
        item.children = [{  name: `${item.name} Total`, children: [], total_count: item.total_count, level: 2, _children: item.children}, ...item.children]
      }
      item.children = m(item.children, item)
    }
  }

  return list
}

const TMap = window.TMap || qq.maps

@Component({
  components: {
    SelectDistrictPopup,
    SelectDaterangePopup,
    SelectStoreListPopup,
    ComExportPopup
  },
})
export default class StoreVisitDashboard extends Vue {
  /** 显示选择店铺列表弹窗 */
  showSelectStoreListPopup = false
  /** 地图实例 */
  map = null;
  /** 中国地图数据 */
  geoJsonDataChina = {
    type: 'FeatureCollection',
    features: []
  }
  /** 地图数据 - 省份 */
  geoJsonData = {
    type: 'FeatureCollection',
    features: []
  }
  /** 地图数据 - 城市 */
  geoJsonDataCity = {
    type: 'FeatureCollection',
    features: []
  }
  /** 门店巡店数据 tree */
  list = []
  /** 一级大区面板数据  */
  dataList = []
  /** 当前选中省份的数据 */
  currentData = {
    name: '',
    children: []
  }
  /** 穿透层级 */
  level = 1
  /** 图例 */
  get colors() {
    return [
      { label: `${this.$t('storeVisit.tips.26')} > 20`, value: '#00256B', min: 20, max: 10000},
      { label: `20 > ${this.$t('storeVisit.tips.26')} > 10`, value: '#0151F3', min: 10, max: 20},
      { label: `${this.$t('storeVisit.tips.26')} < 10`, value: '#00AAFF', min: 0, max: 10},
      { label: `${this.$t('storeVisit.tips.26')} = 0`, value: '#E9ECEF', min: 0, max: 0},
      { label: this.$t('storeVisit.tips.26'), value: '#CED4DA', min: 0, max: 0},
    ]
  }
  /** 默认选中Key */
  defaultArea = 'Retail_China'
  /** 默认搜索条件 */
  searchForm = {
    data_type: 1,
    start_time: 0,
    end_time: 0
  }
  /** 默认时间类型 */
  defaultCurrent = 0
  /**  */
  defaultValue = {
    0: new Date().getFullYear(),
    1: [new Date().getFullYear(), 1],
    2: [new Date().getFullYear(), new Date().getMonth() + 1],
    3: []
  }
  /** 时间选择列 */
  columns = []
  /** 地图覆盖物IDs */
  geometrieIds = []
  labelIds = []
  /** 多边形实例 */
  multiPolygon = null
  multiLabel = null
  multiCityLabel = null
  /** 导出弹框 */
  showExportPopup = false
  
  get defaultTime() {
    const { start_time, end_time } = this.searchForm
    return `${dayjs(start_time*1000).format('YYYY.MM.DD')} - ${dayjs(end_time*1000).format('YYYY.MM.DD')}`
  }
  /** 获取地图数据 */
  get mapData() {
    return (level) => {
      if (level === 1) {
        return this.geoJsonData.features;
      } else if(level === 2) {
        return this.geoJsonDataCity.features
      } else {
        return this.geoJsonDataCity.features;
      }
    }
  }
  /** 缩放级别 */
  get zoom() {
    return (name) => {
      if (this.level === 3) {
        return 8;
      } else if(this.level === 2) {
        const zoomMap = {
          '上海市': 8,
        }
        return zoomMap[name]|| 6.5;
      }
      return 2;
    }
  }

  /** 获取adcode */
  get adcode () {
    return (name) => {
      const list = this.mapData(1)
      const cur = list.find(item => item.properties.name === name)
      return cur.properties.adcode
    }
  }

  /** 获取地图中心点 */
  get center() {
    return (name) => {
      if(this.level === 1) {
        const {centroid, center} = this.geoJsonDataChina.features[0].properties
        const [lng, lat] = centroid || center
        return [lng, lat]
      } else if(this.level === 2) {
        const list = this.mapData(this.level-1)
        const cur = list.find(item => item.properties.name === name)
        // console.log('center 2', list, cur);
        const [lng, lat] = cur ? cur.properties.centroid||cur.properties.center : []
        
        return [ lng, lat ]
      } else if(this.level === 3) {
        const list = this.mapData(this.level-1)
        const cur = list.find(item => item.properties.name === name)
        // console.log('center', list, cur);
        const [lng, lat] = cur ? cur.properties.centroid||cur.properties.center : []
        
        return [ lng, lat ]

      }
    }
  }

  /** 数据映射 */
  get dataMap() {
    return (key) => this.currentList(key).reduce((prev, cur) => {
      //  + parseInt(Math.random() * 30) + 10
      prev[cur.name] = cur.total_count || 0
      return prev;
    }, {})
  }

  /** 颜色映射 */
  get color() {
    return (name, key) => {
      const [a, b, c] = (key||this.defaultArea).split('_')
      const value = this.dataMap(`${a}_${b}`)[name]
      // console.log(this.defaultArea, name, value);
      
      if(value>20) {
        return this.colors[0].value
      } else if(value>10) {
        return this.colors[1].value
      } else if(value>0) {
        return this.colors[2].value
      } else if(value===0) {
        return this.colors[3].value
      } else {
        return this.colors[4].value
      }
    }
  }

  /** 当前数据列表 */
  get currentList () {
    return (key) => {
      if(this.list.length  === 0) return
      const current = findTreeNode(JSON.parse(JSON.stringify(this.list)), (node) => {
        return node.key === key;
      })
      

      // if(this.level === 3) {
      //   return [current]
      // }
      
      return current._children || current.children
    }
  }

  get currentItem () { 
    return (key) => {
      // console.log(key);
      
      if(this.list.length  === 0) return
      const current = findTreeNode(JSON.parse(JSON.stringify(this.list)), (node) => {
        return node.key === key;
      })
      return current || {}
    }
  }

  get trackInfo() {
    const { dataList = [], total_count } = this.currentData
    const total = dataList.reduce((prev, cur) => {
      return prev + cur.total_count
    }, 0)
    
    return {
      view_dimension: this.defaultArea,
      period: this.defaultTime,
      only_initiated_oneself: this.searchForm.data_type === 1 ? '是' : '否',
      SV_cnt: total_count || total,
    }
  }

  get lang () {
    return this.$store.getters["lang/getLang"]
  }
  get getName() {
    return (name) => {
      if(this.lang === 'en-US') {
        // 通过下划线分组
        const arr = name.split('_')
        // console.log(arr);
        
        return arr.map(item => {
          if(item.indexOf(' Total') > -1) {
            return cityNameMap[item.replace(' Total', '')] + ' Total'
          }
          return cityNameMap[item] || item

        }).join('_')
      }
      
      return name
    }
  }

  mounted() {
    this.$nextTick(() => {
      this.init();
    });
  }

  async init() {
    try {
      this.initMap();
      await this.getMonthList()
      await Promise.all([
        this.getGeoJsonChina(),
        this.getGeoJson(),
        this.getGeoJsonCity(),
        this.getData()
      ]);

      // // console.log(this.list[0].children);
      this.addPolygon()
      this.multiLabel.setVisible(false);

      console.log(this.$refs.map);
      // 获取子元素div
      const mapDiv = this.$refs.map;
      if(mapDiv) {
        const child = mapDiv.children[1];
        console.log(child);
        child.style.display = 'none';
      }
      


      
      
    } catch (error) {
      console.error(error);
      
    }
  }

  async getMonthList() {
    const res = await getMonthListApi()
  
    const nowItem = res.data.find(item => item.now == 1)
    // 获取财年数据
    const years = Object.values(res.data.reduce((prev, cur) => {
      if(prev[cur.fiscal_year]) {
        if(cur.now) {
          prev[cur.fiscal_year].now = cur.now
        }
        prev[cur.fiscal_year].data.end_time = cur.end_time
      } else {
        prev[cur.fiscal_year] = {
          label: cur.fiscal_year,
          value: cur.fiscal_year,
          now: cur.now,
          data: {
            start_time: cur.start_time,
            end_time: cur.end_time
          }
        }
      }
      return prev
    }, {}))
    // 财季
    const quarters = Object.values(res.data.reduce((prev, cur) => {
      
      const i = Math.ceil(cur.period/3)
      const key = `${cur.fiscal_year}_${i}`
      
      if(prev[key]) {
        if(cur.now) {
          prev[key].now = cur.now
        }
        prev[key].data.end_time = cur.end_time
      } else {
        prev[key] = {
          label: `P${(i-1)*3+1} - P${(i-1)*3+3}`,
          value: i,
          now: cur.now,
          year: cur.fiscal_year,
          data: {
            start_time: cur.start_time,
            end_time: cur.end_time
          }
        }
      }
      return prev
    }, {}))
    // 财月
    const months = res.data.map(item => {
      // console.log(dayjs.unix(item.start_time).format('YYYY-MM-DD'));
      
      return {
        label: `P${item.period}`,
        value: item.period,
        year: item.fiscal_year,
        now: item.now,
        data: {
          start_time: item.start_time,
          end_time: item.end_time
        }
      }
    })

    this.columns = [
      {
        label: '财年视图',
        value: 0,
        children: years
      },
      {
        label: '财季视图',
        value: 1,
        children: quarters
      },
      {
        label: '财月视图',
        value: 2,
        children: months
      },
      {
        label: '自定义',
        value: 3,
      }
    ]

    this.defaultValue = {
      0: nowItem.fiscal_year,
      1: [nowItem.fiscal_year, Math.ceil(nowItem.period/3)],
      2: [nowItem.fiscal_year, nowItem.period],
      3: []
    }
    
    // 获取默认时间
    const { start_time, end_time } = years.find(item => item.now == 1).data
    this.searchForm.start_time = start_time
    this.searchForm.end_time = end_time
  }

  async getGeoJsonChina() {

    // const topo = await axios('/static/map/china.topo.json');
    // console.log(topo);
    
    // const res = await axios(location.origin +'/static/map/100000.json');
    const res = await axios(API.cosUrl+'/administrative_division/100000.json');
    // console.log(res.data);
    this.geoJsonDataChina = res.data;
  }
  async getGeoJson() {
    // const res = await axios('/static/map/100000_full.json');
    const res = await axios(API.cosUrl+'/administrative_division/100000_full.json');
    // console.log(res.data);
    this.geoJsonData = res.data;
  }
  async getGeoJsonCity() {
    const res = await axios(API.cosUrl+'/administrative_division/100000_full_city.json');
    // console.log(res.data);
    this.geoJsonDataCity = res.data;
  }

  /**
   * 获取数据
   */
  async getData() {
    const res = await getDashboardTreeApi(this.searchForm)
    // console.log(1111, m(JSON.parse(JSON.stringify(res.data))));
    this.list = m(JSON.parse(JSON.stringify(res.data)))
    const [a] = this.defaultArea.split('_')
    const index = this.list.findIndex(item => item.name === a)
    console.log(a, index);
    console.log(this.list[index].dataList);
    
    this.dataList  = this.list[index].dataList
    this.currentData = findTreeNode(JSON.parse(JSON.stringify(this.list)), (node) => {
      return node.key === this.defaultArea;
    })
    
  }

  /** 初始化地图 */
  initMap() {
    // const center = this.center();
    const center = [104.113164, 37.570667];
    console.log(center);
    
    // 初始化地图
    this.map = new TMap.Map(this.$refs.map, {
      center: new TMap.LatLng(center[1], center[0]), // 地图中心点
      zoom: this.zoom(), // 地图缩放级别
      // 缩放过度
      zoomEnable: true,
      pitch: 0, // 俯仰度
      rotation: 0, // 旋转角度，
      viewMode: '2D',
      showControl: false,
      offset: { // 中心点偏移
        x: 0,
        y: -100,
      },
      baseMap: [
        { type: 'vector', features: [''] },
        // {
        //   type: 'traffic',
        //   opacity: 0,
        // }
      ],
      renderOptions: {
        enableBloom: false, // 是否启用泛光效果 注：为true才会有效果
      },
      maxZoom: 12,
    });

    // 监听地图缩放
    this.map.on('zoom', (e) => {
      const z = this.map.getZoom()
      if(z >= 4.5) {
        this.multiLabel.setVisible(true)
      } else {
        this.multiLabel.setVisible(false)
      }

      if(z>=6) {
        this.multiCityLabel.setVisible(true)
      } else {
        this.multiCityLabel.setVisible(false)

      }
    });
  }
  /** 计算paths */
  coordToPaths(coords, type) {
    // console.log('coordToPaths', type, coords.length);
    if(type === 'Polygon') {
      return coords[0].map(coord => {
        return new TMap.LatLng(coord[1], coord[0])
      });
    } else if(type === 'MultiPolygon') {
      return coords.map(item => {
        // console.log(111, item.length);
        return item.map(v => {
          return v.map(coord => {
            return new TMap.LatLng(coord[1], coord[0])
          });
        });
      });
    }
  }
  /**
   * 高亮 中国地图
   */
  highlight() {
    const { type, coordinates} = this.geoJsonDataChina.features[0].geometry

    const paths = this.coordToPaths(coordinates, type)
    
    const a = this.map.enableAreaHighlight({
      paths: paths,
      highlightColor: 'rgba(0,0,0,0)', // 高亮颜色
      shadeColor: 'rgba(244, 246, 247, 1)', // 其余阴影部分
    });
  }
  /**
   * 高亮 中国地图
   */
  highlightProvince({ type, coordinates}) {
    // const { type, coordinates} = this.geoJsonDataChina.features[0].geometry
    // console.log(geometry);

    const paths = this.coordToPaths(coordinates, type)
    
    this.map.enableAreaHighlight({
      paths: paths,
      highlightColor: 'rgba(0,0,0,0)', // 高亮颜色
      shadeColor: 'rgba(244, 246, 247, 1)', // 其余阴影部分
    });
    // this.map.setCenter(new TMap.LatLng(this.center[1], this.center[0]));
  }
  /** 添加覆盖物 */
  addPolygon() {

    const geometries = []
    const styles = {}
    const lableStyles = {}
    const [a, b, c] = this.defaultArea.split('_').map(item => item.replace(' Total', ''))
    const labelList = []
    const cityLabelList = []

    if(!this.multiPolygon) {
      // 1. 创建 MultiPolygon
      this.multiPolygon = new TMap.MultiPolygon({
        map: this.map,
      });
      // 2. 添加事件
      this.multiPolygon.on('click', (e) => {
        
        const props = e.geometry.properties;
        
        if(props.name && props.level >= 2) {
          this.defaultArea = props.key
          if(props.level === 3) {
            this.level = props.level
            this.currentData = props
            const center = this.center(props.name.replace(' Total', ''))
            this.map.setCenter(new TMap.LatLng(center[1], center[0]))
            this.map.setZoom(this.zoom(b))
            // 神策数据上报 -- 巡店筛选
            this.$track.trackSVFilter({
              current_page: '看板',
              filter_type_sv: '查看维度',
              filter_option: this.defaultArea,
              filter_method: '点击地图',
            })
          } else {
            this.handleDistrictConfirm(props, '点击地图')
          }
        }
        return false
      });
    }

    if(!this.multiLabel) {
      this.multiLabel = new TMap.MultiLabel({
        id: 'label-layer',
        map: this.map,
      });
    }
    if(!this.multiCityLabel) {
      this.multiCityLabel = new TMap.MultiLabel({
        id: 'label-city',
        map: this.map,
      });
    }

    Object.assign(styles, this.multiPolygon.getStyles())

    // 3. 添加省份
    const provinces = this.mapData(1)
      .filter(feature => {
        return !!feature.properties.name
      })
      .map(feature => {
        const type = feature.geometry.type
        const coords = feature.geometry.coordinates
        const name = `${feature.properties.name}`
        const key = `${a}_${name} Total`
        
        // if(b === name ) {
        //   console.log(1111);
        //   this.highlightProvince(feature.geometry)
        // }
        const c = feature.properties.centroid || feature.properties.center

        labelList.push({
          id: name, // 点图形数据的标志信息
          styleId: name, // 样式id
          position: new TMap.LatLng(c[1], c[0]), // 标注点位置
          content: this.getName(name), // 标注文本
          properties: {
            // 标注点的属性数据
            title: name,
          },
        })
        
        
        const paths = this.coordToPaths(coords, type)
        const color = this.color(name, `${a}_China`)
        // 添加样式
        styles[name] = new TMap.PolygonStyle({
          color: color, //面填充色
          borderColor: '#ffffff', //边线颜色
          borderWidth: 1, //边线宽度
        })

        lableStyles[name] = new TMap.LabelStyle({
          size: b == name ? 16 : 14, // 文字大小属性
          color: b == name ? 'rgba(0,0,0,1)': getContrastColor(color), //面填充色
        })
        
        return {
          id: name,
          styleId: name,
          paths: paths,
          name: name,
          properties: this.currentItem(key)
        }
      })
    this.multiPolygon.updateGeometries(provinces)


    // 删除已有的城市
    this.multiPolygon.remove([...this.geometrieIds])
    this.multiCityLabel.remove([...this.geometrieIds])
    if((this.level === 2 || this.level === 3) && !['北京市', '天津市', '上海市', '重庆市'].includes(b)) {
      // 4. 添加省份
      const citys = this.mapData(this.level)
        .filter(feature => {
          if(this.level === 1) {
            return !!feature.properties.name
          } else if(this.level === 2) {
            return !!feature.properties.name && feature.properties.parent.adcode === this.adcode(b)
          } else if(this.level === 3) {
            return !!feature.properties.name && feature.properties.parent.adcode === this.adcode(b)
          }
        })
        .map(feature => {
          const type = feature.geometry.type
          const coords = feature.geometry.coordinates
          const name =  feature.properties.name
          const key = this.level == 1 ? `${a}_${name} Total` : `${a}_${b}_${name}`

          // console.log(this.currentItem(key));
          
          const c = feature.properties.centroid || feature.properties.center

          cityLabelList.push({
            id: name, // 点图形数据的标志信息
            styleId: name, // 样式id
            position: new TMap.LatLng(c[1], c[0]), // 标注点位置
            content: this.getName(name), // 标注文本
            properties: {
              // 标注点的属性数据
              title: name,
            },
          })
          
          const paths = this.coordToPaths(coords, type)
          const color = this.color(name)
          // this.addMultiPoly(paths, ;
          styles[name] = new TMap.PolygonStyle({
            color: color, //面填充色
            borderColor: '#ffffff', //边线颜色
            borderWidth: 1, //边线宽度
          })

          lableStyles[name] = new TMap.LabelStyle({
            size: 12, // 文字大小属性
            color: getContrastColor(color), //面填充色
          })
  
          return {
            id: name,
            styleId: name,
            paths: paths,
            name: name,
            properties: this.currentItem(key)
          }
        })
      // 更新地图
      this.multiPolygon.updateGeometries(citys)
      // 移除当前城市所在的省份
      this.multiPolygon.remove([b])
      // 记录城市覆盖物ID
      this.geometrieIds = citys.map((geometry) => geometry.id)
    }
    // 5. 设置样式
    this.multiPolygon.setStyles(styles)
    // 6. 添加label
    
    this.multiLabel.updateGeometries(labelList)
    this.multiCityLabel.updateGeometries(cityLabelList)
    // 6. 添加label样式
    this.multiLabel.setStyles(lableStyles)
    this.multiCityLabel.setStyles(lableStyles)
  }

  /** 切换地区 */
  handleDistrictConfirm(value, filter_method = '切换维度') {
    // console.log('value', value);
    // console.log('level', value.level);
    // console.log('name', value.name);
    // console.log('key', value.key);
    
    this.level = value.level
    this.currentData = value
    
    // this.multiPolygon.remove(this.geometrieIds);
    const name = value.name.replace(' Total', '')
    const center = this.center(name)
    this.map.setCenter(new TMap.LatLng(center[1], center[0]))
    this.map.setZoom(this.zoom(name))
    
    this.addPolygon()

    if(value.level === 1) {
      this.dataList = value.dataList
    }

    // 神策数据上报 -- 巡店筛选
    this.$track.trackSVFilter({
      current_page: '看板',
      filter_type_sv: '查看维度',
      filter_option: this.defaultArea,
      filter_method: filter_method,
    })
  }

  /** 重置地图 */
  resetMap() {
    // this.defaultArea = 'Retail_China'
    // this.level = 1
    
    // const center = this.center()
    // this.map.setCenter(new TMap.LatLng(center[1], center[0]))
    // this.map.setZoom(this.zoom)
    // 删除历史覆盖
    console.log(1111, this.geometrieIds);
    
    // this.multiPolygon.remove(this.geometrieIds);
    // 添加新的覆盖物
    this.addPolygon()
    // this.dataList = this.list[0].dataList
  }

  /** 时间选择器 确定 */
  async handleDaterangeConfirm(e) {
    console.log(e);
    this.defaultCurrent = e.current
    this.defaultValue[e.current] = e.value
    this.searchForm.start_time = e.data.start_time
    this.searchForm.end_time = e.data.end_time
    await this.getData()
    this.resetMap()


    // 神策数据上报 -- 巡店筛选
    this.$track.trackSVFilter({
      current_page: '看板',
      filter_type_sv: '统计范围',
      filter_option: this.defaultTime,
      filter_method: '切换维度',
    })
  }

  /**  */
  async handleDataTypeChange() {
    await this.getData()
    this.resetMap()

    // 神策数据上报 -- 巡店筛选
    this.$track.trackSVFilter({
      current_page: '看板',
      filter_type_sv: '是否我发起的巡店',
      filter_option: this.searchForm.data_type == 1 ? '是' : '否',
      filter_method: '切换维度',
    })
  }

  /** 跳转 */
  handleNav() {
    // 神策数据上报 -- tab切换
    this.$track.trackTabSwitch({
      tab_name: '巡店'
    })
    this.$router.replace({
      path: '/store-visit',
    })
  }

  handleSelectStore() {
    this.showSelectStoreListPopup = true
    // 神策数据上报 -- 看板点击
    this.$track.trackDashboardClick({
      ...this.trackInfo,
      button_type: '查看区域巡店情况',
      store_name: '',
      store_code: ''
    })
  }
  handleSelectStoreConfirm(value) {
    // this.showSelectStoreListPopup = false;
    this.$store.commit("storeVisit/setUpdate", true)
    // 神策数据上报 -- 看板点击
    this.$track.trackDashboardClick({
      ...this.trackInfo,
      button_type: '查看门店巡店情况',
      store_name: value.name,
      store_code: value.store_code
    })
    this.$nextTick(() => {
      this.$router.push({
        path: '/store-visit-record',
        query: {
          third_id: value.third_id,
          data_type: String(this.searchForm.data_type),
          start_time: this.searchForm.start_time,
          end_time: this.searchForm.end_time,
          store_name: value.name
        }
      })
    })
    
  }

  /**
   * 导出
   */
  handleExport() {
    this.showExportPopup = true;
    // 神策数据上报 -- 看板点击
    this.$track.trackDashboardClick({
      ...this.trackInfo,
      button_type: '导出',
      store_name: '',
      store_code: ''
    })
}

  /**
   * 导出确认
   * @param value 
   */
  async handleExportConfirm(value) {
    this.showExportPopup = false;
    this.$toast.loading(this.$t('common.toast.6'))
    const res = await exportDashboardListApi({
      start_time: this.searchForm.start_time,
      end_time: this.searchForm.end_time,
      data_type: this.searchForm.data_type,
      search_level: this.level,
      search_name: this.currentData.name.replace(' Total', ''),
      email: value
    })
    this.$toast.clear()
    if(res.code === 200) {
      this.$toast(this.$t('common.toast.7'))

      // 神策数据上报 -- 巡店记录导出
      this.$track.trackSVExport({
        ...this.trackInfo,
      })
    } else {
      this.$toast(this.$t('common.toast.8'))
    }
  }
  
}
</script>

<style lang="scss" scoped>
.icon-rigth {
  width: 16px;
  height: 16px;
  margin-left: 4px;
  background: url('~@/assets/store-visit/icon_right_black.png') center/100% no-repeat
}

.van-switch {
  width: 26px;
  height: 16px;
  border: none;
  background-color: rgba(173, 181, 189, 1);
  overflow: hidden;
}

::v-deep .van-switch__node {
  width: 12px;
  height: 12px;
  top: 50%;
  transform: translate(2px, -50%);
  box-shadow: none;
}

::v-deep .van-switch--on {
  background-color: rgba(13, 110, 253, 1);
  .van-switch__node {
    transform: translate(12px, -50%);
  }
}
.container {
  height: 100vw;
  height: 100vh;
  background: linear-gradient(180deg, rgba(73, 80, 87, 0.1) 0%, rgba(159, 174, 189, 0.1) 32.47%, rgba(192, 203, 214, 0.1) 100%)，linear-gradient(0deg, #FFFFFF, #FFFFFF);
  display: flex;
  flex-direction: column;

  .header {
    // height: 116px;
    background-color: rgba(255, 255, 255, 1);
    padding: 12px 16px;

    .row  {
      height: 32px;
      display: flex;
      justify-content: space-between;
      ::v-deep .select {
        width: auto;
      }

      .value {
        font-size: 12px;
        font-weight: 700;
        margin-right: 8px;
      }

      .switch {
        display: flex;
        align-items: center;
        .label {
          font-size: 10px;
          font-weight: 700;
          margin-right: 4px;
        }
      }
    }

    .tips {
      line-height: 12px;
      font-size: 10px;
      font-weight: 700;
      color: rgba(173, 181, 189, 1);
    }
    .date-wrap {
      margin-top: 16px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: space-between;

      .date {
        width: 242px;
        height: 32px;
        border-radius: 8px;
        padding: 8px;
        display: flex;
        align-items: center;
        flex-wrap: nowrap;
        background-color: rgba(248, 249, 250, 1);
        .label  {
          font-size: 12px;
          font-weight: 600;
          margin-right: 12px;
          white-space: nowrap;
        }
        .value  {
          font-size: 12px;
          font-weight: 400;
          margin-right: 12px;
          white-space: nowrap;
        }

        .icon-caret-down {
          width: 16px;
          height: 16px;
          background: url('~@/assets/store-visit/icon_caret_down.png') center/100%  no-repeat;
        }
      }

      .btn {
        width: 48px;
        height: 24px;
        border-radius: 4px;
        background-color: rgba(13, 110, 253, 1);
        color: #FFFFFF;
        font-size: 12px;
        font-weight: 600;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

  }
  .legend {
    height: 28px;
    overflow-x: auto;
    &::-webkit-scrollbar {
      display: none;
    }
    .list {
      padding: 0 24px;
      width: max-content;
      display: flex;
      flex-wrap: nowrap;
    }
    .item {
      display: flex;
      align-items: center;
      height: 28px;

      &:not(:last-child) {
        margin-right: 32px;
      }

      .item-color {
        width: 12px;
        height: 8px;
        margin-right: 8px;
      }

      .item-text {
        font-size: 10px;

      }
    }
  }

  .content { 
    flex: 1;
    overflow: hidden;
    position: relative;

    .dashboard {
      position: absolute;
      bottom: 0;
      width: 343px;
      padding: 16px;
      bottom: 16px;
      background-color: #fff;
      border-radius: 12px;
      left: 50%;
      transform: translateX(-50%);
      z-index: 1001;

      .list {
        display: flex;
        justify-content: space-between;
        flex-wrap: wrap;

        .item {
          width: 153px;
          height: 52px;
          border-radius: 8px;
          background-color: rgba(248, 249, 250, 1);
          box-sizing: border-box;
          padding: 8px 12px;

          &:nth-child(1),
          &:nth-child(2) {
            margin-bottom: 4px;
          }

          .item-name {
            font-size: 10px;
            font-weight: 700;
            line-height: 12px;
            color: rgba(73, 80, 87, 1);
          }

          .flex {
            height: 24px;
            display: flex;
            justify-content: space-between;
            align-items: center;

            .item-label {
              font-size: 10px;
              font-weight: 700;
              color: rgba(153, 153, 153, 1);
            }
            .item-value {
              font-size: 16px;
              font-weight: 700;
              color: rgba(0, 0, 0, 1);
            }
          }
        }
      }

      .data {
        .data-head {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 0 8px;

          .data-name {
            font-size: 12px;
            font-weight: 600;
          }

          .data-more {
            display: flex;
            align-items: center;
            font-size: 12px;
          }
        }

        .data-content {
          margin-top: 16px;
          border-radius: 8px;
          background-color: rgba(248, 249, 250, 1);
          display: flex;
          align-items: center;
          justify-content: space-between;
          height: 52px;
          padding: 8px 12px;
          box-sizing: border-box;

          .data-label {
            font-size: 10px;
            font-weight: 700;
            color: rgba(153, 153, 153, 1);
          }
          .data-value {
            font-size: 16px;
            font-weight: 700;
            color: rgba(0, 0, 0, 1);
          }
        }
      }
    }
  }

  .footer {
    background-color: #fff;
    padding-bottom: calc(env(safe-area-inset-bottom));
    padding-bottom: calc(constant(safe-area-inset-bottom));
    z-index: 10;
  }

  .tabbar {
    height: 60px;
    display: flex;
    justify-content: center;

    .tabbar-item {
      width: 40%;
      text-align: center;
      font-size: 16px;
      line-height: 60px;
      font-weight: 400;
      height: 60px;

      &.active {
        color: rgba(0, 0, 0, 1);
        font-weight: 600;
      }
      &:not(.active) {
        color: rgba(0, 0, 0, 0.2);
      }
    }
  }
}

</style>