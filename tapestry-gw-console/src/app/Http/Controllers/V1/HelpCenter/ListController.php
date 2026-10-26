<?php

namespace App\Http\Controllers\V1\HelpCenter;

use App\Models\HelpCenter;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    public $requireLogin = false;

    public function index()
    {
        $helps = HelpCenter::query()
            ->select(['id', 'key', 'type', 'content'])
            ->get()
            ->toArray();
        return success(['list' => $helps]);
    }
}