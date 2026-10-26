<?php
/**
 * EditController.php
 *
 * Created on 2019-05-31 11:15
 * Created by shizhixiong@vchangyi.com
 */

namespace App\Http\Controllers\V1\Template\Category;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Template;

class EditController extends BaseController
{
    /**
     * 模板分类编辑
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     *
     * @author shizhixiong@vchangyi.com
     *
     * @date   2019-05-31 11:51
     */
    public function index()
    {
        $params = request()->input();
        $params['enterprise_id'] = $this->enterpriseId;

        return (new Template())->categoryEdit($params);
    }
}
