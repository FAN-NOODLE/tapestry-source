<?php

namespace App\Http\Controllers\V1\Notification;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

/**
 * 传讯-白名单
 */
class DownloadCrowdTemplateController extends BaseController
{

    protected $requireLogin = false;


    use ExcelTrait;

    /**
     * 下载模板
     * @return BinaryFileResponse
     */
    public function index(): BinaryFileResponse
    {
        $path = resource_path('static/传讯-白名单.xlsx');
        return response()->download($path);
    }



}