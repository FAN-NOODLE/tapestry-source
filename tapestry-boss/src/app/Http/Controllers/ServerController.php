<?php

namespace App\Http\Controllers;

use Illuminate\Support\Facades\Log;
use EasyWeChat\OpenPlatform\Auth\VerifyTicket;
use App\Services\OpenPlatform\EnterpriseService;
use App\Models\EnterpriseOpenPlatform;
use App\Services\OpenPlatform\FastRegisterService;
use Overtrue\LaravelWeChat\Facade as LaravelWechat;

class ServerController extends Controller
{
    /**
     * 接收微信事件消息
     *
     * @param EnterpriseService   $enterpriseService
     * @param VerifyTicket        $verifyTicket
     * @param FastRegisterService $fastRegisterService
     *
     * @return \Symfony\Component\HttpFoundation\Response
     * @throws \EasyWeChat\Kernel\Exceptions\BadRequestException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidArgumentException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \ReflectionException
     */
    public function index(EnterpriseService $enterpriseService, VerifyTicket $verifyTicket, EnterpriseOpenPlatform $enterpriseOpenPlatform,  FastRegisterService $fastRegisterService)
    {
        Log::info('微信回调进来了');

        // 注意此处是大深坑，必须要用Overtrue\LaravelWeChat\Facade来完成初始化,不能用Factory::openPlatform
        $openPlatform = LaravelWechat::openPlatform('default');
        // $openPlatform = Factory::openPlatform(config('wechat.open_platform.default'));
        $openPlatformId = $this->request->input('id');
        $openPlatformConfig = $enterpriseOpenPlatform->where('id', $openPlatformId)->first();
        $openPlatform->config['app_id'] = $openPlatformConfig->app_id;
        $openPlatform->config['secret'] = $openPlatformConfig->secret;
        $openPlatform->config['token'] = $openPlatformConfig->token;
        $openPlatform->config['aes_key'] = $openPlatformConfig->aes_key;

        $server = $openPlatform->server;
        // 处理令牌事件
        $server->on('component_verify_ticket', function ($message) use ($verifyTicket, $openPlatform, $enterpriseOpenPlatform) {
            Log::info('令牌更新');
            Log::info($message);
            $verifyTicket->setTicket($message['ComponentVerifyTicket']);
            Log::info('Cache的ComponentVerifyTicket：' . $verifyTicket->getTicket());
        });

        $server->on('authorized', function ($message) {
            Log::info('授权成功');
            Log::info($message);
        });

        $server->on('unauthorized', function ($message) use ($enterpriseService) {
            Log::info('取消授权');
            Log::info($message);
            $enterpriseService->negativeThirdPartStatus($message['AuthorizerAppid']);
        });

        $server->on('updateauthorizeds', function ($message) {
            Log::info('更新授权');
            Log::info($message);
        });


        $server->on('notify_third_fasteregister', function ($message) use ($fastRegisterService) {
            Log::info('代注册审核');
            Log::info($message);
            $fastRegisterService->updateRegisterAgent($message);
        });

        return $server->serve();
    }
}
