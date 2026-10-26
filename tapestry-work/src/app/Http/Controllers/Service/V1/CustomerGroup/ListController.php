<?php


namespace App\Http\Controllers\Service\V1\CustomerGroup;


use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerGroupService;

class ListController extends BaseController
{
    public function index()
    {
        $input = $this->request->only([
            'enterprise_id',
            'page',
            'size',
            'keyword',
            'select'
        ]);

        // 验证
        if (!isset($input['enterprise_id']) || empty($input['enterprise_id'])) {
            notice(600000);
        }

        $input['keyword'] = isset($input['keyword']) ? trim($input['keyword']) : '';
        $input['select'] = $input['select'] ?? 0; //0-按照分页模式搜索并计算分组人数 1-分页模式搜索不计算分组人数 2-不分页不计算分组人数

        // 查询用户分组列表
        $list = (new CustomerGroupService())->ListCustomerGroup($input);

        // 返回分页的情况
        if ($input['select'] == 0 || $input['select'] == 1) {
            return success([
                'data' => $list,
                'meta' => [
                    'page'  => $input['page'] ?? Constant::PAGE,
                    'size'  => $input['size'] ?? Constant::SIZE,
                    'total' => (new CustomerGroupService())->customerGroupTotal($input)
                ]
            ]);
        }

        return success(['data' => $list]);
    }

}