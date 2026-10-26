<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/29
 * Time: 11:26 AM
 */

namespace App\Http\Controllers\V1\Admin;

use App\Exceptions\Error;
use CrsSdk\Boss\EnterpriseAdmin;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AdminRequest;

class PasswordController extends BaseController
{
    /**
     * 管理员验证器
     *
     * @var AdminRequest
     */
    protected $adminRequest;

    /**
     * PasswordController constructor.
     * @param AdminRequest $adminRequest
     */
    public function __construct(AdminRequest $adminRequest)
    {
        parent::__construct();

        // 管理员验证器赋值
        $this->adminRequest = $adminRequest;
    }

    /**
     * 管理员密码修改
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        $params = $request->only(['old_password', 'password' ,'re_password']);

        // 请求数据验证
        $validator = $this->adminRequest->passwordUpdate($params);
        if ($validator->fails()) {
            throw new Error(451000, $validator->errors()->first());
        }

        //确认密码判断
        if ($params['password'] != $params['re_password']) {
            error(451010);
        }

        try {
            (new EnterpriseAdmin())->password([
                'admin_id' => $this->enterpriseAdminId,
                'old_password' => $params['old_password'],
                'password' => $params['password'],
                'enterprise_id' => $this->enterpriseId
            ]);

            // 日志埋点
            $this->writeOperatorLog("重置管理员密码");

        } catch (Exception $e) {
            throw new Error(450000, $e->getMessage());
        }

        return success();
    }
}
