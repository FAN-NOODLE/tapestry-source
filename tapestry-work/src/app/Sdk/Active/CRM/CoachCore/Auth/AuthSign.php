<?php

namespace App\Sdk\Active\CRM\CoachCore\Auth;

class AuthSign
{

    /**
     * 生成签名 sign
     *
     * @param $params
     * @param $appKey
     *
     * @return string
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/20
     */
    public static function makeSign($params, $appKey): string
    {
        // 签名步骤一：按字典升序排序参数
        ksort($params);
        $string = self::splicingString($params);
        // 签名步骤二：在string后加入KEY
        $string = $string . '&' . $appKey;
        // 签名步骤三：MD5加密
        $sign = md5($string);
        // 签名步骤四：所有字符转为大写
        return strtoupper($sign);
    }

    /**
     * 拼接
     *
     * @param $params
     *
     * @return string
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/20
     */
    public static function splicingString($params): string
    {
        $string = '';
        foreach ($params as $k => $v) {
            if ($k != "sign" && $v != "" && !is_array($v)) {
                $string .= ($v) . '&';
            }
        }
        return trim($string, "&");
    }

}