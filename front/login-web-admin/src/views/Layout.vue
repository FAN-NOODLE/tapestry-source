<template lang="pug">
  el-container(style="height:100vh")
    el-header(style="background: #353943")
      com-header-menu
    router-view
</template>

<script>
// @ is an alias to /src
import ComHeaderMenu from "@/components/common/ComHeaderMenu";

export default {
  name: "layout",
  components: {
    ComHeaderMenu
  }
};
</script>
