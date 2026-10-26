<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Tag Search
 * 添加顾客
 */
class CustomerSaAssignment extends Base
{
    public $errStr = 'query member info error';
    public $api    = '/coachwechatws/ws/customer/saAssignment';


    private $saAccount;
    private $firstbinddate;
    private $latestbinddate;
    private $unionId;
    private $storecode;

    public function __construct(array $params)
    {
        parent::__construct();
        $this->saAccount      = $params['saAccount'];
        $this->firstbinddate  = $params['firstbinddate'];
        $this->latestbinddate = $params['latestbinddate'];
        $this->unionId        = $params['unionId'];
        $this->storecode      = $params['storecode'];
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'saAccount'      => $this->saAccount,
            'firstbinddate'  => $this->firstbinddate,
            'latestbinddate' => $this->latestbinddate,
            'unionId'        => $this->unionId,
            'storecode'      => $this->storecode,
        ];

    }
}