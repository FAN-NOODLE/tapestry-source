<?php

namespace App\Http\Controllers\V1\RedPacket;

use CrsSdk\Work\RedPacket;
use Illuminate\Validation\Rule;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class RedPacketController extends BaseController
{
    // 红包上架中
    const IS_ON_SALE = 2;

    // 红包已下架
    const IS_OFF_SALE = 3;

    /**
     * 更新红包
     * @return JsonResponse
     * @throws \Exception
     */
    public function update()
    {
        $params = $this->request->only([
            'id',
            'cash_coupon_name',
            'amount',
            'quantity',
            'send_name',
            'wishing',
            'remark',
            'members',
        ]);

        $rules = [
            'id'               => 'required|integer|min:1',
            'cash_coupon_name' => 'required',
            'amount'           => 'required|integer|min:1',
            'quantity'         => 'required|integer|min:1',
            'send_name'        => 'required',
            'wishing'          => 'required',
            'remark'           => 'required',
            'members'          => 'required|array',
            'members.*.value'  => 'required',
            'members.*.label'  => 'required',
        ];

        $message = [
            'id.required'               => '红包ID不能为空',
            'id.integer'                => '红包ID类型错误',
            'id.min'                    => '红包ID不能为零',
            'cash_coupon_name.required' => '活动名称不能为空',
            'amount.required'           => '红包金额不能为空',
            'amount.integer'            => '红包金额类型错误',
            'amount.min'                => '红包金额不能为零',
            'quantity.required'         => '红包数量不能为空',
            'quantity.integer'          => '红包数量类型错误',
            'quantity.min'              => '红包数量不能为零',
            'send_name.required'        => '商户名称不能为空',
            'wishing.required'          => '红包祝福语不能为空',
            'remark.required'           => '红包备注不能为空',
            'members.required'          => '发送人列表不能为空',
            'members.*.value.required'  => '发送人uid不能为空',
            'members.*.label.required'  => '发送人名称不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        (new RedPacket())->update($conditions);

        return success();
    }


    /**
     * 创建红包
     * @return JsonResponse
     * @throws \Exception
     */
    public function create()
    {
        $params = $this->request->only([
            'cash_coupon_name',
            'amount',
            'quantity',
            'send_name',
            'wishing',
            'remark',
            'members',
        ]);

        $rules = [
            'cash_coupon_name' => 'required',
            'amount'           => 'required|integer|min:1',
            'quantity'         => 'required|integer|min:1',
            'send_name'        => 'required',
            'wishing'          => 'required',
            'remark'           => 'required',
            'members'          => 'required|array',
            'members.*.value'  => 'required',
            'members.*.label'  => 'required',
        ];

        $message = [
            'cash_coupon_name.required' => '活动名称不能为空',
            'amount.required'           => '红包金额不能为空',
            'amount.integer'            => '红包金额类型错误',
            'amount.min'                => '红包金额不能为零',
            'quantity.required'         => '红包数量不能为空',
            'quantity.integer'          => '红包数量类型错误',
            'quantity.min'              => '红包数量不能为零',
            'send_name.required'        => '商户名称不能为空',
            'wishing.required'          => '红包祝福语不能为空',
            'remark.required'           => '红包备注不能为空',
            'members.required'          => '发送人列表不能为空',
            'members.*.value.required'  => '发送人uid不能为空',
            'members.*.label.required'  => '发送人名称不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        (new RedPacket())->create($conditions);

        return success();
    }

    /**
     * 红包详情
     * @return JsonResponse
     * @throws \Exception
     */
    public function detail()
    {
        $params = $this->request->only(['id']);

        $rules = [
            'id' => 'required',
        ];

        $message = [
            'id.required' => '红包ID不能为空',
            'id.integer'  => '红包ID类型错误',
            'id.min'      => '红包ID不能为零',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        $result = (new RedPacket())->detail($conditions);

        return success($result['data']);
    }

    /**
     * 红包发送人列表
     * @return JsonResponse
     * @throws \Exception
     */
    public function members()
    {
        $params = $this->request->only(['id']);

        $rules = [
            'id' => 'required',
        ];

        $message = [
            'id.required' => '红包ID不能为空',
            'id.integer'  => '红包ID类型错误',
            'id.min'      => '红包ID不能为零',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        $result = (new RedPacket())->members($conditions);

        return success([
            'list' => $result['data']
        ]);
    }

    /**
     * 红包上下架
     * @return JsonResponse
     * @throws \Exception
     */
    public function sale()
    {
        $params = $this->request->only(['id', 'status']);

        $rules = [
            'id'     => 'required',
            'status' => [
                'required',
                Rule::in([self::IS_ON_SALE, self::IS_OFF_SALE])
            ],
        ];

        $message = [
            'id.required'     => '红包ID不能为空',
            'id.integer'      => '红包ID类型错误',
            'id.min'          => '红包ID不能为零',
            'status.required' => '上架状态不能为空',
            'status.in'       => '上架状态不在范围',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        (new RedPacket())->sale($conditions);

        return success();
    }

    /**
     * 开关切换
     * @return JsonResponse
     * @throws \Exception
     */
    public function delete()
    {
        $params = $this->request->only(['id']);

        $rules = [
            'id' => 'required',
        ];

        $message = [
            'id.required' => '红包ID不能为空',
            'id.integer'  => '红包ID类型错误',
            'id.min'      => '红包ID不能为零',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        (new RedPacket())->delete($conditions);

        return success();
    }

    /**
     * 红包列表
     * @return JsonResponse
     * @throws \Exception
     */
    public function list()
    {
        $params                  = $this->request->only(['is_deleted', 'start_time', 'end_time', 'cash_coupon_name', 'size', 'page']);
        $params['is_deleted']    = $params['is_deleted'] ?? 1;
        $params['enterprise_id'] = $this->enterpriseId;

        $list = (new RedPacket())->list($params);

        return success([
            'list'               => $list['data'],
            'page'               => $list['meta']['page'],
            'size'               => $list['meta']['size'],
            'total'              => $list['meta']['total'],
            'cash_coupon_status' => $list['cash_coupon_status'],
        ]);
    }

    /**
     * 红包对账列表
     * @return JsonResponse
     * @throws \Exception
     */
    public function reconciliationList()
    {
        $params                  = $this->request->only(['start_time', 'end_time', 'mem_uid', 'id', 'status', 'size', 'page']);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['domain']        = $this->qyDomain;
        $params['identifier']    = $this->appIdentifier;

        $list = (new RedPacket())->reconciliationList($params);

        return success([
            'list'   => $list['data'],
            'page'   => $list['meta']['page'],
            'size'   => $list['meta']['size'],
            'total'  => $list['meta']['total'],
            'detail' => $list['detail'],
        ]);
    }
}