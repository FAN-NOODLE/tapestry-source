<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\SwBase;

/**
 * Customer Search
 * 根据unionId查询会员信息
 */
class CustomerSearchInfoByUnionId extends SwBase
{
    public $errStr = 'query member info error';
    public $api = '/profile/search/v2';


    private $unionId;
    private $pageNum = 0;
    private $pageSize = 1;

    public function __construct(string $unionId)
    {
        parent::__construct();
        $this->unionId = $unionId;
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
                    'key' => 'union_id',
                    'value' => $this->unionId,
                    'operationType' => 'EQUALS',
                    'valueType' => 'STRING',
                    //                    'operationType'      => 'OR',
                    //                    'subQueryConditions' =>
                    //                        [
                    //                            [
                    //                                'key'           => 'union_id',
                    //                                'value'         => $this->unionId,
                    //                                'operationType' => 'EQUALS',
                    //                                'valueType'     => 'STRING',
                    //                            ],
                    //                        ]
                ]
        ];

    }
}