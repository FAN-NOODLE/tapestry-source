<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldsToSubscribeMessageTouchTemplateTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasColumn('subscribe_message_touch_template', 'scene_type')) {
            Schema::table('subscribe_message_touch_template', function (Blueprint $table) {
                $table->unsignedTinyInteger('tid_used')->default(0)->comment('多类目同场景时只使用一个tid，该字段表示是否被使用，1用 2不用');
                $table->unsignedInteger('category_id')->default(0)->index()->comment('模板所属小程序账号类目id');
                $table->unsignedTinyInteger('scene_type')->default(0)->index()->comment('场景类别（1：直播订阅）');
                $table->string('scene_desc', 50)->default('')->comment('服务场景描述，15个字以内');
                $table->string('kids', 191)->default('')->comment('关键词id集合（按顺序），json存储');
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
        Schema::table('subscribe_message_touch_template', function (Blueprint $table) {
        });
    }
}
