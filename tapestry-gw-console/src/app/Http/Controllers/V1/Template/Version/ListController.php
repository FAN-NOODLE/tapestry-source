<?php
/**
 * 版本更新列表
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/05/30
 * Time: 11:14 AM
 */

namespace App\Http\Controllers\V1\Template\Version;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Version;

class ListController extends BaseController
{

    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("查看版本更新列表");

        return (new Version())->list($input);
    }
}
