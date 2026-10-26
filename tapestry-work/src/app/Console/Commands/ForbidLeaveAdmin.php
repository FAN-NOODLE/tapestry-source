<?php


namespace App\Console\Commands;



use App\Micros\Common;
use Illuminate\Console\Command;

class ForbidLeaveAdmin extends Command
{
    protected $signature = 'forbid-leave-admin';

    protected $description = '禁用离职管理员';

    public function handle()
    {
        Common::authority()
            ->setApi("micro-authority/v1/account/forbid-leave-admin")
            ->post();
    }


}