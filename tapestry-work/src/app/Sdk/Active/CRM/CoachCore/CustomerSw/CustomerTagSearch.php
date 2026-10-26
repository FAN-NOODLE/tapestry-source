<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\SwBase;

/**
 * Customer Tag Search
 * 根据customerId查询会员标签记录
 */
class CustomerTagSearch extends SwBase
{
    public $errStr = 'query member info error';
    public $api = '/customer-tag/searchTags';


    private $pageNum = 0;
    private $pageSize = 10;
    private $customerId;
    private $source;        // 查询数据源: wecom

    public function __construct(string $customerId, string $source)
    {
        parent::__construct();
        $this->customerId = $customerId;
        $this->source = $source;
        $this->api = $this->prefix . $this->api;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [

            'pageNum' => $this->pageNum,
            'pageSize' => $this->pageSize,
            'customerId' => $this->customerId,
            'source' => $this->source,


        ];

    }
}