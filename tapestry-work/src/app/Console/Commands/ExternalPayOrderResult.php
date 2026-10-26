<?php

namespace App\Console\Commands;

use App\Library\Express\Express;
use App\Models\ScrmEnterprise\ExternalPayOrder;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class ExternalPayOrderResult extends Command
{
    /**
     * @var string
     */
    protected $signature = 'auto:external_pay_order_result';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '企业外部订单物流状态更新';

    /**
     * @return void
     */
    public function handle()
    {
        //查询企业
        $enterprise = (new Enterprise())->listNoPage();
        $enterprise = collect($enterprise['data'])->pluck('id')->all();

        //发货后7天没确认收货
        $time = time() - 7 * 24 * 3600;
        foreach ($enterprise as $row) {
            $data = ExternalPayOrder::query()
                ->setConn($row)
                ->where([
                    'order_state'   => ExternalPayOrder::ORDER_STATE_DELIVER,
                    'delivery_type' => ExternalPayOrder::DELIVERY_TYPE_ON,
                ])
                ->where('delivery_time', '<', $time)
                ->get();
            //查询到订单信息处理快递
            foreach ($data as $item) {

                $receiverPhone = $item->contact_info['phone'] ?? '';
                $receiverPhone = $receiverPhone?substr($receiverPhone,-4):'';
                $data          = (new Express())->queryExpress($item->logistics_com, $item->logistics_no, '', $receiverPhone);

                if ($data['data']) {
                    $arr['logistics_track'] = $data['data']['list'];
                    $item->logistics_track = json_encode($arr['logistics_track']);
                    if ($data['data']['status_detail']) {
                        $item->sign_at = strtotime($data['data']['list'][0]['datetime']);
                        $item->order_state = ExternalPayOrder::ORDER_STATE_SUCCESS;
                    }
                    $item->save();
                }
            }
        }
    }
}