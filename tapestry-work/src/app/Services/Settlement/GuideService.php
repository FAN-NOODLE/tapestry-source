<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/16 0016
 * Time: 14:25
 */

namespace App\Services\Settlement;


use App\Library\Uc\Member;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\Enterprise;
use TheSeer\Tokenizer\Exception;

class GuideService
{
    private $enterpriseId;

    private $settlement;

    private $promoteGuideId;

    private $serviceGuideId;

    private $promoteMemUid;

    private $serviceMemUid;

    private $allPromoteGuideInfo;

    private $allServiceGuideInfo;

    private $memberSdk;

    private $showStore = false;

    private $qyDomain;

    private $allMemInfo = [];

    const STORE_VERSION = 2;

    public function __construct($enterpriseId, $version = "")
    {
        $this->enterpriseId = $enterpriseId;
        $this->memberSdk = new Member();
        $this->qyDomain = (new Enterprise())->detail(['enterprise_id' => $this->enterpriseId])['data']['domain'];
        $this->showStore = $version == self::STORE_VERSION ? true : false;
    }

    public function setUser(array $settlement)
    {
        $this->settlement = $settlement;
        $this->promoteGuideId = $settlement['promotion_worker_id'];
        $this->serviceGuideId = $settlement['service_shopping_guid_id'];
        $this->promoteMemUid = $this->allPromoteGuideInfo[$this->promoteGuideId]['mem_uid'] ?? 0;
        $this->serviceMemUid = $this->allServiceGuideInfo[$this->serviceGuideId]['mem_uid'] ?? 0;
    }

    public function setGuideIds($promoteGuideIds, $serviceGuideIds)
    {
        $this->allPromoteGuideInfo = [];
        $this->allServiceGuideInfo = [];
        $this->setAllPromoteGuideInfo($promoteGuideIds);

        $this->setAllServiceGuideInfo($serviceGuideIds);
        $this->setDeptInfo();
    }

    private function setDeptInfo()
    {
        $promoteMemUids = array_column($this->allPromoteGuideInfo, 'mem_uid');
        $serviceMemuids = array_column($this->allPromoteGuideInfo, 'mem_uid');
        $memuids = array_values(array_unique(array_merge($promoteMemUids, $serviceMemuids)));

        $params = [
            'memUids'  => $memuids,
            'pageNum'  => 1,
            'pageSize' => count($memuids)
        ];// 一次取够
        $memberInfo = $this->memberSdk->ignorePermissionsList($this->qyDomain, null, $params);
        foreach ($memberInfo['list'] as $info) {
            $this->allMemInfo[$info['memUid']] = $info;
        }
    }

    private function setAllPromoteGuideInfo($promoteGuideIds)
    {

        $promoteWorkerInfo = StoreWorkerService::list(
            [
                'enterprise_id' => $this->enterpriseId,
                'ids'           => $promoteGuideIds,
                'size'          => count($promoteGuideIds)
            ]
        );
        $this->allPromoteGuideInfo = array_column($promoteWorkerInfo['list'], null, 'id');

    }

    private function setAllServiceGuideInfo($serviceGuideIds)
    {
        $serviceWorkerInfo = StoreWorkerService::list(
            [
                'enterprise_id' => $this->enterpriseId,
                'ids'           => $serviceGuideIds,
                'size'          => count($serviceGuideIds)
            ]
        );

        $this->allServiceGuideInfo = array_column($serviceWorkerInfo['list'], null, 'id');
    }

    private function getPromoteGuideName()
    {
        return $this->allPromoteGuideInfo[$this->promoteGuideId]['name'] ?? "";
    }

    private function getPromoteGuideDept()
    {
        if ($this->showStore && $this->promoteGuideId) {
            $promoteMemDept = $this->getPromoteMemDept();
            $promoteMemStore = $this->allPromoteGuideInfo[$this->promoteGuideId]['store_name'] ?? "";

            return $promoteMemStore . '/' . $promoteMemDept;
        }
        return $this->getPromoteMemDept();
    }

    private function getPromoteGuideAvatar()
    {
        return $this->allPromoteGuideInfo[$this->promoteGuideId]['avatar'] ?? "";
    }

    private function getServiceGuideName()
    {
        return $this->allServiceGuideInfo[$this->serviceGuideId]['name'] ?? "";
    }

    private function getServiceGuideDept()
    {

        if ($this->showStore && $this->serviceGuideId) {
            $serviceMemDept = $this->getServiceMemDept();
            $serviceMemStore = $this->allServiceGuideInfo[$this->serviceGuideId]['store_name'] ?? "";
            return $serviceMemDept ? $serviceMemStore . '/' . $serviceMemDept : $serviceMemStore;
        }
        return $this->getServiceMemDept();
    }

    private function getServiceGuideAvatar()
    {
        return $this->allServiceGuideInfo[$this->serviceGuideId]['avatar'] ?? "";
    }

    private function getPromoteMemDept()
    {

        return $this->allMemInfo[$this->promoteMemUid]['dpName'][0]['dpName'] ?? "";
    }

    private function getServiceMemDept()
    {
        return $this->allMemInfo[$this->serviceMemUid]['dpName'][0]['dpName'] ?? "";
    }

    public function formatSingle()
    {
        return [
            "id"                       => $this->settlement['id'],
            "settlement_status"        => $this->settlement['settlement_status'],
            "order_no"                 => $this->settlement['order_no'],
            "order_status"             => $this->settlement['order_status'],
            "order_pay_time"           => $this->settlement['order_pay_time'],
            "promotion_commission"     => $this->settlement['promotion_commission'],
            "promotion_worker_id"      => $this->settlement['promotion_worker_id'],
            "service_shopping_guid_id" => $this->settlement['service_shopping_guid_id'],
            "service_commission"       => $this->serviceGuideId ? $this->settlement['service_commission'] : 0,
            "service_name"             => $this->getServiceGuideName(),
            "service_dept"             => $this->getServiceGuideDept(),
            "service_avatar"           => $this->getServiceGuideAvatar(),
            "promotion_name"           => $this->getPromoteGuideName(),
            "promotion_dept"           => $this->getPromoteGuideDept(),
            "promotion_avatar"         => $this->getPromoteGuideAvatar(),
        ];
    }

}