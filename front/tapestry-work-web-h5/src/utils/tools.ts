/**
 * @description: 获取视频海报
 */
export function getVideoPoster(videoElement) {
  const video = videoElement
  let resolve = null
  let reject = null
  const promise = new Promise((res, rej) => {
    resolve = res
    reject = rej
  })

  // 检查是否为 iOS 设备
  const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent) && !window.MSStream

  // 根据设备类型选择合适的事件侦听器
  if (isIOS) {
    video.addEventListener("loadeddata", requestThumb)
  } else {
    video.addEventListener("canplay", requestThumb)
  }

  async function requestThumb() {
    // 检查视频是否在恰好播放
    if (video.ended) {
      removeEventListeners()
      return // 如果已暂停或结束，则不做处理
    }

    // 从事件侦听器中移除，避免发生意外调用
    removeEventListeners()

    // 创建画布，将视频的第一帧图像绘制上去
    const canvas = document.createElement("canvas")
    canvas.width = video.videoWidth
    canvas.height = video.videoHeight
    canvas.getContext("2d").drawImage(video, 0, 0, canvas.width, canvas.height)

    // 将画布上的图像数据转换为 base64 格式
    const imageData = canvas.toDataURL()

    resolve({
      base64: imageData,
      file: await base64ToImage(imageData)
    })
  }

  function removeEventListeners() {
    video.removeEventListener("canplay", requestThumb)
    video.removeEventListener("loadeddata", requestThumb)
  }
  video.currentTime = 1 // 视频快进到第1秒
  video.crossOrigin = "Anonymous" // 设置 crossOrigin 属性
  // 静音
  video.muted = true;

  video.load() // 加载视频，

  // 如果视频已经在缓存中，则可以直接调用 requestThumb() 方法
  if (video.readyState >= 2) {
    requestThumb()
  }

  return promise
}

/**
 * @description: base64转图片
 * @param base64string
 * @param filename
 * @returns
 */
export function base64ToImage(base64string, filename = Date.now().toString()) {
  // 获取文件类型
  const match = base64string.match(/data:([a-zA-Z]+\/[a-zA-Z]+)/)
  const type = match ? match[1] : ""
  const ext = type.split("/")[1]
  // 获取 base64 编码的文件内容
  const content = atob(base64string.split(";base64,")[1])
  // 将内容转换为一个 ArrayBuffer
  const buffer = new ArrayBuffer(content.length)
  const view = new Uint8Array(buffer)
  for (let i = 0; i < content.length; i++) {
    view[i] = content.charCodeAt(i)
  }
  // 创建一个 Blob 对象，并将其保存到文件
  const blob = new Blob([buffer], { type })
  const file = new File([blob], filename + "." + ext, { type })
  return new Promise((resolve, reject) => {
    if (file) {
      resolve(file)
    } else {
      reject(Error("Failed to create file."))
    }
  })
}


export async function bufferToImage(buffer: ArrayBuffer, filename: string) {}


/**
 * 返回示例
 * {
 *   platform: 'ios' | 'android' | 'harmony',
 *   brand:    'huawei' | 'xiaomi' | 'oppo' | 'vivo' | 'samsung' | 'apple' | 'unknown'
 * }
 */
export function getDeviceInfo() {
  console.log('navigator', navigator.userAgent);
  
  const ua = navigator.userAgent.toLowerCase();

  // 1. 先判断鸿蒙
  const isHarmony = /harmony/i.test(ua);

  // 2. 再判断 iOS
  const isIOS = /iphone|ipad|ipod/i.test(ua);

  // 3. 剩下就是 Android
  const isAndroid = !isHarmony && !isIOS;

  // 4. 品牌细化
  let brand = 'unknown';
  if (isIOS) {
    brand = 'apple';
  } else {
    // 华为
    if (/huawei|honor|hms|hisuite/i.test(ua)) brand = 'huawei';
    // 小米
    else if (/mi\s|xiaomi|redmi/i.test(ua)) brand = 'xiaomi';
    // OPPO / Realme / OnePlus
    else if (/oppo|realme|oneplus/i.test(ua)) brand = 'oppo';
    // vivo / iQOO
    else if (/vivo|iqoo/i.test(ua)) brand = 'vivo';
    // 三星
    else if (/samsung|sm-/i.test(ua)) brand = 'samsung';
  }

  return {
    platform: isHarmony ? 'harmony' : isIOS ? 'ios' : 'android',
    brand
  };
}

// ===== 用法 =====
const { platform, brand } = getDeviceInfo();
console.log(platform, brand);   // android huawei