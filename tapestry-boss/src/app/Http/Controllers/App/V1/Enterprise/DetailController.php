<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\Enterprise;

use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterpriseService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseAppService;

class DetailController extends BaseController
{
    /**
     * 企业service
     *
     * @var EnterpriseService
     */
    protected $enterpriseService;

    /**
     * 企业admin service
     *
     * @var EnterpriseAdminService
     */
    protected $enterpriseAdminService;

    /**
     * DetailController constructor.
     *
     * @param EnterpriseService      $enterpriseService
     * @param EnterpriseAdminService $enterpriseAdminService
     */
    public function __construct(EnterpriseService $enterpriseService, EnterpriseAdminService $enterpriseAdminService)
    {
        parent::__construct();

        // 企业service赋值
        $this->enterpriseService = $enterpriseService;
        $this->enterpriseAdminService = $enterpriseAdminService;
    }

    /**
     * 企业详情
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @author 赵奇
     * @date   2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id']);
        // 企业ID不能为空
        if (!isset($params['enterprise_id']) || empty($params['enterprise_id'])) {
            return notice(201101);
        }

        // 企业信息查询
        $enterprise = $this->enterpriseService->enterprise($params['enterprise_id']);
        if (!$enterprise) {
            return notice(201002);
        }

        // 查询超管信息
        $enterpriseAdmin = $this->enterpriseAdminService->superAdmin($params['enterprise_id']);
        // 查询企业应用信息
        $enterpriseApp = (new EnterpriseAppService())->enterpriseAppList($params['enterprise_id']);

        return success(['data' => $this->enterpriseDetail($enterprise, $enterpriseAdmin, $enterpriseApp)]);
    }

    /**
     * 企业数据拼装
     *
     * @author 赵奇
     *
     * @param resource $enterprise      企业数据
     * @param resource $enterpriseAdmin 企业超管数据
     * @param mixed    $enterpriseApp   企业应用数据
     *
     * @return array $data 拼装后数据
     */
    private function enterpriseDetail($enterprise, $enterpriseAdmin, $enterpriseApp)
    {
        // 转数组
        $enterprise = !empty($enterprise) ? json_decode($enterprise, true) : [];
        $enterpriseAdmin = !empty($enterpriseAdmin) ? json_decode($enterpriseAdmin, true) : [];
        $enterpriseApp = !empty($enterpriseApp) ? json_decode($enterpriseApp, true) : [];

        // 拼装数据
        $data = [
            'id'             => isset($enterprise['id']) ? $enterprise['id'] : '',
            'hash'           => isset($enterprise['hash']) ? $enterprise['hash'] : '',
            'name'           => isset($enterprise['name']) ? $enterprise['name'] : '',
            'logo'           => isset($enterprise['logo']) ? $enterprise['logo'] : '',
            'contact'        => isset($enterprise['contact']) ? $enterprise['contact'] : '',
            'mobile'         => isset($enterprise['mobile']) ? $enterprise['mobile'] : '',
            'address'        => isset($enterprise['address']) ? $enterprise['address'] : '',
            'crs_status'     => isset($enterprise['crs_status']) ? (int)$enterprise['crs_status'] : '',
            'store_status'   => isset($enterprise['store_status']) ? (int)$enterprise['store_status'] : '',
            'contact_status' => isset($enterprise['contact_status']) ? $enterprise['contact_status'] : '',
            'corp_id'        => isset($enterprise['corp_id']) ? $enterprise['corp_id'] : '',
            'domain'         => isset($enterprise['domain']) ? $enterprise['domain'] : '',
            'status'         => isset($enterprise['status']) ? $enterprise['status'] : '',
            'admin_id'       => isset($enterpriseAdmin['id']) ? $enterpriseAdmin['id'] : '',
            'username'       => isset($enterpriseAdmin['username']) ? $enterpriseAdmin['username'] : '',
            'account'        => isset($enterpriseAdmin['account']) ? $enterpriseAdmin['account'] : '',
            'admin_status'   => isset($enterpriseAdmin['status']) ? (int)$enterpriseAdmin['status'] : '',
            "app"            => !empty($enterpriseApp) ? $enterpriseApp : []
        ];

        return $data;
    }
}
