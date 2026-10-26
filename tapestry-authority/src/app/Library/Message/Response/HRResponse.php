<?php
/**
 * Created By PhpStorm
 * @Date: 2021/4/14
 * @Time: 10:31
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: Response.php
 */

namespace App\Library\Message\Response;
use App\Library\Message\Response;

class HRResponse extends Response
{
    private $respTime = '';
    private $respStatus = '';
    private $msgId = '';
    private $code = [
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
         parent::__construct($response);

         $this->parseResponse();
    }

    /**
     * 获取响应时间
     */
    public function getRespTime()
    {
        return $this->respTime;
    }

    /**
     * 获取响应状态
     */
    public function getRespStatus()
    {
        return $this->respStatus;
    }

    /**
     * 获取msgid
     */
    public function getMsgId()
    {
        return $this->msgId;
    }

    /**
     * 解析响应数据
     */
    public function parseResponse()
    {
        if(is_string($this->response)){
            if(count(explode(',', $this->response)) == 2){
                list($respTime, $respStatus) = explode(',', $this->response, 2);
            }else{
                list($respTime, $respStatus, $msgId) = explode(',', $this->response, 3);
            }
        }

        if(isset($respTime)) {
            $this->respTime = $respTime;
        }

        if(isset($respStatus)) {
            $this->respStatus = $respStatus;
        }

        if(isset($msgId)) {
            $this->msgId = $msgId;
        }
    }

    /**
     * 获取返回状态
     */
    public function getStateCode()
    {
        return $this->getRespStatus();
    }


    /**
     * 获取返回状态的提示语
     */
    public function getCodeMsg()
    {
        if(isset($this->code[$this->getStateCode()])){
            return $this->code[$this->getStateCode()];
        }

        return $this->code[99999];
    }
}