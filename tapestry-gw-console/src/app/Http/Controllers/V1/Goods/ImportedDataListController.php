<?php
/**
 * 批量导入的数据列表
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ImportTrait;
use Illuminate\Http\Request;

class ImportedDataListController extends BaseController
{
    use ImportTrait;

    public function index(Request $request)
    {
        $params = $this->getParams($request);

        $page = $params['page'] ?? 1;
        $size = $params['size'] ?? 5;
        $type = $params['type'] ?? 1;
        $key  = $params['key'];

        $total = $this->getDataTotal($key, $type);
        $list  = $this->getDataList($key, $type, $page, $size);

        return success(['data' => [
            'key'   => $key,
            'page'  => (int)$page,
            'size'  => (int)$size,
            'total' => $total,
            'list' => $list,
        ]]);
    }

    private function getParams(Request $request)
    {
        $input = $request->input();

        // 参数验证规则
        $rules = [
            'key'  => 'required|size:32',
            'type' => 'integer|in:1,2',
            'page' => 'integer|min:0',
            'size' => 'integer|min:0',
        ];

        // 提示语
        $messages = [
            'key.*'  => '键值必填',
            'type.*' => '参数有误',
            'page.*' => '数值类型',
            'size.*' => '数值类型',
        ];

        $this->validate($request, $rules, $messages);

        return $input;
    }

}
