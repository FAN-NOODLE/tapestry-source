<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/3
 * Time: 15:10
 */

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;

class UserRequest
{

    /**
     * 新增员工参数验证
     *
     * @param $params
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function add($params)
    {
        // 参数验证规则
        $rules = [
            'mem_username' => 'required',
            // 'mem_face'     => 'required',
            'dp_id'        => 'required',
            'is_store'     => 'required|in:0,1|integer',
            'mem_gender'   => 'required|in:0,1,2|integer',
            'mem_mobile'   => 'required',
            'mem_role'     => 'nullable|required_if:is_store,1'
        ];

        // 消息提示
        $message = [
            'mem_username.required' => '320152',
            // 'mem_face.required'     => '320153',
            'dp_id.required'        => '320154',
            'is_store.required'     => '320155',
            'is_store.in'           => '320156',
            'mem_gender.required'   => '320157',
            'mem_gender.in'         => '320158',
            'mem_mobile.required'   => '320159',
            'mem_role.required_if'  => '320160',
        ];

        return Validator::make($params, $rules, $message);
    }

    /**
     * 编辑员工参数验证
     *
     * @param $params
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function update($params)
    {
        // 参数验证规则
        $rules = [
            'mem_uid'      => 'required',
            'mem_username' => 'required',
            // 'mem_face'     => 'required',
            'dp_id'        => 'required',
            'is_store'     => 'in:0,1|integer',
            'mem_gender'   => 'required|in:0,1,2|integer',
            'mem_mobile'   => 'required',
            'mem_role'     => 'nullable|required_if:is_store,1'
        ];

        // 消息提示
        $message = [
            'mem_uid.required'      => '320151',
            'mem_username.required' => '320152',
            // 'mem_face.required'     => '320153',
            'dp_id.required'        => '320154',
//            'is_store.required'     => '320155',
            'is_store.in'           => '320156',
            'mem_gender.required'   => '320157',
            'mem_gender.in'         => '320158',
            'mem_mobile.required'   => '320159',
            'mem_role.required_if'  => '320160',
        ];

        return Validator::make($params, $rules, $message);
    }
}
