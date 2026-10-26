<template>
  <layout-form>
    <template slot="content">
      <el-form ref="form" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="职位名称" prop="position">
          <el-select v-model="formData.position" placeholder="请选择职位" :disabled="disabled">
            <el-option v-for="(item, index) in selectOptions" :key="index" :label="item.label" :value="item.label" />
          </el-select>
        </el-form-item>
        <el-form-item label="员工名称" prop="member_id">
          <el-button v-if="memberList.length==0" type="primary" size="small" @click="handleSelectMember">选择员工</el-button>
          <el-table v-if="memberList.length>0" :data="[{data: memberList}]" style="width:400px">
            <el-table-column label="员工" min-width="220">
              <template slot-scope="{ row }">
                <el-tag class="tag-item m-r-10 m-b-10" v-for="(item, index) in row.data" :key="index" type="info">
                  <div class="flex">
                    <svg-icon class="icon" icon-class="user"></svg-icon>
                    <span>{{item.member_name}}</span>
                  </div>
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="!disabled" label="操作" width="140">
              <template>
                <el-button type="text" @click="handleSelectMember">修改</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item label="管理范围" prop="area_ids">
          <!-- 点击按钮 -->
          <el-button v-if="showList.length==0" class="m-b-20" type="primary" size="small" @click="handleSelectStore">选择管理范围</el-button>
          <!-- 选中数据展示 -->
          <el-table v-if="showList.length>0" :data="[{data: showList}]">
            <el-table-column label="管理范围" min-width="220">
              <template slot-scope="{ row }">
                <div id="list" :class="['list', {'expand': isExpand || type!==3}]">
                  <el-tag class="tag-item m-r-10 m-b-10" v-for="(item, index) in row.data" :key="index" type="info">
                    <div class="flex">
                      <svg-icon v-if="item.type===4" class-name="icon" icon-class="store"></svg-icon>
                      <svg-icon v-else class-name="icon" icon-class="organization"></svg-icon>
                      <span>{{item.name}}</span>
                    </div>
                  </el-tag>
                </div>
                <div v-if="type==3" class="flex-center">
                  <el-button class="more" v-if="showExpandBtn" type="text"  @click="handleExpand">
                    <span>{{isExpand?'收起':'展开'}}</span>
                    <i
                      aria-hidden="true"
                      class="el-icon--right"
                      :class="
                        isExpand ? 'el-icon-arrow-up' : 'el-icon-arrow-down'
                      "
                    ></i>
                  </el-button>
                </div>
              </template>
            </el-table-column>
            <el-table-column v-if="!disabled" label="操作" width="140">
              <template>
                <el-button type="text" @click="handleSelectStore">修改</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
      <select-member
        v-if="!disabled && showMemberDialog"
        :visible="showMemberDialog"
        :limit="1"
        :hasSelectArr="hasSelectMemberList"
        :filterMemberMethod="filterNodeMethod"
        @close="handleSelectMemberCancel"
        @submit="handleSelectMemberConfirm"
      ></select-member>

      <com-select-tree 
        v-model="showStoreDialog" 
        :checked="formData.area_ids" 
        :treeList="storeTreeList" 
        @confirm="handleSelectStoreConfirm"
        @cancel="handleSelectStoreCancel"
      ></com-select-tree>
    </template>
    <template slot="footer">
      <el-button v-if="type === 3" type="primary" @click="handleBack">返回</el-button>
      <el-button v-if="type !== 3" @click="handleBack">取消</el-button>
      <el-button v-if="type !== 3" type="primary" @click="handleSave">保存</el-button>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import LayoutForm from "@/components/public/layout-form/index.vue";
import SelectMember from "@/components/business/select-member-new/index.vue";
import ComSelectTree from './com-select-tree.vue'
import Api from '@/api/enterprise/org-structure';
import Api2 from '@/api/enterprise/permission';
import { forEachTree, treeToList } from "@/utils/tree-tool";

interface IFormData {
  id?: number;
  /** 职位名称 */
  position: string;
  /** 员工ID */
  member_id: number;
  /** 员工姓名 */
  member_name: string;
  /** 管理范围 */
  area_ids: any[];
}

@Component({
  components: {
    LayoutForm,
    SelectMember,
    ComSelectTree
  }
})
export default class PermissionCreateForm extends Vue {
  /** 1-add;2-edit;3-detail */
  @Prop({ default: 1 }) type: number;
  selectOptions = [
    {
      label: "GM(General Manager)",
      value: 1
    },
    {
      label: "CL(Channel Leader)",
      value: 2
    },
    {
      label: "RD(Regional Leader)",
      value: 3
    },
    {
      label: "MM/CL(Multi Manager/City Manager)",
      value: 4
    },
    {
      label: "Super Management",
      value: 5
    },
  ]
  /** 表单数据 */
  formData: IFormData = {
    position: "",
    member_id: null, // 清空2
    member_name: "",
    area_ids: []
  }

  showMemberDialog = false;
  showStoreDialog = false;

  /** 表单校验规则 */
  formRules = {
    position: [
      { required: true, message: "请选择职位", trigger: "blur" }
    ],
    member_id: [
      { required: true, message: "请选择员工", trigger: "blur" }
    ],
    area_ids: [
      { required: true, message: "请选择管理范围", trigger: "blur" }
    ]
  }
  /** 门店树形列表 */
  storeTreeList = []

  /** 已有的memberid */
  memberIds = []

  hasSelectMemberList = []

  showExpandBtn= false;
  isExpand = false;

  get disabled() {
    return this.type === 3
  }

  get memberList() {
    if(!this.formData.member_id) return []
    return [{
      member_id: this.formData.member_id,
      member_name: this.formData.member_name
    }]
  }

  get showList() {
    const checked = this.formData.area_ids
    const list = JSON.parse(JSON.stringify(this.storeTreeList))
    const result = []
    forEachTree<any>(list, (item: any) => {
      if(checked.includes(item.id)) {
        result.push(item)
      }
    });
    // 展开
    const nodeList = treeToList<any>(JSON.parse(JSON.stringify(result)))
    // 去重
    return nodeList.filter((item: any, index: number) => {
      return nodeList.findIndex((i: any) => i.id === item.id) === index
    })
  }

  @Watch('showList') 
  handleShowListChange(value) {
    const len = value.length
    this.$nextTick(() => {
      const container = document.getElementById("list");
      const content = document.getElementsByClassName("tag-item");
      const last: HTMLElement = content[len - 1] as HTMLElement;
        if(last) {
          this.showExpandBtn = last.offsetTop > container.clientHeight
        }
      })
  }

  created() {
    this.getStoreTreeList()
    this.getMemberIds()
    if(this.$route.query.id) {
      this.getDetail()
    }
  }

  getStoreTreeList() {
    Api.getDashboardTree({
      has_children: true
    }).then((res: any) => {
      this.storeTreeList = res.data
    })
  }

  getMemberIds() {
    Api2.getMemberIdsApi({}).then((res: any) => {
      this.memberIds = res.data
    })
  }

  getDetail() {
    Api2.getManagerAuthDetailApi({ id: Number(this.$route.query.id) }).then((res: any) => {
      this.formData.id = res.data.id
      this.formData.area_ids = res.data.area_ids
      this.formData.position = res.data.position
      this.formData.member_id = res.data.member_info.id
      this.formData.member_name = res.data.member_info.name
      this.hasSelectMemberList = [{
        member_id: res.data.member_info.id,
        mobile: res.data.member_info.mobile,
        user_id: res.data.member_info.userid,
        name: res.data.member_info.name
      }]
    })
  }

  filterNodeMethod(value: any[]) {
    console.log(value);
    
    const memberIds = this.memberIds.filter(id => id !== this.formData.member_id)
    return value.filter((item: any) => {
      return !memberIds.includes(item.member_id)
    })
  } 

  /**
   * 选择员工
   */
  handleSelectMember() {
    this.showMemberDialog = true;
  }

  /**
   * 选择员工 - 取消
   */
  handleSelectMemberCancel() {
    this.showMemberDialog = false;
  }

  /**
   * 选择员工 - 确认
   */
  handleSelectMemberConfirm(e: any) {
    this.formData.member_id = e[0].member_id;
    this.formData.member_name = e[0].name;
    this.showMemberDialog = false;
    this.hasSelectMemberList = e
  }

  /**
   * 选择管理范围
   */
  handleSelectStore() {
    this.showStoreDialog = true;
  }
  /**
   * 选择管理范围 -- 确定
   * @param value 
   */
  handleSelectStoreConfirm(value: any[]) {
    this.showStoreDialog = false;
    this.formData.area_ids = value
    
  }

  /**
   * 选择管理范围 -- 取消
   */
   handleSelectStoreCancel() {
    this.showStoreDialog = false;
  }

  /**
   * 返回
   */
  handleBack() {
    this.$router.back();
  }

  /**
   * 保存
   */
  handleSave() {
    (this.$refs.form as any).validate(valid => {
      if (valid) {
        Api2.saveManagerAuthApi(this.formData).then((res: any) => {
          this.$message.success("保存成功");
          this.$router.back();
        })
      }
    });
  }

  /**
   * @description 展开收起
   */
   handleExpand() {
    this.isExpand = !this.isExpand
  }
}
</script>
<style lang="scss" scoped>
.el-select {
  width: 400px;
}

.icon {
  margin-right: 5px;
}

.flex {
  display: flex;
  align-items: center;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}
.list {
  max-height: 120px;
  overflow: hidden;

  &.expand {
    max-height: none;
    height: auto;
  }
}
</style>