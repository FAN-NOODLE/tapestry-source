<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddObjectTypeToMarketingPlanTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('marketing_plan')) {
            if (!Schema::hasColumn('marketing_plan', 'object_type')) {
                Schema::table('marketing_plan',function (Blueprint $table) {
                    $table->tinyInteger('object_type')->default(1)->comment('1、按门店所属区域(默认),2、按导购所属部门');
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
            $table->dropColumn('object_type');
        });
    }
}
