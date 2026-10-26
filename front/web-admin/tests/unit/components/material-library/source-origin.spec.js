const fs = require("fs");
const path = require("path");

const materialList = fs.readFileSync(
  path.resolve(
    __dirname,
    "../../../../src/views/myapp/enterprise/material-library/material-list/list/list.vue"
  ),
  "utf8"
);
const linkMaterialItem = fs.readFileSync(
  path.resolve(
    __dirname,
    "../../../../src/views/myapp/enterprise/material-library/components/link-material-item.vue"
  ),
  "utf8"
);

describe("链接素材来源", () => {
  it("搜索区提供服务端来源筛选", () => {
    expect(materialList).toContain('text="素材来源："');
    expect(materialList).toContain("search.source_origin");
    expect(materialList).toContain("filterObj.source_origin = this.search.source_origin");
  });

  it("同步素材卡片展示 CH同步", () => {
    expect(linkMaterialItem).toContain('v-if="item.is_kos_sync === 1"');
    expect(linkMaterialItem).toContain("CH同步");
  });
});
