<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/17
 * Time: 18:38
 */
namespace App\Http\Controllers\V1\MaterialChannel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use App\Exceptions\Error;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;

class ListController extends BaseController
{
    /**
     * 渠道列表
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->validateInput();
        try {

            $data = (new MaterialChannel())->MaterialChannelList($input);

            // 匹配门店名称
            $storeList = (new Store())->lists([
                'enterprise_id' => $this->enterpriseId,
                'store_ids'     => array_column($data['list'], 'source_id'),
                'size'          => 999999
            ])['list'];
            $storeList = array_column($storeList, 'name', 'id');
            foreach ($data['list'] as $key => $value) {
                $data['list'][$key]['store_name'] = $storeList[$value['source_id']] ?? '商城或其他';
            }
        } catch (Error $error) {

            Log::error('获取渠道列表失败', $error->data);
            notice(100000, $error->getMessage());

        } catch (\Exception $e) {

            Log::error('获取渠道列表失败'.$e->getMessage(), $input);
            notice(100000, $e->getMessage());

        }

        return success($data);
    }

    /**
     * 参数验证
     *
     * @param $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    protected function validateInput()
    {
        $rule = [
            'id'                  => 'integer',
            'material_id'         => 'array',
            'sku_id'              => 'integer',
            'source_ids'          => 'array',
            'channel_type'        => 'integer',
            'status'              => 'integer|in:1,2,3,4',
            'time'                => 'array',
            'time.type'           => 'integer|in:1,2',
            'time.start'          => 'integer',
            'time.end'            => 'integer',
            'goods'               => 'array',
            'goods.interval_type' => 'integer|in:1,2,3',
            'page'                => 'integer',
            'size'                => 'integer',
        ];

        $msg = [
            'material_id.*'        => "物料ID必须是数组",
            'sku_id.*'             => "商品不存在",
            'id.integer'           => '渠道ID必须是数字',
            'source_ids.array'     => '来源ID必须是数组',
            'channel_type.integer' => '渠道类型必须是数字',
            'status.*'             => '渠道状态有误',
            'time.*'               => '时间区间搜索有误',
            'goods.*'              => '货品区间搜索有误',
        ];

        $this->validate($this->request, $rule, $msg);
        $input         = $this->request->all();
        $input['page'] = $this->request->input('page', 1);
        $input['size'] = $this->request->input('size', 10);

//        // 货品搜索
//        if (isset($input['goods']) && is_array($input['goods'])) {
//            if (isset($input['goods']['interval_type']) && isset($input['goods']['interval_start']) && isset($input['goods']['interval_end']) && $input['goods']['interval_start'] <= $input['goods']['interval_end']) {
//                switch ($input['goods']['interval_type']) {
//                    case 3:
//                        $money_reg = '/^(([1-9]{1}\d*)|(0{1}))(\.\d{1,2})?$/';
//                        if (!preg_match($money_reg, $input['goods']['interval_start']) || !preg_match($money_reg, $input['goods']['interval_end'])) {
//                            notice(100000, ['渠道价搜索金额有误']);
//                        }
//                        break;
//                    default:
//                        if (!is_numeric($input['goods']['interval_start']) || !is_numeric($input['goods']['interval_end'])) {
//                            notice('interface_exception', ['渠道量搜索必须是数字']);
//                        }
//                        break;
//                }
//            } else {
//                notice('100000', ['货品渠道搜索有误']);
//            }
       // }

        // 时间搜索
        if (isset($input['time']) && is_array($input['time'])) {
            if (isset($input['time']['end']) && isset($input['time']['start'])) {
                if ($input['time']['start'] > $input['time']['end']) {
                    notice('interface_exception', ['开始时间不能大于结束时间']);
                }
            }
        }

        $input['enterprise_id'] = $this->enterpriseId;
        return $input;
    }
}