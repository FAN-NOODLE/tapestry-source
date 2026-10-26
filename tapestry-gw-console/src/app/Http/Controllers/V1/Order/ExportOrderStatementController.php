<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2020/9/16
 * Time: 20:40
 */

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Illuminate\Support\Facades\Validator;

class ExportOrderStatementController extends BaseController
{
    /**
     * 对账导出
     * {domain}/gw-console/v1/order/export-order-statement
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        // 接收参数
        $input = $this->request->only(['start_time', 'end_time']);
        // 参数验证
        $this->validateParams($input);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['operator_id'] = $this->enterpriseAdmin['id'];
        $input['operator_name'] = $this->enterpriseAdmin['username'];

        (new Order())->exportOrderStatement($input);

        return success();
    }

    /**
     * 参数验证
     *
     * @param array $input
     */
    protected function validateParams(array $input)
    {
        // 验证规则
        $rules = [
            'start_time' => 'required | integer',
            'end_time'   => 'required | integer'
        ];
        // 提示语
        $messages = [
            'start_time.required' => '[参数]start_time必须',
            'start_time.integer'  => '[参数]start_time须为integer类型',
            'end_time.required'   => '[参数]end_time必须',
            'end_time.integer'    => '[参数]end_time须为integer类型'
        ];

        // 参数验证
        Validator::make($input, $rules, $messages)->validate();
    }
}
