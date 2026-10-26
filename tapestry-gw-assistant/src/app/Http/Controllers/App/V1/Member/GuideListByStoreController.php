<?php
/**
 * Created by PHPStorm.
 * author: baizhongyang
 * Date: 2020/8/18
 * Time: 17:06
 */


namespace App\Http\Controllers\App\V1\Member;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

//use CrsSdk\Shop\Member;

class GuideListByStoreController extends Controller
{
    /**
     * 企业关联数据范围类型：0-全部门店 1-区域 2-门店  3-商城  4 部门
     */
    const DATA_RANGE_ALL_STORE  = 0;
    const DATA_RANGE_AREA       = 1;
    const DATA_RANGE_STORE      = 2;
    const DATA_RANGE_MALL       = 3;
    const DATA_RANGE_DEPARTMENT = 4;

    /**
     * 导购列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws ValidationException
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/22 16:56
     */
    public function index(Request $request)
    {
        $input    = $request->input();
        $rules    = [
            'dp_ids' => 'array',
        ];
        $messages = [
            'dp_ids.array' => '门店id或部门id为数组',
        ];
        $this->validate($request, $rules, $messages);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['domain']        = $this->loginAdminEnterprise['domain'];
        $admin_id               = $this->loginAdminEnterprise['admin_id'];

        if (isset($input['dp_ids']) && !empty($input['dp_ids'])) {
            $input['dp_ids'] = array_filter($input['dp_ids'], function ($item) {
                return !empty($item);
            });
        }
        if (empty($input['dp_ids'])) {

            $dataRange      = (new Account())->dataRangeList(['admin_id' => $admin_id]);
            $departmentList = [];
            if (!empty($dataRange)) {
                foreach ($dataRange as $item) {
                    $departmentList[] = $item['object_id'];
                }
            }
            $input['dp_ids'] = $departmentList;
        }

        // 获取导购列表
        $guideList = (new StoreWorker())->guideListByStore($input);

        return success(['list' => $guideList]);
    }
}
