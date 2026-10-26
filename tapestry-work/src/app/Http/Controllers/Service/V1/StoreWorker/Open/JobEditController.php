<?php
/*
 * @Author: your name
 * @Date: 2020-03-18 17:45:55
 * @LastEditTime: 2020-03-24 18:27:50
 * @LastEditors: Please set LastEditors
 * @Description: 员工离岗在岗操作
 * @FilePath: \gw-shop\src\app\Http\Controllers\Service\V1\Store\Worker\DetailController.php
 */

namespace App\Http\Controllers\Service\V1\StoreWorker\Open;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Models\UserStoreWorkerRelation;
use App\Services\StoreWorkerService;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;
use \Exception;

class JobEditController extends BaseController
{

    /**
     * @return array|void
     * @throws Exception
     */
    public function index()
    {
        $rules = [
            'mem_uid'       => 'required',
            'enterprise_id' => 'required',
            'store_no'      => 'required'
        ];

        $message = [
            'mem_uid.required'       => 'mem_uid参数必填',
            'enterprise_id.required' => 'enterprise_id参数必填',
            'store_no.required'      => 'store_no参数必填',
        ];
        $this->validate($this->request, $rules, $message);
        $input = $this->request->input();
        switchAppHost(['enterprise_id' => $input['enterprise_id']]);
        $storeInfo = (new Store())->detail([
            'store_no' => $input['store_no'],
            'enterprise_id' => $input['enterprise_id']
        ]);
        $storeId = $storeInfo['id'];
        $worker = StoreWorker::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->where('store_id', $storeId)
            ->where('mem_uid', $input['mem_uid'])
            ->first(['worker_status', 'is_boss', 'id']);
        if (empty($worker)) {
            notice(2300006, [], [], 400);
        }
        // 门店店长无法修改在职状态
        if ($worker->is_boss == StoreWorker::BOSS_YES) {
            notice(2300009, [], [], 400);
        }
        $workerId = $worker->id;
        try {
            // 获取需要修改的员工状态
            $params = [
                'worker_status' => StoreWorker::IN_SERVICE,
                'entry_time'    => time()
            ];
            if ($worker->worker_status == StoreWorker::IN_SERVICE) {
                $params = [
                    'worker_status' => StoreWorker::OUT_SERVICE,
                    'out_time'      => time()
                ];
                // 离职软删除员工关系表
                UserStoreWorkerRelation::query()
                    ->where([
                        'enterprise_id' => $input['enterprise_id'],
                        'worker_id'     => $workerId,
                        'store_id'      => $storeId
                    ])->delete();
            }
            StoreWorker::query()
                ->where('enterprise_id', $input['enterprise_id'])
                ->where('store_id', $storeId)
                ->where('id', $workerId)
                ->where('worker_status', $worker->worker_status)
                ->update($params);
            StoreWorkerService::removeWorkerRedis($workerId);
            $result = ['worker_status' => $params['worker_status']];
        } catch (Exception $e) {
            Log::info('修改员工失败' . $e->getMessage(), $input);
            notice(2300002, [], [], 400);
        }
        return success($result);
    }
}
