<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreatePartnerMemberAuthTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('partner_member_auth')) {

            Schema::create('partner_member_auth', function (Blueprint $table) {
                $table->bigIncrements('id')->comment('自增主键');
                $table->string('app_id', 255)->default('')->comment('小程序app_id');
                $table->string('partner_sn', 20)->default('')->comment('合伙人编码');
                $table->string('partner_name', 40)->default('')->comment('合伙人名称');
                $table->string('partner_mobile', 20)->default('')->comment('手机号');
                $table->tinyInteger('partner_sex')->default(0)->comment('0:未知 1:男 2:女');
                $table->string('open_id', 80)->default('')->comment('微信openid');
                $table->string('union_id', 80)->default('')->comment('微信union_id');
                $table->string('avatar', 255)->default('')->comment('用户头像');
                $table->unsignedInteger('mobile_short')->default(0)->comment('手机号后四位');
                $table->unsignedBigInteger('user_id')->default(0)->comment('商城user_id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("ALTER TABLE `partner_member_auth` COMMENT '合伙人授权基本信息表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('partner_member_auth');
    }
}
