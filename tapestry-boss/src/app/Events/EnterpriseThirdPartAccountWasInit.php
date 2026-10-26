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

use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Broadcasting\PrivateChannel;
use Illuminate\Queue\SerializesModels;

class EnterpriseThirdPartAccountWasInit extends Event
{
    use InteractsWithSockets, SerializesModels;

    /**
     * @var EnterpriseThirdPartAccount
     */
    public $thirdPartAccount;

    /**
     * Create a new event instance.
     */
    public function __construct(EnterpriseThirdPartAccount $thirdPartAccount)
    {
        $this->thirdPartAccount = $thirdPartAccount;
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
