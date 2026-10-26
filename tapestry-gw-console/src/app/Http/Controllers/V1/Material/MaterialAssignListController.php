<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/25
 * Time: 14:17
 */
namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;

class MaterialAssignListController extends BaseController
{
    /**
     * 获取物料分派列表
     *
     * @return array
     */
    public function index()
    {
        $input = $this->validates();
        try {
            $input['enterprise_id'] = $this->enterpriseId;
            $list = (new Material())->assignList($input);

            // 匹配门店名称
            $storeList = (new Store())->lists([
                'enterprise_id' => $this->enterpriseId,
                'store_ids'     => array_column($list['list'], 'source_id'),
                'size'          => 999999
            ])['list'];

            $storeList = array_column($storeList, 'name', 'id');

            foreach ($list['list'] as $key => $value) {
                $list['list'][$key]['store_name'] = $storeList[$value['source_id']] ?? '商城';
            }

        } catch (Error $error) {

            Log::error('获取物料分派列表失败', $error->data);
            notice(450000, $error->getMessage());

        } catch (\Exception $e) {

            Log::error('获取物料分派列表失败' . $e->getMessage(), $input);
            notice(450000, $e->getMessage());

        }

        return success($list);
    }

    protected function validates()
    {
        $input = $this->request->all();

        // 参数验证规则
        $rules = [
            'page'          => 'integer',
            'source_id'     => 'integer',
            'size'          => 'integer|max:1000',
            'time_type'     => 'string|in:updated_at,created_at',
            'start_time'    => 'integer',
            'end_time'      => 'integer'
        ];

        // 提示语
        $messages = [];

        $this->validate($this->request, $rules, $messages);

        if (isset($input['start_time']) && isset($input['end_time'])) {
            if ($input['start_time'] > $input['end_time']) {
                notice('interface_exception', ['开始时间不能大于结束时间']);
            }
        }

        $input['page'] = $this->request->input('page', 1);
        $input['size'] = $this->request->input('size', 10);
        $input['time_type'] = $this->request->input('time_type', '');

        return $input;
    }
}