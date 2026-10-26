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

use App\Services\OpenPlatform\EnterpriseService;
use Overtrue\LaravelWeChat\Events\OpenPlatform\Unauthorized as OpenPlatformUnauthorized;
use Illuminate\Support\Facades\Log;

class Unauthorized
{
    /**
     * Create the event listener.
     *
     */
    public function __construct()
    {
    }

    /**
     * Handle the event.
     *
     * @param OpenPlatformUnauthorized $event
     * @param EnterpriseService        $enterpriseService
     */
    public function handle(OpenPlatformUnauthorized $event, EnterpriseService $enterpriseService)
    {
        // 取消公众号授权授权
        $payload = $event->payload;

        $enterpriseService->negativeThirdPartStatus($payload['AuthorizerAppid']);
        // logger('unauth', ['unauth' => $event]);
    }
}
