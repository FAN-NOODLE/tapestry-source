<template>
  <PermissionCreateForm :type="1"></PermissionCreateForm>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import PermissionCreateForm from "../components/create-form.vue";
@Component({
  components: {
    PermissionCreateForm
  }
})
export default class PermissionAdd extends Vue {
}
</script>