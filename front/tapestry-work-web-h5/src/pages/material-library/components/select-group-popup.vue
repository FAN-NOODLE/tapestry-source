<template>
  <!--分组弹窗 -->
  <van-popup
    v-model="visible"
    :close-on-click-outside="false"
    :close-on-click-action="false"
    position="bottom"
  >
    <div class="popup">
      <div class="popup-content">
        <div class="popup-aside">
          <div
            v-for="(item, index) in list"
            :key="index"
            :class="['popup-aside-item', {
              active: isActive(0, index, selected)
            }]"
            
            @click="handleChange([index, -1, -1, -1], 0, item)"
          >
            <div class="aside-item-text ellipsis-t">{{ item.group_name }}</div>
          </div>
        </div>
        <div class="popup-main" v-if="selected[0] !== -1">
          <!-- 二级分类 -->
          <div :class="[
            'item', 
            {
              expand: item.expand
            }
            ]" 
            v-for="(item, index) in list[selected[0]].children" 
            :key="index">
            <div
              :class="[
                'item-title',
                {
                  active: isActive(1, index, selected),
                  'has-child': item.children && item.children.length > 0, 
                  disabled: !item.is_show
                }
              ]"
              data-key="1"
              @click="handleChange([selected[0], index, -1, -1], 1, item)"
            >
              <div class="item-title__text">{{ item.group_name }}</div>
            </div>
            <!-- 三级分类 -->
            <div v-if="item.children && item.children.length" class="item-content">
              <div :class="['item', {expand: third.expand}]" v-for="(third, thirdIndex) in item.children" :key="thirdIndex">
                <div
                  :class="[
                    'item-title',
                    {
                      active: isActive(2, thirdIndex, selected),
                      'has-child': third.children && third.children.length > 0, 
                      disabled: !third.is_show
                    }
                  ]"
                  data-key="2"
                  @click="handleChange([selected[0], index, thirdIndex, -1], 2, third)"
                >
                  <div class="item-title__text">{{ third.group_name }}</div>
                </div>
                <!-- 四级分类 -->
                <div v-if="third.children && third.children.length" class="item-content">
                  <div class="item" v-for="(four, fourIndex ) in third.children" :key="fourIndex">
                    <div
                      :class="[
                        'item-title',
                        'hide-icon',
                        {
                          active: isActive(3, fourIndex, selected),
                          disabled: !four.is_show
                        }
                      ]"
                      data-key="3"
                      @click="handleChange([selected[0], index, thirdIndex, fourIndex], 3, four)"
                    >
                      <div class="item-title__text">{{ four.group_name }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="popup-footer">
        <div class="btn plain" @click="handleCancel">取消</div>
        <div class="btn" @click="handleConfirm">确定</div>
      </div>
    </div>
  </van-popup>
</template>
<script lang="ts">
import { Component, Vue, PropSync, Prop, Watch } from "vue-property-decorator";
import { Field, Loading, Popup } from "vant";
import { getListForStore } from "@/api/marketing-promotion/index";
function traverseTree(tree) {
  if (tree.length === 0) return
  tree.forEach(item => {
    item.children = item.children || []
    item.expand = false
    traverseTree(item.children)
  })
}

/** 主分组 */
interface IGroup {
  created_at: number;
  deleted_at: number;
  enterprise_id: number;
  group_name: string;
  id: number;
  is_show: number;
  orders: number;
  parent_id: number;
  updated_at: number;
  children: Array<IGroup>;
  expand: boolean;
}

@Component({
  components: {
    [Loading.name]: Loading,
    [Popup.name]: Popup,
    [Field.name]: Field,
  },
})
export default class SelectGroupPopup extends Vue {
  @PropSync("show", { type: Boolean }) visible: boolean;
  @Prop({ type: String, default: '未分组' }) text: string;
  @Prop({ type: Array, default: () => [] }) value: Array<any>;

  /**分组下的树结构子分组 */
  list: Array<IGroup> = [];

  selected = [-1, -1, -1, -1]

  get valueName() {
    const { list, selected } = this
    const value = selected.reduce((prev, item, index) => {
      if (item === -1) {
        return prev
      } else if (index === 0) {
        prev.push(list[item])
        return prev
      } else {
        prev.push(prev[index - 1].children[item])
        return prev
      }
    }, [])
    const name = value.filter(item => item && item.id).map(item => item.group_name).join('/')
    return name
  }

  get isActive() {
    return (index, current, selected) => {
      if (index === 3) return current === selected[3]
      if (index === 2) return current === selected[2] && selected[3] === -1
      if (index === 1) return current === selected[1] && selected[2] === -1
      if (index === 0) return current === selected[0]
    }
  }

  @Watch('value')
  handleValueChange(value) {
    this.init()
  }

  /**生命周期 */
  mounted() {
    this.getMaterialCategory()
  }

  getMaterialCategory() {
    getListForStore().then(res => {
      res.data.unshift({
        id: this.text === '未分组' ? -1 : 0,
        parent_id: 0,
        group_name: this.text,
        children: []
      })
      traverseTree(res.data)
      this.list = res.data.map(item => {
        item.children.unshift({
          id: item.id,
          parent_id: item.parent_id,
          group_name: '全部',
          children: [],
          is_show: item.is_show
        })
        return item
      })
      this.init()
    }).catch(err => {
      console.log(err)
    })
  }
  /**
   * @description 初始化
   */
  init() {
    const { list, value } = this
    if (list.length === 0) return
    // 获取下标 index[]
    const selected = value.reduce((prev, id, index) => {
      if (index === 0) {
        const i = list.findIndex(item => Number(item.id) === Number(id))
        prev.push([i, list[i]])
        return prev
      } else {
        const i = prev[index - 1][1].children.findIndex(item => Number(item.id) === Number(id))
        prev.push([i, prev[index - 1][1].children[i]])
        return prev
      }
    }, []).map(item => item[0])

    // 填充
    while (selected.length < 4) {
      selected.push(-1)
    }

    // 处理空值
    if (selected[0] === -1) {
      selected[0] = 0
    }

    this.selected = selected
    
  }
    
  /**
   * @description 切换
   */
  handleChange(value, key, item) {
    if(key> 0 && item.is_show==0) return
    
    if (key === 1 || key === 2) {
      item.expand = !item.expand
    }
    this.selected = value
  }
  /**
   * @description 确定
   */
  handleConfirm() {
    const { list, selected } = this
    // 获取选中项
    const value = selected.reduce((prev, item, index) => {
      if (item === -1) {
        return prev
      } else if (index === 0) {
        prev.push(list[item])
        return prev
      } else {
        prev.push(prev[index - 1].children[item])
        return prev
      }
    }, [])
    // 如果选中的是is_show == 0的选项直接过滤
    const isEmpty = value.every(item => item.is_show === 0)
    if (isEmpty) return
    
    this.$emit('confirm', {
      value: value.map(item => item.id),
      valueName: this.valueName
    })
  }
  /**
   * @description 取消
   */
  handleCancel() {
    this.visible = false
    this.$emit('cancel')
  }
}
</script>
<style lang="scss" scoped>
.select {
  padding: 15px 0;
  display: flex;
  align-items: center;
  text-align: right;
  &__value {
    flex: 1
  }

  &__icon {
    width: 20px;
    height: 20px;
    background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_right_gray.png") no-repeat center center /100%;
  }  
}
.popup {
  width: 100vw;
  height: 80vh;
  background: #F4F7FB;
  box-shadow: 0px 6px 12px 0px rgba(222, 224, 227, 0.5);
  display: flex;
  flex-direction: column;
  overflow: hidden;

  .popup-content {
    flex: 1;
    display: flex;
    overflow: hidden;
    .popup-aside {
      width: 116px;
      overflow-y: scroll;
  
      .popup-aside-item {
        width: 116px;
        height: 60px;
        padding: 10px 16px;
        box-sizing: border-box;
        font-size: 14px;
        color: #777777;
        line-height: 20px;
        box-sizing: border-box;
        display: flex;
        align-items: center;
        justify-content: flex-start;
  
        &.active {
          color: #404040;
          background: #ffffff;
          font-weight: bold;
        }
        
  
        .aside-item-text {
          word-wrap: break-word;
          text-align: left;
        }
      }
    }
        
      
    
    .popup-main {
      flex: 1;
      overflow: hidden;
      overflow-y: scroll;
      background: #fff;
      box-sizing: border-box;
      padding: 10px 20px 10px 15px;
  
      .item-title {
        height: 40px;
        border-radius: 8px;
        color: #777777;
        padding: 0 15px;
        box-sizing: border-box;
        font-size: 14px;
        position: relative;
        display: flex;
        align-items: center;
        justify-content: flex-start;
        &__text {
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        
  
        &::before {
          content: "";
          width: 14px;
          height: 14px;
          margin-right: 3px;
          display: block;
          visibility: hidden;
          transform-origin: center center;
          transform: rotateZ(-90deg);
          transition: all 0.3s ease-in-out;
          background: url("https://cdn-cos.coach.com.cn/uploads/2023/v4.0.1/material/icon_arrow_down_2.png") center/ 100% no-repeat;
        }
  
        &.has-child {
          &::before {
            visibility: visible;
          }
        }
        
  
        &.hide-icon {
          &::before {
            visibility: hidden;
          }
        }
          
        
  
        &.active {
          background: #F4F7FB;
          color: #777777;
          font-weight: 500;
        }
        
        &[data-key="1"] {
          color: #000;
          font-weight: bold;
        }
        &[data-key="2"] {
          padding-left: 30px;
        }
        &[data-key="3"] {
          padding-left: 45px;
        }
      
        &.disabled {
          color: #ccc;
        }
      }
  
      .item-content {
        height: 0;
        overflow: hidden;
        transition: all 0.3s ease-in-out;
      }
      
      .item.expand {
        &>.item-content {
          height: auto;
        }
        
        &>.item-title::before {
          transform: rotateZ(0deg)
        }
      }
    }
  }
  
  .popup-footer {
    width: 100vw;
    padding: 15px;
    box-sizing: border-box;
    background: #ffffff;
    display: flex;
    justify-content: space-between;
    .btn {
      width: 163px;
      height: 40px;
      background: #000;
      box-sizing: border-box;
      font-size: 13px;
      color: #fff;
      display: flex;
      align-items: center;
      justify-content: center;
      &.plain {
        color: #000;
        background-color: #fff;
        border: 1px solid #000;
      }
    }
  }
}
</style>
