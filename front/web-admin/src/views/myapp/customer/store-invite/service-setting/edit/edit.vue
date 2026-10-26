<template>
  <div>
    <editCreate />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import editCreate from "../components/editCreate.vue";

@Component({
  components: {
    editCreate
  }
})
export default class extends Vue {}
</script>

<style lang="scss" scoped></style>
