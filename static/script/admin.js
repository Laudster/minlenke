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
    document.getElementById("Remove" + num).remove();
    document.getElementById("Title" + num).remove();
    document.getElementById("Link" + num).remove();
    document.getElementById("Break" + num).remove();


    let amount = document.getElementById("amount");
    let numInt = parseInt(num);

    if (amount.value > numInt) {
        const roundabout = ["Remove", "Title", "Link", "Break"];
        for (let i = numInt + 1; i < amount.value; i++) {
            let newNum = String(i - 1);

            for (let j = 0; j < roundabout.length; j++) {
                document.getElementById(roundabout[j] + i).name = roundabout[j] + newNum;
                document.getElementById(roundabout[j] + i).id = roundabout[j] + newNum;
            }
            
            document.getElementById("Remove" + newNum).onclick = () => minus(newNum);
        }
    }

    amount.value = parseInt(document.getElementById("amount").value) - 1;
}

function plus()
{
    let amount = document.getElementById("amount")
    let thisid = amount.value;

    let br = document.createElement("br")
    br.id = "Break" + thisid;
    document.getElementById("links").insertBefore(br, document.getElementById("plusBTN"))

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
    minusButton.onclick = () => minus(thisid)
    minusButton.textContent = "-"

    document.getElementById("links").insertBefore(minusButton, document.getElementById("plusBTN"))

    amount.value = parseInt(amount.value) + 1
}