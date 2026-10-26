<!-- 新增积分规则弹窗 -->
<template>
  <el-dialog :visible.sync="syncVisible" width="1000px" title="">
    <div class="model-container">
      <div class="box" v-for="(model, key) in list" :key="key">
        <div class="box-icon">
          <img :src="model.icon" />
        </div>
        <div class="box-content">
          <div>{{ model.name }}</div>
          <div>{{ model.describe }}</div>
          <div class="buttons">
            <el-button
              type="primary"
              size="mini"
              @click="goAdd(model.id, model.listId)"
              >{{ model.listId ? "编辑" : "新增" }}积分规则</el-button
            >
            <el-button
              type="primary"
              size="mini"
              @click="
                emitAtlas(model.id, model.relationStatus, model.relationId)
              "
              >{{
                model.relationStatus === 1 ? "编辑" : "关联"
              }}能力图谱</el-button
            >
          </div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>
<script lang="ts">
import { Vue, Component, PropSync, Prop } from "vue-property-decorator";
import { integralType } from "../../config";

@Component({})
export default class AddDialog extends Vue {
  /** 判断父属性显示影藏，双向绑定 */
  @PropSync("visible", { type: Boolean }) syncVisible!: boolean;

  /** 渲染list - 维度 */
  @Prop() list!: Array<integralType>;

  /** 跳转到新增能力图谱 */
  private goAdd(id: number, relationId: number) {
    // relationId 判断是否有关联列表id, 如果有，则是编辑
    const name = relationId
      ? "myapp_wisdom_we-journey_integral-set-edit"
      : "myapp_wisdom_we-journey_integral-set-add";
    this.$router.push({
      name,
      query: {
        id: relationId + "",
        type: id + ""
      }
    });
  }

  /** 关联能力图谱 - 触发父元素emit */
  private emitAtlas(modelId: number, status?: number, id?: number) {
    this.syncVisible = false;
    this.$emit("relationAtlas", modelId, status, id);
  }
}
</script>
<style lang="scss" scoped>
.model-container {
  display: grid;
  margin-bottom: 15px;
  grid-template-columns: repeat(3, 1fr);
  grid-gap: 30px;
  .box {
    border: 1px solid #000000;
    height: 132px;
    padding: 20px 16px;
    cursor: pointer;
    display: flex;
    align-items: center;
    box-sizing: border-box;
    &-icon {
      width: 47px;
      height: 47px;
      margin-right: 5px;
      align-self: center;
      img {
        width: 100%;
        height: 100%;
        border-radius: 50%;
      }
    }
    &-content {
      flex: 1;
      font-size: 12px;
      color: #606266;
      div:first-child {
        font-weight: bold;
        color: #212121;
        font-size: 16px;
        margin-bottom: 10px;
      }
      .buttons {
        margin-top: 10px;
        display: flex;
        button {
          padding: 5px 10px !important;
        }
        button:last-child {
          margin-left: 10px !important;
        }
      }
    }
  }
}
</style>
