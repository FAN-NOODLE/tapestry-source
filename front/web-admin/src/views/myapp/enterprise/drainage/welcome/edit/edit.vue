<template>
  <customer-welcome-template></customer-welcome-template>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import customerWelcomeTemplate from "../../../components/customer-welcome-template.vue";

@Component({
  components: {
    customerWelcomeTemplate
  }
})
export default class extends Vue {}
</script>
<style lang="scss"></style>
