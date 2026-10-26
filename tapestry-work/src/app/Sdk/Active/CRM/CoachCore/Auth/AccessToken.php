<?php

namespace App\Sdk\Active\CRM\CoachCore\Auth;

class AccessToken
{
    private $access_token;
    private $token_type;
    private $refresh_token;
    private $expires_in;

    public function __construct(array $token)
    {
        $this->access_token  = $token['access_token'];
        $this->token_type    = $token['token_type'];
        $this->refresh_token = $token['refresh_token'];
        $this->expires_in    = $token['expires_in'];
    }

    public function getAccessToken()
    {
        return $this->access_token;
    }

    public function getTokenType()
    {
        return $this->token_type;
    }

    public function getRefreshToken()
    {
        return $this->refresh_token;
    }

    public function getExpireIn()
    {
        return $this->expires_in;
    }

}