/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */


let data = [
                {
                    comName: "Foo Com",
                    data: [
                        {
                            label: "2012-09-01",
                            value: 11
                        },
                        {
                            label: "2012-07-01",
                            value: 10
                        },
                        {
                            label: "2012-11-01",
                            value: 12
                        },
                        {
                            label: "2013-01-01",
                            value: 7
                        }
                    ]
                },
                {
                    comName: "Lala Com",
                    data: [
                        {
                            label: "2012-09-01",
                            value: 16
                        },
                        {
                            label: "2012-10-01",
                            value: 18
                        },
                        {
                            label: "2012-11-01",
                            value: 19
                        },
                        {
                            label: "2013-01-01",
                            value: 6
                        }
                    ]
                },
                {
                    comName: "Obla Com",
                    data: [
                        {
                            label: "2012-09-01",
                            value: 8
                        },
                        {
                            label: "2012-10-01",
                            value: 16
                        },
                        {
                            label: "2012-11-01",
                            value: 11
                        },
                        {
                            label: "2013-01-01",
                            value: 16
                        }
                    ]
                }
]

let collection = []
let legend = []

data.forEach((com, i) => {
    legend[i] = com.comName

    com.data.forEach((point, n) => {
        let found = false
        collection.forEach((colElem, l) => {
            if (colElem.label == point.label) {
                collection[l][i] = point.value
                found = true
            }
        })
        if (!found) {
            let elem = {
            }
            elem[i] = point.value
            elem["label"] = point.label
            collection.push(elem)
        }
    })
})

console.log(collection)
console.log(legend)


//for(let [i, com] of data.entries()) {
    //for (let [n, point] of com.data.entries()) {
        //let found = false
        //for (let [l, colElem] of ) {
            //if (colElem.label == point.label) {
                //collection[l][i] = point.value
                //found = true
            //}
        //}
        //if (!found) {
            //let elem = {
            //}
            //elem[i] = point.value
            //elem["label"] = point.label
            //collection.push(elem)
        //}
    //}
//}
