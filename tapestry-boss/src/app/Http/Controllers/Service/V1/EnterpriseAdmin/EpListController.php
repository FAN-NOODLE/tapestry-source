<?php
/**
 * Created by PhpStorm.
 * User: zhaoqi
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use App\Services\EnterpriseAppService;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\EnterpriseAdminService;

class EpListController extends BaseController
{
    /**
     * 企业管理员service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * DetailController constructor.
     * @param EnterpriseAdminService $enterpriseAdminService
     * @throws \Exception
     */
    public function __construct(EnterpriseAdminService $enterpriseAdminService)
    {
        parent::__construct();

        // 企业管理员service赋值
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 根据手机号获取管理员所在的企业
     * @param Request $request
     * @return array|void
     */
    public function index(Request $request)
    {
        $params = $request->only(['account']);

        // 查询信息不能为空
        if (isset($params['account']) && !empty($params['account'])) {
            $adminList = $this->enterpriseAdminService->getAdminListByMobile([
                'account' => $params['account']
            ]);
        } else {
            return error(202104);
        }

        $data = [];
        $enterpriseAppService = new EnterpriseAppService();
        foreach ($adminList as $admin) {
            $data[] = [
                'id'             => $admin['enterprise']['id'],
                'hash'           => $admin['enterprise']['hash'],
                'name'           => $admin['enterprise']['name'],
                'logo'           => $admin['enterprise']['logo'],
                'contact'        => $admin['enterprise']->contact,
                'mobile'         => $admin['enterprise']->mobile,
                'address'        => $admin['enterprise']->address,
                'crs_status'     => $admin['enterprise']->crs_status,
                'store_status'   => $admin['enterprise']->store_status,
                'contact_status' => $admin['enterprise']->contact_status,
                'corp_id'        => $admin['enterprise']->corp_id,
                'domain'         => $admin['enterprise']->domain,
                'status'         => $admin['enterprise']->status,
                'is_split_flow'  => $admin['enterprise']->is_split_flow,
                'app_list'       => $enterpriseAppService->enterpriseAppList($admin['enterprise']['id']),
                'username'       => $admin['username'],
            ];
        }

        return success(['data' => $data]);
    }
}
