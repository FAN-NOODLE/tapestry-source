<?php

namespace App\Jobs;

use Illuminate\Support\Facades\Log;
use App\Modules\Message\MessageFactory;

class MessageSendJob extends Job
{
    public $queue = 'message-send';

    private $channel;
    private $enterpriseId;
    private $userIds;
    private $sceneType;
    private $message;
    private $input;

    public function __construct($channel, $enterpriseId, $userIds, $sceneType, $message, $input)
    {
        $this->channel = $channel;
        $this->enterpriseId = $enterpriseId;
        $this->userIds = $userIds;
        $this->sceneType = $sceneType;
        $this->message = $message;
        $this->input = $input;
    }

    public function handle()
    {
        Log::info('消息发送开始消费');

        try {
            $message = MessageFactory::getInstance($this->channel, $this->enterpriseId, $this->userIds, $this->sceneType, $this->message, $this->input);
            return $message->send();
        } catch (\Exception $exception) {
            Log::error('发送消息异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'input'   => $this->input
            ]);
        }

        Log::info('消息发送开始消费完成');
    }
}
