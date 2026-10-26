<?php
/**
 * 小程序开发码回调
 */

namespace App\Http\Controllers\Callback;

use CrsSdk\Shop\DevelopCode;
use Illuminate\Support\Facades\Log;
use App\Services\EnterpriseLogService;

class DevelopCodeController extends BaseController
{
    public function index($hash)
    {
        $param = request()->input();

        EnterpriseLogService::info('开发码回调来了', [
            'input' => $param
        ]);
        Log::info('开发码 回调来了:' . json_encode($param));

        if (empty($param)) {
            return 'failed:empty request';
        }
        if (!isset($param['id']) || !isset($param['status']) || !isset($param['statusDesc'])) {
            return 'failed:lack of parameters';
        }
        if (empty($hash)) {
            return 'failed:not find hash';
        }
        $_SERVER['enterprise_hash'] = $hash;

        // 切换SHOP域名
        switchAppHost($hash);

        try {
            $res = (new DevelopCode())->call($param);
            if ($res) {
                return 'success';
            }
        } catch (\Exception $e) {
            EnterpriseLogService::error('开发码回调通知shop报错', [
                'host'    => config('client.api_url.shop'),
                'hash'    => $hash,
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine(),
                'params'  => $param
            ]);
            return $e->getMessage();
        }

        return 'failed';
    }
}
