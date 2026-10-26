<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Services\ShareChain\StaffShareStrategy;
use App\Services\ShareChain\Substance;
use Illuminate\Support\Facades\Validator;

class StaffShareUnbindingController
{
    /**
     * 解绑逻辑，员工变导购的情况
     * 该解绑只解绑登陆者自己的绑定关系
     *
     * @return array|mixed
     */
    public function unbindingChange()
    {
        Validator::make(request()->all(), [
            'enterprise_id' => 'required',
            'mem_uids'      => 'array',
            'worker_ids'    => 'array',
            'type'          => 'required',
        ], [
            'enterprise_id.required' => '请确认企业',
            'mem_uids.*'             => 'mem_uids类型为数组',
            'worker_ids.*'           => 'worker_ids类型为数组',
            'type.*'                 => 'type必填，1内购2导购',
        ])->validate();

        $params = request()->only(['enterprise_id', 'mem_uids', 'worker_ids', 'type']);

        $staffShareStrategy = new Substance(new StaffShareStrategy);
        $staffShareStrategy->unbindingChange($params['enterprise_id'], $params['mem_uids'] ?? [], $params['worker_ids'] ?? [], $params['type']);

        return success();
    }

    /**
     * 解绑逻辑，员工离职
     * 该解绑会影响该员工的所有被分享者
     * @return array
     */
    public function unbindingLeave()
    {
        Validator::make(request()->all(), [
            'enterprise_id' => 'required',
            'mem_uids'      => 'array',
            'worker_ids'    => 'array',
            'type'          => 'required',
        ], [
            'enterprise_id.required' => '请确认企业',
            'mem_uids.*'             => 'mem_uids类型为数组',
            'worker_ids.*'           => 'worker_ids类型为数组',
            'type.*'                 => 'type必填，1内购2导购',
        ])->validate();

        $params = request()->only(['enterprise_id', 'mem_uids', 'worker_ids', 'type']);

        $staffShareStrategy = new Substance(new StaffShareStrategy);
        $staffShareStrategy->unbindingLeave($params['enterprise_id'], $params['mem_uids'] ?? [], $params['worker_ids'] ?? [], $params['type']);

        return success();
    }
}
