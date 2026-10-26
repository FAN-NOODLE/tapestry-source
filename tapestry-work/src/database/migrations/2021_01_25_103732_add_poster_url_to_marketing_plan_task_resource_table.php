<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddPosterUrlToMarketingPlanTaskResourceTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('marketing_plan_task_resource')) {
            if (!Schema::hasColumn('marketing_plan_task_resource', 'poster_url')) {
                Schema::table('marketing_plan_task_resource',function (Blueprint $table) {
                    $table->string('poster_url',255)->default('')->comment('后台预览海报，没有导购参数的海报');
                    $table->string('sun_code',255)->default('')->comment('预览太阳码');
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
            $table->dropColumn('poster_url');
            $table->dropColumn('sun_code');
        });
    }
}
