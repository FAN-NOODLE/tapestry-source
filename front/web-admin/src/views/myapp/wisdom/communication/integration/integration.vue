<template>
  <keep-alive include="IntegrationList,IntegrationAdd">
    <router-view></router-view>
  </keep-alive>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
@Component({
  name: 'Integration'
})
export default class Integration extends Vue {

}
</script>