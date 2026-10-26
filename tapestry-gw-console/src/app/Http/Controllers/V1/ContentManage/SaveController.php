<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2019/5/6
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\V1\ContentManage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroUser\ConfigGuide;
use CrsSdk\Shop\ContentManage;

class SaveController extends BaseController
{

    /**
     * 保存内容功能
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        // 接值
        $params = $this->request->only([
            'type',
            'content',
            'enable',
            'is_member_new'
        ]);

        $rules    = [
            'type'          => 'required',
            'content'       => 'required',
        ];
        $messages = [
            'type.required'          => '类型不能为空',
            'content.required'       => '内容数据不能为空',
        ];
        $this->validate($this->request, $rules, $messages);

        $this->writeOperatorLog("保存内容配置");

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new ContentManage())->save($params);
        // 通知go-micro-user更新配置引导
        if (isset($result['config_guide'])) {
            $configs = [];
            foreach ($result['config_guide'] as $type => $status) {
                $configs[] = [
                    'type' => $type,
                    'status' => $status
                ];
            }
            $configGuideParams = [
                'enterprise_id' => $this->enterpriseId,
                'configs' => $configs
            ];
            (new ConfigGuide())->set($configGuideParams);
        }
        return success($result ?? []);
    }
}
