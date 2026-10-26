<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/19
 * Time: 18:49
 */
namespace App\Http\Controllers\V1\MaterialChannel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use App\Exceptions\Error;
use Illuminate\Support\Facades\Log;

class UpdateChannelStockController extends BaseController
{
    /**
     * 库存的三种操作方式
     */
    const STOCK_DECREASE = 1;   // 减少
    const STOCK_INCREASE = 2;   // 增加
    const STOCK_REPLACE = 3;    // 覆盖

    const ONLINE_ON = 1; //线上
    const ONLINE_OFF = 2; //线下

    /**
     * 是否要操作物料
     */
    const HANDLE_MATERIAL_YES = 1; // 要操作
    const HANDLE_MATERIAL_NO = 2;  // 不操作

    /**
     * 更改渠道库存
     *
     * @return array|\Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->validateInput();

        try {

            (new MaterialChannel())->updateChannelStock($input);

        } catch (Error $error) {

            Log::error('更改渠道库存失败', $error->data);
            notice(450000, $error->getMessage(), [], 200);

        } catch (\Exception $e) {

            Log::error('更改渠道库存失败' . $e->getMessage(), $input);
            notice(450000, $e->getMessage(), [], 200);

        }

        return success();
    }

    protected function validateInput()
    {
        $input = $this->request->input();

        // 参数验证规则
        $rules = [
            'sku_id'    => 'required|integer',
            'source_id' => 'required|integer',
            'number'    => 'required|integer',
            'type'      => 'required|integer|in:' . self::STOCK_INCREASE . ',' . self::STOCK_DECREASE . ',' . self::STOCK_REPLACE,
            'online'    => 'required|integer|in:' . self::ONLINE_ON . ',' . self::ONLINE_OFF,
        ];

        // 提示语
        $messages = [
            'sku_id.*'    => '请输入sku_id',
            'source_id.*' => '请输入正确的来源ID',
            'number.*'    => '操作数量有误',
            'type.*'      => '请输入正确的操作类型',
            'online.*'    => '请选择要操作的线上线下类型',
        ];

        $this->validate($this->request, $rules, $messages);

        if (in_array($input['type'], [self::STOCK_DECREASE, self::STOCK_INCREASE]) && $input['number'] <= 0) {
            notice(450000, ['操作数量必须大于0']);
        }

        $input['enterprise_id'] = $this->enterpriseId;
        $input['object_id'] = $input['sku_id'];
        $input['is_handle_material'] = self::HANDLE_MATERIAL_YES;
        $input['operate_id'] = $this->enterpriseAdminId;

        return $input;
    }

}
