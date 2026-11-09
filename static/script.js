const guessImg = document.getElementById("guess-img");
const yesButton = document.getElementById("yes-button");
const noButton = document.getElementById("no-button");
const toast = document.getElementById("toast-div");

let isZoomed = false;

guessImg.addEventListener("click", () => {
  isZoomed = !isZoomed;
  if (isZoomed) {
    guessImg.style.width = "90%";
    guessImg.style.height = "90%";
  } else {
    // Does not mean the properties is completely removed, css properties still apply
    guessImg.style.removeProperty("width");
    guessImg.style.removeProperty("height");
  }
});

yesButton.addEventListener("click", async () => {
  const answerUrl = "http://127.0.0.1:4321/answer";
  const res = await fetch(answerUrl);
  if (!res.ok) {
    throw new Error(`Response status: ${res.status}`);
  }
  const answer = await res.json();
  if (answer.IsThisAI == true) {
    toast.append("Yes, this image is AI generated.");
  } else {
    toast.append("No, this image is created by human.");
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
    toast.append("Yes, this image is AI generated.");
  } else {
    toast.append("No, this image is created by human.");
  }
});
