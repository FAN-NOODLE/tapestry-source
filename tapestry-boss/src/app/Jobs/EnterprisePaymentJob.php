<?php


namespace App\Jobs;


use App\Services\EnterprisePaymentService\EnterprisePaymentService;
use Illuminate\Support\Facades\Log;

class EnterprisePaymentJob extends Job
{
    protected $data = [];

    public function __construct($data)
    {
        $this->data = $data;
    }

    public function handle()
    {
        $orderStatus = $this->data['order_status'];
        unset($this->data['order_status']);
        // order_status = 1 增加记录， order_status = 5 申请退款成功，处理失效
        switch ($orderStatus) {
            case 1:
                if (empty($this->data['buy_day']) || empty($this->data['buy_number']) || empty($this->data['payment_time']) || empty($this->data['enterprise_id']) || empty($this->data['payment_order_sn']) || empty($this->data['use_time']) || empty($this->data['fail_time'])) {
                    return;
                }
                $this->add();
                break;
            case 5:
                if (empty($this->data['enterprise_id']) || empty($this->data['payment_order_sn'])) {
                    return;
                }
                $this->failAccount();
                break;
            default:
                break;
        }
    }

    /**
     * 智慧导购新增账号
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function add()
    {
        // 组织数据
        $time = time();
        $data = [
            'enterprise_id'     => $this->data['enterprise_id'],
            'buy_day'           => $this->data['buy_day'],
            'buy_number'        => $this->data['buy_number'],
            'payment_time'      => $this->data['payment_time'],
            'payment_order_sn'  => $this->data['payment_order_sn'],
            'payment_amount'    => $this->data['payment_amount'],
            'use_time'          => $this->data['use_time'],
            'fail_time'         => $this->data['fail_time'],
            'server_charge'     => 0,
            'payment_type'      => 0,
            'payment_source'    => '企业微信',
            'payment_status'    => 0,
            'payment_category'  => 2
        ];
        $result = (new EnterprisePaymentService)->add($data);
        if ($result !== true) {
            Log::info('create error by job', $data);
        }
    }

    /**
     * 智慧导购退款处理账号失效
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function failAccount()
    {
        // 组织数据
        $data = [
            'enterprise_id'     => $this->data['enterprise_id'],
            'payment_order_sn'  => $this->data['payment_order_sn'],
        ];
        $result = (new EnterprisePaymentService)->failAccount($data);
        if ($result !== true) {
            Log::info('create error by job', $data);
        }
    }
}