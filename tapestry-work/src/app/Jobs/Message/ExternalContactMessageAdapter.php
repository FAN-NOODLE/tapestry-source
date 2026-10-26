<?php


namespace App\Jobs\Message;


use App\Common\MessageConstant;

class ExternalContactMessageAdapter
{

    public static function createMessage(array $memberUserRelation)
    {
        $message = (new ExternalContactMessage())
            ->setType(MessageConstant::EXTERNAL_CONTACT_CREATE)
            ->setEnterpriseId($memberUserRelation['enterprise_id'])
            ->setMemUid($memberUserRelation['mem_uid'])
            ->setUnionId($memberUserRelation['union_id'])
            ->setSourceType($memberUserRelation['source_type'])
            ->setCreatedAt($memberUserRelation['created_at'])
            ->setUnionCode(isset($memberUserRelation['union_code']) ? $memberUserRelation['union_code'] : '')
            ->setDeletedAt(0);
        return $message;
    }

    public static function deleteMessage(array $memberUserRelation)
    {
        $message = (new ExternalContactMessage())
            ->setType(MessageConstant::EXTERNAL_CONTACT_DELETE)
            ->setEnterpriseId($memberUserRelation['enterprise_id'])
            ->setMemUid($memberUserRelation['mem_uid'])
            ->setUnionId($memberUserRelation['union_id'])
            ->setSourceType($memberUserRelation['source_type'])
            ->setCreatedAt($memberUserRelation['created_at'])
            ->setUnionCode(isset($memberUserRelation['union_code']) ? $memberUserRelation['union_code'] : '');

        if ($memberUserRelation['deleted_at'] > 0) {
            $message->setIsUnilateral(false);
            $message->setDeletedAt($memberUserRelation['deleted_at']);
        } elseif ($memberUserRelation['delete_type'] > 0) {
            $message->setIsUnilateral(true);
            $message->setDeletedAt($memberUserRelation['cancel_relation_at']);
        }
        return $message;
    }
}