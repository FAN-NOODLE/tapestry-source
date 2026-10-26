<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreatePartnerEnterpriseMapTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('partner_enterprise_map')) {

            Schema::create('partner_enterprise_map', function (Blueprint $table) {
                $table->bigIncrements('id')->comment('自增主键');
                $table->unsignedBigInteger('partner_id')->default(0)->comment('合伙人ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('归属企业');
                $table->unsignedBigInteger('worker_id')->default(0)->comment('导购ID');
                $table->string('worker_name', 60)->default('')->comment('导购名称(暂时冗余,等表迁移了就可以放弃了)');
                $table->string('worker_uid', 150)->default('')->comment('导购用户编码mem_uid');
                $table->unsignedBigInteger('user_id')->default(0)->comment('合伙人的商城用户ID');
                $table->tinyInteger('partner_status')->default(1)->comment('1:启用 0:禁用');
                $table->tinyInteger('mobile_sensitive')->default(1)->comment('1:脱敏显示手机号 0:明文显示手机号');
                $table->string('mobile_show', 20)->default('')->comment('该导购下显示的手机号');
                $table->string('forbidden_reason', 200)->default('')->comment('禁用原因');
                $table->integer('forbidden_code')->default(0)->comment('禁用原因编码');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('updated_at', 'updated_at_index');
                $table->index('partner_id', 'partner_id_index');
                $table->index('created_at', 'created_at_index');
            });
            DB::statement("ALTER TABLE `partner_enterprise_map` COMMENT '合伙人和企业导购关联表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('partner_enterprise_map');
    }
}
