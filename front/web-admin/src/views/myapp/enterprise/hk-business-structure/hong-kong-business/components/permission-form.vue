<template>
  <layout-form>
    <template slot="content">
      <el-form ref="form" :model="formData" :rules="rules" label-width="110px">
        <el-form-item label="职位名称" prop="position">
          <el-select
            v-model="formData.position"
            placeholder="请选择职位"
            :disabled="readonly"
            @change="handlePositionChange"
          >
            <el-option
              v-for="item in positionOptions"
              :key="item.level"
              :label="item.label"
              :value="item.label"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="员工名称" prop="member_id">
          <el-button
            v-if="!memberList.length && !readonly"
            type="primary"
            size="small"
            @click="showMemberDialog = true"
          >
            选择员工
          </el-button>
          <el-table v-if="memberList.length" :data="memberList" class="selected-table">
            <el-table-column prop="member_name" label="员工名称"></el-table-column>
            <el-table-column v-if="!readonly" label="操作" width="100">
              <template>
                <el-button type="text" @click="showMemberDialog = true">修改</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item label="管理范围" prop="area_ids">
          <el-button
            v-if="!selectedAreas.length && !readonly"
            type="primary"
            size="small"
            :loading="areaTreeLoading"
            @click="handleOpenAreaDialog"
          >
            选择管理范围
          </el-button>
          <el-table v-if="selectedAreas.length" :data="selectedAreas" max-height="360">
            <el-table-column label="管理范围" min-width="260">
              <template slot-scope="{ row }">
                {{ getAreaName(row) }}
              </template>
            </el-table-column>
            <el-table-column label="类型" width="100">
              <template slot-scope="{ row }">{{ getAreaTypeName(row) }}</template>
            </el-table-column>
            <el-table-column v-if="!readonly" label="操作" width="100">
              <template slot-scope="{ row }">
                <el-button type="text" @click="removeArea(row)">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button
            v-if="selectedAreas.length && !readonly"
            class="change-area"
            type="text"
            :loading="areaTreeLoading"
            @click="handleOpenAreaDialog"
          >
            修改管理范围
          </el-button>
        </el-form-item>
      </el-form>
      <select-member
        v-if="!readonly && showMemberDialog"
        :limit="1"
        :hasSelectArr="selectedMembers"
        :filterMemberMethod="filterMembers"
        dialogTitle="选择员工"
        @close="showMemberDialog = false"
        @submit="handleMemberConfirm"
      ></select-member>
      <area-tree-dialog
        v-model="showAreaDialog"
        :selected-ids="formData.area_ids"
        :tree-list="areaTreeList"
        @confirm="handleAreaConfirm"
      ></area-tree-dialog>
    </template>
    <template slot="footer">
      <el-button v-if="readonly" type="primary" @click="handleBack">返回</el-button>
      <template v-else>
        <el-button @click="handleBack">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </template>
  </layout-form>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import LayoutForm from "@/components/public/layout-form/index.vue";
import SelectMember from "@/components/business/select-member-new/index.vue";
import AreaTreeDialog from "./area-tree-dialog.vue";
import Api from "@/api/enterprise/hong-kong-business";

@Component({
  components: { LayoutForm, SelectMember, AreaTreeDialog }
})
export default class PermissionForm extends Vue {
  @Prop({ default: "add" }) mode!: "add" | "edit" | "detail";

  positionOptions = [
    { label: "GM(General Manager)", level: 1 }
  ];
  formData: any = {
    id: 0,
    position: "",
    position_level: 0,
    member_id: null,
    member_name: "",
    area_ids: []
  };
  selectedMembers: any[] = [];
  areaTreeList: any[] = [];
  configuredMemberIds: number[] = [];
  showMemberDialog = false;
  showAreaDialog = false;
  areaTreeLoading = false;
  saving = false;
  saved = false;
  initialValue = "";
  rules = {
    position: [{ required: true, message: "请选择职位", trigger: "change" }],
    member_id: [{ required: true, message: "请选择员工", trigger: "change" }],
    area_ids: [
      { type: "array", required: true, min: 1, message: "请选择管理范围", trigger: "change" }
    ]
  };

  get readonly() {
    return this.mode === "detail";
  }

  get memberList() {
    return this.formData.member_id
      ? [{ member_id: this.formData.member_id, member_name: this.formData.member_name }]
      : [];
  }

  get selectedAreas() {
    const selectedIds = new Set(this.formData.area_ids);
    const result: any[] = [];
    this.walkTree(this.areaTreeList, node => {
      if (selectedIds.has(node.id)) {
        result.push(node);
      }
    });
    return result;
  }

  async created() {
    await this.getConfiguredMemberIds();
    if (this.$route.query.id) {
      await this.getAreaTreeList();
      await this.getDetail();
    }
    this.setInitialValue();
  }

  setInitialValue() {
    this.initialValue = JSON.stringify({
      formData: this.formData
    });
  }

  isDirty() {
    if (this.readonly || this.saved) return false;
    const current = JSON.stringify({
      formData: this.formData
    });
    return current !== this.initialValue;
  }

  async confirmLeave() {
    if (!this.isDirty()) return true;
    try {
      await this.$confirm("检测到未保存的内容，确认离开当前页面吗？", "离开页面提示", {
        confirmButtonText: "确认离开",
        cancelButtonText: "继续编辑",
        type: "warning"
      });
      return true;
    } catch (error) {
      return false;
    }
  }

  async getAreaTreeList() {
    this.areaTreeLoading = true;
    try {
      const res: any = await Api.getAreaTree({
        parent_id: 0,
        has_children: true,
        has_virtual: false,
        ignore_close_store: true,
        language: ""
      });
      if (res.code !== 200) {
        this.$message.error(res.message || "获取管理范围失败");
        return false;
      }
      this.areaTreeList = res.data || [];
      return true;
    } finally {
      this.areaTreeLoading = false;
    }
  }

  async getConfiguredMemberIds() {
    const res: any = await Api.getConfiguredMemberIds();
    if (res.code === 200) {
      this.configuredMemberIds = res.data || [];
    }
  }

  async getDetail() {
    const res: any = await Api.getPermissionDetail({
      id: Number(this.$route.query.id)
    });
    if (res.code !== 200) {
      this.$message.error(res.message || "获取职位权限详情失败");
      return;
    }
    const data = res.data || {};
    const member = data.member_info || {};
    const option = this.positionOptions.find(item => item.label === data.position);
    this.formData = {
      id: data.id,
      position: data.position || "",
      position_level: option ? option.level : 0,
      member_id: member.id,
      member_name: member.name || "",
      area_ids: this.normalizeAreaIds(data.area_ids || [])
    };
    this.selectedMembers = member.id
      ? [{
          member_id: member.id,
          name: member.name,
          user_id: member.userid,
          mobile: member.mobile,
          mem_subscribe_status: 1
        }]
      : [];
  }

  filterMembers(list: any[]) {
    const currentId = this.formData.member_id;
    return list.filter(item => {
      return item.member_id === currentId ||
        !this.configuredMemberIds.includes(item.member_id);
    });
  }

  handleMemberConfirm(members: any[]) {
    const member = members[0];
    if (!member) return;
    this.formData.member_id = member.member_id;
    this.formData.member_name = member.name;
    this.selectedMembers = members;
    this.showMemberDialog = false;
  }

  handlePositionChange(value: string) {
    const option = this.positionOptions.find(item => item.label === value);
    this.formData.position_level = option ? option.level : 0;
  }

  async handleOpenAreaDialog() {
    const success = await this.getAreaTreeList();
    if (success) {
      this.showAreaDialog = true;
    }
  }

  handleAreaConfirm(areaIds: number[]) {
    this.formData.area_ids = this.normalizeAreaIds(areaIds);
    (this.$refs.form as any).validateField("area_ids");
  }

  removeArea(area: any) {
    this.formData.area_ids = this.formData.area_ids.filter(
      (id: number) => id !== area.id
    );
  }

  getAreaTypeName(area: any) {
    return area.type === 4 ? "门店" : "区域";
  }

  getAreaName(area: any) {
    return area.store_name || area.name || "--";
  }

  handleBack() {
    this.$router.push({
      name: "myapp_enterprise_hk-business-structure_hong-kong-business_permission-management"
    });
  }

  handleSave() {
    (this.$refs.form as any).validate(async (valid: boolean) => {
      if (!valid) return;
      this.saving = true;
      try {
        const res: any = await Api.savePermission({
          id: this.mode === "add" ? 0 : this.formData.id,
          member_id: this.formData.member_id,
          position: this.formData.position,
          position_level: this.formData.position_level,
          area_ids: this.normalizeAreaIds(this.formData.area_ids)
        });
        if (res.code !== 200) {
          this.$message.error(res.message || "保存失败");
          return;
        }
        this.saved = true;
        this.$message.success("保存成功");
        this.handleBack();
      } finally {
        this.saving = false;
      }
    });
  }

  normalizeAreaIds(ids: number[]) {
    const availableIds = new Set<number>();
    this.walkTree(this.areaTreeList, node => availableIds.add(node.id));
    const selectedIds = new Set(ids.filter(id => availableIds.has(id)));
    return Array.from(selectedIds).filter(id => {
      const node = this.findNode(id);
      let parentId = node ? node.parent_id : 0;
      while (parentId) {
        if (selectedIds.has(parentId)) return false;
        const parent = this.findNode(parentId);
        parentId = parent ? parent.parent_id : 0;
      }
      return true;
    });
  }

  walkTree(list: any[], callback: (node: any) => void) {
    list.forEach(node => {
      callback(node);
      if (node.children && node.children.length) {
        this.walkTree(node.children, callback);
      }
    });
  }

  findNode(id: number) {
    let result = null;
    this.walkTree(this.areaTreeList, node => {
      if (node.id === id) {
        result = node;
      }
    });
    return result as any;
  }
}
</script>

<style lang="scss" scoped>
.el-select {
  width: 400px;
}

.selected-table {
  width: 640px;
}

.change-area {
  margin-top: 10px;
}
</style>
