<?php


namespace App\Http\Controllers\V1\Message;


use App\Http\Controllers\V1\BaseController;
use App\Services\ActivityMessageService;

class MiniSyncMessageNumberController extends BaseController
{
    public function index()
    {
        $result = (new ActivityMessageService())->SyncMessageNumbers();
        return success($result);
    }
}