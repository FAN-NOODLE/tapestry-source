<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateEnterpriseThirdPartAccount extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise_third_part_account')) {
            Schema::create('enterprise_third_part_account', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('app_id',191)->default('')->comment('第三方appid')->index();
                $table->string('app_secret',191)->default('')->comment('第三方app secret');
                $table->string('server_token',191)->default('')->comment('第三方server token');
                $table->string('encoding_aes_key',191)->default('')->comment('第三方encoding aes key');
                $table->string('refresh_token',191)->default('')->comment('第三方refresh token');
                $table->tinyInteger('authorize_type')->default('1')->comment('授权类型（1:手工授权，2:第三方组件授权）');
                $table->string('origin_id',191)->default('')->comment('第三方原始id');
                $table->string('name',191)->default('')->comment('第三方应用的名字');
                $table->tinyInteger('platform')->default(1)->comment('1:服务号，2:小程序');
                $table->tinyInteger('type')->default(1)->comment('区分具体的应用名称(1:商城小程序，2:企业微信小程序)');
                $table->tinyInteger('verify_status')->default(2)->comment('认证状态（1:已认证，2:未认证）');
                $table->string('logo',191)->default('')->comment('第三方应用logo');
                $table->text('business_info')->comment('第三方的商业能力json');
                $table->string('principal_name',191)->default('')->comment('在第三方登记的企业名');
                $table->string('signature',191)->default('')->comment('在第三方记录的应用备注');
                $table->text('other_config')->comment('第三方附属配置json');
                $table->text('authorizer_range')->comment('第三方授权范围json');
                $table->string('main_qrcode',191)->default('')->comment('第三方应用二维码');
                $table->string('tester_qrcode',191)->default('')->comment('第三方应用体验二维码');
                $table->tinyInteger('status')->default(2)->comment('授权状态（1:有效，2:无效）');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `enterprise_third_part_account` comment '小程序和服务号授权记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::drop('enterprise_third_part_account');
    }
}
