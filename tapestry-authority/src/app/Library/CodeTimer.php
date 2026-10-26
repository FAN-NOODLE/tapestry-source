<?php

namespace App\Library;

/**
 * Class 代码计时器
 * 0. $ct = new CodeTimer();
 * 1. $ct->open();
 * 2. $ct->close();
 * 3. echo $ct->getRunTime();// 单位毫秒
 *
 * @package App\Libs
 */
class CodeTimer
{
    private $openTime;
    private $closeTime;

    public function open()
    {
        $this->openTime = microtime(true);
        $this->closeTime = null;
    }

    public function close()
    {
        if ($this->openTime === null) {
            return false;
        }
        $this->closeTime = microtime(true);

        return true;
    }

    public function getRunTime()
    {
        if ($this->openTime === null || $this->closeTime === null) {
            return null;
        }

        return round(($this->closeTime - $this->openTime) * 1000, 2);
    }
}
