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

use Overtrue\LaravelWeChat\Events\OpenPlatform\Authorized as OpenPlatformAuthorized;
use Illuminate\Support\Facades\Log;

class Authorized
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
     * @param OpenPlatformAuthorized $event
     */
    public function handle(OpenPlatformAuthorized $event)
    {
        $payload = $event->payload;
        $monolog = Log::getMonolog();
        $monolog->popHandler();
        Log::useDailyFiles(storage_path('logs/payload.log'));

        Log::info('授权: ' . var_export($payload, true));
        // logger('auth', ['auth' => $event]);
    }
}
