<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateMessageChannelTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('message_scene')) {
            Schema::create('message_scene', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary();
                $table->string('desc', 100)->default('')->comment('消息场景');
                $table->string('page_url')->default('')->comment('');
                $table->unsignedinteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedinteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedinteger('deleted_at')->default(0)->comment('删除时间');
            });
        }
        if (!Schema::hasTable('message_channel')) {
            Schema::create('message_channel', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary();
                $table->unsignedBigInteger('message_scene_id')->default(0)->comment('场景id');
                $table->string('channel')->default('')->comment('渠道');
                $table->unsignedBigInteger('scene_type')->default(0)->comment('实际分发的场景值');
                $table->unsignedTinyInteger('status')->default(1)->comment('开关 1开2关');
                $table->string('page_url')->default('')->comment('');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id');
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
