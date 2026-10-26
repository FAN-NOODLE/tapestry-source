import Api from "@/api/public/public";
interface IHCListItem {
  title: string; //标题
  released_at: number; //时间戳
  content: string; //公告内容
}
interface IHCItem {
  key: string;
  type: string;
  content: string | Array<IHCListItem>;
}

// 帮助中心 text内容配置
const textConfig = {
  "myapp_customer_customer-list_list":
    "查看全渠道客户的基础信息、行为信息、订单信息，方便用户对客户进行数据跟进、信息筛选 ",
  "myapp_customer_customer-list_detail":
    "查看全渠道客户的基础信息、行为信息、订单信息，方便用户对客户进行数据跟进、信息筛选",
  "myapp_customer_customer-groups_list":
    "可根据标签自定义动态分组，并可对不同的客户分组创建营销策略",
  "myapp_customer_customer-label_list":
    "通过静态标签和动态标签构建精准用户画像，为客户分组提供基础，360度洞察客户行为，助力精准营销",
  "myapp_customer_customer-label_edit":
    "通过静态标签和动态标签构建精准用户画像，为客户分组提供基础，360度洞察客户行为，助力精准营销",
  "myapp_customer_customer-groups_create":
    "在标签的基础上构建自定义动态分组，以便更好的聚焦不同分组的客户，高效完成客户转化",
  "myapp_customer_customer-groups_show":
    "在标签的基础上构建自定义动态分组，以便更好的聚焦不同分组的客户，高效完成客户转化",
  "myapp_customer_auto-market_market-strategy_look-stratage":
    "根据客户分组，选择不同营销触点、不同触达渠道创建运营策略",
  "myapp_customer_auto-market_market-strategy_set-stratage":
    "根据客户分组，选择不同营销触点、不同触达渠道创建运营策略",
  "myapp_customer_auto-market_market-strategy_edit-stratage":
    "根据客户分组，选择不同营销触点、不同触达渠道创建运营策略",
  "myapp_customer_auto-market_market-strategy_data-track":
    "根据客户分组，选择不同营销触点、不同触达渠道创建运营策略",
  "myapp_customer_auto-market_insert-market-strategy_edit":
    "已配置的运营策略，用户仅需将素材完善，可自动执行内置策略",
  "myapp_customer_auto-market_insert-market-strategy_detail":
    "已配置的运营策略，用户仅需将素材完善，可自动执行内置策略",
  "myapp_customer_auto-market_sms-mgt_create":
    "创建短信模板，短信模板可供营销策略短信渠道作为素材使用",
  "myapp_customer_auto-market_sms-mgt_set-param":
    "创建短信模板，短信模板可供营销策略短信渠道作为素材使用",
  "myapp_enterprise_statistics_leave-custome_detail":
    "记录客户流失情况，包括员工主动删除或被动删除，且提醒员工已被删除，挽回高质客户",
  myapp_enterprise_drainage_welcome_create:
    "客户添加成功时，员工自动向客户发送一条欢迎语信息，可设置文字+图片/网页/小程序",
  myapp_enterprise_drainage_welcome_edit:
    "客户添加成功时，员工自动向客户发送一条欢迎语信息，可设置文字+图片/网页/小程序",
  "myapp_enterprise_drainage_contact-me_add-contact-me":
    "将多个员工添加至一个二维码内，客户扫码后随机分配员工添加，可设置无需经过验证自动通过好友申请",
  "myapp_enterprise_drainage_contact-me_edit-contact-me":
    "将多个员工添加至一个二维码内，客户扫码后随机分配员工添加，可设置无需经过验证自动通过好友申请",
  "myapp_enterprise_message-control_sensitive-word_detail":
    "当员工与客户的聊天过程中，员工主动发送包含该敏感词时，会被记录下来同时提醒员工",
  "myapp_enterprise_organize-structure_staff-detail":
    "同步企业微信内部通讯录/部门员工相关数据，如未同步数据请至“企微配置”进行配置",
  "myapp_wisdom_staff-distribution_distribution-setting_add-goods":
    "全员分销功能的基础设置，完成设置后，导购可以开始正常使用小程序端分销功能",
  "myapp_wisdom_staff-distribution_settlement-record_detail":
    "可以查询每一笔导购分销订单的结算记录和当前结算状态",
  "myapp_wisdom_staff-distribution_balance-management_detail":
    "可以查询和导出参与分销的导购名下的各状态佣金余额，也可以查看导购名下佣金对应的结算单明细",
  "myapp_customer_auto-market_market-strategy_market-list":
    "根据客户分组，选择不同营销触点、不同触达渠道创建运营策略",
  "myapp_customer_auto-market_insert-market-strategy_list":
    "已配置的运营策略，用户仅需将素材完善，可自动执行内置策略",
  "myapp_customer_auto-market_sms-mgt_list":
    "创建短信模板，短信模板可供营销策略短信渠道作为素材使用",
  "myapp_enterprise_statistics_repeat-custome":
    "记录被多个员工重复添加过的客户，效避免销售撞单、群发消息多次打扰的情况发生",
  "myapp_enterprise_statistics_leave-custome_list":
    "记录客户流失情况，包括员工主动删除或被动删除，且提醒员工已被删除，挽回高质客户",
  "myapp_enterprise_message-control_chat-history":
    "永久存储员工与客户/群聊天记录，管控敏感消息提升有效沟通",
  "myapp_enterprise_message-control_chat-timeout":
    "当客户咨询员工时，超过设定时间内未回复消息会被记录下来同时提醒员工，确保服务质量及时性",
  "myapp_enterprise_message-control_sensitive-word_list":
    "当员工与客户的聊天过程中，员工主动发送包含该敏感词时，会被记录下来同时提醒员工",
  "myapp_enterprise_organize-structure_staff-list":
    "同步企业微信内部通讯录/部门员工相关数据，如未同步数据请至“企微配置”进行配置",
  "myapp_enterprise_enterprise-config":
    "企微功能使用必配项，请确保各项配置参数与“企业微信官方”后台一致",
  myapp_enterprise_drainage_welcome_list:
    "客户添加成功时，员工自动向客户发送一条欢迎语信息，可设置文字+图片/网页/小程序",
  "myapp_enterprise_drainage_contact-me_contact-me-list":
    "将多个员工添加至一个二维码内，客户扫码后随机分配员工添加，可设置无需经过验证自动通过好友申请",
  myapp_enterprise_drainage_channel_list:
    "企业自建渠道，通过该渠道添加的客户，自动打上对应的标签",
  "myapp_enterprise_drainage_statistics_channel-statistics-list":
    "统计各渠道累计添加客户数量、净增客户数量、流失客户数量",
  "myapp_enterprise_material-library_material-list":
    "存放企业常用素材，包含文字、图片、视频、网页、文件、小程序，可通过企业微信聊天侧边栏快速发送给客户",
  "myapp_enterprise_friend-circle_friend-circle-list":
    "为员工创建一个虚拟朋友圈，将该朋友圈网址粘贴至员工企微对外个人信息，达到让客户查看员工的历史朋友圈的效果",
  "myapp_wisdom_data-compass":
    "用于查询weclient pro产生的订单数据、顾客数据和分销数据，包含数据概览、趋势图、明细数据等模块。可以按时间维度、导购所属部门等维度筛选后查询和导出",
  "myapp_wisdom_guide-manage_guide-list":
    "企业微信导购小程序账号管理，展示导购的基本信息。新加入企业微信的员工将自动创建【企业导购】类型账号，后续可在【更多】中转换【企业管理】账号类型；已加入导购账号不可移除，员工从企业微信离职后，导购账号将自动删除。",
  "myapp_wisdom_guide-manage_guide-set":
    "导购业绩设置满足不同业绩计算场景下的需求，品牌可根据自身实际计算方式进行灵活配置",
  "myapp_wisdom_guide-manage_guide-pay":
    "用户在小程序商城购买支付后，微信会收到微信支付的消息推送，点击消息推送中的“联系商家导购员”按钮，可以查看设置的导购员信息，扫码可添加导购员企业微信进行后续服务",
  "myapp_wisdom_staff-distribution_distribution-setting_list":
    "全员分销功能的基础设置，完成设置后，导购可以开始正常使用小程序端分销功能",
  "myapp_wisdom_staff-distribution_settlement-record_list":
    "可以查询每一笔导购分销订单的结算记录和当前结算状态",
  "myapp_wisdom_staff-distribution_withdrawal-approval":
    "导购申请提现佣金的申请单会按不同状态展示，导购提交时为“待审批”状态，可以进行批准和驳回操作，批准后才会由第三方平台转账至导购银行卡",
  "myapp_wisdom_staff-distribution_withdrawal-setting":
    "全员分销佣金提现的基础设置，务必完整设置，以保证导购分销后可以正常结算和提现分销佣金",
  "myapp_wisdom_staff-distribution_balance-management_list":
    "可以查询和导出参与分销的导购名下的各状态佣金余额，也可以查看导购名下佣金对应的结算单明细",
  "myapp_setting_authorization-settings_enterprise-config":
    "企微功能使用必配项，请确保各项配置参数与“企业微信官方”后台一致",
  "myapp_enterprise_drainage_fission-activities":
    "该活动支持裂变好友、裂变群成员和裂变品牌会员，通过欢迎语和企微侧边栏推送让好友带来他的好友进一步提炼现有好友资源价值，同时提高客户活跃度及转化率。",
  "myapp_enterprise_drainage_fission-activities_list":
    "该活动支持裂变好友、裂变群成员和裂变品牌会员，通过欢迎语和企微侧边栏推送让好友带来他的好友进一步提炼现有好友资源价值，同时提高客户活跃度及转化率。",
  "myapp_enterprise_drainage_fission-activities_details":
    "该活动支持裂变好友、裂变群成员和裂变品牌会员，通过欢迎语和企微侧边栏推送让好友带来他的好友进一步提炼现有好友资源价值，同时提高客户活跃度及转化率。",
  "myapp_enterprise_drainage_fission-activities_created":
    "该活动支持裂变好友、裂变群成员和裂变品牌会员，通过欢迎语和企微侧边栏推送让好友带来他的好友进一步提炼现有好友资源价值，同时提高客户活跃度及转化率。",
  "myapp_enterprise_drainage_fission-activities_edit":
    "该活动支持裂变好友、裂变群成员和裂变品牌会员，通过欢迎语和企微侧边栏推送让好友带来他的好友进一步提炼现有好友资源价值，同时提高客户活跃度及转化率。",
  "myapp_customer_auto-market_externalchain-manage_externalchain-list":
    "在带参链接的基础上生成带参短外链，可用于投放，在微信内部/外部均可唤醒小程序，并可配置中间H5页面，提升短外链的打开率，且可追踪点击短外链的客户数据。",
  "myapp_customer_auto-market_externalchain-manage_externalchain-set":
    "在带参链接的基础上生成带参短外链，可用于投放，在微信内部/外部均可唤醒小程序，并可配置中间H5页面，提升短外链的打开率，且可追踪点击短外链的客户数据。",
  "myapp_customer_auto-market_externalchain-manage_externalchain-edit":
    "在带参链接的基础上生成带参短外链，可用于投放，在微信内部/外部均可唤醒小程序，并可配置中间H5页面，提升短外链的打开率，且可追踪点击短外链的客户数据。",
  "myapp_customer_auto-market_externalchain-manage_externalchain-look":
    "在带参链接的基础上生成带参短外链，可用于投放，在微信内部/外部均可唤醒小程序，并可配置中间H5页面，提升短外链的打开率，且可追踪点击短外链的客户数据。"
};

/**
 * 定义要 显示的帮助中心页面 数据
 */
const helpConfig: Array<IHCItem> = [];

for (const k in textConfig) {
  helpConfig.push({
    key: k,
    type: "text",
    content: textConfig[k]
  });
}

// cache 配置
const helpCache = {
  cacheNum: 0,
  keyCache: {
    myapp_home: "loop" // 每次都需要请求
  }
};

export async function getConfig(key: string) {
  if (helpCache.cacheNum === 0 || helpCache.keyCache[key] === "loop") {
    // getNotice();
    const res = await Api.getHelp();
    helpCache.cacheNum++;
    if (res["code"] === 200) {
      res.data.list.forEach((item: any) => {
        if (item.type === "list") {
          item.content = JSON.parse(item.content);
        }
        helpConfig.push(item);
      });
    }
    return helpConfig.find((item: IHCItem) => {
      return item.key === key;
    });
  } else {
    return helpConfig.find((item: IHCItem) => {
      return item.key === key;
    });
  }
}
