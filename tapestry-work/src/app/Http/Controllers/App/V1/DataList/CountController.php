<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/3/19
 * Time: 16:51
 */

namespace App\Http\Controllers\App\V1\DataList;


use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\Order;
use App\Services\MemberUserRelationService;
use App\Common\Constant;

class CountController extends BaseController
{
    /**
     * 本日
     */
    const DAY = 1;
    /**
     * 本月
     */
    const MONTH = 2;
    /**
     * 本季度
     */
    const QUARTER = 3;
    /**
     * 本年
     */
    const YEAR = 4;
    /**
     * 按时间
     */
    const BY_TIME = 5;

    /**
     * @var MemberUserRelationService
     */
    protected $memberUserRelationService;

    /**
     * CountController constructor.
     *
     * @param MemberUserRelationService $memberUserRelationService
     */
    public function __construct(
        MemberUserRelationService $memberUserRelationService
    ) {
        parent::__construct();
        $this->memberUserRelationService = $memberUserRelationService;
    }

    /**
     * 数据统计接口
     *
     * @param Request $request
     * @return mixed|\Psr\Http\Message\StreamInterface
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        // 接值
        $input = $request->only('type', 'start_time', 'end_time');
        $month = intval(date('m'));

        // 参数验证规则
        $rules = [
            'type'  => 'required | integer',
        ];
        // 提示语
        $message = [
            'type.required'  => 'type不能为空',
            'type.integer'  => 'type必须是正整数'
        ];
        $this->validate($request, $rules, $message);

        // 初始化
        $start = 0;
        $end = 0;
        switch ($input['type']) {
            case self::DAY:
                // 本日
                $start = mktime(0,0,0,date('m'),date('d'),date('Y'));
                $end = time();
                break;
            case self::MONTH:
                // 本月
                $start = mktime(0,0,0,date('m'),1,date('Y'));
                $end = time();
                break;
            case self::QUARTER:
                // 本季度
                if ($month < 4) {
                    $start = mktime(0,0,0,1,1,date('Y'));
                }elseif ($month < 7) {
                    $start = mktime(0,0,0,4,1,date('Y'));
                }elseif ($month < 10) {
                    $start = mktime(0,0,0,7,1,date('Y'));
                }else {
                    $start = mktime(0,0,0,10,1,date('Y'));
                }
                $end = time();
                break;
            case self::YEAR:
                // 本年
                $start = mktime(0,0,0,1,1,date('Y'));
                $end = time();
                break;
            case self::BY_TIME:
                // 按时间
                $start = isset($input['start_time']) ? $input['start_time'] / 1000 : 0;
                $end = isset($input['end_time']) ? $input['end_time'] / 1000 : 0;
                break;
        }

        $params = [
            'enterprise_id' => $this->enterpriseId,
            'mem_uid' => $this->memUid,
            'start_time' => $start,
            'end_time' => $end
        ];
        // 开通了商城
        if ($this->storeStatus == Constant::STORE_OPEN) {
            $list = (new Order())->orderByTime($params);
            // 封装数据
            if (empty($list['order_price'])) {
                $list['order_price'] = 0;
            }
            if (empty($list['order_num'])) {
                $list['order_num'] = 0;
            }
            // 百分比
            $scanP = 100;
            if (empty($list['scan_person_num'])) {
                $orderP = 0;
                $doneP = 0;
            }else {
                $orderP = ($list['order_person_num'] / $list['scan_person_num']) * 100 ;
                $doneP = ($list['done_person_num'] / $list['scan_person_num']) * 100 ;
            }

            $data['scan'] = ['num' => $list['scan_person_num'], 'percent' => $scanP];
            $data['order'] = ['num' => $list['order_person_num'], 'percent' => $orderP];
            $data['done'] = ['num' => $list['done_person_num'], 'percent' => $doneP];
            $data['order_price'] = $list['order_price'];
            $data['order_num'] = $list['order_num'];
        }

        // 好友客户数
        $userNum = $this->memberUserRelationService->timeUserCount(
            $this->enterpriseId,
            $this->memUid,
            $start,
            $end
        );
        $data['customer_num'] = $userNum;

        return success(['data' => $data]);
    }
}
