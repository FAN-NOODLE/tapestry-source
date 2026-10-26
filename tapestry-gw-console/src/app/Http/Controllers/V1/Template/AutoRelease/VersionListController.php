<?php
/**
 * 版本管理列表
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-02
 * Tome : 16:31
 */

namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\AutoRelease;

class VersionListController extends BaseController
{
    /**
     * @return mixed|\Psr\Http\Message\StreamInterface
     */
    public function index()
    {
        $list = (new AutoRelease())->versionList(['enterprise_id' => $this->enterpriseId]);

        return success($list);
    }
}
