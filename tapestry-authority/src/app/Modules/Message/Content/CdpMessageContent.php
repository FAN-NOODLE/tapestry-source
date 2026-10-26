<?php

namespace App\Modules\Message\Content;

class CdpMessageContent extends AbstractContent
{
    /**
     * 好友助力失败
     *
     * @return array
     */
    public function helpFail()
    {
        return [
            'activity_name'     => '“' . $this->message['spu_name'] . '”好友助力活动',
            'activity_progress' => '好友助力失败',
            'notice'            => '抱歉您未在助力有效期内完成助力'
        ];
    }

    /**
     * 好友助力完成付邮费
     *
     * @return array
     */
    public function helpOverPayExpressFee()
    {
        return [
            'activity_name'     => '“' . $this->message['spu_name'] . '”好友助力活动',
            'activity_progress' => '好友助力完成',
            'notice'            => '恭喜您好友助力完成，请及时支付邮费'
        ];
    }

    /**
     * 好友助力完成
     *
     * @return array
     */
    public function helpOver()
    {
        return [
            'activity_name'     => '“' . $this->message['spu_name'] . '”好友助力活动',
            'activity_progress' => '好友助力完成',
            'notice'            => '恭喜您好友助力完成，请及时查看'
        ];
    }

    /**
     * 好友助力开始
     *
     * @return array
     */
    public function helpStart()
    {
        return [
            'activity_name' => '您预约的助力抽奖活动，即将开始，点击立即查看',
            'start_at'      => date('Y-m-d H:i:s', $this->message['start_at']),
        ];
    }

    /**
     * 好友助力人数增加
     *
     * @return array
     */
    public function helpPeopleAdd()
    {
        return [
            'activity_name'     => '“' . $this->message['spu_name'] . '”助力抽奖活动',
            'activity_progress' => $this->message['activity_progress'] . '人已经帮你助力',
            'end_at'            => date('Y-m-d H:i:s', $this->message['end_at']),
            'notice'            => '加入官方社群更容易提升助力排名'
        ];
    }

    /**
     * 开奖结果
     *
     * @return array
     */
    public function luckResult()
    {
        return [
            'activity_name' => $this->message['spu_name'],
            'luck_time'     => date('Y-m-d H:i:s', $this->message['luck_time']),
            'luck_result'   => '点击查看开奖结果',
            'notice'        => '如果中奖请及时支付下单，过期视为放弃资格'
        ];
    }

    /**
     * 拼团人数增加
     *
     * @return array
     */
    public function groupPeopleAdd()
    {
        return [
            'spu_name'       => $this->message['spu_name'],
            'need_num'       => $this->message['need_num'],
            'group_progress' => $this->message['group_progress'] . '人已参与拼团',
            'notice'         => '分享至群更容易完成拼团',
        ];
    }

    /**
     * 拼团失败
     *
     * @return array
     */
    public function groupFail()
    {
        return [
            'spu_name'       => $this->message['spu_name'],
            'need_num'       => $this->message['need_num'],
            'group_progress' => $this->message['group_progress'],
            'refund_money'   => $this->message['refund_money'] / 100,
            'notice'         => '未在有效期内完成拼团',
        ];
    }

    /**
     * 砍价人数增加
     *
     * @return array
     */
    public function hagglePeopleAdd()
    {
        return [
            'spu_name'        => $this->message['spu_name'],
            'sell_price'      => $this->message['sell_price'] / 100,
            'activity_price'  => $this->message['activity_price'] / 100,
            'haggle_progress' => $this->message['haggle_progress'] . '人已帮助砍价',
            'end_at'          => $this->message['rest_time'],
        ];
    }

    /**
     * 砍价结果
     *
     * @return array
     */
    public function haggleResult()
    {
        return [
            'spu_name'        => $this->message['spu_name'],
            'sell_price'      => $this->message['sell_price'] / 100,
            'activity_price'  => $this->message['activity_price'] / 100,
            'haggle_progress' => (1 == $this->message['haggle_progress']) ? '砍价成功' : '砍价失败'
        ];
    }

    /**
     * 直播开始
     *
     * @return array
     */
    public function liveStart()
    {
        return [
            'live_title' => '您预约的直播“' . $this->message['live_title'] . '”开始了，点击进入直播间',
            'start_time' => date('Y-m-d H:i:s', $this->message['start_time']),
        ];
    }

    /**
     * 订单等待支付
     *
     * @return array
     */
    public function orderWaitPay()
    {
        return [
            'wait_pay_money'   => '￥' . $this->message['wait_pay_money'] / 100,
            'create_at'        => date('Y-m-d H:i:s', $this->message['create_at']),
            'order_no'         => $this->message['order_no'],
            'order_goods_info' => $this->message['order_goods_info'],
            'notice'           => '订单为您预留到' . date('Y-m-d H:i:s', $this->message['remain_at']) . '，支付金额以实际支付时显示金额为准'
        ];
    }

    /**
     * 订单发货
     *
     * @return array
     */
    public function orderExpress()
    {
        return [
            'order_no'         => $this->message['order_no'],
            'logistic_name'    => $this->message['logistic_name'],
            'logistic_no'      => $this->message['logistic_no'],
            'delivery_address' => $this->message['delivery_address'],
            'order_goods_info' => $this->message['order_goods_info'],
        ];
    }

    /**
     * 订单退款
     *
     * @return array
     */
    public function orderRefund()
    {
        return [
            'order_no'            => $this->message['order_no'],
            'refund_apply_status' => $this->message['refund_apply_status'],
            'refund_fee'          => '￥' . $this->message['refund_fee'] / 100,
            'refund_goods_info'   => $this->message['refund_goods_info']
        ];
    }

    /**
     * 积分变更
     *
     * @return array
     */
    public function scoreChange()
    {
        return [
            'score_change' => $this->message['score_change'],
            'score_total'  => $this->message['score_total'],
            'type'         => $this->message['type']
        ];
    }
}
