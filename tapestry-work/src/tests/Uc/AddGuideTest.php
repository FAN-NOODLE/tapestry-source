<?php

class AddGuideTest extends \PHPUnit\Framework\TestCase
{

    public function testBatchAdd()
    {
        $this->markTestSkipped();
        $appIdentifier = null;
        $qyDomain      = "B1A503057F0000016FD160F1F22E246C";
        $params        = [
            'mobiles'  => [18710839146],
            'pageNum'  => 1,
            'pageSize' => 1,
        ];

        $result = (new \App\Library\Uc\Member())->ignorePermissionsList(
            $qyDomain, $appIdentifier, $params);
        $info   = $result['list'][0];


        $enterpriseId = "3357415263026304";
        $guideMemUids       = \App\Models\StoreWorker::query()
            ->where('enterprise_id', $enterpriseId)
            ->pluck("mem_uid")
            ->toArray();

        if(in_array($info['memUid'],$guideMemUids)){
            $this->assertTrue(false,"已经添加过了");
            return;
        }

        \App\Library\Util\Logger::debug("白信息",$info);
        $members      = [
            [
                "mem_uid"   => $info['memUid'],
                "user_name" => $info["memUsername"],
                "mobile"    => $info["memMobile"],
                "avatar"    => $info["memFace"],
                "user_id"   => $info["memUserid"]
            ]
        ];
        $enterpriseId = "3357415263026304";
        $rules        = [
            'enterprise_id' => $enterpriseId,
            'member'        => $members,
            'type'          => \App\Common\Constant::STORE_WORKER_TYPE_ENTERPRISE
        ];
        \App\Services\StoreWorkerService::addGuide($rules);
        $this->assertTrue(true,"断言");
    }

    public function testAddGuide()
    {

        $data = (new \CrsSdk\Boss\Enterprise())->listsOfficial([]);

        foreach ($data['data'] as $enterprise){
            $this->updateSingleEp($enterprise['id'],$enterprise['domain']);
        }

    }

    private function updateSingleEp($epId,$qyDomain)
    {
        $guideMemUids       = \App\Models\StoreWorker::query()
            ->where('enterprise_id', $epId)
            ->pluck("mem_uid")
            ->toArray();

        $appIdentifier = null;
        $pageNum       = 1;

        do {
            $params = [
                'memSubscribeStatus'  => 1,
                'pageNum'  => $pageNum,
                'pageSize' => 500,
            ];
            $result = (new \App\Library\Uc\Member())->ignorePermissionsList(
                $qyDomain, $appIdentifier, $params);
            if ($result['list']) {
                $members = [];
                foreach ($result['list'] as $info){
                    if(!in_array($info['memUid'],$guideMemUids)){
                        $members[]      =
                            [
                                "mem_uid"   => $info['memUid'],
                                "user_name" => $info["memUsername"],
                                "mobile"    => $info["memMobile"],
                                "avatar"    => $info["memFace"],
                                "user_id"   => $info["memUserid"]
                            ];

                    }
                }
                if($members){
                    $enterpriseId = "3357415263026304";
                    $rules        = [
                        'enterprise_id' => $enterpriseId,
                        'member'        => $members,
                        'type'          => \App\Common\Constant::STORE_WORKER_TYPE_ENTERPRISE
                    ];
                    \App\Library\Util\Logger::debug("需要添加的导购",$rules);
//                    try{
//                        \App\Services\StoreWorkerService::addGuide($rules);
//                    }catch (Throwable $e){
//
//                    }
                }
            } else {
                break;
            }
            $pageNum++;
        } while (true);

    }
}