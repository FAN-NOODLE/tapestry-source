<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\CouponUser;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Models\Task;
use App\Services\MemberUserRelationService;
use App\Services\TaskMemberRecordService;
use CrsSdk\Shop\CouponUser;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use App\Http\Requests\CouponUserRequest;

class UserController extends BaseController
{
    /**
     * 用户礼劵验证器
     *
     * @var CouponUserRequest
     */
    protected $couponUserRequest;

    /**
     * 标记外部联系Service
     *
     * @var MemberUserRelationService
     */
    protected $memberUserRelationService;

    /**
     * 记录Service
     *
     * @var TaskMemberRecordService
     */
    protected $taskMemberRecordService;

    /**
     * UserController constructor.
     * @param CouponUserRequest $couponUserRequest
     * @param MemberUserRelationService $memberUserRelationService
     * @param TaskMemberRecordService $taskMemberRecordService
     */
    public function __construct(
        CouponUserRequest $couponUserRequest,
        MemberUserRelationService $memberUserRelationService,
        TaskMemberRecordService $taskMemberRecordService
    )
    {
        parent::__construct();

        // 用户礼劵验证器赋值
        $this->couponUserRequest = $couponUserRequest;
        // 标记外部联系Service赋值
        $this->memberUserRelationService = $memberUserRelationService;
        // 记录Service赋值
        $this->taskMemberRecordService = $taskMemberRecordService;
    }

    /**
     * 礼劵领取用户
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/25 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->input();

        // 请求数据验证
        $validator = $this->couponUserRequest->user($params);
        if ($validator->fails()) {
            throw new Error(400000, $validator->errors()->first());
        }

        //基本查询条件
        $where['enterprise_id'] = $this->enterpriseId;
        $where['coupon_id']     = $params['coupon_id'];
        $where['mem_uid']       = $this->memUid;
        $where['page']          = isset($params['page']) ? $params['page'] : 1;
        $where['limit']         = isset($params['limit']) ? $params['limit'] : 10;

        //是否是客户条件判断
        $userRelation = $this->memberUserRelationService->userRelationByMemUid($this->enterpriseId, $this->memUid);
        if (!empty($userRelation)) {
            $where['union_id'] = array_column($userRelation, 'union_id');
            try {
                $data = (new CouponUser())->user($where);
            } catch (\Exception $e) {
                return error(710002);
            }

            //列表处理
            $list = [];
            array_walk($data['data'], function ($value, $key) use (&$list, $userRelation, $params) {
                //客户信息查询
                $userRelationList = array_column($userRelation, NULL, 'union_id');

                //基本数据组合
                $list[$key]['id']         = $value['id'];
                $list[$key]['user_id']    = $value['user_id'];
                $list[$key]['union_id']   = $value['union_id'];
                $list[$key]['create']     = date('Y-m-d', $value['created_at']);
                $list[$key]['avatar']     = $value['avatar'];
                $list[$key]['user_level'] = $value['user_level'];
                $list[$key]['level_id']   = $value['level_id'];
                $list[$key]['level_num']  = $value['level_num'];
                $list[$key]['ec_userid']  = $userRelationList[$value['union_id']]['ec_userid'];

                // 用户名称处理
                if (isset($userRelationList[$value['union_id']]['user_remark_name']) && !empty($userRelationList[$value['union_id']]['user_remark_name'])) {
                    $list[$key]['user_name'] = $userRelationList[$value['union_id']]['user_remark_name'];
                } elseif (isset($value['real_name']) && !empty($value['real_name'])) {
                    $list[$key]['user_name'] = $value['real_name'];
                } else {
                    $list[$key]['user_name'] = $value['nick_name'];
                }

                // 状态处理
                $taskMemberRecord     = $this->taskMemberRecordService->detail([
                    'enterprise_id' => $this->enterpriseId,
                    'mem_uid' => $this->memUid,
                    'customer_id' => $list[$key]['ec_userid'],
                    'task_type' => Task::MEMBER_COUPON,
                ]);
                $list[$key]['status'] = empty($taskMemberRecord) ? Constant::TASK_Not_REMIND_STATUS : Constant::TASK_REMIND_STATUS;
            });
            $data['data'] = $list;
        } else {
            $data = [
                'data' => [],
                'meta' => [
                    'total' => 0,
                    'size' => $where['page'],
                    'page' => $where['limit']
                ]
            ];
        }
        return success($data);
    }
}
