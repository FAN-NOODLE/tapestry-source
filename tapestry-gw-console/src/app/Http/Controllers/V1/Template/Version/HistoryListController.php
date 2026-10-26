<?php
/**
 * 版本记录列表
 * Created by PhpStorm.
 * User: 程雨佳
 * Date: 2019/05/30
 * Time: 11:14 AM
 */

namespace App\Http\Controllers\V1\Template\Version;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use App\Common\Constant;
use CrsSdk\Shop\Version;

class HistoryListController extends BaseController
{

    public function index(Request $request)
    {
        // 接受参数
        $input['page']          = $request->get('page', Constant::DEFAULT_PAGE);
        $input['size']          = $request->get('size', Constant::DEFAULT_SIZE);
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("查看版本记录列表");

        $list = (new Version())->historyList($input);
        return success($list);
    }
}
