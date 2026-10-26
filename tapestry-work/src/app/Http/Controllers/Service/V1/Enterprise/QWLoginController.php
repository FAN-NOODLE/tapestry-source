<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Common;
use CrsSdk\Boss\Auth;
use CrsSdk\Boss\Enterprise;

class QWLoginController extends BaseController
{

    /**
     * 标准标识
     */
    const STANDARD_IDENTIFIER = 1;

    /**
     * 灰度标识
     */
    const RELEASE_IDENTIFIER = 2;

    /**
     * 标准(爱点击主体)标识
     */
    const STANDARD_ICLK_IDENTIFIER = 3;

    /**
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        // 应用市场
        $code = $this->request->post('code');
        // 自建用
        $corpId = $this->request->post('corp_id');
        // 是否为灰度企业标识
        $identifier = $this->request->post('identifier');

        if (empty($code)) {
            error(110001);
        }

        if (empty($corpId)) {
            $result = $this->appStore($code, $identifier);
        } else {
            $result = $this->selfPlugin($code, $corpId);
        }

        return success(['data' => $result]);
    }

    /**
     * 自建应用 自带corp_id
     * @param $code
     * @param $corpId
     * @return array
     * @throws \Exception
     */
    public function selfPlugin($code, $corpId)
    {
        $enterprise = (new Auth())->daoGouPlugin(['corp_id' => $corpId]);

        if (empty($enterprise['data'])) {
            error(200020);
        }

        $result = (new Common())->miniProgramLogin($enterprise['data']['domain'], $enterprise['data']['identifier'], ['code' => $code]);

        if (!is_array($result)) {
            error(200021);
        }

        return [
            'hash'          => $enterprise['data']['hash'],
            'mem_uid'       => $result['memUid'],
            'enterprise_id' => $enterprise['data']['enterprise_id'],
            'mobile'        => $result['memMobile'],
            'userid'        => $result['userid']
        ];
    }

    /**
     * 应用市场 使用suite套件
     * @param $code
     * @return array
     * @throws \Exception
     */
    public function appStore($code, $identifier)
    {
        $env = config('app.env');
        switch ($identifier) {
            case self::STANDARD_IDENTIFIER :
                $suiteId = config("work_plugin.crs-store-assistant.{$env}.shoppers_suite_id");
                break;
            case self::RELEASE_IDENTIFIER :
                $suiteId = config("work_plugin.crs-store-assistant.{$env}.release_shoppers_suite_id");
                break;
            case self::STANDARD_ICLK_IDENTIFIER :
                $suiteId = config("work_plugin.crs-store-assistant.{$env}.iclk_shoppers_suite_id");
                break;
            default :
                $suiteId = config("work_plugin.crs-store-assistant.{$env}.shoppers_suite_id");
        }

        $result = (new Common())->hubMiniProgramLogin([
            'code'    => $code,
            'suiteid' => $suiteId
        ]);

        if (!is_array($result)) {
            error(200021);
        }

        $enterpriseData = (new Enterprise())->detail(
            ['corp_id' => $result['corpid']]
        );

        return [
            'hash'          => $enterpriseData['data']['hash'] ?? '',
            'mem_uid'       => $result['memUid'],
            'enterprise_id' => $enterpriseData['data']['id'] ?? 0,
            'mobile'        => $result['memMobile']
        ];
    }
}