<template>
  <div>
    <van-popup v-model="visible" position="bottom" round @click-overlay="handleCancel">
      <div class="popup">
        <div class="popup-header">
          <div class="popup-title">{{ $t('storeVisit.popup.title.8') }}</div>
          <div class="popup-close" @click.stop="handleCancel"></div>
        </div>
        <div class="popup-tips">{{ $t('storeVisit.popup.tips.31') }}</div>
        <div class="popup-search">
          <div class="input" @click="handleInputClick">
            <div class="icon-search"></div>
            <input class="input-inner" ref="input" type="text" :placeholder="String($t('storeVisit.placeholder.5'))" v-model="searchValue" @input="handleInput" />
            <div class="icon-close" v-if="searchValue.length" @click="handleClear"></div>
          </div>
        </div>
        <div v-if="!showSearchResult" class="popup-content">
          <div class="side">
            <div class="side-item" :class="{active: index==current}" v-for="item,index in list" :key="index" @click="handleClick(index)">{{ item.dept_name }}</div>
          </div>
          <div class="main">
            <div class="list">
              <div v-if="current === list.length -1" class="item item-add" @click="handleAddEmail">
                <div class="flex">
                  <div class="item-name">{{ $t('storeVisit.popup.btn.12') }}</div>
                  <div class="icon-right"></div>
                </div>
              </div>
              <div 
                class="item" 
                :id="`id_${current}_${index}`"
                v-for="item,index in memberList" 
                :key="index" 
                :class="{
                  active: checked(item),
                  disabled: disabledSelect(item)
                }"
                >
                <div class="flex" @touchstart="handleTouchStart($event, item)" @touchend="handleTouchEnd($event, item)" @touchmove="handleTouchMove">
                  <div class="item-name">{{ item.name||item.email }}</div>
                  <div class="item-checkbox"></div>
                </div>
                <!-- 编辑 删除 -->
                <div v-if="item.isEdit" class="item-mask" @click="item.isEdit=false"></div>
                <div class="item-opt" :class="{show: item.isEdit}">
                  <div class="delete" @click.stop="handleDelete(index)">
                    <div class="icon-delete"></div>
                    <span>{{ $t('storeVisit.popup.btn.13') }}</span>
                  </div>
                  <div class="edit" @click.stop="handleEdit(index)">
                    <div class="icon-edit"></div>
                    <span>{{ $t('storeVisit.popup.btn.14') }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-if="!showSearchResult" class="popup-footer">
          <div v-if="selection.length" class="popup-selection">
            <div class="label">
              <div>{{ $t('storeVisit.popup.tips.5') }}{{ selection.length }}{{ $t('storeVisit.popup.tips.12') }}：</div>
              <div class="all" @click="handleShowUser">{{ $t('storeVisit.popup.btn.5') }}</div>
            </div>
            <div class="scroll-x">
              <div class="select-list">
                <div class="select-item" v-for="(item, index) in selection" :key="index">
                  <div class="name">{{ item.name||item.email }}</div>
                  <div v-if="item.is_mm" class="dept">MM</div>
                  <div v-else class="dept">{{item.dept_name||'Email' }}</div>
                  <div class="icon-close-circle" @click="handleRemoveSelection(index)"></div>
                </div>
              </div>
            </div>
          </div>
          <div class="flex">
            <van-button color="#000000" plain @click.stop="handleCancel">{{ $t('storeVisit.popup.btn.3') }}</van-button>
            <van-button color="#000000" @click.stop="handleConfirm">{{ $t(confirmText) }}</van-button>
          </div>
        </div>
        <div v-if="showSearchResult" class="popup-content search-result">
          <div v-if="searchResult.length" class="result">
            <div class="list">
              <div class="item" v-for="item,index in searchResult" :key="index" @click="handleSelectItem(item)">
                <div class="item-name">{{item.name||item.email}}</div>
                <div class="item-tag" v-if="item.dept_name">{{item.dept_name}}</div>
              </div>
            </div>
          </div>
          <div v-else class="empty">{{$t('storeVisit.empty.2')}}</div>
        </div>
      </div>
    </van-popup>
    <add-email-popup v-model="showAddEmail" :is-edit="!!currentUser.email" :email="currentUser.email" @back="visible=true" @cancel="handleAddCancel" @confirm="handleAddConfirm"></add-email-popup>
    <!-- 查看已选人员 -->
    <van-popup v-model="showSelectUser" round position="center">
      <div class="popup popup-user">
        <div class="popup-header">
          <div class="popup-title">{{ $t('storeVisit.popup.tips.5') }}{{ selection.length }}{{ $t('storeVisit.popup.tips.12') }}</div>
          <div class="popup-close" @click.stop="handleCloseUser"></div>
        </div>
        
        <scrollbar class="popup-content" max-height="400">
          <div class="user-list">
            <div class="user-item" v-for="(item, index) in selection" :key="index">
              <div class="flex-1">
                <div class="icon-a"></div>
                <div class="user-name">{{item.name||item.email}}</div>
                <div v-if="item.is_mm" class="user-dept">MM</div>
                <div v-else class="user-dept" :class="{ext: item.dept_name === 'Email'}">{{item.dept_name }}</div>
              </div>
              <div class="icon-close-circle-2" @click="handleDeleteUser(index)"></div>
            </div>
          </div>
        </scrollbar>
      </div>
    </van-popup>
  </div>

</template>
<script lang="ts">
import { Component, Model, Prop, Vue, Watch } from 'vue-property-decorator';
import AddEmailPopup from './add-email-popup.vue';
import { deleteEmailApi, getDeptMemberListApi, getEmailContactListApi, saveEmailApi, searchMemberApi } from '@/api/store-visit';
import Scrollbar from './scrollbar.vue'
@Component({
  components: {
    Scrollbar,
    AddEmailPopup
  },
})
export default class SelectVisitorPopup extends Vue {
  @Model('update:show', { type: Boolean, default: false }) show: boolean;
  @Prop({ type: Array, default: () => [] }) defaultSelection: any[];
  @Prop({ type: String, default: 'storeVisit.popup.btn.4' }) confirmText: string;
  
  visible: boolean = false;
  showAddEmail = false;
  showSelectUser = false;

  showSearchResult = false;
  searchValue: string = '';
  searchResult: any[] = []
  
  currentUser:any = {
    email: ''
  };
  
  current = 0;
  list: any[] = []

  selection: any[] = [];

  startX: number = 0;
  startY: number = 0;

  timer:any = null;

  get memberList() {
    if (this.list.length === 0) {
      return [];
    }
    return this.list[this.current].member_list;
  }

  get checked(){
    return (value) => {
      return this.selection.some((item) => {
        if(value.dept_name === 'Email') {
          return item.email === value.email || item.id === value.id
        } else {
          return item.member_id === value.member_id;
        }
      });
    }
  }

  get disabledSelect() {
    return(value) => {
      return false
      // return this.selection.some((item) => {
      //   return item.member_id === value.member_id && item.is_mm
      // });
      // return this.selection.some((item) => {
      //   if(value.dept_name === 'Email') {
      //     return item.email === value.email && !item.is_new
      //   } else {
      //     return item.member_id === value.member_id && !item.is_new
      //   }
      // });
    }
  }

  get disabled() {
    return this.selection.filter(item => item.is_new).length==0
  }

  /** 用户信息 */
  get userInfo():IUserInfo {
    return this.$store.getters["user/getUserInfo"]
  }

  
  @Watch('show', {  immediate: false })
  onShowChange(val: boolean) {
    this.visible = val;
    if (val) {
      this.selection = JSON.parse(JSON.stringify(this.defaultSelection));
    }
  }
  // @Watch('defaultSelection', {  immediate: true })
  // onDefaultSelectionChange(val: any[]) {
  //   this.selection = JSON.parse(JSON.stringify(val));
  // }

  async mounted() {
    await this.getDeptMemberList();
    this.getEmailContactList()
  }

  /**
   * 取消
   */
  handleCancel() {
    this.$emit('update:show', false);
    this.$emit('cancel');
    this.searchValue = '';
    this.showSearchResult = false;
    this.searchResult = [];
  }

  /**
   * 确认
   */
  handleConfirm = this.$helper.debounce(() => {
    // 触发确认事件

    this.emitConfirm();
  }, 500, true)
  
  /**
   * 确认事件
   */
  emitConfirm() {
    // this.$emit('update:show', false);
    this.$emit('confirm', this.selection);
    this.searchValue = '';
    this.showSearchResult = false;
    this.searchResult = [];
  }

  /**
   * 获取部门成员列表
   */
  async getDeptMemberList() {
    const res:any = await getDeptMemberListApi({})
    if (res.code === 200) {
      this.list = res.data.map((item) => {
        return {
          dept_name: item.dept_name,
          member_list: item.member_list.map((member) => {
            return {
              ...member,
              dept_name: item.dept_name,
            }
          })
        }
      });
    }
  }

  /**
   * 获取邮箱联系人
   */
  async getEmailContactList() {
    
    const res:any = await getEmailContactListApi({})
    if (res.code === 200) {
      const last = this.list[this.list.length - 1];
      const data = (res.data||[]).map((item) => {
        return {
          id: item.id,
          member_id: 0,
          name: '',
          email: item.email,
          dept_name: "Email",
          user_id: ""
        }
      })
      if(last && last.dept_name === 'Email') {
        this.list[this.list.length - 1].member_list = data
      } else {
        this.list.push({
          dept_name: 'Email',
          member_list: data
        });
      }
    }
  }

  /**
   * 搜索
   * @param e 
   */
  handleInput = this.$helper.debounce((e) => {
    this.showSearchResult = true;
    this.searchMember(e.target.value)
  }, 500, false)

  /**
   * 搜索成员
   * @param keyword 
   */
  async searchMember(keyword) {
    if (keyword.length === 0) {
      return;
    }
    const res:any = await searchMemberApi({
      str: keyword
    })
    if (res.code === 200) {
      
      this.searchResult = res.data;
      this.showSearchResult = true;
    }
  }

  /**
   * 清空搜索
   */
  handleClear() {
    this.searchValue = '';
    this.showSearchResult = false;
  }
  
  /**
   * 点击获取焦点
   */
  handleInputClick() {
    (this.$refs.input as any).focus()
  }

  /**
   * 切换分组
   * @param index 
   */
  handleClick(index) {
    this.current = index;
  }


  handleRemoveSelection(index) {
    this.selection.splice(index, 1);
  }

  /**
   * 添加邮箱
   */
  handleAddEmail() {
    this.showAddEmail = true;
    this.visible = false;
  }

  /**
   * 选择
   * @param item 
   */
  handleSelect(item) { 
    if(this.disabledSelect(item)) return
    if (this.checked(item)) {
      this.selection = this.selection.filter((i) => {
        if(item.dept_name === 'Email') {
          return i.email !== item.email
        } else {
          return i.member_id !== item.member_id;
        }
      });
    } else {
      // this.$set(item, 'is_new', true)
      this.selection.push({
        is_new: true,
        ...item
      });
      this.$emit('select', {
        track: {
          staff_id: item.user_id,
          staff_name: item.name || item.email,
          staff_department: item.dept_name,
          search_word: '',
          filter_option: item.dept_name,
          seek_method: '筛选',
        },
        item
      })
    }
  }

  handleAddCancel() {
    this.showAddEmail = false;
    this.visible = true
  }
  
  /**
   * 添加邮箱确认
   * @param item 
   */
  async handleAddConfirm(value) {
    const { email, checked, isEdit } = value
    let id = ''
    if(checked) {
      // 保存到通讯录
      const res:any = await saveEmailApi({
        email: email,
        id: this.currentUser.id || undefined,
      })
      if(res.code === 200) {
        id = res.data
        this.getEmailContactList()
        this.currentUser = {
          email: ''
        }
      }
    }
    
    if(isEdit) {
      const index = this.selection.findIndex(item => item.id === id)
      if(index > -1) {
        // 如果已经选中了，则修改邮箱信息
        this.selection[index] = {
          ...this.selection[index],
          email: email
        }
      }
    } else {
      const item = {
        id: id,
        member_id: 0,
        name: '',
        email: email,
        dept_name: 'Email',
        user_id: '',
        is_new: true
      }
      this.selection.push(item);
      this.$emit('select', {
        track: {
          staff_id: '',
          staff_name: email,
          staff_department: 'Email',
          search_word: '',
          filter_option: 'Email',
          seek_method: '添加巡店人邮箱',
          is_save_to_external: checked
        },
        item
      })
      // this.handleConfirm();
      this.visible = true
    }
  }

  /**
   * 选择筛选结果
   * @param item 
   */
  handleSelectItem(item) {
    if (!this.checked(item)) {
      this.selection.push({
        is_new: true,
        ...item
      });
    }
    this.$emit('select', {
      track: {
        staff_id: item.user_id,
        staff_name: item.name || item.email,
        staff_department: item.dept_name,
        search_word: this.searchValue,
        seek_method: '搜索',
        filter_option: ''
      },
      item
    })
    this.showSearchResult = false
    this.searchValue = ''
    this.searchResult = []
    // this.handleConfirm();
  }

  /**
   * 删除
   * @param index 
   */
  async handleDelete(index) {
    const item = this.memberList[index];
    const res:any = await deleteEmailApi({
      id: item.id,
    })
    if(res.code === 200) {
      this.$toast(this.$t('common.toast.21'));
      this.list[this.current].member_list.splice(index, 1);
      // this.getEmailContactList()
    }
  }
  /**
   * 编辑
   * @param index 
   */
  handleEdit(index) {
    const item = this.memberList[index];
    this.currentUser = {
      ...item, 
      isEdit: true,
    };
    this.showAddEmail = true;
    // this.visible = false;
  }

  /**
   * 触摸开始
   * @param item 
   */
   handleTouchStart(e, item) {
    this.startY = e.touches[0].clientY;
    const disabled = this.list[this.current].dept_name !== 'Email';
    
    
    if(disabled || this.disabledSelect(item)) return

    clearTimeout(this.timer);
    this.timer = setTimeout(() => {
      this.$set(item, 'isEdit', true);
    }, 500);
  }

  handleTouchMove(e) {
    const dy = e.changedTouches[0].clientY - this.startY;
    if(Math.abs(dy) > 5) {
      clearTimeout(this.timer);
    }
  }

  handleTouchEnd(e, item) {
    e.preventDefault();
    const dy = e.changedTouches[0].clientY - this.startY;
    clearTimeout(this.timer);
    if(Math.abs(dy) > 5) {
      return;
    }
    if(!item.isEdit) {
      this.handleSelect(item)
    }
  }
  handleShowUser() {
    this.showSelectUser = true;
  }
  handleCloseUser() {
    this.showSelectUser = false;
  }
  handleDeleteUser(index) {
    this.selection.splice(index, 1);
  }

}
</script>
<style lang="scss" scoped>
.icon-close-circle {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_close_circle.png) center/16px no-repeat;
}
.icon-close-circle-2 {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_close_circle_2.png) center/16px no-repeat;
}
.icon-a {
  width: 16px;
  height: 16px;
  background: url(~@/assets/store-visit/icon_a.png) center/16px no-repeat;
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
    }

    .popup-close {
      width: 72px;
      height: 64px;
      position: absolute;
      right: 0;
      top: 0;
      background: url(~@/assets/store-visit/icon_popup_close.png) center/24px no-repeat;
    }

  }

  .popup-tips {
    font-size: 12px;
    color: #96979E;
    line-height: 20px;
    padding:  0 24px;
    margin-top: -10px;
  }

  .popup-search {
    padding: 12px;
    
    .input {
      height: 32px;
      border-radius: 16px;
      background-color: rgba(248, 249, 250, 1);
      padding: 4px 8px;
      box-sizing: border-box;
      display: flex;
      align-items: center;

      .icon-search {
        width: 24px;
        height: 24px;
        background: url(~@/assets/store-visit/icon_search.png) center/24px no-repeat;
      }

      .input-inner {
        flex: 1;
        margin-left: 8px;
        font-size: 12px;
        background-color: transparent;
      }

      .input-inner::placeholder {
        color: rgba(0, 0, 0, 0.2);
      }

      .icon-close {
        width: 24px;
        height: 24px;
        background: url(~@/assets/store-visit/icon_close_circle.png) center/24px no-repeat;

      }
    }
  }

  .popup-content {
    height: 450px;
    overflow-y: auto;
    background-color: rgba(248, 249, 250, 1);
    display: flex;

    .side {
      width: 120px;
      height: 100%;
      overflow-y: auto;
      background-color: #fff;
      &::-webkit-scrollbar {
        display: none;
      }

      .side-item {
        height: 60px;
        line-height: 20px;
        padding: 20px;
        box-sizing: border-box;
        color: rgba(153, 153, 153, 1);
        font-size: 14px;
        font-weight: 400;
        transition: all 0.3s;

        &.active {
          font-weight: bold;
          color: rgba(0, 0, 0, 1);
          background-color: rgba(248, 249, 250, 1);
        }
      }
    }

    .main {
      flex: 1;
      height: 100%;
      overflow-x: hidden;
      overflow-y: auto;
      &::-webkit-scrollbar {
        display: none;
      }

      .list {
        padding: 0 24px;

        .item {
          height: 64px;
          color: #000000;
          position: relative;
          user-select: none;
          .flex {
            height: 100%;
            display: flex;
            align-items: center;
            
          }

          &:not(:last-child) {
            border-bottom: 0.5PX solid rgba(233, 236, 239, 1);
          }

          .item-name {
            font-weight: 400;
            font-size: 14px;
            flex: 1;
            transition: all 0.3s;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            margin-right: 16px;
          }

          .item-checkbox {
            width: 24px;
            height: 24px;
            background: url(~@/assets/store-visit/icon_check.png) center/24px no-repeat;
            opacity: 0;
            transition: all 0.3s;
          }

          .icon-right {
            width: 24px;
            height: 24px;
            background: url(~@/assets/store-visit/icon_right.png) center/24px no-repeat;
            opacity: 1;
          }

          &.active {
            .item-name {
              font-weight: 700;
            }
            .item-checkbox {
              opacity: 1;
            }
          }

          &.disabled {
            opacity: 0.2;
          }

          &.item-add {
            font-weight: 600;
            color: rgba(13, 110, 253, 1);
          }

          .item-mask {
            width: 100vw;
            height: 100vh;
            position: fixed;
            top: 0;
            left: 0;
            z-index: 110;
            // background-color: red;
          }
          .item-opt {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 255px;
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            background-color: rgba(222, 226, 230, 1);
            transition: all 0.3s;
            z-index: 111;
            padding: 0 24px;
            box-sizing: border-box;
            opacity: 0;
            pointer-events: none;

            &.show {
              opacity: 1;
              pointer-events: all;
            }

            .delete {
              display: flex;
              align-items: center;
              justify-content: center;
              margin-right: 12px;
              flex: 1;
              height: 40px;
              background-color: #fff;

              .icon-delete {
                width: 24px;
                height: 24px;
                background: url(~@/assets/store-visit/icon_delete.png) center/24px no-repeat;
              }
              span {
                font-size: 14px;
                color: rgba(0, 0, 0, 1);
                margin-left: 8px;
                font-weight: 600;
              }
            }

            .edit {
              display: flex;
              align-items: center;
              justify-content: center;
              flex: 1;
              height: 40px;
              background-color: #fff;

              .icon-edit {
                width: 24px;
                height: 24px;
                background: url(~@/assets/store-visit/icon_edit.png) center/24px no-repeat;
              }
              span {
                font-size: 14px;
                color: rgba(0, 0, 0, 1);
                margin-left: 8px;
                font-weight: 600;
              }
            }
          }
        }
      }
    }

    &.search-result {
      padding: 12px;
      height: 450px;
      background-color: rgba(255, 255, 255, 1);

      .result {
        
        border-radius: 12px;
        background-color: rgba(248, 249, 250, 1);
        padding: 8px 8px 8px 40px;
        position: relative;
        width: 100%;

        &::before {
          content: "";

          position: absolute;
          left: 8px;
          top: 8px;
          width: 24px;
          height: 24px;
          background: url(~@/assets/store-visit/icon_search_black.png) center/24px no-repeat;
        }

        .list {
          min-height: 166px;
          height: 100%;
          overflow-y: auto;
        }

        .item {
          height: 24px;
          display: flex;
          align-items: center;

          &:not(:last-child) {
            margin-bottom: 12px;
          }
          
          .item-name {
            font-size: 12px;
            color: rgba(0, 0, 0, 1);
          }

          .item-tag {
            height: 16px;
            line-height: 16px;
            background-color: rgba(233, 236, 239, 1);
            font-size: 10px;
            padding: 0 4px;
            margin-left: 4px;
          }
        }
      }

      .empty {
        font-size: 12px;
        color: rgba(153, 153, 153, 1);
        text-align: center;
        padding-top: 50px;
        width: 100%;
        
      }
    }
  }

  .popup-footer {
    box-sizing: border-box;
    z-index: 2;
    background-color: #fff;
    padding-bottom: calc(env(safe-area-inset-bottom) + 12px);
    padding-bottom: calc(constant(safe-area-inset-bottom) + 12px);
    border-top: 1PX solid rgba(233, 236, 239, 1);


    .popup-selection {
      padding: 16px 0 8px;
      .label {
        font-size: 12px;
        color: rgba(0, 0, 0, 1);
        font-weight: 600;
        line-height: 16px;
        padding: 0 24px;
        display: flex;
        align-items: center;
        justify-content: space-between;

        .all {
          font-size: 12px;
          font-weight: bold;
          color: #0D6EFD;
        }
      }

      .scroll-x {
        margin-top: 4px;
        overflow-x: scroll;

        &::-webkit-scrollbar {
          display: none;
        }
      }

      .select-list {
        display: flex;
        flex-wrap: nowrap;
        width: fit-content;
        padding: 0 24px;

        .select-item {
          height: 24px;
          line-height: 24px;
          font-size: 12px;
          color: rgba(0, 0, 0, 0.8);
          border-radius: 3px;
          display: flex;
          align-items: center;

          .name {
            white-space: nowrap;
          }

          .dept {
            background-color: #F8F9FA;
            font-size: 11px;
            line-height: 12px;
            padding: 2px 4px;
            color: #999999;
            margin-left: 4px;
          }

          .icon-close-circle {
            margin-left: 4px;
          }

          &:not(:last-child) {
            margin-right: 8px;

            &::after {
              content: ";";
              margin-left: 8px;
              font-size: 10px;
              color: rgba(153, 153, 153, 1);
            }
          }
        }
      }
    }

    .flex {
      padding: 12px 24px;
      display: flex;
      align-items: center;
    }

    .van-button {
      flex: 1;
      height: 40px;
      border-radius: 4px;
      font-weight: 600;

      &:nth-child(2) {
        margin-left: 8px;
      }

      &.van-button--disabled {
        opacity: 1;
        background: rgba(222, 226, 230, 1) !important;
        border-color: rgba(222, 226, 230, 1) !important;
      }
    }
  }
}

.popup-user {
  width: 311px;

  .popup-content {
    width: 100%;
    height: auto;
    max-height: 400px;
    background: #fff;
    overflow: hidden;
    // overflow-y: scroll;

    // &::-webkit-scrollbar {
    //   display: none;
    //   width: 4px;
    //   height: 4px;
    //   background: rgba(0, 0, 0, 0.1);
    //   border-radius: 2px;
    // }
    // &::-webkit-scrollbar-thumb { 
    //   width: 4px;
    //   background: rgba(0, 0, 0, 0.2);
    //   border-radius: 2px;
    //   &:hover {
    //     background: rgba(0, 0, 0, 0.4);
    //   }
    // }
    
    .user-list {
      box-sizing: border-box;
      width: 100%;
      padding: 16px 24px 24px;
      height: fit-content;
    }
    

    .user-item {
      width: 100%;
      height: 32px;
      background-color: #F8F9FA;
      border-radius: 4px;
      display: flex;
      align-items: center;
      padding-left: 8px;
      .flex-1 {
        flex: 1;
        display: flex;
        align-items: center;
      }

      .user-name {
        font-size: 12px;
        color: rgba(0, 0, 0, 1);
      }

      .user-dept {
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

      .icon-close-circle-2 {
        width: 32px;
        height: 32px;
      }

      &:not(:last-child) {
        margin-bottom: 4px;
      }
    }
  }
}
</style>