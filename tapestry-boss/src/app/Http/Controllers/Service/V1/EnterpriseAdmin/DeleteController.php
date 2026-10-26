<?php
/**
 * Created by PhpStorm.
 * User: zp
 * Date: 2019-02-20
 * Time: 19:13
 */

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseAdmin;
use Illuminate\Support\Facades\Log;
use App\Http\Requests\EnterpriseAdminRequest;

class DeleteController extends BaseController
{
    /**
     * DeleteController constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 增加企业管理员
     *
     * @param Request $request
     *
     * @return array|mixed
     * @author zhaoqi zhaoqi@vchangyi.com
     * @date   2019/02/20 19:13
     *
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'admin_id', 'account']);

        // 数据验证
        if (empty($params['account'])) {
            $validator = (new EnterpriseAdminRequest())->delete($params);
            if ($validator->fails()) {
                return notice($validator->errors()->first(), [], $params);
            }
        } else {
            if (empty($params['enterprise_id'])) {
                notice(201101);
            }
        }

        try {
            if (!empty($params['account'])) {
                EnterpriseAdmin::query()
                    ->where([
                        'enterprise_id' => $params['enterprise_id'],
                        'account'       => $params['account']
                    ])
                    ->delete();
            } else {
                EnterpriseAdmin::query()
                    ->where([
                        'enterprise_id' => $params['enterprise_id'],
                        'id'            => $params['admin_id']
                    ])->where('role', '!=', EnterpriseAdmin::SUPER_ADMIN)
                    ->delete();
            }
        } catch (\Exception $e) {
            Log::error($e->getMessage());

            return notice(202110, [], $params);
        }

        return success([]);
    }
}
