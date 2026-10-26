<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateMemberUserRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('member_user_relation')) {
            Schema::create('member_user_relation', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id')->index();
                $table->char('mem_uid', 32)->default('')->comment('UC人员（导购）id')->index();
                $table->string('ec_id', 50)->default('')->comment('uc外部联系人记录id');
                $table->string('ec_userid', 100)->default('')->comment('uc外部联系人userid');
                $table->string('ec_name')->default('')->comment('uc外部联系人昵称');
                $table->string('ec_avatar')->default('')->comment('uc外部联系人头像');
                $table->tinyInteger('ec_gender')->default(3)->comment('uc外部联系人性别');
                $table->string('user_remark_name')->default('')->comment('客户备注名');
                $table->string('user_name_index')->default('')->comment('客户姓名索引');
                $table->string('user_description')->default('')->comment('客户描述');
                $table->unsignedtinyInteger('stars')->default(2)->comment('是否标星客户（1：是；2：否）')->index();
                $table->string('union_id', 50)->default('')->comment('uc外部联系人union_id')->index();
                $table->tinyInteger('is_exclusive')->default(2)->comment('是否专属导购：(1:是，2:否)');
                $table->unsignedInteger('relation_at')->default(0)->comment('关系标记时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `member_user_relation` comment '员工标记外部联系人信息表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('member_user_relation');
    }
}
