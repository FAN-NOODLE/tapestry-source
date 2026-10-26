<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddTypeToGuideUserRelation extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('guide_user_relation')) {
            Schema::table('guide_user_relation', function (Blueprint $table) {
                $table->unsignedTinyInteger('type')->default(1)->after('relation_at')->comment('分配类型：1、普通分配，2、支付即导购分配');
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('guide_user_relation')) {
            Schema::table('guide_user_relation', function (Blueprint $table) {
                $table->dropColumn('type');
            });
        }
    }
}
