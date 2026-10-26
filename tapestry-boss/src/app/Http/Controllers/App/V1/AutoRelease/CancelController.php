<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2019/11/8
 * Time: 15:58
 */
namespace App\Http\Controllers\App\V1\AutoRelease;

use App\Models\EnterpriseGrayscaleReleaseLog;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\EnterpriseLogService;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class CancelController extends BaseController
{
    /**
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $this->validate($request, [
            'id' => 'integer|required',
        ], [
            'id.integer'  => 'id必须为正整数',
            'id.required' => 'id不能为空',
        ]);

        $id = $request->post('id');

        $log = EnterpriseGrayscaleReleaseLog::query()->where('id', '=', $id)->first();

        if (empty($log)) {
            notice(300505);
        }

        if ($log->can_cancel == EnterpriseGrayscaleReleaseLog::CAN_CANCEL_NO) {
            notice(300514);
        }

        $account = EnterpriseThirdPartAccount::query()->where([
            'enterprise_id' => $log->enterprise_id,
            'type'          => EnterpriseThirdPartAccount::TYPE_SHOP,
            'status'        => EnterpriseThirdPartAccount::STATUS_ACTIVE,
            'verify_status' => EnterpriseThirdPartAccount::VERIFIED
        ])->first();

        // 小程序已撤销灰度, 记录还没变更
        if (empty($account->release_branch)) {
            $log->forceDelete();
            notice(300506);
        }

        $account->update([
            'release_branch' => ''
        ]);

        $log->update([
            'status'     => EnterpriseGrayscaleReleaseLog::CANCEL_YES,
            'can_cancel' => EnterpriseGrayscaleReleaseLog::CAN_CANCEL_NO,
            'cancel_at'  => time()
        ]);

        return success(['data' => '']);
    }
}