<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/13 0013
 * Time: 14:35
 */

namespace App\Http\Controllers\V1\ActivityCommon;

use App\Http\Controllers\V1\BaseController;

class ActivityTypeListController extends BaseController
{

    public function index()
    {
        return success([
            [
                'type' => 'FlashSale',
                'name' => '限时购',
                'value' => '12',
            ],
            [
                'type' => 'Group',
                'name' => '拼团',
                'value' => '14',
            ],
            [
                'type' => 'Haggle',
                'name' => '砍价',
                'value' => '15',
            ],
        ]);
    }

}
