<?php
/**
 * Created by PhpStorm.
 * User: zhaoqi
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseAppService;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * DetailController constructor.
     *
     * @param EnterpriseAdminService $enterpriseAdminService
     *
     * @throws \Exception
     */
    public function __construct(EnterpriseAdminService $enterpriseAdminService)
    {
        parent::__construct();

        // 企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 企业管理员
     *
     * @param Request $request
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @author zhaoqi zhaoqi@vchangyi.com
     * @date   2018/11/25 11:16
     */
    public function index(Request $request)
    {
        $params = $request->only(['admin_id', 'account', 'member_id', 'enterprise_id', 'mobile']);

        // 查询信息不能为空
        if (isset($params['admin_id']) && !empty($params['admin_id'])) {
            $enterpriseAdmin = $this->enterpriseAdminService->adminDetail(['id' => $params['admin_id']]);
        } elseif (isset($params['account']) && !empty($params['account'])) {
            $enterpriseAdmin = $this->enterpriseAdminService->adminDetail(['account' => $params['account']]);
        } elseif (isset($params['enterprise_id']) && isset($params['mobile'])) {
            $enterpriseAdmin = $this->enterpriseAdminService->adminDetail([
                'enterprise_id' => $params['enterprise_id'],
                'mobile'        => $params['mobile']
            ]);
        } elseif (isset($params['enterprise_id']) && isset($params['member_id']) && !empty($params['enterprise_id']) && !empty($params['member_id'])) {
            $enterpriseAdmin = $this->enterpriseAdminService->adminDetail([
                'member_id'     => $params['member_id'],
                'enterprise_id' => $params['enterprise_id']
            ]);
        } else {
            return error(202104);
        }

        if (!$enterpriseAdmin) {
            return error(202007);
        }

        // logo
        $logo = EnterpriseThirdPartAccount::query()
                ->where([
                    'enterprise_id' => $enterpriseAdmin->enterprise_id,
                    'status'        => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 授权状态：有效
                    'platform'      => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT, // 授权状态：有效
                    'type'          => EnterpriseThirdPartAccount::TYPE_SHOP // 授权状态：有效
                ])->value('logo') ?? '';

        $enterpriseAppService = new EnterpriseAppService();
        $data                 = [
            'id'              => $enterpriseAdmin->id,
            'enterprise_id'   => $enterpriseAdmin->enterprise_id,
            'username'        => $enterpriseAdmin->username,
            'account'         => $enterpriseAdmin->account,
            'password'        => $enterpriseAdmin->password,
            'mobile'          => $enterpriseAdmin->mobile,
            'salt'            => $enterpriseAdmin->salt,
            'last_login_time' => $enterpriseAdmin->last_login_time,
            'role'            => $enterpriseAdmin->role,
            'status'          => $enterpriseAdmin->status,
            'mobile'          => $enterpriseAdmin->mobile,
            'is_super'        => $enterpriseAdmin->is_super,
            'enterprise'      => [
                'id'             => $enterpriseAdmin->enterprise->id,
                'hash'           => $enterpriseAdmin->enterprise->hash,
                'name'           => $enterpriseAdmin->enterprise->name,
                'logo'           => $logo,
                'contact'        => $enterpriseAdmin->enterprise->contact,
                'mobile'         => $enterpriseAdmin->enterprise->mobile,
                'address'        => $enterpriseAdmin->enterprise->address,
                'crs_status'     => $enterpriseAdmin->enterprise->crs_status,
                'store_status'   => $enterpriseAdmin->enterprise->store_status,
                'contact_status' => $enterpriseAdmin->enterprise->contact_status,
                'corp_id'        => $enterpriseAdmin->enterprise->corp_id,
                'domain'         => $enterpriseAdmin->enterprise->domain,
                'status'         => $enterpriseAdmin->enterprise->status,
                'is_split_flow'  => $enterpriseAdmin->enterprise->is_split_flow,
                'app_list'       => $enterpriseAppService->enterpriseAppList($enterpriseAdmin->enterprise->id),
            ]
        ];

        return success(['data' => $data]);
    }
}
