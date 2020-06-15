

function randomString() {
    return Math.random().toString(36).replace(/[^a-z]+/g, '').substr(0, 6)
}

function randomId(base) {
    return base + Math.random().toString(36).replace(/[^a-z]+/g, '').substr(0, 6)
}


export default {
    randomString,
    randomId
}

