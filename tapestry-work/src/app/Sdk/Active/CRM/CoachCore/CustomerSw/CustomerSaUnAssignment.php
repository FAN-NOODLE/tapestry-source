<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Tag Search
 * 删除顾客
 */
class CustomerSaUnAssignment extends Base
{
    public $errStr = 'send SW saUnAssignment info error';
    public $api = '/v2/cdp/customer/saUnAssignment';

    public $prodApi = '/cdp-api-sw/customer/saUnAssignment';

    private $saAccount;
    private $unionId;
    private $unBindTime;

    public function __construct(array $params)
    {
        parent::__construct();
        if (config('app.env') == 'prod') {
            $this->api = $this->prodApi;
        }
        $this->saAccount = $params['saAccount'];
        $this->unionId = $params['unionId'];
        $this->unBindTime = $params['unbindTime'];
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'sa_account' => $this->saAccount,
            'union_id' => $this->unionId,
            'unbind_time' => $this->unBindTime,
        ];

    }
}