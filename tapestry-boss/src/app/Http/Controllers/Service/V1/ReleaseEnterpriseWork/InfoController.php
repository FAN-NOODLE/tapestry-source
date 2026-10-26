<?php

namespace App\Http\Controllers\Service\V1\ReleaseEnterpriseWork;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ReleaseEnterpriseWork;
use Illuminate\Support\Facades\Redis;

class InfoController extends BaseController
{
    /**
     * 灰度企业
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required|integer'
        ]);

        $input = $this->request->all();

        //加缓存
        $key = 'crs:boss:releaseEnterpriseWork:Info1:' . $input['enterprise_id'];
        $redis = Redis::connection();
        $result = $redis->get($key);

        if (!empty($result)) {
            return success(json_decode($result, true));
        }

        $data = ReleaseEnterpriseWork::query()->where('enterprise_id', '=', $input['enterprise_id'])->first();

        if (is_object($data)) {
            $result = $data->toArray();
            $redis->set($key, json_encode($result));
        }else{
            $result = [];
        }

        return success($result);
    }
}
