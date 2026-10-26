<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/24
 * Time: 20:28
 */
namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Material;
use Illuminate\Support\Facades\Log;

class UpdateMaterialRuleController extends BaseController
{
    /**
     * 变更渠道物料权限状态
     *
     * @return array
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'status'
        ]);
        $input['enterprise_id'] = $this->enterpriseId;

        try {

            (new Material())->updateRuleStatus($input);

        } catch (Error $error) {

            Log::error('变更渠道物料权限状态失败', $error->data);
            notice(450000, $error->getMessage());

        } catch (\Exception $e) {

            Log::error('变更渠道物料权限状态失败' . $e->getMessage(), $input);
            notice(450000, $e->getMessage());

        }

        return success();
    }
}