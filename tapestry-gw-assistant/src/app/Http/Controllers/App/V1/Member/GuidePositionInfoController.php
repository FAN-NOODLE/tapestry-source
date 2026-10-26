<?php
/**
 * Created by PHPStorm.
 * author: baizhongyang
 * Date: 2020/8/20
 * Time: 17:06
 */

namespace App\Http\Controllers\App\V1\Member;

use CrsSdk\MicroAuthority\Account;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\App\V1\BaseController as Controller;

class GuidePositionInfoController extends Controller
{
    /**
     * 企业关联数据范围类型：0-全部门店 1-区域 2-门店  3-商城  4 部门
     */
    const DATA_RANGE_ALL_STORE = 0;
    const DATA_RANGE_AREA = 1;
    const DATA_RANGE_STORE = 2;
    const DATA_RANGE_MALL = 3;
    const DATA_RANGE_DEPARTMENT = 4;

    /**
     * 导购身份和权限类型判定
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/20 14:21
     */
    public function index()
    {
        $param['admin_id']      = $this->loginAdminEnterprise['admin_id'];
        $param['enterprise_id'] = $this->enterpriseId;
        $param['is_new']        = $this->request->input('is_new', 0);
        $param['mobile']        = $this->mobile;
        $param['store_id']      = $this->request->input('store_id', 0);
        $param['role']          = $this->getRole($param['mobile']);
        Log::info('导购身份判断请求参数', [$param]);

        $storeIdList = $departmentList = [];
        // 如果是企业管理者返回权限类型$range: 0-初始化值，1-无权限 2-门店权限 3-部门权限
        $range = 0;
        if ($param['is_new']) {
            if ($param['role'] == 5 || $param['role'] == 4) {
                $range = 3;
            }
        } else {
            if ($param['role'] == 5) {
                $dataRange = (new Account())->dataRangeList(['admin_id' => $param['admin_id']]);
                if (empty($dataRange)) {
                    $range = 1;
                } else {
                    foreach ($dataRange as $item) {
                        switch ($item['type']) {
                            case self::DATA_RANGE_STORE:
                                $storeIdList[] = $item['object_id'];
                                break;
                            case self::DATA_RANGE_DEPARTMENT:
                                $departmentList[] = $item['object_id'];
                                break;
                        }
                    }
                    if (!empty($storeIdList)) {
                        $range = 2;
                    } else {
                        $range = 3;
                    }
                }
            }
        }

        $data = [
            'role'  => $param['role'],
            'range' => $range
        ];

        return success($data);
    }
}
