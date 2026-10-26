<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Department;

use App\Common\Constant;
use CrsSdk\Work\Department;
use Illuminate\Http\Request;
use App\Services\DepartmentService;
use App\Http\Controllers\V1\BaseController;

class PageListController extends BaseController
{
    /**
     * 部门列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // dp_id:部门id dp_name:部门名称 page:页数(1) size:每页显示记录数(10)
        $input = $request->input();

        if (empty($this->qyDomain)) {
            return success(['data' => []]);
        }
        // 查询参数
        $params = [
            'pageNum'                => Constant::DEFAULT_PAGE,
            'pageSize'               => Constant::DEFAULT_SIZE,
            'departmentChildrenFlag' => 1,
            'orderList'              => [
                [
                    'column'    => 'dpDisplayorder',
                    'orderType' => 'ASC'
                ]
            ]
        ];

        // 存在部门id，查子级
        if (isset($input['dp_id']) && !empty($input['dp_id'])) {
            $params['dpId'] = $input['dp_id'];
        }
        // 名称搜索
        if (isset($input['dp_name']) && !empty($input['dp_name'])) {
            $params['dpName'] = $input['dp_name'] . '%';
        }
        // 分页
        if (isset($input['page']) && !empty($input['page'])) {
            $params['pageNum'] = $input['page'];
        }
        if (isset($input['size']) && $input['size'] > 0) {
            $params['pageSize'] = $input['size'];
        }

        // 调用SDK
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => $params
        ];

        $Department = new Department();
        $response   = $Department->departmentList($condition);

        // 调用SDK获取门店dptId
        $dptIdData = $Department->storeDptId([
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier
        ]);

        // 部门门店类型id
        $dptId = $dptIdData['data']['dptId'] ?? '';

        // 数据格式化
        $data = (new DepartmentService())->formatData($response['list'], $dptId);

        return success([
            'list'  => $data,
            'page'  => $response['pageNum'],
            'size'  => $response['pageSize'],
            'total' => $response['total'],

        ]);
    }

}
