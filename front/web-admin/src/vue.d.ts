import "vue";
declare module "vue/types/vue" {
  interface Vue {
    $md5: any;
    $axios: any;
    $echarts: any;
    $VUE_APP_API_ROOT: any;
    $VUE_APP_API_VERSION: any;
  }
}
