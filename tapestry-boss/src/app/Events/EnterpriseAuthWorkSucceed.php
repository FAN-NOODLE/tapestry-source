<?php

/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Events;

use Illuminate\Broadcasting\Channel;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Broadcasting\PrivateChannel;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

class EnterpriseAuthWorkSucceed extends Event
{
    use InteractsWithSockets, SerializesModels;

    public $weChatCorpAuthResponse;
    public $enterpriseId;
    public $suiteId;

    /**
     * Create a new event instance.
     *
     * @param $weChatCorpAuthResponse 企业微信返回的信息
     * @param null $enterpriseId 区分是否已经存在企业的标示
     */
    public function __construct($weChatCorpAuthResponse, $suiteId, $enterpriseId = null)
    {
        $this->weChatCorpAuthResponse = $weChatCorpAuthResponse;
        $this->enterpriseId           = $enterpriseId;
        $this->suiteId                = $suiteId;
    }

    /**
     * Get the channels the event should broadcast on.
     *
     * @return \Illuminate\Broadcasting\Channel|array
     */
    public function broadcastOn()
    {
        return new PrivateChannel('channel-name');
    }
}
