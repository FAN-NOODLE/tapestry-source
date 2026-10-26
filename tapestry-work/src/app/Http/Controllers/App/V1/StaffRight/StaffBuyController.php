<?php


namespace App\Http\Controllers\App\V1\StaffRight;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\StaffManagerService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\StaffRight;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class StaffBuyController extends BaseController
{
    /**
     * 判断员工是否有内购权限
     * @api /app/v1/staff-right/is-right
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {

        $enterpriseId = $this->enterpriseId;
        $qyDomain = $this->qyDomain;
        $memUid = $request->input('mem_uid', '');

        if (empty($memUid)) {
            notice(720005);
        }

        try {

            $res = (new StaffManagerService())->staffIsBuy($enterpriseId, $memUid, $qyDomain);

        } catch (\Exception $exception) {

            notice(720006);

        }

        return success(['data' => ['is_right' => $res]]);
    }

}