<i18n>
{
  "en": {
    "stock": "Stock",
    "order": "order",
    "to order": "to order"
  },
  "ru": {
    "stock": "Склад",
    "order": "Заказ",
    "to order": "в заказ"
  }
}
</i18n>


<template>

    <div id="seller" class="container">

        <div class="row justify-content-center">
            <div class="col-sm-9 ">

                    <h4><i class="fas fa-home"></i> {{ $t('stock') }}</h4>
                    <table class="table table-striped table-hover table-sm">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>pn</th>
                                <th>name</th>
                                <th>count</th>
                                <th>cost</th>
                                <th>{{ $t('to order') }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(item, index) in stock">
                                <td>{{ index }}</td>
                                <td>{{ item.pn }}</td>
                                <td>{{ item.name }}</td>
                                <td>{{ item.count }}</td>
                                <td>{{ item.cost }}</td>
                                <td><button class="btn btn-primary" v-on:click="addOrder(item, index)">{{ $t('to order') }}</button></td>

                            </tr>
                        </tbody>

                    </table>

                    <h4><i class="fas fa-truck-loading"></i> {{ $t('order') }}</h4>

                    <table class="table  table-striped table-hover table-sm">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>pn</th>
                                <th>name</th>
                                <th>count</th>
                                <th>cost</th>
                                <th>sum</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(item, index) in order">
                                <td>{{ index }}</td>
                                <td>{{ item.pn }}</td>
                                <td>{{ item.name }}</td>
                                <td>{{ item.count }}</td>
                                <td>{{ item.cost }}</td>
                                <td>{{ item.cost * item.count }}</td>
                            </tr>
                        </tbody>

                    </table>
                    <h5>Total price: {{ sum }}</h5>

                </div>
            </div>
    </div>
</template>

<script>
export default {

    data() {
        return {
            stock:  [],
            order:  [],
            sum:    0
        }
    },
    methods: {
        addStock(item) {
            this.stock.push(item)
        },
        addOrder(stockItem, index) {
            if(this.stock[index].count > 0) {

                let orderItem = Object.assign({}, stockItem)
                orderItem.count = 1

                let i = this.order.findIndex((item) => {
                    return item.id === orderItem.id
                })

                console.log(i)

                if ( i >= 0)
                    this.order[i].count += 1
                else
                    this.order.push(orderItem)

                this.sum = 0
                this.order.forEach((item) => {
                    this.sum += item.count * item.cost
                })

                this.stock[index].count -= 1
            }
        }
    },
    mounted() {
        this.addStock({
            id: 0,
            pn: "SN189",
            name: "GoodsN",
            count: 3,
            cost: 12

        })
        this.addStock({
            id: 1,
            pn: "SN832",
            name: "GoodsA",
            count: 7,
            cost: 46
        })

        this.addStock({
            id: 2,
            pn: "SN652",
            name: "GoodsB",
            count: 4,
            cost: 4
        })
    }
}
</script>
