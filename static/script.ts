const guessImg = document.querySelector<HTMLImageElement>("#guess-img");
const questDiv = document.querySelector<HTMLDivElement>(".quest-content");
const yesButton = document.querySelector<HTMLButtonElement>("#yes-button");
const noButton = document.querySelector<HTMLButtonElement>("#no-button");
const overlay = document.querySelector<HTMLDivElement>("#overlay");
const banner = document.querySelector<HTMLDivElement>(".is-this-ai b");

type Answer = {
  ID: string;
  IsThisAI: boolean;
};

if (
  yesButton == null ||
  noButton == null ||
  guessImg == null ||
  overlay == null ||
  questDiv == null ||
  banner == null
) {
  document.body.innerText = "Some of the elements declared don't exist.";
} else {
  guessImg.addEventListener("click", () => {
    guessImg.classList.toggle("guess-img-zoomin");
    guessImg.classList.toggle("guess-img-zoomout");
  });

  yesButton.addEventListener("click", async () => {
    const answerUrl = "http://127.0.0.1:4321/answer";
    const res = await fetch(answerUrl);
    if (!res.ok) {
      throw new Error(`Response status: ${res.status}`);
    }
    const answer: Answer = await res.json();
    if (answer.IsThisAI == true) {
      banner.style.color = "rgb(190, 201, 126)";
      banner.textContent = "🎉 You are right! This image is AI generated.";
      yesButton.disabled = true;
      noButton.disabled = true;
      await new Promise((r) => setTimeout(r, 3000));
      window.location.reload();
    } else {
      // flipBoxBack.style.cssText = "background-color: rgb(248, 154, 138) ;";
      // flipBoxBack.append("No, this image is created by human.");
      //
      // flipBoxInner.style.cssText = "transform: rotateY(180deg);";

      banner.style.color = "rgb(248, 154, 138)";
      banner.textContent = "No, this image is created by human.";
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
    const answer: Answer = await res.json();
    if (answer.IsThisAI == true) {
      // flipBoxBack.style.cssText = "background-color: rgb(248, 154, 138) ;";
      // flipBoxBack.append("Yes, this image is AI generated.");
      //
      // flipBoxInner.style.cssText = "transform: rotateY(180deg);";

      banner.style.color = "rgb(248, 154, 138)";
      banner.textContent = "Yes, this image is AI generated.";
      yesButton.disabled = true;
      noButton.disabled = true;
      await new Promise((r) => setTimeout(r, 3000));
      window.location.reload();
    } else {
      // flipBoxBack.style.cssText = "background-color: rgb(190, 201, 126);";
      // flipBoxBack.append("🎉 You are right, this image is created by human.");
      //
      // flipBoxInner.style.cssText = "transform: rotateY(180deg);";

      banner.style.color = "rgb(190, 201, 126)";
      banner.textContent = "🎉 You are right, this image is created by human.";
      yesButton.disabled = true;
      noButton.disabled = true;
      await new Promise((r) => setTimeout(r, 3000));
      window.location.reload();
    }
  });
}
