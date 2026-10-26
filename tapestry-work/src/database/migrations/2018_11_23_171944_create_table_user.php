<?php

use \Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableUser extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('user')) {
            Schema::create('user', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('union_id', 80)->default('')->comment('用户统一id')->index();
                $table->string('open_id', 80)->default('')->comment('用户识别id')->index();
                $table->string('nick_name', 50)->default('')->comment('昵称');
                $table->string('real_name')->default('')->comment('真实姓名');
                $table->string('real_name_index')->default('')->comment('真实姓名索引');
                $table->string('mobile', 20)->default('')->comment('手机号码');
                $table->date('birthday')->nullable()->comment('生日');
                $table->string('avatar')->default('')->comment('头像');
                $table->string('country')->default('')->comment('国家');
                $table->string('province')->default('')->comment('地区');
                $table->string('city')->default('')->comment('城市');
                $table->unsignedInteger('register_at')->default(0)->comment('注册时间');
                $table->unsignedInteger('auth_at')->default(0)->comment('授权时间');
                $table->unsignedInteger('type')->default(0)->comment('0、未授权，1、已授权，2、未成交，3、已成交');
                $table->tinyInteger('sex')->default(3)->comment('性别（1:男，2:女，3:未知）');
                $table->unsignedInteger('user_amount')->default(0)->comment('会员消费金额，单位：分');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('修改时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("ALTER TABLE `user` COMMENT '用户信息表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('user');
    }
}
