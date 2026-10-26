<?php
/**
 * Created by PhpStorm.
 * User: guokai
 * Date: 2019/04/18
 * Time: 17:55
 */
namespace App\Http\Controllers\V1\Page\ReduceMoney;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PageReduce;
use Illuminate\Http\Request;

class AddController extends BaseController
{
    /**
     * 添加立减金
     *
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $data = $request->input();
        $rules = [
            'name' => 'required',
            'condition_value' => 'required',
            'discount_amount' => 'required',
            'start_time' => 'required',
            'end_time' => 'required',
            'receive_count' => 'required',
        ];
        $messages = [
            'name.required' => '立减金名称不能为空',
            'condition_value.required' => '立减金条件值不能为空',
            'discount_amount.required' => '立减金额不能为空',
            'start_time.required' => '有效期开始时间不能为空',
            'end_time.required' => '有效期结束时间不能为空',
            'receive_count.required' => '领取次数不能为空',

        ];
        $this->validate($request, $rules, $messages);

        $data['enterprise_id'] = $this->enterpriseId;

        try{
            $this->writeOperatorLog('添加立减金');
        }catch (\Exception $e){
        }

        $result = (new PageReduce())->add($data);

        return success($result);
    }
}
