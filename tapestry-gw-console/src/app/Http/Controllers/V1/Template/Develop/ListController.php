<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-09-26
 * Tome : 17:30
 */

namespace App\Http\Controllers\V1\Template\Develop;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\DevelopCode;

class ListController extends BaseController
{
    public function index()
    {
        $page = $this->request->input('page', 1);
        $size = $this->request->input('size', 10);
        $list = (new DevelopCode())->list([
            'page'          => $page,
            'size'          => $size,
            'enterprise_id' => $this->enterpriseId
        ]);
        return success($list);
    }
}
