<?php

namespace App\Sdk\Active\CRM\CoachCore\Auth;

use App\Sdk\Active\CRM\CoachCore\Base;
use App\Sdk\Active\CRM\CoachCore\Config;
use App\Sdk\Active\CRM\CoachCore\ConfigSW;
use App\Sdk\Active\CrmHelper;

class RefreshAccessToken extends Base
{
    public $errStr = 'refresh AccessToken error';
    public $api = '/v2/oauth/token?';

//    protected $type = 'form_data';

    protected $tokenUrl = true;
    private $refreshToken;
    /**
     * @var bool|mixed
     */
    private $enterpriseConfig;

    public function __construct($enterpriseId, $refreshToken)
    {
        $this->refreshToken = $refreshToken;
        $this->enterpriseConfig = CrmHelper::getEnterpriseConfig($enterpriseId);
        if (isset($this->enterpriseConfig['name'])) {

            if ($this->enterpriseConfig['name'] == 'coach') {
                $this->api .= 'grant_type=refresh_token&refresh_token=' . $this->refreshToken . '&client_id=' . Config::getClientId();

            } else {
                $this->api .= 'grant_type=refresh_token&client_id=' . ConfigSW::getClientId() . '&refresh_token=' . $this->refreshToken;
            }
        }
        parent::__construct();

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
        $this->queryParam = [
            "grant_type" => "refresh_token",
            "refresh_token" => $this->refreshToken,
            "client_id" => Config::getClientId()
        ];
    }


}