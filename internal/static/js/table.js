document.addEventListener("DOMContentLoaded", () => {
  const elements = document.querySelectorAll("[data-harga]");
  const formatter = new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  });
  elements.forEach((element) => {
    const harga = Number(element.dataset.harga);
    const formatted = formatter.format(harga);
    element.textContent = formatted;
  });
});

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

// Modal detail produk
const produkCards = document.querySelectorAll("[data-produk]");
produkCards.forEach((card) => {
  card.addEventListener("click", () => {
    const encoded = card.dataset.produk;

    const json = atob(encoded);
    const produk = JSON.parse(json);

    // console.log(produk);
  });
});

const modal = document.getElementById("modal_produk");
const form = document.getElementById("add_menu_form");

const closeButton = document.getElementById("button_close_modal");

const produkNama = document.getElementById("modal_produk_nama");
const produkDeskripsi = document.getElementById("modal_produk_deskripsi");
const produkHarga = document.getElementById("modal_produk_harga");
const quantityElement = document.getElementById("modal_produk_quantity");
const produkIdInput = document.getElementById("produk_id_input");
const quantityInput = document.getElementById("quantity_input");
const notesInput = document.getElementById("notes_input");
const hapusButton = document.getElementById("hapus_button");
const simpanButton = document.getElementById("simpan_button");

const totalQuantityElement = document.getElementById("total_quantity");
const totalAmountElement = document.getElementById("total_amount");

let selectedProduk = null;

// Format Rupiah
const formatRupiah = (harga) => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(harga);
};

let quantity = 0;

// Open Modal
// function openModal(produk) {

//   modal.classList.remove("hidden");
//   document.body.classList.add("overflow-hidden");
// }

function closeModal() {
  modal.classList.add("hidden");
  document.body.classList.remove("overflow-hidden");
  selectedProduk = null;
}

// Product Card Click
// document.querySelectorAll(".produk-card").forEach((card) => {
//   card.addEventListener("click", () => {
//     const encoded = card.dataset.produk;
//     const produk = JSON.parse(atob(encoded));

//     openModal(produk);
//   });
// });

form.addEventListener("submit", () => {
  simpanButton.disabled = true;
});

// Close Button
// closeButton.addEventListener("click", closeModal);
// simpanButton.addEventListener("click", closeModal);

// Click Outside
// modal.addEventListener("click", (event) => {
//   if (event.target === modal) {
//     closeModal();
//   }
// });

// Escape
// document.addEventListener("keydown", (event) => {
//   if (event.key === "Escape" && !modal.classList.contains("hidden")) {
//     closeModal();
//   }
// });

// Increment
document.getElementById("button_increment").addEventListener("click", () => {
  quantity++;

  quantityElement.textContent = quantity;
  quantityInput.value = quantity;

  simpanButton.disabled = quantity === selectedProduk.qty;
});

// Decrement
document.getElementById("button_decrement").addEventListener("click", () => {
  if (quantity <= 0) return;

  quantity--;
  quantityElement.textContent = quantity;
  quantityInput.value = quantity;

  simpanButton.disabled = quantity === selectedProduk.qty || quantity === 0;
});

notesInput.addEventListener("input", (e) => {
  if (selectedProduk.qty === 0) return;
  simpanButton.disabled = e.target.value === selectedProduk.catatan;
});

hapusButton.addEventListener("click", async () => {
  hapusButton.disabled = true;
  const noMeja = hapusButton.dataset.meja;
  try {
    const response = await fetch(`/meja/${noMeja}/${selectedProduk.id}`, {
      method: "DELETE",
    });
    if (!response.ok) {
      throw new Error("Gagal menghapus pesanan");
    }
    window.location.href = response.url;
  } catch (error) {
    console.error(error);
    alert("Gagal menghapus pesanan. Silakan coba kembali.");
  } finally {
    hapusButton.disabled = false;
  }
});

// modal cart
const modalCart = document.getElementById("modal_cart");
const closeButtonCart = document.getElementById("button_close_modal_cart");

function openModalCart() {
  modalCart.classList.remove("hidden");
  document.body.classList.add("overflow-hidden");
}

function closeModalCart() {
  modalCart.classList.add("hidden");
  document.body.classList.remove("overflow-hidden");
}

document.getElementById("cart_button").addEventListener("click", () => {
  openModalCart();
});

closeButtonCart.addEventListener("click", closeModalCart);

// Click Outside
modalCart.addEventListener("click", (event) => {
  if (event.target === modalCart) {
    closeModalCart();
  }
});

//
const addMenu = document.getElementById("add-menu");

function openAddMenu(produk) {
  const url = new URL(window.location.href);

  url.searchParams.set("id", produk.id);
  url.hash = "add-menu";

  history.pushState({}, "", url);

  addMenu.classList.add("active");
  document.body.style.overflow = "hidden";

  selectedProduk = produk;
  quantity = selectedProduk.qty;

  // Isi informasi produk
  produkNama.textContent = selectedProduk.nama;
  produkDeskripsi.textContent = selectedProduk.deskripsi;
  produkHarga.textContent = formatRupiah(selectedProduk.harga);
  produkIdInput.value = selectedProduk.id;

  quantityElement.textContent = selectedProduk.qty;
  quantityInput.value = selectedProduk.qty;
  notesInput.value = selectedProduk.catatan;

  if (selectedProduk.qty == 0) {
    hapusButton.disabled = true;
  }
}

function closeAddMenu() {
  document.body.style.overflow = "";
  history.back();
}

document.getElementById("close_page").addEventListener("click", closeAddMenu);

window.addEventListener("popstate", () => {
  const url = new URL(window.location.href);

  const isOpen = url.hash === "#add-menu";

  document.getElementById("add-menu").classList.toggle("active", isOpen);
});

document.querySelectorAll(".produk-card").forEach((card) => {
  card.addEventListener("click", () => {
    try {
      const encoded = card.dataset.produk;
      const produk = JSON.parse(atob(encoded));

      openAddMenu(produk);
    } catch (error) {
      console.error("Gagal membaca data produk:", error);
    }
  });
});
