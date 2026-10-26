<?php

namespace App\Http\Controllers\V1\HelpCenter;

use App\Models\HelpCenter;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    public $requireLogin = false;

    public function index()
    {
        $id = $this->request->input('id');
        $route = $this->request->input('route');

        $query = HelpCenter::query()->select(['id', 'route', 'content']);

        if ($id) {
            $query = $query->where('id', '=', $id);
        }
        if ($route) {
            $query = $query->where('route', '=', $route);
        }

        $help = $query->first();
        if ($help) {
            $res = $help->toArray();
        } else {
            $res = [];
        }

        return success($res);
    }
}