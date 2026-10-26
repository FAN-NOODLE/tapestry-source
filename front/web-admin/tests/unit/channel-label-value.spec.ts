import { buildChannelLabelValue } from "@/views/myapp/enterprise/components/channel-label-value";

describe("buildChannelLabelValue", () => {
  it("keeps WeCom group and tag ids when channel detail is loaded for editing", () => {
    expect(
      buildChannelLabelValue({
        label_group_id: 101,
        label_id: 202,
        wx_tag_group_id: "etg_group",
        wx_tag_id: "et_tag"
      })
    ).toBe("101,202,etg_group,et_tag");
  });
});
