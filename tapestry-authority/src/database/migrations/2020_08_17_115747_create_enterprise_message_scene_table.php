<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseMessageSceneTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_message_scene')) {
            Schema::create('enterprise_message_scene', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary();
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id');
                $table->unsignedBigInteger('message_scene_id')->default(0)->comment('场景id');
                $table->unsignedInteger('scene_type')->default(0)->comment('小程序订阅消息场景id');
                $table->string('config')->default('')->comment('配置项');
                $table->unsignedTinyInteger('status')->default(1)->comment('开关 1启用 2停用');
                $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('log')) {
            Schema::dropIfExists('log');
        }
    }
}
