<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 4:34 PM
 */

namespace App\Http\Controllers\V1\Auth;

use App\Exceptions\Error;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;
use CrsSdk\Work\EnterprisePlugin;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class WorkCompleteNewController extends BaseController
{
    /**
     * 企业微信信息完善接口
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->all();

        // 去除参数首尾空格
        foreach ($params as $key => $value) {
            $params[$key] = trim($value);
        }

        if (empty($params['external_secret']) || strlen($params['external_secret']) < 10) {
            throw new Error(100000, '外部联系人secret长度错误，请重新填写');
        }

        $params['qy_domain'] = $this->qyDomain;

        // 关键参数在企业信息副表中存储一下
        $params['enterprise_id'] = $this->enterpriseId;
        $params['corp_id'] = $this->enterprise['corp_id'];

        // 自建应用安装
        // 小程序secret、外部联系人secret需要单独提交给UC
        $result = (new EnterprisePlugin())->extInstall($params);

        (new Enterprise())->workSetting($params);

        // 日志埋点
        $this->writeOperatorLog("完善企业微信小程序信息-新");

        return success($result);
    }
}
