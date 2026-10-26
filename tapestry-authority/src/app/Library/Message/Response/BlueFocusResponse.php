<?php


namespace App\Library\Message\Response;


use App\Library\Message\Response;

class BlueFocusResponse extends Response
{
    private $respStatus = 0;
    private $msgId = '';
    private $code = [
        0       => '提交成功',
        101 => '系统忙（因平台侧原因，暂时无法处理提交的短信）',
        102 => '客户IP地址认证错,请求调用的IP地址不在客户绑定的IP',
        201 => '无此客户',
        202 => '密码错',
        301 => '消息长度错（>500或<=0）',
        302 => '包含错误的手机号码',
        303 => '手机号码个数超限（群发>10000或<=0）',
        304 => 'extno格式错（非数字或者长度不对）',
        305 => '手机号码个数和短信内容个数不一致',
        306 => '消息内容为空',
        307 => '手机号码为空',
        308 => '内容编码错',
        99999 => '未知错误'
    ];

    public function __construct($response)
    {
        $this->response = explode(",", $response);
    }


    /**
     * 获取返回状态
     */
    public function getStateCode()
    {
        return $this->response[1];
    }

    /**
     * 获取返回状态的提示语
     */
    public function getCodeMsg() : string
    {
        if(isset($this->code[$this->getStateCode()])){
            return $this->code[$this->getStateCode()];
        }

        return $this->code[99999];
    }

}