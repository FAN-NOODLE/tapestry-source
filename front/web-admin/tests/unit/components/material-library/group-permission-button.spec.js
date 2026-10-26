const fs = require("fs");
const path = require("path");

const editGroup = fs.readFileSync(
  path.resolve(
    __dirname,
    "../../../../src/views/myapp/enterprise/material-library/components/edit-group.vue"
  ),
  "utf8"
);

describe("素材分组批量权限按钮", () => {
  it("始终展示，并在未选择分组时禁用", () => {
    expect(editGroup).toContain(':disabled="!chooseCateList.length"');
    expect(editGroup).toContain("批量设置分组权限");
    expect(editGroup).not.toContain('v-if="chooseCateList.length"');
  });
});
