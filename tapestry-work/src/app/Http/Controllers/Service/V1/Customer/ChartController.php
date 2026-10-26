<?php
/**
 * 员工图表
 * User: 朱虹达
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\MemberUserRelationService;
use Illuminate\Http\Request;

class ChartController extends BaseController
{
    /**
     * 周
     */
    const WEEK = 1;

    /**
     * 月
     */
    const MONTH = 2;

    /**
     * @var MemberUserRelationService
     */
    protected $memberUserRelationService;

    /**
     * ChartController constructor.
     * @param MemberUserRelationService $memberUserRelationService
     */
    public function __construct(MemberUserRelationService $memberUserRelationService)
    {
        parent::__construct();
        $this->memberUserRelationService = $memberUserRelationService;
    }

    /**
     * 入口方法
     *
     * @param Request $request
     * @return array|mixed
     * @throws \Illuminate\Validation\ValidationException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:18 AM
     */
    public function index(Request $request)
    {
        $input = $request->input();

        // 参数验证规则
        $rules = [
            'enterprise_id' => 'required',
            'date_type'     => 'required'
        ];

        // 提示语
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'date_type.required'     => '时间类型不能为空'
        ];

        $this->validate($request, $rules, $message);

        $x = $y = [];

        try {
            if ($input['date_type'] == self::WEEK) {
                // 获取六周之前的周一日期开始时间
                $input['start'] = strtotime("-6 monday", time());

                // 数据列表
                $list = $this->memberUserRelationService->weekChart($input);

                // 获取最近六个周数据
                $date = getWeek();

                if (!empty($list)) {
                    foreach ($list as $value) {
                        $data[$value->xtime] = $value->number;
                    }
                }

                // 循环周组装数据
                foreach ($date as $value) {
                    $sales = $data[$value['date']] ?? 0;
                    $v['x'] = $value['start'].'-'.$value['end'];
                    $v['y'] = $sales;
                    $lists[] = $v;
                }
            } else {
                // 获取六个月之前的1号日期
                $start = strtotime("-5 month");
                $input['start'] = strtotime(date('Y-m-01', $start));

                // 数据列表
                $list = $this->memberUserRelationService->monthChart($input);

                // 获取最近六个月数据
                $date = getMonth();

                if (!empty($list)) {
                    foreach ($list as $value) {
                        $data[$value->xtime] = $value->number;
                    }
                }

                // 循环月份组装数据
                foreach ($date as $value) {
                    $sales = $data[$value['start']] ?? 0;
                    $v['x'] = $value['date'];
                    $v['y'] = $sales;
                    $lists[] = $v;
                }
            }

            $x = array_column($lists, 'x');
            $y = array_column($lists, 'y');

            return success([
                'data' => [
                    'x' => $x,
                    'y' => $y
                ]
            ]);

        } catch (\Exception $e) {
            error(600022);
        }
    }
}
