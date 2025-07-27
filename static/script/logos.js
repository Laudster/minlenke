let favicons = document.getElementsByClassName("favicon")

for (let i = 0; i < favicons.length; i++) {
    let favicon = ""

    let splits = favicons[i].src.split("/")

    favicon += "https://img.logo.dev/" + splits[2] + "?token=pk_VROKU-lyQAWLHET6MxycrA&retina=true"

    favicons[i].src = favicon
}
