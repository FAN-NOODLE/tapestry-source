<?php

namespace App\Http\Controllers\Service\V1\WithdrawaSetting;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\WithdrawalAutomaticReview;
use Illuminate\Http\Request;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Log;

//use Illuminate\Support\Facades\Log;

class AutoDetailController extends BaseController
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
        $memType = $request->get('mem_type', WithdrawalAutomaticReview::MEM_TYPE_WORKER);

        if (!$enterpriseId) {
            throw new Error(700000, '企业ID不能为空');
        }

        $setInfo = WithdrawalAutomaticReview::query()
            ->where('mem_type', $memType)
            ->where('enterprise_id', $enterpriseId)
            ->first(['auto_status','enterprise_id','auto_at','last_at']);

        if(empty($setInfo)){
            return success(['auto_status' => '','enterprise_id' => $enterpriseId,'auto_at' => '','last_at' => '']);
        }
        return success($setInfo->toArray());
    }
}
