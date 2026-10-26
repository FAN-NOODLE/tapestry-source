<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Services\ShareChain\StaffShareStrategy;
use App\Services\ShareChain\Substance;
use App\Services\StoreWorkerService;
use CrsSdk\Work\Member;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use \Exception;

class JobEditController extends BaseController
{

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $rules   = [
            'id'            => 'required|integer',
            'enterprise_id' => 'required',
            'store_id'      => 'required|integer',
            'worker_status' => 'required|in:' . StoreWorker::IN_SERVICE . ',' . StoreWorker::OUT_SERVICE,
            'domain'        => 'max:32'
        ];
        $message = [
            'id.integer'             => 'id参数必须为数字',
            'id.required'            => 'id参数必填',
            'store_id.required'      => 'store_id参数必填',
            'worker_status.required' => '状态参数必填',
            'domain.max'             => 'domain参数不正确'
        ];

        $input  = $this->validate($this->request, $rules, $message);
        $worker = StoreWorkerService::detail([
            'enterprise_id' => $input['enterprise_id'],
            'store_id'      => $input['store_id'],
            'worker_id'     => $input['id']
        ]);
        if (empty($worker)) {
            return success();
        }
        // 门店店长无法修改在职状态
        if ($worker['is_boss'] == StoreWorker::BOSS_YES) {
            notice(2300009);
        }
        if ($worker['worker_status'] == $input['worker_status']) {
            $string = $input['worker_status'] == StoreWorker::IN_SERVICE ? '在职' : '离职';
            notice(2300012, [$string]);
        }
        if ($input['worker_status'] == StoreWorker::IN_SERVICE) {
            // 调用sdk
            $condition = [
                'qy_domain'      => $input['domain'],
                'app_identifier' => 'contact',
                'params'         => ['memUid' => $worker['mem_uid'] ?? '']
            ];
            $detail    = (new Member())->detail($condition);
            if (!$detail) {
                notice(2300014);
            }
        }
        try {
            DB::beginTransaction();
            // 获取需要修改的员工状态
            $params = [
                'worker_status' => StoreWorker::IN_SERVICE,
                'entry_time'    => time()
            ];
            if ($input['worker_status'] == StoreWorker::OUT_SERVICE) {
                $params = [
                    'worker_status' => StoreWorker::OUT_SERVICE,
                    'out_time'      => time()
                ];
                try {
                    // 更新用户分享记录
                    $staffShareStrategy = new Substance(new StaffShareStrategy());
                    $staffShareStrategy->unbindingLeave($input['enterprise_id'], [], [$input['id']], Constant::GUIDE_TYPE);
                } catch (Exception $e) {
                }
                // 离职软删除员工关系表
                StoreWorkerService::delUserStoreWorkerRelation([
                    'enterprise_id' => $input['enterprise_id'],
                    'worker_id'     => $input['id'],
                    'store_id'      => $input['store_id']
                ]);
            }
            StoreWorkerService::update([
                'enterprise_id' => $input['enterprise_id'],
                'store_id'      => $input['store_id'],
                'id'            => $input['id'],
                'worker_status' => $worker['worker_status']
            ], $params);
            // 删除用户导购关联的缓存数据
            $redisKey = redisKey('USER_STORE_GUIDE_INFO', ['user_id' => $worker['user_id']]);
            delRedis([$redisKey]);
            StoreWorkerService::removeWorkerRedis($input['id']);
            // 管理员变更
            StoreWorkerService::changeAdmin($worker, $params['worker_status']);
            DB::commit();
        } catch (Exception $e) {
            Log::info('修改员工状态失败' . $e->getMessage(), $input);
            DB::rollBack();
            notice(2300002);
        }
        return success();
    }


}
