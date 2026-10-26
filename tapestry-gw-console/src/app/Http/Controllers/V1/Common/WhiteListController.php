<?php


namespace App\Http\Controllers\V1\Common;


use App\Http\Controllers\V1\BaseController;
use App\Models\WhiteList;

class WhiteListController extends BaseController
{
    protected $requireLogin = false;
    public function index()
    {
        $type = $this->request->post('type', WhiteList::TYPE_CARD);
        $list = WhiteList::query()->where('type', $type)->get()->toArray();

        return success($list);
    }
}
