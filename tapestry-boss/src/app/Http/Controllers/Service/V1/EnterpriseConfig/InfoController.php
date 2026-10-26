<?php

namespace App\Http\Controllers\Service\V1\EnterpriseConfig;

use App\Models\EnterpriseConfig;
use App\Http\Controllers\Service\V1\BaseController;

class InfoController extends BaseController
{
    /**
     * 获取配置详情
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required|integer',
        ]);
        $input = $this->request->only(['enterprise_id', 'type', 'key']);

        $data = EnterpriseConfig::query()
            ->select(['type', 'key', 'value', 'name'])
            ->where('enterprise_id', '=', $input['enterprise_id'])
            ->when(isset($input['type']), function ($query) use ($input) {
                return $query->where('type', '=', $input['type']);
            })
            ->when(isset($input['key']), function ($query) use ($input) {
                if(is_array($input['key'])){
                    return $query->whereIn('key', $input['key']);
                }else{
                    return $query->where('key', '=', $input['key']);
                }
            })
            ->get()
            ->toArray();

        return success(['data' => $data]);
    }
}
