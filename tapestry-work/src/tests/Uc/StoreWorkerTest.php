<?php

class StoreWorkerTest extends \PHPUnit\Framework\TestCase
{
    public function testGuide()
    {
        $enterpriseId = "3357415263026304";
        $memUid       = \App\Models\StoreWorker::query()
            ->where('enterprise_id', $enterpriseId)
            ->pluck("mem_uid")
            ->toArray();
        \App\Library\Util\Logger::debug("导购memuid",$memUid);
        $this->assertTrue(count($memUid) > 0, "没有导购");

    }

}