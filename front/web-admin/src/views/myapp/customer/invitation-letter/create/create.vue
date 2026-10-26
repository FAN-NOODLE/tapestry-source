<template>
  <create-invitation-letter type="create"></create-invitation-letter>
</template>
<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
import createInvitationLetter from "../components/create-invitation-letter.vue";

@Component({
  components: {
    createInvitationLetter,
  },
})
export default class extends Vue {}
</script>
<style lang="scss"></style>
