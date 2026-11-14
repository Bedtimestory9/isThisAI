"use strict";
const guessImg = document.querySelector("#guess-img");
const questDiv = document.querySelector(".quest-content");
const yesButton = document.querySelector("#yes-button");
const noButton = document.querySelector("#no-button");
const flipBoxBack = document.querySelector(".flip-box-back");
const flipBoxInner = document.querySelector(".flip-box-inner");
if (yesButton == null ||
    noButton == null ||
    flipBoxInner == null ||
    flipBoxBack == null) {
    throw new Error(`Some or all of the elements inside quest-div don't exist!`);
}
else {
    yesButton.addEventListener("click", async () => {
        const answerUrl = "http://127.0.0.1:4321/answer";
        const res = await fetch(answerUrl);
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }
        const answer = await res.json();
        if (answer.IsThisAI == true) {
            flipBoxBack.style.cssText = "background-color: rgb(190, 201, 126);";
            flipBoxBack.append("🎉 You are right! This image is AI generated.");
            flipBoxInner.style.cssText = "transform: rotateY(180deg);";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
        else {
            flipBoxBack.style.cssText = "background-color: rgb(248, 154, 138) ;";
            flipBoxBack.append("No, this image is created by human.");
            flipBoxInner.style.cssText = "transform: rotateY(180deg);";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
    });
    noButton.addEventListener("click", async () => {
        const answerUrl = "http://127.0.0.1:4321/answer";
        const res = await fetch(answerUrl);
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }
        const answer = await res.json();
        if (answer.IsThisAI == true) {
            flipBoxBack.style.cssText = "background-color: rgb(248, 154, 138) ;";
            flipBoxBack.append("Yes, this image is AI generated.");
            flipBoxInner.style.cssText = "transform: rotateY(180deg);";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
        else {
            flipBoxBack.style.cssText = "background-color: rgb(190, 201, 126);";
            flipBoxBack.append("🎉 You are right, this image is created by human.");
            flipBoxInner.style.cssText = "transform: rotateY(180deg);";
            yesButton.disabled = true;
            noButton.disabled = true;
            await new Promise((r) => setTimeout(r, 3000));
            window.location.reload();
        }
    });
}
