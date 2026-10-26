<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateEnterpriseWorkPluginTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_work_plugin')) {
            Schema::create('enterprise_work_plugin', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('corp_id', 255)->default('')->comment('企业corpid');
                $table->string('agent_id', 255)->default('')->comment('应用id');
                $table->string('secret', 255)->default('')->comment('应用secret');
                $table->string('identifier', 255)->default('')->comment('应用标识');
                $table->tinyInteger('check_file')->default(2)->comment('域名校验文件是否上传 1 是 2否');
                $table->unsignedInteger('updated_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->unsignedInteger('deleted_at')->default(0);
            });

            DB::statement("ALTER TABLE `enterprise_work_plugin` COMMENT '企业微信自建应用表'");
        }

    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_work_plugin');
    }
}
