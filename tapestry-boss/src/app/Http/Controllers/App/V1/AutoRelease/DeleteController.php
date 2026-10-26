<?php

namespace App\Http\Controllers\App\V1\AutoRelease;

use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use App\Models\EnterpriseGrayscaleReleaseLog;
use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Validation\ValidationException;

class DeleteController extends BaseController
{
    /**
     * @param Request $request
     * @return JsonResponse
     * @throws ValidationException
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

        // 当前记录还没撤销, 不可删除
        if ($log->can_cancel == EnterpriseGrayscaleReleaseLog::CAN_CANCEL_YES) {
            notice(300508);
        }

        $log->forceDelete();

        return success(['data' => '']);
    }
}

