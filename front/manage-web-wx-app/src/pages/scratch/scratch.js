import Scratch from "./scratch-canvas.js";
import {
  getWxlogin,
  getScratchActivity,
  scratchPrize,
  couponList,
  checkCoupon,
  shareInfo
} from "@/api/scratch-api";
const { enterpriseid } = require("@/lib/config");
const btnTxts = ["线下核销", "已核销", "已过期"];
Page({
  // 页面的初始数据
  data: {
    firstTouch: 0,
    isOver: 0,
    top: 0,
    list: [],
    activityId: "1",
    isPrize: false,
    prizeInfo: {}, //获奖信息
    times: -1, // 抽奖次数
    btnTxt: "", // 核销文案
    couponObj: {}, // 核销信息
    writeOffModal: false,
    pageDetails: {}, // 页面详情
    can_scratch: true, // 判断能否参与刮奖
    pageHeight: 900,
    writeOffstatus: -1, //核销状态
    loading: false,
    prizeType: "", // 抽奖类型
    source: "", // 场景值
  },
  // 生命周期函数--监听页面加载
  async onLoad(options) {
    const { activityId, source } = options;
    const res = wx.getMenuButtonBoundingClientRect();
    this.setData({
      top: res.top + 8,
      activityId,
      loading: true,
      source,
    });
  },
  async onShow() {
    await this.getWxCode();
  },
  onHide() {
    clearInterval(this.timer);
  },
  /**
   * 获取code
   */
  getWxCode() {
    wx.login({
      success: (res) => {
        this.getWxlogin(res.code);
      },
      fail(res) {
        console.log("errpr", res);
      },
    });
  },
  /**
   * 获取token
   */
  async getWxlogin(code) {
    try {
      const params = {
        enterprise_id: enterpriseid,
        code,
      };
      const res = await getWxlogin(params);
      if (res.code === 200) {
        wx.setStorageSync("customer_token", res.data.token);
        this.getActivity();
        this.getStartPrize();
      }
    } catch (error) {}
  },
  /**
   * 获取详情
   */
  async getActivity() {
    try {
      const { activityId, source } = this.data;
      const params = {
        activity_id: activityId,
        source: source ? source : "qr_share",
      };
      const { data, code } = await getScratchActivity(params);
      if (code === 200) {
        data.activity_content = data.activity_content.split("；");
        this.setData(
          {
            pageDetails: data,
            can_scratch: data.can_scratch,
          },
          () => {
            this.draw();
          }
        );
        this.getCouponList();
      }
    } catch (error) {}
  },
  /**
   * 获取券列表
   */
  async getCouponList() {
    try {
      const params = {
        page: 1,
        size: 100,
      };
      const {
        data: { check_id, data, status, meta },
      } = await couponList(params);
      let couponObj = {};
      if (status === 0) {
        couponObj = data.find((item) => item.id === check_id);
        couponObj.qr_code = "data:image/png;base64," + couponObj.qr_code;
      }
      this.getPageHeight(meta.total);
      this.setData({
        list: data,
        couponObj,
        btnTxt: btnTxts[status],
        writeOffstatus: status,
        loading: false,
      });
    } catch (error) {
      this.setData({
        loading: false,
      });
    }
  },
  // 开始抽奖
  getStartPrize() {
    const params = {
      activity_id: this.data.activityId,
    };
    setTimeout(() => {
      this.setData({
        prizeType: "start",
      });
    }, 1500);

    this.getPrize(params, "start");
  },
  //结束抽奖
  getEndPrize() {
    const params = {
      activity_id: this.data.activityId,
      record_id: this.data.prizeInfo.record_id,
    };
    this.setData({
      prizeType: "end",
    });
    this.getPrize(params, "end");
  },

  /**
   * 获取刮奖详情
   */
  async getPrize(params, type) {
    try {
      const { data, code } = await scratchPrize(params);
      if (code === 200) {
        this.setData({
          prizeInfo: data.coupon,
          times: data.times,
        });
        if (type === "end") {
          this.getCouponList();
        }
      }
    } catch (error) {}
  },
  /**
   * 再来一次
   */
  bindAgain() {
    this.setData(
      {
        prizeType: "",
      },
      () => {
        this.draw();
        this.getStartPrize();
      }
    );
  },
  /**
   * 刮奖页面初始化
   */
  draw() {
    let self = this;
    this.setData({
      isPrize: this.data.can_scratch,
    });
    this.coverCanvas = new Scratch(self, {
      canvasId: "#coverCanvas",
      width: 630,
      height: 340,
      // maskColor:"",  // 封面颜色
      bgImg:
        "https://cdn-cos.coach.com.cn/uploads/2022/scratch/scratch-guajiang.png", // 封面图片
      overCallBack: () => {
        this.setData({
          isPrize: false,
        });
        this.getEndPrize();
      },
      startCallBack: () => {
        console.log("开始");
      },
    });
  },
  /**
   * 点击核销
   */
  bindWriteOff() {
    const { writeOffstatus } = this.data;
    if (writeOffstatus === 0) {
      this.setData({
        writeOffModal: true,
      });
      this.timer = setInterval(() => {
        this.getCheck();
      }, 3000);
    }
  },
  /**
   * 轮询
   */
  async getCheck() {
    try {
      const params = {
        id: this.data.couponObj.id,
      };
      const { data, code } = await checkCoupon(params);
      if (code === 200) {
        if (data.is_check === 1) {
          clearInterval(this.timer);
          this.setData({
            btnTxt: data.is_check ? "已核销" : "线下核销",
            writeOffstatus: data.is_check,
            writeOffModal: false,
          });
        }
      }
    } catch (error) {}
  },
  /**
   * 分享
   */
  onShareAppMessage: function () {
    let newTime = new Date().valueOf();
    this.onShareReport();
    return {
      title: "幸运刮刮卡",
      path: `/pages/scratch/scratch?activityId=${this.data.activityId}&source=user_share`,
      imageUrl: `https://cdn-cos.coach.com.cn/uploads/2022/scratch/share_coach.png?time=${newTime}`,
    };
  },
  /**
   * 分享上报
   */
  async onShareReport() {
    try {
      const params = {
        activity_id: this.data.activityId,
      };
      await shareInfo(params);
    } catch (error) {}
  },

  getPageHeight(total) {
    switch (total) {
      case 1:
        this.setData({
          pageHeight: 832,
        });
        break;
      case 2:
        this.setData({
          pageHeight: 932,
        });
        break;
      case 3:
        this.setData({
          pageHeight: 1063,
        });
        break;
      default:
        break;
    }
  },
});
