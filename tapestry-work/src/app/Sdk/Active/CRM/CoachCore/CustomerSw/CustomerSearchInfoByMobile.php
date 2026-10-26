<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\SwBase;

/**
 * Customer Search
 * 根据手机号查询会员详细信息
 */
class CustomerSearchInfoByMobile extends SwBase
{
    public $errStr = 'query member info error';
    public $api = '/profile/search/v2';


    private $mobile;
    private $pageNum = 0;
    private $pageSize = 1;

    public function __construct(string $mobile)
    {
        parent::__construct();
        $this->mobile = $mobile;
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
                    'key' => 'mobile',
                    'value' => $this->mobile,
                    'operationType' => 'EQUALS',
                    'valueType' => 'STRING',
                    //                    'operationType'      => 'OR',
                    //                    'subQueryConditions' =>
                    //                        [
                    //                            [
                    //                                'key'           => 'mobile',
                    //                                'value'         => $this->mobile,
                    //                                'operationType' => 'EQUALS',
                    //                                'valueType'     => 'STRING',
                    //                            ],
                    //                        ]
                ]
        ];

    }
}