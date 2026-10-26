// 所属分组group
// 类型type 1-单选 2-多选 3-时间范围 4-输入数字范围
// 判断是不是到2026年了
import { update_time } from '@/utils/const'

export const options = [
  {
    key: 'card_type',
    label: '会员类型',
    group: 1,
    type: 2,
    tips: []
  },
  {
    key: 'birth_month',
    label: '生日月份',
    group: 1,
    type: 2,
    tips: []
  },
  {
    key: 'manage_type',
    label: '顾客筛选',
    group: 1,
    type: 1,
    tips: []
  },
  {
    key: 'member_ids',
    label: '员工筛选',
    group: 1,
    type: 2,
    tips: []
  },
  {
    key: 'reserve_type',
    label: '近 30 天邀约状态',
    group: 2,
    type: 1,
    tips: [
      {
        title: '一、邀约中',
        content: '企业内，近30天存在进行中的邀约'
      },
      {
        title: '二、邀约成功',
        content: '企业内，近30天邀约成功的顾客'
      },
      {
        title: '三、无邀约',
        content: '企业内，近30天无邀约记录的顾客'
      }
    ]
  },
  {
    key: 'task_type',
    label: '近 30 天任务触达',
    group: 2,
    type: 1,
    tips: [
      {
        title: '一、已触达',
        content: '企业内，近30天通过任务沟通被触达的顾客'
      },
      {
        title: '二、未触达',
        content: '企业内，近30天未出现在沟通任务名单中的顾客，或出现在任务名单中但未被触达的顾客'
      }
    ]
  },
  {
    key: 'task_have',
    label: '近 30 天任务状况',
    group: 2,
    type: 1,
    tips: [
      {
        title: '一、无任务',
        content: '企业内，近30天无企微沟通任务关联匹配的顾客'
      },
      {
        title: '二、有任务',
        content: '企业内，近30天有企微沟通任务关联匹配的顾客'
      }
    ]
  },
  {
    key: 'coupon_type',
    label: '优惠券筛选',
    group: 3,
    type: 2,
    tips: [
      {
        title: '一、30天内临期券',
        content: '持有有效期少于30天的优惠券的顾客'
      },
      {
        title: '二、待使用券',
        content: '持有待使用优惠券的顾客'
      }
    ]
  },
  {
    key: 'purchased_date',
    label: '历史消费时间',
    group: 3,
    type: 3,
    tips: []
  },
  {
    key: 'not_purchased_date',
    label: '近期未消费时间',
    group: 3,
    type: 1,
    tips: []
  },
  {
    key: 'history_amount_total',
    label: '历史购买总金额 LTV',
    group: 3,
    type: 4,
    tips: []
  },
  {
    key: 'amount_total',
    label: 'L12M 总金额',
    group: 3,
    type: 4,
    tips: []
  },
  {
    key: 'rel_create_time',
    label: '添加企微好友时间',
    group: 4,
    type: 2,
    tips: []
  },
  {
    key: 'rel_interact_time',
    label: '上次触达时间',
    group: 4,
    type: 2,
    tips: [
      {
        title: '一、我的',
        content: '客户与导购本人的最近一次消息互动时间（包括SA主动发起和客户主动发起）'
      },
      {
        title: '二、企业',
        content: '客户与企业中所有导购的最近一次消息互动时间（包括SA主动发起和客户主动发起）'
      }
    ]
  },
  {
    key: 'rel_active_time',
    label: '活跃沟通顾客',
    group: 4,
    type: 2,
    tips: [
      {
        title: '一、我的',
        content: '客户与导购本人最近一次产生两次以上来回会话的沟通时间'
      },
      {
        title: '二、企业',
        content: '客户与企业中所有导购最近一次产生两次以上来回会话的沟通时间'
      }
    ]
  },
  {
    key: 'chat_tag',
    label: '顾客沟通意愿',
    group: 4,
    type: 1,
    tips: []
  },
  {
    key: 'high_willing',
    label: '高意愿用户',
    group: 5,
    type: 1,
    tips: []
  },
  {
    key: 'price_sensitivity',
    label: '价格敏感度',
    group: 5,
    type: 1,
    tips: []
  },
  {
    key: 'category_preference',
    label: '商品偏好',
    group: 5,
    type: 2,
    tips: []
  },
  {
    key: 'age_tag',
    label: '顾客类型 AGE',
    group: 5,
    type: 1,
    tips: []
  }
]
if (new Date() >= new Date(update_time) && new Date() <= new Date('2026-12-31 23:59:59')) {
  options[0].tips = [
    {
      title: '',
      content: '因会员卡等体系调整，2026年起，银卡/金卡/钻卡会员会籍到期后将更新为【已购会员】，会籍到期时间可在顾客详情中查看。\n\n如需筛选全部已经购买的会员可全选Retention中的选项。'
    }
  ]
}
export default options
