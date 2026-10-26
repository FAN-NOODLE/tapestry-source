<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-14
 * Time: 17:06
 */

namespace App\Library\Message\Handlers;

use App\Library\Message\Contracts\SendHandlerInterface;

class MiniSubscriptionHandlers implements SendHandlerInterface
{
    public function checkRequestParams(array $params): array
    {
        // TODO: Implement checkRequestParams() method.
    }

    public function send($subject, array $params, string $templateCode = ''): bool
    {
        // TODO: Implement send() method.
    }

    public function getSendInstant(array $config = [])
    {

    }
}
