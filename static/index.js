const z = {
    name:"salih",
    hello: "massage"
}

console.log(z)
const zJson = JSON.stringify(z)

document.onclick = (e) => {
    alert(zJson)
}
