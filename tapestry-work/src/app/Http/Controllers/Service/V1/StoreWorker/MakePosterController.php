<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Jobs\StoreWorkerGenerateQrCodeJob;

class MakePosterController extends BaseController
{

    /**
     * 批量生成导购海报
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input        = $this->getParams();
        $posterParams = [];
        foreach ($input['poster'] as $poster) {
            if (!isset($posterParams[$poster['store_id']])) {
                $posterParams[$poster['store_id']]['store_id']      = $poster['store_id'];
                $posterParams[$poster['store_id']]['enterprise_id'] = $input['enterprise_id'];
                $posterParams[$poster['store_id']]['domain']        = $input['domain'];
                $posterParams[$poster['store_id']]['operator_id']   = $input['operator_id'];
                $posterParams[$poster['store_id']]['operator_name'] = $input['operator_name'];
            }
            $posterParams[$poster['store_id']]['worker_ids'][] = $poster['worker_id'];
        }
        foreach ($posterParams as $posterItem) {
            $this->dispatch(new StoreWorkerGenerateQrCodeJob($posterItem));
        }
        return success();
    }

    /**
     * 获取参数
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     *
     */
    private function getParams(): array
    {
        $rules = [
            'operator_id'        => 'integer',
            'domain'             => 'required|max:32',
            'poster.*.worker_id' => 'required|integer',
            'poster.*.store_id'  => 'required|integer',
            'enterprise_id'      => 'required|integer',
            'operator_name'      => 'max:125'
        ];

        $messages = [
            'enterprise_id.*'      => '企业ID不正确',
            'domain.*'             => '企业domain不正确',
            'poster.*.worker_id.*' => '导购ID不正确',
            'poster.*.store_id.*'  => '导购所属门店不正确',
        ];

        $input                  = $this->validate($this->request, $rules, $messages);
        $input['operator_id']   = $input['operator_id'] ?? 0;
        $input['operator_name'] = $input['operator_name'] ?? '系统管理员';
        return $input;
    }
}
