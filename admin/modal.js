window.onload = () => {
    const enhancer = document.querySelector("#modalEnhanced");
    const r = document.querySelector(':root');

    if (enhancer) {
        r.style.setProperty('--modalCloseSpeed', '0.1s');
    }
}