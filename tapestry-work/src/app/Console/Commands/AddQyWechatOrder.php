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
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class AddQyWechatOrder extends Command
{

    const MAX_ONCE_DATS = 30;


    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'add-qy-wechat-order';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '每日企微对外收款记录（企微获取）';

    /**
     * @throws \Exception
     */
    public function handle()
    {
        $enterpriseId = env('COACH_ENTERPRISE_ID');
        if (empty($enterpriseId)) {
            return;
        }
        $day = 10;
        // 查询企业列表
        $enterprises = [(new EnterpriseService())->detail($enterpriseId)['data']];
//        $enterprises = (new Enterprise())->lists()['data']; //所有企业

        $now = time();
        $qyWechat = new QyWechat();
        foreach ($enterprises as $enterprise) {
            $dataOrm = ExternalPayOrder::query()->setConn($enterprise['id']);
            $lastPayTime = $dataOrm->orderByDesc('pay_time')->value('pay_time') ?: 0;
            //取消时间限制
            /*if ($now - $lastPayTime < 86400) {
                // 如果最新数据是24内的，就忽略
                continue;
            }*/
            $memberIdByUserId = Members::query()->where('enterprise_id', $enterprise['id'])
                ->get(['id','userid'])->keyBy('userid')->toArray();
            $insertArr = [];
            $refundArr = [];
            $nextCursor = ''; //查询游标
            while (true) {
                // 最多获取10天的数据
                $params = [
                    'begin_time' => ($now - $lastPayTime > $day * 86400) ? ($now - $day * 86400) : ($lastPayTime + 1),
                    'end_time' => $now,
                    'limit' => 10,
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
                    if(isset($val['bill_type']) && $val['bill_type']){
                        continue;
                    }
                    $insertArr[] = [
                        'id' => idGenerate(),
                        'transaction_id' => $val['transaction_id'],
                        'pay_time' => $val['pay_time'],
                        'trade_state' => $val['trade_state'],
                        'out_trade_no' => $val['out_trade_no'],
                        'external_userid' => $val['external_userid'],
                        'total_fee' => $val['total_fee'],
                        'payee_userid' => $val['payee_userid'],
                        'member_id' => isset($memberIdByUserId[$val['payee_userid']]) ? $memberIdByUserId[$val['payee_userid']]['id'] : '' ,
                        'payment_type' => $val['payment_type'],
                        'mch_id' => $val['mch_id'],
                        'remark' => $val['remark'],
                        'commodity_list' => json_encode($val['commodity_list']),
                        'contact_info' => json_encode($val['contact_info']),
                        'delivery_type' => isset($val['contact_info']['phone']) && $val['contact_info']['phone']  ? 2 : 1,
                        'total_refund_fee' => $val['total_refund_fee'],
                        'created_at' => $now,
                        'updated_at' => $now
                    ];
                    $useridArr[] = $val['payee_userid'];
                    if (!empty($val['refund_list'])) {
                        foreach ($val['refund_list'] as $refund) {
                            $refundArr = [
                                'id' => idGenerate(),
                                'out_trade_no'  => $val['out_trade_no'],
                                'out_refund_no' => $refund['out_refund_no'],
                                'external_userid' => $val['external_userid'],
                                'refund_userid' => $refund['refund_userid'],
                                'refund_member_id' => isset($memberIdByUserId[$refund['refund_userid']]) ? $memberIdByUserId[$refund['refund_userid']]['id'] : '' ,
                                'refund_comment' => $refund['refund_comment'],
                                'refund_reqtime' => $refund['refund_reqtime'],
                                'refund_status' => $refund['refund_status'],
                                'refund_fee' => $refund['refund_fee'],
                                'created_at' => $now,
                                'updated_at' => $now
                            ];
                        }
                    }
                }
                if (empty($res['next_cursor'])) {
                    break;
                }

                $nextCursor = $res['next_cursor'];
            }
//            dd($refundArr,$insertArr);
            $dataOrm->insert($insertArr);
            ExternalPayRefundOrder::query()->setConn($enterprise['id'])->insert($refundArr);

        }
    }


}
