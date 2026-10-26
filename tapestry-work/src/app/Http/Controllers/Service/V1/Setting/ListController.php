<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 16:09
 */

namespace App\Http\Controllers\Service\V1\Setting;


use App\Common\Constant;

class ListController extends DataLakeConfigBase
{
    public function checkParam()
    {
        $rules    = [
            'enterprise_id' => 'required',
            'page'          => 'integer',
            'size'          => 'integer',

        ];
        $messages = [
            'enterprise_id.required' => '企业id缺失',
            'page.integer'           => '分页为整型',
            'size.integer'           => '每页数量为整型',
        ];
        $this->validate($this->request, $rules, $messages);
    }

    public function index()
    {
        $page = $this->request->input('page', Constant::DEFAULT_PAGE);
        $size = $this->request->input('size', Constant::DEFAULT_SIZE);
        $enterpriseId = $this->request->input('enterprise_id');
        return success($this->dataLakeConfigService->list($page, $size,['enterprise_id'=>$enterpriseId]));
    }
}