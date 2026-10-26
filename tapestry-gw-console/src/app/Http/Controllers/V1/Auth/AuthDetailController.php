<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 19-3-5
 * Time: 上午11:45
 */

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;
use Exception;
use Illuminate\Http\Request;

class AuthDetailController extends BaseController
{
    /**
     * AuthDetailController constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * @param Request $request
     * @return mixed
     * @throws Exception
     */
    public function index(Request $request)
    {
        $params = $request->only('type');

        if (empty($params['type']) || !is_int($params['type']) || $params['type'] < 1) {
            notice(451025, [], $params);
        }

        $params['enterprise_id'] = $this->enterpriseId;

        // 日志埋点
        $this->writeOperatorLog("查看小程序信息详情");

        $result = (new Enterprise())->authDetail($params);

        return success($result['data'] ?? []);
    }
}
