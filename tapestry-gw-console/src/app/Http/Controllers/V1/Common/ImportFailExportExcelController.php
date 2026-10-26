<?php

namespace App\Http\Controllers\V1\Common;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;

class ImportFailExportExcelController extends BaseController
{
    use ExcelTrait;

    /**
     * 导出Excel 用于导出失败前端把失败数据传过来进行导出
     *
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->all();

        $this->validate($this->request, [
            'title' => 'required',
            'data'  => 'required'
        ]);

        $this->excelExport($input['title'], $input['data']);
    }
}
