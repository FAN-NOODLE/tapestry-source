<?php
/**
 * DetailController.php
 *
 * Created on 2019-05-31 11:16
 * Created by shizhixiong@vchangyi.com
 */

namespace App\Http\Controllers\V1\Template\Category;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Template;

class DetailController extends BaseController
{
    /**
     * 模板分类详情
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     *
     * @author shizhixiong@vchangyi.com
     *
     * @date   2019-05-31 11:54
     */
    public function index()
    {
        $params = request()->input();
        $params['enterprise_id'] = $this->enterpriseId;

        return (new Template())->categoryDetail($params);
    }
}
