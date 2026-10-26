<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/28
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\OrderService;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;

class RankingController extends BaseController
{
    /**
     * @var OrderService
     */
    protected $orderService;

    /**
     * RankingController constructor.
     * @param OrderService $orderService
     */
    public function __construct(OrderService $orderService)
    {
        parent::__construct();
        $this->orderService = $orderService;
    }

    /**
     * 获取当前部门人员列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     * @throws \Illuminate\Validation\ValidationException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:25 AM
     */
    public function index(Request $request)
    {
        // 接值
        $input = $request->only(['type', 'date_type', 'sort', 'page', 'size']);

        // 参数验证规则
        $rules = [
            'type' => Rule::in([Constant::ORDER_SORT, Constant::ORDER_NUMBER]),
            'date_type' => Rule::in([Constant::DAY, Constant::WEEK, Constant::MONTH]),
            'sort_type' => Rule::in(['desc', 'asc'])
        ];

        // 提示语
        $message = [
            'type.in' => 'type参数不合法',
            'date_type.in' => 'date_type参数不合法',
            'sort_type.in' => 'sort参数不合法'
        ];

        $this->validate($request, $rules, $message);

        // 分页参数
        $page = $input['page'] ?? Constant::DEFAULT_PAGE;
        $size = $input['size'] ?? Constant::DEFAULT_SIZE;

        // 组装查询UC参数
        $input['qyDomain'] = $this->qyDomain;
        $input['appIdentifier'] = $this->appIdentifier;
        $input['memUid'] = $this->memUid;
        $input['enterpriseId'] = $this->enterpriseId;

        // 获取导购销售排行列表
        $list = $this->orderService->orderRanking($input);
   
        // 总数
        $total = count($list['list']);

        // 数组分页
        $list['list'] = array_slice($list['list'], ($page - 1) * $size, $size);

        // 数据返回
        return success([
            'data' => [
                'list' => array_values($list['list']),
                'order_total' => $list['order_total'],
                'price_total' => $list['price_total'],
                'start_time' => $list['start_time'],
                'end_time' => $list['end_time'],
            ],
            'meta' => [
                'total' => $total,
                'page' => intval($page),
                'size' => intval($size)
            ]
        ]);
    }
}
