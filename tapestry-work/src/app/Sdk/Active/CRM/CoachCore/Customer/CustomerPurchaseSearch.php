<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Purchase Search
 * 根据customerId查询会员购买记录
 */
class CustomerPurchaseSearch extends Base
{
    public $errStr = 'query member info error';
    public $api = '/v2/cdp/purchase/history';

    // 备注 pageNumber 0开始

    private $pageNum;
    private $pageSize;
    private $customerId;
    private $startDateTime;
    private $endDateTime;

    public function __construct(string $customerId, $startDateTime, $endDateTime, $page = 0, $size = 10)
    {
        parent::__construct();
        $this->customerId = $customerId;
        $this->startDateTime = $startDateTime;
        $this->endDateTime = $endDateTime;
        $this->pageNum = $page;
        $this->pageSize = $size;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'pageNumber' => $this->pageNum,
            'pageSize' => $this->pageSize,
            'customerId' => $this->customerId,
            'startDateTime' => $this->startDateTime,
            'endDateTime' => $this->endDateTime,
        ];

    }
}