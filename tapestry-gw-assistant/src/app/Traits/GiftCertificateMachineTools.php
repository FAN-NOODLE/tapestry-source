<?php

namespace App\Traits;

use App\Models\FirstPurchaseGiftCertificate;
use App\Models\GiftCertificateTrace;

/**
 * 首购礼资格状态机工具类
 * Trait OrderMachineTools
 *
 * @package App\Traits
 */
trait GiftCertificateMachineTools
{
    /**
     * 更新首购礼资格状态
     *
     * @param $giftCertificateId //首购礼资格id
     * @param $update //要更新的字段
     */
    public function updateGiftCertificate($giftCertificateId, $update)
    {
        FirstPurchaseGiftCertificate::query()->where('id', $giftCertificateId)->update($update);
    }

    /**
     * 插入首购礼资格操作记录
     *
     * @param $certificateId //首购礼资格id
     * @param $from //当前状态
     * @param $to //要转换成的状态
     * @param $event //事件名
     * @param $eventDesc //事件描述
     * @param $update //要更新的字段
     */
    public function recordGiftCertificateTrace($certificateId, $from, $to, $event, $eventDesc, $update = [])
    {
        $insert = [
            'gift_certificate_id'     => $certificateId,
            'from_certificate_status' => $from,
            'to_certificate_status'   => $to,
            'event'                   => $event,
            'event_desc'              => $eventDesc,
            'remark'                  => empty($update) ? '' : json_encode($update, JSON_UNESCAPED_UNICODE),
            'created_at'              => time(),
        ];

        GiftCertificateTrace::query()->insert($insert);
    }
}
