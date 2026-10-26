<?php

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Queue\SerializesModels;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Support\Facades\Config;
use Illuminate\Support\Facades\DB;

abstract class Job implements ShouldQueue
{
    /*
    |--------------------------------------------------------------------------
    | Queueable Jobs
    |--------------------------------------------------------------------------
    |
    | This job base class provides a central location to place any logic that
    | is shared across all of your jobs. The trait included with the class
    | provides access to the "queueOn" and "delay" queue helper methods.
    |
    */

    use InteractsWithQueue, Queueable, SerializesModels;

    protected $databaseName = null;

    protected $enterpriseHash = null;


    public function __construct($hash = '')
    {
        if (!empty($hash)) {
            $this->databaseName = config('database.database_hash.' . $hash, '');
        } elseif (!empty($_SERVER['enterprise_hash'])) {
            $hash               = $_SERVER['enterprise_hash'];
            $this->databaseName = config('database.database_hash.' . $hash, '');
        }

        $this->enterpriseHash = $hash;
    }

    public function resetDb()
    {
        if (!empty($this->databaseName)) {

            Config::set('database.connections.mysql.database', $this->databaseName);
        } else {
            $databaseName = env('DB_DATABASE');
            Config::set('database.connections.mysql.database', $databaseName);
        }

        DB::reconnect();
    }
}
