<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/9/23
 * Time: 11:28
 */

namespace App\Http\Controllers\V1\ContentManage;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\ContentManage;

class ActivityListController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        return [];
//        return (new ContentManage())->activityList($input);
    }
}
