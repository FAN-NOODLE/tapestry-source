<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波、赵杰
 * Date: 2020/4/10
 * Time: 17:06
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

/**
 * 店员详情
 * Class DetailController
 * @package App\Http\Controllers\App\V1\Member
 */
class WorkerDetailController extends Controller
{
    /**
     * 获取店员详情
     *
     * @param Request $request
     * @param StoreWorker $storeWorker
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request, StoreWorker $storeWorker)
    {
        $id      = $request->get('id', 0);

        $wheres = [
            'enterprise_id' => $this->enterpriseId,
            'id'            => $id
        ];

        try {
            $detail = $storeWorker->detail($wheres);
        } catch (\Exception $e) {
            Log::error("调用接口Shop:StoreWork->detail出错，错误信息：{$e->getMessage()}, 错误code：{$e->getCode()}", $wheres);
            error($e->getCode());
        }
        if (empty($detail)) {
            notice('store_no_authority');
        }

        $result = [
            'id'                => $detail['id'],
            'enterprise_id'     => $this->enterpriseId,
            'store_id'          => $detail['store_id'],
            'name'              => $detail['name'],
            'mobile'            => $detail['mobile'],
            'worker_no'         => $detail['worker_no'],
            'qr_code'           => $detail['qr_code'],
            'is_boss'           => $detail['is_boss'],
            'worker_status'     => $detail['worker_status'],
            'avatar'            => $detail['avatar'],
            'is_send_sms'       => $detail['is_send_sms'],
            'sms_mobile'        => $detail['sms_mobile'],
            'store_sun_code'    => $detail['store_sun_code'],
        ];

        return success($result);
    }
}
