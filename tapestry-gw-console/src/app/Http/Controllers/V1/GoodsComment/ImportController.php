<?php


namespace App\Http\Controllers\V1\GoodsComment;


use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\MicroGoodsV1\Comment;

class ImportController extends BaseController
{
    use ExcelTrait;

    /**
     * 是否要保存草稿 1：是
     */
    const DRAFT_YES = 1;

    /**
     * 是否要保存草稿 2：否
     */
    const DRAFT_NO = 2;

    /**
     * 模板表头
     */
    const EXCEL_TITLE = [
        '评论者昵称（最多10个字）',
        '评论时间（例：2020-1-1 00:00:01）',
        '评论内容（最多200个字）'
    ];

    /**
     * 批量导入
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        if (!$this->request->hasFile('excel')) {
            notice(100000, ['未上传excel文件']);
        }

        $file = $this->request->file('excel');

        // 保存临时文件
        $saveName = uniqid('excel-') . '.' . $file->getClientOriginalExtension();
        $file->move(storage_path('logs/gw-console'), $saveName);
        $savePath = storage_path('logs/gw-console/' . $saveName);

        // 解析excel
        $excelData = $this->excelRead($savePath);

        // 删除临时文件
        unlink($savePath);

        if (count($excelData) > 500) {
            notice(100000, ['批量处理最多支持500条']);
        }

        // 判断excel模板是否被修改
        if (array_slice($excelData[0], 0, count(self::EXCEL_TITLE)) !== self::EXCEL_TITLE) {
            notice(100000, ['模板表头内容错误,请重新下载模板']);
        }

        // 开始处理数据
        $result = [
            'success_count' => 0,
            'error_count'   => 0,
            'error_list'    => [],
            'title'         => array_merge(self::EXCEL_TITLE, ['失败原因'])
        ];

        $commentSdk = new Comment();
        foreach ($excelData as $key => $value) {

            // 过滤表头
            if ($key === 0) {
                continue;
            }

            try {

                if (!empty($value[1]) && !strtotime($value[1])) {
                    notice(100000, ['评论时间格式有误']);
                }

                $params = [
                    'enterprise_id' => $this->enterpriseId,
                    'nickname'      => $value[0],
                    'comment_at'    => empty($value[1]) ? 0 : strtotime($value[1]),
                    'content'       => $value[2],
                    'is_draft'      => self::DRAFT_YES
                ];

                $params = collect($params)->map(function ($value) {
                    return trim($value) == $value ? $value : trim($value);
                })->toArray();

                $commentSdk->addComment($params);

                $result['success_count']++;

            } catch (\Exception $e) {
                $value = array_slice($value, 0, count(self::EXCEL_TITLE));
                $value[] = $e->getMessage();
                $result['error_list'][] = $value;
                $result['error_count']++;
            }
        }

        return success($result);
    }

    /**
     * 下载模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function downloadAddTemplate()
    {
        $path = resource_path('static/商品评论批量导入模板.xlsx');

        return response()->download($path);
    }
}