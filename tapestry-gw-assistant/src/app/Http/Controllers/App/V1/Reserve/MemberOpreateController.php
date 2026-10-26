<?php

namespace App\Http\Controllers\App\V1\Reserve;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Http\{JsonResponse, Request};

class MemberOpreateController extends BaseController
{
    /**
     * @var Sdk
     */
    protected $microScrm;
    protected $work;

    public function __construct()
    {
        parent::__construct();
        $this->microScrm = Common::microScrm();
        $this->work      = Common::work();
    }

    /**
     * 创建预约单
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function createOrder(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        if (empty($params['user_unionid'])) {
            error(10002, "该用户身份不是微信用户");
        }
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/create')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 时间段邀约列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function timeList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/time-list')->post($params);
        } catch (Exception $e) {
            error(900003);
        }
        return success($result);
    }

    /**
     * 时间段邀约列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function rankList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/rank-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }


    /**
     * 预约单列表
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function orderList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $memberInfo = $this->microScrm->setApi('micro-scrm/v1/member/info')->post([
            'enterprise_id' => $this->enterpriseId,
            'mobile'        => $this->mobile
        ]);
        $params['member_id'] = $memberInfo['id'];
        if (empty($memberInfo)) {
            throw new Error(10003, "用户部门信息丢失");
        }

        // 导购搜索
        // 如果是店员则只能看到自己的数据
        if (!$this->isManager) {
            $params['member_ids'] =(array)$memberInfo['id'];
        }

        // 店长可以筛选门店和门店下的员工数据
        if ($this->isManager) {
            // 初始化标志变量
            $searchAllStores = false;

            // 处理门店筛选逻辑
            if (empty($params['store_ids'])) {
                // 默认展示店长管理的所有门店
                $searchAllStores = true;
                $dataRange = $this->dataRangeList;
                $params['store_ids'] = array_column($dataRange, "object_id");
            }

            // 处理员工筛选逻辑
            $hasMemberFilter = !empty($params['member_ids']);

            if (!$hasMemberFilter) {
                // 没有指定员工筛选时
                if ($searchAllStores) {
                    // 查询所有门店时，添加店长标识
                    $params['is_sm'] = 1;
                }
                $params['member_ids'][] = (int)$params['member_id'];
            } else {
                // 有指定员工筛选时
                if ($searchAllStores) {
                    // 如果筛选了具体员工，则移除门店限制
                    unset($params['store_ids']);
                }
            }
            // 去重处理
            $params['member_ids'] = array_unique($params['member_ids']);
        }
        if (isset($params['start_time']) && $params['start_time']) {
            // 获取当前时间戳的凌晨时间
            $start                        = strtotime(date("Y-m-d", $params['start_time']));
            $params['reserve_start_time'] = $start;
            $params['reserve_end_time']   = $start + 86399;
            unset($params['start_time']);
        }
        if (!isset($params['status'])) {
            $params['status'] = -1;
        }
        $params["order_sort"] = 1;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/lists')->post($params);
        } catch (Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
        return success($result);
    }


    /**
     * 预约单列表
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function orderListByCoach(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        //门店搜索
        if (empty($params['third_ids'])) {
            $params = array_merge($params, $this->getRoleRange());
        }
        //导购搜索 $params['member_ids']
        //时间搜索
        if (!empty($params['start_time'])) {
            // 获取当前时间戳的凌晨时间
            $start = strtotime(date("Y-m-d", $params['start_time']));
            $params['reserve_start_time'] = $start;
            $params['reserve_end_time']   = $start + 86399;
            unset($params['start_time']);
        }

        $params['status'] = $params['status'] ?? -1; //所有邀约

        $params["order_sort"] = 1;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/lists-by-coach')->post($params);
        } catch (Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
        return success($result);
    }

    /**
     * 预约单编辑
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function orderEdit(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/edit')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 预约单取消
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function orderCancel(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/cancel')->post($params);
        } catch (Exception $e) {
            error(900003);
        }
        return success($result);
    }

    /**
     * 预约单分配
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function orderAssign(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/assign')->post($params);
        } catch (Exception $e) {
            error(900003);
        }
        return success($result);
    }


    /**
     * 分享预约单
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function orderShare(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/share')->post($params);
        } catch (Exception $e) {
            error(900003);
        }
        return success($result);
    }


    /**
     * 获取用户的基本信息
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function userSimpleInfo(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/user-simple-info')->post($params);
        } catch (Exception $e) {
            error(900003);
        }
        return success($result);
    }

    /**
     * 获取导购门店信息
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function memberInfo(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $info                    = [];
        try {
            $info = $this->work->setApi('service/v1/store/department-detail')->post($params);
            $info = current($info);

            if (empty($info)) {
                $result = $this->microScrm->setApi('micro-scrm/v1/member/get-detail')->post($params);
                $info   = [
                    'third_store_code'     => $result['department_name'],
                    'wechat_department_id' => $result['department_third_id'],
                    'store_name'           => $result['department_name'],
                    'address'              => $result['address'],
                ];
            }

        } catch (Exception $e) {
            error(900003);
        }
        $info['enterprise_id'] = $params['enterprise_id'];

        return success($info);
    }


    /**
     * 服务列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function serviceList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        $params['status']        = 1; // 有效的
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/service/member-list')->post($params);
        } catch (Exception $e) {
            error(900003);
        }
        return success($result);
    }

    /**
     * 导购侧预约单详情
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function orderDetail(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/order-info')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 完成预约
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function orderFinish(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/finish')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 获取与外部联系人有关的导购列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     */
    public function memberLists(Request $request)
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/member-lists')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 导购确认订单是否是到店购买
     * @param Request $request
     * @return JsonResponse
     */
    public function changeStatus(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/chang-status')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 导购确认用户到店
     * @param Request $request
     * @return JsonResponse
     */
    public function confirmStore(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/confirm-store')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 导购分享邀约用户评论
     * @param Request $request
     * @return JsonResponse
     */
    public function commentShare(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/comment-share')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 验证导购是否可以给用户发送邀请评价
     * @param Request $request
     * @return JsonResponse
     */
    public function verifyShare(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/verify-share')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * SW&&ks验证同一时间段的邀约限制
     * @param Request $request
     * @return JsonResponse
     */
    public function checkReserveTime(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/check-reserve-time')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 客户详情的邀约单列表展示客户的所有邀约单数据
     * @param Request $request
     * @return JsonResponse
     */
    public function customerReserveLists(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        if (empty($params['external_userid'])){
            \error(10002,"无效的外部联系人");
        }
        $result                  = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/customer-reserve-lists')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }
}