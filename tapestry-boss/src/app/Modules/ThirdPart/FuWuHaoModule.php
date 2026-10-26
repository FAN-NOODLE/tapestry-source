<?php
/**
 * 服务号相关公共
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/8/13
 * Time: 16:20
 */

namespace App\Modules\ThirdPart;

use App\Models\EnterpriseThirdPartAccount;
use App\Modules\BaseModule;
use EasyWeChat\OpenPlatform\Application;
use Illuminate\Support\Facades\Log;

class FuWuHaoModule extends BaseModule
{
    protected $openPlatform;

    protected $thirdPartAccount;

    public function __construct(Application $openPlatform, EnterpriseThirdPartAccount $thirdPartAccount)
    {
        $this->openPlatform = $openPlatform;
        $this->thirdPartAccount = $thirdPartAccount;
    }

    /**
     * 设置自定义菜单
     *
     * @param EnterpriseThirdPartAccount $thirdPartAccount 服务号信息
     *
     * @param string                     $miniAppId
     *
     * @return bool
     */
    public function createMenu(EnterpriseThirdPartAccount $thirdPartAccount, $miniAppId = '')
    {
        // 小程序Id为空时，获取小程序Id
        if (!$miniAppId) {
            $miniAccount = $this->thirdPartAccount
                ->where([
                    'enterprise_id' => $thirdPartAccount->enterprise_id,
                    'platform'      => EnterpriseThirdPartAccount::OFFICIAL_ACCOUNT,
                ])->first();

            if ($miniAccount) {
                $miniAppId = $miniAccount->app_id;
            }
        }

        if (!$miniAppId) {
            return true;
        }

        $basicUrl = config('wechat_setting.fwh.menu_update_remind_host')
            . '/html/'
            . config('wechat_setting.fwh.menu_update_remind_html')
            . '?id='
            . $thirdPartAccount->enterprise_id;

        $miniUrls = config('wechat_setting.shop_program.page_url');
        $baseMenu = [
            [
                'type'     => 'miniprogram',
                'name'     => '商城',
                'url'      => $basicUrl,
                'appid'    => $miniAppId,
                'pagepath' => $miniUrls['index']
            ],
            [
                'name'       => '我的信息',
                'sub_button' => [
                    [
                        'type'     => 'miniprogram',
                        'name'     => '我的车辆信息',
                        'url'      => $basicUrl,
                        'appid'    => $miniAppId,
                        'pagepath' => $miniUrls['vehicle_info']
                    ],
                    [
                        'type'     => 'miniprogram',
                        'name'     => '我的车检报告',
                        'url'      => $basicUrl,
                        'appid'    => $miniAppId,
                        'pagepath' => $miniUrls['vehicle_inspect']
                    ],
                    [
                        'type'     => 'miniprogram',
                        'name'     => '我的施工录像',
                        'url'      => $basicUrl,
                        'appid'    => $miniAppId,
                        'pagepath' => $miniUrls['build_video']
                    ],
                ]
            ]
        ];

        // 替换企业公众平台配置
        changeOpenPlatformConfig($thirdPartAccount->enterprise_id, $this->openPlatform);

        $res = $this->openPlatform
            ->officialAccount(
                $thirdPartAccount->app_id,
                $thirdPartAccount->refresh_token ?? ''
            )->menu
            ->create($baseMenu);

        Log::info('create menu for official_account : ' . $thirdPartAccount->app_id . '' . var_export($res, true));

        return true;
    }
}
