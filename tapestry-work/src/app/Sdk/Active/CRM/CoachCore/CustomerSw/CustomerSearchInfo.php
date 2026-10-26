<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\SwBase;

/**
 * Customer Search
 * 通用查询会员信息
 */
class CustomerSearchInfo extends SwBase
{
    public $errStr = 'query member info error';
    public $api = '/profile/search/v2';


    private $stringCode;
    private $stringValue;
    private $pageNum = 0;
    private $pageSize = 1;

    public function __construct(string $stringCode, string $stringValue)
    {
        parent::__construct();
        $this->stringCode = $stringCode;
        $this->stringValue = $stringValue;
        $this->api = $this->prefix . $this->api;
        $this->setQueryParam();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'pageNum' => $this->pageNum,
            'pageSize' => $this->pageSize,
            'queryCondition' =>
                [
                    'key' => $this->stringCode,
                    'value' => $this->stringValue,
                    'operationType' => 'EQUALS',
                    'valueType' => 'STRING',
                ],
        ];
    }
}