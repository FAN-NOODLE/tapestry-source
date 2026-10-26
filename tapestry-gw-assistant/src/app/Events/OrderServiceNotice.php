<?php

namespace App\Events;

class OrderServiceNotice extends Event
{
    /**
     * 企业ID
     *
     * @var int
     */
    private $enterpriseId = 0;

    /**
     * 服务人员ID
     *
     * @var int
     */
    private $workerId = 0;

    /**
     * 订单编号
     *
     * @var int
     */
    private $orderNo = 0;

    /**
     * OrderServiceNotice constructor.
     *
     * @param int $enterpriseId
     * @param int $workerId
     * @param int $orderNo
     */
    public function __construct(int $enterpriseId, int $workerId, int $orderNo)
    {
        $this->enterpriseId = $enterpriseId;
        $this->workerId = $workerId;
        $this->orderNo = $orderNo;
    }


    /**
     *
     * 导购ID
     *
     * @return int
     */
    public function getWorkerId()
    {
        return $this->workerId;
    }

    /**
     * 订单编号
     *
     * @return int
     */
    public function getOrderNo()
    {
        return $this->orderNo;
    }

    /**
     * 企业ID
     *
     * @return int
     */
    public function getEnterpriseId()
    {
        return $this->enterpriseId;
    }
}
