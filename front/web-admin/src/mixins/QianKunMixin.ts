import { Vue, Component } from "vue-property-decorator";
@Component
export default class QianKunMixin extends Vue {
  // 是否是在非微前端运行
  get _isOpenQIANKUN() {
    return !window.__POWERED_BY_QIANKUN__;
  }
}
