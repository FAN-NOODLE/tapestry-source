<?php


namespace App\Library;


use CrsSdk\MicroAuthority\WX;

class AccessToken
{
    /**
     * @param $enterpriseId
     * @param string $appId
     * @param string $appSecret
     * @return mixed
     * @throws \Exception
     */
    public function getAccessToken($enterpriseId, $appId = '', $appSecret = '')
    {
        $token = (new WX())->accessToken(['enterprise_id' => $enterpriseId]);

        return $token['access_token'];
    }
}
