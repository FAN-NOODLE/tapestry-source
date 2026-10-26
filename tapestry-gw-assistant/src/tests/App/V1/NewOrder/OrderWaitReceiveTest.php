<?php

namespace Tests\App\V1\NewOrder;

use App\Services\Order\OrderStateConfig;
use ExampleTest;
class OrderWaitReceiveTest extends ExampleTest
{
    /** @test 创建订单 */
    public function test_create_order()
    {
        $response = $this->common_create_order();
        $this->assertEquals(200, $response->status());
    }

    /** @test 用户签收 */

    public function test_customer_sign(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_RECEIVE);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_SIGN
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 用户申请退货 */

    public function test_open_send_out_order(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_RECEIVE);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_APPLY_RETURN,
            'reason'=> '不想要了',
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 用户申请仅退款 */

    public function test_customer_apply_refund_only_money(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_RECEIVE);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_APPLY_REFUND,
            'reason'=> '不想要了',
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 超期自动签收 */
    public function test_auto_sign(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_RECEIVE);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_AUTO_SIGN
        ];
        $response = $this->call('POST', ExampleTest::SERVICE_BASE_URL . ExampleTest::SERVICE_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

}
