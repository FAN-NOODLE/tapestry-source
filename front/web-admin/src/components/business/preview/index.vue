<template>
  <div class="preview">
    <component
      :is="type"
      :index="index"
      :urls="urls"
      @close="close"
    ></component>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Prop } from "vue-property-decorator";
import PreviewImage from "./preview-image.vue";
@Component({
  components: {
    Images: PreviewImage
  }
})
export default class PreviewFile extends Vue {
  @Prop({ default: "Images" }) type?: string;
  @Prop() urls!: Array<string>;
  @Prop({ default: 0 }) index?: number;
  private close() {
    this.$emit("close");
  }
}
</script>
