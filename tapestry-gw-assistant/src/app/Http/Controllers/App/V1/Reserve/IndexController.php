<?php

namespace App\Http\Controllers\App\V1\Reserve;

use App\Http\Controllers\Controller;
use App\Micros\Common;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Support\Facades\Redis;
use Illuminate\Http\{JsonResponse, Request};

class IndexController extends Controller
{
    /**
     * @var Sdk
     */
    protected $microScrm;

    public function __construct()
    {
        parent::__construct();
        $this->microScrm = Common::microScrm();
    }

    /**
     * 获取预约单详情信息
     * @param Request $request
     * @return JsonResponse
     */
    public function getInfo(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/user-order-info')->post($params);
        } catch (Exception $e) {
            if ($e->getMessage() == "无权访问") {
                error(10001, "抱歉，您非邀约客户，暂无法查看预约内容");
            }
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 编辑预约单
     * @param Request $request
     * @return JsonResponse
     */
    public function editOrder(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/user-edit-order')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 取消预约单
     * @param Request $request
     * @return JsonResponse
     */
    public function cancelOrder(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/user-cancel-order')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 微信用户打分
     * @param Request $request
     * @return JsonResponse
     */
    public function makeScore(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/user-make-score')->post($params);
        } catch (Exception $e) {
            if ($e->getMessage() == "无权访问") {
                error(10001, "权限不足");
            }
            error(10002, $e->getMessage());
        }
        return success($result);
    }


    /**
     * 预约单列表
     * @return JsonResponse
     */
    public function orderList(): JsonResponse
    {
        $params                  = $this->request->all();
        $tokenData               = $this->request->get('tokenData');
        $params['unionid']       = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['status']        = $params['status'] ?? -1;
        $params['order_sort']    = 1;
        if (isset($params['start_time']) && $params['start_time']) {
            // 获取当前时间戳的凌晨时间
            $start                        = strtotime(date("Y-m-d", $params['start_time']));
            $params['reserve_start_time'] = $start;
            $params['reserve_end_time']   = $start + 86399;
            unset($params['start_time']);
        }
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/lists')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 预约单分配
     * @return JsonResponse
     */
    public function orderAssign(): JsonResponse
    {
        $params                  = $this->request->all();
        $tokenData               = $this->request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order-assign')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 时间段邀约列表
     * @param Request $request
     * @return JsonResponse
     */
    public function timeList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $this->request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        if (empty($params['member_id'])) {
            error(900000, "导购ID不能为空");
        }
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/time-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 获取用户的评分详情
     * @param Request $request
     * @return JsonResponse
     */
    public function scoreInfo(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $this->request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['union_id']      = $tokenData['union_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/user-score-info')->post($params);
        } catch (Exception $e) {
            if ($e->getMessage() == "无权访问") {
                error(10001, "您还没有创建到店邀约服务");
            }
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function serviceList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $this->request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $result                  = [];
        $params['status']        = 1; // 有效的
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/service/list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * @param Request $request
     * @return JsonResponse
     */
    public function signIn(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/user-sign-in')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 用户编辑评论
     * @param Request $request
     * @return JsonResponse
     */
    public function editScore(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/user-edit-score')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 针对于SW的活动，做一些特定处理，
     * @param Request $request
     * @return JsonResponse
     */
    public function getReserveActivity(Request $request) : JsonResponse
    {
        $reserve_info= Redis::connection()->get("sw_reserve_activity_id");

        $info=[];
        if (!empty($reserve_info)){
            $info=json_decode($reserve_info,true);
        }
        $data=["data"=>$info];
        return success($data);
    }
}