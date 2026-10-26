<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 19-2-25
 * Time: 下午2:25
 */

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use Illuminate\Http\Request;
use App\Models\Enterprise;
use App\Services\EnterpriseAppService;
use App\Http\Requests\EnterpriseRequest;

class EnterpriseInfo extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    public function index(Request $request)
    {
        $data = $request->only('enterprise_id');

        // 请求数据验证
        $validator = (new EnterpriseRequest())->enterpriseInfo($data);
        if ($validator->fails()) {
            return notice($validator->errors()->first(), [], $data);
        }

        $enterprise = Enterprise::query()
            ->where('id', $data['enterprise_id'])
            ->first(['id', 'crs_status', 'store_status']);

        if (!empty($enterprise)) {
            $enterprise = $enterprise->toArray();
        }

        $installApp = (new EnterpriseAppService())->enterpriseAppList($data['enterprise_id']);
        if (!empty($installApp)) {
            $installApp = array_column($installApp->toArray(), 'status', 'identity');
        }

        $enterprise = array_merge($enterprise, $installApp);

        return success(['data' => $enterprise]);
    }
}
