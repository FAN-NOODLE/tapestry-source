<template>
  <div class="add-activity">
    <activityCompent :statusType="1" />
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import activityCompent from "../components/activityCompent.vue";
@Component({
  components: {
    activityCompent,
  },
})
export default class extends Vue {}
</script>
<style lang="scss" scoped></style>
