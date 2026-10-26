<?php

namespace Tests\App\V1\NewOrder;

use App\Services\Order\OrderStateConfig;
use ExampleTest;
class OrderWaitDeliverTest extends ExampleTest
{
    /** @test 创建订单 */
    public function test_create_order()
    {
        $response = $this->common_create_order();
        $this->assertEquals(200, $response->status());
    }

    /** @test 拼团失败 */

    public function test_group_fail_refund(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_GROUP_FAIL_REFUND
        ];
        $response = $this->call('POST', ExampleTest::SERVICE_BASE_URL . ExampleTest::SERVICE_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test open通知发货订单 */

    public function test_open_send_out_order(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_OPEN_SEND_OUT_ORDER,
            'logistic_identity'=> 'YTO',
            'logistic_no'=> $this->faker->randomNumber(5, true)
        ];
        $response = $this->call('POST', ExampleTest::SERVICE_BASE_URL . ExampleTest::SERVICE_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 用户申请退款 */

    public function test_customer_apply_refund_only_money(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_APPLY_REFUND,
            'reason'=> '不想要了',
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

}
