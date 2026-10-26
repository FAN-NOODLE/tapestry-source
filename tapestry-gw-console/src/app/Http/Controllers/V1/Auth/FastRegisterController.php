<?php

namespace App\Http\Controllers\V1\Auth;

use CrsSdk\Boss\FastRegisterWeApp;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class FastRegisterController extends BaseController
{
    /**
     * 申请
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function apply(Request $request)
    {
        $params = $request->only([
            'id',
            'legal_persona_name',
            'legal_persona_wechat',
            'name',
            'code_type',
            'code',
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        $res = (new FastRegisterWeApp())->apply($params);

        return success($res);
    }

    /**
     * 获取详情
     *
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function detail(Request $request)
    {
        $params = $request->only([
            'app_id',
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        $res = (new FastRegisterWeApp())->detail($params);

        return success($res);
    }

    /**
     * 更新信息
     *
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function updateInfo(Request $request)
    {

        $params['enterprise_id'] = $this->enterpriseId;

        $res = (new FastRegisterWeApp())->updateInfo($params);

        return success($res);
    }

}
