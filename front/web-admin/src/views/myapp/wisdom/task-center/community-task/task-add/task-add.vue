<template>
  <div class="task-center-add">
    <com-task-center></com-task-center>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import comTaskCenter from "../components/com-task-center.vue";
@Component({
  components: {
    comTaskCenter
  }
})
export default class extends Vue {}
</script>

<style lang="scss" scoped></style>
