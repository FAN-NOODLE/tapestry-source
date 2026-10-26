<?php

namespace App\Http\Controllers\V1\Excitation;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

/**
 * 手动积分
 */
class DownloadManualCoinsTemplateController extends BaseController
{

    protected $requireLogin = false;


    use ExcelTrait;

    /**
     * 下载模板
     * @return BinaryFileResponse
     */
    public function index(): BinaryFileResponse
    {
        $path = resource_path('static/导购金币-手动补发金币导入模板.xlsx');
        return response()->download($path);
    }



}