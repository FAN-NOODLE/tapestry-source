<?php
/**
 * Created by PhpStorm.
 * User: chengyujia
 * Date: 2019/11/20
 * Time: 09:44
 */

namespace App\Http\Controllers\V1\LiveBroadcast;


use App\Http\Controllers\V1\BaseController;
use App\Common\Constant;
use CrsSdk\MicroActivity\Live;
use Illuminate\Http\Request;

class MiNiProgramLiveListController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $input['page'] = $request->get('page', Constant::DEFAULT_PAGE);
        $input['size'] = $request->get('size', Constant::DEFAULT_SIZE);
        $input['enterprise_id'] = $this->enterpriseId;
        try {
            $this->writeOperatorLog('小程序直播列表');
        } catch (\Exception $e) {

        }
        $return = (new Live())->miNiProgramLiveList($input);
        return success($return);
    }
}
