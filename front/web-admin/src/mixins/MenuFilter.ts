import { Vue, Component } from "vue-property-decorator";
@Component
export default class MenuFilter extends Vue {
  public _menuFilter(keys: any, menu: Array<any>) {
    if (!menu || menu.length === 0) {
      return;
    }
    for (let i = 0, len = menu.length; i < len; i++) {
      if (!menu[i]) {
        return;
      }
      if (!keys.has(menu[i].key)) {
        this.$delete(menu, i);
        i--;
      } else {
        this._menuFilter(keys, menu[i].children);
      }
    }
  }
  public menuButton(keys: any, menu: Array<any>) {
    if (!menu || menu.length === 0) {
      return;
    }
    for (let i = 0, len = menu.length; i < len; i++) {
      if (!menu[i]) {
        return;
      }
      if (!keys.has(menu[i].key)) {
        this.$delete(menu, i);
        i--;
      } else {
        this._menuFilter(keys, menu[i].buttons);
      }
    }
  }
}
