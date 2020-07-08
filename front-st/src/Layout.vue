<template>

    <div id="layout">

        <el-menu class="el-menu-demo" mode="horizontal" background-color="#545c64" text-color="#fff" active-text-color="#ffd04b">
            <el-menu-item>M5APP</el-menu-item>
        </el-menu>

        <el-main>
            <slot/>
        </el-main>

        <el-footer>
            <hr/>
        </el-footer>

    </div>

</template>
