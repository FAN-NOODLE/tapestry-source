<?php

namespace App\Sdk\Active\CRM\CoachCore\CustomerSw;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Tag Search
 * 添加顾客
 */
class CustomerSaAssignment extends Base
{
    public $errStr = 'send Sw saAssignment info error';
    public $api = '/v2/cdp/customer/saAssignment';

    public $prodApi = '/cdp-api-sw/customer/saAssignment';

    private $saAccount;
    private $bindTime;
    private $unionId;
    private $storeCode;
    private $source;

    public function __construct(array $params)
    {
        parent::__construct();
        if (config('app.env') == 'prod') {
            $this->api = $this->prodApi;
        }
        $this->saAccount = $params['saAccount'];
        $this->bindTime = $params['bindTime'];
        $this->unionId = $params['unionId'];
        $this->storeCode = $params['storecode'];
        $this->source = $params['source'];
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'sa_account' => $this->saAccount,
            'bind_time' => $this->bindTime,
            'union_id' => $this->unionId,
            'store_code' => $this->storeCode,
            'source' => $this->source,
        ];

    }
}