const navContainer = document.getElementById("nav-container");
const navLinks = navContainer.querySelectorAll('a[href^="#grup-"]');
const sections = document.querySelectorAll('[id^="grup-"]');

let active = "";

const observer = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        active = entry.target.id;

        const activeLink = navContainer.querySelector(`[href="#${active}"]`);

        activeLink?.scrollIntoView({
          behavior: "smooth",
          block: "nearest",
          inline: "center",
        });

        navLinks.forEach((link) => {
          link.classList.toggle(
            "active",
            link.getAttribute("href") === `#${active}`,
          );
        });
      }
    });
  },
  {
    rootMargin: "-64px 0px -60% 0px",
    threshold: 0,
  },
);

sections.forEach((section) => {
  observer.observe(section);
});

const buttonBackToTop = document.getElementById("button_back_to_top");
buttonBackToTop.addEventListener("click", () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
});

const backToTop = document.getElementById("back_to_top");

const checkScroll = () => {
  const show =
    window.innerHeight + window.scrollY >=
    document.documentElement.scrollHeight - 50;

  backToTop.classList.toggle("hidden", !show);
};

window.addEventListener("scroll", checkScroll);
window.addEventListener("resize", checkScroll);

checkScroll();
