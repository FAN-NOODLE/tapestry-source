<?php

namespace App\Http\Controllers\V1\Contact;

use CrsSdk\Work\Member;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class ContactController extends BaseController
{
    const TYPE_MEMBER = 1; // 员工类型

    const TYPE_DEPART = 2; // 部门类型

    const QR_SINGLE = 1; // 单人二维码

    const QR_MANY = 2; // 多人二维码

    // 对象类型数组
    const TYPE_ARRAY = [
        self::TYPE_DEPART,
        self::TYPE_MEMBER
    ];

    // 二维码类型数组
    const QR_TYPE_ARRAY = [
        self::QR_SINGLE,
        self::QR_MANY
    ];

    /**
     * 选择员工tree
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function memberTree(Request $request)
    {
        $params = $request->only(['id', 'type']);

        $rules = [
            'id' => 'sometimes|required'
        ];

        $message = [
            'id.required' => '联系人ID不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => [
                'id'   => $params['id'] ?? 0,
                'type' => $params['type'] ?? 1 // 默认包含员工
            ]
        ];

        $result = (new Member())->memberTree($conditions);

        return success([
            'list' => $result['data']
        ]);
    }

    /**
     * 新建联系方式
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function createQrCode(Request $request)
    {
        $params = $request->only(['channel', 'type', 'is_verify', 'members', 'channel_id', 'name']);
        $rules  = [
            'channel'         => 'required',
            'type'            => [
                'required',
                Rule::in(self::QR_TYPE_ARRAY),
            ],
            'is_verify'       => 'required',
            'members'         => 'required|array',
            'members.*.value' => 'required',
            'members.*.label' => 'required',
            'members.*.type'  => 'required',
        ];

        $message = [
            'channel.required'         => '二维码渠道名称不能为空',
            'type.required'            => '二维码类型不能为空',
            'type.in'                  => '二维码类型不在范围',
            'is_verify.required'       => '二维码添加设置不能为空',
            'members.required'         => '二维码[members]不能为空',
            'members.array'            => '二维码[members]类型错误',
            'members.*.value.required' => '二维码标识数值[value]不能为空',
            'members.*.label.required' => '二维码标识名称[label]不能为空',
            'members.*.type.required'  => '二维码标识类型[type]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        $result = (new Member())->createQrCode($conditions);

        return success($result['data']);
    }

    /**
     * 编辑修改联系方式
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function updateQrCode(Request $request)
    {
        $params = $request->only(['id', 'channel_id', 'name', 'channel', 'type', 'is_verify', 'members']);
        $rules  = [
            'id'              => 'required|integer|min:1',
            'channel'         => 'required',
            'type'            => [
                'required',
                Rule::in(self::QR_TYPE_ARRAY)
            ],
            'is_verify'       => 'required',
            'members'         => 'required|array',
            'members.*.value' => 'required',
            'members.*.label' => 'required',
            'members.*.type'  => 'required',
        ];

        $message = [
            'id.required'              => '二维码[id]不能为空',
            'id.integer'               => '二维码[id]类型错误',
            'id.min'                   => '二维码[id]不能为零',
            'channel.required'         => '二维码渠道名称不能为空',
            'type.required'            => '二维码类型不能为空',
            'type.in'                  => '二维码类型不在范围',
            'is_verify.required'       => '二维码添加设置不能为空',
            'members.required'         => '二维码[members]不能为空',
            'members.array'            => '二维码[members]类型错误',
            'members.*.value.required' => '二维码标识数值[value]不能为空',
            'members.*.label.required' => '二维码标识名称[label]不能为空',
            'members.*.type.required'  => '二维码标识类型[type]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => $params
        ];

        $result = (new Member())->updateQrCode($conditions);

        return success($result['data']);
    }

    /**
     * 查看联系我
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function selectQrCode(Request $request)
    {
        $params = $request->only(['id']);

        $rules = [
            'id' => 'required|integer|min:1'
        ];

        $message = [
            'id.required' => '二维码[id]不能为空',
            'id.integer'  => '二维码[id]类型错误',
            'id.min'      => '二维码[id]不能为零',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => [
                'id' => $params['id'],
            ]
        ];

        $result = (new Member())->selectQrCode($conditions);

        return success($result['data']);
    }

    /**
     * 删除联系我
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function deleteQrCode(Request $request)
    {
        $params = $request->only(['id']);

        $rules = [
            'id' => 'required|integer|min:1'
        ];

        $message = [
            'id.required' => '二维码[id]不能为空',
            'id.integer'  => '二维码[id]类型错误',
            'id.min'      => '二维码[id]不能为零',
        ];

        Validator::make($params, $rules, $message)->validate();

        $conditions = [
            'enterprise_id' => $this->enterpriseId,
            'domain'        => $this->qyDomain,
            'identifier'    => $this->appIdentifier,
            'extend'        => [
                'id' => $params['id'],
            ]
        ];

        (new Member())->deleteQrCode($conditions);

        return success();
    }
}