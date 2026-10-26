<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddRoleFieldsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('role')) {
            Schema::table('role', function (Blueprint $table) {
                if (!Schema::hasColumn('role', 'type')) {
                    $table->unsignedTinyInteger('type')->default(1)->comment('类型,1:自然流量,2:门店,3:自然流量+门店')->after('powers');
                }
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
        if (Schema::hasTable('role')) {
            Schema::table('role', function (Blueprint $table) {
                $table->dropColumn('type');
            });
        }
    }
}
