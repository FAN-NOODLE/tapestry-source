<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddAddObjectIdsToMarketingPlanTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('marketing_plan')) {
            if (!Schema::hasColumn('marketing_plan', 'add_object_ids')) {
                Schema::table('marketing_plan',function (Blueprint $table) {
                    $table->text('add_object_ids')->default('')->comment('新增的可见范围');
                    $table->string('plan_sn',50)->default('')->comment('计划编号');
                });
            }

        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('marketing_plan', function (Blueprint $table) {
            $table->dropColumn('add_object_ids');
            $table->dropColumn('plan_sn');
        });
    }
}
