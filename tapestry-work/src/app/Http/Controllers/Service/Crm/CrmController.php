<?php

namespace App\Http\Controllers\Service\Crm;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Sdk\Active\CRM\CoachCrm;
use App\Sdk\Active\CRM\SwCoupon\SwCouponClient;
use App\Services\Crm\CrmService;
use Exception;
use Illuminate\Http\JsonResponse;

class CrmController extends BaseController
{

    private $enterpriseId;

    public function __construct()
    {
        parent::__construct();
        $this->enterpriseId = $this->request->input('enterprise_id');

    }

    /**
     * 测试获取token
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getToken(): JsonResponse
    {
        $result['token'] = (new CrmService($this->enterpriseId))->getToken();

        return successJson($result);
    }

    /**
     * 通用获取第三方会员信息
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2023/3/17
     */
    public function getCustomerInfo(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->getCustomerInfo($input);

        return successJson($result);
    }

    /**
     * 获取第三方会员信息接口
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getMemberInfo(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->getMemberInfo($input);

        return successJson($result);
    }

    /**
     * 查询会员积分记录
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerPointHistory(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->getCustomerPointHistory($input);

        return successJson($result);
    }

    /**
     * 查询会员优惠卷记录
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerCoupon(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->getCustomerCoupon($input);

        return successJson($result);

    }

    /**
     * 查询Sw会员优惠卷记录
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getSwCustomerCoupon(): JsonResponse
    {
        // 接收参数
        $input = $this->request->all();
        if (empty($input['union_id'])) {
            throw new Error(100000, '参数错误');
        }
        $unionId = $input['union_id'];
        $page = $input['page'] ?? 1;
        $size = $input['size'] ?? 999;
        $result = (new SwCouponClient())->getUserCouponList($unionId, $page, $size);
        return successJson($result);

    }

    /**
     * 查询会员购买记录\订单
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerPurchase(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->getCustomerPurchase($input);

        return successJson($result);

    }

    /**
     * 查询会员标签记录
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerTag(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->getCustomerTag($input);

        return successJson($result);

    }

    /**
     * 用户是否是会员
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/29
     */
    public function checkMember(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->checkMember($input);

        return successJson($result);

    }

    /**
     * coach发放优惠券
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionWelfareGive(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->actionWelfareGive($input);

        return successJson($result);

    }

    /**
     * coach是否发放过优惠券
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionWelfareHaved(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->actionWelfareHaved($input);

        return successJson($result);

    }

    /**
     * ks是否发放过优惠券
     *
     * @return JsonResponse
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionUserJoin(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->actionUserJoin($input);

        return successJson($result);

    }

    /**
     * coach导购添加顾客
     *
     * @return JsonResponse
     * @Time: 2022/11/18
     */
    public function saAssignment(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->saAssignment($input);

        return successJson($result);
    }

    /**
     * coach导购删除顾客
     *
     * @return JsonResponse
     * @Time: 2022/11/18
     */
    public function saUnAssignment(): JsonResponse
    {
        // 接收参数
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->saUnAssignment($input);

        return successJson($result);
    }

    /**
     * coach聚合查询会员列表筛选条件
     *
     * @return JsonResponse
     */
    public function customerOption(): JsonResponse
    {
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->customerOption($input);

        return successJson($result);
    }

    /**
     * coach根据门店或门店导购批量查询会员列表
     *
     * @return JsonResponse
     */
    public function customerList(): JsonResponse
    {
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->customerListWeclienteling($input);

        return successJson($result);
    }

    /**
     * coach查询会员详情
     *
     * @return JsonResponse
     */
    public function customerDetail(): JsonResponse
    {
        $input  = $this->request->all();
        $result = (new CrmService($this->enterpriseId))->customerDetail($input);

        return successJson($result);
    }


    /**
     * 获取CoachCrm第三方会员信息
     *
     * @return array
     * @throws Exception
     * @author: wangle@vchangyi.com
     * @Time: 2023/3/17
     */
    public function getCoachCustomerInfo(): array
    {
        // 接收参数
        $input  = $this->request->all();
        if (empty($input['union_id'])) {
            error(110010);
        }
        $coachCrm = new CoachCrm(['enterprise_id' => $this->enterpriseId]);
        $result =  $coachCrm->getCustomerInfoBySearchV2($input);
        $result['vipStatus'] = [];
        if (!empty($result['membership'])) {
            $result['vipStatus'] = $coachCrm->getCustomerVip($result['membership']);
        }
        return $result;
    }

}