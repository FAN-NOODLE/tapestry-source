<?php
/**
 * 小程序自动打包回调
 */

namespace App\Http\Controllers\Callback;

use App\Models\Enterprise;
use App\Services\OpenPlatform\MiniProgramService;
use CrsSdk\Shop\AutoRelease;
use App\Jobs\AutoPackAuditJob;
use Illuminate\Support\Facades\Log;
use App\Services\EnterpriseLogService;
use App\Models\EnterpriseMiniProgramAutoPack;

class AutoPackAuditController extends BaseController
{
    public function index(MiniProgramService $miniProgramService)
    {
        $param = request()->input();

        EnterpriseLogService::info('自动化打包回调来了', ['input' => $param]);

        Log::info('自动化打包回调来了', ['input' => $param]);

        if (empty($param)) {
            return 'failed:empty request';
        }

        if (!isset($param['id']) || !isset($param['status']) || !isset($param['statusDesc'])) {
            return 'failed:lack of parameters';
        }

        $autoPackInfo = EnterpriseMiniProgramAutoPack::query()
            ->where('task_id', '=', $param['id'])
            ->orderBy('id', 'DESC')
            ->first(['status', 'task_config', 'version_id', 'enterprise_id']);

        if (empty($autoPackInfo)) {
            Log::info('自动化打包，提审失败： 提交代码记录不存在');
            return 'failed:not found pack record';
        }

        $param['hash'] = '';

        $enterpriseInfo = Enterprise::query()->find((int)$autoPackInfo->enterprise_id);
        if (!empty($enterpriseInfo)) {

            // 切换SHOP域名
            switchAppHost($autoPackInfo->enterprise_id);

            $param['hash']              = $enterpriseInfo->hash ?? '';
            $_SERVER['enterprise_hash'] = $param['hash'];
        }

        $version = (new AutoRelease())->versionInfo(['version_id' => $autoPackInfo->version_id]);
        EnterpriseLogService::info('自动化打包版本记录', [
            'version_id' => $autoPackInfo->version_id,
            'version'    => $version
        ]);
        if (empty($version) || empty($version['data'])) {
            Log::info('自动化打包，提审失败： 版本记录不存在');
            return 'failed:not found version record';
        }

        if ($version['data']['status'] != EnterpriseMiniProgramAutoPack::SHOP_VERSION_PACK_ING) {
            Log::info('自动化打包，提审失败： 版本状态不符合提审条件：' . json_encode($version));
            return 'failed:version status not allowed';
        }
        $param['version_id'] = $autoPackInfo->version_id;
        //开始进入队列
        Log::info('自动化打包，开始异步处理：' . json_encode($param));
        EnterpriseLogService::info('自动化打包，开始异步处理', $param);

        // 替换开放平台配置
        changeOpenPlatformConfig((int)$autoPackInfo->enterprise_id, $miniProgramService->openPlatform);
        (new AutoPackAuditJob($param))->handle($miniProgramService);

        return 'success:the request in job';
    }
}
