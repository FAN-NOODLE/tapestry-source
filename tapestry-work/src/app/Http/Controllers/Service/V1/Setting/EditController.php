<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 15:34
 */

namespace App\Http\Controllers\Service\V1\Setting;


class EditController extends DataLakeConfigBase
{
    public function checkParam()
    {
        $rules    = [
            'enterprise_id' => 'required',
            'key'=>'required',
            'value'=>'required'
        ];
        $messages = [
            'enterprise_id.required' => '企业id缺失',
            'key.required' => '编辑字段名称缺失',
            'value.required' => '编辑字段值缺失',

        ];
        $this->validate($this->request, $rules, $messages);
    }

    public function index()
    {
        $key = $this->request->input('key');
        $value = $this->request->input('value');
        $this->dataLakeConfigService->edit([$key=>$value],['enterprise_id'=>$this->enterpriseId]);
        return success([]);
    }
}