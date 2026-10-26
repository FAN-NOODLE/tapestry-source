<template>
  <childTask />
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import childTask from "../components/child-task.vue";
@Component({
  components: {
    childTask
  }
})
export default class extends Vue {}
</script>
<style lang="scss" scoped></style>
