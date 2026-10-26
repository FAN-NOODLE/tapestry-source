<?php
/**
 * 营销计划可见范围表
 */
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateMarketingPlanVisibleRangeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $tableName = 'marketing_plan_visible_range';
        if (!Schema::hasTable($tableName)) {
            Schema::create($tableName, function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('plan_id')->default(0)->comment('营销计划表ID');
                $table->tinyInteger('object_type')->default(1)->comment('1、按门店所属区域(默认),2、按导购所属部门');
                $table->string('object_id', 255)->default('')->comment('对象ID(区域或者部门的ID)');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('plan_id', 'plan_id_index');
                $table->index(['object_type', 'object_id'], 'object_type_id_index');
            });

            DB::statement("alter table `{$tableName}` COMMENT '营销计划可见范围表' ");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //Schema::dropIfExists('marketing_plan');
    }
}
