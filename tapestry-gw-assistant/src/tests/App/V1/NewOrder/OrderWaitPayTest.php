<?php

namespace Tests\App\V1\NewOrder;

use App\Services\Order\OrderStateConfig;
use ExampleTest;
class OrderWaitPayTest extends ExampleTest
{
    /** @test 创建订单 */
    public function test_create_order()
    {
        $response = $this->common_create_order();
        $this->assertEquals(200, $response->status());
    }

    /** @test 用户取消订单 */

    public function test_customer_cancel(){
        $response = $this->common_create_order();
        $order_id = $this->getOrderId($response->getContent());
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_CANCEL
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 自动取消订单 */

    public function test_auto_cancel(){
        $response = $this->common_create_order();
        $order_id = $this->getOrderId($response->getContent());
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_AUTO_CANCEL
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 支付成功需确认 */

    public function test_pay_success_confirm(){
        $response = $this->common_create_order();
        $order_id = $this->getOrderId($response->getContent());
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_WX_PAY_SUCCESS_CONFIRM
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 支付成功无需确认 */

    public function test_pay_success_no_confirm(){
        $response = $this->common_create_order();
        $order_id = $this->getOrderId($response->getContent());
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_WX_PAY_SUCCESS_NO_CONFIRM
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }


}
