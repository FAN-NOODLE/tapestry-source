<?php
/**
 * 版本更新列表
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-02
 * Tome : 16:28
 */

namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\AutoRelease;

class UpdateListController extends BaseController
{
    /**
     * @return array|object
     */
    public function index()
    {
        $list = (new AutoRelease())->updateList(['enterprise_id' => $this->enterpriseId]);

        return success($list);
    }
}
