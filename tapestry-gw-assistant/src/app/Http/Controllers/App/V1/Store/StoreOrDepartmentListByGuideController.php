<?php
/**
 * Created by PHPStorm.
 * author: baizhongyang
 * Date: 2020/8/20
 * Time: 17:04
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Micros\Work\Work;
use CrsSdk\MicroAuthority\BaseClient;
use CrsSdk\Shop\Store;
use CrsSdk\Work\Department;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\App\V1\BaseController as Controller;

class StoreOrDepartmentListByGuideController extends Controller
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
     * 根据导购返回对应的门店列表或部门列表
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/20 14:40
     */
    public function index()
    {
        $param['name']          = $this->request->input('name', '');
        $param['is_new']        = $this->request->input('is_new', 0);
        $param['enterprise_id'] = $this->loginAdminEnterprise['enterprise_id'];
        $param['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $param['admin_id']      = $this->loginAdminEnterprise['admin_id'];
        $param['mobile']        = $this->mobile;
        $param['store_id']      = $this->request->input('store_id', 0);
        $param['pid']           = $this->request->input('pid', 0);

        Log::info('导购身份判断请求参数', [$param]);


        // 获取企业管理者权限范围数据
        $dataRange = (new Account())->dataRangeList(['admin_id' => $param['admin_id']]);
        if (empty($dataRange)) {
            return success(['list' => []]);
        }
        $dpIds = array_column($dataRange, 'object_id');

        if (empty($param['pid'])) {
            $param['dp_ids'] = $dpIds;
        } else {
            $param['dp_ids'] = [$param['pid']];
        }
        $list = (new Work())->thirdIdsList($param);

        return success(['list' => $list]);

//         获取部门列表
//        $condition = [
//            'qy_domain'      => $param['qy_domain'],
//            'app_identifier' => 'contact',
//            'params'         => [
//                'pageSize' => 999,
//                'dpName'   => $param['name']
//            ],
//        ];
//        if (empty($param['pid'])) {
//            $condition['params']['thirds'] = $dpIds;
//        } else {
//            $condition['params']['thirds'] = [$param['pid']];
//        }
//
//        $list = (new Department())->departmentList($condition);

//        $data = [];
//        foreach ($list['list'] as $value) {
//            $data[] = [
//                'id'       => $value['dpId'],
//                'name'     => $value['dpName'],
//                'third_id' => (int)$value['dpThirdid']
//            ];
//        }
//
//        return success(['list' => $data]);

    }

    /**
     * 查询部門列表
     *
     * @param $condition
     *
     * @return array
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/23 18:36
     */
    protected function departmentList($condition)
    {
        $page           = 1;
        $maxPages       = 0;
        $departmentList = [];
        do {
            $condition['params']['pageNum']  = $page;
            $condition['params']['pageSize'] = 200;

            // 结果处理
            $result = (new Department())->departmentList($condition);

            $departmentList = array_merge($departmentList, $result['list']);
            $maxPages       = $maxPages ?: ceil($result['total'] / 200);

            // 步程自增
            $page++;
        } while ($page <= $maxPages);

        return $departmentList;
    }
}
