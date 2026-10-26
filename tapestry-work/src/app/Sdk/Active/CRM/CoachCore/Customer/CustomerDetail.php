<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * 会员详情
 */
class CustomerDetail extends Base
{
    public $errStr = 'query member info error';
    public $api    = '/weclientws/ws/rest/customer/showCustomerDetail';


    private $storeCode;
    private $regionCode;
    private $curPage;
    private $pageSize;

    private $saAccount;
    private $memberId;
    // Update comment based on flag value Example: 1 (Update)  0: Read Only , 1: Update
    private $updateFlag;
    // 如果updateFlag的值为"1"系统需要更新注释部分
    private $comment;


    public function __construct(array $params)
    {
        parent::__construct();

        $this->storeCode  = $params['storeCode'] ?? '';
        $this->curPage    = $params['curPage'] ?? 1;
        $this->pageSize   = $params['pageSize'] ?? 10;
        $this->regionCode = !empty($params['regionCode']) ? $params['regionCode'] : 'CN';

        $this->saAccount  = $params['saAccount'] ?? '';
        $this->memberId   = $params['memberId'] ?? '';
        $this->updateFlag = $params['updateFlag'] ?? 0;
        $this->comment    = $params['comment'] ?? '';


        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        // 门店 storeCode=""  curPage="1" pageSize="100"
        $this->queryParam = '<?xml version="1.0" encoding="UTF-8" standalone="no"?>' .
            '<request storeCode="' . $this->storeCode . '" regionCode="' . $this->regionCode . '" curPage="' . $this->curPage . '" pageSize="' . $this->pageSize . '">' .
            '<customerDetailReq memberId="' . $this->memberId . '" updateFlag="' . $this->updateFlag . '" comment="' . $this->comment . '"/>' .
            '</request>';


    }
}