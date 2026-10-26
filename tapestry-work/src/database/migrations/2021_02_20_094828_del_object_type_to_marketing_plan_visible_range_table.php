<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class DelObjectTypeToMarketingPlanVisibleRangeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('marketing_plan_visible_range')) {
            if (Schema::hasColumn('marketing_plan_visible_range', 'object_type')) {
                Schema::table('marketing_plan_visible_range',function (Blueprint $table) {
                    $table->dropColumn(['object_type']);
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
            $table->tinyInteger('object_type')->default(1)->comment('1、按门店所属区域(默认),2、按导购所属部门');
        });
    }
}
