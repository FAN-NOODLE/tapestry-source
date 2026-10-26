<?php

namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;
use App\Common\Constant;

class EnterpriseRequest
{
    /**
     * 企业添加验证
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     * @param array $data 验证的数据
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function common($data = [])
    {
        return Validator::make($data, [
            'name'         => 'required',
            'contact'      => 'required',
            'mobile'       => 'required',
            'address'      => 'required',
            'crs_status'   => 'required | in:' . Constant::ENTERPRISE_INSTALL_STATUS . ','
                . Constant::ENTERPRISE_UNINSTALLED_STATUS,
            'store_status' => 'required | in:' . Constant::ENTERPRISE_INSTALL_STATUS . ','
                . Constant::ENTERPRISE_UNINSTALLED_STATUS
        ], [
            'name.required'         => '企业名称不能为空',
            'contact.required'      => '联系人不能为空',
            'mobile.required'       => '联系人电话不能为空',
            'address.required'      => '详细地址不能为空',
            'crs_status.required'   => 'crs服务传值不能为空',
            'crs_status.in'         => 'crs服务传值类型错误',
            'store_status.required' => '商城服务传值不能为空',
            'store_status.in'       => '商城服务传值类型错误'
        ]);
    }

    /**
     * 企业信息获取验证（shop用于获取企业应用开通查询，谨慎修改）
     *
     * @author zhaoqi zhaoqi@vchangyi.com
     * @date   2019/02/25 15:35
     *
     * @param array $data 验证的数据
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function enterpriseInfo($data = [])
    {
        return Validator::make($data, [
            'enterprise_id' => 'required|integer|min:1'
        ], [
            'enterprise_id.required' => 600001,
            'enterprise_id.integer'  => 201104,
            'enterprise_id.min'      => 201104
        ]);
    }
}
