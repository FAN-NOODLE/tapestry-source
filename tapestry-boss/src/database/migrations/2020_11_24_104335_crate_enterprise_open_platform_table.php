<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CrateEnterpriseOpenPlatformTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_open_platform_config')) {
            Schema::create('enterprise_open_platform_config', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->string('app_id', 150)->default('')->comment('开放平台app_id');
                $table->string('secret', 150)->default('')->comment('开放平台secret');
                $table->string('token', 150)->default('')->comment('开放平台token');
                $table->string('aes_key', 150)->default('')->comment('开放平台aes_key');
                $table->unsignedInteger('updated_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->unsignedInteger('deleted_at')->default(0);
            });

            DB::statement("ALTER TABLE `enterprise_open_platform_config` COMMENT '企业开放平台配置'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_open_platform_config');
    }
}
