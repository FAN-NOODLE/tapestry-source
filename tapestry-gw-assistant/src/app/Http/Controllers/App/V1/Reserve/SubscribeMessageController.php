<?php

namespace App\Http\Controllers\App\V1\Reserve;

use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use CrsSdk\Sdk;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SubscribeMessageController extends Controller
{

    /**
     * 获取订阅消息配置
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function getConfig(Request $request): JsonResponse
    {
        $tokenData = $request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        $res = Sdk::work()->setApi("service/v1/subscribe-message/get-config")->post($params);

        return success($res);
    }


    /**
     * 订阅
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function subscribe(Request $request): JsonResponse
    {
        $rules = [
            'tmpl_list' => "required",
            'tmpl_list.*.url' => "required",
            'tmpl_list.*.pri_tmpl_id' => "required",
            'rel_id' => "required",
        ];
        $params = $this->validate($this->request, $rules);
        $tokenData = $request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['union_id'] = $tokenData['union_id'];
        $params['open_id'] = $tokenData['open_id'];
        $res = Sdk::work()->setApi("service/v1/subscribe-message/subscribe")->post($params);

        return success($res);
    }

}