<template>
  <div class="copy-marketing-plan">
    <com-marketing-plan></com-marketing-plan>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import comMarketingPlan from "../components/com-marketing-plan.vue";
@Component({
  components: {
    comMarketingPlan
  }
})
export default class extends Vue {}
</script>

<style lang="scss" scoped></style>
