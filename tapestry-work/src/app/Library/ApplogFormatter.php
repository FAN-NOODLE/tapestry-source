<?php
/**
 * Created by PhpStorm.
 * User: guoxingong
 * Date: 2019/4/18
 * Time: 10:29
 */

namespace App\Library;

use App\Library\JsonFormatter;

class ApplogFormatter
{
    /**
     * Customize the given logger instance.
     *
     * @param  \Illuminate\Log\Logger  $logger
     * @return void
     */
    public function __invoke($logger)
    {
        foreach ($logger->getHandlers() as $handler) {
            $handler->setFormatter(new JsonFormatter());
        }
    }
}