<?php

namespace App\Http\Controllers\App\V1\MarketingPlan;

use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Micros\Common;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class PlansController extends Controller
{
    /**
     * coach 任务中心的日历列表接口
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/month-lists')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 任务中心的主任务列表
     * @param Request $request
     * @return JsonResponse
     */
    public function lists(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/lists')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 获取主任务详情
     * @param Request $request
     * @return JsonResponse
     */
    public function detail(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/detail')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 的任务中心的未完成的任务数提示
     * @param Request $request
     * @return JsonResponse
     */
    public function tipsInfo(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/tips')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 计划任务退群的顾客列表
     * @param Request $request
     * @return JsonResponse
     */
    public function leaveGroupChat(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/leave-group-chat')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach sacode&&unionid|customerid的导入客户列表
     * @param Request $request
     * @return JsonResponse
     */
    public function saCustomerList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/sa-customer-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 计划任务上报点击电话按钮
     * @param Request $request
     * @return JsonResponse
     */
    public function confirmCall(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/confirm-call')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 计划任务点击发消息按钮
     * @param Request $request
     * @return JsonResponse
     */
    public function confirmContact(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/confirm-contact')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 计划任务点击标记完成
     * @param Request $request
     * @return JsonResponse
     */
    public function confirmFinish(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/confirm-finish')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 计划任务设置联系说明
     * @param Request $request
     * @return JsonResponse
     */
    public function setRemark(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/set-remark')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 计划任务子任务详情点击转发给顾客按钮的上报
     * @param Request $request
     * @return JsonResponse
     */
    public function setForward(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/set-forward')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }


    /**
     * 点击添加好友按钮的上报
     * @param Request $request
     * @return JsonResponse
     */
    public function confirmFriend(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/confirm-friend')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 主任务关闭未触达的提示
     * @param Request $request
     * @return JsonResponse
     */
    public function closeTouchTips(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/close_touch_tips')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 子任务已完成关闭未触达状态
     * @param Request $request
     * @return JsonResponse
     */
    public function closeTaskTouchTips(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/task/close_touch_tips')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }


    /**
     * coach 新增计划任务顾客列表接口
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function searchCustomerList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/search-customer-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 新增计划任务顾客列表接口
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function searchExternalList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/search-external-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach升级的计划任务顾客列表,全选获取openid
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function searchExternalAllOpenId(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/search-external-all-openid')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     *
     * coach 计划任务顾客列表移动接口
     * @param Request $request
     * @return JsonResponse
     */
    public function customerMove(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/customer-move')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 企业沟通任务转发上报接口
     * @param Request $request
     * @return JsonResponse
     */
    public function sendFinish(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/send-finish')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    public function customerTaskList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['user_id'] = $this->memberAuthInfo['userid'];
        $result = [];
        try {
            $result = Common::microScrm()->setApi('micro-scrm/v1/plan/customer-task-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }
}