function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);

    if (parts.length === 2)
        return parts.pop().split(";").shift();
}

document.addEventListener("DOMContentLoaded", function () {
    const csrfToken = getCookie("csrf_token");

    if (csrfToken) {
        let csrfElements = document.getElementsByClassName("csrf_token")

        for (let i = 0; i < csrfElements.length; i++) {
            csrfElements[i].value = csrfToken
        }
    }

    let modal = document.getElementById("modal");

    document.getElementById("åpne").onclick = () => {
        modal.style.display = "block";
    }

    document.getElementById("lukk").onclick = () => {
        modal.style.display = "none";
    }
});

function minus(num)
{
    document.getElementById("Remove" + num).remove()
    document.getElementById("Title" + num).remove()
    document.getElementById("Link" + num).remove()

    document.getElementById("amount").value = parseInt(document.getElementById("amount").value) - 1
}

function plus()
{
    let amount = document.getElementById("amount")

    amount.value = parseInt(amount.value) + 1

    let title = document.createElement("input")
    title.id = "Title" + amount.value
    title.name = "Title" + amount.value
    title.placeholder = "Tittel"
    title.style.marginRight = "4px"

    document.getElementById("links").insertBefore(title, document.getElementById("plusBTN"))

    let link = document.createElement("input")
    link.id = "Link" + amount.value
    link.name = "Link" + amount.value
    link.placeholder = "Lenke"
    link.style.marginRight = "4px"

    document.getElementById("links").insertBefore(link, document.getElementById("plusBTN"))

    let minusButton = document.createElement("button")
    minusButton.id = "Remove" + amount.value
    minusButton.type = "button"
    minusButton.onclick = () => minus(amount.value)
    minusButton.textContent = "-"

    document.getElementById("links").insertBefore(minusButton, document.getElementById("plusBTN"))

    let br = document.createElement("br")

    document.getElementById("links").insertBefore(br, document.getElementById("plusBTN"))
}