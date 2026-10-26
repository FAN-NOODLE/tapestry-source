<?php

namespace App\Sdk\Active\CRM\CoachCore\Auth;

use App\Sdk\Active\CRM\CoachCore\Base;
use App\Sdk\Active\CRM\CoachCore\Config;
use App\Sdk\Active\CRM\CoachCore\ConfigSW;
use App\Sdk\Active\CrmHelper;

/**
 * Class GetAccessToken
 */
class GetAccessToken extends Base
{

    public $errStr = 'get AccessToken error';
    public $api = '/v2/oauth/token?';

    protected $tokenUrl = true;

    private $enterpriseConfig;

    public function __construct($enterpriseId)
    {
        parent::__construct();
        $this->enterpriseConfig = CrmHelper::getEnterpriseConfig($enterpriseId);
        if (isset($this->enterpriseConfig['name'])) {

            if ($this->enterpriseConfig['name'] == 'coach') {
                $this->api .= 'grant_type=password&username=' . Config::getUsername() . '&password=' . Config::getPassword() . '&client_id=' . Config::getClientId() . '&client_secret=' . Config::getClientSecret();

            } else {
                $this->api .= 'grant_type=password&username=' . ConfigSW::getUsername() . '&password=' . ConfigSW::getPassword() . '&client_id=' . ConfigSW::getClientId() . '&client_secret=' . ConfigSW::getClientSecret();

            }
        }

    }

    public function setHeader($enterpriseId, $xCountryCode = 'CN')
    {
        if ($this->enterpriseConfig['name'] == 'coach') {
            $this->header = [
                'auth' => [Config::getClientId(), Config::getClientSecret()]
            ];
        } else {
            $this->header = [
                'auth' => [ConfigSW::getClientId(), ConfigSW::getClientSecret()]
            ];
        }

    }

    protected function setQueryParam()
    {

    }


}