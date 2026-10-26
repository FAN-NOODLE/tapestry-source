<template>
  <div><createc-com type="edit"></createc-com></div>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import createcCom from "../components/create-com.vue";
@Component({
  components: { createcCom }
})
export default class extends Vue {}
</script>
<style lang="scss" scoped></style>
