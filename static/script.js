"use strict";
const guessImg = document.querySelector("#guess-img");
const questDiv = document.querySelector(".quest-content");
const yesButton = document.querySelector("#yes-button");
const noButton = document.querySelector("#no-button");
const overlay = document.querySelector("#overlay");
const banner = document.querySelector(".is-this-ai b");
if (yesButton == null ||
    noButton == null ||
    guessImg == null ||
    overlay == null ||
    questDiv == null ||
    banner == null) {
    document.body.innerText = "Some of the elements declared don't exist.";
}
else {
    guessImg.addEventListener("click", () => {
        guessImg.classList.toggle("guess-img-zoomin");
        guessImg.classList.toggle("guess-img-zoomout");
    });
    yesButton.addEventListener("click", async () => {
        const answerUrl = "http://127.0.0.1:4321/answer"; // production: http://xreality.space/answer
        const res = await fetch(answerUrl, {
            method: "POST",
            body: JSON.stringify({ isThisAI: true }),
        });
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }
        const answer = await res.json();
        if (answer.IsThisAI == true) {
            banner.style.color = "rgb(190, 201, 126)";
            banner.textContent = "🎉 You are right! This image is AI generated.";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
        else {
            banner.style.color = "rgb(248, 154, 138)";
            banner.textContent = "No, this image is created by human.";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
    });
    noButton.addEventListener("click", async () => {
        const answerUrl = "http://127.0.0.1:4321/answer"; // production: http://xreality.space/answer
        const res = await fetch(answerUrl, {
            method: "POST",
            body: JSON.stringify({ isThisAI: false }),
        });
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }
        const answer = await res.json();
        if (answer.IsThisAI == true) {
            banner.style.color = "rgb(248, 154, 138)";
            banner.textContent = "No, this image is AI generated.";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
        else {
            banner.style.color = "rgb(190, 201, 126)";
            banner.textContent = "🎉 You are right, this image is created by human.";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
    });
}
