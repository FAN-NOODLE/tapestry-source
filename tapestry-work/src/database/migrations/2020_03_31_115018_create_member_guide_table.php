<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateMemberGuideTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('member_guide')) {
            Schema::create('member_guide', function (Blueprint $table) {
                $table->increments('id')->comment('自增主键');
                $table->unsignedInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->string('guide_id')->default('')->comment('支付即导购id');
                $table->string('mem_uid')->default('')->comment('memberUid');
                $table->unsignedInteger('allocation_status')->default(2)->comment('分配状态：1、开启，2、关闭，3、强制开启，4、永久关闭');
                $table->unsignedInteger('trade_guide_status')->default(2)->comment('支付即导购分配状态：1、开启，2、关闭');
                $table->unsignedInteger('limit_user_count')->default(0)->comment('分配客户上限');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index('enterprise_id', 'enterprise_id');
            });

            // 表注释
            DB::statement("ALTER TABLE `member_guide` comment '支付即导购-员工表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('member_guide');
    }
}
