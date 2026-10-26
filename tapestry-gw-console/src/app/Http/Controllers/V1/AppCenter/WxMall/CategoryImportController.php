<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2019-07-17
 * Time: 10:37
 */

namespace App\Http\Controllers\V1\AppCenter\WxMall;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxMall;
use Illuminate\Http\Request;

class CategoryImportController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        // 参数验证
        $this->validate($request, ['file' => 'required'], ['file.required' => '请上传好物圈类目文件']);

        $filePath = $request->file->path();
        // 好物圈类目数据
        $list = [];
        // 读取文件，列表赋值
        if (($handle = fopen($filePath, 'r')) !== false) {
            while (($data = fgetcsv($handle, 1000, ',')) !== false) {
                $num = count($data);
                $arr = [];
                for ($i = 0; $i < $num; $i++) {
                    $arr[] = $data[$i];
                }

                if (!empty($arr)) {
                    $list[] = $arr;
                }
            }

            fclose($handle);
        }

        // 移除标题行
        if (!empty($list)) {
            unset($list[0]);
        }

        // 导入数据为空，请确认导入表格数据
        if (empty($list)) {
            notice(920101);
        }

        // 导入好物圈物品类目列表
        $result = (new WxMall())->wxMallCategoryImport(['data' => $list]);

        return success($result);
    }
}
