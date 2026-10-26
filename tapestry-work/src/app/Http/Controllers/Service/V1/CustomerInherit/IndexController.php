<?php

namespace App\Http\Controllers\Service\V1\CustomerInherit;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerInheritService;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Redis;
use Illuminate\Validation\ValidationException;

class IndexController extends BaseController
{

    protected $customer;

    /**
     * @throws Exception
     */
    public function __construct(CustomerInheritService $customer)
    {
        parent::__construct();
        $this->customer = $customer;
    }

    /**
     * @return void
     * 任务列表
     */
    public function inheritList(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required',
            'status' => 'sometimes|int|min:0|max:3',
            'created_at' => 'sometimes|array',
            'page' => 'sometimes|int',
            'size' => 'sometimes|int',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required' => '手机号不能为空',
            'status.int' => '状态类型不正确',
            'created_at.array' => '时间类型不正确必须是一个数组',
            'status.min' => '状态不正确',
            'status.max' => '状态不正确',
        ];
        $this->validate($this->request, $rules, $message);

        //查询SM信息
        $memberIds = $this->customer->getMemberIds($params);
        if (!$memberIds) {
            error("导购身份不正确");
        }
        //查列表信息
        $memberIds = array_merge($memberIds, [$params['member_id']]);
        $res = $this->customer->inheritList($params, $memberIds);

        return successJson($res);
    }


    /**
     * @return void
     * 任务详情
     * @throws ValidationException
     */
    public function inheritInfo(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required',
            'inherit_id' => 'required',
            'page' => 'sometimes|int',
            'size' => 'sometimes|int',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required' => '手机号不能为空',
            'inherit_id.required' => '任务id不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        $res = $this->customer->inheritInfo($params);

        return successJson($res);
    }

    /**
     * @return void
     * 在职转移
     */
    public function inheritTransfer(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required',
            'transfer_success_msg' => 'required',
            'takeover_member_id' => 'required',//接替导购成员的id
            'relation_id' => 'array',//work_external_user_relation 里的id

        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required' => '手机号不能为空',
            'inherit_id.required' => '任务id不能为空',
        ];
        $this->validate($this->request, $rules, $message);
        //避免连续重复提交
        $cacheKey = md5(json_encode($params));
        $paramsKey = redisKey('CUSTOMER_INHERIT_PARAMS', ['hash' => $cacheKey]);
        $cache = getOrSetData($paramsKey);
        if ($cache) {
            return successJson();
        }
        getOrSetData($paramsKey, true, 60);

        //避免上次异步还未处理结束又一次提交
        $timeKey = redisKey('CUSTOMER_INHERIT_WAIT_TIME');
        $timeCache = getOrSetData($timeKey);
        if ($timeCache) {
            $ttl = Redis::connection()->ttl($timeKey);
            error("您上次转移的顾客较多，请等待" . $ttl . '秒后再试');
        }
        if (isset($params['is_all_checked']) && $params['is_all_checked'] == 1) { // 全选在职继承
            getOrSetData($timeKey, true, 600);
        } else {
            getOrSetData($timeKey, count($params['relation_id']), count($params['relation_id']));
        }
        $res = $this->customer->inheritTransfer($params);
        if ($res) {
            return successJson(["id" => $res]);
        } else {
            error("提交失败");
        }
    }


    /**
     * @return JsonResponse
     * @throws ValidationException
     * 导购列表
     */
    public function memberList(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required',
            'search' => 'sometimes',
            'dp_ids' => 'array',
            'page' => 'sometimes|int',
            'size' => 'sometimes|int',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required' => '手机号不能为空',
            'dp_ids.array' => '门店信息格式不正确',
        ];
        $this->validate($this->request, $rules, $message);
        $res = $this->customer->memberList($params);

        return successJson($res);
    }

    /**
     * @return void
     * 查询外部联系人转移次数以及转移中的用户
     */
    public function customerInheritRecord(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'external_id' => 'array',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'external_id.array' => '外部联系人格式不正确',
        ];
        $this->validate($this->request, $rules, $message);
        //获取剩余转移次数
        $externalToTotal = $this->customer->customerInheritRecord($params);
        //
        $waitExternalArr = $this->customer->customerInheritWait($params);

        return successJson(['external_to_total' => $externalToTotal, 'wait_external_arr' => $waitExternalArr]);
    }

    /**
     * @return JsonResponse
     * @throws ValidationException
     * 外部联系人
     */
    public function memberRecord(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'external_id' => 'array',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'external_id.array' => '外部联系人格式不正确',
        ];
        $this->validate($this->request, $rules, $message);

        $res = $this->customer->memberRecord($params);

        return successJson($res);
    }

    /**
     * 查询分配顾客成功的记录
     *
     * @return JsonResponse
     * @throws ValidationException
     * @Time: 2022/11/17
     */
    public function recordForExternalIds(): JsonResponse
    {
        $params = $this->request->all();
        $rules = [
            'enterprise_id' => 'required',
            'external_ids' => 'array',
        ];
        $message = [
            'enterprise_ids.required' => '企业ID不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        $res = $this->customer->customerInheritRecordForExternalIds($params);

        return successJson($res);
    }
}