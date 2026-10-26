import { Toast } from "vant";
import "vue";
import SensorsTrack from "./utils/sensors/track";
declare module "vue/types/vue" {
  interface Vue {
    $helper: helper;
    $toast: Toast;
    $track: SensorsTrack;
  }
}

declare module 'TMap' {
  export default class TMap {
  }
}