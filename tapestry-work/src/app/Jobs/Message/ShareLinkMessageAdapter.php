<?php


namespace App\Jobs\Message;


use App\Common\MessageConstant;

class ShareLinkMessageAdapter
{

    public static function createMessage(array $afterBinding)
    {
        $message = new ShareLinkMessage();
        $message->setType(MessageConstant::SHARE_LINK_CREATE);
        $message->setUserId($afterBinding['user_id']);
        $message->setEnterpriseId($afterBinding['enterprise_id']);
        $message->setCurrWorkerId($afterBinding['worker_id']);
        $message->setCurrWorkerMemUid($afterBinding['mem_uid']);
        $message->setCreatedAt($afterBinding['created_at']);
        $message->setUpdatedAt($afterBinding['created_at']);
        return $message;
    }

    public static function deleteMessage(array $beforeBinding, $deletedAt)
    {
        $message = new ShareLinkMessage();
        $message->setType(MessageConstant::SHARE_LINK_DELETE);
        $message->setUserId($beforeBinding['user_id']);
        $message->setEnterpriseId($beforeBinding['enterprise_id']);
        $message->setOldWorkerId($beforeBinding['worker_id']);
        $message->setOldWorkerMemUid($beforeBinding['mem_uid']);
        $message->setCreatedAt($beforeBinding['created_at']);
        $message->setUpdatedAt($deletedAt);
        return $message;
    }

    public static function updateMessage(array $beforeBinding, array $afterBinding)
    {
        $message = new ShareLinkMessage();
        $message->setType(MessageConstant::SHARE_LINK_UPDATE);
        $message->setUserId($afterBinding['user_id']);
        $message->setEnterpriseId($afterBinding['enterprise_id']);
        $message->setCurrWorkerId($afterBinding['worker_id']);
        $message->setCurrWorkerMemUid($afterBinding['mem_uid']);
        $message->setOldWorkerId($beforeBinding['worker_id']);
        $message->setOldWorkerMemUid($beforeBinding['mem_uid']);
        $message->setCreatedAt($afterBinding['created_at']);
        $message->setUpdatedAt($afterBinding['created_at']);
        return $message;
    }
}