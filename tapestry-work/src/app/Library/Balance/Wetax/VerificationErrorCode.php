<?php


namespace App\Library\Balance\Wetax;


class VerificationErrorCode
{

    /**
     * 结算单审核失败错误码
     */
    const ERROR_CODE_LIST = [
        11 => '姓名不能为空',
        12 => '姓名格式不正确',
        13 => '证件类型不能为空',
        14 => '证件类型格式不正确',
        15 => '证件号不能为空',
        16 => '证件号格式不正确',
        17 => '姓名与身份证不匹配',
        18 => '该人员已被加入黑名单',
        19 => '未进行个体户登记注册',
        21 => '开户行名称不能为空',
        22 => '银行卡不能为空',
        23 => '银行卡格式不正确',
        24 => '非本人银行卡',
        25 => '收款账户不能为空',
        26 => '收款账户格式不正确',
        27 => '银行卡类型必须是借记卡',
        31 => '手机号不能为空',
        32 => '手机号格式不正确',
        33 => '非本人手机号',
        41 => '绩效费不能为空',
        42 => '绩效费格式不正确',
        43 => '绩效费超过限制',
        50 => '疑似董监高',
        51 => '疑似内部员工',
        52 => '结算用户年龄低于16岁',
        99 => '客户订单号已存在',
    ];

    public static function getErrorByCodeList(array $errorCodeList)
    {
        $msgList = [];
        foreach ($errorCodeList as $errorCode) {
            array_push($msgList, self::ERROR_CODE_LIST[$errorCode] ?? "未知错误-" . $errorCode);
        }
        return implode(",", $msgList);
    }
}