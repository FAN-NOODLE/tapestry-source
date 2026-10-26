<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateEnterpriseVirtualOpenPlatform extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_virtual_open_platform')) {
            Schema::create('enterprise_virtual_open_platform', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('open_app_id',191)->default('')->comment('虚拟开放平台id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_virtual_open_platform` comment '虚拟开放平台信息表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('enterprise_virtual_open_platform');
    }
}
