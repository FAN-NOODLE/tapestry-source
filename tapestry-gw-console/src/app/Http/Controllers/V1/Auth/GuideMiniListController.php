<?php
/*
* +------------------------------------+
* | Created: phpstorm                  
* +------------------------------------+
* | Author : wyg                       
* +------------------------------------+
* | Date   : 2019/11/20                   
* +------------------------------------+
* | Time   : 11:49                   
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\EnterprisePlugin;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class GuideMiniListController extends BaseController
{
    const CHANNEL_TYPE = 5;

    /**
     * 企业微信小程序列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 查询参数组装
        $params = [
            'qy_domain' => $this->qyDomain,
            'params'    => [
                'enterprise_id' => $this->enterpriseId,
                'channel_type'  => self::CHANNEL_TYPE
            ],
        ];

        Log::debug('guide-mini-params:' . var_export($params, true));

        $result = (new EnterprisePlugin())->miniAppList($params);

        Log::debug('guide-mini-result:' . var_export($result, true));

        return success(['list' => $result['data'] ?? []]);
    }
}