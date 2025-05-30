let favicons = document.getElementsByClassName("favicon")

for (let i = 0; i < favicons.length; i++) {
    let favicon = ""

    let splits = favicons[i].src.split("/")

    favicon += "https://" + splits[2] + "/favicon.ico"

    favicons[i].src = favicon
}