<?php

namespace App\Http\Controllers\Service\V1\EnterprisePlugin;

use CrsSdk\Boss\Enterprise;
use App\Library\Uc\EnterprisePlugin;
use App\Http\Controllers\Service\V1\BaseController;

class DaoGouInstallController extends BaseController
{
    /**
     * 安装导购小程序
     * @return array|void
     * @throws \Exception
     */
    public function daoGouInstall()
    {
        $params = $this->request->all();
        if (empty($params['enterprise_id'])) {
            error(600000);
        }

        $enterprise = (new Enterprise())->detail(['enterprise_id' => $params['enterprise_id']]);
        if (empty($enterprise['data'])) {
            error(110008);
        }

        if (empty($params['agent_id'] || empty($params['mini_program_secret']))) {
            error(110007);
        }

        $originSecret  = $this->request->input('origin_secret');
        $originAgentId = $this->request->input('origin_agent_id');

        if (empty($originSecret) && empty($originAgentId)) {
            // 安装小程序应用
            $miniInstallResult = (new EnterprisePlugin())->appInstall([
                'corpId'            => $enterprise['data']['corp_id'],
                'installType'       => 1,
                'agentId'           => $params['agent_id'], // 小程序agentId因为没有，所以是0
                'plName'            => '商城小程序',
                'plIdentifier'      => 'crs_store_assistant_self',
                'token'             => $this->randCode(6),
                'secret'            => $params['mini_program_secret'],
                'plHost'            => 'http://crs.vchangyi.com',
                'encodingAeskey'    => $this->randCode(43),
                'imgName'           => md5(time()),
                'enterpriseChannel' => 4  // 渠道是crs
            ]);
        } else {
            // 安装小程序应用
            $miniInstallResult = (new EnterprisePlugin())->appUpdate($params['qy_domain'], [
                'agentid'    => $params['agent_id'],
                'plSecret'   => $params['mini_program_secret'],
                'plPluginid' => $params['pl_plugin_id'],
            ]);
        }


        return success([
            'data' => [
                'mini_program_install_result' => $miniInstallResult
            ]
        ]);
    }

    private function randCode($length = 5, $type = 0)
    {
        // 生成字符串的值
        $arr = [
            1 => "0123456789",
            2 => "abcdefghijklmnopqrstuvwxyz",
            3 => "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
            4 => "~@#$%^&*(){}[]|"
        ];
        // 初始化
        $code = '';
        // 当前搭配
        if ($type == 0) {
            array_pop($arr);
            $string = implode("", $arr);
        } else if ($type == "-1") {
            $string = implode("", $arr);
        } else {
            $string = $arr[$type];
        }
        // 生成字符串
        $count = strlen($string) - 1;
        for ($i = 0; $i < $length; $i++) {
            $str[$i] = $string[rand(0, $count)];
            $code    .= $str[$i];
        }

        return $code;
    }
}