<?php

namespace App\Services\ShareChain;

/**
 * 实体类
 *
 * 依赖外部不同策略的实体类
 */
class Substance
{
    /**
     * 策略实例
     * @var object
     */
    private $_strategy;

    /**
     * Substance constructor.
     * @param StrategyInterface $strategy
     */
    public function __construct(StrategyInterface $strategy)
    {
        $this->_strategy = $strategy;
    }

    /**
     * 员工信息初始化，该接口的主要作用是记录初始分享者（员工）的memUid和worker_id
     *
     * @param $enterpriseId
     * @param $userId
     * @param $memUid
     * @param $workerId
     * @param $type
     */
    public function staffInitial($enterpriseId, $userId, $memUid, $workerId, $type)
    {
        $this->_strategy->staffInitial($enterpriseId, $userId, $memUid, $workerId, $type);
    }

    /**
     * 获取当前用户分享链路
     *
     * @param $enterpriseId
     * @param $userId
     * @param $type
     */
    public function getStaffChain($enterpriseId, $userId, $type)
    {
        $this->_strategy->getStaffChain($enterpriseId, $userId, $type);
    }

    /**
     * 用户分享关系绑定，包含是否抢客等判断
     *
     * @param $enterpriseId
     * @param $userId
     * @param $parentUserId
     * @param $workerId
     * @param $type
     */
    public function staffBinding($enterpriseId, $userId, $parentUserId, $workerId, $type)
    {
        $this->_strategy->staffBinding($enterpriseId, $userId, $parentUserId, $workerId, $type);
    }

    /**
     * 解绑逻辑，员工离职
     * 该解绑会影响该员工的所有被分享者
     *
     * @param int $enterpriseId
     * @param array $memUids
     * @param array $workerIds
     * @param int $type
     * @return void
     */
    public function unbindingLeave(int $enterpriseId, array $memUids, array $workerIds, int $type)
    {
        $this->_strategy->unbindingLeave($enterpriseId, $memUids, $workerIds, $type);
    }

    /**
     * 解绑逻辑，员工变导购的情况
     * 该解绑只解绑登陆者自己的绑定关系
     *
     * @param int $enterpriseId
     * @param array $memUids
     * @param array $workerIds
     * @param int $type
     * @return void
     */
    public function unbindingChange(int $enterpriseId, array $memUids, array $workerIds, int $type)
    {
        $this->_strategy->unbindingChange($enterpriseId, $memUids, $workerIds, $type);
    }
}
