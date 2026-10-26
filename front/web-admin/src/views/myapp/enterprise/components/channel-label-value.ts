interface ChannelLabelValueSource {
  label_group_id: number;
  label_id: number;
  wx_tag_group_id?: string;
  wx_tag_id?: string;
}

export function buildChannelLabelValue(label: ChannelLabelValueSource): string {
  return [
    label.label_group_id,
    label.label_id,
    label.wx_tag_group_id || "",
    label.wx_tag_id || ""
  ].join(",");
}
