<?php
/*
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Http\Controllers\Callback;

use App\Events\WeappAudit\FailureEvent;
use App\Events\WeappAudit\SuccessEvent;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\EnterpriseLogService;
use App\Services\OpenPlatform\OpenPlatformService;
use Carbon\Carbon;
use CrsSdk\MicroGoodsV1\Video;
use EasyWeChat\OpenPlatform\Application;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AccountController extends BaseController
{
    protected $openPlatform;
    protected $thirdPartAccount;
    protected $openPlatformService;
    protected $request;
    protected $userImport;

    /**
     * AccountController constructor.
     *
     * @param Application $openPlatform
     * @param EnterpriseThirdPartAccount $thirdPartAccount
     * @param OpenPlatformService $openPlatformService
     * @param Request $request
     */
    public function __construct(
        Application $openPlatform, EnterpriseThirdPartAccount $thirdPartAccount, OpenPlatformService $openPlatformService,
        Request $request
    ) {
        $this->openPlatform = $openPlatform;
        $this->thirdPartAccount = $thirdPartAccount;
        $this->openPlatformService = $openPlatformService;
        $this->request = $request;
    }

    /**
     * 第三方平台统一通知:
     *
     * @param $appId
     * @return \Symfony\Component\HttpFoundation\Response|bool
     * @throws \EasyWeChat\Kernel\Exceptions\BadRequestException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidArgumentException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \ReflectionException
     * @author
     */
    public function index($appId)
    {
        Log::info('第三方平台统一通知:start', [
            'app_id' => $appId
        ]);

        // 确定通知是何种类型的平台
        $config = $this->thirdPartAccount->where('app_id', $appId)
            ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('verify_status', '=', EnterpriseThirdPartAccount::VERIFIED)
            ->with('enterprise:id,hash')
            ->first();
        if (!$config) { // 未在企业授权内的 appId 排除掉
            return 'success';
        }

        Log::info('第三方平台统一通知:step-1', [
            'app_id'               => $appId,
            'open_platform_app_id' => $this->openPlatform->config['app_id'],
        ]);

        // 替换企业开放平台配置
        changeOpenPlatformConfig($config->enterprise_id, $this->openPlatform);

        if ($config->platform == EnterpriseThirdPartAccount::OFFICIAL_ACCOUNT) {
            $app = $this->openPlatform->officialAccount($appId, $config->refresh_token)->server;
        } else {
            $app = $this->openPlatform->miniProgram($appId, $config->refresh_token)->server;
        }

        Log::info('第三方平台统一通知:step-2', [
            'app_id'               => $appId,
            'open_platform_app_id' => $this->openPlatform->config['app_id'],
        ]);

        // 处理推送过来的消息
        $openPlatformService = $this->openPlatformService;
        $app->push(function ($message) use ($appId, $openPlatformService, $config) {
            EnterpriseLogService::info('第三方平台统一通知:事件处理', [
                'app_id'  => $appId,
                'message' => $message,
                'config'  => $config
            ]);
            // 微信服务器测试处理
            $info = $openPlatformService->wechatWholeNetworkRelease($appId, $message, $config->enterprise_id);
            if (!is_bool($info)) {
                return $info;
            }

            // 真实用户处理
            switch ($message['MsgType']) {
                case 'event':  // 事件消息处理
                    $eventRes = $this->handleWechatEvent($message, $config);
                    if ('success' == $eventRes) {
                        return 'success';
                    }
                    break;
                case 'text':
                    if ($message['Content'] == '获取身份') {
                        return 'Hello ' . $message['FromUserName'] . ' , Type ' . $message['MsgType'];
                    }
                    break;
            }

            return '欢迎关注公众号！';
        });

        return $app->serve();
    }

    /**
     * 事件处理.
     *
     * @param $message
     * @param $thirdPartAccount
     * @return string
     */
    protected function handleWechatEvent($message, $thirdPartAccount)
    {
        // 处理场景值
        // ==========================================================================================================
        // array('Event' => 'subscribe', 'EventKey' => null, 'FromUserName' => openId, 'ToUserName' => originId)
        // Event包含subscribe/unsubscribe/SCAN，EventKey包含null/qrscene_?/?，?是scene_id
        // FromUserName为openId，ToUserName为公众号原始id
        // ===========================================================================================================
        Log::info('第三方平台统一通知:step-3', [
            'app_id'  => $thirdPartAccount->app_id,
            'message' => $message,
        ]);
        switch ($message['Event']) {
            case 'weapp_audit_success': // 微信小程序审核成功
                $this->weAppAudit($message, $thirdPartAccount->app_id);
                break;

            case 'weapp_audit_fail': // 微信小程序审核失败
                $this->weAppAudit($message, $thirdPartAccount->app_id);
                break;

            case 'subscribe':  // 订阅事件处理 客户授权表增加相关信息
                //$this->userImport->import($thirdPartAccount, 1, [$message['FromUserName']]);
                break;

            case 'unsubscribe':  // 取消订阅事件
                break;

            case 'SCAN':  // 扫码关注事件
                break;

            case 'open_product_account_audit':  // 视频号商家接入回调
            case 'open_product_spu_audit':  // 视频号商品审核回调
            case 'open_product_category_audit':  // 视频号类目审核结果
            case 'open_product_brand_audit':  // 视频号品牌审核结果
                try {
                    $hash = '';
                    if (!empty($thirdPartAccount->enterprise->hash)) {
                        $hash = $thirdPartAccount->enterprise->hash;
                        $_SERVER['enterprise_hash'] = $hash;
                    }
                    Log::info('企业[' . $hash . ']视频号回调', $message);
                    $res = (new Video())->handlerCallBack($message);
                    Log::info("视频号回调转发Goods成功", $res);
                    return 'success';
                } catch (\Exception $exception) {
                    Log::info("视频号回调转发Goods失败" . $exception->getMessage());
                }
                break;

            default:
                break;
        }
    }

    /**
     * 小程序代码审核事件处理.
     *
     * @param $message
     * @param $appId
     */
    protected function weAppAudit($message, $appId)
    {
        Log::info('第三方平台统一通知:step-4', [
            'app_id'  => $appId,
            'message' => $message,
        ]);
        if (key_exists('SuccTime', $message)) {
            $info = [
                'app_id'       => $appId,
                'to'           => $message['ToUserName'],
                'from'         => $message['FromUserName'],
                'created_at'   => Carbon::createFromTimestamp($message['CreateTime']),
                'message_type' => $message['MsgType'],
                'event'        => $message['Event'],
                'succeeded_at' => Carbon::createFromTimestamp($message['SuccTime']),
            ];
            EnterpriseLogService::info('小程序代码审核事件处理-Succ', $info);
            event(new SuccessEvent($info));
        } else {
            $info = [
                'app_id'       => $appId,
                'to'           => $message['ToUserName'],
                'from'         => $message['FromUserName'],
                'created_at'   => Carbon::createFromTimestamp($message['CreateTime']),
                'message_type' => $message['MsgType'],
                'event'        => $message['Event'],
                'reason'       => $message['Reason'],
                'failed_at'    => Carbon::createFromTimestamp($message['FailTime']),
            ];
            EnterpriseLogService::info('小程序代码审核事件处理-Fail', $info);
            event(new FailureEvent($info));
        }
    }
}
