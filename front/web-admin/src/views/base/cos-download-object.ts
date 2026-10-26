interface CosDownloadCredential {
  bucket: string;
  private_bucket?: string;
  region: string;
}

export function resolveCosDownloadObject(
  fileUrl: string,
  credential: CosDownloadCredential,
  publicOrigins: string[] = []
): { bucket: string; key: string } {
  const decodedUrl = decodeURIComponent(fileUrl);
  const privateOrigin = credential.private_bucket
    ? `https://${credential.private_bucket}.cos.${credential.region}.myqcloud.com/`
    : "";

  if (privateOrigin && decodedUrl.startsWith(privateOrigin)) {
    return {
      bucket: credential.private_bucket as string,
      key: decodedUrl.slice(privateOrigin.length)
    };
  }

  const publicOrigin = `https://${credential.bucket}.cos.${credential.region}.myqcloud.com/`;
  let key = decodedUrl.startsWith(publicOrigin)
    ? decodedUrl.slice(publicOrigin.length)
    : decodedUrl;
  for (const origin of publicOrigins) {
    if (key.startsWith(origin)) {
      key = key.slice(origin.length);
      break;
    }
  }

  return { bucket: credential.bucket, key };
}
