<?php

namespace Tests\App\V1\NewOrder;

use App\Services\Order\OrderStateConfig;
use ExampleTest;
class WaitDeliverAfterSaleTest extends ExampleTest
{
    /** @test 创建订单 */
    public function test_create_order()
    {
        $response = $this->common_create_order();
        $this->assertEquals(200, $response->status());
    }

    /** @test  用户取消售后*/

    public function test_customer_cancel(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER_AFTER_SALE);
        $this->common_apply_refund($order_id);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_CANCEL
        ];
        $response = $this->call('POST', ExampleTest::APP_BASE_URL . ExampleTest::APP_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 人工拒绝 */

    public function test_manual_refuse_apply_refund(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER_AFTER_SALE);
        $this->common_apply_refund($order_id);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_MANUAL_REFUSE_APPLY_REFUND
        ];
        $response = $this->call('POST', ExampleTest::SERVICE_BASE_URL . ExampleTest::SERVICE_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test open通知发货 */

    public function test_open_send_out_order(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER_AFTER_SALE);
        $this->common_apply_refund($order_id);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_OPEN_SEND_OUT_ORDER,
            'logistic_identity'=> 'YTO',
            'logistic_no'=> $this->faker->randomNumber(5, true)
        ];
        $response = $this->call('POST', ExampleTest::SERVICE_BASE_URL . ExampleTest::SERVICE_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test open通知仅退款 */

    public function test_open_agree_only_refund(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER_AFTER_SALE);
        $this->common_apply_refund($order_id);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_OPEN_AGREE_ONLY_REFUND,
        ];
        $response = $this->call('POST', ExampleTest::SERVICE_BASE_URL . ExampleTest::SERVICE_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /** @test 人工同意退款 */

    public function test_manual_agree_refund(){
        $response = $this->common_create_order();
        $order_id = $this->updateOrder($response->getContent(),OrderStateConfig::WAIT_DELIVER_AFTER_SALE);
        $this->common_apply_refund($order_id);
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_MANUAL_AGREE_REFUND,
        ];
        $response = $this->call('POST', ExampleTest::SERVICE_BASE_URL . ExampleTest::SERVICE_MACHINE_STATE_URL, $data,[],[],$this->buildHeader());
        $this->assertEquals(200, $response->status());
    }
}
