<?php

use Illuminate\Database\Seeder;
use App\Models\Admin;

class DatabaseSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        $username = 'admin';
        $password = 'REDACTED_FOR_SOURCE_SHARING';
        $account = '15249279779';

        Admin::create([
            'username' => $username,
            'account' => $account,
            'password' => md5($account.$password),
            'status' => 1
        ]);
    }
}
