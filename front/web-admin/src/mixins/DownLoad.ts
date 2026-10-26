import { Vue, Component } from "vue-property-decorator";
@Component
export default class DownLoad extends Vue {
  /** 文件导出，按照bolb 对象导出 */
  public downLoadFile(res: any, name: string) {
    const blob = new Blob([res], {
      type:
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf-8"
    });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = name;
    a.click();
    window.URL.revokeObjectURL(url);
    // this.$message.success("导出成功");
  }
}
