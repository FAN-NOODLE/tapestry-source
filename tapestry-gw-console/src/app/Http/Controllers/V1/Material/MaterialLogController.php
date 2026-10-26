<?php


namespace App\Http\Controllers\V1\Material;


use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\EnterpriseAdmin;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;

class MaterialLogController extends BaseController
{
    /**
     * 总库日志筛选类型 1-操作来源 2-操作内容
     */
    const SEARCH_TYPE_SOURCE = 'from_where';
    const SEARCH_TYPE_CONTENT = 'content';

    /**
     * 用于查询门店id既是enterprise_id得情况，此ID在门店无法搜索到
     */
    const SHOP_ID_STRING = "库存管理";

    /**
     * 请求获取总库日志
     * @api {domain}/gw-console/v1/material/log-list
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->only([
            'material_id',
            'page',
            'size',
            'search_type',
            'keyword',
            'operate_name',
            'start_time',
            'end_time',
            'stock_type',
            'min_value',
            'max_value',
        ]);

        if (!isset($input['material_id']) || empty($input['material_id'])) {
            notice(451518);
        }

        // 库存筛选类型被选择时，取值范围必须被定义且为数值
        if (isset($input['stock_type']) && $input['stock_type'] != 'all') {
            if (!isset($input['min_value']) || !is_numeric($input['min_value']) || !isset($input['max_value']) || !is_numeric($input['max_value'])) {
                notice(451511);
            }
        }

        $input['page']          = $input['page'] ?? Constant::PAGE;
        $input['size']          = $input['size'] ?? Constant::DEFAULT_SIZE;
        $input['search_type']   = $input['search_type'] ?? '';
        $input['keyword']       = isset($input['keyword']) ? trim($input['keyword']) : '';
        $input['operate_name']  = isset($input['operate_name']) ? trim($input['operate_name']) : '';
        $input['start_time']    = $input['start_time'] ?? 0;
        $input['end_time']      = $input['end_time'] ?? 0;
        $input['stock_type']    = $input['stock_type'] ?? '';
        $input['min_value']     = isset($input['min_value']) ? max($input['min_value'], 0) : 0;
        $input['max_value']     = $input['max_value'] ?? 0;
        $input['enterprise_id'] = $this->enterpriseId;

        if ($input['min_value'] > $input['max_value']) {
            notice(451510);
        }
        if (!empty($input['search_type']) && $input['search_type'] != 'all' && empty($input['keyword'])) {
            notice(451519);
        }

        // 搜索类型为按照来源搜索，使用门店名称获取门店列表
        $sourceIds = [];
        if (isset($input['search_type']) && $input['search_type'] == self::SEARCH_TYPE_SOURCE) {
            $storeParam = [
                'name'          => $input['keyword'],
                'enterprise_id' => $this->enterpriseId,
                'page'          => Constant::PAGE,
                'size'          => 1000,
            ];

            $storeList = [];
            try {
                $storeList = (new Store())->lists($storeParam)['list'];
            } catch (\Exception $exception) {
                Log::info('获取渠道归属列表失败' . $exception->getMessage(), $input);
                notice(451520);
            }

            $sourceIds = array_column($storeList, 'id');

            // 企业ID为门店ID得类型--商城
            if (mb_strpos(self::SHOP_ID_STRING, $input['keyword']) !== false) {
                $sourceIds[] = $this->enterpriseId;
                $sourceIds[] = 0; // FIXME：可能存在source_id为0的数据，兼容默认为库存管理操作
            }

            if (empty($sourceIds)) {
                return success([
                    'list' => [],
                    'page'  => $input['page'],
                    'size'  => $input['size'],
                    'total' => 0

                ]);
            }

        }

        $input['source_ids'] = $sourceIds;

        // 请求获取总库日志
        try {
            $logList = (new Material())->logList($input);
        } catch (Error $error) {
            Log::info('获取总库日志失败' . $error->getMessage(), $error->data);
            notice(451517, [], $error->data);
        } catch (\Exception $exception) {
            Log::info('获取总库日志失败' . $exception->getMessage(), $input);
            notice(451517, [], ['message' => $exception->getMessage()]);
        }

        if (empty($logList['list'])) {
            return success([
                'list'  => [],
                'page'  => $logList['page'] ?? $input['page'],
                'size'  => $logList['size'] ?? $input['size'],
                'total' => $logList['total'] ?? 0
            ]);
        }

//        $sourceIds  = array_column($logList['list'], 'source_id');
//        $storeParam = [
//            'store_ids'     => $sourceIds,
//            'enterprise_id' => $this->enterpriseId,
//            'page'          => Constant::PAGE,
//            'size'          => 1000,
//        ];
//        $storeList  = (new Store())->lists($storeParam) ?? [];
//        $storeList  = array_column($storeList['list'], null, 'id');
//
//        $enterpriseAdminSdk = new EnterpriseAdmin();
//        foreach ($logList['list'] as &$value) {
//            $value['operate_from'] = isset($storeList[$value['source_id']]) ? $storeList[$value['source_id']]['name'] : '库存管理';
//            // 查操作人名称
//            if (empty($value['operate_name'])) {
//                $adminInfo = $enterpriseAdminSdk->detail([
//                    'enterprise_id' => $this->enterpriseId,
//                    'admin_id'      => $value['operate_id']
//                ])['data'];
//                $value['operate_name'] = $adminInfo['username'] ?? '';
//            }
//        }

//        unset($value);
        return success([
            'list'  => $logList['list'] ?? [],
            'page'  => $logList['page'] ?? $input['page'],
            'size'  => $logList['size'] ?? $input['size'],
            'total' => $logList['total'] ?? 0
        ]);
    }

}