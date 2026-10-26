<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 4:47 PM
 */

namespace App\Http\Controllers\Service\V1\EnterprisePlugin;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\EnterprisePlugin;
use App\Library\Uc\Member;
use Illuminate\Http\Request;

class ExtShopInstallController extends BaseController
{
    /**
     * 扩展应用安装接口
     *
     * @param Request $request
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 4:54 PM
     */
    public function index(Request $request)
    {
        $miniProgramSecret = $request->post('mini_program_secret');
        $corpId            = $request->post('corp_id');
        $qyDomain          = $request->post('qy_domain');
        $appIdentifier     = $request->post('app_identifier');
        $agentId           = $request->post('agent_id', 0);
        if (empty($qyDomain) || empty($miniProgramSecret) || empty($corpId)) {
            throw new Error(100000, '必填参数错误');
        }

        $enterprisePlugin = new EnterprisePlugin();

        // 安装小程序应用
        $miniInstallResult = $enterprisePlugin->appInstall(
            [
                'corpId'            => $corpId,
                'installType'       => 1,
                'agentId'           => $agentId, // 小程序agentId因为没有，所以是0
                'plName'            => '商城小程序',
                'plIdentifier'      => $appIdentifier,
                'token'             => $this->randCode(6),
                'secret'            => $miniProgramSecret,
                'plHost'            => 'http://crs.vchangyi.com',
                'encodingAeskey'    => $this->randCode(43),
                'imgName'           => md5(time()),
                'enterpriseChannel' => 4  // 渠道是crs
            ]
        );

        //调UC同步通讯录
        (new Member())->syncContact($qyDomain);

        return success([
            'data' => [
                'mini_program_install_result' => $miniInstallResult
            ]
        ]);
    }

    /**
     * 生成随机字符串
     *
     * @param  integer $length [生成的长度]
     * @param  integer $type   [生成的类型]
     * @return [type]   str       [description]
     * @php 随机码类型：0，数字+大写字母；1，数字；2，小写字母；3，大写字母；4，特殊字符；-1，数字+大小写字母+特殊字符
     */
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
            $code .= $str[$i];
        }

        return $code;
    }
}
