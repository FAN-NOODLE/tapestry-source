<?php

namespace App\Http\Controllers\V1\Auth;

use CrsSdk\Boss\Auth;
use App\Exceptions\Error;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class InstallAppController extends BaseController
{
    /**
     * 手动安装应用
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function install(Request $request)
    {
        $params = $request->only(['agent_id', 'secret', 'check_file']);
        if (empty($params['agent_id']) || empty($params['secret'])) {
            throw new Error(1400, '参数不完整!');
        }

        $params['enterprise_id'] = $this->enterpriseId;

        (new Auth())->installApp($params);

        return success([]);
    }

    /**
     * 获取自建应用信息
     * @return mixed
     * @throws \Exception
     */
    public function getWorkPlugin()
    {
        $plugin = (new Auth())->getWorkPlugin(['enterprise_id' => $this->enterpriseId]);

        return success($plugin);
    }

    // 手动安装应用
    public function update(Request $request)
    {
        $params = $request->only(['agent_id', 'secret', 'check_file']);
        if (empty($params['agent_id']) || empty($params['secret'])) {
            throw new Error(1400, '参数不完整!');
        }

        $params['enterprise_id'] = $this->enterpriseId;

        (new Auth())->updateApp($params);

        return success([]);
    }

    /**
     * 安装第三方应用
     *
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function installThirdApp(Request $request)
    {
        $data = [
            'enterprise_id' => $this->enterpriseId,
            'back_url'      => $request->input('back_url', ''),
            "wisdom_set"    => $request->input('wisdom_set', 0),
        ];

        $plugin = (new Auth())->installThirdApp($data);

        $this->writeOperatorLog("企业微信第三方小程序授权");

        return success($plugin);
    }

}
