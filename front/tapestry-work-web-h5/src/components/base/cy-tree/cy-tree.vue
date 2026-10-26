<template>
  <div class="tree-menu">
    <ul v-for="(item, index) in options" :key="index">
      <cy-tree-item :data="item"></cy-tree-item>
    </ul>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop, Watch, Emit } from "vue-property-decorator";
import cyTreeItem from "./cy-tree-item.vue";
@Component({
  name: "cy-tree",
  components: { cyTreeItem },
})
export default class extends Vue {
  @Prop() options: Array<any>;
}
</script>
<style lang="scss" scoped>
.tree-menu {
  width: 100%;
  padding-right: 10px;
}
</style>
