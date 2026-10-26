<?php


namespace App\Library\Balance\Wetax;

/**
 * AES数据加密/解密
 * Class EncryptUtil
 * @package App\Library\Balance\Wetax
 */
class EncryptUtil
{

    /**
     * 加密模式
     */
    const METHOD = "AES-256-CBC";

    const IV = "ff465fdecc764337";

    /**
     * 填充方式：用PKCS#7进行补位
     */
    const OPTION = OPENSSL_RAW_DATA;

    /**
     * 加密：padding->CBC加密->base64编码
     * @param string $context 加密内容
     * @param string $secret 加密密钥
     * @return string
     */
    public static function encrypt(string $context, string $secret): string
    {
        $encrypted = openssl_encrypt($context,self::METHOD, $secret, self::OPTION, self::IV);
        return base64_encode($encrypted);
    }

    /**
     * 解密：base64解码->CBC解密->unpadding
     * @param string $context
     * @param string $secret 加密密钥
     * @return string
     */
    public static function decrypt(string $context, string $secret): string
    {
        $decodeData = base64_decode($context);
        return openssl_decrypt($decodeData,self::METHOD, $secret, self::OPTION, self::IV);
    }
}