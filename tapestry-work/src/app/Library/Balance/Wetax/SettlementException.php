<?php


namespace App\Library\Balance\Wetax;


use Throwable;

class SettlementException extends \RuntimeException
{

    private $errorCode;

    /**
     * 失败原因
     * @var string
     */
    private $failReason;

    const ILLEGAL_CONFIG = 1;

    const HTTP_REQUEST_FAILURE = 2;

    const CODE_ERROR = 3;

    const DATA_ERROR = 4;

    /**
     * 失败原因列表
     * @var array
     */
    public static $failReasonList = [
        self::ILLEGAL_CONFIG => '结算平台配置错误',
        self::HTTP_REQUEST_FAILURE => 'HTTP请求失败',
        self::CODE_ERROR => 'code参数错误',
        self::DATA_ERROR => 'data内容错误'
    ];

    /**
     * SettlementException constructor.
     * @param int $errorCode
     * @param string $failReason
     * @param Throwable|null $previous
     */
    public function __construct(int $errorCode, string $failReason = '', Throwable $previous = null)
    {
        $this->errorCode = $errorCode;
        $this->failReason = $failReason != '' ? $failReason : self::$failReasonList[$errorCode] ?? '未知原因';
        parent::__construct($this->failReason, $this->errorCode, $previous);
    }

    /**
     * @return int
     */
    public function getErrorCode(): int
    {
        return $this->errorCode;
    }

    /**
     * @return string
     */
    public function getFailReason(): string
    {
        return $this->failReason;
    }



}