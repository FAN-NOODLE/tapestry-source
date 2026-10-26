<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Search
 * 根据手机号查询会员详细信息
 */
class CustomerSearchInfoByMobile extends Base
{
    public $errStr = 'query member info error';
    public $api = '/v2/cdp/profile/search/v2';


    private $mobile;
    private $pageNum = 1;
    private $pageSize = 1;

    public function __construct(string $mobile)
    {
        parent::__construct();
        $this->mobile = $mobile;
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
                ],

        ];

    }
}