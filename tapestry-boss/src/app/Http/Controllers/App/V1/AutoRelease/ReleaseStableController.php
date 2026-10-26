<?php

namespace App\Http\Controllers\App\V1\AutoRelease;

use App\Common\RedisKey;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class ReleaseStableController extends BaseController
{
    /**
     * 新版
     */
    const NEW_TAG = 1;

    /**
     * 旧版
     */
    const OLD_TAG = 2;

    /**
     * 设置发布版本Tag
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {

        $this->paramsValidate();

        $type = $this->request->post('type');
        $releaseTag = $this->request->post('release_tag');

        // type 1 新版 2 旧版
        if (!in_array($type, [self::NEW_TAG, self::OLD_TAG])) {
            error(300515);
        }

        // 新版本
        if ($type == self::NEW_TAG) {
            $ids = Enterprise::query()
                ->where('version', '>', Enterprise::NORMAL_VERSION)
                ->where('status', '=', Enterprise::STATUS_ENABLE)
                ->pluck('id')
                ->toArray();

            EnterpriseThirdPartAccount::query()
                ->whereIn('enterprise_id', $ids)
                ->where('enterprise_type', '=', EnterpriseThirdPartAccount::RELEASE_STABLE)
                ->update([
                    'release_tag' => $releaseTag
                ]);
        }

        // 旧版本
        if ($type == self::OLD_TAG) {
            $ids = Enterprise::query()
                ->where('version', '=', Enterprise::NORMAL_VERSION)
                ->where('status', '=', Enterprise::STATUS_ENABLE)
                ->pluck('id')
                ->toArray();

            EnterpriseThirdPartAccount::query()
                ->whereIn('enterprise_id', $ids)
                ->where('enterprise_type', '=', EnterpriseThirdPartAccount::RELEASE_STABLE)
                ->update([
                    'release_tag' => $releaseTag
                ]);
        }
        $authKeys = [];
        foreach ($ids as $id) {
            $authKeys[] = RedisKey::MINI_PROGRAM_AUTH_INFO_BY_ENTERPRISE_ID . $id;
        }
        if (!empty($authKeys)) {
            delRedis($authKeys);
        }
        return success([
            'data' => ''
        ]);
    }

    /**
     * 数据验证
     *
     * @throws ValidationException
     */
    private function paramsValidate()
    {
        $rules = [
            'type'        => 'required',
            'release_tag' => 'required',
        ];

        $errMessage = [
            'type.required'        => '类型不能为空',
            'release_tag.required' => 'tag号不能为空',
        ];

        $this->validate($this->request, $rules, $errMessage);
    }
}
