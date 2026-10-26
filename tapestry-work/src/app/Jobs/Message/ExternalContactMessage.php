<?php


namespace App\Jobs\Message;

/**
 * 外部联系人变更消息
 * Class ExternalContactMessage
 * @package App\Jobs\Message
 */
class ExternalContactMessage extends Message
{
    /**
     * 消息类型
     * @var int
     */
    private $type;

    /**
     * @var int
     */
    private $enterpriseId;

    /**
     * 导购memUid
     * @var string
     */
    private $memUid;

    /**
     * 用户unionId
     * @var string
     */
    private $unionId;

    /**
     * 好友添加方式
     * @var int
     */
    private $sourceType;

    /**
     * 好友数据添加时间
     * @var int
     */
    private $createdAt;

    /**
     * 删除时间
     * @var int
     */
    private $deletedAt;

    /**
     * 是否单方面删除
     * @var bool
     */
    private $isUnilateral;

    /**
     * 会员卡号
     * @var int
     */
    private $unionCode;

    /**
     * ExternalContactMessage constructor.
     */
    public function __construct()
    {

    }

    /**
     * @return int
     */
    public function getType(): int
    {
        return $this->type;
    }

    /**
     * @param int $type
     * @return ExternalContactMessage
     */
    public function setType(int $type): ExternalContactMessage
    {
        $this->type = $type;
        return $this;
    }

    /**
     * @return int
     */
    public function getEnterpriseId(): int
    {
        return $this->enterpriseId;
    }

    /**
     * @param string $enterpriseId
     * @return ExternalContactMessage
     */
    public function setEnterpriseId(string $enterpriseId): ExternalContactMessage
    {
        $this->enterpriseId = $enterpriseId;
        return $this;
    }

    /**
     * @return string
     */
    public function getMemUid(): string
    {
        return $this->memUid;
    }

    /**
     * @param string $memUid
     * @return ExternalContactMessage
     */
    public function setMemUid(string $memUid): ExternalContactMessage
    {
        $this->memUid = $memUid;
        return $this;
    }

    /**
     * @return string
     */
    public function getUnionId(): string
    {
        return $this->unionId;
    }

    /**
     * @param string $unionId
     * @return ExternalContactMessage
     */
    public function setUnionId(string $unionId): ExternalContactMessage
    {
        $this->unionId = $unionId;
        return $this;
    }

    /**
     * @return int
     */
    public function getSourceType(): int
    {
        return $this->sourceType;
    }

    /**
     * @param int $sourceType
     * @return ExternalContactMessage
     */
    public function setSourceType(int $sourceType): ExternalContactMessage
    {
        $this->sourceType = $sourceType;
        return $this;
    }

    /**
     * @return int
     */
    public function getCreatedAt(): ?int
    {
        return $this->createdAt;
    }

    /**
     * @param int $createdAt
     * @return ExternalContactMessage
     */
    public function setCreatedAt(?int $createdAt): ExternalContactMessage
    {
        $this->createdAt = $createdAt;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getDeletedAt(): ?int
    {
        return $this->deletedAt;
    }

    /**
     * @param int|null $deletedAt
     * @return ExternalContactMessage
     */
    public function setDeletedAt(?int $deletedAt): ExternalContactMessage
    {
        $this->deletedAt = $deletedAt;
        return $this;
    }

    /**
     * @return bool|null
     */
    public function isUnilateral(): ? bool
    {
        return $this->isUnilateral;
    }

    /**
     * @param bool|null $isUnilateral
     * @return ExternalContactMessage
     */
    public function setIsUnilateral(?bool $isUnilateral): ExternalContactMessage
    {
        $this->isUnilateral = $isUnilateral;
        return $this;
    }

    /**
     * @return string
     */
    public function getUnionCode(): ?string
    {
        return $this->unionCode;
    }

    /**
     * @param string $unionCode
     * @return ExternalContactMessage
     */
    public function setUnionCode(?string $unionCode): ExternalContactMessage
    {
        $this->unionCode = $unionCode;
        return $this;
    }
}