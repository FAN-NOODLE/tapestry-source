import { VueConstructor } from "vue";
import FormatDate from "./FormatDate";
import FormatMoney from "./FormatMoney";
import FormatMoneyByDot from "./FormatMoneyByDot";
import FormatYuanMoney from "./FormatYuanMoney";
import FormatText from "./FormatText";
import FormatFileSize from "./FormatFileSize";
import FormatNumber from "./FormatNumber";
const filters = [FormatDate, FormatMoney, FormatMoneyByDot, FormatYuanMoney, FormatText, FormatFileSize, FormatNumber];

export default {
  install: (Vue: VueConstructor) => {
    if (filters.length && filters.length > 0) {
      filters.forEach((item) => {
        Vue.filter(item.key, item.func);
      });
    }
  },
};
