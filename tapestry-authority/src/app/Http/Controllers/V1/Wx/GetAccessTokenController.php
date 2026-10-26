<?php

namespace App\Http\Controllers\V1\Wx;

use App\Common\RedisKey;
use App\Http\Controllers\V1\BaseController;
use App\Models\Enterprise;
use App\Traits\RedisMutexLock;
use App\WxServices\AccessTokenWxService;
use Illuminate\Http\Request;

class GetAccessTokenController extends BaseController
{
    use RedisMutexLock;

    public function index(Request $request)
    {
        $hash = $request->input('hash', '');
        $type = $request->input('type', 1);
        $refresh = $request->input('refresh', false);
        $enterpriseId = $request->input('enterprise_id', 0);

        if (empty($enterpriseId) && !empty($hash)) {
            $enterpriseId = Enterprise::query()->where('hash', $hash)->value('id') ?? 0;
        }

        $isLock = $this->lock(sprintf(RedisKey::ACCESS_TOKEN_LOCK, $enterpriseId, $type));
        if ($isLock === false) {
            error(100404);
        }

        if (!empty($enterpriseId)) {
            $accessToken = (new AccessTokenWxService())->getAccessToken($enterpriseId, $type, $refresh);
        }

        $this->release(sprintf(RedisKey::ACCESS_TOKEN_LOCK, $enterpriseId, $type));

        return success([
            'access_token' => $accessToken ?? ''
        ]);
    }
}
