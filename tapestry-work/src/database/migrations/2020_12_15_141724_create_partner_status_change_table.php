<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreatePartnerStatusChangeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('partner_status_change')) {
            Schema::create('partner_status_change', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->comment('自增主键');
                $table->unsignedBigInteger('partner_enterprise_id')->default(0)->comment('合伙人在关联企业导购的ID');
                $table->tinyInteger('partner_status')->default(1)->comment('1:启用 0:禁用');
                $table->string('forbidden_reason', 200)->default('')->comment('禁用原因');
                $table->unsignedInteger('forbidden_code')->default(0)->comment('禁用编码');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });
            DB::statement("ALTER TABLE `partner_status_change` COMMENT '合伙人状态变更流水表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('partner_status_change');
    }
}
