<template>
  <create-stratage-template></create-stratage-template>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import createStratageTemplate from "../../components/create-stratage-template.vue";

@Component({
  components: {
    createStratageTemplate
  }
})
export default class extends Vue {}
</script>
<style lang="scss"></style>
