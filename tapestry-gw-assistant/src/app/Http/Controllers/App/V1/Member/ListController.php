<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 17:06
 */


namespace App\Http\Controllers\App\V1\Member;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Log;

/**
 * 店员列表
 * Class ListController
 * @package App\Http\Controllers\App\V1\Member
 */
class ListController extends Controller
{
    /**
     * 获取店员列表
     * @param Request $request
     * @return array|array[]|\Illuminate\Http\JsonResponse|\Psr\Http\Message\ResponseInterface|void
     * @throws \App\Exceptions\WorkFlowException
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request)
    {
        // 参数校验
        $rules   = [
            'store_id'      => 'required|integer',
            'worker_status' => 'integer',
            'page'          => 'integer',
            'size'          => 'integer',
        ];
        $message = [
            'enterprise_id.integer' => '企业ID参数类型有误',
            'store_id.required'     => '门店ID参数不能为空',
            'store_id.integer'      => '门店ID参数类型有误',
            'worker_status.integer' => 'worker_status参数类型有误',
        ];
        $this->validate($request, $rules, $message);
        $enterpriseId = Auth::user()['enterprise_id'];
        $storeId = $request->get('store_id');
        $workerStatus = $request->get('worker_status');
        $page = $request->get('page', Constant::PAGE);
        $size = $request->get('size', Constant::LIMIT);
        $keywords = $request->get('keywords', '');
        $searchMobile = $searchName = '';
        if (!empty($keywords) && is_numeric($keywords) && strpos($keywords, '.') === false) {
            $searchMobile = $keywords;
        } else {
            $searchName = $keywords;
        }
        try{
            return success((new StoreWorker())->getList([
                'enterprise_id' => $enterpriseId,
                'store_id'      => $storeId,
                'worker_status' => $workerStatus,
                'mobile' => $searchMobile,
                'name' => $searchName,
                'page' => $page,
                'size' => $size,
            ]));
        } catch (\Exception $e) {
            Log::error("获取员工列表错误，错误信息为：{$e->getMessage()}");
            error(655002);
        }
    }
}
