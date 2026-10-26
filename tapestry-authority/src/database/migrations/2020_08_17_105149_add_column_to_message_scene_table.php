<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddColumnToMessageSceneTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasColumn('message_scene', 'configure')) {
            Schema::table('message_scene', function (Blueprint $table) {
                $table->unsignedInteger('scene_type')->default(0)->comment('小程序scene_type');
                $table->unsignedInteger('configure')->default(0)->comment('后台需要配置策略');
                $table->unsignedInteger('sort')->default(0)->comment('配置策略的顺序');
                $table->string('touch_scene',500)->default('')->comment('用户订阅触点');
                $table->string('touch_page',100)->default('')->comment('触点页面');
                $table->string('message_scene',500)->default('')->comment('消息发送场景');
                $table->string('default_send_time',500)->default('')->comment('默认发送时间');
                $table->string('available_category_ids',500)->default('')->comment('可用的类目ids');
                $table->text('available_categories')->comment('可用的类目json');
                $table->string('cdp_chan',100)->default('')->comment('cdp chan');
                $table->string('cdp_params',100)->default('')->comment('cdp 需要拼接的参数');
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
        Schema::table('message_scene', function (Blueprint $table) {
        });
    }
}
