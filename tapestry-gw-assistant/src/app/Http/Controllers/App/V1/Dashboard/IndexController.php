<?php
/**
 *小程序首页dashboard
 *
 * @author  纪鑫鑫 <jixinxin@vchangyi.com>
 * @date    2020-08-03 18:57
 */

namespace App\Http\Controllers\App\V1\Dashboard;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\StoreMemberViewTab;
use CrsSdk\Work\GuideManage;
use Illuminate\Http\Request;

/**
 * Class IndexController
 * @package App\Http\Controllers\App\V1\Dashboard
 * @author  纪鑫鑫 <jixinxin@vchangyi.com>
 * @date    2020-08-03 18:57
 */
class IndexController extends BaseController
{
    /**
     * 小程序首页数据获取
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @author  纪鑫鑫 <jixinxin@vchangyi.com>
     * @date    2020-08-10 11:19
     */
    public function index(Request $request)
    {
        //根据手机号获取员工信息
        $storeId    = $request->get('store_id', 0);
        $memberInfo = $this->storeWorkerDetail($this->enterpriseId, $storeId);
        $tabTime = StoreMemberViewTab::query()
            ->where(
                [
                    'enterprise_id'   => $this->enterpriseId,
                    'store_id'        => $storeId,
                    'store_worker_id' => $memberInfo['id']??0,
                ]
            )->get(['type', 'extension_view_time'])->toArray();
        if (!empty($tabTime)) {
            $tabTime = array_column($tabTime, 'extension_view_time', 'type');
        }

        $role = in_array($this->role, [Constant::LOGIN_ROLE_STORE_MANAGER, Constant::LOGIN_ROLE_STORE_MEMBER]) ? true : false;
        $params = [
            'role'               => $role,
            'enterprise_id'      => $this->enterpriseId,
            'worker_id'          => $memberInfo['id']??0,
            'mem_uid'            => $memberInfo['mem_uid'] ?? '',
            'store_id'           => $storeId,
            'contact_qr_code_id' => $memberInfo['contact_qr_code_id']??0,
            'tab_time'           => $tabTime[StoreMemberViewTab::CUSTOMER_TAB_TYPE] ?? 0,
            'service_tab_time'   => $tabTime[StoreMemberViewTab::SERVICE_ORDER_TO_BE_RECEIVED_TAB_TYPE] ?? 0,
            'extension_tab_time' => $tabTime[StoreMemberViewTab::EXTENSION_TAB_TYPE] ?? 0
        ];

        $data = (new GuideManage())->getDashboard($params);
        $data['avatar'] = $memberInfo['avatar'];
        $data['customer_view_tab'] = isset($tabTime[StoreMemberViewTab::CUSTOMER_TOTAL_TAB_TYPE]) ? false :true;
        return success($data);
    }
}
