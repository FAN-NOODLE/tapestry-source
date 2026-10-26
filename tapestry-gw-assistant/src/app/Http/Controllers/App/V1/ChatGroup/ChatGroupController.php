<?php

namespace App\Http\Controllers\App\V1\ChatGroup;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class ChatGroupController extends BaseController
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
     *社群成员明细
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function membersInfo(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $memberInfo              = $this->microScrm->setApi('micro-scrm/v1/members/list')->post([
            'enterprise_id' => $this->enterpriseId,
            'mobile'        => $this->mobile
        ]);
        $params['member_id']     = current($memberInfo['member_list'])['id'];
        $params['user_id']       = current($memberInfo['member_list'])['userid'];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/chat-group/members_info')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 社群详情
     *
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function indexList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/chat-group/index_list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     *
     * 数据更新时间
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function dataUpdateTime(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/chat-group/data-update-time')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }

        return success([$result]);
    }


}