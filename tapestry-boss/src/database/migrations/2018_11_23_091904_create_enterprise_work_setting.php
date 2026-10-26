<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateEnterpriseWorkSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_work_setting')) {
            Schema::create('enterprise_work_setting', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('external_secret',100)->default('')->comment('外部联系人secret');
                $table->string('mini_program_secret',100)->default('')->comment('企业微信小程序secret');
                $table->string('corp_id',100)->default('')->comment('企业微信企业id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_work_setting` comment '企业微信配置信息表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('enterprise_work_setting');
    }
}
