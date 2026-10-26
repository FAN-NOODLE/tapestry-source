import { shouldRequestSourceGroupMove } from "@/views/myapp/enterprise/material-library/material-list/list/source-group-drag";

describe("shouldRequestSourceGroupMove", () => {
  it("requests the move when a valid drop changes sibling order", () => {
    expect(
      shouldRequestSourceGroupMove(
        { data: { id: 2 } },
        "before",
        [1, 2, 3],
        [2, 1, 3]
      )
    ).toBe(true);
  });

  it("skips the request when sibling order is unchanged", () => {
    expect(
      shouldRequestSourceGroupMove(
        { data: { id: 2 } },
        "before",
        [1, 2, 3],
        [1, 2, 3]
      )
    ).toBe(false);
  });

  it("skips the request when Element UI reports no valid drop", () => {
    expect(
      shouldRequestSourceGroupMove(
        { data: { id: 1 } },
        "none",
        [1, 2, 3],
        [2, 1, 3]
      )
    ).toBe(false);
  });
});
