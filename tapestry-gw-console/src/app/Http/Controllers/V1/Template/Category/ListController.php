<?php
/**
 * ListController.php
 *
 * Created on 2019-05-31 11:15
 * Created by shizhixiong@vchangyi.com
 */

namespace App\Http\Controllers\V1\Template\Category;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Template;

class ListController extends BaseController
{
    /**
     * 模板分类列表
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     *
     * @author shizhixiong@vchangyi.com
     *
     * @date   2019-05-31 11:27
     */
    public function index()
    {
        $page = request()->post('page', 1);
        $size = request()->post('size', 10);
        $search = request()->post('search', []);
        $search['enterprise_id'] = $this->enterpriseId;

        return (new Template())->categoryList([
            'page' => $page,
            'size' => $size,
            'search' => $search
        ]);
    }
}
