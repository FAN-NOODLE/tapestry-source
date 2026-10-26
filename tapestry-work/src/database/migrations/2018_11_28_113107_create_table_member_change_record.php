<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateTableMemberChangeRecord extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('member_change_record')) {
            Schema::create('member_change_record', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->index()->comment('企业id')->index();
                $table->string('old_mem_uid', 50)->default('')->comment('离职员工uid');
                $table->string('ec_userid', 80)->default('')->comment('uc外部联系人userid');
                $table->string('mem_uid', 50)->default('')->comment('新导购uid');
                $table->unsignedtinyInteger('status')->default(2)->index()->comment('外部联系人确认状态（1：已确认；2：待确认）');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            // 表注释
            DB::statement("ALTER TABLE `member_change_record` comment '外部联系人再分配记录表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('member_change_record');
    }
}
