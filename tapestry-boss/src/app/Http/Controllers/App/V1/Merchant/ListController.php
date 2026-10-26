<?php


namespace App\Http\Controllers\App\V1\Merchant;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\MicroAuthority\Merchant;

class ListController extends BaseController
{
    public function index()
    {
        $list = (new Merchant())->merchantList();
        return success([
            'data' => $list
        ]);
    }
}