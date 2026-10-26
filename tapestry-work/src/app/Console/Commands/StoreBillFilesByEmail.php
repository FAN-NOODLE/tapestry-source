<?php
/**
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月20日17:9:30
 */

namespace App\Console\Commands;

use App\Library\WechatNotice\Chat;
use App\Models\BillStore;
use App\Models\Scrm\Departments;
use App\Models\Scrm\MemberDepartments;
use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\ExternalPayOrder;
use App\Models\ScrmEnterprise\ExternalPayRefundOrder;
use App\Models\StoreDepartment;
use Exception;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Mail;
use Throwable;

class StoreBillFilesByEmail extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'store-bill-files-by-email {--end_time=}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '邮件发送对外收款附件给各个门店';


    public function handle()
    {
        $enterpriseId = env('COACH_ENTERPRISE_ID', 3363603667487424);
        $endTime = time();
        try{
            $paramEndTime = intval($this->option('end_time'));
        } catch (Throwable $exception ){
        }
        $endTime = !empty($paramEndTime) ? $paramEndTime : $endTime;
        $endTime = $endTime - 3600; // 避免脚本执行的时候 订单数据还没拉到，所有数据获取时间往前推一个小时
        $startTime = $endTime - 86400;
        $dateStr = date("Ymd", $endTime);
        $arr = BillStore::query()->where('enterprise_id',$enterpriseId)->pluck('store_code')->toArray();
//        $arr = ['OCF38','OC109','OC185','OC151','OCF29','OC48','OCF39','OC220','OC241','OC90','OC134','OC201','OC248','OCF25','OC97','OC183','OC39','OC49','OC239','OC61','OC246','OC174','OC228','OC135','OC233','OC230','OC91','OC209','OC218','OC320','OC332','OC245'];
        //查询配置了邮件的门店
        $storeList = StoreDepartment::query()->where("enterprise_id", $enterpriseId)
            ->where('email', '!=', "")
            ->whereIn('third_store_code', $arr)
            ->get(['wechat_department_id', 'email', 'third_store_code'])
            ->keyBy('wechat_department_id')
            ->toArray();
        //获取各个门店的导购id
        $dpList = Departments::query()->with("membersRel:department_id,member_id")
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('third_id', array_keys($storeList))
            ->get(['third_id', 'id'])
            ->keyBy('third_id')
            ->toArray();
        $trueArr = [];
        foreach ($dpList as $thirdId => $val) {
            $trueArr[] = [
                'members' => array_column($val['members_rel'], 'member_id'),
                'email' => trim($storeList[$thirdId]['email']),
                'store_code' => $storeList[$thirdId]['third_store_code']
            ];
        }
        $res = [];
        //查询各个导购集合的统计并发送邮件
        foreach ($trueArr as $val) {
            //收款单
            $payOrder = $this->getPayOrderData($enterpriseId, $val['members'], $startTime, $endTime);
            $payFileName = $val['store_code'] . '收款WeclientPro小程序日结单' . $dateStr . '.csv';
            $this->saveCsv($payFileName, $payOrder);
            //退款单
            $refundOrder = $this->getRefundOrderData($enterpriseId, $val['members'], $startTime, $endTime);
            $refundFileName = $val['store_code'] . '退款WeclientPro小程序日结单' . $dateStr . '.csv';
            $this->saveCsv($refundFileName, $refundOrder);
            try {
                $res[$val['email']] = $this->sendemail([
                    'to' => $val['email'],
                    'payFileName' => storage_path($payFileName),
                    'refundFileName' => storage_path($refundFileName),
                    'storeCode' => $val['store_code'],
                    'date' => $dateStr,
                    'payTotal' => count($payOrder) - 1,
                    'refundTotal' => count($refundOrder) - 1,
                ]);
                sleep(10);
            } catch (Exception $e) {
                print_r($e->getMessage());
                Chat::getInstance()->send('企业收款发送邮件失败', [
                    'to' => $val['email'],
                    'storeCode' => $val['store_code'],
                    'date' => $dateStr,
                    'payTotal' => count($payOrder) - 1,
                    'refundTotal' => count($refundOrder) - 1,
                    'error'=>$e->getMessage(),
                ]);
            }

        }
        print_r($res);
        return true;
    }

    private function saveCsv($fileName, $data)
    {
        $file = fopen(storage_path($fileName), "w");
        fwrite($file, chr(239) . chr(187) . chr(191));// 加上bom头，系统自动默认为UTF-8编码
        foreach ($data as $line) {
            fputcsv($file, $line);
        }
        fclose($file);
    }

    //查询各个导购集合的统计
    private function getPayOrderData($enterpriseId, $memberIds, $startTime, $endTime)
    {
        $payList = ExternalPayOrder::query()
            ->setConn($enterpriseId)
            ->with('user:openid,nick,union_code')
            ->whereIn('member_id', $memberIds)
            ->whereBetween('pay_time', [$startTime, $endTime])
            ->get()
            ->toArray();
        $members = Members::query()->where('enterprise_id', $enterpriseId)
            ->whereIn('id', array_unique(array_column($payList, 'member_id')))
            ->get(['id', 'name'])
            ->keyBy('id')
            ->toArray();
        $departments = MemberDepartments::query()->with('departments:id,name')->where('enterprise_id', $enterpriseId)
            ->whereIn('member_id', array_column($members, 'id'))
            ->get()
            ->keyBy('member_id')
            ->toArray();
        $arr[] = ['订单编号', '支付时间', '支付金额', '顾客昵称', '顾客memberid', 'SA工号', 'SA姓名', '收款门店'];

        foreach ($payList as $val) {
            $arr[] = [
                $val['out_trade_no'],
                date("Y-m-d H:i:s", $val['pay_time']),
                $val['total_fee'] / 100,
                $val['user']['nick'],//顾客昵称
                $val['user']['union_code'], //顾客memberid
                $val['payee_userid'],
                isset($members[$val['member_id']]) ? $members[$val['member_id']]['name'] : '',
                isset($members[$val['member_id']]) && isset($departments[$members[$val['member_id']]['id']]) ?
                    $departments[$members[$val['member_id']]['id']]['departments']['name'] : '',
            ];
        }
        return $arr;

    }

    private function getRefundOrderData($enterpriseId, $memberIds, $startTime, $endTime): array
    {
        $refundList = ExternalPayRefundOrder::query()
            ->setConn($enterpriseId)
            ->with('order')
            ->whereIn('refund_member_id', $memberIds)
            ->where('refund_status', 2)
            ->whereBetween('updated_at', [$startTime, $endTime])
            ->get()
            ->toArray();
        $members = Members::query()->where('enterprise_id', $enterpriseId)
            ->whereIn('id', array_unique(array_column($refundList, 'refund_member_id')))
            ->get(['id', 'name'])
            ->keyBy('id')
            ->toArray();
        $departments = MemberDepartments::query()->with('departments:id,name')
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('member_id', array_column($members, 'id'))
            ->get()
            ->keyBy('member_id')
            ->toArray();
        $arr[] = ['退款编号', '退款门店', '订单所属顾问', '申请时间', '退款单号', '退款金额', '原订单金额', '原销售单号'];
        foreach ($refundList as $val) {
            $arr[] = [
                $val['out_refund_no'],
                isset($members[$val['refund_member_id']]) && isset($departments[$members[$val['refund_member_id']]['id']]) ?
                    $departments[$members[$val['refund_member_id']]['id']]['departments']['name'] : '',
                $val['refund_userid'],
                date("Y-m-d H:i:s", $val['refund_reqtime']),
                $val['out_refund_no'],
                $val['refund_fee'] / 100,
                ($val['order']['total_fee'] ?? 0) / 100,
                $val['out_trade_no'],
            ];
        }
        return $arr;
    }

    //发送邮件
    private function sendEmail($params)
    {
        //发送邮件
        Mail::send("emails.storeBillFiles", $params, function ($message) use ($params) {
            $message->to($params['to'])->subject($params['storeCode'] . ' WeClient Pro小程序日结单' . $params['date']);
            if (is_file($params['payFileName'])) {
                $message->attach($params['payFileName']);
            }
            if (is_file($params['refundFileName'])) {
                $message->attach($params['refundFileName']);
            }
        });
        return Mail::failures();
    }


}
