<template>
  <div>
    <van-popup v-model="visible" position="bottom" round :close-on-click-overlay="false">
      <div class="popup">
        <div class="popup-header">
          <div class="popup-title">
            {{ $t('storeVisit.tips.11') }}
            <com-popover placement="bottom" theme="light">
              <template #reference>
                <div class="popup-tip"></div>
              </template>
              <div class="tip-content">
                {{ $t('storeVisit.tips.44') }}
              </div>
            </com-popover>
          </div>
          <!-- <div class="popup-close" @click.stop="handleCancel"></div> -->
          <div class="value" @click="handleAddVisitor">
            <div>{{ $t('storeVisit.btn.28') }}</div>
            <div class="icon-right"></div>
          </div>
        </div>
        <div class="popup-content">
          <div class="row">
            <div class="list">
              <div class="item" v-for="(item, index) in managerList" :key="index">
                <div class="flex-1">
                  <div class="item-name">@ {{item.name||item.email}}</div>
                  <div class="item-dept">{{ item.position_abbr }}</div>
                </div>
                <div class="item-action">
                  <div class="icon-delete" @click="handleDelete(index, 'manager')"></div>
                </div>
              </div>
              <div class="item" v-for="(item, index) in copyMemberList" :key="index">
                <div class="flex-1">
                  <div class="item-name">@ {{item.name||item.email}}</div>
                  <div class="item-dept" :class="{ext: item.dept_name === '' || item.dept_name === 'Email'}">{{ item.dept_name|| 'Email' }}</div>
                </div>
                <div class="item-action">
                  <div class="icon-delete" @click="handleDelete(index, 'copyMember')"></div>
                </div>
              </div>
              <div class="item" v-for="(item, index) in visitorMemberList" :key="index">
                <div class="flex-1">
                  <div class="item-name">@ {{item.name||item.email}}</div>
                  <div v-if="item.is_mm" class="item-dept">MM</div>
                  <div v-else class="item-dept" :class="{ext: item.dept_name === '' || item.dept_name === 'Email'}">{{item.dept_name || 'Email' }}</div>
                </div>
                <div class="item-action">
                  <div class="icon-delete" @click="handleDelete(index, 'visitor')"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="popup-footer">
          <van-button color="#000000" plain @click.stop="handleCancel">{{ $t('storeVisit.popup.btn.3') }}</van-button>
          <van-button color="#000000" @click.stop="handleConfirm">{{ $t('storeVisit.popup.btn.4') }}</van-button>
        </div>
      </div>
    </van-popup>
    <!-- 选择 -->
    <select-visitor-popup v-model="showSelectVisitorPopup" :default-selection="selection" @select="handleVisitorSelect" @cancel="handleSelectVisitorCancel" @confirm="handleSelectVisitorConfirm"></select-visitor-popup>
  </div>
</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import SelectVisitorPopup from "./select-visitor-popup.vue";
import ComPopover from './com-popover.vue';

@Component({
  components: {
    SelectVisitorPopup,
    ComPopover,
  }
})
export default class EditVisitorPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: Object, default: () =>({}) }) detail: StoreVisit.IDetail;
  @Prop({ type: String, default: '' }) entrance: string;
  
  visible: boolean = false;
  
  useStoreVisitor = 0;
  useOtherVisitor = 0;
  showSelectVisitorPopup = false;

  showPopover = false;

  /** 默认选中 */
  managerList: any[] = []
  copyMemberList: any[] = []
  visitorMemberList: any[] = []

  get selection() {
    return this.managerList.concat([...this.copyMemberList, ...this.visitorMemberList])
  }
  
  @Watch('show', { immediate: true })
  onShowChange(val: boolean) {
    this.visible = val;
    if(val) {
      this.useStoreVisitor = this.detail.store_is_show;
      this.useOtherVisitor = this.detail.visitor_is_show;
      this.visitorMemberList = JSON.parse(JSON.stringify(this.detail.visitor_member || []));
      this.copyMemberList = JSON.parse(JSON.stringify(this.detail.copy_member || []));
      this.managerList = JSON.parse(JSON.stringify(this.detail.store_manager || []));
      
    }
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
  }

  /**
   * 确认
   */
  handleConfirm = this.$helper.debounce(() => {
    // 触发确认事件

    this.emitConfirm();
  }, 0, true)
  
  /**
   * 确认事件
   */
  emitConfirm() {
    // this.$emit('update:show', false);

    // 删除的巡店人员
    const deleteVisitor = this.detail.visitor_member.filter(item => !this.visitorMemberList.some(sel => 
      sel.member_id === item.member_id
      && sel.email === item.email
    )).map(item=> {
      return {
        member_id: item.member_id,
        email: item.email,
        dept_name: item.dept_name,
        name: item.name
      }
    })
    // 新增的巡店人员
    const addVisitor = this.visitorMemberList.filter(item => !this.detail.visitor_member.some(sel => 
      sel.member_id === item.member_id
      && sel.email === item.email
    )).map(item=> {
      return {
        member_id: item.member_id,
        email: item.email,
        dept_name: item.dept_name,
        name: item.name
      }
    })

    // 删除的门店店长
    const deleteManager = this.detail.store_manager.filter(item => !this.managerList.some(sel => 
      sel.member_id === item.member_id
    )).map(item=> {
      return {
        member_id: item.member_id,
        email: item.email,
        dept_name: item.dept_name,
        name: item.name
      }
    })
    // 新增的门店店长
    const addManager = this.managerList.filter(item => !this.detail.store_manager.some(sel => 
      sel.member_id === item.member_id
    )).map(item=> {
      return {
        member_id: item.member_id,
        email: item.email,
        dept_name: item.dept_name,
        name: item.name
      }
    })

    this.$emit('confirm', {
      store_is_show: this.useStoreVisitor,
      visitor_is_show: this.useOtherVisitor,
      del_visitor_list: [...deleteVisitor, ...deleteManager],
      add_visitor_list: addVisitor,
      add_manager_list: addManager,
    });
  }
  
  /**
   * 添加巡店人员
   */
  handleAddVisitor() {
    this.showSelectVisitorPopup = true;
  }

  handleSelectVisitorCancel() {
    this.showSelectVisitorPopup = false;
  }

  handleSelectVisitorConfirm(data) {
    const copyIds = this.copyMemberList.map(item => item.member_id)
    const managers = data.filter(item => item.position_abbr && this.detail.third_id ===item.third_id)
    const visitors = data.filter(item => !(item.position_abbr && this.detail.third_id ===item.third_id))
    const copyMember = visitors.filter(item => copyIds.includes(item.member_id))
    
    this.showSelectVisitorPopup = false;
    this.$emit('update:list', data);
    this.visitorMemberList = visitors.filter(item => !copyIds.includes(item.member_id));
    this.managerList = managers;
    this.copyMemberList = copyMember;
  }

  handleDelete(index, type) {
    if(type === 'manager') {
      this.managerList.splice(index, 1);
    } else if(type === 'copyMember') {
      this.copyMemberList.splice(index, 1);
    } else {
      this.visitorMemberList.splice(index, 1);
     }
  }

  handleVisitorSelect(e) {
    const { track } = e
    // 神策数据上报 -- 选择门店
    this.$track.trackChooseStaff({
      SV_id: this.detail.id,
      store_name: this.detail.store_name,
      store_code: this.detail.store_code,
      SV_initiator: this.detail.creator.name,
      SV_initiator_id: this.detail.creator.user_id,
      SV_initiator_department: this.detail.creator.dept_name,
      visit_time: this.detail.date + ' ' + this.detail.time_span,
      choose_scene: '报告查看权限',
      choose_entrance: this.detail.is_report ? '新增':'添加',
      ...track,
    })
  }

}
</script>
<style lang="scss" scoped>
.icon-right {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_right.png) center/16px no-repeat;
}
.icon-delete {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_popup_close.png) center/100% no-repeat;
}

.van-switch {
  width: 26px;
  height: 16px;
  border: none;
  background-color: rgba(173, 181, 189, 1);
  overflow: hidden;
  // display: block;
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

.popup {
  position: relative;
  width: 100vw;
  background: #FFFFFF;
  .popup-header {
    padding: 20px 24px;
    position: relative;
    
    .popup-title {
      font-weight: 600;
      font-size: 16px;
      line-height: 24px;
      color: #000000;
      display: flex;
      align-items: center;
    }

    .popup-close {
      width: 72px;
      height: 64px;
      position: absolute;
      right: 0;
      top: 0;
      background: url(~@/assets/store-visit/icon_popup_close.png) center/24px no-repeat;
    }
    .value {
      font-size: 14px;
      color: rgba(13, 110, 253, 1);
      display: flex;
      align-items: center;
      position: absolute;
      right: 24px;
      top: 50%;
      transform: translateY(-50%);
    }
  }

  .popup-tips {
    padding: 12px 24px;
    font-size: 12px;
    line-height: 130%;
  }

  .popup-content {
    padding: 8px 0;
    height: 350px;
    overflow-y: auto;
    &::-webkit-scrollbar {
      width: 0;
    }

    .row { 
      padding: 0 24px;

      .row-head {
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 28px;
        .label {
          font-size: 14px;
        }
      }

      .add {
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 28px;
        margin-top: 12px;
        .label {
          font-size: 14px;
          color: rgba(13, 110, 253, 1);
        }

      }

      .list {
        margin-top: 12px;

        .item {
          height: 32px;
          background-color: rgba(248, 249, 250, 1);
          border-radius: 4px;
          display: flex;
          align-items: center;
          padding: 0 8px;

          &:not(:last-child) {
            margin-bottom: 4px;
          }

          .flex-1 {
            flex: 1;
            display: flex;
            align-items: center;
          }


          .item-name {
            font-size: 12px;
            color: rgba(0, 0, 0, 1);
          }

          .item-dept {
            font-weight: 700;
            font-size: 10px;
            color: rgba(0, 0, 0, 1);
            background-color: rgba(206, 212, 218, 1);
            line-height: 16px;
            padding: 0 4px;
            border-radius: 2px;
            margin-left: 4px;

            &.ext {
              color: rgba(153, 153, 153, 1);
              background-color: rgba(255, 255, 255, 1);
            }
          }
        }
      }
    }
  }

  .popup-footer {
    box-sizing: border-box;
    z-index: 2;
    display: flex;
    align-items: center;
    background-color: #fff;
    padding: 12px 24px;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);
    border-top: 1px solid rgba(233, 236, 239, 1);

    .van-button {
      height: 40px;
      border-radius: 4px;
      flex: 1;

      &:not(:last-child) {
        margin-right: 8px;
      }

      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }
    }
  }
}
</style>

<style lang="scss"> 


.popup-tip {
  width: 24px;
  height: 24px;
  background: url(~@/assets/store-visit/icon_wen.png) center/16px no-repeat;
}
.tip-content {
  font-size: 12px;
  font-weight: 400;
  width: 300px;
  box-sizing: border-box;
  color: #666;
}

</style>