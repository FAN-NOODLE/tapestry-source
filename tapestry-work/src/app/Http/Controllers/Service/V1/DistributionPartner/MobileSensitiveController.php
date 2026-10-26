<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/12/31 16:07
 */

namespace App\Http\Controllers\Service\V1\DistributionPartner;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseMap;
use Illuminate\Http\Request;

class MobileSensitiveController extends BaseController
{
    /**
     * @param Request $request
     * @return array
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/6 17:07
     */
    const SENSITIVE = 1; //脱敏显示
    public function index(Request $request)
    {
        $params = $this->getParams($request);
        $partner_mobile = $params['mobile'];
        if ($params['mobile_sensitive'] == self::SENSITIVE) {
            $partner_mobile = $this->dataDesensitization($params['mobile'], 3, 4);
        }
        
        PartnerEnterpriseMap::query()
            ->where('enterprise_id',$params['enterprise_id'])
            ->where('partner_id',$params['partner_id'])
            ->update([
                'mobile_sensitive'=>$params['mobile_sensitive'],
                "mobile_show"=>$partner_mobile
            ]);
        return success();
    }

    private function getParams(Request $request): array
    {
        $rules = [
            'mobile_sensitive' => 'required|in:0,1',
            'partner_id' => 'required',
            'enterprise_id' => 'required',
            'mobile'=>'required'
        ];
        $message = [
            'mobile_sensitive.*' => '手机号码保护不合法',
            'enterprise_id.*' => '企业ID不合法',
            'partner_id.*'=>'合伙人ID不合法',
            'mobile.*'=>'手机号不合法'
        ];
        $input = $this->validate($request, $rules, $message);
        return $input;
    }

    function dataDesensitization($string, $start = 0, $length = 0, $re = '*')
    {
        if (empty($string)) {
            return false;
        }
        $strarr = array();
        $mb_strlen = mb_strlen($string);
        while ($mb_strlen) {//循环把字符串变为数组
            $strarr[] = mb_substr($string, 0, 1, 'utf8');
            $string = mb_substr($string, 1, $mb_strlen, 'utf8');
            $mb_strlen = mb_strlen($string);
        }
        $strlen = count($strarr);
        $begin = $start >= 0 ? $start : ($strlen - abs($start));
        $end = $last = $strlen - 1;
        if ($length > 0) {
            $end = $begin + $length - 1;
        } elseif ($length < 0) {
            $end -= abs($length);
        }
        for ($i = $begin; $i <= $end; $i++) {
            $strarr[$i] = $re;
        }
        if ($begin >= $end || $begin >= $last || $end > $last) return false;
        return implode('', $strarr);
    }
}