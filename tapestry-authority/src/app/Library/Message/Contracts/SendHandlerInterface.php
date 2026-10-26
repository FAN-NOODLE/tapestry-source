<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-14
 * Time: 16:56
 */

namespace App\Library\Message\Contracts;

interface SendHandlerInterface
{
    /**
     * 发送
     *
     * @param array $params
     *             + subject string 发送对象
     *             + templateCode string 发送模板ID
     *             + data array 消息数据
     * @return bool
     */
    public function send(array $params) : bool ;

    /**
     * 获取当前调用接口的账号
     * @return string
     */
    public function getAccount() : string;
}
