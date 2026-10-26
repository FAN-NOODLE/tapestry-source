<?php

namespace App\WxServices;

use App\Library\AccessToken;
use App\Services\EnterpriseThirdPartAccountService;
use App\Traits\OpenPlatformTrait;
use App\Traits\RequestRecord;
use Psr\SimpleCache\InvalidArgumentException;

class AccessTokenWxService
{
    use OpenPlatformTrait, RequestRecord;

    public function getAccessToken($enterpriseId, $type = 1, $refresh = false)
    {
        $accessToken = $type == 1 ? $this->getAuthorizeAccessToken($enterpriseId, $refresh) : $this->getAccessTokenByAppSecret($enterpriseId, $refresh);

        return $accessToken ?? '';
    }

    public function getAuthorizeAccessToken($enterpriseId, $refresh)
    {
        try {
            $miniProgram = $this->miniProgram($enterpriseId);

            $tokenArr = $miniProgram->access_token->getToken($refresh);

            $this->wxRecord($tokenArr, [
                'url'  => 'cgi-bin/token?grant_type=client_credential',
                'file' => __FILE__,
                'line' => __LINE__,
                'data' => [
                    'eid'    => $enterpriseId,
                    'result' => $tokenArr
                ],
            ]);

            $accessToken = $tokenArr['authorizer_access_token'] ?? '';

        } catch (\Exception $e) {

            $this->wxRecord(['errmsg' => $e->getMessage()], [
                'url'  => 'cgi-bin/token?grant_type=client_credential',
                'file' => __FILE__,
                'line' => __LINE__,
                'data' => [
                    'eid' => $enterpriseId,
                ],
            ], $e);

            $accessToken = '';
        }

        return $accessToken;
    }

    /**
     * @param $enterpriseId
     * @param $refresh
     * @return string
     * @throws InvalidArgumentException
     */
    public function getAccessTokenByAppSecret($enterpriseId, $refresh)
    {
        $miniProgramInfo = (new EnterpriseThirdPartAccountService())->getMiniProgram($enterpriseId);

        if (empty($miniProgramInfo)) {
            error(700021);
        }

        $appId = $miniProgramInfo['app_id'];
        $appSecret = $miniProgramInfo['app_secret'];

        return (new AccessToken())->getAccessToken($appId, $appSecret, $refresh);
    }
}
