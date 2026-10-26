<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/8/16
 * Time: 16:17
 */

namespace App\Http\Controllers\V1\LiveBroadcast;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\MicroActivity\Live;
use Illuminate\Support\Facades\Log;

class AddController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();

        try {
            $miniProgramData = (new MiniProgram())->getOnlineVersion([
                'enterprise_id' => $this->enterpriseId
            ]);
            Log::info('生成太阳码-获取版本号：' . var_export([
                'miniProgramData' => $miniProgramData,
                'time'            => date('Y-m-d H:i:s')
            ],
                true));
            // $apiTag = explode('_', $miniProgramData['use_tag']);
            // $apiTag = $apiTag[0];
            if (!empty($miniProgramData['use_branch'])) {
                $_SERVER['api_version'] = 'v1.0';
            }
        } catch (\Exception $exception) {
            Log::error('直播生成太阳码使用版本信息失败' . $exception->getMessage());
        }

        $input['enterprise_id'] = $this->enterpriseId;

        $data = (new Live())->add($input);
        

        $this->writeOperatorLog('添加直播');

        return success($data);
    }
}
