<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddAppletPathUpdateToMarketingPlanTaskResourceTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('marketing_plan_task_resource')) {
            if (!Schema::hasColumn('marketing_plan_task_resource', 'applet_path_update')) {
                Schema::table('marketing_plan_task_resource',function (Blueprint $table) {
                    $table->tinyInteger('applet_path_update')->default('1')->comment('小程序码跳转路径是否更新:1、是  2、否');
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
        Schema::table('marketing_plan_task_resource', function (Blueprint $table) {
            $table->dropColumn('applet_path_update');
        });
    }
}
