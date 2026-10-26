<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 2019-02-20
 * Time: 18:47
 */

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseAdmin;
use Illuminate\Support\Facades\Log;
use App\Http\Requests\EnterpriseAdminRequest;

class StatusController extends BaseController
{
    /**
     * StatusController constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 增加企业管理员
     *
     * @author zhaoqi zhaoqi@vchangyi.com
     * @date   2019/02/20 15:26
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'admin_id', 'status']);

        $validator = (new EnterpriseAdminRequest())->status($params);
        if ($validator->fails()) {
            return notice($validator->errors()->first(), [], $params);
        }

        try {
            EnterpriseAdmin::query()
                ->where([
                    'id'            => $params['admin_id'],
                    'enterprise_id' => $params['enterprise_id']
                ])->where('role', '!=', EnterpriseAdmin::SUPER_ADMIN)
                ->update(['status' => $params['status']]);
        } catch (\Exception $e) {
            Log::error($e->getMessage());

            return notice(202109, [], $params);
        }

        return success([]);
    }
}
