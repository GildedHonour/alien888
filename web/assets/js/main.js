"use sctrict";

const burger = document.querySelector(".burger");
const menu = document.getElementById("menu");

if (burger) {
  burger.addEventListener("click", () => {
    burger.classList.toggle("active");
    menu.classList.toggle("active");
  });
}


// Language

const langItem = document.querySelector(".header-lang-btn");
const langItemBtn = document.querySelector(".header-lang-btn button");
const langDropBtn = document.querySelectorAll(".header-lang-drop button");
const menuLangItem = document.querySelector(".menu-lang-btn");
const menuLangItemBtn = document.querySelector(".menu-lang-btn button");
const menuLangDropBtn = document.querySelectorAll(".menu-lang-drop button");

if (langItem) {
  langItemBtn.addEventListener("click", () => {
    langItem.classList.toggle("active");
  });

  langDropBtn.forEach((el) => {
    el.addEventListener("click", (e) => {
      let lang = e.target.getAttribute("data-lang");
      langItemBtn.setAttribute("value", lang);
      langItemBtn.querySelector("p").textContent = lang;
      langItem.classList.remove("active");
    });
  });
}

if (menuLangItem) {
  menuLangItemBtn.addEventListener("click", () => {
    menuLangItem.classList.toggle("active");
  });

  menuLangDropBtn.forEach((el) => {
    el.addEventListener("click", (e) => {
      let lang = e.target.getAttribute("data-lang");
      menuLangItemBtn.setAttribute("value", lang);
      menuLangItemBtn.querySelector("p").textContent = lang;
      menuLangItem.classList.remove("active");
    });
  });
}
