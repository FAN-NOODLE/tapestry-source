<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/24
 * Time: 19:39
 */
namespace App\Http\Controllers\V1\MaterialChannel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\MaterialChannel;

class TypeListController extends BaseController
{
    public function index()
    {
        $list = (new MaterialChannel())->channelTypeList();

        return success($list);
    }
}