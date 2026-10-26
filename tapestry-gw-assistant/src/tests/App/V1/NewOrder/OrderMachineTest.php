<?php

namespace Tests\App\V1\NewOrder;

use App\Services\Order\OrderMachine;
use App\Services\Order\OrderStateConfig;
use ExampleTest;
class OrderMachineTest extends ExampleTest
{
    /** @test run */

    public function test_run(){
        $response = $this->common_create_order();
        $order_id = $this->getOrderId($response->getContent());
        $orderMachine = new OrderMachine();
        $data = [
            'id'=>$order_id,
            'event'=>OrderStateConfig::EVENT_CUSTOMER_CANCEL
        ];
        $result = $orderMachine->run($data,OrderStateConfig::CONSOLE_OPERATOR);
        $this->assertCount(2, $result);
    }
}
