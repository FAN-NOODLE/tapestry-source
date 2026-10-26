<?php

namespace App\Console\Commands;

use App\Services\Crm\CrmService;
use Illuminate\Console\Command;

class TestCrmSdk extends Command
{
    /**
     * --说明--
     *  根据enterprise_id识别调用crm接口
     *  根据action_id识别调用crm接口的方法 action_id:1,查询会员;2,查询优惠卷;3,查询订单
     *  php artisan test-crm-sdk --enterprise_id=$v --union_id=$v --action_id=$v
     */

    /**
     * @var string
     */
    protected $signature = 'test-crm-sdk {--enterprise_id=} {--union_id=} {--action_id=}';

    /**
     * @var string
     */
    protected $description = '测试调用对接的crm接口';


    public function handle(): array
    {
        $this->info("测试调用对接的crm接口");

        $input['enterprise_id'] = $enterpriseId = intval($this->option('enterprise_id'));
        $input['union_id']      = strval($this->option('union_id'));
        $actionId               = intval($this->option('action_id'));

//        dd($enterpriseId,$input,$actionId);
        // 1,查询会员;2,查询优惠卷;3,查询订单
        switch ($actionId) {
            case 1:
                $result = (new CrmService($enterpriseId))->getMemberInfo($input);
                break;
            case 2:
                $input['status']       = 0;
                $input['channel_type'] = '';
                $input['source']       = 1;
                $result                = (new CrmService($enterpriseId))->getCustomerCoupon($input);

                break;
            case 3:
                $input['startDateTime'] = 0;
                $input['endDateTime']   = time() * 1000;

                $result = (new CrmService($enterpriseId))->getCustomerPurchase($input);

                break;
            default:
                $result = [];

        }

        $this->info("测试调用对接的crm接口返回:result" . json_encode($result, JSON_UNESCAPED_UNICODE));

        return $result;

    }

}