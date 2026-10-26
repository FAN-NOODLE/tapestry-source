import { resolveCosDownloadObject } from "@/views/base/cos-download-object";

describe("resolveCosDownloadObject", () => {
  it("uses the private bucket returned by the credential for a private COS URL", () => {
    const got = resolveCosDownloadObject(
      "https://tapestry-uat-private-1308608494.cos.ap-shanghai.myqcloud.com/private/attachment/report.pdf",
      {
        bucket: "tapestry-uat-1308608494",
        private_bucket: "tapestry-uat-private-1308608494",
        region: "ap-shanghai"
      },
      ["https://cdn-cos-uat.coach.com.cn/"]
    );

    expect(got).toEqual({
      bucket: "tapestry-uat-private-1308608494",
      key: "private/attachment/report.pdf"
    });
  });
});
