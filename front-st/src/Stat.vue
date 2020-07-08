
<template>
    <layout>
        <el-tabs v-model="role" @tab-click="onTab">
            <el-tab-pane label="как ТК" name="transportCompany"></el-tab-pane>
            <el-tab-pane label="как СК" name="insuranceCompany"></el-tab-pane>
        </el-tabs>

        <el-card class="box-card">
            <h3>Статистика по претензиям</h3>

            <el-form :inline="true" label-position="top">

                <el-form-item label="Начало периода">
                    <el-date-picker v-model="begin" type="date" :clearable="true" placeholder="Укажите начало периода">
                    </el-date-picker>
                </el-form-item>

                <el-form-item label="Конец периода">
                    <el-date-picker v-model="end" type="date" :clearable="true" placeholder="Укажите конец периода">
                    </el-date-picker>
                </el-form-item>

                <el-form-item label="Тип страхования">
                    <el-select v-model="insuranceType" multiple :clearable="true" placeholder="Укажите тип страхования">
                        <el-option v-for="item in insuranceTypeOptions" :key="item.value" :label="item.label" :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item label="Статус претензии">
                    <el-select v-model="claimStatus" multiple :clearable="true" placeholder="Укажите статус претензии">
                        <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item v-if="isTransportCompany" label="Список страхователей">
                    <el-select v-model="policyHolders" multiple :clearable="true" placeholder="Укажите страхователей">
                        <el-option v-for="item in policyHoldersOptions" :key="item.value" :label="item.label" :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item v-if="isInsuranceCompany" label="Список страховщиков">
                    <el-select v-model="insurers" multiple :clearable="true" placeholder="Укажите страховщиков">
                        <el-option v-for="item in insurersOptions" :key="item.value" :label="item.label" :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item label="Валюта">
                    <el-select v-model="currency" :clearable="true" placeholder="Укажите валюту">
                        <el-option v-for="item in currencyOptions" :key="item.value" :label="item.label" :value="item.value">
                        </el-option>
                    </el-select>
                </el-form-item>

                <el-checkbox v-model="vip">VIP</el-checkbox>

                <el-button size="medium">Принять</el-button>

                <el-button v-on:click="cleanForm" size="medium" native-type="reset">Сбросить</el-button>

            </el-form>

        </el-card>
    </layout>
</template>

<script>

import Layout from './Layout.vue'

export default {
    components: {
        Layout
    },
    methods: {
        onTab(tab, event) {
            console.log(tab._props.name);
        },
        cleanForm() {
            console.log("cleanForm")
            this.insuranceType =  []
            this.claimStatus =    []
            this.claimSubject =   []
            this.policyHolders =  []
            this.insurers =       []
            this.currency =       ""
        },
        acceptForm() {
            console.log("acceptForm")
        }
    },
    computed: {
        isTransportCompany() {
            return this.role == 'transportCompany'
        },
        isInsuranceCompany() {
            return this.role == 'insuranceCompany'
        }
    },
    data() {
        return {
            role:           "transportCompany",
            roleName:       "",
            begin:          "",
            end:            "",
            insuranceType:  [],
            claimStatus:    [],
            claimSubject:   [],
            policyHolders:  [],
            insurers:       [],
            currency:       "",
            vip:            false,

            insuranceTypeOptions: [
                { label: "Сроки",       value: "Сроки" },
                { label: "Грузы",       value: "Грузы" },
                { label: "Неизвестно",  value: "Неизвестно" },
            ],
            statusOptions: [
                { label: "Удовлетворена",              value: "Удовлетворена" },
                { label: "Отозвана",                   value: "Отозвана" },
                { label: "Рассмотрение страховщиком",  value: "Рассмотрение страховщиком" },
                { label: "Принята",                    value: "Принята" },
                { label: "Корректировка страхователем", value: "Корректировка страхователем" },
                { label: "Новая претензия",            value: "Новая претензия" },
                { label: "Отказано" ,                  value: "Отказано" },
            ],
            policyHoldersOptions: [
                { label: "ООО «МИС»",                  value: "MisMSP" },
                { label: "ООО «Деловые линии»",        value: "DellinMSP" },
                { label: "ООО «ЦАП 2015»",             value: "AvtoTransitMSP" },
                { label: "ООО «Компания Скиф-Карго»",  value: "SkifMSP" },
            ],
            insurersOptions: [
                { label: "AО «Группа Ренессанс Страхование»",  value: "RenaissanceMSP" },
                { label: "СПАО «Ингосстрах»",                  value: "IngosMSP" },
                { label: "АО «АльфаСтрахование»",              value: "AlfastrahMSP" },
            ],
            currencyOptions: [
                { label: "rub", value: "rub" },
            ],
        }
    }
}
</script>
