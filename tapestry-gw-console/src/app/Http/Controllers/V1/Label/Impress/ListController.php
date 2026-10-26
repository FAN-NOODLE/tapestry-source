<?php
/**
 * Created by sublime.
 * User: zhao qi
 * Date: 2019/01/17
 * Time: 10:00
 */
namespace App\Http\Controllers\V1\Label\Impress;

use CrsSdk\Shop\ImpressLabel;
use App\Http\Controllers\App\V1\Label\LabelBaseController;

class ListController extends LabelBaseController
{
    /**
     * 印象标签列表
     *
     * @return array
     */
    public function index()
    {
        $params['enterprise_id'] = $this->enterpriseId;
        return (new ImpressLabel())->list($params);
    }
}
