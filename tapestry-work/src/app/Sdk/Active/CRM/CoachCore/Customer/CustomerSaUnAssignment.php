<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Tag Search
 * 删除顾客
 */
class CustomerSaUnAssignment extends Base
{
    public $errStr = 'query member info error';
    public $api    = '/coachwechatws/ws/customer/saUnAssignment';


    private $saAccount;
    private $unionId;

    public function __construct(array $params)
    {
        parent::__construct();
        $this->saAccount      = $params['saAccount'];
        $this->unionId        = $params['unionId'];
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'saAccount'      => $this->saAccount,
            'unionId'        => $this->unionId,
        ];

    }
}