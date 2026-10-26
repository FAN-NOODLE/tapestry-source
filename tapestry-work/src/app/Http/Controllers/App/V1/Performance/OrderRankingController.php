<?php
/**
 * 业绩榜单
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Performance;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\OrderService;
use \Illuminate\Http\Request;
use Illuminate\Validation\Rule;

class OrderRankingController extends BaseController
{
    /**
     * 是否有用户信息 1：是，2：否 2
     */
    const ORDER_USER_FALSE = 2;

    public function index(Request $request, OrderService $orderService)
    {
        // 接值
        $input = $request->only(['date_type', 'page', 'size']);

        // 参数验证规则
        $rules = [
            'date_type' => Rule::in([Constant::DAY, Constant::WEEK, Constant::MONTH]),
            'page' => 'integer',
            'size' => 'integer',
        ];
        // 提示语
        $message = [
            'date_type.in' => 'date_type参数不合法',
            'page.integer' => '页码必须为正整数',
            'size.integer' => '每页数量必须为正整数'
        ];
        $this->validate($request, $rules, $message);

        // 分页
        $page = $input['page'] ?? Constant::DEFAULT_PAGE;
        $size = $input['size'] ?? Constant::DEFAULT_SIZE;

        // 组装查询UC参数
        $input['qyDomain'] = $this->qyDomain;
        $input['appIdentifier'] = $this->appIdentifier;
        $input['enterpriseId'] = $this->enterpriseId;
        $input['memUid'] = $this->memUid;

        // 获取导购销售排行列表
        $list = $orderService->orderRanking($input);

        // 数据返回
        return success([
            'data' => [
                'list' => array_values(array_slice($list['list'], ($page - 1) * $size, $size)),
                'meta' => [
                    'total' => count($list['list']),
                    'page' => intval($page),
                    'size' => intval($size)
                ],
                'current_data' => $list['list'][$this->memUid] ?? []
            ],
        ]);
    }
}
