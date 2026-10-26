<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/26 17:16
 */

namespace App\Http\Controllers\V1\MarketingPlan;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\AppCode;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class GetSunCodeController extends BaseController
{
    public function index(Request $request)
    {
        $rules = [
            'applet_path' => 'required'
            ];
        $message = [
            'plan_id.*' => '小程序路径不能为空'
        ];
        $this->validate($request, $rules, $message);

        $enterprise_Id = $this->enterpriseId;
        $applet_path = trim($request->input('applet_path'));
        [$path,$params] = $this->convertUrlQuery($applet_path);
        $sunCode = '';

        if(isset($params['scene']) && !empty($params['scene'])){
            $decryptData = decryptShortLink($params['scene']);
            $params       = [];
            foreach ($decryptData['args'] ?? [] as $key => $value) {
                $params[$value['name']] = $value['value'];
            }
        }

        try{
            $sunCode = (new AppCode())->create(['enterprise_id' => $enterprise_Id, 'params' => $params, 'path' => $path])['data']['code_url'] ?? '';
        } catch (\Exception $exception) {
            Log::info('生成营销预览码失败：' . $exception->getFile() . ',行数：' . $exception->getLine() . ',错误原因：' . $exception->getMessage() .'.getTrace'.json_encode($exception->getTrace()));
            notice(2200002);
        }

        if(empty($sunCode)){
            notice(2200001);
        }
        return success(['url'=>$sunCode]);
    }

    /**
     * 小程序路径解参
     * @param $applet_path
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/13 15:03
     */
    public function convertUrlQuery($applet_path)
    {
        $path = '';
        $params = [];
        if(empty($applet_path)){
            return [$path,$params];
        }

        $tmp = explode("?",$applet_path);
        $path = $tmp[0];
        $query = $tmp[1] ?? '';
        if(!empty($query)){
            parse_str($query,$params);
        }
        return [$path,$params];
    }
}