<?php


namespace App\Http\Requests;


class StoreRequest
{
    /**
     * @param $params
     *
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function add($params)
    {
        // 参数验证规则
        $rules = [
            'enterprise_id' => 'required',
            'name'          => 'required|between:1,50',
            'short_name'    => 'required|between:1,50',
            'type'          => 'required',
            'area_path'     => 'required',
            'province'      => 'required',
            'city'          => 'required',
            'district'      => 'required',
            'address'       => 'required',
            'code'          => 'required',
            'lon_no'        => 'required',
            'lat_no'        => 'required',
            'bs_hours'      => 'required',
            'work_id'       => 'required',
            'work_phone'    => 'required',
            'post_tpl'      => 'required',
        ];

        // 提示语
        $messages = [
            'enterprise_id.*' => '企业Id不为空',
            'name.*'          => '门店名称不能为空',
            'short_name.*'    => '门店简称不能为空',
            'type.*'          => '门店类型不能为空',
            'area_path.*'     => '门店所属区域必选',
            'province.*'      => '门店位置必选',
            'city.*'          => '门店位置必选',
            'district.*'      => '门店位置必选',
            'address.*'       => '门店位置必选',
            'code.*'          => '省市区编码有误',
            'lon_no.*'        => '经度必填',
            'lat_no.*'        => '纬度必填',
            'bs_hours.*'      => '门店营业时间',
            'work_id.*'       => '店长id不能为空',
            'work_phone.*'    => '店长手机号不能为空',
            'post_tpl.*'      => '邮费模板不能为空',
        ];

        validate($params, $rules, $messages);
    }

    /**
     * @param $params
     *
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function edit($params)
    {
        // 参数验证规则
        $rules = [
            'id'            => 'required',
            'enterprise_id' => 'required',
            'name'          => 'required|between:1,50',
            'short_name'    => 'required|between:1,50',
            'type'          => 'required',
            'area_path'     => 'required',
            'province'      => 'required',
            'city'          => 'required',
            'district'      => 'required',
            'address'       => 'required',
            'code'          => 'required',
            'lon_no'        => 'required',
            'lat_no'        => 'required',
            'bs_hours'      => 'required',
            'work_id'       => 'required',
            'work_phone'    => 'required',
            'post_tpl'      => 'required',
        ];

        // 提示语
        $messages = [
            'id.*'            => '门店Id不能为空',
            'enterprise_id.*' => '企业Id不为空',
            'name.*'          => '门店名称不能为空',
            'short_name.*'    => '门店简称不能为空',
            'type.*'          => '门店类型不能为空',
            'area_path.*'     => '门店所属区域必选',
            'province.*'      => '门店位置必选',
            'city.*'          => '门店位置必选',
            'district.*'      => '门店位置必选',
            'address.*'       => '门店位置必选',
            'code.*'          => '省市区编码有误',
            'lon_no.*'        => '经度必填',
            'lat_no.*'        => '纬度必填',
            'bs_hours.*'      => '门店营业时间',
            'work_id.*'       => '店长id不能为空',
            'work_phone.*'    => '店长手机号不能为空',
            'post_tpl.*'      => '邮费模板不能为空',
        ];

        validate($params, $rules, $messages);
    }

    /**
     * @param $params
     *
     * @throws \Illuminate\Validation\ValidationException
     */
    public static function status($params)
    {
        // 参数验证规则
        $rules = [
            'id'            => 'required',
            'enterprise_id' => 'required',
        ];

        // 提示语
        $messages = [
            'id.*'            => '门店Id不能为空',
            'enterprise_id.*' => '企业Id不为空',
        ];

        validate($params, $rules, $messages);
    }
}
