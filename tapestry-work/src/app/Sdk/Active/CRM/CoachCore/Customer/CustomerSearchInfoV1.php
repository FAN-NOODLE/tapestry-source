<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Search
 * 通用查询会员信息
 */
class CustomerSearchInfoV1 extends Base
{
    public $errStr = 'query member info error';
    public $api = '/v2/cdp/profile/search';


    private $stringCode;
    private $stringValue;
    private $pageNum = 1;
    private $pageSize = 1;

    public function __construct(string $stringCode, string $stringValue)
    {
        parent::__construct();
        $this->stringCode = $stringCode;
        $this->stringValue = $stringValue;
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