<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('enterprise')) {
            Schema::create('enterprise', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('自增id');
                $table->string('hash', 50)->default('')->comment('企业HASH');
                $table->string('name', 100)->default('')->comment('企业名称');
                $table->string('logo', 255)->default('')->comment('企业LOGO');
                $table->string('contact', 40)->default('')->comment('企业联系人');
                $table->string('mobile', 20)->default('')->comment('企业联系手机号');
                $table->string('address', 255)->default('')->comment('详细地址');
                $table->unsignedTinyInteger('crs_status')->default(1)->comment('CRS小程序状态（1:正常，2:禁用）');
                $table->unsignedTinyInteger('store_status')->default(1)->comment('商城小程序状态（1:正常，2:禁用）');
                $table->unsignedTinyInteger('contact_status')->default(1)->comment('企业微信通讯录应用状态（1:正常，2:禁用）');
                $table->string('corp_id', 100)->default('')->comment('企业Corpid');
                $table->string('domain', 32)->default('')->comment('企业域名');
                $table->unsignedTinyInteger('status')->default(1)->comment('状态（1:正常，2:禁用）');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->unique('hash', 'hash_unique');
            });
            \Illuminate\Support\Facades\DB::statement("alter table `enterprise` comment '企业信息表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise');
    }
}
