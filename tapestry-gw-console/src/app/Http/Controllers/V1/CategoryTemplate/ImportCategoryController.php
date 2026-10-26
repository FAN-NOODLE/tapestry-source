<?php
/**
 * 分类导入
 * Created by PhpStorm.
 * User: 赵洋
 * Date: 2020/3/21
 * Time: 14：47
 */

namespace App\Http\Controllers\V1\CategoryTemplate;

use App\Http\Controllers\V1\BaseController;
use App\Services\CategoryTemplateService;
use App\Traits\ExcelTrait;
use Illuminate\Http\Request;

class ImportCategoryController extends BaseController
{
    use ExcelTrait;

    /**
     * 导入支持模板文件格式
     */
    CONST IMPORT_FILE_EXT = 'xlsx';

    /**
     * 分类导入
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\UnsupportedTypeException
     * @throws \Box\Spout\Reader\Exception\ReaderNotOpenedException
     * @throws \Illuminate\Validation\ValidationException
     * User: zhaoyang@vchangyi.com
     * Date: 2020/3/26 15:07
     */
    public function index(Request $request)
    {

        if (!$this->request->hasFile('excel')) {
            notice(100000, ['未上传excel文件']);
        }
        // 验证参数
        $rule = [
            'cate_tpl_id' => 'required|integer',
        ];

        $msg = [
            'cate_tpl_id.*' => '分类模版id不存在'
        ];

        $input = $this->validate($request, $rule, $msg);
        $file  = $this->request->file('excel');
        $ext = $file->getClientOriginalExtension();
        if ($ext != self::IMPORT_FILE_EXT) {
            notice(100000, ['请使用下载的模板文件']);
        }
        // 保存临时文件
        $saveName = uniqid('excel-') . '.' . $ext;
        $file->move(storage_path('logs/gw-console'), $saveName);
        $savePath = storage_path('logs/gw-console/' . $saveName);

        // 解析excel
        $excelData = $this->excelRead($savePath);
        // 删除临时文件
        unlink($savePath);
        if (count($excelData) > 500) {
            notice(100000, ['批量处理最多支持500条']);
        }
        $input['enterprise_id'] = $this->enterpriseId;

        $result = (new CategoryTemplateService())->importCategory($excelData, $input);

        return success($result);

    }

    /**
     * 下载分类导入模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function downloadTemplate()
    {
        $path = resource_path('static/分类导入模板.xlsx');

        return response()->download($path);
    }


}
