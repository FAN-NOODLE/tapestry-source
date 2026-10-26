<?php

namespace App\Http\Controllers\Service\V1\WithdrawaSetting;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\WithdrawalSetting;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     *
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->get('enterprise_id', 0);

        if (!$enterpriseId) {
            throw new Error(700000, '企业ID不能为空');
        }

        $setInfo = WithdrawalSetting::query()->where('enterprise_id', $enterpriseId)
            ->first(['enterprise_id','transfer_no','type']);

        if(empty($setInfo)){
            return success([]);

        }
        $setInfo = $setInfo->toArray();
        $setInfo['type'] = array_map('intval', explode(',',$setInfo['type']));
        return success($setInfo);
    }
}
