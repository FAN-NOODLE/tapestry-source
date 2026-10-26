<?php


namespace App\Http\Controllers\V1\NewTemplate;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroUser\ConfigGuide;
use CrsSdk\Shop\NewTemplate;

class AuthorizationTemplateController extends BaseController
{
    public function authorizationDetail()
    {
        $detail = (new NewTemplate())->AuthorizationDetail(['enterprise_id' => $this->enterpriseId]);

        return success($detail);
    }

    public function authorizationSave()
    {
        $result = (new NewTemplate())->AuthorizationSave([
            'enterprise_id' => $this->enterpriseId,
            'content' => $this->request->input('content')
        ]);

        // 通知go-micro-user更新配置引导
        if (isset($result['config_guide'])) {
            $configs = [];
            foreach ($result['config_guide'] as $type => $status) {
                $configs[] = [
                    'type' => $type,
                    'status' => $status
                ];
            }
            $params = [
                'enterprise_id' => $this->enterpriseId,
                'configs' => $configs
            ];
            (new ConfigGuide())->set($params);
        }

        return success([]);
    }
}