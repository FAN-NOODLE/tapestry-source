<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\SwBase;

/**
 * Customer Point History
 * 根据customerId查询会员积分记录
 */
class CustomerPointHistory extends SwBase
{
    public $errStr = 'query member info error';
    public $api = '/profile/point/history';


    private $pageNum = 0;
    private $pageSize = 10;
    private $customerId;
    private $startTime;
    private $endTime;

    public function __construct(string $customerId, $startTime, $endTime)
    {
        parent::__construct();
        $this->customerId = $customerId;
        $this->startTime = $startTime;
        $this->endTime = $endTime;
        $this->api = $this->prefix . $this->api;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'pageNum' => $this->pageNum,
            'pageSize' => $this->pageSize,
            'customerId' => $this->customerId,
            'startTime' => $this->startTime,
            'endTime' => $this->endTime,
        ];

    }
}