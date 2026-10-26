<template>
  <permission-form ref="form" mode="add"></permission-form>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import PermissionForm from "../components/permission-form.vue";

@Component({ components: { PermissionForm } })
export default class PermissionAdd extends Vue {
  beforeRouteLeave(to: any, from: any, next: (allow?: boolean) => void) {
    (this.$refs.form as any).confirmLeave().then((allow: boolean) => next(allow));
  }
}
</script>
