<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\Service\V1\Enterprise;

use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\EnterpriseService;

class SaveController extends BaseController
{
    /**
     * 企业service
     *
     * @var EnterpriseService
     */
    protected $enterpriseService;

    /**
     * SaveController constructor.
     *
     * @param EnterpriseService $enterpriseService
     */
    public function __construct(EnterpriseService $enterpriseService)
    {
        parent::__construct();

        //企业service赋值
        $this->enterpriseService = $enterpriseService;
    }

    /**
     * 编辑企业
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'corp_id', 'domain', 'member_id']);

        // 企业ID不能为空
        if (!isset($params['enterprise_id']) || empty($params['enterprise_id'])) {
            return error(201101);
        }

        // 企业信息参数处理
        $enterpriseId = $params['enterprise_id'];
        unset($params['enterprise_id']);
        if (empty($params)) {
            return error(201003);
        }

        // 企业信息保存
        $result = $this->enterpriseService->enterpriseSave($enterpriseId, $params);
        if (!$result) {
            return error(201004);
        }

        return success(['data' => []]);
    }
}
