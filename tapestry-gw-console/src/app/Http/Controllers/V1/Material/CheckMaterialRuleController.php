<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/24
 * Time: 20:37
 */
namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Material;
use Illuminate\Support\Facades\Log;

class CheckMaterialRuleController extends BaseController
{
    /**
     * 验证门店的物料权限
     *
     * @return array
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        try {
            $input['source_id'] = $input['source_id'] ?? $this->enterpriseId;
            $result = (new Material())->checkMaterialRule($input);

        } catch (Error $error) {

            Log::error('验证门店的物料权限失败', $error->data);
            notice(450000, $error->getMessage());

        } catch (\Exception $e) {

            Log::error('验证门店的物料权限失败' . $e->getMessage(), $input);
            notice(450000, $e->getMessage());

        }

        return success($result);
    }
}