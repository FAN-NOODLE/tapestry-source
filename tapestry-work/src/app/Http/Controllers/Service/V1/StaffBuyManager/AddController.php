<?php


namespace App\Http\Controllers\Service\V1\StaffBuyManager;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\StaffRightRequest;
use App\Jobs\SyncContactCodeJob;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class AddController extends BaseController
{
    /**
     * 添加允许参加内购的员工
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->post();

        $input['app_identifier'] = $input['app_identifier'] ?? 'contact';

        // 验证参数
        StaffRightRequest::add($input);

        try {
            DB::beginTransaction();

            try {
                // 处理员工内购分享绑定关系确定或解除
                (new StaffManagerService())->staffRightUnbindStaffShare($input);
            } catch (\Exception $exception) {
                Log::error('添加内购员工--绑定关系：' . $exception->getMessage());
            }

            // 插入数据(允许参加内购的员工)
            $insertStaffRight = (new StaffManagerService())->createForStaffRight($input);

            if (!$insertStaffRight) {
                notice(720001);
            }

            // 插入数据(内购员工设置)
            $insertStaffSetting = (new StaffManagerService())->createOrUpdateForStaffSetting($input);

            if (!$insertStaffSetting) {
                notice(720007);
            }

            //为内购员工自建联系人二维码
//            dispatch(new SyncContactCodeJob($input));

            DB::commit();
        } catch (\Exception $e) {

            DB::rollBack();
            Log::info('内购员工数据新增错误' . json_encode($e->getMessage()), $input);
            notice($e->getCode(), [$e->getMessage()]);
        }

        return success(['data' => []]);
    }
}
