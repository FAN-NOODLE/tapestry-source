<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateQyWxEmsg extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('qy_wx_emsg')) {
            Schema::create('qy_wx_emsg', function (Blueprint $table) {
                $table->string('tsid', 20)->primary('tsid');
                $table->string('uc_eid',60)->default('')->comment('uc平台的enterprise企业唯一标识');
                $table->tinyInteger('send_type')->default(0)->comment('发送类型 0默认发送文本消息 1图片素材|2图文链接|3小程序卡片');
                if(version_compare(DB::select("select version() AS V")[0]->V,'5.7','gt')) {
                    $table->json('send_object')->comment('发送对象 - {"type":"0全部客户|1所属员工|2渠道|3标签","list":[ID列表]}');
                }else{
                    $table->text('send_object')->default('')->comment('发送对象 - {"type":"0全部客户|1所属员工|2渠道|3标签","list":[ID列表]}');
                }
                $table->unsignedInteger('send_time')->default(0)->comment('发送时间 0为立即发送');
                $table->string('name',30)->default('')->comment('推送名称');
                $table->string('text_content',4000)->default('')->comment('文本消息内容');
                $table->string('image_media_url',255)->default('')->comment('图文素材封面图片|(图片素材路径|小程序卡片素材路径|发送消息时取创建临时media_id,有效期3天)');
                $table->string('link_title',64)->default('')->comment('链接标题');
                $table->string('link_desc',512)->default('')->comment('消息简介');
                $table->text('link_content')->default('')->comment('链接富文本信息的copy，生成动态链接读取');
                $table->string('minip_title',64)->default('')->comment('小程序卡片标题');
                $table->string('minip_appid',60)->default('')->comment('小程序appid');
                $table->string('minip_page',255)->default('')->comment('小程序路径|图片原始名称|图文消息的作者');
                $table->unsignedInteger('num_est')->default(0)->comment('预计触达人数');
                $table->unsignedInteger('num_fact')->default(0)->comment('预计触达人数');
                $table->unsignedInteger('user_id')->default(0)->comment('发送用户id');
                $table->tinyInteger('status')->default(0)->comment('0 草稿状态 1推送中 2已完成 3定时未开始');
                $table->unsignedInteger('pv')->default(0)->comment('点击量');
                $table->unsignedInteger('uv')->default(0)->comment('用户数');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("ALTER TABLE `qy_wx_emsg` COMMENT '企业微信外部联系人消息群发表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (!Schema::hasTable('qy_wx_emsg')) {
            Schema::dropIfExists('qy_wx_emsg');
        }
    }
}
