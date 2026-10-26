<?php
/**
 * 保存配置
 * User: 何岳龙
 * Date: 2018/11/26
 * Time: 10:04
 */
namespace App\Http\Controllers\V1\Config;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\Shop\Config;

class SaveController extends BaseController
{
    /**
     * 更新企业配置数据,然后再更新配置信息对应的redis缓存数据
     * @return mixed
     */
    public function index()
    {
        $params = $this->request->only(['data']);

        $rules = [
            'data'         => 'required|array',
            'data.*.type'  => 'required|integer',
            'data.*.key'   => 'required',
            'data.*.value' => 'required',
            'data.*.name'  => 'string',
        ];

        $message = [
            'data.*'                => '设置数据不符合规范',
            'data.*.type.*'         => '类型不符合规范',
            'data.*.key.required'   => '字段不能为空',
            'data.*.value.required' => '字段值不能为空',
        ];

        $this->validate(request(), $rules, $message);

        $params['enterprise_id'] = $this->enterpriseId;

        foreach ($params['data'] as &$v) {
            if (!isset(Constant::CONFIG_KEY_DESC[$v['key']])) error(600116);
            $v['desc'] = Constant::CONFIG_KEY_DESC[$v['key']] ?? ' ';
        }

        // 兼容新输出格式
        $info           = (new EnterpriseConfig())->save($params);
        $result['list'] = isset($info['data']) ? $info['data'] : $info;

        // 调用shop端接口,更新配置对应的redis数据
        (new Config())->updateRedisKey(['enterprise_id'=>$this->enterpriseId]);

        return success($result);
    }
}

