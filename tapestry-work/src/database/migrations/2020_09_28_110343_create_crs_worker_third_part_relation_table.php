<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCrsWorkerThirdPartRelationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('crs_worker_third_part_relation')) {
            Schema::create('crs_worker_third_part_relation', function (Blueprint $table) {
                $table->unsignedBigInteger('id', true);
                $table->unsignedTinyInteger('sync_source')->default(1)->comment('订单同步来源：1-crs, 2-hwj');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('third_part_user_id')->default(0)->comment('会员ID');
                $table->unsignedTinyInteger('store_staff_id')->default(0)->comment('分享人id');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
                $table->index('store_staff_id', 'crs_worker_third_part_relation_index_storeStaffId');
                $table->index('enterprise_id', 'crs_worker_third_part_relation_index_enterpriseId');
                $table->index('third_part_user_id', 'crs_worker_third_part_relation_index_thirdPartUserId');
            });

            DB::statement("alter table `crs_worker_third_part_relation` comment '员工与三方员工关联表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('crs_worker_third_part_relation');
    }
}
