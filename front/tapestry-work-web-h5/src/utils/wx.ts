import { getConfigApi, getAgentConfigApi } from "@/api/common";
import Store from "@/store/index";
import { getDeviceInfo } from "./tools";
const store: any = Store;
/**
 * 需要注册的 API
 */
// prettier-ignore
const jsApiList = [
  "hideMenuItems", // 批量隐藏功能按钮
  "showMenuItems", // 批量显示功能按钮
  "hideAllNonBaseMenuItem", // 隐藏所有非基础按钮
  "showAllNonBaseMenuItem", // 显示所有功能按钮
  "showOptionMenu", // 显示菜单项
  "hideOptionMenu", // 隐藏菜单项
  "onMenuShareAppMessage", // 分享菜单 - app信息
  "onMenuShareWechat", // 分享菜单 - 微信
  "onMenuShareTimeline", // 分享菜单 - 时间线
  "chooseImage", // 选择图片
  "chooseVideo", // 选择图片
  "previewImage", // 预览图片
  "uploadImage", // 上传图片
  "downloadImage", // 下载图片
  "getLocalImgData", // 该接口获取到的是图片的数据
  "scanQRCode", // 扫描二维码
  "selectExternalContact", // 获取通讯录
  "openUserProfile", //打开个人信息页接口
  "shareToExternalContact", // 此接口支持企业成员把小程序，传递到群发助手进行发送
  "sendChatMessage", // 通过聊天工具栏向当前会话发送消息，支持多种消息格式
  "getCurExternalContact", // 获取当前外部联系人userid
  "getCurExternalChat", // 获取当前客户群的群ID
  "getContext", // 获取当前客户群的群ID
  "openEnterpriseChat", // 打开会话
  "onHistoryBack", // 监听返回上一页
  "launchMiniprogram", // 打开小程序
  'setClipboardData', // 设置剪贴板内容
  'shareToExternalMoments', // 发送朋友圈
  'makePhoneCall', // 打电话
  "onNetworkStatusChange", // 监听网络状态变化
  "getNetworkType", // 获取网络状态
  "previewFile", // 预览文件
  "getClipboardData", // 粘贴板
  "getLocation",
]

/**
 * 根据code获取用户信息  wx.config 进行权限验证配置, 注入的是企业的身份与权限
 * @param {*} url
 * @param {*} path
 */
const wxConfig = async (url = window.location.href) => {
  console.log(url);
  
  const promise1 = getConfigApi({
    url: url,
    hash: store.state.common.ENTERPRISE_HASH,
  });
  const promise2 = getAgentConfigApi({
    url: url,
    hash: store.state.common.ENTERPRISE_HASH,
  });
  return new Promise((resolve: Function, reject: Function) => {
    Promise.all([promise1, promise2]).then((res) => {
      const res1 = res[0];
      const res2 = res[1];
      const { appid, timestamp, noncestr, signature } = res1.data;
      window.wx.config({
        beta: true,
        debug: false,
        appId: appid,
        timestamp: timestamp,
        nonceStr: noncestr,
        signature: signature,
        jsApiList: jsApiList, // 获取当前外部联系人userid
        success: (res: any) => {
          const str = JSON.stringify(res);
          console.log("config调用成功" + str, jsApiList);
        },
        fail: (err: any) => {
          const error = JSON.stringify(err);
          console.log("config调用失败" + error);
        },
      });
      window.wx.ready(async () => {
        console.log("调用微信sdk成功", jsApiList);
        const { corp_id, agent_id, timestamp, noncestr, signature } = res2.data;
        if(window.wx.agentConfig) {
          window.wx.agentConfig({
            corpid: corp_id, // 必填，企业微信的corpid，必须与当前登录的企业一致
            agentid: agent_id, // 必填，企业微信的应用id
            timestamp: timestamp, // 必填，生成签名的时间戳
            nonceStr: noncestr, // 必填，生成签名的随机串
            signature: signature, // 必填，签名，见附录1
            jsApiList: jsApiList, //必填
            success: (res: any) => {
              const str = JSON.stringify(res);
              resolve();
              console.log("agentConfig调用成功" + str, jsApiList);
            },
            fail: (err: any) => {
              const error = JSON.stringify(err);
              reject();
              console.log("agentConfig调用失败" + error);
            },
          });
        } else {
          resolve()
        }
      });
    });
  });
};

const wxAgentConfig = async () => {
  return new Promise((resolve: Function, reject: Function) => {
    getAgentConfigApi({
      url: window.location.href,
      hash: store.state.common.ENTERPRISE_HASH,
    }).then((res) => {
      const { corp_id, agent_id, timestamp, noncestr, signature } = res.data;
      window.wx.agentConfig && window.wx.agentConfig({
        corpid: corp_id, // 必填，企业微信的corpid，必须与当前登录的企业一致
        agentid: agent_id, // 必填，企业微信的应用id
        timestamp: timestamp, // 必填，生成签名的时间戳
        nonceStr: noncestr, // 必填，生成签名的随机串
        signature: signature, // 必填，签名，见附录1
        jsApiList: jsApiList, //必填
        success: (res: any) => {
          const str = JSON.stringify(res);
          console.log("agentConfig调用成功" + str);
          resolve();
        },
        fail: (err: any) => {
          const error = JSON.stringify(err);
          console.log("agentConfig调用失败" + error);
          reject();
        },
      });
    });
  });
};

function getContext() {
  return new Promise((resolve, reject) => {
    window.wx.invoke("getContext", {}, function (res) {
      if (res.err_msg == "getContext:ok") {
        resolve(res)
      } else {
        //错误处理
        console.log("getContext:error" + JSON.stringify(res));
        reject(res)
      }
    });
  })
}

function getCurExternalContact() {
  return new Promise((resolve, reject) => {
    window.wx.invoke("getCurExternalContact", {}, (res) => {
      if (res.err_msg == "getCurExternalContact:ok") {
        Store.commit("user/setExternalId", res.userId);
        
        resolve(res)
      } else {
        //错误处理
        console.log("getCurExternalContact:error" + JSON.stringify(res));
        reject(res)
      }
    });
  })
}

function getCurExternalChat() {
  return new Promise((resolve, reject) => {
    // 从群聊会话的工具栏进入-打开客户画像-社群页面
    window.wx.invoke("getCurExternalChat", {}, function (res) {
      if (res.err_msg == "getCurExternalChat:ok") {
        Store.commit("common/setGroupChatId", res.chatId);
        resolve(res)
      } else {
        //错误处理
        console.log("getCurExternalChat:error" + JSON.stringify(res));
        reject(res)
      }
    });
  })
}

async function openEnterpriseChat(message) {
  await wxConfig()
  await wxAgentConfig()
  return new Promise((resolve, reject) => {
    // 调用企微接口
    window.wx.invoke("openEnterpriseChat", {
      userIds: '', //参与会话的企业成员列表，格式为userid1;userid2;...，用分号隔开。
      externalUserIds: '', // 参与会话的外部联系人列表，格式为userId1;userId2;…，用分号隔开。
      groupName: "", // 会话名称。单聊时该参数传入空字符串""即可。
      chatId: "", // 若要打开已有会话，需指定此参数。如果是新建会话，chatId必须为空串
      ...message
    }, (res) => {
      console.log(res);
      if (res.err_msg == "openEnterpriseChat:ok") {
        resolve(res)
      } else {
        reject(res)
      }
    });

  })
}
async function chooseVideo(params = {}) {
  await wxConfig()
  await wxAgentConfig()
  return new Promise((resolve, reject) => {
    window.wx.invoke('chooseVideo', {
      sourceType: ['album', 'camera'], // 同时允许相册和相机
      maxDuration: 60, // 视频最大时长（秒）
      ...params,
      success: function(res) {
          resolve(res);
      },
      fail: function(err) {
          reject(err);
      }
    });

  })
}

async function chooseImage(params = {}) {
  await wxConfig()
  await wxAgentConfig()
  return new Promise((resolve, reject) => {
    window.wx.invoke('chooseImage', {
      count: 1,
      sizeType: ['original', 'compressed'],
      sourceType: ['album', 'camera'],
      defaultCameraMode: 'batch',
      isSaveToAlbum: true,
      ...params,
      success: function(res) {
          resolve(res);
      },
      fail: function(err) {
          reject(err);
      }
    });

  })
}

async function getLocation(url = '') {
  // 获取当前系统型号区分安卓、苹果、华为
  const  {platform} = getDeviceInfo()
  console.log(platform);
  
  // if(platform === 'android') {
  // }
  // await wxConfig(url)
  return new Promise((resolve, reject) => {
    (window.wx as any).getLocation({
      type: 'wgs84', // 默认为wgs84的gps坐标，如果要返回直接给openLocation用的火星坐标，可传入'gcj02'
      success: function (res) {
        console.log('getLocation success', res);
        resolve(res);
      },
      fail: function (err) {
        console.log('getLocation fail', err);
        reject(err);
      },
      complete: function (res) {
        console.log('getLocation complete', res);
        if (res.errMsg === 'getLocation:ok') {
          resolve(res);
        } else {
          reject(res);
        }
      }
    });
  })
}

async function shareToExternalMoments() {
  await wxConfig()
  return new Promise((resolve, reject) => {
    (window.wx as any).invoke("shareToExternalMoments", {
      text: {
        content: '企业微信'
      },
    }, (res) => {
      console.log(res);
      resolve(res);
    });
  })
}

export default { wxConfig, wxAgentConfig, getContext, getCurExternalContact, getCurExternalChat, openEnterpriseChat, chooseVideo, chooseImage, getLocation, shareToExternalMoments };
