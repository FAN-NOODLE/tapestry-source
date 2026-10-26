<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 16:38
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Cache;

/**
 * 店员离岗复岗
 * Class EditStatusController
 * @package App\Http\Controllers\App\V1\Member
 */
class EditStatusController extends Controller
{
    /**
     * 店员离岗复岗
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request)
    {
        // 参数校验
        $rules   = [
            'store_id'      => 'required|integer',
            'id'            => 'required|integer',
            'worker_status' => 'required|in:' . implode(',', array_keys(Constant::MEMBER_WORK_STATUS))
        ];
        $message = [
            'store_id.required'      => '店铺ID参数不能为空',
            'store_id.integer'       => '店铺ID参数类型有误',
            'id.required'            => '员工ID参数类型有误',
            'id.integer'             => '员工ID参数类型有误',
            'worker_status.required' => '员工状态不能为空',
            'worker_status.in'       => '员工状态变更不正确'
        ];
        $input   = $this->validate($request, $rules, $message);
        try {
            $storeWorker = new StoreWorker();
            $input['enterprise_id'] = Auth::user()['enterprise_id'];
            $input['domain']=$this->loginAdminEnterprise['domain'];
            $storeWorker->jobEdit($input);
            $detail = $storeWorker->detail($input);
            if (!empty($detail)) {
                Cache::forget($detail['mobile']);
            }
            return success();
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
    }
}
