<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddWorkerIdAndMemUidToMarketingPlanVisibleRangeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('marketing_plan_visible_range')) {
            if (!Schema::hasColumn('marketing_plan_visible_range', 'worker_id')) {
                Schema::table('marketing_plan_visible_range',function (Blueprint $table) {
                    $table->string('mem_uid', 150)->default('')->comment('导购uid');
                    $table->unsignedBigInteger('worker_id')->default(0)->comment('导购ID');
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
        Schema::table('marketing_plan_visible_range', function (Blueprint $table) {
            $table->dropColumn('worker_id');
            $table->dropColumn('mem_uid');
        });
    }
}
