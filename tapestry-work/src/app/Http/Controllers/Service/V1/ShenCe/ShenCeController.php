<?php

namespace App\Http\Controllers\Service\V1\ShenCe;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Sdk\ShenCe\Client;
use App\Services\ShenCe\ShenCeService;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Validator;

class ShenCeController extends BaseController
{
    /**
     * @var ShenCeService
     */
    private $service;

    public function __construct()
    {
        parent::__construct();

        $enterpriseId  = $this->request->input('enterprise_id');
        $this->service = new ShenCeService($enterpriseId, Client::SyncConsumer);
    }

    public function itemSet(): JsonResponse
    {
        $validator = Validator::make($this->request->all(), [
            'item_type'  => 'required|string',
            'item_id'    => 'required|string',
            'properties' => 'required|array',
        ], [
            'item_type.required'  => 'item_type 字段为必填项',
            'item_type.string'    => 'item_type 字段必须为字符串类型',
            'item_id.required'    => 'item_id 字段为必填项',
            'item_id.string'      => 'item_id 字段必须为字符串类型',
            'properties.required' => 'properties 字段为必填项',
            'properties.array'    => 'properties 字段必须为对象类型',
        ]);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }

        $params     = $this->request->all();
        $item_type  = $params['item_type'];
        $item_id    = $params['item_id'];
        $properties = $params['properties'];

        $result = $this->service->getClient()->item_set($item_type, $item_id, $properties);

        return successJson(compact('result'));
    }

    public function profileSet(): JsonResponse
    {
        $validator = Validator::make($this->request->all(), [
            'distinct_id' => 'required|string',
            'profiles'    => 'required|array',
        ], [
            'distinct_id.required' => 'distinct_id 字段为必填项',
            'distinct_id.string'   => 'distinct_id 字段必须为字符串类型',
            'properties.required'  => 'properties 字段为必填项',
            'properties.array'     => 'properties 字段必须为对象类型',
        ]);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }

        $params      = $this->request->all();
        $distinct_id = $params['distinct_id'];
        $is_login_id = isset($params['is_login_id']) && $params['is_login_id'];
        $profiles    = $params['profiles'];

        $result = $this->service->getClient()->profile_set($distinct_id, $is_login_id, $profiles);

        return successJson(compact('result'));
    }

    public function track(): JsonResponse
    {
        $validator = Validator::make($this->request->all(), [
            'distinct_id' => 'required|string',
            'event_name'  => 'required|string',
            'profiles'    => 'required|array',
        ], [
            'distinct_id.required' => 'distinct_id 字段为必填项',
            'distinct_id.string'   => 'distinct_id 字段必须为字符串类型',
            'event_name.required'  => 'event_name 字段为必填项',
            'event_name.string'    => 'event_name 字段必须为字符串类型',
            'properties.required'  => 'properties 字段为必填项',
            'properties.array'     => 'properties 字段必须为对象类型',
        ]);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }

        $params      = $this->request->all();
        $distinct_id = $params['distinct_id'];
        $is_login_id = isset($params['is_login_id']) && $params['is_login_id'];
        $profiles    = $params['profiles'];
        $eventName   = $params['event_name'];

        $result = $this->service->getClient()->track($distinct_id, $is_login_id, $eventName, $profiles);

        return successJson(compact('result'));
    }
}