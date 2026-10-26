<?php
/**
 * Created by PhpStorm.
 * User: Leon Liu
 * Date: 2020/6/28
 * Time: 20:17
 */


namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\NotSaleListExportJob;
use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\Shop\DownloadCenter;
use Illuminate\Support\Facades\Log;

class NotSaleListController extends BaseController
{

    /**
     *  导出文件的类型  1：excel
     */
    const EXPORT_TYPE_EXCEL = 1;

    /**
     * 导出文件的后缀
     */
    const EXTENSION = 'xlsx';

    /**
     * 导出文件的表头
     */
    const HEADER = [
        'SKUid',
        '物料编码',
        'sku标题',
        '吊牌价',
        '售卖价'
    ];

    /**
     * 导出文件的主文件名
     */
    const FILENAME = '待上架商品信息';

    /**
     * 导出对应的路由
     */
    const ROUTE_URL = 'goods/not-sale-list';

    /**
     * 标签
     */
    const TAG = 'goods-not-sale-list';

    /**
     * 按名称搜索
     */
    const TYPE_SEARCH_NAME = 1;

    /**
     * 下载未上架商品列表（下载中心）
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        // 获取对应商品数据
        $input['enterprise_id'] = $this->enterpriseId;

        // 待上架或已下架
        $input['type'] = 4;
        $input['size'] = 1000;

        // 去除活动商品
        $input['search']['channel_type'] = [1];
        try {
            $lists = (new Goods())->list($input);
        } catch (Error $error) {
            Log::error('获取商品列表失败', $error->data);
            notice(451508, [], $error->data);
        } catch (\Exception $exception) {
            Log::error('获取商品列表失败' . $exception->getMessage(), $input);
            notice(451508);
        }

        // 数据拼接
        $data = [];
        $goods = $lists['list'];
        foreach ($goods as $good) {
            $data[] = [
                $good['sku_id'], $good['material_no'], $good['sku_name'], $good['sign_price']/100, $good['info'][0]['price']/100
            ];
        }

        // 生成对应下载中心的文件名，格式为：待上架商品信息202007030014.xlsx
        $name = self::FILENAME . date('Ymd');
        $search['page'] = 1;
        $search['size'] = 1;
        $search['enterprise_id'] = $this->enterpriseId;
        $search['search']['type'] = self::TYPE_SEARCH_NAME;
        $search['search']['name'] = self::FILENAME;
        $fileName = (new DownloadCenter())->lists($search)['list'][0]['name'] ?? "";
        preg_match('/^[^\d]+\d{8}(\d{4})/', $fileName, $match);

        // 后四位是文件序号
        if ($match) {
            $maxNo = $match[1];
            $nowNo = sprintf("%04d", ($maxNo + 1));
        } else {
            $nowNo = '0001';
        }
        $name .= $nowNo;

        // 通知下载中心要下载的文件名
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'operator_id' => $this->enterpriseAdmin['id'],
            'operator_name' => $this->enterpriseAdmin['username'],
            'type' => self::EXPORT_TYPE_EXCEL,
            'name' => $name . '.' . self::EXTENSION,
            'route_url' => self::ROUTE_URL,
            'tag' => self::TAG,
            'param_json' => json_encode([])
        ];

        // 发布下载中心生成文件的Job
        // (new NotSaleListExportJob($data, self::HEADER, $params,$name))->handle();
        try {
            dispatch(new NotSaleListExportJob($data, self::HEADER, $params, $name));
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }

        return success();
    }
}
