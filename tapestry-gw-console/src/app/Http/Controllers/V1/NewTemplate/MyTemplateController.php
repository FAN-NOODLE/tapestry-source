<?php

namespace App\Http\Controllers\V1\NewTemplate;

use CrsSdk\MicroUser\ConfigGuide;
use CrsSdk\Shop\NewTemplate;
use App\Http\Controllers\V1\BaseController;

class MyTemplateController extends BaseController
{
    public function myDetail()
    {
        $detail = (new NewTemplate())->MyDetail(['enterprise_id' => $this->enterpriseId]);

        return success($detail);
    }

    public function mySave()
    {
        $detail = (new NewTemplate())->MySave([
            'content'       => $this->request->input('content'),
            'enterprise_id' => $this->enterpriseId,
            'is_member_new' => $this->request->input('is_member_new', 0)
        ]);

        // 通知go-micro-user更新配置引导
        if (isset($detail['config_guide'])) {
            $configs = [];
            foreach ($detail['config_guide'] as $type => $status) {
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
        return success($detail);
    }
}