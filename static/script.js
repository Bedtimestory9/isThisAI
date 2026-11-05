const guessImg = document.getElementById("guess-img");

let isZoomed = false;

guessImg.addEventListener("click", () => {
  isZoomed = !isZoomed;
  if (isZoomed) {
    guessImg.style.width = "90vw";
  } else {
    guessImg.style.width = "400px";
  }
});
