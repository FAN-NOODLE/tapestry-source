<?php

namespace App\Http\Controllers\V1\Invite;

use App\Micros\Work\Work;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;

class InviteController extends BaseController
{
    /**
     * 生成太阳码
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $this->validateData($request);

        /**
         * @param string    path            小程序页面路径
         * @param int       enterprise_id   企业ID
         * @param array     params          太阳码携带的参数
         */
        $params = $request->only('path', 'params');
        $params['enterprise_id'] = $this->enterpriseId;

        $data = Work::createQrcode($params);

        return success($data);

    }

    public function validateData($request)
    {
        $rules    = [
            'path'          => 'required',
        ];
        $messages = [
            'path.required'          => 'Sorry,二维码路径不能为空',
        ];
        $this->validate($request, $rules, $messages);
    }
}