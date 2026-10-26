<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\SwBase;

/**
 * Customer Purchase Search
 * 根据customerId查询会员购买记录
 */
class CustomerPurchaseSearch extends SwBase
{
    public $errStr = 'query member info error';
    public $api = '/purchase/history';


    private $pageNum = 0;
    private $pageSize = 10;
    private $customerId;
    private $startDateTime;
    private $endDateTime;

    public function __construct(string $customerId, $startDateTime, $endDateTime)
    {
        parent::__construct();
        $this->customerId = $customerId;
        $this->startDateTime = $startDateTime;
        $this->endDateTime = $endDateTime;
        $this->api = $this->prefix . $this->api;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'pageNum' => $this->pageNum,
            'pageSize' => $this->pageSize,
            'customerId' => $this->customerId,
            'startDateTime' => $this->startDateTime,
            'endDateTime' => $this->endDateTime,
        ];

    }
}