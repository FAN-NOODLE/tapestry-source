<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;

class RemarkController extends BaseController
{
    /**
     * 备注新增接口
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function add()
    {
        $input = $this->request->all();

        $this->validate($this->request, [
            'order_no' => 'required',
            'type'     => 'required|integer',
            'content'  => 'required'
        ], [
            'order_no.required' => 'order_no必填',
            'type.required'     => 'type必填',
            'type.integer'      => 'type需为integer',
            'content.required'  => 'content必填'
        ]);

        $data = (new Order())->remarkAdd([
            'enterprise_id' => $this->enterpriseId,
            'order_no'      => $input['order_no'],
            'type'          => $input['type'],
            'admin_id'      => $this->enterpriseAdmin['id'],
            'admin_name'    => $this->enterpriseAdmin['username'],
            'content'       => $input['content']
        ])['data'];

        $this->writeOperatorLog('管理员备注订单：'.$input['order_no']);

        return success(['data' => $data]);
    }

    /**
     * 更新用户订单备注为已读状态（已废弃，改为在运营添加备注接口更新）
     * {domain}/app/v1/order/remark-to-read
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function remarkToRead()
    {
        $this->validate($this->request, [
            'order_no'      => 'required'
        ], [
            'order_no.required'      => 'order_no必填'
        ]);

        $input = $this->request->only(['enterprise_id', 'order_no']);

        $input['enterprise_id'] = $this->enterpriseId;
        $response = (new Order())->remarkToRead($input);

        return success($response);
    }
}
