<?php


namespace App\Library\Balance\Wetax;


use Illuminate\Support\Facades\Log;

/**
 * 结算配置
 * Class SettlementConfig
 * @package App\Library\Balance\Wetax
 */
class SettlementConfig
{
    private static $instance;

    /**
     * 基础域名
     * @var string
     */
    private $host;

    /**
     * 商户结算平台id
     * @var string
     */
    private $appKey;

    /**
     * 商户结算平台secret
     * @var string
     */
    private $appSecret;

    /**
     * 结算单回调地址
     * @var string
     */
    private $settlementCallbackUrl;

    /**
     * 结算单退款回调地址
     * @var string
     */
    private $settlementRefundCallbackUrl;

    /**
     * SettlementConfig constructor.
     */
    private function __construct()
    {
        $this->host = env('WETAX_HOST');
        $this->appKey = env('WETAX_APP_KEY', '');
        $this->appSecret = env('WETAX_APP_SECRET', '');
        $this->settlementCallbackUrl = env('WETAX_SETTLEMENT_CALLBACK_URL', '');
        $this->settlementRefundCallbackUrl = env('WETAX_SETTLEMENT_REFUND_CALLBACK_URL', '');
        Log::info("settlement config:", [
            'host' => $this->host,
            'appKey' => $this->appKey,
            'appSecret' => $this->appSecret,
            'settlementCallbackUrl' => $this->settlementCallbackUrl,
            'settlementRefundCallbackUrl' => $this->settlementRefundCallbackUrl
        ]);
        if (empty($this->host) || empty($this->appKey) || empty($this->appSecret)
            || empty($this->settlementCallbackUrl) || empty($this->settlementRefundCallbackUrl)) {
            throw new SettlementException(SettlementException::ILLEGAL_CONFIG);
        }
    }

    /**
     * 初始化入口
     * @return SettlementConfig
     */
    public static function getInstance()
    {
        if (is_null(self::$instance)) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    /**
     * @return string
     */
    public function getHost(): string
    {
        return $this->host;
    }

    /**
     * @return string
     */
    public function getAppKey(): string
    {
        return $this->appKey;
    }

    /**
     * @return string
     */
    public function getAppSecret(): string
    {
        return $this->appSecret;
    }

    /**
     * @return string
     */
    public function getSettlementCallbackUrl(): string
    {
        return $this->settlementCallbackUrl;
    }

    /**
     * @return string
     */
    public function getSettlementRefundCallbackUrl(): string
    {
        return $this->settlementRefundCallbackUrl;
    }
    
    private function __clone()
    {

    }
}