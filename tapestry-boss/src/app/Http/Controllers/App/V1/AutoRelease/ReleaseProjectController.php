<?php

namespace App\Http\Controllers\App\V1\AutoRelease;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class ReleaseProjectController extends BaseController
{
    /**
     * 发布项目
     * @param Request $request
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(Request $request): JsonResponse
    {
        $this->paramsValidate();

        $tag = $request->post('release_tag');
        $enterpriseId = $request->post('enterprise_id');

        $account = EnterpriseThirdPartAccount::query()
            ->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('verify_status', '=', EnterpriseThirdPartAccount::VERIFIED)
            ->where('enterprise_id', '=', $enterpriseId)
            ->first();

        if (empty($account)) {
            notice(300502);
        }
        if ($account->release_branch == $tag) {
            notice(300504);
        }
        if ($account->enterprise_type == EnterpriseThirdPartAccount::RELEASE_STABLE) {
            notice(300512);
        }

        $account->update(['release_tag' => $tag]);

        return success(['data' => '']);
    }

    /**
     * 数据验证
     *
     * @throws ValidationException
     */
    private function paramsValidate()
    {
        $rules = [
            'release_tag'   => 'required',
            'enterprise_id' => 'required|integer',
        ];
        $errMessage = [
            'release_tag.required'   => 'tag不能为空',
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer'  => '企业ID必须为正整数',
        ];

        $this->validate($this->request, $rules, $errMessage);
    }
}