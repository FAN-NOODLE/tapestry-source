<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Point History
 * 根据customerId查询会员积分记录
 */
class CustomerPointHistory extends Base
{
    public $errStr = 'query member info error';
    public $api = '/v2/cdp/profile/point/history';


    private $pageNum = 1;
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