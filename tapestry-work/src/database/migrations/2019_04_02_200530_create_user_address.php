<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateUserAddress extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('user_address')) {
            Schema::create('user_address', function (Blueprint $table) {
                $table->increments('id');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedInteger('user_id')->default(0)->comment('用户ID')->index();
                $table->string('province')->default('')->comment('省');
                $table->string('city')->default('')->comment('市');
                $table->string('area')->default('')->comment('区');
                $table->string('address')->default('')->comment('详细地址');
                $table->string('province_code')->default('')->comment('省编码');
                $table->string('city_code')->default('')->comment('市编码');
                $table->string('area_code')->default('')->comment('区编码');
                $table->unsignedInteger('zip_code')->default(0)->comment('邮编');
                $table->unsignedTinyInteger('default')->default(2)->comment('是否默认地址(1:是，2:否)');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            Db::statement("ALTER TABLE `user_address` COMMENT '用户地址管理表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
        Schema::dropIfExists('user_address');
    }
}
