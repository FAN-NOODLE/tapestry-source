<?php


namespace App\Library\Message\Response;


use App\Library\Message\Response;

class MWResponse extends Response
{
    private $respStatus = 0;
    private $msgId = '';
    private $code = [
        0       => '发送成功',
        -100001 => '鉴权不通过,请检查账号,密码,时间戳,固定串,以及MD5算法是否按照文档要求进行设置',
        -100002 => '用户多次鉴权不通过,请检查帐号,密码,时间戳,固定串,以及MD5算法是否按照文档要求进行设置',
        -100003 => '用户欠费',
        -100004 => 'custid或者exdata字段填写不合法',
        -100011 => '短信内容超长',
        -100012 => '手机号码不合法',
        -100014 => '手机号码超过最大支持数量（1000）',
        -100029 => '端口绑定失败',
        -100056 => '用户账号登录的连接数超限',
        -100057 => '用户账号登录的IP错误',
        -100999 => '平台数据库内部错误',
        -100500 => '富信接口发非富信',
        -100501 => '富信接口内容过长',
        -100502 => '短信模板被禁用',
        99999     => '未知错误'
    ];

    public function __construct($response)
    {
        parent::__construct($response);

        $this->parseResponse();
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
        if (isset($this->response['result'])) {
            $this->respStatus = $this->response['result'];
        }
        if (isset($this->response['msgid'])) {
            $this->msgId = $this->response['msgid'];
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
        if (isset($this->code[$this->getStateCode()])) {
            return $this->code[$this->getStateCode()];
        }

        return $this->code[99999];
    }
}