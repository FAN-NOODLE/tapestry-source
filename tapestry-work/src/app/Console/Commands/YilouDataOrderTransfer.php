<?php

namespace App\Console\Commands;

use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\ExternalPayOrder;
use App\Models\ScrmEnterprise\ExternalPayRefundOrder;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Services\CommonService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class YilouDataOrderTransfer extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'yilou-data-order-transfer';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '一楼数据订单数据迁移';

    protected $size = 20;

    protected $lastPage = 0;

    protected $page = 1;

    public function handle()
    {
        $enterpriseId = env('COACH_ENTERPRISE_ID');
        $lab = true;

        while ($lab) {
            //排除 待付款和已过期订单
            $order = DB::connection('yilou_data')->table('coach_order')
                ->whereNotIn('Status', [0,1,5])
                //->where("OrderID",793)
                //->where('OrderNO','1585292732845723067')
                ->paginate($this->size, ['*'], 'page', $this->page);

            $this->lastPage = $order->lastPage();
            if ($this->lastPage == $this->page) {
                $lab = false;

            }
            $this->page++;

            //业务
            $outTradeNo = collect($order->items())->pluck('OrderNO')->all();
            //去查询有没有写入过数据表
            $payOrder = ExternalPayOrder::query()->setConn($enterpriseId)->whereIn('out_trade_no', $outTradeNo)
                ->pluck('out_trade_no')
                ->toArray();
            //差集订单编号
            $payOrderDiff = collect($outTradeNo)->diff($payOrder);

            //一楼订单数据
            $orderArray = collect($order->items())->keyBy('OrderNO')->all();

            //查询退款订单
            $orderRefund = DB::connection('yilou_data')->table('coach_order_refund')
                ->whereIn('OrderNO', $payOrderDiff)
                ->get();

            //拼接退款金额
            foreach ($orderRefund as $refund) {
                if (isset($orderArray[$refund->OrderNO]->total_refund_fee)) {
                    $orderArray[$refund->OrderNO]->total_refund_fee += $refund->RefundMoney;
                } else {
                    $orderArray[$refund->OrderNO]->total_refund_fee = $refund->RefundMoney;
                }
            }

            //外部联系人用户信息
            $unionid        = collect($order->items())->pluck('UnionID')->all();
            $externalUserId = WxUserCorpExternal::query()->setConn($enterpriseId)
                ->whereIn('unionid', $unionid)
                ->select('openid', 'unionid')
                ->get()
                ->keyBy('unionid');
            //导购信息
            $memberID = collect($order->items())->pluck('EmpnoCode')->all();
            $membersId   = Members::query()->where(['enterprise_id' => $enterpriseId])
                ->whereIn('userid',$memberID)
                ->select('id', 'userid')
                ->get()
                ->keyBy('userid');

            //定义插入数据容器
            $orderAddArr = [];

            foreach ($payOrderDiff as $row) {
                $orderItem     = collect($orderArray[$row])->toArray();
                $orderAddArr[] = [
                    'id'               => idGenerate(),
                    'transaction_id'   => $orderItem['OrderWXNO'],
                    'trade_state'      => $this->setTradeState($orderItem['RefundStatus']),
                    'order_state'      => $this->setOrderState($orderItem['Status']),
                    'pay_time'         => strtotime($orderItem['PayTime']),
                    'out_trade_no'     => $row,
                    'external_userid'  => $externalUserId[$orderItem['UnionID']]['openid'] ?? '',
                    'total_fee'        => $orderItem['PayMoney'],
                    'payee_userid'     => $orderItem['EmpnoCode'],
                    'member_id'        => $membersId[$orderItem['EmpnoCode']]['id'] ?? 0,
                    'remark'           => $orderItem['SARemark'],
                    'total_refund_fee' => $orderItem['total_refund_fee'] ?? 0,
                    'contact_info'     => json_encode([
                        'name'    => $orderItem['Consignee'],
                        'phone'   => $orderItem['ConsigneeMobile'],
                        'address' => $orderItem['ConsigneeAddress']
                    ]),
                    'delivery_time'    => $orderItem['DeliverTime'] ? strtotime($orderItem['DeliverTime']) : 0,
                    'delivery_type'    => $orderItem['DeliverType'] == '自提' ? 1 : 2,
                    'sales_order_no'   => $orderItem['SaleNo'],
                    'logistics_no'     => $orderItem['DeliverExpressNO'],
                    'logistics_com'    => $this->companyNo($orderItem['DeliverCompanyNo']),
                    'logistics_name'   => $orderItem['DeliverCompany'],
                    'logistics_track'  => '',
                    'sign_at'          => 0,
                    'tickets'          => $orderItem['SmallTicketImgs']?$this->ticketsToCos($orderItem['SmallTicketImgs']):'',
                    'created_at'       => strtotime($orderItem['PayTime'])
                ];
            }

            //退款信息
            $orderRefundAddArr = [];
            foreach ($orderRefund as $refunds) {
                $refunds             = collect($refunds)->toArray();
                $orderRefundAddArr[] = [
                    'id'               => idGenerate(),
                    'out_trade_no'     => $refunds['OrderNO'],
                    'out_refund_no'    => $refunds['RefundNumber'],
                    'external_userid'  => $externalUserId[$refunds['UnionID']]['openid'] ?? '', //用户企微id
                    'refund_userid'    => $refunds['EmpnoCode'],//导购员工号
                    'refund_member_id' => $membersId[$refunds['EmpnoCode']]['id'] ?? 0,//导购用户id
                    'refund_comment'   => $refunds['RefundRemark'],
                    'refund_reqtime'   => strtotime($refunds['WXRefundTime']),
                    'refund_time'      => strtotime($refunds['WXRefundTime']),
                    'refund_status'    => $this->setRefundStatus($refunds['RefundStatus']),
                    'refund_fee'       => $refunds['RefundMoney'],
                    'created_at'       => strtotime($refunds['RefundTime'])
                ];
            }

            //插入
            if(count($orderAddArr)){
                ExternalPayOrder::query()->setConn($enterpriseId)->insert($orderAddArr);
            }

            if(count($orderRefundAddArr)){
                ExternalPayRefundOrder::query()->setConn($enterpriseId)->insert($orderRefundAddArr);
            }


            //dd($orderAddArr, $orderRefundAddArr);
        }

        if($lab == false){
            dd("完成");
        }

    }

    /**
     * @param $refundStatus
     *
     * @return int
     */
    public function setTradeState($refundStatus)
    {
        if ($refundStatus) {
            return 3;
        } else {
            return 1;
        }
    }

    /**
     * @param $status
     *
     * @return int|mixed
     */
    public function setOrderState($status)
    {
        switch ($status) {
            case 2:
                $status = 1;
                break;
            case 3:
                $status = 2;
                break;
            case 4:
                $status = 3;
                break;
        }
        return $status;
    }

    /**
     * 退款状态
     * @param $status
     *  //一楼状态   未退款 = 0,退款中 = 1,退款驳回 = 2,部分退款 = 3,全额退款 = 4
     * //我们的 退款状态。0:已申请退款；1:退款处理中；2:退款成功；3:退款关闭；4:退款异常；5:审批中；6:审批失败；7:审批取消',
     * @return int|mixed
     */
    public function setRefundStatus($status)
    {
        switch ($status) {
            case 0:
                $status = 0;
                break;
            case 1:
                $status = 1;
                break;
            case 2:
                $status = 6;
                break;
            case 3:
            case 4:
                $status = 2;
                break;
            default:
                $status = 0;
                break;
        }
        return $status;
    }

    /**
     * @param $com
     *
     * @return string
     */
    public function companyNo($com)
    {
        switch ($com) {
            case 'SFEXPRESS':
                $com = 'sf';
                break;
            case 'HQKY':
                $com = 'huaqikuaiyun';
                break;
            default:
                $com = '';
        }
        return $com;
    }

    /**
     * @param $img
     *
     * @return mixed|string
     */
    public function ticketsToCos($img)
    {
        try {
            if($img){
                $img = json_decode($img,true);
                if(count($img)){
                    $img = CommonService::uploadImageUrl($img[0]);
                }else{
                    return '';
                }
            }
            return $img;
        }catch (\Exception $exception){
            return '';
        }
    }

}
