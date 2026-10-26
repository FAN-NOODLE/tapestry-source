<template>
  <new-message></new-message>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import newMessage from "../../components/new-message.vue";

@Component({
  components: {
    newMessage
  }
})
export default class extends Vue {}
</script>
<style lang="scss"></style>
