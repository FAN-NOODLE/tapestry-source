<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 16:12
 */

namespace App\Http\Controllers\Service\V1\Setting;


class DeleteController extends DataLakeConfigBase
{
    public function checkParam()
    {
        $rules    = [
            'enterprise_id' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id缺失',

        ];
        $this->validate($this->request, $rules, $messages);
    }

    public function index()
    {
        $this->dataLakeConfigService->delete(['enterprise_id'=>$this->enterpriseId]);
        return success([]);
    }
}