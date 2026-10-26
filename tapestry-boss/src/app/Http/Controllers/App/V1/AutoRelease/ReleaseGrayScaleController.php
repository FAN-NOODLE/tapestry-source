<?php

namespace App\Http\Controllers\App\V1\AutoRelease;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\EnterpriseGrayscaleReleaseLog;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Validation\ValidationException;

class ReleaseGrayScaleController extends BaseController
{
    /**
     * @param Request $request
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(Request $request): JsonResponse
    {
        $this->paramsValidate();

        $remark = $request->post('remark', '');
        $branch = $request->post('release_branch');
        $enterpriseId = $request->post('enterprise_id');

        $account = EnterpriseThirdPartAccount::query()->where([
            'enterprise_id' => $enterpriseId,
            'type'          => EnterpriseThirdPartAccount::TYPE_SHOP,
            'status'        => EnterpriseThirdPartAccount::STATUS_ACTIVE,
            'verify_status' => EnterpriseThirdPartAccount::VERIFIED
        ])->first();

        if (empty($account)) {
            notice(300502);
        }

        if ($account->release_branch == $branch) {
            notice(300504);
        }

        if ($account->enterprise_type == EnterpriseThirdPartAccount::RELEASE_PROJECT) {
            notice(300510);
        }

        // 检查是否正在灰度某个版本, 未撤销状态且允许撤销状态
        $releasing = EnterpriseGrayscaleReleaseLog::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->where('status', '=', EnterpriseGrayscaleReleaseLog::CANCEL_NO)
            ->where('can_cancel', '=', EnterpriseGrayscaleReleaseLog::CAN_CANCEL_YES)
            ->first();
        if ($releasing) {
            notice(300516, [$releasing->release_branch]);
        }

        try {
            // 开始事务
            DB::beginTransaction();

            $account->update(['release_branch' => $branch]);

            EnterpriseGrayscaleReleaseLog::query()->where('enterprise_id', '=', $enterpriseId)->update([
                'status'     => EnterpriseGrayscaleReleaseLog::CANCEL_YES,
                'can_cancel' => EnterpriseGrayscaleReleaseLog::CAN_CANCEL_NO,
            ]);

            EnterpriseGrayscaleReleaseLog::query()->create([
                'enterprise_id'  => $enterpriseId,
                'release_branch' => $branch,
                'status'         => EnterpriseGrayscaleReleaseLog::CANCEL_NO,
                'can_cancel'     => EnterpriseGrayscaleReleaseLog::CAN_CANCEL_YES,
                'remark'         => $remark
            ]);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            notice(300503);
        }

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
            'release_branch' => 'required',
            'enterprise_id'  => 'required|integer',
            'remark'         => 'max:90'
        ];

        $errMessage = [
            'release_branch.required' => '分支名称不能为空',
            'enterprise_id.required'  => '企业ID不能为空',
            'enterprise_id.integer'   => '企业ID必须为正整数',
            'remark.max'              => '备注过长',
        ];

        $this->validate($this->request, $rules, $errMessage);
    }
}