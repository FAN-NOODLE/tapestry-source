<?php

namespace App\Modules\Message;

class MessageFactory
{
    /**
     * CDP
     */
    const CDP = 'cdp';

    /**
     * 小程序
     */
    const MINI_PROGRAM = 'mini_program';

    /**
     * 获取Message实例
     *
     * @param string $channel 渠道标识
     * @param int $enterpriseId 企业ID
     * @param array $userIds 用户ID列表
     * @param string $sceneType 场景值
     * @param array $message 消息内容
     * @param array $option 额外参数
     * @return CdpMessage|MiniProgramMessage
     * @throws \Exception
     */
    public static function getInstance($channel, $enterpriseId, $userIds, $sceneType, $message, $option)
    {
        $class = null;
        switch ($channel) {
            case self::CDP:
                $class = new CdpMessage();
                break;
            case self::MINI_PROGRAM:
                $class = new MiniProgramMessage();
                break;
            default:
                throw new \Exception("未定义的平台标识: {$channel}");
        }

        $class->setChannel($channel);
        $class->setEnterpriseId($enterpriseId);
        $class->setSceneType($sceneType);
        $class->setMessage($message);
        $class->setOption($option);
        $class->setRecords($userIds);
        $class->setUserIds($userIds);

        return $class;
    }
}
