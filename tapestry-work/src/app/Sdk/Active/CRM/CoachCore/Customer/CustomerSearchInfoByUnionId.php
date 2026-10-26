<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Search
 * 根据unionId查询会员信息
 */
class CustomerSearchInfoByUnionId extends Base
{
    public $errStr = 'query member info error';
    public $api = '/v2/cdp/profile/search/v2';


    private $unionId;
    private $pageNum = 1;
    private $pageSize = 1;

    public function __construct(string $unionId)
    {
        parent::__construct();
        $this->unionId = $unionId;
        $this->setQueryParam();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'pageNum' => $this->pageNum,
            'pageSize' => $this->pageSize,
            'queryCondition' =>
                [
                    'key' => 'union_id',
                    'value' => $this->unionId,
                    'operationType' => 'EQUALS',
                    'valueType' => 'STRING',
                ],

        ];

    }
}