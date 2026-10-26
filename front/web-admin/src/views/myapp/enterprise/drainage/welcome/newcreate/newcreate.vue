<template>
  <personalise-welcome-template></personalise-welcome-template>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import personaliseWelcomeTemplate from "../components/personalise-welcome-template.vue";
@Component({
  components: {
    personaliseWelcomeTemplate
  }
})
export default class extends Vue {}
</script>
<style lang="scss" scoped></style>
