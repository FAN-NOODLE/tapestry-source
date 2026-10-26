<template>
  <customer-welcome-form></customer-welcome-form>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import customerWelcomeForm from "@/views/myapp/enterprise/components/customer-welcome-form.vue";

@Component({
  components: {
    customerWelcomeForm
  }
})
export default class extends Vue {}
</script>
<style lang="scss"></style>
