<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateSubscribeMessageSettingTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('subscribe_message_setting')) {
            Schema::create('subscribe_message_setting', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('主键id');
                $table->unsignedInteger('category_id')->default(0)->index()->comment('模板所属小程序账号类目id');
                $table->unsignedInteger('tid')->default(0)->comment('微信模板标题id');
                $table->unsignedTinyInteger('type')->default(0)->comment('模版类型（2：一次性订阅，3：长期订阅）');
                $table->string('touch_ids')->default('')->comment('触发事件id集合，json存储');
                $table->unsignedTinyInteger('scene_type')->default(0)->index()->comment('场景类别（1：直播订阅）');
                $table->string('scene_desc', 50)->default('')->comment('服务场景描述，15个字以内');
                $table->string('kids', 191)->default('')->comment('关键词id集合（按顺序），json存储');
                $table->string('key_rules', 191)->default('')->comment('模板参数key列表（按顺序），json存储');
                $table->string('page', 191)->default('')->comment('跳转的页面地址');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("alter table `subscribe_message_setting` comment '订阅消息-配置表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('subscribe_message_setting');
    }
}
