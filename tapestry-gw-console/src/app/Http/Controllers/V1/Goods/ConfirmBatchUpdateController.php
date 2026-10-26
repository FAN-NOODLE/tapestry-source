<?php
/**
 * 确认批量修改商品
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2020/3/16
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Micros\GoodsMicro;
use App\Traits\ImportTrait;
use Illuminate\Http\Request;

class ConfirmBatchUpdateController extends BaseController
{
    use ImportTrait;

    public function index(Request $request)
    {
        $params = $this->getParams($request);
        $key    = $params['key'];

        $data = $this->getDataList($key, 1);

        $successTotal = 0;
        $errorTotal   = 0;
        $errorList    = [];

        foreach ($data as $v) {
            try {
                if (!isset($v['sku_id'])) {
                    continue;
                }
                array_walk($v, function (&$val, $key) {
                    $val = trim($val);
                });

                $skuId    = $v['sku_id'];
                $skuName  = $v['sku_name'];
                $skuPrice = $v['sku_price'];

                $editChannelParams = [
                    'enterprise_id' => $this->enterpriseId,
                    'source_id'     => $params['source_id'],
                    'channel_type'  => 1,
                    'object_id'     => $skuId,
                    'object_name'   => $skuName,
                    'operate_id'    => $this->enterpriseAdminId,
                    'data'          => [
                        [
                            'sku_id'         => $skuId,
                            'on_line_price'  => $skuPrice,
                            'off_line_price' => $skuPrice,
                            'type'           => 'update',
                        ]
                    ]
                ];

                // to fix  目前只支持单个修改 (后期规划协程)
                (new GoodsMicro($editChannelParams))->editChannel();

                $successTotal++;
            } catch (Error $e) {
                $errorTotal++;
                $v['reason'] = $e->getMessage();
                $errorList[] = $v;
            } catch (\Exception $e) {
                $errorTotal++;
                $v['reason'] = $e->getMessage();
                $errorList[] = $v;
            }
        }

        return success(
            [
                'success_total' => $successTotal,
                'error_total'   => $errorTotal,
                'error_list'    => $errorList,
            ]
        );
    }

    private function getParams(Request $request)
    {
        $input = $request->input();

        // 参数验证规则
        $rules = [
            'key'       => 'required',
            'source_id' => 'required|integer',
        ];

        // 提示语
        $messages = [
            'key.*'       => '键值必填',
            'source_id.*' => '渠道参数必填',
        ];

        $this->validate($request, $rules, $messages);

        return $input;
    }
}
