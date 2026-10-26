<?php

use PHPUnit\Framework\TestCase;

class RefundTest extends TestCase
{
    public function testMachineRun()
    {
        $refundMachine = new \App\StateMachine\Refund\RefundMachine();


        $data['event'] = \App\StateMachine\Refund\RefundMachineConfig::EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND;
        $data['order_id'] = '2359';
        $data['reason'] = '无理由';
        $res = $refundMachine->run($data);

        $this->assertTrue($res);
    }
}
