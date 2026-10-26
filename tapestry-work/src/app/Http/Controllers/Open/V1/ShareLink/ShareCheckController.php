<?php
/**
 * 合伙人分享链路 - 检查
 */

namespace App\Http\Controllers\Open\V1\ShareLink;

use App\Http\Controllers\Open\V1\BaseController;
use App\Models\Distribution;
use App\Models\PartnerEnterpriseSettings;
use App\Services\PartnerMemberService;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\Log;
use App\Services\StaffManagerService;

class ShareCheckController extends BaseController
{

    public function index()
    {
        $input = $this->getParams();

        list($storeStaffId) = StoreWorkerService::getShareWorkerIdByUserId($input['enterprise_id'], $input['user_id']);
        $returnResult = [
            'store_staff_id' => $storeStaffId ?? 0
        ];

        return $this->success($returnResult);
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {
        Log::info('分享链路检查ShareCheckController', [$this->request->input()]);

        $input = $this->request->only([
            'store_id',
            'user_id'
        ]);
        $input['enterprise_id'] = $this->enterpriseInfo['id'];
        $rules = [
            'store_id' => 'integer',
            'user_id' => 'required',
        ];
        $messages = [
            'store_id.*' => '门店id参数不合法',
            'user_id.*' => '用户id参数不合法'
        ];
        $this->validate($this->request, $rules, $messages);

        return $input;
    }
}