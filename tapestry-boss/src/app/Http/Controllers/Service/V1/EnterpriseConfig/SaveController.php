<?php
/**
 * 企业配置设置
 * User: zhaoqi
 * Date: 2019/02/20
 * Time: 15:26
 */

namespace App\Http\Controllers\Service\V1\EnterpriseConfig;

use App\Models\EnterpriseConfig;
use App\Services\EnterpriseConfigService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;
use App\Http\Controllers\Service\V1\BaseController;

class SaveController extends BaseController
{

    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'data']);

        $rules = [
            'data'         => 'required',
            'data.*.type'  => 'required|integer',
            'data.*.key'   => 'required',
            'data.*.value' => 'required',
            'data.*.desc'  => 'required',
            'data.*.name'  => '',
        ];

        $message = [
            'data.required'         => '设置数据不符合规范',
            'data.*.type.*'         => '类型不符合规范',
            'data.*.key.required'   => '字段不能为空',
            'data.*.value.required' => '字段值不能为空',
            'data.*.desc.required'  => '字段说明不能为空',
        ];

        $this->validate(request(), $rules, $message);

        $enterpriseId = $params['enterprise_id'] ?? 0;

        foreach ($params['data'] as $v) {
            $where = [
                'type'          => $v['type'],
                'key'           => $v['key'],
                'enterprise_id' => $enterpriseId
            ];

            // 如果开启好物圈则验证
            if ($v['key'] == "mall_open" && $v['value'] == 1) {
                $isConfig = (new EnterpriseConfigService())->isSettingConfig($enterpriseId);
                if (!$isConfig) {
                    return notice(620001, [], [], 210);
                }
            }

            // 企业ID赋值
            $v['enterprise_id'] = $enterpriseId;

            $query = EnterpriseConfig::query()->where($where)->first();

            if (!empty($query)) {
                $query->update($v);
            } else {
                EnterpriseConfig::query()->create($v);
            }

        }

        // redis键名
        $redisKey = redisKey('ENTERPRISE_CONFIG', ['enterprise_id' => $enterpriseId]);
        Redis::del($redisKey);

        return success(['data' => []]);
    }
}
