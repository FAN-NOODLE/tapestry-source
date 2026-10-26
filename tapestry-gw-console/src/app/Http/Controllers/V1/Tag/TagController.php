<?php

namespace App\Http\Controllers\V1\Tag;

use App\Http\Controllers\V1\BaseController;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

class TagController extends BaseController
{
    protected $requireLogin = false;
    public function exportCrowdTemplate(): BinaryFileResponse
    {
        $path = resource_path('static/标签人群包.xlsx');
        return response()->download($path);
    }
}