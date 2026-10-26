<?php

namespace App\Http\Controllers\V1\Department;

use App\Micros\Work\Work;
use CrsSdk\Work\Member;
use App\Common\Constant;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\Department;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Services\DepartmentService;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 部门列表
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        // dp_id:部门id dp_level:部门级别 dp_name:部门名称
        $input = $request->input();
        if (empty($this->qyDomain)) {
            return success(['data' => []]);
        }
        // 查询参数
        $params = [
            'pageNum'   => Constant::DEFAULT_PAGE,
            'pageSize'  => 1,
            'dpLevel'   => 1,
            'orderList' => [
                [
                    'column'    => 'dpDisplayorder',
                    'orderType' => 'ASC'
                ]
            ]
        ];

        // 存在部门id，查子级
        if (isset($input['dp_id']) && !empty($input['dp_id'])) {
            $params['dpId'] = $input['dp_id'];
            unset($params['dpLevel']);
        }
        // 名称搜索
        if (isset($input['dp_name']) && !empty($input['dp_name'])) {
            $params['dpName'] = $input['dp_name'] . '%';
            unset($params['dpLevel']);
        }
        // 分页
        if (isset($input['page_num']) && !empty($input['page_num'])) {
            $params['pageNum'] = $input['page_num'];
        }

        // 调用SDK
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => $params
        ];
        $noDpData  = [];
        try {
            $Department = new Department();
            $response   = $Department->departmentList($condition);

            $memberCondition = [
                'qy_domain'      => $this->qyDomain,
                'app_identifier' => $this->appIdentifier,
                'enterprise_id'  => $this->enterpriseId,
                'params'         => []
            ];

            // 如果是查询顶级部门且没有顶级部门
            if (!empty($input['is_top']) && ($response['total'] > 1 || $response['total'] == 0)) {

                $memberList = (new Member())->memberList($memberCondition);

                // 总数
                $total = intval($memberList['total']);
                if ($response['total'] == 0 && empty($total)) {
                    return success(['data' => []]);
                }

                $paramsEnterprise = [
                    'enterprise_id' => $this->enterpriseId
                ];
                // 获取企业信息
                $detail = (new Enterprise())->detail($paramsEnterprise);

                $condition = [
                    'qy_domain'      => $this->qyDomain,
                    'app_identifier' => $this->appIdentifier,
                    'params'         => ['dpLevel' => 1]
                ];
                // 查询部门信息
                $response = $Department->departmentList($condition);

                $data[] = [
                    'dp_id'               => 0,
                    'dp_name'             => $detail['data']['name'],
                    'dp_parent_id'        => '',
                    'is_child_department' => !empty($response['total']) ? 1 : 0,
                    'is_store'            => 0,
                    'dp_user_count'       => 0,
                    'dp_display_order'    => 0,
                ];

                return success(['list' => $data]);
            }
            if ($response['total'] > 1) {
                // 查询数量重新赋值
                $params['pageSize'] = $response['total'];
                $condition          = [
                    'qy_domain'      => $this->qyDomain,
                    'app_identifier' => $this->appIdentifier,
                    'params'         => $params
                ];
                // 调用SDK
                $response = $Department->departmentList($condition);
            }


        } catch (\Exception $e) {
            return success(['data' => []]);
        }

        // 调用SDK获取门店dptId
        $dptIdData = $Department->storeDptId([
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier
        ]);

        // 部门门店类型id
        $dptId = $dptIdData['data']['dptId'] ?? '';

        // 数据格式化
        $DepartmentService = new DepartmentService();
        $data              = $DepartmentService->formatData($response['list'], $dptId);

        $data = array_merge($data, $noDpData);

        return success(['list' => $data]);
    }

    /**
     * 查询部门列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function thirdList(Request $request): JsonResponse
    {
        $params = $request->input();
        if (empty($params['dp_id']) && empty($params['dp_name'])) {
            $params['third_parent_id'] = $params['third_parent_id'] ?? 0;
        }
        $departmentLists = Work::lists(['enterprise_id' => $this->enterpriseId, 'params' => $params]);

        return success(['list' => $departmentLists]);
    }
}
