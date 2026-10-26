<?php

/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Biz;

use App\Models\Enterprise;
use App\Models\EnterpriseOpenPlatform;
use Exception;
use CrsSdk\Shop\Store;
use App\Exceptions\Error;
use Illuminate\Support\Facades\Log;
use EasyWeChat\OpenPlatform\Application;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\OpenPlatform\EnterpriseService;
use App\Services\OpenPlatform\MiniProgramService;
use App\Services\OpenPlatform\OpenPlatformService;

class OpenPlatformBiz extends BaseBiz
{
    /**
     * @var Application
     */
    protected $openPlatform;

    /**
     * @var OpenPlatformService
     */
    protected $openPlatformService;

    /**
     * @var EnterpriseService
     */
    protected $enterpriseService;

    /**
     * @var MiniProgramService
     */
    protected $miniProgramService;

    /**
     * @var EnterpriseThirdPartAccount
     */
    protected $thirdPartAccount;


    /**
     * OpenPlatformBiz constructor.
     *
     * @param Application $openPlatform
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     * @param OpenPlatformService $openPlatformService
     * @param EnterpriseService $enterpriseService
     * @param MiniProgramService $miniProgramService
     */
    public function __construct(
        Application $openPlatform,
        EnterpriseThirdPartAccount $thirdPartAccount,
        OpenPlatformService $openPlatformService,
        EnterpriseService $enterpriseService,
        MiniProgramService $miniProgramService
    )
    {
        $this->openPlatform        = $openPlatform;
        $this->enterpriseService   = $enterpriseService;
        $this->openPlatformService = $openPlatformService;
        $this->miniProgramService  = $miniProgramService;
        $this->thirdPartAccount    = $thirdPartAccount;
    }

    /**
     * 授权并获得授权信息.
     * @param $authCode
     * @param $enterpriseId
     * @param string $authFor
     * @return mixed
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function authorize($authCode, $enterpriseId, $authFor = '')
    {
        // ===================================================
        // 授权信息返回格式：
        // [
        //    "authorization_info" => [
        //          "authorizer_appid" => ,
        //          "authorizer_access_token" => ,
        //          "expires_in" => ,
        //          "authorizer_refresh_token" => ,
        //          "func_info" => [
        //              [
        //                  "funcscope_category" => [
        //                      "id" => 1
        //                  ]
        //              ],
        //              [
        //                  "funcscope_category" => [
        //                      "id" => 1
        //                  ]
        //              ]
        //          ],
        //    ]
        // ]
        // ===================================================
        // 替换企业开放平台配置
        changeOpenPlatformConfig($enterpriseId, $this->openPlatform);
        $authorizeInfo = $this->openPlatform->handleAuthorize($authCode);


        // 由于微信的授权信息错误会返回200 & errcode & errmsg, 需要额外处理
        if (key_exists('errcode', $authorizeInfo) || key_exists('errmsg', $authorizeInfo)) {
            throw new Error(100000, '微信授权错误：' . json_encode($authorizeInfo));
        }

        Log::info('authorization_info-' . json_encode($authorizeInfo));
        // appid是否存在，如果不是一个企业就报错，防止重复授权
        $thirdPartAccount = $this->thirdPartAccount->where('app_id', $authorizeInfo['authorization_info']['authorizer_appid'])->first();
        if ($thirdPartAccount && $thirdPartAccount->status == EnterpriseThirdPartAccount::STATUS_ACTIVE && $thirdPartAccount->enterprise_id != $enterpriseId) {
            throw new Error(100000, '已经被其他企业绑定，请去公众号后台取消授权！');
        }

        // 获取微信账户信息，以判断是小程序还是公众号
        $authorizer = $this->openPlatformService->getAuthorizerConfigByAppId($authorizeInfo['authorization_info']['authorizer_appid'], $enterpriseId);
        $authType   = $authorizer['authorizer_info']['type']; // 可能情况MINI_PROGRAM，FU_WU_HAO，DING_YUE_HAO

        Log::info('authorizer' . var_export($authorizer, true));

        // 绑定类型
        if ($authFor && $authType != EnterpriseThirdPartAccount::MINI_PROGRAM) {
            throw new Error(100000, '绑定类型错误！');
        }

        // 不支持订阅号
        if ($authType == EnterpriseThirdPartAccount::DING_YUE_HAO) {
            throw new Error(100000, '暂不支持订阅号的绑定！');
        }

        // 代开小程序、公众号的第三方的开放平台.
        $openAppId = $this->openPlatformService->setVirtualOpenPlatform($enterpriseId,
            $authorizeInfo['authorization_info']['authorizer_appid'],
            $authorizeInfo['authorization_info']['authorizer_refresh_token'],
            $authType);

        // 代小程序绑定请求域名
        if ($authType == EnterpriseThirdPartAccount::MINI_PROGRAM) {
            $this->miniProgramService->setDomain($authorizeInfo['authorization_info']['authorizer_appid'],
                'set',
                $authorizeInfo['authorization_info']['authorizer_refresh_token'],
                $enterpriseId
                );
        }

        // 保存信息
        return $this->transaction(function () use ($authorizer, $enterpriseId, $openAppId, $authFor) {
            // 企业微信开放平台信息配置
            $openConfig = [
                'enterprise_id' => $enterpriseId,
                'open_app_id'   => $openAppId
            ];
            $this->enterpriseService->updateOrCreateEnterpriseConfig($openConfig);

            // logo http转https
            $logo = !empty($authorizer['authorizer_info']['head_img']) ? str_replace('http:', 'https:', $authorizer['authorizer_info']['head_img']) : '';

            // 获取企业版本, 用于区别企业是标准版还是网关版, 然后找到对应的tag
            // 获取当前企业的信息
            $enterprise = Enterprise::query()
                ->where('id', '=', $enterpriseId)
                ->first();
            // 根据当前企业的版本拿到同版本的最早的企业ID
            $ids = Enterprise::query()
                ->where('version', '=', $enterprise->version)
                ->where('status', '=', Enterprise::STATUS_ENABLE)
                ->orderBy('id', 'asc')
                ->pluck('id')
                ->toArray();

            $release_tag = '';
            foreach ($ids as $id) {
                $account = EnterpriseThirdPartAccount::query()
                    ->where('enterprise_id', '=', $id)
                    ->where('enterprise_type', '=', EnterpriseThirdPartAccount::RELEASE_STABLE)
                    ->first();
                if ($account && $account->release_tag) {
                    $release_tag = $account->release_tag;
                    break;
                }
            }

            $openPaltformConfig = EnterpriseOpenPlatform::query()
                ->where('id', $enterprise->open_platform_id)
                ->first();
            if (empty($openPaltformConfig)) {
                Log::info('开放平台配置不存在，企业id：' . $enterprise->id . "开放平台id" . $enterprise->open_platform_id);
                throw new Error(100000, '开放平台配置不存在！');
            }
            if (empty($openPaltformConfig->ext_build_json)) {
                Log::info('开放平台配置ext_build_json不存在，企业id：' . $enterprise->id . "开放平台id" . $enterprise->open_platform_id);
                throw new Error(100000, '开放平台配置ext_build_json不存在！');
            }
            // 第三方配置
            $thirdPartConfig = [
                'enterprise_id'    => $enterpriseId,
                'app_id'           => $authorizer['authorization_info']['authorizer_appid'],
                'origin_id'        => $authorizer['authorizer_info']['user_name'],
                'name'             => $authorizer['authorizer_info']['nick_name'],
                'platform'         => $authorizer['authorizer_info']['type'] == EnterpriseThirdPartAccount::MINI_PROGRAM ? EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT : EnterpriseThirdPartAccount::OFFICIAL_ACCOUNT,
                'verify_status'    => $authorizer['authorizer_info']['verify_type_info']['id'] == 0 ? EnterpriseThirdPartAccount::VERIFIED : EnterpriseThirdPartAccount::UNVERIFIED,
                'refresh_token'    => $authorizer['authorization_info']['authorizer_refresh_token'],
                'logo'             => $logo,
                'business_info'    => $authorizer['authorizer_info']['business_info'] ? json_encode($authorizer['authorizer_info']['business_info']) : '',
                'principal_name'   => $authorizer['authorizer_info']['principal_name'] ?? '',
                'signature'        => $authorizer['authorizer_info']['signature'] ?? '',
                'other_config'     => isset($authorizer['authorizer_info']['MiniProgramInfo']) ? json_encode($authorizer['authorizer_info']['MiniProgramInfo']) : '',
                'authorizer_range' => isset($authorizer['authorization_info']['func_info']) ? json_encode($authorizer['authorization_info']['func_info']) : '',
                'status'           => EnterpriseThirdPartAccount::STATUS_ACTIVE,
                'authorize_type'   => 2, // 第三方授权
                'type'             => $authFor == 'SHOP' ? EnterpriseThirdPartAccount::TYPE_SHOP : EnterpriseThirdPartAccount::TYPE_QY_WX,
                'ext_build_json' => $openPaltformConfig->ext_build_json ?? "",
            ];

            // 根据当前企业的版本设置对应的 release_tag
            if (!empty($release_tag)) {
                $thirdPartConfig = array_merge($thirdPartConfig, ['release_tag' => $release_tag]);
            }

            // 同步商城名称到商城门店
            try {
                switchAppHost($enterpriseId);

                $paramsStore = [
                    'enterprise_id' => $enterpriseId,
                    'name'          => $authorizer['authorizer_info']['nick_name'],
                    'is_store'      => 2
                ];

                (new Store())->add($paramsStore);

            } catch (Exception $e) {
                Log::info('小程序授权同步门店名称异常' . $e->getMessage());
            }

            // 将 refresh_token 存入缓存
            // Cache::put($thirdPartConfig['app_id'], $thirdPartConfig['refresh_token'], 7200);

            return $this->enterpriseService->updateOrCreateThirdPartConfig($thirdPartConfig);

            // authorizer_access_token 会在下次实际使用中重新获取，这里就不用保存了
        });
    }
}
