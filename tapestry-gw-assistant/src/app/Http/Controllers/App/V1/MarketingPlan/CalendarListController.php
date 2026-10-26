<?php
/**
 * 获取营销日历
 */
namespace App\Http\Controllers\App\V1\MarketingPlan;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Work\MarketingPlan;
use Illuminate\Http\Request;

class CalendarListController extends Controller
{
    /**
     * 获取营销日历
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/13 17:44
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input = $request->input();
        $input['current_date'] = $request->input('current_date', '');
        $input['domain'] = $this->loginAdminEnterprise['domain'];
        $input['enterprise_id'] = $this->enterpriseId;
        $storeId = $this->request->post('store_id', 0);
        $input['store_id'] = $storeId;
        $input['mobile'] = $this->mobile;

        $result = (new MarketingPlan())->taskCalendar($input);

        return success($result['data'] ?? $result);
    }
}
