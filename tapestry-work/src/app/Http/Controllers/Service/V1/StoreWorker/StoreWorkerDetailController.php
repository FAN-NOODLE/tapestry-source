<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;

class StoreWorkerDetailController extends BaseController
{
    public function index()
    {
        $rules = [
            'contact_qr_code_id' => 'required',
            'mem_uid'            => 'string',
            'enterprise_id'      => 'integer'
        ];
        $message = [
            'contact_qr_code_id.required' => '联系我二维码id不能为空',
            'mem_uid.string'              => '员工参数有误',
            'enterprise_id.integer'       => '企业ID有误'
        ];
        $this->validate($this->request, $rules, $message);
        $params = $this->request->only([
            'contact_qr_code_id',
            'mem_uid',
            'enterprise_id'
        ]);
        $params['contact_qr_code_id'] = $params['contact_qr_code_id'] ?? 0;
        $params['mem_uid'] = $params['mem_uid'] ?? '';
        $params['enterprise_id'] = $params['enterprise_id'] ?? 0;
        $data = [];
        $storeWorkerInfo = StoreWorker::query()->where(['contact_qr_code_id' => $params['contact_qr_code_id']])->first();
        if (empty($storeWorkerInfo) && !empty($params['mem_uid'])) {
            // 优先获取在岗里最新更新的导购信息
            $storeWorkerInfo = StoreWorker::query()
                ->when(!empty($params['enterprise_id']), function ($query) use ($params) {
                    return $query->where('enterprise_id', $params['enterprise_id']);
                })
                ->where('mem_uid', $params['mem_uid'])
                ->where('worker_status', StoreWorker::IN_SERVICE)
                ->orderBy('updated_at', 'DESC')
                ->first();
        }
        if ($storeWorkerInfo) {
            $data = $storeWorkerInfo->toArray();
        }
        return success($data);
    }
}
