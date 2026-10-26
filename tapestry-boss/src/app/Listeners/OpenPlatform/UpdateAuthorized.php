<?php
/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Listeners\OpenPlatform;

use Overtrue\LaravelWeChat\Events\OpenPlatform\UpdateAuthorized as OpenPlatformUpdateAuthorized;

class UpdateAuthorized
{
    /**
     * Create the event listener.
     */
    public function __construct()
    {
    }

    /**
     * Handle the event.
     *
     * @param OpenPlatformUpdateAuthorized $event
     */
    public function handle(OpenPlatformUpdateAuthorized $event)
    {
        // logger('reauth', ['reauth' => $event]);
    }
}
