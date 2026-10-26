<template>
  <div class="select-item">
    <div class="input-box" @click="showModel = true">
      <span class="tag" v-if="list.length === 0">选择导购</span>
      <div v-else class="list">
        <el-tag class="tag" type="info">
          {{ list[0].name }}
        </el-tag>
        <el-tag class="tag" type="info" v-if="list.length > 1">
          +{{ list.length }}
        </el-tag>
      </div>
    </div>
    <SelectMember
      dialogTitle="选择导购"
      :hasSelectArr="list"
      @close="showModel = false"
      @submit="selectList"
      v-if="showModel"
    ></SelectMember>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import SelectMember from "@/components/business/select-member/index.vue";
@Component({
  components: {
    SelectMember
  }
})
export default class LayoutList extends Vue {
  private list = [];
  private showModel = false;
  private selectList(list) {
    this.list = list;
    this.showModel = false;
  }
  private clear() {
    this.list = [];
  }
}
</script>
<style lang="scss" scoped>
.select-item {
  display: inline-block;
}
.input-box {
  width: 250px;
  min-height: 36px;
  border: 1px solid #dcdfe6;
  background: white;
  cursor: pointer;
  padding: 3px 10px;
  box-sizing: border-box;
  color: rgb(214, 216, 220);
  .list {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
  }
}
.number {
  width: 15px;
  height: 15px;
  border-radius: 20px;
  background: rgba(0, 0, 0, 0.2);
  color: white;
  line-height: 15px;
  font-size: 12px;
  text-align: center;
}
.tag {
  margin-right: 10px;
  height: 28px;
  line-height: 28px;
}
</style>
