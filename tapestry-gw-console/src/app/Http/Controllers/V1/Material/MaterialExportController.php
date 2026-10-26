<?php


namespace App\Http\Controllers\V1\Material;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\MaterialExportJob;
use App\Services\MaterialService;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use Illuminate\Support\Facades\Log;

class MaterialExportController extends BaseController
{
    /**
     * 搜索类型：free_stock-可分配库存 cost_price-成本价
     */
    const SEARCH_FREE_STOCK = 'free_stock';
    const SEARCH_COST_PRICE = 'cost_price';

    /**
     * 筛选条件类型 :spu_no-货品编号 material_no-物料编号 name-品名
     */
    const SPU_NO = 'spu_no';
    const MATERIAL_NO = 'material_no';
    const NAME = 'name';

    /**
     * 物料权限状态 1-启用 2-禁用
     */
    const material_auth_status = 1;

    /**
     * 文件名称（列表显示和导出名称）
     * @var
     */
    const FILE_NAME = '库存管理-物料列表';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'gw-console/v1/material/list';

    /**
     * 总库列表出tag
     */
    const TAG = 'material-list';

    /**
     * 总库管理列表导出
     * @return \Illuminate\Http\JsonResponse
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\InvalidSheetNameException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function index()
    {
        $input = $this->request->only([
            'first_type', // 第一个筛选条件类型 :spu_no-货品编号 material_no-物料编号 name-品名（ 'spu_no','material_no','name',）
            'first_value',
            'second_type', // 第二个筛选条件类型：total_stock-库存总量 free_stock-可分配库存 cost_price-固定成本价
            'min_value',
            'max_value',
            'source_id', // 归属渠道id
            'time_type', // 时间筛选条件类型: create-创建时间 update-更新时间
            'start_time',
            'end_time',
            'tag',
            'route_url'
        ]);

        $input['tag'] = $input['tag'] ?? self::TAG;
        $input['route_url'] = $input['route_url'] ?? self::ROUTE_URL;

        // 判断筛选条件货品编号/物料编号/品名
        if (
            isset($input['first_type'])
            && !empty($input['first_type'])
            && in_array($input['first_type'], [self::SPU_NO, self::MATERIAL_NO, self::NAME])
        ) {

            if (empty(trim($input['first_value']))) {
                notice(451519);
            }

            $params['material_type']  = $input['first_type'];
            $params['material_value'] = $input['first_value'];
        }

        // 可分配库存/成本价筛选
        if (isset($input['second_type']) && !empty($input['second_type']) && in_array($input['second_type'], [self::SEARCH_FREE_STOCK, self::SEARCH_COST_PRICE])) {
            // 输入大于零的数字
            if (!is_numeric($input['min_value']) || !is_numeric($input['max_value'])) {
                notice(451511);
            }

            $minValue = isset($input['min_value']) ? max($input['min_value'], 0) : '';
            $maxValue = $input['max_value'] ?? '';

            // 输入值范围最大值比最小值小
            if ($minValue > $maxValue) {
                notice(451510);
            }

            switch ($input['second_type']) {
                case 'total_stock':
                    // ~~~~~~~
                    break;
                case self::SEARCH_FREE_STOCK:
                    $params['start_stock'] = $minValue;
                    $params['stop_stock'] = $maxValue;
                    break;
                case self::SEARCH_COST_PRICE:
                    $params['start_price'] = $minValue;
                    $params['stop_price'] = $maxValue;
                    break;
                default:
                    break;
            }
        }

        // 判断筛选条件创建时间/更新时间
        if (
            isset($input['time_type']) && !empty($input['time_type'])
            && !empty($input['start_time']) && !empty($input['end_time'])
        ) {
            switch ($input['time_type']) {
                case 'create':
                    $params['time_type'] = 'created_at';
                    break;
                case 'update':
                    $params['time_type'] = 'updated_at';
                    break;
                default:
                    break;
            }
            $params['start_time'] = $input['start_time'] ?? 0;
            $params['end_time']   = $input['end_time'] ?? 0;
        }

        // 特殊情况可以获取所有门店物料权限--需要后期根据业务调整
        if (!empty($this->storeIds)) {
            $params['source_id'] = $this->storeIds;
        }

        // 不是admin则只能获取启用的物料
        if ($this->isSuperAdmin == 0) {
            $params['status'] = self::material_auth_status;

        }

        // 超管source_id等于enterprise_id
        if ($this->isSuperAdmin == 1) {
            $params['source_id'] = [$this->enterpriseId];
        }

        // 渠道归属Id--筛选优先
        if (isset($input['source_id']) && $input['source_id'] != 0) {
            //调用获取渠道下所绑定的物料id
            $params['source_id'] = [$input['source_id']];
        }

        $params['size']          = Constant::MAX_DEFAULT_SIZE;
        $params['page']          = Constant::PAGE;
        $params['enterprise_id'] = $this->enterpriseId;
        $params['operate_name']  = $this->enterpriseAdmin['username'];
        $params['operate_id']    = $this->enterpriseAdminId;
        $params['tag']           = $input['tag'];
        $params['route_url'] = $input['route_url'];

        // 开始导出
        //dispatch(new MaterialExportJob($params, $input, '库存管理-物料列表'));
        (new MaterialExportJob($params, $input, self::FILE_NAME))->handle();

        return success();
    }

}