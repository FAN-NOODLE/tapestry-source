<?php


namespace App\Http\Controllers\Service\V1\StaffBuyManager;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\StaffRightRequest;
use App\Jobs\SyncContactCodeJob;
use App\Models\StaffSetting;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class UpdateController extends BaseController
{
    /**
     * 修改允许参与内购的员工或部门
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->post();

        $input['app_identifier'] = $input['app_identifier'] ?? 'contact';

        // 验证和添加一样
        StaffRightRequest::add($input);

        try {
            DB::beginTransaction();

            try {
                // 处理员工内购分享绑定关系确定或解除
                (new StaffManagerService())->staffRightUnbindStaffShare($input);
            } catch (\Exception $exception) {
                Log::error('修改内购员工--绑定关系：' . $exception->getMessage());
            }

            // 软删除旧数据
            (new StaffManagerService())->deleteStaffRight($input['enterprise_id']);

            // 添加新数据
            (new StaffManagerService())->createForStaffRight($input);

            // 更新数据
            // (new StaffManagerService())->updateForStaffRight($input);

            (new StaffManagerService())->createOrUpdateForStaffSetting($input);

            // 清除缓存
            $guideSettingKey = redisKey('ENTERPRISE_GUIDE_SETTING', ['enterprise_id' => $input['enterprise_id'], 'type' => StaffSetting::STAFF_TYPE]);
            delRedis($guideSettingKey);

            //为内购员工自建联系人二维码
//            $this->dispatch(new SyncContactCodeJob($input));

            DB::commit();
        } catch (\Exception $e) {

            DB::rollBack();
            Log::info('内购员工数据修改错误' . json_encode($e->getMessage()), $input);
            Log::info('file:' . $e->getFile());
            Log::info('line:' . $e->getLine());

            notice($e->getCode() ?? 720000, [$e->getMessage()]);
        }

        return success(['data' => []]);
    }
}
