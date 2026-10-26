import { Loading } from "element-ui";
import { Vue, Component } from "vue-property-decorator";
interface IValidate {
  type: string;
  limit: number;
  suggest?: number;
}
interface ILimit {
  width?: number;
  height?: number;
}
const file_type: any = {
  png: ["image/png", "图片"],
  jpg: ["image/jpeg", "图片"],
  gif: ["image/gif", "图片"],
  xls: ["", "文件"],
  csv: ["", "文件"]
};
@Component
export default class UploadValidate extends Vue {
  public _uploadValidate(file: any, validate: Array<IValidate>, lim?: ILimit) {
    let target = null;
    // 选中验证规则
    for (let i = 0, len = validate.length; i < len; i++) {
      if (file_type[validate[i].type][0] === file.type) {
        target = validate[i];
        break;
      }
    }
    if (!target) {
      this.$message.error(
        `文件类型错误，请上传 ${validate
          .map(item => item.type)
          .join(" 或 ")} 类型的文件！`
      );
      return false;
    }
    if (target.limit) {
      if (parseInt(`${file.size / 1024}`) > target.limit) {
        if (target.type === "gif") {
          this.$message.error("为保证小程序流畅度，GIF图片建议不超过1.5M");
          return false;
        }
        this.$message.error(
          `${file_type[target.type][1]}超出 ${
            target.limit
          }kb 上传限制，请重新上传！`
        );
        return false;
      }
    }
    if (target.suggest) {
      if (parseInt(`${file.size / 1024}`) > target.suggest) {
        if (target.type === "gif") {
          this.$message.warning("为保证小程序流畅度，GIF图片建议不超过1.5M");
        } else {
          this.$message.warning(
            `建议${file_type[target.type][1]}上传大小为 ${target.suggest}kb！`
          );
        }
      }
    }
    const obj = { width: 0, height: 0 };
    const limit = lim || obj;
    const width = limit.width || 0;
    const height = limit.height || 0;
    const is_size = new Promise((resolve, reject) => {
      let _validate = true;
      let msg = "";
      const img = new Image();
      img.src = window.URL.createObjectURL(file);
      img.onload = () => {
        if (width === 0 && height === 0) {
          _validate = true;
        } else if (width === 0) {
          _validate = img.height === height;
          msg = `请上传高为${height}px的图片`;
        } else if (height === 0) {
          _validate = img.width === width;
          msg = `请上传宽为${width}px的图片`;
        } else {
          _validate = img.width === width && img.height === height;
          msg = `请上传${width}*${height}大小的图片`;
        }
        if (_validate) {
          Loading.service({
            fullscreen: true,
            text: "图片上传中，请稍后"
          });
          resolve(file);
        } else {
          this.$message.error(msg);
          reject();
        }
      };
    });
    return is_size;
  }
}
