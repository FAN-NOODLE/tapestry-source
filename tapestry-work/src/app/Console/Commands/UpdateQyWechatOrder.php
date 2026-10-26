<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\QyWechat;
use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\ExternalPayOrder;
use App\Models\ScrmEnterprise\ExternalPayRefundOrder;
use App\Services\EnterpriseService;
use Illuminate\Console\Command;

class UpdateQyWechatOrder extends Command
{

    const MAX_ONCE_DATS = 30;


    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:update-qy-wechat-order';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新30天的收款记录（企微获取）';


    public function handle()
    {
        $enterpriseId = env('COACH_ENTERPRISE_ID');
        if (!empty($enterpriseId)) {
            // 查询企业列表
            $enterprise = (new EnterpriseService())->detail($enterpriseId)['data'];
            $qyWechat = new QyWechat();
            $now = time();
            $memberIdByUserId = Members::query()->where('enterprise_id', $enterprise['id'])
                ->get(['id', 'userid'])->keyBy('userid')->toArray();
            $dataOrm = ExternalPayOrder::query()->setConn($enterprise['id']);
            $end = $dataOrm->orderByDesc('pay_time')->value('pay_time');
            if (!empty($end)) {
                // 查30天的数据
                $newOrderData = [];  //新的订单数据
                $newRefundData = []; //新的退款数据
                $orderNoToExternalUseridArr = [];// 订单号对应的外部联系人id
                $start = $end - self::MAX_ONCE_DATS * 86400;  //开始时间是结束时间之前30天
                //获取老的订单数据(字段顺序需新老数据一致)
                $oldOrderData = ExternalPayOrder::query()->setConn($enterprise['id'])
                    ->whereBetween('pay_time', [$start, $end])
                    ->get(['out_trade_no', 'trade_state', 'total_refund_fee'])
                    ->toArray();
                $nextCursor = '';
                // 企微查询新数据（企微一次只能获取1000条数据，需要循环获取）
                while (true) {
                    $params = [
                        'begin_time' => $start,
                        'end_time' => $end,
                        'limit' => 1000,
                    ];
                    if (!empty($nextCursor)) {
                        $params['cursor'] = $nextCursor;
                    }
                    $res = $qyWechat->getQyResult($params, [
                        'path' => 'cgi-bin/externalpay/get_bill_list',
                        'method' => 'POST',
                        'enterprise_id' => $enterprise['id'],
                        'reqPL' => 'crs_store_assistant_self'
                    ]);
                    foreach ($res['bill_list'] as $val) {
                        //企微接口更新我们只接收bill_type为0的信息
                        if (isset($val['bill_type']) && $val['bill_type']) {
                            continue;
                        }
                        $newOrderData[] = [
                            'out_trade_no' => $val['out_trade_no'],
                            'trade_state' => $val['trade_state'],
                            'total_refund_fee' => $val['total_refund_fee'],
                        ];
                        if (!empty($val['refund_list'])) {
                            $newRefundData[] = [
                                'refund_list' => $val['refund_list'],
                                'out_trade_no' => $val['out_trade_no']
                            ];
                        }
                        $orderNoToExternalUseridArr[$val['out_trade_no']] = $val['external_userid'];
                    }
                    if (!isset($res['next_cursor']) || empty($res['next_cursor'])) {
                        break;
                    }
                    $nextCursor = $res['next_cursor'];
                }
                //处理需要更新的主表数据
                $updateData = arrayDiffTwoDimension($newOrderData, $oldOrderData);
                foreach ($updateData as $val) {
                    $dataOrm->where('out_trade_no', $val['out_trade_no'])->update($val);
                }

                $deleteRefundOrderNo = [];
                $insertRefundOrderData = [];
                foreach ($newRefundData as $val) {
                    foreach ($val['refund_list'] as $refund) {
                        $deleteRefundOrderNo[] = $refund['out_refund_no'];
                        $insertRefundOrderData[] = [
                            'id' => idGenerate(),
                            'out_trade_no' => $val['out_trade_no'],
                            'out_refund_no' => $refund['out_refund_no'],
                            'external_userid' => $orderNoToExternalUseridArr[$val['out_trade_no']] ?? '',
                            'refund_userid' => $refund['refund_userid'],
                            'refund_member_id' => isset($memberIdByUserId[$refund['refund_userid']]) ? $memberIdByUserId[$refund['refund_userid']]['id'] : '',
                            'refund_comment' => $refund['refund_comment'],
                            'refund_reqtime' => $refund['refund_reqtime'],
                            'refund_status' => $refund['refund_status'],
                            'refund_fee' => $refund['refund_fee'],
                            'created_at' => $now,
                            'updated_at' => $now
                        ];
                    }
                }

                //查询要插入的退款数据中已经退款成功的
                $deleteRefundOrderNo = array_unique($deleteRefundOrderNo);
                $refundSuccessOrderNo = ExternalPayRefundOrder::query()
                    ->setConn($enterprise['id'])
                    ->whereIn('out_refund_no', $deleteRefundOrderNo)
                    ->where('refund_status', 2)
                    ->pluck('out_refund_no')
                    ->toArray();
                //删除非退款成功的原始数据
                if (!empty($deleteOrderNo = array_diff($deleteRefundOrderNo, $refundSuccessOrderNo))) {
                    ExternalPayRefundOrder::query()->setConn($enterprise['id'])
                        ->whereIn('out_refund_no', $deleteOrderNo)
                        ->forceDelete();
                }

                foreach ($insertRefundOrderData as $key => $val) {
                    // 如果要新增的退款数据中存在之前已经退款成功的订单，则忽略
                    if (in_array($val['out_refund_no'], $refundSuccessOrderNo)) {
                        unset($insertRefundOrderData[$key]);
                    }
                }
                if (!empty($insertRefundOrderData)) {
                    //插入数据有变动及新增的退款信息
                    ExternalPayRefundOrder::query()->setConn($enterprise['id'])->insert(array_values($insertRefundOrderData));
                }
            }
        }
    }

}
