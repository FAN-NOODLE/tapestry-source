<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseSdkConfigTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_sdk_config')) {

            Schema::create('enterprise_sdk_config', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('app_type')->default(1)->comment('小程序类型：1-商城 2-企业微信');
                $table->string('sdk_name', 30)->default('')->comment('sdk名称');
                $table->string('request_domain', 255)->default('')->comment('sdk域名');
                $table->text('pack_config')->comment('自动打包sdk-配置');
                $table->text('ext_config')->comment('自动打包sdk-扩展配置');
                $table->unsignedInteger('updated_at')->default(0);
                $table->unsignedInteger('created_at')->default(0);
                $table->unsignedInteger('deleted_at')->default(0);
            });
        }

        DB::statement("ALTER TABLE `enterprise_sdk_config` COMMENT '小程序数据收集sdk配置'");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_sdk_config');
    }
}
