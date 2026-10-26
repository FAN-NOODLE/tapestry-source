<?php

namespace App\Modules;

use App\Modules\Login\ScanCodeLogin;
use App\Modules\Login\WorkWechatLogin;

class LoginFactory
{
    /**
     * 扫码
     */
    const SCAN_CODE = 'scan_code';

    /**
     * 服务商后台直接跳转
     */
    const WORK_WECHAT = 'work_wechat';

    /**
     * 获取处理类
     *
     * @param $type
     * @param $data
     * @return ScanCodeLogin|WorkWechatLogin
     * @throws \Exception
     */
    public static function getInstance($type, $data)
    {
        $class = null;

        switch ($type) {
            case self::SCAN_CODE:
                $class = new ScanCodeLogin($data['auth_code']);
                break;
            case self::WORK_WECHAT:
                $class = new WorkWechatLogin($data['token'], $data['corpid'], $data);
                break;
            default:
                throw new \Exception('处理方式不存在');
        }

        return $class;
    }
}
