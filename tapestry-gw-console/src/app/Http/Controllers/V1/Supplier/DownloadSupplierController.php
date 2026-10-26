<?php


namespace App\Http\Controllers\V1\Supplier;

use App\Http\Controllers\V1\BaseController;

class DownloadSupplierController extends BaseController
{
    public function index()
    {
        $path = resource_path('static/导入供货物料模板.xlsx');

        return response()->download($path);
    }
}