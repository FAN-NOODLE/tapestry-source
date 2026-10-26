<?php


namespace App\Jobs\Message;


class ShareLinkMessage extends Message
{
    /**
     * @var int
     */
    private $type;

    /**
     * @var int
     */
    private $enterpriseId;

    /**
     * @var int
     */
    private $userId;

    /**
     * @var int|null
     */
    private $currWorkerId;

    /**
     * @var string|null
     */
    private $currWorkerMemUid;

    /**
     * @var int|null
     */
    private $oldWorkerId;

    /**
     * @var string|null
     */
    private $oldWorkerMemUid;

    /**
     * @var int
     */
    private $createdAt;

    /**
     * @var int
     */
    private $updatedAt;

    /**
     * @return int
     */
    public function getType(): int
    {
        return $this->type;
    }

    /**
     * @param int $type
     */
    public function setType(int $type): void
    {
        $this->type = $type;
    }

    /**
     * @return int
     */
    public function getEnterpriseId(): int
    {
        return $this->enterpriseId;
    }

    /**
     * @param int $enterpriseId
     */
    public function setEnterpriseId(int $enterpriseId): void
    {
        $this->enterpriseId = $enterpriseId;
    }

    /**
     * @return int
     */
    public function getUserId(): int
    {
        return $this->userId;
    }

    /**
     * @param int $userId
     */
    public function setUserId(int $userId): void
    {
        $this->userId = $userId;
    }

    /**
     * @return int
     */
    public function getCurrWorkerId(): ?int
    {
        return $this->currWorkerId;
    }

    /**
     * @param int $currWorkerId
     */
    public function setCurrWorkerId(?int $currWorkerId): void
    {
        $this->currWorkerId = $currWorkerId;
    }

    /**
     * @return string
     */
    public function getCurrWorkerMemUid(): ?string
    {
        return $this->currWorkerMemUid;
    }

    /**
     * @param string $currWorkerMemUid
     */
    public function setCurrWorkerMemUid(?string $currWorkerMemUid): void
    {
        $this->currWorkerMemUid = $currWorkerMemUid;
    }

    /**
     * @return int|null
     */
    public function getOldWorkerId(): ?int
    {
        return $this->oldWorkerId;
    }

    /**
     * @param int|null $oldWorkerId
     */
    public function setOldWorkerId(?int $oldWorkerId): void
    {
        $this->oldWorkerId = $oldWorkerId;
    }

    /**
     * @return string|null
     */
    public function getOldWorkerMemUid(): ?string
    {
        return $this->oldWorkerMemUid;
    }

    /**
     * @param string|null $oldWorkerMemUid
     */
    public function setOldWorkerMemUid(?string $oldWorkerMemUid): void
    {
        $this->oldWorkerMemUid = $oldWorkerMemUid;
    }

    /**
     * @return int
     */
    public function getCreatedAt(): int
    {
        return $this->createdAt;
    }

    /**
     * @param int $createdAt
     */
    public function setCreatedAt(int $createdAt): void
    {
        $this->createdAt = $createdAt;
    }

    /**
     * @return int
     */
    public function getUpdatedAt(): int
    {
        return $this->updatedAt;
    }

    /**
     * @param int $updatedAt
     */
    public function setUpdatedAt(int $updatedAt): void
    {
        $this->updatedAt = $updatedAt;
    }



}