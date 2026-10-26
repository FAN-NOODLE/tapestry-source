<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateServiceMerchantSettingTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('service_merchant_setting', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary()->comment('自增id');
            $table->string('merchant_id')->default('')->comment('商户id');
            $table->string('app_id')->default('')->comment('商户关联的公众号appid');
            $table->string('merchant_key')->default('')->comment('商户密钥');
            $table->string('api_v3_key')->default('')->comment('商户后台api v3 密钥');
            $table->string('serial_no')->default('')->comment('商户后台api序列号');
            $table->string('platform_serial_no')->default('')->comment('平台序列号');
            $table->string('private_key_path')->default('')->comment('商户私钥目录');
            $table->string('public_key_path')->default('')->comment('商户公钥目录');
            $table->string('platform_cert')->default('')->comment('平台证书目录');
            $table->tinyInteger('status')->default(0)->comment('状态值 1=启用enable，2=禁用disable');
            $table->unsignedInteger('updated_at')->default(0);
            $table->unsignedInteger('created_at')->default(0);
            $table->unsignedInteger('deleted_at')->default(0);
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('service_merchant_setting');
    }
}
