document.addEventListener("DOMContentLoaded", function() {
    // After the opening animation, hide the animation container
    const openingAnimation = document.querySelector(".opening-animation");

    openingAnimation.addEventListener("animationend", function() {
        openingAnimation.style.display = "none";
    });
});