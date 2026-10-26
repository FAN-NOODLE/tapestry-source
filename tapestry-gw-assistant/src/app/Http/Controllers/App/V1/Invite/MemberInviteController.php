<?php

namespace App\Http\Controllers\App\V1\Invite;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Validation\ValidationException;
use Illuminate\Http\{JsonResponse, Request};

/**
 * 导购端邀请函
 */
class MemberInviteController extends BaseController
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
     * 邀请函配置页列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/11/9
     */
    public function configList(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/invite/template/config-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 邀请函模板配置详情
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/11/9
     */
    public function configDetail(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/invite/template/config-detail')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 发送邀约单记录
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @Time: 2022/11/9
     */
    public function createRecord(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/invite/create-record')->post($params);
        } catch (Exception $e) {

            error(10002, $e->getMessage());
        }

        return success($result);

    }

    /**
     * 创建预约单
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     */
    public function createOrder(Request $request): JsonResponse
    {
        // 参数校验
        $rules   = [
            'style_type'     => 'required|integer',
            'background_url' => 'required|string',
        ];
        $message = [
            'style_type.integer'      => 'style_type参数类型有误',
            'background_url.required' => 'background_url参数不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/reserve/order/create')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * 导购查看邀请函信息
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/9
     */
    public function getRecordInfo(Request $request): JsonResponse
    {
        // 参数校验
        $rules   = [
            'invite_record_id' => 'required|integer',
        ];
        $message = [
            'invite_record_id.required' => 'invite_record_id参数不能为空',
        ];
        $this->validate($request, $rules, $message);
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;

        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/invite/get-record-info')->post($params);
        } catch (Exception $e) {
            if ($e->getMessage() == "无权访问") {
                error(10001, "权限不足");
            }
            error(10002, $e->getMessage());
        }

        return success($result);

    }

}