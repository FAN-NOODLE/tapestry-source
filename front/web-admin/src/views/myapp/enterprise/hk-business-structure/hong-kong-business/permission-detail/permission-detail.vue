<template>
  <permission-form mode="detail"></permission-form>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import PermissionForm from "../components/permission-form.vue";

@Component({ components: { PermissionForm } })
export default class PermissionDetail extends Vue {}
</script>
