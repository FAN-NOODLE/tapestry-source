<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateQyWxMedia extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('qy_wx_media')) {
            Schema::create('qy_wx_media', function (Blueprint $table) {
                $table->string('scid', 20)->primary('scid');
                $table->string('uc_eid',60)->default('')->comment('uc平台的enterprise企业唯一标识');
                $table->tinyInteger('type')->default(0)->comment('素材类型 0文本素材 1图片素材|2图文链接|3小程序卡片');
                $table->string('name',255)->default('')->comment('图片素材名称|图文素材标题|小程序卡片标题');
                $table->string('text_content',4000)->default('')->comment('文本消息内容');
                $table->string('image_media_url',255)->default('')->comment('图文素材封面图片|(图片素材路径|小程序卡片素材路径|发送消息时取创建临时media_id,有效期3天)');
                $table->string('link_desc',512)->default('')->comment('图文消息简介');
                $table->text('link_content')->default('')->comment('链接富文本信息的copy，生成动态链接读取');
                $table->string('minip_appid',60)->default('')->comment('小程序appid');
                $table->string('minip_page',255)->default('')->comment('小程序路径');
                $table->unsignedInteger('user_id')->default(0)->comment('作者id');
                $table->string('user_name',60)->default('')->comment('作者名称');
                $table->tinyInteger('status')->default(0)->comment('-1删除');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            \Illuminate\Support\Facades\DB::statement("ALTER TABLE `qy_wx_media` COMMENT '企业微信消息发送素材列表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (!Schema::hasTable('qy_wx_media')) {
            Schema::dropIfExists('qy_wx_media');
        }
    }
}
