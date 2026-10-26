

export function getImageData(url: string, type: 'base64'|'blob' = 'blob'): Promise<string|Blob> {
  return new Promise((resolve) => {
    const canvas = document.createElement('canvas');
    const img = new Image();
    img.crossOrigin = 'anonymous';
    img.src = url + (url.indexOf('?') > -1 ? '&' : '?') + 'timestamp=' + new Date().getTime();
    img.onload = function () {
      canvas.width = img.width;
      canvas.height = img.height;
      canvas.getContext('2d')?.drawImage(img, 0, 0);
      if(type === 'base64') {
        resolve(canvas.toDataURL());
      } else {
        canvas.toBlob((blob) => {
          resolve(blob);
        })
      }
    };
  })
}
