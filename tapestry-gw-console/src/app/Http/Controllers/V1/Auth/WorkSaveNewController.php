<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 19-2-27
 * Time: 下午3:37
 */

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\EnterprisePlugin;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class WorkSaveNewController extends BaseController
{
    /**
     * WorkSaveController constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * work授权信息编辑
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only(['mini_program_secret', 'external_secret', 'pl_plugin_id', 'corp_id']);

        if (empty(trim($params['external_secret'])) || strlen(trim($params['external_secret'])) < 10) {
            notice(451020, [], $params);
        }

        if (empty(trim($params['pl_plugin_id']))) {
            notice(451012, [], $params);
        }

        try {
            // work小程序
            $params['qy_domain'] = $this->qyDomain;
            $params['app_identifier'] = 'crs-mini-program';
            $params['enterprise_id'] = $this->enterpriseId;

            // 小程序secret、外部联系人secret需要单独提交给UC
            (new EnterprisePlugin())->extUpdate($params);

            (new Enterprise())->workSetting($params);

            // 日志埋点
            $this->writeOperatorLog("编辑企业微信小程序信息-新");

            return success();
        } catch (\Exception $e) {
            Log::error($e->getMessage());
            notice(451000, [], $params);
        }
    }
}