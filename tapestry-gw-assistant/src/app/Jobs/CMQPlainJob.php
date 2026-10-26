<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-15
 * Time: 12:01
 */

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Freyo\LaravelQueueCMQ\Queue\Contracts\PlainPayload;

class CMQPlainJob implements ShouldQueue, PlainPayload
{
    use Dispatchable, InteractsWithQueue, Queueable;

    protected $payload;

    /**
     * Create a new job instance.
     *
     * @return void
     */
    public function __construct($payload)
    {
        $this->payload = $payload;
    }

    /**
     * Get the plain payload of the job.
     *
     * @return string
     */
    public function getPayload()
    {
        return $this->payload;
    }
}
