<?php
/**
 * 同步订阅模板消息
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019/09/23
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\Message;

use App\Services\SubscribeMessageService;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Validator;

class MiniSyncTemplatesController extends BaseController
{
    protected $params;

    /**
     * 手动更新小程序模板消息
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        Validator::make($this->request->all(), [
            'enterprise_id' => 'required',
        ], [
            'enterprise_id.required' => '请确认企业',
        ])->validate();

        $enterpriseId = $this->request->input('enterprise_id');

        Log::info('手动更新小程序模板消息：开始');
        $res =  (new SubscribeMessageService($enterpriseId))->sync();
        Log::info('手动更新小程序模板消息：结束');

        return success($res);
    }
}