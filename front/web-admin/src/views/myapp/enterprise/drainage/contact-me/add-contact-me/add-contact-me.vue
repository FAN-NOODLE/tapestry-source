<template>
  <edit-contact></edit-contact>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import editContact from "../../components/edit-contact.vue";
@Component({
  components: {
    editContact
  }
})
export default class extends Vue {}
</script>
<style lang="scss" scoped></style>
