<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSourceTypeToSyncOrgLogTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('sync_org_log')) {
            Schema::table('sync_org_log', function (Blueprint $table) {
                $table->unsignedTinyInteger('source_type')->default(1)->after('status')->comment('操作类型，1=自动，2=手动');
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
        if (Schema::hasTable('sync_org_log')) {
            Schema::table('sync_org_log', function (Blueprint $table) {
                $table->dropColumn('source_type');
            });
        }
    }
}
