<template>
  <div class="add-activity">
    <activityCompent :statusType="3" :activityId="activityId" />
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
export default class extends Vue {
  private activityId = 0;
  private created() {
    this.activityId = Number(this.$route.query.id);
  }
}
</script>
<style lang="scss" scoped></style>
