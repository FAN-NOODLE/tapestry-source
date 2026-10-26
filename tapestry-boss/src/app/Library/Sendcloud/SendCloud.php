<?php


namespace App\Library\Sendcloud;



class SendCloud extends Abstracts
{

    /**
     * 发送模板邮件地址
     *
     * @var string
     */
    const SEND_TEMPLATE_MAIL = "%s/apiv2/mail/sendtemplate";

    /**
     * 发送模板邮件
     *
     * @param array $params
     * @return mixed
     */
    public function sendTemplateEmail(array $params)
    {
        $url = sprintf(self::SEND_TEMPLATE_MAIL, $this->baseUri);
        return $this->post($url, $params);
    }
}